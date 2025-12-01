package websocket

import (
	"log"
	"sync"
)

// BetRequest represents a bet request from a client
type BetRequest struct {
	Client  *Client
	Request *PlaceBetRequest
}

// StatusQuery represents a status query from a client
type StatusQuery struct {
	Client  *Client
	Request *QueryBetStatusRequest
}

// Hub maintains the set of active clients and broadcasts messages to clients
type Hub struct {
	// Registered clients
	clients map[*Client]bool

	// Client lookup by userID
	clientsByUser map[string]*Client

	// Register requests from clients
	register chan *Client

	// Unregister requests from clients
	unregister chan *Client

	// Bet requests from clients
	betRequests chan *BetRequest

	// Status queries from clients
	statusQueries chan *StatusQuery

	// Mutex for thread-safe operations
	mu sync.RWMutex
}

// NewHub creates a new Hub
func NewHub() *Hub {
	return &Hub{
		clients:       make(map[*Client]bool),
		clientsByUser: make(map[string]*Client),
		register:      make(chan *Client),
		unregister:    make(chan *Client),
		betRequests:   make(chan *BetRequest, 256),
		statusQueries: make(chan *StatusQuery, 256),
	}
}

// Run starts the hub's main loop
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			// If user already has a connection, close the old one
			if oldClient, exists := h.clientsByUser[client.userID]; exists {
				log.Printf("User %s already connected, closing old connection", client.userID)
				close(oldClient.send)
				delete(h.clients, oldClient)
			}
			h.clients[client] = true
			h.clientsByUser[client.userID] = client
			h.mu.Unlock()
			log.Printf("Client registered: user=%s, total clients=%d", client.userID, len(h.clients))

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				delete(h.clientsByUser, client.userID)
				close(client.send)
				log.Printf("Client unregistered: user=%s, total clients=%d", client.userID, len(h.clients))
			}
			h.mu.Unlock()
		}
	}
}

// GetClient returns the client for a given userID
func (h *Hub) GetClient(userID string) (*Client, bool) {
	h.mu.RLock()
	defer h.mu.RUnlock()
	client, ok := h.clientsByUser[userID]
	return client, ok
}

// GetClientByTicketID returns the client associated with a ticketID
// Note: This requires maintaining a ticketID -> userID mapping elsewhere
func (h *Hub) SendToUser(userID string, message interface{}) bool {
	client, ok := h.GetClient(userID)
	if !ok {
		return false
	}
	return client.SendMessage(message) == nil
}

// ClientCount returns the number of connected clients
func (h *Hub) ClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}
