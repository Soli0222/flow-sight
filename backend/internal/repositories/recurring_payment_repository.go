package repositories

import (
	"database/sql"
	"github.com/Soli0222/flow-sight/backend/internal/models"

	"github.com/google/uuid"
)

type RecurringPaymentRepository struct {
	db *sql.DB
}

func NewRecurringPaymentRepository(db *sql.DB) *RecurringPaymentRepository {
	return &RecurringPaymentRepository{db: db}
}

func (r *RecurringPaymentRepository) GetAll(userID uuid.UUID) ([]models.RecurringPayment, error) {
	query := `
		SELECT id, user_id, name, amount, payment_day, start_year_month, 
		       total_payments, remaining_payments, bank_account, is_active, 
		       note, created_at, updated_at
		FROM recurring_payments 
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return []models.RecurringPayment{}, err
	}
	defer rows.Close()

	payments := make([]models.RecurringPayment, 0)
	for rows.Next() {
		var payment models.RecurringPayment
		err := rows.Scan(
			&payment.ID, &payment.UserID, &payment.Name, &payment.Amount,
			&payment.PaymentDay, &payment.StartYearMonth, &payment.TotalPayments,
			&payment.RemainingPayments, &payment.BankAccount, &payment.IsActive,
			&payment.Note, &payment.CreatedAt, &payment.UpdatedAt,
		)
		if err != nil {
			return []models.RecurringPayment{}, err
		}
		payments = append(payments, payment)
	}

	return payments, nil
}

func (r *RecurringPaymentRepository) GetActiveByUserID(userID uuid.UUID) ([]models.RecurringPayment, error) {
	query := `
		SELECT id, user_id, name, amount, payment_day, start_year_month, 
		       total_payments, remaining_payments, bank_account, is_active, 
		       note, created_at, updated_at
		FROM recurring_payments 
		WHERE user_id = $1 AND is_active = true
		ORDER BY payment_day ASC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var payments []models.RecurringPayment
	for rows.Next() {
		var payment models.RecurringPayment
		err := rows.Scan(
			&payment.ID, &payment.UserID, &payment.Name, &payment.Amount,
			&payment.PaymentDay, &payment.StartYearMonth, &payment.TotalPayments,
			&payment.RemainingPayments, &payment.BankAccount, &payment.IsActive,
			&payment.Note, &payment.CreatedAt, &payment.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		payments = append(payments, payment)
	}

	return payments, nil
}

func (r *RecurringPaymentRepository) GetByID(id uuid.UUID) (*models.RecurringPayment, error) {
	query := `
		SELECT id, user_id, name, amount, payment_day, start_year_month, 
		       total_payments, remaining_payments, bank_account, is_active, 
		       note, created_at, updated_at
		FROM recurring_payments 
		WHERE id = $1
	`

	var payment models.RecurringPayment
	err := r.db.QueryRow(query, id).Scan(
		&payment.ID, &payment.UserID, &payment.Name, &payment.Amount,
		&payment.PaymentDay, &payment.StartYearMonth, &payment.TotalPayments,
		&payment.RemainingPayments, &payment.BankAccount, &payment.IsActive,
		&payment.Note, &payment.CreatedAt, &payment.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &payment, nil
}

func (r *RecurringPaymentRepository) Create(payment *models.RecurringPayment) error {
	query := `
		INSERT INTO recurring_payments (id, user_id, name, amount, payment_day, 
		                               start_year_month, total_payments, remaining_payments, 
		                               bank_account, is_active, note, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13)
	`

	_, err := r.db.Exec(query,
		payment.ID, payment.UserID, payment.Name, payment.Amount,
		payment.PaymentDay, payment.StartYearMonth, payment.TotalPayments,
		payment.RemainingPayments, payment.BankAccount, payment.IsActive,
		payment.Note, payment.CreatedAt, payment.UpdatedAt,
	)

	return err
}

func (r *RecurringPaymentRepository) Update(payment *models.RecurringPayment) error {
	query := `
		UPDATE recurring_payments 
		SET name = $2, amount = $3, payment_day = $4, start_year_month = $5,
		    total_payments = $6, remaining_payments = $7, bank_account = $8,
		    is_active = $9, note = $10, updated_at = $11
		WHERE id = $1
	`

	_, err := r.db.Exec(query,
		payment.ID, payment.Name, payment.Amount, payment.PaymentDay,
		payment.StartYearMonth, payment.TotalPayments, payment.RemainingPayments,
		payment.BankAccount, payment.IsActive, payment.Note, payment.UpdatedAt,
	)

	return err
}

func (r *RecurringPaymentRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM recurring_payments WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
