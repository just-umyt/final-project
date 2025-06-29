package app

import (
	"cart/internal/config"
	myHttp "cart/internal/controller/http"
	"log"

	"cart/internal/repository"
	"cart/internal/services"
	"cart/internal/usecase"
	"cart/pkg/postgres"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"
)

var (
	ErrLoadEnv               = "error loading .env file: %v"
	ErrDBConnect             = "error connecting to database: %v"
	ErrLoadClientTimeOut     = "error loading CLIENT_TIMEOUT: %v"
	ErrLoadServerReadTimeOut = "error loading SERVER_READ_HEADER_TIMEOUT: %v"
	ErrLoadServerShutdown    = "error loading SERVER_SHUTDOWN_TIMEOUT: %v"
	ErrShutdown              = "shutdown error: %v"
)

func RunApp() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := config.LoadConfig(".env"); err != nil {
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

	dbPool, err := postgres.NewDBPool(ctx, dbConfig)
	if err != nil {
		err = fmt.Errorf(ErrDBConnect, err)
		return err
	}
	defer dbPool.Close()

	cartRepo := repository.NewCartRepository(dbPool)

	trxManager := repository.NewPgTxManager(dbPool)

	timeOut, err := strconv.Atoi(os.Getenv("CLIENT_TIMEOUT"))
	if err != nil {
		err = fmt.Errorf(ErrLoadClientTimeOut, err)
		return err
	}

	stockService := services.NewStockService(time.Duration(timeOut)*time.Second, os.Getenv("CLIENT_URL"))

	cartUsecase := usecase.NewCartUsecase(cartRepo, trxManager, stockService)

	controller := myHttp.NewCartController(cartUsecase)

	newMux := myHttp.NewMux(controller)

	serverAddress := fmt.Sprintf("%s:%s", os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT"))

	serverTimeOut, err := strconv.Atoi(os.Getenv("SERVER_READ_HEADER_TIMEOUT"))
	if err != nil {
		err = fmt.Errorf(ErrLoadServerReadTimeOut, err)
		return err
	}

	serverConfig := &myHttp.ServerConfig{
		Address:           serverAddress,
		Handler:           newMux,
		ReadHeaderTimeout: time.Duration(serverTimeOut) * time.Second,
	}

	server := myHttp.NewServer(serverConfig)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Error listen and serve : %v", err)
		}
	}()

	log.Printf("listening in  %s", serverAddress)

	<-ctx.Done()

	log.Println("shutting down server gracefully...")

	shutdownTimer, err := strconv.Atoi(os.Getenv("SERVER_SHUTDOWN_TIMEOUT"))
	if err != nil {
		err = fmt.Errorf(ErrLoadServerShutdown, err)
		return err
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Duration(shutdownTimer)*time.Second)
	defer cancel()

	if err = server.Shutdown(shutdownCtx); err != nil {
		err = fmt.Errorf(ErrShutdown, err)
		return err
	}

	return nil
}
