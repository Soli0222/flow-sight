package repositories

import (
	"database/sql"
	"flow-sight-backend/internal/models"

	"github.com/google/uuid"
)

type CreditCardRepository struct {
	db *sql.DB
}

func NewCreditCardRepository(db *sql.DB) *CreditCardRepository {
	return &CreditCardRepository{db: db}
}

func (r *CreditCardRepository) GetAll(userID uuid.UUID) ([]models.CreditCard, error) {
	query := `
		SELECT id, user_id, name, closing_day, payment_day, bank_account, created_at, updated_at
		FROM credit_cards 
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return []models.CreditCard{}, err
	}
	defer rows.Close()

	creditCards := make([]models.CreditCard, 0)
	for rows.Next() {
		var creditCard models.CreditCard
		err := rows.Scan(
			&creditCard.ID, &creditCard.UserID, &creditCard.Name,
			&creditCard.ClosingDay, &creditCard.PaymentDay, &creditCard.BankAccount,
			&creditCard.CreatedAt, &creditCard.UpdatedAt,
		)
		if err != nil {
			return []models.CreditCard{}, err
		}
		creditCards = append(creditCards, creditCard)
	}

	return creditCards, nil
}

func (r *CreditCardRepository) GetByID(id uuid.UUID) (*models.CreditCard, error) {
	query := `
		SELECT id, user_id, name, closing_day, payment_day, bank_account, created_at, updated_at
		FROM credit_cards 
		WHERE id = $1
	`

	var creditCard models.CreditCard
	err := r.db.QueryRow(query, id).Scan(
		&creditCard.ID, &creditCard.UserID, &creditCard.Name,
		&creditCard.ClosingDay, &creditCard.PaymentDay, &creditCard.BankAccount,
		&creditCard.CreatedAt, &creditCard.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &creditCard, nil
}

func (r *CreditCardRepository) Create(creditCard *models.CreditCard) error {
	query := `
		INSERT INTO credit_cards (id, user_id, name, closing_day, payment_day, bank_account, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.Exec(query,
		creditCard.ID, creditCard.UserID, creditCard.Name,
		creditCard.ClosingDay, creditCard.PaymentDay, creditCard.BankAccount,
		creditCard.CreatedAt, creditCard.UpdatedAt,
	)

	return err
}

func (r *CreditCardRepository) Update(creditCard *models.CreditCard) error {
	query := `
		UPDATE credit_cards 
		SET name = $2, closing_day = $3, payment_day = $4, 
		    bank_account = $5, updated_at = $6
		WHERE id = $1
	`

	_, err := r.db.Exec(query,
		creditCard.ID, creditCard.Name, creditCard.ClosingDay,
		creditCard.PaymentDay, creditCard.BankAccount, creditCard.UpdatedAt,
	)

	return err
}

func (r *CreditCardRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM credit_cards WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
