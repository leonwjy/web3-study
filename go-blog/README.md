# Go Blog - 个人博客系统后端

基于 Go + Gin + GORM + MySQL 构建的个人博客系统后端 API。

## 功能特性

- 用户认证：注册、登录（JWT Token）
- 文章管理：创建、查看、更新、删除（CRUD）
- 评论功能：发表评论、查看评论列表
- Swagger API 文档
- 多环境配置（local/test/prd）

## 技术栈

| 类型 | 技术 |
|------|------|
| Web 框架 | [Gin](https://github.com/gin-gonic/gin) |
| ORM | [GORM](https://gorm.io/) |
| 数据库 | MySQL |
| 认证 | JWT ([golang-jwt](https://github.com/golang-jwt/jwt)) |
| 配置管理 | [Viper](https://github.com/spf13/viper) |
| API 文档 | [Swagger](https://github.com/swaggo/swag) |
| 日志 | log/slog (Go 标准库) |

## 项目结构

```
go_blog/
├── api/
│   └── v1/             # API 接口层 (类似 Java Controller)
│       ├── user.go
│       ├── post.go
│       └── comment.go
├── cmd/
│   ├── blog/           # 应用入口
│   │   └── main.go
│   └── migrate/        # 数据库迁移工具
│       └── main.go
├── configs/            # 配置文件
│   ├── config.local.yaml
│   ├── config.test.yaml
│   └── config.prd.yaml
├── docs/               # Swagger 文档 (自动生成)
├── internal/
│   ├── app/
│   │   ├── dto/        # 请求参数结构体
│   │   ├── vo/         # 响应数据结构体
│   │   ├── model/      # 数据库模型
│   │   ├── service/    # 业务逻辑层
│   │   ├── repository/ # 数据访问层
│   │   ├── middleware/ # 中间件
│   │   └── router/     # 路由配置
│   └── pkg/
│       ├── config/     # 配置管理
│       ├── database/   # 数据库连接
│       ├── response/   # 统一响应格式
│       └── auth/       # JWT 认证
└── README.md
```

## 环境要求

- Go 1.21+
- MySQL 5.7+

## 快速开始

### 1. 克隆项目

```bash
git clone <repository-url>
cd go_blog
```

### 2. 安装依赖

```bash
make tidy
```

### 3. 配置数据库

编辑 `configs/config.local.yaml`，配置你的 MySQL 连接信息：

```yaml
database:
  host: localhost
  port: 3306
  username: root
  password: your_password
  dbname: go_blog
```

### 4. 创建数据库

```sql
CREATE DATABASE go_blog CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 5. 执行数据库迁移

```bash
make migrate
```

### 6. 启动服务

```bash
make run
```

服务启动后访问：
- API: http://localhost:8080
- Swagger 文档: http://localhost:8080/swagger/index.html

## Makefile 命令

```bash
make help  # 查看所有可用命令
```

### 常用命令

| 命令 | 说明 |
|------|------|
| `make run` | 运行应用（本地环境） |
| `make migrate` | 执行数据库迁移 |
| `make swagger` | 生成 Swagger 文档 |
| `make build` | 编译应用到 `bin/` 目录 |
| `make dev` | 开发模式（生成文档 + 运行） |
| `make tidy` | 整理 Go 依赖 |
| `make test` | 运行测试 |
| `make clean` | 清理编译文件 |

### 多环境支持

| 命令 | 环境 |
|------|------|
| `make run` | 本地 (local) |
| `make run-test` | 测试 (test) |
| `make run-prd` | 生产 (prd) |
| `make migrate` | 本地迁移 |
| `make migrate-test` | 测试环境迁移 |
| `make migrate-prd` | 生产环境迁移 |

## API 接口

### 用户相关

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| POST | /api/v1/register | 用户注册 | 否 |
| POST | /api/v1/login | 用户登录 | 否 |

### 文章相关

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| GET | /api/v1/posts | 获取文章列表 | 否 |
| GET | /api/v1/posts/:id | 获取文章详情 | 否 |
| POST | /api/v1/posts | 创建文章 | 是 |
| PUT | /api/v1/posts/:id | 更新文章 | 是（作者） |
| DELETE | /api/v1/posts/:id | 删除文章 | 是（作者） |

### 评论相关

| 方法 | 路径 | 描述 | 认证 |
|------|------|------|------|
| GET | /api/v1/posts/:id/comments | 获取文章评论 | 否 |
| POST | /api/v1/posts/:id/comments | 发表评论 | 是 |

## 统一响应格式

```json
{
  "code": 200,
  "msg": "success",
  "data": {}
}
```

状态码说明：
- 200: 成功
- 400: 请求参数错误
- 401: 未授权
- 403: 禁止访问
- 404: 资源不存在
- 500: 服务器内部错误

## 认证方式

请求需要认证的接口时，需要在 Header 中携带 JWT Token：

```
Authorization: Bearer <your_token>
```

## 更新 Swagger 文档

修改 API 后，重新生成 Swagger 文档：

```bash
# 安装 swag 工具（首次）
go install github.com/swaggo/swag/cmd/swag@latest

# 生成文档
make swagger
```

## 开发说明

### 架构分层

```
Handler (接收请求) 
    ↓ 使用 DTO 绑定参数
Service (业务逻辑)
    ↓ 使用 Model 操作数据
Repository (数据访问)
    ↓
Database (MySQL)
```

### 数据流转

1. Handler 接收请求，使用 DTO 绑定和验证参数
2. Service 处理业务逻辑，操作 Model
3. Repository 执行数据库操作
4. Handler 将结果转换为 VO 返回

## License

MIT License

