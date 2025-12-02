# MTS WebSocket 交互协议与前端实现方案

**版本**: 1.0  
**日期**: 2025-12-01  
**作者**: Manus AI

---

## 1. 概述

为了解决 MTS (Managed Trading Services) 投注存在的延迟问题（通常为 2-8 秒），并为用户提供流畅、实时的投注体验，我们设计了一套基于 WebSocket 的前后端交互协议。该协议旨在取代传统的 HTTP 轮询方式，实现从投注提交到结果返回的全程异步、实时通信。

本文档详细定义了 WebSocket 的连接生命周期、核心交互流程、数据交换格式以及异常处理机制，为前后端开发提供统一的实现蓝图。

---

## 2. WebSocket 连接生命周期

```mermaid
sequenceDiagram
    participant F as 前端 Frontend
    participant W as WebSocket 服务
    participant M as MTS Service
    participant S as Sportradar MTS
    
    Note over F,S: WebSocket 连接生命周期
    
    F->>W: 建立 WebSocket 连接<br/>ws://mts-service/ws?userId=xxx&token=xxx
    W->>W: 验证 token 和用户身份
    alt 验证成功
        W-->>F: 连接成功<br/>type: connection_established
        Note over F,W: 连接已建立，开始心跳
    else 验证失败
        W-->>F: 连接拒绝<br/>type: connection_rejected
        W->>W: 关闭连接
    end
    
    loop 每 30 秒
        F->>W: 发送心跳<br/>type: ping
        W-->>F: 响应心跳<br/>type: pong
    end
    
    Note over F,W: 如果 60 秒内未收到心跳，服务端主动断开连接
    
    alt 网络异常
        W->>W: 检测到连接断开
        W->>W: 清理用户会话
        F->>F: 检测到连接断开
        F->>F: 尝试重连（指数退避）
        F->>W: 重新建立连接
    end
    
    Note over F,W: 用户主动关闭页面或登出
    F->>W: 关闭 WebSocket 连接
    W->>W: 清理用户会话
    W-->>F: 连接关闭确认
```

**图 1：WebSocket 连接生命周期**

1.  **连接建立**：
    -   前端在用户登录并进入投注页面后，向 `ws://<mts-service-host>/ws` 发起连接请求。
    -   请求参数中必须包含用户身份凭证（如 `userId` 和 `token`）以供后端验证。

2.  **心跳维持**：
    -   连接建立后，前端每 30 秒发送一次 `ping` 消息。
    -   服务端收到 `ping` 后，立即回复 `pong` 消息，以确认连接活跃。
    -   如果服务端在 60 秒内未收到任何心跳，将主动断开连接。

3.  **断线重连**：
    -   前端应监听连接断开事件。
    -   一旦断开，采用指数退避策略（例如，等待 1s, 2s, 4s, 8s...）自动尝试重连，直到连接恢复。
    -   重连成功后，前端应主动查询之前提交但未收到最终结果的投注状态。

---

## 3. 核心交互流程与数据格式

所有通过 WebSocket 的通信都使用统一的 JSON 格式。每个消息体都包含一个 `type` 字段，用于标识消息类型。

### 3.1. 投注请求 (客户端 → 服务端)

当用户点击“下单”按钮时，前端应**立即禁用该按钮**并显示加载状态，防止用户重复点击造成多个 WebSocket 请求。然后，前端发送 `place_bet` 消息。

```json
{
  "type": "place_bet",
  "requestId": "unique-client-generated-id-123",
  "betType": "single" | "multi" | "accumulator" | "system" | "banker",
  "oddsChangePreference": "accept_any" | "accept_higher" | "accept_none",
  "payload": { ... } // 对应不同投注类型的具体数据
}
```

-   `oddsChangePreference`: **(新增)** 用户的赔率变化偏好设置，用于后端处理赔率变化时的决策。

当用户点击“下单”按钮时，前端发送 `place_bet` 消息。

```json
{
  "type": "place_bet",
  "requestId": "unique-client-generated-id-123",
  "betType": "single" | "multi" | "accumulator" | "system" | "banker",
  "payload": { ... } // 对应不同投注类型的具体数据
}
```

-   `requestId`: 前端生成的唯一 ID，用于追踪整个投注生命周期。
-   `betType`: 明确告知后端本次投注的类型。
-   `payload`: 包含投注所需的所有信息（如 selections, stake 等）。

### 3.2. 投注接收确认 (服务端 → 客户端)

