# `Selection` 对象详解

**返回索引**: [主文档](./main.md)

---

## 1. 概述

`Selection` 是 MTS 注单体系中**最核心、最复杂的对象**。它代表一个用户的投注选择，可以是单一的投注项，也可以是嵌套的、包含其他 `Selection` 的复杂结构（如系统串）。

`Selection` 的行为和字段构成由其 `type` 字段决定。

## 2. `type` 字段详解

`type` 字段决定了 `Selection` 的类型和结构。以下是所有支持的类型：

| `type` 值 | 描述 | 使用场景 |
|:---|:---|:---|
| `uf` | **Unified Feed (标准)** | 最常见的类型，用于标准的体育赛事投注。 |
| `external` | **外部** | 用于非 Sportradar 提供的赛事或自定义事件。 |
| `uf-custom-bet` | **自定义串关** | 用于将多个来自同一赛事的投注项组合成一个自定义赔率的串关。 |
| `system` | **系统串** | 用于构建系统串或 Banker，其 `selections` 字段会包含其他 `Selection` 对象。 |

## 3. 字段构成

`Selection` 对象的字段根据其 `type` 的不同而变化。

### 3.1. 标准类型 (`uf`, `external`, `uf-custom-bet`)

当 `type` 为 `uf`、`external` 或 `uf-custom-bet` 时，`Selection` 代表一个具体的投注项。

#### 字段表

| 字段 | Go 类型 | JSON 类型 | 是否必需 | 描述与业务意义 |
|:---|:---|:---|:---:|:---|
| `type` | `string` | `string` | ✅ 是 | **类型标识**。必须是 `uf`, `external` 或 `uf-custom-bet`。 |
| `productId` | `string` | `string` | 否 | **产品 ID**。通常为 `"3"`，代表体育博彩。 |
| `eventId` | `string` | `string` | ✅ 是 | **赛事 ID**。全局唯一的赛事标识，如 `"sr:match:12345"`。 |
| `marketId` | `string` | `string` | ✅ 是 | **盘口 ID**。赛事内唯一的盘口标识，如 `"1"` (1X2)。 |
| `outcomeId` | `string` | `string` | ✅ 是 | **结果 ID**。盘口内唯一的结果标识，如 `"1"` (主胜)。 |
| `specifiers` | `string` | `string` | 否 | **盘口说明符**。用于参数化盘口，如大小球的 `"total=2.5"` 或让分盘的 `"hcp=-1.5"`。 |
| `odds` | `*Odds` | `object` | ✅ 是 | **赔率对象**。指向一个 `Odds` 对象，包含赔率类型和值。详情见 [Odds.md](./Odds.md)。 |

#### JSON 示例

```json
{
  "type": "uf",
  "productId": "3",
  "eventId": "sr:match:12345",
  "marketId": "1",
  "outcomeId": "1",
  "specifiers": "hcp=-1.5",
  "odds": {
    "type": "decimal",
    "value": "1.95"
  }
}
```

### 3.2. 系统串类型 (`system`)

当 `type` 为 `system` 时，`Selection` 作为一个容器，用于构建复杂的系统串或 Banker 投注。

#### 字段表

| 字段 | Go 类型 | JSON 类型 | 是否必需 | 描述与业务意义 |
|:---|:---|:---|:---:|:---|
| `type` | `string` | `string` | ✅ 是 | **类型标识**。必须是 `system`。 |
| `size` | `[]int` | `array` | ✅ 是 | **组合大小**。定义了系统串的组合方式。例如 `[2, 3]` 表示这是一个 2/N 和 3/N 的系统串。 |
| `selections` | `[]Selection` | `array` | ✅ 是 | **嵌套的 Selections**。一个 `Selection` 对象数组，包含了构成系统串的所有投注项。 |

#### JSON 示例：一个 2/3 系统串 (Doubles from 3)

```json
{
  "type": "system",
  "size": [2],
  "selections": [
    {
      "type": "uf",
      "eventId": "sr:match:11111",
      "marketId": "1",
      "outcomeId": "1",
      "odds": { "type": "decimal", "value": "1.50" }
    },
    {
      "type": "uf",
      "eventId": "sr:match:22222",
      "marketId": "1",
      "outcomeId": "2",
      "odds": { "type": "decimal", "value": "2.00" }
    },
    {
      "type": "uf",
      "eventId": "sr:match:33333",
      "marketId": "1",
      "outcomeId": "1",
      "odds": { "type": "decimal", "value": "1.80" }
    }
  ]
}
```

### 3.3. Banker 系统串的特殊结构

Banker 注单通过嵌套的 `system` 类型 `Selection` 实现。结构如下：

- **外层 `system`**: 定义了普通 `selections` 的组合方式。
- **内层 `system`**: 作为 `banker` 的容器，其 `size` 永远是 `[1]`。

#### JSON 示例：Banker 2/3 系统串

```json
{
  "type": "system",
  "size": [3], // Banker + 2/3 = 3-folds
  "selections": [
    {
      "type": "system", // Banker 容器
      "size": [1], // Banker 的 size 永远是 [1]
      "selections": [
        { // Banker 投注项
          "type": "uf",
          "eventId": "sr:match:99999",
          "marketId": "1",
          "outcomeId": "1",
          "odds": { "type": "decimal", "value": "1.20" }
        }
      ]
    },
    {
      "type": "system", // 普通 selections 容器
      "size": [2], // 2/3 系统串
      "selections": [
        { ...普通 selection 1... },
        { ...普通 selection 2... },
        { ...普通 selection 3... }
      ]
    }
  ]
}
```

## 4. 最佳实践

- **始终验证 `type`**: 在处理 `Selection` 对象时，首先检查其 `type` 字段。
- **递归处理**: 处理 `system` 类型的 `Selection` 时，需要使用递归或循环来遍历其嵌套的 `selections`。
- **`specifiers` 的使用**: 仔细阅读 MTS 文档，了解不同盘口支持的 `specifiers` 格式。
- **Banker 结构**: 实现 Banker 功能时，务必遵循双层 `system` 嵌套的结构。

---

**相关文档**:
- [Odds.md](./Odds.md)
- [Bet.md](./Bet.md)
