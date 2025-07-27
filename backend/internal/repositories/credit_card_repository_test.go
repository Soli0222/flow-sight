package repositories

import (
	"database/sql"
	"github.com/Soli0222/flow-sight/backend/internal/models"
	"github.com/Soli0222/flow-sight/backend/test/helpers"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCreditCardRepository_GetAll(t *testing.T) {
	db, mock := helpers.SetupMockDB(t)
	defer helpers.TeardownMockDB(db)

	repo := NewCreditCardRepository(db)
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
				closingDay := 25
				rows := sqlmock.NewRows([]string{
					"id", "user_id", "name", "closing_day", "payment_day", "bank_account", "created_at", "updated_at",
				}).
					AddRow(
						uuid.New(), userID, "Main Credit Card", &closingDay, 10, bankAccountID,
						time.Now(), time.Now(),
					).
					AddRow(
						uuid.New(), userID, "Sub Credit Card", &closingDay, 15, bankAccountID,
						time.Now(), time.Now(),
					)

				mock.ExpectQuery(`SELECT id, user_id, name, closing_day, payment_day, bank_account, created_at, updated_at FROM credit_cards WHERE user_id = \$1 ORDER BY created_at DESC`).
					WithArgs(userID).
					WillReturnRows(rows)
			},
			expectedCount: 2,
			expectedError: false,
		},
		{
			name:   "no credit cards found",
			userID: userID,
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "user_id", "name", "closing_day", "payment_day", "bank_account", "created_at", "updated_at",
				})

				mock.ExpectQuery(`SELECT id, user_id, name, closing_day, payment_day, bank_account, created_at, updated_at FROM credit_cards WHERE user_id = \$1 ORDER BY created_at DESC`).
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
				mock.ExpectQuery(`SELECT id, user_id, name, closing_day, payment_day, bank_account, created_at, updated_at FROM credit_cards WHERE user_id = \$1 ORDER BY created_at DESC`).
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

			creditCards, err := repo.GetAll(tt.userID)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, creditCards, tt.expectedCount)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestCreditCardRepository_GetByID(t *testing.T) {
	db, mock := helpers.SetupMockDB(t)
	defer helpers.TeardownMockDB(db)

	repo := NewCreditCardRepository(db)
	creditCardID := uuid.New()
	userID := uuid.New()
	bankAccountID := uuid.New()

	tests := []struct {
		name          string
		creditCardID  uuid.UUID
		setupMock     func(sqlmock.Sqlmock)
		expectedFound bool
		expectedError bool
	}{
		{
			name:         "credit card found",
			creditCardID: creditCardID,
			setupMock: func(mock sqlmock.Sqlmock) {
				closingDay := 25
				rows := sqlmock.NewRows([]string{
					"id", "user_id", "name", "closing_day", "payment_day", "bank_account", "created_at", "updated_at",
				}).
					AddRow(
						creditCardID, userID, "Main Credit Card", &closingDay, 10, bankAccountID,
						time.Now(), time.Now(),
					)

				mock.ExpectQuery(`SELECT id, user_id, name, closing_day, payment_day, bank_account, created_at, updated_at FROM credit_cards WHERE id = \$1`).
					WithArgs(creditCardID).
					WillReturnRows(rows)
			},
			expectedFound: true,
			expectedError: false,
		},
		{
			name:         "credit card not found",
			creditCardID: creditCardID,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, user_id, name, closing_day, payment_day, bank_account, created_at, updated_at FROM credit_cards WHERE id = \$1`).
					WithArgs(creditCardID).
					WillReturnError(sql.ErrNoRows)
			},
			expectedFound: false,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(mock)

			creditCard, err := repo.GetByID(tt.creditCardID)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, creditCard)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, creditCard)
				assert.Equal(t, tt.creditCardID, creditCard.ID)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestCreditCardRepository_Create(t *testing.T) {
	db, mock := helpers.SetupMockDB(t)
	defer helpers.TeardownMockDB(db)

	repo := NewCreditCardRepository(db)
	userID := uuid.New()
	bankAccountID := uuid.New()
	creditCard := helpers.CreateTestCreditCard(userID, bankAccountID)

	tests := []struct {
		name          string
		creditCard    *models.CreditCard
		setupMock     func(sqlmock.Sqlmock)
		expectedError bool
	}{
		{
			name:       "successful creation",
			creditCard: creditCard,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO credit_cards \(id, user_id, name, closing_day, payment_day, bank_account, created_at, updated_at\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8\)`).
					WithArgs(creditCard.ID, creditCard.UserID, creditCard.Name, creditCard.ClosingDay, creditCard.PaymentDay, creditCard.BankAccount, sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: false,
		},
		{
			name:       "database error",
			creditCard: creditCard,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO credit_cards \(id, user_id, name, closing_day, payment_day, bank_account, created_at, updated_at\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8\)`).
					WithArgs(creditCard.ID, creditCard.UserID, creditCard.Name, creditCard.ClosingDay, creditCard.PaymentDay, creditCard.BankAccount, sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(sql.ErrConnDone)
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(mock)

			err := repo.Create(tt.creditCard)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestCreditCardRepository_Update(t *testing.T) {
	db, mock := helpers.SetupMockDB(t)
	defer helpers.TeardownMockDB(db)

	repo := NewCreditCardRepository(db)
	userID := uuid.New()
	bankAccountID := uuid.New()
	creditCard := helpers.CreateTestCreditCard(userID, bankAccountID)
	creditCard.Name = "Updated Credit Card Name"

	tests := []struct {
		name          string
		creditCard    *models.CreditCard
		setupMock     func(sqlmock.Sqlmock)
		expectedError bool
	}{
		{
			name:       "successful update",
			creditCard: creditCard,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE credit_cards SET name = \$2, closing_day = \$3, payment_day = \$4, bank_account = \$5, updated_at = \$6 WHERE id = \$1`).
					WithArgs(creditCard.ID, creditCard.Name, creditCard.ClosingDay, creditCard.PaymentDay, creditCard.BankAccount, sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: false,
		},
		{
			name:       "database error",
			creditCard: creditCard,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE credit_cards SET name = \$2, closing_day = \$3, payment_day = \$4, bank_account = \$5, updated_at = \$6 WHERE id = \$1`).
					WithArgs(creditCard.ID, creditCard.Name, creditCard.ClosingDay, creditCard.PaymentDay, creditCard.BankAccount, sqlmock.AnyArg()).
					WillReturnError(sql.ErrConnDone)
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(mock)

			err := repo.Update(tt.creditCard)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestCreditCardRepository_Delete(t *testing.T) {
	db, mock := helpers.SetupMockDB(t)
	defer helpers.TeardownMockDB(db)

	repo := NewCreditCardRepository(db)
	creditCardID := uuid.New()

	tests := []struct {
		name          string
		creditCardID  uuid.UUID
		setupMock     func(sqlmock.Sqlmock)
		expectedError bool
	}{
		{
			name:         "successful deletion",
			creditCardID: creditCardID,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM credit_cards WHERE id = \$1`).
					WithArgs(creditCardID).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: false,
		},
		{
			name:         "database error",
			creditCardID: creditCardID,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM credit_cards WHERE id = \$1`).
					WithArgs(creditCardID).
					WillReturnError(sql.ErrConnDone)
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(mock)

			err := repo.Delete(tt.creditCardID)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
