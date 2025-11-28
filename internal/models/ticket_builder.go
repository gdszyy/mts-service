package models

import (
	"fmt"
	"time"
)

// TicketBuilder helps construct MTS ticket requests
type TicketBuilder struct {
	operatorID int64
	ticketID   string
	bets       []Bet
	context    *Context
}

// NewTicketBuilder creates a new ticket builder
func NewTicketBuilder(operatorID int64, ticketID string) *TicketBuilder {
	return &TicketBuilder{
		operatorID: operatorID,
		ticketID:   ticketID,
		bets:       []Bet{},
	}
}

// SetContext sets the context for the ticket
func (tb *TicketBuilder) SetContext(ctx *Context) *TicketBuilder {
	tb.context = ctx
	return tb
}

// AddSingleBet adds a single bet to the ticket
func (tb *TicketBuilder) AddSingleBet(selection Selection, stake Stake) *TicketBuilder {
	bet := Bet{
		Selections: []Selection{selection},
		Stake:      []Stake{stake},
	}
	tb.bets = append(tb.bets, bet)
	return tb
}

// AddAccumulatorBet adds an accumulator bet (multiple selections, all must win)
func (tb *TicketBuilder) AddAccumulatorBet(selections []Selection, stake Stake) *TicketBuilder {
	if len(selections) < 2 {
		panic("accumulator requires at least 2 selections")
	}
	bet := Bet{
		Selections: selections,
		Stake:      []Stake{stake},
	}
	tb.bets = append(tb.bets, bet)
	return tb
}

// AddSystemBet adds a system bet (e.g., 2/3, 3/5)
// size: array of combination sizes (e.g., [2] for doubles, [2,3] for doubles and trebles)
// selections: the selections to combine
// stake: unit stake (mode should be "unit" for system bets)
func (tb *TicketBuilder) AddSystemBet(size []int, selections []Selection, stake Stake) *TicketBuilder {
	if len(selections) < 2 {
		panic("system bet requires at least 2 selections")
	}
	
	// Validate size values
	for _, s := range size {
		if s < 1 || s > len(selections) {
			panic(fmt.Sprintf("invalid size %d for %d selections", s, len(selections)))
		}
	}
	
	systemSelection := Selection{
		Type:       "system",
		Size:       size,
		Selections: selections,
	}
	
	bet := Bet{
		Selections: []Selection{systemSelection},
		Stake:      []Stake{stake},
	}
	tb.bets = append(tb.bets, bet)
	return tb
}

// AddBankerSystemBet adds a system bet with banker selections
// This creates a proper MTS banker structure with nested system selections:
// - Outer system: combines banker system and main system
// - Banker system: type="system", size=[1], contains banker selections
// - Main system: type="system", size=size, contains non-banker selections
// 
// bankers: selections that must be in every combination
// size: array of combination sizes for non-banker selections
// selections: non-banker selections to combine
// stake: unit stake
func (tb *TicketBuilder) AddBankerSystemBet(bankers []Selection, size []int, selections []Selection, stake Stake) *TicketBuilder {
	if len(bankers) < 1 {
		panic("banker system bet requires at least 1 banker")
	}
	if len(selections) < 1 {
		panic("banker system bet requires at least 1 non-banker selection")
	}
	
	// Validate size values
	for _, s := range size {
		if s < 1 || s > len(selections) {
			panic(fmt.Sprintf("invalid size %d for %d non-banker selections", s, len(selections)))
		}
	}
	
	// Create banker system (size=[1] means all bankers must be included)
	bankerSystem := Selection{
		Type:       "system",
		Size:       []int{1},
		Selections: bankers,
	}
	
	// Create main system for non-banker selections
	mainSystem := Selection{
		Type:       "system",
		Size:       size,
		Selections: selections,
	}
	
	// Calculate outer system size: banker count + each size value
	outerSize := make([]int, len(size))
	for i, s := range size {
		outerSize[i] = len(bankers) + s
	}
	
	// Create outer system that combines banker and main systems
	rootSystem := Selection{
		Type:       "system",
		Size:       outerSize,
		Selections: []Selection{bankerSystem, mainSystem},
	}
	
	bet := Bet{
		Selections: []Selection{rootSystem},
		Stake:      []Stake{stake},
	}
	tb.bets = append(tb.bets, bet)
	return tb
}

// AddTrixieBet adds a Trixie bet (3 selections: 3 doubles + 1 treble = 4 bets)
func (tb *TicketBuilder) AddTrixieBet(selections []Selection, stake Stake) *TicketBuilder {
	if len(selections) != 3 {
		panic("trixie requires exactly 3 selections")
	}
	return tb.AddSystemBet([]int{2, 3}, selections, stake)
}

// AddPatentBet adds a Patent bet (3 selections: 3 singles + 3 doubles + 1 treble = 7 bets)
func (tb *TicketBuilder) AddPatentBet(selections []Selection, stake Stake) *TicketBuilder {
	if len(selections) != 3 {
		panic("patent requires exactly 3 selections")
	}
	return tb.AddSystemBet([]int{1, 2, 3}, selections, stake)
}

// AddYankeeBet adds a Yankee bet (4 selections: 6 doubles + 4 trebles + 1 four-fold = 11 bets)
func (tb *TicketBuilder) AddYankeeBet(selections []Selection, stake Stake) *TicketBuilder {
	if len(selections) != 4 {
		panic("yankee requires exactly 4 selections")
	}
	return tb.AddSystemBet([]int{2, 3, 4}, selections, stake)
}

