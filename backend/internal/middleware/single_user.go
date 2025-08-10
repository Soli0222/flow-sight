package middleware

import (
	"github.com/gin-gonic/gin"
)

// SingleUserMiddleware keeps compatibility but does nothing in single-user mode
func SingleUserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// No user context required in single-user mode
		c.Next()
	}
}
