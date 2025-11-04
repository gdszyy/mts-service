package service

import (
	"fmt"
	"time"

	"github.com/gdsZyy/betting-system/internal/database"
	"github.com/gdsZyy/betting-system/internal/engine"
	"github.com/gdsZyy/betting-system/internal/models"
	"gorm.io/gorm"
)

// BetService 投注服务
type BetService struct {
	engine *engine.BetEngine
}

// NewBetService 创建新的投注服务
func NewBetService() *BetService {
	return &BetService{
		engine: engine.NewBetEngine(),
	}
}

// PlaceBetRequest 下注请求
type PlaceBetRequest struct {
	UserID     int64                   `json:"user_id" binding:"required"`
	BetType    string                  `json:"bet_type" binding:"required"`
	UnitStake  float64                 `json:"unit_stake" binding:"required,gt=0"`
	Selections []PlaceBetSelectionRequest `json:"selections" binding:"required,min=1"`
}

// PlaceBetSelectionRequest 下注选项请求
type PlaceBetSelectionRequest struct {
	OutcomeID int64   `json:"outcome_id" binding:"required"`
	Odds      float64 `json:"odds" binding:"required,gt=1"`
	IsBanker  bool    `json:"is_banker"`
}

// PlaceBet 下注
func (s *BetService) PlaceBet(req *PlaceBetRequest) (*models.Bet, error) {
	// 1. 验证用户
	var user models.User
	if err := database.DB.First(&user, req.UserID).Error; err != nil {
		return nil, fmt.Errorf("user not found: %w", err)
	}

	if user.Status != "active" {
		return nil, fmt.Errorf("user account is not active")
	}

	// 2. 验证所有 outcomes 存在且可用
	outcomeIDs := make([]int64, len(req.Selections))
	for i, sel := range req.Selections {
		outcomeIDs[i] = sel.OutcomeID
	}

	var outcomes []models.Outcome
	if err := database.DB.Preload("Market.Event").Where("id IN ?", outcomeIDs).Find(&outcomes).Error; err != nil {
		return nil, fmt.Errorf("failed to fetch outcomes: %w", err)
	}

	if len(outcomes) != len(outcomeIDs) {
		return nil, fmt.Errorf("some outcomes not found")
	}

	// 验证所有 outcomes 都是 active 状态
	for _, outcome := range outcomes {
		if outcome.Status != "active" {
			return nil, fmt.Errorf("outcome %d is not active", outcome.ID)
		}
		if outcome.Market.Status != "active" {
			return nil, fmt.Errorf("market %d is not active", outcome.MarketID)
		}
		if outcome.Market.Event.Status != "scheduled" && outcome.Market.Event.Status != "live" {
			return nil, fmt.Errorf("event %d is not available for betting", outcome.Market.EventID)
		}
	}

	// 3. 构建引擎选项
	engineSelections := make([]engine.Selection, len(req.Selections))
	for i, sel := range req.Selections {
		engineSelections[i] = engine.Selection{
			ID:       sel.OutcomeID,
			Odds:     sel.Odds,
			IsBanker: sel.IsBanker,
		}
	}

	// 4. 生成投注组合
	combinations, err := s.engine.GenerateCombinations(engine.BetType(req.BetType), engineSelections, req.UnitStake)
	if err != nil {
		return nil, fmt.Errorf("failed to generate combinations: %w", err)
	}

	// 5. 计算总投注金额和潜在回报
	totalStake := s.engine.CalculateTotalStake(combinations)
	totalPotentialReturn := s.engine.CalculateTotalPotentialReturn(combinations)

	// 6. 验证用户余额
	if user.Balance < totalStake {
		return nil, fmt.Errorf("insufficient balance: have %.2f, need %.2f", user.Balance, totalStake)
	}

	// 7. 开始数据库事务
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 8. 扣除用户余额
	if err := tx.Model(&user).Update("balance", gorm.Expr("balance - ?", totalStake)).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to deduct balance: %w", err)
	}

	// 9. 创建投注记录
	bet := &models.Bet{
		UserID:          req.UserID,
		BetType:         req.BetType,
		TotalStake:      totalStake,
		PotentialReturn: totalPotentialReturn,
		Status:          "pending",
		PlacedAt:        time.Now(),
	}

	if err := tx.Create(bet).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create bet: %w", err)
	}

	// 10. 创建投注选项
	betSelections := make([]models.BetSelection, len(req.Selections))
	for i, sel := range req.Selections {
		betSelections[i] = models.BetSelection{
			BetID:     bet.ID,
			OutcomeID: sel.OutcomeID,
			Odds:      sel.Odds,
			IsBanker:  sel.IsBanker,
			Status:    "pending",
		}
	}

	if err := tx.Create(&betSelections).Error; err != nil {
		tx.Rollback()
		return nil, fmt.Errorf("failed to create bet selections: %w", err)
	}

	// 11. 创建投注组合腿
	for _, combo := range combinations {
		betLeg := models.BetLeg{
			BetID:           bet.ID,
			LegType:         combo.Type,
			Stake:           combo.Stake,
			Odds:            combo.TotalOdds,
			PotentialReturn: combo.PotentialReturn,
			Status:          "pending",
		}

		if err := tx.Create(&betLeg).Error; err != nil {
			tx.Rollback()
			return nil, fmt.Errorf("failed to create bet leg: %w", err)
		}

		// 12. 创建组合腿选项关联
		for _, sel := range combo.Selections {
			// 找到对应的 BetSelection
			var betSelectionID int64
			for _, bs := range betSelections {
				if bs.OutcomeID == sel.ID {
					betSelectionID = bs.ID
					break
				}
			}

			betLegSelection := models.BetLegSelection{
				LegID:       betLeg.ID,
				SelectionID: betSelectionID,
			}

			if err := tx.Create(&betLegSelection).Error; err != nil {
				tx.Rollback()
				return nil, fmt.Errorf("failed to create bet leg selection: %w", err)
			}
		}
	}

	// 13. 提交事务
	if err := tx.Commit().Error; err != nil {
		return nil, fmt.Errorf("failed to commit transaction: %w", err)
	}

	// 14. 重新加载完整的投注信息
	var fullBet models.Bet
	if err := database.DB.Preload("Selections.Outcome.Market.Event").
		Preload("Legs.Selections.Selection.Outcome").
		First(&fullBet, bet.ID).Error; err != nil {
		return nil, fmt.Errorf("failed to load bet: %w", err)
	}

	return &fullBet, nil
}

