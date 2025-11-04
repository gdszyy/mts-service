package engine

import (
	"fmt"
	"math"
)

// BetType 定义投注类型
type BetType string

const (
	BetTypeSingle       BetType = "single"
	BetTypeAccumulator  BetType = "accumulator"
	BetTypeTrixie       BetType = "trixie"
	BetTypePatent       BetType = "patent"
	BetTypeYankee       BetType = "yankee"
	BetTypeLucky15      BetType = "lucky_15"
	BetTypeSuperYankee  BetType = "super_yankee"
	BetTypeLucky31      BetType = "lucky_31"
	BetTypeHeinz        BetType = "heinz"
	BetTypeLucky63      BetType = "lucky_63"
	BetTypeSuperHeinz   BetType = "super_heinz"
	BetTypeGoliath      BetType = "goliath"
)

// Selection 表示一个投注选项
type Selection struct {
	ID       int64
	Odds     float64
	IsBanker bool
}

// Combination 表示一个投注组合
type Combination struct {
	Type            string      // single/double/treble/4-fold等
	Selections      []Selection
	Stake           float64
	TotalOdds       float64
	PotentialReturn float64
}

// BetEngine 投注引擎
type BetEngine struct{}

// NewBetEngine 创建新的投注引擎
func NewBetEngine() *BetEngine {
	return &BetEngine{}
}

// GenerateCombinations 生成投注组合
func (e *BetEngine) GenerateCombinations(betType BetType, selections []Selection, unitStake float64) ([]Combination, error) {
	// 验证选项数量
	if err := e.validateSelections(betType, selections); err != nil {
		return nil, err
	}

	// 分离 Banker 和非 Banker 选项
	bankers, nonBankers := e.separateBankers(selections)

	var combinations []Combination

	switch betType {
	case BetTypeSingle:
		combinations = e.generateSingles(selections, unitStake)
	case BetTypeAccumulator:
		combinations = e.generateAccumulator(selections, unitStake)
	case BetTypeTrixie:
		combinations = e.generateTrixie(bankers, nonBankers, unitStake)
	case BetTypePatent:
		combinations = e.generatePatent(bankers, nonBankers, unitStake)
	case BetTypeYankee:
		combinations = e.generateYankee(bankers, nonBankers, unitStake)
	case BetTypeLucky15:
		combinations = e.generateLucky15(bankers, nonBankers, unitStake)
	case BetTypeSuperYankee:
		combinations = e.generateSuperYankee(bankers, nonBankers, unitStake)
	case BetTypeLucky31:
		combinations = e.generateLucky31(bankers, nonBankers, unitStake)
	case BetTypeHeinz:
		combinations = e.generateHeinz(bankers, nonBankers, unitStake)
	case BetTypeLucky63:
		combinations = e.generateLucky63(bankers, nonBankers, unitStake)
	case BetTypeSuperHeinz:
		combinations = e.generateSuperHeinz(bankers, nonBankers, unitStake)
	case BetTypeGoliath:
		combinations = e.generateGoliath(bankers, nonBankers, unitStake)
	default:
		return nil, fmt.Errorf("unsupported bet type: %s", betType)
	}

	return combinations, nil
}

// validateSelections 验证选项数量
func (e *BetEngine) validateSelections(betType BetType, selections []Selection) error {
	required := map[BetType]int{
		BetTypeSingle:      1,
		BetTypeAccumulator: 2,
		BetTypeTrixie:      3,
		BetTypePatent:      3,
		BetTypeYankee:      4,
		BetTypeLucky15:     4,
		BetTypeSuperYankee: 5,
		BetTypeLucky31:     5,
		BetTypeHeinz:       6,
		BetTypeLucky63:     6,
		BetTypeSuperHeinz:  7,
		BetTypeGoliath:     8,
	}

	minRequired, ok := required[betType]
	if !ok {
		return fmt.Errorf("unknown bet type: %s", betType)
	}

	if len(selections) < minRequired {
		return fmt.Errorf("bet type %s requires at least %d selections, got %d", betType, minRequired, len(selections))
	}

	return nil
}

// separateBankers 分离 Banker 和非 Banker 选项
func (e *BetEngine) separateBankers(selections []Selection) (bankers, nonBankers []Selection) {
	for _, sel := range selections {
		if sel.IsBanker {
			bankers = append(bankers, sel)
		} else {
			nonBankers = append(nonBankers, sel)
		}
	}
	return
}