服务端收到并成功验证 `place_bet` 请求后，立即返回 `bet_received` 消息。

```json
{
  "type": "bet_received",
  "requestId": "unique-client-generated-id-123",
  "ticketId": "server-generated-ticket-id-456"
}
```

-   **前端操作**：收到此消息后，前端应立即将界面更新为“处理中”状态，并禁用相关的下单按钮，防止重复提交。

### 3.3. 投注最终结果 (服务端 → 客户端)

MTS 处理完成后，服务端通过 `bet_result` 消息推送最终结果。

```json
{
  "type": "bet_result",
  "requestId": "unique-client-generated-id-123",
  "ticketId": "server-generated-ticket-id-456",
  "status": "accepted" | "rejected",
  "details": { ... } // 包含赔率、派彩、拒绝原因等详细信息
}
```

-   **前端操作**：
    -   `accepted`: 显示成功提示，清空投注单，更新用户余额。
    -   `rejected`: 显示拒绝原因，高亮问题选项，保留投注单以便用户修改。

### 3.4. 批量单注的流程 (与 MTS 保持一致)

根据 MTS 规范，批量单注应通过一次 `multi` 类型的投注请求提交，而不是发送多个独立的 `single` 请求。

服务端在收到批量单注请求后，会一次性处理所有注单，然后返回一个包含所有结果的最终响应。不再推送部分结果。

```mermaid
sequenceDiagram
    participant F as 前端 Frontend
    participant W as WebSocket 服务
    participant M as MTS Service
    participant S as Sportradar MTS

    Note over F,S: 场景：用户在单关模式下批量提交多个单注

    F->>F: 用户为 3 个选项输入金额<br/>点击"下单"按钮
    F->>F: (禁用下单按钮，显示加载状态)
    F->>W: 发送批量投注请求<br/>type: place_bet<br/>requestId: xxx<br/>betType: multi<br/>payload: {bets: [bet1, bet2, bet3]}

    W->>W: 验证请求格式
    alt 验证失败
        W-->>F: 返回错误<br/>type: bet_error
        F->>F: 显示错误提示, 启用下单按钮
    else 验证成功
        W-->>F: 确认接收<br/>type: bet_received<br/>requestId: xxx
        F->>F: 界面保持"处理中..."

        W->>M: 调用 POST /api/bets/multi
        M->>S: 发送包含多个 bet 的 Ticket

        Note over M,S: MTS 一次性处理所有 bet

        S-->>M: 返回包含多个结果的 Response
        M-->>W: 推送最终结果

        W-->>F: 推送最终结果<br/>type: bet_result<br/>requestId: xxx<br/>summary: {accepted: 2, rejected: 1}<br/>details: [...]
        F->>F: 显示汇总结果<br/>清空成功的投注<br/>保留失败的投注<br/>更新余额, 启用下单按钮
    end
```

**图 2：批量单注提交流程**

1.  **部分结果推送** (`bet_partial_result`)：每当一个单注处理完成，服务端就推送一次部分结果。
2.  **最终结果推送** (`bet_result`)：所有单注处理完成后，推送一个包含汇总信息的最终结果。

---

## 4. 投注类型业务说明

| 投注类型 | 业务意义 | 场景示例 |
| :--- | :--- | :--- |
| **单注 (Single)** | 最基础的投注形式，对单个比赛的单个结果进行投注。 | 投注 A 队获胜。 |
| **串关 (Accumulator)** | 将多个不同比赛的选项串联成一个投注，所有选项必须全部猜对才能获胜，赔率相乘，风险和回报都较高。 | 同时投注 A 队胜、B 队胜、C 队平局，三场结果必须全对。 |
| **系统投注 (System)** | 同样是选择多个选项，但允许其中有部分选项猜错。系统会自动将这些选项组合成多个串关。 | 选择 4 个选项，进行 "2/4 系统投注"，系统会将其组合成 6 个不同的二串一。只要猜对其中任意 2 场，就能获得部分奖金。 |
| **Banker 系统投注** | 系统投注的变种，用户可以指定一个或多个选项为 "胆" (Banker)。这个"胆"必须猜对，然后与其他非"胆"选项组合成系统投注。 | 在 4 个选项中，指定 1 个为"胆"，然后进行 "2/3 系统投注"。如果"胆"猜对了，另外 3 个选项中只要猜对 2 个，就能获胜。 |
| **批量单注 (Multi)** | 在单关模式下，一次性提交多个独立的单注。每个单注的输赢独立计算。 | 同时提交 A 队胜、B 队胜、C 队胜三个独立的单注，每个投注金额 10 元。 |

