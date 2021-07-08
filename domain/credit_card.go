package domain

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type CreditCard struct {
	ID              string
	Name            string
	Number          string
	ExpirationMonth int32
	ExpirationYear  int32
	CVV             int32
	Balance         int32
	Limit           int32
	CreatedAt       time.Time
}

func NewCreditCard() *CreditCard {
	c := &CreditCard{}
	c.ID = uuid.NewV4().String()
	c.CreatedAt = time.Now()
	return c
}
