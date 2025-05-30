package utils

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"time"
)

func BuildCommentCacheKey(videoID, lastID, limit uint) string {
	return fmt.Sprintf("video:%d:comments:lastid:%d:limit:%d", videoID, lastID, limit)
}

func RateLimit(rdb *redis.Client, limitType string, id uint, limit time.Duration) bool {
	key := fmt.Sprintf("%s:limit:%d", limitType, id)

	// 使用 SETNX 设置 key，如果设置成功，说明允许评论
	ok, err := rdb.SetNX(ctx, key, 1, limit).Result()
	if err != nil {
		return false // Redis 异常时默认限制（更安全）
	}
	return ok // true = 允许评论，false = 限流中
}
