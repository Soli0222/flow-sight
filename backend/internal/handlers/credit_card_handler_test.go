package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/Soli0222/flow-sight/backend/internal/models"
	"github.com/Soli0222/flow-sight/backend/test/helpers"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestCreditCardHandler_GetCreditCards(t *testing.T) {
	tests := []struct {
		name           string
		setupMock      func(*MockCreditCardServiceInterface)
		expectedStatus int
		expectedCount  int
	}{
		{
			name: "successful retrieval",
			setupMock: func(m *MockCreditCardServiceInterface) {
				testCreditCards := []models.CreditCard{
					{
						ID:          uuid.MustParse("11223344-5566-7788-99aa-bbccddeeff00"),
						Name:        "My Credit Card",
						ClosingDay:  func(i int) *int { return &i }(15),
						PaymentDay:  25,
						BankAccount: uuid.MustParse("aabbccdd-eeff-1122-3344-556677889900"),
						CreatedAt:   time.Now(),
						UpdatedAt:   time.Now(),
					},
				}
				m.On("GetCreditCards").Return(testCreditCards, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCount:  1,
		},
		{
			name: "service error",
			setupMock: func(m *MockCreditCardServiceInterface) {
				m.On("GetCreditCards").Return([]models.CreditCard{}, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMockCreditCardServiceInterface(t)
			handler := NewCreditCardHandler(mockService)

			// Setup mock
			tt.setupMock(mockService)

			// Create test context
			c, w := helpers.CreateTestContext(t, "GET", "/credit-cards", nil, true)

			// Call handler
			handler.GetCreditCards(c)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var creditCards []models.CreditCard
				err := json.Unmarshal(w.Body.Bytes(), &creditCards)
				assert.NoError(t, err)
				assert.Len(t, creditCards, tt.expectedCount)
			}
		})
	}
}

func TestCreditCardHandler_GetCreditCard(t *testing.T) {
	tests := []struct {
		name           string
		creditCardID   string
		setupMock      func(*MockCreditCardServiceInterface)
		expectedStatus int
	}{
		{
			name:         "successful retrieval",
			creditCardID: "11223344-5566-7788-99aa-bbccddeeff00",
			setupMock: func(m *MockCreditCardServiceInterface) {
				testCreditCard := &models.CreditCard{
					ID:          uuid.MustParse("11223344-5566-7788-99aa-bbccddeeff00"),
					Name:        "My Credit Card",
					ClosingDay:  func(i int) *int { return &i }(15),
					PaymentDay:  25,
					BankAccount: uuid.MustParse("aabbccdd-eeff-1122-3344-556677889900"),
					CreatedAt:   time.Now(),
					UpdatedAt:   time.Now(),
				}
				m.On("GetCreditCard", uuid.MustParse("11223344-5566-7788-99aa-bbccddeeff00")).Return(testCreditCard, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:         "invalid credit card ID",
			creditCardID: "invalid-uuid",
			setupMock: func(m *MockCreditCardServiceInterface) {
				// No mock setup needed for invalid UUID
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:         "credit card not found",
			creditCardID: "99999999-9999-9999-9999-999999999999",
			setupMock: func(m *MockCreditCardServiceInterface) {
				m.On("GetCreditCard", uuid.MustParse("99999999-9999-9999-9999-999999999999")).Return((*models.CreditCard)(nil), assert.AnError)
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMockCreditCardServiceInterface(t)
			handler := NewCreditCardHandler(mockService)

			// Setup mock
			tt.setupMock(mockService)

			// Create test context
			c, w := helpers.CreateTestContext(t, "GET", fmt.Sprintf("/credit-cards/%s", tt.creditCardID), nil, true)
			c.Params = gin.Params{{Key: "id", Value: tt.creditCardID}}

			// Call handler
			handler.GetCreditCard(c)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestCreditCardHandler_CreateCreditCard(t *testing.T) {
	tests := []struct {
		name           string
		requestBody    map[string]interface{}
		setupMock      func(*MockCreditCardServiceInterface)
		expectedStatus int
	}{
		{
			name: "successful creation",
			requestBody: map[string]interface{}{
				"name":         "My Credit Card",
				"closing_day":  15,
				"payment_day":  25,
				"bank_account": "aabbccdd-eeff-1122-3344-556677889900",
			},
			setupMock: func(m *MockCreditCardServiceInterface) {
				m.On("CreateCreditCard", mock.AnythingOfType("*models.CreditCard")).Return(nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name: "service validation error",
			requestBody: map[string]interface{}{
				"name": "", // 空の名前でサービスエラーが発生することを期待
			},
			setupMock: func(m *MockCreditCardServiceInterface) {
				m.On("CreateCreditCard", mock.AnythingOfType("*models.CreditCard")).Return(assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name: "service error",
			requestBody: map[string]interface{}{
				"name":         "My Credit Card",
				"closing_day":  15,
				"payment_day":  25,
				"bank_account": "aabbccdd-eeff-1122-3344-556677889900",
			},
			setupMock: func(m *MockCreditCardServiceInterface) {
				m.On("CreateCreditCard", mock.AnythingOfType("*models.CreditCard")).Return(assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMockCreditCardServiceInterface(t)
			handler := NewCreditCardHandler(mockService)

			// Setup mock
			tt.setupMock(mockService)

			// Create test context
			c, w := helpers.CreateTestContext(t, "POST", "/credit-cards", tt.requestBody, true)

			// Call handler
			handler.CreateCreditCard(c)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestCreditCardHandler_UpdateCreditCard(t *testing.T) {
	tests := []struct {
		name           string
		creditCardID   string
		requestBody    map[string]interface{}
		setupMock      func(*MockCreditCardServiceInterface)
		expectedStatus int
	}{
		{
			name:         "successful update",
			creditCardID: "11223344-5566-7788-99aa-bbccddeeff00",
			requestBody: map[string]interface{}{
				"name":         "Updated Credit Card",
				"closing_day":  20,
				"payment_day":  30,
				"bank_account": "aabbccdd-eeff-1122-3344-556677889900",
			},
			setupMock: func(m *MockCreditCardServiceInterface) {
				m.On("UpdateCreditCard", mock.AnythingOfType("*models.CreditCard")).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:         "invalid credit card ID",
			creditCardID: "invalid-uuid",
			requestBody: map[string]interface{}{
				"name":         "Updated Credit Card",
				"closing_day":  20,
				"payment_day":  30,
				"bank_account": "aabbccdd-eeff-1122-3344-556677889900",
			},
			setupMock: func(m *MockCreditCardServiceInterface) {
				// No mock setup needed for invalid UUID
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:         "service error",
			creditCardID: "11223344-5566-7788-99aa-bbccddeeff00",
			requestBody: map[string]interface{}{
				"name":         "Updated Credit Card",
				"closing_day":  20,
				"payment_day":  30,
				"bank_account": "aabbccdd-eeff-1122-3344-556677889900",
			},
			setupMock: func(m *MockCreditCardServiceInterface) {
				m.On("UpdateCreditCard", mock.AnythingOfType("*models.CreditCard")).Return(assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMockCreditCardServiceInterface(t)
			handler := NewCreditCardHandler(mockService)

			// Setup mock
			tt.setupMock(mockService)

			// Create test context
			c, w := helpers.CreateTestContext(t, "PUT", fmt.Sprintf("/credit-cards/%s", tt.creditCardID), tt.requestBody, true)
			c.Params = gin.Params{{Key: "id", Value: tt.creditCardID}}

			// Call handler
			handler.UpdateCreditCard(c)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestCreditCardHandler_DeleteCreditCard(t *testing.T) {
	tests := []struct {
		name           string
		creditCardID   string
		setupMock      func(*MockCreditCardServiceInterface)
		expectedStatus int
	}{
		{
			name:         "successful deletion",
			creditCardID: "11223344-5566-7788-99aa-bbccddeeff00",
			setupMock: func(m *MockCreditCardServiceInterface) {
				creditCardUUID, _ := uuid.Parse("11223344-5566-7788-99aa-bbccddeeff00")
				m.On("DeleteCreditCard", creditCardUUID).Return(nil)
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:         "invalid credit card ID",
			creditCardID: "invalid-uuid",
			setupMock: func(m *MockCreditCardServiceInterface) {
				// No mock setup needed for invalid UUID
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:         "service error",
			creditCardID: "99999999-9999-9999-9999-999999999999",
			setupMock: func(m *MockCreditCardServiceInterface) {
				creditCardUUID, _ := uuid.Parse("99999999-9999-9999-9999-999999999999")
				m.On("DeleteCreditCard", creditCardUUID).Return(assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMockCreditCardServiceInterface(t)
			handler := NewCreditCardHandler(mockService)

			// Setup mock
			tt.setupMock(mockService)

			// Create test context
			c, w := helpers.CreateTestContext(t, "DELETE", fmt.Sprintf("/credit-cards/%s", tt.creditCardID), nil, true)
			c.Params = gin.Params{{Key: "id", Value: tt.creditCardID}}

			// Call handler
			handler.DeleteCreditCard(c)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
