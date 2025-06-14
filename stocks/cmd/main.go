package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
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
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

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

	// item := models.SKU{
	// 	SkuId:    2,
	// 	Name:     "new item2",
	// 	Price:    100,
	// 	Count:    34,
	// 	Type:     "item",
	// 	Location: "Aisle 3",
	// 	UserId:   2,
	// }

	//controllers
	controller := myHttp.NewStockController(stockUsecase)

	newMux := http.NewServeMux()
	newMux.HandleFunc("POST /stocks/item/add", controller.AddSkuController)
	newMux.HandleFunc("POST /stocks/item/get", controller.GetSkuBySkuIdControlller)

	//server
	serverAddr := fmt.Sprintf("%s:%d", viper.GetString("server.host"), viper.GetInt("server.port"))
	readHeaderTimeOut := time.Duration(viper.GetInt("server.readheadertimeout")) * time.Second

	serverConfig := &myHttp.SeverConfig{
		Addr:              serverAddr,
		Handler:           newMux,
		ReadHeaderTimeout: readHeaderTimeOut,
	}
	fmt.Println(serverAddr)

	server := myHttp.NewServer(serverConfig)

	logger.Log.Fatal(server.ListenAndServe())

	logger.Log.Info("Starting stocks server")
}
