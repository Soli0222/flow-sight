package repositories

import (
	"flow-sight-backend/internal/models"
	"flow-sight-backend/test/helpers"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestCardMonthlyTotalRepository_GetByCreditCardID(t *testing.T) {
	t.Run("successful retrieval", func(t *testing.T) {
		db, mock := helpers.SetupMockDB(t)
		defer helpers.TeardownMockDB(db)

		repo := NewCardMonthlyTotalRepository(db)
		creditCardID := uuid.New()

		expected := []*models.CardMonthlyTotal{
			helpers.CreateTestCardMonthlyTotal(),
			helpers.CreateTestCardMonthlyTotal(),
		}
		expected[0].CreditCardID = creditCardID
		expected[1].CreditCardID = creditCardID

		rows := sqlmock.NewRows([]string{
			"id", "credit_card_id", "year_month", "total_amount",
			"is_confirmed", "created_at", "updated_at",
		})
		for _, total := range expected {
			rows.AddRow(
				total.ID, total.CreditCardID, total.YearMonth, total.TotalAmount,
				total.IsConfirmed, total.CreatedAt, total.UpdatedAt,
			)
		}

		mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, credit_card_id, year_month, total_amount, is_confirmed, created_at, updated_at
		FROM card_monthly_totals 
		WHERE credit_card_id = $1
		ORDER BY year_month DESC
	`)).WithArgs(creditCardID).WillReturnRows(rows)

		result, err := repo.GetByCreditCardID(creditCardID)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, expected[0].ID, result[0].ID)
		assert.Equal(t, expected[1].ID, result[1].ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("no totals found", func(t *testing.T) {
		db, mock := helpers.SetupMockDB(t)
		defer helpers.TeardownMockDB(db)

		repo := NewCardMonthlyTotalRepository(db)
		creditCardID := uuid.New()

		rows := sqlmock.NewRows([]string{
			"id", "credit_card_id", "year_month", "total_amount",
			"is_confirmed", "created_at", "updated_at",
		})

		mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, credit_card_id, year_month, total_amount, is_confirmed, created_at, updated_at
		FROM card_monthly_totals 
		WHERE credit_card_id = $1
		ORDER BY year_month DESC
	`)).WithArgs(creditCardID).WillReturnRows(rows)

		result, err := repo.GetByCreditCardID(creditCardID)

		assert.NoError(t, err)
		assert.Len(t, result, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		db, mock := helpers.SetupMockDB(t)
		defer helpers.TeardownMockDB(db)

		repo := NewCardMonthlyTotalRepository(db)
		creditCardID := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, credit_card_id, year_month, total_amount, is_confirmed, created_at, updated_at
		FROM card_monthly_totals 
		WHERE credit_card_id = $1
		ORDER BY year_month DESC
	`)).WithArgs(creditCardID).WillReturnError(assert.AnError)

		result, err := repo.GetByCreditCardID(creditCardID)

		assert.Error(t, err)
		assert.Empty(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCardMonthlyTotalRepository_GetByYearMonth(t *testing.T) {
	t.Run("successful retrieval", func(t *testing.T) {
		db, mock := helpers.SetupMockDB(t)
		defer helpers.TeardownMockDB(db)

		repo := NewCardMonthlyTotalRepository(db)
		yearMonth := "2024-01"

		expected := []*models.CardMonthlyTotal{
			helpers.CreateTestCardMonthlyTotal(),
			helpers.CreateTestCardMonthlyTotal(),
		}
		expected[0].YearMonth = yearMonth
		expected[1].YearMonth = yearMonth

		rows := sqlmock.NewRows([]string{
			"id", "credit_card_id", "year_month", "total_amount",
			"is_confirmed", "created_at", "updated_at",
		})
		for _, total := range expected {
			rows.AddRow(
				total.ID, total.CreditCardID, total.YearMonth, total.TotalAmount,
				total.IsConfirmed, total.CreatedAt, total.UpdatedAt,
			)
		}

		mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, credit_card_id, year_month, total_amount, is_confirmed, created_at, updated_at
		FROM card_monthly_totals 
		WHERE year_month = $1
		ORDER BY created_at DESC
	`)).WithArgs(yearMonth).WillReturnRows(rows)

		result, err := repo.GetByYearMonth(yearMonth)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
		assert.Equal(t, expected[0].ID, result[0].ID)
		assert.Equal(t, expected[1].ID, result[1].ID)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("no totals found", func(t *testing.T) {
		db, mock := helpers.SetupMockDB(t)
		defer helpers.TeardownMockDB(db)

		repo := NewCardMonthlyTotalRepository(db)
		yearMonth := "2024-01"

		rows := sqlmock.NewRows([]string{
			"id", "credit_card_id", "year_month", "total_amount",
			"is_confirmed", "created_at", "updated_at",
		})

		mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, credit_card_id, year_month, total_amount, is_confirmed, created_at, updated_at
		FROM card_monthly_totals 
		WHERE year_month = $1
		ORDER BY created_at DESC
	`)).WithArgs(yearMonth).WillReturnRows(rows)

		result, err := repo.GetByYearMonth(yearMonth)

		assert.NoError(t, err)
		assert.Len(t, result, 0)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCardMonthlyTotalRepository_GetByID(t *testing.T) {
	t.Run("total found", func(t *testing.T) {
		db, mock := helpers.SetupMockDB(t)
		defer helpers.TeardownMockDB(db)

		repo := NewCardMonthlyTotalRepository(db)
		expected := helpers.CreateTestCardMonthlyTotal()

		rows := sqlmock.NewRows([]string{
			"id", "credit_card_id", "year_month", "total_amount",
			"is_confirmed", "created_at", "updated_at",
		}).AddRow(
			expected.ID, expected.CreditCardID, expected.YearMonth, expected.TotalAmount,
			expected.IsConfirmed, expected.CreatedAt, expected.UpdatedAt,
		)

		mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, credit_card_id, year_month, total_amount, is_confirmed, created_at, updated_at
		FROM card_monthly_totals 
		WHERE id = $1
	`)).WithArgs(expected.ID).WillReturnRows(rows)

		result, err := repo.GetByID(expected.ID)

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, expected.ID, result.ID)
		assert.Equal(t, expected.YearMonth, result.YearMonth)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("total not found", func(t *testing.T) {
		db, mock := helpers.SetupMockDB(t)
		defer helpers.TeardownMockDB(db)

		repo := NewCardMonthlyTotalRepository(db)
		id := uuid.New()

		mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, credit_card_id, year_month, total_amount, is_confirmed, created_at, updated_at
		FROM card_monthly_totals 
		WHERE id = $1
	`)).WithArgs(id).WillReturnError(assert.AnError)

		result, err := repo.GetByID(id)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCardMonthlyTotalRepository_Create(t *testing.T) {
	t.Run("successful creation", func(t *testing.T) {
		db, mock := helpers.SetupMockDB(t)
		defer helpers.TeardownMockDB(db)

		repo := NewCardMonthlyTotalRepository(db)
		total := helpers.CreateTestCardMonthlyTotal()

		mock.ExpectExec(regexp.QuoteMeta(`
		INSERT INTO card_monthly_totals (id, credit_card_id, year_month, total_amount, 
		                                is_confirmed, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`)).WithArgs(
			total.ID, total.CreditCardID, total.YearMonth, total.TotalAmount,
			total.IsConfirmed, sqlmock.AnyArg(), sqlmock.AnyArg(),
		).WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Create(total)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		db, mock := helpers.SetupMockDB(t)
		defer helpers.TeardownMockDB(db)

		repo := NewCardMonthlyTotalRepository(db)
		total := helpers.CreateTestCardMonthlyTotal()

		mock.ExpectExec(regexp.QuoteMeta(`
		INSERT INTO card_monthly_totals (id, credit_card_id, year_month, total_amount, 
		                                is_confirmed, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`)).WithArgs(
			total.ID, total.CreditCardID, total.YearMonth, total.TotalAmount,
			total.IsConfirmed, sqlmock.AnyArg(), sqlmock.AnyArg(),
		).WillReturnError(assert.AnError)

		err := repo.Create(total)

		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCardMonthlyTotalRepository_Update(t *testing.T) {
	t.Run("successful update", func(t *testing.T) {
		db, mock := helpers.SetupMockDB(t)
		defer helpers.TeardownMockDB(db)

		repo := NewCardMonthlyTotalRepository(db)
		total := helpers.CreateTestCardMonthlyTotal()
		total.TotalAmount = 200000
		total.IsConfirmed = true
		total.UpdatedAt = time.Now()

		mock.ExpectExec(regexp.QuoteMeta(`
		UPDATE card_monthly_totals 
		SET year_month = $2, total_amount = $3, is_confirmed = $4, updated_at = $5
		WHERE id = $1
	`)).WithArgs(
			total.ID, total.YearMonth, total.TotalAmount, total.IsConfirmed, sqlmock.AnyArg(),
		).WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Update(total)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		db, mock := helpers.SetupMockDB(t)
		defer helpers.TeardownMockDB(db)

		repo := NewCardMonthlyTotalRepository(db)
		total := helpers.CreateTestCardMonthlyTotal()

		mock.ExpectExec(regexp.QuoteMeta(`
		UPDATE card_monthly_totals 
		SET year_month = $2, total_amount = $3, is_confirmed = $4, updated_at = $5
		WHERE id = $1
	`)).WithArgs(
			total.ID, total.YearMonth, total.TotalAmount, total.IsConfirmed, sqlmock.AnyArg(),
		).WillReturnError(assert.AnError)

		err := repo.Update(total)

		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}

func TestCardMonthlyTotalRepository_Delete(t *testing.T) {
	t.Run("successful deletion", func(t *testing.T) {
		db, mock := helpers.SetupMockDB(t)
		defer helpers.TeardownMockDB(db)

		repo := NewCardMonthlyTotalRepository(db)
		id := uuid.New()

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM card_monthly_totals WHERE id = $1`)).
			WithArgs(id).
			WillReturnResult(sqlmock.NewResult(1, 1))

		err := repo.Delete(id)

		assert.NoError(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("database error", func(t *testing.T) {
		db, mock := helpers.SetupMockDB(t)
		defer helpers.TeardownMockDB(db)

		repo := NewCardMonthlyTotalRepository(db)
		id := uuid.New()

		mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM card_monthly_totals WHERE id = $1`)).
			WithArgs(id).
			WillReturnError(assert.AnError)

		err := repo.Delete(id)

		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})
}
