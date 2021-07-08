package dto

import "time"

type Transaction struct {
	ID              string    `json:"transaction_id"`
	Name            string    `json:"-"`
	Amount          float64   `json:"amount"`
	Description     string    `json:"description"`
	Store           string    `json:"store"`
	Number          string    `json:"credit_card_number"`
	ExpirationMonth int32     `json:"-"`
	ExpirationYear  int32     `json:"-"`
	CVV             int32     `json:"-"`
	CreatedAt       time.Time `json:"payment_date"`
}
