# Category 和 Tournament API 实现总结

## 项目信息

- **项目名称**: Betradar UOF Service
- **仓库**: https://github.com/gdszyy/mts-service
- **部署地址**: https://betradar-uof-service-production.up.railway.app
- **实现日期**: 2024-11-25

## 实现内容

### 1. 数据模型扩展

新增了三个核心数据模型以支持 SportRader 的层级结构:

#### Sport (运动类型)
- **字段**: id, external_id, name, created_at, updated_at
- **示例**: `sr:sport:1` (Football)

#### Category (分类)
- **字段**: id, external_id, sport_id, name, created_at, updated_at
- **关系**: 属于一个 Sport
- **示例**: `sr:category:1` (England)

#### Tournament (联赛)
- **字段**: id, external_id, category_id, name, scheduled, scheduled_end, created_at, updated_at
- **关系**: 属于一个 Category
- **示例**: `sr:tournament:17` (Premier League)

#### Event (比赛) - 已修改
- **新增字段**: tournament_id
- **关系**: 属于一个 Tournament

### 2. API 接口

#### 2.1 获取分类数据

**端点**: `GET /api/categories`

**功能特性**:
- ✅ 支持按体育类型过滤 (`sport_ids` 参数)
- ✅ 支持分页 (`page`, `page_size` 参数)
- ✅ 支持排序 (`sort` 参数: name_asc, name_desc, match_count_desc)
- ✅ 返回每个分类下的比赛数量统计

**查询参数**:
| 参数 | 类型 | 必填 | 默认值 | 说明 |
|:---|:---|:---:|:---|:---|
| sport_ids | string | 否 | - | 体育类型列表,逗号分隔 |
| page | integer | 否 | 1 | 页码 |
| page_size | integer | 否 | 20 | 每页大小 (最大100) |
| sort | string | 否 | name_asc | 排序方式 |

**响应示例**:
```json
{
  "data": [
    {
      "id": 1,
      "external_id": "sr:category:1",
      "sport_id": 1,
      "sport_name": "Football",
      "name": "England",
      "match_count": 150
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 50,
    "total_page": 3
  }
}
```

#### 2.2 获取联赛数据

**端点**: `GET /api/tournaments`

**功能特性**:
- ✅ 支持按分类过滤 (`category_id` 参数，必填)
- ✅ 支持分页 (`page`, `page_size` 参数)
- ✅ 支持排序 (`sort` 参数: name_asc, name_desc, match_count_desc)
- ✅ 返回每个联赛下的比赛数量统计

**查询参数**:
| 参数 | 类型 | 必填 | 默认值 | 说明 |
|:---|:---|:---:|:---|:---|
| category_id | integer | **是** | - | 分类ID |
| page | integer | 否 | 1 | 页码 |
| page_size | integer | 否 | 20 | 每页大小 (最大100) |
| sort | string | 否 | name_asc | 排序方式 |

**响应示例**:
```json
{
  "data": [
    {
      "id": 1,
      "external_id": "sr:tournament:17",
      "category_id": 1,
      "category_name": "England",
      "name": "Premier League",
      "scheduled": "2024-08-01T00:00:00Z",
      "scheduled_end": "2025-05-31T23:59:59Z",
      "match_count": 380
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 15,
    "total_page": 1
  }
}
```

### 3. 技术实现

#### 3.1 文件结构
```
betting-system/
├── internal/
│   ├── models/
│   │   └── models.go                    # 新增 Sport, Category, Tournament 模型
│   ├── handlers/
│   │   ├── category_handler.go          # 新增分类处理器
│   │   └── tournament_handler.go        # 新增联赛处理器
│   ├── routes/
│   │   └── routes.go                    # 更新路由配置
│   └── database/
│       └── database.go                  # 更新数据库迁移
├── scripts/
│   ├── seed_test_data.sql              # 测试数据脚本
│   └── test_api.sh                     # API 测试脚本
├── API_DOCUMENTATION_CATEGORIES_TOURNAMENTS.md  # API 文档
└── DEPLOYMENT_GUIDE.md                 # 部署指南
```

#### 3.2 核心技术点

**数据库查询优化**:
- 使用子查询统计比赛数量，避免 N+1 查询问题
- 使用 LEFT JOIN 关联表，确保数据完整性
- 使用 COALESCE 处理空值

**分页实现**:
- 限制 `page_size` 最大值为 100，防止过载
- 返回总页数和总记录数，方便前端实现分页控件

**排序功能**:
- 支持按名称升序/降序
- 支持按比赛数量降序（热门优先）

