package helpers

import (
	"flow-sight-backend/internal/models"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// AssertBankAccount validates BankAccount model properties
func AssertBankAccount(t *testing.T, expected, actual *models.BankAccount) {
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.UserID, actual.UserID)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.Balance, actual.Balance)
	assert.WithinDuration(t, expected.CreatedAt, actual.CreatedAt, 0)
	assert.WithinDuration(t, expected.UpdatedAt, actual.UpdatedAt, 0)
}

// AssertCreditCard validates CreditCard model properties
func AssertCreditCard(t *testing.T, expected, actual *models.CreditCard) {
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.UserID, actual.UserID)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.ClosingDay, actual.ClosingDay)
	assert.Equal(t, expected.PaymentDay, actual.PaymentDay)
	assert.Equal(t, expected.BankAccount, actual.BankAccount)
	assert.WithinDuration(t, expected.CreatedAt, actual.CreatedAt, 0)
	assert.WithinDuration(t, expected.UpdatedAt, actual.UpdatedAt, 0)
}

// AssertIncomeSource validates IncomeSource model properties
func AssertIncomeSource(t *testing.T, expected, actual *models.IncomeSource) {
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.UserID, actual.UserID)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.IncomeType, actual.IncomeType)
	assert.Equal(t, expected.BaseAmount, actual.BaseAmount)
	assert.Equal(t, expected.BankAccount, actual.BankAccount)
	assert.Equal(t, expected.PaymentDay, actual.PaymentDay)
	assert.Equal(t, expected.ScheduledDate, actual.ScheduledDate)
	assert.Equal(t, expected.ScheduledYearMonth, actual.ScheduledYearMonth)
	assert.Equal(t, expected.IsActive, actual.IsActive)
	assert.WithinDuration(t, expected.CreatedAt, actual.CreatedAt, 0)
	assert.WithinDuration(t, expected.UpdatedAt, actual.UpdatedAt, 0)
}

// AssertRecurringPayment validates RecurringPayment model properties
func AssertRecurringPayment(t *testing.T, expected, actual *models.RecurringPayment) {
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.UserID, actual.UserID)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.Amount, actual.Amount)
	assert.Equal(t, expected.PaymentDay, actual.PaymentDay)
	assert.Equal(t, expected.StartYearMonth, actual.StartYearMonth)
	assert.Equal(t, expected.TotalPayments, actual.TotalPayments)
	assert.Equal(t, expected.RemainingPayments, actual.RemainingPayments)
	assert.Equal(t, expected.BankAccount, actual.BankAccount)
	assert.Equal(t, expected.IsActive, actual.IsActive)
	assert.Equal(t, expected.Note, actual.Note)
	assert.WithinDuration(t, expected.CreatedAt, actual.CreatedAt, 0)
	assert.WithinDuration(t, expected.UpdatedAt, actual.UpdatedAt, 0)
}

// AssertMonthlyIncomeRecord validates MonthlyIncomeRecord model properties
func AssertMonthlyIncomeRecord(t *testing.T, expected, actual *models.MonthlyIncomeRecord) {
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.IncomeSourceID, actual.IncomeSourceID)
	assert.Equal(t, expected.YearMonth, actual.YearMonth)
	assert.Equal(t, expected.ActualAmount, actual.ActualAmount)
	assert.Equal(t, expected.IsConfirmed, actual.IsConfirmed)
	assert.Equal(t, expected.Note, actual.Note)
	assert.WithinDuration(t, expected.CreatedAt, actual.CreatedAt, 0)
	assert.WithinDuration(t, expected.UpdatedAt, actual.UpdatedAt, 0)
}

// AssertCardMonthlyTotal validates CardMonthlyTotal model properties
func AssertCardMonthlyTotal(t *testing.T, expected, actual *models.CardMonthlyTotal) {
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.CreditCardID, actual.CreditCardID)
	assert.Equal(t, expected.YearMonth, actual.YearMonth)
	assert.Equal(t, expected.TotalAmount, actual.TotalAmount)
	assert.Equal(t, expected.IsConfirmed, actual.IsConfirmed)
	assert.WithinDuration(t, expected.CreatedAt, actual.CreatedAt, 0)
	assert.WithinDuration(t, expected.UpdatedAt, actual.UpdatedAt, 0)
}

// AssertUser validates User model properties
func AssertUser(t *testing.T, expected, actual *models.User) {
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.Email, actual.Email)
	assert.Equal(t, expected.Name, actual.Name)
	assert.Equal(t, expected.Picture, actual.Picture)
	assert.Equal(t, expected.GoogleID, actual.GoogleID)
	assert.WithinDuration(t, expected.CreatedAt, actual.CreatedAt, 0)
	assert.WithinDuration(t, expected.UpdatedAt, actual.UpdatedAt, 0)
}

// AssertAppSetting validates AppSetting model properties
func AssertAppSetting(t *testing.T, expected, actual *models.AppSetting) {
	assert.Equal(t, expected.ID, actual.ID)
	assert.Equal(t, expected.UserID, actual.UserID)
	assert.Equal(t, expected.Key, actual.Key)
	assert.Equal(t, expected.Value, actual.Value)
	assert.WithinDuration(t, expected.CreatedAt, actual.CreatedAt, 0)
	assert.WithinDuration(t, expected.UpdatedAt, actual.UpdatedAt, 0)
}

// ValidateUUID validates that a string is a valid UUID
func ValidateUUID(t *testing.T, uuidStr string) {
	_, err := uuid.Parse(uuidStr)
	assert.NoError(t, err, "should be a valid UUID")
}

// ValidateNonEmptyUUID validates that a UUID is not empty
func ValidateNonEmptyUUID(t *testing.T, id uuid.UUID) {
	assert.NotEqual(t, uuid.Nil, id, "UUID should not be empty")
}
