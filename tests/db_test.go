package tests

import (
	"log"
	"user-system/models"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// SetupTestDB 初始化 SQLite 内存数据库，用于测试
func SetupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect test db: %v", err)
	}

	// 自动迁移模型（包含评论、用户、视频）
	if err := db.AutoMigrate(&models.User{}, &models.Video{}, &models.Comment{}); err != nil {
		log.Fatalf("failed to migrate test db: %v", err)
	}

	// 清空表中的数据
	db.Exec("DELETE FROM users")
	db.Exec("DELETE FROM videos")
	db.Exec("DELETE FROM comments")
	
	return db
}
