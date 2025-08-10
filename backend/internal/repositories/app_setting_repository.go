package repositories

import (
	"database/sql"

	"github.com/Soli0222/flow-sight/backend/internal/models"
)

type AppSettingRepository struct {
	db *sql.DB
}

func NewAppSettingRepository(db *sql.DB) *AppSettingRepository {
	return &AppSettingRepository{db: db}
}

func (r *AppSettingRepository) GetAll() ([]models.AppSetting, error) {
	query := `
		SELECT id, key, value, created_at, updated_at
		FROM app_settings 
		ORDER BY key ASC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return []models.AppSetting{}, err
	}
	defer rows.Close()

	settings := make([]models.AppSetting, 0)
	for rows.Next() {
		var setting models.AppSetting
		err := rows.Scan(
			&setting.ID, &setting.Key, &setting.Value,
			&setting.CreatedAt, &setting.UpdatedAt,
		)
		if err != nil {
			return []models.AppSetting{}, err
		}
		settings = append(settings, setting)
	}

	return settings, nil
}

func (r *AppSettingRepository) GetByKey(key string) (*models.AppSetting, error) {
	query := `
		SELECT id, key, value, created_at, updated_at
		FROM app_settings 
		WHERE key = $1
	`

	var setting models.AppSetting
	err := r.db.QueryRow(query, key).Scan(
		&setting.ID, &setting.Key, &setting.Value,
		&setting.CreatedAt, &setting.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &setting, nil
}

func (r *AppSettingRepository) Upsert(setting *models.AppSetting) error {
	query := `
		INSERT INTO app_settings (id, key, value, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5)
		ON CONFLICT (key) 
		DO UPDATE SET value = EXCLUDED.value, updated_at = EXCLUDED.updated_at
	`

	_, err := r.db.Exec(query,
		setting.ID, setting.Key, setting.Value,
		setting.CreatedAt, setting.UpdatedAt,
	)

	return err
}

func (r *AppSettingRepository) Delete(key string) error {
	query := `DELETE FROM app_settings WHERE key = $1`
	_, err := r.db.Exec(query, key)
	return err
}
