package router

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"aiweekly/backend/internal/config"
	"aiweekly/backend/internal/handler"
	"aiweekly/backend/internal/llm"
	"aiweekly/backend/internal/middleware"
	"aiweekly/backend/internal/response"
)

func Register(r *gin.Engine, db *gorm.DB, cfg config.Config) {
	r.GET("/healthz", func(c *gin.Context) {
		response.OK(c, gin.H{"status": "ok"})
	})

	api := r.Group("/api")
	{
		auth := &handler.AuthHandler{DB: db, Cfg: cfg}
		me := &handler.MeHandler{DB: db}

		// 创建 LLM 服务
		llmService := llm.NewService(cfg.DeepseekAPIKey)

		work := &handler.WorkHandler{DB: db, LLM: llmService}

		api.POST("/auth/register", auth.Register)
		api.POST("/auth/login", auth.Login)
		api.POST("/auth/refresh", auth.Refresh)
		api.POST("/auth/logout", auth.Logout)

		api.GET("/auth/me", middleware.RequireAccessToken(db, cfg.JWTSecret), me.Me)

		workGroup := api.Group("/work")
		workGroup.Use(middleware.RequireAccessToken(db, cfg.JWTSecret))
		{
			workGroup.GET("/records", work.GetRecords)
			workGroup.POST("/records", work.CreateRecord)
			workGroup.PUT("/records/:id", work.UpdateRecord)
			workGroup.DELETE("/records/:id", work.DeleteRecord)
			workGroup.GET("/records/monthly", work.GetMonthlyStats)

			// 日报相关
			workGroup.POST("/daily-report/generate", work.GenerateDailyReport)

			// 周报相关
			workGroup.GET("/weekly-report", work.GetWeeklyReport)
			workGroup.POST("/weekly-report/generate", work.GenerateWeeklyReport)
			workGroup.GET("/monthly-summary", work.GetMonthlySummary)

			// 用户设置相关
			workGroup.GET("/settings", work.GetUserSetting)
			workGroup.PUT("/settings", work.UpdateUserSetting)
		}
	}
}
