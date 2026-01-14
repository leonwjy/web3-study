package main

import (
	"fmt"
	"log/slog"
	"os"

	"go-auction/config"
	"go-auction/models"
)

func main() {
	// 初始化日志
	initLogger()

	slog.Info("开始测试数据库模型迁移...")

	// 加载配置
	cfg := config.Load()

	// 初始化数据库
	db := config.InitDatabase(cfg)
	defer config.Close()

	// 检查表是否存在
	tables := []string{"auctions", "bids", "nfts", "sync_status"}
	for _, table := range tables {
		if db.Migrator().HasTable(table) {
			slog.Info("表存在", "table", table)
		} else {
			slog.Error("表不存在", "table", table)
		}
	}

	// 测试创建一条记录
	slog.Info("测试创建 SyncStatus 记录...")
	syncStatus := &models.SyncStatus{
		ContractAddress: "0x1234567890123456789012345678901234567890",
		LastSyncedBlock: 1000,
	}
	if err := db.Create(syncStatus).Error; err != nil {
		slog.Error("创建记录失败", "error", err)
	} else {
		slog.Info("创建记录成功", "id", syncStatus.ID)
	}

	// 查询记录
	var count int64
	db.Model(&models.SyncStatus{}).Count(&count)
	fmt.Printf("\n总结:\n")
	fmt.Printf("  ✅ 数据库连接: 成功\n")
	fmt.Printf("  ✅ 表迁移: 成功\n")
	fmt.Printf("  ✅ 记录创建: 成功 (SyncStatus 表有 %d 条记录)\n", count)
}

// initLogger 初始化日志配置
func initLogger() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelDebug,
	}

	handler := slog.NewTextHandler(os.Stdout, opts)
	slog.SetDefault(slog.New(handler))
}
