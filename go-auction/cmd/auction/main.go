package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go-auction/blockchain"
	"go-auction/config"
	"go-auction/repositories"
	"go-auction/routes"

	_ "go-auction/docs" // 导入 Swagger 文档
)

// @title NFT Auction API
// @version 1.0
// @description NFT 拍卖系统后端 API 文档
// @termsOfService http://swagger.io/terms/

// @contact.name API Support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8080
// @BasePath /

func main() {
	// 初始化日志
	initLogger()

	slog.Info("启动 NFT Auction 服务端...")

	// 加载配置
	cfg := config.Load()
	env := os.Getenv("GO_ENV")
	if env == "" {
		env = "local"
	}
	slog.Info("配置加载成功", "env", env)

	// 初始化数据库
	config.InitDatabase(cfg)
	defer config.Close()
	slog.Info("数据库初始化成功")

	// 初始化Redis
	config.InitRedis(cfg)
	defer config.CloseRedis()
	slog.Info("Redis初始化成功")

	// 初始化区块链客户端和事件监听器（如果配置了合约地址）
	var eventListener *blockchain.EventListener
	if cfg.Blockchain.ContractAddress != "" && cfg.Blockchain.ContractAddress != "0x..." {
		bcClient, err := blockchain.NewClient(&cfg.Blockchain)
		if err != nil {
			slog.Error("区块链客户端初始化失败", "error", err)
		} else {
			defer bcClient.Close()
			slog.Info("区块链客户端初始化成功")

			// 创建 Repository 实例
			syncRepo := repositories.NewSyncStatusRepository()
			auctionRepo := repositories.NewAuctionRepository()
			bidRepo := repositories.NewBidRepository()
			nftRepo := repositories.NewNFTRepository()

			// 创建事件监听器
			eventListener, err = blockchain.NewEventListener(
				bcClient,
				syncRepo,
				auctionRepo,
				bidRepo,
				nftRepo,
				cfg.Blockchain.StartBlock,
			)
			if err != nil {
				slog.Error("事件监听器创建失败", "error", err)
			} else {
				// 启动事件监听（在后台 goroutine）
				go func() {
					if err := eventListener.Start(context.Background()); err != nil {
						slog.Error("事件监听器启动失败", "error", err)
					}
				}()
				slog.Info("事件监听器启动成功")
			}
		}
	} else {
		slog.Warn("未配置合约地址，跳过区块链事件监听")
	}

	// 设置路由
	router := routes.SetupRouter()
	slog.Info("路由配置完成")

	// 启动HTTP服务器
	port := cfg.Server.Port
	if port == 0 {
		port = 8080
	}

	srv := &http.Server{
		Addr:    fmt.Sprintf(":%d", port),
		Handler: router,
	}

	// 在goroutine中启动服务器
	go func() {
		slog.Info("HTTP服务器启动", "port", port, "swagger", fmt.Sprintf("http://localhost:%d/swagger/index.html", port))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("HTTP服务器启动失败", "error", err)
			os.Exit(1)
		}
	}()

	// 等待中断信号以优雅地关闭服务器
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	slog.Info("正在关闭服务器...")

	// 停止事件监听器
	if eventListener != nil {
		eventListener.Stop()
		slog.Info("事件监听器已停止")
	}

	// 设置5秒超时上下文
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// 优雅关闭服务器
	if err := srv.Shutdown(ctx); err != nil {
		slog.Error("服务器关闭失败", "error", err)
		os.Exit(1)
	}

	slog.Info("服务器已关闭")
}

// initLogger 初始化日志配置
func initLogger() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}

	var handler slog.Handler
	if os.Getenv("GO_ENV") == "prd" {
		// 生产环境使用JSON格式
		handler = slog.NewJSONHandler(os.Stdout, opts)
	} else {
		// 开发环境使用文本格式
		handler = slog.NewTextHandler(os.Stdout, opts)
	}

	slog.SetDefault(slog.New(handler))
}
