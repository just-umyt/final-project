package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"stocks/internal/config"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	AddItemHttpReqURL    = "/stocks/item/add"
	DeleteItemHttpReqURL = "/stocks/item/delete"
	ListItemHttpReqURL   = "/stocks/list"
	GetItemHttpReqURL    = "/stocks/get"

	TestSuccessName = "Succes"
	TesNotFoundName = "NotFound"

	envPath = "../.env"
)

func TestIntegration_AddItem(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST") == "" {
		t.Skip("integration test is not set")
	}

	err := config.LoadConfig(envPath)
	require.NoError(t, err)

	init := testAppConfig{}

	err = init.Setup(t.Context())
	require.NoError(t, err)

	t.Cleanup(func() {
		err := init.Close()
		require.NoError(t, err)
	})

	tests := []struct {
		name     string
		body     any
		wantCode int
	}{
		{
			name: TestSuccessName,
			body: AddStockRequest{

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
			body: AddStockRequest{

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

			resp, err := http.Post(init.Gateway.URL+AddItemHttpReqURL, "application/json", reqBody)
			require.NoError(t, err)

			defer resp.Body.Close()

			if resp.StatusCode != tt.wantCode {
				t.Errorf("status is not correct: %d, want code: %d", resp.StatusCode, tt.wantCode)
			}
		})
	}
}

func TestIntegration_ListItems(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST") == "" {
		t.Skip("integration test is not set")
	}

	err := config.LoadConfig(envPath)
	require.NoError(t, err)

	init := testAppConfig{}

	err = init.Setup(t.Context())
	require.NoError(t, err)

	t.Cleanup(func() {
		err := init.Close()
		require.NoError(t, err)
	})

	tests := []struct {
		name     string
		body     any
		wantCode int
	}{
		{
			name: TestSuccessName,
			body: GetItemsByLocRequest{
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

			resp, err := http.Post(init.Gateway.URL+ListItemHttpReqURL, "application/json", reqBody)
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
	if os.Getenv("INTEGRATION_TEST") == "" {
		t.Skip("integration test is not set")
	}

	err := config.LoadConfig(envPath)
	require.NoError(t, err)

	init := testAppConfig{}

	err = init.Setup(t.Context())
	require.NoError(t, err)

	t.Cleanup(func() {
		err := init.Close()
		require.NoError(t, err)
	})

	tests := []struct {
		name     string
		body     any
		wantCode int
	}{
		{
			name: TestSuccessName,
			body: GetItemBySKURequest{
				SKU: 1001,
			},
			wantCode: http.StatusOK,
		},
		{
			name: TesNotFoundName,
			body: GetItemBySKURequest{
				SKU: 1000,
			},
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, err := createReqBody(tt.body)
			require.NoError(t, err)

			resp, err := http.Post(init.Gateway.URL+GetItemHttpReqURL, "application/json", reqBody)
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
	if os.Getenv("INTEGRATION_TEST") == "" {
		t.Skip("integration test is not set")
	}

	err := config.LoadConfig(envPath)
	require.NoError(t, err)

	init := testAppConfig{}

	err = init.Setup(t.Context())
	require.NoError(t, err)

	t.Cleanup(func() {
		err := init.Close()
		require.NoError(t, err)
	})

	tests := []struct {
		name     string
		body     any
		reqURL   string
		wantCode int
	}{
		{
			name: TestSuccessName,
			body: AddStockRequest{

				SKUID:    1001,
				UserID:   1,
				Count:    10,
				Price:    100,
				Location: "AG",
			},
			reqURL:   AddItemHttpReqURL,
			wantCode: http.StatusOK,
		},
		{
			name: "Succes",
			body: DeleteStockRequest{

				SKUID:  1001,
				UserID: 1,
			},
			reqURL:   DeleteItemHttpReqURL,
			wantCode: http.StatusOK,
		},
		{
			name: "NotFound",
			body: DeleteStockRequest{

				SKUID:  1001,
				UserID: 1,
			},
			reqURL:   DeleteItemHttpReqURL,
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, err := createReqBody(tt.body)
			require.NoError(t, err)

			resp, err := http.Post(init.Gateway.URL+tt.reqURL, "application/json", reqBody)
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
