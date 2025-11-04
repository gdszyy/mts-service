# 投注系统数据库架构设计

## 核心表结构

### 1. users（用户表）

| 字段名 | 类型 | 约束 | 说明 |
|:---|:---|:---|:---|
| id | BIGSERIAL | PRIMARY KEY | 用户ID |
| username | VARCHAR(100) | UNIQUE NOT NULL | 用户名 |
| email | VARCHAR(255) | UNIQUE NOT NULL | 邮箱 |
| balance | DECIMAL(15,2) | NOT NULL DEFAULT 0 | 账户余额（单位：元） |
| currency | VARCHAR(3) | NOT NULL DEFAULT 'CNY' | 货币类型 |
| status | VARCHAR(20) | NOT NULL DEFAULT 'active' | 状态（active/suspended/closed） |
| created_at | TIMESTAMP | NOT NULL DEFAULT NOW() | 创建时间 |
| updated_at | TIMESTAMP | NOT NULL DEFAULT NOW() | 更新时间 |

### 2. events（赛事表）

| 字段名 | 类型 | 约束 | 说明 |
|:---|:---|:---|:---|
| id | BIGSERIAL | PRIMARY KEY | 赛事ID |
| external_id | VARCHAR(100) | UNIQUE | 外部赛事ID（如UOF的sr:match:xxx） |
| sport_id | VARCHAR(50) | NOT NULL | 运动类型ID |
| home_team | VARCHAR(255) | NOT NULL | 主队名称 |
| away_team | VARCHAR(255) | NOT NULL | 客队名称 |
| start_time | TIMESTAMP | NOT NULL | 开赛时间 |
| status | VARCHAR(20) | NOT NULL DEFAULT 'scheduled' | 状态（scheduled/live/finished/cancelled） |
| home_score | INTEGER | | 主队得分 |
| away_score | INTEGER | | 客队得分 |
| created_at | TIMESTAMP | NOT NULL DEFAULT NOW() | 创建时间 |
| updated_at | TIMESTAMP | NOT NULL DEFAULT NOW() | 更新时间 |

### 3. markets（盘口表）

| 字段名 | 类型 | 约束 | 说明 |
|:---|:---|:---|:---|
| id | BIGSERIAL | PRIMARY KEY | 盘口ID |
| event_id | BIGINT | FOREIGN KEY REFERENCES events(id) | 赛事ID |
| market_type | VARCHAR(50) | NOT NULL | 盘口类型（1x2/handicap/totals等） |
| specifier | VARCHAR(100) | | 盘口参数（如handicap=-1.5） |
| status | VARCHAR(20) | NOT NULL DEFAULT 'active' | 状态（active/suspended/settled/cancelled） |
| created_at | TIMESTAMP | NOT NULL DEFAULT NOW() | 创建时间 |
| updated_at | TIMESTAMP | NOT NULL DEFAULT NOW() | 更新时间 |

### 4. outcomes（结果选项表）

| 字段名 | 类型 | 约束 | 说明 |
|:---|:---|:---|:---|
| id | BIGSERIAL | PRIMARY KEY | 选项ID |
| market_id | BIGINT | FOREIGN KEY REFERENCES markets(id) | 盘口ID |
| outcome_id | VARCHAR(50) | NOT NULL | 结果ID（如1/x/2，over/under等） |
| odds | DECIMAL(10,4) | NOT NULL | 赔率 |
| status | VARCHAR(20) | NOT NULL DEFAULT 'active' | 状态（active/suspended/won/lost/void） |
| created_at | TIMESTAMP | NOT NULL DEFAULT NOW() | 创建时间 |
| updated_at | TIMESTAMP | NOT NULL DEFAULT NOW() | 更新时间 |

### 5. bets（投注表）

| 字段名 | 类型 | 约束 | 说明 |
|:---|:---|:---|:---|
| id | BIGSERIAL | PRIMARY KEY | 投注ID |
| user_id | BIGINT | FOREIGN KEY REFERENCES users(id) | 用户ID |
| bet_type | VARCHAR(50) | NOT NULL | 投注类型（single/accumulator/system/patent/trixie等） |
| total_stake | DECIMAL(15,2) | NOT NULL | 总投注金额 |
| potential_return | DECIMAL(15,2) | NOT NULL | 潜在回报 |
| actual_return | DECIMAL(15,2) | | 实际回报 |
| status | VARCHAR(20) | NOT NULL DEFAULT 'pending' | 状态（pending/won/lost/void/partially_won） |
| placed_at | TIMESTAMP | NOT NULL DEFAULT NOW() | 下注时间 |
| settled_at | TIMESTAMP | | 结算时间 |
| created_at | TIMESTAMP | NOT NULL DEFAULT NOW() | 创建时间 |
| updated_at | TIMESTAMP | NOT NULL DEFAULT NOW() | 更新时间 |

### 6. bet_selections（投注选项表）

| 字段名 | 类型 | 约束 | 说明 |
|:---|:---|:---|:---|
| id | BIGSERIAL | PRIMARY KEY | 选项ID |
| bet_id | BIGINT | FOREIGN KEY REFERENCES bets(id) | 投注ID |
| outcome_id | BIGINT | FOREIGN KEY REFERENCES outcomes(id) | 结果选项ID |
| odds | DECIMAL(10,4) | NOT NULL | 下注时的赔率 |
| is_banker | BOOLEAN | NOT NULL DEFAULT FALSE | 是否为Banker |
| status | VARCHAR(20) | NOT NULL DEFAULT 'pending' | 状态（pending/won/lost/void） |
| created_at | TIMESTAMP | NOT NULL DEFAULT NOW() | 创建时间 |

