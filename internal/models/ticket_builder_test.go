package models

import (
	"encoding/json"
	"testing"
)

func TestSingleBet(t *testing.T) {
	builder := NewTicketBuilder(45426, "test-single-001")
	
	selection := NewSelection("3", "sr:match:12345", "1", "1", 2.50)
	stake := NewStake("cash", "EUR", 10.00, "total")
	
	builder.AddSingleBet(selection, stake)
	builder.SetContext(&Context{
		Channel: &Channel{Type: "internet", Lang: "EN"},
		LimitID: 4268,
	})
	
	ticket := builder.Build("corr-single-001")
	
	// Validate structure
	if ticket.Content.Type != "ticket" {
		t.Errorf("Expected type 'ticket', got '%s'", ticket.Content.Type)
	}
	
	if len(ticket.Content.Bets) != 1 {
		t.Fatalf("Expected 1 bet, got %d", len(ticket.Content.Bets))
	}
	
	bet := ticket.Content.Bets[0]
	if len(bet.Selections) != 1 {
		t.Errorf("Expected 1 selection, got %d", len(bet.Selections))
	}
	
	if bet.Selections[0].Type != "uf" {
		t.Errorf("Expected selection type 'uf', got '%s'", bet.Selections[0].Type)
	}
	
	// Test JSON serialization
	jsonData, err := json.MarshalIndent(ticket, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal ticket: %v", err)
	}
	
	t.Logf("Single Bet JSON:\n%s", string(jsonData))
}

func TestAccumulatorBet(t *testing.T) {
	builder := NewTicketBuilder(45426, "test-accumulator-001")
	
	selections := []Selection{
		NewSelection("3", "sr:match:12345", "1", "1", 2.50),
		NewSelection("3", "sr:match:12346", "1", "2", 1.80),
		NewSelection("3", "sr:match:12347", "1", "1", 3.00),
	}
	stake := NewStake("cash", "EUR", 10.00, "total")
	
	builder.AddAccumulatorBet(selections, stake)
	builder.SetContext(&Context{
		Channel: &Channel{Type: "internet", Lang: "EN"},
		LimitID: 4268,
	})
	
	ticket := builder.Build("corr-accumulator-001")
	
	// Validate structure
	if len(ticket.Content.Bets) != 1 {
		t.Fatalf("Expected 1 bet, got %d", len(ticket.Content.Bets))
	}
	
	bet := ticket.Content.Bets[0]
	if len(bet.Selections) != 3 {
		t.Errorf("Expected 3 selections, got %d", len(bet.Selections))
	}
	
	for i, sel := range bet.Selections {
		if sel.Type != "uf" {
			t.Errorf("Selection %d: expected type 'uf', got '%s'", i, sel.Type)
		}
	}
	
	// Test JSON serialization
	jsonData, err := json.MarshalIndent(ticket, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal ticket: %v", err)
	}
	
	t.Logf("Accumulator Bet JSON:\n%s", string(jsonData))
}

func TestSystemBet(t *testing.T) {
	builder := NewTicketBuilder(45426, "test-system-001")
	
	selections := []Selection{
		NewSelection("3", "sr:match:12345", "1", "1", 2.50),
		NewSelection("3", "sr:match:12346", "1", "2", 1.80),
		NewSelection("3", "sr:match:12347", "1", "1", 3.00),
		NewSelection("3", "sr:match:12348", "1", "2", 2.20),
	}
	stake := NewStake("cash", "EUR", 1.00, "unit") // Unit stake for system bets
	
	// 2/4 system (all doubles)
	builder.AddSystemBet([]int{2}, selections, stake)
	builder.SetContext(&Context{
		Channel: &Channel{Type: "internet", Lang: "EN"},
		LimitID: 4268,
	})
	
	ticket := builder.Build("corr-system-001")
	
	// Validate structure
	if len(ticket.Content.Bets) != 1 {
		t.Fatalf("Expected 1 bet, got %d", len(ticket.Content.Bets))
	}
	
	bet := ticket.Content.Bets[0]
	if len(bet.Selections) != 1 {
		t.Fatalf("Expected 1 selection (system), got %d", len(bet.Selections))
	}
	
	systemSel := bet.Selections[0]
	if systemSel.Type != "system" {
		t.Errorf("Expected selection type 'system', got '%s'", systemSel.Type)
	}
	
	if len(systemSel.Size) != 1 || systemSel.Size[0] != 2 {
		t.Errorf("Expected size [2], got %v", systemSel.Size)
	}
	
	if len(systemSel.Selections) != 4 {
		t.Errorf("Expected 4 nested selections, got %d", len(systemSel.Selections))
	}
	
	// Test JSON serialization
	jsonData, err := json.MarshalIndent(ticket, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal ticket: %v", err)
	}
	
	t.Logf("System Bet (2/4) JSON:\n%s", string(jsonData))
}

