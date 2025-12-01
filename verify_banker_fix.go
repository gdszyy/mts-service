package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/gdsZyy/mts-service/internal/models"
)

func main() {
	log.Println("Verifying Banker System Bet Fix...")
	log.Println("=" + "==========================================================")

	// Create a banker system bet using the fixed method
	builder := models.NewTicketBuilder(9985, "test-banker-fix-001")

	// Create banker selections
	banker := models.NewSelection("3", "sr:match:12345", "1", "1", "1.50")

	// Create non-banker selections
	sel1 := models.NewSelection("3", "sr:match:12346", "1", "2", "2.50")
	sel2 := models.NewSelection("3", "sr:match:12347", "1", "1", "3.00")
	sel3 := models.NewSelection("3", "sr:match:12348", "1", "2", "2.20")

	// Create stake
	stake := models.NewStake("cash", "EUR", "1.00", "unit")

	// Add banker system bet (3/4 with 1 banker = 2/3 system + 1 banker)
	builder.AddBankerSystemBet(
		[]models.Selection{banker},
		[]int{2}, // 2/3 system
		[]models.Selection{sel1, sel2, sel3},
		stake,
	)

	// Build ticket
	ticket := builder.Build("test-correlation-id")

	// Marshal to JSON
	jsonData, err := json.MarshalIndent(ticket, "", "  ")
	if err != nil {
		log.Fatalf("Failed to marshal ticket: %v", err)
	}

	fmt.Println("\nGenerated Ticket JSON:")
	fmt.Println("=" + "==========================================================")
	fmt.Println(string(jsonData))
	fmt.Println("=" + "==========================================================")

	// Verify structure
	log.Println("\nVerifying structure...")
	
	if len(ticket.Content.Bets) != 1 {
		log.Fatalf("Expected 1 bet, got %d", len(ticket.Content.Bets))
	}
	
	bet := ticket.Content.Bets[0]
	if len(bet.Selections) != 2 {
		log.Fatalf("Expected 2 top-level selections, got %d", len(bet.Selections))
	}
	
	// First selection should be system type
	if bet.Selections[0].Type != "system" {
		log.Fatalf("Expected first selection to be 'system', got '%s'", bet.Selections[0].Type)
	}
	
	// Second selection should be uf type (banker)
	if bet.Selections[1].Type != "uf" {
		log.Fatalf("Expected second selection to be 'uf', got '%s'", bet.Selections[1].Type)
	}
	
	// System selection should have 3 nested selections
	if len(bet.Selections[0].Selections) != 3 {
		log.Fatalf("Expected system selection to have 3 nested selections, got %d", len(bet.Selections[0].Selections))
	}
	
	// System selection should have size [2]
	if len(bet.Selections[0].Size) != 1 || bet.Selections[0].Size[0] != 2 {
		log.Fatalf("Expected system selection size to be [2], got %v", bet.Selections[0].Size)
	}
	
	log.Println("✓ Structure verification passed!")
	log.Println("✓ Banker system bet structure matches MTS documentation!")
	log.Println("\nThe fix is correct. Deploy the updated code to production to test with real MTS.")
}
