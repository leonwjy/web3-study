# NFT Auction Go 服务端 - 部署说明文档

## 目录

- [项目概述](#项目概述)
- [环境要求](#环境要求)
- [快速开始](#快速开始)
- [配置说明](#配置说明)
- [部署步骤](#部署步骤)
- [运行和维护](#运行和维护)
- [故障排查](#故障排查)

---

## 项目概述

NFT Auction Go 服务端是一个基于 Go + Gin + GORM + MySQL + Redis 构建的 NFT 拍卖系统后端 API。

### 主要功能

- NFT 拍卖管理：创建拍卖、查看拍卖列表、查看拍卖详情
- 出价管理：出价、查看出价历史、撤回出价
- 区块链事件监听：实时同步链上事件到数据库
- 数据缓存：使用 Redis 缓存热点数据
- 多环境配置：支持 local/test/prd 环境

### 技术栈

| 类型 | 技术 |
|------|------|
| Web 框架 | Gin |
| ORM | GORM |
| 数据库 | MySQL 8.0 |
| 缓存 | Redis 7 |
| 配置管理 | Viper |
| 日志 | log/slog |
| 区块链 | go-ethereum |

---

## 环境要求

### 必需环境

- **Go**: 1.21 或更高版本
- **Docker**: 20.10+ (用于本地开发)
- **Docker Compose**: 2.0+ (用于本地开发)

### 可选环境

- **MySQL**: 8.0+ (生产环境)
- **Redis**: 7.0+ (生产环境)
- **Ethereum RPC节点**: Sepolia测试网或主网

### 系统要求

- **CPU**: 2核心或以上
- **内存**: 4GB或以上
- **磁盘**: 10GB或以上可用空间
- **网络**: 可访问以太坊RPC节点

---

## 快速开始

### 1. 克隆项目

```bash
git clone <repository-url>
cd go-auction
```

### 2. 安装依赖

```bash
# 使用Makefile（推荐）
make setup

# 或手动安装
go mod download
```

### 3. 启动Docker服务

```bash
# 使用Makefile
make docker-up

# 或手动启动
docker-compose up -d
```

### 4. 配置环境变量

```bash
# 设置环境（local/test/prd）
export GO_ENV=local
```

### 5. 运行应用

```bash
# 使用Makefile
make run

# 或手动运行
go run cmd/auction/main.go
```

### 6. 验证服务

```bash
# 健康检查
curl http://localhost:8080/health

# Swagger文档
open http://localhost:8080/swagger/index.html
```

---

## 配置说明

### 配置文件位置

配置文件位于 `config/` 目录下：

- `config.local.yaml` - 本地开发环境
- `config.test.yaml` - 测试环境
- `config.prd.yaml` - 生产环境

### 环境变量

通过 `GO_ENV` 环境变量指定使用的配置文件：

```bash
export GO_ENV=local  # 使用 config.local.yaml
export GO_ENV=test   # 使用 config.test.yaml
export GO_ENV=prd    # 使用 config.prd.yaml
```

### 配置项说明

#### 服务器配置

```yaml
server:
  port: 8080          # HTTP服务端口
  mode: debug        # 运行模式：debug/release/test
```

#### 数据库配置

```yaml
database:
  host: localhost
  port: 3306
  username: root
  password: password
  dbname: auction_db
  max_idle_conns: 10    # 最大空闲连接数
  max_open_conns: 100   # 最大打开连接数
```

#### Redis配置

```yaml
redis:
  host: localhost
  port: 6380
  password: ""          # Redis密码（可选）
  db: 0                 # Redis数据库编号
```

#### 区块链配置

```yaml
blockchain:
  rpc_url: "https://sepolia.infura.io/v3/YOUR_API_KEY"
  ws_url: "wss://sepolia.infura.io/v3/YOUR_API_KEY"  # WebSocket URL（可选）
  contract_address: "0x..."  # 拍卖合约地址
  start_block: 0              # 开始同步的区块号
  chain_id: 11155111          # 链ID（Sepolia: 11155111）
```

#### 日志配置

```yaml
log:
  level: info  # 日志级别：debug/info/warn/error
```

### 生产环境配置建议

1. **数据库**：
   - 使用独立的MySQL服务器
   - 配置连接池大小（根据并发量调整）
   - 启用SSL连接

2. **Redis**：
   - 使用独立的Redis服务器
   - 配置密码认证
   - 启用持久化

3. **区块链**：
   - 使用可靠的RPC节点（如Infura、Alchemy）
   - 配置WebSocket URL以提高实时性
   - 设置合理的start_block避免同步过多历史数据

4. **服务器**：
   - 设置 `mode: release` 以提高性能
   - 配置反向代理（Nginx）
   - 启用HTTPS

---

## 部署步骤

### 方式一：使用Docker Compose（推荐用于开发/测试）

#### 1. 准备配置文件

```bash
# 复制并修改配置文件
cp config/config.local.yaml config/config.prd.yaml
vim config/config.prd.yaml
```

#### 2. 启动服务

```bash
# 启动Docker服务
make docker-up

# 设置环境变量
export GO_ENV=prd

# 运行应用
make run-prod
```

### 方式二：直接部署（生产环境）

#### 1. 准备服务器

```bash
# 安装Go
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# 安装MySQL和Redis（或使用Docker）
sudo apt-get update
sudo apt-get install mysql-server redis-server
```

#### 2. 构建应用

```bash
# 克隆代码
git clone <repository-url>
cd go-auction

# 构建
make build

# 或手动构建
go build -o bin/go-auction cmd/auction/main.go
```

#### 3. 配置数据库

```bash
# 创建数据库
mysql -u root -p
CREATE DATABASE auction_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
EXIT;

# 应用会自动运行数据库迁移（首次启动时）
```

#### 4. 配置应用

```bash
# 编辑生产环境配置
vim config/config.prd.yaml

# 设置环境变量
export GO_ENV=prd
```

#### 5. 运行应用

```bash
# 使用systemd管理（推荐）
sudo vim /etc/systemd/system/auction.service
```

systemd服务文件示例：

```ini
[Unit]
Description=NFT Auction Go Service
After=network.target mysql.service redis.service

[Service]
Type=simple
User=www-data
WorkingDirectory=/opt/go-auction
Environment="GO_ENV=prd"
ExecStart=/opt/go-auction/bin/go-auction
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

启动服务：

```bash
sudo systemctl daemon-reload
sudo systemctl enable auction
sudo systemctl start auction
sudo systemctl status auction
```

### 方式三：使用Docker（生产环境）

#### 1. 创建Dockerfile

```dockerfile
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o bin/go-auction cmd/auction/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/bin/go-auction .
COPY --from=builder /app/config ./config

CMD ["./go-auction"]
```

#### 2. 构建镜像

```bash
docker build -t go-auction:latest .
```

#### 3. 运行容器

```bash
docker run -d \
  --name go-auction \
  -p 8080:8080 \
  -e GO_ENV=prd \
  -v $(pwd)/config:/root/config \
  go-auction:latest
```

---

## 运行和维护

### 常用命令

```bash
# 查看应用日志
tail -f /var/log/auction.log  # 如果配置了日志文件
# 或使用systemd
journalctl -u auction -f

# 重启服务
sudo systemctl restart auction

# 停止服务
sudo systemctl stop auction

# 查看服务状态
sudo systemctl status auction
```

### 健康检查

```bash
# 检查服务健康状态
curl http://localhost:8080/health

# 预期响应
{
  "code": 200,
  "msg": "success",
  "data": {
    "status": "ok",
    "database": {
      "status": "ok",
      "open_connections": 2,
      "in_use": 1,
      "idle": 1
    },
    "redis": {
      "status": "ok"
    },
    "timestamp": 1234567890
  }
}
```

### 监控指标

建议监控以下指标：

1. **服务健康**：定期检查 `/health` 端点
2. **数据库连接**：监控连接池使用情况
3. **Redis连接**：监控Redis连接状态
4. **区块链同步**：监控事件同步进度
5. **API响应时间**：监控API响应时间
6. **错误率**：监控错误日志

### 日志管理

应用使用 `log/slog` 进行日志记录，日志级别可通过配置文件设置。

建议：

1. **开发环境**：使用 `debug` 级别
2. **测试环境**：使用 `info` 级别
3. **生产环境**：使用 `warn` 或 `error` 级别

日志输出位置：

- 标准输出（stdout）
- 可配置日志文件（需要额外配置）

### 数据库备份

```bash
# 备份数据库
mysqldump -u root -p auction_db > backup_$(date +%Y%m%d).sql

# 恢复数据库
mysql -u root -p auction_db < backup_20240101.sql
```

### 缓存清理

```bash
# 连接Redis
redis-cli

# 查看所有键
KEYS *

# 清空缓存（谨慎使用）
FLUSHDB
```

---

## 故障排查

### 常见问题

#### 1. 数据库连接失败

**症状**：应用启动失败，提示数据库连接错误

**排查步骤**：
```bash
# 检查MySQL服务状态
sudo systemctl status mysql

# 检查数据库配置
cat config/config.prd.yaml | grep -A 10 database

# 测试数据库连接
mysql -h localhost -u root -p -e "SHOW DATABASES;"
```

**解决方案**：
- 确认MySQL服务正在运行
- 检查数据库配置是否正确
- 确认数据库用户权限
- 检查防火墙设置

#### 2. Redis连接失败

**症状**：应用启动失败，提示Redis连接错误

**排查步骤**：
```bash
# 检查Redis服务状态
sudo systemctl status redis

# 测试Redis连接
redis-cli ping

# 检查Redis配置
cat config/config.prd.yaml | grep -A 5 redis
```

**解决方案**：
- 确认Redis服务正在运行
- 检查Redis配置（host、port、password）
- 检查防火墙设置

#### 3. 区块链RPC连接失败

**症状**：事件监听器无法连接区块链节点

**排查步骤**：
```bash
# 检查RPC URL配置
cat config/config.prd.yaml | grep rpc_url

# 测试RPC连接（需要curl）
curl -X POST -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
  https://sepolia.infura.io/v3/YOUR_API_KEY
```

**解决方案**：
- 检查RPC URL是否正确
- 确认API Key是否有效
- 检查网络连接
- 确认链ID配置正确

#### 4. 事件同步卡住

**症状**：事件监听器无法同步新区块

**排查步骤**：
```bash
# 查看应用日志
journalctl -u auction -f | grep "同步"

# 检查同步状态表
mysql -u root -p auction_db -e "SELECT * FROM sync_status;"
```

**解决方案**：
- 检查RPC节点是否正常
- 检查网络连接
- 重启事件监听器
- 手动更新同步状态（如需要）

#### 5. 端口被占用

**症状**：应用启动失败，提示端口被占用

**排查步骤**：
```bash
# 检查端口占用
lsof -i :8080
# 或
netstat -tulpn | grep 8080
```

**解决方案**：
- 修改配置文件中的端口号
- 停止占用端口的进程
- 使用不同的端口

### 日志分析

应用日志包含以下关键信息：

1. **启动日志**：服务启动、配置加载、连接初始化
2. **请求日志**：HTTP请求和响应（如果启用）
3. **事件日志**：区块链事件处理
4. **错误日志**：错误和异常信息

查看日志：

```bash
# systemd服务日志
journalctl -u auction -f

# Docker容器日志
docker logs -f go-auction

# 文件日志（如果配置）
tail -f /var/log/auction.log
```

### 性能优化

1. **数据库优化**：
   - 调整连接池大小
   - 添加适当的索引
   - 定期优化表

2. **Redis优化**：
   - 调整缓存TTL
   - 使用Redis集群（高并发场景）

3. **应用优化**：
   - 调整Worker Pool大小
   - 优化批量查询大小
   - 启用HTTP/2

---

## 联系和支持

如有问题，请查看：

- 项目README：`README.md`
- 测试文档：`TEST.md`
- 项目Issues：GitHub Issues

---

**最后更新**：2024-01-01
