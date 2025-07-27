package repositories

import (
	"database/sql"
	"flow-sight-backend/internal/models"
	"flow-sight-backend/test/helpers"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestRecurringPaymentRepository_GetAll(t *testing.T) {
	db, mock := helpers.SetupMockDB(t)
	defer helpers.TeardownMockDB(db)

	repo := NewRecurringPaymentRepository(db)
	userID := uuid.New()

	tests := []struct {
		name          string
		userID        uuid.UUID
		setupMock     func(sqlmock.Sqlmock)
		expectedCount int
		expectedError bool
	}{
		{
			name:   "successful retrieval",
			userID: userID,
			setupMock: func(mock sqlmock.Sqlmock) {
				bankAccountID := uuid.New()
				totalPayments := 12
				remainingPayments := 8
				rows := sqlmock.NewRows([]string{
					"id", "user_id", "name", "amount", "payment_day", "start_year_month",
					"total_payments", "remaining_payments", "bank_account", "is_active",
					"note", "created_at", "updated_at",
				}).
					AddRow(
						uuid.New(), userID, "Monthly Rent", int64(120000), 1, "2024-01",
						&totalPayments, &remainingPayments, bankAccountID, true,
						"Monthly rent payment", time.Now(), time.Now(),
					).
					AddRow(
						uuid.New(), userID, "Insurance", int64(8000), 15, "2024-01",
						nil, nil, bankAccountID, true,
						"Monthly insurance", time.Now(), time.Now(),
					)

				mock.ExpectQuery(`SELECT id, user_id, name, amount, payment_day, start_year_month, total_payments, remaining_payments, bank_account, is_active, note, created_at, updated_at FROM recurring_payments WHERE user_id = \$1 ORDER BY created_at DESC`).
					WithArgs(userID).
					WillReturnRows(rows)
			},
			expectedCount: 2,
			expectedError: false,
		},
		{
			name:   "no payments found",
			userID: userID,
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "user_id", "name", "amount", "payment_day", "start_year_month",
					"total_payments", "remaining_payments", "bank_account", "is_active",
					"note", "created_at", "updated_at",
				})

				mock.ExpectQuery(`SELECT id, user_id, name, amount, payment_day, start_year_month, total_payments, remaining_payments, bank_account, is_active, note, created_at, updated_at FROM recurring_payments WHERE user_id = \$1 ORDER BY created_at DESC`).
					WithArgs(userID).
					WillReturnRows(rows)
			},
			expectedCount: 0,
			expectedError: false,
		},
		{
			name:   "database error",
			userID: userID,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, user_id, name, amount, payment_day, start_year_month, total_payments, remaining_payments, bank_account, is_active, note, created_at, updated_at FROM recurring_payments WHERE user_id = \$1 ORDER BY created_at DESC`).
					WithArgs(userID).
					WillReturnError(sql.ErrConnDone)
			},
			expectedCount: 0,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(mock)

			payments, err := repo.GetAll(tt.userID)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, payments, tt.expectedCount)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRecurringPaymentRepository_GetByID(t *testing.T) {
	db, mock := helpers.SetupMockDB(t)
	defer helpers.TeardownMockDB(db)

	repo := NewRecurringPaymentRepository(db)
	paymentID := uuid.New()
	userID := uuid.New()
	bankAccountID := uuid.New()

	tests := []struct {
		name          string
		paymentID     uuid.UUID
		setupMock     func(sqlmock.Sqlmock)
		expectedFound bool
		expectedError bool
	}{
		{
			name:      "payment found",
			paymentID: paymentID,
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "user_id", "name", "amount", "payment_day", "start_year_month",
					"total_payments", "remaining_payments", "bank_account", "is_active",
					"note", "created_at", "updated_at",
				}).
					AddRow(
						paymentID, userID, "Monthly Rent", int64(120000), 1, "2024-01",
						nil, nil, bankAccountID, true,
						"Monthly rent payment", time.Now(), time.Now(),
					)

				mock.ExpectQuery(`SELECT id, user_id, name, amount, payment_day, start_year_month, total_payments, remaining_payments, bank_account, is_active, note, created_at, updated_at FROM recurring_payments WHERE id = \$1`).
					WithArgs(paymentID).
					WillReturnRows(rows)
			},
			expectedFound: true,
			expectedError: false,
		},
		{
			name:      "payment not found",
			paymentID: paymentID,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, user_id, name, amount, payment_day, start_year_month, total_payments, remaining_payments, bank_account, is_active, note, created_at, updated_at FROM recurring_payments WHERE id = \$1`).
					WithArgs(paymentID).
					WillReturnError(sql.ErrNoRows)
			},
			expectedFound: false,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(mock)

			payment, err := repo.GetByID(tt.paymentID)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, payment)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, payment)
				assert.Equal(t, tt.paymentID, payment.ID)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestRecurringPaymentRepository_Create(t *testing.T) {
	db, mock := helpers.SetupMockDB(t)
	defer helpers.TeardownMockDB(db)

	repo := NewRecurringPaymentRepository(db)
	userID := uuid.New()
	bankAccountID := uuid.New()
	payment := helpers.CreateTestRecurringPayment(userID, bankAccountID)

	tests := []struct {
		name          string
		payment       *models.RecurringPayment
		setupMock     func(sqlmock.Sqlmock)
		expectedError bool
	}{
		{
			name:    "successful creation",
			payment: payment,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO recurring_payments \(id, user_id, name, amount, payment_day, start_year_month, total_payments, remaining_payments, bank_account, is_active, note, created_at, updated_at\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8, \$9, \$10, \$11, \$12, \$13\)`).
					WithArgs(payment.ID, payment.UserID, payment.Name, payment.Amount, payment.PaymentDay, payment.StartYearMonth, payment.TotalPayments, payment.RemainingPayments, payment.BankAccount, payment.IsActive, payment.Note, sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: false,
		},
		{
			name:    "database error",
			payment: payment,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO recurring_payments \(id, user_id, name, amount, payment_day, start_year_month, total_payments, remaining_payments, bank_account, is_active, note, created_at, updated_at\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8, \$9, \$10, \$11, \$12, \$13\)`).
					WithArgs(payment.ID, payment.UserID, payment.Name, payment.Amount, payment.PaymentDay, payment.StartYearMonth, payment.TotalPayments, payment.RemainingPayments, payment.BankAccount, payment.IsActive, payment.Note, sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(sql.ErrConnDone)
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(mock)

			err := repo.Create(tt.payment)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
