package repositories

import (
	"database/sql"
	"github.com/Soli0222/flow-sight/backend/test/helpers"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestMonthlyIncomeRepository_GetByIncomeSourceID(t *testing.T) {
	db, mock := helpers.SetupMockDB(t)
	defer helpers.TeardownMockDB(db)

	repo := NewMonthlyIncomeRepository(db)
	incomeSourceID := uuid.New()

	tests := []struct {
		name           string
		incomeSourceID uuid.UUID
		setupMock      func(sqlmock.Sqlmock)
		expectedCount  int
		expectedError  bool
	}{
		{
			name:           "successful retrieval",
			incomeSourceID: incomeSourceID,
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "income_source_id", "year_month", "actual_amount", "is_confirmed", "note", "created_at", "updated_at",
				}).
					AddRow(
						uuid.New(), incomeSourceID, "2024-01", int64(300000), true, "January salary", time.Now(), time.Now(),
					).
					AddRow(
						uuid.New(), incomeSourceID, "2024-02", int64(320000), false, "February salary", time.Now(), time.Now(),
					)

				mock.ExpectQuery(`SELECT id, income_source_id, year_month, actual_amount, is_confirmed, note, created_at, updated_at FROM monthly_income_records WHERE income_source_id = \$1 ORDER BY year_month DESC`).
					WithArgs(incomeSourceID).
					WillReturnRows(rows)
			},
			expectedCount: 2,
			expectedError: false,
		},
		{
			name:           "no records found",
			incomeSourceID: incomeSourceID,
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "income_source_id", "year_month", "actual_amount", "is_confirmed", "note", "created_at", "updated_at",
				})

				mock.ExpectQuery(`SELECT id, income_source_id, year_month, actual_amount, is_confirmed, note, created_at, updated_at FROM monthly_income_records WHERE income_source_id = \$1 ORDER BY year_month DESC`).
					WithArgs(incomeSourceID).
					WillReturnRows(rows)
			},
			expectedCount: 0,
			expectedError: false,
		},
		{
			name:           "database error",
			incomeSourceID: incomeSourceID,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, income_source_id, year_month, actual_amount, is_confirmed, note, created_at, updated_at FROM monthly_income_records WHERE income_source_id = \$1 ORDER BY year_month DESC`).
					WithArgs(incomeSourceID).
					WillReturnError(sql.ErrConnDone)
			},
			expectedCount: 0,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(mock)

			records, err := repo.GetByIncomeSourceID(tt.incomeSourceID)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, records, tt.expectedCount)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestMonthlyIncomeRepository_GetByYearMonth(t *testing.T) {
	db, mock := helpers.SetupMockDB(t)
	defer helpers.TeardownMockDB(db)

	repo := NewMonthlyIncomeRepository(db)
	yearMonth := "2024-01"

	tests := []struct {
		name          string
		yearMonth     string
		setupMock     func(sqlmock.Sqlmock)
		expectedCount int
		expectedError bool
	}{
		{
			name:      "successful retrieval",
			yearMonth: yearMonth,
			setupMock: func(mock sqlmock.Sqlmock) {
				incomeSourceID1 := uuid.New()
				incomeSourceID2 := uuid.New()
				rows := sqlmock.NewRows([]string{
					"id", "income_source_id", "year_month", "actual_amount", "is_confirmed", "note", "created_at", "updated_at",
				}).
					AddRow(
						uuid.New(), incomeSourceID1, yearMonth, int64(300000), true, "January salary", time.Now(), time.Now(),
					).
					AddRow(
						uuid.New(), incomeSourceID2, yearMonth, int64(50000), true, "January bonus", time.Now(), time.Now(),
					)

				mock.ExpectQuery(`SELECT id, income_source_id, year_month, actual_amount, is_confirmed, note, created_at, updated_at FROM monthly_income_records WHERE year_month = \$1 ORDER BY created_at DESC`).
					WithArgs(yearMonth).
					WillReturnRows(rows)
			},
			expectedCount: 2,
			expectedError: false,
		},
		{
			name:      "no records found",
			yearMonth: yearMonth,
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "income_source_id", "year_month", "actual_amount", "is_confirmed", "note", "created_at", "updated_at",
				})

				mock.ExpectQuery(`SELECT id, income_source_id, year_month, actual_amount, is_confirmed, note, created_at, updated_at FROM monthly_income_records WHERE year_month = \$1 ORDER BY created_at DESC`).
					WithArgs(yearMonth).
					WillReturnRows(rows)
			},
			expectedCount: 0,
			expectedError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(mock)

			records, err := repo.GetByYearMonth(tt.yearMonth)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, records, tt.expectedCount)
				// Verify all returned records have the correct year_month
				for _, record := range records {
					assert.Equal(t, tt.yearMonth, record.YearMonth)
				}
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
