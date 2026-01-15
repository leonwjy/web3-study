package config

import (
	"context"
	"log/slog"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisClient 全局 Redis 客户端实例
var RedisClient *redis.Client

// InitRedis 初始化 Redis 连接
func InitRedis(cfg *Config) *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.GetAddr(),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		slog.Error("Redis 连接失败", "error", err, "addr", cfg.Redis.GetAddr())
		os.Exit(1)
	}

	RedisClient = rdb
	slog.Info("Redis 连接成功", "addr", cfg.Redis.GetAddr(), "db", cfg.Redis.DB)
	return rdb
}

// GetRedis 获取 Redis 客户端实例
func GetRedis() *redis.Client {
	return RedisClient
}

// CloseRedis 关闭 Redis 连接
func CloseRedis() error {
	if RedisClient != nil {
		return RedisClient.Close()
	}
	return nil
}

// RedisGet 获取缓存值
func RedisGet(ctx context.Context, key string) (string, error) {
	if RedisClient == nil {
		return "", redis.Nil
	}
	return RedisClient.Get(ctx, key).Result()
}

// RedisSet 设置缓存值
func RedisSet(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	if RedisClient == nil {
		return nil
	}
	return RedisClient.Set(ctx, key, value, expiration).Err()
}

// RedisDelete 删除缓存值
func RedisDelete(ctx context.Context, keys ...string) error {
	if RedisClient == nil {
		return nil
	}
	return RedisClient.Del(ctx, keys...).Err()
}

// RedisExists 检查 key 是否存在
func RedisExists(ctx context.Context, key string) (bool, error) {
	if RedisClient == nil {
		return false, nil
	}
	count, err := RedisClient.Exists(ctx, key).Result()
	return count > 0, err
}

// RedisDeletePattern 根据模式删除多个key（使用SCAN命令）
func RedisDeletePattern(ctx context.Context, pattern string) error {
	if RedisClient == nil {
		return nil
	}

	iter := RedisClient.Scan(ctx, 0, pattern, 0).Iterator()
	var keys []string
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return err
	}

	if len(keys) > 0 {
		return RedisClient.Del(ctx, keys...).Err()
	}
	return nil
}