func TestTrixieBet(t *testing.T) {
	builder := NewTicketBuilder(45426, "test-trixie-001")
	
	selections := []Selection{
		NewSelection("3", "sr:match:12345", "1", "1", 2.50),
		NewSelection("3", "sr:match:12346", "1", "2", 1.80),
		NewSelection("3", "sr:match:12347", "1", "1", 3.00),
	}
	stake := NewStake("cash", "EUR", 1.00, "unit")
	
	builder.AddTrixieBet(selections, stake)
	builder.SetContext(&Context{
		Channel: &Channel{Type: "internet", Lang: "EN"},
		LimitID: 4268,
	})
	
	ticket := builder.Build("corr-trixie-001")
	
	// Validate structure
	bet := ticket.Content.Bets[0]
	systemSel := bet.Selections[0]
	
	if systemSel.Type != "system" {
		t.Errorf("Expected selection type 'system', got '%s'", systemSel.Type)
	}
	
	// Trixie: size [2,3] (3 doubles + 1 treble)
	if len(systemSel.Size) != 2 || systemSel.Size[0] != 2 || systemSel.Size[1] != 3 {
		t.Errorf("Expected size [2,3], got %v", systemSel.Size)
	}
	
	// Test JSON serialization
	jsonData, err := json.MarshalIndent(ticket, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal ticket: %v", err)
	}
	
	t.Logf("Trixie Bet JSON:\n%s", string(jsonData))
}

func TestYankeeBet(t *testing.T) {
	builder := NewTicketBuilder(45426, "test-yankee-001")
	
	selections := []Selection{
		NewSelection("3", "sr:match:12345", "1", "1", 2.50),
		NewSelection("3", "sr:match:12346", "1", "2", 1.80),
		NewSelection("3", "sr:match:12347", "1", "1", 3.00),
		NewSelection("3", "sr:match:12348", "1", "2", 2.20),
	}
	stake := NewStake("cash", "EUR", 1.00, "unit")
	
	builder.AddYankeeBet(selections, stake)
	builder.SetContext(&Context{
		Channel: &Channel{Type: "internet", Lang: "EN"},
		LimitID: 4268,
	})
	
	ticket := builder.Build("corr-yankee-001")
	
	// Validate structure
	bet := ticket.Content.Bets[0]
	systemSel := bet.Selections[0]
	
	// Yankee: size [2,3,4] (6 doubles + 4 trebles + 1 four-fold)
	if len(systemSel.Size) != 3 || systemSel.Size[0] != 2 || systemSel.Size[1] != 3 || systemSel.Size[2] != 4 {
		t.Errorf("Expected size [2,3,4], got %v", systemSel.Size)
	}
	
	// Test JSON serialization
	jsonData, err := json.MarshalIndent(ticket, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal ticket: %v", err)
	}
	
	t.Logf("Yankee Bet JSON:\n%s", string(jsonData))
}

