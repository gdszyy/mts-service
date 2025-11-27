# 子弹装填流程 - 代码示例

本文档提供了在 `mts-service` 项目中实现子弹装填流程的完整代码示例。

---

## 示例 1: 单注 (Single Bet)

### 阶段一：构建 (Build)

```go
package main

import (
    "github.com/gdsZyy/mts-service/internal/models"
    "github.com/gdsZyy/mts-service/internal/config"
)

func buildSingleBet() (*models.TicketRequest, error) {
    // 1. 收集输入
    selection := models.Selection{
        Type:      "uf",
        ProductID: "3",
        EventID:   "sr:match:12345",
        MarketID:  "1",
        OutcomeID: "1",
        Odds: &models.Odds{
            Type:  "decimal",
            Value: "2.50",
        },
    }
    
    stake := models.Stake{
        Type:     "cash",
        Currency: "EUR",
        Amount:   "10.00",
        Mode:     "total",
    }
    
    context := &models.Context{
        Channel: &models.Channel{
            Type: "internet",
            Lang: "EN",
        },
        LimitID: 4268,
    }
    
    // 2. 选择构建器
    cfg := config.LoadConfig()
    builder := models.NewTicketBuilder(cfg.OperatorID)
    
    // 3. 生成注单结构
    ticket, err := builder.NewSingleBetTicket(
        "single-001",
        selection,
        stake,
        context,
    )
    
    return ticket, err
}
```

### 阶段二：验证 (Validate)

```go
func validateSingleBet(ticket *models.TicketRequest) error {
    // 1. 结构验证
    if ticket.TicketID == "" {
        return fmt.Errorf("ticketId is required")
    }
    
    if len(ticket.Content.Bets) == 0 {
        return fmt.Errorf("at least one bet is required")
    }
    
    bet := ticket.Content.Bets[0]
    if len(bet.Selections) == 0 {
        return fmt.Errorf("at least one selection is required")
    }
    
    // 2. 逻辑验证
    if len(bet.Stake) == 0 {
        return fmt.Errorf("stake is required")
    }
    
    stake := bet.Stake[0]
    amount, err := strconv.ParseFloat(stake.Amount, 64)
    if err != nil || amount <= 0 {
        return fmt.Errorf("stake amount must be positive")
    }
    
    // 3. 格式验证
    selection := bet.Selections[0]
    if selection.EventID == "" || selection.MarketID == "" {
        return fmt.Errorf("eventId and marketId are required")
    }
    
    return nil
}
```

### 阶段三：发射 (Fire)

```go
func fireBet(mtsService *service.MTSService, ticket *models.TicketRequest) (*models.TicketResponse, error) {
    // 1. 序列化 (在 SendTicket 内部完成)
    // 2. 发送消息
    // 3. 等待响应
    response, err := mtsService.SendTicket(ticket)
    if err != nil {
        return nil, fmt.Errorf("failed to send ticket: %w", err)
    }
    
    return response, nil
}
```

### 阶段四：确认 (Confirm)

```go
func confirmBet(response *models.TicketResponse) (string, error) {
    // 1-2. 接收和解析响应 (已由 MTSService 完成)
    
    // 3. 状态判断
    if response.Content.Status == "accepted" {
        log.Printf("✓ Ticket %s accepted", response.Content.TicketID)
        return "accepted", nil
    } else {
        log.Printf("✗ Ticket %s rejected: [%d] %s", 
            response.Content.TicketID,
            response.Content.Code,
            response.Content.Message)
        return "rejected", nil
    }
    
    // 4. 发送 ACK (由 MTSService 自动处理)
    // 5. 更新状态 (应在此处实现数据库更新)
    // 6. 返回结果
}
```

### 完整流程

```go
func processSingleBet() {
    // 阶段一：构建
    ticket, err := buildSingleBet()
    if err != nil {
        log.Printf("Build failed: %v", err)
        return
    }
    
    // 阶段二：验证
    if err := validateSingleBet(ticket); err != nil {
        log.Printf("Validation failed: %v", err)
        return
    }
    
    // 阶段三：发射
    mtsService := getMTSService() // 获取 MTS 服务实例
    response, err := fireBet(mtsService, ticket)
    if err != nil {
        log.Printf("Fire failed: %v", err)
        return
    }
    
    // 阶段四：确认
    status, err := confirmBet(response)
    if err != nil {
        log.Printf("Confirm failed: %v", err)
        return
    }
    
    log.Printf("Final status: %s", status)
}
```

---

## 示例 2: Trixie 系统串 (Preset System Bet)

### 阶段一：构建 (Build)

