# MTS Service API Examples

本文档提供了所有 API 端点的实际使用示例，包括 cURL 命令、Python 代码和 JavaScript 代码。

## 目录

1. [单注 (Single Bet)](#单注-single-bet)
2. [串关 (Accumulator)](#串关-accumulator)
3. [系统串 (System Bet)](#系统串-system-bet)
4. [Banker 系统串](#banker-系统串)
5. [预设系统串 (Trixie, Yankee 等)](#预设系统串)
6. [混合注单 (Multi Bet)](#混合注单-multi-bet)
7. [Cashout](#cashout)

---

## 单注 (Single Bet)

### cURL

```bash
curl -X POST http://localhost:8080/api/bets/single \
  -H "Content-Type: application/json" \
  -d '{
    "ticketId": "single-001",
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

### Python

```python
import requests
import json

url = "http://localhost:8080/api/bets/single"

payload = {
    "ticketId": "single-001",
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
}

response = requests.post(url, json=payload)
print(json.dumps(response.json(), indent=2))
```

### JavaScript (Node.js)

```javascript
const axios = require('axios');

const url = 'http://localhost:8080/api/bets/single';

const payload = {
  ticketId: 'single-001',
  selection: {
    productId: '3',
    eventId: 'sr:match:12345',
    marketId: '1',
    outcomeId: '1',
    odds: 2.50
  },
  stake: {
    type: 'cash',
    currency: 'EUR',
    amount: 10.00,
    mode: 'total'
  }
};

axios.post(url, payload)
  .then(response => console.log(JSON.stringify(response.data, null, 2)))
  .catch(error => console.error('Error:', error.message));
```

---

## 串关 (Accumulator)

### cURL

```bash
curl -X POST http://localhost:8080/api/bets/accumulator \
  -H "Content-Type: application/json" \
  -d '{
    "ticketId": "acc-001",
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
  }'
```

### Python

```python
import requests

url = "http://localhost:8080/api/bets/accumulator"

payload = {
    "ticketId": "acc-001",
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

response = requests.post(url, json=payload)
print(response.json())
```

---

## 系统串 (System Bet)

### 2/4 系统串 (6 注)

```bash
curl -X POST http://localhost:8080/api/bets/system \
  -H "Content-Type: application/json" \
  -d '{
    "ticketId": "sys-2-4-001",
    "size": [2],
    "selections": [
      {"productId": "3", "eventId": "sr:match:12345", "marketId": "1", "outcomeId": "1", "odds": 2.50},
      {"productId": "3", "eventId": "sr:match:12346", "marketId": "1", "outcomeId": "2", "odds": 1.80},
      {"productId": "3", "eventId": "sr:match:12347", "marketId": "1", "outcomeId": "1", "odds": 3.00},
      {"productId": "3", "eventId": "sr:match:12348", "marketId": "1", "outcomeId": "2", "odds": 2.20}
    ],
    "stake": {
      "type": "cash",
      "currency": "EUR",
      "amount": 1.00,
      "mode": "unit"
    }
  }'
```

**说明**: 
- `size: [2]` = 所有双式组合（C(4,2) = 6 注）
- 总投注金额 = 1.00 EUR × 6 = 6.00 EUR

### 2/3 + 3/3 系统串 (4 注)

```bash
curl -X POST http://localhost:8080/api/bets/system \
  -H "Content-Type: application/json" \
  -d '{
    "ticketId": "sys-2-3-3-3-001",
    "size": [2, 3],
    "selections": [
      {"productId": "3", "eventId": "sr:match:12345", "marketId": "1", "outcomeId": "1", "odds": 2.50},
      {"productId": "3", "eventId": "sr:match:12346", "marketId": "1", "outcomeId": "2", "odds": 1.80},
      {"productId": "3", "eventId": "sr:match:12347", "marketId": "1", "outcomeId": "1", "odds": 3.00}
    ],
    "stake": {
      "type": "cash",
      "currency": "EUR",
      "amount": 1.00,
      "mode": "unit"
    }
  }'
```

**说明**:
- `size: [2, 3]` = 所有双式 + 所有三串一
- C(3,2) + C(3,3) = 3 + 1 = 4 注
- 总投注金额 = 1.00 EUR × 4 = 4.00 EUR

---

## Banker 系统串

### 1 Banker + 2/3 系统串

```bash
curl -X POST http://localhost:8080/api/bets/banker-system \
  -H "Content-Type: application/json" \
  -d '{
    "ticketId": "banker-001",
    "bankers": [
      {"productId": "3", "eventId": "sr:match:12345", "marketId": "1", "outcomeId": "1", "odds": 1.50}
    ],
    "size": [2],
    "selections": [
      {"productId": "3", "eventId": "sr:match:12346", "marketId": "1", "outcomeId": "2", "odds": 2.50},
      {"productId": "3", "eventId": "sr:match:12347", "marketId": "1", "outcomeId": "1", "odds": 3.00},
      {"productId": "3", "eventId": "sr:match:12348", "marketId": "1", "outcomeId": "2", "odds": 2.20}
    ],
    "stake": {
      "type": "cash",
      "currency": "EUR",
      "amount": 1.00,
      "mode": "unit"
    }
  }'
```

**说明**:
- Banker 选项会出现在每个组合中
- 从 3 个非 Banker 选项中选 2 个 = C(3,2) = 3 注
- 每注都包含 Banker 选项
- 总投注金额 = 1.00 EUR × 3 = 3.00 EUR

---

## 预设系统串

### Trixie (3 选项, 4 注)

```bash
curl -X POST http://localhost:8080/api/bets/preset \
  -H "Content-Type: application/json" \
  -d '{
    "ticketId": "trixie-001",
    "type": "trixie",
    "selections": [
      {"productId": "3", "eventId": "sr:match:12345", "marketId": "1", "outcomeId": "1", "odds": 2.50},
      {"productId": "3", "eventId": "sr:match:12346", "marketId": "1", "outcomeId": "2", "odds": 1.80},
      {"productId": "3", "eventId": "sr:match:12347", "marketId": "1", "outcomeId": "1", "odds": 3.00}
    ],
    "stake": {
      "type": "cash",
      "currency": "EUR",
      "amount": 1.00,
      "mode": "unit"
    }
  }'
```

**组成**: 3 双式 + 1 三串一 = 4 注

### Yankee (4 选项, 11 注)

```bash
curl -X POST http://localhost:8080/api/bets/preset \
  -H "Content-Type: application/json" \
  -d '{
    "ticketId": "yankee-001",
    "type": "yankee",
    "selections": [
      {"productId": "3", "eventId": "sr:match:12345", "marketId": "1", "outcomeId": "1", "odds": 2.50},
      {"productId": "3", "eventId": "sr:match:12346", "marketId": "1", "outcomeId": "2", "odds": 1.80},
      {"productId": "3", "eventId": "sr:match:12347", "marketId": "1", "outcomeId": "1", "odds": 3.00},
      {"productId": "3", "eventId": "sr:match:12348", "marketId": "1", "outcomeId": "2", "odds": 2.20}
    ],
    "stake": {
      "type": "cash",
      "currency": "EUR",
      "amount": 1.00,
      "mode": "unit"
    }
  }'
```

**组成**: 6 双式 + 4 三串一 + 1 四串一 = 11 注

### Lucky 15 (4 选项, 15 注)

```bash
curl -X POST http://localhost:8080/api/bets/preset \
  -H "Content-Type: application/json" \
  -d '{
    "ticketId": "lucky15-001",
    "type": "lucky15",
    "selections": [
      {"productId": "3", "eventId": "sr:match:12345", "marketId": "1", "outcomeId": "1", "odds": 2.50},
      {"productId": "3", "eventId": "sr:match:12346", "marketId": "1", "outcomeId": "2", "odds": 1.80},
      {"productId": "3", "eventId": "sr:match:12347", "marketId": "1", "outcomeId": "1", "odds": 3.00},
      {"productId": "3", "eventId": "sr:match:12348", "marketId": "1", "outcomeId": "2", "odds": 2.20}
    ],
    "stake": {
      "type": "cash",
      "currency": "EUR",
      "amount": 1.00,
      "mode": "unit"
    }
  }'
```

**组成**: 4 单式 + 6 双式 + 4 三串一 + 1 四串一 = 15 注

### 所有支持的预设类型

| Type | Selections | Bets | Command |
|:---|:---:|:---:|:---|
| `trixie` | 3 | 4 | 见上文 |
| `patent` | 3 | 7 | `"type": "patent"` |
| `yankee` | 4 | 11 | 见上文 |
| `lucky15` | 4 | 15 | 见上文 |
| `super_yankee` | 5 | 26 | `"type": "super_yankee"` |
| `lucky31` | 5 | 31 | `"type": "lucky31"` |
| `heinz` | 6 | 57 | `"type": "heinz"` |
| `lucky63` | 6 | 63 | `"type": "lucky63"` |
| `super_heinz` | 7 | 120 | `"type": "super_heinz"` |
| `goliath` | 8 | 247 | `"type": "goliath"` |

---

## 混合注单 (Multi Bet)

在一个注单中包含多种类型的投注：

```bash
curl -X POST http://localhost:8080/api/bets/multi \
  -H "Content-Type: application/json" \
  -d '{
    "ticketId": "multi-001",
    "bets": [
      {
        "type": "single",
        "selections": [
          {"productId": "3", "eventId": "sr:match:12345", "marketId": "1", "outcomeId": "1", "odds": 2.50}
        ],
        "stake": {"type": "cash", "currency": "EUR", "amount": 5.00, "mode": "total"}
      },
      {
        "type": "accumulator",
        "selections": [
          {"productId": "3", "eventId": "sr:match:12346", "marketId": "1", "outcomeId": "2", "odds": 1.80},
          {"productId": "3", "eventId": "sr:match:12347", "marketId": "1", "outcomeId": "1", "odds": 3.00}
        ],
        "stake": {"type": "cash", "currency": "EUR", "amount": 10.00, "mode": "total"}
      },
      {
        "type": "trixie",
        "selections": [
          {"productId": "3", "eventId": "sr:match:12348", "marketId": "1", "outcomeId": "2", "odds": 2.20},
          {"productId": "3", "eventId": "sr:match:12349", "marketId": "1", "outcomeId": "1", "odds": 1.95},
          {"productId": "3", "eventId": "sr:match:12350", "marketId": "1", "outcomeId": "2", "odds": 2.75}
        ],
        "stake": {"type": "cash", "currency": "EUR", "amount": 1.00, "mode": "unit"}
      }
    ]
  }'
