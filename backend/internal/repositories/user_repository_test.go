package repositories

import (
	"database/sql"
	"flow-sight-backend/test/helpers"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestUserRepository_GetByEmail(t *testing.T) {
	db, mock := helpers.SetupMockDB(t)
	defer helpers.TeardownMockDB(db)

	repo := NewUserRepository(db)
	userID := uuid.New()
	email := "test@example.com"

	tests := []struct {
		name          string
		email         string
		setupMock     func(sqlmock.Sqlmock)
		expectedFound bool
		expectedError bool
	}{
		{
			name:  "user found by email",
			email: email,
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "email", "name", "picture", "google_id", "password", "created_at", "updated_at",
				}).
					AddRow(
						userID, email, "Test User", "https://example.com/pic.jpg",
						"google-123", "hashed-password", time.Now(), time.Now(),
					)

				mock.ExpectQuery(`SELECT id, email, name, picture, google_id, password, created_at, updated_at FROM users WHERE email = \$1`).
					WithArgs(email).
					WillReturnRows(rows)
			},
			expectedFound: true,
			expectedError: false,
		},
		{
			name:  "user not found",
			email: email,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, email, name, picture, google_id, password, created_at, updated_at FROM users WHERE email = \$1`).
					WithArgs(email).
					WillReturnError(sql.ErrNoRows)
			},
			expectedFound: false,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(mock)

			user, err := repo.GetByEmail(tt.email)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.email, user.Email)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_GetByGoogleID(t *testing.T) {
	db, mock := helpers.SetupMockDB(t)
	defer helpers.TeardownMockDB(db)

	repo := NewUserRepository(db)
	userID := uuid.New()
	googleID := "google-123"

	tests := []struct {
		name          string
		googleID      string
		setupMock     func(sqlmock.Sqlmock)
		expectedFound bool
		expectedError bool
	}{
		{
			name:     "user found by google ID",
			googleID: googleID,
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "email", "name", "picture", "google_id", "password", "created_at", "updated_at",
				}).
					AddRow(
						userID, "test@example.com", "Test User", "https://example.com/pic.jpg",
						googleID, "hashed-password", time.Now(), time.Now(),
					)

				mock.ExpectQuery(`SELECT id, email, name, picture, google_id, password, created_at, updated_at FROM users WHERE google_id = \$1`).
					WithArgs(googleID).
					WillReturnRows(rows)
			},
			expectedFound: true,
			expectedError: false,
		},
		{
			name:     "user not found",
			googleID: googleID,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, email, name, picture, google_id, password, created_at, updated_at FROM users WHERE google_id = \$1`).
					WithArgs(googleID).
					WillReturnError(sql.ErrNoRows)
			},
			expectedFound: false,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(mock)

			user, err := repo.GetByGoogleID(tt.googleID)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.googleID, user.GoogleID)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_GetByID(t *testing.T) {
	db, mock := helpers.SetupMockDB(t)
	defer helpers.TeardownMockDB(db)

	repo := NewUserRepository(db)
	userID := uuid.New()

	tests := []struct {
		name          string
		userID        uuid.UUID
		setupMock     func(sqlmock.Sqlmock)
		expectedFound bool
		expectedError bool
	}{
		{
			name:   "user found by ID",
			userID: userID,
			setupMock: func(mock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{
					"id", "email", "name", "picture", "google_id", "password", "created_at", "updated_at",
				}).
					AddRow(
						userID, "test@example.com", "Test User", "https://example.com/pic.jpg",
						"google-123", "hashed-password", time.Now(), time.Now(),
					)

				mock.ExpectQuery(`SELECT id, email, name, picture, google_id, password, created_at, updated_at FROM users WHERE id = \$1`).
					WithArgs(userID).
					WillReturnRows(rows)
			},
			expectedFound: true,
			expectedError: false,
		},
		{
			name:   "user not found",
			userID: userID,
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectQuery(`SELECT id, email, name, picture, google_id, password, created_at, updated_at FROM users WHERE id = \$1`).
					WithArgs(userID).
					WillReturnError(sql.ErrNoRows)
			},
			expectedFound: false,
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.setupMock(mock)

			user, err := repo.GetByID(tt.userID)

			if tt.expectedError {
				assert.Error(t, err)
				assert.Nil(t, user)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, user)
				assert.Equal(t, tt.userID, user.ID)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_Create(t *testing.T) {
	db, mock := helpers.SetupMockDB(t)
	defer helpers.TeardownMockDB(db)

	repo := NewUserRepository(db)

	tests := []struct {
		name          string
		setupMock     func(sqlmock.Sqlmock)
		expectedError bool
	}{
		{
			name: "successful creation",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO users \(id, email, name, picture, google_id, password, created_at, updated_at\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8\)`).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: false,
		},
		{
			name: "database error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`INSERT INTO users \(id, email, name, picture, google_id, password, created_at, updated_at\) VALUES \(\$1, \$2, \$3, \$4, \$5, \$6, \$7, \$8\)`).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(sql.ErrConnDone)
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := helpers.CreateTestUser()
			tt.setupMock(mock)

			err := repo.Create(user)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}

func TestUserRepository_Update(t *testing.T) {
	db, mock := helpers.SetupMockDB(t)
	defer helpers.TeardownMockDB(db)

	repo := NewUserRepository(db)

	tests := []struct {
		name          string
		setupMock     func(sqlmock.Sqlmock)
		expectedError bool
	}{
		{
			name: "successful update",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE users SET email = \$2, name = \$3, picture = \$4, google_id = \$5, password = \$6, updated_at = \$7 WHERE id = \$1`).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnResult(sqlmock.NewResult(1, 1))
			},
			expectedError: false,
		},
		{
			name: "database error",
			setupMock: func(mock sqlmock.Sqlmock) {
				mock.ExpectExec(`UPDATE users SET email = \$2, name = \$3, picture = \$4, google_id = \$5, password = \$6, updated_at = \$7 WHERE id = \$1`).
					WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg()).
					WillReturnError(sql.ErrConnDone)
			},
			expectedError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			user := helpers.CreateTestUser()
			user.Name = "Updated Name"
			tt.setupMock(mock)

			err := repo.Update(user)

			if tt.expectedError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}

			assert.NoError(t, mock.ExpectationsWereMet())
		})
	}
}