// generateSingles 生成单注
func (e *BetEngine) generateSingles(selections []Selection, unitStake float64) []Combination {
	var combinations []Combination
	for _, sel := range selections {
		combinations = append(combinations, Combination{
			Type:            "single",
			Selections:      []Selection{sel},
			Stake:           unitStake,
			TotalOdds:       sel.Odds,
			PotentialReturn: unitStake * sel.Odds,
		})
	}
	return combinations
}

// generateAccumulator 生成串关
func (e *BetEngine) generateAccumulator(selections []Selection, unitStake float64) []Combination {
	totalOdds := 1.0
	for _, sel := range selections {
		totalOdds *= sel.Odds
	}

	return []Combination{{
		Type:            fmt.Sprintf("%d-fold", len(selections)),
		Selections:      selections,
		Stake:           unitStake,
		TotalOdds:       totalOdds,
		PotentialReturn: unitStake * totalOdds,
	}}
}

// generateTrixie 生成 Trixie (3 doubles + 1 treble = 4 bets)
func (e *BetEngine) generateTrixie(bankers, nonBankers []Selection, unitStake float64) []Combination {
	var combinations []Combination

	// 生成所有可能的选项（Banker + 非Banker）
	allSelections := append(bankers, nonBankers...)

	// 如果有 Banker，则只对非 Banker 进行组合，并在每个组合中包含所有 Banker
	if len(bankers) > 0 {
		// 生成所有包含 Banker 的 2-fold 组合
		if len(nonBankers) >= 1 {
			combinations = append(combinations, e.generateCombinationsWithBankers(bankers, nonBankers, 1, unitStake)...)
		}
		// 生成所有包含 Banker 的 3-fold 组合
		if len(nonBankers) >= 2 {
			combinations = append(combinations, e.generateCombinationsWithBankers(bankers, nonBankers, 2, unitStake)...)
		}
	} else {
		// 无 Banker，正常生成组合
		// 3 doubles
		combinations = append(combinations, e.generateNCombinations(allSelections, 2, unitStake)...)
		// 1 treble
		combinations = append(combinations, e.generateNCombinations(allSelections, 3, unitStake)...)
	}

	return combinations
}

// generatePatent 生成 Patent (3 singles + 3 doubles + 1 treble = 7 bets)
func (e *BetEngine) generatePatent(bankers, nonBankers []Selection, unitStake float64) []Combination {
	var combinations []Combination

	allSelections := append(bankers, nonBankers...)

	if len(bankers) > 0 {
		// 生成 singles（每个非Banker + 所有Banker）
		for _, nb := range nonBankers {
			sels := append([]Selection{}, bankers...)
			sels = append(sels, nb)
			combinations = append(combinations, e.createCombination(sels, unitStake))
		}
		// 生成 doubles
		combinations = append(combinations, e.generateCombinationsWithBankers(bankers, nonBankers, 2, unitStake)...)
		// 生成 treble
		combinations = append(combinations, e.generateCombinationsWithBankers(bankers, nonBankers, 3, unitStake)...)
	} else {
		// 3 singles
		combinations = append(combinations, e.generateSingles(allSelections, unitStake)...)
		// 3 doubles
		combinations = append(combinations, e.generateNCombinations(allSelections, 2, unitStake)...)
		// 1 treble
		combinations = append(combinations, e.generateNCombinations(allSelections, 3, unitStake)...)
	}

	return combinations
}

// generateYankee 生成 Yankee (6 doubles + 4 trebles + 1 4-fold = 11 bets)
func (e *BetEngine) generateYankee(bankers, nonBankers []Selection, unitStake float64) []Combination {
	var combinations []Combination

	allSelections := append(bankers, nonBankers...)

	if len(bankers) > 0 {
		combinations = append(combinations, e.generateCombinationsWithBankers(bankers, nonBankers, 2, unitStake)...)
		combinations = append(combinations, e.generateCombinationsWithBankers(bankers, nonBankers, 3, unitStake)...)
		combinations = append(combinations, e.generateCombinationsWithBankers(bankers, nonBankers, 4, unitStake)...)
	} else {
		combinations = append(combinations, e.generateNCombinations(allSelections, 2, unitStake)...)
		combinations = append(combinations, e.generateNCombinations(allSelections, 3, unitStake)...)
		combinations = append(combinations, e.generateNCombinations(allSelections, 4, unitStake)...)
	}

	return combinations
}

