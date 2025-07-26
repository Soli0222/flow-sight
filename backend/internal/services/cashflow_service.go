package services

import (
	"flow-sight-backend/internal/models"
	"flow-sight-backend/internal/repositories"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type CashflowService struct {
	bankAccountRepo      *repositories.BankAccountRepository
	incomeSourceRepo     *repositories.IncomeSourceRepository
	monthlyIncomeRepo    *repositories.MonthlyIncomeRepository
	recurringPaymentRepo *repositories.RecurringPaymentRepository
	cardMonthlyTotalRepo *repositories.CardMonthlyTotalRepository
	assetRepo            *repositories.AssetRepository
}

func NewCashflowService(
	bankAccountRepo *repositories.BankAccountRepository,
	incomeSourceRepo *repositories.IncomeSourceRepository,
	monthlyIncomeRepo *repositories.MonthlyIncomeRepository,
	recurringPaymentRepo *repositories.RecurringPaymentRepository,
	cardMonthlyTotalRepo *repositories.CardMonthlyTotalRepository,
	assetRepo *repositories.AssetRepository,
) *CashflowService {
	return &CashflowService{
		bankAccountRepo:      bankAccountRepo,
		incomeSourceRepo:     incomeSourceRepo,
		monthlyIncomeRepo:    monthlyIncomeRepo,
		recurringPaymentRepo: recurringPaymentRepo,
		cardMonthlyTotalRepo: cardMonthlyTotalRepo,
		assetRepo:            assetRepo,
	}
}

func (s *CashflowService) GetCashflowProjection(userID uuid.UUID, months int) ([]models.CashflowProjection, error) {
	// Get initial balance from all bank accounts
	bankAccounts, err := s.bankAccountRepo.GetAll(userID)
	if err != nil {
		return nil, err
	}

	totalBalance := int64(0)
	for _, account := range bankAccounts {
		totalBalance += account.Balance
	}

	// Get active income sources
	incomeSources, err := s.incomeSourceRepo.GetActiveByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Get active recurring payments
	recurringPayments, err := s.recurringPaymentRepo.GetActiveByUserID(userID)
	if err != nil {
		return nil, err
	}

	// Get assets (for card payments calculation)
	assets, err := s.assetRepo.GetAll(userID)
	if err != nil {
		return nil, err
	}

	// Generate cashflow projection for the specified months
	projections := make([]models.CashflowProjection, 0)
	currentBalance := totalBalance
	startDate := time.Now()

	for monthOffset := 0; monthOffset < months; monthOffset++ {
		projectionMonth := startDate.AddDate(0, monthOffset, 0)
		yearMonth := projectionMonth.Format("2006-01")

		// Get days in this month
		daysInMonth := time.Date(projectionMonth.Year(), projectionMonth.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()

		// Process each day in the month
		for day := 1; day <= daysInMonth; day++ {
			currentDate := time.Date(projectionMonth.Year(), projectionMonth.Month(), day, 0, 0, 0, 0, time.UTC)
			dayIncome := int64(0)
			dayExpense := int64(0)
			details := make([]models.CashflowProjectionDetail, 0)

			// Calculate income for this day
			for _, incomeSource := range incomeSources {
				if incomeSource.IncomeType == "monthly_fixed" {
					// Assume income comes on the 25th of each month (salary day)
					if day == 25 {
						// Check if there's a specific record for this month
						records, err := s.monthlyIncomeRepo.GetByYearMonth(yearMonth)
						if err == nil {
							for _, record := range records {
								if record.IncomeSourceID == incomeSource.ID {
									dayIncome += record.ActualAmount
									details = append(details, models.CashflowProjectionDetail{
										Type:        "income",
										Description: fmt.Sprintf("収入: %s", incomeSource.Name),
										Amount:      record.ActualAmount,
									})
									break
								}
							}
						} else {
							// Use base amount if no specific record
							dayIncome += incomeSource.BaseAmount
							details = append(details, models.CashflowProjectionDetail{
								Type:        "income",
								Description: fmt.Sprintf("収入: %s", incomeSource.Name),
								Amount:      incomeSource.BaseAmount,
							})
						}
					}
				} else if incomeSource.IncomeType == "one_time" {
					// Check if this is the scheduled month for one-time income
					if incomeSource.ScheduledYearMonth != nil && *incomeSource.ScheduledYearMonth == yearMonth {
						if day == 1 {
							dayIncome += incomeSource.BaseAmount
							details = append(details, models.CashflowProjectionDetail{
								Type:        "income",
								Description: fmt.Sprintf("臨時収入: %s", incomeSource.Name),
								Amount:      incomeSource.BaseAmount,
							})
						}
					}
				}
			}

			// Calculate recurring payments for this day
			for _, payment := range recurringPayments {
				if payment.PaymentDay == day {
					// Check if this payment is still active (for loans with remaining payments)
					if payment.RemainingPayments == nil || *payment.RemainingPayments > 0 {
						dayExpense += payment.Amount
						details = append(details, models.CashflowProjectionDetail{
							Type:        "recurring_payment",
							Description: fmt.Sprintf("固定支出: %s", payment.Name),
							Amount:      payment.Amount,
						})
					}
				}
			}

			// Calculate card payments for this day
			for _, asset := range assets {
				if asset.AssetType == "card" && asset.PaymentDay == day {
					// Calculate payment based on closing date and card usage
					paymentAmount := s.calculateCardPayment(asset, yearMonth)
					if paymentAmount > 0 {
						dayExpense += paymentAmount
						details = append(details, models.CashflowProjectionDetail{
							Type:        "card_payment",
							Description: fmt.Sprintf("カード支払い: %s", asset.Name),
							Amount:      paymentAmount,
						})
					}
				}
			}

			// Update balance
			currentBalance = currentBalance + dayIncome - dayExpense

			// Create projection for this day
			projection := models.CashflowProjection{
				Date:    currentDate.Format("2006-01-02"),
				Income:  dayIncome,
				Expense: dayExpense,
				Balance: currentBalance,
				Details: details,
			}

			projections = append(projections, projection)
		}
	}

	return projections, nil
}

func (s *CashflowService) calculateCardPayment(asset models.Asset, currentYearMonth string) int64 {
	// Parse current year-month
	parts := strings.Split(currentYearMonth, "-")
	if len(parts) != 2 {
		return 0
	}

	year, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0
	}

	month, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0
	}

	// Calculate which month's usage should be paid
	// For example, if closing day is 15 and payment day is 10:
	// - Usage from 16th of previous month to 15th of current month is paid on 10th of next month

	var targetYearMonth string
	if asset.ClosingDay != nil {
		// If current date is before closing day, we pay for previous month's usage
		// If current date is after closing day, we pay for current month's usage
		// For simplicity, let's assume we pay previous month's usage
		prevMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC).AddDate(0, -1, 0)
		targetYearMonth = prevMonth.Format("2006-01")
	} else {
		// For loans, use current month
		targetYearMonth = currentYearMonth
	}

	// Get card usage for the target month
	totals, err := s.cardMonthlyTotalRepo.GetByAssetID(asset.ID)
	if err != nil {
		return 0
	}

	for _, total := range totals {
		if total.YearMonth == targetYearMonth {
			return total.TotalAmount
		}
	}

	return 0
}
