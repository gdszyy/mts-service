package models

// TicketRequest represents a ticket placement request
type TicketRequest struct {
	Operation     string        `json:"operation"`
	Content       TicketContent `json:"content"`
	CorrelationID string        `json:"correlationId"`
	TimestampUTC  int64         `json:"timestampUtc"`
	Version       string        `json:"version"`
}

// TicketContent represents the content of a ticket
type TicketContent struct {
	Type       string      `json:"type"`
	TicketID   string      `json:"ticketId"`
	Sender     Sender      `json:"sender"`
	Bets       []Bet       `json:"bets"`
	TotalStake int64       `json:"totalStake"`
	TestSource bool        `json:"testSource,omitempty"`
	OddsChange string      `json:"oddsChange,omitempty"`
}

// Sender represents the bookmaker sending the ticket
type Sender struct {
	Bookmaker    string      `json:"bookmaker"`
	Currency     string      `json:"currency"`
	Channel      string      `json:"channel,omitempty"`
	EndCustomer  EndCustomer `json:"endCustomer"`
	SuggestedCCF float64     `json:"suggestedCcf,omitempty"`
}

// EndCustomer represents the end customer placing the bet
type EndCustomer struct {
	ID         string `json:"id"`
	IP         string `json:"ip,omitempty"`
	LanguageID string `json:"languageId,omitempty"`
	DeviceID   string `json:"deviceId,omitempty"`
	Confidence int64  `json:"confidence,omitempty"`
}

// Bet represents a single bet within a ticket
type Bet struct {
	ID         string      `json:"id"`
	Stake      int64       `json:"stake"`
	Selections []Selection `json:"selections"`
	BetBonus   int64       `json:"betBonus,omitempty"`
	CustomBet  bool        `json:"customBet,omitempty"`
}

// Selection represents a single selection within a bet
type Selection struct {
	ID      string `json:"id"`
	EventID string `json:"eventId"`
	Odds    int    `json:"odds"`
	Banker  bool   `json:"banker,omitempty"`
}

// TicketResponse represents the response from MTS
type TicketResponse struct {
	Operation     string                `json:"operation"`
	Content       TicketResponseContent `json:"content"`
	CorrelationID string                `json:"correlationId"`
	TimestampUTC  int64                 `json:"timestampUtc"`
	Version       string                `json:"version"`
}

// TicketResponseContent represents the content of a ticket response
type TicketResponseContent struct {
	Type         string      `json:"type"`
	TicketID     string      `json:"ticketId"`
	Status       string      `json:"status"`
	Reason       *Reason     `json:"reason,omitempty"`
	BetDetails   []BetDetail `json:"betDetails,omitempty"`
	Signature    string      `json:"signature"`
	ExchangeRate float64     `json:"exchangeRate,omitempty"`
}

// Reason represents the reason for ticket rejection
type Reason struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

// BetDetail represents details of a bet in the response
type BetDetail struct {
	BetID            string            `json:"betId"`
	Status           string            `json:"status"`
	Reason           *Reason           `json:"reason,omitempty"`
	AlternativeStake *AlternativeStake `json:"alternativeStake,omitempty"`
	SelectionDetails []SelectionDetail `json:"selectionDetails,omitempty"`
}

// AlternativeStake represents an alternative stake suggestion
type AlternativeStake struct {
	Stake int64 `json:"stake"`
}

// SelectionDetail represents details of a selection in the response
type SelectionDetail struct {
	SelectionID string  `json:"selectionId"`
	Odds        int     `json:"odds"`
	Reason      *Reason `json:"reason,omitempty"`
}

// TicketAck represents a ticket acknowledgement message
type TicketAck struct {
	Operation     string           `json:"operation"`
	Content       TicketAckContent `json:"content"`
	CorrelationID string           `json:"correlationId"`
	TimestampUTC  int64            `json:"timestampUtc"`
	Version       string           `json:"version"`
}

// TicketAckContent represents the content of a ticket acknowledgement
type TicketAckContent struct {
	Type            string `json:"type"`
	TicketID        string `json:"ticketId"`
	TicketSignature string `json:"ticketSignature"`
}

