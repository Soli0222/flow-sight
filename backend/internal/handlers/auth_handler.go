package handlers

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Soli0222/flow-sight/backend/internal/config"
	"github.com/Soli0222/flow-sight/backend/internal/middleware"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authService AuthServiceInterface
	config      *config.Config
}

func NewAuthHandler(authService AuthServiceInterface, cfg *config.Config) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		config:      cfg,
	}
}

// GoogleLogin godoc
// @Summary Start Google OAuth login
// @Description Redirect to Google OAuth login page
// @Tags auth
// @Accept json
// @Produce json
// @Success 200 {object} map[string]string
// @Router /auth/google [get]
func (h *AuthHandler) GoogleLogin(c *gin.Context) {
	logger := middleware.GetLogger(c)
	ctx := context.Background()

	state := generateState()
	http.SetCookie(c.Writer, &http.Cookie{
		Name:     "oauthstate",
		Value:    state,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteLaxMode,
		// Optionally set MaxAge or Expires
	})
	url := h.authService.GetGoogleAuthURL(state)

	logger.InfoContext(ctx, "Google OAuth login initiated",
		"ip_address", c.ClientIP(),
		"user_agent", c.Request.UserAgent(),
	)

	c.JSON(http.StatusOK, gin.H{
		"url": url,
	})
}

// generateState creates a cryptographically secure random string for OAuth2 state
func generateState() string {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		// fallback: use uuid if crypto/rand fails
		return uuid.New().String()
	}
	return base64.URLEncoding.EncodeToString(b)
}

// GoogleCallback godoc
// @Summary Handle Google OAuth callback
// @Description Handle the callback from Google OAuth and create/login user
// @Tags auth
// @Accept json
// @Produce json
// @Param code query string true "Authorization code from Google"
// @Param state query string true "State parameter"
// @Success 302 {string} string "Redirect to frontend"
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/google/callback [get]
func (h *AuthHandler) GoogleCallback(c *gin.Context) {
	logger := middleware.GetLogger(c)
	ctx := context.Background()

	code := c.Query("code")
	if code == "" {
		logger.WarnContext(ctx, "OAuth callback missing code parameter",
			"ip_address", c.ClientIP(),
		)
		// エラー時はフロントエンドのログインページにリダイレクト
		c.Redirect(http.StatusFound, fmt.Sprintf("%s/login?error=no_code", h.config.Host))
		return
	}

	user, token, err := h.authService.HandleGoogleCallback(code)
	if err != nil {
		logger.ErrorContext(ctx, "OAuth callback failed",
			"error", err.Error(),
			"ip_address", c.ClientIP(),
		)
		// エラー時はフロントエンドのログインページにリダイレクト
		c.Redirect(http.StatusFound, fmt.Sprintf("%s/login?error=callback_failed", h.config.Host))
		return
	}

	logger.BusinessOperation(ctx, "user_login", user.ID.String(), map[string]interface{}{
		"email":      user.Email,
		"login_type": "google_oauth",
		"ip_address": c.ClientIP(),
	})

	// 成功時はフロントエンドのコールバックページにトークンとユーザー情報を渡してリダイレクト
	// 本来はより安全な方法（HTTPOnly cookieなど）を使うべきですが、
	// 簡単な実装としてURLパラメータを使用
	userJSON := fmt.Sprintf(`{"id":"%s","email":"%s","name":"%s","picture":"%s"}`,
		user.ID, user.Email, user.Name, user.Picture)
	encodedUser := url.QueryEscape(userJSON)

	redirectURL := fmt.Sprintf("%s/auth/callback?token=%s&user=%s",
		h.config.Host, token, encodedUser)

	c.Redirect(http.StatusFound, redirectURL)
}

// GetMe godoc
// @Summary Get current user information
// @Description Get the current authenticated user's information
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.User
// @Failure 401 {object} map[string]string
// @Router /auth/me [get]
func (h *AuthHandler) GetMe(c *gin.Context) {
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	userUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user_id format in context"})
		return
	}

	user, err := h.authService.GetUserByID(userUUID.String())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}
