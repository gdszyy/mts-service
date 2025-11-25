# Events API - Market Types 过滤功能

## 概述

`/api/events` 接口现在支持通过 `market_types` 参数指定返回的 market 类型，而不是总是返回所有 market。

## 接口详情

### 获取赛事列表

**请求**

```
GET /api/events
```

**查询参数**

| 参数 | 类型 | 必填 | 说明 | 示例 |
|:---|:---|:---:|:---|:---|
| `status` | string | 否 | 状态筛选 | `scheduled`, `live`, `finished` |
| `limit` | integer | 否 | 每页数量，默认 20 | `20` |
| `offset` | integer | 否 | 偏移量，默认 0 | `0` |
| `market_types` | string | 否 | Market 类型列表，逗号分隔 | `1x2,handicap,totals` |

### 行为说明

#### 1. 不传 `market_types` 参数

当不传 `market_types` 参数时，接口**不会返回** `markets` 字段，保持原有的轻量级响应。

**请求示例**:
```bash
curl "https://betradar-uof-service-production.up.railway.app/api/events?limit=2"
```

**响应示例**:
```json
{
  "events": [
    {
      "id": 1,
      "external_id": "sr:match:1001",
      "sport_id": "sr:sport:1",
      "home_team": "Manchester United",
      "away_team": "Liverpool",
      "start_time": "2024-12-01T15:00:00Z",
      "status": "scheduled"
    }
  ],
  "total": 5,
  "limit": 2,
  "offset": 0
}
```

#### 2. 传入 `market_types` 参数

当传入 `market_types` 参数时，接口会返回指定类型的 markets 及其 outcomes。

**请求示例**:
```bash
# 只返回 1x2 market
curl "https://betradar-uof-service-production.up.railway.app/api/events?limit=2&market_types=1x2"

# 返回多个 market 类型
curl "https://betradar-uof-service-production.up.railway.app/api/events?limit=2&market_types=1x2,handicap,totals"
```

**响应示例**:
```json
{
  "events": [
    {
      "id": 1,
      "external_id": "sr:match:1001",
      "sport_id": "sr:sport:1",
      "home_team": "Manchester United",
      "away_team": "Liverpool",
      "start_time": "2024-12-01T15:00:00Z",
      "status": "scheduled",
      "markets": [
        {
          "id": 1,
          "event_id": 1,
          "market_type": "1x2",
          "specifier": null,
          "status": "active",
          "outcomes": [
            {
              "id": 1,
              "market_id": 1,
              "outcome_id": "1",
              "odds": 2.50,
              "status": "active"
            },
            {
              "id": 2,
              "market_id": 1,
              "outcome_id": "x",
              "odds": 3.20,
              "status": "active"
            },
            {
              "id": 3,
              "market_id": 1,
              "outcome_id": "2",
              "odds": 2.80,
              "status": "active"
            }
          ]
        },
        {
          "id": 2,
          "event_id": 1,
          "market_type": "handicap",
          "specifier": "handicap=-1.5",
          "status": "active",
          "outcomes": [
            {
              "id": 4,
              "market_id": 2,
              "outcome_id": "1",
              "odds": 1.85,
              "status": "active"
            },
            {
              "id": 5,
              "market_id": 2,
              "outcome_id": "2",
              "odds": 1.95,
              "status": "active"
            }
          ]
        }
      ]
    }
  ],
  "total": 5,
  "limit": 2,
  "offset": 0
}
```

## 常见 Market 类型

| Market Type | 说明 | 示例 Specifier |
|:---|:---|:---|
| `1x2` | 胜平负 | - |
| `handicap` | 让球盘 | `handicap=-1.5` |
| `totals` | 大小球 | `total=2.5` |
| `both_teams_to_score` | 双方都进球 | - |
| `correct_score` | 正确比分 | - |
| `double_chance` | 双重机会 | - |
| `first_half_1x2` | 上半场胜平负 | - |
| `over_under` | 进球数大小 | `total=2.5` |

## 使用场景

### 场景 1: 只获取赛事基本信息（不需要 markets）

```bash
curl "https://betradar-uof-service-production.up.railway.app/api/events?status=scheduled&limit=10"
```

**优点**: 响应快，数据量小，适合列表展示。

### 场景 2: 获取赛事及特定 market 信息

```bash
# 只获取 1x2 赔率
curl "https://betradar-uof-service-production.up.railway.app/api/events?status=scheduled&limit=10&market_types=1x2"

# 获取多个 market 类型
curl "https://betradar-uof-service-production.up.railway.app/api/events?status=scheduled&limit=10&market_types=1x2,handicap,totals"
```

**优点**: 按需加载，减少不必要的数据传输。

### 场景 3: 结合分页使用

```bash
# 第一页，每页 20 条，只返回 1x2 market
curl "https://betradar-uof-service-production.up.railway.app/api/events?limit=20&offset=0&market_types=1x2"

# 第二页
curl "https://betradar-uof-service-production.up.railway.app/api/events?limit=20&offset=20&market_types=1x2"
```

## 性能优化

### 1. 按需加载

- 不传 `market_types` 参数时，不会加载任何 market 数据
- 只传需要的 market 类型，避免加载不必要的数据

### 2. 数据库查询优化

- 使用 `IN` 查询批量过滤 market 类型
- 使用 `Preload` 预加载关联数据，避免 N+1 查询

### 3. 建议

- 列表页面：不传 `market_types`，只显示赛事基本信息
- 详情页面：传入需要的 `market_types`，显示完整赔率信息
- 限制 `limit` 参数，避免一次性加载过多数据

## 错误处理

### 无效的 market_types

如果传入的 `market_types` 在数据库中不存在，接口会返回空的 `markets` 数组，不会报错。

```bash
curl "https://betradar-uof-service-production.up.railway.app/api/events?limit=1&market_types=invalid_market_type"
```

响应:
```json
{
  "events": [
    {
      "id": 1,
      "external_id": "sr:match:1001",
      "home_team": "Manchester United",
      "away_team": "Liverpool",
      "markets": []
    }
  ],
  "total": 5,
  "limit": 1,
  "offset": 0
}
```

## 测试

使用提供的测试脚本:

```bash
cd betting-system/scripts
./test_events_market_filter.sh https://betradar-uof-service-production.up.railway.app
```

## 数据准备

确保数据库中有测试数据:

```bash
# 1. 创建基础数据（sports, categories, tournaments, events）
psql $DATABASE_URL -f betting-system/scripts/seed_test_data.sql

# 2. 创建 markets 和 outcomes 数据
psql $DATABASE_URL -f betting-system/scripts/seed_markets_outcomes.sql
```

## 版本历史

| 版本 | 日期 | 变更说明 |
|:---|:---|:---|
| 1.1.0 | 2024-11-25 | 添加 `market_types` 参数支持，支持按需加载 markets |
| 1.0.0 | 2024-11-24 | 初始版本 |