**错误处理**:
- 参数验证 (必填参数、类型检查)
- 资源不存在检查 (404)
- 数据库错误处理 (500)

### 4. 数据库 Schema

#### 新增表

```sql
-- sports 表
CREATE TABLE sports (
  id BIGSERIAL PRIMARY KEY,
  external_id VARCHAR(100) UNIQUE NOT NULL,
  name VARCHAR(255) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMP
);

-- categories 表
CREATE TABLE categories (
  id BIGSERIAL PRIMARY KEY,
  external_id VARCHAR(100) UNIQUE NOT NULL,
  sport_id BIGINT NOT NULL REFERENCES sports(id),
  name VARCHAR(255) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
  deleted_at TIMESTAMP
);

-- tournaments 表
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

#### 修改表

```sql
-- events 表添加 tournament_id 字段
ALTER TABLE events ADD COLUMN tournament_id BIGINT REFERENCES tournaments(id);
CREATE INDEX idx_events_tournament_id ON events(tournament_id);
```

### 5. 测试

#### 5.1 测试数据

提供了完整的测试数据脚本 (`scripts/seed_test_data.sql`):
- 3 个 Sports (Football, Basketball, Tennis)
- 5 个 Categories (England, Spain, Germany, NBA, ATP)
- 4 个 Tournaments (Premier League, La Liga, Bundesliga, NBA Regular Season)
- 5 个 Events (测试比赛)

#### 5.2 测试脚本

提供了自动化测试脚本 (`scripts/test_api.sh`):
- 健康检查测试
- 获取所有分类测试
- 按体育类型过滤测试
- 排序功能测试
- 分页功能测试
- 获取联赛测试
- 错误处理测试

### 6. 部署

#### 6.1 Git 提交

```bash
git add -A
git commit -m "feat: add category and tournament API endpoints with pagination and sorting"
git push origin main
```

#### 6.2 Railway 自动部署

- Railway 会自动检测代码变更并重新部署
- GORM AutoMigrate 会自动创建新表
- 数据库连接使用现有的 DATABASE_URL

#### 6.3 数据初始化

连接到 Railway 数据库并执行:
```bash
psql postgresql://postgres:qcriEvdpsnxvfPLaGuCuTqtivHpKoodg@turntable.proxy.rlwy.net:48608/railway -f betting-system/scripts/seed_test_data.sql
```

### 7. 使用示例

#### 示例 1: 获取所有足球分类

```bash
curl "https://betradar-uof-service-production.up.railway.app/api/categories?sport_ids=sr:sport:1"
```

#### 示例 2: 获取比赛最多的分类

```bash
curl "https://betradar-uof-service-production.up.railway.app/api/categories?sort=match_count_desc&page_size=10"
```

#### 示例 3: 获取英格兰的联赛

```bash
# 先获取英格兰分类的 ID
CATEGORY_ID=$(curl -s "https://betradar-uof-service-production.up.railway.app/api/categories?sport_ids=sr:sport:1" | jq -r '.data[] | select(.name == "England") | .id')

# 获取该分类下的联赛
curl "https://betradar-uof-service-production.up.railway.app/api/tournaments?category_id=$CATEGORY_ID"
```

### 8. 文档

提供了完整的文档:
- **API_DOCUMENTATION_CATEGORIES_TOURNAMENTS.md**: 详细的 API 接口文档
- **DEPLOYMENT_GUIDE.md**: 部署指南
- **sportradar_structure_analysis.md**: SportRader 数据结构分析

### 9. 后续优化建议

1. **性能优化**:
   - 添加 Redis 缓存层
   - 优化数据库索引
   - 实现查询结果缓存

2. **功能增强**:
   - 添加全文搜索功能
   - 支持多语言
   - 添加数据同步机制 (从 SportRader API)

3. **安全性**:
   - 添加 API 认证 (JWT)
   - 添加请求限流
   - 添加 CORS 配置

4. **监控和日志**:
   - 添加 APM 监控
   - 添加结构化日志
   - 添加错误追踪

5. **测试**:
   - 添加单元测试
   - 添加集成测试
   - 添加性能测试

## 总结

本次实现完成了两个核心接口的开发:

1. **获取分类数据接口** (`/api/categories`):
   - 支持按体育类型过滤
   - 支持分页和排序
   - 返回比赛数量统计

2. **获取联赛数据接口** (`/api/tournaments`):
   - 支持按分类过滤
   - 支持分页和排序
   - 返回比赛数量统计

所有代码已提交到 GitHub 仓库，并准备好部署到 Railway 生产环境。提供了完整的文档、测试脚本和部署指南。