func TestBankerSystemBet(t *testing.T) {
	builder := NewTicketBuilder(45426, "test-banker-001")
	
	bankers := []Selection{
		NewSelection("3", "sr:match:12345", "1", "1", 1.50),
	}
	
	selections := []Selection{
		NewSelection("3", "sr:match:12346", "1", "2", 2.50),
		NewSelection("3", "sr:match:12347", "1", "1", 3.00),
		NewSelection("3", "sr:match:12348", "1", "2", 2.20),
	}
	stake := NewStake("cash", "EUR", 1.00, "unit")
	
	// 2/3 system with 1 banker
	builder.AddBankerSystemBet(bankers, []int{2}, selections, stake)
	builder.SetContext(&Context{
		Channel: &Channel{Type: "internet", Lang: "EN"},
		LimitID: 4268,
	})
	
	ticket := builder.Build("corr-banker-001")
	
	// Validate structure
	bet := ticket.Content.Bets[0]
	
	// Should have 2 selections: 1 banker + 1 system
	if len(bet.Selections) != 2 {
		t.Fatalf("Expected 2 selections (1 banker + 1 system), got %d", len(bet.Selections))
	}
	
	// First selection should be banker (standard "uf" type)
	if bet.Selections[0].Type != "uf" {
		t.Errorf("Expected first selection type 'uf' (banker), got '%s'", bet.Selections[0].Type)
	}
	
	// Second selection should be system
	if bet.Selections[1].Type != "system" {
		t.Errorf("Expected second selection type 'system', got '%s'", bet.Selections[1].Type)
	}
	
	systemSel := bet.Selections[1]
	if len(systemSel.Size) != 1 || systemSel.Size[0] != 2 {
		t.Errorf("Expected size [2], got %v", systemSel.Size)
	}
	
	if len(systemSel.Selections) != 3 {
		t.Errorf("Expected 3 nested selections, got %d", len(systemSel.Selections))
	}
	
	// Test JSON serialization
	jsonData, err := json.MarshalIndent(ticket, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal ticket: %v", err)
	}
	
	t.Logf("Banker System Bet (1 banker + 2/3 system) JSON:\n%s", string(jsonData))
}

func TestLucky15Bet(t *testing.T) {
	builder := NewTicketBuilder(45426, "test-lucky15-001")
	
	selections := []Selection{
		NewSelection("3", "sr:match:12345", "1", "1", 2.50),
		NewSelection("3", "sr:match:12346", "1", "2", 1.80),
		NewSelection("3", "sr:match:12347", "1", "1", 3.00),
		NewSelection("3", "sr:match:12348", "1", "2", 2.20),
	}
	stake := NewStake("cash", "EUR", 1.00, "unit")
	
	builder.AddLucky15Bet(selections, stake)
	builder.SetContext(&Context{
		Channel: &Channel{Type: "internet", Lang: "EN"},
		LimitID: 4268,
	})
	
	ticket := builder.Build("corr-lucky15-001")
	
	// Validate structure
	bet := ticket.Content.Bets[0]
	systemSel := bet.Selections[0]
	
	// Lucky 15: size [1,2,3,4] (4 singles + 6 doubles + 4 trebles + 1 four-fold)
	if len(systemSel.Size) != 4 {
		t.Errorf("Expected size array length 4, got %d", len(systemSel.Size))
	}
	
	// Test JSON serialization
	jsonData, err := json.MarshalIndent(ticket, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal ticket: %v", err)
	}
	
	t.Logf("Lucky 15 Bet JSON:\n%s", string(jsonData))
}

func TestMultipleBetsInOneTicket(t *testing.T) {
	builder := NewTicketBuilder(45426, "test-multi-001")
	
	// Add a single bet
	single := NewSelection("3", "sr:match:12345", "1", "1", 2.50)
	builder.AddSingleBet(single, NewStake("cash", "EUR", 5.00, "total"))
	
	// Add an accumulator
	accSelections := []Selection{
		NewSelection("3", "sr:match:12346", "1", "2", 1.80),
		NewSelection("3", "sr:match:12347", "1", "1", 3.00),
	}
	builder.AddAccumulatorBet(accSelections, NewStake("cash", "EUR", 10.00, "total"))
	
	// Add a system bet
	sysSelections := []Selection{
		NewSelection("3", "sr:match:12348", "1", "2", 2.20),
		NewSelection("3", "sr:match:12349", "1", "1", 1.95),
		NewSelection("3", "sr:match:12350", "1", "2", 2.75),
	}
	builder.AddTrixieBet(sysSelections, NewStake("cash", "EUR", 1.00, "unit"))
	
	builder.SetContext(&Context{
		Channel: &Channel{Type: "internet", Lang: "EN"},
		LimitID: 4268,
	})
	
	ticket := builder.Build("corr-multi-001")
	
	// Validate structure
	if len(ticket.Content.Bets) != 3 {
		t.Fatalf("Expected 3 bets, got %d", len(ticket.Content.Bets))
	}
	
	// Test JSON serialization
	jsonData, err := json.MarshalIndent(ticket, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal ticket: %v", err)
	}
	
	t.Logf("Multiple Bets in One Ticket JSON:\n%s", string(jsonData))
}
