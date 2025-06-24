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

	"stocks/internal/config"
	myHttp "stocks/internal/controller/http"
	"stocks/internal/repository"
	"stocks/internal/usecase"
	"stocks/pkg/logger"
	"stocks/pkg/postgres"

	"github.com/spf13/viper"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	logger.InitLogger()

	err := config.InitConfig()
	if err != nil {
		logger.Log.Fatal("Error loading config:", err)
		return
	}

	dbPort, err := strconv.Atoi(os.Getenv("DB_PORT"))
	if err != nil {
		logger.Log.Fatal("Error load port of Database:", err)
	}
	dbConfig := &postgres.PostgresConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     dbPort,
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

	serverAddr := fmt.Sprintf("%s:%d", viper.GetString("server.host"), viper.GetInt("server.port"))
	readHeaderTimeOut := time.Duration(viper.GetInt("server.readheadertimeout")) * time.Second

	serverConfig := &myHttp.SeverConfig{
		Addr:              serverAddr,
		Handler:           newMux,
		ReadHeaderTimeout: readHeaderTimeOut,
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

	shutdownCtx, cancel := context.WithTimeout(context.Background(), time.Duration(viper.GetInt("server.shutdowntimeout"))*time.Second)
	defer cancel()

	if err = server.Shutdown(shutdownCtx); err != nil {
		logger.Log.Warnf("shutdown: %v", err)
	} else {
		logger.Log.Info("shutdown succes!")
	}
}
