package integration

import (
	"bytes"
	"cart/internal/config"
	"cart/internal/repository"
	myHttp "cart/internal/router/http"
	"cart/internal/router/http/controller"
	"cart/internal/services"
	"cart/internal/usecase"
	"cart/pkg/postgres"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

var (
	AddItemHttpReqURL    = "/cart/item/add"
	DeleteItemHttpReqURL = "/cart/item/delete"
	ListItemHttpReqURL   = "/cart/list"
	ClearCartHttpReqURL  = "/cart/clear"

	TestSuccessName = "Succes"
	TesNotFoundName = "NotFound"
)

func TestIntegration_AddItem(t *testing.T) {
	err := config.LoadConfig("../.env")
	require.NoError(t, err)

	if os.Getenv("INTEGRATION_TEST") == "" {
		t.Skip("integration test is not set")
	}

	dbConfig := &postgres.PostgresConfig{
		Host:     os.Getenv("TEST_DB_HOST"),
		Port:     os.Getenv("TEST_DB_PORT"),
		User:     os.Getenv("TEST_DB_USER"),
		Password: os.Getenv("TEST_DB_PASSWORD"),
		Dbname:   os.Getenv("TEST_DB_NAME"),
		SSLMode:  os.Getenv("TEST_DB_SSLMODE"),
	}

	db, err := postgres.NewDB(dbConfig)
	require.NoError(t, err)

	migration, err := postgres.NewMigration(db, os.Getenv("TEST_MIGRATION_SOURCE_URL"))
	require.NoError(t, err)

	err = postgres.MigrationUp(migration)
	require.NoError(t, err)

	dbPool, err := postgres.NewDBPool(t.Context(), dbConfig)
	require.NoError(t, err)

	trxManager := postgres.NewPgTxManager(dbPool)

	cartRepo := repository.NewCartRepository(dbPool)

	timeOut, err := strconv.Atoi(os.Getenv("TEST_CLIENT_TIMEOUT"))
	require.NoError(t, err)

	stockServer := testStockService()

	stockService := services.NewStockService(time.Duration(timeOut)*time.Second, stockServer.URL)

	cartUsecase := usecase.NewCartUsecase(cartRepo, trxManager, stockService)

	cartController := controller.NewCartController(cartUsecase)

	newMux := myHttp.NewMux(cartController)

	server := httptest.NewServer(newMux)

	t.Cleanup(func() {
		db.Close()
		dbPool.Close()
		server.Close()
		stockServer.Close()
	})

	tests := []struct {
		name     string
		body     any
		wantCode int
	}{
		{
			name: TestSuccessName,
			body: controller.AddItemRequest{
				UserID: 1,
				SKUID:  1001,
				Count:  9,
			},
			wantCode: http.StatusOK,
		},
		{
			name: TestSuccessName,
			body: controller.AddItemRequest{
				UserID: 2,
				SKUID:  1001,
				Count:  9,
			},
			wantCode: http.StatusOK,
		},
		{
			name: "NotEnoughStock",
			body: controller.AddItemRequest{
				UserID: 1,
				SKUID:  1001,
				Count:  11,
			},
			wantCode: http.StatusPreconditionFailed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, err := createReqBody(tt.body)
			require.NoError(t, err)

			resp, err := http.Post(server.URL+AddItemHttpReqURL, "application/json", reqBody)
			require.NoError(t, err)

			defer resp.Body.Close()

			if resp.StatusCode != tt.wantCode {
				t.Errorf("status is not correct: %d, want code: %d", resp.StatusCode, tt.wantCode)
			}
		})
	}
}

