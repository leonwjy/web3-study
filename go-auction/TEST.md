# NFT Auction Go 服务端 - 测试文档

## 目录

- [测试概述](#测试概述)
- [测试环境准备](#测试环境准备)
- [测试方案](#测试方案)
- [测试用例](#测试用例)
- [测试步骤](#测试步骤)
- [性能测试](#性能测试)
- [集成测试](#集成测试)

---

## 测试概述

### 测试目标

确保 NFT Auction Go 服务端的功能正确性、稳定性和性能满足要求。

### 测试范围

1. **功能测试**：API接口功能验证
2. **集成测试**：数据库、Redis、区块链集成验证
3. **性能测试**：并发处理能力、响应时间
4. **稳定性测试**：长时间运行、错误恢复
5. **安全测试**：输入验证、SQL注入防护

### 测试环境

- **开发环境**：本地开发测试
- **测试环境**：独立测试服务器
- **生产环境**：生产环境验证（谨慎操作）

---

## 测试环境准备

### 1. 环境要求

```bash
# Go版本
go version  # 需要 >= 1.21

# Docker版本
docker --version  # 需要 >= 20.10
docker-compose --version  # 需要 >= 2.0
```

### 2. 启动测试环境

```bash
# 克隆项目
git clone <repository-url>
cd go-auction

# 启动Docker服务（MySQL + Redis）
make docker-up
# 或
docker-compose up -d

# 验证服务启动
docker-compose ps
```

### 3. 配置测试环境

```bash
# 设置测试环境
export GO_ENV=test

# 检查配置文件
cat config/config.test.yaml
```

### 4. 准备测试数据

```bash
# 运行应用（会自动创建数据库表）
make run

# 或手动运行数据库迁移
make migrate
```

---

## 测试方案

### 测试策略

1. **单元测试**：测试各个模块的独立功能
2. **集成测试**：测试模块间的交互
3. **端到端测试**：测试完整的业务流程
4. **性能测试**：测试系统负载能力
5. **压力测试**：测试系统极限

### 测试工具

- **API测试**：curl、Postman、httpie
- **性能测试**：Apache Bench (ab)、wrk
- **数据库测试**：MySQL客户端
- **Redis测试**：redis-cli

### 测试数据准备

#### 测试NFT数据

```sql
-- 插入测试NFT
INSERT INTO nfts (contract_address, token_id, owner, name, image_url, description, created_at, updated_at)
VALUES 
  ('0x1234567890123456789012345678901234567890', 1, '0xabcdefabcdefabcdefabcdefabcdefabcdefabcd', 'Test NFT #1', '', '', NOW(), NOW()),
  ('0x1234567890123456789012345678901234567890', 2, '0xabcdefabcdefabcdefabcdefabcdefabcdefabcd', 'Test NFT #2', '', '', NOW(), NOW());
```

#### 测试拍卖数据

```sql
-- 插入测试拍卖
INSERT INTO auctions (id, nft_contract, token_id, seller, starting_price, current_highest_bid, highest_bidder, start_time, end_time, payment_token, status, created_at, updated_at)
VALUES 
  (1, '0x1234567890123456789012345678901234567890', 1, '0xabcdefabcdefabcdefabcdefabcdefabcdefabcd', '100.000000000000000000', '0', '', UNIX_TIMESTAMP(), UNIX_TIMESTAMP() + 86400, '0x0000000000000000000000000000000000000000', 'active', NOW(), NOW());
```

---

## 测试用例

### 1. 健康检查API测试

#### TC-001: 健康检查接口

**测试目标**：验证健康检查接口返回正确的服务状态

**前置条件**：
- 服务已启动
- MySQL和Redis服务正常运行

**测试步骤**：
```bash
curl -X GET http://localhost:8080/health
```

**预期结果**：
```json
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

**验证点**：
- HTTP状态码为200
- 返回JSON格式正确
- database和redis状态为"ok"

---

### 2. 拍卖API测试

#### TC-002: 获取拍卖列表

**测试目标**：验证获取拍卖列表接口

**前置条件**：
- 数据库中有测试拍卖数据

**测试步骤**：
```bash
# 获取所有拍卖
curl -X GET "http://localhost:8080/api/v1/auctions?page=1&page_size=20"

# 获取活跃拍卖
curl -X GET "http://localhost:8080/api/v1/auctions/active?page=1&page_size=20"

# 获取已结束拍卖
curl -X GET "http://localhost:8080/api/v1/auctions/ended?page=1&page_size=20"

# 按卖家筛选
curl -X GET "http://localhost:8080/api/v1/auctions?seller=0xabcdefabcdefabcdefabcdefabcdefabcdefabcd"
```

**预期结果**：
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "total": 1,
    "list": [
      {
        "id": 1,
        "nft_contract": "0x1234567890123456789012345678901234567890",
        "token_id": 1,
        "seller": "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd",
        "starting_price": "100.000000000000000000",
        "current_highest_bid": "0",
        "highest_bidder": "",
        "start_time": 1234567890,
        "end_time": 1234567890,
        "payment_token": "0x0000000000000000000000000000000000000000",
        "status": "active",
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    ]
  }
}
```

**验证点**：
- 返回正确的分页数据
- 列表数据格式正确
- 筛选功能正常

#### TC-003: 获取拍卖详情

**测试目标**：验证获取拍卖详情接口

**测试步骤**：
```bash
curl -X GET "http://localhost:8080/api/v1/auctions/1"
```

**预期结果**：
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "id": 1,
    "nft_contract": "0x1234567890123456789012345678901234567890",
    "token_id": 1,
    "seller": "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd",
    "starting_price": "100.000000000000000000",
    "current_highest_bid": "0",
    "highest_bidder": "",
    "start_time": 1234567890,
    "end_time": 1234567890,
    "payment_token": "0x0000000000000000000000000000000000000000",
    "status": "active",
    "nft": {
      "id": 1,
      "contract_address": "0x1234567890123456789012345678901234567890",
      "token_id": 1,
      "owner": "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd",
      "name": "Test NFT #1",
      "image_url": "",
      "description": "",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

**验证点**：
- 返回正确的拍卖详情
- 包含关联的NFT信息
- 数据格式正确

#### TC-004: 获取不存在的拍卖

**测试目标**：验证错误处理

**测试步骤**：
```bash
curl -X GET "http://localhost:8080/api/v1/auctions/99999"
```

**预期结果**：
```json
{
  "code": 404,
  "msg": "拍卖不存在",
  "data": null
}
```

**验证点**：
- HTTP状态码为404
- 错误消息正确

---

### 3. 出价API测试

#### TC-005: 获取拍卖的出价列表

**测试目标**：验证获取出价列表接口

**前置条件**：
- 数据库中有测试出价数据

**测试步骤**：
```bash
curl -X GET "http://localhost:8080/api/v1/bids/auction/1?page=1&page_size=20"
```

**预期结果**：
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "total": 1,
    "list": [
      {
        "id": 1,
        "auction_id": 1,
        "bidder": "0x1111111111111111111111111111111111111111",
        "bid_amount_usd": "150.000000000000000000",
        "original_amount": "150.000000000000000000",
        "payment_token": "0x0000000000000000000000000000000000000000",
        "block_number": 12345678,
        "tx_hash": "0x2222222222222222222222222222222222222222222222222222222222222222",
        "log_index": 0,
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    ]
  }
}
```

**验证点**：
- 返回正确的出价列表
- 分页功能正常

#### TC-006: 获取最高出价

**测试目标**：验证获取最高出价接口

**测试步骤**：
```bash
curl -X GET "http://localhost:8080/api/v1/bids/highest/1"
```

**预期结果**：
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "id": 1,
    "auction_id": 1,
    "bidder": "0x1111111111111111111111111111111111111111",
    "bid_amount_usd": "150.000000000000000000",
    "original_amount": "150.000000000000000000",
    "payment_token": "0x0000000000000000000000000000000000000000",
    "block_number": 12345678,
    "tx_hash": "0x2222222222222222222222222222222222222222222222222222222222222222",
    "log_index": 0,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

**验证点**：
- 返回最高出价
- 数据正确

---

### 4. NFT API测试

#### TC-007: 获取NFT列表

**测试目标**：验证获取NFT列表接口

**测试步骤**：
```bash
curl -X GET "http://localhost:8080/api/v1/nfts?owner=0xabcdefabcdefabcdefabcdefabcdefabcdefabcd&page=1&page_size=20"
```

**预期结果**：
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "total": 2,
    "list": [
      {
        "id": 1,
        "contract_address": "0x1234567890123456789012345678901234567890",
        "token_id": 1,
        "owner": "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd",
        "name": "Test NFT #1",
        "image_url": "",
        "description": "",
        "created_at": "2024-01-01T00:00:00Z",
        "updated_at": "2024-01-01T00:00:00Z"
      }
    ]
  }
}
```

**验证点**：
- 返回正确的NFT列表
- 按所有者筛选正常

#### TC-008: 获取NFT详情

**测试目标**：验证获取NFT详情接口

**测试步骤**：
```bash
# 按ID获取
curl -X GET "http://localhost:8080/api/v1/nfts/1"

# 按合约地址和TokenID获取
curl -X GET "http://localhost:8080/api/v1/nfts/by-contract?contract_address=0x1234567890123456789012345678901234567890&token_id=1"
```

**预期结果**：
```json
{
  "code": 200,
  "msg": "success",
  "data": {
    "id": 1,
    "contract_address": "0x1234567890123456789012345678901234567890",
    "token_id": 1,
    "owner": "0xabcdefabcdefabcdefabcdefabcdefabcdefabcd",
    "name": "Test NFT #1",
    "image_url": "",
    "description": "",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

**验证点**：
- 返回正确的NFT详情
- 两种查询方式都正常

---

### 5. 参数验证测试

#### TC-009: 无效参数测试

**测试目标**：验证参数验证功能

**测试步骤**：
```bash
# 无效的拍卖ID
curl -X GET "http://localhost:8080/api/v1/auctions/0"

# 无效的分页参数
curl -X GET "http://localhost:8080/api/v1/auctions?page=-1&page_size=0"

# 缺少必需参数
curl -X GET "http://localhost:8080/api/v1/nfts/by-contract"
```

**预期结果**：
```json
{
  "code": 400,
  "msg": "请求参数错误",
  "data": null
}
```

**验证点**：
- 返回400错误
- 错误消息正确

---

### 6. 缓存测试

#### TC-010: Redis缓存功能测试

**测试目标**：验证Redis缓存是否正常工作

**测试步骤**：
```bash
# 第一次请求（应该查询数据库）
time curl -X GET "http://localhost:8080/api/v1/auctions/1"

# 第二次请求（应该从缓存读取）
time curl -X GET "http://localhost:8080/api/v1/auctions/1"

# 检查Redis中的缓存
redis-cli -p 6380
> KEYS auction:*
> GET auction:1
```

**预期结果**：
- 第一次请求耗时较长（数据库查询）
- 第二次请求耗时明显减少（缓存命中）
- Redis中存在缓存键

**验证点**：
- 缓存正常工作
- 缓存命中率正常

---

## 测试步骤

### 完整测试流程

#### 1. 环境准备阶段

```bash
# 1. 启动Docker服务
make docker-up

# 2. 等待服务就绪
sleep 10

# 3. 验证服务状态
docker-compose ps
curl http://localhost:8080/health
```

#### 2. 功能测试阶段

```bash
# 1. 健康检查测试
curl -X GET http://localhost:8080/health

# 2. 拍卖API测试
curl -X GET "http://localhost:8080/api/v1/auctions?page=1&page_size=20"
curl -X GET "http://localhost:8080/api/v1/auctions/1"

# 3. 出价API测试
curl -X GET "http://localhost:8080/api/v1/bids/auction/1"
curl -X GET "http://localhost:8080/api/v1/bids/highest/1"

# 4. NFT API测试
curl -X GET "http://localhost:8080/api/v1/nfts?owner=0x..."
curl -X GET "http://localhost:8080/api/v1/nfts/1"
```

#### 3. 错误处理测试

```bash
# 1. 404错误测试
curl -X GET "http://localhost:8080/api/v1/auctions/99999"

# 2. 400错误测试
curl -X GET "http://localhost:8080/api/v1/auctions/0"

# 3. 参数验证测试
curl -X GET "http://localhost:8080/api/v1/auctions?page=-1"
```

#### 4. 性能测试阶段

```bash
# 1. 使用ab进行压力测试
ab -n 1000 -c 10 http://localhost:8080/api/v1/auctions

# 2. 使用wrk进行性能测试
wrk -t4 -c100 -d30s http://localhost:8080/api/v1/auctions
```

#### 5. 清理阶段

```bash
# 清理测试数据（可选）
# 停止服务
make docker-down
```

---

## 性能测试

### 测试指标

1. **响应时间**：API平均响应时间
2. **吞吐量**：每秒处理的请求数（QPS）
3. **并发能力**：同时处理的请求数
4. **错误率**：失败请求的比例

### 测试场景

#### 场景1：获取拍卖列表（高并发）

```bash
# 使用Apache Bench
ab -n 10000 -c 100 http://localhost:8080/api/v1/auctions

# 预期结果
# Requests per second: >= 500
# Time per request: <= 200ms
```

#### 场景2：获取拍卖详情（缓存测试）

```bash
# 第一次请求（数据库查询）
time curl -X GET "http://localhost:8080/api/v1/auctions/1"

# 缓存预热后测试
ab -n 10000 -c 100 http://localhost:8080/api/v1/auctions/1

# 预期结果
# Requests per second: >= 1000（缓存命中）
# Time per request: <= 10ms
```

#### 场景3：混合负载测试

```bash
# 使用wrk进行混合负载测试
wrk -t4 -c100 -d60s --script=test.lua http://localhost:8080

# test.lua脚本示例
request = function()
  paths = {
    "/api/v1/auctions",
    "/api/v1/auctions/1",
    "/api/v1/bids/auction/1",
    "/api/v1/nfts"
  }
  path = paths[math.random(#paths)]
  return wrk.format("GET", path)
end
```

### 性能基准

| API端点 | 目标QPS | 目标响应时间 | 缓存命中率 |
|---------|---------|--------------|------------|
| GET /api/v1/auctions | >= 500 | <= 200ms | N/A |
| GET /api/v1/auctions/:id | >= 1000 | <= 10ms | >= 80% |
| GET /api/v1/bids/auction/:id | >= 500 | <= 200ms | N/A |
| GET /api/v1/bids/highest/:id | >= 1000 | <= 10ms | >= 80% |
| GET /api/v1/nfts | >= 500 | <= 200ms | N/A |

---

## 集成测试

### 区块链事件监听测试

#### TC-011: 事件同步测试

**测试目标**：验证区块链事件是否正确同步到数据库

**前置条件**：
- 配置了有效的区块链RPC节点
- 配置了合约地址
- 合约上有事件产生

**测试步骤**：
```bash
# 1. 启动应用（事件监听器会自动启动）
make run

# 2. 观察日志
# 应该看到事件同步日志

# 3. 检查数据库
mysql -u root -p auction_db -e "SELECT * FROM auctions ORDER BY created_at DESC LIMIT 10;"
mysql -u root -p auction_db -e "SELECT * FROM bids ORDER BY created_at DESC LIMIT 10;"

# 4. 检查同步状态
mysql -u root -p auction_db -e "SELECT * FROM sync_status;"
```

**预期结果**：
- 事件监听器正常启动
- 事件正确同步到数据库
- 同步状态正确更新

**验证点**：
- 事件处理无错误
- 数据一致性正确
- 去重功能正常

### 数据库集成测试

#### TC-012: 数据库连接池测试

**测试目标**：验证数据库连接池是否正常工作

**测试步骤**：
```bash
# 1. 并发请求测试
ab -n 1000 -c 50 http://localhost:8080/api/v1/auctions

# 2. 检查数据库连接数
mysql -u root -p -e "SHOW PROCESSLIST;"

# 3. 检查连接池配置
# 查看应用日志中的连接池信息
```

**预期结果**：
- 连接数不超过配置的最大值
- 连接正常复用
- 无连接泄漏

---

## 自动化测试脚本

### 测试脚本示例

创建 `scripts/test.sh`：

```bash
#!/bin/bash

set -e

BASE_URL="http://localhost:8080"

echo "=== 开始API测试 ==="

# 健康检查
echo "1. 测试健康检查..."
curl -f -s "$BASE_URL/health" > /dev/null && echo "✓ 健康检查通过" || echo "✗ 健康检查失败"

# 拍卖列表
echo "2. 测试获取拍卖列表..."
curl -f -s "$BASE_URL/api/v1/auctions?page=1&page_size=20" > /dev/null && echo "✓ 拍卖列表测试通过" || echo "✗ 拍卖列表测试失败"

# 拍卖详情
echo "3. 测试获取拍卖详情..."
curl -f -s "$BASE_URL/api/v1/auctions/1" > /dev/null && echo "✓ 拍卖详情测试通过" || echo "✗ 拍卖详情测试失败"

echo "=== 测试完成 ==="
```

运行测试：

```bash
chmod +x scripts/test.sh
./scripts/test.sh
```

---

## 测试报告模板

### 测试执行报告

```
测试日期：2024-01-01
测试人员：XXX
测试环境：test

测试结果汇总：
- 总用例数：12
- 通过：10
- 失败：2
- 通过率：83.3%

详细结果：
[列出每个测试用例的执行结果]

问题记录：
1. TC-004: 错误消息格式不一致
2. TC-010: 缓存TTL设置不合理

性能测试结果：
- 平均响应时间：150ms
- 峰值QPS：800
- 错误率：0.1%
```

---

## 持续集成测试

### GitHub Actions示例

创建 `.github/workflows/test.yml`：

```yaml
name: Test

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    
    services:
      mysql:
        image: mysql:8.0
        env:
          MYSQL_ROOT_PASSWORD: password
          MYSQL_DATABASE: auction_db
        ports:
          - 3306:3306
      redis:
        image: redis:7-alpine
        ports:
          - 6380:6379
    
    steps:
      - uses: actions/checkout@v2
      
      - name: Set up Go
        uses: actions/setup-go@v2
        with:
          go-version: 1.21
      
      - name: Install dependencies
        run: go mod download
      
      - name: Run tests
        run: make test
      
      - name: Run linter
        run: make lint
```

---

**最后更新**：2024-01-01
