package config

import (
	"log/slog"
	"os"
	"time"

	"go-auction/models"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DB 全局数据库连接实例
var DB *gorm.DB

// InitDatabase 初始化数据库连接
func InitDatabase(cfg *Config) *gorm.DB {
	dsn := cfg.Database.GetDSN()

	// 配置 GORM 日志级别
	var logLevel logger.LogLevel
	switch cfg.Log.Level {
	case "debug":
		logLevel = logger.Info
	case "info":
		logLevel = logger.Warn
	default:
		logLevel = logger.Error
	}

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logLevel),
	}

	db, err := gorm.Open(mysql.Open(dsn), gormConfig)
	if err != nil {
		slog.Error("数据库连接失败", "error", err, "dsn", dsn)
		os.Exit(1)
	}

	// 获取底层 sql.DB 以配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		slog.Error("获取数据库实例失败", "error", err)
		os.Exit(1)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(cfg.Database.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.Database.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// 自动迁移数据库表结构
	err = db.AutoMigrate(
		&models.Auction{},
		&models.Bid{},
		&models.NFT{},
		&models.SyncStatus{},
	)
	if err != nil {
		slog.Error("数据库迁移失败", "error", err)
		os.Exit(1)
	}

	DB = db
	slog.Info("数据库连接成功", "host", cfg.Database.Host, "database", cfg.Database.DBName)
	return db
}

// GetDB 获取数据库连接实例
func GetDB() *gorm.DB {
	return DB
}

// Close 关闭数据库连接
func Close() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
