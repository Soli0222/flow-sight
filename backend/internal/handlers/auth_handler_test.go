package handlers

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Soli0222/flow-sight/backend/internal/config"
	"github.com/Soli0222/flow-sight/backend/internal/models"
	"github.com/Soli0222/flow-sight/backend/test/helpers"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthHandler_GoogleLogin(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*MockAuthServiceInterface)
		expectedStatus int
	}{
		{
			name: "successful login initiation",
			setupMock: func(m *MockAuthServiceInterface) {
				m.On("GetGoogleAuthURL", mock.AnythingOfType("string")).Return("https://accounts.google.com/oauth/authorize?...")
			},
			expectedStatus: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMockAuthServiceInterface(t)
			cfg := &config.Config{Host: "http://localhost:3000"}
			handler := NewAuthHandler(mockService, cfg)

			// Setup mock
			tt.setupMock(mockService)

			// Create test context
			c, w := helpers.CreateTestContext(t, "GET", "/auth/google", nil, false)

			// Call handler
			handler.GoogleLogin(c)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var response map[string]string
				err := json.Unmarshal(w.Body.Bytes(), &response)
				assert.NoError(t, err)
				assert.Contains(t, response, "url")
				assert.NotEmpty(t, response["url"])
			}
		})
	}
}

func TestAuthHandler_GoogleCallback(t *testing.T) {
	tests := []struct {
		name             string
		code             string
		setupMock        func(*MockAuthServiceInterface)
		expectedStatus   int
		expectedRedirect string
	}{
		{
			name: "successful callback",
			code: "test-auth-code",
			setupMock: func(m *MockAuthServiceInterface) {
				testUser := &models.User{
					ID:      uuid.MustParse("cbf3d545-d81d-450d-acb3-c5c49a597d6d"),
					Email:   "test@example.com",
					Name:    "Test User",
					Picture: "https://example.com/picture.jpg",
				}
				testToken := "test-jwt-token"
				m.On("HandleGoogleCallback", "test-auth-code").Return(testUser, testToken, nil)
			},
			expectedStatus:   http.StatusFound,
			expectedRedirect: "http://localhost:3000/auth/callback",
		},
		{
			name: "missing code parameter",
			code: "",
			setupMock: func(m *MockAuthServiceInterface) {
				// No mock setup needed for missing code
			},
			expectedStatus:   http.StatusFound,
			expectedRedirect: "http://localhost:3000/login?error=no_code",
		},
		{
			name: "callback service error",
			code: "invalid-code",
			setupMock: func(m *MockAuthServiceInterface) {
				m.On("HandleGoogleCallback", "invalid-code").Return((*models.User)(nil), "", assert.AnError)
			},
			expectedStatus:   http.StatusFound,
			expectedRedirect: "http://localhost:3000/login?error=callback_failed",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMockAuthServiceInterface(t)
			cfg := &config.Config{Host: "http://localhost:3000"}
			handler := NewAuthHandler(mockService, cfg)

			// Setup mock
			tt.setupMock(mockService)

			// Create test context
			url := "/auth/google/callback"
			if tt.code != "" {
				url += "?code=" + tt.code
			}
			c, w := helpers.CreateTestContext(t, "GET", url, nil, false)

			// Set query parameters
			if tt.code != "" {
				c.Request.URL.RawQuery = "code=" + tt.code
			}

			// Call handler
			handler.GoogleCallback(c)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusFound {
				location := w.Header().Get("Location")
				assert.Contains(t, location, tt.expectedRedirect)
			}
		})
	}
}

func TestAuthHandler_GetMe(t *testing.T) {
	tests := []struct {
		name           string
		authenticated  bool
		setupMock      func(*MockAuthServiceInterface)
		expectedStatus int
	}{
		{
			name:          "successful user retrieval",
			authenticated: true,
			setupMock: func(m *MockAuthServiceInterface) {
				testUser := &models.User{
					ID:      uuid.MustParse("cbf3d545-d81d-450d-acb3-c5c49a597d6d"),
					Email:   "test@example.com",
					Name:    "Test User",
					Picture: "https://example.com/picture.jpg",
				}
				m.On("GetUserByID", "cbf3d545-d81d-450d-acb3-c5c49a597d6d").Return(testUser, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:          "unauthenticated user",
			authenticated: false,
			setupMock: func(m *MockAuthServiceInterface) {
				// No mock setup needed for unauthenticated request
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:          "user not found",
			authenticated: true,
			setupMock: func(m *MockAuthServiceInterface) {
				m.On("GetUserByID", "cbf3d545-d81d-450d-acb3-c5c49a597d6d").Return((*models.User)(nil), assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMockAuthServiceInterface(t)
			cfg := &config.Config{Host: "http://localhost:3000"}
			handler := NewAuthHandler(mockService, cfg)

			// Setup mock
			tt.setupMock(mockService)

			// Create test context
			var c *gin.Context
			var w *httptest.ResponseRecorder
			if tt.authenticated {
				userID := uuid.MustParse("cbf3d545-d81d-450d-acb3-c5c49a597d6d")
				c, w = helpers.CreateTestContextWithUserID(t, "GET", "/auth/me", nil, userID)
			} else {
				c, w = helpers.CreateTestContext(t, "GET", "/auth/me", nil, false)
			}

			// Call handler
			handler.GetMe(c)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var user models.User
				err := json.Unmarshal(w.Body.Bytes(), &user)
				assert.NoError(t, err)
				assert.Equal(t, "test@example.com", user.Email)
				assert.Equal(t, "Test User", user.Name)
			}
		})
	}
}
