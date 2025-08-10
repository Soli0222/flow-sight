package services

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Soli0222/flow-sight/backend/internal/models"
	"github.com/Soli0222/flow-sight/backend/internal/repositories"
)

type CashflowService struct {
	bankAccountRepo      *repositories.BankAccountRepository
	incomeSourceRepo     *repositories.IncomeSourceRepository
	monthlyIncomeRepo    *repositories.MonthlyIncomeRepository
	recurringPaymentRepo *repositories.RecurringPaymentRepository
	cardMonthlyTotalRepo *repositories.CardMonthlyTotalRepository
	creditCardRepo       *repositories.CreditCardRepository
	appSettingRepo       *repositories.AppSettingRepository
}

func NewCashflowService(
	bankAccountRepo *repositories.BankAccountRepository,
	incomeSourceRepo *repositories.IncomeSourceRepository,
	monthlyIncomeRepo *repositories.MonthlyIncomeRepository,
	recurringPaymentRepo *repositories.RecurringPaymentRepository,
	cardMonthlyTotalRepo *repositories.CardMonthlyTotalRepository,
	creditCardRepo *repositories.CreditCardRepository,
	appSettingRepo *repositories.AppSettingRepository,
) *CashflowService {
	return &CashflowService{
		bankAccountRepo:      bankAccountRepo,
		incomeSourceRepo:     incomeSourceRepo,
		monthlyIncomeRepo:    monthlyIncomeRepo,
		recurringPaymentRepo: recurringPaymentRepo,
		cardMonthlyTotalRepo: cardMonthlyTotalRepo,
		creditCardRepo:       creditCardRepo,
		appSettingRepo:       appSettingRepo,
	}
}

