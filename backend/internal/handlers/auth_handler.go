package handlers

import (
	"flow-sight-backend/internal/services"
	"fmt"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type AuthHandler struct {
	authService *services.AuthService
}

func NewAuthHandler(authService *services.AuthService) *AuthHandler {
	return &AuthHandler{
		authService: authService,
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
	state := "random-state-string" // In production, use a secure random state
	url := h.authService.GetGoogleAuthURL(state)

	c.JSON(http.StatusOK, gin.H{
		"url": url,
	})
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
	code := c.Query("code")
	if code == "" {
		// エラー時はフロントエンドのログインページにリダイレクト
		c.Redirect(http.StatusFound, "http://localhost:4000/login?error=no_code")
		return
	}

	user, token, err := h.authService.HandleGoogleCallback(code)
	if err != nil {
		// エラー時はフロントエンドのログインページにリダイレクト
		c.Redirect(http.StatusFound, "http://localhost:4000/login?error=callback_failed")
		return
	}

	// 成功時はフロントエンドのコールバックページにトークンとユーザー情報を渡してリダイレクト
	// 本来はより安全な方法（HTTPOnly cookieなど）を使うべきですが、
	// 簡単な実装としてURLパラメータを使用
	userJSON := fmt.Sprintf(`{"id":"%s","email":"%s","name":"%s","picture":"%s"}`,
		user.ID, user.Email, user.Name, user.Picture)
	encodedUser := url.QueryEscape(userJSON)

	redirectURL := fmt.Sprintf("http://localhost:4000/auth/callback?token=%s&user=%s",
		token, encodedUser)

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