func TestIntegration_CartList(t *testing.T) {
	err := config.LoadConfig("../.env")
	require.NoError(t, err)

	if os.Getenv("INTEGRATION_TEST") == "" {
		t.Skip("integration test is not set")
	}

	dbConfig := &postgres.PostgresConfig{
		Host:     os.Getenv("TEST_DB_HOST"),
		Port:     os.Getenv("TEST_DB_PORT"),
		User:     os.Getenv("TEST_DB_USER"),
		Password: os.Getenv("TEST_DB_PASSWORD"),
		Dbname:   os.Getenv("TEST_DB_NAME"),
		SSLMode:  os.Getenv("TEST_DB_SSLMODE"),
	}

	db, err := postgres.NewDB(dbConfig)
	require.NoError(t, err)

	migration, err := postgres.NewMigration(db, os.Getenv("TEST_MIGRATION_SOURCE_URL"))
	require.NoError(t, err)

	err = postgres.MigrationUp(migration)
	require.NoError(t, err)

	dbPool, err := postgres.NewDBPool(t.Context(), dbConfig)
	require.NoError(t, err)

	trxManager := postgres.NewPgTxManager(dbPool)

	cartRepo := repository.NewCartRepository(dbPool)

	timeOut, err := strconv.Atoi(os.Getenv("TEST_CLIENT_TIMEOUT"))
	require.NoError(t, err)

	stockServer := testStockService()

	stockService := services.NewStockService(time.Duration(timeOut)*time.Second, stockServer.URL)

	cartUsecase := usecase.NewCartUsecase(cartRepo, trxManager, stockService)

	cartController := controller.NewCartController(cartUsecase)

	newMux := myHttp.NewMux(cartController)

	server := httptest.NewServer(newMux)

	t.Cleanup(func() {
		db.Close()
		dbPool.Close()
		server.Close()
		stockServer.Close()
	})

	tests := []struct {
		name     string
		body     any
		wantMsg  bool
		wantCode int
	}{
		{
			name: TestSuccessName,
			body: controller.DeleteItemRequest{
				UserID: 1,
				SKUID:  1001,
			},
			wantCode: http.StatusOK,
			wantMsg:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, err := createReqBody(tt.body)
			require.NoError(t, err)

			resp, err := http.Post(server.URL+ListItemHttpReqURL, "application/json", reqBody)
			require.NoError(t, err)

			defer resp.Body.Close()

			if resp.StatusCode != tt.wantCode {
				t.Errorf("status is not correct: %d, want code: %d", resp.StatusCode, tt.wantCode)
			}

			if (resp.Body != nil) != tt.wantMsg {
				t.Errorf("wantMsg is not correct: %v, want msg: %v", resp.Body, tt.wantMsg)
			}
		})
	}
}

func TestIntegration_DeleteItem(t *testing.T) {
	err := config.LoadConfig("../.env")
	require.NoError(t, err)

	if os.Getenv("INTEGRATION_TEST") == "" {
		t.Skip("integration test is not set")
	}

	dbConfig := &postgres.PostgresConfig{
		Host:     os.Getenv("TEST_DB_HOST"),
		Port:     os.Getenv("TEST_DB_PORT"),
		User:     os.Getenv("TEST_DB_USER"),
		Password: os.Getenv("TEST_DB_PASSWORD"),
		Dbname:   os.Getenv("TEST_DB_NAME"),
		SSLMode:  os.Getenv("TEST_DB_SSLMODE"),
	}

	db, err := postgres.NewDB(dbConfig)
	require.NoError(t, err)

	migration, err := postgres.NewMigration(db, os.Getenv("TEST_MIGRATION_SOURCE_URL"))
	require.NoError(t, err)

	err = postgres.MigrationUp(migration)
	require.NoError(t, err)

	dbPool, err := postgres.NewDBPool(t.Context(), dbConfig)
	require.NoError(t, err)

	trxManager := postgres.NewPgTxManager(dbPool)

	cartRepo := repository.NewCartRepository(dbPool)

	timeOut, err := strconv.Atoi(os.Getenv("TEST_CLIENT_TIMEOUT"))
	require.NoError(t, err)

	stockServer := testStockService()

	stockService := services.NewStockService(time.Duration(timeOut)*time.Second, stockServer.URL)

	cartUsecase := usecase.NewCartUsecase(cartRepo, trxManager, stockService)

	cartController := controller.NewCartController(cartUsecase)

	newMux := myHttp.NewMux(cartController)

	server := httptest.NewServer(newMux)

	t.Cleanup(func() {
		db.Close()
		dbPool.Close()
		server.Close()
		stockServer.Close()
	})

	tests := []struct {
		name     string
		body     any
		wantCode int
	}{
		{
			name: TestSuccessName,
			body: controller.DeleteItemRequest{
				UserID: 1,
				SKUID:  1001,
			},
			wantCode: http.StatusOK,
		},
		{
			name: TesNotFoundName,
			body: controller.DeleteItemRequest{
				UserID: 3,
				SKUID:  1001,
			},
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, err := createReqBody(tt.body)
			require.NoError(t, err)

			resp, err := http.Post(server.URL+DeleteItemHttpReqURL, "application/json", reqBody)
			require.NoError(t, err)

			defer resp.Body.Close()

			if resp.StatusCode != tt.wantCode {
				t.Errorf("status is not correct: %d, want code: %d", resp.StatusCode, tt.wantCode)
			}

		})
	}
}

