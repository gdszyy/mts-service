# MTS Service & Betting System - 项目总览

本仓库包含两个独立但相关的模块：

## 1. MTS Service（Sportradar MTS 集成服务）

位于根目录，提供与 Sportradar MTS（Managed Trading Services）的集成。

### 功能特性
- ✅ MTS WebSocket 连接管理
- ✅ OAuth 2.0 认证
- ✅ UOF API `whoami.xml` 自动配置
- ✅ 票据放置（Ticket Placement）
- ✅ 环境变量配置（支持 CI/生产环境切换）

### 文档
- [API 文档](./API_DOCUMENTATION.md)
- [Railway 部署指南](./RAILWAY_DEPLOYMENT.md)
- [README](./README.md)

### 状态
⚠️ **需要有效的 Sportradar MTS 凭证才能运行**

---

## 2. Betting System（独立投注管理系统）

位于 `betting-system/` 目录，提供完整的投注管理功能，**不依赖 MTS**。

### 功能特性
- ✅ 支持 12 种复杂投注类型（Trixie、Yankee、Patent、Lucky 系列等）
- ✅ Banker 支持（所有系统投注类型）
- ✅ 用户管理（创建、充值、提现）
- ✅ 赛事和盘口管理
- ✅ 投注结算和回报计算
- ✅ PostgreSQL 数据库持久化
- ✅ RESTful API

### 文档
- [Betting System README](./betting-system/README.md)
- [数据库架构](./betting-system/database_schema.md)

### 状态
✅ **完全独立运行，已通过测试**

---

## 快速开始

### MTS Service
```bash
# 设置环境变量
cp .env.example .env
vi .env

# 运行服务
go run cmd/server/main.go
```

### Betting System
```bash
# 进入 betting-system 目录
cd betting-system

# 设置环境变量
cp .env.example .env
vi .env

# 运行服务
go run cmd/server/main.go
```

---

## 技术栈

| 模块 | 技术栈 |
| :--- | :--- |
| **MTS Service** | Go, Gorilla WebSocket, OAuth 2.0 |
| **Betting System** | Go, Gin, GORM, PostgreSQL |

---

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License

