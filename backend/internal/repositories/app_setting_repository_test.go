package repositories

import (
	"database/sql"
	"testing"
	"time"

	"github.com/Soli0222/flow-sight/backend/internal/models"
	"github.com/Soli0222/flow-sight/backend/test/helpers"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAppSettingRepository_GetAll(t *testing.T) {
	db, mock := helpers.SetupMockDB(t)
	defer helpers.TeardownMockDB(db)

	repo := NewAppSettingRepository(db)

	tests := []struct {
		name          string
		setupMock     func(sqlmock.Sqlmock)
		expectedCount int
		expectedError bool
	}{
		{
			name: "successful retrieval",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "key", "value", "created_at", "updated_at",
				}).
					AddRow(
						uuid.New(), "currency", "JPY", time.Now(), time.Now(),
					).
					AddRow(
						uuid.New(), "timezone", "Asia/Tokyo", time.Now(), time.Now(),
					)

				mock.ExpectQuery(`SELECT id, key, value, created_at, updated_at FROM app_settings ORDER BY key ASC`).
					WillReturnRows(rows)
			},
			expectedCount: 2,
			expectedError: false,
		},
		{
			name: "no settings found",
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "key", "value", "created_at", "updated_at",
				})

				mock.ExpectQuery(`SELECT id, key, value, created_at, updated_at FROM app_settings ORDER BY key ASC`).
					WillReturnRows(rows)
			},
			expectedCount: 0,
			expectedError: false,
		},
		{
			name: "database error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, key, value, created_at, updated_at FROM app_settings ORDER BY key ASC`).
					WillReturnError(sql.ErrConnDone)
			},
			expectedCount: 0,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			tt.setupMock(mock)

			settings, err := repo.GetAll()

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.Len(t, settings, tt.expectedCount)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAppSettingRepository_GetByKey(t *testing.T) {
	db, mock := helpers.SetupMockDB(t)
	defer helpers.TeardownMockDB(db)

	repo := NewAppSettingRepository(db)
	settingID := uuid.New()
	key := "currency"

	tests := []struct {
		name          string
		key           string
		setupMock     func(sqlmock.Sqlmock)
		expectedFound bool
		expectedError bool
	}{
		{
			name: "setting found",
			key:  key,
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "key", "value", "created_at", "updated_at",
				}).
					AddRow(
						settingID, key, "JPY", time.Now(), time.Now(),
					)

				mock.ExpectQuery(`SELECT id, key, value, created_at, updated_at FROM app_settings WHERE key = \$1`).
					WithArgs(key).
					WillReturnRows(rows)
			},
			expectedFound: true,
			expectedError: false,
		},
		{
			name: "setting not found",
			key:  key,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, key, value, created_at, updated_at FROM app_settings WHERE key = \$1`).
					WithArgs(key).
					WillReturnError(sql.ErrNoRows)
			},
			expectedFound: false,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt := tt
			tt.setupMock(mock)

			setting, err := repo.GetByKey(tt.key)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, setting)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, setting)
				assert.Equal(t, tt.key, setting.Key)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAppSettingRepository_Upsert(t *testing.T) {
	db, mock := helpers.SetupMockDB(t)
	defer helpers.TeardownMockDB(db)

	repo := NewAppSettingRepository(db)
	settingID := uuid.New()

	tests := []struct {
		name          string
		setupMock     func(sqlmock.Sqlmock)
		expectedError bool
	}{
		{
			name: "successful upsert",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO app_settings \(id, key, value, created_at, updated_at\) VALUES \(\$1, \$2, \$3, \$4, \$5\) ON CONFLICT \(key\) DO UPDATE SET value = EXCLUDED\.value, updated_at = EXCLUDED\.updated_at`).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: false,
		},
		{
			name: "database error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO app_settings \(id, key, value, created_at, updated_at\) VALUES \(\$1, \$2, \$3, \$4, \$5\) ON CONFLICT \(key\) DO UPDATE SET value = EXCLUDED\.value, updated_at = EXCLUDED\.updated_at`).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(sql.ErrConnDone)
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appSetting := &models.AppSetting{
				ID:        settingID,
				Key:       "currency",
				Value:     "JPY",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

			tt := tt
			tt.setupMock(mock)

			err := repo.Upsert(appSetting)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
