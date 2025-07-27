package services

import (
	"database/sql"
	"github.com/Soli0222/flow-sight/backend/internal/config"
	"github.com/Soli0222/flow-sight/backend/internal/services/mocks"
	"github.com/Soli0222/flow-sight/backend/test/helpers"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestAuthService_GenerateJWT(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{}
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret: "test-secret",
		},
	}
	service := NewAuthService(mockRepo, cfg)

	user := helpers.CreateTestUser()

	token, err := service.GenerateJWT(user)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestAuthService_ValidateJWT(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{}
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret: "test-secret",
		},
	}
	service := NewAuthService(mockRepo, cfg)

	user := helpers.CreateTestUser()

	// Generate a valid token
	token, err := service.GenerateJWT(user)
	assert.NoError(t, err)

	// Validate the token
	claims, err := service.ValidateJWT(token)
	assert.NoError(t, err)
	assert.NotNil(t, claims)
	assert.Equal(t, user.ID.String(), claims.UserID)
	assert.Equal(t, user.Email, claims.Email)

	// Test invalid token
	_, err = service.ValidateJWT("invalid-token")
	assert.Error(t, err)
}

func TestAuthService_GetUserByID(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{}
	cfg := &config.Config{
		JWT: config.JWTConfig{
			Secret: "test-secret",
		},
	}
	service := NewAuthService(mockRepo, cfg)

	tests := []struct {
		name          string
		userID        string
		setupMock     func(*mocks.MockUserRepository)
		expectedError bool
		expectedNil   bool
	}{
		{
			name:   "successful retrieval",
			userID: uuid.New().String(),
			setupMock: func(m *mocks.MockUserRepository) {
				user := helpers.CreateTestUser()
				m.On("GetByID", mock.AnythingOfType("uuid.UUID")).Return(user, nil)
			},
			expectedError: false,
			expectedNil:   false,
		},
		{
			name:   "user not found",
			userID: uuid.New().String(),
			setupMock: func(m *mocks.MockUserRepository) {
				m.On("GetByID", mock.AnythingOfType("uuid.UUID")).Return(nil, sql.ErrNoRows)
			},
			expectedError: true,
			expectedNil:   true,
		},
		{
			name:   "invalid UUID",
			userID: "invalid-uuid",
			setupMock: func(m *mocks.MockUserRepository) {
				// No mock setup needed for invalid UUID
			},
			expectedError: true,
			expectedNil:   true,
		},
		{
			name:   "repository error",
			userID: uuid.New().String(),
			setupMock: func(m *mocks.MockUserRepository) {
				m.On("GetByID", mock.AnythingOfType("uuid.UUID")).Return(nil, assert.AnError)
			},
			expectedError: true,
			expectedNil:   true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockRepo.ExpectedCalls = nil
			tt.setupMock(mockRepo)

			result, err := service.GetUserByID(tt.userID)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			if tt.expectedNil {
				assert.Nil(t, result)
			} else {
				assert.NotNil(t, result)
			}

			mockRepo.AssertExpectations(t)
		})
	}
}

func TestAuthService_GetGoogleAuthURL(t *testing.T) {
	mockRepo := &mocks.MockUserRepository{}
	cfg := &config.Config{
		OAuth: config.OAuthConfig{
			GoogleClientID:     "test-client-id",
			GoogleClientSecret: "test-client-secret",
			RedirectURL:        "http://localhost:3000/auth/callback",
		},
		JWT: config.JWTConfig{
			Secret: "test-secret",
		},
	}
	service := NewAuthService(mockRepo, cfg)

	state := "test-state"
	url := service.GetGoogleAuthURL(state)

	assert.NotEmpty(t, url)
	assert.Contains(t, url, "accounts.google.com")
	assert.Contains(t, url, "test-client-id")
	assert.Contains(t, url, state)
}
