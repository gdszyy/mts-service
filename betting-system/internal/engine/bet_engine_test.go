package engine

import (
	"testing"
)

func TestGenerateSingles(t *testing.T) {
	engine := NewBetEngine()

	selections := []Selection{
		{ID: 1, Odds: 2.0, IsBanker: false},
		{ID: 2, Odds: 3.0, IsBanker: false},
	}

	combinations, err := engine.GenerateCombinations(BetTypeSingle, selections, 10.0)
	if err != nil {
		t.Fatalf("Failed to generate combinations: %v", err)
	}

	if len(combinations) != 2 {
		t.Errorf("Expected 2 combinations, got %d", len(combinations))
	}

	// 验证第一个组合
	if combinations[0].Type != "single" {
		t.Errorf("Expected type 'single', got '%s'", combinations[0].Type)
	}
	if combinations[0].Stake != 10.0 {
		t.Errorf("Expected stake 10.0, got %.2f", combinations[0].Stake)
	}
	if combinations[0].TotalOdds != 2.0 {
		t.Errorf("Expected odds 2.0, got %.2f", combinations[0].TotalOdds)
	}
	if combinations[0].PotentialReturn != 20.0 {
		t.Errorf("Expected potential return 20.0, got %.2f", combinations[0].PotentialReturn)
	}
}

func TestGenerateAccumulator(t *testing.T) {
	engine := NewBetEngine()

	selections := []Selection{
		{ID: 1, Odds: 2.0, IsBanker: false},
		{ID: 2, Odds: 3.0, IsBanker: false},
		{ID: 3, Odds: 2.5, IsBanker: false},
	}

	combinations, err := engine.GenerateCombinations(BetTypeAccumulator, selections, 10.0)
	if err != nil {
		t.Fatalf("Failed to generate combinations: %v", err)
	}

	if len(combinations) != 1 {
		t.Errorf("Expected 1 combination, got %d", len(combinations))
	}

	// 验证组合
	if combinations[0].Type != "3-fold" {
		t.Errorf("Expected type '3-fold', got '%s'", combinations[0].Type)
	}
	if combinations[0].TotalOdds != 15.0 { // 2.0 * 3.0 * 2.5
		t.Errorf("Expected odds 15.0, got %.2f", combinations[0].TotalOdds)
	}
	if combinations[0].PotentialReturn != 150.0 {
		t.Errorf("Expected potential return 150.0, got %.2f", combinations[0].PotentialReturn)
	}
}

func TestGenerateTrixie(t *testing.T) {
	engine := NewBetEngine()

	selections := []Selection{
		{ID: 1, Odds: 2.0, IsBanker: false},
		{ID: 2, Odds: 3.0, IsBanker: false},
		{ID: 3, Odds: 2.5, IsBanker: false},
	}

	combinations, err := engine.GenerateCombinations(BetTypeTrixie, selections, 10.0)
	if err != nil {
		t.Fatalf("Failed to generate combinations: %v", err)
	}

	// Trixie = 3 doubles + 1 treble = 4 bets
	if len(combinations) != 4 {
		t.Errorf("Expected 4 combinations, got %d", len(combinations))
	}

	// 验证总投注金额
	totalStake := engine.CalculateTotalStake(combinations)
	if totalStake != 40.0 {
		t.Errorf("Expected total stake 40.0, got %.2f", totalStake)
	}

	// 验证组合类型
	doubleCount := 0
	trebleCount := 0
	for _, combo := range combinations {
		if combo.Type == "double" {
			doubleCount++
		} else if combo.Type == "treble" {
			trebleCount++
		}
	}

	if doubleCount != 3 {
		t.Errorf("Expected 3 doubles, got %d", doubleCount)
	}
	if trebleCount != 1 {
		t.Errorf("Expected 1 treble, got %d", trebleCount)
	}
}

func TestGenerateTrixieWithBanker(t *testing.T) {
	engine := NewBetEngine()

	selections := []Selection{
		{ID: 1, Odds: 2.0, IsBanker: true},  // Banker
		{ID: 2, Odds: 3.0, IsBanker: false},
		{ID: 3, Odds: 2.5, IsBanker: false},
	}

	combinations, err := engine.GenerateCombinations(BetTypeTrixie, selections, 10.0)
	if err != nil {
		t.Fatalf("Failed to generate combinations: %v", err)
	}

	// Trixie with 1 Banker = 2 doubles + 1 treble = 3 bets
	// (Banker+2), (Banker+3), (Banker+2+3)
	if len(combinations) != 3 {
		t.Errorf("Expected 3 combinations, got %d", len(combinations))
	}

	// 验证所有组合都包含 Banker
	for _, combo := range combinations {
		hasBanker := false
		for _, sel := range combo.Selections {
			if sel.ID == 1 {
				hasBanker = true
				break
			}
		}
		if !hasBanker {
			t.Errorf("Combination does not include Banker: %+v", combo)
		}
	}
}

