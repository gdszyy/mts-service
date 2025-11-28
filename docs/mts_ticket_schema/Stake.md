# `Stake` 对象详解

**返回索引**: [主文档](./main.md)

---

## 1. 概述

`Stake` 对象定义了投注的**金额、货币和计算模式**。它是 `Bet` 对象中的一个必需字段，并且是一个数组，尽管在大多数情况下，这个数组只包含一个元素。

## 2. 字段构成

| 字段 | Go 类型 | JSON 类型 | 是否必需 | 描述与业务意义 |
|:---|:---|:---|:---:|:---|
| `type` | `string` | `string` | ✅ 是 | **金额类型**。定义了这笔金额的性质。常见值为 `"cash"` (现金)。 |
| `currency` | `string` | `string` | ✅ 是 | **货币代码**。遵循 ISO 4217 标准的 3 位字母代码，如 `"EUR"`、`"USD"`。 |
| `amount` | `string` | `string` | ✅ 是 | **金额**。以字符串形式表示，以避免浮点数精度问题。例如 `"10.00"`。 |
| `mode` | `string` | `string` | 否 | **计算模式**。决定了 `amount` 如何应用于投注。这是**区分单注和系统串投注金额计算方式的关键字段**。 |

## 3. `mode` 字段详解

`mode` 字段是 `Stake` 对象中最具业务意义的字段。它有两个主要值：`total` 和 `unit`。

### 3.1. `mode: "total"` (总金额模式)

- **描述**: `amount` 是这笔 `Bet` 的**总投注金额**。
- **使用场景**: 
  - **单注 (Single Bet)**
  - **串关 (Accumulator)**
- **业务意义**: 无论这笔 `Bet` 内部有多少组合（对于单注和串关，只有一个组合），总的投注成本就是 `amount` 的值。

#### JSON 示例：一个 10 EUR 的串关

```json
{
  "stake": [
    {
      "type": "cash",
      "currency": "EUR",
      "amount": "10.00",
      "mode": "total"
    }
  ]
}
```

### 3.2. `mode: "unit"` (单位金额模式)

- **描述**: `amount` 是这笔 `Bet` 中**每一个组合 (combination) 的投注金额**。
- **使用场景**: 
  - **系统串 (System Bet)**
  - **Banker 系统串**
- **业务意义**: 总投注金额需要通过计算得出：**总金额 = `amount` × 组合数量**。

#### JSON 示例：一个 2/3 系统串 (Trixie)，每组合 5 EUR

一个 Trixie 包含 4 个组合（3 个 doubles + 1 个 treble）。

```json
{
  "stake": [
    {
      "type": "cash",
      "currency": "EUR",
      "amount": "5.00",
      "mode": "unit"
    }
  ]
}
```

**总投注金额计算**: `5.00 EUR` (amount) × `4` (组合) = `20.00 EUR`。

## 4. 为什么 `stake` 是一个数组？

MTS API 将 `stake` 设计为一个数组，是为了支持未来可能出现的、更复杂的投注金额构成，例如：

- **混合支付**: 一部分用现金，一部分用赠金（bonus）。
- **多种货币**: 尽管当前 MTS 不支持单笔 `Bet` 中使用多种货币，但该结构为此留下了扩展空间。

在当前实现中，我们通常只向这个数组中添加一个 `Stake` 对象。

#### JSON 示例：混合支付（理论上）

```json
{
  "stake": [
    {
      "type": "cash",
      "currency": "EUR",
      "amount": "8.00",
      "mode": "total"
    },
    {
      "type": "bonus",
      "currency": "EUR",
      "amount": "2.00",
      "mode": "total"
    }
  ]
}
```

## 5. 最佳实践

- **使用字符串表示金额**: 始终使用字符串来处理 `amount`，以避免浮点数计算带来的精度问题。
- **明确 `mode`**: 在构建 `Bet` 对象时，根据注单类型明确设置 `mode`。单注和串关使用 `total`，系统串使用 `unit`。
- **计算总金额**: 在前端向用户显示总投注金额时，如果是 `unit` 模式，务必正确计算组合数量并乘以单位金额。
- **默认行为**: 如果 `mode` 字段被省略，MTS 会默认将其视为 `total`。

---

**相关文档**:
- [Bet.md](./Bet.md)
- [Selection.md](./Selection.md)
