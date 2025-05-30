package models

import (
	"gorm.io/gorm"
)

type Comment struct {
	gorm.Model
	Content string `gorm:"size:500" json:"content"`
	UserID  uint   `json:"user_id"`
	VideoID uint   `json:"video_id"`

	User User `gorm:"foreignKey:UserID"` // 关联用户，方便获取昵称头像
}