---

## 5. 通用注单 (Ticket) 格式与字段说明

所有投注请求最终都会被 `mts-service` 构造成符合 Sportradar MTS 规范的 Ticket。以下是通用的 Ticket 结构和字段说明。

### 5.1. 通用 Ticket 字段说明

| 字段路径 | 类型 | 是否必须 | 描述 |
| :--- | :--- | :--- | :--- |
| `ticket.bets[].selections[].id` | string | 是 | 投注选项的唯一标识符，例如 `sr:match:12345|3`。 |
| `ticket.bets[].selections[].odds` | string | 是 | 投注时的赔率，**必须为字符串格式**。 |
| `ticket.bets[].stake.amount` | string | 是 | 投注金额，**必须为字符串格式**。 |
| `ticket.bets[].stake.type` | string | 是 | 金额类型，通常为 `per_bet` (每个投注) 或 `total` (总计)。 |
| `ticket.bets[].type` | string | 是 | 投注类型，如 `single`, `accumulator`, `system`。 |
| `ticket.bets[].system.size` | array | 否 | 系统投注的组合方式，例如 `[2, 3]` 表示二串一和三串一。 |
| `ticket.bets[].system.type` | string | 否 | 系统投注的预设名称，例如 `yankee`。 |

### 5.2. 各类型注单格式示例

#### 单注 (Single)

```json
// POST /api/bets/single
{
  "type": "single",
  "selections": [
    {
      "id": "sr:match:38102865|3",
      "odds": "1.85"
    }
  ],
  "stake": {
    "amount": "10.00",
    "type": "per_bet"
  }
}
```

#### 串关 (Accumulator)

```json
// POST /api/bets/accumulator
{
  "type": "accumulator",
  "selections": [
    {
      "id": "sr:match:38102865|3",
      "odds": "1.85"
    },
    {
      "id": "sr:match:38102867|1",
      "odds": "2.10"
    }
  ],
  "stake": {
    "amount": "10.00",
    "type": "total"
  }
}
```

#### 系统投注 (System)

```json
// POST /api/bets/system
{
  "type": "system",
  "selections": [
    { "id": "sr:match:38102865|3", "odds": "1.85" },
    { "id": "sr:match:38102867|1", "odds": "2.10" },
    { "id": "sr:match:38102869|2", "odds": "3.50" }
  ],
  "stake": {
    "amount": "5.00",
    "type": "per_bet"
  },
  "system": {
    "size": [2] // 表示 2/3 系统投注
  }
}
```

#### Banker 系统投注

```json
// POST /api/bets/banker-system
{
  "type": "banker-system",
  "bankers": [
    { "id": "sr:match:38102865|3", "odds": "1.85" }
  ],
  "selections": [
    { "id": "sr:match:38102867|1", "odds": "2.10" },
    { "id": "sr:match:38102869|2", "odds": "3.50" },
    { "id": "sr:match:38102871|1", "odds": "1.50" }
  ],
  "stake": {
    "amount": "5.00",
    "type": "per_bet"
  },
  "system": {
    "size": [2] // 表示 1 个 banker + 2/3 系统投注
  }
}
```

#### 批量单注 (Multi)

```json
// POST /api/bets/multi
{
  "bets": [
    {
      "type": "single",
      "selections": [{ "id": "sr:match:38102865|3", "odds": "1.85" }],
      "stake": { "amount": "10.00", "type": "per_bet" }
    },
    {
      "type": "single",
      "selections": [{ "id": "sr:match:38102867|1", "odds": "2.10" }],
      "stake": { "amount": "20.00", "type": "per_bet" }
    }
  ]
}
```

---

## 6. 投注场景与交互序列图

### 4.1. 单注投注

