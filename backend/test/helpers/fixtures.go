package helpers

import (
	"flow-sight-backend/internal/models"
	"time"

	"github.com/google/uuid"
)

// CreateTestUser creates a test user with default values
func CreateTestUser() *models.User {
	return &models.User{
		ID:        uuid.New(),
		Email:     "test@example.com",
		Name:      "Test User",
		GoogleID:  "test-google-id",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// CreateTestBankAccount creates a test bank account with default values
func CreateTestBankAccount(userID uuid.UUID) *models.BankAccount {
	return &models.BankAccount{
		ID:        uuid.New(),
		UserID:    userID,
		Name:      "Test Bank Account",
		Balance:   100000, // 1000.00 in cents
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// CreateTestCreditCard creates a test credit card with default values
func CreateTestCreditCard(userID, bankAccountID uuid.UUID) *models.CreditCard {
	closingDay := 25
	return &models.CreditCard{
		ID:          uuid.New(),
		UserID:      userID,
		Name:        "Test Credit Card",
		ClosingDay:  &closingDay,
		PaymentDay:  10,
		BankAccount: bankAccountID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}
}

// CreateTestIncomeSource creates a test income source with default values
func CreateTestIncomeSource(userID, bankAccountID uuid.UUID) *models.IncomeSource {
	paymentDay := 25
	return &models.IncomeSource{
		ID:          uuid.New(),
		UserID:      userID,
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
func CreateTestRecurringPayment(userID, bankAccountID uuid.UUID) *models.RecurringPayment {
	return &models.RecurringPayment{
		ID:             uuid.New(),
		UserID:         userID,
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
func CreateTestAppSetting(userID uuid.UUID) *models.AppSetting {
	return &models.AppSetting{
		ID:        uuid.New(),
		UserID:    userID,
		Key:       "test_setting",
		Value:     "test_value",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// CreateTestAppSettingWithKeyValue creates a test AppSetting with specified key and value
func CreateTestAppSettingWithKeyValue(userID uuid.UUID, key, value string) *models.AppSetting {
	return &models.AppSetting{
		ID:        uuid.New(),
		UserID:    userID,
		Key:       key,
		Value:     value,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// CreateTestUserWithEmail creates a test user with specified email
func CreateTestUserWithEmail(email string) *models.User {
	return &models.User{
		ID:        uuid.New(),
		Email:     email,
		Name:      "Test User",
		Picture:   "https://example.com/picture.jpg",
		GoogleID:  "test-google-id",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}

// CreateTestIncomeSourceOneTime creates a test one-time income source
func CreateTestIncomeSourceOneTime(userID, bankAccountID uuid.UUID) *models.IncomeSource {
	scheduledDate := "2024-12-25"
	return &models.IncomeSource{
		ID:            uuid.New(),
		UserID:        userID,
		Name:          "Test One-time Income",
		IncomeType:    "one_time",
		BaseAmount:    500000, // 5000.00 in cents
		BankAccount:   bankAccountID,
		ScheduledDate: &scheduledDate,
		IsActive:      true,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}
}

// CreateTestRecurringPaymentWithLoan creates a test recurring payment with loan details
func CreateTestRecurringPaymentWithLoan(userID, bankAccountID uuid.UUID) *models.RecurringPayment {
	totalPayments := 36
	remainingPayments := 24
	return &models.RecurringPayment{
		ID:                uuid.New(),
		UserID:            userID,
		Name:              "Test Loan Payment",
		Amount:            120000, // 1200.00 in cents
		PaymentDay:        15,
		StartYearMonth:    "2024-01",
		TotalPayments:     &totalPayments,
		RemainingPayments: &remainingPayments,
		BankAccount:       bankAccountID,
		IsActive:          true,
		Note:              "Test loan payment with installments",
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
}
