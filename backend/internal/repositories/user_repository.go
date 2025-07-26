package repositories

import (
	"database/sql"
	"flow-sight-backend/internal/models"
	"time"

	"github.com/google/uuid"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetByEmail(email string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, name, picture, google_id, password, created_at, updated_at
		FROM users
		WHERE email = $1
	`
	err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Picture,
		&user.GoogleID,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByGoogleID(googleID string) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, name, picture, google_id, password, created_at, updated_at
		FROM users
		WHERE google_id = $1
	`
	err := r.db.QueryRow(query, googleID).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Picture,
		&user.GoogleID,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByID(id uuid.UUID) (*models.User, error) {
	user := &models.User{}
	query := `
		SELECT id, email, name, picture, google_id, password, created_at, updated_at
		FROM users
		WHERE id = $1
	`
	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Email,
		&user.Name,
		&user.Picture,
		&user.GoogleID,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) Create(user *models.User) error {
	user.ID = uuid.New()
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()

	query := `
		INSERT INTO users (id, email, name, picture, google_id, password, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.Exec(query,
		user.ID,
		user.Email,
		user.Name,
		user.Picture,
		user.GoogleID,
		user.Password,
		user.CreatedAt,
		user.UpdatedAt,
	)
	return err
}

func (r *UserRepository) Update(user *models.User) error {
	user.UpdatedAt = time.Now()

	query := `
		UPDATE users
		SET email = $2, name = $3, picture = $4, google_id = $5, password = $6, updated_at = $7
		WHERE id = $1
	`
	_, err := r.db.Exec(query,
		user.ID,
		user.Email,
		user.Name,
		user.Picture,
		user.GoogleID,
		user.Password,
		user.UpdatedAt,
	)
	return err
}
