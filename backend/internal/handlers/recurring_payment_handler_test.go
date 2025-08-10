package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Soli0222/flow-sight/backend/internal/models"
	"github.com/Soli0222/flow-sight/backend/test/helpers"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestRecurringPaymentHandler_GetRecurringPayments(t *testing.T) {
	tests := []struct {
		name           string
		authenticated  bool
		setupMock      func(*MockRecurringPaymentServiceInterface)
		expectedStatus int
		expectedCount  int
	}{
		{
			name:          "successful retrieval",
			authenticated: true,
			setupMock: func(m *MockRecurringPaymentServiceInterface) {
				testPayments := []models.RecurringPayment{
					{
						ID:                uuid.MustParse("11223344-5566-7788-99aa-bbccddeeff00"),
						Name:              "Monthly Subscription",
						Amount:            int64(99900), // $999.00 in cents
						PaymentDay:        15,
						StartYearMonth:    "2024-01",
						TotalPayments:     nil,
						RemainingPayments: nil,
						BankAccount:       uuid.MustParse("aabbccdd-eeff-1122-3344-556677889900"),
						IsActive:          true,
						Note:              "Netflix subscription",
						CreatedAt:         time.Now(),
						UpdatedAt:         time.Now(),
					},
				}
				m.On("GetRecurringPayments").Return(testPayments, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCount:  1,
		},
		{
			name:          "unauthenticated user - still processes request",
			authenticated: false,
			setupMock: func(m *MockRecurringPaymentServiceInterface) {
				// Handler doesn't check auth in single-user mode; still calls service
				m.On("GetRecurringPayments").Return([]models.RecurringPayment{}, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCount:  0,
		},
		{
			name:          "service error",
			authenticated: true,
			setupMock: func(m *MockRecurringPaymentServiceInterface) {
				m.On("GetRecurringPayments").Return([]models.RecurringPayment{}, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMockRecurringPaymentServiceInterface(t)
			handler := NewRecurringPaymentHandler(mockService)

			// Setup mock
			tt.setupMock(mockService)

			// Create test context
			var c *gin.Context
			var w *httptest.ResponseRecorder
			if tt.authenticated {
				userID := uuid.MustParse("cbf3d545-d81d-450d-acb3-c5c49a597d6d")
				c, w = helpers.CreateTestContextWithUserID(t, "GET", "/recurring-payments", nil, userID)
			} else {
				c, w = helpers.CreateTestContext(t, "GET", "/recurring-payments", nil, false)
			}

			// Call handler
			handler.GetRecurringPayments(c)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var payments []models.RecurringPayment
				err := json.Unmarshal(w.Body.Bytes(), &payments)
				assert.NoError(t, err)
				assert.Len(t, payments, tt.expectedCount)
			}
		})
	}
}

func TestRecurringPaymentHandler_GetRecurringPayment(t *testing.T) {
	tests := []struct {
		name           string
		paymentID      string
		setupMock      func(*MockRecurringPaymentServiceInterface)
		expectedStatus int
	}{
		{
			name:      "successful retrieval",
			paymentID: "11223344-5566-7788-99aa-bbccddeeff00",
			setupMock: func(m *MockRecurringPaymentServiceInterface) {
				testPayment := &models.RecurringPayment{
					ID:                uuid.MustParse("11223344-5566-7788-99aa-bbccddeeff00"),
					Name:              "Monthly Subscription",
					Amount:            int64(99900),
					PaymentDay:        15,
					StartYearMonth:    "2024-01",
					TotalPayments:     nil,
					RemainingPayments: nil,
					BankAccount:       uuid.MustParse("aabbccdd-eeff-1122-3344-556677889900"),
					IsActive:          true,
					Note:              "Netflix subscription",
					CreatedAt:         time.Now(),
					UpdatedAt:         time.Now(),
				}
				m.On("GetRecurringPayment", uuid.MustParse("11223344-5566-7788-99aa-bbccddeeff00")).Return(testPayment, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:      "invalid payment ID",
			paymentID: "invalid-uuid",
			setupMock: func(m *MockRecurringPaymentServiceInterface) {
				// No mock setup needed for invalid UUID
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:      "payment not found",
			paymentID: "99999999-9999-9999-9999-999999999999",
			setupMock: func(m *MockRecurringPaymentServiceInterface) {
				m.On("GetRecurringPayment", uuid.MustParse("99999999-9999-9999-9999-999999999999")).Return((*models.RecurringPayment)(nil), assert.AnError)
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMockRecurringPaymentServiceInterface(t)
			handler := NewRecurringPaymentHandler(mockService)

			// Setup mock
			tt.setupMock(mockService)

			// Create test context
			c, w := helpers.CreateTestContext(t, "GET", fmt.Sprintf("/recurring-payments/%s", tt.paymentID), nil, true)
			c.Params = gin.Params{{Key: "id", Value: tt.paymentID}}

			// Call handler
			handler.GetRecurringPayment(c)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestRecurringPaymentHandler_CreateRecurringPayment(t *testing.T) {
	tests := []struct {
		name           string
		authenticated  bool
		requestBody    map[string]interface{}
		setupMock      func(*MockRecurringPaymentServiceInterface)
		expectedStatus int
	}{
		{
			name:          "successful creation",
			authenticated: true,
			requestBody: map[string]interface{}{
				"name":             "Monthly Subscription",
				"amount":           int64(99900),
				"payment_day":      15,
				"start_year_month": "2024-01",
				"bank_account":     "aabbccdd-eeff-1122-3344-556677889900",
				"is_active":        true,
				"note":             "Netflix subscription",
			},
			setupMock: func(m *MockRecurringPaymentServiceInterface) {
				m.On("CreateRecurringPayment", mock.AnythingOfType("*models.RecurringPayment")).Return(nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:          "unauthenticated user - still processes request",
			authenticated: false,
			requestBody: map[string]interface{}{
				"name":             "Monthly Subscription",
				"amount":           int64(99900),
				"payment_day":      15,
				"start_year_month": "2024-01",
				"bank_account":     "aabbccdd-eeff-1122-3344-556677889900",
				"is_active":        true,
				"note":             "Netflix subscription",
			},
			setupMock: func(m *MockRecurringPaymentServiceInterface) {
				// Handler doesn't check auth; still calls service
				m.On("CreateRecurringPayment", mock.AnythingOfType("*models.RecurringPayment")).Return(nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:          "service validation error",
			authenticated: true,
			requestBody: map[string]interface{}{
				"name": "",
			},
			setupMock: func(m *MockRecurringPaymentServiceInterface) {
				m.On("CreateRecurringPayment", mock.AnythingOfType("*models.RecurringPayment")).Return(assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
		},
		{
			name:          "service error",
			authenticated: true,
			requestBody: map[string]interface{}{
				"name":             "Monthly Subscription",
				"amount":           int64(99900),
				"payment_day":      15,
				"start_year_month": "2024-01",
				"bank_account":     "aabbccdd-eeff-1122-3344-556677889900",
				"is_active":        true,
				"note":             "Netflix subscription",
			},
			setupMock: func(m *MockRecurringPaymentServiceInterface) {
				m.On("CreateRecurringPayment", mock.AnythingOfType("*models.RecurringPayment")).Return(assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMockRecurringPaymentServiceInterface(t)
			handler := NewRecurringPaymentHandler(mockService)

			// Setup mock
			tt.setupMock(mockService)

			// Create test context
			var c *gin.Context
			var w *httptest.ResponseRecorder
			if tt.authenticated {
				userID := uuid.MustParse("cbf3d545-d81d-450d-acb3-c5c49a597d6d")
				c, w = helpers.CreateTestContextWithUserID(t, "POST", "/recurring-payments", tt.requestBody, userID)
			} else {
				c, w = helpers.CreateTestContext(t, "POST", "/recurring-payments", tt.requestBody, false)
			}

			// Call handler
			handler.CreateRecurringPayment(c)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestRecurringPaymentHandler_UpdateRecurringPayment(t *testing.T) {
	tests := []struct {
		name           string
		paymentID      string
		requestBody    map[string]interface{}
		setupMock      func(*MockRecurringPaymentServiceInterface)
		expectedStatus int
	}{
		{
			name:      "successful update",
			paymentID: "11223344-5566-7788-99aa-bbccddeeff00",
			requestBody: map[string]interface{}{
				"name":             "Updated Subscription",
				"amount":           int64(119900),
				"payment_day":      20,
				"start_year_month": "2024-02",
				"bank_account":     "aabbccdd-eeff-1122-3344-556677889900",
				"is_active":        true,
				"note":             "Updated Netflix subscription",
			},
			setupMock: func(m *MockRecurringPaymentServiceInterface) {
				m.On("UpdateRecurringPayment", mock.AnythingOfType("*models.RecurringPayment")).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:      "invalid payment ID",
			paymentID: "invalid-uuid",
			requestBody: map[string]interface{}{
				"name":             "Updated Subscription",
				"amount":           int64(119900),
				"payment_day":      20,
				"start_year_month": "2024-02",
				"bank_account":     "aabbccdd-eeff-1122-3344-556677889900",
				"is_active":        true,
				"note":             "Updated Netflix subscription",
			},
			setupMock: func(m *MockRecurringPaymentServiceInterface) {
				// No mock setup needed for invalid UUID
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:      "service error",
			paymentID: "11223344-5566-7788-99aa-bbccddeeff00",
			requestBody: map[string]interface{}{
				"name":             "Updated Subscription",
				"amount":           int64(119900),
				"payment_day":      20,
				"start_year_month": "2024-02",
				"bank_account":     "aabbccdd-eeff-1122-3344-556677889900",
				"is_active":        true,
				"note":             "Updated Netflix subscription",
			},
			setupMock: func(m *MockRecurringPaymentServiceInterface) {
				m.On("UpdateRecurringPayment", mock.AnythingOfType("*models.RecurringPayment")).Return(assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMockRecurringPaymentServiceInterface(t)
			handler := NewRecurringPaymentHandler(mockService)

			// Setup mock
			tt.setupMock(mockService)

			// Create test context
			c, w := helpers.CreateTestContext(t, "PUT", fmt.Sprintf("/recurring-payments/%s", tt.paymentID), tt.requestBody, true)
			c.Params = gin.Params{{Key: "id", Value: tt.paymentID}}

			// Call handler
			handler.UpdateRecurringPayment(c)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestRecurringPaymentHandler_DeleteRecurringPayment(t *testing.T) {
	tests := []struct {
		name           string
		paymentID      string
		setupMock      func(*MockRecurringPaymentServiceInterface)
		expectedStatus int
	}{
		{
			name:      "successful deletion",
			paymentID: "11223344-5566-7788-99aa-bbccddeeff00",
			setupMock: func(m *MockRecurringPaymentServiceInterface) {
				paymentUUID, _ := uuid.Parse("11223344-5566-7788-99aa-bbccddeeff00")
				m.On("DeleteRecurringPayment", paymentUUID).Return(nil)
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:      "invalid payment ID",
			paymentID: "invalid-uuid",
			setupMock: func(m *MockRecurringPaymentServiceInterface) {
				// No mock setup needed for invalid UUID
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:      "service error",
			paymentID: "99999999-9999-9999-9999-999999999999",
			setupMock: func(m *MockRecurringPaymentServiceInterface) {
				paymentUUID, _ := uuid.Parse("99999999-9999-9999-9999-999999999999")
				m.On("DeleteRecurringPayment", paymentUUID).Return(assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMockRecurringPaymentServiceInterface(t)
			handler := NewRecurringPaymentHandler(mockService)

			// Setup mock
			tt.setupMock(mockService)

			// Create test context
			c, w := helpers.CreateTestContext(t, "DELETE", fmt.Sprintf("/recurring-payments/%s", tt.paymentID), nil, true)
			c.Params = gin.Params{{Key: "id", Value: tt.paymentID}}

			// Call handler
			handler.DeleteRecurringPayment(c)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
