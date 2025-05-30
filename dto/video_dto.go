package dto

import (
	"user-system/models"
	"user-system/utils"
)

type VideoRequest struct {
	Title       string `json:"title" binding:"required,min=2,max=100"`
	Description string `json:"description" binding:"max=500"`
	URL         string `json:"url" binding:"required,url"`
	Duration    int    `json:"duration"`
	Tags        string `json:"tags"`
}

// VideoResponse 是返回给前端的视频结构，时间格式化、剔除敏感字段
type VideoResponse struct {
	ID          uint   `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	URL         string `json:"url"`
	Duration    int    `json:"duration"`
	Tags        string `json:"tags"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// 转换模型到响应结构，带时间格式化
func BuildVideoResponse(video *models.Video) VideoResponse {
	return VideoResponse{
		ID:          video.ID,
		Title:       video.Title,
		Description: video.Description,
		URL:         video.URL,
		Duration:    video.Duration,
		Tags:        video.Tags,
		CreatedAt:   utils.FormatTime(video.CreatedAt),
		UpdatedAt:   utils.FormatTime(video.UpdatedAt),
	}
}

func BuildVideoResponseList(videos []models.Video) []VideoResponse {
	list := make([]VideoResponse, len(videos))
	for i, v := range videos {
		list[i] = BuildVideoResponse(&v)
	}
	return list
}
