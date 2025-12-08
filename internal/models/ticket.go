package models

// TicketRequest represents a ticket placement request conforming to MTS Transaction 3.0 API standard
type TicketRequest struct {
	OperatorID    int64         `json:"operatorId"`    // Operator ID provided by Sportradar
	CorrelationID string        `json:"correlationId"` // Client-defined string for request-response pairing
	TimestampUTC  int64         `json:"timestampUtc"`  // Client submission timestamp in Unix milliseconds
	Operation     string        `json:"operation"`     // Should be "ticket-placement" for ticket placement requests
	Version       string        `json:"version"`       // Protocol version, should be "3.0"
	Content       TicketContent `json:"content"`       // Message body containing transaction details
}

// TicketContent represents the content of a ticket placement request
type TicketContent struct {
	Type     string   `json:"type"`               // Content type, should be "ticket"
	TicketID string   `json:"ticketId"`           // Client-defined ticket ID for unique identification
	Bets     []Bet    `json:"bets"`               // Array of bets, must contain at least one bet
	Context  *Context `json:"context,omitempty"` // Optional transaction context information
}

// Context represents transaction context information
type Context struct {
	Channel     *Channel     `json:"channel,omitempty"`     // Channel information
	IP          string       `json:"ip,omitempty"`          // IP address
	EndCustomer *EndCustomer `json:"endCustomer,omitempty"` // End customer information
	LimitID     int64        `json:"limitId,omitempty"`     // Limit ID
}

// Channel represents channel information within context
type Channel struct {
	Type string `json:"type"` // Channel type (e.g., "internet", "agent", "mobile")
	Lang string `json:"lang"` // Language code (e.g., "EN")
}

// Bet represents a single bet within a ticket
type Bet struct {
	Selections []Selection `json:"selections"` // Array of selections, must contain at least one
	Stake      []Stake     `json:"stake"`      // Array of stake objects, must contain at least one
}

// Selection represents a single selection within a bet
// For standard selections: type="uf", "external", or "uf-custom-bet"
// For system bets: type="system" with nested selections
type Selection struct {
	// Common fields
	Type string `json:"type"` // Selection type: "uf", "external", "uf-custom-bet", or "system"

	// Fields for standard selections (type="uf", "external", "uf-custom-bet")
	ProductID  string `json:"productId,omitempty"`  // Product ID (e.g., "3")
	EventID    string `json:"eventId,omitempty"`    // Event ID (e.g., "sr:match:14950205")
	MarketID   string `json:"marketId,omitempty"`   // Market ID (e.g., "14")
	OutcomeID  string `json:"outcomeId,omitempty"`  // Outcome ID (e.g., "1712")
	Specifiers string `json:"specifiers,omitempty"` // Optional specifiers (e.g., "hcp=1:0")
	Odds       *Odds  `json:"odds,omitempty"`       // Odds object (not used for system type)

	// Fields for system bets (type="system")
	Size       []int       `json:"size,omitempty"`       // Array of combination sizes (e.g., [2,3] for doubles and trebles)
	Selections []Selection `json:"selections,omitempty"` // Nested selections for system bets
}

// Odds represents the odds for a selection
type Odds struct {
	Type  string `json:"type"`  // Odds type (e.g., "decimal")
	Value string `json:"value"` // Odds value as a string (e.g., "1.59")
}

// Stake represents a stake object within a bet
type Stake struct {
	Type     string `json:"type"`     // Stake type (e.g., "cash", "free")
	Currency string `json:"currency"` // Currency code (e.g., "EUR", "mBTC")
	Amount   string `json:"amount"`   // Amount as a string (e.g., "10")
	Mode     string `json:"mode,omitempty"` // Optional mode (e.g., "total")
}

// ExchangeRate represents currency exchange rate information
type ExchangeRate struct {
	FromCurrency string `json:"fromCurrency"` // Original currency (e.g., "EUR", "USD")
	ToCurrency   string `json:"toCurrency"`   // Target currency (e.g., "EUR")
	Rate         string `json:"rate"`         // Exchange rate as string (e.g., "1.00000000")
}

// TicketResponse represents the response from MTS
type TicketResponse struct {
	OperatorID    int64                  `json:"operatorId,omitempty"`
	Operation     string                 `json:"operation"`
	Content       TicketResponseContent  `json:"content"`
	CorrelationID string                 `json:"correlationId"`
	TimestampUTC  int64                  `json:"timestampUtc"`
	Version       string                 `json:"version"`
}

// TicketResponseContent represents the content of a ticket response
type TicketResponseContent struct {
	Type         string         `json:"type"`
	TicketID     string         `json:"ticketId"`
	Status       string         `json:"status"`
	Code         int            `json:"code,omitempty"`         // Validation code (0=success, negative=error)
	Message      string         `json:"message,omitempty"`      // Readable validation message
	Reason       *Reason        `json:"reason,omitempty"`       // Deprecated, use Code/Message instead
	BetDetails   []BetDetail    `json:"betDetails,omitempty"`
	Signature    string         `json:"signature"`              // Server-returned signature for acknowledgement
	ExchangeRate []ExchangeRate `json:"exchangeRate,omitempty"` // Array of exchange rates
}

// Reason represents the reason for ticket rejection
type Reason struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// BetDetail represents details of a bet in the response
type BetDetail struct {
	BetID            string            `json:"betId,omitempty"`
	Code             int               `json:"code,omitempty"`    // Bet-level validation code
	Message          string            `json:"message,omitempty"` // Bet-level validation message
	Status           string            `json:"status,omitempty"`
	Reason           *Reason           `json:"reason,omitempty"` // Deprecated
	AlternativeStake *AlternativeStake `json:"alternativeStake,omitempty"`
	SelectionDetails []SelectionDetail `json:"selectionDetails,omitempty"`
}

// AlternativeStake represents an alternative stake suggestion
type AlternativeStake struct {
	Stake int64 `json:"stake"`
}

// SelectionDetail represents details of a selection in the response
type SelectionDetail struct {
	Selection Selection `json:"selection"`         // Complete selection object
	Code      int       `json:"code,omitempty"`    // Selection-level validation code
	Message   string    `json:"message,omitempty"` // Selection-level validation message
	Reason    *Reason   `json:"reason,omitempty"`  // Deprecated
}

// TicketAck represents a ticket acknowledgement message
type TicketAck struct {
	OperatorID    int64            `json:"operatorId"`
	CorrelationID string           `json:"correlationId"`
	TimestampUTC  int64            `json:"timestampUtc"`
	Operation     string           `json:"operation"`
	Version       string           `json:"version"`
	Content       TicketAckContent `json:"content"`
}

// TicketAckContent represents the content of a ticket acknowledgement
type TicketAckContent struct {
	Type                  string `json:"type"`
	TicketID              string `json:"ticketId"`
	Acknowledged          bool   `json:"acknowledged"`
	TicketSignature       string `json:"ticketSignature,omitempty"`
	CancellationSignature string `json:"cancellationSignature,omitempty"`
	CashoutSignature      string `json:"cashoutSignature,omitempty"`
	SettlementSignature   string `json:"settlementSignature,omitempty"`
}

// ErrorReplyContent represents the content of an error reply
type ErrorReplyContent struct {
	Type    string `json:"type"`
	Code    int    `json:"code"`
	Message string `json:"message"`
}
