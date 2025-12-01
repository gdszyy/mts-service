package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/gdsZyy/mts-service/internal/config"
	"github.com/gdsZyy/mts-service/internal/models"
	"github.com/gdsZyy/mts-service/internal/service"
	"github.com/google/uuid"
)

// BetProcessor handles bet requests from WebSocket clients
type BetProcessor struct {
	hub        *Hub
	mtsService *service.MTSService
	cfg        *config.Config
	
	// Track pending tickets for status queries
	pendingTickets map[string]string // ticketID -> userID
}

// NewBetProcessor creates a new BetProcessor
func NewBetProcessor(hub *Hub, mtsService *service.MTSService, cfg *config.Config) *BetProcessor {
	return &BetProcessor{
		hub:            hub,
		mtsService:     mtsService,
		cfg:            cfg,
		pendingTickets: make(map[string]string),
	}
}

// Start begins processing bet requests
func (bp *BetProcessor) Start() {
	go bp.processBetRequests()
	go bp.processStatusQueries()
}

// processBetRequests handles incoming bet requests
func (bp *BetProcessor) processBetRequests() {
	for betReq := range bp.hub.betRequests {
		go bp.handleBetRequest(betReq)
	}
}

// processStatusQueries handles status query requests
func (bp *BetProcessor) processStatusQueries() {
	for query := range bp.hub.statusQueries {
		go bp.handleStatusQuery(query)
	}
}

// handleBetRequest processes a single bet request
func (bp *BetProcessor) handleBetRequest(betReq *BetRequest) {
	client := betReq.Client
	req := betReq.Request

	log.Printf("Processing bet request: requestID=%s, betType=%s, userID=%s", 
		req.RequestID, req.BetType, client.userID)

	// Generate ticket ID(s)
	var ticketIDs []string
	var tickets []*models.TicketRequest

	switch req.BetType {
	case "single":
		ticket, err := bp.buildSingleBet(req)
		if err != nil {
			client.SendError(req.RequestID, fmt.Sprintf("Failed to build ticket: %v", err), nil)
			return
		}
		tickets = append(tickets, ticket)
		ticketIDs = append(ticketIDs, ticket.Content.TicketID)

	case "multi":
		// Handle multiple single bets
		bets, ok := req.Payload["bets"].([]interface{})
		if !ok {
			client.SendError(req.RequestID, "Invalid multi bet payload", nil)
			return
		}
		for _, betData := range bets {
			ticket, err := bp.buildSingleBetFromPayload(betData)
			if err != nil {
				client.SendError(req.RequestID, fmt.Sprintf("Failed to build ticket: %v", err), nil)
				return
			}
			tickets = append(tickets, ticket)
			ticketIDs = append(ticketIDs, ticket.Content.TicketID)
		}

	case "accumulator":
		ticket, err := bp.buildAccumulatorBet(req)
		if err != nil {
			client.SendError(req.RequestID, fmt.Sprintf("Failed to build ticket: %v", err), nil)
			return
		}
		tickets = append(tickets, ticket)
		ticketIDs = append(ticketIDs, ticket.Content.TicketID)

	case "system":
		ticket, err := bp.buildSystemBet(req)
		if err != nil {
			client.SendError(req.RequestID, fmt.Sprintf("Failed to build ticket: %v", err), nil)
			return
		}
		tickets = append(tickets, ticket)
		ticketIDs = append(ticketIDs, ticket.Content.TicketID)

	case "banker":
		ticket, err := bp.buildBankerBet(req)
		if err != nil {
			client.SendError(req.RequestID, fmt.Sprintf("Failed to build ticket: %v", err), nil)
			return
		}
		tickets = append(tickets, ticket)
		ticketIDs = append(ticketIDs, ticket.Content.TicketID)

	default:
		client.SendError(req.RequestID, fmt.Sprintf("Unknown bet type: %s", req.BetType), nil)
		return
	}

	// Send bet received confirmation
	if len(ticketIDs) == 1 {
		client.SendMessage(&BetReceivedResponse{
			BaseMessage: BaseMessage{
				Type:      MessageTypeBetReceived,
				Timestamp: time.Now(),
			},
			RequestID: req.RequestID,
			TicketID:  ticketIDs[0],
		})
	} else {
		client.SendMessage(&BetReceivedResponse{
			BaseMessage: BaseMessage{
				Type:      MessageTypeBetReceived,
				Timestamp: time.Now(),
			},
			RequestID: req.RequestID,
			TicketIDs: ticketIDs,
		})
	}

	// Track pending tickets
	for _, ticketID := range ticketIDs {
		bp.pendingTickets[ticketID] = client.userID
	}

	// Process tickets
	if len(tickets) == 1 {
		bp.processSingleTicket(client, req.RequestID, tickets[0])
	} else {
		bp.processMultipleTickets(client, req.RequestID, tickets)
	}
}

