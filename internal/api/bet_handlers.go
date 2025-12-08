package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gdsZyy/mts-service/internal/models"
)

// PlaceSingleBet handles single bet requests
func (h *Handler) PlaceSingleBet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondJSON(w, http.StatusMethodNotAllowed, APIResponse{
			Success: false,
			Error:   &APIError{Code: 405, Message: "Method not allowed"},
		})
		return
	}

	var req SingleBetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   &APIError{Code: 400, Message: "Invalid request body", Details: err.Error()},
		})
		return
	}

	// Validate request
	if err := validateSingleBetRequest(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   &APIError{Code: 400, Message: "Validation failed", Details: err.Error()},
		})
		return
	}

	// Build ticket using TicketBuilder
	builder := models.NewTicketBuilder(h.cfg.OperatorID, req.TicketID)
	
	selection := convertSelectionRequest(req.Selection)
	stake := convertStakeRequest(req.Stake)
	
	builder.AddSingleBet(selection, stake)
	
	if req.Context != nil {
		builder.SetContext(convertContextRequest(req.Context, h.cfg))
	} else {
		builder.SetContext(getDefaultContext(h.cfg))
	}

	ticket := builder.Build(generateCorrelationID())
	
	// Send to MTS
	response, err := h.mtsService.SendTicket(ticket)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   &APIError{Code: 500, Message: "Failed to send ticket", Details: err.Error()},
		})
		return
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    response,
	})
}

// PlaceAccumulatorBet handles accumulator bet requests
func (h *Handler) PlaceAccumulatorBet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondJSON(w, http.StatusMethodNotAllowed, APIResponse{
			Success: false,
			Error:   &APIError{Code: 405, Message: "Method not allowed"},
		})
		return
	}

	var req AccumulatorBetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   &APIError{Code: 400, Message: "Invalid request body", Details: err.Error()},
		})
		return
	}

	// Validate request
	if err := validateAccumulatorBetRequest(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   &APIError{Code: 400, Message: "Validation failed", Details: err.Error()},
		})
		return
	}

	// Build ticket using TicketBuilder
	builder := models.NewTicketBuilder(h.cfg.OperatorID, req.TicketID)
	
	
	selections := make([]models.Selection, len(req.Selections))
	for i, sel := range req.Selections {
		selections[i] = convertSelectionRequest(sel)
	}
	stake := convertStakeRequest(req.Stake)
	
	builder.AddAccumulatorBet(selections, stake)
	
	if req.Context != nil {
		builder.SetContext(convertContextRequest(req.Context, h.cfg))
	} else {
		builder.SetContext(getDefaultContext(h.cfg))
	}
	

	ticket := builder.Build(generateCorrelationID())
	
	// Send to MTS
	response, err := h.mtsService.SendTicket(ticket)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   &APIError{Code: 500, Message: "Failed to send ticket", Details: err.Error()},
		})
		return
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    response,
	})
}

// PlaceSystemBet handles system bet requests
func (h *Handler) PlaceSystemBet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondJSON(w, http.StatusMethodNotAllowed, APIResponse{
			Success: false,
			Error:   &APIError{Code: 405, Message: "Method not allowed"},
		})
		return
	}

	var req SystemBetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   &APIError{Code: 400, Message: "Invalid request body", Details: err.Error()},
		})
		return
	}

	// Validate request
	if err := validateSystemBetRequest(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   &APIError{Code: 400, Message: "Validation failed", Details: err.Error()},
		})
		return
	}

	// Build ticket using TicketBuilder
	builder := models.NewTicketBuilder(h.cfg.OperatorID, req.TicketID)
	
	
	selections := make([]models.Selection, len(req.Selections))
	for i, sel := range req.Selections {
		selections[i] = convertSelectionRequest(sel)
	}
	stake := convertStakeRequest(req.Stake)
	
	builder.AddSystemBet(req.Size, selections, stake)
	
	if req.Context != nil {
		builder.SetContext(convertContextRequest(req.Context, h.cfg))
	} else {
		builder.SetContext(getDefaultContext(h.cfg))
	}
	

	ticket := builder.Build(generateCorrelationID())
	
	// Send to MTS
	response, err := h.mtsService.SendTicket(ticket)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   &APIError{Code: 500, Message: "Failed to send ticket", Details: err.Error()},
		})
		return
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    response,
	})
}

