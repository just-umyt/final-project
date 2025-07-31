package app

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"
	"os"
	"os/signal"
	"stocks/internal/config"
	"stocks/internal/producer"
	"stocks/internal/repository"
	"stocks/internal/usecase"
	"stocks/pkg/postgres"
	"syscall"
	"time"

	myLog "stocks/internal/observability/log"
	"stocks/internal/observability/metrics"
	"stocks/internal/observability/tracer"
	myGrpc "stocks/internal/router/grpc"

	pb "stocks/pkg/api/stock"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	ErrLoadEnv        = "error loading .env file: %v"
	ErrDBConnect      = "error connecting to database: %v"
	ErrMigration      = "error migration: %v"
	ErrMigrationUp    = "error migration up: %v"
	ErrShutdown       = "shutdown error: %v"
	ErrListener       = "failed to listen: %v"
	ErrTracerShutdown = "failed to shutdown tracer: %v"
	ErrListenGRPC     = "failed to serve grpc server"
	ErrListenGateway  = "failed to serve gateway server"
	ErrListenMetrics  = "failed to serve metrics server"

	tracingServiceName = "stock-service"

	metricsTimeout           = 5 * time.Second
	gatewayReadHeaderTimeout = 3 * time.Second
	gatewayShutdownTimeout   = 5 * time.Second
)

func RunApp(env string, logger myLog.Logger) error {
	//context
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	//load config
	err := config.LoadConfig(env)
	if err != nil {
		return fmt.Errorf(ErrLoadEnv, err)
	}

	//tracer
	tracing, err := tracer.InitTracer(ctx, os.Getenv("JAEGER_ENDPOINT"), tracingServiceName)
	if err != nil {
		return err
	}

	defer func() {
		err := tracing.Shutdown(ctx)
		if err != nil && !errors.Is(err, context.Canceled) {
			logger.Errorf(ErrTracerShutdown, err)
		}
	}()

	//database config
	dbConfig := &postgres.PostgresConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Dbname:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	//migration up
	err = migrationUp(dbConfig)
	if err != nil {
		return err
	}

	//database Pool
	dbPool, err := postgres.NewDBPool(ctx, dbConfig)
	if err != nil {
		return fmt.Errorf(ErrDBConnect, err)
	}
	defer dbPool.Close()

	//kafka
	address := os.Getenv("KAFKA_BROKERS")

	kafkaProducer, err := producer.NewProducer(address)
	if err != nil {
		return err
	}

	defer kafkaProducer.Close()

	//grpc listener
	grpcServerAddress := fmt.Sprintf("%s:%s", os.Getenv("GRPC_HOST"), os.Getenv("GRPC_PORT"))

	lis, err := net.Listen(os.Getenv("GRPC_NETWORK"), grpcServerAddress)
	if err != nil {
		return fmt.Errorf(ErrListener, err)
	}

	trxManager := postgres.NewPgTxManager(dbPool)
	stockRepo := repository.NewStockRepository(dbPool)
	stockUsecase := usecase.NewStockUsecase(stockRepo, trxManager, kafkaProducer, logger)
	stockService := myGrpc.NewStockServer(stockUsecase)
	metric := metrics.RegisterMetrics()
	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(
		myGrpc.LoggingInterceptor(
			logger,
			metric,
			tracing.Tracer(tracingServiceName),
		),
	))

	//grpc register
	reflection.Register(grpcServer)
	pb.RegisterStockServiceServer(grpcServer, stockService)

	//gateway listener
	gatewayAddr := fmt.Sprintf("%s:%s", os.Getenv("GATEWAY_SERVER_HOST"), os.Getenv("GATEWAY_SERVER_PORT"))

	mux, err := myGrpc.NewMux(ctx, grpcServerAddress)
	if err != nil {
		return err
	}

	serverConfig := &myGrpc.ServerConfig{
		Address:           gatewayAddr,
		Handler:           mux,
		ReadHeaderTimeout: gatewayReadHeaderTimeout,
	}

	gatewayServer := myGrpc.NewGatewayServer(serverConfig)

	//grpc ListenAndServe
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			logger.Errorf(ErrListenGRPC, err)
		}
	}()

	//gateway ListenAndServe
	go func() {
		if err := gatewayServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatalf(ErrListenGateway, err)
		}
	}()

	//metrics ListenAndServe
	go func() {
		if err := metrics.ListenAndServe(os.Getenv("PROMETHEUS"), metricsTimeout); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Fatal(ErrListenMetrics, myLog.Error(err))
		}
	}()

	logger.Infof("listening in %s\n", gatewayAddr)

	//gracefull shutdowns
	<-ctx.Done()

	logger.Info("shutting down server gracefully...")

	grpcServer.GracefulStop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), gatewayShutdownTimeout)

	defer cancel()

	if err = gatewayServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf(ErrShutdown, err)
	}

	return nil
}

func migrationUp(dbConfig *postgres.PostgresConfig) error {
	db, err := postgres.NewDB(dbConfig)
	if err != nil {
		return fmt.Errorf(ErrDBConnect, err)
	}
	defer db.Close()

	migration, err := postgres.NewMigration(db, os.Getenv("MIGRATION_SOURCE_URL"))
	if err != nil {
		return fmt.Errorf(ErrMigration, err)
	}

	err = postgres.MigrationUp(migration)
	if err != nil {
		return fmt.Errorf(ErrMigrationUp, err)
	}

	return nil
}
