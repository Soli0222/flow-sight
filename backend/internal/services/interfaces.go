package services

import (
	"github.com/Soli0222/flow-sight/backend/internal/models"

	"github.com/google/uuid"
)

// BankAccountRepositoryInterface defines the interface for bank account repository
type BankAccountRepositoryInterface interface {
	GetAll(userID uuid.UUID) ([]models.BankAccount, error)
	GetByID(id uuid.UUID) (*models.BankAccount, error)
	Create(account *models.BankAccount) error
	Update(account *models.BankAccount) error
	Delete(id uuid.UUID) error
}

// CreditCardRepositoryInterface defines the interface for credit card repository
type CreditCardRepositoryInterface interface {
	GetAll(userID uuid.UUID) ([]models.CreditCard, error)
	GetByID(id uuid.UUID) (*models.CreditCard, error)
	Create(creditCard *models.CreditCard) error
	Update(creditCard *models.CreditCard) error
	Delete(id uuid.UUID) error
}

// RecurringPaymentRepositoryInterface defines the interface for recurring payment repository
type RecurringPaymentRepositoryInterface interface {
	GetAll(userID uuid.UUID) ([]models.RecurringPayment, error)
	GetByID(id uuid.UUID) (*models.RecurringPayment, error)
	Create(payment *models.RecurringPayment) error
	Update(payment *models.RecurringPayment) error
	Delete(id uuid.UUID) error
}

// UserRepositoryInterface defines the interface for user repository
type UserRepositoryInterface interface {
	GetByID(id uuid.UUID) (*models.User, error)
	GetByEmail(email string) (*models.User, error)
	GetByGoogleID(googleID string) (*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
}

// IncomeSourceRepositoryInterface defines the interface for income source repository
type IncomeSourceRepositoryInterface interface {
	GetAll(userID uuid.UUID) ([]models.IncomeSource, error)
	GetByID(id uuid.UUID) (*models.IncomeSource, error)
	GetActiveByUserID(userID uuid.UUID) ([]models.IncomeSource, error)
	Create(incomeSource *models.IncomeSource) error
	Update(incomeSource *models.IncomeSource) error
	Delete(id uuid.UUID) error
}

// MonthlyIncomeRepositoryInterface defines the interface for monthly income repository
type MonthlyIncomeRepositoryInterface interface {
	GetByIncomeSourceID(incomeSourceID uuid.UUID) ([]models.MonthlyIncomeRecord, error)
	GetByYearMonth(userID uuid.UUID, yearMonth string) ([]models.MonthlyIncomeRecord, error)
}

// AppSettingRepositoryInterface defines the interface for app setting repository
type AppSettingRepositoryInterface interface {
	GetByUserID(userID uuid.UUID) ([]models.AppSetting, error)
	GetByKey(userID uuid.UUID, key string) (*models.AppSetting, error)
	Upsert(setting *models.AppSetting) error
}

// CardMonthlyTotalRepositoryInterface defines the interface for card monthly total repository
type CardMonthlyTotalRepositoryInterface interface {
	GetByCreditCardID(creditCardID uuid.UUID) ([]models.CardMonthlyTotal, error)
	GetByYearMonth(yearMonth string) ([]models.CardMonthlyTotal, error)
	GetByID(id uuid.UUID) (*models.CardMonthlyTotal, error)
	Create(total *models.CardMonthlyTotal) error
	Update(total *models.CardMonthlyTotal) error
	Delete(id uuid.UUID) error
}
