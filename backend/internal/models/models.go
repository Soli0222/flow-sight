package models

import (
	"time"

	"github.com/google/uuid"
)

// CreditCard represents a credit card
type CreditCard struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Name        string    `json:"name" db:"name"`
	ClosingDay  *int      `json:"closing_day,omitempty" db:"closing_day"` // Closing day of the month
	PaymentDay  int       `json:"payment_day" db:"payment_day"`
	BankAccount uuid.UUID `json:"bank_account" db:"bank_account"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
}

// BankAccount represents a user's bank account
type BankAccount struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Balance   int64     `json:"balance" db:"balance"` // Amount in cents
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// IncomeSource represents a source of income
type IncomeSource struct {
	ID                 uuid.UUID `json:"id" db:"id"`
	Name               string    `json:"name" db:"name"`
	IncomeType         string    `json:"income_type" db:"income_type"` // "monthly_fixed" or "one_time"
	BaseAmount         int64     `json:"base_amount" db:"base_amount"` // Amount in cents
	BankAccount        uuid.UUID `json:"bank_account" db:"bank_account"`
	PaymentDay         *int      `json:"payment_day,omitempty" db:"payment_day"`                   // For monthly_fixed income (1-31)
	ScheduledDate      *string   `json:"scheduled_date,omitempty" db:"scheduled_date"`             // For one_time income (YYYY-MM-DD format)
	ScheduledYearMonth *string   `json:"scheduled_year_month,omitempty" db:"scheduled_year_month"` // For one-time income (backward compatibility)
	IsActive           bool      `json:"is_active" db:"is_active"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	UpdatedAt          time.Time `json:"updated_at" db:"updated_at"`
}

// MonthlyIncomeRecord represents actual income for a specific month
type MonthlyIncomeRecord struct {
	ID             uuid.UUID `json:"id" db:"id"`
	IncomeSourceID uuid.UUID `json:"income_source_id" db:"income_source_id"`
	YearMonth      string    `json:"year_month" db:"year_month"`       // Format: "2024-01"
	ActualAmount   int64     `json:"actual_amount" db:"actual_amount"` // Amount in cents
	IsConfirmed    bool      `json:"is_confirmed" db:"is_confirmed"`
	Note           string    `json:"note" db:"note"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
	UpdatedAt      time.Time `json:"updated_at" db:"updated_at"`
}

// RecurringPayment represents a fixed recurring payment
type RecurringPayment struct {
	ID                uuid.UUID `json:"id" db:"id"`
	Name              string    `json:"name" db:"name"`
	Amount            int64     `json:"amount" db:"amount"` // Amount in cents
	PaymentDay        int       `json:"payment_day" db:"payment_day"`
	StartYearMonth    string    `json:"start_year_month" db:"start_year_month"`       // Format: "2024-01"
	TotalPayments     *int      `json:"total_payments,omitempty" db:"total_payments"` // For loans
	RemainingPayments *int      `json:"remaining_payments,omitempty" db:"remaining_payments"`
	BankAccount       uuid.UUID `json:"bank_account" db:"bank_account"`
	IsActive          bool      `json:"is_active" db:"is_active"`
	Note              string    `json:"note" db:"note"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	UpdatedAt         time.Time `json:"updated_at" db:"updated_at"`
}

// CardMonthlyTotal represents monthly credit card usage
type CardMonthlyTotal struct {
	ID           uuid.UUID `json:"id" db:"id"`
	CreditCardID uuid.UUID `json:"credit_card_id" db:"credit_card_id"`
	YearMonth    string    `json:"year_month" db:"year_month"`     // Format: "2024-01"
	TotalAmount  int64     `json:"total_amount" db:"total_amount"` // Amount in cents
	IsConfirmed  bool      `json:"is_confirmed" db:"is_confirmed"`
	CreatedAt    time.Time `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time `json:"updated_at" db:"updated_at"`
}

// AppSetting represents application settings
type AppSetting struct {
	ID        uuid.UUID `json:"id" db:"id"`
	Key       string    `json:"key" db:"key"`
	Value     string    `json:"value" db:"value"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// CashflowProjection represents a cashflow projection result
type CashflowProjection struct {
	Date    string                     `json:"date"`
	Income  int64                      `json:"income"`
	Expense int64                      `json:"expense"`
	Balance int64                      `json:"balance"`
	Details []CashflowProjectionDetail `json:"details"`
}

// CashflowProjectionDetail represents details of a cashflow projection
type CashflowProjectionDetail struct {
	Type        string `json:"type"` // "income", "recurring_payment", "card_payment"
	Description string `json:"description"`
	Amount      int64  `json:"amount"`
}

// DashboardSummary represents dashboard summary data
type DashboardSummary struct {
	TotalBalance     int64                `json:"total_balance"`
	MonthlyIncome    int64                `json:"monthly_income"`
	MonthlyExpense   int64                `json:"monthly_expense"`
	TotalAssets      int                  `json:"total_assets"`
	RecentActivities []CashflowProjection `json:"recent_activities"`
}
