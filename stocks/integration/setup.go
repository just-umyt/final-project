package integration

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http/httptest"
	"os"
	"stocks/internal/producer"
	"stocks/internal/repository"
	"stocks/internal/usecase"
	"stocks/pkg/postgres"

	myGrpc "stocks/internal/router/grpc"
	pb "stocks/pkg/api"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

type testAppConfig struct {
	DB        *sql.DB
	Migration *migrate.Migrate
	DBPool    *pgxpool.Pool
	StockGRPC *grpc.Server
	Gateway   *httptest.Server
}

func (t *testAppConfig) Setup(ctx context.Context) error {
	var err error

	dbConfig := &postgres.PostgresConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Dbname:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	t.DB, err = postgres.NewDB(dbConfig)
	if err != nil {
		return err
	}

	t.Migration, err = postgres.NewMigration(t.DB, os.Getenv("MIGRATION_SOURCE_URL"))
	if err != nil {
		return err
	}

	err = postgres.MigrationUp(t.Migration)
	if err != nil {
		return err
	}

	t.DBPool, err = postgres.NewDBPool(ctx, dbConfig)
	if err != nil {
		return err
	}

	trxManager := postgres.NewPgTxManager(t.DBPool)

	stockRepo := repository.NewStockRepository(t.DBPool)

	//kafka
	kafkaProducer, err := producer.NewProducer(os.Getenv("KAFKA_BROKERS"))
	if err != nil {
		return err
	}

	stockUsecae := usecase.NewStockUsecase(stockRepo, trxManager, kafkaProducer)

	//grpc
	grpcServerAddress := fmt.Sprintf("%s:%s", os.Getenv("GRPC_HOST"), os.Getenv("GRPC_PORT"))

	lis, err := net.Listen(os.Getenv("GRPC_NETWORK"), grpcServerAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := myGrpc.NewStockServer(stockUsecae)

	t.StockGRPC = grpc.NewServer()

	pb.RegisterStockServiceServer(t.StockGRPC, srv)

	go func() {
		if err := t.StockGRPC.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	//gateway
	mux, err := myGrpc.NewMux(ctx, grpcServerAddress)
	if err != nil {
		return err
	}

	t.Gateway = httptest.NewServer(mux)

	return nil
}

func (t *testAppConfig) Close() error {
	t.DB.Close()
	t.DBPool.Close()
	t.StockGRPC.GracefulStop()
	t.Gateway.Close()

	return t.Migration.Down()
}
