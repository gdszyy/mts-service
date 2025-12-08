package models

// CashoutRequest represents a cashout request conforming to MTS Transaction 3.0 API
type CashoutRequest struct {
	OperatorID    int64          `json:"operatorId"`
	CorrelationID string         `json:"correlationId"`
	TimestampUTC  int64          `json:"timestampUtc"`
	Operation     string         `json:"operation"` // "cashout-inform", "cashout-build", or "cashout-placement"
	Version       string         `json:"version"`   // Should be "3.0"
	Content       CashoutContent `json:"content"`
}

// CashoutContent represents the content of a cashout request
type CashoutContent struct {
	Type       string             `json:"type"`       // Should be "cashout-inform", "cashout-build", or "cashout-placement"
	Cashout    CashoutInfo        `json:"cashout"`
	Validation *CashoutValidation `json:"validation,omitempty"` // Only for cashout-inform
}

// CashoutInfo represents cashout information
type CashoutInfo struct {
	Type      string        `json:"type"`      // Should be "cashout"
	CashoutID string        `json:"cashoutId"` // Unique cashout identifier
	Details   CashoutDetail `json:"details"`
}

// CashoutDetail represents cashout details
type CashoutDetail struct {
	Type            string         `json:"type"`            // "ticket", "ticket-partial", "bet", "bet-partial"
	TicketID        string         `json:"ticketId"`        // Original ticket ID
	TicketSignature string         `json:"ticketSignature"` // Signature from original ticket response
	Code            int            `json:"code"`            // Cashout reason code (e.g., 100, 101)
	Percentage      string         `json:"percentage,omitempty"` // For partial cashout (e.g., "0.5" for 50%)
	BetID           string         `json:"betId,omitempty"`      // For bet-level cashout
	Payout          []CashoutPayout `json:"payout"`
}

// CashoutPayout represents payout information for cashout
type CashoutPayout struct {
	Type     string `json:"type"`     // "cash" or "free"
	Currency string `json:"currency"` // Currency code (e.g., "EUR")
	Amount   string `json:"amount"`   // Amount as string
	// Source field is not supported in MTS schema
	// TraceID  string `json:"traceId,omitempty"` // Optional trace ID
}

// CashoutValidation represents validation information for cashout-inform
type CashoutValidation struct {
	Code    int    `json:"code"`    // Validation code (e.g., 1100)
	Message string `json:"message"` // Validation message
}

// CashoutResponse represents the response from MTS for cashout requests
type CashoutResponse struct {
	OperatorID    int64                  `json:"operatorId"`
	Content       CashoutResponseContent `json:"content"`
	CorrelationID string                 `json:"correlationId"`
	TimestampUTC  int64                  `json:"timestampUtc"`
	Operation     string                 `json:"operation"` // "cashout-inform", "cashout-build", or "cashout-placement"
	Version       string                 `json:"version"`
}

// CashoutResponseContent represents the content of a cashout response
type CashoutResponseContent struct {
	Type      string `json:"type"`      // "cashout-inform-reply", "cashout-build-reply", or "cashout-placement-reply"
	CashoutID string `json:"cashoutId"` // Cashout ID from request
	Signature string `json:"signature"` // Server signature for acknowledgement
	Status    string `json:"status"`    // "accepted" or "rejected"
	TicketID  string `json:"ticketId"`  // Original ticket ID
	Code      int    `json:"code"`      // Response code (0=success, negative=error)
	Message   string `json:"message,omitempty"` // Response message

	// Fields for cashout-build-reply
	LTD                    *LTDInfo                `json:"ltd,omitempty"`
	Cashout                *CashoutAmountInfo      `json:"cashout,omitempty"`
	EndCustomerSuggestions *EndCustomerSuggestions `json:"endCustomerSuggestions,omitempty"`
	ChannelSuggestions     *ChannelSuggestions     `json:"channelSuggestions,omitempty"`
	BetDetails             []CashoutBetDetail      `json:"betDetails,omitempty"`
	ExchangeRate           []ExchangeRate          `json:"exchangeRate,omitempty"`
}

// LTDInfo represents Late Ticket Detection information
type LTDInfo struct {
	ModelSuggestedLTD       string `json:"modelSuggestedLtd,omitempty"`
	ConfiguredLTD           int    `json:"configuredLtd,omitempty"`
	SuggestedLTD            int    `json:"suggestedLtd,omitempty"`
	AccountLbsLtdOffset     int    `json:"accountLbsLtdOffset,omitempty"`
	LiveSelectionLtdOffset  int    `json:"liveSelectionLtdOffset,omitempty"`
	AppliedLTD              int    `json:"appliedLtd,omitempty"`
}

