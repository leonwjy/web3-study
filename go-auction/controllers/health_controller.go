package controllers

import (
	"context"
	"time"

	"go-auction/config"
	"go-auction/utils"

	"github.com/gin-gonic/gin"
)

// HealthController 健康检查控制器
type HealthController struct{}

// NewHealthController 创建健康检查控制器实例
func NewHealthController() *HealthController {
	return &HealthController{}
}

// HealthCheck 健康检查接口
// @Summary 健康检查
// @Description 检查服务健康状态（数据库、Redis连接）
// @Tags 系统
// @Accept json
// @Produce json
// @Success 200 {object} response.Response{data=map[string]interface{}} "服务正常"
// @Router /health [get]
func (c *HealthController) HealthCheck(ctx *gin.Context) {
	health := make(map[string]interface{})
	status := "ok"

	// 检查数据库连接
	dbStatus := checkDatabase()
	health["database"] = dbStatus
	if dbStatus["status"] != "ok" {
		status = "degraded"
	}

	// 检查Redis连接
	redisStatus := checkRedis()
	health["redis"] = redisStatus
	if redisStatus["status"] != "ok" {
		status = "degraded"
	}

	health["status"] = status
	health["timestamp"] = time.Now().Unix()

	if status == "ok" {
		utils.Success(ctx, health)
	} else {
		utils.Fail(ctx, 503, "服务降级")
	}
}

// checkDatabase 检查数据库连接
func checkDatabase() map[string]interface{} {
	result := make(map[string]interface{})

	db := config.GetDB()
	if db == nil {
		result["status"] = "error"
		result["message"] = "数据库未初始化"
		return result
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	sqlDB, err := db.DB()
	if err != nil {
		result["status"] = "error"
		result["message"] = err.Error()
		return result
	}

	if err := sqlDB.PingContext(ctx); err != nil {
		result["status"] = "error"
		result["message"] = err.Error()
		return result
	}

	stats := sqlDB.Stats()
	result["status"] = "ok"
	result["open_connections"] = stats.OpenConnections
	result["in_use"] = stats.InUse
	result["idle"] = stats.Idle
	return result
}

// checkRedis 检查Redis连接
func checkRedis() map[string]interface{} {
	result := make(map[string]interface{})

	redisClient := config.GetRedis()
	if redisClient == nil {
		result["status"] = "error"
		result["message"] = "Redis未初始化"
		return result
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := redisClient.Ping(ctx).Err(); err != nil {
		result["status"] = "error"
		result["message"] = err.Error()
		return result
	}

	result["status"] = "ok"
	return result
}
