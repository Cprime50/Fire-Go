package quote

import (
	"errors"
	"log"
)

type QuoteService interface {
	CreateQuote(userId string, quote string) error
	UpdateQuote(userId string, role string, quoteId string, quote string) error
	DeleteQuote(userId string, role string, quoteId string) error
	GetQuotes(role string) ([]*Quote, error)
	GetQuotesByUserId(userId string, role string, requestedUserId string) ([]*Quote, error)
	ApproveQuote(userId string, role string, quoteId string) error
	GetUnapprovedQuotes() ([]*Quote, error)
}

type QuoteServiceImpl struct{}

func (s *QuoteServiceImpl) CreateQuote(userId string, quote string) error {
	if userId == "" || quote == "" {
		log.Println("Error: Invalid request body")
		return ErrInvalidRequestBody
	}
	err := createQuote(&Quote{
		UserId:   userId,
		Quote:    quote,
		Approved: false,
	})
	if err != nil {
		log.Println("Error creating quote:", err)
		return ErrCreateQuote
	}
	return nil
}

func (s *QuoteServiceImpl) UpdateQuote(userId string, role string, quoteId string, quote string) error {
	if userId == "" || quoteId == "" || role == "" || quote == "" {
		log.Println("Error: Invalid request body")
		return ErrInvalidRequestBody
	}

	quoteGotten, err := getQuoteById(quoteId)
	if err != nil {
		if errors.Is(err, ErrQuoteNotFound) {
			log.Println("Error: Quote not found")
			return ErrQuoteNotFound
		}
		log.Println("Error updating quote:", err)
		return ErrUpdateQuote
	}

	if role != "admin" {
		if userId != quoteGotten.UserId {
			log.Println("Error: Not authorized")
			return ErrNotAuthorized
		}
	}

	err = updateQuote(&Quote{
		Id:       quoteId,
		Quote:    quote,
		Approved: false,
	})

	if err != nil {
		log.Println("Error updating quote:", err)
		return ErrUpdateQuote
	}

	return nil
}

func (s *QuoteServiceImpl) DeleteQuote(userId string, role string, quoteId string) error {
	if userId == "" || quoteId == "" || role == "" {
		log.Println("Error: Invalid request body")
		return ErrInvalidRequestBody
	}

	quoteGotten, err := getQuoteById(quoteId)
	if err != nil {
		if errors.Is(err, ErrQuoteNotFound) {
			log.Println("Error: Quote not found")
			return ErrQuoteNotFound
		}
		log.Println("Error deleting quote:", err)
		return ErrDeletingQuote
	}

	if role != "admin" {
		if userId != quoteGotten.UserId {
			log.Println("Error: Not authorized")
			return ErrNotAuthorized
		}
	}

	err = deleteQuote(quoteId)
	if err != nil {
		log.Println("Error deleting quote:", err)
		return ErrDeletingQuote
	}

	return nil
}

func (s *QuoteServiceImpl) GetQuotes(role string) ([]*Quote, error) {
	var quotes []*Quote
	var err error

	if role == "admin" {
		quotes, err = getAllQuotes()
	} else {
		quotes, err = getAllApprovedQuotes()
	}

	if err != nil {
		if errors.Is(err, ErrQuoteNotFound) {
			log.Println("Error: No quotes found")
			return nil, ErrQuoteNotFound
		}
		log.Println("Error getting quotes:", err)
		return nil, err
	}

	return quotes, nil
}

func (s *QuoteServiceImpl) GetQuotesByUserId(userId string, role string, requestedUserId string) ([]*Quote, error) {
	if userId == "" || requestedUserId == "" {
		log.Println("Error: Invalid request body")
		return nil, ErrInvalidRequestBody
	}

	var quotes []*Quote
	var err error

	if role == "admin" || userId == requestedUserId {
		quotes, err = getQuotesByUserId(requestedUserId)
	} else {
		quotes, err = getApprovedQuotesByUserId(requestedUserId)
	}

	if err != nil {
		if errors.Is(err, ErrQuoteNotFound) {
			log.Println("Error: No quotes found")
			return nil, ErrQuoteNotFound
		}
		log.Println("Error getting quotes:", err)
		return nil, ErrGettingQuote
	}

	return quotes, nil
}

func (s *QuoteServiceImpl) ApproveQuote(userId string, role string, quoteId string) error {
	if userId == "" || quoteId == "" || role == "" {
		log.Println("Error: Invalid request body")
		return ErrInvalidRequestBody
	}

	if role != "admin" {
		log.Println("Error: Not authorized")
		return ErrNotAuthorized
	}

	err := approveQuote(quoteId)
	if err != nil {
		if errors.Is(err, ErrQuoteNotFound) {
			log.Println("Error: Quote not found")
			return ErrQuoteNotFound
		}
		log.Println("Error approving quote:", err)
		return ErrApprovingQuote
	}

	return nil
}

func (s *QuoteServiceImpl) GetUnapprovedQuotes() ([]*Quote, error) {
	unapprovedQuotes, err := getUnapprovedQuotes()
	if err != nil {
		if errors.Is(err, ErrQuoteNotFound) {
			log.Println("Error: No unapproved quotes found")
			return nil, ErrQuoteNotFound
		}
		log.Println("Error getting unapproved quotes:", err)
		return nil, ErrGettingQuote
	}
	return unapprovedQuotes, nil
}
