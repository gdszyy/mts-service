# MTS Service 整改报告

## 概述

本次整改旨在将项目 `mts-service` 与您提供的 **MTS AWS INTEGRATION** 环境凭证和配置对接，解决 WebSocket 连接成功后因缺少 `operatorId` 导致的初始化消息失败问题。

## 整改内容

### 1. 配置结构更新 (`internal/config/config.go`)

在配置结构中新增了以下字段，用于加载 MTS AWS INTEGRATION 环境所需的凭证和 URL：

| 字段 | 环境变量 | 描述 |
| :--- | :--- | :--- |
| `LimitID` | `MTS_LIMIT_ID` | 交易消息体所需的 Limit ID (`4268`) |
| `OperatorID` | `MTS_OPERATOR_ID` | WebSocket 初始化消息所需的 Operator ID (`45426`) |
| `WSURL` | `MTS_WS_URL` | WebSocket 连接地址 (`wss://wss.dataplane-nonprod.sportradar.dev`) |
| `WSAudience` | `MTS_WS_AUDIENCE` | OAuth2 Audience (`mbs-dp-non-prod-wss`) |

### 2. MTS 服务逻辑更新 (`internal/service/mts.go`)

*   **OAuth2 认证：** 移除了硬编码的 `IntegrationAudience`，改为从配置中加载 `WSAudience`，确保使用正确的 OAuth2 Audience。
*   **WebSocket 连接：** 移除了硬编码的 `IntegrationWSURL`，改为从配置中加载 `WSURL`，确保连接到正确的 WebSocket 端点。
*   **初始化消息发送：** 新增了 `sendInitializationMessage()` 方法，在 WebSocket 连接成功后立即发送初始化订阅消息，解决了缺少 `operatorId` 的问题。

    *   **初始化消息结构 (JSON):**
        ```json
        {
            "type": "subscribe",
            "bookmaker_id": "45426", // 从 MTS_BOOKMAKER_ID 加载
            "limit_id": "4268",      // 从 MTS_LIMIT_ID 加载
            "operatorId": "45426",   // 从 MTS_OPERATOR_ID 加载
            "correlationId": "...",
            "timestampUtc": "...",
            "operation": "initialization", // 假设的操作名
            "version": "3.0"
        }
        ```
    *   **注意：** `operation` 字段被假设为 `"initialization"`，`version` 被假设为 `"3.0"`。如果 Sportradar MTS 服务对初始化消息有更精确的要求，可能需要微调此结构。

## 测试指导

为了测试整改后的服务，您需要设置以下环境变量，并在 `mts-service` 目录下运行项目。

### 1. 设置环境变量

请使用您提供的凭证设置以下环境变量：

```bash
# 核心 MTS 凭证
export MTS_CLIENT_ID="rOCi1mVU6UhH7I2lu3t1ME4dyfkiPFVk"
export MTS_CLIENT_SECRET="IWGVLxyH9XcoJpZNEDNjCSFSvDAp49c_7kOi_iCuxQzitOsMfY8X4HMmw3Dcydcr"

# MTS AWS INTEGRATION 环境配置
export MTS_BOOKMAKER_ID="45426"
export MTS_LIMIT_ID="4268"
export MTS_OPERATOR_ID="45426" # 经确认，使用 Bookmaker ID 作为 Operator ID

# WebSocket 连接信息
export MTS_WS_URL="wss://wss.dataplane-nonprod.sportradar.dev"
export MTS_WS_AUDIENCE="mbs-dp-non-prod-wss"

# 可选：如果需要运行在非生产模式（默认就是非生产）
export MTS_PRODUCTION="false" 
```

### 2. 运行项目

在 `mts-service` 目录下执行：

```bash
go run main.go
```

### 3. 预期结果

*   服务启动日志将显示配置加载成功。
*   服务将尝试获取 OAuth2 Token。
*   服务将成功连接到 `wss://wss.dataplane-nonprod.sportradar.dev`。
*   服务将发送初始化消息，日志中会打印出消息内容，例如：
    `Sending initialization message: map[bookmaker_id:45426 correlationId:init-xxxxxx limit_id:4268 operatorId:45426 operation:initialization timestampUtc:xxxxxx type:subscribe version:3.0]`
*   如果初始化消息被 MTS 服务接受，服务将进入正常运行状态，等待接收和处理交易消息。

如果服务在发送初始化消息后仍然失败，请检查日志中是否有 MTS 服务返回的错误信息，并提供给我，以便进一步调试初始化消息的结构。
