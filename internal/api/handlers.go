package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
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
// This structure accepts flexible input and maps it to MTS Transaction 3.0 API standard
type PlaceTicketRequest struct {
	TicketID   string      `json:"ticketId"`
	CustomerID string      `json:"customerId"`
	Currency   string      `json:"currency"`
	TotalStake string      `json:"totalStake"` // Changed to string to support 8 decimal places
	Bets       []BetInput  `json:"bets"`
	CustomerIP string      `json:"customerIp,omitempty"`
	DeviceID   string      `json:"deviceId,omitempty"`
	LanguageID string      `json:"languageId,omitempty"`
	Channel    string      `json:"channel,omitempty"`
	ProductID  string      `json:"productId,omitempty"` // Default product ID for selections
	MarketID   string      `json:"marketId,omitempty"`  // Default market ID for selections
	BetType    string      `json:"betType,omitempty"`   // "single", "system", "parlay" etc. For system/串关, use "system"
}

// BetInput represents a bet in the API request
type BetInput struct {
	Selections []SelectionInput `json:"selections"`
	Amount     string           `json:"amount"` // Stake amount as string with up to 8 decimal places
}

// SelectionInput represents a selection in the API request
type SelectionInput struct {
	EventID   string `json:"eventId"`
	OutcomeID string `json:"outcomeId"`
	Odds      string `json:"odds"` // Odds as string (e.g., "1.59")
	ProductID string `json:"productId,omitempty"`
	MarketID  string `json:"marketId,omitempty"`
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
	if req.TotalStake == "" {
		return fmt.Errorf("totalStake is required")
	}
	if len(req.Bets) == 0 {
		return fmt.Errorf("at least one bet is required")
	}

	for i, bet := range req.Bets {
		if bet.Amount == "" {
			return fmt.Errorf("bet[%d].amount is required", i)
		}
		if len(bet.Selections) == 0 {
			return fmt.Errorf("bet[%d] must have at least one selection", i)
		}

		for j, sel := range bet.Selections {
			if sel.EventID == "" {
				return fmt.Errorf("bet[%d].selection[%d].eventId is required", i, j)
			}
			if sel.OutcomeID == "" {
				return fmt.Errorf("bet[%d].selection[%d].outcomeId is required", i, j)
			}
			if sel.Odds == "" {
				return fmt.Errorf("bet[%d].selection[%d].odds is required", i, j)
			}
		}
	}

	return nil
}

// formatAmountTo8Decimals formats an amount string to have exactly 8 decimal places
func formatAmountTo8Decimals(amount string) string {
	// Parse the amount as a float
	val, err := strconv.ParseFloat(amount, 64)
	if err != nil {
		log.Printf("Warning: Failed to parse amount '%s': %v, using as-is", amount, err)
		return amount
	}

	// Format to 8 decimal places
	return fmt.Sprintf("%.8f", val)
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

	// Set default product ID and market ID if not provided
	defaultProductID := req.ProductID
	if defaultProductID == "" {
		defaultProductID = "3" // Default product ID
	}

	defaultMarketID := req.MarketID
	if defaultMarketID == "" {
		defaultMarketID = "14" // Default market ID
	}

	// Determine if this is a system/parlay bet (only one betId for all selections)
	isSystemBet := strings.ToLower(req.BetType) == "system" || strings.ToLower(req.BetType) == "parlay"

	var bets []models.Bet

	if isSystemBet {
		// For system/parlay bets: all selections go into a single bet with a single betId
		// betId format: ticketId-1
		allSelections := []models.Selection{}

		for _, betInput := range req.Bets {
			for _, selInput := range betInput.Selections {
				// Use provided product/market ID or fall back to defaults
				productID := selInput.ProductID
				if productID == "" {
					productID = defaultProductID
				}

				marketID := selInput.MarketID
				if marketID == "" {
					marketID = defaultMarketID
				}

				allSelections = append(allSelections, models.Selection{
					Type:       "uf", // Unified Feed binding type
					ProductID:  productID,
					EventID:    selInput.EventID,
					MarketID:   marketID,
					OutcomeID:  selInput.OutcomeID,
					Odds: models.Odds{
						Type:  "decimal",
						Value: selInput.Odds,
					},
				})
			}
		}

		// For system bet, use the first bet's amount or total stake
		stakeAmount := req.Bets[0].Amount
		if stakeAmount == "" {
			stakeAmount = req.TotalStake
		}
		stakeAmount = formatAmountTo8Decimals(stakeAmount)

		// Default mode to "total" if not specified
		stakeMode := "total"

		bets = []models.Bet{
			{
				Selections: allSelections,
				Stake: []models.Stake{
					{
						Type:     "cash",
						Currency: req.Currency,
						Amount:   stakeAmount,
						Mode:     stakeMode,
					},
				},
			},
		}
	} else {
		// For regular/single bets: each bet gets its own betId
		bets = make([]models.Bet, len(req.Bets))
		for i, betInput := range req.Bets {
			selections := make([]models.Selection, len(betInput.Selections))
			for j, selInput := range betInput.Selections {
				// Use provided product/market ID or fall back to defaults
				productID := selInput.ProductID
				if productID == "" {
					productID = defaultProductID
				}

				marketID := selInput.MarketID
				if marketID == "" {
					marketID = defaultMarketID
				}

				selections[j] = models.Selection{
					Type:       "uf", // Unified Feed binding type
					ProductID:  productID,
					EventID:    selInput.EventID,
					MarketID:   marketID,
					OutcomeID:  selInput.OutcomeID,
					Odds: models.Odds{
						Type:  "decimal",
						Value: selInput.Odds,
					},
				}
			}

			// Convert stake amount to string with 8 decimal places
			stakeAmount := betInput.Amount
			if stakeAmount == "" {
				// Calculate stake from total stake if not provided
				totalStakeVal, _ := strconv.ParseFloat(req.TotalStake, 64)
				stakeAmount = fmt.Sprintf("%.8f", totalStakeVal/float64(len(req.Bets)))
			} else {
				stakeAmount = formatAmountTo8Decimals(stakeAmount)
			}

			// Default mode to "total" if not specified
			stakeMode := "total"

			bets[i] = models.Bet{
				Selections: selections,
				Stake: []models.Stake{
					{
						Type:     "cash",
						Currency: req.Currency,
						Amount:   stakeAmount,
						Mode:     stakeMode,
					},
				},
			}
		}
	}

	// Set default channel and language
	channel := req.Channel
	if channel == "" {
		channel = "internet"
	}

	languageID := req.LanguageID
	if languageID == "" {
		languageID = "EN"
	}

	return &models.TicketRequest{
		OperatorID:    operatorID,
		CorrelationID: correlationID,
		TimestampUTC:  time.Now().UnixMilli(),
		Operation:     "ticket-placement", // Standard MTS Transaction 3.0 operation
		Version:       "3.0",              // Standard MTS Transaction 3.0 version
		Content: models.TicketContent{
			Type:     "ticket", // Standard MTS Transaction 3.0 content type
			TicketID: req.TicketID,
			Bets:     bets,
			Context: &models.Context{
				Channel: &models.Channel{
					Type: channel,
					Lang: languageID,
				},
				IP: req.CustomerIP,
			},
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
