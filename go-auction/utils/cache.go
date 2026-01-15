package utils

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"go-auction/config"
)

// CacheGet 从缓存获取值并反序列化
func CacheGet(ctx context.Context, key string, dest interface{}) error {
	val, err := config.RedisGet(ctx, key)
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), dest)
}

// CacheSet 将值序列化后存入缓存
func CacheSet(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal cache value: %w", err)
	}
	return config.RedisSet(ctx, key, string(data), expiration)
}

// CacheDelete 删除缓存
func CacheDelete(ctx context.Context, keys ...string) error {
	return config.RedisDelete(ctx, keys...)
}

// CacheDeletePattern 根据模式删除缓存
func CacheDeletePattern(ctx context.Context, pattern string) error {
	return config.RedisDeletePattern(ctx, pattern)
}

// CacheKey 生成缓存键
func CacheKey(parts ...interface{}) string {
	key := ""
	for i, part := range parts {
		if i > 0 {
			key += ":"
		}
		key += fmt.Sprintf("%v", part)
	}
	return key
}
