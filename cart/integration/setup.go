package integration

import (
	"cart/internal/producer"
	"cart/internal/repository"
	"cart/internal/services"
	"errors"
	"fmt"
	"log"
	"net"

	"cart/internal/usecase"
	"cart/pkg/postgres"
	"context"
	"database/sql"
	"net/http/httptest"
	"os"

	myLog "cart/internal/observability/log"
	"cart/internal/observability/tracer"
	myGrpc "cart/internal/router/grpc"
	pb "cart/pkg/api/cart"
	myZap "cart/pkg/zap"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

const (
	tracingServiceName = "cart-service"
	appLogPath         = "../app.log"
)

type testAppConfig struct {
	DB            *sql.DB
	Migration     *migrate.Migrate
	DBPool        *pgxpool.Pool
	Gateway       *httptest.Server
	CartGRPC      *grpc.Server
	FakeStockGRPC *grpc.Server
	StockClient   *grpc.ClientConn
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

	//db config
	dbConfig := &postgres.PostgresConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		DBname:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	//db
	t.DB, err = postgres.NewDB(dbConfig)
	if err != nil {
		return err
	}

	//migration
	t.Migration, err = postgres.NewMigration(t.DB, os.Getenv("MIGRATION_SOURCE_URL"))
	if err != nil {
		return err
	}

	err = postgres.MigrationUp(t.Migration)
	if err != nil {
		return err
	}

	//db pool
	t.DBPool, err = postgres.NewDBPool(ctx, dbConfig)
	if err != nil {
		return err
	}

	//stock client
	t.StockClient, err = grpc.NewClient(os.Getenv("CLIENT_URL"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}

	//kafka
	kafkaProducer, err := producer.NewProducer(os.Getenv("KAFKA_BROKERS"))
	if err != nil {
		return err
	}

	//grpc listener
	grpcServerAddress := fmt.Sprintf("%s:%s", os.Getenv("GRPC_HOST"), os.Getenv("GRPC_PORT"))

	lis, err := net.Listen(os.Getenv("GRPC_NETWORK"), grpcServerAddress)
	if err != nil {
		return fmt.Errorf("failed to listen: %w", err)
	}

	trxManager := postgres.NewPgTxManager(t.DBPool)
	cartRepo := repository.NewCartRepository(t.DBPool)
	stockService := services.NewStockClient(t.StockClient)
	cartUsecase := usecase.NewCartUsecase(cartRepo, trxManager, stockService, kafkaProducer, t.Logger)
	srv := myGrpc.NewCartServer(cartUsecase, t.Tracer.Tracer("cart-service"))

	t.CartGRPC = grpc.NewServer()
	reflection.Register(t.CartGRPC)

	pb.RegisterCartServiceServer(t.CartGRPC, srv)

	go func() {
		if err := t.CartGRPC.Serve(lis); err != nil {
			log.Printf("failed to serve: %v", err)
		}
	}()

	//gateway
	mux, err := myGrpc.NewMux(ctx, grpcServerAddress)
	if err != nil {
		return err
	}

	t.Gateway = httptest.NewServer(mux)

	//fake stock Service
	t.FakeStockGRPC, err = startFakeStockService(os.Getenv("CLIENT_URL"))
	if err != nil {
		return err
	}

	return nil
}

func (t *testAppConfig) Close() error {
	t.DB.Close()
	t.DBPool.Close()
	t.Gateway.Close()
	t.CartGRPC.GracefulStop()
	t.FakeStockGRPC.GracefulStop()
	t.StockClient.Close()
	t.LoggerCleanup()

	err := t.Tracer.Shutdown(context.Background())
	if err != nil && !errors.Is(err, context.Canceled) {
		t.Logger.Errorf("failed to shutdown tracer: %v", err)
	}

	return t.Migration.Down()
}
