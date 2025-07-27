package handlers

import (
	"github.com/Soli0222/flow-sight/backend/internal/models"
	"github.com/Soli0222/flow-sight/backend/test/helpers"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestBankAccountHandler_GetBankAccounts_New(t *testing.T) {
	tests := []struct {
		name           string
		authenticated  bool
		setupMock      func(*MockBankAccountServiceInterface, uuid.UUID)
		expectedStatus int
		expectedCount  int
	}{
		{
			name:          "successful retrieval",
			authenticated: true,
			setupMock: func(m *MockBankAccountServiceInterface, userID uuid.UUID) {
				accounts := []models.BankAccount{
					*helpers.CreateTestBankAccount(userID),
					*helpers.CreateTestBankAccount(userID),
				}
				m.On("GetBankAccounts", userID).Return(accounts, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCount:  2,
		},
		{
			name:          "unauthenticated user",
			authenticated: false,
			setupMock: func(m *MockBankAccountServiceInterface, userID uuid.UUID) {
				// No mock setup needed for unauthenticated request
			},
			expectedStatus: http.StatusUnauthorized,
			expectedCount:  0,
		},
		{
			name:          "service error",
			authenticated: true,
			setupMock: func(m *MockBankAccountServiceInterface, userID uuid.UUID) {
				m.On("GetBankAccounts", userID).Return([]models.BankAccount{}, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMockBankAccountServiceInterface(t)
			handler := NewBankAccountHandler(mockService)

			// Create test context
			c, w := helpers.CreateTestContext(t, "GET", "/bank-accounts", nil, tt.authenticated)

			userID := uuid.New()
			if tt.authenticated {
				c.Set("user_id", userID)
			}

			tt.setupMock(mockService, userID)

			// Execute handler
			handler.GetBankAccounts(c)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var accounts []models.BankAccount
				helpers.ParseJSONResponse(t, w, &accounts)
				assert.Len(t, accounts, tt.expectedCount)
			}
		})
	}
}

func TestBankAccountHandler_GetBankAccount_New(t *testing.T) {
	tests := []struct {
		name           string
		accountIDStr   string
		authenticated  bool
		setupMock      func(*MockBankAccountServiceInterface, string)
		expectedStatus int
	}{
		{
			name:          "successful retrieval",
			accountIDStr:  uuid.New().String(),
			authenticated: true,
			setupMock: func(m *MockBankAccountServiceInterface, accountIDStr string) {
				accountID, _ := uuid.Parse(accountIDStr)
				testAccount := helpers.CreateTestBankAccount(uuid.New())
				m.On("GetBankAccount", accountID).Return(testAccount, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:          "invalid account ID",
			accountIDStr:  "invalid-uuid",
			authenticated: true,
			setupMock: func(m *MockBankAccountServiceInterface, accountIDStr string) {
				// No mock setup needed for invalid UUID
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:          "account not found",
			accountIDStr:  uuid.New().String(),
			authenticated: true,
			setupMock: func(m *MockBankAccountServiceInterface, accountIDStr string) {
				accountID, _ := uuid.Parse(accountIDStr)
				m.On("GetBankAccount", accountID).Return(nil, assert.AnError)
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMockBankAccountServiceInterface(t)
			handler := NewBankAccountHandler(mockService)

			// Create test context
			c, w := helpers.CreateTestContext(t, "GET", "/bank-accounts/"+tt.accountIDStr, nil, tt.authenticated)
			c.Params = []gin.Param{{Key: "id", Value: tt.accountIDStr}}

			if tt.authenticated {
				c.Set("user_id", uuid.New())
			}

			tt.setupMock(mockService, tt.accountIDStr)

			// Execute handler
			handler.GetBankAccount(c)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestBankAccountHandler_CreateBankAccount_New(t *testing.T) {
	tests := []struct {
		name           string
		authenticated  bool
		requestBody    interface{}
		setupMock      func(*MockBankAccountServiceInterface)
		expectedStatus int
	}{
		{
			name:          "successful creation",
			authenticated: true,
			requestBody: map[string]interface{}{
				"name":    "Test Bank Account",
				"balance": int64(100000), // 1000.00 in cents
			},
			setupMock: func(m *MockBankAccountServiceInterface) {
				m.On("CreateBankAccount", mock.AnythingOfType("*models.BankAccount")).Return(nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:          "unauthenticated user",
			authenticated: false,
			requestBody: map[string]interface{}{
				"name":    "Test Bank Account",
				"balance": int64(100000),
			},
			setupMock: func(m *MockBankAccountServiceInterface) {
				// No mock setup needed for unauthenticated request
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:          "invalid JSON body",
			authenticated: true,
			requestBody:   "invalid-json",
			setupMock: func(m *MockBankAccountServiceInterface) {
				// No mock setup needed for invalid JSON
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:          "service error",
			authenticated: true,
			requestBody: map[string]interface{}{
				"name":    "Test Bank Account",
				"balance": int64(100000),
			},
			setupMock: func(m *MockBankAccountServiceInterface) {
				m.On("CreateBankAccount", mock.AnythingOfType("*models.BankAccount")).Return(assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMockBankAccountServiceInterface(t)
			handler := NewBankAccountHandler(mockService)

			// Create test context
			c, w := helpers.CreateTestContext(t, "POST", "/bank-accounts", tt.requestBody, tt.authenticated)

			if tt.authenticated {
				c.Set("user_id", uuid.New())
			}

			tt.setupMock(mockService)

			// Execute handler
			handler.CreateBankAccount(c)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestBankAccountHandler_UpdateBankAccount_New(t *testing.T) {
	tests := []struct {
		name           string
		accountID      string
		authenticated  bool
		requestBody    interface{}
		setupMock      func(*MockBankAccountServiceInterface)
		expectedStatus int
	}{
		{
			name:          "successful update",
			accountID:     uuid.New().String(),
			authenticated: true,
			requestBody: map[string]interface{}{
				"name":    "Updated Bank Account",
				"balance": int64(200000), // 2000.00 in cents
			},
			setupMock: func(m *MockBankAccountServiceInterface) {
				m.On("UpdateBankAccount", mock.AnythingOfType("*models.BankAccount")).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:          "unauthenticated user",
			accountID:     uuid.New().String(),
			authenticated: false,
			requestBody: map[string]interface{}{
				"name":    "Updated Bank Account",
				"balance": int64(200000),
			},
			setupMock: func(m *MockBankAccountServiceInterface) {
				// No mock setup needed for unauthenticated request
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:          "invalid account ID",
			accountID:     "invalid-uuid",
			authenticated: true,
			requestBody: map[string]interface{}{
				"name":    "Updated Bank Account",
				"balance": int64(200000),
			},
			setupMock: func(m *MockBankAccountServiceInterface) {
				// No mock setup needed for invalid UUID
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:          "service error",
			accountID:     uuid.New().String(),
			authenticated: true,
			requestBody: map[string]interface{}{
				"name":    "Updated Bank Account",
				"balance": int64(200000),
			},
			setupMock: func(m *MockBankAccountServiceInterface) {
				m.On("UpdateBankAccount", mock.AnythingOfType("*models.BankAccount")).Return(assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMockBankAccountServiceInterface(t)
			handler := NewBankAccountHandler(mockService)

			// Create test context
			c, w := helpers.CreateTestContext(t, "PUT", "/bank-accounts/"+tt.accountID, tt.requestBody, tt.authenticated)
			c.Params = []gin.Param{{Key: "id", Value: tt.accountID}}

			if tt.authenticated {
				c.Set("user_id", uuid.New())
			}

			tt.setupMock(mockService)

			// Execute handler
			handler.UpdateBankAccount(c)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestBankAccountHandler_DeleteBankAccount_New(t *testing.T) {
	tests := []struct {
		name           string
		accountID      string
		authenticated  bool
		setupMock      func(*MockBankAccountServiceInterface)
		expectedStatus int
	}{
		{
			name:          "successful deletion",
			accountID:     "354a4ccc-1ac2-44ea-9d52-a9b76b9a7518",
			authenticated: true,
			setupMock: func(m *MockBankAccountServiceInterface) {
				accountUUID, _ := uuid.Parse("354a4ccc-1ac2-44ea-9d52-a9b76b9a7518")
				m.On("DeleteBankAccount", accountUUID).Return(nil)
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:          "unauthenticated user",
			accountID:     "unknown-uuid",
			authenticated: false,
			setupMock: func(m *MockBankAccountServiceInterface) {
				// No mock setup needed for unauthenticated request
			},
			expectedStatus: http.StatusUnauthorized,
		},
		{
			name:          "invalid account ID",
			accountID:     "invalid-uuid",
			authenticated: true,
			setupMock: func(m *MockBankAccountServiceInterface) {
				// No mock setup needed for invalid UUID
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:          "service error",
			accountID:     "e8149fec-e1be-4512-8acc-3437222b581a",
			authenticated: true,
			setupMock: func(m *MockBankAccountServiceInterface) {
				accountUUID, _ := uuid.Parse("e8149fec-e1be-4512-8acc-3437222b581a")
				m.On("DeleteBankAccount", accountUUID).Return(assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMockBankAccountServiceInterface(t)
			handler := NewBankAccountHandler(mockService)

			// Create test context
			c, w := helpers.CreateTestContext(t, "DELETE", "/bank-accounts/"+tt.accountID, nil, tt.authenticated)
			c.Params = []gin.Param{{Key: "id", Value: tt.accountID}}

			if tt.authenticated {
				c.Set("user_id", uuid.New())
			}

			tt.setupMock(mockService)

			// Execute handler
			handler.DeleteBankAccount(c)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
