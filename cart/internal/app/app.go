package app

import (
	"cart/internal/config"
	"cart/internal/producer"
	myGrpc "cart/internal/router/grpc"
	"cart/internal/services"
	"errors"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"

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

var (
	ErrLoadEnv               = "error loading .env file: %v"
	ErrDBConnect             = "error connecting to database: %v"
	ErrMigration             = "error migration: %v"
	ErrMigrationUp           = "error migration up: %v"
	ErrLoadClientTimeOut     = "error loading CLIENT_TIMEOUT: %v"
	ErrLoadServerReadTimeOut = "error loading SERVER_READ_HEADER_TIMEOUT: %v"
	ErrLoadServerShutdown    = "error loading SERVER_SHUTDOWN_TIMEOUT: %v"
	ErrShutdown              = "shutdown error: %v"
	ErrListener              = "failed to listen: %v"
)

func RunApp(env string) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := config.LoadConfig(env); err != nil {
		err = fmt.Errorf(ErrLoadEnv, err)
		return err
	}

	dbConfig := &postgres.PostgresConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Dbname:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	db, err := postgres.NewDB(dbConfig)
	if err != nil {
		return fmt.Errorf(ErrDBConnect, err)
	}

	migration, err := postgres.NewMigration(db, os.Getenv("MIGRATION_SOURCE_URL"))
	if err != nil {
		return fmt.Errorf(ErrMigration, err)
	}

	err = postgres.MigrationUp(migration)
	if err != nil {
		return fmt.Errorf(ErrMigrationUp, err)
	}

	dbPool, err := postgres.NewDBPool(ctx, dbConfig)
	if err != nil {
		return fmt.Errorf(ErrDBConnect, err)
	}
	defer dbPool.Close()

	cartRepo := repository.NewCartRepository(dbPool)

	trxManager := postgres.NewPgTxManager(dbPool)

	conn, err := grpc.NewClient(os.Getenv("CLIENT_URL"), grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return err
	}
	defer conn.Close()

	stockService := services.NewStockClient(conn)

	kafkaBrokers := os.Getenv("KAFKA_BROKERS")

	kafkaProducer, err := producer.NewProducer(kafkaBrokers)
	if err != nil {
		return err
	}

	defer kafkaProducer.Close()

	cartUsecase := usecase.NewCartUsecase(cartRepo, trxManager, stockService, kafkaProducer)

	//grpc listener
	grpcServerAddress := fmt.Sprintf("%s:%s", os.Getenv("GRPC_HOST"), os.Getenv("GRPC_PORT"))

	lis, err := net.Listen(os.Getenv("GRPC_NETWORK"), grpcServerAddress)
	if err != nil {
		return fmt.Errorf(ErrListener, err)
	}

	cartService := myGrpc.NewCartServer(cartUsecase)

	grpcServer := grpc.NewServer()

	reflection.Register(grpcServer)

	pb.RegisterCartServiceServer(grpcServer, cartService)

	go func() {
		if err := grpcServer.Serve(lis); err != nil {
			log.Printf("failed to serve: %v", err)
		}
	}()

	//gateway listener
	gatewayAddr := fmt.Sprintf("%s:%s", os.Getenv("GATEWAY_SERVER_HOST"), os.Getenv("GATEWAY_SERVER_PORT"))

	mux, err := myGrpc.NewMux(ctx, grpcServerAddress)
	if err != nil {
		return err
	}

	serverConfig := &myGrpc.ServerConfig{
		Address: gatewayAddr,
		Handler: mux,
	}

	gatewayServer := myGrpc.NewGatewayServer(serverConfig)

	go func() {
		if err := gatewayServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Fatalf("error listen and serve : %v", err)
		}
	}()

	log.Printf("listening in  %s", gatewayAddr)

	<-ctx.Done()

	log.Println("shutting down server gracefully...")

	grpcServer.GracefulStop()

	shutdownTimer, err := strconv.Atoi(os.Getenv("GATEWAY_SERVER_SHUTDOWN_TIMEOUT"))
	if err != nil {
		return fmt.Errorf(ErrLoadServerShutdown, err)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Duration(shutdownTimer)*time.Second)
	defer cancel()

	if err = gatewayServer.Shutdown(shutdownCtx); err != nil {
		return fmt.Errorf(ErrShutdown, err)
	}

	return nil
}
