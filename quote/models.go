package quote

import "time"

type Quote struct {
	Id        string    `json:"id"`
	UserId    string    `json:"user_id"`
	Quote     string    `json:"quote"`
	Approved  bool      `json:"approved"`
	CreatedAt time.Time `json:"created_at"`
}

type QuoteRequest struct {
	Quote string `json:"quote"`
}

type QuoteUpdateRequest struct {
	Id    string `json:"id"`
	Quote string `json:"quote"`
}

type QuoteResponse struct {
	Quote   *Quote
	Message string
}
