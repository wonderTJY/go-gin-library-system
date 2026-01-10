# go-gin-library-system

一个使用 Go + Gin + GORM 实现的简易图书馆管理练习项目，用来练习：

- 基本的 RESTful API 设计（图书 / 学生增删改查、借书 / 还书）
- GORM 模型关系设计（Book、Student、借阅中间表、库存 BookCopy）
- Gin 中间件（日志、Request ID、中间件链路）
- Git 基本工作流（本地开发、提交、rebase、推送到 GitHub）

---

## 技术栈

- Go 1.22+
- Gin Web Framework
- GORM
- SQLite（本地开发环境的简单数据库）

---

## 运行环境准备

1. 安装 Go

   - 到官网下载安装包并安装：<https://go.dev/dl/>
   - 安装完成后，在终端执行确认版本：

     ```bash
     go version
     ```

2. 克隆仓库

   ```bash
   git clone https://github.com/wonderTJY/go-gin-library-system.git
   cd go-gin-library-system
   ```

3. 安装依赖

   在项目根目录执行：

   ```bash
   go mod tidy
   ```

   用于根据 go.mod 自动下载并整理依赖。

---

## 启动项目

在项目根目录执行：

```bash
go run main.go
```

默认会：

- 初始化 SQLite 数据库并自动迁移模型
- 在本地启动 HTTP 服务：`http://127.0.0.1:8080`

看到类似输出（含 Gin 日志和自定义日志中间件）说明启动成功。

---

## 主要功能

### 图书相关 API

前缀：`/api/v1/books`

- `GET /api/v1/books`  
  列出所有图书。

- `GET /api/v1/books/:id`  
  根据 ID 获取单本图书详情。

- `POST /api/v1/books`  
  创建图书，示例请求体：

  ```json
  {
    "title": "The Go Programming Language",
    "author": "Alan A. A. Donovan",
    "stock": 10
  }
  ```

- `PUT /api/v1/books/:id`  
  更新图书信息（目前以“全量更新”方式处理），示例请求体：

  ```json
  {
    "title": "The Go Programming Language (2nd Edition)",
    "author": "Alan A. A. Donovan",
    "stock": 12
  }
  ```

- `DELETE /api/v1/books/:id`  
  删除图书。

### 学生相关 API

前缀：`/api/v1/students`

- `GET /api/v1/students`  
  列出所有学生。

- `GET /api/v1/students/:id`  
  根据 ID 获取学生详情。

- `POST /api/v1/students`  
  创建学生，示例请求体：

  ```json
  {
    "name": "Tom",
    "email": "tom@example.com"
  }
  ```

- `PUT /api/v1/students/:id`  
  更新学生信息（目前以“全量更新”方式处理），示例请求体：

  ```json
  {
    "name": "Tom Updated",
    "email": "tom.updated@example.com"
  }
  ```

- `DELETE /api/v1/students/:id`  
  删除学生。

### 学生借阅相关 API

前缀：`/api/v1/students`

- `GET /api/v1/students/:id/books`  
  查询某个学生的借阅记录。

- `POST /api/v1/students/:student_id/books/:book_id/borrow`  
  借书：

  - 检查学生是否存在
  - 检查图书是否存在且库存 `stock > 0`
  - 在中间表 `Book_Student` 记录一条借阅记录（状态为 `borrowed`）
  - 图书库存 `stock - 1`

- `POST /api/v1/students/:student_id/books/:book_id/return`  
  还书：

  - 根据 `student_id + book_id` 查找借阅记录（状态为 `borrowed`）
  - 将状态更新为 `returned`，记录归还时间
  - 图书库存 `stock + 1`

---

## 中间件

项目在 `middleware` 目录中实现并全局挂载了几个基础中间件（在 `router/router.go` 中统一配置）：

- Request ID 中间件（`requestID.go`）：
  - 为每个请求生成或透传 `X-Request-ID` 请求头
  - 将 `request_id` 写入 Gin 上下文，供日志等使用
- 请求计数中间件（`RequestCount.go`）：
  - 使用 `sync/atomic` 对全局请求总数做并发安全自增
  - 当前计数通过 `request_count` 存入 Gin 上下文
- 日志中间件（`logging.go`）：
  - 记录请求时间
  - HTTP 方法（GET/POST/...）
  - 路径（URL Path）
  - 状态码
  - 耗时（latency）
  - 全局请求计数（`request_count`）
  - Request ID（从 header 中读取 `X-Request-ID`）

日志输出到终端，配合 Gin 默认的访问日志，可以清楚看到每个请求的处理情况。

---

## 开发建议

- 修改代码后，可以先本地检查：

  ```bash
  go build ./...
  go test ./...   # 如果后续添加了测试
  ```

- 本地确认无误后：

  ```bash
  git add <改动的文件>
  git commit -m "feat: xxx"  # 或 fix: xxx
  git push
  ```

- 推荐一类改动一个提交，提交信息尽量说明“做了什么”和“为什么做”。

---

## TODO（练习方向）

- [x] 增加 Request ID 中间件，并在日志中打印 request_id
- [ ] 增加简单 Token 校验中间件，保护部分路由
- [ ] 增加统一错误响应中间件，规范错误返回结构
- [ ] 为核心 handler 编写单元测试 / 集成测试

