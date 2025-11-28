# `Bet` 对象详解

**返回索引**: [主文档](./main.md)

---

## 1. 概述

`Bet` 对象代表一个独立的投注单元。一个注单 (`Ticket`) 可以包含一个或多个 `Bet` 对象。这种设计允许在一次请求中提交多种不同类型的投注，例如同时提交一个单注和一个系统串。

## 2. 字段构成

| 字段 | Go 类型 | JSON 类型 | 是否必需 | 描述与业务意义 |
|:---|:---|:---|:---:|:---|
| `selections` | `[]Selection` | `array` | ✅ 是 | **投注选项数组**。包含一个或多个 `Selection` 对象。这是 `Bet` 的核心。详情见 [Selection.md](./Selection.md)。 |
| `stake` | `[]Stake` | `array` | ✅ 是 | **投注金额数组**。包含一个或多个 `Stake` 对象，定义了这笔 `Bet` 的金额。详情见 [Stake.md](./Stake.md)。 |

## 3. `Bet` 的结构与注单类型

`Bet` 的结构，特别是其 `selections` 数组的内容，直接决定了这笔投注的类型。

### 3.1. 单注 (Single Bet)

- `selections` 数组包含**一个**标准的 `Selection` 对象 (`type="uf"`)。
- `stake` 数组中的 `mode` 通常是 `"total"`。

#### JSON 示例

```json
{
  "bets": [
    {
      "selections": [
        {
          "type": "uf",
          "eventId": "sr:match:12345",
          // ... 其他字段
        }
      ],
      "stake": [
        {
          "type": "cash",
          "amount": "10.00",
          "mode": "total"
        }
      ]
    }
  ]
}
```

### 3.2. 串关 (Accumulator)

- `selections` 数组包含**多个**标准的 `Selection` 对象 (`type="uf"`)。
- `stake` 数组中的 `mode` 是 `"total"`。

#### JSON 示例

```json
{
  "bets": [
    {
      "selections": [
        { "type": "uf", "eventId": "sr:match:11111", ... },
        { "type": "uf", "eventId": "sr:match:22222", ... },
        { "type": "uf", "eventId": "sr:match:33333", ... }
      ],
      "stake": [
        {
          "type": "cash",
          "amount": "10.00",
          "mode": "total"
        }
      ]
    }
  ]
}
```

### 3.3. 系统串 (System Bet)

- `selections` 数组只包含**一个** `Selection` 对象，且其 `type` 是 `"system"`。
- `stake` 数组中的 `mode` 是 `"unit"`。

#### JSON 示例

```json
{
  "bets": [
    {
      "selections": [
        {
          "type": "system",
          "size": [2],
          "selections": [
            { "type": "uf", "eventId": "sr:match:11111", ... },
            { "type": "uf", "eventId": "sr:match:22222", ... },
            { "type": "uf", "eventId": "sr:match:33333", ... }
          ]
        }
      ],
      "stake": [
        {
          "type": "cash",
          "amount": "5.00",
          "mode": "unit"
        }
      ]
    }
  ]
}
```

## 4. 为什么 `bets` 是一个数组？

将 `bets` 设计为一个数组，允许在一次 `Ticket` 请求中提交多个独立的投注。这对于需要同时下多种类型注单的用户非常有用。

### JSON 示例：同时提交一个单注和一个 Trixie

```json
{
  "ticketId": "multi-bet-001",
  "bets": [
    // 第一个 Bet: 单注
    {
      "selections": [
        { "type": "uf", "eventId": "sr:match:99999", ... }
      ],
      "stake": [
        { "type": "cash", "amount": "50.00", "mode": "total" }
      ]
    },
    // 第二个 Bet: Trixie 系统串
    {
      "selections": [
        {
          "type": "system",
          "size": [2, 3],
          "selections": [
            { "type": "uf", "eventId": "sr:match:11111", ... },
            { "type": "uf", "eventId": "sr:match:22222", ... },
            { "type": "uf", "eventId": "sr:match:33333", ... }
          ]
        }
      ],
      "stake": [
        { "type": "cash", "amount": "10.00", "mode": "unit" }
      ]
    }
  ]
}
```

## 5. 最佳实践

- **单一职责**: 尽量让每个 `Bet` 对象只代表一种投注类型（单注、串关或一种系统串）。
- **使用 `bets` 数组**: 当用户在投注单中添加了多种投注类型时（例如，一个 2/3 系统串和一个 3/3 的串关），应将它们创建为两个独立的 `Bet` 对象，并放入同一个 `bets` 数组中。
- **清晰的 API 设计**: 在设计对外 API 时，可以为不同类型的 `Bet` 创建不同的端点（如 `/bets/single`, `/bets/system`），然后在后端将它们统一封装成 `Bet` 对象。

---

**相关文档**:
- [Selection.md](./Selection.md)
- [Stake.md](./Stake.md)
- [TicketContent.md](./TicketContent.md)
