package repositories

import (
	"database/sql"
	"github.com/Soli0222/flow-sight/backend/internal/models"

	"github.com/google/uuid"
)

type AppSettingRepository struct {
	db *sql.DB
}

func NewAppSettingRepository(db *sql.DB) *AppSettingRepository {
	return &AppSettingRepository{db: db}
}

func (r *AppSettingRepository) GetByUserID(userID uuid.UUID) ([]models.AppSetting, error) {
	query := `
		SELECT id, user_id, key, value, created_at, updated_at
		FROM app_settings 
		WHERE user_id = $1
		ORDER BY key ASC
	`

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return []models.AppSetting{}, err
	}
	defer rows.Close()

	settings := make([]models.AppSetting, 0)
	for rows.Next() {
		var setting models.AppSetting
		err := rows.Scan(
			&setting.ID, &setting.UserID, &setting.Key, &setting.Value,
			&setting.CreatedAt, &setting.UpdatedAt,
		)
		if err != nil {
			return []models.AppSetting{}, err
		}
		settings = append(settings, setting)
	}

	return settings, nil
}

func (r *AppSettingRepository) GetByKey(userID uuid.UUID, key string) (*models.AppSetting, error) {
	query := `
		SELECT id, user_id, key, value, created_at, updated_at
		FROM app_settings 
		WHERE user_id = $1 AND key = $2
	`

	var setting models.AppSetting
	err := r.db.QueryRow(query, userID, key).Scan(
		&setting.ID, &setting.UserID, &setting.Key, &setting.Value,
		&setting.CreatedAt, &setting.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &setting, nil
}

func (r *AppSettingRepository) Upsert(setting *models.AppSetting) error {
	query := `
		INSERT INTO app_settings (id, user_id, key, value, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		ON CONFLICT (user_id, key) 
		DO UPDATE SET value = EXCLUDED.value, updated_at = EXCLUDED.updated_at
	`

	_, err := r.db.Exec(query,
		setting.ID, setting.UserID, setting.Key, setting.Value,
		setting.CreatedAt, setting.UpdatedAt,
	)

	return err
}

func (r *AppSettingRepository) Delete(userID uuid.UUID, key string) error {
	query := `DELETE FROM app_settings WHERE user_id = $1 AND key = $2`
	_, err := r.db.Exec(query, userID, key)
	return err
}
