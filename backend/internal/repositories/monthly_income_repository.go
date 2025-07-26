package repositories

import (
	"database/sql"
	"flow-sight-backend/internal/models"

	"github.com/google/uuid"
)

type MonthlyIncomeRepository struct {
	db *sql.DB
}

func NewMonthlyIncomeRepository(db *sql.DB) *MonthlyIncomeRepository {
	return &MonthlyIncomeRepository{db: db}
}

func (r *MonthlyIncomeRepository) GetByIncomeSourceID(incomeSourceID uuid.UUID) ([]models.MonthlyIncomeRecord, error) {
	query := `
		SELECT id, income_source_id, year_month, actual_amount, is_confirmed, note, created_at, updated_at
		FROM monthly_income_records 
		WHERE income_source_id = $1
		ORDER BY year_month DESC
	`

	rows, err := r.db.Query(query, incomeSourceID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []models.MonthlyIncomeRecord
	for rows.Next() {
		var record models.MonthlyIncomeRecord
		err := rows.Scan(
			&record.ID, &record.IncomeSourceID, &record.YearMonth,
			&record.ActualAmount, &record.IsConfirmed, &record.Note,
			&record.CreatedAt, &record.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return records, nil
}

func (r *MonthlyIncomeRepository) GetByYearMonth(yearMonth string) ([]models.MonthlyIncomeRecord, error) {
	query := `
		SELECT id, income_source_id, year_month, actual_amount, is_confirmed, note, created_at, updated_at
		FROM monthly_income_records 
		WHERE year_month = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, yearMonth)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var records []models.MonthlyIncomeRecord
	for rows.Next() {
		var record models.MonthlyIncomeRecord
		err := rows.Scan(
			&record.ID, &record.IncomeSourceID, &record.YearMonth,
			&record.ActualAmount, &record.IsConfirmed, &record.Note,
			&record.CreatedAt, &record.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		records = append(records, record)
	}

	return records, nil
}

func (r *MonthlyIncomeRepository) GetByID(id uuid.UUID) (*models.MonthlyIncomeRecord, error) {
	query := `
		SELECT id, income_source_id, year_month, actual_amount, is_confirmed, note, created_at, updated_at
		FROM monthly_income_records 
		WHERE id = $1
	`

	var record models.MonthlyIncomeRecord
	err := r.db.QueryRow(query, id).Scan(
		&record.ID, &record.IncomeSourceID, &record.YearMonth,
		&record.ActualAmount, &record.IsConfirmed, &record.Note,
		&record.CreatedAt, &record.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &record, nil
}

func (r *MonthlyIncomeRepository) Create(record *models.MonthlyIncomeRecord) error {
	query := `
		INSERT INTO monthly_income_records (id, income_source_id, year_month, actual_amount, 
		                                   is_confirmed, note, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`

	_, err := r.db.Exec(query,
		record.ID, record.IncomeSourceID, record.YearMonth, record.ActualAmount,
		record.IsConfirmed, record.Note, record.CreatedAt, record.UpdatedAt,
	)

	return err
}

func (r *MonthlyIncomeRepository) Update(record *models.MonthlyIncomeRecord) error {
	query := `
		UPDATE monthly_income_records 
		SET year_month = $2, actual_amount = $3, is_confirmed = $4, note = $5, updated_at = $6
		WHERE id = $1
	`

	_, err := r.db.Exec(query,
		record.ID, record.YearMonth, record.ActualAmount, record.IsConfirmed,
		record.Note, record.UpdatedAt,
	)

	return err
}

func (r *MonthlyIncomeRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM monthly_income_records WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
