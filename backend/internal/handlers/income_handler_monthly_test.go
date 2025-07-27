package handlers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"flow-sight-backend/internal/models"
	"flow-sight-backend/test/helpers"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestIncomeHandler_GetMonthlyIncomeRecords(t *testing.T) {
	tests := []struct {
		name           string
		authenticated  bool
		incomeSourceID string
		setupMock      func(*MockIncomeServiceInterface)
		expectedStatus int
		expectedCount  int
	}{
		{
			name:           "successful retrieval",
			authenticated:  true,
			incomeSourceID: "11223344-5566-7788-99aa-bbccddeeff00",
			setupMock: func(m *MockIncomeServiceInterface) {
				testRecords := []models.MonthlyIncomeRecord{
					{
						ID:             uuid.MustParse("ffffffff-ffff-ffff-ffff-ffffffffffff"),
						IncomeSourceID: uuid.MustParse("11223344-5566-7788-99aa-bbccddeeff00"),
						YearMonth:      "2024-12",
						ActualAmount:   int64(500000),
						IsConfirmed:    true,
						Note:           "December salary",
						CreatedAt:      time.Now(),
						UpdatedAt:      time.Now(),
					},
				}
				m.On("GetMonthlyIncomeRecords", uuid.MustParse("11223344-5566-7788-99aa-bbccddeeff00")).Return(testRecords, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCount:  1,
		},
		{
			name:           "unauthenticated user - still processes request",
			authenticated:  false,
			incomeSourceID: "11223344-5566-7788-99aa-bbccddeeff00",
			setupMock: func(m *MockIncomeServiceInterface) {
				// Handler doesn't check auth, so it will still call service
				testRecords := []models.MonthlyIncomeRecord{
					{
						ID:             uuid.MustParse("ffffffff-ffff-ffff-ffff-ffffffffffff"),
						IncomeSourceID: uuid.MustParse("11223344-5566-7788-99aa-bbccddeeff00"),
						YearMonth:      "2024-12",
						ActualAmount:   int64(500000),
						IsConfirmed:    true,
						Note:           "December salary",
						CreatedAt:      time.Now(),
						UpdatedAt:      time.Now(),
					},
				}
				m.On("GetMonthlyIncomeRecords", uuid.MustParse("11223344-5566-7788-99aa-bbccddeeff00")).Return(testRecords, nil)
			},
			expectedStatus: http.StatusOK,
			expectedCount:  1,
		},
		{
			name:           "invalid income source ID",
			authenticated:  true,
			incomeSourceID: "invalid-uuid",
			setupMock: func(m *MockIncomeServiceInterface) {
				// No mock setup needed for invalid UUID
			},
			expectedStatus: http.StatusBadRequest,
			expectedCount:  0,
		},
		{
			name:           "service error",
			authenticated:  true,
			incomeSourceID: "11223344-5566-7788-99aa-bbccddeeff00",
			setupMock: func(m *MockIncomeServiceInterface) {
				m.On("GetMonthlyIncomeRecords", uuid.MustParse("11223344-5566-7788-99aa-bbccddeeff00")).Return([]models.MonthlyIncomeRecord{}, assert.AnError)
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
				c, w = helpers.CreateTestContextWithUserID(t, "GET", fmt.Sprintf("/monthly-income-records?income_source_id=%s", tt.incomeSourceID), nil, userID)
			} else {
				c, w = helpers.CreateTestContext(t, "GET", fmt.Sprintf("/monthly-income-records?income_source_id=%s", tt.incomeSourceID), nil, false)
			}

			// Call handler
			handler.GetMonthlyIncomeRecords(c)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)

			if tt.expectedStatus == http.StatusOK {
				var records []models.MonthlyIncomeRecord
				err := json.Unmarshal(w.Body.Bytes(), &records)
				assert.NoError(t, err)
				assert.Len(t, records, tt.expectedCount)
			}
		})
	}
}

