package repositories

import (
	"database/sql"
	"flow-sight-backend/internal/models"

	"github.com/google/uuid"
)

type CardMonthlyTotalRepository struct {
	db *sql.DB
}

func NewCardMonthlyTotalRepository(db *sql.DB) *CardMonthlyTotalRepository {
	return &CardMonthlyTotalRepository{db: db}
}

func (r *CardMonthlyTotalRepository) GetByAssetID(assetID uuid.UUID) ([]models.CardMonthlyTotal, error) {
	query := `
		SELECT id, asset_id, year_month, total_amount, is_confirmed, created_at, updated_at
		FROM card_monthly_totals 
		WHERE asset_id = $1
		ORDER BY year_month DESC
	`

	rows, err := r.db.Query(query, assetID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var totals []models.CardMonthlyTotal
	for rows.Next() {
		var total models.CardMonthlyTotal
		err := rows.Scan(
			&total.ID, &total.AssetID, &total.YearMonth, &total.TotalAmount,
			&total.IsConfirmed, &total.CreatedAt, &total.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		totals = append(totals, total)
	}

	return totals, nil
}

func (r *CardMonthlyTotalRepository) GetByYearMonth(yearMonth string) ([]models.CardMonthlyTotal, error) {
	query := `
		SELECT id, asset_id, year_month, total_amount, is_confirmed, created_at, updated_at
		FROM card_monthly_totals 
		WHERE year_month = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, yearMonth)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var totals []models.CardMonthlyTotal
	for rows.Next() {
		var total models.CardMonthlyTotal
		err := rows.Scan(
			&total.ID, &total.AssetID, &total.YearMonth, &total.TotalAmount,
			&total.IsConfirmed, &total.CreatedAt, &total.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		totals = append(totals, total)
	}

	return totals, nil
}

func (r *CardMonthlyTotalRepository) GetByID(id uuid.UUID) (*models.CardMonthlyTotal, error) {
	query := `
		SELECT id, asset_id, year_month, total_amount, is_confirmed, created_at, updated_at
		FROM card_monthly_totals 
		WHERE id = $1
	`

	var total models.CardMonthlyTotal
	err := r.db.QueryRow(query, id).Scan(
		&total.ID, &total.AssetID, &total.YearMonth, &total.TotalAmount,
		&total.IsConfirmed, &total.CreatedAt, &total.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &total, nil
}

func (r *CardMonthlyTotalRepository) Create(total *models.CardMonthlyTotal) error {
	query := `
		INSERT INTO card_monthly_totals (id, asset_id, year_month, total_amount, 
		                                is_confirmed, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`

	_, err := r.db.Exec(query,
		total.ID, total.AssetID, total.YearMonth, total.TotalAmount,
		total.IsConfirmed, total.CreatedAt, total.UpdatedAt,
	)

	return err
}

func (r *CardMonthlyTotalRepository) Update(total *models.CardMonthlyTotal) error {
	query := `
		UPDATE card_monthly_totals 
		SET year_month = $2, total_amount = $3, is_confirmed = $4, updated_at = $5
		WHERE id = $1
	`

	_, err := r.db.Exec(query,
		total.ID, total.YearMonth, total.TotalAmount, total.IsConfirmed, total.UpdatedAt,
	)

	return err
}

func (r *CardMonthlyTotalRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM card_monthly_totals WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
