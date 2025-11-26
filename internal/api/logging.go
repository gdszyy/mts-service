package api

import (
	"log"
	"time"
)

// LogRequest logs incoming API requests
func LogRequest(endpoint, ticketID string, details ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	if len(details) > 0 {
		log.Printf("[%s] [%s] TicketID=%s, Details=%v", timestamp, endpoint, ticketID, details)
	} else {
		log.Printf("[%s] [%s] TicketID=%s", timestamp, endpoint, ticketID)
	}
}

// LogValidationError logs validation errors
func LogValidationError(endpoint, ticketID string, err error) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	log.Printf("[%s] [%s] Validation failed: TicketID=%s, Error=%v", timestamp, endpoint, ticketID, err)
}

// LogMTSRequest logs MTS request sending
func LogMTSRequest(endpoint, ticketID string, details ...interface{}) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	if len(details) > 0 {
		log.Printf("[%s] [%s] Sending to MTS: TicketID=%s, Details=%v", timestamp, endpoint, ticketID, details)
	} else {
		log.Printf("[%s] [%s] Sending to MTS: TicketID=%s", timestamp, endpoint, ticketID)
	}
}

// LogMTSResponse logs MTS response
func LogMTSResponse(endpoint, ticketID, status string, accepted bool) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	if accepted {
		log.Printf("[%s] [%s] ✓ Ticket ACCEPTED: TicketID=%s, Status=%s", timestamp, endpoint, ticketID, status)
	} else {
		log.Printf("[%s] [%s] ✗ Ticket REJECTED: TicketID=%s, Status=%s", timestamp, endpoint, ticketID, status)
	}
}

// LogMTSError logs MTS errors
func LogMTSError(endpoint, ticketID string, err error) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	log.Printf("[%s] [%s] ✗ MTS Error: TicketID=%s, Error=%v", timestamp, endpoint, ticketID, err)
}

// LogCashoutRequest logs cashout requests
func LogCashoutRequest(cashoutID, ticketID string, cashoutType string, amount float64) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	log.Printf("[%s] [Cashout] Request: CashoutID=%s, TicketID=%s, Type=%s, Amount=%.2f", 
		timestamp, cashoutID, ticketID, cashoutType, amount)
}

// LogCashoutResponse logs cashout responses
func LogCashoutResponse(cashoutID, status string, code int) {
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	if status == "accepted" {
		log.Printf("[%s] [Cashout] ✓ ACCEPTED: CashoutID=%s, Status=%s, Code=%d", 
			timestamp, cashoutID, status, code)
	} else {
		log.Printf("[%s] [Cashout] ✗ REJECTED: CashoutID=%s, Status=%s, Code=%d", 
			timestamp, cashoutID, status, code)
	}
}
