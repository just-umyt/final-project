package app

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"stocks/internal/config"
	"stocks/internal/producer"
	"stocks/internal/repository"
	myHttp "stocks/internal/router/http"
	"stocks/internal/router/http/controller"
	"stocks/internal/usecase"
	"stocks/pkg/postgres"
	"strconv"
	"syscall"
	"time"
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

	address := os.Getenv("KAFKA_BROKERS")

	kafkaProducer, err := producer.NewProducer(address)
	if err != nil {
		return err
	}

	defer kafkaProducer.Close()

	stockUsecase := usecase.NewStockUsecase(stockRepo, trxManager, kafkaProducer)

	stockController := controller.NewStockController(stockUsecase)

	newMux := myHttp.NewMux(stockController)

	serverAddress := fmt.Sprintf("%s:%s", os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT"))

	readHeaderTimeOut, err := strconv.Atoi(os.Getenv("SERVER_READ_HEADER_TIMEOUT"))
	if err != nil {
		err = fmt.Errorf(ErrLoadServerReadTimeOut, err)
		return err
	}

	serverConfig := &myHttp.ServerConfig{
		Address:           serverAddress,
		Handler:           newMux,
		ReadHeaderTimeout: time.Duration(readHeaderTimeOut) * time.Second,
	}

	server := myHttp.NewServer(serverConfig)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error listen and serve : %v", err)
		}
	}()

	log.Printf("listening in  %s\n", serverAddress)

	<-ctx.Done()

	log.Println("shutting down server gracefully...")

	shutdown, err := strconv.Atoi(os.Getenv("SERVER_SHUTDOWN_TIMEOUT"))
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
