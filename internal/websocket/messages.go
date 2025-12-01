package websocket

import "time"

// MessageType represents the type of WebSocket message
type MessageType string

const (
	// Client to Server
	MessageTypePlaceBet       MessageType = "place_bet"
	MessageTypeQueryBetStatus MessageType = "query_bet_status"
	MessageTypePing           MessageType = "ping"

	// Server to Client
	MessageTypeConnectionEstablished MessageType = "connection_established"
	MessageTypeConnectionRejected    MessageType = "connection_rejected"
	MessageTypeBetReceived           MessageType = "bet_received"
	MessageTypeBetResult             MessageType = "bet_result"
	MessageTypeBetPartialResult      MessageType = "bet_partial_result"
	MessageTypeBetTimeout            MessageType = "bet_timeout"
	MessageTypeBetResultDelayed      MessageType = "bet_result_delayed"
	MessageTypeBetStatus             MessageType = "bet_status"
	MessageTypeBetError              MessageType = "bet_error"
	MessageTypePong                  MessageType = "pong"
)

// BaseMessage is the base structure for all WebSocket messages
type BaseMessage struct {
	Type      MessageType `json:"type"`
	Timestamp time.Time   `json:"timestamp"`
}

// PlaceBetRequest represents a bet placement request from client
type PlaceBetRequest struct {
	BaseMessage
	RequestID string                 `json:"requestId"`
	BetType   string                 `json:"betType"` // single, multi, accumulator, system, banker
	Payload   map[string]interface{} `json:"payload"`
}

// QueryBetStatusRequest represents a bet status query from client
type QueryBetStatusRequest struct {
	BaseMessage
	TicketID string `json:"ticketId"`
}

// PingMessage represents a heartbeat ping from client
type PingMessage struct {
	BaseMessage
}

// ConnectionEstablishedResponse sent when connection is successfully established
type ConnectionEstablishedResponse struct {
	BaseMessage
	UserID  string `json:"userId"`
	Message string `json:"message"`
}

// ConnectionRejectedResponse sent when connection is rejected
type ConnectionRejectedResponse struct {
	BaseMessage
	Reason string `json:"reason"`
}

// BetReceivedResponse sent immediately after receiving a bet request
type BetReceivedResponse struct {
	BaseMessage
	RequestID string   `json:"requestId"`
	TicketID  string   `json:"ticketId,omitempty"`
	TicketIDs []string `json:"ticketIds,omitempty"` // For multi bets
}

// BetResultResponse sent when MTS returns the final result
type BetResultResponse struct {
	BaseMessage
	RequestID string                 `json:"requestId"`
	TicketID  string                 `json:"ticketId,omitempty"`
	Status    string                 `json:"status"` // accepted, rejected
	Details   map[string]interface{} `json:"details"`
	Summary   *BetSummary            `json:"summary,omitempty"` // For multi bets
}

// BetPartialResultResponse sent for each completed bet in a multi-bet request
type BetPartialResultResponse struct {
	BaseMessage
	RequestID string                 `json:"requestId"`
	Completed string                 `json:"completed"` // e.g., "2/3"
	TicketID  string                 `json:"ticketId"`
	Status    string                 `json:"status"`
	Details   map[string]interface{} `json:"details"`
}

// BetTimeoutResponse sent when MTS doesn't respond within timeout period
type BetTimeoutResponse struct {
	BaseMessage
	RequestID string `json:"requestId"`
	TicketID  string `json:"ticketId"`
	Message   string `json:"message"`
}

// BetResultDelayedResponse sent when delayed result arrives after timeout
type BetResultDelayedResponse struct {
	BaseMessage
	TicketID string                 `json:"ticketId"`
	Status   string                 `json:"status"`
	Details  map[string]interface{} `json:"details"`
}

// BetStatusResponse sent in response to status query
type BetStatusResponse struct {
	BaseMessage
	TicketID string                 `json:"ticketId"`
	Status   string                 `json:"status"` // accepted, rejected, pending, not_found
	Details  map[string]interface{} `json:"details,omitempty"`
}

// BetErrorResponse sent when there's an error processing the bet
type BetErrorResponse struct {
	BaseMessage
	RequestID string                 `json:"requestId,omitempty"`
	Error     string                 `json:"error"`
	Details   map[string]interface{} `json:"details,omitempty"`
}

// PongMessage sent in response to ping
type PongMessage struct {
	BaseMessage
}

// BetSummary contains summary information for multi-bet results
type BetSummary struct {
	Total    int                      `json:"total"`
	Accepted int                      `json:"accepted"`
	Rejected int                      `json:"rejected"`
	Details  []map[string]interface{} `json:"details"`
}