```

**说明**:
- 1 单注 (5 EUR)
- 1 串关 (10 EUR)
- 1 Trixie (1 EUR × 4 = 4 EUR)
- 总投注金额 = 19 EUR

---

## Cashout

### 完整注单 Cashout

```bash
curl -X POST http://localhost:8080/api/cashout \
  -H "Content-Type: application/json" \
  -d '{
    "cashoutId": "cashout-001",
    "ticketId": "single-001",
    "ticketSignature": "signature-from-original-response",
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
  }'
```

### 部分注单 Cashout (50%)

```bash
curl -X POST http://localhost:8080/api/cashout \
  -H "Content-Type: application/json" \
  -d '{
    "cashoutId": "cashout-002",
    "ticketId": "acc-001",
    "ticketSignature": "signature-from-original-response",
    "type": "ticket-partial",
    "code": 101,
    "percentage": 0.50,
    "payout": [
      {
        "type": "cash",
        "currency": "EUR",
        "amount": 7.75,
        "source": "cashout"
      }
    ]
  }'
```

### Bet 级别 Cashout

```bash
curl -X POST http://localhost:8080/api/cashout \
  -H "Content-Type: application/json" \
  -d '{
    "cashoutId": "cashout-003",
    "ticketId": "multi-001",
    "ticketSignature": "signature-from-original-response",
    "type": "bet",
    "code": 100,
    "betId": "bet-1",
    "payout": [
      {
        "type": "cash",
        "currency": "EUR",
        "amount": 12.00,
        "source": "cashout"
      }
    ]
  }'
