package main

import (
	myHttp "cart/internal/controller/http"
	"fmt"
	"strconv"

	"cart/internal/repository"
	"cart/internal/services"
	"cart/internal/usecase"
	"cart/pkg/logger"
	"cart/pkg/postgres"
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger.InitLogger()

	err := godotenv.Load(".env")
	if err != nil {
		logger.Log.Fatal("Error loading .env file:", err)
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
		logger.Log.Fatal("Error connecting to database: ", err)
	}
	defer dbPool.Close()

	transaction := repository.NewPgTxManager(dbPool)

	timeOut, err := strconv.Atoi(os.Getenv("CLIENT_TIMEOUT"))
	if err != nil {
		logger.Log.Fatal("Error loading CLIENT_TIMEOUT: ", err)
	}

	getSkuService := services.NewStockService(time.Duration(timeOut)*time.Second, os.Getenv("CLIENT_URL"))

	cartUsecase := usecase.NewCartUsecase(*transaction, getSkuService)

	controller := myHttp.NewCartController(cartUsecase)

	newMux := http.NewServeMux()
	newMux.HandleFunc("POST /cart/item/add", controller.AddItem)
	newMux.HandleFunc("POST /cart/item/delete", controller.DeleteItem)
	newMux.HandleFunc("POST /cart/list", controller.CartList)
	newMux.HandleFunc("POST /cart/clear", controller.CartClear)

	serverAddress := fmt.Sprintf("%s:%s", os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT"))

	serverTimeOut, err := strconv.Atoi(os.Getenv("SERVER_READ_HEADER_TIMEOUT"))
	if err != nil {
		logger.Log.Fatal("Error loading SERVER_READ_HEADER_TIMEOUT: ", err)
	}

	serverConfig := &myHttp.ServerConfig{
		Address:           serverAddress,
		Handler:           newMux,
		ReadHeaderTimeout: time.Duration(serverTimeOut) * time.Second,
	}

	server := myHttp.NewServer(serverConfig)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatalf("Error listen and serve : %v", err)
		}
	}()

	logger.Log.Infof("listening in  %s", serverAddress)

	<-ctx.Done()

	logger.Log.Info("shutting down server gracefully...")

	shutdownTimer, err := strconv.Atoi(os.Getenv("SERVER_SHUTDOWN_TIMEOUT"))
	if err != nil {
		logger.Log.Fatal("Error loading SERVER_SHUTDOWN_TIMEOUT: ", err)
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Duration(shutdownTimer)*time.Second)
	defer cancel()

	if err = server.Shutdown(shutdownCtx); err != nil {
		logger.Log.Warnf("shutdown error: %v", err)
	} else {
		logger.Log.Info("shutdown succes!")
	}
}
