package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	myHttp "stocks/internal/controller/http"
	"stocks/internal/repository"
	"stocks/internal/usecase"
	"stocks/pkg/logger"
	"stocks/pkg/postgres"

	"github.com/joho/godotenv"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger.InitLogger()

	err := godotenv.Load(".env")
	if err != nil {
		logger.Log.Fatal("Error loading env file:", err)
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

	stockUsecase := usecase.NewStockUsecase(*transaction)

	controller := myHttp.NewStockController(stockUsecase)

	newMux := http.NewServeMux()
	newMux.HandleFunc("POST /stocks/item/add", controller.AddStock)
	newMux.HandleFunc("POST /stocks/item/get", controller.GetSkuStocksBySkuId)
	newMux.HandleFunc("POST /stocks/item/delete", controller.DeleteStockBySkuId)
	newMux.HandleFunc("POST /stocks/list/location", controller.GetSkusByLocation)

	serverAddr := fmt.Sprintf("%s:%s", os.Getenv("SERVER_HOST"), os.Getenv("SERVER_PORT"))

	readHeaderTimeOut, err := strconv.Atoi(os.Getenv("SERVER_READ_HEADER_TIMEOUT"))
	if err != nil {
		logger.Log.Fatal("Error loading SERVER_READ_HEADER_TIMEOUT: ", err)
	}

	serverConfig := &myHttp.ServerConfig{
		Addr:              serverAddr,
		Handler:           newMux,
		ReadHeaderTimeout: time.Duration(readHeaderTimeOut) * time.Second,
	}

	server := myHttp.NewServer(serverConfig)

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatalf("listen and serve: %v", err)
		}
	}()

	logger.Log.Infof("listening in  %s", serverAddr)

	<-ctx.Done()

	logger.Log.Info("shutting down server gracefully...")

	shutdown, err := strconv.Atoi(os.Getenv("SERVER_SHUTDOWN_TIMEOUT"))
	if err != nil {
		logger.Log.Fatal("Error loading SERVER_SHUTDOWN_TIMEOUT: ", err)
	}
	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Duration(shutdown)*time.Second)

	defer cancel()

	if err = server.Shutdown(shutdownCtx); err != nil {
		logger.Log.Warnf("shutdown error: %v", err)
	} else {
		logger.Log.Info("shutdown succes!")
	}
}
