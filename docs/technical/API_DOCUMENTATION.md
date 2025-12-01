# MTS Service API Documentation

## Overview

MTS Service 提供了完整的 RESTful API 接口，允许外部系统通过 HTTP 请求提交各种类型的投注注单到 Sportradar MTS (Managed Trading Services)。

**Base URL**: `http://your-server:8080`

**版本**: 2.0.0

## Authentication

目前 API 不需要身份验证。所有请求直接发送到相应的端点即可。

## Common Response Format

所有 API 端点返回统一的响应格式：

```json
{
  "success": true,
  "data": { ... },
  "error": null
}
```

错误响应：

```json
{
  "success": false,
  "data": null,
  "error": {
    "code": 400,
    "message": "Validation failed",
    "details": "ticketId is required"
  }
}
```

## Endpoints

### 1. Health Check

检查服务健康状态。

**Endpoint**: `GET /health`

**Response**:
```json
{
  "status": "healthy",
  "timestamp": 1732612345,
  "service": "mts-service"
}
```

---

### 2. Place Single Bet

提交单注投注。

**Endpoint**: `POST /api/bets/single`

**Request Body**:
```json
{
  "ticketId": "ticket-001",
  "selection": {
    "productId": "3",
    "eventId": "sr:match:12345",
    "marketId": "1",
    "outcomeId": "1",
    "odds": 2.50,
    "specifiers": ""
  },
  "stake": {
    "type": "cash",
    "currency": "EUR",
    "amount": 10.00,
    "mode": "total"
  },
  "context": {
    "channelType": "internet",
    "language": "EN",
    "ip": "192.168.1.1"
  }
}
```

**Field Descriptions**:

| Field | Type | Required | Description |
|:---|:---|:---:|:---|
| `ticketId` | string | ✅ | 唯一的注单 ID |
| `selection.productId` | string | ✅ | 产品 ID（通常为 "3"） |
| `selection.eventId` | string | ✅ | 赛事 ID（如 "sr:match:12345"） |
| `selection.marketId` | string | ✅ | 市场 ID |
| `selection.outcomeId` | string | ✅ | 结果 ID |
| `selection.odds` | number | ✅ | 赔率（十进制格式） |
| `selection.specifiers` | string | ❌ | 可选的说明符（如 "hcp=1:0"） |
| `stake.type` | string | ✅ | "cash" 或 "free" |
| `stake.currency` | string | ✅ | 货币代码（如 "EUR", "USD"） |
| `stake.amount` | number | ✅ | 投注金额 |
| `stake.mode` | string | ✅ | "total" 或 "unit" |
| `context.channelType` | string | ❌ | 渠道类型（默认 "internet"） |
| `context.language` | string | ❌ | 语言代码（默认 "EN"） |
| `context.ip` | string | ❌ | 客户 IP 地址 |

**Response**:
```json
{
  "success": true,
  "data": {
    "content": {
      "type": "ticket-reply",
      "ticketId": "ticket-001",
      "status": "accepted",
      "signature": "...",
      "betDetails": [...]
    },
    "correlationId": "corr-1732612345000",
    "timestampUtc": 1732612345000,
    "operation": "ticket-placement",
    "version": "3.0"
  }
}
```

---

### 3. Place Accumulator Bet

提交串关投注（所有选项必须全部正确才能赢）。

**Endpoint**: `POST /api/bets/accumulator`

**Request Body**:
```json
{
  "ticketId": "ticket-acc-001",
  "selections": [
    {
      "productId": "3",
      "eventId": "sr:match:12345",
      "marketId": "1",
      "outcomeId": "1",
      "odds": 2.50
    },
    {
      "productId": "3",
      "eventId": "sr:match:12346",
      "marketId": "1",
      "outcomeId": "2",
      "odds": 1.80
    },
    {
      "productId": "3",
      "eventId": "sr:match:12347",
      "marketId": "1",
      "outcomeId": "1",
      "odds": 3.00
    }
  ],
  "stake": {
    "type": "cash",
    "currency": "EUR",
    "amount": 10.00,
    "mode": "total"
  }
}
```

**Requirements**:
- 至少需要 2 个选项
- `stake.mode` 应为 "total"

