package middleware

import (
	"context"
	"time"

	"github.com/Soli0222/flow-sight/backend/internal/logger"

	"github.com/gin-gonic/gin"
)

const (
	RequestIDKey = "request_id"
	LoggerKey    = "logger"
)

// RequestLogger middleware logs HTTP requests with structured logging
func RequestLogger(baseLogger *logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		requestID := logger.GenerateRequestID()

		// Set request ID in context
		c.Set(RequestIDKey, requestID)

		// Create logger with request ID
		reqLogger := baseLogger.WithRequestID(requestID)
		c.Set(LoggerKey, reqLogger)

		// Add request ID to response header for debugging
		c.Header("X-Request-ID", requestID)

		// Log request start (debug level)
		ctx := context.Background()
		reqLogger.DebugContext(ctx, "Request started",
			"method", c.Request.Method,
			"path", c.Request.URL.Path,
			"remote_addr", c.ClientIP(),
			"user_agent", c.Request.UserAgent(),
		)

		// Process request
		c.Next()

		// Calculate duration
		duration := time.Since(start)

		// Get user ID if available
		userID := ""
		if uid, exists := c.Get("user_id"); exists {
			if uuidVal, ok := uid.(interface{ String() string }); ok {
				userID = uuidVal.String()
			}
		}

		// Log request completion
		reqLogger.Request(ctx, c.Request.Method, c.Request.URL.Path, c.Writer.Status(), duration, userID)

		// Log errors if any
		if len(c.Errors) > 0 {
			for _, ginErr := range c.Errors {
				reqLogger.ErrorContext(ctx, "Request error",
					"error", ginErr.Error(),
					"type", int(ginErr.Type),
				)
			}
		}
	}
}

// GetLogger extracts logger from gin context
func GetLogger(c *gin.Context) *logger.Logger {
	if loggerValue, exists := c.Get(LoggerKey); exists {
		if l, ok := loggerValue.(*logger.Logger); ok {
			return l
		}
	}
	// Fallback: this should not happen in normal operation
	panic("logger not found in context")
}

// GetRequestID extracts request ID from gin context
func GetRequestID(c *gin.Context) string {
	if requestID, exists := c.Get(RequestIDKey); exists {
		if id, ok := requestID.(string); ok {
			return id
		}
	}
	return ""
}