```go
func buildTrixieBet() (*models.TicketRequest, error) {
    // 1. 收集输入 - 3 个 selections
    selections := []models.Selection{
        {
            Type:      "uf",
            ProductID: "3",
            EventID:   "sr:match:11111",
            MarketID:  "1",
            OutcomeID: "1",
            Odds:      &models.Odds{Type: "decimal", Value: "1.50"},
        },
        {
            Type:      "uf",
            ProductID: "3",
            EventID:   "sr:match:22222",
            MarketID:  "1",
            OutcomeID: "2",
            Odds:      &models.Odds{Type: "decimal", Value: "2.00"},
        },
        {
            Type:      "uf",
            ProductID: "3",
            EventID:   "sr:match:33333",
            MarketID:  "1",
            OutcomeID: "1",
            Odds:      &models.Odds{Type: "decimal", Value: "1.80"},
        },
    }
    
    stake := models.Stake{
        Type:     "cash",
        Currency: "EUR",
        Amount:   "5.00",  // 每个组合 5 EUR
        Mode:     "unit",  // unit 模式
    }
    
    context := &models.Context{
        Channel: &models.Channel{Type: "internet", Lang: "EN"},
        LimitID: 4268,
    }
    
    // 2. 选择构建器
    cfg := config.LoadConfig()
    builder := models.NewTicketBuilder(cfg.OperatorID)
    
    // 3. 生成注单结构
    ticket, err := builder.NewPresetSystemBetTicket(
        "trixie-001",
        models.Trixie,  // 3 个 doubles + 1 个 treble = 4 个组合
        selections,
        stake,
        context,
    )
    
    return ticket, err
}
```

### 阶段二：验证 (Validate)

```go
func validateTrixieBet(ticket *models.TicketRequest) error {
    // 1. 结构验证
    if len(ticket.Content.Bets) == 0 {
        return fmt.Errorf("no bets found")
    }
    
    bet := ticket.Content.Bets[0]
    if len(bet.Selections) == 0 {
        return fmt.Errorf("no selections found")
    }
    
    // 2. 逻辑验证 - Trixie 特定检查
    systemSelection := bet.Selections[0]
    if systemSelection.Type != "system" {
        return fmt.Errorf("expected system selection")
    }
    
    // Trixie 需要正好 3 个 selections
    if len(systemSelection.Selections) != 3 {
        return fmt.Errorf("trixie requires exactly 3 selections, got %d", 
            len(systemSelection.Selections))
    }
    
    // Trixie 的 size 应该是 [2, 3]
    expectedSize := []int{2, 3}
    if !reflect.DeepEqual(systemSelection.Size, expectedSize) {
        return fmt.Errorf("trixie size should be [2, 3], got %v", 
            systemSelection.Size)
    }
    
    // 3. Stake 模式检查
    if bet.Stake[0].Mode != "unit" {
        log.Printf("Warning: Trixie typically uses 'unit' stake mode")
    }
    
    return nil
}
```

### 完整流程（与示例 1 类似）

```go
func processTrixieBet() {
    ticket, err := buildTrixieBet()
    if err != nil {
        log.Printf("Build failed: %v", err)
        return
    }
    
    if err := validateTrixieBet(ticket); err != nil {
        log.Printf("Validation failed: %v", err)
        return
    }
    
    mtsService := getMTSService()
    response, err := fireBet(mtsService, ticket)
    if err != nil {
        log.Printf("Fire failed: %v", err)
        return
    }
    
    status, err := confirmBet(response)
    log.Printf("Trixie bet final status: %s", status)
}
```

---

## 示例 3: Banker 系统串

### 阶段一：构建 (Build)

```go
func buildBankerBet() (*models.TicketRequest, error) {
    // 1. 收集输入
    // Banker: 1 个必中项
    bankerSelections := []models.Selection{
        {
            Type:      "uf",
            ProductID: "3",
            EventID:   "sr:match:99999",
            MarketID:  "1",
            OutcomeID: "1",
            Odds:      &models.Odds{Type: "decimal", Value: "1.20"},
        },
    }
    
    // 普通 selections: 3 个选项
    regularSelections := []models.Selection{
        {
            Type:      "uf",
            ProductID: "3",
            EventID:   "sr:match:11111",
            MarketID:  "1",
            OutcomeID: "2",
            Odds:      &models.Odds{Type: "decimal", Value: "1.80"},
        },
        {
            Type:      "uf",
            ProductID: "3",
            EventID:   "sr:match:22222",
            MarketID:  "1",
            OutcomeID: "1",
            Odds:      &models.Odds{Type: "decimal", Value: "2.20"},
        },
        {
            Type:      "uf",
            ProductID: "3",
            EventID:   "sr:match:33333",
            MarketID:  "1",
            OutcomeID: "2",
            Odds:      &models.Odds{Type: "decimal", Value: "1.50"},
        },
    }
    
    stake := models.Stake{
        Type:     "cash",
        Currency: "EUR",
        Amount:   "2.00",
        Mode:     "unit",
    }
    
    context := &models.Context{
        Channel: &models.Channel{Type: "internet", Lang: "EN"},
        LimitID: 4268,
    }
    
    // 2. 选择构建器
    cfg := config.LoadConfig()
    builder := models.NewTicketBuilder(cfg.OperatorID)
    
    // 3. 生成注单结构
    // Banker 2/3 系统串: 从 3 个普通 selections 中选 2 个，每个组合都包含 Banker
    ticket, err := builder.NewBankerSystemBetTicket(
        "banker-001",
        []int{2},  // 2/3 系统串
        bankerSelections,
        regularSelections,
        stake,
        context,
    )
    
    return ticket, err
}
```