// generateLucky15 生成 Lucky 15 (4 singles + 6 doubles + 4 trebles + 1 4-fold = 15 bets)
func (e *BetEngine) generateLucky15(bankers, nonBankers []Selection, unitStake float64) []Combination {
	var combinations []Combination

	allSelections := append(bankers, nonBankers...)

	if len(bankers) > 0 {
		for _, nb := range nonBankers {
			sels := append([]Selection{}, bankers...)
			sels = append(sels, nb)
			combinations = append(combinations, e.createCombination(sels, unitStake))
		}
		combinations = append(combinations, e.generateCombinationsWithBankers(bankers, nonBankers, 2, unitStake)...)
		combinations = append(combinations, e.generateCombinationsWithBankers(bankers, nonBankers, 3, unitStake)...)
		combinations = append(combinations, e.generateCombinationsWithBankers(bankers, nonBankers, 4, unitStake)...)
	} else {
		combinations = append(combinations, e.generateSingles(allSelections, unitStake)...)
		combinations = append(combinations, e.generateNCombinations(allSelections, 2, unitStake)...)
		combinations = append(combinations, e.generateNCombinations(allSelections, 3, unitStake)...)
		combinations = append(combinations, e.generateNCombinations(allSelections, 4, unitStake)...)
	}

	return combinations
}

// generateSuperYankee 生成 Super Yankee/Canadian (10 doubles + 10 trebles + 5 4-folds + 1 5-fold = 26 bets)
func (e *BetEngine) generateSuperYankee(bankers, nonBankers []Selection, unitStake float64) []Combination {
	var combinations []Combination

	allSelections := append(bankers, nonBankers...)

	if len(bankers) > 0 {
		combinations = append(combinations, e.generateCombinationsWithBankers(bankers, nonBankers, 2, unitStake)...)
		combinations = append(combinations, e.generateCombinationsWithBankers(bankers, nonBankers, 3, unitStake)...)
		combinations = append(combinations, e.generateCombinationsWithBankers(bankers, nonBankers, 4, unitStake)...)
		combinations = append(combinations, e.generateCombinationsWithBankers(bankers, nonBankers, 5, unitStake)...)
	} else {
		combinations = append(combinations, e.generateNCombinations(allSelections, 2, unitStake)...)
		combinations = append(combinations, e.generateNCombinations(allSelections, 3, unitStake)...)
		combinations = append(combinations, e.generateNCombinations(allSelections, 4, unitStake)...)
		combinations = append(combinations, e.generateNCombinations(allSelections, 5, unitStake)...)
	}

	return combinations
}

// generateLucky31 生成 Lucky 31
func (e *BetEngine) generateLucky31(bankers, nonBankers []Selection, unitStake float64) []Combination {
	var combinations []Combination

	allSelections := append(bankers, nonBankers...)

	if len(bankers) > 0 {
		for _, nb := range nonBankers {
			sels := append([]Selection{}, bankers...)
			sels = append(sels, nb)
			combinations = append(combinations, e.createCombination(sels, unitStake))
		}
		combinations = append(combinations, e.generateCombinationsWithBankers(bankers, nonBankers, 2, unitStake)...)
		combinations = append(combinations, e.generateCombinationsWithBankers(bankers, nonBankers, 3, unitStake)...)
		combinations = append(combinations, e.generateCombinationsWithBankers(bankers, nonBankers, 4, unitStake)...)
		combinations = append(combinations, e.generateCombinationsWithBankers(bankers, nonBankers, 5, unitStake)...)
	} else {
		combinations = append(combinations, e.generateSingles(allSelections, unitStake)...)
		combinations = append(combinations, e.generateNCombinations(allSelections, 2, unitStake)...)
		combinations = append(combinations, e.generateNCombinations(allSelections, 3, unitStake)...)
		combinations = append(combinations, e.generateNCombinations(allSelections, 4, unitStake)...)
		combinations = append(combinations, e.generateNCombinations(allSelections, 5, unitStake)...)
	}

	return combinations
}

// generateHeinz 生成 Heinz
func (e *BetEngine) generateHeinz(bankers, nonBankers []Selection, unitStake float64) []Combination {
	var combinations []Combination

	allSelections := append(bankers, nonBankers...)

	if len(bankers) > 0 {
		combinations = append(combinations, e.generateCombinationsWithBankers(bankers, nonBankers, 2, unitStake)...)
		combinations = append(combinations, e.generateCombinationsWithBankers(bankers, nonBankers, 3, unitStake)...)
		combinations = append(combinations, e.generateCombinationsWithBankers(bankers, nonBankers, 4, unitStake)...)
		combinations = append(combinations, e.generateCombinationsWithBankers(bankers, nonBankers, 5, unitStake)...)
		combinations = append(combinations, e.generateCombinationsWithBankers(bankers, nonBankers, 6, unitStake)...)
	} else {
		combinations = append(combinations, e.generateNCombinations(allSelections, 2, unitStake)...)
		combinations = append(combinations, e.generateNCombinations(allSelections, 3, unitStake)...)
		combinations = append(combinations, e.generateNCombinations(allSelections, 4, unitStake)...)
		combinations = append(combinations, e.generateNCombinations(allSelections, 5, unitStake)...)
		combinations = append(combinations, e.generateNCombinations(allSelections, 6, unitStake)...)
	}

	return combinations
}

