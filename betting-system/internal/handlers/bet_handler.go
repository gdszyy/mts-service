package handlers

import (
	"net/http"
	"strconv"

	"github.com/gdsZyy/betting-system/internal/service"
	"github.com/gin-gonic/gin"
)

// BetHandler 投注处理器
type BetHandler struct {
	betService *service.BetService
}

// NewBetHandler 创建新的投注处理器
func NewBetHandler() *BetHandler {
	return &BetHandler{
		betService: service.NewBetService(),
	}
}

// PlaceBet 下注
// @Summary 下注
// @Description 创建新的投注
// @Tags bets
// @Accept json
// @Produce json
// @Param request body service.PlaceBetRequest true "下注请求"
// @Success 200 {object} models.Bet
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/bets [post]
func (h *BetHandler) PlaceBet(c *gin.Context) {
	var req service.PlaceBetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	bet, err := h.betService.PlaceBet(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Failed to place bet",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, bet)
}

// GetBet 获取投注详情
// @Summary 获取投注详情
// @Description 根据投注ID获取投注详情
// @Tags bets
// @Produce json
// @Param id path int true "投注ID"
// @Success 200 {object} models.Bet
// @Failure 404 {object} ErrorResponse
// @Router /api/bets/{id} [get]
func (h *BetHandler) GetBet(c *gin.Context) {
	betID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid bet ID",
			Message: err.Error(),
		})
		return
	}

	bet, err := h.betService.GetBet(betID)
	if err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Bet not found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, bet)
}

// GetUserBets 获取用户的投注列表
// @Summary 获取用户的投注列表
// @Description 根据用户ID获取投注列表
// @Tags bets
// @Produce json
// @Param user_id path int true "用户ID"
// @Param limit query int false "每页数量" default(20)
// @Param offset query int false "偏移量" default(0)
// @Success 200 {object} BetListResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/users/{user_id}/bets [get]
func (h *BetHandler) GetUserBets(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("user_id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid user ID",
			Message: err.Error(),
		})
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	bets, total, err := h.betService.GetUserBets(userID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to fetch bets",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, BetListResponse{
		Bets:   bets,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}

// SettleBet 结算投注
// @Summary 结算投注
// @Description 根据投注ID结算投注
// @Tags bets
// @Produce json
// @Param id path int true "投注ID"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Router /api/bets/{id}/settle [post]
func (h *BetHandler) SettleBet(c *gin.Context) {
	betID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid bet ID",
			Message: err.Error(),
		})
		return
	}

	if err := h.betService.SettleBet(betID); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Failed to settle bet",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, SuccessResponse{
		Message: "Bet settled successfully",
	})
}

// ErrorResponse 错误响应
type ErrorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// SuccessResponse 成功响应
type SuccessResponse struct {
	Message string `json:"message"`
}

// BetListResponse 投注列表响应
type BetListResponse struct {
	Bets   interface{} `json:"bets"`
	Total  int64       `json:"total"`
	Limit  int         `json:"limit"`
	Offset int         `json:"offset"`
}

