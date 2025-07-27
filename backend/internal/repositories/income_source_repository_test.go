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

func TestIncomeSourceRepository_GetAll(t *testing.T) {
	db, mock := helpers.SetupMockDB(t)
	defer helpers.TeardownMockDB(db)

	repo := NewIncomeSourceRepository(db)
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
				paymentDay := 25
				scheduledDate := "2024-12-25"
				rows := sqlmock.NewRows([]string{
					"id", "user_id", "name", "income_type", "base_amount", "bank_account",
					"payment_day", "scheduled_date", "scheduled_year_month", "is_active", "created_at", "updated_at",
				}).
					AddRow(
						uuid.New(), userID, "Salary", "monthly_fixed", int64(300000), bankAccountID,
						&paymentDay, nil, nil, true, time.Now(), time.Now(),
					).
					AddRow(
						uuid.New(), userID, "Bonus", "one_time", int64(100000), bankAccountID,
						nil, &scheduledDate, nil, true, time.Now(), time.Now(),
					)

				mock.ExpectQuery(`SELECT id, user_id, name, income_type, base_amount, bank_account, payment_day, scheduled_date::text, scheduled_year_month, is_active, created_at, updated_at FROM income_sources WHERE user_id = \$1 ORDER BY created_at DESC`).
					WithArgs(userID).
					WillReturnRows(rows)
			},
			expectedCount: 2,
			expectedError: false,
		},
		{
			name:   "no income sources found",
			userID: userID,
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "user_id", "name", "income_type", "base_amount", "bank_account",
					"payment_day", "scheduled_date", "scheduled_year_month", "is_active", "created_at", "updated_at",
				})

				mock.ExpectQuery(`SELECT id, user_id, name, income_type, base_amount, bank_account, payment_day, scheduled_date::text, scheduled_year_month, is_active, created_at, updated_at FROM income_sources WHERE user_id = \$1 ORDER BY created_at DESC`).
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
				mock.ExpectQuery(`SELECT id, user_id, name, income_type, base_amount, bank_account, payment_day, scheduled_date::text, scheduled_year_month, is_active, created_at, updated_at FROM income_sources WHERE user_id = \$1 ORDER BY created_at DESC`).
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

			sources, err := repo.GetAll(tt.userID)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, sources, tt.expectedCount)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestIncomeSourceRepository_GetByID(t *testing.T) {
	db, mock := helpers.SetupMockDB(t)
	defer helpers.TeardownMockDB(db)

	repo := NewIncomeSourceRepository(db)
	sourceID := uuid.New()
	userID := uuid.New()
	bankAccountID := uuid.New()

	tests := []struct {
		name          string
		sourceID      uuid.UUID
		setupMock     func(sqlmock.Sqlmock)
		expectedFound bool
		expectedError bool
	}{
		{
			name:     "income source found",
			sourceID: sourceID,
			setupMock: func(mock sqlmock.Sqlmock) {
				paymentDay := 25
				rows := sqlmock.NewRows([]string{
					"id", "user_id", "name", "income_type", "base_amount", "bank_account",
					"payment_day", "scheduled_date", "scheduled_year_month", "is_active", "created_at", "updated_at",
				}).
					AddRow(
						sourceID, userID, "Salary", "monthly_fixed", int64(300000), bankAccountID,
						&paymentDay, nil, nil, true, time.Now(), time.Now(),
					)

				mock.ExpectQuery(`SELECT id, user_id, name, income_type, base_amount, bank_account, payment_day, scheduled_date::text, scheduled_year_month, is_active, created_at, updated_at FROM income_sources WHERE id = \$1`).
					WithArgs(sourceID).
					WillReturnRows(rows)
			},
			expectedFound: true,
			expectedError: false,
		},
		{
			name:     "income source not found",
			sourceID: sourceID,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, user_id, name, income_type, base_amount, bank_account, payment_day, scheduled_date::text, scheduled_year_month, is_active, created_at, updated_at FROM income_sources WHERE id = \$1`).
					WithArgs(sourceID).
					WillReturnError(sql.ErrNoRows)
			},
			expectedFound: false,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(mock)

			source, err := repo.GetByID(tt.sourceID)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, source)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, source)
				assert.Equal(t, tt.sourceID, source.ID)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestIncomeSourceRepository_GetActiveByUserID(t *testing.T) {
	db, mock := helpers.SetupMockDB(t)
	defer helpers.TeardownMockDB(db)

	repo := NewIncomeSourceRepository(db)
	userID := uuid.New()

	tests := []struct {
		name          string
		userID        uuid.UUID
		setupMock     func(sqlmock.Sqlmock)
		expectedCount int
		expectedError bool
	}{
		{
			name:   "successful retrieval of active sources",
			userID: userID,
			setupMock: func(mock sqlmock.Sqlmock) {
				bankAccountID := uuid.New()
				paymentDay := 25
				rows := sqlmock.NewRows([]string{
					"id", "user_id", "name", "income_type", "base_amount", "bank_account",
					"payment_day", "scheduled_date", "scheduled_year_month", "is_active", "created_at", "updated_at",
				}).
					AddRow(
						uuid.New(), userID, "Salary", "monthly_fixed", int64(300000), bankAccountID,
						&paymentDay, nil, nil, true, time.Now(), time.Now(),
					)

				mock.ExpectQuery(`SELECT id, user_id, name, income_type, base_amount, bank_account, payment_day, scheduled_date::text, scheduled_year_month, is_active, created_at, updated_at FROM income_sources WHERE user_id = \$1 AND is_active = true ORDER BY created_at DESC`).
					WithArgs(userID).
					WillReturnRows(rows)
			},
			expectedCount: 1,
			expectedError: false,
		},
		{
			name:   "no active sources found",
			userID: userID,
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "user_id", "name", "income_type", "base_amount", "bank_account",
					"payment_day", "scheduled_date", "scheduled_year_month", "is_active", "created_at", "updated_at",
				})

				mock.ExpectQuery(`SELECT id, user_id, name, income_type, base_amount, bank_account, payment_day, scheduled_date::text, scheduled_year_month, is_active, created_at, updated_at FROM income_sources WHERE user_id = \$1 AND is_active = true ORDER BY created_at DESC`).
					WithArgs(userID).
					WillReturnRows(rows)
			},
			expectedCount: 0,
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(mock)

			sources, err := repo.GetActiveByUserID(tt.userID)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, sources, tt.expectedCount)
				// Verify all returned sources are active
				for _, source := range sources {
					assert.True(t, source.IsActive)
				}
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestIncomeSourceRepository_Create(t *testing.T) {
	db, mock := helpers.SetupMockDB(t)
	defer helpers.TeardownMockDB(db)

	repo := NewIncomeSourceRepository(db)
	userID := uuid.New()
	bankAccountID := uuid.New()
	source := helpers.CreateTestIncomeSource(userID, bankAccountID)

	tests := []struct {
		name          string
		source        *models.IncomeSource
		setupMock     func(sqlmock.Sqlmock)
		expectedError bool
	}{
		{
			name:   "successful creation",
			source: source,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO income_sources \(id, user_id, name, income_type, base_amount, bank_account, payment_day, scheduled_date, scheduled_year_month, is_active, created_at, updated_at\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8, \$9, \$10, \$11, \$12\)`).
					WithArgs(source.ID, source.UserID, source.Name, source.IncomeType, source.BaseAmount, source.BankAccount, source.PaymentDay, source.ScheduledDate, source.ScheduledYearMonth, source.IsActive, sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: false,
		},
		{
			name:   "database error",
			source: source,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO income_sources \(id, user_id, name, income_type, base_amount, bank_account, payment_day, scheduled_date, scheduled_year_month, is_active, created_at, updated_at\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8, \$9, \$10, \$11, \$12\)`).
					WithArgs(source.ID, source.UserID, source.Name, source.IncomeType, source.BaseAmount, source.BankAccount, source.PaymentDay, source.ScheduledDate, source.ScheduledYearMonth, source.IsActive, sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(sql.ErrConnDone)
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(mock)

			err := repo.Create(tt.source)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestIncomeSourceRepository_Update(t *testing.T) {
	db, mock := helpers.SetupMockDB(t)
	defer helpers.TeardownMockDB(db)

	repo := NewIncomeSourceRepository(db)
	userID := uuid.New()
	bankAccountID := uuid.New()
	source := helpers.CreateTestIncomeSource(userID, bankAccountID)
	source.Name = "Updated Income Source"

	tests := []struct {
		name          string
		source        *models.IncomeSource
		setupMock     func(sqlmock.Sqlmock)
		expectedError bool
	}{
		{
			name:   "successful update",
			source: source,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE income_sources SET name = \$2, income_type = \$3, base_amount = \$4, bank_account = \$5, payment_day = \$6, scheduled_date = \$7, scheduled_year_month = \$8, is_active = \$9, updated_at = \$10 WHERE id = \$1`).
					WithArgs(source.ID, source.Name, source.IncomeType, source.BaseAmount, source.BankAccount, source.PaymentDay, source.ScheduledDate, source.ScheduledYearMonth, source.IsActive, sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: false,
		},
		{
			name:   "database error",
			source: source,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE income_sources SET name = \$2, income_type = \$3, base_amount = \$4, bank_account = \$5, payment_day = \$6, scheduled_date = \$7, scheduled_year_month = \$8, is_active = \$9, updated_at = \$10 WHERE id = \$1`).
					WithArgs(source.ID, source.Name, source.IncomeType, source.BaseAmount, source.BankAccount, source.PaymentDay, source.ScheduledDate, source.ScheduledYearMonth, source.IsActive, sqlmock.AnyArg()).
					WillReturnError(sql.ErrConnDone)
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(mock)

			err := repo.Update(tt.source)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestIncomeSourceRepository_Delete(t *testing.T) {
	db, mock := helpers.SetupMockDB(t)
	defer helpers.TeardownMockDB(db)

	repo := NewIncomeSourceRepository(db)
	sourceID := uuid.New()

	tests := []struct {
		name          string
		sourceID      uuid.UUID
		setupMock     func(sqlmock.Sqlmock)
		expectedError bool
	}{
		{
			name:     "successful deletion",
			sourceID: sourceID,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM income_sources WHERE id = \$1`).
					WithArgs(sourceID).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: false,
		},
		{
			name:     "database error",
			sourceID: sourceID,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`DELETE FROM income_sources WHERE id = \$1`).
					WithArgs(sourceID).
					WillReturnError(sql.ErrConnDone)
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(mock)

			err := repo.Delete(tt.sourceID)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
