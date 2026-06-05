package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"aiweekly/backend/internal/config"
	"aiweekly/backend/internal/model"
	"aiweekly/backend/internal/response"
	"aiweekly/backend/internal/service"
)

type AuthHandler struct {
	DB  *gorm.DB
	Cfg config.Config
}

type registerReq struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Phone    string `json:"phone"`
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req registerReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}

	req.Username = strings.TrimSpace(req.Username)
	req.Phone = strings.TrimSpace(req.Phone)
	if req.Username == "" || req.Phone == "" || len(req.Password) < 6 || len(req.Phone) != 11 {
		response.Fail(c, 400, "用户名/手机号不能为空，密码至少 6 位，手机号必须为 11 位")
		return
	}

	var exists int64
	h.DB.Model(&model.User{}).Where("username = ?", req.Username).Count(&exists)
	if exists > 0 {
		response.Fail(c, 409, "用户名已存在")
		return
	}

	h.DB.Model(&model.User{}).Where("phone = ?", req.Phone).Count(&exists)
	if exists > 0 {
		response.Fail(c, 409, "手机号已存在")
		return
	}

	hash, err := service.HashPassword(req.Password)
	if err != nil {
		response.Fail(c, 500, "密码处理失败")
		return
	}

	u := model.User{
		Username:     req.Username,
		PasswordHash: hash,
		Phone:        req.Phone,
		TokenVersion: 0,
	}
	if err := h.DB.Create(&u).Error; err != nil {
		response.Fail(c, 500, "创建用户失败")
		return
	}

	response.OK(c, gin.H{
		"id":       u.ID,
		"username": u.Username,
		"phone":    u.Phone,
	})
}

type loginReq struct {
	Phone    string `json:"phone"`
	Password string `json:"password"`
	Remember bool   `json:"remember"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Fail(c, 400, "参数错误")
		return
	}

	req.Phone = strings.TrimSpace(req.Phone)
	if req.Phone == "" || req.Password == "" {
		response.Fail(c, 400, "请输入手机号和密码")
		return
	}

	var u model.User
	if err := h.DB.Where("phone = ?", req.Phone).First(&u).Error; err != nil {
		response.Fail(c, 401, "手机号或密码错误")
		return
	}

	if !service.ComparePassword(u.PasswordHash, req.Password) {
		response.Fail(c, 401, "手机号或密码错误")
		return
	}

	access, err := service.SignToken(h.Cfg.JWTSecret, u.ID, u.TokenVersion, service.TokenTypeAccess, h.Cfg.AccessTokenTTL)
	if err != nil {
		response.Fail(c, 500, "生成 token 失败")
		return
	}

	refreshTTL := h.Cfg.RefreshTokenTTL
	if !req.Remember {
		// keep it shorter when not remembering; still allows a brief refresh window
		refreshTTL = 24 * time.Hour
	}
	refresh, err := service.SignToken(h.Cfg.JWTSecret, u.ID, u.TokenVersion, service.TokenTypeRefresh, refreshTTL)
	if err != nil {
		response.Fail(c, 500, "生成 token 失败")
		return
	}

	setAuthCookies(c, h.Cfg, access, refresh, h.Cfg.AccessTokenTTL, refreshTTL)

	response.OK(c, gin.H{
		"access_token":  access,
		"refresh_token": refresh,
		"user": gin.H{
			"id":       u.ID,
			"username": u.Username,
			"phone":    u.Phone,
		},
	})
}

func (h *AuthHandler) Refresh(c *gin.Context) {
	rt, ok := getTokenFromCookieOrBody(c, "refresh_token")
	if !ok {
		response.Fail(c, 401, "缺少 refresh token")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	claims, err := service.ParseToken(h.Cfg.JWTSecret, rt)
	if err != nil || claims.TokenType != service.TokenTypeRefresh {
		response.Fail(c, 401, "refresh token 无效")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	var u model.User
	if err := h.DB.Where("id = ?", claims.UserID).First(&u).Error; err != nil {
		response.Fail(c, 401, "用户不存在")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	if u.TokenVersion != claims.TokenVersion {
		response.Fail(c, 401, "登录状态已失效")
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	access, err := service.SignToken(h.Cfg.JWTSecret, u.ID, u.TokenVersion, service.TokenTypeAccess, h.Cfg.AccessTokenTTL)
	if err != nil {
		response.Fail(c, 500, "生成 token 失败")
		return
	}

	// Extend refresh by issuing a new refresh token with the default TTL.
	refresh, err := service.SignToken(h.Cfg.JWTSecret, u.ID, u.TokenVersion, service.TokenTypeRefresh, h.Cfg.RefreshTokenTTL)
	if err != nil {
		response.Fail(c, 500, "生成 token 失败")
		return
	}

	setAuthCookies(c, h.Cfg, access, refresh, h.Cfg.AccessTokenTTL, h.Cfg.RefreshTokenTTL)

	response.OK(c, gin.H{
		"access_token":  access,
		"refresh_token": refresh,
		"user": gin.H{
			"id":       u.ID,
			"username": u.Username,
			"phone":    u.Phone,
		},
	})
}

func (h *AuthHandler) Logout(c *gin.Context) {
	// invalidate refresh tokens by bumping token version (best-effort)
	at, _ := getTokenFromCookieOrBody(c, "access_token")
	if at != "" {
		if claims, err := service.ParseToken(h.Cfg.JWTSecret, at); err == nil {
			h.DB.Model(&model.User{}).Where("id = ?", claims.UserID).UpdateColumn("token_version", gorm.Expr("token_version + 1"))
		}
	}

	clearAuthCookies(c, h.Cfg)
	response.OK(c, gin.H{})
}

type tokenBody struct {
	Token string `json:"token"`
}

func getTokenFromCookieOrBody(c *gin.Context, cookieName string) (string, bool) {
	if v, err := c.Cookie(cookieName); err == nil && strings.TrimSpace(v) != "" {
		return v, true
	}
	var body tokenBody
	if err := c.ShouldBindJSON(&body); err == nil && strings.TrimSpace(body.Token) != "" {
		return strings.TrimSpace(body.Token), true
	}
	return "", false
}

func setAuthCookies(c *gin.Context, cfg config.Config, access, refresh string, accessTTL, refreshTTL time.Duration) {
	accessMaxAge := int(accessTTL.Seconds())
	refreshMaxAge := int(refreshTTL.Seconds())

	c.SetSameSite(http.SameSiteLaxMode)

	if cfg.CookieDomain == "" {
		c.SetCookie("access_token", access, accessMaxAge, "/", "", cfg.CookieSecure, true)
		c.SetCookie("refresh_token", refresh, refreshMaxAge, "/", "", cfg.CookieSecure, true)
	} else {
		c.SetCookie("access_token", access, accessMaxAge, "/", cfg.CookieDomain, cfg.CookieSecure, true)
		c.SetCookie("refresh_token", refresh, refreshMaxAge, "/", cfg.CookieDomain, cfg.CookieSecure, true)
	}
}

func clearAuthCookies(c *gin.Context, cfg config.Config) {
	c.SetSameSite(http.SameSiteLaxMode)
	if cfg.CookieDomain == "" {
		c.SetCookie("access_token", "", -1, "/", "", cfg.CookieSecure, true)
		c.SetCookie("refresh_token", "", -1, "/", "", cfg.CookieSecure, true)
	} else {
		c.SetCookie("access_token", "", -1, "/", cfg.CookieDomain, cfg.CookieSecure, true)
		c.SetCookie("refresh_token", "", -1, "/", cfg.CookieDomain, cfg.CookieSecure, true)
	}
}
