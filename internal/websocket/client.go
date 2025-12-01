package websocket

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period (must be less than pongWait)
	pingPeriod = 30 * time.Second

	// Maximum message size allowed from peer
	maxMessageSize = 512 * 1024 // 512 KB
)

// Client represents a WebSocket client connection
type Client struct {
	hub    *Hub
	conn   *websocket.Conn
	send   chan []byte
	userID string
	mu     sync.Mutex
}

// NewClient creates a new WebSocket client
func NewClient(hub *Hub, conn *websocket.Conn, userID string) *Client {
	return &Client{
		hub:    hub,
		conn:   conn,
		send:   make(chan []byte, 256),
		userID: userID,
	}
}

// readPump pumps messages from the WebSocket connection to the hub
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error for user %s: %v", c.userID, err)
			}
			break
		}

		// Parse base message to determine type
		var baseMsg BaseMessage
		if err := json.Unmarshal(message, &baseMsg); err != nil {
			log.Printf("Failed to parse message from user %s: %v", c.userID, err)
			c.SendError("", "Invalid message format", nil)
			continue
		}

		// Handle different message types
		switch baseMsg.Type {
		case MessageTypePlaceBet:
			var req PlaceBetRequest
			if err := json.Unmarshal(message, &req); err != nil {
				log.Printf("Failed to parse place_bet request: %v", err)
				c.SendError("", "Invalid place_bet request", nil)
				continue
			}
			c.hub.betRequests <- &BetRequest{
				Client:  c,
				Request: &req,
			}

		case MessageTypeQueryBetStatus:
			var req QueryBetStatusRequest
			if err := json.Unmarshal(message, &req); err != nil {
				log.Printf("Failed to parse query_bet_status request: %v", err)
				c.SendError("", "Invalid query_bet_status request", nil)
				continue
			}
			c.hub.statusQueries <- &StatusQuery{
				Client:  c,
				Request: &req,
			}

		case MessageTypePing:
			// Respond with pong
			c.SendPong()

		default:
			log.Printf("Unknown message type from user %s: %s", c.userID, baseMsg.Type)
			c.SendError("", "Unknown message type", nil)
		}
	}
}

// writePump pumps messages from the hub to the WebSocket connection
func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// Hub closed the channel
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// Add queued messages to the current WebSocket message
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// SendMessage sends a message to the client
func (c *Client) SendMessage(msg interface{}) error {
	c.mu.Lock()
	defer c.mu.Unlock()

	data, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	select {
	case c.send <- data:
	default:
		log.Printf("Client %s send buffer full, closing connection", c.userID)
		close(c.send)
	}

	return nil
}

// SendError sends an error message to the client
func (c *Client) SendError(requestID, errorMsg string, details map[string]interface{}) {
	c.SendMessage(&BetErrorResponse{
		BaseMessage: BaseMessage{
			Type:      MessageTypeBetError,
			Timestamp: time.Now(),
		},
		RequestID: requestID,
		Error:     errorMsg,
		Details:   details,
	})
}

// SendPong sends a pong message to the client
func (c *Client) SendPong() {
	c.SendMessage(&PongMessage{
		BaseMessage: BaseMessage{
			Type:      MessageTypePong,
			Timestamp: time.Now(),
		},
	})
}
