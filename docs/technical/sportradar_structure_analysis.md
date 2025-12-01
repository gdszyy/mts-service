# SportRader UOF 数据结构分析

## 层级结构

根据 SportRader 文档,数据层级结构如下:

```
Sport (运动类型)
  └─ Category (分类/地区)
      └─ Tournament/Simple_Tournament (联赛/赛事)
          └─ Season (赛季)
              └─ Match/Event (比赛)
```

## 关键实体

### Sport (运动类型)
- ID 格式: `sr:sport:{id}`
- 示例: `sr:sport:1` (足球), `sr:sport:9` (高尔夫)

### Category (分类)
- ID 格式: `sr:category:{id}`
- 示例: `sr:category:28` (Men), `sr:category:96` (International)
- 通常代表地区或性别分类

### Tournament (联赛)
- ID 格式: `sr:stage:{id}` 或 `sr:simple_tournament:{id}`
- 示例: `sr:stage:607727` (PGA Tour 2021)
- 包含多个赛季或比赛

### Match/Event (比赛)
- ID 格式: `sr:match:{id}` 或 `sr:sport_event:{id}`
- 示例: `sr:match:12345678`

## 数据库设计建议

根据现有的 `events` 表,我们需要添加以下表来支持完整的层级结构:

### 1. sports 表
```sql
CREATE TABLE sports (
  id BIGSERIAL PRIMARY KEY,
  external_id VARCHAR(100) UNIQUE NOT NULL,  -- sr:sport:1
  name VARCHAR(255) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### 2. categories 表
```sql
CREATE TABLE categories (
  id BIGSERIAL PRIMARY KEY,
  external_id VARCHAR(100) UNIQUE NOT NULL,  -- sr:category:28
  sport_id BIGINT NOT NULL REFERENCES sports(id),
  name VARCHAR(255) NOT NULL,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### 3. tournaments 表
```sql
CREATE TABLE tournaments (
  id BIGSERIAL PRIMARY KEY,
  external_id VARCHAR(100) UNIQUE NOT NULL,  -- sr:stage:607727
  category_id BIGINT NOT NULL REFERENCES categories(id),
  name VARCHAR(255) NOT NULL,
  scheduled TIMESTAMP,
  scheduled_end TIMESTAMP,
  created_at TIMESTAMP NOT NULL DEFAULT NOW(),
  updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);
```

### 4. 修改 events 表
需要添加 `tournament_id` 字段:
```sql
ALTER TABLE events ADD COLUMN tournament_id BIGINT REFERENCES tournaments(id);
```

## API 接口设计

### 1. 获取分类数据
**请求**: `GET /api/categories`

**参数**:
- `sport_ids`: 体育类型列表 (逗号分隔)
- `page`: 页码 (默认 1)
- `page_size`: 每页大小 (默认 20)
- `sort`: 排序方式 (name_asc, name_desc, match_count_desc)

**响应**:
```json
{
  "data": [
    {
      "id": 1,
      "external_id": "sr:category:28",
      "sport_id": 9,
      "sport_name": "Golf",
      "name": "Men",
      "match_count": 150
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 100
  }
}
```

### 2. 获取联赛数据
**请求**: `GET /api/tournaments`

**参数**:
- `category_id`: 分类 ID (必填)
- `page`: 页码 (默认 1)
- `page_size`: 每页大小 (默认 20)
- `sort`: 排序方式 (name_asc, name_desc, match_count_desc)

**响应**:
```json
{
  "data": [
    {
      "id": 1,
      "external_id": "sr:stage:607727",
      "category_id": 28,
      "category_name": "Men",
      "name": "PGA Tour 2021",
      "scheduled": "2020-09-10T07:00:00Z",
      "scheduled_end": "2021-09-06T04:00:00Z",
      "match_count": 45
    }
  ],
  "pagination": {
    "page": 1,
    "page_size": 20,
    "total": 50
  }
}
```
