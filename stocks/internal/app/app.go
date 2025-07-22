package app

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"stocks/internal/config"
	"stocks/internal/repository"
	"stocks/internal/usecase"
	"stocks/pkg/postgres"
	"strconv"
	"syscall"
	"time"

	myGrpc "stocks/internal/router/grpc"

	pb "stocks/pkg/api"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

var (
	ErrLoadEnv               = "error loading .env file: %v"
	ErrDBConnect             = "error connecting to database: %v"
	ErrMigration             = "error migration: %v"
	ErrMigrationUp           = "error migration up: %v"
	ErrLoadServerReadTimeOut = "error loading SERVER_READ_HEADER_TIMEOUT: %v"
	ErrLoadServerShutdown    = "error loading SERVER_SHUTDOWN_TIMEOUT: %v"
	ErrShutdown              = "shutdown error: %v"
)

func RunApp(env string) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	err := config.LoadConfig(env)
	if err != nil {
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
		err = fmt.Errorf(ErrDBConnect, err)
		return err
	}
	defer db.Close()

	migration, err := postgres.NewMigration(db, os.Getenv("MIGRATION_SOURCE_URL"))
	if err != nil {
		err = fmt.Errorf(ErrMigration, err)
		return err
	}

	err = postgres.MigrationUp(migration)
	if err != nil {
		err = fmt.Errorf(ErrMigrationUp, err)
		return err
	}

	dbPool, err := postgres.NewDBPool(ctx, dbConfig)
	if err != nil {
		err = fmt.Errorf(ErrDBConnect, err)
		return err
	}
	defer dbPool.Close()

	trxManager := postgres.NewPgTxManager(dbPool)

	stockRepo := repository.NewStockRepository(dbPool)

	// address := os.Getenv("KAFKA_BROKERS")

	// kafkaProducer, err := producer.NewProducer(address)
	// if err != nil {
	// 	return err
	// }

	// defer kafkaProducer.Close()
	var kafkaProducer usecase.IProducer

	stockUsecase := usecase.NewStockUsecase(stockRepo, trxManager, kafkaProducer)

	// stockController := controller.NewStockController(stockUsecase)

	grpcServerAddress := fmt.Sprintf("%s:%s", os.Getenv("GRPC_HOST"), os.Getenv("GRPC_PORT"))

	lis, err := net.Listen(os.Getenv("GRPC_NETWORK"), grpcServerAddress)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	srv := myGrpc.NewStockServer(stockUsecase)

	s := grpc.NewServer()

	reflection.Register(s)

	pb.RegisterStockServiceServer(s, srv)

	go func() {
		if err := s.Serve(lis); err != nil {
			log.Fatalf("failed to serve: %v", err)
		}
	}()

	serverAddress := fmt.Sprintf("%s:%s", os.Getenv("GATEWAY_SERVER_HOST"), os.Getenv("GATEWAY_SERVER_PORT"))

	mux, err := myGrpc.NewMux(ctx, grpcServerAddress)
	if err != nil {
		return err
	}

	serverConfig := &myGrpc.ServerConfig{
		Address: serverAddress,
		Handler: mux,
	}

	server := myGrpc.NewGatewayServer(serverConfig)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error listen and serve : %v", err)
		}
	}()

	log.Printf("listening in  %s\n", serverAddress)

	<-ctx.Done()

	log.Println("shutting down server gracefully...")

	s.GracefulStop()

	shutdown, err := strconv.Atoi(os.Getenv("GATEWAY_SERVER_SHUTDOWN_TIMEOUT"))
	if err != nil {
		err = fmt.Errorf(ErrLoadServerShutdown, err)
		return err
	}
	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Duration(shutdown)*time.Second)

	defer cancel()

	if err = server.Shutdown(shutdownCtx); err != nil {
		err = fmt.Errorf(ErrShutdown, err)

		return err
	}

	return nil
}