---

### 4. Place System Bet

提交系统串投注（自定义组合大小）。

**Endpoint**: `POST /api/bets/system`

**Request Body**:
```json
{
  "ticketId": "ticket-sys-001",
  "size": [2],
  "selections": [
    {
      "productId": "3",
      "eventId": "sr:match:12345",
      "marketId": "1",
      "outcomeId": "1",
      "odds": 2.50
    },
    {
      "productId": "3",
      "eventId": "sr:match:12346",
      "marketId": "1",
      "outcomeId": "2",
      "odds": 1.80
    },
    {
      "productId": "3",
      "eventId": "sr:match:12347",
      "marketId": "1",
      "outcomeId": "1",
      "odds": 3.00
    },
    {
      "productId": "3",
      "eventId": "sr:match:12348",
      "marketId": "1",
      "outcomeId": "2",
      "odds": 2.20
    }
  ],
  "stake": {
    "type": "cash",
    "currency": "EUR",
    "amount": 1.00,
    "mode": "unit"
  }
}
```

**Field Descriptions**:

| Field | Description | Example |
|:---|:---|:---|
| `size` | 组合大小数组 | `[2]` = 所有双式<br>`[2, 3]` = 所有双式 + 所有三串一<br>`[2, 3, 4]` = 双式 + 三串一 + 四串一 |

**Requirements**:
- 至少需要 2 个选项
- `stake.mode` 必须为 "unit"（单位投注）
- `size` 中的每个值必须在 1 到选项数量之间

**Examples**:
- **2/4 系统** (6 注): `size: [2]`, 4 个选项
- **2/3 系统** (3 注): `size: [2]`, 3 个选项
- **2/3 + 3/3 系统** (4 注): `size: [2, 3]`, 3 个选项

---

### 5. Place Banker System Bet

提交 Banker 系统串投注（某些选项必须在每个组合中）。

**Endpoint**: `POST /api/bets/banker-system`

**Request Body**:
```json
{
  "ticketId": "ticket-banker-001",
  "bankers": [
    {
      "productId": "3",
      "eventId": "sr:match:12345",
      "marketId": "1",
      "outcomeId": "1",
      "odds": 1.50
    }
  ],
  "size": [2],
  "selections": [
    {
      "productId": "3",
      "eventId": "sr:match:12346",
      "marketId": "1",
      "outcomeId": "2",
      "odds": 2.50
    },
    {
      "productId": "3",
      "eventId": "sr:match:12347",
      "marketId": "1",
      "outcomeId": "1",
      "odds": 3.00
    },
    {
      "productId": "3",
      "eventId": "sr:match:12348",
      "marketId": "1",
      "outcomeId": "2",
      "odds": 2.20
    }
  ],
  "stake": {
    "type": "cash",
    "currency": "EUR",
    "amount": 1.00,
    "mode": "unit"
  }
}
```

**Field Descriptions**:

| Field | Description |
|:---|:---|
| `bankers` | 必须出现在每个组合中的选项 |
| `size` | 非 Banker 选项的组合大小 |
| `selections` | 非 Banker 选项 |

**Example**:
- 1 个 Banker + 2/3 系统 = 3 注（Banker 出现在每注中）

---

### 6. Place Preset System Bet

提交预设系统串投注（Trixie, Yankee, Lucky 15 等）。

**Endpoint**: `POST /api/bets/preset`

**Request Body**:
```json
{
  "ticketId": "ticket-trixie-001",
  "type": "trixie",
  "selections": [
    {
      "productId": "3",
      "eventId": "sr:match:12345",
      "marketId": "1",
      "outcomeId": "1",
      "odds": 2.50
    },
    {
      "productId": "3",
      "eventId": "sr:match:12346",
      "marketId": "1",
      "outcomeId": "2",
      "odds": 1.80
    },
    {
      "productId": "3",
      "eventId": "sr:match:12347",
      "marketId": "1",
      "outcomeId": "1",
      "odds": 3.00
    }
  ],
  "stake": {
    "type": "cash",
    "currency": "EUR",
    "amount": 1.00,
    "mode": "unit"
  }
}
```

**Supported Preset Types**:

