package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/langsdorf/gobank/infrastructure/grpc/server"
	"github.com/langsdorf/gobank/infrastructure/kafka"
	"github.com/langsdorf/gobank/infrastructure/repository"
	"github.com/langsdorf/gobank/usecase"
	_ "github.com/lib/pq"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("error loading .env file")
	}
}

func main() {
	fmt.Println("GOBANK INICIADO")
	db := setupDb()
	defer db.Close()

	// cc := domain.NewCreditCard()
	// cc.Number = "1234"
	// cc.Balance = 0
	// cc.CVV = 123
	// cc.CreatedAt = time.Now()
	// cc.ExpirationMonth = 10
	// cc.ExpirationYear = 2025
	// cc.Limit = 1000
	// cc.Name = "gobank"

	// repo := repository.NewTransactionDb(db)

	// err := repo.CreateCreditCard(*cc)

	// if err != nil {
	// 	log.Fatal(err)
	// }

	producer := setupKafkaProducer()
	processTransactionUseCase := setupUseCase(db, producer)

	serveGRPC(processTransactionUseCase)
}

func serveGRPC(processTransactionUseCase usecase.UseCaseTransaction) {
	grpcServer := server.NewGRPCServer()

	grpcServer.ProcessTransactionUseCase = processTransactionUseCase

	grpcServer.Serve()
}

func setupKafkaProducer() kafka.KafkaProducer {
	producer := kafka.NewKafkaProducer()

	producer.SetupProducer(os.Getenv("KafkaBootstrapServers"))

	return producer
}

func setupUseCase(db *sql.DB, producer kafka.KafkaProducer) usecase.UseCaseTransaction {
	transactionRepository := repository.NewTransactionDb(db)

	useCase := usecase.NewUseCaseTransaction(transactionRepository)

	useCase.KafkaProducer = producer

	return useCase
}

func setupDb() *sql.DB {
	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		os.Getenv("host"),
		os.Getenv("port"),
		os.Getenv("user"),
		os.Getenv("password"),
		os.Getenv("dbname"),
	)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		log.Fatal("Erro na conexão com o banco de dados")
	}

	return db
}
