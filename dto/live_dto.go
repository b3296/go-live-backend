package dto

import (
	"user-system/models"
)

type LiveListRequest struct {
	LastID uint `form:"last_id"`
	Limit  int  `form:"limit"`
}

type StopLiveRequest struct {
	StreamKey string `json:"stream"`
}

type CreateLiveRequest struct {
	Title    string `json:"title" binding:"required"`
	CoverURL string `json:"cover_url"`
}

type UpdateLiveStatusRequest struct {
	Status int `json:"status" binding:"required,oneof=2 3 4"` // 只能是结束、强制、异常
}

type LiveStartRequest struct {
	App    string `json:"app"`
	Stream string `json:"stream"` // 这个就是我们之前生成的 stream_key
}

type LiveResponse struct {
	ID        uint   `json:"id"`
	Title     string `json:"title"`
	CoverURL  string `json:"cover_url"`
	UserID    uint   `json:"user_id"`
	Nickname  string `json:"nickname"` // 关联用户昵称
	StreamKey string `json:"stream_key"`
	StartTime string `json:"start_time"` // 格式化时间
}

func BuildLiveResponse(l *models.Live) LiveResponse {
	return LiveResponse{
		ID:        l.ID,
		Title:     l.Title,
		CoverURL:  l.CoverURL,
		Nickname:  l.User.Name,
		StreamKey: l.StreamKey,
		//StartTime: utils.FormatTime(*l.StartTime),
	}
}

func BuildLiveResponseList(lives []models.Live) []LiveResponse {
	list := make([]LiveResponse, len(lives))
	for i, l := range lives {
		list[i] = BuildLiveResponse(&l)
	}
	return list
}
