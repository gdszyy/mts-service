package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gdsZyy/betting-system/internal/database"
	"github.com/gdsZyy/betting-system/internal/models"
	"github.com/gin-gonic/gin"
)

// EventHandler 赛事处理器
type EventHandler struct{}

// NewEventHandler 创建新的赛事处理器
func NewEventHandler() *EventHandler {
	return &EventHandler{}
}

// CreateEventRequest 创建赛事请求
type CreateEventRequest struct {
	ExternalID string    `json:"external_id"`
	SportID    string    `json:"sport_id" binding:"required"`
	HomeTeam   string    `json:"home_team" binding:"required"`
	AwayTeam   string    `json:"away_team" binding:"required"`
	StartTime  time.Time `json:"start_time" binding:"required"`
}

// CreateEvent 创建赛事
// @Summary 创建赛事
// @Description 创建新的赛事
// @Tags events
// @Accept json
// @Produce json
// @Param request body CreateEventRequest true "创建赛事请求"
// @Success 200 {object} models.Event
// @Failure 400 {object} ErrorResponse
// @Router /api/events [post]
func (h *EventHandler) CreateEvent(c *gin.Context) {
	var req CreateEventRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	event := models.Event{
		ExternalID: req.ExternalID,
		SportID:    req.SportID,
		HomeTeam:   req.HomeTeam,
		AwayTeam:   req.AwayTeam,
		StartTime:  req.StartTime,
		Status:     "scheduled",
	}

	if err := database.DB.Create(&event).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to create event",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, event)
}

// GetEvent 获取赛事详情
// @Summary 获取赛事详情
// @Description 根据赛事ID获取赛事详情
// @Tags events
// @Produce json
// @Param id path int true "赛事ID"
// @Success 200 {object} models.Event
// @Failure 404 {object} ErrorResponse
// @Router /api/events/{id} [get]
func (h *EventHandler) GetEvent(c *gin.Context) {
	eventID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid event ID",
			Message: err.Error(),
		})
		return
	}

	var event models.Event
	if err := database.DB.First(&event, eventID).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Event not found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, event)
}

// ListEvents 获取赛事列表
// @Summary 获取赛事列表
// @Description 获取赛事列表，支持指定返回的 market 类型
// @Tags events
// @Produce json
// @Param status query string false "状态筛选"
// @Param limit query int false "每页数量" default(20)
// @Param offset query int false "偏移量" default(0)
// @Param market_types query string false "Market类型列表，逗号分隔，如: 1x2,handicap,totals"
// @Success 200 {object} EventListResponse
// @Router /api/events [get]
func (h *EventHandler) ListEvents(c *gin.Context) {
	status := c.Query("status")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))
	marketTypesParam := c.Query("market_types")

	query := database.DB.Model(&models.Event{})
	if status != "" {
		query = query.Where("status = ?", status)
	}

	var total int64
	query.Count(&total)

	var events []models.Event
	if err := query.Order("start_time ASC").Limit(limit).Offset(offset).Find(&events).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to fetch events",
			Message: err.Error(),
		})
		return
	}

	// 如果指定了 market_types，则加载过滤后的 markets 和 outcomes
	if marketTypesParam != "" {
		marketTypes := strings.Split(marketTypesParam, ",")
		// 清理空格
		for i, mt := range marketTypes {
			marketTypes[i] = strings.TrimSpace(mt)
		}

		// 为每个 event 加载指定类型的 markets
		for i := range events {
			var markets []models.Market
			if err := database.DB.Where("event_id = ? AND market_type IN ?", events[i].ID, marketTypes).
				Preload("Event").
				Find(&markets).Error; err != nil {
				// 如果加载失败，继续处理其他 events
				continue
			}

			// 为每个 market 加载 outcomes
			for j := range markets {
				var outcomes []models.Outcome
				if err := database.DB.Where("market_id = ?", markets[j].ID).Find(&outcomes).Error; err == nil {
					markets[j].Outcomes = outcomes
				}
			}

			events[i].Markets = markets
		}
	}

	c.JSON(http.StatusOK, EventListResponse{
		Events: events,
		Total:  total,
		Limit:  limit,
		Offset: offset,
	})
}

// CreateMarketRequest 创建盘口请求
type CreateMarketRequest struct {
	EventID    int64  `json:"event_id" binding:"required"`
	MarketType string `json:"market_type" binding:"required"`
	Specifier  string `json:"specifier"`
}

