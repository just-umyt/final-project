package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"stocks/internal/config"
	"stocks/internal/repository"
	"stocks/internal/router/http/controller"
	"stocks/internal/usecase"
	"stocks/pkg/postgres"
	"testing"

	myHttp "stocks/internal/router/http"

	"github.com/stretchr/testify/require"
)

var (
	AddItemHttpReqURL    = "/stocks/item/add"
	DeleteItemHttpReqURL = "/stocks/item/delete"
	ListItemHttpReqURL   = "/stocks/list/location"
	GetItemHttpReqURL    = "/stocks/item/get"

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

	stockRepo := repository.NewStockRepository(dbPool)

	stockUsecae := usecase.NewStockUsecase(stockRepo, trxManager)

	stockController := controller.NewStockController(stockUsecae)

	newMux := myHttp.NewMux(stockController)

	server := httptest.NewServer(newMux)

	t.Cleanup(func() {
		db.Close()
		dbPool.Close()
		server.Close()
	})

	tests := []struct {
		name     string
		body     any
		wantCode int
	}{
		{
			name: TestSuccessName,
			body: controller.AddStockRequest{

				SKUID:    1001,
				UserID:   1,
				Count:    10,
				Price:    100,
				Location: "AG",
			},
			wantCode: http.StatusOK,
		},
		{
			name: TesNotFoundName,
			body: controller.AddStockRequest{

				SKUID:    1000,
				UserID:   1,
				Count:    10,
				Price:    100,
				Location: "AG",
			},
			wantCode: http.StatusNotFound,
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

func TestIntegration_ListItems(t *testing.T) {
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

	stockRepo := repository.NewStockRepository(dbPool)

	stockUsecae := usecase.NewStockUsecase(stockRepo, trxManager)

	stockController := controller.NewStockController(stockUsecae)

	newMux := myHttp.NewMux(stockController)

	server := httptest.NewServer(newMux)

	t.Cleanup(func() {
		db.Close()
		dbPool.Close()
		server.Close()
	})

	tests := []struct {
		name     string
		body     any
		wantCode int
	}{
		{
			name: TestSuccessName,
			body: controller.GetItemsByLocRequest{
				UserID:      1,
				Location:    "AG",
				PageSize:    1,
				CurrentPage: 1,
			},
			wantCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, err := createReqBody(tt.body)
			require.NoError(t, err)

			resp, err := http.Post(server.URL+ListItemHttpReqURL, "application/json", reqBody)
			require.NoError(t, err)
			defer resp.Body.Close()

			if resp.Body == nil {
				t.Errorf("response body is nil")
			}

			if resp.StatusCode != tt.wantCode {
				t.Errorf("status is not correct: %d, want code: %d", resp.StatusCode, tt.wantCode)
			}
		})
	}
}

func TestIntegration_GetItem(t *testing.T) {
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

	stockRepo := repository.NewStockRepository(dbPool)

	stockUsecae := usecase.NewStockUsecase(stockRepo, trxManager)

	stockController := controller.NewStockController(stockUsecae)

	newMux := myHttp.NewMux(stockController)

	server := httptest.NewServer(newMux)

	t.Cleanup(func() {
		db.Close()
		dbPool.Close()
		server.Close()
	})

	tests := []struct {
		name     string
		body     any
		wantCode int
	}{
		{
			name: TestSuccessName,
			body: controller.GetItemBySKURequest{
				SKU: 1001,
			},
			wantCode: http.StatusOK,
		},
		{
			name: TesNotFoundName,
			body: controller.GetItemBySKURequest{
				SKU: 1000,
			},
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, err := createReqBody(tt.body)
			require.NoError(t, err)

			resp, err := http.Post(server.URL+GetItemHttpReqURL, "application/json", reqBody)
			require.NoError(t, err)
			defer resp.Body.Close()

			if resp.Body == nil {
				t.Errorf("response body is nil")
			}

			if resp.StatusCode != tt.wantCode {
				t.Errorf("status is not correct: %d, want code: %d", resp.StatusCode, tt.wantCode)
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

	stockRepo := repository.NewStockRepository(dbPool)

	stockUsecae := usecase.NewStockUsecase(stockRepo, trxManager)

	stockController := controller.NewStockController(stockUsecae)

	newMux := myHttp.NewMux(stockController)

	server := httptest.NewServer(newMux)

	t.Cleanup(func() {
		db.Close()
		dbPool.Close()
		err := migration.Down()
		require.NoError(t, err)
		server.Close()
	})

	tests := []struct {
		name     string
		body     any
		wantCode int
	}{
		{
			name: "Succes",
			body: controller.DeleteStockRequest{

				SKUID:  1001,
				UserID: 1,
			},
			wantCode: http.StatusOK,
		},
		{
			name: "NotFound",
			body: controller.DeleteStockRequest{

				SKUID:  1001,
				UserID: 1,
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

func createReqBody(data any) (io.Reader, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	return bytes.NewBuffer(body), nil
}
