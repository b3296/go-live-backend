package models

import (
	"gorm.io/gorm"
)

// User 用户模型，自动迁移建表
// gorm.Model 会自动包含 ID, CreatedAt, UpdatedAt, DeletedAt
// Email 为唯一字段

type User struct {
	gorm.Model
	Email    string `gorm:"unique" json:"email"`
	Password string `json:"password"`
	Name     string `gorm:"size:100" json:"name"`
}

// UserBaseFields 需要预加载的字段
var UserBaseFields = []string{"id", "name"}
