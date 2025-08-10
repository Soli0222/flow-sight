package helpers

import (
	"time"

	"github.com/Soli0222/flow-sight/backend/internal/models"

	"github.com/google/uuid"
)

// CreateTestBankAccount creates a test bank account with default values
func CreateTestBankAccount() *models.BankAccount {
	return &models.BankAccount{
		ID:        uuid.New(),
		Name:      "Test Bank Account",
		Balance:   100000, // 1000.00 in cents
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// CreateTestCreditCard creates a test credit card with default values
func CreateTestCreditCard(bankAccountID uuid.UUID) *models.CreditCard {
	closingDay := 25
	return &models.CreditCard{
		ID:          uuid.New(),
		Name:        "Test Credit Card",
		ClosingDay:  &closingDay,
		PaymentDay:  10,
		BankAccount: bankAccountID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// CreateTestIncomeSource creates a test income source with default values
func CreateTestIncomeSource(bankAccountID uuid.UUID) *models.IncomeSource {
	paymentDay := 25
	return &models.IncomeSource{
		ID:          uuid.New(),
		Name:        "Test Income Source",
		IncomeType:  "monthly_fixed",
		BaseAmount:  300000, // 3000.00 in cents
		BankAccount: bankAccountID,
		PaymentDay:  &paymentDay,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// CreateTestRecurringPayment creates a test recurring payment with default values
func CreateTestRecurringPayment(bankAccountID uuid.UUID) *models.RecurringPayment {
	return &models.RecurringPayment{
		ID:             uuid.New(),
		Name:           "Test Recurring Payment",
		Amount:         50000, // 500.00 in cents
		PaymentDay:     15,
		StartYearMonth: "2024-01",
		BankAccount:    bankAccountID,
		IsActive:       true,
		Note:           "Test recurring payment",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}

// CreateTestCardMonthlyTotal creates a test CardMonthlyTotal with default values
func CreateTestCardMonthlyTotal() *models.CardMonthlyTotal {
	return &models.CardMonthlyTotal{
		ID:           uuid.New(),
		CreditCardID: uuid.New(),
		YearMonth:    "2024-01",
		TotalAmount:  150000,
		IsConfirmed:  false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

// CreateTestCardMonthlyTotalWithCreditCard creates a test CardMonthlyTotal with specified credit card
func CreateTestCardMonthlyTotalWithCreditCard(creditCardID uuid.UUID) *models.CardMonthlyTotal {
	return &models.CardMonthlyTotal{
		ID:           uuid.New(),
		CreditCardID: creditCardID,
		YearMonth:    "2024-01",
		TotalAmount:  150000,
		IsConfirmed:  false,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
}

// CreateTestMonthlyIncomeRecord creates a test MonthlyIncomeRecord with default values
func CreateTestMonthlyIncomeRecord(incomeSourceID uuid.UUID) *models.MonthlyIncomeRecord {
	return &models.MonthlyIncomeRecord{
		ID:             uuid.New(),
		IncomeSourceID: incomeSourceID,
		YearMonth:      "2024-01",
		ActualAmount:   280000, // 2800.00 in cents
		IsConfirmed:    true,
		Note:           "Test monthly income record",
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
}

// CreateTestAppSetting creates a test AppSetting with default values
func CreateTestAppSetting() *models.AppSetting {
	return &models.AppSetting{
		ID:        uuid.New(),
		Key:       "test_setting",
		Value:     "test_value",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// CreateTestAppSettingWithKeyValue creates a test AppSetting with specified key and value
func CreateTestAppSettingWithKeyValue(key, value string) *models.AppSetting {
	return &models.AppSetting{
		ID:        uuid.New(),
		Key:       key,
		Value:     value,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