```mermaid
sequenceDiagram
    participant F as 前端 Frontend
    participant W as WebSocket 服务
    participant M as MTS Service
    participant S as Sportradar MTS
    
    Note over F,S: 场景：用户提交单注投注
    
    F->>F: 用户点击"下单"按钮
    F->>F: 生成唯一 requestId
    F->>W: 发送投注请求<br/>type: place_bet<br/>requestId: xxx<br/>betType: single<br/>payload: {...}
    
    W->>W: 验证请求格式和用户权限
    alt 验证失败
        W-->>F: 返回错误<br/>type: bet_error<br/>requestId: xxx<br/>error: {...}
        F->>F: 显示错误提示
    else 验证成功
        W-->>F: 确认接收<br/>type: bet_received<br/>requestId: xxx<br/>ticketId: yyy
        F->>F: 显示"投注处理中..."<br/>禁用下单按钮
        
        W->>M: 调用 POST /api/bets/single
        M->>M: 构建 MTS Ticket
        M->>S: 发送 Ticket 到 Sportradar
        
        Note over M,S: MTS 处理延迟（通常 2-5 秒）
        
        S-->>M: 返回 Ticket Response
        M->>M: 解析响应结果
        
        alt MTS 接受投注
            M-->>W: 投注成功<br/>ticketId: yyy<br/>status: accepted
            W-->>F: 推送结果<br/>type: bet_result<br/>requestId: xxx<br/>ticketId: yyy<br/>status: accepted<br/>details: {...}
            F->>F: 显示成功提示<br/>清空投注单<br/>更新余额
        else MTS 拒绝投注
            M-->>W: 投注被拒<br/>ticketId: yyy<br/>status: rejected<br/>reason: {...}
            W-->>F: 推送结果<br/>type: bet_result<br/>requestId: xxx<br/>ticketId: yyy<br/>status: rejected<br/>reason: {...}
            F->>F: 显示拒绝原因<br/>保留投注单<br/>启用下单按钮
        end
    end
```

**图 3：单注投注交互序列图**

### 4.2. 串关投注 (Accumulator/System/Banker)

```mermaid
sequenceDiagram
    participant F as 前端 Frontend
    participant W as WebSocket 服务
    participant M as MTS Service
    participant S as Sportradar MTS
    
    Note over F,S: 场景：用户提交串关投注（Accumulator/System/Banker）
    
    F->>F: 用户在串关模式下点击"下单"
    F->>F: 生成唯一 requestId
    F->>W: 发送投注请求<br/>type: place_bet<br/>requestId: xxx<br/>betType: accumulator/system/banker<br/>payload: {...}
    
    W->>W: 验证请求格式和选项有效性
    alt 验证失败
        W-->>F: 返回错误<br/>type: bet_error<br/>requestId: xxx<br/>error: {...}
        F->>F: 显示错误提示
    else 验证成功
        W-->>F: 确认接收<br/>type: bet_received<br/>requestId: xxx<br/>ticketId: yyy
        F->>F: 显示"投注处理中..."<br/>禁用下单按钮<br/>显示加载动画
        
        W->>M: 调用对应 API<br/>POST /api/bets/accumulator<br/>POST /api/bets/system<br/>POST /api/bets/banker-system
        M->>M: 构建复杂 MTS Ticket<br/>处理多个 selections
        M->>S: 发送 Ticket 到 Sportradar
        
        Note over M,S: MTS 处理延迟（串关通常 3-8 秒）
        
        S-->>M: 返回 Ticket Response
        M->>M: 解析响应结果<br/>提取详细信息
        
        alt MTS 接受投注
            M-->>W: 投注成功<br/>ticketId: yyy<br/>status: accepted<br/>totalOdds: xxx<br/>potentialPayout: xxx
            W-->>F: 推送结果<br/>type: bet_result<br/>requestId: xxx<br/>ticketId: yyy<br/>status: accepted<br/>details: {...}
            F->>F: 显示成功提示<br/>展示投注详情<br/>清空投注单<br/>更新余额
        else MTS 拒绝投注
            M-->>W: 投注被拒<br/>ticketId: yyy<br/>status: rejected<br/>reason: {...}
            W-->>F: 推送结果<br/>type: bet_result<br/>requestId: xxx<br/>ticketId: yyy<br/>status: rejected<br/>reason: {...}
            F->>F: 显示拒绝原因<br/>高亮问题选项<br/>保留投注单<br/>启用下单按钮
        end
    end
```

**图 4：串关投注交互序列图**

---

## 5. 异常处理机制