// processSingleTicket sends a single ticket to MTS and pushes result
func (bp *BetProcessor) processSingleTicket(client *Client, requestID string, ticket *models.TicketRequest) {
	// Send to MTS
	response, err := bp.mtsService.SendTicket(ticket)
	if err != nil {
		client.SendError(requestID, fmt.Sprintf("Failed to send ticket: %v", err), nil)
		delete(bp.pendingTickets, ticket.Content.TicketID)
		return
	}

	// Convert response to map for details
	details := make(map[string]interface{})
	responseBytes, _ := json.Marshal(response)
	json.Unmarshal(responseBytes, &details)

	// Determine status
	status := "rejected"
	if response.Content.Status == "accepted" {
		status = "accepted"
	}

	// Send result
	client.SendMessage(&BetResultResponse{
		BaseMessage: BaseMessage{
			Type:      MessageTypeBetResult,
			Timestamp: time.Now(),
		},
		RequestID: requestID,
		TicketID:  ticket.Content.TicketID,
		Status:    status,
		Details:   details,
	})

	delete(bp.pendingTickets, ticket.Content.TicketID)
}

// processMultipleTickets sends multiple tickets to MTS and pushes partial/final results
func (bp *BetProcessor) processMultipleTickets(client *Client, requestID string, tickets []*models.TicketRequest) {
	total := len(tickets)
	completed := 0
	accepted := 0
	rejected := 0
	var allDetails []map[string]interface{}

	for _, ticket := range tickets {
		// Send to MTS
		response, err := bp.mtsService.SendTicket(ticket)
		
		var details map[string]interface{}
		var status string
		
		if err != nil {
			status = "rejected"
			details = map[string]interface{}{
				"error": err.Error(),
			}
			rejected++
		} else {
			responseBytes, _ := json.Marshal(response)
			json.Unmarshal(responseBytes, &details)
			
			if response.Content.Status == "accepted" {
				status = "accepted"
				accepted++
			} else {
				status = "rejected"
				rejected++
			}
		}

		completed++
		allDetails = append(allDetails, details)

		// Send partial result
		client.SendMessage(&BetPartialResultResponse{
			BaseMessage: BaseMessage{
				Type:      MessageTypeBetPartialResult,
				Timestamp: time.Now(),
			},
			RequestID: requestID,
			Completed: fmt.Sprintf("%d/%d", completed, total),
			TicketID:  ticket.Content.TicketID,
			Status:    status,
			Details:   details,
		})

		delete(bp.pendingTickets, ticket.Content.TicketID)
	}

	// Send final result with summary
	client.SendMessage(&BetResultResponse{
		BaseMessage: BaseMessage{
			Type:      MessageTypeBetResult,
			Timestamp: time.Now(),
		},
		RequestID: requestID,
		Status:    "completed",
		Details:   map[string]interface{}{},
		Summary: &BetSummary{
			Total:    total,
			Accepted: accepted,
			Rejected: rejected,
			Details:  allDetails,
		},
	})
}

// handleStatusQuery processes a status query request
func (bp *BetProcessor) handleStatusQuery(query *StatusQuery) {
	client := query.Client
	req := query.Request

	// In a real implementation, you would query the ticket status from a database or cache
	// For now, we'll just check if it's in pending tickets
	_, isPending := bp.pendingTickets[req.TicketID]
	
	status := "not_found"
	if isPending {
		status = "pending"
	}

	client.SendMessage(&BetStatusResponse{
		BaseMessage: BaseMessage{
			Type:      MessageTypeBetStatus,
			Timestamp: time.Now(),
		},
		TicketID: req.TicketID,
		Status:   status,
	})
}

// Helper functions to build tickets from WebSocket requests

func (bp *BetProcessor) buildSingleBet(req *PlaceBetRequest) (*models.TicketRequest, error) {
	ticketID := uuid.New().String()
	builder := models.NewTicketBuilder(bp.cfg.OperatorID, ticketID)

	// Extract selection and stake from payload
	selectionData, ok := req.Payload["selection"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid selection data")
	}
	
	stakeData, ok := req.Payload["stake"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid stake data")
	}

	selection := convertSelection(selectionData)
	stake := convertStake(stakeData)

	builder.AddSingleBet(selection, stake)
	builder.SetContext(getDefaultContext(bp.cfg))

	return builder.Build(uuid.New().String()), nil
}

