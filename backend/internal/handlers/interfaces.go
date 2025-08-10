package handlers

import (
	"github.com/Soli0222/flow-sight/backend/internal/models"

	"github.com/google/uuid"
)

// BankAccountServiceInterface defines the interface for bank account service
type BankAccountServiceInterface interface {
	GetBankAccounts() ([]models.BankAccount, error)
	GetBankAccount(id uuid.UUID) (*models.BankAccount, error)
	CreateBankAccount(account *models.BankAccount) error
	UpdateBankAccount(account *models.BankAccount) error
	DeleteBankAccount(id uuid.UUID) error
}

// CreditCardServiceInterface defines the interface for credit card service
type CreditCardServiceInterface interface {
	GetCreditCards() ([]models.CreditCard, error)
	GetCreditCard(id uuid.UUID) (*models.CreditCard, error)
	CreateCreditCard(creditCard *models.CreditCard) error
	UpdateCreditCard(creditCard *models.CreditCard) error
	DeleteCreditCard(id uuid.UUID) error
}

// RecurringPaymentServiceInterface defines the interface for recurring payment service
type RecurringPaymentServiceInterface interface {
	GetRecurringPayments() ([]models.RecurringPayment, error)
	GetRecurringPayment(id uuid.UUID) (*models.RecurringPayment, error)
	CreateRecurringPayment(payment *models.RecurringPayment) error
	UpdateRecurringPayment(payment *models.RecurringPayment) error
	DeleteRecurringPayment(id uuid.UUID) error
}

// IncomeServiceInterface defines the interface for income service
type IncomeServiceInterface interface {
	GetIncomeSources() ([]models.IncomeSource, error)
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
