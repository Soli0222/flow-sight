package helpers

import (
	"bytes"
	"encoding/json"
	"log/slog"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/Soli0222/flow-sight/backend/internal/logger"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// SetupGin creates a gin engine in test mode
func SetupGin() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

// CreateTestLogger creates a test logger
func CreateTestLogger() *logger.Logger {
	slogger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	return &logger.Logger{Logger: slogger}
}

// CreateTestContext creates a test context with optional user authentication
func CreateTestContext(t *testing.T, method, path string, body interface{}, authenticated bool) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)

	var bodyBytes []byte
	var err error

	if body != nil {
		bodyBytes, err = json.Marshal(body)
		require.NoError(t, err)
	}

	req := httptest.NewRequest(method, path, bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Add test logger to context
	testLogger := CreateTestLogger()
	c.Set("logger", testLogger)

	// Add authenticated user context if needed
	if authenticated {
		userID := uuid.New()
		c.Set("user_id", userID)
	}

	return c, w
}

// CreateTestContextWithUserID creates a test context with a specific user ID
func CreateTestContextWithUserID(t *testing.T, method, path string, body interface{}, userID uuid.UUID) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)

	var bodyBytes []byte
	var err error

	if body != nil {
		bodyBytes, err = json.Marshal(body)
		require.NoError(t, err)
	}

	req := httptest.NewRequest(method, path, bytes.NewBuffer(bodyBytes))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = req

	// Add test logger to context
	testLogger := CreateTestLogger()
	c.Set("logger", testLogger)
	c.Set("user_id", userID)

	return c, w
}

// ParseJSONResponse parses the response body as JSON
func ParseJSONResponse(t *testing.T, w *httptest.ResponseRecorder, target interface{}) {
	err := json.Unmarshal(w.Body.Bytes(), target)
	require.NoError(t, err)
}
