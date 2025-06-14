package main

import (
	"context"
	"os"

	"stocks/internal/config"
	"stocks/internal/models"
	"stocks/internal/repository"
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

	//controllers
	// newMux := myHttp.InitControllers()

	//server

	// serverAddr := fmt.Sprintf("%s:%d", viper.GetString("server.host"), viper.GetInt("server.port"))

	// readHeaderTimeOut := time.Duration(viper.GetInt("server.readheadertimeout")) * time.Second

	// serverConfig := &myHttp.SeverConfig{
	// 	Addr:              serverAddr,
	// 	Handler:           newMux,
	// 	ReadHeaderTimeout: readHeaderTimeOut,
	// }

	// server := myHttp.NewServer(serverConfig)

	logger.Log.Info("Starting stocks server")
}
