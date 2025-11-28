# 其他核心对象详解

**返回索引**: [主文档](./main.md)

---

本文档包含了除 `Selection`, `Stake`, `Bet` 之外的其他核心对象的详细说明。

## 1. `TicketRequest` 对象

这是注单请求的顶层结构，包含了信封和内容。

| 字段 | Go 类型 | JSON 类型 | 是否必需 | 描述与业务意义 |
|:---|:---|:---|:---:|:---|
| `operatorId` | `int64` | `number` | ✅ 是 | **操作员 ID**。由 Sportradar 提供，用于标识您的账户。 |
| `correlationId` | `string` | `string` | ✅ 是 | **关联 ID**。客户端生成的唯一字符串，用于将请求与响应配对。建议使用 UUID。 |
| `timestampUtc` | `int64` | `number` | ✅ 是 | **时间戳**。客户端提交请求时的 Unix 毫秒时间戳。 |
| `operation` | `string` | `string` | ✅ 是 | **操作类型**。对于注单请求，必须是 `"ticket-placement"`。 |
| `version` | `string` | `string` | ✅ 是 | **协议版本**。当前必须是 `"3.0"`。 |
| `content` | `TicketContent` | `object` | ✅ 是 | **内容对象**。包含了注单的实际内容。 |

---

## 2. `TicketContent` 对象

这是注单的核心内容。

| 字段 | Go 类型 | JSON 类型 | 是否必需 | 描述与业务意义 |
|:---|:---|:---|:---:|:---|
| `type` | `string` | `string` | ✅ 是 | **内容类型**。对于注单请求，必须是 `"ticket"`。 |
| `ticketId` | `string` | `string` | ✅ 是 | **注单 ID**。客户端生成的唯一 ID，用于标识这张注单。**MTS 会使用此 ID 进行幂等性检查**，防止重复提交。 |
| `bets` | `[]Bet` | `array` | ✅ 是 | **投注数组**。包含一个或多个 `Bet` 对象。详情见 [Bet.md](./Bet.md)。 |
| `context` | `*Context` | `object` | 否 | **上下文信息**。可选，但建议提供，用于风控和分析。 |

---

## 3. `Odds` 对象

定义了赔率的类型和值。

| 字段 | Go 类型 | JSON 类型 | 是否必需 | 描述与业务意义 |
|:---|:---|:---|:---:|:---|
| `type` | `string` | `string` | ✅ 是 | **赔率类型**。最常用的是 `"decimal"` (欧盘)。 |
| `value` | `string` | `string` | ✅ 是 | **赔率值**。以字符串形式表示，以避免浮点数精度问题。例如 `"1.95"`。 |

#### JSON 示例

```json
{
  "odds": {
    "type": "decimal",
    "value": "1.95"
  }
}
```

---

## 4. `Context` 对象

提供了关于投注来源的上下文信息。

| 字段 | Go 类型 | JSON 类型 | 是否必需 | 描述与业务意义 |
|:---|:---|:---|:---:|:---|
| `channel` | `*Channel` | `object` | 否 | **渠道对象**。描述了投注的来源渠道。 |
| `ip` | `string` | `string` | 否 | **IP 地址**。用户的公网 IP 地址。 |
| `limitId` | `int64` | `number` | 否 | **限额 ID**。用于关联特定的用户限额策略。 |

### 4.1. `Channel` 对象

| 字段 | Go 类型 | JSON 类型 | 是否必需 | 描述与业务意义 |
|:---|:---|:---|:---:|:---|
| `type` | `string` | `string` | ✅ 是 | **渠道类型**。例如 `"internet"` (PC 网站), `"mobile"` (移动应用), `"agent"` (代理)。 |
| `lang` | `string` | `string` | ✅ 是 | **语言代码**。用户的语言偏好，如 `"EN"`, `"ZH"`。 |

#### JSON 示例

```json
{
  "context": {
    "channel": {
      "type": "internet",
      "lang": "EN"
    },
    "ip": "123.45.67.89",
    "limitId": 4268
  }
}
```

---

**相关文档**:
- [主文档](./main.md)