// CreateMarket 创建盘口
// @Summary 创建盘口
// @Description 创建新的盘口
// @Tags markets
// @Accept json
// @Produce json
// @Param request body CreateMarketRequest true "创建盘口请求"
// @Success 200 {object} models.Market
// @Failure 400 {object} ErrorResponse
// @Router /api/markets [post]
func (h *EventHandler) CreateMarket(c *gin.Context) {
	var req CreateMarketRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	// 验证赛事存在
	var event models.Event
	if err := database.DB.First(&event, req.EventID).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Event not found",
			Message: err.Error(),
		})
		return
	}

	market := models.Market{
		EventID:    req.EventID,
		MarketType: req.MarketType,
		Specifier:  req.Specifier,
		Status:     "active",
	}

	if err := database.DB.Create(&market).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to create market",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, market)
}

// CreateOutcomeRequest 创建结果选项请求
type CreateOutcomeRequest struct {
	MarketID  int64   `json:"market_id" binding:"required"`
	OutcomeID string  `json:"outcome_id" binding:"required"`
	Odds      float64 `json:"odds" binding:"required,gt=1"`
}

// CreateOutcome 创建结果选项
// @Summary 创建结果选项
// @Description 创建新的结果选项
// @Tags outcomes
// @Accept json
// @Produce json
// @Param request body CreateOutcomeRequest true "创建结果选项请求"
// @Success 200 {object} models.Outcome
// @Failure 400 {object} ErrorResponse
// @Router /api/outcomes [post]
func (h *EventHandler) CreateOutcome(c *gin.Context) {
	var req CreateOutcomeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	// 验证盘口存在
	var market models.Market
	if err := database.DB.First(&market, req.MarketID).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Market not found",
			Message: err.Error(),
		})
		return
	}

	outcome := models.Outcome{
		MarketID:  req.MarketID,
		OutcomeID: req.OutcomeID,
		Odds:      req.Odds,
		Status:    "active",
	}

	if err := database.DB.Create(&outcome).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to create outcome",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, outcome)
}

// UpdateOutcomeOddsRequest 更新赔率请求
type UpdateOutcomeOddsRequest struct {
	Odds float64 `json:"odds" binding:"required,gt=1"`
}

// UpdateOutcomeOdds 更新赔率
// @Summary 更新赔率
// @Description 更新结果选项的赔率
// @Tags outcomes
// @Accept json
// @Produce json
// @Param id path int true "结果选项ID"
// @Param request body UpdateOutcomeOddsRequest true "更新赔率请求"
// @Success 200 {object} models.Outcome
// @Failure 400 {object} ErrorResponse
// @Router /api/outcomes/{id}/odds [put]
func (h *EventHandler) UpdateOutcomeOdds(c *gin.Context) {
	outcomeID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid outcome ID",
			Message: err.Error(),
		})
		return
	}

	var req UpdateOutcomeOddsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	var outcome models.Outcome
	if err := database.DB.First(&outcome, outcomeID).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Outcome not found",
			Message: err.Error(),
		})
		return
	}

	if err := database.DB.Model(&outcome).Update("odds", req.Odds).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to update odds",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, outcome)
}

// SettleOutcomeRequest 结算结果选项请求
type SettleOutcomeRequest struct {
	Status string `json:"status" binding:"required,oneof=won lost void"`
}

// SettleOutcome 结算结果选项
// @Summary 结算结果选项
// @Description 结算结果选项（设置为won/lost/void）
// @Tags outcomes
// @Accept json
// @Produce json
// @Param id path int true "结果选项ID"
// @Param request body SettleOutcomeRequest true "结算请求"
// @Success 200 {object} models.Outcome
// @Failure 400 {object} ErrorResponse
// @Router /api/outcomes/{id}/settle [post]
func (h *EventHandler) SettleOutcome(c *gin.Context) {
	outcomeID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid outcome ID",
			Message: err.Error(),
		})
		return
	}

	var req SettleOutcomeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	var outcome models.Outcome
	if err := database.DB.First(&outcome, outcomeID).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "Outcome not found",
			Message: err.Error(),
		})
		return
	}

	if err := database.DB.Model(&outcome).Update("status", req.Status).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to settle outcome",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, outcome)
}

// EventListResponse 赛事列表响应
type EventListResponse struct {
	Events []models.Event `json:"events"`
	Total  int64          `json:"total"`
	Limit  int            `json:"limit"`
	Offset int            `json:"offset"`
}