func TestGeneratePatent(t *testing.T) {
	engine := NewBetEngine()

	selections := []Selection{
		{ID: 1, Odds: 2.0, IsBanker: false},
		{ID: 2, Odds: 3.0, IsBanker: false},
		{ID: 3, Odds: 2.5, IsBanker: false},
	}

	combinations, err := engine.GenerateCombinations(BetTypePatent, selections, 10.0)
	if err != nil {
		t.Fatalf("Failed to generate combinations: %v", err)
	}

	// Patent = 3 singles + 3 doubles + 1 treble = 7 bets
	if len(combinations) != 7 {
		t.Errorf("Expected 7 combinations, got %d", len(combinations))
	}

	// 验证总投注金额
	totalStake := engine.CalculateTotalStake(combinations)
	if totalStake != 70.0 {
		t.Errorf("Expected total stake 70.0, got %.2f", totalStake)
	}

	// 验证组合类型
	singleCount := 0
	doubleCount := 0
	trebleCount := 0
	for _, combo := range combinations {
		if combo.Type == "single" {
			singleCount++
		} else if combo.Type == "double" {
			doubleCount++
		} else if combo.Type == "treble" {
			trebleCount++
		}
	}

	if singleCount != 3 {
		t.Errorf("Expected 3 singles, got %d", singleCount)
	}
	if doubleCount != 3 {
		t.Errorf("Expected 3 doubles, got %d", doubleCount)
	}
	if trebleCount != 1 {
		t.Errorf("Expected 1 treble, got %d", trebleCount)
	}
}

func TestGenerateYankee(t *testing.T) {
	engine := NewBetEngine()

	selections := []Selection{
		{ID: 1, Odds: 2.0, IsBanker: false},
		{ID: 2, Odds: 3.0, IsBanker: false},
		{ID: 3, Odds: 2.5, IsBanker: false},
		{ID: 4, Odds: 1.8, IsBanker: false},
	}

	combinations, err := engine.GenerateCombinations(BetTypeYankee, selections, 10.0)
	if err != nil {
		t.Fatalf("Failed to generate combinations: %v", err)
	}

	// Yankee = 6 doubles + 4 trebles + 1 4-fold = 11 bets
	if len(combinations) != 11 {
		t.Errorf("Expected 11 combinations, got %d", len(combinations))
	}

	// 验证总投注金额
	totalStake := engine.CalculateTotalStake(combinations)
	if totalStake != 110.0 {
		t.Errorf("Expected total stake 110.0, got %.2f", totalStake)
	}
}

func TestGenerateLucky15(t *testing.T) {
	engine := NewBetEngine()

	selections := []Selection{
		{ID: 1, Odds: 2.0, IsBanker: false},
		{ID: 2, Odds: 3.0, IsBanker: false},
		{ID: 3, Odds: 2.5, IsBanker: false},
		{ID: 4, Odds: 1.8, IsBanker: false},
	}

	combinations, err := engine.GenerateCombinations(BetTypeLucky15, selections, 10.0)
	if err != nil {
		t.Fatalf("Failed to generate combinations: %v", err)
	}

	// Lucky 15 = 4 singles + 6 doubles + 4 trebles + 1 4-fold = 15 bets
	if len(combinations) != 15 {
		t.Errorf("Expected 15 combinations, got %d", len(combinations))
	}

	// 验证总投注金额
	totalStake := engine.CalculateTotalStake(combinations)
	if totalStake != 150.0 {
		t.Errorf("Expected total stake 150.0, got %.2f", totalStake)
	}
}

func TestValidateSelections(t *testing.T) {
	engine := NewBetEngine()

	// 测试选项数量不足
	selections := []Selection{
		{ID: 1, Odds: 2.0, IsBanker: false},
		{ID: 2, Odds: 3.0, IsBanker: false},
	}

	_, err := engine.GenerateCombinations(BetTypeTrixie, selections, 10.0)
	if err == nil {
		t.Error("Expected error for insufficient selections, got nil")
	}

	// 测试未知投注类型
	_, err = engine.GenerateCombinations(BetType("unknown"), selections, 10.0)
	if err == nil {
		t.Error("Expected error for unknown bet type, got nil")
	}
}

func TestCalculateTotalPotentialReturn(t *testing.T) {
	engine := NewBetEngine()

	selections := []Selection{
		{ID: 1, Odds: 2.0, IsBanker: false},
		{ID: 2, Odds: 3.0, IsBanker: false},
		{ID: 3, Odds: 2.5, IsBanker: false},
	}

	combinations, err := engine.GenerateCombinations(BetTypeTrixie, selections, 10.0)
	if err != nil {
		t.Fatalf("Failed to generate combinations: %v", err)
	}

	totalReturn := engine.CalculateTotalPotentialReturn(combinations)

	// Trixie 潜在回报 = (2*3 + 2*2.5 + 3*2.5 + 2*3*2.5) * 10
	// = (6 + 5 + 7.5 + 15) * 10 = 335
	expectedReturn := 335.0
	if totalReturn != expectedReturn {
		t.Errorf("Expected total potential return %.2f, got %.2f", expectedReturn, totalReturn)
	}
}

