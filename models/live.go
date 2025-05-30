package models

import (
	"gorm.io/gorm"
	"time"
)

type LiveStatus int

const (
	LivePending LiveStatus = 0 // 未开始
	LiveOngoing LiveStatus = 1 // 直播中
	LiveStopped LiveStatus = 2 // 正常结束
	LiveForced  LiveStatus = 3 // 管理员强制结束
	LiveCrashed LiveStatus = 4 // 异常断开
)

type Live struct {
	gorm.Model
	UserID    uint       `gorm:"not null" json:"user_id"`
	Title     string     `gorm:"size:100" json:"title"`
	CoverURL  string     `gorm:"size:255" json:"cover_url"`
	StreamKey string     `gorm:"size:100;uniqueIndex" json:"stream_key"`
	Status    LiveStatus `gorm:"default:0" json:"status"` // 使用 tinyint 存储
	StartTime *time.Time `json:"start_time"`
	EndTime   *time.Time `json:"end_time"`

	User User `gorm:"foreignKey:UserID"` // 关联用户，方便获取昵称头像
}
