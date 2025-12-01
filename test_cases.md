# MTS Service API 测试用例

## 1. 概述

本文档定义了 `mts-service` 的 API 测试用例，旨在验证其与 Sportradar MTS 的集成是否符合规范，并确保所有投注类型的正确处理。测试将覆盖以下几个方面：

- **API 端点功能**：验证每个 API 端点的请求和响应是否正确。
- **投注类型支持**：确保所有支持的投注类型（单式、复式、系统等）都能被正确处理。
- **数据验证**：测试无效或格式错误的输入是否能被正确拒绝。
- **异步流程**：通过提交 ticket 并轮询结果来验证整个异步处理流程。

## 2. 测试环境

- **API Base URL**: `http://localhost:8080`
- **MTS 环境**: Sportradar 测试环境
- **测试数据**: 使用 Sportradar 提供的测试赛事 ID。如果找不到，将使用占位符 ID，例如 `sr:match:12345`。

## 3. 测试用例

### 3.1 健康检查 (`GET /health`)

| 用例 ID | 描述 | 预期结果 |
| :--- | :--- | :--- |
| TC-HEALTH-001 | 发送 GET 请求到 `/health` | 返回 `200 OK`，响应体包含 `{"status": "healthy", ...}`。 |

### 3.2 单注 (`POST /api/bets/single`)

| 用例 ID | 描述 | 请求体 | 预期结果 |
| :--- | :--- | :--- | :--- |
| TC-SINGLE-001 | **成功**：提交一个有效的单注 | 合法的 `ticketId`, `selection`, `stake` | 返回 `200 OK`，`success: true`，`data.content.status` 为 `accepted`。 |
| TC-SINGLE-002 | **失败**：缺少 `ticketId` | 请求体中缺少 `ticketId` | 返回 `400 Bad Request`，`success: false`，并有明确的错误信息。 |
| TC-SINGLE-003 | **失败**：无效的赔率 | `selection.odds` 为负数或零 | 返回 `400 Bad Request`，`success: false`，并有明确的错误信息。 |
| TC-SINGLE-004 | **失败**：无效的赛事 ID | `selection.eventId` 格式不正确 | 返回 `200 OK`，但 `data.content.status` 为 `rejected`，并包含 MTS 的错误代码。 |

### 3.3 串关 (`POST /api/bets/accumulator`)

| 用例 ID | 描述 | 请求体 | 预期结果 |
| :--- | :--- | :--- | :--- |
| TC-ACC-001 | **成功**：提交一个有效的两场串关 | 包含 2 个 `selections` 的合法请求 | 返回 `200 OK`，`success: true`，`data.content.status` 为 `accepted`。 |
| TC-ACC-002 | **失败**：少于 2 个选项 | `selections` 数组只包含 1 个元素 | 返回 `400 Bad Request`，`success: false`，并有明确的错误信息。 |
| TC-ACC-003 | **失败**：`stake.mode` 不为 `total` | `stake.mode` 设置为 `unit` | 返回 `400 Bad Request`，`success: false`，并有明确的错误信息。 |

### 3.4 系统投注 (`POST /api/bets/system`)

| 用例 ID | 描述 | 请求体 | 预期结果 |
| :--- | :--- | :--- | :--- |
| TC-SYS-001 | **成功**：提交一个有效的 2/3 系统投注 | 3 个 `selections`，`size: [2]`，`stake.mode: "unit"` | 返回 `200 OK`，`success: true`，`data.content.status` 为 `accepted`。 |
| TC-SYS-002 | **失败**：`stake.mode` 不为 `unit` | `stake.mode` 设置为 `total` | 返回 `400 Bad Request`，`success: false`，并有明确的错误信息。 |
| TC-SYS-003 | **失败**：`size` 无效 | `size` 中的值大于 `selections` 的数量 | 返回 `400 Bad Request`，`success: false`，并有明确的错误信息。 |

### 3.5 Banker 系统投注 (`POST /api/bets/banker-system`)

| 用例 ID | 描述 | 请求体 | 预期结果 |
| :--- | :--- | :--- | :--- |
| TC-BANKER-001 | **成功**：提交一个有效的 1 个 Banker + 2/3 系统投注 | 1 个 `banker`，3 个 `selections`，`size: [2]` | 返回 `200 OK`，`success: true`，`data.content.status` 为 `accepted`。 |
| TC-BANKER-002 | **失败**：缺少 `bankers` | `bankers` 数组为空 | 返回 `400 Bad Request`，`success: false`，并有明确的错误信息。 |

### 3.6 预设系统投注 (`POST /api/bets/preset`)

| 用例 ID | 描述 | 请求体 | 预期结果 |
| :--- | :--- | :--- | :--- |
| TC-PRESET-001 | **成功**：提交一个有效的 Trixie 投注 | `type: "trixie"`，3 个 `selections` | 返回 `200 OK`，`success: true`，`data.content.status` 为 `accepted`。 |
| TC-PRESET-002 | **失败**：选项数量与类型不匹配 | `type: "trixie"`，但只有 2 个 `selections` | 返回 `400 Bad Request`，`success: false`，并有明确的错误信息。 |
| TC-PRESET-003 | **失败**：无效的预设类型 | `type` 设置为 `"invalid_type"` | 返回 `400 Bad Request`，`success: false`，并有明确的错误信息。 |

### 3.7 多重彩 (`POST /api/bets/multi`)

| 用例 ID | 描述 | 请求体 | 预期结果 |
| :--- | :--- | :--- | :--- |
| TC-MULTI-001 | **成功**：提交一个包含单注和串关的多重彩 | `bets` 数组包含一个 `single` 和一个 `accumulator` | 返回 `200 OK`，`success: true`，`data.content.status` 为 `accepted`。 |
| TC-MULTI-002 | **失败**：`bets` 数组为空 | `bets` 数组为空 | 返回 `400 Bad Request`，`success: false`，并有明确的错误信息。 |

### 3.8 异步结果验证

| 用例 ID | 描述 | 步骤 | 预期结果 |
| :--- | :--- | :--- | :--- |
| TC-ASYNC-001 | **验证**：成功提交的 ticket 能在后续查询中获取 | 1. 成功提交一个 ticket。<br>2. 保存返回的 `ticketId`。<br>3. 开发一个独立的脚本或函数，使用 `ticketId` 查询 MTS 状态（需要模拟或找到查询接口）。 | 查询结果应与提交时返回的 `status` 一致。 |

## 4. 异步验证策略

由于 `mts-service` 本身没有提供查询 ticket 状态的接口，我们将采取以下策略来验证异步流程：

1.  **日志监控**：在执行测试时，监控 `mts-service` 的日志输出。成功的 ticket 提交和 MTS 的 `ticket-ack` 消息应该会被记录下来。
2.  **模拟查询服务**：如果条件允许，可以创建一个模拟的 MTS 查询服务，该服务可以根据 `ticketId` 返回预设的状态。但这超出了本次测试的范围。
3.  **依赖 MTS 回调**：在实际生产环境中，MTS 会通过回调将最终的结算结果发送给指定的端点。在本次测试中，我们将重点关注 `ticket-reply` 和 `ticket-ack` 的即时响应。

基于以上分析，我们将主要通过 API 的即时响应和日志来验证测试结果。
