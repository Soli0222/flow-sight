package services

import (
	"github.com/Soli0222/flow-sight/backend/internal/models"
	"github.com/Soli0222/flow-sight/backend/internal/repositories"
	"time"

	"github.com/google/uuid"
)

type DashboardService struct {
	bankAccountRepo      *repositories.BankAccountRepository
	creditCardRepo       *repositories.CreditCardRepository
	incomeSourceRepo     *repositories.IncomeSourceRepository
	monthlyIncomeRepo    *repositories.MonthlyIncomeRepository
	recurringPaymentRepo *repositories.RecurringPaymentRepository
	cashflowService      *CashflowService
}

func NewDashboardService(
	bankAccountRepo *repositories.BankAccountRepository,
	creditCardRepo *repositories.CreditCardRepository,
	incomeSourceRepo *repositories.IncomeSourceRepository,
	monthlyIncomeRepo *repositories.MonthlyIncomeRepository,
	recurringPaymentRepo *repositories.RecurringPaymentRepository,
	cashflowService *CashflowService,
) *DashboardService {
	return &DashboardService{
		bankAccountRepo:      bankAccountRepo,
		creditCardRepo:       creditCardRepo,
		incomeSourceRepo:     incomeSourceRepo,
		monthlyIncomeRepo:    monthlyIncomeRepo,
		recurringPaymentRepo: recurringPaymentRepo,
		cashflowService:      cashflowService,
	}
}

func (s *DashboardService) GetDashboardSummary(userID uuid.UUID) (*models.DashboardSummary, error) {
	// Get total balance from all bank accounts
	bankAccounts, err := s.bankAccountRepo.GetAll(userID)
	if err != nil {
		return nil, err
	}

	totalBalance := int64(0)
	for _, account := range bankAccounts {
		totalBalance += account.Balance
	}

	// Get credit cards count
	creditCards, err := s.creditCardRepo.GetAll(userID)
	if err != nil {
		return nil, err
	}

	// Calculate total assets (bank accounts + credit cards)
	totalAssets := len(bankAccounts) + len(creditCards)

	// Get current month string (YYYY-MM)
	currentTime := time.Now()
	currentYearMonth := currentTime.Format("2006-01")

	// Calculate monthly income and expense
	monthlyIncome, monthlyExpense, err := s.calculateMonthlyIncomeExpense(userID, currentYearMonth)
	if err != nil {
		return nil, err
	}

	// Get recent cashflow activities
	recentActivities, err := s.cashflowService.GetCashflowProjection(userID, 1, true)
	if err != nil {
		// If cashflow fails, continue with empty activities
		recentActivities = make([]models.CashflowProjection, 0)
	}

	// Ensure recentActivities is not nil
	if recentActivities == nil {
		recentActivities = make([]models.CashflowProjection, 0)
	}

	// Filter recent activities to only current month and limit to 5
	var filteredActivities []models.CashflowProjection
	for _, activity := range recentActivities {
		if len(activity.Date) >= 7 && activity.Date[:7] == currentYearMonth {
			if activity.Income > 0 || activity.Expense > 0 {
				filteredActivities = append(filteredActivities, activity)
				if len(filteredActivities) >= 5 {
					break
				}
			}
		}
	}

	// Ensure filteredActivities is not nil
	if filteredActivities == nil {
		filteredActivities = make([]models.CashflowProjection, 0)
	}

	return &models.DashboardSummary{
		TotalBalance:     totalBalance,
		MonthlyIncome:    monthlyIncome,
		MonthlyExpense:   monthlyExpense,
		TotalAssets:      totalAssets,
		RecentActivities: filteredActivities,
	}, nil
}

func (s *DashboardService) calculateMonthlyIncomeExpense(userID uuid.UUID, yearMonth string) (int64, int64, error) {
	var totalIncome int64 = 0
	var totalExpense int64 = 0

	// Get active income sources
	incomeSources, err := s.incomeSourceRepo.GetActiveByUserID(userID)
	if err != nil {
		return 0, 0, err
	}

	// Calculate monthly income
	for _, source := range incomeSources {
		if source.IncomeType == "monthly_fixed" {
			// Check if there's a specific record for this month
			records, err := s.monthlyIncomeRepo.GetByUserIDAndYearMonth(userID, yearMonth)
			if err == nil {
				recordFound := false
				for _, record := range records {
					if record.IncomeSourceID == source.ID {
						totalIncome += record.ActualAmount
						recordFound = true
						break
					}
				}
				// Use base amount if no specific record found
				if !recordFound {
					totalIncome += source.BaseAmount
				}
			} else {
				// Use base amount if query failed
				totalIncome += source.BaseAmount
			}
		} else if source.IncomeType == "one_time" {
			// Check if this one-time income is scheduled for current month
			if source.ScheduledYearMonth != nil && *source.ScheduledYearMonth == yearMonth {
				totalIncome += source.BaseAmount
			}
		}
	}

	// Get active recurring payments
	recurringPayments, err := s.recurringPaymentRepo.GetActiveByUserID(userID)
	if err != nil {
		return totalIncome, 0, err
	}

	// Calculate monthly expense from recurring payments
	for _, payment := range recurringPayments {
		// Check if this payment is still active (for loans with remaining payments)
		if payment.RemainingPayments == nil || *payment.RemainingPayments > 0 {
			totalExpense += payment.Amount
		}
	}

	return totalIncome, totalExpense, nil
}
