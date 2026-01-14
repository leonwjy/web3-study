# Go Blog Makefile
# 使用方法: make [command]

.PHONY: all build run clean swagger migrate tidy test help

# 默认目标
all: swagger build

# 变量定义
APP_NAME := blog
MIGRATE_NAME := migrate
BUILD_DIR := ./bin
MAIN_FILE := ./cmd/blog/main.go
MIGRATE_FILE := ./cmd/migrate/main.go

# 编译应用
build:
	@echo ">>> 编译应用..."
	@go build -o $(BUILD_DIR)/$(APP_NAME) $(MAIN_FILE)
	@go build -o $(BUILD_DIR)/$(MIGRATE_NAME) $(MIGRATE_FILE)
	@echo ">>> 编译完成: $(BUILD_DIR)/$(APP_NAME), $(BUILD_DIR)/$(MIGRATE_NAME)"

# 运行应用 (开发模式)
run:
	@echo ">>> 启动应用..."
	@GO_ENV=local go run $(MAIN_FILE)

# 运行应用 (测试环境)
run-test:
	@echo ">>> 启动应用 (测试环境)..."
	@GO_ENV=test go run $(MAIN_FILE)

# 运行应用 (生产环境)
run-prd:
	@echo ">>> 启动应用 (生产环境)..."
	@GO_ENV=prd go run $(MAIN_FILE)

# 数据库迁移
migrate:
	@echo ">>> 执行数据库迁移..."
	@GO_ENV=local go run $(MIGRATE_FILE)

# 数据库迁移 (测试环境)
migrate-test:
	@echo ">>> 执行数据库迁移 (测试环境)..."
	@GO_ENV=test go run $(MIGRATE_FILE)

# 数据库迁移 (生产环境)
migrate-prd:
	@echo ">>> 执行数据库迁移 (生产环境)..."
	@GO_ENV=prd go run $(MIGRATE_FILE)

# 生成 Swagger 文档
swagger:
	@echo ">>> 生成 Swagger 文档..."
	@swag init -g $(MAIN_FILE) -o docs
	@echo ">>> Swagger 文档已生成: docs/"

# 整理依赖
tidy:
	@echo ">>> 整理依赖..."
	@go mod tidy
	@echo ">>> 依赖整理完成"

# 运行测试
test:
	@echo ">>> 运行测试..."
	@go test -v ./...

# 清理编译文件
clean:
	@echo ">>> 清理编译文件..."
	@rm -rf $(BUILD_DIR)
	@echo ">>> 清理完成"

# 一键启动 (生成文档 + 运行)
dev: swagger run

# 显示帮助信息
help:
	@echo "Go Blog Makefile 命令说明:"
	@echo ""
	@echo "  make build        - 编译应用到 bin/ 目录"
	@echo "  make run          - 运行应用 (本地开发环境)"
	@echo "  make run-test     - 运行应用 (测试环境)"
	@echo "  make run-prd      - 运行应用 (生产环境)"
	@echo "  make migrate      - 执行数据库迁移 (本地环境)"
	@echo "  make migrate-test - 执行数据库迁移 (测试环境)"
	@echo "  make migrate-prd  - 执行数据库迁移 (生产环境)"
	@echo "  make swagger      - 生成 Swagger 文档"
	@echo "  make tidy         - 整理 Go 依赖"
	@echo "  make test         - 运行测试"
	@echo "  make clean        - 清理编译文件"
	@echo "  make dev          - 开发模式 (生成文档 + 运行)"
	@echo "  make help         - 显示帮助信息"
	@echo ""
	@echo "环境变量:"
	@echo "  GO_ENV=local      - 本地开发环境 (默认)"
	@echo "  GO_ENV=test       - 测试环境"
	@echo "  GO_ENV=prd        - 生产环境"