func (s *CashflowService) GetCashflowProjection(months int, onlyChanges bool) ([]models.CashflowProjection, error) {
	// Get initial balance from all bank accounts
	bankAccounts, err := s.bankAccountRepo.GetAll()
	if err != nil {
		return nil, err
	}

	totalBalance := int64(0)
	for _, account := range bankAccounts {
		totalBalance += account.Balance
	}

	// Get active income sources
	incomeSources, err := s.incomeSourceRepo.GetActive()
	if err != nil {
		return nil, err
	}

	// Get active recurring payments
	recurringPayments, err := s.recurringPaymentRepo.GetActive()
	if err != nil {
		return nil, err
	}

	// Get credit cards (for card payments calculation)
	creditCards, err := s.creditCardRepo.GetAll()
	if err != nil {
		return nil, err
	}

	// Get minimum monthly expense setting
	minimumMonthlyExpense := s.getMinimumMonthlyExpense()

	// Generate cashflow projection for the specified months
	projections := make([]models.CashflowProjection, 0)
	currentBalance := totalBalance
	startDate := time.Now()

	for monthOffset := 0; monthOffset < months; monthOffset++ {
		projectionMonth := startDate.AddDate(0, monthOffset, 0)
		yearMonth := projectionMonth.Format("2006-01")

		// Get days in this month
		daysInMonth := time.Date(projectionMonth.Year(), projectionMonth.Month()+1, 0, 0, 0, 0, 0, time.UTC).Day()

		// Track monthly expense total for minimum expense calculation
		monthlyExpenseTotal := int64(0)

		// Process each day in the month
		for day := 1; day <= daysInMonth; day++ {
			currentDate := time.Date(projectionMonth.Year(), projectionMonth.Month(), day, 0, 0, 0, 0, time.UTC)
			dayIncome := int64(0)
			dayExpense := int64(0)
			details := make([]models.CashflowProjectionDetail, 0)

			// Calculate income for this day
			for _, incomeSource := range incomeSources {
				if incomeSource.IncomeType == "monthly_fixed" {
					// Use payment_day if available, otherwise default to 25th
					paymentDay := 25
					if incomeSource.PaymentDay != nil {
						paymentDay = *incomeSource.PaymentDay
					}

					if day == paymentDay {

						// Check if there's a specific record for this month
						records, err := s.monthlyIncomeRepo.GetByYearMonth(yearMonth)
						if err == nil {
							recordFound := false
							for _, record := range records {
								if record.IncomeSourceID == incomeSource.ID {
									dayIncome += record.ActualAmount
									details = append(details, models.CashflowProjectionDetail{
										Type:        "income",
										Description: fmt.Sprintf("収入: %s", incomeSource.Name),
										Amount:      record.ActualAmount,
									})
									recordFound = true
									break
								}
							}
							// Use base amount if no specific record found
							if !recordFound {
								dayIncome += incomeSource.BaseAmount
								details = append(details, models.CashflowProjectionDetail{
									Type:        "income",
									Description: fmt.Sprintf("収入: %s", incomeSource.Name),
									Amount:      incomeSource.BaseAmount,
								})
							}
						} else {
							// Use base amount if query failed
							dayIncome += incomeSource.BaseAmount
							details = append(details, models.CashflowProjectionDetail{
								Type:        "income",
								Description: fmt.Sprintf("収入: %s", incomeSource.Name),
								Amount:      incomeSource.BaseAmount,
							})
						}
					}
				} else if incomeSource.IncomeType == "one_time" {
					// Check if this is the scheduled date for one-time income
					if incomeSource.ScheduledDate != nil {

						// Try multiple date formats to parse the scheduled date
						var scheduledDate time.Time
						var err error

						// First try "YYYY-MM-DD" format
						scheduledDate, err = time.Parse("2006-01-02", *incomeSource.ScheduledDate)
						if err != nil {
							// If that fails, try with time component
							scheduledDate, err = time.Parse("2006-01-02T15:04:05Z07:00", *incomeSource.ScheduledDate)
							if err != nil {
								// Try ISO format without timezone
								scheduledDate, err = time.Parse("2006-01-02T15:04:05", *incomeSource.ScheduledDate)
								if err != nil {
									continue
								}
							}
						}

						// Compare only the date part
						if scheduledDate.Year() == currentDate.Year() &&
							scheduledDate.Month() == currentDate.Month() &&
							scheduledDate.Day() == currentDate.Day() {
							dayIncome += incomeSource.BaseAmount
							details = append(details, models.CashflowProjectionDetail{
								Type:        "income",
								Description: fmt.Sprintf("臨時収入: %s", incomeSource.Name),
								Amount:      incomeSource.BaseAmount,
							})
						}
					} else if incomeSource.ScheduledYearMonth != nil && *incomeSource.ScheduledYearMonth == yearMonth {
						// Fallback to first day of month for backward compatibility
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
					// Check if this payment should be applied in this month
					shouldApplyPayment := s.shouldApplyRecurringPayment(payment, yearMonth)

					if shouldApplyPayment {
						dayExpense += payment.Amount
						monthlyExpenseTotal += payment.Amount
						details = append(details, models.CashflowProjectionDetail{
							Type:        "recurring_payment",
							Description: fmt.Sprintf("固定支出: %s", payment.Name),
							Amount:      payment.Amount,
						})
					}
				}
			}

			// Calculate card payments for this day
			for _, creditCard := range creditCards {
				if creditCard.PaymentDay == day {
					// Calculate payment based on closing date and card usage
					paymentAmount := s.calculateCardPayment(creditCard, yearMonth)
					if paymentAmount > 0 {
						dayExpense += paymentAmount
						monthlyExpenseTotal += paymentAmount
						details = append(details, models.CashflowProjectionDetail{
							Type:        "card_payment",
							Description: fmt.Sprintf("カード支払い: %s", creditCard.Name),
							Amount:      paymentAmount,
						})
					}
				}
			}

			// Check for minimum monthly expense on the 26th day of 3rd month onwards
			if monthOffset >= 2 && day == 26 && minimumMonthlyExpense > 0 {
				if monthlyExpenseTotal < minimumMonthlyExpense {
					shortfall := minimumMonthlyExpense - monthlyExpenseTotal
					dayExpense += shortfall
					monthlyExpenseTotal += shortfall
					details = append(details, models.CashflowProjectionDetail{
						Type:        "recurring_payment",
						Description: "最低月支出調整",
						Amount:      shortfall,
					})
				}
			}

			// Update balance
			currentBalance = currentBalance + dayIncome - dayExpense

			// Create projection for this day only if there are changes or if onlyChanges is false
			if !onlyChanges || dayIncome > 0 || dayExpense > 0 {
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
	}

	return projections, nil
}

// getMinimumMonthlyExpense retrieves the minimum monthly expense setting
func (s *CashflowService) getMinimumMonthlyExpense() int64 {
	settings, err := s.appSettingRepo.GetAll()
	if err != nil {
		return 0
	}

	for _, setting := range settings {
		if setting.Key == "minimum_monthly_expense" {
			amount, err := strconv.ParseInt(setting.Value, 10, 64)
			if err != nil {
				return 0
			}
			return amount
		}
	}

	return 0
}

func (s *CashflowService) calculateCardPayment(creditCard models.CreditCard, currentYearMonth string) int64 {
	// Parse year and month
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
	if creditCard.ClosingDay != nil {
		// If current date is before closing day, we pay for previous month's usage
		// If current date is after closing day, we pay for current month's usage
		// For simplicity, let's assume we pay previous month's usage
		prevMonth := time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.UTC).AddDate(0, -1, 0)
		targetYearMonth = prevMonth.Format("2006-01")
	} else {
		// For credit cards without closing day, use current month
		targetYearMonth = currentYearMonth
	}

	// Get card usage for the target month
	totals, err := s.cardMonthlyTotalRepo.GetByCreditCardID(creditCard.ID)
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

// shouldApplyRecurringPayment determines if a recurring payment should be applied in the given month
func (s *CashflowService) shouldApplyRecurringPayment(payment models.RecurringPayment, targetYearMonth string) bool {
	// If payment is not active, don't apply
	if !payment.IsActive {
		return false
	}

	// Parse start year-month and target year-month
	startYear, startMonth, err := parseYearMonth(payment.StartYearMonth)
	if err != nil {
		return false
	}

	targetYear, targetMonth, err := parseYearMonth(targetYearMonth)
	if err != nil {
		return false
	}

	// Calculate months elapsed since start
	monthsElapsed := (targetYear-startYear)*12 + (targetMonth - startMonth)

	// If target month is before start month, don't apply
	if monthsElapsed < 0 {
		return false
	}

	// If no total payments specified (infinite payments), apply if active
	if payment.TotalPayments == nil || *payment.TotalPayments == 0 {
		return true
	}

	// Calculate current payment number (1-based)
	currentPaymentNumber := monthsElapsed + 1

	// Check if we haven't exceeded the total payment count
	return currentPaymentNumber <= *payment.TotalPayments
}

// parseYearMonth parses a year-month string like "2024-01" into year and month integers
func parseYearMonth(yearMonth string) (int, int, error) {
	parts := strings.Split(yearMonth, "-")
	if len(parts) != 2 {
		return 0, 0, fmt.Errorf("invalid year-month format: %s", yearMonth)
	}

	year, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid year in year-month: %s", yearMonth)
	}

	month, err := strconv.Atoi(parts[1])
	if err != nil {
		return 0, 0, fmt.Errorf("invalid month in year-month: %s", yearMonth)
	}

	return year, month, nil
}
