package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"stocks/internal/config"
	myHttp "stocks/internal/controller/http"
	"stocks/internal/repository"
	"stocks/internal/usecase"
	"stocks/pkg/logger"
	"stocks/pkg/postgres"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

func main() {
	//context
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	//logger
	logger.InitLogger()

	//config
	err := config.InitConfig()
	if err != nil {
		logger.Log.Fatal("Error loading config:", err)
		return
	}

	//env
	err = godotenv.Load()
	if err != nil {
		logger.Log.Fatal("Error loading .env file:", err)
	}

	//database
	dbConfig := &postgres.PostgresConfig{
		Host:     viper.GetString("database.dbhost"),
		Port:     viper.GetInt("database.dbport"),
		User:     viper.GetString("database.dbuser"),
		Password: os.Getenv("DB_PASSWORD"),
		Dbname:   viper.GetString("database.dbname"),
		SSLMode:  viper.GetString("database.dbsslmode"),
	}

	dbPool, err := postgres.NewDBPool(ctx, dbConfig)
	if err != nil {
		logger.Log.Fatal("Error connecting to database: ", err)
	}
	defer dbPool.Close()

	//transaction
	transaction := repository.NewPgTxManager(dbPool)

	//usecase
	stockUsecase := usecase.NewStockUsecase(*transaction)

	//controllers
	controller := myHttp.NewStockController(stockUsecase)

	newMux := http.NewServeMux()
	newMux.HandleFunc("POST /stocks/item/add", controller.AddStockController)
	newMux.HandleFunc("POST /stocks/item/get", controller.GetSkuStocksBySkuIdControlller)
	newMux.HandleFunc("POST /stocks/item/delete", controller.DeleteStockBySkuIdController)
	newMux.HandleFunc("POST /stocks/list/location", controller.GetSkusByLocationController)

	//server
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
