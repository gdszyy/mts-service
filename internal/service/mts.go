package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gdsZyy/mts-service/internal/config"
	"github.com/gdsZyy/mts-service/internal/models"
	"github.com/gorilla/websocket"
)

const (
	WriteWait              = 10 * time.Second
	PongWait               = 60 * time.Second
	PingPeriod             = 54 * time.Second
	MaxMessageSize         = 512 * 1024
	ConnectionRefreshTime  = 110 * time.Minute // Refresh connection every 110 minutes (within 100-120 minute window)
	GracefulShutdownTime   = 5 * time.Minute   // Wait up to 5 minutes for old connection to drain
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

// ConnectionState represents the state of a WebSocket connection
type ConnectionState struct {
	conn              *websocket.Conn
	connectedAt       time.Time
	isActive          bool // Whether this connection should accept new requests
	pendingResponses  int32 // Number of requests awaiting responses on this connection
	mu                sync.RWMutex
}

type MTSService struct {
	cfg          *config.Config
	wsURL        string
	wsAudience   string
	
	// Connection management
	activeConn    *ConnectionState
	oldConn       *ConnectionState
	connMu        sync.RWMutex
	
	token        *TokenResponse
	tokenExpiry  time.Time
	tokenMu      sync.RWMutex
	
	responses    map[string]chan *models.TicketResponse
	responseMu   sync.RWMutex
	
	// Idempotency: store sent messages and their responses
	sentMessages map[string]*models.TicketResponse // Key: JSON hash of the message
	sentMsgMu    sync.RWMutex
	
	ctx          context.Context
	cancel       context.CancelFunc
	connected    int32 // atomic flag for connection status
	reconnecting int32 // atomic flag for reconnection status
	httpClient   *http.Client
}

func NewMTSService(cfg *config.Config) *MTSService {
	wsURL := cfg.WSURL
	wsAudience := cfg.WSAudience
	if cfg.Production {
		wsURL = fmt.Sprintf("wss://%s", cfg.VirtualHost)
		wsAudience = "mbs-dp-production-wss"
	}

	ctx, cancel := context.WithCancel(context.Background())

	return &MTSService{
		cfg:          cfg,
		wsURL:        wsURL,
		wsAudience:   wsAudience,
		responses:    make(map[string]chan *models.TicketResponse),
		sentMessages: make(map[string]*models.TicketResponse),
		ctx:          ctx,
		cancel:       cancel,
		httpClient:   &http.Client{Timeout: 30 * time.Second},
	}
}

func (s *MTSService) Start() error {
	if err := s.connect(); err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	log.Println("MTS Service started successfully")
	
	// Start connection refresh monitor
	go s.connectionRefreshMonitor()
	
	return nil
}

func (s *MTSService) Stop() error {
	s.cancel()
	
	// Close active connection
	s.connMu.Lock()
	if s.activeConn != nil && s.activeConn.conn != nil {
		s.activeConn.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		s.activeConn.conn.Close()
		s.activeConn = nil
	}
	
	// Close old connection if still exists
	if s.oldConn != nil && s.oldConn.conn != nil {
		s.oldConn.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		s.oldConn.conn.Close()
		s.oldConn = nil
	}
	s.connMu.Unlock()
	
	atomic.StoreInt32(&s.connected, 0)
	log.Println("MTS Service stopped")
	return nil
}

func (s *MTSService) getToken() (string, error) {
	s.tokenMu.RLock()
	if s.token != nil && time.Now().Before(s.tokenExpiry) {
		token := s.token.AccessToken
		s.tokenMu.RUnlock()
		return token, nil
	}
	s.tokenMu.RUnlock()

	return s.refreshToken()
}

func (s *MTSService) refreshToken() (string, error) {
	s.tokenMu.Lock()
	defer s.tokenMu.Unlock()

	if s.token != nil && time.Now().Before(s.tokenExpiry) {
		return s.token.AccessToken, nil
	}

	data := url.Values{}
	data.Set("grant_type", "client_credentials")
	data.Set("client_id", s.cfg.ClientID)
	data.Set("client_secret", s.cfg.ClientSecret)
	data.Set("audience", s.wsAudience)

	req, err := http.NewRequest("POST", s.cfg.AuthURL, bytes.NewBufferString(data.Encode()))
	if err != nil {
		return "", fmt.Errorf("failed to create auth request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", "Sportradar-MTS-Client/1.0 (Go)")

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return "", fmt.Errorf("failed to execute auth request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("auth request failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp TokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return "", fmt.Errorf("failed to decode token response: %w", err)
	}

	s.token = &tokenResp
	s.tokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second).Add(-30 * time.Second)

	return tokenResp.AccessToken, nil
}

func (s *MTSService) connect() error {
	token, err := s.getToken()
	if err != nil {
		return fmt.Errorf("failed to get auth token: %w", err)
	}

	header := make(map[string][]string)
	header["Authorization"] = []string{"Bearer " + token}

	dialer := websocket.Dialer{
		HandshakeTimeout: 45 * time.Second,
	}

	conn, _, err := dialer.Dial(s.wsURL, header)
	if err != nil {
		return fmt.Errorf("failed to connect to WebSocket: %w", err)
	}

	newConnState := &ConnectionState{
		conn:        conn,
		connectedAt: time.Now(),
		isActive:    true,
	}

	s.connMu.Lock()
	s.activeConn = newConnState
	s.connMu.Unlock()

	atomic.StoreInt32(&s.connected, 1)

	conn.SetReadLimit(MaxMessageSize)
	conn.SetReadDeadline(time.Now().Add(PongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(PongWait))
		return nil
	})

	go s.readPump(newConnState)
	go s.pingPump(newConnState)

	log.Println("Connected to MTS WebSocket")

	// Send initialization message
	if err := s.sendInitializationMessage(); err != nil {
		return fmt.Errorf("failed to send initialization message: %w", err)
	}

	return nil
}