// CashoutAmountInfo represents cashout amount information
type CashoutAmountInfo struct {
	CashoutType string            `json:"cashoutType"` // "ticket", "bet"
	CashoutID   string            `json:"cashoutId"`
	MaxPayout   []CashoutPayout   `json:"maxPayout,omitempty"`
	FairCashout []CashoutPayout   `json:"fairCashout,omitempty"`
	Cashout     []CashoutPayout   `json:"cashout,omitempty"`
}

// EndCustomerSuggestions represents end customer suggestions
type EndCustomerSuggestions struct {
	EndCustomer          *EndCustomer `json:"endCustomer,omitempty"`
	AppliedConfidence    string       `json:"appliedConfidence,omitempty"`
	SuggestedConfidence  string       `json:"suggestedConfidence,omitempty"`
	SuggestedLateBetScore string      `json:"suggestedLateBetScore,omitempty"`
	SuggestedMarkerScore string       `json:"suggestedMarkerScore,omitempty"`
	SuggestedBotScore    string       `json:"suggestedBotScore,omitempty"`
}

// EndCustomer represents end customer information
type EndCustomer struct {
	ID         string `json:"id"`
	Confidence string `json:"confidence"`
}

// ChannelSuggestions represents channel suggestions
type ChannelSuggestions struct {
	Channel              *Channel `json:"channel,omitempty"`
	AppliedConfidence    string   `json:"appliedConfidence,omitempty"`
	SuggestedConfidence  string   `json:"suggestedConfidence,omitempty"`
	SuggestedLateBetScore string  `json:"suggestedLateBetScore,omitempty"`
}

// CashoutBetDetail represents bet details in cashout response
type CashoutBetDetail struct {
	BetID             string                      `json:"betId"`
	SelectionDetails  []CashoutSelectionDetail    `json:"selectionDetails,omitempty"`
	Payout            []CashoutPayout             `json:"payout,omitempty"`
	SettledPercentage string                      `json:"settledPercentage,omitempty"`
}

// CashoutSelectionDetail represents selection details in cashout response
type CashoutSelectionDetail struct {
	Selection            Selection            `json:"selection"`
	AppliedEventRating   int                  `json:"appliedEventRating,omitempty"`
	SuggestedEventRating int                  `json:"suggestedEventRating,omitempty"`
	ConfiguredLTD        int                  `json:"configuredLtd,omitempty"`
	SuggestedLTD         int                  `json:"suggestedLtd,omitempty"`
	AppliedMarketFactor  string               `json:"appliedMarketFactor,omitempty"`
	CurrentProbability   *CurrentProbability  `json:"currentProbability,omitempty"`
	CurrentResult        *CurrentResult       `json:"currentResult,omitempty"`
}

// CurrentProbability represents current probability information
type CurrentProbability struct {
	Type     string `json:"type"`     // "push", "normal"
	Win      string `json:"win,omitempty"`
	Refund   string `json:"refund,omitempty"`
	HalfWin  string `json:"halfWin,omitempty"`
	HalfLose string `json:"halfLose,omitempty"`
}

// CurrentResult represents current result information
type CurrentResult struct {
	Type           string `json:"type"` // "unsettled", "win", "lose", "void"
	DeadHeatFactor string `json:"deadHeatFactor,omitempty"`
	VoidFactor     string `json:"voidFactor,omitempty"`
}

// CashoutAck represents an acknowledgement for a cashout response
type CashoutAck struct {
	OperatorID    int64              `json:"operatorId"`
	CorrelationID string             `json:"correlationId"`
	TimestampUTC  int64              `json:"timestampUtc"`
	Operation     string             `json:"operation"` // "cashout-inform-ack", "cashout-build-ack", or "cashout-placement-ack"
	Version       string             `json:"version"`   // Should be "3.0"
	Content       CashoutAckContent  `json:"content"`
}

// CashoutAckContent represents the content of a cashout acknowledgement
type CashoutAckContent struct {
	Type              string `json:"type"`              // "cashout-inform-ack", "cashout-build-ack", or "cashout-placement-ack"
	CashoutID         string `json:"cashoutId"`         // Cashout ID from response
	CashoutSignature  string `json:"cashoutSignature"`  // Signature from cashout response
	Acknowledged      bool   `json:"acknowledged"`      // Should be true
}
