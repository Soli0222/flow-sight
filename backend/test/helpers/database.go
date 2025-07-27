package helpers

import (
	"database/sql"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

// SetupMockDB creates a mock database connection for testing
func SetupMockDB(t *testing.T) (*sql.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	require.NoError(t, err)
	return db, mock
}

// TeardownMockDB closes the mock database connection
func TeardownMockDB(db *sql.DB) {
	db.Close()
}

// ExpectBankAccountRows creates expected rows for bank account queries
func ExpectBankAccountRows(mock sqlmock.Sqlmock, accounts []MockBankAccountData) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{
		"id", "user_id", "name", "balance", "created_at", "updated_at",
	})

	for _, account := range accounts {
		rows.AddRow(
			account.ID,
			account.UserID,
			account.Name,
			account.Balance,
			account.CreatedAt,
			account.UpdatedAt,
		)
	}
	return rows
}

// ExpectCreditCardRows creates expected rows for credit card queries
func ExpectCreditCardRows(mock sqlmock.Sqlmock, cards []MockCreditCardData) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{
		"id", "user_id", "name", "closing_day", "payment_day", "bank_account", "created_at", "updated_at",
	})

	for _, card := range cards {
		rows.AddRow(
			card.ID,
			card.UserID,
			card.Name,
			card.ClosingDay,
			card.PaymentDay,
			card.BankAccount,
			card.CreatedAt,
			card.UpdatedAt,
		)
	}
	return rows
}

// ExpectIncomeSourceRows creates expected rows for income source queries
func ExpectIncomeSourceRows(mock sqlmock.Sqlmock, sources []MockIncomeSourceData) *sqlmock.Rows {
	rows := sqlmock.NewRows([]string{
		"id", "user_id", "name", "income_type", "base_amount", "bank_account",
		"payment_day", "scheduled_date", "scheduled_year_month", "is_active", "created_at", "updated_at",
	})

	for _, source := range sources {
		rows.AddRow(
			source.ID,
			source.UserID,
			source.Name,
			source.IncomeType,
			source.BaseAmount,
			source.BankAccount,
			source.PaymentDay,
			source.ScheduledDate,
			source.ScheduledYearMonth,
			source.IsActive,
			source.CreatedAt,
			source.UpdatedAt,
		)
	}
	return rows
}

// MockBankAccountData represents test data for bank account
type MockBankAccountData struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	Name      string
	Balance   int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

// MockCreditCardData represents test data for credit card
type MockCreditCardData struct {
	ID          uuid.UUID
	UserID      uuid.UUID
	Name        string
	ClosingDay  *int
	PaymentDay  int
	BankAccount uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// MockIncomeSourceData represents test data for income source
type MockIncomeSourceData struct {
	ID                 uuid.UUID
	UserID             uuid.UUID
	Name               string
	IncomeType         string
	BaseAmount         int64
	BankAccount        uuid.UUID
	PaymentDay         *int
	ScheduledDate      *string
	ScheduledYearMonth *string
	IsActive           bool
	CreatedAt          time.Time
	UpdatedAt          time.Time
}
