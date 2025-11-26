# MTS 功能实现与重构总结

## 1. 项目目标

本次重构的核心目标是**完全移除原有的本地投注系统 (`betting-system`)**，并基于 Sportradar MTS Transaction 3.0 API 规范，从核心模型开始，重新实现一个纯粹、完整且功能强大的 MTS 对接服务。这包括支持所有关键的注单类型和 Cashout 功能。

## 2. 主要变更

### 2.1. 移除本地系统

- **已完成**: 整个 `betting-system` 目录已被彻底移除。
- **影响**: 项目现在不再依赖任何本地数据库或独立的投注引擎，所有业务逻辑直接围绕 MTS API 构建，成为一个轻量级、无状态的代理服务。

### 2.2. 重构核心数据模型

为了支持 MTS 的所有功能，我们对数据模型进行了全面重构和扩展。

**文件: `internal/models/ticket.go`**

- **`Selection` 结构体**: 这是本次重构的核心。该结构体被重新设计，以支持**嵌套**，从而实现对系统串的支持。

  ```go
  type Selection struct {
      // ... 标准选项字段 (productId, eventId, etc.)

      // --- 新增字段用于系统串 ---
      Size       []int       `json:"size,omitempty"`       // 组合大小, e.g., [2, 3]
      Selections []Selection `json:"selections,omitempty"` // 嵌套的原始选项
  }
  ```

**文件: `internal/models/cashout.go` (新增)**

- **已完成**: 创建了全新的 `cashout.go` 文件，其中包含了与 Cashout 流程相关的**所有请求和响应结构体**，例如 `CashoutRequest`, `CashoutResponse`, `CashoutDetail` 等。这为后续实现 Cashout 功能奠定了基础。

### 2.3. 引入注单构建器 (Ticket Builder)

为了简化和规范化复杂注单的创建过程，我们引入了一个强大的构建器。

**文件: `internal/models/ticket_builder.go` (新增)**

- **目的**: 提供一个流畅、易于使用的 API 来创建所有类型的 MTS 注单请求。
- **核心功能**:
  - 支持链式调用: `NewTicketBuilder(...).AddSingleBet(...).Build(...)`
  - **封装了所有注单类型的复杂逻辑**，提供了简单的方法调用：
    - `AddSingleBet()`
    - `AddAccumulatorBet()`
    - `AddSystemBet()`
    - `AddBankerSystemBet()`
    - **预设的系统串方法**: `AddTrixieBet()`, `AddYankeeBet()`, `AddLucky15Bet()` 等。

- **示例用法**:

  ```go
  // 创建一个 2/3 系统串
  builder := NewTicketBuilder(operatorID, ticketID)
  builder.AddSystemBet([]int{2}, selections, unitStake)
  ticketRequest := builder.Build(correlationID)
  ```

### 2.4. 增加全面的单元测试

**文件: `internal/models/ticket_builder_test.go` (新增)**

- **已完成**: 创建了针对 `TicketBuilder` 的完整单元测试套件。
- **测试覆盖范围**:
  - ✅ 单注 (Single Bet)
  - ✅ 串关 (Accumulator)
  - ✅ 标准系统串 (System Bet 2/4)
  - ✅ 预设系统串 (Trixie, Yankee, Lucky 15)
  - ✅ Banker 注单
  - ✅ 单个票据中包含多种类型的注单
- **目的**: 确保所有通过构建器生成的 JSON 请求都**严格符合 MTS API 规范**，从根本上杜绝格式错误。

## 3. 当前实现状态

| 功能 | 状态 | 文件 | 备注 |
| :--- | :--- | :--- | :--- |
| **移除本地系统** | ✅ 完成 | N/A | 项目结构已清理。 |
| **单注模型** | ✅ 完成 | `ticket.go` | 结构经验证，功能完整。 |
| **串关模型** | ✅ 完成 | `ticket.go` | 通过 `Selection` 数组实现。 |
| **系统串模型** | ✅ 完成 | `ticket.go` | 通过嵌套 `Selection` 和 `size` 字段实现。 |
| **Banker 模型** | ✅ 完成 | `ticket.go` | 通过将标准 `Selection` 和 `system` 类型的 `Selection` 并列实现。 |
| **Cashout 模型** | ✅ 完成 | `cashout.go` | 已定义所有请求和响应结构。 |
| **注单构建器** | ✅ 完成 | `ticket_builder.go` | 已实现所有注单类型的构建逻辑。 |
| **构建器单元测试** | ✅ 完成 | `ticket_builder_test.go` | 已验证所有注单类型的 JSON 输出。 |
| **服务层集成** | ❌ **未开始** | `mts.go` | `TicketBuilder` 尚未集成到 WebSocket 服务中。 |
| **Cashout 逻辑** | ❌ **未开始** | `mts.go` | Cashout 请求和响应逻辑尚未实现。 |
| **API 端点** | ❌ **未开始** | `handlers.go` | 需要创建新的 API 端点来暴露这些新功能。 |

## 4. 后续步骤

现在，我们已经完成了坚实的基础模型和构建器。下一步是按顺序完成服务层和 API 层的实现：

1.  **集成 `TicketBuilder`**: 修改 `mts.go`，创建一个统一的 `PlaceTicket` 方法，该方法接收一个 `TicketRequest` 对象（由 `TicketBuilder` 生成）并将其发送到 MTS。

2.  **实现串关功能**: 创建一个新的 API 端点（例如 `POST /api/bets/accumulator`），该端点调用 `TicketBuilder` 的 `AddAccumulatorBet` 方法，并通过 `PlaceTicket` 发送请求。

3.  **实现系统串功能**: 创建 `POST /api/bets/system` 端点，调用 `AddSystemBet` 或其他预设方法。

4.  **实现 Banker 功能**: 扩展系统串端点，使其能够接收 Banker 选项，并调用 `AddBankerSystemBet`。

5.  **实现 Cashout 功能**: 
    - 创建 `POST /api/tickets/{id}/cashout` 端点。
    - 在 `mts.go` 中实现 `RequestCashout` 方法，用于发送 `cashout-inform` 请求。
    - 在 WebSocket 消息处理循环中添加对 `cashout-inform-reply` 的处理逻辑。

6.  **最终测试与验证**: 编写集成测试，模拟完整的 API 调用流程，确保所有功能按预期工作。

7.  **提交与部署**: 将所有修改提交到 `main` 分支，并准备部署。
