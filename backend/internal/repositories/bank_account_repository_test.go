package repositories

import (
	"database/sql"
	"testing"
	"time"

	"github.com/Soli0222/flow-sight/backend/internal/models"
	"github.com/Soli0222/flow-sight/backend/test/helpers"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBankAccountRepository_GetAll(t *testing.T) {
	db, mock := helpers.SetupMockDB(t)
	defer helpers.TeardownMockDB(db)

	repo := NewBankAccountRepository(db)

	tests := []struct {
		name          string
		setupMock     func(sqlmock.Sqlmock)
		expectedCount int
		expectedError bool
		errorType     string
	}{
		{
			name: "successful retrieval with multiple accounts",
			setupMock: func(mock sqlmock.Sqlmock) {
				accounts := []helpers.MockBankAccountData{
					{
						ID:        uuid.New(),
						Name:      "Main Account",
						Balance:   100000,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					{
						ID:        uuid.New(),
						Name:      "Savings Account",
						Balance:   50000,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
				}
				rows := sqlmock.NewRows([]string{"id", "name", "balance", "created_at", "updated_at"})
				for _, a := range accounts {
					rows.AddRow(a.ID, a.Name, a.Balance, a.CreatedAt, a.UpdatedAt)
				}

				mock.ExpectQuery(`SELECT id, name, balance, created_at, updated_at FROM bank_accounts\s+ORDER BY created_at DESC`).
					WillReturnRows(rows)
			},
			expectedCount: 2,
			expectedError: false,
		},
		{
			name: "successful retrieval with single account",
			setupMock: func(mock sqlmock.Sqlmock) {
				accounts := []helpers.MockBankAccountData{
					{ID: uuid.New(), Name: "Main Account", Balance: 100000, CreatedAt: time.Now(), UpdatedAt: time.Now()},
				}
				rows := sqlmock.NewRows([]string{"id", "name", "balance", "created_at", "updated_at"})
				for _, a := range accounts {
					rows.AddRow(a.ID, a.Name, a.Balance, a.CreatedAt, a.UpdatedAt)
				}
				mock.ExpectQuery(`SELECT id, name, balance, created_at, updated_at FROM bank_accounts\s+ORDER BY created_at DESC`).
					WillReturnRows(rows)
			},
			expectedCount: 1,
			expectedError: false,
		},
		{
			name: "no accounts found - empty result",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "name", "balance", "created_at", "updated_at"})
				mock.ExpectQuery(`SELECT id, name, balance, created_at, updated_at FROM bank_accounts\s+ORDER BY created_at DESC`).
					WillReturnRows(rows)
			},
			expectedCount: 0,
			expectedError: false,
		},
		{
			name: "database connection error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, name, balance, created_at, updated_at FROM bank_accounts\s+ORDER BY created_at DESC`).
					WillReturnError(sql.ErrConnDone)
			},
			expectedCount: 0,
			expectedError: true,
			errorType:     "connection",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt // capture range var

			tt.setupMock(mock)

			accounts, err := repo.GetAll()

			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorType == "connection" {
					assert.True(t, helpers.IsDatabaseError(err))
				}
			} else {
				assert.NoError(t, err)
				assert.Len(t, accounts, tt.expectedCount)

				// Validate that returned accounts have sane fields
				for _, account := range accounts {
					helpers.ValidateNonEmptyUUID(t, account.ID)
					assert.NotEmpty(t, account.Name)
					assert.GreaterOrEqual(t, account.Balance, int64(0))
				}
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestBankAccountRepository_GetByID(t *testing.T) {
	db, mock := helpers.SetupMockDB(t)
	defer helpers.TeardownMockDB(db)

	repo := NewBankAccountRepository(db)

	tests := []struct {
		name          string
		accountID     uuid.UUID
		setupMock     func(sqlmock.Sqlmock, uuid.UUID)
		expectedError bool
		errorType     string
	}{
		{
			name:      "successful retrieval",
			accountID: uuid.New(),
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				row := sqlmock.NewRows([]string{
					"id", "name", "balance", "created_at", "updated_at",
				}).AddRow(
					id, "Test Account", int64(100000),
					time.Now(), time.Now(),
				)

				mock.ExpectQuery(`SELECT id, name, balance, created_at, updated_at FROM bank_accounts WHERE id = \$1`).
					WithArgs(id).
					WillReturnRows(row)
			},
			expectedError: false,
		},
		{
			name:      "account not found",
			accountID: uuid.New(),
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				mock.ExpectQuery(`SELECT id, name, balance, created_at, updated_at FROM bank_accounts WHERE id = \$1`).
					WithArgs(id).
					WillReturnError(sql.ErrNoRows)
			},
			expectedError: true,
			errorType:     "notfound",
		},
		{
			name:      "database connection error",
			accountID: uuid.New(),
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				mock.ExpectQuery(`SELECT id, name, balance, created_at, updated_at FROM bank_accounts WHERE id = \$1`).
					WithArgs(id).
					WillReturnError(sql.ErrConnDone)
			},
			expectedError: true,
			errorType:     "connection",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(mock, tt.accountID)

			account, err := repo.GetByID(tt.accountID)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, account)

				switch tt.errorType {
				case "notfound":
					assert.True(t, helpers.IsNotFoundError(err))
				case "connection":
					assert.True(t, helpers.IsDatabaseError(err))
				}
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, account)
				assert.Equal(t, tt.accountID, account.ID)
				assert.NotEmpty(t, account.Name)
				assert.GreaterOrEqual(t, account.Balance, int64(0))
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestBankAccountRepository_Create(t *testing.T) {
	db, mock := helpers.SetupMockDB(t)
	defer helpers.TeardownMockDB(db)

	repo := NewBankAccountRepository(db)

	tests := []struct {
		name          string
		account       *models.BankAccount
		setupMock     func(sqlmock.Sqlmock)
		expectedError bool
		errorType     string
	}{
		{
			name:    "successful creation",
			account: helpers.CreateTestBankAccount(),
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO bank_accounts \(id, name, balance, created_at, updated_at\) VALUES \(\$1, \$2, \$3, \$4, \$5\)`).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: false,
		},
		{
			name:    "database error",
			account: helpers.CreateTestBankAccount(),
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO bank_accounts \(id, name, balance, created_at, updated_at\) VALUES \(\$1, \$2, \$3, \$4, \$5\)`).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(sql.ErrConnDone)
			},
			expectedError: true,
			errorType:     "connection",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(mock)

			err := repo.Create(tt.account)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorType == "connection" {
					assert.True(t, helpers.IsDatabaseError(err))
				}
			} else {
				assert.NoError(t, err)
				// Verify that the account has required fields set
				helpers.ValidateNonEmptyUUID(t, tt.account.ID)
				assert.NotEmpty(t, tt.account.Name)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestBankAccountRepository_Update(t *testing.T) {
	db, mock := helpers.SetupMockDB(t)
	defer helpers.TeardownMockDB(db)

	repo := NewBankAccountRepository(db)

	tests := []struct {
		name          string
		account       *models.BankAccount
		setupMock     func(sqlmock.Sqlmock)
		expectedError bool
		errorType     string
	}{
		{
			name:    "successful update",
			account: helpers.CreateTestBankAccount(),
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE bank_accounts SET name = \$2, balance = \$3, updated_at = \$4 WHERE id = \$1`).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedError: false,
		},
		{
			name:    "account not found",
			account: helpers.CreateTestBankAccount(),
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE bank_accounts SET name = \$2, balance = \$3, updated_at = \$4 WHERE id = \$1`).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectedError: false, // Repository doesn't return error for no rows affected
		},
		{
			name:    "database error",
			account: helpers.CreateTestBankAccount(),
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE bank_accounts SET name = \$2, balance = \$3, updated_at = \$4 WHERE id = \$1`).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(sql.ErrConnDone)
			},
			expectedError: true,
			errorType:     "connection",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(mock)

			err := repo.Update(tt.account)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorType == "connection" {
					assert.True(t, helpers.IsDatabaseError(err))
				}
			} else {
				assert.NoError(t, err)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestBankAccountRepository_Delete(t *testing.T) {
	db, mock := helpers.SetupMockDB(t)
	defer helpers.TeardownMockDB(db)

	repo := NewBankAccountRepository(db)

	tests := []struct {
		name          string
		accountID     uuid.UUID
		setupMock     func(sqlmock.Sqlmock, uuid.UUID)
		expectedError bool
		errorType     string
	}{
		{
			name:      "successful deletion",
			accountID: uuid.New(),
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				mock.ExpectExec(`DELETE FROM bank_accounts WHERE id = \$1`).
					WithArgs(id).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedError: false,
		},
		{
			name:      "account not found",
			accountID: uuid.New(),
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				mock.ExpectExec(`DELETE FROM bank_accounts WHERE id = \$1`).
					WithArgs(id).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectedError: false, // Repository doesn't return error for no rows affected
		},
		{
			name:      "database error",
			accountID: uuid.New(),
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				mock.ExpectExec(`DELETE FROM bank_accounts WHERE id = \$1`).
					WithArgs(id).
					WillReturnError(sql.ErrConnDone)
			},
			expectedError: true,
			errorType:     "connection",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(mock, tt.accountID)

			err := repo.Delete(tt.accountID)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorType == "connection" {
					assert.True(t, helpers.IsDatabaseError(err))
				}
			} else {
				assert.NoError(t, err)
			}

			require.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
