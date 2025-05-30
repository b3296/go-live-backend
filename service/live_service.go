package service

import (
	"errors"
	"fmt"
	"time"
	"user-system/dto"
	"user-system/models"

	"gorm.io/gorm"
)

type LiveService struct {
	DB *gorm.DB
}

func NewLiveService(db *gorm.DB) *LiveService {
	return &LiveService{DB: db}
}

func (s *LiveService) CreateLive(userID uint, title, coverURL string) (*models.Live, error) {
	streamKey := generateStreamKey(userID)
	live := &models.Live{
		UserID:    userID,
		Title:     title,
		CoverURL:  coverURL,
		StreamKey: streamKey,
		Status:    models.LivePending,
	}
	return live, s.DB.Create(live).Error
}

func (s *LiveService) StartLive(streamKey string) error {
	var live models.Live
	if err := s.DB.Where("stream_key = ?", streamKey).First(&live).Error; err != nil {
		return err
	}
	now := time.Now()
	live.Status = models.LiveOngoing
	live.StartTime = &now
	return s.DB.Save(&live).Error
}

func (s *LiveService) StopLive(streamKey string, forced, crashed bool) error {
	var live models.Live
	if err := s.DB.Where("stream_key = ?", streamKey).First(&live).Error; err != nil {
		return err
	}
	if live.Status != models.LiveOngoing {
		return errors.New("直播未在进行中")
	}

	now := time.Now()
	live.Status = models.LiveStopped
	if forced {
		live.Status = models.LiveForced
	}
	if crashed {
		live.Status = models.LiveCrashed
	}
	live.EndTime = &now
	return s.DB.Save(&live).Error
}

func generateStreamKey(userID uint) string {
	return fmt.Sprintf("stream_%d_%d", userID, time.Now().UnixNano())
}

// GetList 获取评论分页列表（包含用户昵称）
func (l *LiveService) GetList(req dto.LiveListRequest) ([]models.Live, error) {
	if req.Limit <= 0 {
		req.Limit = 10
	}

	var lives []models.Live

	query := l.DB.Model(&models.Live{}).Where("status = ?", models.LivePending)
	if req.LastID > 0 {
		query = query.Where("id < ?", req.LastID)
	}
	err := query.Preload("User").
		Order("id desc").
		Limit(req.Limit).
		Find(&lives).Error
	if err != nil {
		return nil, err
	}

	return lives, nil
}
