# User System Backend (Go + GORM + Gin)

本项目为一个用户系统的后端服务，基于 Go 语言开发，使用 Gin Web 框架和 GORM ORM，支持用户注册、登录、直播间管理、弹幕互动等功能。

## 技术栈

- Go 1.20+
- Gin
- GORM
- MySQL
- Sqlite(用于单元测试)
- Redis
- WebSocket
- Docker（可选）
- SRS（推流服务）

## 项目结构简要说明

```
├── controller     # 路由处理
├── service        # 业务逻辑封装
├── dto            # 请求/响应结构体定义
├── response       # 响应封装
├── model          # GORM 数据模型
├── middleware     # 中间件
├── router         # 路由配置
├── tests          # 单元测试
├── utils          # 工具类
├── config         # 配置加载
├── main.go        # 程序入口
└── go.mod
```

## 启动方式

```bash
go run main.go
```

## 功能模块

- 用户注册/登录（JWT 认证）
- 创建/关闭直播间
- 弹幕发送与接收（WebSocket）
- 推流回调集成（SRS）
- 分页评论（游标分页）
- Redis 缓存与限流处理


