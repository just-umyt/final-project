package main

import (
	"cart/internal/config"
	myHttp "cart/internal/controller/http"
	"fmt"

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
	"github.com/spf13/viper"
)

func main() {
	//context
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	//logger
	logger.InitLogger()

	// config
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

	//database pool
	dbPool, err := postgres.NewDBPool(ctx, dbConfig)
	if err != nil {
		logger.Log.Fatal("Error connecting to database: ", err)
	}
	defer dbPool.Close()

	//transaction
	transaction := repository.NewPgTxManager(dbPool)

	//service
	httpClient := http.Client{
		Timeout: viper.GetDuration("client.timeout") * time.Second,
	}
	getSkuService := services.NewSkuGetService(&httpClient, viper.GetString("client.url"))

	//usecase
	cartUsecase := usecase.NewCartUsecase(*transaction, getSkuService)

	//controllers
	controller := myHttp.NewCartController(cartUsecase)

	newMux := http.NewServeMux()
	newMux.HandleFunc("POST /cart/item/add", controller.CartAddItemController)
	newMux.HandleFunc("POST /cart/item/delete", controller.DeleteItemController)
	newMux.HandleFunc("POST /cart/list", controller.CartListController)
	newMux.HandleFunc("POST /cart/clear", controller.CartClearController)

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
