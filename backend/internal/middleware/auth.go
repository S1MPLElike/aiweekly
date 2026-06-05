package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"aiweekly/backend/internal/model"
	"aiweekly/backend/internal/response"
	"aiweekly/backend/internal/service"
)

const (
	ctxUserKey   = "auth_user"
	ctxClaimsKey = "auth_claims"
)

func RequireAccessToken(db *gorm.DB, jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString := extractBearer(c.GetHeader("Authorization"))
		if tokenString == "" {
			// fallback to cookie for cookie-based auth
			if cookie, err := c.Cookie("access_token"); err == nil {
				tokenString = cookie
			}
		}

		if tokenString == "" {
			response.Fail(c, 401, "未登录")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		claims, err := service.ParseToken(jwtSecret, tokenString)
		if err != nil || claims.TokenType != service.TokenTypeAccess {
			response.Fail(c, 401, "登录已过期")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		var user model.User
		if err := db.Where("id = ?", claims.UserID).First(&user).Error; err != nil {
			response.Fail(c, 401, "用户不存在")
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		c.Set(ctxUserKey, &user)
		c.Set(ctxClaimsKey, claims)
		c.Next()
	}
}

func CurrentUser(c *gin.Context) (*model.User, bool) {
	v, ok := c.Get(ctxUserKey)
	if !ok {
		return nil, false
	}
	u, ok := v.(*model.User)
	return u, ok
}

func extractBearer(h string) string {
	h = strings.TrimSpace(h)
	if h == "" {
		return ""
	}
	const prefix = "Bearer "
	if !strings.HasPrefix(h, prefix) {
		return ""
	}
	return strings.TrimSpace(h[len(prefix):])
}

