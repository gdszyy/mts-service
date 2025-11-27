# MTS 子弹装填流程 (Bullet Loading Workflow)

**版本**: 1.0  
**日期**: 2025年11月27日  
**作者**: Manus AI

---

## 1. 概念介绍

**子弹装填 (Bullet Loading)** 是一个行业术语，用于描述在投注系统中准备和提交注单的完整流程。它将注单视为一颗“子弹”，将提交流程视为“装填”动作。这个过程确保了每一笔投注在发送给上游系统（如 Sportradar MTS）之前，都经过了充分的**构建、验证和格式化**，以最大限度地提高接受率并减少错误。

一个健壮的子弹装填流程是任何成功投注平台的核心，它直接影响到系统的**可靠性、响应速度和用户体验**。

## 2. 核心原则

子弹装填流程遵循以下核心原则：

- **数据准确性**: 所有提交的数据（如赛事 ID、赔率、金额）必须准确无误。
- **结构标准化**: 注单的 JSON 结构必须严格遵守 MTS Transaction 3.0 API 规范 [1]。
- **流程原子性**: 一次完整的装填和发射（提交）应被视为一个原子操作。
- **状态可追溯**: 从构建到最终确认，注单的每一个状态都应可被追踪。

## 3. 子弹装填的四个阶段

完整的子弹装填流程可以分为四个主要阶段：**构建 (Build)**、**验证 (Validate)**、**发射 (Fire)** 和 **确认 (Confirm)**。

### 阶段一：构建 (Build) - 准备弹药

这是流程的起点，目标是根据用户输入或系统逻辑，创建一个符合 MTS 规范的、结构化的注单对象。

#### 关键步骤

1.  **收集输入**: 从前端 API 请求中获取所有必要信息，包括：
    - **Selections**: 用户选择的投注项（赛事、市场、结果）。
    - **Stake**: 投注金额和模式（`unit` 或 `total`）。
    - **Bet Type**: 注单类型（单注、串关、系统串等）。
    - **User Context**: 用户信息、渠道、IP 地址等。

2.  **选择构建器**: 根据注单类型，使用相应的构建器函数。在我们的 `mts-service` 项目中，这由 `TicketBuilder` 完成。

    | 注单类型 | 构建器方法 |
    |:---|:---|
    | 单注 | `NewSingleBetTicket()` |
    | 串关 | `NewAccumulatorTicket()` |
    | 系统串 | `NewSystemBetTicket()` |
    | Banker | `NewBankerSystemBetTicket()` |
    | 预设系统串 | `NewPresetSystemBetTicket()` |

3.  **生成注单结构**: 构建器将输入信息转换为一个完整的 `TicketRequest` 对象。这个过程包括：
    - 创建唯一的 `ticketId` 和 `correlationId`。
    - 嵌套 `selections` 以符合系统串或 Banker 的要求。
    - 设置正确的 `stake` 模式。
    - 填充所有必要的信封（envelope）字段，如 `operatorId`、`timestampUtc` 等。

#### 示例：构建一个 Trixie 系统串

```go
// 1. 收集输入 (3 个 selections, 10 EUR unit stake)

// 2. 选择构建器
builder := models.NewTicketBuilder(cfg.OperatorID)

// 3. 生成注单结构
ticket, err := builder.NewPresetSystemBetTicket(
    "trixie-001",         // Ticket ID
    models.Trixie,        // Bet Type
    selections,           // 3 selections
    stake,                // 10 EUR unit stake
    context,              // User context
)
```

### 阶段二：验证 (Validate) - 检查膛线

在发射之前，必须对构建好的注单进行严格的本地验证，以捕捉明显的错误，避免不必要的 API 调用。

#### 关键步骤

1.  **结构验证**: 检查注单对象是否包含所有必需字段。
    - `ticketId` 是否为空？
    - `bets` 和 `selections` 数组是否至少有一个元素？

2.  **逻辑验证**: 检查注单是否符合业务规则。
    - **系统串检查**: 一个 Trixie（2/3 系统串）是否正好有 3 个 `selections`？
    - **Banker 检查**: Banker 注单是否至少包含 1 个 Banker 和 2 个普通 `selections`？
    - **金额检查**: 投注金额是否为正数？

3.  **赔率和 ID 格式验证**: 检查所有 ID 和赔率是否为有效的字符串格式。

在我们的项目中，这些验证逻辑封装在 `internal/api/helpers.go` 中的 `validate...Request` 函数中。

