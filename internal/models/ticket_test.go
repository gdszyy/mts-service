package models

import (
	"encoding/json"
	"testing"
)

func TestTicketResponseUnmarshal(t *testing.T) {
	// MTS 实际返回的 JSON（基于项目日志）
	jsonData := `{
		"content": {
			"type": "ticket-reply",
			"message": "Match is not found in MTS (selection: uof:3/sr:match:16470657/534/pre:outcometext:9919, match: sr:match:16470657)",
			"code": -401,
			"signature": "EekSnMy+stHi2PPSYStrZz6pntqVN3gkRGryKe72ot7Vaz1zSHYBg5HY/N1F50HSjpierCiFcnwaRgG8",
			"ticketId": "init-ticket-1764059527727917670",
			"status": "rejected",
			"betDetails": [
				{
					"code": -401,
					"message": "Match is not found in MTS (selection: uof:3/sr:match:16470657/534/pre:outcometext:9919, match: sr:match:16470657)",
					"betId": "",
					"selectionDetails": [
						{
							"code": -401,
							"message": "Match is not found in MTS (selection: uof:3/sr:match:16470657/534/pre:outcometext:9919, match: sr:match:16470657)",
							"selection": {
								"type": "uf",
								"productId": "3",
								"eventId": "sr:match:16470657",
								"marketId": "534",
								"specifiers": "",
								"outcomeId": "pre:outcometext:9919",
								"odds": {
									"type": "decimal",
									"value": "2.10"
								}
							}
						}
					]
				}
			],
			"exchangeRate": [
				{
					"fromCurrency": "EUR",
					"toCurrency": "EUR",
					"rate": "1.00000000"
				}
			]
		},
		"correlationId": "init-1764059527727914658",
		"timestampUtc": 1764059528040,
		"operation": "ticket-placement",
		"version": "3.0"
	}`

	var response TicketResponse
	err := json.Unmarshal([]byte(jsonData), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response: %v", err)
	}

	// 验证基本字段
	if response.Content.Type != "ticket-reply" {
		t.Errorf("Expected type 'ticket-reply', got '%s'", response.Content.Type)
	}

	if response.Content.Status != "rejected" {
		t.Errorf("Expected status 'rejected', got '%s'", response.Content.Status)
	}

	if response.Content.Code != -401 {
		t.Errorf("Expected code -401, got %d", response.Content.Code)
	}

	if response.Content.Message == "" {
		t.Error("Expected non-empty message")
	}

	// 验证 exchangeRate 数组
	if len(response.Content.ExchangeRate) != 1 {
		t.Fatalf("Expected 1 exchange rate, got %d", len(response.Content.ExchangeRate))
	}

	rate := response.Content.ExchangeRate[0]
	if rate.FromCurrency != "EUR" {
		t.Errorf("Expected fromCurrency 'EUR', got '%s'", rate.FromCurrency)
	}

	if rate.ToCurrency != "EUR" {
		t.Errorf("Expected toCurrency 'EUR', got '%s'", rate.ToCurrency)
	}

	if rate.Rate != "1.00000000" {
		t.Errorf("Expected rate '1.00000000', got '%s'", rate.Rate)
	}

	// 验证 betDetails
	if len(response.Content.BetDetails) != 1 {
		t.Fatalf("Expected 1 bet detail, got %d", len(response.Content.BetDetails))
	}

	betDetail := response.Content.BetDetails[0]
	if betDetail.Code != -401 {
		t.Errorf("Expected betDetail code -401, got %d", betDetail.Code)
	}

	if betDetail.Message == "" {
		t.Error("Expected non-empty betDetail message")
	}

	// 验证 selectionDetails
	if len(betDetail.SelectionDetails) != 1 {
		t.Fatalf("Expected 1 selection detail, got %d", len(betDetail.SelectionDetails))
	}

	selDetail := betDetail.SelectionDetails[0]
	if selDetail.Code != -401 {
		t.Errorf("Expected selectionDetail code -401, got %d", selDetail.Code)
	}

	if selDetail.Selection.EventID != "sr:match:16470657" {
		t.Errorf("Expected eventId 'sr:match:16470657', got '%s'", selDetail.Selection.EventID)
	}

	if selDetail.Selection.MarketID != "534" {
		t.Errorf("Expected marketId '534', got '%s'", selDetail.Selection.MarketID)
	}

	if selDetail.Selection.OutcomeID != "pre:outcometext:9919" {
		t.Errorf("Expected outcomeId 'pre:outcometext:9919', got '%s'", selDetail.Selection.OutcomeID)
	}
}

