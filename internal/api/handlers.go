package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gdsZyy/mts-service/internal/config"
	"github.com/gdsZyy/mts-service/internal/models"
	"github.com/gdsZyy/mts-service/internal/service"
	"github.com/google/uuid"
)

type Handler struct {
	mtsService *service.MTSService
	cfg        *config.Config
}

func NewHandler(mtsService *service.MTSService, cfg *config.Config) *Handler {
	return &Handler{
		mtsService: mtsService,
		cfg:        cfg,
	}
}

// HealthCheck handles health check requests
func (h *Handler) HealthCheck(w http.ResponseWriter, r *http.Request) {
	status := "healthy"
	if !h.mtsService.IsConnected() {
		status = "disconnected"
	}

	response := map[string]interface{}{
		"status":    status,
		"timestamp": time.Now().Unix(),
		"service":   "mts-service",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// PlaceTicket handles ticket placement requests
func (h *Handler) PlaceTicket(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req PlaceTicketRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Invalid request body", err)
		return
	}

	// Validate request
	if err := validatePlaceTicketRequest(&req); err != nil {
		respondError(w, http.StatusBadRequest, "Validation failed", err)
		return
	}

	// Build MTS ticket request
	ticket := h.buildTicketRequest(&req)

	log.Printf("Sending ticket: %s (correlation: %s)", ticket.Content.TicketID, ticket.CorrelationID)

	// Send to MTS
	response, err := h.mtsService.SendTicket(ticket)
	if err != nil {
		respondError(w, http.StatusInternalServerError, "Failed to send ticket", err)
		return
	}

	log.Printf("Received response for ticket: %s, status: %s", response.Content.TicketID, response.Content.Status)

	// Return response
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

// PlaceTicketRequest represents the API request for placing a ticket
type PlaceTicketRequest struct {
	TicketID    string      `json:"ticketId"`
	CustomerID  string      `json:"customerId"`
	Currency    string      `json:"currency"`
	TotalStake  int64       `json:"totalStake"`
	TestSource  bool        `json:"testSource"`
	OddsChange  string      `json:"oddsChange"`
	Bets        []BetInput  `json:"bets"`
	CustomerIP  string      `json:"customerIp,omitempty"`
	DeviceID    string      `json:"deviceId,omitempty"`
	LanguageID  string      `json:"languageId,omitempty"`
	Channel     string      `json:"channel,omitempty"`
}

// BetInput represents a bet in the API request
type BetInput struct {
	ID         string           `json:"id"`
	Stake      int64            `json:"stake"`
	CustomBet  bool             `json:"customBet"`
	Selections []SelectionInput `json:"selections"`
}

// SelectionInput represents a selection in the API request
type SelectionInput struct {
	ID      string `json:"id"`
	EventID string `json:"eventId"`
	Odds    int    `json:"odds"`
	Banker  bool   `json:"banker"`
}

func validatePlaceTicketRequest(req *PlaceTicketRequest) error {
	if req.TicketID == "" {
		return fmt.Errorf("ticketId is required")
	}
	if req.CustomerID == "" {
		return fmt.Errorf("customerId is required")
	}
	if req.Currency == "" {
		return fmt.Errorf("currency is required")
	}
	if req.TotalStake <= 0 {
		return fmt.Errorf("totalStake must be positive")
	}
	if len(req.Bets) == 0 {
		return fmt.Errorf("at least one bet is required")
	}

	for i, bet := range req.Bets {
		if bet.ID == "" {
			return fmt.Errorf("bet[%d].id is required", i)
		}
		if bet.Stake <= 0 {
			return fmt.Errorf("bet[%d].stake must be positive", i)
		}
		if len(bet.Selections) == 0 {
			return fmt.Errorf("bet[%d] must have at least one selection", i)
		}

		for j, sel := range bet.Selections {
			if sel.ID == "" {
				return fmt.Errorf("bet[%d].selection[%d].id is required", i, j)
			}
			if sel.EventID == "" {
				return fmt.Errorf("bet[%d].selection[%d].eventId is required", i, j)
			}
			if sel.Odds <= 0 {
				return fmt.Errorf("bet[%d].selection[%d].odds must be positive", i, j)
			}
		}
	}

	return nil
}

func (h *Handler) buildTicketRequest(req *PlaceTicketRequest) *models.TicketRequest {
		// Generate correlation ID
		correlationID := uuid.New().String()

		// Operator ID is mandatory
		operatorID := h.cfg.OperatorID
		if operatorID == 0 {
			log.Println("Warning: OperatorID is not set in config. Using default 9985.")
			operatorID = 9985 // Fallback or a known test ID
		}

	// Build bets
	bets := make([]models.Bet, len(req.Bets))
	for i, betInput := range req.Bets {
		selections := make([]models.Selection, len(betInput.Selections))
		for j, selInput := range betInput.Selections {
			selections[j] = models.Selection{
				ID:      selInput.ID,
				EventID: selInput.EventID,
				Odds:    selInput.Odds,
				Banker:  selInput.Banker,
			}
		}

		bets[i] = models.Bet{
			ID:         betInput.ID,
			Stake:      betInput.Stake,
			CustomBet:  betInput.CustomBet,
			Selections: selections,
		}
	}

	// Default values
	oddsChange := req.OddsChange
	if oddsChange == "" {
		oddsChange = "any"
	}

	channel := req.Channel
	if channel == "" {
		channel = "internet"
	}

	languageID := req.LanguageID
	if languageID == "" {
		languageID = "en"
	}

		return &models.TicketRequest{
			OperatorID:    operatorID, // Added operatorId
			Operation:     "ticket-placement",
			CorrelationID: correlationID,
			TimestampUTC:  time.Now().UnixMilli(),
			Version:       "2.4",
			Content: models.TicketContent{
				Type:            "ticket",
				TicketID:        req.TicketID,
				TicketSignature: "", // For ticket-placement, signature is usually optional/empty unless required by config
				TotalStake:      req.TotalStake,
				TestSource:      req.TestSource,
				OddsChange:      oddsChange,
				Sender: models.Sender{
					Bookmaker: h.cfg.BookmakerID,
					Currency:  req.Currency,
					Channel:   channel,
					EndCustomer: models.EndCustomer{
						ID:         req.CustomerID,
						IP:         req.CustomerIP,
						LanguageID: languageID,
						DeviceID:   req.DeviceID,
						Confidence: 12092, // Default CCF
					},
				},
				Bets: bets,
			},
		}
}

func respondError(w http.ResponseWriter, status int, message string, err error) {
	log.Printf("Error: %s - %v", message, err)

	response := map[string]interface{}{
		"error":   message,
		"details": err.Error(),
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(response)
}