### 阶段三：发射 (Fire) - 扣动扳机

这是将经过验证的注单发送到 MTS 的核心步骤。

#### 关键步骤

1.  **序列化**: 将 `TicketRequest` 对象序列化为 JSON 字符串。

2.  **建立连接**: 确保与 MTS WebSocket 的连接是活动的。

3.  **发送消息**: 通过 WebSocket 发送 JSON 消息。

4.  **设置超时**: 启动一个计时器，等待 MTS 的响应。如果超时，应将注单标记为“未知”状态并进行后续处理。

在 `mts-service` 中，此逻辑由 `internal/service/mts.go` 中的 `SendTicket` 方法处理。

```go
// 1. 序列化
msg, err := json.Marshal(ticket)

// 2. 发送
err := s.conn.WriteMessage(websocket.TextMessage, msg)

// 3. 等待响应
select {
case response := <-s.responses[correlationId]:
    // ... process response
case <-time.After(s.cfg.ResponseTimeout):
    // ... handle timeout
}
```

### 阶段四：确认 (Confirm) - 命中目标

接收并处理来自 MTS 的响应是流程的最后一步，它决定了注单的最终状态。

#### 关键步骤

1.  **接收响应**: 从 WebSocket 读取 `ticket-reply` 消息。

2.  **解析响应**: 将 JSON 响应反序列化为 `TicketResponse` 对象。

3.  **状态判断**: 检查响应中的 `status` 字段：
    - **`accepted`**: 注单被接受。这是**成功**的最终状态。
    - **`rejected`**: 注单被拒绝。需要分析 `code` 和 `message` 字段以找出原因（如赔率变化、赛事关闭、超出限额等）。

4.  **发送确认 (Acknowledgement)**: 无论注单被接受还是拒绝，都必须向 MTS 发送一个 `ticket-ack` 消息，告知 MTS 你已经收到了响应。**这是 MTS 流程中强制性的一步** [2]。

    ```go
    // 创建并发送 ticket-ack
    ack := createTicketAck(response, true) // true = acknowledged
    sendAck(ack)
    ```

5.  **更新状态**: 在本地数据库中更新注单的最终状态（`Accepted`, `Rejected`, `Timeout`）。

6.  **返回结果**: 将最终结果返回给前端用户。

## 4. 流程图

```mermaid
graph TD
    A[开始: 用户提交投注] --> B{构建注单 (Build)};
    B --> C{本地验证 (Validate)};
    C -- 验证失败 --> D[返回错误给用户];
    C -- 验证成功 --> E{发射注单 (Fire)};
    E --> F[通过 WebSocket 发送给 MTS];
    F --> G{等待 MTS 响应};
    G -- 超时 --> H[状态: 未知/超时];
    G -- 收到响应 --> I{解析响应 (Confirm)};
    I --> J{状态是 'accepted' ?};
    J -- 是 --> K[状态: 接受];
    J -- 否 --> L[状态: 拒绝];
    K --> M{发送 ACK (acknowledged=true)};
    L --> M;
    M --> N[更新本地数据库];
    N --> O[返回最终结果给用户];
    H --> N;
```

## 5. 错误处理和特殊情况

- **赔率变化**: 如果 MTS 返回的赔率与提交的不同，可以根据预设策略（`oddsChange` 字段）决定是接受新赔率还是拒绝注单。
- **网络问题**: 如果 WebSocket 连接断开，需要实现重连和状态同步逻辑。
- **重复提交**: 使用唯一的 `ticketId` 可以防止重复处理同一张注单。

## 6. 总结

“子弹装填”流程是一个系统化、标准化的方法，用于处理投注生命周期中的每一步。通过将流程分解为**构建、验证、发射、确认**四个阶段，我们可以确保系统的健壮性、可靠性和可维护性，从而为用户提供流畅、可靠的投注体验。

---

### 参考文献

[1] Sportradar MTS. "Ticket Placement Request". *Sportradar Developer Portal*. [https://docs.sportradar.com/transaction30api/api-description/ticket-json-format-description/ticket-placement-request](https://docs.sportradar.com/transaction30api/api-description/ticket-json-format-description/ticket-placement-request)

[2] Sportradar MTS. "Ticket Acceptance Flow". *Sportradar Developer Portal*. [https://docs.sportradar.com/transaction30api/api-description/ticket-acceptance-flow](https://docs.sportradar.com/transaction30api/api-description/ticket-acceptance-flow)
