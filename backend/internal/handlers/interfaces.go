package handlers

import (
	"github.com/Soli0222/flow-sight/backend/internal/models"
	"github.com/Soli0222/flow-sight/backend/internal/services"

	"github.com/google/uuid"
)

// BankAccountServiceInterface defines the interface for bank account service
type BankAccountServiceInterface interface {
	GetBankAccounts(userID uuid.UUID) ([]models.BankAccount, error)
	GetBankAccount(id uuid.UUID) (*models.BankAccount, error)
	CreateBankAccount(account *models.BankAccount) error
	UpdateBankAccount(account *models.BankAccount) error
	DeleteBankAccount(id uuid.UUID) error
}

// CreditCardServiceInterface defines the interface for credit card service
type CreditCardServiceInterface interface {
	GetCreditCards(userID uuid.UUID) ([]models.CreditCard, error)
	GetCreditCard(id uuid.UUID) (*models.CreditCard, error)
	CreateCreditCard(creditCard *models.CreditCard) error
	UpdateCreditCard(creditCard *models.CreditCard) error
	DeleteCreditCard(id uuid.UUID) error
}

// AuthServiceInterface defines the interface for auth service
type AuthServiceInterface interface {
	GetGoogleAuthURL(state string) string
	HandleGoogleCallback(code string) (*models.User, string, error)
	GenerateJWT(user *models.User) (string, error)
	ValidateJWT(tokenString string) (*services.Claims, error)
	GetUserByID(userID string) (*models.User, error)
}

// RecurringPaymentServiceInterface defines the interface for recurring payment service
type RecurringPaymentServiceInterface interface {
	GetRecurringPayments(userID uuid.UUID) ([]models.RecurringPayment, error)
	GetRecurringPayment(id uuid.UUID) (*models.RecurringPayment, error)
	CreateRecurringPayment(payment *models.RecurringPayment) error
	UpdateRecurringPayment(payment *models.RecurringPayment) error
	DeleteRecurringPayment(id uuid.UUID) error
}

// IncomeServiceInterface defines the interface for income service
type IncomeServiceInterface interface {
	GetIncomeSources(userID uuid.UUID) ([]models.IncomeSource, error)
	GetIncomeSource(id uuid.UUID) (*models.IncomeSource, error)
	CreateIncomeSource(source *models.IncomeSource) error
	UpdateIncomeSource(source *models.IncomeSource) error
	DeleteIncomeSource(id uuid.UUID) error
	GetMonthlyIncomeRecords(incomeSourceID uuid.UUID) ([]models.MonthlyIncomeRecord, error)
	GetMonthlyIncomeRecord(id uuid.UUID) (*models.MonthlyIncomeRecord, error)
	CreateMonthlyIncomeRecord(record *models.MonthlyIncomeRecord) error
	UpdateMonthlyIncomeRecord(record *models.MonthlyIncomeRecord) error
	DeleteMonthlyIncomeRecord(id uuid.UUID) error
}
