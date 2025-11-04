# Betting System - 投注管理系统

一个功能完整的投注管理系统，支持复杂投注类型（Banker、Trixie、Yankee、Patent、Lucky 系列等）和数据库持久化。

## 功能特性

### 投注类型支持

#### 基础投注类型
- **Single**: 单注
- **Accumulator**: 串关（2选项及以上）

#### Full Cover Bets（全覆盖投注）
- **Trixie**: 3选项 → 3 doubles + 1 treble = 4注
- **Yankee**: 4选项 → 6 doubles + 4 trebles + 1 4-fold = 11注
- **Super Yankee (Canadian)**: 5选项 → 26注
- **Heinz**: 6选项 → 57注
- **Super Heinz**: 7选项 → 120注
- **Goliath**: 8选项 → 247注

#### Full Cover with Singles（包含单注的全覆盖）
- **Patent**: 3选项 → 3 singles + 3 doubles + 1 treble = 7注
- **Lucky 15**: 4选项 → 15注
- **Lucky 31**: 5选项 → 31注
- **Lucky 63**: 6选项 → 63注

### Banker 支持

所有系统投注类型都支持 Banker 标记：
- Banker 选项会自动包含在所有组合中
- 非 Banker 选项进行正常组合
- 例如：4个选项（1个Banker + 3个普通），Trixie = 所有组合都包含Banker

### 核心功能

- ✅ 用户管理（创建、查询、充值、提现）
- ✅ 赛事管理（创建、查询、状态更新）
- ✅ 盘口管理（创建、更新赔率）
- ✅ 投注管理（下注、查询、结算）
- ✅ 自动投注组合生成
- ✅ 余额验证和扣除
- ✅ 投注结算和回报计算
- ✅ 数据库事务支持

## 技术栈

- **语言**: Go 1.21+
- **Web框架**: Gin
- **ORM**: GORM
- **数据库**: PostgreSQL
- **架构**: RESTful API

## 快速开始

### 1. 安装依赖

```bash
# 安装 PostgreSQL
# Ubuntu/Debian
sudo apt-get install postgresql postgresql-contrib

# macOS
brew install postgresql

# 启动 PostgreSQL
sudo service postgresql start  # Linux
brew services start postgresql # macOS
```

### 2. 创建数据库

```bash
# 登录 PostgreSQL
sudo -u postgres psql

# 创建数据库
CREATE DATABASE betting_system;

# 创建用户（可选）
CREATE USER betting_user WITH PASSWORD 'your_password';
GRANT ALL PRIVILEGES ON DATABASE betting_system TO betting_user;

# 退出
\q
```

### 3. 配置环境变量

```bash
# 复制环境变量示例文件
cp .env.example .env

# 编辑 .env 文件，填入数据库配置
vi .env
```

### 4. 安装 Go 依赖

```bash
go mod download
```

### 5. 运行服务

```bash
go run cmd/server/main.go
```

服务将在 `http://localhost:8080` 启动。

## API 文档

### 用户管理

#### 创建用户
```bash
POST /api/users
Content-Type: application/json

{
  "username": "john_doe",
  "email": "john@example.com",
  "balance": 1000.00,
  "currency": "CNY"
}
```

#### 获取用户详情
```bash
GET /api/users/:id
```

#### 充值
```bash
POST /api/users/:id/deposit
Content-Type: application/json

{
  "amount": 500.00
}
```

#### 提现
```bash
POST /api/users/:id/withdraw
Content-Type: application/json

{
  "amount": 200.00
}
```

### 赛事管理

#### 创建赛事
```bash
POST /api/events
Content-Type: application/json

{
  "external_id": "sr:match:12345",
  "sport_id": "sr:sport:1",
  "home_team": "Manchester United",
  "away_team": "Liverpool",
  "start_time": "2025-11-01T15:00:00Z"
}
```

#### 获取赛事列表
```bash
GET /api/events?status=scheduled&limit=20&offset=0
```

#### 获取赛事详情
```bash
GET /api/events/:id
```

### 盘口管理

