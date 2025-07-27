package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Soli0222/flow-sight/backend/internal/config"
	"github.com/Soli0222/flow-sight/backend/internal/models"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type AuthService struct {
	userRepo    UserRepositoryInterface
	config      *config.Config
	oauthConfig *oauth2.Config
}

type GoogleUserInfo struct {
	ID      string `json:"id"`
	Email   string `json:"email"`
	Name    string `json:"name"`
	Picture string `json:"picture"`
}

type Claims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func NewAuthService(userRepo UserRepositoryInterface, cfg *config.Config) *AuthService {
	oauthConfig := &oauth2.Config{
		ClientID:     cfg.OAuth.GoogleClientID,
		ClientSecret: cfg.OAuth.GoogleClientSecret,
		RedirectURL:  cfg.OAuth.RedirectURL,
		Scopes:       []string{"openid", "profile", "email"},
		Endpoint:     google.Endpoint,
	}

	return &AuthService{
		userRepo:    userRepo,
		config:      cfg,
		oauthConfig: oauthConfig,
	}
}

func (s *AuthService) GetGoogleAuthURL(state string) string {
	return s.oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline)
}

func (s *AuthService) HandleGoogleCallback(code string) (*models.User, string, error) {
	// Exchange code for token
	token, err := s.oauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, "", fmt.Errorf("failed to exchange token: %w", err)
	}

	// Get user info from Google
	client := s.oauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return nil, "", fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	var googleUser GoogleUserInfo
	if err := json.NewDecoder(resp.Body).Decode(&googleUser); err != nil {
		return nil, "", fmt.Errorf("failed to decode user info: %w", err)
	}

	// Check if user exists by Google ID first
	user, err := s.userRepo.GetByGoogleID(googleUser.ID)
	if err != nil && err != sql.ErrNoRows {
		return nil, "", fmt.Errorf("failed to get user by google id: %w", err)
	}

	// If user doesn't exist by Google ID, check by email
	if user == nil {
		user, err = s.userRepo.GetByEmail(googleUser.Email)
		if err != nil && err != sql.ErrNoRows {
			return nil, "", fmt.Errorf("failed to get user by email: %w", err)
		}

		// If user exists by email but doesn't have Google ID, update it
		if user != nil && user.GoogleID == "" {
			user.GoogleID = googleUser.ID
			user.Name = googleUser.Name
			user.Picture = googleUser.Picture
			if err := s.userRepo.Update(user); err != nil {
				return nil, "", fmt.Errorf("failed to update user: %w", err)
			}
		}
	}

	// If user doesn't exist, create new user
	if user == nil {
		user = &models.User{
			Email:    googleUser.Email,
			Name:     googleUser.Name,
			Picture:  googleUser.Picture,
			GoogleID: googleUser.ID,
		}
		if err := s.userRepo.Create(user); err != nil {
			return nil, "", fmt.Errorf("failed to create user: %w", err)
		}

		// Verify user was created successfully by re-fetching
		user, err = s.userRepo.GetByGoogleID(googleUser.ID)
		if err != nil {
			return nil, "", fmt.Errorf("failed to verify user creation: %w", err)
		}
	}

	// Generate JWT token
	jwtToken, err := s.GenerateJWT(user)
	if err != nil {
		return nil, "", fmt.Errorf("failed to generate jwt: %w", err)
	}

	return user, jwtToken, nil
}

func (s *AuthService) GenerateJWT(user *models.User) (string, error) {
	claims := &Claims{
		UserID: user.ID.String(),
		Email:  user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWT.Secret))
}

func (s *AuthService) ValidateJWT(tokenString string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.config.JWT.Secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, errors.New("invalid token")
}

func (s *AuthService) GetUserByID(userID string) (*models.User, error) {
	// Parse UUID string directly
	uid, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user id: %w", err)
	}

	user, err := s.userRepo.GetByID(uid)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("user not found with id %s: %w", userID, err)
		}
		return nil, fmt.Errorf("failed to get user by id %s: %w", userID, err)
	}

	return user, nil
}
