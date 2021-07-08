package service

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/langsdorf/gobank/dto"
	"github.com/langsdorf/gobank/infrastructure/grpc/pb"
	"github.com/langsdorf/gobank/usecase"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type TransactionService struct {
	ProcessTransactionUseCase usecase.UseCaseTransaction
	pb.UnimplementedPaymentServiceServer
}

func NewTransactionServce() *TransactionService {
	return &TransactionService{}
}

func (t *TransactionService) Payment(ctx context.Context, in *pb.PaymentRequest) (*empty.Empty, error) {
	transactionDto := dto.Transaction{
		Name:            in.GetCreditCard().Name,
		Number:          in.GetCreditCard().Number,
		ExpirationMonth: in.GetCreditCard().ExpirationMonth,
		ExpirationYear:  in.GetCreditCard().ExpirationYear,
		CVV:             in.GetCreditCard().Cvv,
		Store:           in.GetStore(),
		Amount:          in.GetAmount(),
		Description:     in.GetDescription(),
	}

	transaction, err := t.ProcessTransactionUseCase.ProcessTransaction(transactionDto)

	if err != nil {
		return &empty.Empty{}, status.Error(codes.FailedPrecondition, err.Error())
	}

	if transaction.Status != "APPROVED" {
		return &empty.Empty{}, status.Error(codes.FailedPrecondition, "transaction rejected by the bank")
	}

	return &empty.Empty{}, nil
}