| Type | Selections | Total Bets | Description |
|:---|:---:|:---:|:---|
| `trixie` | 3 | 4 | 3 双式 + 1 三串一 |
| `patent` | 3 | 7 | 3 单式 + 3 双式 + 1 三串一 |
| `yankee` | 4 | 11 | 6 双式 + 4 三串一 + 1 四串一 |
| `lucky15` | 4 | 15 | 4 单式 + 6 双式 + 4 三串一 + 1 四串一 |
| `super_yankee` | 5 | 26 | 10 双式 + 10 三串一 + 5 四串一 + 1 五串一 |
| `lucky31` | 5 | 31 | 5 单式 + Super Yankee |
| `heinz` | 6 | 57 | 15 双式 + 20 三串一 + 15 四串一 + 6 五串一 + 1 六串一 |
| `lucky63` | 6 | 63 | 6 单式 + Heinz |
| `super_heinz` | 7 | 120 | 21 双式 + 35 三串一 + 35 四串一 + 21 五串一 + 7 六串一 + 1 七串一 |
| `goliath` | 8 | 247 | 28 双式 + 56 三串一 + 70 四串一 + 56 五串一 + 28 六串一 + 8 七串一 + 1 八串一 |

**Requirements**:
- 选项数量必须与类型要求完全匹配
- `stake.mode` 必须为 "unit"

---

### 7. Place Multi Bet

在一个注单中提交多种类型的投注。

**Endpoint**: `POST /api/bets/multi`

**Request Body**:
```json
{
  "ticketId": "ticket-multi-001",
  "bets": [
    {
      "type": "single",
      "selections": [
        {
          "productId": "3",
          "eventId": "sr:match:12345",
          "marketId": "1",
          "outcomeId": "1",
          "odds": 2.50
        }
      ],
      "stake": {
        "type": "cash",
        "currency": "EUR",
        "amount": 5.00,
        "mode": "total"
      }
    },
    {
      "type": "accumulator",
      "selections": [
        {
          "productId": "3",
          "eventId": "sr:match:12346",
          "marketId": "1",
          "outcomeId": "2",
          "odds": 1.80
        },
        {
          "productId": "3",
          "eventId": "sr:match:12347",
          "marketId": "1",
          "outcomeId": "1",
          "odds": 3.00
        }
      ],
      "stake": {
        "type": "cash",
        "currency": "EUR",
        "amount": 10.00,
        "mode": "total"
      }
    },
    {
      "type": "trixie",
      "selections": [
        {
          "productId": "3",
          "eventId": "sr:match:12348",
          "marketId": "1",
          "outcomeId": "2",
          "odds": 2.20
        },
        {
          "productId": "3",
          "eventId": "sr:match:12349",
          "marketId": "1",
          "outcomeId": "1",
          "odds": 1.95
        },
        {
          "productId": "3",
          "eventId": "sr:match:12350",
          "marketId": "1",
          "outcomeId": "2",
          "odds": 2.75
        }
      ],
      "stake": {
        "type": "cash",
        "currency": "EUR",
        "amount": 1.00,
        "mode": "unit"
      }
    }
  ]
}
```

**Supported Bet Types**:
- `single`
- `accumulator`
- `system`
- `banker_system`
- All preset types (trixie, yankee, etc.)

---

### 8. Request Cashout

请求提前结算（Cashout）。

**Endpoint**: `POST /api/cashout`

**Request Body**:
```json
{
  "cashoutId": "cashout-001",
  "ticketId": "ticket-001",
  "ticketSignature": "signature-from-ticket-response",
  "type": "ticket",
  "code": 100,
  "payout": [
    {
      "type": "cash",
      "currency": "EUR",
      "amount": 15.50,
      "source": "cashout"
    }
  ]
}
```

**Field Descriptions**:

| Field | Type | Required | Description |
|:---|:---|:---:|:---|
| `cashoutId` | string | ✅ | 唯一的 Cashout ID |
| `ticketId` | string | ✅ | 原始注单 ID |
| `ticketSignature` | string | ✅ | 原始注单响应中的签名 |
| `type` | string | ✅ | "ticket", "ticket-partial", "bet", "bet-partial" |
| `code` | number | ✅ | Cashout 原因代码（100, 101 等） |
| `percentage` | number | ❌ | 部分 Cashout 的百分比（0.0-1.0） |
| `betId` | string | ❌ | Bet 级别 Cashout 的 Bet ID |
| `payout` | array | ✅ | 支付信息数组 |

