package models

import (
	"encoding/json"
	"testing"
)

// TestBankerSystemBet tests the corrected banker system bet implementation
func TestBankerSystemBet(t *testing.T) {
	builder := NewTicketBuilder(45426, "test-banker-001")
	
	// Create 2 banker selections
	banker1 := NewSelection("3", "sr:match:11111", "1", "1", 1.50)
	banker2 := NewSelection("3", "sr:match:22222", "1", "1", 1.80)
	bankers := []Selection{banker1, banker2}
	
	// Create 3 non-banker selections
	sel1 := NewSelection("3", "sr:match:33333", "1", "1", 2.00)
	sel2 := NewSelection("3", "sr:match:44444", "1", "1", 2.20)
	sel3 := NewSelection("3", "sr:match:55555", "1", "1", 2.50)
	selections := []Selection{sel1, sel2, sel3}
	
	// Create banker system: 2 bankers + 2/3 system (doubles and trebles from non-bankers)
	// This should create combinations of size 4 (2 bankers + 2 non-bankers) and 5 (2 bankers + 3 non-bankers)
	stake := NewStake("cash", "EUR", 1.00, "unit")
	
	builder.AddBankerSystemBet(bankers, []int{2, 3}, selections, stake)
	ticket := builder.Build("corr-banker-001")
	
	// Validate structure
	if len(ticket.Content.Bets) != 1 {
		t.Fatalf("Expected 1 bet, got %d", len(ticket.Content.Bets))
	}
	
	bet := ticket.Content.Bets[0]
	if len(bet.Selections) != 1 {
		t.Fatalf("Expected 1 root selection, got %d", len(bet.Selections))
	}
	
	rootSelection := bet.Selections[0]
	
	// Validate root system structure
	if rootSelection.Type != "system" {
		t.Errorf("Expected root selection type 'system', got '%s'", rootSelection.Type)
	}
	
	// Validate outer size: should be [4, 5] (2 bankers + 2, 2 bankers + 3)
	expectedOuterSize := []int{4, 5}
	if len(rootSelection.Size) != len(expectedOuterSize) {
		t.Errorf("Expected outer size length %d, got %d", len(expectedOuterSize), len(rootSelection.Size))
	} else {
		for i, expected := range expectedOuterSize {
			if rootSelection.Size[i] != expected {
				t.Errorf("Expected outer size[%d] = %d, got %d", i, expected, rootSelection.Size[i])
			}
		}
	}
	
	// Validate nested structure: should have 2 nested systems
	if len(rootSelection.Selections) != 2 {
		t.Fatalf("Expected 2 nested systems, got %d", len(rootSelection.Selections))
	}
	
	// Validate banker system
	bankerSystem := rootSelection.Selections[0]
	if bankerSystem.Type != "system" {
		t.Errorf("Expected banker system type 'system', got '%s'", bankerSystem.Type)
	}
	if len(bankerSystem.Size) != 1 || bankerSystem.Size[0] != 1 {
		t.Errorf("Expected banker system size [1], got %v", bankerSystem.Size)
	}
	if len(bankerSystem.Selections) != 2 {
		t.Errorf("Expected 2 banker selections, got %d", len(bankerSystem.Selections))
	}
	
	// Validate main system
	mainSystem := rootSelection.Selections[1]
	if mainSystem.Type != "system" {
		t.Errorf("Expected main system type 'system', got '%s'", mainSystem.Type)
	}
	expectedMainSize := []int{2, 3}
	if len(mainSystem.Size) != len(expectedMainSize) {
		t.Errorf("Expected main system size length %d, got %d", len(expectedMainSize), len(mainSystem.Size))
	}
	if len(mainSystem.Selections) != 3 {
		t.Errorf("Expected 3 non-banker selections, got %d", len(mainSystem.Selections))
	}
	
	// Test JSON serialization
	jsonData, err := json.MarshalIndent(ticket, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal ticket: %v", err)
	}
	
	t.Logf("Banker System Bet JSON:\n%s", string(jsonData))
	
	// Validate JSON structure
	var parsed map[string]interface{}
	if err := json.Unmarshal(jsonData, &parsed); err != nil {
		t.Fatalf("Failed to parse JSON: %v", err)
	}
	
	// Navigate to root selection
	content := parsed["content"].(map[string]interface{})
	bets := content["bets"].([]interface{})
	bet0 := bets[0].(map[string]interface{})
	selections := bet0["selections"].([]interface{})
	rootSel := selections[0].(map[string]interface{})
	
	// Verify nested selections exist
	nestedSels := rootSel["selections"].([]interface{})
	if len(nestedSels) != 2 {
		t.Errorf("Expected 2 nested selections in JSON, got %d", len(nestedSels))
	}
}

// TestBankerSystemBetWithSingleBanker tests banker system with 1 banker
func TestBankerSystemBetWithSingleBanker(t *testing.T) {
	builder := NewTicketBuilder(45426, "test-banker-single-001")
	
	// Create 1 banker
	banker := NewSelection("3", "sr:match:11111", "1", "1", 1.50)
	bankers := []Selection{banker}
	
	// Create 4 non-banker selections
	sel1 := NewSelection("3", "sr:match:22222", "1", "1", 2.00)
	sel2 := NewSelection("3", "sr:match:33333", "1", "1", 2.20)
	sel3 := NewSelection("3", "sr:match:44444", "1", "1", 2.50)
	sel4 := NewSelection("3", "sr:match:55555", "1", "1", 2.80)
	selections := []Selection{sel1, sel2, sel3, sel4}
	
	// Create Yankee with 1 banker: 1 banker + 2/3/4 system
	// Outer size should be [3, 4, 5] (1 banker + 2, 1 banker + 3, 1 banker + 4)
	stake := NewStake("cash", "EUR", 1.00, "unit")
	
	builder.AddBankerSystemBet(bankers, []int{2, 3, 4}, selections, stake)
	ticket := builder.Build("corr-banker-single-001")
	
	// Validate outer size
	bet := ticket.Content.Bets[0]
	rootSelection := bet.Selections[0]
	
	expectedOuterSize := []int{3, 4, 5}
	if len(rootSelection.Size) != len(expectedOuterSize) {
		t.Errorf("Expected outer size length %d, got %d", len(expectedOuterSize), len(rootSelection.Size))
	} else {
		for i, expected := range expectedOuterSize {
			if rootSelection.Size[i] != expected {
				t.Errorf("Expected outer size[%d] = %d, got %d", i, expected, rootSelection.Size[i])
			}
		}
	}
	
	// Test JSON serialization
	jsonData, err := json.MarshalIndent(ticket, "", "  ")
	if err != nil {
		t.Fatalf("Failed to marshal ticket: %v", err)
	}
	
	t.Logf("Single Banker System Bet JSON:\n%s", string(jsonData))
}
