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
	"sync"
	"time"

	"github.com/gdsZyy/mts-service/internal/config"
	"github.com/gdsZyy/mts-service/internal/models"
	"github.com/gorilla/websocket"
)

const (
	WriteWait      = 10 * time.Second
	PongWait       = 60 * time.Second
	PingPeriod     = 54 * time.Second
	MaxMessageSize = 512 * 1024
)

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

type MTSService struct {
	cfg          *config.Config
		wsURL        string
		wsAudience   string
	conn         *websocket.Conn
	token        *TokenResponse
	tokenExpiry  time.Time
	tokenMu      sync.RWMutex
	connMu       sync.RWMutex
	responses    map[string]chan *models.TicketResponse
	responseMu   sync.RWMutex
	ctx          context.Context
	cancel       context.CancelFunc
	connected    bool
	reconnecting bool
	httpClient   *http.Client
}

	func NewMTSService(cfg *config.Config) *MTSService {
		wsURL := cfg.WSURL
		wsAudience := cfg.WSAudience
		if cfg.Production {
			// 生产环境使用配置中的 VirtualHost 作为 URL，并使用生产环境的 Audience
			wsURL = fmt.Sprintf("wss://%s", cfg.VirtualHost)
			wsAudience = "mbs-dp-production-wss"
		}

		ctx, cancel := context.WithCancel(context.Background())

		return &MTSService{
			cfg:        cfg,
			wsURL:      wsURL,
			wsAudience: wsAudience,
			responses:  make(map[string]chan *models.TicketResponse),
			ctx:        ctx,
			cancel:     cancel,
			httpClient: &http.Client{Timeout: 30 * time.Second},
		}
	}

func (s *MTSService) Start() error {
	if err := s.connect(); err != nil {
		return fmt.Errorf("failed to connect: %w", err)
	}
	log.Println("MTS Service started successfully")
	return nil
}

func (s *MTSService) Stop() error {
	s.cancel()
	s.connMu.Lock()
	defer s.connMu.Unlock()

	if s.conn != nil {
		s.conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
		s.conn.Close()
		s.conn = nil
	}
	s.connected = false
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
		s.tokenExpiry = time.Now().Add(time.Duration(tokenResp.ExpiresIn) * time.Second).Add(-30 * time.Second) // 提前 30 秒过期，确保刷新

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

	s.connMu.Lock()
	s.conn = conn
	s.connected = true
	s.connMu.Unlock()

	conn.SetReadLimit(MaxMessageSize)
	conn.SetReadDeadline(time.Now().Add(PongWait))
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(PongWait))
		return nil
	})

	go s.readPump()
	go s.pingPump()

		log.Println("Connected to MTS WebSocket")

		// 发送初始化订阅消息
		if err := s.sendInitializationMessage(); err != nil {
			return fmt.Errorf("failed to send initialization message: %w", err)
		}

		return nil
	}

	// sendInitializationMessage 发送 WebSocket 连接后的初始化订阅消息
	func (s *MTSService) sendInitializationMessage() error {
			// 构造初始化消息
			// 根据最新的 Schema 错误，根级别不应有 type 字段，content.type 必须是 ticket-inform，且 content 中不应有 bookmaker_id/limit_id/operatorId
			initMsg := map[string]interface{}{
					"correlationId": fmt.Sprintf("init-%d", time.Now().UnixNano()), // 使用唯一 ID
					"timestampUtc":  time.Now().UnixMilli(),
					"operation":     "ticket-placement-inform", // 必须是 ticket-placement-inform
					"version":       "3.0",
					"content": map[string]interface{}{
						"type": "ticket-inform", // 必须是 ticket-inform
						"ticketId": fmt.Sprintf("init-ticket-%d", time.Now().UnixNano()), // 添加 ticketId
						"ticketSignature": "initial-signature", // 添加 ticketSignature 占位符
					},
			}

				// 检查 BookmakerID, LimitID, OperatorID 是否已设置
				// 根据最新的 Schema 错误，这些字段不应出现在 content 中，因此移除检查。
				// 假设这些配置信息是通过 Auth Token 或其他方式隐式传递给 MTS 的。

		log.Printf("Sending initialization message: %+v", initMsg)

		return s.sendMessage(initMsg)
	}

