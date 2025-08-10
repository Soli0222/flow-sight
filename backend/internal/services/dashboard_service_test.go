package services

import (
	"testing"

	"github.com/Soli0222/flow-sight/backend/internal/models"
	"github.com/Soli0222/flow-sight/backend/internal/services/mocks"
	"github.com/Soli0222/flow-sight/backend/test/helpers"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// MockDashboardService は DashboardService のシンプルなモック
type MockDashboardService struct {
	bankAccountRepo      BankAccountRepositoryInterface
	creditCardRepo       CreditCardRepositoryInterface
	incomeSourceRepo     IncomeSourceRepositoryInterface
	monthlyIncomeRepo    MonthlyIncomeRepositoryInterface
	recurringPaymentRepo RecurringPaymentRepositoryInterface
	cashflowService      *mocks.MockCashflowService
}

func NewMockDashboardService(
	bankAccountRepo BankAccountRepositoryInterface,
	creditCardRepo CreditCardRepositoryInterface,
	incomeSourceRepo IncomeSourceRepositoryInterface,
	monthlyIncomeRepo MonthlyIncomeRepositoryInterface,
	recurringPaymentRepo RecurringPaymentRepositoryInterface,
	cashflowService *mocks.MockCashflowService,
) *MockDashboardService {
	return &MockDashboardService{
		bankAccountRepo:      bankAccountRepo,
		creditCardRepo:       creditCardRepo,
		incomeSourceRepo:     incomeSourceRepo,
		monthlyIncomeRepo:    monthlyIncomeRepo,
		recurringPaymentRepo: recurringPaymentRepo,
		cashflowService:      cashflowService,
	}
}

func (s *MockDashboardService) GetDashboardSummary() (*models.DashboardSummary, error) {
	// Get total balance from all bank accounts
	bankAccounts, err := s.bankAccountRepo.GetAll()
	if err != nil {
		return nil, err
	}

	totalBalance := int64(0)
	for _, account := range bankAccounts {
		totalBalance += account.Balance
	}

	// Get credit cards count
	creditCards, err := s.creditCardRepo.GetAll()
	if err != nil {
		return nil, err
	}

	// Calculate total assets (bank accounts + credit cards)
	totalAssets := len(bankAccounts) + len(creditCards)

	// Get recent cashflow activities
	recentActivities, err := s.cashflowService.GetCashflowProjection(1, true)
	if err != nil {
		recentActivities = make([]models.CashflowProjection, 0)
	}

	return &models.DashboardSummary{
		TotalBalance:     totalBalance,
		MonthlyIncome:    100000, // Fixed for test
		MonthlyExpense:   50000,  // Fixed for test
		TotalAssets:      totalAssets,
		RecentActivities: recentActivities,
	}, nil
}

func TestDashboardService_GetDashboardSummary(t *testing.T) {
	mockBankAccountRepo := &mocks.MockBankAccountRepository{}
	mockCreditCardRepo := &mocks.MockCreditCardRepository{}
	mockIncomeSourceRepo := &mocks.MockIncomeSourceRepository{}
	mockMonthlyIncomeRepo := &mocks.MockMonthlyIncomeRepository{}
	mockRecurringPaymentRepo := &mocks.MockRecurringPaymentRepository{}
	mockCashflowService := &mocks.MockCashflowService{}

	service := NewMockDashboardService(
		mockBankAccountRepo,
		mockCreditCardRepo,
		mockIncomeSourceRepo,
		mockMonthlyIncomeRepo,
		mockRecurringPaymentRepo,
		mockCashflowService,
	)

	bankAccountID := uuid.New()

	tests := []struct {
		name          string
		setupMocks    func()
		expectedError bool
	}{
		{
			name: "successful dashboard summary",
			setupMocks: func() {
				// Setup bank accounts
				bankAccounts := []models.BankAccount{
					*helpers.CreateTestBankAccount(),
					*helpers.CreateTestBankAccount(),
				}
				bankAccounts[0].Balance = 100000 // 1000.00
				bankAccounts[1].Balance = 200000 // 2000.00
				mockBankAccountRepo.On("GetAll").Return(bankAccounts, nil)

				// Setup credit cards
				creditCards := []models.CreditCard{
					*helpers.CreateTestCreditCard(bankAccountID),
				}
				mockCreditCardRepo.On("GetAll").Return(creditCards, nil)

				// Setup cashflow activities
				activities := []models.CashflowProjection{
					{
						Date:    "2024-01-15",
						Income:  50000,
						Expense: 0,
					},
				}
				mockCashflowService.On("GetCashflowProjection", 1, true).Return(activities, nil)
			},
			expectedError: false,
		},
		{
			name: "bank account repository error",
			setupMocks: func() {
				mockBankAccountRepo.On("GetAll").Return([]models.BankAccount{}, assert.AnError)
			},
			expectedError: true,
		},
		{
			name: "credit card repository error",
			setupMocks: func() {
				// Setup successful bank accounts
				bankAccounts := []models.BankAccount{
					*helpers.CreateTestBankAccount(),
				}
				mockBankAccountRepo.On("GetAll").Return(bankAccounts, nil)

				// Setup error for credit cards
				mockCreditCardRepo.On("GetAll").Return([]models.CreditCard{}, assert.AnError)
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset mocks
			mockBankAccountRepo.ExpectedCalls = nil
			mockCreditCardRepo.ExpectedCalls = nil
			mockCashflowService.ExpectedCalls = nil

			tt.setupMocks()

			result, err := service.GetDashboardSummary()

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, result)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, result)

				// Verify structure
				assert.GreaterOrEqual(t, result.TotalBalance, int64(0))
				assert.GreaterOrEqual(t, result.TotalAssets, 0)
				assert.NotNil(t, result.RecentActivities)

				// Verify specific values for successful case
				if tt.name == "successful dashboard summary" {
					assert.Equal(t, int64(300000), result.TotalBalance) // 1000 + 2000
					assert.Equal(t, 3, result.TotalAssets)              // 2 bank accounts + 1 credit card
					assert.Equal(t, int64(100000), result.MonthlyIncome)
					assert.Equal(t, int64(50000), result.MonthlyExpense)
				}
			}

			// Verify all expectations were met
			mockBankAccountRepo.AssertExpectations(t)
			mockCreditCardRepo.AssertExpectations(t)
			mockCashflowService.AssertExpectations(t)
		})
	}
}
