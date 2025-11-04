package handlers

import (
	"net/http"
	"strconv"

	"github.com/gdsZyy/betting-system/internal/database"
	"github.com/gdsZyy/betting-system/internal/models"
	"github.com/gin-gonic/gin"
)

// UserHandler 用户处理器
type UserHandler struct{}

// NewUserHandler 创建新的用户处理器
func NewUserHandler() *UserHandler {
	return &UserHandler{}
}

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string  `json:"username" binding:"required"`
	Email    string  `json:"email" binding:"required,email"`
	Balance  float64 `json:"balance"`
	Currency string  `json:"currency"`
}

// CreateUser 创建用户
// @Summary 创建用户
// @Description 创建新用户
// @Tags users
// @Accept json
// @Produce json
// @Param request body CreateUserRequest true "创建用户请求"
// @Success 200 {object} models.User
// @Failure 400 {object} ErrorResponse
// @Router /api/users [post]
func (h *UserHandler) CreateUser(c *gin.Context) {
	var req CreateUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	currency := req.Currency
	if currency == "" {
		currency = "CNY"
	}

	user := models.User{
		Username: req.Username,
		Email:    req.Email,
		Balance:  req.Balance,
		Currency: currency,
		Status:   "active",
	}

	if err := database.DB.Create(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to create user",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetUser 获取用户详情
// @Summary 获取用户详情
// @Description 根据用户ID获取用户详情
// @Tags users
// @Produce json
// @Param id path int true "用户ID"
// @Success 200 {object} models.User
// @Failure 404 {object} ErrorResponse
// @Router /api/users/{id} [get]
func (h *UserHandler) GetUser(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid user ID",
			Message: err.Error(),
		})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "User not found",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateBalanceRequest 更新余额请求
type UpdateBalanceRequest struct {
	Amount float64 `json:"amount" binding:"required"`
}

// DepositBalance 充值
// @Summary 充值
// @Description 增加用户余额
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Param request body UpdateBalanceRequest true "充值请求"
// @Success 200 {object} models.User
// @Failure 400 {object} ErrorResponse
// @Router /api/users/{id}/deposit [post]
func (h *UserHandler) DepositBalance(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid user ID",
			Message: err.Error(),
		})
		return
	}

	var req UpdateBalanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	if req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid amount",
			Message: "Amount must be greater than 0",
		})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "User not found",
			Message: err.Error(),
		})
		return
	}

	user.Balance += req.Amount
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to update balance",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

// WithdrawBalance 提现
// @Summary 提现
// @Description 减少用户余额
// @Tags users
// @Accept json
// @Produce json
// @Param id path int true "用户ID"
// @Param request body UpdateBalanceRequest true "提现请求"
// @Success 200 {object} models.User
// @Failure 400 {object} ErrorResponse
// @Router /api/users/{id}/withdraw [post]
func (h *UserHandler) WithdrawBalance(c *gin.Context) {
	userID, err := strconv.ParseInt(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid user ID",
			Message: err.Error(),
		})
		return
	}

	var req UpdateBalanceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid request",
			Message: err.Error(),
		})
		return
	}

	if req.Amount <= 0 {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Invalid amount",
			Message: "Amount must be greater than 0",
		})
		return
	}

	var user models.User
	if err := database.DB.First(&user, userID).Error; err != nil {
		c.JSON(http.StatusNotFound, ErrorResponse{
			Error:   "User not found",
			Message: err.Error(),
		})
		return
	}

	if user.Balance < req.Amount {
		c.JSON(http.StatusBadRequest, ErrorResponse{
			Error:   "Insufficient balance",
			Message: "User does not have enough balance",
		})
		return
	}

	user.Balance -= req.Amount
	if err := database.DB.Save(&user).Error; err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{
			Error:   "Failed to update balance",
			Message: err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, user)
}

