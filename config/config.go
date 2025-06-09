package config

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"

	"user-system/models"
)

var DB *gorm.DB
var RedisClient *redis.Client

func Init() {
	// 加载 .env 文件
	err := godotenv.Load()
	if err != nil {
		log.Println(".env 文件加载失败，使用默认环境变量")
	}

	// 从环境变量读取数据库和 Redis 配置
	mysqlDSN := os.Getenv("MYSQL_DSN")

	fmt.Println(">>> 正在连接 MySQL...")
	db, err := gorm.Open(mysql.Open(mysqlDSN), &gorm.Config{})
	if err != nil {
		log.Fatal("❌ MySQL 连接失败：", err)
	}
	fmt.Println("✅ MySQL 连接成功")

	// 自动迁移建表
	db.AutoMigrate(&models.User{}, &models.Video{}, &models.Comment{}, &models.Live{})
	DB = db

	InitRedis()

	fmt.Println("✅ 所有服务初始化完成")
}

func InitRedis() {
	fmt.Println(">>> 正在连接 Redis...")
	redisAddr := os.Getenv("REDIS_ADDR")
	RedisClient = redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})
	_, err := RedisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal("❌ Redis 连接失败：", err)
	}
	fmt.Println("✅ Redis 连接成功")
}
