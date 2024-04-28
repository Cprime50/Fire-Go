package profile

import "errors"

var (
	ErrQuoteNotFound             = errors.New("quote not found")
	ErrProfileNotFound           = errors.New("profile not found")
	ErrInvalidRequestBody        = errors.New("invalid request body")
	ErrCreateProfile             = errors.New("failed to create profile")
	ErrUpdateProfile             = errors.New("failed to update profile")
	ErrProfileAlreadyExists      = errors.New("profile already exists")
	ErrForeignKeyViolation       = errors.New("foreign key violation")
	ErrUniqueConstraintViolation = errors.New("unique constraint violation")
)
