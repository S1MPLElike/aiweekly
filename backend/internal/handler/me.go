package handler

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"aiweekly/backend/internal/middleware"
	"aiweekly/backend/internal/response"
)

type MeHandler struct {
	DB *gorm.DB
}

func (h *MeHandler) Me(c *gin.Context) {
	u, ok := middleware.CurrentUser(c)
	if !ok {
		response.Fail(c, 401, "未登录")
		return
	}

	response.OK(c, gin.H{
		"id":       u.ID,
		"username": u.Username,
		"phone":    u.Phone,
	})
}

