package models

import (
	"gorm.io/gorm"
)

// Video 数据库模型
type Video struct {
	gorm.Model
	Title       string `gorm:"size:100" json:"title"`
	Description string `gorm:"size:500" json:"description"`
	URL         string `json:"url"`
	Duration    int    `json:"duration"` // 秒
	Tags        string `json:"tags"`     // 逗号分隔标签
}