// connectionRefreshMonitor monitors connection age and initiates refresh when needed
func (s *MTSService) connectionRefreshMonitor() {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			s.connMu.RLock()
			activeConn := s.activeConn
			s.connMu.RUnlock()

			if activeConn != nil {
				age := time.Since(activeConn.connectedAt)
				if age >= ConnectionRefreshTime {
					log.Printf("Connection age (%v) reached refresh threshold, initiating smooth refresh", age)
					s.initiateConnectionRefresh()
				}
			}
		}
	}
}

// initiateConnectionRefresh implements Option 1: Smooth connection switch
// 1. Open new connection
// 2. Divert new traffic to new connection
// 3. Keep old connection alive until all responses received
// 4. Close old connection
func (s *MTSService) initiateConnectionRefresh() {
	log.Println("Initiating smooth connection refresh (Option 1)...")

	// Step 1: Open new connection
	if err := s.connect(); err != nil {
		log.Printf("Failed to open new connection during refresh: %v", err)
		return
	}

	// Step 2: Mark old connection as inactive (new traffic goes to new connection)
	s.connMu.Lock()
	if s.oldConn != nil && s.oldConn.conn != nil {
		// Close previous old connection if it still exists
		s.oldConn.conn.Close()
	}
	
	// Move current active connection to old connection
	s.oldConn = s.activeConn
	if s.oldConn != nil {
		s.oldConn.isActive = false
		log.Printf("Old connection marked as inactive. Pending responses: %d", atomic.LoadInt32(&s.oldConn.pendingResponses))
	}
	s.connMu.Unlock()

	// Step 3 & 4: Monitor old connection and close when all responses received
	go s.drainOldConnection()
}

// drainOldConnection waits for all pending responses on the old connection, then closes it
func (s *MTSService) drainOldConnection() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	timeout := time.NewTimer(GracefulShutdownTime)
	defer timeout.Stop()

	for {
		select {
		case <-s.ctx.Done():
			// Service is shutting down
			s.connMu.Lock()
			if s.oldConn != nil && s.oldConn.conn != nil {
				s.oldConn.conn.Close()
				s.oldConn = nil
			}
			s.connMu.Unlock()
			return

		case <-timeout.C:
			// Timeout reached, force close old connection
			log.Println("Graceful shutdown timeout reached, force closing old connection")
			s.connMu.Lock()
			if s.oldConn != nil && s.oldConn.conn != nil {
				s.oldConn.conn.Close()
				s.oldConn = nil
			}
			s.connMu.Unlock()
			return

		case <-ticker.C:
			s.connMu.RLock()
			oldConn := s.oldConn
			s.connMu.RUnlock()

			if oldConn == nil {
				return
			}

			pendingCount := atomic.LoadInt32(&oldConn.pendingResponses)
			if pendingCount == 0 {
				log.Println("All responses received on old connection, closing it")
				s.connMu.Lock()
				if s.oldConn != nil && s.oldConn.conn != nil {
					s.oldConn.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
					s.oldConn.conn.Close()
					s.oldConn = nil
				}
				s.connMu.Unlock()
				return
			}
		}
	}
}