// generateLucky63 生成 Lucky 63
func (e *BetEngine) generateLucky63(bankers, nonBankers []Selection, unitStake float64) []Combination {
	var combinations []Combination

	allSelections := append(bankers, nonBankers...)

	if len(bankers) > 0 {
		for _, nb := range nonBankers {
			sels := append([]Selection{}, bankers...)
			sels = append(sels, nb)
			combinations = append(combinations, e.createCombination(sels, unitStake))
		}
		combinations = append(combinations, e.generateCombinationsWithBankers(bankers, nonBankers, 2, unitStake)...)
		combinations = append(combinations, e.generateCombinationsWithBankers(bankers, nonBankers, 3, unitStake)...)
		combinations = append(combinations, e.generateCombinationsWithBankers(bankers, nonBankers, 4, unitStake)...)
		combinations = append(combinations, e.generateCombinationsWithBankers(bankers, nonBankers, 5, unitStake)...)
		combinations = append(combinations, e.generateCombinationsWithBankers(bankers, nonBankers, 6, unitStake)...)
	} else {
		combinations = append(combinations, e.generateSingles(allSelections, unitStake)...)
		combinations = append(combinations, e.generateNCombinations(allSelections, 2, unitStake)...)
		combinations = append(combinations, e.generateNCombinations(allSelections, 3, unitStake)...)
		combinations = append(combinations, e.generateNCombinations(allSelections, 4, unitStake)...)
		combinations = append(combinations, e.generateNCombinations(allSelections, 5, unitStake)...)
		combinations = append(combinations, e.generateNCombinations(allSelections, 6, unitStake)...)
	}

	return combinations
}

// generateSuperHeinz 生成 Super Heinz
func (e *BetEngine) generateSuperHeinz(bankers, nonBankers []Selection, unitStake float64) []Combination {
	var combinations []Combination

	allSelections := append(bankers, nonBankers...)

	if len(bankers) > 0 {
		combinations = append(combinations, e.generateCombinationsWithBankers(bankers, nonBankers, 2, unitStake)...)
		combinations = append(combinations, e.generateCombinationsWithBankers(bankers, nonBankers, 3, unitStake)...)
		combinations = append(combinations, e.generateCombinationsWithBankers(bankers, nonBankers, 4, unitStake)...)
		combinations = append(combinations, e.generateCombinationsWithBankers(bankers, nonBankers, 5, unitStake)...)
		combinations = append(combinations, e.generateCombinationsWithBankers(bankers, nonBankers, 6, unitStake)...)
		combinations = append(combinations, e.generateCombinationsWithBankers(bankers, nonBankers, 7, unitStake)...)
	} else {
		combinations = append(combinations, e.generateNCombinations(allSelections, 2, unitStake)...)
		combinations = append(combinations, e.generateNCombinations(allSelections, 3, unitStake)...)
		combinations = append(combinations, e.generateNCombinations(allSelections, 4, unitStake)...)
		combinations = append(combinations, e.generateNCombinations(allSelections, 5, unitStake)...)
		combinations = append(combinations, e.generateNCombinations(allSelections, 6, unitStake)...)
		combinations = append(combinations, e.generateNCombinations(allSelections, 7, unitStake)...)
	}

	return combinations
}

