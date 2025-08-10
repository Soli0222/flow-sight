package services

import (
	"testing"
	"time"

	"github.com/Soli0222/flow-sight/backend/internal/models"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCashflowService_shouldApplyRecurringPayment(t *testing.T) {
	service := &CashflowService{}
	bankAccountID := uuid.New()

	tests := []struct {
		name            string
		payment         models.RecurringPayment
		targetYearMonth string
		expectedResult  bool
		description     string
	}{
		{
			name: "infinite payments - should apply",
			payment: models.RecurringPayment{
				ID:                uuid.New(),
				Name:              "Monthly Rent",
				Amount:            100000,
				PaymentDay:        1,
				StartYearMonth:    "2024-01",
				TotalPayments:     nil, // infinite
				RemainingPayments: nil,
				BankAccount:       bankAccountID,
				IsActive:          true,
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
			},
			targetYearMonth: "2024-06",
			expectedResult:  true,
			description:     "Infinite payments should always apply when active",
		},
		{
			name: "zero total payments - should apply",
			payment: models.RecurringPayment{
				ID:                uuid.New(),
				Name:              "Monthly Rent",
				Amount:            100000,
				PaymentDay:        1,
				StartYearMonth:    "2024-01",
				TotalPayments:     func() *int { i := 0; return &i }(), // 0 means infinite
				RemainingPayments: nil,
				BankAccount:       bankAccountID,
				IsActive:          true,
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
			},
			targetYearMonth: "2024-06",
			expectedResult:  true,
			description:     "Zero total payments should apply (means infinite)",
		},
		{
			name: "within payment period - should apply",
			payment: models.RecurringPayment{
				ID:                uuid.New(),
				Name:              "Loan Payment",
				Amount:            50000,
				PaymentDay:        15,
				StartYearMonth:    "2024-01",
				TotalPayments:     func() *int { i := 12; return &i }(), // 12 months
				RemainingPayments: func() *int { i := 8; return &i }(),
				BankAccount:       bankAccountID,
				IsActive:          true,
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
			},
			targetYearMonth: "2024-06", // 6th month, within 12 payments
			expectedResult:  true,
			description:     "Payment within the total payment period should apply",
		},
		{
			name: "exceeds payment period - should not apply",
			payment: models.RecurringPayment{
				ID:                uuid.New(),
				Name:              "Loan Payment",
				Amount:            50000,
				PaymentDay:        15,
				StartYearMonth:    "2024-01",
				TotalPayments:     func() *int { i := 6; return &i }(), // 6 months only
				RemainingPayments: func() *int { i := 0; return &i }(),
				BankAccount:       bankAccountID,
				IsActive:          true,
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
			},
			targetYearMonth: "2024-08", // 8th month, exceeds 6 payments
			expectedResult:  false,
			description:     "Payment beyond the total payment period should not apply",
		},
		{
			name: "before start month - should not apply",
			payment: models.RecurringPayment{
				ID:                uuid.New(),
				Name:              "Future Payment",
				Amount:            30000,
				PaymentDay:        10,
				StartYearMonth:    "2024-06",
				TotalPayments:     func() *int { i := 12; return &i }(),
				RemainingPayments: func() *int { i := 12; return &i }(),
				BankAccount:       bankAccountID,
				IsActive:          true,
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
			},
			targetYearMonth: "2024-03", // Before start month
			expectedResult:  false,
			description:     "Payment before start month should not apply",
		},
		{
			name: "inactive payment - should not apply",
			payment: models.RecurringPayment{
				ID:                uuid.New(),
				Name:              "Inactive Payment",
				Amount:            25000,
				PaymentDay:        5,
				StartYearMonth:    "2024-01",
				TotalPayments:     nil,
				RemainingPayments: nil,
				BankAccount:       bankAccountID,
				IsActive:          false, // Inactive
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
			},
			targetYearMonth: "2024-06",
			expectedResult:  false,
			description:     "Inactive payments should not apply",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := service.shouldApplyRecurringPayment(tt.payment, tt.targetYearMonth)
			assert.Equal(t, tt.expectedResult, result, tt.description)
		})
	}
}

func TestParseYearMonth(t *testing.T) {
	tests := []struct {
		name          string
		input         string
		expectedYear  int
		expectedMonth int
		expectedError bool
	}{
		{
			name:          "valid format",
			input:         "2024-06",
			expectedYear:  2024,
			expectedMonth: 6,
			expectedError: false,
		},
		{
			name:          "invalid format - no dash",
			input:         "202406",
			expectedYear:  0,
			expectedMonth: 0,
			expectedError: true,
		},
		{
			name:          "invalid format - too many parts",
			input:         "2024-06-15",
			expectedYear:  0,
			expectedMonth: 0,
			expectedError: true,
		},
		{
			name:          "invalid year",
			input:         "abc-06",
			expectedYear:  0,
			expectedMonth: 0,
			expectedError: true,
		},
		{
			name:          "invalid month",
			input:         "2024-abc",
			expectedYear:  0,
			expectedMonth: 0,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			year, month, err := parseYearMonth(tt.input)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Equal(t, tt.expectedYear, year)
				assert.Equal(t, tt.expectedMonth, month)
			}
		})
	}
}
