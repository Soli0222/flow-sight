package repositories

import (
	"database/sql"

	"github.com/Soli0222/flow-sight/backend/internal/models"

	"github.com/google/uuid"
)

type BankAccountRepository struct {
	db *sql.DB
}

func NewBankAccountRepository(db *sql.DB) *BankAccountRepository {
	return &BankAccountRepository{db: db}
}

func (r *BankAccountRepository) GetAll() ([]models.BankAccount, error) {
	query := `
		SELECT id, name, balance, created_at, updated_at
		FROM bank_accounts 
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return []models.BankAccount{}, err
	}
	defer rows.Close()

	accounts := make([]models.BankAccount, 0)
	for rows.Next() {
		var account models.BankAccount
		err := rows.Scan(
			&account.ID, &account.Name, &account.Balance,
			&account.CreatedAt, &account.UpdatedAt,
		)
		if err != nil {
			return []models.BankAccount{}, err
		}
		accounts = append(accounts, account)
	}

	return accounts, nil
}

func (r *BankAccountRepository) GetByID(id uuid.UUID) (*models.BankAccount, error) {
	query := `
		SELECT id, name, balance, created_at, updated_at
		FROM bank_accounts 
		WHERE id = $1
	`

	var account models.BankAccount
	err := r.db.QueryRow(query, id).Scan(
		&account.ID, &account.Name, &account.Balance,
		&account.CreatedAt, &account.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &account, nil
}

func (r *BankAccountRepository) Create(account *models.BankAccount) error {
	query := `
		INSERT INTO bank_accounts (id, name, balance, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
	`

	_, err := r.db.Exec(query,
		account.ID, account.Name, account.Balance,
		account.CreatedAt, account.UpdatedAt,
	)

	return err
}

func (r *BankAccountRepository) Update(account *models.BankAccount) error {
	query := `
		UPDATE bank_accounts 
		SET name = $2, balance = $3, updated_at = $4
		WHERE id = $1
	`

	_, err := r.db.Exec(query,
		account.ID, account.Name, account.Balance, account.UpdatedAt,
	)

	return err
}

func (r *BankAccountRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM bank_accounts WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
