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

func TestRecurringPaymentRepository_GetAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewRecurringPaymentRepository(db)

	rows := sqlmock.NewRows([]string{"id", "name", "amount", "payment_day", "start_year_month", "total_payments", "remaining_payments", "bank_account", "is_active", "note", "created_at", "updated_at"}).
		AddRow(uuid.New(), "Payment1", 1000, 10, "2024-01", nil, nil, uuid.New(), true, "note", time.Now(), time.Now())

	mock.ExpectQuery(regexp.QuoteMeta(`SELECT id, name, amount, payment_day, start_year_month, total_payments, remaining_payments, bank_account, is_active, note, created_at, updated_at FROM recurring_payments ORDER BY created_at DESC`)).
		WillReturnRows(rows)

	payments, err := repo.GetAll()
	assert.NoError(t, err)
	assert.Len(t, payments, 1)
}

func TestRecurringPaymentRepository_CRUD(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewRecurringPaymentRepository(db)

	bankAccountID := uuid.New()
	payment := helpers.CreateTestRecurringPayment(bankAccountID)

	mock.ExpectExec(regexp.QuoteMeta(`INSERT INTO recurring_payments`)).
		WithArgs(sqlmock.AnyArg(), payment.Name, payment.Amount, payment.PaymentDay, payment.StartYearMonth, payment.TotalPayments, payment.RemainingPayments, payment.BankAccount, payment.IsActive, payment.Note, sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Create(payment)
	assert.NoError(t, err)

	payment.Name = "Updated Payment"
	mock.ExpectExec(regexp.QuoteMeta(`UPDATE recurring_payments SET`)).
		WithArgs(payment.ID, payment.Name, payment.Amount, payment.PaymentDay, payment.StartYearMonth, payment.TotalPayments, payment.RemainingPayments, payment.BankAccount, payment.IsActive, payment.Note, sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Update(payment)
	assert.NoError(t, err)

	mock.ExpectExec(regexp.QuoteMeta(`DELETE FROM recurring_payments WHERE id = $1`)).
		WithArgs(payment.ID).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.Delete(payment.ID)
	assert.NoError(t, err)
}
