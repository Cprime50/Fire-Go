package quote

import "errors"

var (
	ErrQuoteNotFound             = errors.New("quote not found")
	ErrquoteNotFound             = errors.New("quote not found")
	ErrInvalidRequestBody        = errors.New("invalid request body")
	ErrNotAuthorized             = errors.New("unauthorized access")
	ErrCreateQuote               = errors.New("failed to create quote")
	ErrUpdateQuote               = errors.New("failed to update quote")
	ErrDeletingQuote             = errors.New("failed to delete quote")
	ErrGettingQuote              = errors.New("failed to get quote")
	ErrApprovingQuote            = errors.New("error approving quote")
	ErrquoteAlreadyExists        = errors.New("quote already exists")
	ErrForeignKeyViolation       = errors.New("foreign key violation")
	ErrUniqueConstraintViolation = errors.New("unique constraint violation")
)