```mermaid
sequenceDiagram
    participant F as 前端 Frontend
    participant W as WebSocket 服务
    participant M as MTS Service
    participant S as Sportradar MTS
    
    Note over F,S: 异常场景 1: 投注提交后赔率变化
    
    F->>W: 发送投注请求<br/>selection odds: 2.50
    W-->>F: 确认接收<br/>ticketId: yyy
    W->>M: 调用 API
    M->>S: 发送 Ticket
    
    alt Sportradar 检测到赔率变化
        S-->>M: 返回 rejected<br/>reason: odds_changed<br/>new_odds: 2.30
        M-->>W: 投注被拒<br/>原因: 赔率变化
        W-->>F: 推送结果<br/>type: bet_result<br/>status: rejected<br/>reason: odds_changed<br/>old_odds: 2.50<br/>new_odds: 2.30
        
        F->>F: 检查用户设置
        alt 用户设置: 接受任何赔率变化
            F->>F: 自动更新赔率为 2.30
            F->>W: 重新提交投注<br/>新 requestId
        else 用户设置: 不接受赔率变化
            F->>F: 高亮显示赔率变化<br/>显示确认对话框
            F->>F: 等待用户手动确认
        end
    end
    
    Note over F,S: 异常场景 2: 投注提交超时
    
    F->>W: 发送投注请求
    W-->>F: 确认接收
    W->>M: 调用 API
    M->>S: 发送 Ticket
    
    Note over M,S: 网络延迟或 MTS 响应慢
    
    alt 超过 15 秒未收到响应
        M->>M: 检测超时
        M-->>W: 返回超时状态<br/>status: timeout
        W-->>F: 推送超时通知<br/>type: bet_timeout<br/>requestId: xxx<br/>ticketId: yyy
        F->>F: 显示提示: 投注处理中，请稍后查看投注记录
        
        Note over M,S: 后台继续等待 MTS 响应
        
        S-->>M: 延迟返回结果
        M-->>W: 推送延迟结果
        W-->>F: 推送结果<br/>type: bet_result_delayed<br/>ticketId: yyy<br/>status: accepted/rejected
        F->>F: 显示通知: 您的投注已处理
    end
    
    Note over F,S: 异常场景 3: WebSocket 连接断开
    
    F->>W: 发送投注请求
    W-->>F: 确认接收<br/>ticketId: yyy
    W->>M: 调用 API
    
    Note over F,W: WebSocket 连接意外断开
    
    F->>F: 检测到连接断开
    F->>F: 尝试重连
    F->>W: 重新建立连接
    W-->>F: 连接成功
    
    F->>W: 查询投注状态<br/>type: query_bet_status<br/>ticketId: yyy
    W->>M: 查询 Ticket 状态
    M-->>W: 返回状态
    W-->>F: 推送状态<br/>type: bet_status<br/>ticketId: yyy<br/>status: accepted/rejected/pending
    F->>F: 根据状态更新 UI
```

**图 5：异常处理流程图**

-   **赔率变化**：服务端在 `bet_result` 的 `details` 中返回新旧赔率，前端根据用户设置决定是自动重提还是弹窗确认。
-   **投注超时**：如果服务端在预设时间（例如 15 秒）内未收到 MTS 的响应，会先返回 `bet_timeout` 消息。前端提示用户稍后查看结果，后台继续等待。最终结果通过 `bet_result_delayed` 推送。
-   **连接断开**：前端重连成功后，应主动发送 `query_bet_status` 消息，查询之前提交但未收到结果的 `ticketId` 的状态。

---

## 6. 结论

该 WebSocket 协议通过实时、双向的通信机制，优雅地解决了 MTS 投注延迟带来的用户体验问题。它不仅提供了清晰的交互流程和数据格式，还为各种异常情况设计了健壮的处理方案。我们强烈建议前后端开发团队以此文档为核心依据，协同完成开发，从而为用户打造一个真正流畅、可靠的实时投注平台。

---

## 7. 附录：Go 实现参考

为了便于后端开发，以下提供了本次 WebSocket 改造的核心 Go 代码文件。

### 7.1. 消息定义 (`internal/websocket/messages.go`)

```go
// ... (内容见附件)
```

### 7.2. 客户端连接管理 (`internal/websocket/client.go`)

```go
// ... (内容见附件)
```

### 7.3. 连接池管理 (`internal/websocket/hub.go`)

```go
// ... (内容见附件)
```

### 7.4. 投注处理器 (`internal/websocket/bet_processor.go`)

```go
// ... (内容见附件)
```

### 7.5. HTTP 处理器 (`internal/websocket/handler.go`)

```go
// ... (内容见附件)
```

### 7.6. 主程序集成 (`cmd/server/mts_main.go`)

```go
// ... (内容见附件)
```
