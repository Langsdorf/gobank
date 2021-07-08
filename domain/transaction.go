package domain

import (
	"time"

	uuid "github.com/satori/go.uuid"
)

type TransactionRepository interface {
	Save(transaction Transaction, creditCard CreditCard) error
	GetCreditCard(creditCard CreditCard) (CreditCard, error)
	CreateCreditCard(creditCard CreditCard) error
}

type Transaction struct {
	ID           string
	Name         string
	Amount       float64
	Status       string
	Description  string
	Store        string
	CreditCardId string
	CreatedAt    time.Time
}

func NewTransaction() *Transaction {
	t := &Transaction{}

	t.ID = uuid.NewV4().String()
	t.CreatedAt = time.Now()

	return t
}

func (t *Transaction) Validate(creditCard *CreditCard) {
	if t.Amount+float64(creditCard.Balance) > float64(creditCard.Limit) {
		t.Status = "REJECTED"
	} else {
		t.Status = "APPROVED"
		creditCard.Balance += int32(t.Amount)
	}
}