```

---

## Python 完整示例

```python
import requests
import json
from datetime import datetime

class MTSClient:
    def __init__(self, base_url="http://localhost:8080"):
        self.base_url = base_url
        self.session = requests.Session()
    
    def generate_ticket_id(self, prefix="ticket"):
        timestamp = int(datetime.now().timestamp())
        return f"{prefix}-{timestamp}"
    
    def place_single_bet(self, event_id, market_id, outcome_id, odds, amount, currency="EUR"):
        url = f"{self.base_url}/api/bets/single"
        payload = {
            "ticketId": self.generate_ticket_id("single"),
            "selection": {
                "productId": "3",
                "eventId": event_id,
                "marketId": market_id,
                "outcomeId": outcome_id,
                "odds": odds
            },
            "stake": {
                "type": "cash",
                "currency": currency,
                "amount": amount,
                "mode": "total"
            }
        }
        response = self.session.post(url, json=payload)
        return response.json()
    
    def place_accumulator(self, selections, amount, currency="EUR"):
        url = f"{self.base_url}/api/bets/accumulator"
        payload = {
            "ticketId": self.generate_ticket_id("acc"),
            "selections": selections,
            "stake": {
                "type": "cash",
                "currency": currency,
                "amount": amount,
                "mode": "total"
            }
        }
        response = self.session.post(url, json=payload)
        return response.json()
    
    def place_trixie(self, selections, unit_stake, currency="EUR"):
        url = f"{self.base_url}/api/bets/preset"
        payload = {
            "ticketId": self.generate_ticket_id("trixie"),
            "type": "trixie",
            "selections": selections,
            "stake": {
                "type": "cash",
                "currency": currency,
                "amount": unit_stake,
                "mode": "unit"
            }
        }
        response = self.session.post(url, json=payload)
        return response.json()

