package service

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"

	"github.com/gdsZyy/mts-service/internal/models"
)

// computeMessageHash computes a SHA256 hash of the message JSON for idempotency
func (s *MTSService) computeMessageHash(ticket *models.TicketRequest) string {
	data, err := json.Marshal(ticket)
	if err != nil {
		log.Printf("Warning: Failed to marshal ticket for hashing: %v", err)
		return ""
	}
	hash := sha256.Sum256(data)
	return hex.EncodeToString(hash[:])
}

// SendTicketWithIdempotency sends a ticket with idempotency support
// If the same ticket (same JSON) is sent multiple times, the original response is returned
func (s *MTSService) SendTicketWithIdempotency(ticket *models.TicketRequest) (*models.TicketResponse, error) {
	// Compute message hash for idempotency
	msgHash := s.computeMessageHash(ticket)

	// Check if we have already sent this exact message and have a cached response
	s.sentMsgMu.RLock()
	if cachedResponse, exists := s.sentMessages[msgHash]; exists {
		s.sentMsgMu.RUnlock()
		log.Printf("Returning cached response for duplicate message (hash: %s)", msgHash)
		return cachedResponse, nil
	}
	s.sentMsgMu.RUnlock()

	// Send the ticket normally
	response, err := s.SendTicket(ticket)
	if err != nil {
		return nil, err
	}

	// Cache the response for future identical requests
	s.sentMsgMu.Lock()
	s.sentMessages[msgHash] = response
	s.sentMsgMu.Unlock()

	return response, nil
}
