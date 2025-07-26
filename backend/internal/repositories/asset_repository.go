package repositories

import (
	"database/sql"
	"flow-sight-backend/internal/models"

	"github.com/google/uuid"
)

type AssetRepository struct {
	db *sql.DB
}

func NewAssetRepository(db *sql.DB) *AssetRepository {
	return &AssetRepository{db: db}
}

func (r *AssetRepository) GetAll(userID uuid.UUID) ([]models.Asset, error) {
	query := `
		SELECT id, user_id, name, asset_type, closing_day, payment_day, bank_account, created_at, updated_at
		FROM assets 
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var assets []models.Asset
	for rows.Next() {
		var asset models.Asset
		err := rows.Scan(
			&asset.ID, &asset.UserID, &asset.Name, &asset.AssetType,
			&asset.ClosingDay, &asset.PaymentDay, &asset.BankAccount,
			&asset.CreatedAt, &asset.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		assets = append(assets, asset)
	}

	return assets, nil
}

func (r *AssetRepository) GetByID(id uuid.UUID) (*models.Asset, error) {
	query := `
		SELECT id, user_id, name, asset_type, closing_day, payment_day, bank_account, created_at, updated_at
		FROM assets 
		WHERE id = $1
	`

	var asset models.Asset
	err := r.db.QueryRow(query, id).Scan(
		&asset.ID, &asset.UserID, &asset.Name, &asset.AssetType,
		&asset.ClosingDay, &asset.PaymentDay, &asset.BankAccount,
		&asset.CreatedAt, &asset.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &asset, nil
}

func (r *AssetRepository) Create(asset *models.Asset) error {
	query := `
		INSERT INTO assets (id, user_id, name, asset_type, closing_day, payment_day, bank_account, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`

	_, err := r.db.Exec(query,
		asset.ID, asset.UserID, asset.Name, asset.AssetType,
		asset.ClosingDay, asset.PaymentDay, asset.BankAccount,
		asset.CreatedAt, asset.UpdatedAt,
	)

	return err
}

func (r *AssetRepository) Update(asset *models.Asset) error {
	query := `
		UPDATE assets 
		SET name = $2, asset_type = $3, closing_day = $4, payment_day = $5, 
		    bank_account = $6, updated_at = $7
		WHERE id = $1
	`

	_, err := r.db.Exec(query,
		asset.ID, asset.Name, asset.AssetType, asset.ClosingDay,
		asset.PaymentDay, asset.BankAccount, asset.UpdatedAt,
	)

	return err
}

func (r *AssetRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM assets WHERE id = $1`
	_, err := r.db.Exec(query, id)
	return err
}
