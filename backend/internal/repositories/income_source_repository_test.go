package repositories

import (
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/Soli0222/flow-sight/backend/test/helpers"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestIncomeSourceRepository_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewIncomeSourceRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name", "income_type", "base_amount", "bank_account", "payment_day", "scheduled_date", "scheduled_year_month", "is_active", "created_at", "updated_at"}).
		AddRow(uuid.New(), "Source1", "monthly_fixed", 1000, uuid.New(), 10, nil, nil, true, time.Now(), time.Now()).
		AddRow(uuid.New(), "Source2", "one_time", 2000, uuid.New(), nil, "2024-07-01", nil, true, time.Now(), time.Now())

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, income_type, base_amount, bank_account, payment_day, scheduled_date::text, scheduled_year_month, is_active, created_at, updated_at FROM income_sources ORDER BY created_at DESC`)).
		WillReturnRows(rows)

	sources, err := repo.GetAll()
	assert.NoError(t, err)
	assert.Len(t, sources, 2)
}

func TestIncomeSourceRepository_GetActive(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewIncomeSourceRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name", "income_type", "base_amount", "bank_account", "payment_day", "scheduled_date", "scheduled_year_month", "is_active", "created_at", "updated_at"}).
		AddRow(uuid.New(), "Source1", "monthly_fixed", 1000, uuid.New(), 10, nil, nil, true, time.Now(), time.Now())

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, income_type, base_amount, bank_account, payment_day, scheduled_date::text, scheduled_year_month, is_active, created_at, updated_at FROM income_sources WHERE is_active = true ORDER BY created_at DESC`)).
		WillReturnRows(rows)

	sources, err := repo.GetActive()
	assert.NoError(t, err)
	assert.Len(t, sources, 1)
}

func TestIncomeSourceRepository_CRUD(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewIncomeSourceRepository(db)

	bankAccountID := uuid.New()
	source := helpers.CreateTestIncomeSource(bankAccountID)

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO income_sources`)).
		WithArgs(sqlmock.AnyArg(), source.Name, source.IncomeType, source.BaseAmount, source.BankAccount, source.PaymentDay, source.ScheduledDate, source.ScheduledYearMonth, source.IsActive, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(source)
	assert.NoError(t, err)

	source.Name = "Updated Name"
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE income_sources SET`)).
		WithArgs(source.ID, source.Name, source.IncomeType, source.BaseAmount, source.BankAccount, source.PaymentDay, source.ScheduledDate, source.ScheduledYearMonth, source.IsActive, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Update(source)
	assert.NoError(t, err)

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM income_sources WHERE id = $1`)).
		WithArgs(source.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Delete(source.ID)
	assert.NoError(t, err)
}
