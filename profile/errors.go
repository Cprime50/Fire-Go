package profile

import "errors"

var (
	ErrQuoteNotFound             = errors.New("quote not found")
	ErrProfileNotFound           = errors.New("profile not found")
	ErrInvalidRequestBody        = errors.New("invalid request body")
	ErrNotAuthorized             = errors.New("unauthorized access")
	ErrCreateProfile             = errors.New("failed to create profile")
	ErrUpdateProfile             = errors.New("failed to update profile")
	ErrDeletingProfile           = errors.New("failed to delete profile")
	ErrGettingProfile            = errors.New("failed to get profile")
	ErrProfileAlreadyExists      = errors.New("profile already exists")
	ErrForeignKeyViolation       = errors.New("foreign key violation")
	ErrUniqueConstraintViolation = errors.New("unique constraint violation")
)
