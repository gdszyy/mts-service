package api

import (
	"fmt"
	"strconv"

	"github.com/gdsZyy/mts-service/internal/config"
	"github.com/gdsZyy/mts-service/internal/models"
)

// Validation functions

func validateSingleBetRequest(req *SingleBetRequest) error {
	if req.TicketID == "" {
		return fmt.Errorf("ticketId is required")
	}
	if err := validateSelectionRequest(&req.Selection); err != nil {
		return fmt.Errorf("selection: %w", err)
	}
	if err := validateStakeRequest(&req.Stake); err != nil {
		return fmt.Errorf("stake: %w", err)
	}
	return nil
}

func validateAccumulatorBetRequest(req *AccumulatorBetRequest) error {
	if req.TicketID == "" {
		return fmt.Errorf("ticketId is required")
	}
	if len(req.Selections) < 2 {
		return fmt.Errorf("accumulator requires at least 2 selections")
	}
	for i, sel := range req.Selections {
		if err := validateSelectionRequest(&sel); err != nil {
			return fmt.Errorf("selection[%d]: %w", i, err)
		}
	}
	if err := validateStakeRequest(&req.Stake); err != nil {
		return fmt.Errorf("stake: %w", err)
	}
	return nil
}

func validateSystemBetRequest(req *SystemBetRequest) error {
	if req.TicketID == "" {
		return fmt.Errorf("ticketId is required")
	}
	if len(req.Size) == 0 {
		return fmt.Errorf("size is required")
	}
	if len(req.Selections) < 2 {
		return fmt.Errorf("system bet requires at least 2 selections")
	}
	for _, s := range req.Size {
		if s < 1 || s > len(req.Selections) {
			return fmt.Errorf("invalid size %d for %d selections", s, len(req.Selections))
		}
	}
	for i, sel := range req.Selections {
		if err := validateSelectionRequest(&sel); err != nil {
			return fmt.Errorf("selection[%d]: %w", i, err)
		}
	}
	if err := validateStakeRequest(&req.Stake); err != nil {
		return fmt.Errorf("stake: %w", err)
	}
	// System bets should use "unit" mode
	if req.Stake.Mode != "unit" {
		return fmt.Errorf("system bet stake mode must be 'unit'")
	}
	return nil
}

func validateBankerSystemBetRequest(req *BankerSystemBetRequest) error {
	if req.TicketID == "" {
		return fmt.Errorf("ticketId is required")
	}
	if len(req.Bankers) < 1 {
		return fmt.Errorf("banker system bet requires at least 1 banker")
	}
	if len(req.Selections) < 1 {
		return fmt.Errorf("banker system bet requires at least 1 non-banker selection")
	}
	if len(req.Size) == 0 {
		return fmt.Errorf("size is required")
	}
	for _, s := range req.Size {
		if s < 1 || s > len(req.Selections) {
			return fmt.Errorf("invalid size %d for %d non-banker selections", s, len(req.Selections))
		}
	}
	for i, sel := range req.Bankers {
		if err := validateSelectionRequest(&sel); err != nil {
			return fmt.Errorf("banker[%d]: %w", i, err)
		}
	}
	for i, sel := range req.Selections {
		if err := validateSelectionRequest(&sel); err != nil {
			return fmt.Errorf("selection[%d]: %w", i, err)
		}
	}
	if err := validateStakeRequest(&req.Stake); err != nil {
		return fmt.Errorf("stake: %w", err)
	}
	if req.Stake.Mode != "unit" {
		return fmt.Errorf("banker system bet stake mode must be 'unit'")
	}
	return nil
}

func validatePresetSystemBetRequest(req *PresetSystemBetRequest) error {
	if req.TicketID == "" {
		return fmt.Errorf("ticketId is required")
	}
	if req.Type == "" {
		return fmt.Errorf("type is required")
	}
	
	// Validate selection count based on type
	requiredCount := map[string]int{
		"trixie":       3,
		"patent":       3,
		"yankee":       4,
		"lucky15":      4,
		"lucky_15":     4,
		"super_yankee": 5,
		"canadian":     5,
		"lucky31":      5,
		"lucky_31":     5,
		"heinz":        6,
		"lucky63":      6,
		"lucky_63":     6,
		"super_heinz":  7,
		"goliath":      8,
	}
	
	required, ok := requiredCount[req.Type]
	if !ok {
		return fmt.Errorf("unknown preset type: %s", req.Type)
	}
	
	if len(req.Selections) != required {
		return fmt.Errorf("%s requires exactly %d selections, got %d", req.Type, required, len(req.Selections))
	}
	
	for i, sel := range req.Selections {
		if err := validateSelectionRequest(&sel); err != nil {
			return fmt.Errorf("selection[%d]: %w", i, err)
		}
	}
	if err := validateStakeRequest(&req.Stake); err != nil {
		return fmt.Errorf("stake: %w", err)
	}
	if req.Stake.Mode != "unit" {
		return fmt.Errorf("preset system bet stake mode must be 'unit'")
	}
	return nil
}