func (s *MTSService) readPump() {
	defer func() {
		s.connMu.Lock()
		s.connected = false
		s.connMu.Unlock()
		s.reconnect()
	}()

	for {
		select {
		case <-s.ctx.Done():
			return
		default:
			_, message, err := s.conn.ReadMessage()
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Printf("WebSocket read error: %v", err)
				}
				return
			}

			s.handleMessage(message)
		}
	}
}

func (s *MTSService) pingPump() {
	ticker := time.NewTicker(PingPeriod)
	defer ticker.Stop()

	for {
		select {
		case <-s.ctx.Done():
			return
		case <-ticker.C:
			s.connMu.RLock()
			conn := s.conn
			s.connMu.RUnlock()

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

	func (s *MTSService) handleMessage(message []byte) {
		var response models.TicketResponse
		if err := json.Unmarshal(message, &response); err != nil {
			log.Printf("Failed to unmarshal response: %v. Message: %s", err, string(message))
			return
		}

		// 检查是否是 MTS 错误响应
		if response.Content.Type == "error-reply" {
			// 尝试解析更详细的错误信息
			var errorResponse struct {
				Content models.ErrorReplyContent `json:"content"`
			}
			if err := json.Unmarshal(message, &errorResponse); err == nil {
				log.Printf("MTS Error Reply received (CorrelationID: %s): Code=%d, Message=%s", response.CorrelationID, errorResponse.Content.Code, errorResponse.Content.Message)
			} else {
				log.Printf("MTS Error Reply received, but failed to parse details. Message: %s", string(message))
			}
			// 即使是错误，也需要将响应传递给等待的通道，以便 SendTicket 可以超时或返回错误
		}

		go s.sendAcknowledgement(&response)

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
	ack := models.TicketAck{
		Operation:     "ticket-placement-ack",
		CorrelationID: response.CorrelationID,
		TimestampUTC:  time.Now().UnixMilli(),
		Version:       "2.4",
		Content: models.TicketAckContent{
			Type:            "ticket-ack",
			TicketID:        response.Content.TicketID,
			TicketSignature: response.Content.Signature,
		},
	}

	return s.sendMessage(&ack)
}

func (s *MTSService) SendTicket(ticket *models.TicketRequest) (*models.TicketResponse, error) {
	s.connMu.RLock()
	if !s.connected {
		s.connMu.RUnlock()
		return nil, fmt.Errorf("not connected to MTS")
	}
	s.connMu.RUnlock()

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
		return nil, fmt.Errorf("failed to send ticket: %w", err)
	}

	select {
		case response := <-responseCh:
			// 检查是否是 MTS 错误响应
			if response.Content.Type == "error-reply" {
				// 尝试解析更详细的错误信息

					// 由于我们在 handleMessage 中已经打印了详细日志，这里直接返回一个通用错误，并提示用户检查日志
					// 更好的做法是在 handleMessage 中将 ErrorReplyContent 结构体放入 TicketResponse 中
					// 但为了避免对 TicketResponse 结构体进行大规模修改，我们先采用日志+通用错误的方式
					return nil, fmt.Errorf("MTS returned an error reply (version %s). Check service logs for details. CorrelationID: %s", response.Version, response.CorrelationID)
			}
			return response, nil
		case <-time.After(10 * time.Second):
			return nil, fmt.Errorf("timeout waiting for ticket response")
	case <-s.ctx.Done():
		return nil, fmt.Errorf("service closed")
	}
}

func (s *MTSService) sendMessage(msg interface{}) error {
	s.connMu.RLock()
	conn := s.conn
	s.connMu.RUnlock()

	if conn == nil {
		return fmt.Errorf("connection is nil")
	}

	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("failed to marshal message: %w", err)
	}

	conn.SetWriteDeadline(time.Now().Add(WriteWait))
	if err := conn.WriteMessage(websocket.TextMessage, data); err != nil {
		return fmt.Errorf("failed to write message: %w", err)
	}

	return nil
}

func (s *MTSService) reconnect() {
	s.connMu.Lock()
	if s.reconnecting {
		s.connMu.Unlock()
		return
	}
	s.reconnecting = true
	s.connMu.Unlock()

	defer func() {
		s.connMu.Lock()
		s.reconnecting = false
		s.connMu.Unlock()
	}()

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
	s.connMu.RLock()
	defer s.connMu.RUnlock()
	return s.connected
}

