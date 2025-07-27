package middleware

import (
	"context"
	"flow-sight-backend/internal/services"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func AuthMiddleware(authService *services.AuthService) gin.HandlerFunc {
	return func(c *gin.Context) {
		logger := GetLogger(c)
		ctx := context.Background()

		authHeader := c.GetHeader("Authorization")

		if authHeader == "" {
			logger.Security(ctx, "auth_missing_header", "", c.ClientIP(), false)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "authorization header is required"})
			c.Abort()
			return
		}

		// Extract token from "Bearer <token>"
		tokenParts := strings.Split(authHeader, " ")

		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			logger.Security(ctx, "auth_invalid_header", "", c.ClientIP(), false)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid authorization header format"})
			c.Abort()
			return
		}

		token := tokenParts[1]

		claims, err := authService.ValidateJWT(token)
		if err != nil {
			logger.Security(ctx, "auth_invalid_token", "", c.ClientIP(), false)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		// Set user information in context
		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			logger.Security(ctx, "auth_invalid_user_id", claims.UserID, c.ClientIP(), false)
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user_id in token"})
			c.Abort()
			return
		}

		logger.Security(ctx, "auth_success", claims.UserID, c.ClientIP(), true)

		c.Set("user_id", userID)
		c.Set("user_email", claims.Email)
		c.Next()
	}
}