// GetBet 获取投注详情
func (s *BetService) GetBet(betID int64) (*models.Bet, error) {
	var bet models.Bet
	if err := database.DB.Preload("Selections.Outcome.Market.Event").
		Preload("Legs.Selections.Selection.Outcome").
		First(&bet, betID).Error; err != nil {
		return nil, fmt.Errorf("bet not found: %w", err)
	}

	return &bet, nil
}

// GetUserBets 获取用户的投注列表
func (s *BetService) GetUserBets(userID int64, limit, offset int) ([]models.Bet, int64, error) {
	var bets []models.Bet
	var total int64

	// 查询总数
	if err := database.DB.Model(&models.Bet{}).Where("user_id = ?", userID).Count(&total).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to count bets: %w", err)
	}

	// 查询列表
	if err := database.DB.Preload("Selections.Outcome.Market.Event").
		Where("user_id = ?", userID).
		Order("placed_at DESC").
		Limit(limit).
		Offset(offset).
		Find(&bets).Error; err != nil {
		return nil, 0, fmt.Errorf("failed to fetch bets: %w", err)
	}

	return bets, total, nil
}

// SettleBet 结算投注
func (s *BetService) SettleBet(betID int64) error {
	// 1. 获取投注信息
	var bet models.Bet
	if err := database.DB.Preload("Selections.Outcome").
		Preload("Legs.Selections.Selection.Outcome").
		First(&bet, betID).Error; err != nil {
		return fmt.Errorf("bet not found: %w", err)
	}

	if bet.Status != "pending" {
		return fmt.Errorf("bet is already settled")
	}

	// 2. 检查所有选项的状态
	allWon := true
	anyLost := false
	anyVoid := false

	for _, sel := range bet.Selections {
		switch sel.Outcome.Status {
		case "won":
			// 继续
		case "lost":
			anyLost = true
			allWon = false
		case "void":
			anyVoid = true
			allWon = false
		default:
			return fmt.Errorf("outcome %d is not settled yet", sel.OutcomeID)
		}
	}

	// 3. 开始数据库事务
	tx := database.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// 4. 更新选项状态
	for _, sel := range bet.Selections {
		if err := tx.Model(&sel).Update("status", sel.Outcome.Status).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update selection status: %w", err)
		}
	}

	// 5. 结算每个组合腿
	totalActualReturn := 0.0

	for _, leg := range bet.Legs {
		legWon := true
		legVoid := false

		// 检查组合腿中的所有选项
		for _, legSel := range leg.Selections {
			switch legSel.Selection.Outcome.Status {
			case "won":
				// 继续
			case "lost":
				legWon = false
			case "void":
				legVoid = true
			}
		}

		var legStatus string
		var legActualReturn float64

		if legVoid {
			// 如果有 void，则退还本金
			legStatus = "void"
			legActualReturn = leg.Stake
		} else if legWon {
			// 所有选项都赢，计算回报
			legStatus = "won"
			legActualReturn = leg.PotentialReturn
		} else {
			// 有选项输了
			legStatus = "lost"
			legActualReturn = 0
		}

		totalActualReturn += legActualReturn

		if err := tx.Model(&leg).Updates(map[string]interface{}{
			"status":        legStatus,
			"actual_return": legActualReturn,
		}).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update leg status: %w", err)
		}
	}

	// 6. 更新投注状态
	var betStatus string
	if anyVoid && !anyLost {
		betStatus = "void"
	} else if allWon {
		betStatus = "won"
	} else if anyLost {
		if totalActualReturn > 0 {
			betStatus = "partially_won"
		} else {
			betStatus = "lost"
		}
	} else {
		betStatus = "lost"
	}

	settledAt := time.Now()
	if err := tx.Model(&bet).Updates(map[string]interface{}{
		"status":        betStatus,
		"actual_return": totalActualReturn,
		"settled_at":    settledAt,
	}).Error; err != nil {
		tx.Rollback()
		return fmt.Errorf("failed to update bet status: %w", err)
	}

	// 7. 如果有回报，增加用户余额
	if totalActualReturn > 0 {
		if err := tx.Model(&models.User{}).Where("id = ?", bet.UserID).
			Update("balance", gorm.Expr("balance + ?", totalActualReturn)).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("failed to update user balance: %w", err)
		}
	}

	// 8. 提交事务
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