func (bp *BetProcessor) buildSingleBetFromPayload(betData interface{}) (*models.TicketRequest, error) {
	betMap, ok := betData.(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid bet data")
	}

	ticketID := uuid.New().String()
	builder := models.NewTicketBuilder(bp.cfg.OperatorID, ticketID)

	selectionData, ok := betMap["selection"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid selection data")
	}
	
	stakeData, ok := betMap["stake"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid stake data")
	}

	selection := convertSelection(selectionData)
	stake := convertStake(stakeData)

	builder.AddSingleBet(selection, stake)
	builder.SetContext(getDefaultContext(bp.cfg))

	return builder.Build(uuid.New().String()), nil
}

func (bp *BetProcessor) buildAccumulatorBet(req *PlaceBetRequest) (*models.TicketRequest, error) {
	ticketID := uuid.New().String()
	builder := models.NewTicketBuilder(bp.cfg.OperatorID, ticketID)

	selectionsData, ok := req.Payload["selections"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid selections data")
	}
	
	stakeData, ok := req.Payload["stake"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid stake data")
	}

	var selections []models.Selection
	for _, selData := range selectionsData {
		selMap, ok := selData.(map[string]interface{})
		if !ok {
			continue
		}
		selections = append(selections, convertSelection(selMap))
	}

	stake := convertStake(stakeData)

	builder.AddAccumulatorBet(selections, stake)
	builder.SetContext(getDefaultContext(bp.cfg))

	return builder.Build(uuid.New().String()), nil
}

func (bp *BetProcessor) buildSystemBet(req *PlaceBetRequest) (*models.TicketRequest, error) {
	ticketID := uuid.New().String()
	builder := models.NewTicketBuilder(bp.cfg.OperatorID, ticketID)

	selectionsData, ok := req.Payload["selections"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid selections data")
	}
	
	stakeData, ok := req.Payload["stake"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid stake data")
	}

	systemSize, ok := req.Payload["systemSize"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid systemSize")
	}

	var selections []models.Selection
	for _, selData := range selectionsData {
		selMap, ok := selData.(map[string]interface{})
		if !ok {
			continue
		}
		selections = append(selections, convertSelection(selMap))
	}

	stake := convertStake(stakeData)

	builder.AddSystemBet([]int{int(systemSize)}, selections, stake)
	builder.SetContext(getDefaultContext(bp.cfg))

	return builder.Build(uuid.New().String()), nil
}

func (bp *BetProcessor) buildBankerBet(req *PlaceBetRequest) (*models.TicketRequest, error) {
	ticketID := uuid.New().String()
	builder := models.NewTicketBuilder(bp.cfg.OperatorID, ticketID)

	selectionsData, ok := req.Payload["selections"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid selections data")
	}
	
	bankerSelectionsData, ok := req.Payload["bankerSelections"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid bankerSelections data")
	}
	
	stakeData, ok := req.Payload["stake"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid stake data")
	}

	systemSize, ok := req.Payload["systemSize"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid systemSize")
	}

	var selections []models.Selection
	for _, selData := range selectionsData {
		selMap, ok := selData.(map[string]interface{})
		if !ok {
			continue
		}
		selections = append(selections, convertSelection(selMap))
	}

	var bankerSelections []models.Selection
	for _, selData := range bankerSelectionsData {
		selMap, ok := selData.(map[string]interface{})
		if !ok {
			continue
		}
		bankerSelections = append(bankerSelections, convertSelection(selMap))
	}

	stake := convertStake(stakeData)

	builder.AddBankerSystemBet(bankerSelections, []int{int(systemSize)}, selections, stake)
	builder.SetContext(getDefaultContext(bp.cfg))

	return builder.Build(uuid.New().String()), nil
}

// Helper functions

func convertSelection(data map[string]interface{}) models.Selection {
	return models.Selection{
		Type:      "uf",
		ProductID: "3",
		EventID:   getStringValue(data, "eventId"),
		MarketID:  getStringValue(data, "marketId"),
		OutcomeID: getStringValue(data, "outcomeId"),
		Odds: &models.Odds{
			Type:  "decimal",
			Value: getStringValue(data, "odds"),
		},
	}
}

func convertStake(data map[string]interface{}) models.Stake {
	return models.Stake{
		Type:     "cash",
		Amount:   getStringValue(data, "amount"),
		Currency: getStringValue(data, "currency"),
	}
}

func getStringValue(data map[string]interface{}, key string) string {
	if val, ok := data[key]; ok {
		if str, ok := val.(string); ok {
			return str
		}
	}
	return ""
}

func getDefaultContext(cfg *config.Config) *models.Context {
	return &models.Context{
		LimitID: 1,
	}
}