// sendInitializationMessage sends initialization message after connection
func (s *MTSService) sendInitializationMessage() error {
	operatorID := s.cfg.OperatorID
	if operatorID == 0 {
		log.Println("Warning: OperatorID is not set in config. Using default 9985 for initialization message.")
		operatorID = 9985
	}

	initMsg := &models.TicketRequest{
		OperatorID:    operatorID,
		CorrelationID: fmt.Sprintf("init-%d", time.Now().UnixNano()),
		TimestampUTC:  time.Now().UnixMilli(),
		Operation:     "ticket-placement",
		Version:       "3.0",
		Content: models.TicketContent{
			Type:     "ticket",
			TicketID: fmt.Sprintf("init-ticket-%d", time.Now().UnixNano()),
			Bets: []models.Bet{
				{
					Selections: []models.Selection{
						{
							Type:       "uf",
							ProductID:  "3",
							EventID:    "sr:match:16470657",
							MarketID:   "534",
							OutcomeID:  "pre:outcometext:9919",
							Odds: models.Odds{
								Type:  "decimal",
								Value: "2.10",
							},
						},
					},
					Stake: []models.Stake{
							{
								Type:     "cash",
								Currency: "EUR",
								Amount:   "10.00000000", // 8 decimal places precision
								Mode:     "total", // Default mode for system bets
							},
					},
				},
			},
			Context: &models.Context{
				Channel: &models.Channel{
					Type: "agent",
					Lang: "EN",
				},
					LimitID: getLimitIDFromConfig(s.cfg),
			},
		},
	}

	log.Printf("Sending initialization message: %+v", initMsg)

	return s.sendMessage(initMsg)
}

func getLimitIDFromConfig(cfg *config.Config) int64 {
	if cfg.LimitID != "" {
		limitID, err := strconv.ParseInt(cfg.LimitID, 10, 64)
		if err != nil {
			log.Printf("Warning: Failed to parse LimitID from config '%s' for initialization message: %v, using 0.", cfg.LimitID, err)
			return 0
		}
		return limitID
	}
	// Fallback to a default or 0 if not set
	return 0
}

func (s *MTSService) readPump(connState *ConnectionState) {
	defer func() {
		connState.mu.Lock()
		if connState.conn != nil {
			connState.conn.Close()
		}
		connState.mu.Unlock()

		// Check if this was the active connection
		s.connMu.RLock()
		isActive := (s.activeConn == connState)
		s.connMu.RUnlock()

		if isActive {
			atomic.StoreInt32(&s.connected, 0)
			s.reconnect()
		}
	}()

	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			connState.mu.RLock()
			conn := connState.conn
			connState.mu.RUnlock()

			if conn == nil {
				return
			}

			_, message, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("WebSocket read error: %v", err)
				}
				return
			}

			s.handleMessage(message, connState)
		}
	}
}

func (s *MTSService) pingPump(connState *ConnectionState) {
	ticker := time.NewTicker(PingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			connState.mu.RLock()
			conn := connState.conn
			connState.mu.RUnlock()

			if conn == nil {
				continue
			}

			conn.SetWriteDeadline(time.Now().Add(WriteWait))
			if err := conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				log.Printf("Failed to send ping: %v", err)
				return
			}
		}
	}
}

func (s *MTSService) handleMessage(message []byte, connState *ConnectionState) {
	var response models.TicketResponse
	if err := json.Unmarshal(message, &response); err != nil {
		log.Printf("Failed to unmarshal response: %v. Message: %s", err, string(message))
		return
	}

	// Fill in OperatorID if missing
	if response.OperatorID == 0 {
		response.OperatorID = s.cfg.OperatorID
		if response.OperatorID == 0 {
			response.OperatorID = 9985
		}
	}

	// Decrement pending responses counter
	atomic.AddInt32(&connState.pendingResponses, -1)

	// Check if this is an error response
	if response.Content.Type == "error-reply" {
		var errorResponse struct {
			Content models.ErrorReplyContent `json:"content"`
		}
		if err := json.Unmarshal(message, &errorResponse); err == nil {
			log.Printf("MTS Error Reply received (CorrelationID: %s): Code=%d, Message=%s", response.CorrelationID, errorResponse.Content.Code, errorResponse.Content.Message)
		} else {
			log.Printf("MTS Error Reply received, but failed to parse details. Message: %s", string(message))
		}
	} else {
		// Send ACK for non-error responses
		go s.sendAcknowledgement(&response)
	}

	// Deliver response to waiting channel
	s.responseMu.RLock()
	ch, ok := s.responses[response.CorrelationID]
	s.responseMu.RUnlock()

	if ok {
		select {
		case ch <- &response:
		case <-time.After(5 * time.Second):
			log.Printf("Timeout delivering response for correlation ID: %s", response.CorrelationID)
		}
	}
}