### 阶段二：验证 (Validate)

```go
func validateBankerBet(ticket *models.TicketRequest) error {
    bet := ticket.Content.Bets[0]
    systemSelection := bet.Selections[0]
    
    // Banker 系统串应该有嵌套的 system selection
    if systemSelection.Type != "system" {
        return fmt.Errorf("expected system selection for banker bet")
    }
    
    // 应该至少有 2 个 selections (1 banker + 1 regular)
    if len(systemSelection.Selections) < 2 {
        return fmt.Errorf("banker bet requires at least 1 banker and 1 regular selection")
    }
    
    // 第一个 selection 应该是 banker (type="system")
    bankerSel := systemSelection.Selections[0]
    if bankerSel.Type != "system" {
        return fmt.Errorf("first selection should be banker (type=system)")
    }
    
    return nil
}
```

---

## 示例 4: 通过 API 端点提交

### HTTP 请求示例

```bash
# 单注
curl -X POST http://localhost:8080/api/bets/single \
  -H "Content-Type: application/json" \
  -d '{
    "ticketId": "api-single-001",
    "selection": {
      "productId": "3",
      "eventId": "sr:match:12345",
      "marketId": "1",
      "outcomeId": "1",
      "odds": "2.50"
    },
    "stake": {
      "type": "cash",
      "currency": "EUR",
      "amount": "10.00",
      "mode": "total"
    }
  }'

# Trixie
curl -X POST http://localhost:8080/api/bets/preset \
  -H "Content-Type: application/json" \
  -d '{
    "ticketId": "api-trixie-001",
    "type": "trixie",
    "selections": [
      {
        "productId": "3",
        "eventId": "sr:match:11111",
        "marketId": "1",
        "outcomeId": "1",
        "odds": "1.50"
      },
      {
        "productId": "3",
        "eventId": "sr:match:22222",
        "marketId": "1",
        "outcomeId": "2",
        "odds": "2.00"
      },
      {
        "productId": "3",
        "eventId": "sr:match:33333",
        "marketId": "1",
        "outcomeId": "1",
        "odds": "1.80"
      }
    ],
    "stake": {
      "type": "cash",
      "currency": "EUR",
      "amount": "5.00",
      "mode": "unit"
    }
  }'
```

### API Handler 中的流程

```go
func (h *BetHandler) PlaceSingleBet(w http.ResponseWriter, r *http.Request) {
    // 阶段一：构建 (从 HTTP 请求)
    var req SingleBetRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Invalid request", http.StatusBadRequest)
        return
    }
    
    // 阶段二：验证
    if err := validateSingleBetRequest(&req); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }
    
    // 转换为 TicketRequest
    ticket := convertToTicketRequest(&req, h.cfg.OperatorID)
    
    // 阶段三：发射
    response, err := h.mtsService.SendTicket(ticket)
    if err != nil {
        http.Error(w, "Failed to send ticket", http.StatusInternalServerError)
        return
    }
    
    // 阶段四：确认并返回
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(response)
}
```

---

## 错误处理示例

### 处理 MTS 拒绝

```go
func handleRejectedTicket(response *models.TicketResponse) {
    code := response.Content.Code
    message := response.Content.Message
    
    switch code {
        case -401:
            log.Printf("Match not found: %s", message)
            // 通知用户赛事不可用
        case -701:
            log.Printf("Liability limit exceeded: %s", message)
            // 建议用户降低投注金额
        case -1001:
            log.Printf("Odds changed: %s", message)
            // 提供新赔率让用户确认
        default:
            log.Printf("Ticket rejected [%d]: %s", code, message)
    }
}
```

### 处理超时

```go
func handleTimeout(ticketID string) {
    log.Printf("Ticket %s timed out, marking as unknown", ticketID)
    
    // 1. 在数据库中标记为 "pending" 或 "unknown"
    // 2. 设置后台任务定期查询状态
    // 3. 通知用户稍后查看结果
}
```

---

## 总结

这些示例展示了如何在 `mts-service` 项目中实现完整的子弹装填流程。关键要点：

1. **使用 TicketBuilder**: 确保注单结构正确。
2. **本地验证**: 在发送前捕捉错误。
3. **错误处理**: 优雅地处理拒绝和超时。
4. **状态管理**: 追踪注单的完整生命周期。

通过遵循这些模式，您可以构建一个健壮、可靠的投注系统。