// PlaceBankerSystemBet handles banker system bet requests
func (h *Handler) PlaceBankerSystemBet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondJSON(w, http.StatusMethodNotAllowed, APIResponse{
			Success: false,
			Error:   &APIError{Code: 405, Message: "Method not allowed"},
		})
		return
	}

	var req BankerSystemBetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   &APIError{Code: 400, Message: "Invalid request body", Details: err.Error()},
		})
		return
	}

	// Validate request
	if err := validateBankerSystemBetRequest(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   &APIError{Code: 400, Message: "Validation failed", Details: err.Error()},
		})
		return
	}

	// Build ticket using TicketBuilder
	builder := models.NewTicketBuilder(h.cfg.OperatorID, req.TicketID)
	
	
	bankers := make([]models.Selection, len(req.Bankers))
	for i, sel := range req.Bankers {
		bankers[i] = convertSelectionRequest(sel)
	}
	
	selections := make([]models.Selection, len(req.Selections))
	for i, sel := range req.Selections {
		selections[i] = convertSelectionRequest(sel)
	}
	stake := convertStakeRequest(req.Stake)
	
	builder.AddBankerSystemBet(bankers, req.Size, selections, stake)
	
	if req.Context != nil {
		builder.SetContext(convertContextRequest(req.Context, h.cfg))
	} else {
		builder.SetContext(getDefaultContext(h.cfg))
	}
	

	ticket := builder.Build(generateCorrelationID())
	
	// Send to MTS
	response, err := h.mtsService.SendTicket(ticket)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   &APIError{Code: 500, Message: "Failed to send ticket", Details: err.Error()},
		})
		return
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    response,
	})
}

// PlacePresetSystemBet handles preset system bet requests (Trixie, Yankee, etc.)
func (h *Handler) PlacePresetSystemBet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondJSON(w, http.StatusMethodNotAllowed, APIResponse{
			Success: false,
			Error:   &APIError{Code: 405, Message: "Method not allowed"},
		})
		return
	}

	var req PresetSystemBetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   &APIError{Code: 400, Message: "Invalid request body", Details: err.Error()},
		})
		return
	}

	// Validate request
	if err := validatePresetSystemBetRequest(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   &APIError{Code: 400, Message: "Validation failed", Details: err.Error()},
		})
		return
	}

	// Build ticket using TicketBuilder
	builder := models.NewTicketBuilder(h.cfg.OperatorID, req.TicketID)
	
	
	selections := make([]models.Selection, len(req.Selections))
	for i, sel := range req.Selections {
		selections[i] = convertSelectionRequest(sel)
	}
	stake := convertStakeRequest(req.Stake)
	
	// Call appropriate method based on type
	switch strings.ToLower(req.Type) {
	case "trixie":
		builder.AddTrixieBet(selections, stake)
	case "patent":
		builder.AddPatentBet(selections, stake)
	case "yankee":
		builder.AddYankeeBet(selections, stake)
	case "lucky15", "lucky_15":
		builder.AddLucky15Bet(selections, stake)
	case "super_yankee", "canadian":
		builder.AddSuperYankeeBet(selections, stake)
	case "lucky31", "lucky_31":
		builder.AddLucky31Bet(selections, stake)
	case "heinz":
		builder.AddHeinzBet(selections, stake)
	case "lucky63", "lucky_63":
		builder.AddLucky63Bet(selections, stake)
	case "super_heinz":
		builder.AddSuperHeinzBet(selections, stake)
	case "goliath":
		builder.AddGoliathBet(selections, stake)
	default:
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   &APIError{Code: 400, Message: "Invalid preset type", Details: fmt.Sprintf("Unknown type: %s", req.Type)},
		})
		return
	}
	
	if req.Context != nil {
		builder.SetContext(convertContextRequest(req.Context, h.cfg))
	} else {
		builder.SetContext(getDefaultContext(h.cfg))
	}
	

	ticket := builder.Build(generateCorrelationID())
	
	// Send to MTS
	response, err := h.mtsService.SendTicket(ticket)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   &APIError{Code: 500, Message: "Failed to send ticket", Details: err.Error()},
		})
		return
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    response,
	})
}

