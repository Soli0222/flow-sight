package repositories

import (
	"database/sql"
	"flow-sight-backend/internal/models"

	"github.com/google/uuid"
)

type IncomeSourceRepository struct {
	db *sql.DB
}

func NewIncomeSourceRepository(db *sql.DB) *IncomeSourceRepository {
	return &IncomeSourceRepository{db: db}
}

func (r *IncomeSourceRepository) GetAll(userID uuid.UUID) ([]models.IncomeSource, error) {
	query := `
		SELECT id, user_id, name, income_type, base_amount, bank_account, 
		       scheduled_year_month, is_active, created_at, updated_at
		FROM income_sources 
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sources []models.IncomeSource
	for rows.Next() {
		var source models.IncomeSource
		err := rows.Scan(
			&source.ID, &source.UserID, &source.Name, &source.IncomeType,
			&source.BaseAmount, &source.BankAccount, &source.ScheduledYearMonth,
			&source.IsActive, &source.CreatedAt, &source.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		sources = append(sources, source)
	}

	return sources, nil
}

func (r *IncomeSourceRepository) GetByID(id uuid.UUID) (*models.IncomeSource, error) {
	query := `
		SELECT id, user_id, name, income_type, base_amount, bank_account, 
		       scheduled_year_month, is_active, created_at, updated_at
		FROM income_sources 
		WHERE id = $1
	`

	var source models.IncomeSource
	err := r.db.QueryRow(query, id).Scan(
		&source.ID, &source.UserID, &source.Name, &source.IncomeType,
		&source.BaseAmount, &source.BankAccount, &source.ScheduledYearMonth,
		&source.IsActive, &source.CreatedAt, &source.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &source, nil
}

func (r *IncomeSourceRepository) GetActiveByUserID(userID uuid.UUID) ([]models.IncomeSource, error) {
	query := `
		SELECT id, user_id, name, income_type, base_amount, bank_account, 
		       scheduled_year_month, is_active, created_at, updated_at
		FROM income_sources 
		WHERE user_id = $1 AND is_active = true
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var sources []models.IncomeSource
	for rows.Next() {
		var source models.IncomeSource
		err := rows.Scan(
			&source.ID, &source.UserID, &source.Name, &source.IncomeType,
			&source.BaseAmount, &source.BankAccount, &source.ScheduledYearMonth,
			&source.IsActive, &source.CreatedAt, &source.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		sources = append(sources, source)
	}

	return sources, nil
}

func (r *IncomeSourceRepository) Create(source *models.IncomeSource) error {
	query := `
		INSERT INTO income_sources (id, user_id, name, income_type, base_amount, 
		                           bank_account, scheduled_year_month, is_active, 
		                           created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)
	`

	_, err := r.db.Exec(query,
		source.ID, source.UserID, source.Name, source.IncomeType,
		source.BaseAmount, source.BankAccount, source.ScheduledYearMonth,
		source.IsActive, source.CreatedAt, source.UpdatedAt,
	)

	return err
}

func (r *IncomeSourceRepository) Update(source *models.IncomeSource) error {
	query := `
		UPDATE income_sources 
		SET name = $2, income_type = $3, base_amount = $4, bank_account = $5,
		    scheduled_year_month = $6, is_active = $7, updated_at = $8
		WHERE id = $1
	`

	_, err := r.db.Exec(query,
		source.ID, source.Name, source.IncomeType, source.BaseAmount,
		source.BankAccount, source.ScheduledYearMonth, source.IsActive,
		source.UpdatedAt,
	)

	return err
}

func (r *IncomeSourceRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM income_sources WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
