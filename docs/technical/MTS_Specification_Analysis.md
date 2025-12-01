# MTS 消息规范及错误分析

## 错误信息回顾

用户提供的错误信息如下：

> `2025/11/17 02:26:53 MTS Error Reply received (CorrelationID: init-1763346071798579366): Code=401, Message=Schema validation failed: $.version: must be the constant value '3.0';$.content.ticketId: must be at least 1 characters long;$.content.ticketSignature: must be at least 1 characters long`

## MTS 消息结构概述

根据 Sportradar 的 MTS 3.0 API 文档，MTS 消息（包括请求和回复）通常包含一个**信封（Envelope）**和一个**内容（Content）**对象。

| 字段名 | 位置 | 类型 | 描述 | 规范要求 |
| :--- | :--- | :--- | :--- | :--- |
| `version` | 信封（Envelope） | String | 指示票据格式版本 | 必须是常量值 `"3.0"` |
| `correlationId` | 信封（Envelope） | String | 客户端定义的字符串，用于请求-响应配对 | 最小长度 >= 1 |
| `timestampUtc` | 信封（Envelope） | Integer | 票据放置时间戳（Unix 毫秒） | |
| `operation` | 信封（Envelope） | String | 操作类型（例如：`ticket-placement`） | |
| `content` | 信封（Envelope） | Object | 包含实际消息内容的 JSON 对象 | |
| `ticketId` | 内容（Content） | String | 原始票据的 ID | 最小长度 >= 1 |
| `ticketSignature` | 内容（Content） | String | 签名（仅在回复消息中出现） | 最小长度 >= 1 |

## 错误分析与规范解读

错误信息 `Schema validation failed` 表明客户端发送的消息（或 MTS 返回的回复消息）不符合 MTS 3.0 的 JSON 模式定义。

| 错误字段 | 错误描述 | 规范要求 | 客户端/服务端的潜在问题 |
| :--- | :--- | :--- | :--- |
| `$.version` | `must be the constant value '3.0'` | **信封**中的 `version` 字段必须是字符串 `"3.0"`。 | 客户端发送的请求中，`version` 字段可能缺失、使用了错误的类型（如数字 `3.0`）或使用了错误的字符串值（如 `"3.1"`）。 |
| `$.content.ticketId` | `must be at least 1 characters long` | **内容**对象中的 `ticketId` 字段必须是一个长度至少为 1 个字符的字符串。 | 客户端发送的请求中，`content.ticketId` 字段可能缺失或为空字符串 `""`。 |
| `$.content.ticketSignature` | `must be at least 1 characters long` | **内容**对象中的 `ticketSignature` 字段必须是一个长度至少为 1 个字符的字符串。 | `ticketSignature` 是 MTS 在**回复消息**中发送的签名。如果这是客户端收到的 **MTS Error Reply**，则说明 MTS 在构建这个错误回复时，可能在 `content` 对象中包含了 `ticketSignature` 字段，但其值不符合最小长度要求。这通常发生在 MTS 内部处理回复时。**更可能的情况是，客户端正在解析一个 MTS 回复，而该回复的 `content` 结构不符合预期的模式。** |

## 关键规范链接

以下是 Sportradar 文档中与 MTS 消息结构相关的关键部分：

1.  **错误回复规范 (`Error-reply Response`)**:
    *   **链接**: [https://docs.sportradar.com/transaction30api/api-description/ticket-json-format-description/error-reply-response](https://docs.sportradar.com/transaction30api/api-description/ticket-json-format-description/error-reply-response)
    *   **内容**: 描述了当请求被 Data Plane Router 拒绝或发生内部错误时，MTS 返回的 `error-reply` 消息结构。其 `content` 对象只包含 `type: "error-reply"`, `code`, 和 `message`。

2.  **回复消息信封 (`Acknowledgement Reply Messages`)**:
    *   **链接**: [https://docs.sportradar.com/transaction30api/api-description/ticket-json-format-description/acknowledgement-reply-messages](https://docs.sportradar.com/transaction30api/api-description/ticket-json-format-description/acknowledgement-reply-messages)
    *   **内容**: 描述了 MTS 回复消息的信封结构，其中明确提到了 `version` 字段必须是 `"3.0"`。

## 建议排查方向

根据错误信息，建议检查以下几点：

1.  **检查请求的 `version` 字段**: 确保客户端发送给 MTS 的请求中，顶层的 `version` 字段是**字符串** `"3.0"`。
2.  **检查请求的 `content.ticketId` 字段**: 确保客户端发送的请求中，`content` 对象内的 `ticketId` 字段存在且非空。
3.  **检查客户端对 `MTS Error Reply` 的解析逻辑**:
    *   `MTS Error Reply` 的 `content` 结构应该只包含 `type: "error-reply"`, `code`, 和 `message`。
    *   如果客户端的解析器期望在所有回复的 `content` 中都找到 `ticketId` 和 `ticketSignature`，那么当收到 `error-reply` 时就会因为字段缺失或结构不匹配而失败。
    *   请参考 [Error-reply Response](https://docs.sportradar.com/transaction30api/api-description/ticket-json-format-description/error-reply-response) 链接，确认 `error-reply` 的结构，并相应调整客户端的解析逻辑。

---
*此分析基于 Sportradar MTS 3.0 API 文档，并结合用户提供的错误信息得出。*