func (s *MTSService) sendAcknowledgement(response *models.TicketResponse) error {
	operatorID := s.cfg.OperatorID
	if operatorID == 0 {
		log.Println("Warning: OperatorID is not set in config. Using default 9985 for ACK.")
		operatorID = 9985
	}

	ack := models.TicketAck{
		OperatorID:    operatorID,
		CorrelationID: response.CorrelationID,
		TimestampUTC:  time.Now().UnixMilli(),
		Operation:     "ticket-placement-ack",
		Version:       "3.0",
		Content: models.TicketAckContent{
			Type:         "ticket-ack",
			TicketID:     response.Content.TicketID,
			Acknowledged: true,
		},
	}

	// Set correct signature based on operation
	switch ack.Operation {
	case "ticket-placement-ack":
		ack.Content.TicketSignature = response.Content.Signature
	case "ticket-cancel-ack":
		ack.Content.CancellationSignature = response.Content.Signature
	case "ticket-cashout-ack":
		ack.Content.CashoutSignature = response.Content.Signature
	case "ticket-ext-settlement-ack":
		ack.Content.SettlementSignature = response.Content.Signature
	}

	return s.sendMessage(&ack)
}

func (s *MTSService) SendTicket(ticket *models.TicketRequest) (*models.TicketResponse, error) {
	s.connMu.RLock()
	activeConn := s.activeConn
	s.connMu.RUnlock()

	if activeConn == nil || atomic.LoadInt32(&s.connected) != 1 {
		return nil, fmt.Errorf("not connected to MTS")
	}

	// Increment pending responses counter
	atomic.AddInt32(&activeConn.pendingResponses, 1)

	responseCh := make(chan *models.TicketResponse, 1)
	s.responseMu.Lock()
	s.responses[ticket.CorrelationID] = responseCh
	s.responseMu.Unlock()

	defer func() {
		s.responseMu.Lock()
		delete(s.responses, ticket.CorrelationID)
		s.responseMu.Unlock()
		close(responseCh)
	}()

	if err := s.sendMessage(ticket); err != nil {
		atomic.AddInt32(&activeConn.pendingResponses, -1)
		return nil, fmt.Errorf("failed to send ticket: %w", err)
	}

	select {
	case response := <-responseCh:
		if response.Content.Type == "error-reply" {
			return nil, fmt.Errorf("MTS returned an error reply (version %s). Check service logs for details. CorrelationID: %s", response.Version, response.CorrelationID)
		}
		return response, nil
	case <-time.After(10 * time.Second):
		atomic.AddInt32(&activeConn.pendingResponses, -1)
		return nil, fmt.Errorf("timeout waiting for ticket response")
	case <-s.ctx.Done():
		atomic.AddInt32(&activeConn.pendingResponses, -1)
		return nil, fmt.Errorf("service closed")
	}
}

func (s *MTSService) sendMessage(msg interface{}) error {
	s.connMu.RLock()
	activeConn := s.activeConn
	s.connMu.RUnlock()

	if activeConn == nil {
		return fmt.Errorf("connection is nil")
	}

	activeConn.mu.RLock()
	conn := activeConn.conn
	activeConn.mu.RUnlock()

	if conn == nil {
		return fmt.Errorf("connection is nil")
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	// Log the message content
	logMessage := string(data)
	logMessage = strings.ReplaceAll(logMessage, "\n", "\t")
	log.Printf("Sending MTS message: %s", logMessage)

	conn.SetWriteDeadline(time.Now().Add(WriteWait))
	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	return nil
}

func (s *MTSService) reconnect() {
	if !atomic.CompareAndSwapInt32(&s.reconnecting, 0, 1) {
		return
	}

	defer atomic.StoreInt32(&s.reconnecting, 0)

	backoff := time.Second
	maxBackoff := 60 * time.Second

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-time.After(backoff):
			log.Println("Attempting to reconnect to MTS...")
			if err := s.connect(); err != nil {
				log.Printf("Reconnection failed: %v", err)
				backoff *= 2
				if backoff > maxBackoff {
					backoff = maxBackoff
				}
			} else {
				log.Println("Reconnected successfully")
				return
			}
		}
	}
}

func (s *MTSService) IsConnected() bool {
	return atomic.LoadInt32(&s.connected) == 1
}
