# Category 和 Tournament API 文档

## 概述

本文档描述了用于获取体育分类和联赛数据的 RESTful API 接口。这些接口支持分页、排序和比赛数量统计功能。

## 基础 URL

```
https://betradar-uof-service-production.up.railway.app
```

或本地开发环境:
```
http://localhost:8080
```

## 接口列表

### 1. 获取分类数据

获取体育分类列表，支持按体育类型过滤、分页和排序。

**请求**

```
GET /api/categories
```

**查询参数**

| 参数 | 类型 | 必填 | 说明 | 示例 |
|:---|:---|:---:|:---|:---|
| `sport_ids` | string | 否 | 体育类型ID列表，逗号分隔 | `sr:sport:1,sr:sport:2` |
| `page` | integer | 否 | 页码，默认为 1 | `1` |
| `page_size` | integer | 否 | 每页大小，默认为 20，最大 100 | `20` |
| `sort` | string | 否 | 排序方式，可选值见下表 | `name_asc` |

**排序方式**

| 值 | 说明 |
|:---|:---|
| `name_asc` | 按名称升序（默认） |
| `name_desc` | 按名称降序 |
| `match_count_desc` | 按比赛数量降序 |

**响应示例**

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
    },
    {
      "id": 2,
      "external_id": "sr:category:7",
      "sport_id": 1,
      "sport_name": "Football",
      "name": "Spain",
      "match_count": 120
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

**字段说明**

| 字段 | 类型 | 说明 |
|:---|:---|:---|
| `id` | integer | 分类内部ID |
| `external_id` | string | SportRader 分类ID |
| `sport_id` | integer | 体育类型内部ID |
| `sport_name` | string | 体育类型名称 |
| `name` | string | 分类名称 |
| `match_count` | integer | 该分类下的比赛数量 |

**请求示例**

```bash
# 获取所有分类
curl "https://betradar-uof-service-production.up.railway.app/api/categories"

# 获取足球分类
curl "https://betradar-uof-service-production.up.railway.app/api/categories?sport_ids=sr:sport:1"

# 按比赛数量降序排列
curl "https://betradar-uof-service-production.up.railway.app/api/categories?sort=match_count_desc"

# 分页获取
curl "https://betradar-uof-service-production.up.railway.app/api/categories?page=1&page_size=10"
```

---

### 2. 获取联赛数据

获取指定分类下的联赛列表，支持分页和排序。

**请求**

```
GET /api/tournaments
```

**查询参数**

| 参数 | 类型 | 必填 | 说明 | 示例 |
|:---|:---|:---:|:---|:---|
| `category_id` | integer | **是** | 分类ID | `1` |
| `page` | integer | 否 | 页码，默认为 1 | `1` |
| `page_size` | integer | 否 | 每页大小，默认为 20，最大 100 | `20` |
| `sort` | string | 否 | 排序方式，可选值见下表 | `name_asc` |

**排序方式**

| 值 | 说明 |
|:---|:---|
| `name_asc` | 按名称升序（默认） |
| `name_desc` | 按名称降序 |
| `match_count_desc` | 按比赛数量降序 |

**响应示例**

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
    },
    {
      "id": 2,
      "external_id": "sr:tournament:18",
      "category_id": 1,
      "category_name": "England",
      "name": "Championship",
      "scheduled": "2024-08-01T00:00:00Z",
      "scheduled_end": "2025-05-31T23:59:59Z",
      "match_count": 552
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

**字段说明**

| 字段 | 类型 | 说明 |
|:---|:---|:---|
| `id` | integer | 联赛内部ID |
| `external_id` | string | SportRader 联赛ID |
| `category_id` | integer | 分类ID |
| `category_name` | string | 分类名称 |
| `name` | string | 联赛名称 |
| `scheduled` | string | 联赛开始时间 (ISO 8601 格式) |
| `scheduled_end` | string | 联赛结束时间 (ISO 8601 格式) |
| `match_count` | integer | 该联赛下的比赛数量 |

**请求示例**

```bash
# 获取指定分类的联赛
curl "https://betradar-uof-service-production.up.railway.app/api/tournaments?category_id=1"

# 按比赛数量降序排列
curl "https://betradar-uof-service-production.up.railway.app/api/tournaments?category_id=1&sort=match_count_desc"

# 分页获取
curl "https://betradar-uof-service-production.up.railway.app/api/tournaments?category_id=1&page=1&page_size=5"
```

---

## 错误响应

### 400 Bad Request

请求参数无效或缺少必填参数。

```json
{
  "error": "category_id is required"
}
```

### 404 Not Found

请求的资源不存在。

```json
{
  "error": "Category not found"
}
```

### 500 Internal Server Error

服务器内部错误。

```json
{
  "error": "Failed to fetch categories"
}
```

---

## 数据模型

### Sport (体育类型)

```
Sport
├── id: integer
├── external_id: string (sr:sport:{id})
└── name: string
```

### Category (分类)

```
Category
├── id: integer
├── external_id: string (sr:category:{id})
├── sport_id: integer (外键 -> Sport)
└── name: string
```

### Tournament (联赛)

```
Tournament
├── id: integer
├── external_id: string (sr:tournament:{id} 或 sr:stage:{id})
├── category_id: integer (外键 -> Category)
├── name: string
├── scheduled: timestamp
└── scheduled_end: timestamp
```

### Event (比赛)

```
Event
├── id: integer
├── external_id: string (sr:match:{id})
├── sport_id: string
├── tournament_id: integer (外键 -> Tournament)
├── home_team: string
├── away_team: string
├── start_time: timestamp
└── status: string
```

---

## 使用场景

### 场景 1: 获取所有足球联赛

1. 获取足球分类
```bash
curl "https://betradar-uof-service-production.up.railway.app/api/categories?sport_ids=sr:sport:1"
```

2. 从响应中提取 `category_id`

3. 获取该分类下的联赛
```bash
curl "https://betradar-uof-service-production.up.railway.app/api/tournaments?category_id={category_id}"
```

### 场景 2: 获取比赛最多的分类

```bash
curl "https://betradar-uof-service-production.up.railway.app/api/categories?sort=match_count_desc&page_size=10"
```

### 场景 3: 获取英格兰足球联赛

1. 获取英格兰分类
```bash
curl "https://betradar-uof-service-production.up.railway.app/api/categories?sport_ids=sr:sport:1" | jq '.data[] | select(.name == "England")'
```

2. 获取该分类下的联赛
```bash
curl "https://betradar-uof-service-production.up.railway.app/api/tournaments?category_id={category_id}&sort=match_count_desc"
```

---

## 注意事项

1. **分页限制**: `page_size` 最大值为 100，超过将自动设置为 100
2. **比赛数量统计**: `match_count` 是实时统计的，反映当前数据库中的比赛数量
3. **外部 ID 格式**: 
   - Sport: `sr:sport:{id}`
   - Category: `sr:category:{id}`
   - Tournament: `sr:tournament:{id}` 或 `sr:stage:{id}`
   - Event: `sr:match:{id}`
4. **时区**: 所有时间字段使用 UTC 时区

---

## 版本历史

| 版本 | 日期 | 变更说明 |
|:---|:---|:---|
| 1.0.0 | 2024-11-25 | 初始版本，添加 Category 和 Tournament 接口 |
