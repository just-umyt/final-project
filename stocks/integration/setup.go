package integration

import (
	"context"
	"database/sql"
	"errors"
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
	pb "stocks/pkg/api/stock"
	myZap "stocks/pkg/zap"

	myLog "stocks/internal/observability/log"
	"stocks/internal/observability/tracer"

	"go.opentelemetry.io/otel/sdk/trace"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
)

const (
	tracingServiceName = "stock-service"
	appLogPath         = "../app.log"
)

type testAppConfig struct {
	DB            *sql.DB
	Migration     *migrate.Migrate
	DBPool        *pgxpool.Pool
	StockGRPC     *grpc.Server
	Gateway       *httptest.Server
	Logger        myLog.Logger
	LoggerCleanup func()
	Tracer        *trace.TracerProvider
}

func (t *testAppConfig) Setup(ctx context.Context) error {
	var err error

	//logger
	t.Logger, t.LoggerCleanup, err = myZap.NewLogger(appLogPath)
	if err != nil {
		return err
	}

	//tracing
	t.Tracer, err = tracer.InitTracer(ctx, os.Getenv("JAEGER_ENDPOINT"), tracingServiceName)
	if err != nil {
		return err
	}

	//database config
	dbConfig := &postgres.PostgresConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Dbname:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	//database
	t.DB, err = postgres.NewDB(dbConfig)
	if err != nil {
		return err
	}

	//migration setup
	t.Migration, err = postgres.NewMigration(t.DB, os.Getenv("MIGRATION_SOURCE_URL"))
	if err != nil {
		return err
	}

	//migration up
	err = postgres.MigrationUp(t.Migration)
	if err != nil {
		return err
	}

	//database pool
	t.DBPool, err = postgres.NewDBPool(ctx, dbConfig)
	if err != nil {
		return err
	}

	//kafka
	kafkaProducer, err := producer.NewProducer(os.Getenv("KAFKA_BROKERS"))
	if err != nil {
		return err
	}

	//grpc
	grpcServerAddress := fmt.Sprintf("%s:%s", os.Getenv("GRPC_HOST"), os.Getenv("GRPC_PORT"))

	//grpc listener
	lis, err := net.Listen(os.Getenv("GRPC_NETWORK"), grpcServerAddress)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	trxManager := postgres.NewPgTxManager(t.DBPool)
	stockRepo := repository.NewStockRepository(t.DBPool)
	stockUsecase := usecase.NewStockUsecase(stockRepo, trxManager, kafkaProducer, t.Logger)
	srv := myGrpc.NewStockServer(stockUsecase)

	t.StockGRPC = grpc.NewServer()
	pb.RegisterStockServiceServer(t.StockGRPC, srv)

	go func() {
		if err := t.StockGRPC.Serve(lis); err != nil {
			log.Printf("failed to serve: %v", err)
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
	t.LoggerCleanup()

	err := t.Tracer.Shutdown(context.Background())
	if err != nil && !errors.Is(err, context.Canceled) {
		t.Logger.Errorf("failed to shutdown tracer: %v", err)
	}

	return t.Migration.Down()
}