// generateGoliath 生成 Goliath
func (e *BetEngine) generateGoliath(bankers, nonBankers []Selection, unitStake float64) []Combination {
	var combinations []Combination

	allSelections := append(bankers, nonBankers...)

	if len(bankers) > 0 {
		combinations = append(combinations, e.generateCombinationsWithBankers(bankers, nonBankers, 2, unitStake)...)
		combinations = append(combinations, e.generateCombinationsWithBankers(bankers, nonBankers, 3, unitStake)...)
		combinations = append(combinations, e.generateCombinationsWithBankers(bankers, nonBankers, 4, unitStake)...)
		combinations = append(combinations, e.generateCombinationsWithBankers(bankers, nonBankers, 5, unitStake)...)
		combinations = append(combinations, e.generateCombinationsWithBankers(bankers, nonBankers, 6, unitStake)...)
		combinations = append(combinations, e.generateCombinationsWithBankers(bankers, nonBankers, 7, unitStake)...)
		combinations = append(combinations, e.generateCombinationsWithBankers(bankers, nonBankers, 8, unitStake)...)
	} else {
		combinations = append(combinations, e.generateNCombinations(allSelections, 2, unitStake)...)
		combinations = append(combinations, e.generateNCombinations(allSelections, 3, unitStake)...)
		combinations = append(combinations, e.generateNCombinations(allSelections, 4, unitStake)...)
		combinations = append(combinations, e.generateNCombinations(allSelections, 5, unitStake)...)
		combinations = append(combinations, e.generateNCombinations(allSelections, 6, unitStake)...)
		combinations = append(combinations, e.generateNCombinations(allSelections, 7, unitStake)...)
		combinations = append(combinations, e.generateNCombinations(allSelections, 8, unitStake)...)
	}

	return combinations
}

// generateNCombinations 生成 n 个选项的所有组合
func (e *BetEngine) generateNCombinations(selections []Selection, n int, unitStake float64) []Combination {
	var combinations []Combination

	// 生成所有 n 个选项的组合
	combos := e.getCombinations(selections, n)

	for _, combo := range combos {
		combinations = append(combinations, e.createCombination(combo, unitStake))
	}

	return combinations
}

// generateCombinationsWithBankers 生成包含 Banker 的组合
func (e *BetEngine) generateCombinationsWithBankers(bankers, nonBankers []Selection, n int, unitStake float64) []Combination {
	var combinations []Combination

	// 从非 Banker 中选择 n 个，然后加上所有 Banker
	combos := e.getCombinations(nonBankers, n)

	for _, combo := range combos {
		// 将所有 Banker 加入组合
		fullCombo := append([]Selection{}, bankers...)
		fullCombo = append(fullCombo, combo...)
		combinations = append(combinations, e.createCombination(fullCombo, unitStake))
	}

	return combinations
}

// getCombinations 获取所有 n 个元素的组合
func (e *BetEngine) getCombinations(selections []Selection, n int) [][]Selection {
	var result [][]Selection

	if n == 0 {
		return [][]Selection{{}}
	}

	if len(selections) < n {
		return result
	}

	// 递归生成组合
	e.generateCombinationsRecursive(selections, n, 0, []Selection{}, &result)

	return result
}

// generateCombinationsRecursive 递归生成组合
func (e *BetEngine) generateCombinationsRecursive(selections []Selection, n, start int, current []Selection, result *[][]Selection) {
	if len(current) == n {
		combo := make([]Selection, len(current))
		copy(combo, current)
		*result = append(*result, combo)
		return
	}

	for i := start; i < len(selections); i++ {
		current = append(current, selections[i])
		e.generateCombinationsRecursive(selections, n, i+1, current, result)
		current = current[:len(current)-1]
	}
}

// createCombination 创建一个组合
func (e *BetEngine) createCombination(selections []Selection, unitStake float64) Combination {
	totalOdds := 1.0
	for _, sel := range selections {
		totalOdds *= sel.Odds
	}

	// 四舍五入到小数点后2位
	totalOdds = math.Round(totalOdds*100) / 100
	potentialReturn := math.Round(unitStake*totalOdds*100) / 100

	combType := ""
	switch len(selections) {
	case 1:
		combType = "single"
	case 2:
		combType = "double"
	case 3:
		combType = "treble"
	default:
		combType = fmt.Sprintf("%d-fold", len(selections))
	}

	return Combination{
		Type:            combType,
		Selections:      selections,
		Stake:           unitStake,
		TotalOdds:       totalOdds,
		PotentialReturn: potentialReturn,
	}
}

// CalculateTotalStake 计算总投注金额
func (e *BetEngine) CalculateTotalStake(combinations []Combination) float64 {
	total := 0.0
	for _, combo := range combinations {
		total += combo.Stake
	}
	return math.Round(total*100) / 100
}

// CalculateTotalPotentialReturn 计算总潜在回报
func (e *BetEngine) CalculateTotalPotentialReturn(combinations []Combination) float64 {
	total := 0.0
	for _, combo := range combinations {
		total += combo.PotentialReturn
	}
	return math.Round(total*100) / 100
}

