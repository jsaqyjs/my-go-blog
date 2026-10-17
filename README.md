# My Go Blog

基于 Go、Gin、GORM、MySQL 和 Redis 实现的 RESTful 博客后端项目。

项目包含用户注册登录、JWT 身份认证、文章管理、文章权限控制和 Redis 缓存等功能，用于实践 Go Web 后端开发、数据库设计、并发编程与容器化部署。

## 技术栈

| 分类 | 技术 |
|---|---|
| 开发语言 | Go |
| Web 框架 | Gin |
| ORM | GORM |
| 数据库 | MySQL |
| 缓存 | Redis |
| 身份认证 | JWT |
| 密码加密 | bcrypt |
| 日志 | slog |
| 部署 | Docker / Docker Compose |

## 核心功能

- 用户注册、登录及个人信息查询
- 使用 bcrypt 存储密码哈希
- 使用 JWT 完成接口身份认证
- 文章创建、查询、修改与删除
- 支持分页查询文章
- 根据文章所属用户进行权限校验
- 使用 Redis 缓存文章详情
- 文章修改与删除后执行缓存失效处理
- 对不存在的文章进行短期空值缓存
- 使用 Worker Pool 控制批量任务并发度
- 使用环境变量管理服务配置
- 支持 HTTP 服务优雅退出

以上功能请以仓库中实际实现并测试通过的代码为准。

## 项目架构

```text
HTTP Request
    |
    v
Gin Router / Middleware
    |
    v
Service
    |
    +---- MySQL (GORM)
    |
    +---- Redis
    |
    v
JSON Response
```

## 主要接口

| 方法 | 接口 | 说明 |
|---|---|---|
| POST | /users/register | 用户注册 |
| POST | /users/login | 用户登录 |
| GET | /users/me | 查询当前用户 |
| GET | /articles | 文章分页列表 |
| GET | /articles/:id | 文章详情 |
| POST | /articles | 创建文章 |
| PUT | /articles/:id | 修改文章 |
| DELETE | /articles/:id | 删除文章 |
| GET | /articles/mine | 我的文章 |
| POST | /articles/batch-process | 批量处理文章 |

接口路径及请求格式应与项目实际路由定义保持一致。

## 本地运行

### 环境要求

- Go 1.25 或与 go.mod 兼容的版本
- MySQL
- Redis
- Docker Desktop（可选）

### 下载代码

```bash
git clone <你的仓库地址>
cd my-go-blog
```

### 配置环境

参考 `.env.example` 设置以下环境变量：

```text
DB_DSN
REDIS_ADDR
JWT_SECRET
SERVER_PORT
```

注意：Go 程序不会自动读取 `.env` 文件，需要使用运行配置、环境加载工具或 Docker Compose 提供对应变量。

### 运行项目

```bash
go mod download
go run .
```

### Docker Compose 启动

如果已完成容器化配置：

```bash
docker compose up -d --build
```

查看状态：

```bash
docker compose ps
```

查看日志：

```bash
docker compose logs -f app
```

## 技术设计

### JWT 鉴权

用户登录后签发 JWT，受保护接口通过中间件解析 Token，并把用户 ID 保存到当前请求的 Gin Context 中。

### 文章权限控制

修改和删除文章时，验证文章是否属于当前登录用户，防止用户操作他人的文章。

### Redis 缓存

文章详情查询优先读取 Redis。缓存未命中时查询 MySQL，并将查询结果写入 Redis。文章修改或删除后主动删除对应缓存。

### 并发任务处理

使用 Goroutine、Channel 和 WaitGroup 实现 Worker Pool，通过固定数量的 Worker 处理批量文章任务。

## 后续改进方向

- 完善单元测试和集成测试
- 增加接口限流
- 增加统一指标监控
- 优化 Redis 缓存一致性
- 增加 CI 自动化检查
- 完善压力测试和性能分析

## 项目说明

本项目主要用于 Go 后端技术学习和工程实践，尚需根据部署环境进一步完善安全、监控和运维能力。
