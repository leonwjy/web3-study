package main

import (
	"fmt"
	"log/slog"
	"os"

	"go_blog/internal/app/router"
	"go_blog/internal/pkg/config"
	"go_blog/internal/pkg/database"

	_ "go_blog/docs" // 导入 Swagger 文档
)

// @title Go Blog API
// @version 1.0
// @description Go 个人博客系统 API 文档
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description 输入格式: Bearer {token}

func main() {
	// 初始化日志
	initLogger()

	slog.Info("启动 Go Blog 服务...")

	// 加载配置
	cfg := config.Load()

	// 初始化数据库连接
	database.NewMySQL(cfg)
	defer database.Close()

	// 配置路由
	r := router.SetupRouter(cfg.Server.Mode)

	// 启动服务
	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	slog.Info("服务启动成功", "addr", addr, "swagger", fmt.Sprintf("http://localhost:%d/swagger/index.html", cfg.Server.Port))

	if err := r.Run(addr); err != nil {
		slog.Error("服务启动失败", "error", err)
		os.Exit(1)
	}
}

// initLogger 初始化日志配置
func initLogger() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	// 根据环境变量设置日志级别
	env := os.Getenv("GO_ENV")
	if env == "local" || env == "" {
		opts.Level = slog.LevelDebug
	}

	handler := slog.NewTextHandler(os.Stdout, opts)
	slog.SetDefault(slog.New(handler))
}

