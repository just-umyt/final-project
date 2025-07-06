package integration

import (
	"context"
	"database/sql"
	"net/http/httptest"
	"os"
	"stocks/internal/repository"
	"stocks/internal/router/http/controller"
	"stocks/internal/usecase"
	"stocks/pkg/postgres"

	myHttp "stocks/internal/router/http"

	"github.com/golang-migrate/migrate/v4"
	"github.com/jackc/pgx/v5/pgxpool"
)

type testAppConfig struct {
	DB        *sql.DB
	Migration *migrate.Migrate
	DBPool    *pgxpool.Pool
	Server    *httptest.Server
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

	stockRepo := repository.NewStockRepository(t.DBPool)

	stockUsecae := usecase.NewStockUsecase(stockRepo, trxManager)

	stockController := controller.NewStockController(stockUsecae)

	newMux := myHttp.NewMux(stockController)

	t.Server = httptest.NewServer(newMux)

	return nil
}

func (t *testAppConfig) Close() error {
	t.DB.Close()
	t.DBPool.Close()
	t.Server.Close()

	return t.Migration.Down()
}
