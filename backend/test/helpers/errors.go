package helpers

import (
	"database/sql"
	"errors"
)

// Common test errors
var (
	ErrTestDatabaseConnection = errors.New("test database connection error")
	ErrTestNotFound           = errors.New("test record not found")
	ErrTestValidation         = errors.New("test validation error")
	ErrTestConflict           = errors.New("test conflict error")
	ErrTestUnauthorized       = errors.New("test unauthorized error")
	ErrTestPermissionDenied   = errors.New("test permission denied error")
	ErrTestInternalServer     = errors.New("test internal server error")
)

// IsDatabaseError checks if error is a database-related error
func IsDatabaseError(err error) bool {
	return errors.Is(err, sql.ErrConnDone) ||
		errors.Is(err, sql.ErrNoRows) ||
		errors.Is(err, sql.ErrTxDone) ||
		errors.Is(err, ErrTestDatabaseConnection)
}

// IsNotFoundError checks if error indicates a record was not found
func IsNotFoundError(err error) bool {
	return errors.Is(err, sql.ErrNoRows) ||
		errors.Is(err, ErrTestNotFound)
}

// IsValidationError checks if error is a validation error
func IsValidationError(err error) bool {
	return errors.Is(err, ErrTestValidation)
}

// IsConflictError checks if error indicates a conflict (e.g., duplicate key)
func IsConflictError(err error) bool {
	return errors.Is(err, ErrTestConflict)
}

// IsAuthorizationError checks if error is related to authorization
func IsAuthorizationError(err error) bool {
	return errors.Is(err, ErrTestUnauthorized) ||
		errors.Is(err, ErrTestPermissionDenied)
}