### 7. bet_legs（投注组合腿表）

| 字段名 | 类型 | 约束 | 说明 |
|:---|:---|:---|:---|
| id | BIGSERIAL | PRIMARY KEY | 组合腿ID |
| bet_id | BIGINT | FOREIGN KEY REFERENCES bets(id) | 投注ID |
| leg_type | VARCHAR(50) | NOT NULL | 组合类型（single/double/treble/4-fold等） |
| stake | DECIMAL(15,2) | NOT NULL | 单腿投注金额 |
| odds | DECIMAL(10,4) | NOT NULL | 组合赔率 |
| potential_return | DECIMAL(15,2) | NOT NULL | 潜在回报 |
| actual_return | DECIMAL(15,2) | | 实际回报 |
| status | VARCHAR(20) | NOT NULL DEFAULT 'pending' | 状态（pending/won/lost/void） |
| created_at | TIMESTAMP | NOT NULL DEFAULT NOW() | 创建时间 |

### 8. bet_leg_selections（组合腿选项关联表）

| 字段名 | 类型 | 约束 | 说明 |
|:---|:---|:---|:---|
| id | BIGSERIAL | PRIMARY KEY | 关联ID |
| leg_id | BIGINT | FOREIGN KEY REFERENCES bet_legs(id) | 组合腿ID |
| selection_id | BIGINT | FOREIGN KEY REFERENCES bet_selections(id) | 选项ID |
| created_at | TIMESTAMP | NOT NULL DEFAULT NOW() | 创建时间 |

## 投注类型说明

### 系统投注类型定义

| 投注类型 | 选项数量 | 组合说明 | Banker支持 |
|:---|:---|:---|:---|
| **Trixie** | 3 | 3 doubles + 1 treble = 4 bets | ✅ |
| **Patent** | 3 | 3 singles + 3 doubles + 1 treble = 7 bets | ✅ |
| **Yankee** | 4 | 6 doubles + 4 trebles + 1 4-fold = 11 bets | ✅ |
| **Lucky 15** | 4 | 4 singles + 6 doubles + 4 trebles + 1 4-fold = 15 bets | ✅ |
| **Super Yankee (Canadian)** | 5 | 10 doubles + 10 trebles + 5 4-folds + 1 5-fold = 26 bets | ✅ |
| **Lucky 31** | 5 | 5 singles + 10 doubles + 10 trebles + 5 4-folds + 1 5-fold = 31 bets | ✅ |
| **Heinz** | 6 | 15 doubles + 20 trebles + 15 4-folds + 6 5-folds + 1 6-fold = 57 bets | ✅ |
| **Lucky 63** | 6 | 6 singles + 15 doubles + 20 trebles + 15 4-folds + 6 5-folds + 1 6-fold = 63 bets | ✅ |
| **Super Heinz** | 7 | 21 doubles + 35 trebles + 35 4-folds + 21 5-folds + 7 6-folds + 1 7-fold = 120 bets | ✅ |
| **Goliath** | 8 | 28 doubles + 56 trebles + 70 4-folds + 56 5-folds + 28 6-folds + 8 7-folds + 1 8-fold = 247 bets | ✅ |

### Banker 逻辑

当某些选项被标记为 **Banker** 时：
- Banker 选项必须出现在所有组合中
- 非 Banker 选项进行正常组合
- 例如：4个选项，1个Banker + 3个普通选项的Trixie = 3 doubles + 1 treble（所有组合都包含Banker）

## 索引设计

```sql
-- 用户表索引
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_email ON users(email);

-- 赛事表索引
CREATE INDEX idx_events_external_id ON events(external_id);
CREATE INDEX idx_events_start_time ON events(start_time);
CREATE INDEX idx_events_status ON events(status);

-- 盘口表索引
CREATE INDEX idx_markets_event_id ON markets(event_id);
CREATE INDEX idx_markets_status ON markets(status);

-- 结果选项表索引
CREATE INDEX idx_outcomes_market_id ON outcomes(market_id);
CREATE INDEX idx_outcomes_status ON outcomes(status);

-- 投注表索引
CREATE INDEX idx_bets_user_id ON bets(user_id);
CREATE INDEX idx_bets_status ON bets(status);
CREATE INDEX idx_bets_placed_at ON bets(placed_at);

-- 投注选项表索引
CREATE INDEX idx_bet_selections_bet_id ON bet_selections(bet_id);
CREATE INDEX idx_bet_selections_outcome_id ON bet_selections(outcome_id);

-- 投注组合腿表索引
CREATE INDEX idx_bet_legs_bet_id ON bet_legs(bet_id);

-- 组合腿选项关联表索引
CREATE INDEX idx_bet_leg_selections_leg_id ON bet_leg_selections(leg_id);
CREATE INDEX idx_bet_leg_selections_selection_id ON bet_leg_selections(selection_id);
```

## 数据完整性约束

1. **用户余额约束**: `balance >= 0`（不允许负余额）
2. **赔率约束**: `odds > 1.0`（赔率必须大于1）
3. **投注金额约束**: `total_stake > 0`（投注金额必须大于0）
4. **状态转换约束**: 通过应用层逻辑确保状态转换的合法性

