package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gdsZyy/mts-service/internal/models"
)

// RequestCashout handles cashout-inform requests
func (h *Handler) RequestCashout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondJSON(w, http.StatusMethodNotAllowed, APIResponse{
			Success: false,
			Error:   &APIError{Code: 405, Message: "Method not allowed"},
		})
		return
	}

	var req CashoutRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   &APIError{Code: 400, Message: "Invalid request body", Details: err.Error()},
		})
		return
	}

	// Validate request
	if err := validateCashoutRequest(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   &APIError{Code: 400, Message: "Validation failed", Details: err.Error()},
		})
		return
	}

	// Build cashout request
	cashoutReq := buildCashoutRequest(&req, h.cfg.OperatorID)

	// Send to MTS (需要在 service 层实现 SendCashout 方法)
	// TODO: Implement SendCashout in MTSService
	// response, err := h.mtsService.SendCashout(cashoutReq)
	
	// 临时返回错误，提示需要实现
	respondJSON(w, http.StatusNotImplemented, APIResponse{
		Success: false,
		Error:   &APIError{
			Code:    501,
			Message: "Cashout functionality not yet implemented in service layer",
			Details: fmt.Sprintf("Cashout request built successfully: %+v", cashoutReq),
		},
	})
}

func validateCashoutRequest(req *CashoutRequest) error {
	if req.CashoutID == "" {
		return fmt.Errorf("cashoutId is required")
	}
	if req.TicketID == "" {
		return fmt.Errorf("ticketId is required")
	}
	if req.TicketSignature == "" {
		return fmt.Errorf("ticketSignature is required")
	}
	if req.Type == "" {
		return fmt.Errorf("type is required")
	}
	validTypes := map[string]bool{
		"ticket":         true,
		"ticket-partial": true,
		"bet":            true,
		"bet-partial":    true,
	}
	if !validTypes[req.Type] {
		return fmt.Errorf("invalid type: %s", req.Type)
	}
	if req.Code == 0 {
		return fmt.Errorf("code is required")
	}
	if len(req.Payout) == 0 {
		return fmt.Errorf("at least one payout is required")
	}
	
	// Validate partial cashout
	if req.Type == "ticket-partial" || req.Type == "bet-partial" {
		if req.Percentage <= 0 || req.Percentage > 1 {
			return fmt.Errorf("percentage must be between 0 and 1 for partial cashout")
		}
	}
	
	// Validate bet-level cashout
	if req.Type == "bet" || req.Type == "bet-partial" {
		if req.BetID == "" {
			return fmt.Errorf("betId is required for bet-level cashout")
		}
	}
	
	for i, payout := range req.Payout {
		if payout.Type == "" {
			return fmt.Errorf("payout[%d].type is required", i)
		}
		if payout.Currency == "" {
			return fmt.Errorf("payout[%d].currency is required", i)
		}
		if payout.Amount <= 0 {
			return fmt.Errorf("payout[%d].amount must be greater than 0", i)
		}
	}
	
	return nil
}

func buildCashoutRequest(req *CashoutRequest, operatorID int64) *models.CashoutRequest {
	correlationID := fmt.Sprintf("cashout-corr-%d", time.Now().UnixNano())
	
	// Convert payout requests to models
	payouts := make([]models.CashoutPayout, len(req.Payout))
	for i, p := range req.Payout {
		payouts[i] = models.CashoutPayout{
			Type:     p.Type,
			Currency: p.Currency,
			Amount:   fmt.Sprintf("%.8f", p.Amount),
			Source:   p.Source,
		}
	}
	
	// Build cashout detail
	detail := models.CashoutDetail{
		Type:            req.Type,
		TicketID:        req.TicketID,
		TicketSignature: req.TicketSignature,
		Code:            req.Code,
		Payout:          payouts,
	}
	
	if req.Percentage > 0 {
		detail.Percentage = fmt.Sprintf("%.2f", req.Percentage)
	}
	
	if req.BetID != "" {
		detail.BetID = req.BetID
	}
	
	return &models.CashoutRequest{
		OperatorID:    operatorID,
		CorrelationID: correlationID,
		TimestampUTC:  time.Now().UnixMilli(),
		Operation:     "cashout-inform",
		Version:       "3.0",
		Content: models.CashoutContent{
			Type: "cashout-inform",
			Cashout: models.CashoutInfo{
				Type:      "cashout",
				CashoutID: req.CashoutID,
				Details:   detail,
			},
		},
	}
}
