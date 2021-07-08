package repository

import (
	"database/sql"
	"errors"

	"github.com/langsdorf/gobank/domain"
)

type TransactionDb struct {
	db *sql.DB
}

func NewTransactionDb(db *sql.DB) *TransactionDb {
	return &TransactionDb{db: db}
}

func (t *TransactionDb) Save(transaction domain.Transaction, creditCard domain.CreditCard) error {
	stmt, err := t.db.Prepare(`INSERT INTO transactions (id, credit_card_id, amount, status, description, store, created_at) VALUES ($1, $2, $3, $4, $5, $6, $7)`)

	if err != nil {
		return err
	}

	_, err = stmt.Exec(
		transaction.ID,
		creditCard.ID,
		transaction.Amount,
		transaction.Status,
		transaction.Description,
		transaction.Store,
		transaction.CreatedAt,
	)

	if err != nil {
		return err
	}

	if transaction.Status == "APPROVED" {
		err = t.UpdateBalance(creditCard)

		if err != nil {
			return err
		}
	}

	err = stmt.Close()

	if err != nil {
		return err
	}

	return nil
}

func (t *TransactionDb) UpdateBalance(c domain.CreditCard) error {
	stmt, err := t.db.Prepare(`UPDATE credit_cards SET balance = $1 WHERE id = $2`)

	if err != nil {
		return err
	}

	_, err = stmt.Exec(
		c.Balance,
		c.ID,
	)

	if err != nil {
		return err
	}

	err = stmt.Close()

	if err != nil {
		return err
	}

	return nil
}

func (t *TransactionDb) CreateCreditCard(c domain.CreditCard) error {
	stmt, err := t.db.Prepare(`INSERT INTO credit_cards (id, name, number, expiration_month, expiration_year, cvv, balance, balance_limit) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`)

	if err != nil {
		return err
	}

	_, err = stmt.Exec(
		c.ID,
		c.Name,
		c.Number,
		c.ExpirationMonth,
		c.ExpirationYear,
		c.CVV,
		c.Balance,
		c.Limit,
	)

	if err != nil {
		return err
	}

	err = stmt.Close()

	if err != nil {
		return err
	}

	return nil
}

func (t *TransactionDb) GetCreditCard(creditCard domain.CreditCard) (domain.CreditCard, error) {
	var c domain.CreditCard

	stmt, err := t.db.Prepare("SELECT id, balance, balance_limit FROM credit_cards WHERE number=$1")

	if err != nil {
		return c, err
	}

	if err = stmt.QueryRow(creditCard.Number).Scan(&c.ID, &c.Balance, &c.Limit); err != nil {
		return c, errors.New("credit card does not exists")
	}

	return c, nil
}