func validateMultiBetRequest(req *MultiBetRequest) error {
	if req.TicketID == "" {
		return fmt.Errorf("ticketId is required")
	}
	if len(req.Bets) == 0 {
		return fmt.Errorf("at least one bet is required")
	}
	for i, bet := range req.Bets {
		if bet.Type == "" {
			return fmt.Errorf("bet[%d].type is required", i)
		}
		if len(bet.Selections) == 0 {
			return fmt.Errorf("bet[%d] must have at least one selection", i)
		}
		for j, sel := range bet.Selections {
			if err := validateSelectionRequest(&sel); err != nil {
				return fmt.Errorf("bet[%d].selection[%d]: %w", i, j, err)
			}
		}
		if err := validateStakeRequest(&bet.Stake); err != nil {
			return fmt.Errorf("bet[%d].stake: %w", i, err)
		}
	}
	return nil
}

func validateSelectionRequest(sel *SelectionRequest) error {
	if sel.ProductID == "" {
		return fmt.Errorf("productId is required")
	}
	if sel.EventID == "" {
		return fmt.Errorf("eventId is required")
	}
	if sel.MarketID == "" {
		return fmt.Errorf("marketId is required")
	}
	if sel.OutcomeID == "" {
		return fmt.Errorf("outcomeId is required")
	}
	if sel.Odds == "" {
		return fmt.Errorf("odds is required")
	}
	// Validate odds is a valid number
	if odds, err := strconv.ParseFloat(sel.Odds, 64); err != nil || odds <= 0 {
		return fmt.Errorf("odds must be a valid number greater than 0")
	}
	return nil
}

func validateStakeRequest(stake *StakeRequest) error {
	if stake.Type == "" {
		return fmt.Errorf("type is required")
	}
	if stake.Type != "cash" && stake.Type != "free" {
		return fmt.Errorf("type must be 'cash' or 'free'")
	}
	if stake.Currency == "" {
		return fmt.Errorf("currency is required")
	}
	if stake.Amount == "" {
		return fmt.Errorf("amount is required")
	}
	// Validate amount is a valid number
	if amount, err := strconv.ParseFloat(stake.Amount, 64); err != nil || amount <= 0 {
		return fmt.Errorf("amount must be a valid number greater than 0")
	}
	if stake.Mode == "" {
		return fmt.Errorf("mode is required")
	}
	if stake.Mode != "total" && stake.Mode != "unit" {
		return fmt.Errorf("mode must be 'total' or 'unit'")
	}
	return nil
}

// Conversion functions

func convertSelectionRequest(req SelectionRequest) models.Selection {
	return models.NewSelection(
		req.ProductID,
		req.EventID,
		req.MarketID,
		req.OutcomeID,
		req.Odds,
		req.Specifiers,
	)
}

func convertStakeRequest(req StakeRequest) models.Stake {
	return models.NewStake(
		req.Type,
		req.Currency,
		req.Amount,
		req.Mode,
	)
}

func convertContextRequest(req *ContextRequest, cfg *config.Config) *models.Context {
	var channel *models.Channel
	
	if req.Channel != nil {
		channelType := req.Channel.Type
		if channelType == "" {
			channelType = "internet"
		}
		
		language := req.Channel.Lang
		if language == "" {
			language = "EN"
		}
		
		channel = &models.Channel{
			Type: channelType,
			Lang: language,
		}
	} else {
		// Default channel if not provided
		channel = &models.Channel{
			Type: "internet",
			Lang: "EN",
		}
	}
	
	var limitID int64
	if cfg.LimitID != "" {
		limitID, _ = strconv.ParseInt(cfg.LimitID, 10, 64)
	}
	
	return &models.Context{
		Channel: channel,
		IP:      req.IP,
		LimitID: limitID,
	}
}

func getDefaultContext(cfg *config.Config) *models.Context {
	var limitID int64
	if cfg.LimitID != "" {
		limitID, _ = strconv.ParseInt(cfg.LimitID, 10, 64)
	}
	
	return &models.Context{
		Channel: &models.Channel{
			Type: "internet",
			Lang: "EN",
		},
		LimitID: limitID,
	}
}
