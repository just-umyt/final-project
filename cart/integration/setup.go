package integration

import (
	"cart/internal/repository"
	"cart/internal/router/http/controller"
	"cart/internal/services"
	"cart/internal/usecase"
	"cart/pkg/postgres"
	"context"
	"database/sql"
	"net/http/httptest"
	"os"
	"strconv"
	"time"

	myHttp "cart/internal/router/http"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"
)

type testAppConfig struct {
	DB          *sql.DB
	Migration   *migrate.Migrate
	DBPool      *pgxpool.Pool
	Server      *httptest.Server
	StockServer *httptest.Server
}

func (t *testAppConfig) Setup(ctx context.Context) error {
	var err error

	dbConfig := &postgres.PostgresConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Dbname:   os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
	}

	t.DB, err = postgres.NewDB(dbConfig)
	if err != nil {
		return err
	}

	t.Migration, err = postgres.NewMigration(t.DB, os.Getenv("MIGRATION_SOURCE_URL"))
	if err != nil {
		return err
	}

	err = postgres.MigrationUp(t.Migration)
	if err != nil {
		return err
	}

	t.DBPool, err = postgres.NewDBPool(ctx, dbConfig)
	if err != nil {
		return err
	}

	trxManager := postgres.NewPgTxManager(t.DBPool)

	cartRepo := repository.NewCartRepository(t.DBPool)

	var timeOut int

	timeOut, err = strconv.Atoi(os.Getenv("CLIENT_TIMEOUT"))
	if err != nil {
		return err
	}

	t.StockServer = testStockService()

	stockService := services.NewStockService(time.Duration(timeOut)*time.Second, t.StockServer.URL)

	cartUsecase := usecase.NewCartUsecase(cartRepo, trxManager, stockService)

	cartController := controller.NewCartController(cartUsecase)

	newMux := myHttp.NewMux(cartController)

	t.Server = httptest.NewServer(newMux)

	return nil
}

func (t *testAppConfig) Close() error {
	t.DB.Close()
	t.DBPool.Close()
	t.Server.Close()
	t.StockServer.Close()

	return t.Migration.Down()
}
