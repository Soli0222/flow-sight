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

func TestIncomeHandler_GetIncomeSources(t *testing.T) {
	tests := []struct {
		name           string
		authenticated  bool
		setupMock      func(*MockIncomeServiceInterface)
		expectedStatus int
		expectedCount  int
	}{
		{
			name:          "successful retrieval",
			authenticated: true,
			setupMock: func(m *MockIncomeServiceInterface) {
				testSources := []models.IncomeSource{
					{
						ID:                 uuid.MustParse("11223344-5566-7788-99aa-bbccddeeff00"),
						UserID:             uuid.MustParse("cbf3d545-d81d-450d-acb3-c5c49a597d6d"),
						Name:               "Monthly Salary",
						IncomeType:         "monthly_fixed",
						BaseAmount:         int64(500000), // $5000.00 in cents
						BankAccount:        uuid.MustParse("aabbccdd-eeff-1122-3344-556677889900"),
						PaymentDay:         func(i int) *int { return &i }(25),
						ScheduledDate:      nil,
						ScheduledYearMonth: nil,
						IsActive:           true,
						CreatedAt:          time.Now(),
						UpdatedAt:          time.Now(),
					},
				}
				m.On("GetIncomeSources", uuid.MustParse("cbf3d545-d81d-450d-acb3-c5c49a597d6d")).Return(testSources, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCount:  1,
		},
		{
			name:          "unauthenticated user",
			authenticated: false,
			setupMock: func(m *MockIncomeServiceInterface) {
				// No mock setup needed for unauthenticated request
			},
			expectedStatus: http.StatusForbidden,
			expectedCount:  0,
		},
		{
			name:          "service error",
			authenticated: true,
			setupMock: func(m *MockIncomeServiceInterface) {
				m.On("GetIncomeSources", uuid.MustParse("cbf3d545-d81d-450d-acb3-c5c49a597d6d")).Return([]models.IncomeSource{}, assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
			expectedCount:  0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMockIncomeServiceInterface(t)
			handler := NewIncomeHandler(mockService)

			// Setup mock
			tt.setupMock(mockService)

			// Create test context
			var c *gin.Context
			var w *httptest.ResponseRecorder
			if tt.authenticated {
				userID := uuid.MustParse("cbf3d545-d81d-450d-acb3-c5c49a597d6d")
				c, w = helpers.CreateTestContextWithUserID(t, "GET", "/income-sources", nil, userID)
			} else {
				c, w = helpers.CreateTestContext(t, "GET", "/income-sources", nil, false)
			}

			// Call handler
			handler.GetIncomeSources(c)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var sources []models.IncomeSource
				err := json.Unmarshal(w.Body.Bytes(), &sources)
				assert.NoError(t, err)
				assert.Len(t, sources, tt.expectedCount)
			}
		})
	}
}

func TestIncomeHandler_GetIncomeSource(t *testing.T) {
	tests := []struct {
		name           string
		sourceID       string
		setupMock      func(*MockIncomeServiceInterface)
		expectedStatus int
	}{
		{
			name:     "successful retrieval",
			sourceID: "11223344-5566-7788-99aa-bbccddeeff00",
			setupMock: func(m *MockIncomeServiceInterface) {
				testSource := &models.IncomeSource{
					ID:                 uuid.MustParse("11223344-5566-7788-99aa-bbccddeeff00"),
					UserID:             uuid.MustParse("cbf3d545-d81d-450d-acb3-c5c49a597d6d"),
					Name:               "Monthly Salary",
					IncomeType:         "monthly_fixed",
					BaseAmount:         int64(500000),
					BankAccount:        uuid.MustParse("aabbccdd-eeff-1122-3344-556677889900"),
					PaymentDay:         func(i int) *int { return &i }(25),
					ScheduledDate:      nil,
					ScheduledYearMonth: nil,
					IsActive:           true,
					CreatedAt:          time.Now(),
					UpdatedAt:          time.Now(),
				}
				m.On("GetIncomeSource", uuid.MustParse("11223344-5566-7788-99aa-bbccddeeff00")).Return(testSource, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:     "invalid source ID",
			sourceID: "invalid-uuid",
			setupMock: func(m *MockIncomeServiceInterface) {
				// No mock setup needed for invalid UUID
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:     "source not found",
			sourceID: "99999999-9999-9999-9999-999999999999",
			setupMock: func(m *MockIncomeServiceInterface) {
				m.On("GetIncomeSource", uuid.MustParse("99999999-9999-9999-9999-999999999999")).Return((*models.IncomeSource)(nil), assert.AnError)
			},
			expectedStatus: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMockIncomeServiceInterface(t)
			handler := NewIncomeHandler(mockService)

			// Setup mock
			tt.setupMock(mockService)

			// Create test context
			c, w := helpers.CreateTestContext(t, "GET", fmt.Sprintf("/income-sources/%s", tt.sourceID), nil, true)
			c.Params = gin.Params{{Key: "id", Value: tt.sourceID}}

			// Call handler
			handler.GetIncomeSource(c)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestIncomeHandler_CreateIncomeSource(t *testing.T) {
	tests := []struct {
		name           string
		authenticated  bool
		requestBody    map[string]interface{}
		setupMock      func(*MockIncomeServiceInterface)
		expectedStatus int
	}{
		{
			name:          "successful creation - monthly fixed",
			authenticated: true,
			requestBody: map[string]interface{}{
				"name":         "Monthly Salary",
				"income_type":  "monthly_fixed",
				"base_amount":  int64(500000),
				"bank_account": "aabbccdd-eeff-1122-3344-556677889900",
				"payment_day":  25,
				"is_active":    true,
			},
			setupMock: func(m *MockIncomeServiceInterface) {
				m.On("CreateIncomeSource", mock.AnythingOfType("*models.IncomeSource")).Return(nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:          "successful creation - one time",
			authenticated: true,
			requestBody: map[string]interface{}{
				"name":           "Bonus Payment",
				"income_type":    "one_time",
				"base_amount":    int64(100000),
				"bank_account":   "aabbccdd-eeff-1122-3344-556677889900",
				"scheduled_date": "2024-12-25",
				"is_active":      true,
			},
			setupMock: func(m *MockIncomeServiceInterface) {
				m.On("CreateIncomeSource", mock.AnythingOfType("*models.IncomeSource")).Return(nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:          "unauthenticated user",
			authenticated: false,
			requestBody: map[string]interface{}{
				"name":         "Monthly Salary",
				"income_type":  "monthly_fixed",
				"base_amount":  int64(500000),
				"bank_account": "aabbccdd-eeff-1122-3344-556677889900",
				"payment_day":  25,
				"is_active":    true,
			},
			setupMock: func(m *MockIncomeServiceInterface) {
				// No mock setup needed for unauthenticated request
			},
			expectedStatus: http.StatusForbidden,
		},
		{
			name:          "service error",
			authenticated: true,
			requestBody: map[string]interface{}{
				"name":         "Monthly Salary",
				"income_type":  "monthly_fixed",
				"base_amount":  int64(500000),
				"bank_account": "aabbccdd-eeff-1122-3344-556677889900",
				"payment_day":  25,
				"is_active":    true,
			},
			setupMock: func(m *MockIncomeServiceInterface) {
				m.On("CreateIncomeSource", mock.AnythingOfType("*models.IncomeSource")).Return(assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMockIncomeServiceInterface(t)
			handler := NewIncomeHandler(mockService)

			// Setup mock
			tt.setupMock(mockService)

			// Create test context
			var c *gin.Context
			var w *httptest.ResponseRecorder
			if tt.authenticated {
				userID := uuid.MustParse("cbf3d545-d81d-450d-acb3-c5c49a597d6d")
				c, w = helpers.CreateTestContextWithUserID(t, "POST", "/income-sources", tt.requestBody, userID)
			} else {
				c, w = helpers.CreateTestContext(t, "POST", "/income-sources", tt.requestBody, false)
			}

			// Call handler
			handler.CreateIncomeSource(c)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestIncomeHandler_UpdateIncomeSource(t *testing.T) {
	tests := []struct {
		name           string
		sourceID       string
		requestBody    map[string]interface{}
		setupMock      func(*MockIncomeServiceInterface)
		expectedStatus int
	}{
		{
			name:     "successful update",
			sourceID: "11223344-5566-7788-99aa-bbccddeeff00",
			requestBody: map[string]interface{}{
				"name":         "Updated Salary",
				"income_type":  "monthly_fixed",
				"base_amount":  int64(550000),
				"bank_account": "aabbccdd-eeff-1122-3344-556677889900",
				"payment_day":  28,
				"is_active":    true,
			},
			setupMock: func(m *MockIncomeServiceInterface) {
				m.On("UpdateIncomeSource", mock.AnythingOfType("*models.IncomeSource")).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:     "invalid source ID",
			sourceID: "invalid-uuid",
			requestBody: map[string]interface{}{
				"name":         "Updated Salary",
				"income_type":  "monthly_fixed",
				"base_amount":  int64(550000),
				"bank_account": "aabbccdd-eeff-1122-3344-556677889900",
				"payment_day":  28,
				"is_active":    true,
			},
			setupMock: func(m *MockIncomeServiceInterface) {
				// No mock setup needed for invalid UUID
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:     "service error",
			sourceID: "11223344-5566-7788-99aa-bbccddeeff00",
			requestBody: map[string]interface{}{
				"name":         "Updated Salary",
				"income_type":  "monthly_fixed",
				"base_amount":  int64(550000),
				"bank_account": "aabbccdd-eeff-1122-3344-556677889900",
				"payment_day":  28,
				"is_active":    true,
			},
			setupMock: func(m *MockIncomeServiceInterface) {
				m.On("UpdateIncomeSource", mock.AnythingOfType("*models.IncomeSource")).Return(assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMockIncomeServiceInterface(t)
			handler := NewIncomeHandler(mockService)

			// Setup mock
			tt.setupMock(mockService)

			// Create test context
			c, w := helpers.CreateTestContext(t, "PUT", fmt.Sprintf("/income-sources/%s", tt.sourceID), tt.requestBody, true)
			c.Params = gin.Params{{Key: "id", Value: tt.sourceID}}

			// Call handler
			handler.UpdateIncomeSource(c)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestIncomeHandler_DeleteIncomeSource(t *testing.T) {
	tests := []struct {
		name           string
		sourceID       string
		setupMock      func(*MockIncomeServiceInterface)
		expectedStatus int
	}{
		{
			name:     "successful deletion",
			sourceID: "11223344-5566-7788-99aa-bbccddeeff00",
			setupMock: func(m *MockIncomeServiceInterface) {
				sourceUUID, _ := uuid.Parse("11223344-5566-7788-99aa-bbccddeeff00")
				m.On("DeleteIncomeSource", sourceUUID).Return(nil)
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:     "invalid source ID",
			sourceID: "invalid-uuid",
			setupMock: func(m *MockIncomeServiceInterface) {
				// No mock setup needed for invalid UUID
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:     "service error",
			sourceID: "99999999-9999-9999-9999-999999999999",
			setupMock: func(m *MockIncomeServiceInterface) {
				sourceUUID, _ := uuid.Parse("99999999-9999-9999-9999-999999999999")
				m.On("DeleteIncomeSource", sourceUUID).Return(assert.AnError)
			},
			expectedStatus: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockService := NewMockIncomeServiceInterface(t)
			handler := NewIncomeHandler(mockService)

			// Setup mock
			tt.setupMock(mockService)

			// Create test context
			c, w := helpers.CreateTestContext(t, "DELETE", fmt.Sprintf("/income-sources/%s", tt.sourceID), nil, true)
			c.Params = gin.Params{{Key: "id", Value: tt.sourceID}}

			// Call handler
			handler.DeleteIncomeSource(c)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
