# NFT Auction Go 服务端

基于 Go + Gin + GORM + MySQL + Redis 构建的 NFT 拍卖系统后端 API。

## 功能特性

- NFT 拍卖管理：创建拍卖、查看拍卖列表、查看拍卖详情
- 出价管理：出价、查看出价历史、撤回出价
- 区块链事件监听：实时同步链上事件到数据库
- 数据缓存：使用 Redis 缓存热点数据
- 多环境配置：支持 local/test/prd 环境

## 技术栈

| 类型 | 技术 |
|------|------|
| Web 框架 | [Gin](https://github.com/gin-gonic/gin) |
| ORM | [GORM](https://gorm.io/) |
| 数据库 | MySQL 8.0 |
| 缓存 | Redis 7 |
| 配置管理 | [Viper](https://github.com/spf13/viper) |
| 日志 | log/slog (Go 标准库) |
| 区块链 | go-ethereum (ethclient, abigen) |
| 认证 | JWT (可选) |

## 项目结构

```
go-auction/
├── cmd/
│   └── auction/
│       └── main.go              # 应用入口
├── config/
│   ├── config.go                # Viper 配置管理
│   ├── database.go              # 数据库配置
│   ├── redis.go                 # Redis 配置
│   ├── config.local.yaml        # 本地环境配置
│   ├── config.test.yaml         # 测试环境配置
│   └── config.prd.yaml          # 生产环境配置
├── controllers/                 # 控制器层
├── middleware/                  # 中间件
├── models/                      # 数据模型（GORM）
├── routes/                      # 路由配置
├── services/                    # 业务逻辑层
├── repositories/                # 数据访问层
├── dto/                         # 数据传输对象
├── vo/                          # 视图对象
├── blockchain/                 # 区块链交互层
│   ├── client.go
│   ├── event_listener.go
│   ├── event_handler.go
│   └── contract/               # 合约绑定代码（abigen 生成）
├── utils/                       # 工具函数
├── docker-compose.yml           # Docker 环境
├── .env.example                 # 环境变量示例
└── README.md
```

## 环境要求

- Go 1.21+
- MySQL 8.0+ (或使用 Docker)
- Redis 7+ (或使用 Docker)
- Docker & Docker Compose (可选，用于本地开发环境)

## 快速开始

### 1. 克隆项目

```bash
git clone <repository-url>
cd go-auction
```

### 2. 配置环境变量

```bash
cp .env.example .env
# 编辑 .env 文件，填入你的配置信息
```

### 3. 启动 Docker 环境（MySQL + Redis）

```bash
docker-compose up -d
```

### 4. 安装依赖

```bash
go mod tidy
```

### 5. 运行项目

```bash
go run cmd/auction/main.go
```

## 配置说明

### 环境变量

- `GO_ENV`: 运行环境 (local/test/prd)，默认 local
- `PORT`: 服务端口，默认 8080
- `DB_*`: MySQL 数据库配置
- `REDIS_*`: Redis 配置
- `INFURA_PROJECT_ID`: Infura 项目 ID（区块链 RPC）
- `BLOCKCHAIN_CONTRACT_ADDRESS`: Auction 合约地址

### 配置文件

配置文件位于 `config/` 目录，根据 `GO_ENV` 环境变量自动加载对应的配置文件：
- `config.local.yaml` - 本地开发环境
- `config.test.yaml` - 测试环境
- `config.prd.yaml` - 生产环境

## 开发指南

### 数据库迁移

使用 GORM 的 AutoMigrate 功能自动创建表结构：

```go
// 在 main.go 中初始化数据库时会自动迁移
```

### 区块链合约绑定代码生成

```bash
# 从 hardhat-auction 编译产物生成 Go 绑定代码
abigen --abi=../hardhat-auction/artifacts/contracts/Auction.sol/Auction.json \
       --pkg=contract \
       --type=Auction \
       --out=blockchain/contract/auction.go
```

## API 文档

（待补充 Swagger 文档）

## License

MIT
