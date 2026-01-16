.PHONY: help build run test clean docker-up docker-down swagger migrate

# 项目配置
APP_NAME=auction
BINARY_NAME=go-auction
CMD_DIR=cmd/auction
BUILD_DIR=bin
GO_VERSION=1.21

# 颜色定义
GREEN=\033[0;32m
YELLOW=\033[1;33m
NC=\033[0m # No Color

help: ## 显示帮助信息
	@echo "$(GREEN)可用命令:$(NC)"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  $(YELLOW)%-15s$(NC) %s\n", $$1, $$2}'

build: ## 构建应用
	@echo "$(GREEN)构建应用...$(NC)"
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(BINARY_NAME) ./$(CMD_DIR)
	@echo "$(GREEN)构建完成: $(BUILD_DIR)/$(BINARY_NAME)$(NC)"

run: ## 运行应用（本地开发）
	@echo "$(GREEN)启动应用...$(NC)"
	@go run ./$(CMD_DIR)

run-prod: build ## 运行生产版本
	@echo "$(GREEN)启动生产版本...$(NC)"
	@GO_ENV=prd ./$(BUILD_DIR)/$(BINARY_NAME)

test: ## 运行测试
	@echo "$(GREEN)运行测试...$(NC)"
	@go test -v ./...

test-coverage: ## 运行测试并生成覆盖率报告
	@echo "$(GREEN)运行测试覆盖率...$(NC)"
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "$(GREEN)覆盖率报告已生成: coverage.html$(NC)"

clean: ## 清理构建文件
	@echo "$(GREEN)清理构建文件...$(NC)"
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@go clean

docker-up: ## 启动Docker服务（MySQL + Redis）
	@echo "$(GREEN)启动Docker服务...$(NC)"
	@docker-compose up -d
	@echo "$(GREEN)Docker服务已启动$(NC)"

docker-down: ## 停止Docker服务
	@echo "$(GREEN)停止Docker服务...$(NC)"
	@docker-compose down
	@echo "$(GREEN)Docker服务已停止$(NC)"

docker-logs: ## 查看Docker日志
	@docker-compose logs -f

docker-restart: docker-down docker-up ## 重启Docker服务

swagger: ## 生成Swagger文档
	@echo "$(GREEN)生成Swagger文档...$(NC)"
	@swag init -g cmd/auction/main.go -o docs
	@echo "$(GREEN)Swagger文档已生成$(NC)"

migrate: ## 运行数据库迁移
	@echo "$(GREEN)运行数据库迁移...$(NC)"
	@go run ./cmd/auction --migrate

lint: ## 运行代码检查
	@echo "$(GREEN)运行代码检查...$(NC)"
	@golangci-lint run ./...

fmt: ## 格式化代码
	@echo "$(GREEN)格式化代码...$(NC)"
	@go fmt ./...

vet: ## 运行go vet
	@echo "$(GREEN)运行go vet...$(NC)"
	@go vet ./...

mod-tidy: ## 整理依赖
	@echo "$(GREEN)整理依赖...$(NC)"
	@go mod tidy
	@go mod verify

mod-download: ## 下载依赖
	@echo "$(GREEN)下载依赖...$(NC)"
	@go mod download

install-tools: ## 安装开发工具
	@echo "$(GREEN)安装开发工具...$(NC)"
	@go install github.com/swaggo/swag/cmd/swag@latest
	@go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "$(GREEN)工具安装完成$(NC)"

dev: docker-up run ## 开发模式：启动Docker并运行应用

setup: mod-download install-tools docker-up ## 初始化开发环境

sync-metadata: ## 同步NFT元数据（历史数据）
	@echo "$(GREEN)同步NFT元数据...$(NC)"
	@go run ./cmd/sync-nft-metadata

sync-metadata-dry: ## 同步NFT元数据（仅显示，不更新）
	@echo "$(GREEN)同步NFT元数据（Dry Run）...$(NC)"
	@go run ./cmd/sync-nft-metadata -dry-run

all: clean build test ## 执行完整构建流程