**Cashout Types**:
- `ticket`: 完整注单 Cashout
- `ticket-partial`: 部分注单 Cashout（需要 `percentage`）
- `bet`: 单个 Bet Cashout（需要 `betId`）
- `bet-partial`: 部分 Bet Cashout（需要 `betId` 和 `percentage`）

**Note**: Cashout 功能目前在服务层尚未完全实现，API 端点会返回 501 Not Implemented。

---

## Error Codes

| Code | Message | Description |
|:---|:---|:---|
| 400 | Bad Request | 请求格式错误或参数验证失败 |
| 405 | Method Not Allowed | HTTP 方法不允许 |
| 500 | Internal Server Error | 服务器内部错误 |
| 501 | Not Implemented | 功能尚未实现 |

---

## Examples

### cURL Examples

**Single Bet**:
```bash
curl -X POST http://localhost:8080/api/bets/single \
  -H "Content-Type: application/json" \
  -d '{
    "ticketId": "test-single-001",
    "selection": {
      "productId": "3",
      "eventId": "sr:match:12345",
      "marketId": "1",
      "outcomeId": "1",
      "odds": 2.50
    },
    "stake": {
      "type": "cash",
      "currency": "EUR",
      "amount": 10.00,
      "mode": "total"
    }
  }'
```

**Accumulator Bet**:
```bash
curl -X POST http://localhost:8080/api/bets/accumulator \
  -H "Content-Type: application/json" \
  -d '{
    "ticketId": "test-acc-001",
    "selections": [
      {
        "productId": "3",
        "eventId": "sr:match:12345",
        "marketId": "1",
        "outcomeId": "1",
        "odds": 2.50
      },
      {
        "productId": "3",
        "eventId": "sr:match:12346",
        "marketId": "1",
        "outcomeId": "2",
        "odds": 1.80
      }
    ],
    "stake": {
      "type": "cash",
      "currency": "EUR",
      "amount": 10.00,
      "mode": "total"
    }
  }'
```

**Trixie Bet**:
```bash
curl -X POST http://localhost:8080/api/bets/preset \
  -H "Content-Type: application/json" \
  -d '{
    "ticketId": "test-trixie-001",
    "type": "trixie",
    "selections": [
      {
        "productId": "3",
        "eventId": "sr:match:12345",
        "marketId": "1",
        "outcomeId": "1",
        "odds": 2.50
      },
      {
        "productId": "3",
        "eventId": "sr:match:12346",
        "marketId": "1",
        "outcomeId": "2",
        "odds": 1.80
      },
      {
        "productId": "3",
        "eventId": "sr:match:12347",
        "marketId": "1",
        "outcomeId": "1",
        "odds": 3.00
      }
    ],
    "stake": {
      "type": "cash",
      "currency": "EUR",
      "amount": 1.00,
      "mode": "unit"
    }
  }'
```

---

## Notes

1. **TicketID 唯一性**: 每个 `ticketId` 必须是全局唯一的，建议使用 UUID 或时间戳组合。

2. **Odds 格式**: 所有赔率使用十进制格式（Decimal），例如 2.50 表示 2.5 倍赔率。

3. **Stake Mode**:
   - `total`: 总投注金额（用于单注和串关）
   - `unit`: 单位投注金额（用于系统串，总金额 = 单位金额 × 组合数）

4. **Currency**: 支持的货币代码包括 EUR, USD, GBP, mBTC 等，具体取决于 MTS 配置。

5. **Context**: 如果不提供 `context`，系统会使用默认值（internet 渠道，EN 语言）。

6. **MTS Response**: 所有成功的请求都会返回 MTS 的原始响应，包括 `status`（accepted/rejected）、`signature`、`betDetails` 等。

---

## Support

如有问题或需要帮助，请联系技术支持或查阅 [MTS Transaction 3.0 API 官方文档](https://docs.betradar.com/display/BD/MTS+-+Transaction+3.0)。
