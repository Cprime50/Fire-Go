package quote

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/cprime50/fire-go/db"
	"github.com/google/uuid"
)

type Quote struct {
	Id        string    `json:"id"`
	UserId    string    `json:"user_id"`
	Quote     string    `json:"quote"`
	Approved  bool      `json:"approved"`
	CreatedAt time.Time `json:"created_at"`
}

var (
	ErrQuoteNotFound             = errors.New("quote not found")
	ErrProfileNotFound           = errors.New("profile not found")
	ErrDuplicateEntry            = errors.New("duplicate entry")
	ErrForeignKeyViolation       = errors.New("foreign key violation")
	ErrUniqueConstraintViolation = errors.New("unique constraint violation")
)

func createQuote(quote *Quote) error {
	id, err := uuid.NewRandom()
	if err != nil {
		return fmt.Errorf("CreateQuote uuid.NewRandom: %w", err)
	}
	_, err = db.Db.Exec(
		"INSERT INTO quotes (id, user_id, quote, approved, created_at) VALUES ($1, $2, $3, $4, $5)",
		id.String(),
		quote.UserId,
		quote.Quote,
		quote.Approved,
		time.Now(),
	)
	if err != nil {
		return fmt.Errorf("CreateQuote error: %w", err)
	}
	return nil
}

func updateQuote(quote *Quote) error {
	_, err := db.Db.Exec(
		"UPDATE quotes SET quote = $1, approved = FALSE, updated_at = $2 WHERE id = $3",
		quote.Quote,
		time.Now(),
		quote.Id,
	)
	if err != nil {
		return fmt.Errorf("UpdateQuote error: %w", err)
	}
	return nil
}

// deletequote
func deleteQuote(quoteId string) error {
	_, err := db.Db.Exec(
		"DELETE FROM quotes WHERE id = $1",
		quoteId,
	)
	if err != nil {
		return fmt.Errorf("DeleteQuote error: %w", err)
	}
	return nil
}

func getQuoteById(quoteId string) (*Quote, error) {
	var quote Quote
	err := db.Db.QueryRow(
		"SELECT * FROM quotes WHERE id = $1",
		quoteId,
	).Scan(
		&quote.Id,
		&quote.UserId,
		&quote.Quote,
		&quote.Approved,
		&quote.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrQuoteNotFound
		}
		return nil, fmt.Errorf("GetQuoteById error: %w", err)
	}
	return &quote, nil
}

// GetQuotesByProfileId retrieves a user quote by user ID.
func getQuotesByUserId(userId string) ([]*Quote, error) {
	rows, err := db.Db.Query("SELECT * FROM quotes WHERE user_id = $1", userId)
	if err != nil {
		return nil, fmt.Errorf("db.Query: %w", err)
	}
	defer rows.Close()
	quotes, _ := queryQuotes(rows)
	return quotes, nil
}

func getAllQuotes() ([]*Quote, error) {
	rows, err := db.Db.Query("SELECT * FROM quotes")
	if err != nil {
		return nil, fmt.Errorf("db.Query: %w", err)
	}
	defer rows.Close()

	quotes, err := queryQuotes(rows)
	if err != nil {
		return nil, fmt.Errorf("queryQuotes: %w", err)
	}
	if len(quotes) == 0 {
		return nil, ErrQuoteNotFound
	}
	return quotes, nil
}

func getAllApprovedQuotes() ([]*Quote, error) {
	rows, err := db.Db.Query("SELECT * FROM quotes WHERE approved = true")
	if err != nil {
		return nil, fmt.Errorf("db.Query: %w", err)
	}
	defer rows.Close()
	quotes, err := queryQuotes(rows)
	if err != nil {
		return nil, fmt.Errorf("queryQuotes: %w", err)
	}
	if len(quotes) == 0 {
		return nil, ErrQuoteNotFound
	}
	return quotes, nil
}

func getApprovedQuotesByUserId(userId string) ([]*Quote, error) {
	rows, err := db.Db.Query("SELECT * FROM quotes WHERE user_id = $1 AND approved = true", userId)
	if err != nil {
		return nil, fmt.Errorf("db.Query: %w", err)
	}
	defer rows.Close()
	quotes, err := queryQuotes(rows)
	if err != nil {
		return nil, fmt.Errorf("queryQuotes: %w", err)
	}
	if len(quotes) == 0 {
		return nil, ErrQuoteNotFound
	}
	return quotes, nil
}

func approveQuote(quoteId string) error {
	result, err := db.Db.Exec("UPDATE quotes SET approved = TRUE WHERE id = $1", quoteId)
	if err != nil {
		return fmt.Errorf("ApproveQuote error: %w", err)
	}
	rowsAffected, _ := result.RowsAffected()
	if rowsAffected == 0 {
		return ErrQuoteNotFound
	}
	return nil
}

func queryQuotes(rows *sql.Rows) ([]*Quote, error) {
	var quotes []*Quote
	for rows.Next() {
		quote := &Quote{}
		quote.CreatedAt = time.Now()

		if err := rows.Scan(&quote.Id, &quote.UserId, &quote.Quote, &quote.Approved, &quote.CreatedAt); err != nil {
			return nil, fmt.Errorf("rows.Scan: %w", err)
		}
		quotes = append(quotes, quote)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows.Err: %w", err)
	}
	return quotes, nil
}

// Get UnapprovedQuote
func getUnapprovedQuotes() ([]*Quote, error) {
	rows, err := db.Db.Query("SELECT * FROM quotes WHERE approved = FALSE")
	if err != nil {
		return nil, fmt.Errorf("db.Query: %w", err)
	}
	defer rows.Close()
	quotes, err := queryQuotes(rows)
	if err != nil {
		return nil, fmt.Errorf("queryQuotes: %w", err)
	}
	if len(quotes) == 0 {
		return nil, ErrQuoteNotFound
	}
	return quotes, nil
}
