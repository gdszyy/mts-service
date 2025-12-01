package websocket

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		// Allow all origins for now (should be restricted in production)
		return true
	},
}

// Handler handles WebSocket HTTP requests
type Handler struct {
	hub *Hub
}

// NewHandler creates a new WebSocket handler
func NewHandler(hub *Hub) *Handler {
	return &Handler{
		hub: hub,
	}
}

// ServeWS handles WebSocket upgrade requests
func (h *Handler) ServeWS(w http.ResponseWriter, r *http.Request) {
	// Get userID from query parameters (in production, use proper authentication)
	userID := r.URL.Query().Get("userId")
	if userID == "" {
		log.Println("WebSocket connection rejected: missing userId")
		http.Error(w, "Missing userId parameter", http.StatusBadRequest)
		return
	}

	// Optional: Validate token
	token := r.URL.Query().Get("token")
	if token == "" {
		log.Printf("WebSocket connection rejected for user %s: missing token", userID)
		http.Error(w, "Missing token parameter", http.StatusUnauthorized)
		return
	}

	// TODO: Validate token against your authentication system
	// For now, we'll accept any non-empty token

	// Upgrade HTTP connection to WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade WebSocket connection for user %s: %v", userID, err)
		return
	}

	// Create new client
	client := NewClient(h.hub, conn, userID)

	// Register client with hub
	h.hub.register <- client

	// Send connection established message
	client.SendMessage(&ConnectionEstablishedResponse{
		BaseMessage: BaseMessage{
			Type:      MessageTypeConnectionEstablished,
			Timestamp: time.Now(),
		},
		UserID:  userID,
		Message: "WebSocket connection established successfully",
	})

	// Start client goroutines
	go client.writePump()
	go client.readPump()

	log.Printf("WebSocket connection established for user %s", userID)
}