// AddLucky15Bet adds a Lucky 15 bet (4 selections: 4 singles + 6 doubles + 4 trebles + 1 four-fold = 15 bets)
func (tb *TicketBuilder) AddLucky15Bet(selections []Selection, stake Stake) *TicketBuilder {
	if len(selections) != 4 {
		panic("lucky 15 requires exactly 4 selections")
	}
	return tb.AddSystemBet([]int{1, 2, 3, 4}, selections, stake)
}

// AddSuperYankeeBet adds a Super Yankee/Canadian bet (5 selections: 10 doubles + 10 trebles + 5 four-folds + 1 five-fold = 26 bets)
func (tb *TicketBuilder) AddSuperYankeeBet(selections []Selection, stake Stake) *TicketBuilder {
	if len(selections) != 5 {
		panic("super yankee requires exactly 5 selections")
	}
	return tb.AddSystemBet([]int{2, 3, 4, 5}, selections, stake)
}

// AddLucky31Bet adds a Lucky 31 bet (5 selections: 5 singles + 10 doubles + 10 trebles + 5 four-folds + 1 five-fold = 31 bets)
func (tb *TicketBuilder) AddLucky31Bet(selections []Selection, stake Stake) *TicketBuilder {
	if len(selections) != 5 {
		panic("lucky 31 requires exactly 5 selections")
	}
	return tb.AddSystemBet([]int{1, 2, 3, 4, 5}, selections, stake)
}

// AddHeinzBet adds a Heinz bet (6 selections: 15 doubles + 20 trebles + 15 four-folds + 6 five-folds + 1 six-fold = 57 bets)
func (tb *TicketBuilder) AddHeinzBet(selections []Selection, stake Stake) *TicketBuilder {
	if len(selections) != 6 {
		panic("heinz requires exactly 6 selections")
	}
	return tb.AddSystemBet([]int{2, 3, 4, 5, 6}, selections, stake)
}

// AddLucky63Bet adds a Lucky 63 bet (6 selections: 6 singles + 15 doubles + 20 trebles + 15 four-folds + 6 five-folds + 1 six-fold = 63 bets)
func (tb *TicketBuilder) AddLucky63Bet(selections []Selection, stake Stake) *TicketBuilder {
	if len(selections) != 6 {
		panic("lucky 63 requires exactly 6 selections")
	}
	return tb.AddSystemBet([]int{1, 2, 3, 4, 5, 6}, selections, stake)
}

// AddSuperHeinzBet adds a Super Heinz bet (7 selections: 21 doubles + 35 trebles + 35 four-folds + 21 five-folds + 7 six-folds + 1 seven-fold = 120 bets)
func (tb *TicketBuilder) AddSuperHeinzBet(selections []Selection, stake Stake) *TicketBuilder {
	if len(selections) != 7 {
		panic("super heinz requires exactly 7 selections")
	}
	return tb.AddSystemBet([]int{2, 3, 4, 5, 6, 7}, selections, stake)
}

// AddGoliathBet adds a Goliath bet (8 selections: 28 doubles + 56 trebles + 70 four-folds + 56 five-folds + 28 six-folds + 8 seven-folds + 1 eight-fold = 247 bets)
func (tb *TicketBuilder) AddGoliathBet(selections []Selection, stake Stake) *TicketBuilder {
	if len(selections) != 8 {
		panic("goliath requires exactly 8 selections")
	}
	return tb.AddSystemBet([]int{2, 3, 4, 5, 6, 7, 8}, selections, stake)
}

// Build constructs the final TicketRequest
func (tb *TicketBuilder) Build(correlationID string) *TicketRequest {
	if len(tb.bets) == 0 {
		panic("ticket must contain at least one bet")
	}
	
	return &TicketRequest{
		OperatorID:    tb.operatorID,
		CorrelationID: correlationID,
		TimestampUTC:  time.Now().UnixMilli(),
		Operation:     "ticket-placement",
		Version:       "3.0",
		Content: TicketContent{
			Type:     "ticket",
			TicketID: tb.ticketID,
			Bets:     tb.bets,
			Context:  tb.context,
		},
	}
}

// NewSelection creates a new standard selection
// odds can be either float64 or string
func NewSelection(productID, eventID, marketID, outcomeID string, odds interface{}, specifiers ...string) Selection {
	spec := ""
	if len(specifiers) > 0 {
		spec = specifiers[0]
	}
	
	// Convert odds to string
	var oddsStr string
	switch v := odds.(type) {
	case float64:
		oddsStr = fmt.Sprintf("%.2f", v)
	case string:
		oddsStr = v
	default:
		panic(fmt.Sprintf("odds must be float64 or string, got %T", odds))
	}
	
	return Selection{
		Type:       "uf",
		ProductID:  productID,
		EventID:    eventID,
		MarketID:   marketID,
		OutcomeID:  outcomeID,
		Specifiers: spec,
		Odds: &Odds{
			Type:  "decimal",
			Value: oddsStr,
		},
	}
}

// NewStake creates a new stake object
// amount can be either float64 or string
func NewStake(stakeType, currency string, amount interface{}, mode string) Stake {
	// Convert amount to string
	var amountStr string
	switch v := amount.(type) {
	case float64:
		amountStr = fmt.Sprintf("%.2f", v)
	case string:
		amountStr = v
	default:
		panic(fmt.Sprintf("amount must be float64 or string, got %T", amount))
	}
	
	return Stake{
		Type:     stakeType,
		Currency: currency,
		Amount:   amountStr,
		Mode:     mode,
	}
}