func TestIntegration_ClearCart(t *testing.T) {
	err := config.LoadConfig("../.env")
	require.NoError(t, err)

	if os.Getenv("INTEGRATION_TEST") == "" {
		t.Skip("integration test is not set")
	}

	dbConfig := &postgres.PostgresConfig{
		Host:     os.Getenv("TEST_DB_HOST"),
		Port:     os.Getenv("TEST_DB_PORT"),
		User:     os.Getenv("TEST_DB_USER"),
		Password: os.Getenv("TEST_DB_PASSWORD"),
		Dbname:   os.Getenv("TEST_DB_NAME"),
		SSLMode:  os.Getenv("TEST_DB_SSLMODE"),
	}

	db, err := postgres.NewDB(dbConfig)
	require.NoError(t, err)

	migration, err := postgres.NewMigration(db, os.Getenv("TEST_MIGRATION_SOURCE_URL"))
	require.NoError(t, err)

	err = postgres.MigrationUp(migration)
	require.NoError(t, err)

	dbPool, err := postgres.NewDBPool(t.Context(), dbConfig)
	require.NoError(t, err)

	trxManager := postgres.NewPgTxManager(dbPool)

	cartRepo := repository.NewCartRepository(dbPool)

	timeOut, err := strconv.Atoi(os.Getenv("TEST_CLIENT_TIMEOUT"))
	require.NoError(t, err)

	stockServer := testStockService()

	stockService := services.NewStockService(time.Duration(timeOut)*time.Second, stockServer.URL)

	cartUsecase := usecase.NewCartUsecase(cartRepo, trxManager, stockService)

	cartController := controller.NewCartController(cartUsecase)

	newMux := myHttp.NewMux(cartController)

	server := httptest.NewServer(newMux)

	t.Cleanup(func() {
		db.Close()
		dbPool.Close()
		server.Close()
		stockServer.Close()
		require.NoError(t, migration.Down())
	})

	tests := []struct {
		name     string
		body     any
		wantCode int
	}{
		{
			name: TestSuccessName,
			body: controller.UserIDRequest{
				UserID: 2,
			},
			wantCode: http.StatusOK,
		},
		{
			name: TesNotFoundName,
			body: controller.UserIDRequest{
				UserID: 3,
			},
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, err := createReqBody(tt.body)
			require.NoError(t, err)

			resp, err := http.Post(server.URL+ClearCartHttpReqURL, "application/json", reqBody)
			require.NoError(t, err)

			defer resp.Body.Close()

			if resp.StatusCode != tt.wantCode {
				t.Errorf("status is not correct: %d, want code: %d", resp.StatusCode, tt.wantCode)
			}

		})
	}
}

func createReqBody(data any) (io.Reader, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(body), nil
}
