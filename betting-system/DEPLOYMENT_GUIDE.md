# 部署指南

## 概述

本指南说明如何将新添加的 Category 和 Tournament 接口部署到 Railway 生产环境。

## 前置条件

- Railway 账号和项目已配置
- 数据库连接已配置 (DATABASE_URL)
- Go 1.24+ 环境

## 部署步骤

### 1. 提交代码到 Git

```bash
cd /path/to/mts-service
git add .
git commit -m "feat: add category and tournament API endpoints"
git push origin main
```

### 2. 数据库迁移

新添加的表会通过 GORM AutoMigrate 自动创建。当应用启动时，会自动执行以下迁移:

- `sports` 表
- `categories` 表
- `tournaments` 表
- `events` 表添加 `tournament_id` 字段

### 3. 插入测试数据

连接到 Railway 数据库并执行测试数据脚本:

```bash
# 方式 1: 通过 Railway CLI
railway connect postgres
\i betting-system/scripts/seed_test_data.sql

# 方式 2: 通过 psql
psql postgresql://postgres:qcriEvdpsnxvfPLaGuCuTqtivHpKoodg@turntable.proxy.rlwy.net:48608/railway -f betting-system/scripts/seed_test_data.sql
```

### 4. 验证部署

部署完成后，使用以下命令验证接口:

```bash
# 测试健康检查
curl https://betradar-uof-service-production.up.railway.app/health

# 测试获取分类
curl https://betradar-uof-service-production.up.railway.app/api/categories

# 测试获取联赛 (替换 {category_id} 为实际值)
curl "https://betradar-uof-service-production.up.railway.app/api/tournaments?category_id=1"
```

或使用测试脚本:

```bash
cd betting-system/scripts
./test_api.sh https://betradar-uof-service-production.up.railway.app
```

## 数据库 Schema 变更

### 新增表

#### sports
```sql
CREATE TABLE sports (
  id BIGSERIAL PRIMARY KEY,
  external_id VARCHAR(100) UNIQUE NOT NULL,
  name VARCHAR(255) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMP
);
```

#### categories
```sql
CREATE TABLE categories (
  id BIGSERIAL PRIMARY KEY,
  external_id VARCHAR(100) UNIQUE NOT NULL,
  sport_id BIGINT NOT NULL REFERENCES sports(id),
  name VARCHAR(255) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMP
);
```

#### tournaments
```sql
CREATE TABLE tournaments (
  id BIGSERIAL PRIMARY KEY,
  external_id VARCHAR(100) UNIQUE NOT NULL,
  category_id BIGINT NOT NULL REFERENCES categories(id),
  name VARCHAR(255) NOT NULL,
  scheduled TIMESTAMP,
  scheduled_end TIMESTAMP,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMP
);
```

### 修改表

#### events
```sql
ALTER TABLE events ADD COLUMN tournament_id BIGINT REFERENCES tournaments(id);
CREATE INDEX idx_events_tournament_id ON events(tournament_id);
```

## 环境变量

无需新增环境变量，使用现有的数据库配置即可:

```
DATABASE_URL=postgresql://postgres:qcriEvdpsnxvfPLaGuCuTqtivHpKoodg@turntable.proxy.rlwy.net:48608/railway
```

## 回滚计划

如果需要回滚，执行以下步骤:

1. 回滚代码到上一个版本
```bash
git revert HEAD
git push origin main
```

2. (可选) 删除新增的表
```sql
DROP TABLE IF EXISTS tournaments CASCADE;
DROP TABLE IF EXISTS categories CASCADE;
DROP TABLE IF EXISTS sports CASCADE;
ALTER TABLE events DROP COLUMN IF EXISTS tournament_id;
```

## 监控和日志

- 使用 Railway 控制台查看应用日志
- 监控数据库连接和查询性能
- 检查 API 响应时间和错误率

## 常见问题

### Q1: 数据库迁移失败

**A**: 检查 DATABASE_URL 是否正确配置，确保数据库连接正常。

### Q2: 接口返回空数据

**A**: 确保已执行测试数据脚本 `seed_test_data.sql`。

### Q3: 比赛数量统计不准确

**A**: 确保 `events` 表中的 `tournament_id` 字段已正确关联到 `tournaments` 表。

## 性能优化建议

1. **索引优化**: 已自动创建必要的索引
2. **查询优化**: 使用子查询统计比赛数量，避免 N+1 查询
3. **缓存**: 考虑添加 Redis 缓存热门分类和联赛数据
4. **分页**: 限制 `page_size` 最大值为 100

## 后续改进

1. 添加缓存层 (Redis)
2. 添加 API 文档 (Swagger)
3. 添加单元测试和集成测试
4. 添加 API 限流和认证
5. 优化数据库查询性能