# 使用示例
if __name__ == "__main__":
    client = MTSClient()
    
    # 单注
    result = client.place_single_bet(
        event_id="sr:match:12345",
        market_id="1",
        outcome_id="1",
        odds=2.50,
        amount=10.00
    )
    print("Single Bet Result:")
    print(json.dumps(result, indent=2))
    
    # 串关
    selections = [
        {"productId": "3", "eventId": "sr:match:12345", "marketId": "1", "outcomeId": "1", "odds": 2.50},
        {"productId": "3", "eventId": "sr:match:12346", "marketId": "1", "outcomeId": "2", "odds": 1.80},
        {"productId": "3", "eventId": "sr:match:12347", "marketId": "1", "outcomeId": "1", "odds": 3.00}
    ]
    result = client.place_accumulator(selections, 10.00)
    print("\nAccumulator Result:")
    print(json.dumps(result, indent=2))
```

---

## 测试脚本

项目包含一个自动化测试脚本，可以测试所有端点：

```bash
# 运行测试脚本
./scripts/test_api.sh

# 指定自定义 URL
BASE_URL=http://your-server:8080 ./scripts/test_api.sh
```

测试脚本会自动测试所有端点并显示结果。

---

## 错误处理

### 验证错误示例

**请求**:
```bash
curl -X POST http://localhost:8080/api/bets/single \
  -H "Content-Type: application/json" \
  -d '{
    "selection": {
      "productId": "3",
      "eventId": "sr:match:12345",
      "marketId": "1",
      "outcomeId": "1",
      "odds": 2.50
    }
  }'
```

**响应** (HTTP 400):
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

### MTS 拒绝示例

**响应** (HTTP 200, 但 MTS 拒绝):
```json
{
  "success": true,
  "data": {
    "content": {
      "type": "ticket-reply",
      "ticketId": "single-001",
      "status": "rejected",
      "code": -401,
      "message": "Match is not found in MTS"
    }
  }
}
```

---

## 注意事项

1. **TicketID 唯一性**: 每个 `ticketId` 必须全局唯一
2. **Odds 格式**: 使用十进制格式（如 2.50）
3. **Stake Mode**: 
   - `total` 用于单注和串关
   - `unit` 用于系统串
4. **预设类型选项数**: 必须与类型要求完全匹配
5. **Cashout Signature**: 必须使用原始注单响应中的 `signature` 字段

---

## 更多信息

完整的 API 文档请参阅 [API_DOCUMENTATION.md](API_DOCUMENTATION.md)。
