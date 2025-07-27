package repositories

import (
	"database/sql"
	"flow-sight-backend/internal/models"
	"flow-sight-backend/test/helpers"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestAppSettingRepository_GetByUserID(t *testing.T) {
	db, mock := helpers.SetupMockDB(t)
	defer helpers.TeardownMockDB(db)

	repo := NewAppSettingRepository(db)
	userID := uuid.New()

	tests := []struct {
		name          string
		userID        uuid.UUID
		setupMock     func(sqlmock.Sqlmock)
		expectedCount int
		expectedError bool
	}{
		{
			name:   "successful retrieval",
			userID: userID,
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "user_id", "key", "value", "created_at", "updated_at",
				}).
					AddRow(
						uuid.New(), userID, "currency", "JPY", time.Now(), time.Now(),
					).
					AddRow(
						uuid.New(), userID, "timezone", "Asia/Tokyo", time.Now(), time.Now(),
					)

				mock.ExpectQuery(`SELECT id, user_id, key, value, created_at, updated_at FROM app_settings WHERE user_id = \$1 ORDER BY key ASC`).
					WithArgs(userID).
					WillReturnRows(rows)
			},
			expectedCount: 2,
			expectedError: false,
		},
		{
			name:   "no settings found",
			userID: userID,
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "user_id", "key", "value", "created_at", "updated_at",
				})

				mock.ExpectQuery(`SELECT id, user_id, key, value, created_at, updated_at FROM app_settings WHERE user_id = \$1 ORDER BY key ASC`).
					WithArgs(userID).
					WillReturnRows(rows)
			},
			expectedCount: 0,
			expectedError: false,
		},
		{
			name:   "database error",
			userID: userID,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, user_id, key, value, created_at, updated_at FROM app_settings WHERE user_id = \$1 ORDER BY key ASC`).
					WithArgs(userID).
					WillReturnError(sql.ErrConnDone)
			},
			expectedCount: 0,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(mock)

			settings, err := repo.GetByUserID(tt.userID)

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
	userID := uuid.New()
	settingID := uuid.New()
	key := "currency"

	tests := []struct {
		name          string
		userID        uuid.UUID
		key           string
		setupMock     func(sqlmock.Sqlmock)
		expectedFound bool
		expectedError bool
	}{
		{
			name:   "setting found",
			userID: userID,
			key:    key,
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "user_id", "key", "value", "created_at", "updated_at",
				}).
					AddRow(
						settingID, userID, key, "JPY", time.Now(), time.Now(),
					)

				mock.ExpectQuery(`SELECT id, user_id, key, value, created_at, updated_at FROM app_settings WHERE user_id = \$1 AND key = \$2`).
					WithArgs(userID, key).
					WillReturnRows(rows)
			},
			expectedFound: true,
			expectedError: false,
		},
		{
			name:   "setting not found",
			userID: userID,
			key:    key,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, user_id, key, value, created_at, updated_at FROM app_settings WHERE user_id = \$1 AND key = \$2`).
					WithArgs(userID, key).
					WillReturnError(sql.ErrNoRows)
			},
			expectedFound: false,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(mock)

			setting, err := repo.GetByKey(tt.userID, tt.key)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, setting)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, setting)
				assert.Equal(t, tt.key, setting.Key)
				assert.Equal(t, tt.userID, setting.UserID)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestAppSettingRepository_Upsert(t *testing.T) {
	db, mock := helpers.SetupMockDB(t)
	defer helpers.TeardownMockDB(db)

	repo := NewAppSettingRepository(db)
	userID := uuid.New()
	settingID := uuid.New()

	tests := []struct {
		name          string
		setupMock     func(sqlmock.Sqlmock)
		expectedError bool
	}{
		{
			name: "successful upsert",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO app_settings \(id, user_id, key, value, created_at, updated_at\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6\) ON CONFLICT \(user_id, key\) DO UPDATE SET value = EXCLUDED\.value, updated_at = EXCLUDED\.updated_at`).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: false,
		},
		{
			name: "database error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO app_settings \(id, user_id, key, value, created_at, updated_at\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6\) ON CONFLICT \(user_id, key\) DO UPDATE SET value = EXCLUDED\.value, updated_at = EXCLUDED\.updated_at`).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(sql.ErrConnDone)
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			appSetting := &models.AppSetting{
				ID:        settingID,
				UserID:    userID,
				Key:       "currency",
				Value:     "JPY",
				CreatedAt: time.Now(),
				UpdatedAt: time.Now(),
			}

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
