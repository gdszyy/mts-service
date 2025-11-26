package api

// Common structures for all bet types

// SelectionRequest represents a selection in API requests
type SelectionRequest struct {
	ProductID  string  `json:"productId"`            // Product ID (e.g., "3")
	EventID    string  `json:"eventId"`              // Event ID (e.g., "sr:match:12345")
	MarketID   string  `json:"marketId"`             // Market ID (e.g., "1")
	OutcomeID  string  `json:"outcomeId"`            // Outcome ID (e.g., "1712")
	Odds       float64 `json:"odds"`                 // Odds as decimal (e.g., 2.50)
	Specifiers string  `json:"specifiers,omitempty"` // Optional specifiers (e.g., "hcp=1:0")
}

// StakeRequest represents stake information
type StakeRequest struct {
	Type     string  `json:"type"`     // "cash" or "free"
	Currency string  `json:"currency"` // Currency code (e.g., "EUR", "USD")
	Amount   float64 `json:"amount"`   // Amount as decimal
	Mode     string  `json:"mode"`     // "total" or "unit"
}

// ContextRequest represents context information
type ContextRequest struct {
	ChannelType string `json:"channelType,omitempty"` // "internet", "mobile", "agent"
	Language    string `json:"language,omitempty"`    // Language code (e.g., "EN")
	IP          string `json:"ip,omitempty"`          // Customer IP address
}

// SingleBetRequest represents a single bet request
type SingleBetRequest struct {
	TicketID  string           `json:"ticketId"`  // Unique ticket ID
	Selection SelectionRequest `json:"selection"` // The selection
	Stake     StakeRequest     `json:"stake"`     // Stake information
	Context   *ContextRequest  `json:"context,omitempty"`
}

// AccumulatorBetRequest represents an accumulator bet request
type AccumulatorBetRequest struct {
	TicketID   string             `json:"ticketId"`   // Unique ticket ID
	Selections []SelectionRequest `json:"selections"` // Multiple selections (all must win)
	Stake      StakeRequest       `json:"stake"`      // Stake information
	Context    *ContextRequest    `json:"context,omitempty"`
}

// SystemBetRequest represents a system bet request
type SystemBetRequest struct {
	TicketID   string             `json:"ticketId"`   // Unique ticket ID
	Size       []int              `json:"size"`       // Combination sizes (e.g., [2] for doubles, [2,3] for doubles and trebles)
	Selections []SelectionRequest `json:"selections"` // Selections to combine
	Stake      StakeRequest       `json:"stake"`      // Unit stake
	Context    *ContextRequest    `json:"context,omitempty"`
}

// BankerSystemBetRequest represents a banker system bet request
type BankerSystemBetRequest struct {
	TicketID   string             `json:"ticketId"`   // Unique ticket ID
	Bankers    []SelectionRequest `json:"bankers"`    // Banker selections (must be in every combination)
	Size       []int              `json:"size"`       // Combination sizes for non-banker selections
	Selections []SelectionRequest `json:"selections"` // Non-banker selections to combine
	Stake      StakeRequest       `json:"stake"`      // Unit stake
	Context    *ContextRequest    `json:"context,omitempty"`
}

// PresetSystemBetRequest represents a preset system bet request (Trixie, Yankee, etc.)
type PresetSystemBetRequest struct {
	TicketID   string             `json:"ticketId"`   // Unique ticket ID
	Type       string             `json:"type"`       // "trixie", "patent", "yankee", "lucky15", "lucky31", "lucky63", "super_yankee", "heinz", "super_heinz", "goliath"
	Selections []SelectionRequest `json:"selections"` // Selections (count must match type requirement)
	Stake      StakeRequest       `json:"stake"`      // Unit stake
	Context    *ContextRequest    `json:"context,omitempty"`
}

// MultiBetRequest represents a ticket with multiple bets
type MultiBetRequest struct {
	TicketID string          `json:"ticketId"` // Unique ticket ID
	Bets     []BetDefinition `json:"bets"`     // Multiple bets in one ticket
	Context  *ContextRequest `json:"context,omitempty"`
}

// BetDefinition represents a single bet in a multi-bet ticket
type BetDefinition struct {
	Type       string             `json:"type"`                 // "single", "accumulator", "system", "banker_system", or preset type
	Selections []SelectionRequest `json:"selections"`           // Selections for this bet
	Stake      StakeRequest       `json:"stake"`                // Stake for this bet
	Size       []int              `json:"size,omitempty"`       // For system bets
	Bankers    []SelectionRequest `json:"bankers,omitempty"`    // For banker system bets
}

// CashoutRequest represents a cashout request
type CashoutRequest struct {
	CashoutID       string  `json:"cashoutId"`       // Unique cashout identifier
	TicketID        string  `json:"ticketId"`        // Original ticket ID
	TicketSignature string  `json:"ticketSignature"` // Signature from original ticket response
	Type            string  `json:"type"`            // "ticket", "ticket-partial", "bet", "bet-partial"
	Code            int     `json:"code"`            // Cashout reason code (e.g., 100, 101)
	Percentage      float64 `json:"percentage,omitempty"` // For partial cashout (0.0-1.0)
	BetID           string  `json:"betId,omitempty"`      // For bet-level cashout
	Payout          []PayoutRequest `json:"payout"`  // Payout information
}

// PayoutRequest represents payout information
type PayoutRequest struct {
	Type     string  `json:"type"`     // "cash" or "free"
	Currency string  `json:"currency"` // Currency code
	Amount   float64 `json:"amount"`   // Amount as decimal
	Source   string  `json:"source,omitempty"` // "cashout" or "bonus"
}

// APIResponse represents a standard API response
type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Error   *APIError   `json:"error,omitempty"`
}

// APIError represents an error in API response
type APIError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Details string `json:"details,omitempty"`
}
