package logger

import (
	"context"
	"io"
	"log/slog"
	"os"
	"time"

	"github.com/Soli0222/flow-sight/backend/internal/config"
)

// Logger wraps slog.Logger with additional functionality
type Logger struct {
	*slog.Logger
}

// LogConfig holds logging configuration
type LogConfig struct {
	Level       slog.Level
	Environment string
	Service     string
	Version     string
	Output      io.Writer
}

// New creates a new structured logger based on environment
func New(cfg *config.Config, version string) *Logger {
	logConfig := &LogConfig{
		Environment: cfg.Env,
		Service:     "github.com/Soli0222/flow-sight/backend",
		Version:     version,
		Output:      os.Stdout,
	}

	// Set log level based on environment
	if cfg.Env == "production" {
		logConfig.Level = slog.LevelInfo
	} else {
		logConfig.Level = slog.LevelDebug
	}

	var handler slog.Handler

	// Configure handler based on environment
	if cfg.Env == "production" {
		// JSON format for production
		handler = slog.NewJSONHandler(logConfig.Output, &slog.HandlerOptions{
			Level: logConfig.Level,
			ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
				// Add custom attributes
				if a.Key == slog.TimeKey {
					a.Value = slog.StringValue(time.Now().UTC().Format(time.RFC3339))
				}
				return a
			},
		})
	} else {
		// Text format for development
		handler = slog.NewTextHandler(logConfig.Output, &slog.HandlerOptions{
			Level: logConfig.Level,
		})
	}

	// Add default attributes
	logger := slog.New(handler).With(
		"service", logConfig.Service,
		"version", logConfig.Version,
		"environment", logConfig.Environment,
	)

	return &Logger{Logger: logger}
}

// WithRequestID adds request ID to logger context
func (l *Logger) WithRequestID(requestID string) *Logger {
	return &Logger{
		Logger: l.Logger.With("request_id", requestID),
	}
}

// WithUserID adds user ID to logger context
func (l *Logger) WithUserID(userID string) *Logger {
	return &Logger{
		Logger: l.Logger.With("user_id", userID),
	}
}

// WithError adds error information to logger context
func (l *Logger) WithError(err error) *Logger {
	if err == nil {
		return l
	}
	return &Logger{
		Logger: l.Logger.With("error", err.Error()),
	}
}

// Request logs HTTP request information
func (l *Logger) Request(ctx context.Context, method, path string, statusCode int, duration time.Duration, userID string) {
	attrs := []any{
		"method", method,
		"path", path,
		"status_code", statusCode,
		"duration", duration,
	}

	if userID != "" {
		attrs = append(attrs, "user_id", userID)
	}

	l.InfoContext(ctx, "HTTP request", attrs...)
}

// Error logs error with structured information
func (l *Logger) Error(ctx context.Context, msg string, err error, attrs ...any) {
	if err != nil {
		attrs = append(attrs, "error", err.Error())
	}
	l.ErrorContext(ctx, msg, attrs...)
}

// BusinessOperation logs important business operations
func (l *Logger) BusinessOperation(ctx context.Context, operation string, userID string, details map[string]interface{}) {
	attrs := []any{
		"operation", operation,
		"user_id", userID,
	}

	for key, value := range details {
		attrs = append(attrs, key, value)
	}

	l.InfoContext(ctx, "Business operation", attrs...)
}

// DatabaseOperation logs database operations
func (l *Logger) DatabaseOperation(ctx context.Context, operation string, table string, duration time.Duration, err error) {
	attrs := []any{
		"operation", operation,
		"table", table,
		"duration", duration,
	}

	if err != nil {
		attrs = append(attrs, "error", err.Error())
		l.ErrorContext(ctx, "Database operation failed", attrs...)
	} else {
		l.DebugContext(ctx, "Database operation", attrs...)
	}
}

// Security logs security-related events
func (l *Logger) Security(ctx context.Context, event string, userID string, ipAddress string, success bool) {
	attrs := []any{
		"event", event,
		"ip_address", ipAddress,
		"success", success,
	}

	if userID != "" {
		attrs = append(attrs, "user_id", userID)
	}

	if success {
		l.InfoContext(ctx, "Security event", attrs...)
	} else {
		l.WarnContext(ctx, "Security event", attrs...)
	}
}
