package models

import (
	"time"
)

type Danmaku struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	VideoID   uint      `json:"video_id"` // 关联视频ID
	UserID    uint      `json:"user_id"`  // 发送弹幕的用户ID
	Content   string    `json:"content"`  // 弹幕内容
	Position  string    `json:"position"` // 弹幕位置，例如 "top", "bottom", "middle"
	Color     string    `json:"color"`    // 弹幕颜色，HEX或字符串
	CreatedAt time.Time `json:"created_at"`
}
