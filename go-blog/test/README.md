# 测试用例

本目录包含 Go Blog 项目的集成测试用例。

## 测试文件说明

| 文件 | 说明 |
|------|------|
| `setup_test.go` | 测试环境初始化和辅助函数 |
| `user_test.go` | 用户模块测试（注册、登录） |
| `post_test.go` | 文章模块测试（CRUD） |
| `comment_test.go` | 评论模块测试（创建、列表） |

## 运行测试

### 前置条件

1. 确保 MySQL 数据库已启动
2. 确保 `configs/config.local.yaml` 中的数据库配置正确
3. 执行数据库迁移：`make migrate`

### 运行所有测试

```bash
make test
```

或者

```bash
GO_ENV=local go test -v ./test/...
```

### 运行单个模块测试

```bash
# 用户模块
go test -v ./test/... -run TestUser

# 文章模块
go test -v ./test/... -run TestPost

# 评论模块
go test -v ./test/... -run TestComment
```

### 运行单个测试用例

```bash
# 测试用户注册
go test -v ./test/... -run TestUserRegister

# 测试用户登录
go test -v ./test/... -run TestUserLogin

# 测试创建文章
go test -v ./test/... -run TestPostCreate
```

## 测试覆盖率

```bash
go test -v -cover ./test/...
```

## 测试用例列表

### 用户模块 (user_test.go)

| 测试用例 | 说明 |
|---------|------|
| TestUserRegister/注册成功 | 正常注册流程 |
| TestUserRegister/用户名已存在 | 重复用户名注册 |
| TestUserRegister/邮箱已存在 | 重复邮箱注册 |
| TestUserRegister/参数验证失败-用户名太短 | 用户名少于3字符 |
| TestUserRegister/参数验证失败-密码太短 | 密码少于6字符 |
| TestUserRegister/参数验证失败-邮箱格式错误 | 无效的邮箱格式 |
| TestUserLogin/登录成功 | 正常登录流程 |
| TestUserLogin/用户名不存在 | 不存在的用户 |
| TestUserLogin/密码错误 | 错误的密码 |
| TestUserLogin/参数验证失败-缺少用户名 | 缺少必填字段 |

### 文章模块 (post_test.go)

| 测试用例 | 说明 |
|---------|------|
| TestPostCreate/创建文章成功 | 正常创建文章 |
| TestPostCreate/未认证创建文章 | 未登录创建文章 |
| TestPostCreate/参数验证失败-标题为空 | 缺少标题 |
| TestPostCreate/参数验证失败-内容为空 | 缺少内容 |
| TestPostGetList/获取文章列表 | 获取列表 |
| TestPostGetList/分页参数 | 分页查询 |
| TestPostGetByID/获取文章详情成功 | 获取详情 |
| TestPostGetByID/文章不存在 | 不存在的文章 |
| TestPostGetByID/无效的文章ID | 非法ID |
| TestPostUpdate/更新文章成功 | 正常更新 |
| TestPostUpdate/未认证更新文章 | 未登录更新 |
| TestPostUpdate/非作者更新文章 | 权限校验 |
| TestPostDelete/删除文章成功 | 正常删除 |
| TestPostDelete/非作者删除文章 | 权限校验 |
| TestPostDelete/删除不存在的文章 | 不存在的文章 |

### 评论模块 (comment_test.go)

| 测试用例 | 说明 |
|---------|------|
| TestCommentCreate/创建评论成功 | 正常创建评论 |
| TestCommentCreate/未认证创建评论 | 未登录评论 |
| TestCommentCreate/参数验证失败-内容为空 | 缺少内容 |
| TestCommentCreate/文章不存在 | 对不存在的文章评论 |
| TestCommentCreate/其他用户可以评论 | 多用户评论 |
| TestCommentGetByPostID/获取评论列表成功 | 获取列表 |
| TestCommentGetByPostID/分页参数 | 分页查询 |
| TestCommentGetByPostID/文章不存在 | 不存在的文章 |
| TestCommentGetByPostID/无效的文章ID | 非法ID |
| TestCommentMultipleUsers | 多用户评论场景 |

## 注意事项

1. 测试会在数据库中创建测试数据，建议使用独立的测试数据库
2. 每个测试用例使用唯一的用户名和邮箱，避免冲突
3. 测试完成后，测试数据会保留在数据库中
