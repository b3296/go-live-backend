package utils

import (
	"context"
	"github.com/go-redis/redis/v8"
	"time"
)

var ctx = context.Background()

func AcquireLock(rdb *redis.Client, key string, expiration time.Duration) bool {
	return rdb.SetNX(ctx, key, "locked", expiration).Val()
}

func ReleaseLock(rdb *redis.Client, key string) {
	rdb.Del(ctx, key)
}
