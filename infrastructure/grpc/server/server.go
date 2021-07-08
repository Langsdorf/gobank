package server

import (
	"log"
	"net"

	"github.com/langsdorf/gobank/infrastructure/grpc/pb"
	"github.com/langsdorf/gobank/infrastructure/grpc/service"
	"github.com/langsdorf/gobank/usecase"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type GRPCServer struct {
	ProcessTransactionUseCase usecase.UseCaseTransaction
}

func NewGRPCServer() GRPCServer {
	return GRPCServer{}
}

func (s GRPCServer) Serve() {

	lis, err := net.Listen("tcp", "0.0.0.0:50052")

	if err != nil {
		log.Fatal("TCP Port busy")
	}

	transactionService := service.NewTransactionServce()
	transactionService.ProcessTransactionUseCase = s.ProcessTransactionUseCase

	grpcServer := grpc.NewServer()

	reflection.Register(grpcServer)

	pb.RegisterPaymentServiceServer(grpcServer, transactionService)

	grpcServer.Serve(lis)
}
