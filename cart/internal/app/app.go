package app

import (
	"cart/internal/config"
	"cart/internal/producer"
	myGrpc "cart/internal/router/grpc"
	"cart/internal/services"
	"errors"
	"net"
	"net/http"
	"time"

	myLog "cart/internal/observability/log"
	"cart/internal/observability/metrics"
	"cart/internal/observability/tracer"

	"cart/internal/repository"
	"cart/internal/usecase"
	pb "cart/pkg/api/cart"
	"cart/pkg/postgres"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

const (
	ErrLoadEnv           = "error loading .env file: %v"
	ErrDBConnect         = "error connecting to database: %v"
	ErrMigration         = "error migration: %v"
	ErrMigrationUp       = "error migration up: %v"
	ErrLoadClientTimeOut = "error loading CLIENT_TIMEOUT: %v"
	ErrShutdown          = "shutdown error: %v"
	ErrListener          = "failed to listen: %v"
	ErrListenGRPC        = "failed to serve grpc server"
	ErrListenGateway     = "failed to serve gateway server"
	ErrListenMetrics     = "failed to serve metrics server"
	ErrTracerShutdown    = "failed to shutdown tracer: %v"

	tracingServiceName = "cart-service"

	metricsTimeout           = 5 * time.Second
	gatewayShutdownTimeout   = 5 * time.Second
	gatewayReadHeaderTimeout = 3 * time.Second
)

func RunApp(env string, logger myLog.Logger) error {
	//context
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	//load config
	if err := config.LoadConfig(env); err != nil {
		err = fmt.Errorf(ErrLoadEnv, err)
		return err
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
		DBname:   os.Getenv("DB_NAME"),
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

	//stock service
	conn, err := grpc.NewClient(os.Getenv("CLIENT_URL"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()

	//kafka
	kafkaBrokers := os.Getenv("KAFKA_BROKERS")

	kafkaProducer, err := producer.NewProducer(kafkaBrokers)
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

	cartRepo := repository.NewCartRepository(dbPool)
	trxManager := postgres.NewPgTxManager(dbPool)
	stockService := services.NewStockClient(conn)
	cartUsecase := usecase.NewCartUsecase(cartRepo, trxManager, stockService, kafkaProducer, logger)
	cartService := myGrpc.NewCartServer(cartUsecase, tracing.Tracer(tracingServiceName))
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
	pb.RegisterCartServiceServer(grpcServer, cartService)

	//gateway listener
	gatewayAddr := fmt.Sprintf("%s:%s", os.Getenv("GATEWAY_SERVER_HOST"), os.Getenv("GATEWAY_SERVER_PORT"))

	mux, err := myGrpc.NewMux(ctx, grpcServerAddress)
	if err != nil {
		return err
	}

	serverConfig := &myGrpc.ServerConfig{
		Address:             gatewayAddr,
		Handler:             mux,
		ReaderHeaderTimeout: gatewayReadHeaderTimeout,
	}

	gatewayServer := myGrpc.NewGatewayServer(serverConfig)

	//grpc ListenAndServe
	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			logger.Error(ErrListenGRPC, myLog.Error(err))
		}
	}()

	//gateway ListenAndServe
	go func() {
		if err := gatewayServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error(ErrListenGateway, myLog.Error(err))
		}
	}()

	//metrics ListenAndServe
	go func() {
		if err := metrics.ListenAndServe(os.Getenv("PROMETHEUS"), metricsTimeout); err != nil && !errors.Is(err, http.ErrServerClosed) {
			logger.Error(ErrListenMetrics, myLog.Error(err))
		}
	}()

	logger.Infof("gateway listening in %s", gatewayAddr)

	//gracefull shutdown
	<-ctx.Done()

	logger.Info("shutting down server  gracefully...")

	grpcServer.GracefulStop()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), gatewayShutdownTimeout)
	defer cancel()

	if err = gatewayServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf(ErrShutdown, err)
	}

	return nil
}

func migrationUp(dbConfig *postgres.PostgresConfig) error {
	//database for migration
	db, err := postgres.NewDB(dbConfig)
	if err != nil {
		return fmt.Errorf(ErrDBConnect, err)
	}

	//migration
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
