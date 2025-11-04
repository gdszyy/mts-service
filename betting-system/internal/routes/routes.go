package routes

import (
	"github.com/gdsZyy/betting-system/internal/handlers"
	"github.com/gin-gonic/gin"
)

// SetupRoutes 设置路由
func SetupRoutes(r *gin.Engine) {
	// 创建处理器
	betHandler := handlers.NewBetHandler()
	eventHandler := handlers.NewEventHandler()
	userHandler := handlers.NewUserHandler()

	// API 路由组
	api := r.Group("/api")
	{
		// 用户路由
		users := api.Group("/users")
		{
			users.POST("", userHandler.CreateUser)
			users.GET("/:id", userHandler.GetUser)
			users.POST("/:id/deposit", userHandler.DepositBalance)
			users.POST("/:id/withdraw", userHandler.WithdrawBalance)
			users.GET("/:user_id/bets", betHandler.GetUserBets)
		}

		// 赛事路由
		events := api.Group("/events")
		{
			events.POST("", eventHandler.CreateEvent)
			events.GET("", eventHandler.ListEvents)
			events.GET("/:id", eventHandler.GetEvent)
		}

		// 盘口路由
		markets := api.Group("/markets")
		{
			markets.POST("", eventHandler.CreateMarket)
		}

		// 结果选项路由
		outcomes := api.Group("/outcomes")
		{
			outcomes.POST("", eventHandler.CreateOutcome)
			outcomes.PUT("/:id/odds", eventHandler.UpdateOutcomeOdds)
			outcomes.POST("/:id/settle", eventHandler.SettleOutcome)
		}

		// 投注路由
		bets := api.Group("/bets")
		{
			bets.POST("", betHandler.PlaceBet)
			bets.GET("/:id", betHandler.GetBet)
			bets.POST("/:id/settle", betHandler.SettleBet)
		}
	}

	// 健康检查
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "ok",
		})
	})
}

