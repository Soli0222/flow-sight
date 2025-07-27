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
	"github.com/stretchr/testify/require"
)

func TestBankAccountRepository_GetAll(t *testing.T) {
	db, mock := helpers.SetupMockDB(t)
	defer helpers.TeardownMockDB(db)

	repo := NewBankAccountRepository(db)
	userID := uuid.New()

	tests := []struct {
		name          string
		userID        uuid.UUID
		setupMock     func(sqlmock.Sqlmock)
		expectedCount int
		expectedError bool
		errorType     string
	}{
		{
			name:   "successful retrieval with multiple accounts",
			userID: userID,
			setupMock: func(mock sqlmock.Sqlmock) {
				accounts := []helpers.MockBankAccountData{
					{
						ID:        uuid.New(),
						UserID:    userID,
						Name:      "Main Account",
						Balance:   100000,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
					{
						ID:        uuid.New(),
						UserID:    userID,
						Name:      "Savings Account",
						Balance:   50000,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
				}
				rows := helpers.ExpectBankAccountRows(mock, accounts)

				mock.ExpectQuery(`SELECT id, user_id, name, balance, created_at, updated_at FROM bank_accounts WHERE user_id = \$1 ORDER BY created_at DESC`).
					WithArgs(userID).
					WillReturnRows(rows)
			},
			expectedCount: 2,
			expectedError: false,
		},
		{
			name:   "successful retrieval with single account",
			userID: userID,
			setupMock: func(mock sqlmock.Sqlmock) {
				accounts := []helpers.MockBankAccountData{
					{
						ID:        uuid.New(),
						UserID:    userID,
						Name:      "Main Account",
						Balance:   100000,
						CreatedAt: time.Now(),
						UpdatedAt: time.Now(),
					},
				}
				rows := helpers.ExpectBankAccountRows(mock, accounts)

				mock.ExpectQuery(`SELECT id, user_id, name, balance, created_at, updated_at FROM bank_accounts WHERE user_id = \$1 ORDER BY created_at DESC`).
					WithArgs(userID).
					WillReturnRows(rows)
			},
			expectedCount: 1,
			expectedError: false,
		},
		{
			name:   "no accounts found - empty result",
			userID: userID,
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "user_id", "name", "balance", "created_at", "updated_at",
				})

				mock.ExpectQuery(`SELECT id, user_id, name, balance, created_at, updated_at FROM bank_accounts WHERE user_id = \$1 ORDER BY created_at DESC`).
					WithArgs(userID).
					WillReturnRows(rows)
			},
			expectedCount: 0,
			expectedError: false,
		},
		{
			name:   "database connection error",
			userID: userID,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, user_id, name, balance, created_at, updated_at FROM bank_accounts WHERE user_id = \$1 ORDER BY created_at DESC`).
					WithArgs(userID).
					WillReturnError(sql.ErrConnDone)
			},
			expectedCount: 0,
			expectedError: true,
			errorType:     "connection",
		},
		{
			name:   "invalid user ID",
			userID: uuid.Nil,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, user_id, name, balance, created_at, updated_at FROM bank_accounts WHERE user_id = \$1 ORDER BY created_at DESC`).
					WithArgs(uuid.Nil).
					WillReturnRows(sqlmock.NewRows([]string{
						"id", "user_id", "name", "balance", "created_at", "updated_at",
					}))
			},
			expectedCount: 0,
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(mock)

			accounts, err := repo.GetAll(tt.userID)

			if tt.expectedError {
				assert.Error(t, err)
				if tt.errorType == "connection" {
					assert.True(t, helpers.IsDatabaseError(err))
				}
			} else {
				assert.NoError(t, err)
				assert.Len(t, accounts, tt.expectedCount)

				// Validate that all returned accounts belong to the correct user
				for _, account := range accounts {
					assert.Equal(t, tt.userID, account.UserID)
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
					"id", "user_id", "name", "balance", "created_at", "updated_at",
				}).AddRow(
					id, uuid.New(), "Test Account", int64(100000),
					time.Now(), time.Now(),
				)

				mock.ExpectQuery(`SELECT id, user_id, name, balance, created_at, updated_at FROM bank_accounts WHERE id = \$1`).
					WithArgs(id).
					WillReturnRows(row)
			},
			expectedError: false,
		},
		{
			name:      "account not found",
			accountID: uuid.New(),
			setupMock: func(mock sqlmock.Sqlmock, id uuid.UUID) {
				mock.ExpectQuery(`SELECT id, user_id, name, balance, created_at, updated_at FROM bank_accounts WHERE id = \$1`).
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
				mock.ExpectQuery(`SELECT id, user_id, name, balance, created_at, updated_at FROM bank_accounts WHERE id = \$1`).
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
				helpers.ValidateNonEmptyUUID(t, account.UserID)
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
			account: helpers.CreateTestBankAccount(uuid.New()),
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO bank_accounts \(id, user_id, name, balance, created_at, updated_at\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6\)`).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: false,
		},
		{
			name:    "database error",
			account: helpers.CreateTestBankAccount(uuid.New()),
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO bank_accounts \(id, user_id, name, balance, created_at, updated_at\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6\)`).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
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
				helpers.ValidateNonEmptyUUID(t, tt.account.UserID)
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
			account: helpers.CreateTestBankAccount(uuid.New()),
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE bank_accounts SET name = \$2, balance = \$3, updated_at = \$4 WHERE id = \$1`).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(0, 1))
			},
			expectedError: false,
		},
		{
			name:    "account not found",
			account: helpers.CreateTestBankAccount(uuid.New()),
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE bank_accounts SET name = \$2, balance = \$3, updated_at = \$4 WHERE id = \$1`).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(0, 0))
			},
			expectedError: false, // Repository doesn't return error for no rows affected
		},
		{
			name:    "database error",
			account: helpers.CreateTestBankAccount(uuid.New()),
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
