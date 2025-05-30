package service

import (
	"context"
	"fmt"
	"time"
	"user-system/config"
)

func SaveDanmakuToRedis(userID uint, videoID uint, content string) {
	key := fmt.Sprintf("danmaku:%d", videoID)
	msg := fmt.Sprintf("[%d] %s", userID, content)
	_ = config.RedisClient.RPush(context.Background(), key, msg).Err()

	// 可选：设置过期时间
	config.RedisClient.Expire(context.Background(), key, time.Hour*6)
}