#### 创建盘口
```bash
POST /api/markets
Content-Type: application/json

{
  "event_id": 1,
  "market_type": "1x2",
  "specifier": ""
}
```

### 结果选项管理

#### 创建结果选项
```bash
POST /api/outcomes
Content-Type: application/json

{
  "market_id": 1,
  "outcome_id": "1",
  "odds": 2.50
}
```

#### 更新赔率
```bash
PUT /api/outcomes/:id/odds
Content-Type: application/json

{
  "odds": 2.75
}
```

#### 结算结果选项
```bash
POST /api/outcomes/:id/settle
Content-Type: application/json

{
  "status": "won"
}
```

### 投注管理

#### 下注
```bash
POST /api/bets
Content-Type: application/json

{
  "user_id": 1,
  "bet_type": "trixie",
  "unit_stake": 10.00,
  "selections": [
    {
      "outcome_id": 1,
      "odds": 2.50,
      "is_banker": true
    },
    {
      "outcome_id": 2,
      "odds": 3.00,
      "is_banker": false
    },
    {
      "outcome_id": 3,
      "odds": 2.20,
      "is_banker": false
    }
  ]
}
```

#### 获取投注详情
```bash
GET /api/bets/:id
```

#### 获取用户投注列表
```bash
GET /api/users/:user_id/bets?limit=20&offset=0
```

#### 结算投注
```bash
POST /api/bets/:id/settle
```

## 投注类型示例

### Trixie 示例

**选项**:
- 选项 A: 赔率 2.0
- 选项 B: 赔率 3.0
- 选项 C: 赔率 2.5

**单位投注**: 10元

**生成的组合**:
1. Double (A+B): 10元 × 2.0 × 3.0 = 60元
2. Double (A+C): 10元 × 2.0 × 2.5 = 50元
3. Double (B+C): 10元 × 3.0 × 2.5 = 75元
4. Treble (A+B+C): 10元 × 2.0 × 3.0 × 2.5 = 150元

**总投注**: 40元
**潜在回报**: 335元

### Trixie with Banker 示例

**选项**:
- 选项 A (Banker): 赔率 2.0
- 选项 B: 赔率 3.0
- 选项 C: 赔率 2.5

**单位投注**: 10元

**生成的组合**:
1. Double (A+B): 10元 × 2.0 × 3.0 = 60元
2. Double (A+C): 10元 × 2.0 × 2.5 = 50元
3. Treble (A+B+C): 10元 × 2.0 × 3.0 × 2.5 = 150元

**总投注**: 30元（注意：由于A是Banker，B+C的组合不会生成）
**潜在回报**: 260元

## 数据库架构

详细的数据库架构设计请参考 [database_schema.md](database_schema.md)。

## 项目结构

```
betting-system/
├── cmd/
│   └── server/
│       └── main.go              # 主程序入口
├── internal/
│   ├── database/
│   │   └── database.go          # 数据库连接和迁移
│   ├── engine/
│   │   └── bet_engine.go        # 投注引擎（组合生成）
│   ├── handlers/
│   │   ├── bet_handler.go       # 投注处理器
│   │   ├── event_handler.go     # 赛事处理器
│   │   └── user_handler.go      # 用户处理器
│   ├── models/
│   │   └── models.go            # 数据模型
│   ├── routes/
│   │   └── routes.go            # 路由配置
│   └── service/
│       └── bet_service.go       # 投注服务
├── .env.example                 # 环境变量示例
├── go.mod                       # Go 模块定义
├── go.sum                       # Go 依赖锁定
└── README.md                    # 项目文档
```

## 开发计划

- [ ] 添加单元测试
- [ ] 添加集成测试
- [ ] 添加 API 文档（Swagger）
- [ ] 添加日志系统
- [ ] 添加监控和指标
- [ ] 添加缓存层（Redis）
- [ ] 添加消息队列（RabbitMQ/Kafka）
- [ ] 添加 Docker 支持
- [ ] 添加 CI/CD 配置

## 许可证

MIT License

## 贡献

欢迎提交 Issue 和 Pull Request！

