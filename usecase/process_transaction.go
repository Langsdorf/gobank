package usecase

import (
	"encoding/json"
	"os"
	"time"

	"github.com/langsdorf/gobank/domain"
	"github.com/langsdorf/gobank/dto"
	"github.com/langsdorf/gobank/infrastructure/kafka"
)

type UseCaseTransaction struct {
	TransactionRepository domain.TransactionRepository
	KafkaProducer         kafka.KafkaProducer
}

func NewUseCaseTransaction(transactionRepository domain.TransactionRepository) UseCaseTransaction {
	return UseCaseTransaction{TransactionRepository: transactionRepository}
}

func (u *UseCaseTransaction) ProcessTransaction(transactionDto dto.Transaction) (domain.Transaction, error) {
	creditCard := domain.NewCreditCard()
	creditCard.Name = transactionDto.Name
	creditCard.Number = transactionDto.Number
	creditCard.ExpirationMonth = transactionDto.ExpirationMonth
	creditCard.ExpirationYear = transactionDto.ExpirationYear
	creditCard.CVV = transactionDto.CVV

	ccBalanceAndLimit, err := u.TransactionRepository.GetCreditCard(*creditCard)

	if err != nil {
		return domain.Transaction{}, err
	}

	creditCard.ID = ccBalanceAndLimit.ID
	creditCard.Balance = ccBalanceAndLimit.Balance
	creditCard.Limit = ccBalanceAndLimit.Limit

	t := domain.NewTransaction()

	t.CreditCardId = creditCard.ID
	t.Amount = transactionDto.Amount
	t.Store = transactionDto.Store
	t.Description = transactionDto.Description
	t.CreatedAt = time.Now()

	t.Validate(creditCard)

	err = u.TransactionRepository.Save(*t, *creditCard)

	if err != nil {
		return domain.Transaction{}, err
	}

	transactionDto.ID = t.ID
	transactionDto.CreatedAt = t.CreatedAt

	transactionJson, err := json.Marshal(transactionDto)

	if err != nil {
		return domain.Transaction{}, err
	}

	err = u.KafkaProducer.Publish(string(transactionJson), os.Getenv("KafkaTransactionsTopic"))

	if err != nil {
		return domain.Transaction{}, err
	}

	return *t, nil
}
