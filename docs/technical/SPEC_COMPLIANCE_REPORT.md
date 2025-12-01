# MTS 规范符合性审查报告

**版本**: 1.0  
**日期**: 2025年11月27日  
**审查范围**: `gdszyy/mts-service` 项目 `main` 分支
**审查依据**: [MTS 注单结构化文档体系 (Notion)](https://www.notion.so/2b9bacf5f267813ba42ed3f60f1c0356)

---

## 1. 总体评估

项目在**核心数据模型 (`internal/models`)** 和 **TicketBuilder (`internal/models/ticket_builder.go`)** 层面**高度符合** MTS 规范。`Selection` 的递归结构、`Stake` 的 `mode` 字段、以及 Banker 的双层嵌套实现都非常正确。

主要问题集中在 **API 层 (`internal/api`)**，存在一些与规范不一致的地方，以及可以改进的细节。

## 2. 详细审查结果

### 2.1. 数据模型 (`internal/models`) - ✅ 高度符合

| 检查项 | 状态 | 备注 |
|:---|:---:|:---|
| `TicketRequest` 结构 | ✅ 符合 | 包含了所有必需的信封和内容字段。 |
| `Selection` 结构 | ✅ 符合 | 正确实现了 `type` 驱动的联合结构，支持标准和 `system` 类型。 |
| `Stake` 结构 | ✅ 符合 | 正确包含了 `mode` 字段，用于区分 `total` 和 `unit`。 |
| `Bet` 结构 | ✅ 符合 | 正确包含了 `selections` 和 `stake` 数组。 |
| Banker 实现 | ✅ 符合 | `TicketBuilder` 中正确实现了双层 `system` 嵌套。 |

**结论**: 数据模型层是健全的，完全遵循了 MTS 规范。

### 2.2. API 层 (`internal/api`) - ⚠️ 中度符合 (需要修改)

API 层是外部系统与 `mts-service` 交互的入口，其设计直接影响了系统的易用性和健壮性。当前实现存在以下主要问题：

#### 问题 1：金额和赔率使用 `float64` 而非 `string`

- **位置**: `internal/api/request_models.go`
- **不符合项**: 
  - `SelectionRequest.Odds` 是 `float64`
  - `StakeRequest.Amount` 是 `float64`
- **风险**: 使用浮点数处理货币和赔率是**极其危险**的，会导致精度丢失和计算错误。例如，`0.1 + 0.2` 在二进制浮点数中不等于 `0.3`。
- **规范要求**: 金额和赔率必须使用**字符串**来保证精度。
- **修改建议**: 
  - 将 `SelectionRequest.Odds` 修改为 `string` 类型。
  - 将 `StakeRequest.Amount` 修改为 `string` 类型。
  - 在 API handler 中，使用 `strconv.ParseFloat` 将字符串转换为 `float64` 进行内部计算，或直接传递字符串给 `NewStake` 和 `NewSelection`。

#### 问题 2：Banker 实现不符合 MTS 规范

- **位置**: `internal/models/ticket_builder.go` (AddBankerSystemBet)
- **不符合项**: 当前的 `AddBankerSystemBet` 方法将 `bankers` 和 `systemSelection` 直接放在了同一个 `selections` 数组中。这不符合 MTS 对 Banker 的双层嵌套 `system` 要求。
- **规范要求**: Banker 注单需要一个外层的 `system` selection，其 `selections` 包含两个内层的 `system` selection：一个用于 `bankers` (size=[1])，另一个用于普通 `selections`。
- **修改建议**: 
  - 重写 `AddBankerSystemBet` 方法。
  - 创建一个 `bankerSystemSelection` (`type="system", size=[1]`) 来包裹 `bankers`。
  - 创建一个 `mainSystemSelection` (`type="system", size=size`) 来包裹 `selections`。
  - 创建一个顶层的 `rootSystemSelection` (`type="system"`)，其 `selections` 包含 `bankerSystemSelection` 和 `mainSystemSelection`。
  - 更新 `size` 的计算逻辑，应该是 `len(bankers) + s` 对于 `size` 中的每个 `s`。

#### 问题 3：Context 对象结构不一致

- **位置**: `internal/api/request_models.go`
- **不符合项**: `ContextRequest` 是一个扁平结构 (`ChannelType`, `Language`, `IP`)。
- **规范要求**: `Context` 对象包含一个嵌套的 `Channel` 对象 (`{ "type": "...", "lang": "..." }`)。
- **修改建议**: 
  - 创建一个新的 `ChannelRequest` 结构体。
  - 将 `ContextRequest` 修改为包含 `*ChannelRequest` 和 `IP` 字段。

### 2.3. TicketBuilder (`internal/models/ticket_builder.go`) - ✅ 高度符合 (除 Banker 外)

| 检查项 | 状态 | 备注 |
|:---|:---:|:---|
| 单注/串关/系统串 | ✅ 符合 | 实现正确，逻辑清晰。 |
| 预设系统串 | ✅ 符合 | 正确调用了 `AddSystemBet`。 |
| Banker 系统串 | ❌ **不符合** | **严重问题**。实现方式与 MTS 规范不符，无法被 MTS 正确解析。 |
| `NewSelection` / `NewStake` | ⚠️ 中度符合 | 接受 `float64` 作为参数，应考虑接受 `string` 以强制 API 层使用字符串。 |

**结论**: `TicketBuilder` 的大部分实现是优秀的，但 Banker 的实现是错误的，需要立即修正。

## 3. 修改建议总结

### 优先级 1：关键性修复 (必须修改)

1.  **修复 Banker 实现**: 重写 `AddBankerSystemBet` 方法，实现正确的双层 `system` 嵌套结构。
2.  **API 金额和赔率改用字符串**: 将 `request_models.go` 中的 `float64` 全部修改为 `string`，并在 handler 中进行转换。

### 优先级 2：规范一致性 (建议修改)

3.  **重构 ContextRequest**: 使其包含嵌套的 `ChannelRequest` 对象，与 MTS 规范保持一致。

### 优先级 3：最佳实践 (可选优化)

4.  **更新构建器辅助函数**: 修改 `NewSelection` 和 `NewStake`，使其接受 `string` 类型的金额和赔率，从源头上杜绝浮点数问题。

## 4. 下一步行动

建议按照以上优先级顺序进行代码修改。修复完成后，应重新运行测试并部署。

---
**审查人机协作，使命必达。
