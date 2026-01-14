package main

import (
	"log/slog"
	"os"

	"go_blog/internal/app/model"
	"go_blog/internal/pkg/config"
	"go_blog/internal/pkg/database"
)

func main() {
	// 初始化日志
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	slog.Info("开始数据库迁移...")

	// 加载配置
	cfg := config.Load()

	// 连接数据库
	db := database.NewMySQL(cfg)

	// 执行迁移
	slog.Info("正在创建/更新表结构...")
	err := db.AutoMigrate(
		&model.User{},
		&model.Post{},
		&model.Comment{},
	)
	if err != nil {
		slog.Error("数据库迁移失败", "error", err)
		os.Exit(1)
	}

	slog.Info("数据库迁移完成!")
}
