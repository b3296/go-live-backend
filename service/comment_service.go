package service

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/go-redis/redis/v8"
	"time"
	"user-system/dto"
	"user-system/models"
	"user-system/utils"

	"gorm.io/gorm"
)

type CommentService struct {
	DB    *gorm.DB
	Redis *redis.Client
	ctx   context.Context
}

// NewCommentService 构造函数，注入数据库依赖
func NewCommentService(db *gorm.DB, rdb *redis.Client) *CommentService {
	return &CommentService{
		DB:    db,
		Redis: rdb,
		ctx:   context.Background(),
	}
}

// Create 创建评论（不允许匿名，强制绑定 userID）
func (s *CommentService) Create(req dto.CommentRequest, userID uint) (*models.Comment, error) {
	if userID == 0 {
		return nil, errors.New("未登录用户不允许评论")
	}

	// 限流判断
	if !utils.RateLimit(s.Redis, "comment_create", userID, 5*time.Second) {
		return nil, errors.New("评论太频繁，请稍后再试")
	}

	comment := models.Comment{
		Content: req.Content,
		VideoID: req.VideoID,
		UserID:  userID,
	}
	if err := s.DB.Create(&comment).Error; err != nil {
		return nil, err
	}
	// 预加载用户信息
	s.DB.Preload("User",
		func(db *gorm.DB) *gorm.DB {
			return db.Select(models.UserBaseFields)
		}).
		First(&comment, comment.ID)
	return &comment, nil
}

// GetComments 获取评论分页列表（包含用户昵称）
func (s *CommentService) GetComments(req dto.CommentListRequest) ([]models.Comment, error) {
	if req.Limit <= 0 {
		req.Limit = 10
	}

	key := utils.BuildCommentCacheKey(req.VideoID, req.LastID, uint(req.Limit))
	var comments []models.Comment

	// 1. 尝试从 Redis 缓存中获取
	cached, err := s.Redis.Get(s.ctx, key).Result()
	if err == nil {
		_ = json.Unmarshal([]byte(cached), &comments)
		return comments, nil
	}

	// 2. 加互斥锁防击穿
	lockKey := key + ":lock"
	if !utils.AcquireLock(s.Redis, lockKey, 5*time.Second) {
		time.Sleep(100 * time.Millisecond) // 让出 CPU
		return s.GetComments(req)          // 重试
	}
	defer utils.ReleaseLock(s.Redis, lockKey)

	// 3. 数据库查询（按 last_id 分页）
	query := s.DB.Model(&models.Comment{}).Where("video_id = ?", req.VideoID)
	if req.LastID > 0 {
		query = query.Where("id < ?", req.LastID)
	}
	err = query.Preload("User").
		Order("id desc").
		Limit(req.Limit).
		Find(&comments).Error
	if err != nil {
		return nil, err
	}

	// 4. 存入缓存（缓存空结果也写入防穿透）
	data, _ := json.Marshal(comments)
	s.Redis.Set(s.ctx, key, data, 5*time.Minute)

	return comments, nil
}