// PlaceMultiBet handles multi-bet requests (multiple bets in one ticket)
func (h *Handler) PlaceMultiBet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		respondJSON(w, http.StatusMethodNotAllowed, APIResponse{
			Success: false,
			Error:   &APIError{Code: 405, Message: "Method not allowed"},
		})
		return
	}

	var req MultiBetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   &APIError{Code: 400, Message: "Invalid request body", Details: err.Error()},
		})
		return
	}

	// Validate request
	if err := validateMultiBetRequest(&req); err != nil {
		respondJSON(w, http.StatusBadRequest, APIResponse{
			Success: false,
			Error:   &APIError{Code: 400, Message: "Validation failed", Details: err.Error()},
		})
		return
	}

	// Build ticket using TicketBuilder
	builder := models.NewTicketBuilder(h.cfg.OperatorID, req.TicketID)
	
	for _, bet := range req.Bets {
		selections := make([]models.Selection, len(bet.Selections))
		for i, sel := range bet.Selections {
			selections[i] = convertSelectionRequest(sel)
		}
		stake := convertStakeRequest(bet.Stake)
		
		switch strings.ToLower(bet.Type) {
		case "single":
			if len(selections) != 1 {
				respondJSON(w, http.StatusBadRequest, APIResponse{
					Success: false,
					Error:   &APIError{Code: 400, Message: "Single bet must have exactly 1 selection"},
				})
				return
			}
			builder.AddSingleBet(selections[0], stake)
		case "accumulator":
			builder.AddAccumulatorBet(selections, stake)
		case "system":
			builder.AddSystemBet(bet.Size, selections, stake)
		case "banker_system":
			bankers := make([]models.Selection, len(bet.Bankers))
			for i, sel := range bet.Bankers {
				bankers[i] = convertSelectionRequest(sel)
			}
			builder.AddBankerSystemBet(bankers, bet.Size, selections, stake)
		default:
			// Try preset types
			switch strings.ToLower(bet.Type) {
			case "trixie":
				builder.AddTrixieBet(selections, stake)
			case "patent":
				builder.AddPatentBet(selections, stake)
			case "yankee":
				builder.AddYankeeBet(selections, stake)
			case "lucky15", "lucky_15":
				builder.AddLucky15Bet(selections, stake)
			case "super_yankee", "canadian":
				builder.AddSuperYankeeBet(selections, stake)
			case "lucky31", "lucky_31":
				builder.AddLucky31Bet(selections, stake)
			case "heinz":
				builder.AddHeinzBet(selections, stake)
			case "lucky63", "lucky_63":
				builder.AddLucky63Bet(selections, stake)
			case "super_heinz":
				builder.AddSuperHeinzBet(selections, stake)
			case "goliath":
				builder.AddGoliathBet(selections, stake)
			default:
				respondJSON(w, http.StatusBadRequest, APIResponse{
					Success: false,
					Error:   &APIError{Code: 400, Message: "Invalid bet type", Details: fmt.Sprintf("Unknown type: %s", bet.Type)},
				})
				return
			}
		}
	}
	
	if req.Context != nil {
		builder.SetContext(convertContextRequest(req.Context, h.cfg))
	} else {
		builder.SetContext(getDefaultContext(h.cfg))
	}
	

	ticket := builder.Build(generateCorrelationID())
	
	// Send to MTS
	response, err := h.mtsService.SendTicket(ticket)
	if err != nil {
		respondJSON(w, http.StatusInternalServerError, APIResponse{
			Success: false,
			Error:   &APIError{Code: 500, Message: "Failed to send ticket", Details: err.Error()},
		})
		return
	}

	respondJSON(w, http.StatusOK, APIResponse{
		Success: true,
		Data:    response,
	})
}

// Helper function to respond with JSON
func respondJSON(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

// Helper function to generate correlation ID
func generateCorrelationID() string {
	return fmt.Sprintf("corr-%d", time.Now().UnixNano())
}