func TestIncomeHandler_GetMonthlyIncomeRecord(t *testing.T) {
	tests := []struct {
		name           string
		recordID       string
		setupMock      func(*MockIncomeServiceInterface)
		expectedStatus int
	}{
		{
			name:     "successful retrieval",
			recordID: "ffffffff-ffff-ffff-ffff-ffffffffffff",
			setupMock: func(m *MockIncomeServiceInterface) {
				testRecord := &models.MonthlyIncomeRecord{
					ID:             uuid.MustParse("ffffffff-ffff-ffff-ffff-ffffffffffff"),
					IncomeSourceID: uuid.MustParse("11223344-5566-7788-99aa-bbccddeeff00"),
					YearMonth:      "2024-12",
					ActualAmount:   int64(500000),
					IsConfirmed:    true,
					Note:           "December salary",
					CreatedAt:      time.Now(),
					UpdatedAt:      time.Now(),
				}
				m.On("GetMonthlyIncomeRecord", uuid.MustParse("ffffffff-ffff-ffff-ffff-ffffffffffff")).Return(testRecord, nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:     "invalid record ID",
			recordID: "invalid-uuid",
			setupMock: func(m *MockIncomeServiceInterface) {
				// No mock setup needed for invalid UUID
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:     "record not found",
			recordID: "99999999-9999-9999-9999-999999999999",
			setupMock: func(m *MockIncomeServiceInterface) {
				m.On("GetMonthlyIncomeRecord", uuid.MustParse("99999999-9999-9999-9999-999999999999")).Return((*models.MonthlyIncomeRecord)(nil), assert.AnError)
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
			c, w := helpers.CreateTestContext(t, "GET", fmt.Sprintf("/monthly-income-records/%s", tt.recordID), nil, true)
			c.Params = gin.Params{{Key: "id", Value: tt.recordID}}

			// Call handler
			handler.GetMonthlyIncomeRecord(c)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestIncomeHandler_CreateMonthlyIncomeRecord(t *testing.T) {
	tests := []struct {
		name           string
		authenticated  bool
		requestBody    map[string]interface{}
		setupMock      func(*MockIncomeServiceInterface)
		expectedStatus int
	}{
		{
			name:          "successful creation",
			authenticated: true,
			requestBody: map[string]interface{}{
				"income_source_id": "11223344-5566-7788-99aa-bbccddeeff00",
				"year_month":       "2024-12",
				"actual_amount":    int64(500000),
				"is_confirmed":     true,
				"note":             "December salary",
			},
			setupMock: func(m *MockIncomeServiceInterface) {
				m.On("CreateMonthlyIncomeRecord", mock.AnythingOfType("*models.MonthlyIncomeRecord")).Return(nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:          "unauthenticated user - still processes request",
			authenticated: false,
			requestBody: map[string]interface{}{
				"income_source_id": "11223344-5566-7788-99aa-bbccddeeff00",
				"year_month":       "2024-12",
				"actual_amount":    int64(500000),
				"is_confirmed":     true,
				"note":             "December salary",
			},
			setupMock: func(m *MockIncomeServiceInterface) {
				// Handler doesn't check auth, so it will still call service
				m.On("CreateMonthlyIncomeRecord", mock.AnythingOfType("*models.MonthlyIncomeRecord")).Return(nil)
			},
			expectedStatus: http.StatusCreated,
		},
		{
			name:          "service error",
			authenticated: true,
			requestBody: map[string]interface{}{
				"income_source_id": "11223344-5566-7788-99aa-bbccddeeff00",
				"year_month":       "2024-12",
				"actual_amount":    int64(500000),
				"is_confirmed":     true,
				"note":             "December salary",
			},
			setupMock: func(m *MockIncomeServiceInterface) {
				m.On("CreateMonthlyIncomeRecord", mock.AnythingOfType("*models.MonthlyIncomeRecord")).Return(assert.AnError)
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
				c, w = helpers.CreateTestContextWithUserID(t, "POST", "/monthly-income-records", tt.requestBody, userID)
			} else {
				c, w = helpers.CreateTestContext(t, "POST", "/monthly-income-records", tt.requestBody, false)
			}

			// Call handler
			handler.CreateMonthlyIncomeRecord(c)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestIncomeHandler_UpdateMonthlyIncomeRecord(t *testing.T) {
	tests := []struct {
		name           string
		recordID       string
		requestBody    map[string]interface{}
		setupMock      func(*MockIncomeServiceInterface)
		expectedStatus int
	}{
		{
			name:     "successful update",
			recordID: "ffffffff-ffff-ffff-ffff-ffffffffffff",
			requestBody: map[string]interface{}{
				"income_source_id": "11223344-5566-7788-99aa-bbccddeeff00",
				"year_month":       "2024-12",
				"actual_amount":    int64(550000),
				"is_confirmed":     true,
				"note":             "Updated December salary",
			},
			setupMock: func(m *MockIncomeServiceInterface) {
				m.On("UpdateMonthlyIncomeRecord", mock.AnythingOfType("*models.MonthlyIncomeRecord")).Return(nil)
			},
			expectedStatus: http.StatusOK,
		},
		{
			name:     "invalid record ID",
			recordID: "invalid-uuid",
			requestBody: map[string]interface{}{
				"income_source_id": "11223344-5566-7788-99aa-bbccddeeff00",
				"year_month":       "2024-12",
				"actual_amount":    int64(550000),
				"is_confirmed":     true,
				"note":             "Updated December salary",
			},
			setupMock: func(m *MockIncomeServiceInterface) {
				// No mock setup needed for invalid UUID
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:     "service error",
			recordID: "ffffffff-ffff-ffff-ffff-ffffffffffff",
			requestBody: map[string]interface{}{
				"income_source_id": "11223344-5566-7788-99aa-bbccddeeff00",
				"year_month":       "2024-12",
				"actual_amount":    int64(550000),
				"is_confirmed":     true,
				"note":             "Updated December salary",
			},
			setupMock: func(m *MockIncomeServiceInterface) {
				m.On("UpdateMonthlyIncomeRecord", mock.AnythingOfType("*models.MonthlyIncomeRecord")).Return(assert.AnError)
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
			c, w := helpers.CreateTestContext(t, "PUT", fmt.Sprintf("/monthly-income-records/%s", tt.recordID), tt.requestBody, true)
			c.Params = gin.Params{{Key: "id", Value: tt.recordID}}

			// Call handler
			handler.UpdateMonthlyIncomeRecord(c)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}

func TestIncomeHandler_DeleteMonthlyIncomeRecord(t *testing.T) {
	tests := []struct {
		name           string
		recordID       string
		setupMock      func(*MockIncomeServiceInterface)
		expectedStatus int
	}{
		{
			name:     "successful deletion",
			recordID: "ffffffff-ffff-ffff-ffff-ffffffffffff",
			setupMock: func(m *MockIncomeServiceInterface) {
				recordUUID, _ := uuid.Parse("ffffffff-ffff-ffff-ffff-ffffffffffff")
				m.On("DeleteMonthlyIncomeRecord", recordUUID).Return(nil)
			},
			expectedStatus: http.StatusNoContent,
		},
		{
			name:     "invalid record ID",
			recordID: "invalid-uuid",
			setupMock: func(m *MockIncomeServiceInterface) {
				// No mock setup needed for invalid UUID
			},
			expectedStatus: http.StatusBadRequest,
		},
		{
			name:     "service error",
			recordID: "99999999-9999-9999-9999-999999999999",
			setupMock: func(m *MockIncomeServiceInterface) {
				recordUUID, _ := uuid.Parse("99999999-9999-9999-9999-999999999999")
				m.On("DeleteMonthlyIncomeRecord", recordUUID).Return(assert.AnError)
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
			c, w := helpers.CreateTestContext(t, "DELETE", fmt.Sprintf("/monthly-income-records/%s", tt.recordID), nil, true)
			c.Params = gin.Params{{Key: "id", Value: tt.recordID}}

			// Call handler
			handler.DeleteMonthlyIncomeRecord(c)

			// Assert response
			assert.Equal(t, tt.expectedStatus, w.Code)
		})
	}
}
