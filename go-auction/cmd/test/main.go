package main

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"time"

	"go-auction/config"

	"github.com/redis/go-redis/v9"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func main() {
	// 初始化日志
	initLogger()

	slog.Info("开始测试配置...")

	// 1. 测试配置加载
	slog.Info("=== 测试配置加载 ===")
	cfg := config.Load()
	if cfg == nil {
		slog.Error("配置加载失败")
		os.Exit(1)
	}

	fmt.Printf("配置加载成功:\n")
	fmt.Printf("  Server Port: %d\n", cfg.Server.Port)
	fmt.Printf("  Server Mode: %s\n", cfg.Server.Mode)
	fmt.Printf("  Database Host: %s\n", cfg.Database.Host)
	fmt.Printf("  Database Name: %s\n", cfg.Database.DBName)
	fmt.Printf("  Redis Host: %s:%d\n", cfg.Redis.Host, cfg.Redis.Port)
	fmt.Printf("  Log Level: %s\n", cfg.Log.Level)
	fmt.Printf("  Blockchain ChainID: %d\n", cfg.Blockchain.ChainID)
	fmt.Println()

	// 2. 测试数据库连接（可选，如果 Docker 未启动会跳过）
	slog.Info("=== 测试数据库连接 ===")
	fmt.Println("提示: 如果 Docker 服务未启动，数据库连接测试会失败（这是正常的）")

	// 尝试连接数据库，但不退出程序
	db, err := tryConnectDatabase(cfg)
	if db != nil {
		sqlDB, err := db.DB()
		if err == nil {
			if err := sqlDB.Ping(); err == nil {
				slog.Info("数据库连接测试成功")
				defer config.Close()
			} else {
				slog.Warn("数据库连接测试失败", "error", err, "提示", "请确保 Docker MySQL 服务已启动")
			}
		}
	} else {
		slog.Warn("数据库连接失败", "error", err, "提示", "请确保 Docker MySQL 服务已启动: docker-compose up -d")
	}
	fmt.Println()

	// 3. 测试 Redis 连接（可选，如果 Docker 未启动会跳过）
	slog.Info("=== 测试 Redis 连接 ===")
	fmt.Println("提示: 如果 Docker 服务未启动，Redis 连接测试会失败（这是正常的）")

	rdb, err := tryConnectRedis(cfg)
	if rdb != nil && err == nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		if err := rdb.Ping(ctx).Err(); err == nil {
			slog.Info("Redis 连接测试成功")

			// 测试基本操作
			testKey := "test:config"
			testValue := "test_value"
			if err := config.RedisSet(ctx, testKey, testValue, time.Minute); err == nil {
				slog.Info("Redis Set 操作成功")
			}

			if val, err := config.RedisGet(ctx, testKey); err == nil {
				slog.Info("Redis Get 操作成功", "value", val)
			}

			if exists, err := config.RedisExists(ctx, testKey); err == nil {
				slog.Info("Redis Exists 操作成功", "exists", exists)
			}

			config.RedisDelete(ctx, testKey)
			slog.Info("Redis Delete 操作成功")
			defer config.CloseRedis()
		} else {
			slog.Warn("Redis 连接测试失败", "error", err, "提示", "请确保 Docker Redis 服务已启动")
		}
	} else {
		slog.Warn("Redis 连接失败", "error", err, "提示", "请确保 Docker Redis 服务已启动: docker-compose up -d")
	}
	fmt.Println()

	slog.Info("配置测试完成！")
	fmt.Println("\n总结:")
	fmt.Println("  ✅ 配置加载: 成功")
	if db != nil {
		fmt.Println("  ✅ 数据库连接: 成功")
	} else {
		fmt.Println("  ⚠️  数据库连接: 失败（需要启动 Docker MySQL）")
	}
	if rdb != nil {
		fmt.Println("  ✅ Redis 连接: 成功")
	} else {
		fmt.Println("  ⚠️  Redis 连接: 失败（需要启动 Docker Redis）")
	}
}

// tryConnectDatabase 尝试连接数据库（不退出程序）
func tryConnectDatabase(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.Database.GetDSN()

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent), // 静默模式，避免输出错误日志
	}

	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		return nil, err
	}

	// 设置连接池参数
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	config.DB = db
	return db, nil
}

// tryConnectRedis 尝试连接 Redis（不退出程序）
func tryConnectRedis(cfg *config.Config) (*redis.Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.GetAddr(),
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	config.RedisClient = rdb
	return rdb, nil
}

// initLogger 初始化日志配置
func initLogger() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	handler := slog.NewTextHandler(os.Stdout, opts)
	slog.SetDefault(slog.New(handler))
}
