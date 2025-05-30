package dto

import (
	"user-system/models"
	"user-system/utils"
)

// CommentRequest 用户提交的评论数据
type CommentRequest struct {
	VideoID uint   `json:"video_id" binding:"required"`
	Content string `json:"content" binding:"required,min=1,max=500"`
}

type CommentListRequest struct {
	VideoID uint `form:"video_id" binding:"required"`
	LastID  uint `form:"last_id"` // 上一次的最后一条评论 ID，游标
	Limit   int  `form:"limit"`   // 每页数量，默认值后端处理
}

// CommentResponse 返回给前端的评论结构，包含用户昵称和格式化时间
type CommentResponse struct {
	ID        uint   `json:"id"`
	VideoID   uint   `json:"video_id"`
	UserID    uint   `json:"user_id"`
	Nickname  string `json:"nickname"` // 关联用户昵称
	Content   string `json:"content"`
	CreatedAt string `json:"created_at"` // 格式化时间
	UpdatedAt string `json:"updated_at"`
}

// BuildCommentResponse 将模型转换为响应结构体，格式化时间
func BuildCommentResponse(c *models.Comment) CommentResponse {
	return CommentResponse{
		ID:        c.ID,
		VideoID:   c.VideoID,
		UserID:    c.UserID,
		Nickname:  c.User.Name,
		Content:   c.Content,
		CreatedAt: utils.FormatTime(c.CreatedAt),
		UpdatedAt: utils.FormatTime(c.UpdatedAt),
	}
}

// BuildCommentResponseList 批量转换
func BuildCommentResponseList(comments []models.Comment) []CommentResponse {
	list := make([]CommentResponse, len(comments))
	for i, c := range comments {
		list[i] = BuildCommentResponse(&c)
	}
	return list
}