func TestExchangeRateEmptyArray(t *testing.T) {
	// 测试空数组情况
	jsonData := `{
		"content": {
			"type": "ticket-reply",
			"ticketId": "test-123",
			"status": "accepted",
			"code": 0,
			"signature": "test-sig",
			"exchangeRate": []
		},
		"correlationId": "test-corr",
		"timestampUtc": 1234567890,
		"operation": "ticket-placement",
		"version": "3.0"
	}`

	var response TicketResponse
	err := json.Unmarshal([]byte(jsonData), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response with empty exchangeRate: %v", err)
	}

	if len(response.Content.ExchangeRate) != 0 {
		t.Errorf("Expected empty exchange rate array, got %d items", len(response.Content.ExchangeRate))
	}
}

func TestTicketResponseAccepted(t *testing.T) {
	// 测试接受的票据响应
	jsonData := `{
		"content": {
			"type": "ticket-reply",
			"message": "Transaction processed",
			"code": 0,
			"signature": "X18sOrIkGst6JoZBYSgfbCzsGHdncjIWksAMkfPsfNWoxxwtvVIoNhj3ceYkHExPmKBxeRAeQAyocK7s",
			"ticketId": "Ticket_701941",
			"status": "accepted",
			"betDetails": [
				{
					"code": 0,
					"message": "Transaction processed",
					"betId": "Ticket_701941-0",
					"selectionDetails": [
						{
							"code": 0,
							"message": "Transaction processed",
							"selection": {
								"type": "uf",
								"productId": "3",
								"eventId": "sr:match:52102093",
								"marketId": "304",
								"specifiers": "quarternr=1",
								"outcomeId": "70",
								"odds": {
									"type": "decimal",
									"value": "1.5"
								}
							}
						}
					]
				}
			],
			"exchangeRate": [
				{
					"fromCurrency": "EUR",
					"toCurrency": "EUR",
					"rate": "1.0"
				}
			]
		},
		"correlationId": "MTS_Ticket_701941",
		"timestampUtc": 1739373264316,
		"operation": "ticket-placement",
		"version": "3.0"
	}`

	var response TicketResponse
	err := json.Unmarshal([]byte(jsonData), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal accepted response: %v", err)
	}

	if response.Content.Status != "accepted" {
		t.Errorf("Expected status 'accepted', got '%s'", response.Content.Status)
	}

	if response.Content.Code != 0 {
		t.Errorf("Expected code 0, got %d", response.Content.Code)
	}

	if response.Content.Message != "Transaction processed" {
		t.Errorf("Expected message 'Transaction processed', got '%s'", response.Content.Message)
	}
}

func TestTicketResponseMultipleExchangeRates(t *testing.T) {
	// 测试多个汇率的情况
	jsonData := `{
		"content": {
			"type": "ticket-reply",
			"ticketId": "test-multi-rate",
			"status": "accepted",
			"code": 0,
			"signature": "test-sig",
			"exchangeRate": [
				{
					"fromCurrency": "USD",
					"toCurrency": "EUR",
					"rate": "0.931741"
				},
				{
					"fromCurrency": "GBP",
					"toCurrency": "EUR",
					"rate": "1.15234"
				}
			]
		},
		"correlationId": "test-corr",
		"timestampUtc": 1234567890,
		"operation": "ticket-placement",
		"version": "3.0"
	}`

	var response TicketResponse
	err := json.Unmarshal([]byte(jsonData), &response)
	if err != nil {
		t.Fatalf("Failed to unmarshal response with multiple exchange rates: %v", err)
	}

	if len(response.Content.ExchangeRate) != 2 {
		t.Fatalf("Expected 2 exchange rates, got %d", len(response.Content.ExchangeRate))
	}

	// 验证第一个汇率
	if response.Content.ExchangeRate[0].FromCurrency != "USD" {
		t.Errorf("Expected first fromCurrency 'USD', got '%s'", response.Content.ExchangeRate[0].FromCurrency)
	}

	// 验证第二个汇率
	if response.Content.ExchangeRate[1].FromCurrency != "GBP" {
		t.Errorf("Expected second fromCurrency 'GBP', got '%s'", response.Content.ExchangeRate[1].FromCurrency)
	}
}
