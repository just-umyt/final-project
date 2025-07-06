package integration

import (
	"bytes"
	"cart/internal/config"
	"cart/internal/router/http/controller"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	AddItemHttpReqURL    = "/cart/item/add"
	DeleteItemHttpReqURL = "/cart/item/delete"
	ListItemHttpReqURL   = "/cart/list"
	ClearCartHttpReqURL  = "/cart/clear"

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
			body: controller.AddItemRequest{
				UserID: 1,
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

			resp, err := http.Post(init.Server.URL+AddItemHttpReqURL, "application/json", reqBody)
			require.NoError(t, err)

			defer resp.Body.Close()

			if resp.StatusCode != tt.wantCode {
				t.Errorf("status is not correct: %d, want code: %d", resp.StatusCode, tt.wantCode)
			}
		})
	}
}

func TestIntegration_CartList(t *testing.T) {
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
		wantMsg  bool
		wantCode int
	}{
		{
			name: TestSuccessName,
			body: controller.AddItemRequest{
				UserID: 1,
				SKUID:  1001,
				Count:  9,
			},
			reqURL:   AddItemHttpReqURL,
			wantMsg:  true,
			wantCode: http.StatusOK,
		},
		{
			name: TestSuccessName,
			body: controller.DeleteItemRequest{
				UserID: 1,
				SKUID:  1001,
			},
			reqURL:   ListItemHttpReqURL,
			wantCode: http.StatusOK,
			wantMsg:  true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, err := createReqBody(tt.body)
			require.NoError(t, err)

			resp, err := http.Post(init.Server.URL+tt.reqURL, "application/json", reqBody)
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
			body: controller.AddItemRequest{
				UserID: 1,
				SKUID:  1001,
				Count:  9,
			},
			reqURL:   AddItemHttpReqURL,
			wantCode: http.StatusOK,
		},
		{
			name: TestSuccessName,
			body: controller.DeleteItemRequest{
				UserID: 1,
				SKUID:  1001,
			},
			reqURL:   DeleteItemHttpReqURL,
			wantCode: http.StatusOK,
		},
		{
			name: TesNotFoundName,
			body: controller.DeleteItemRequest{
				UserID: 2,
				SKUID:  1001,
			},
			reqURL:   DeleteItemHttpReqURL,
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, err := createReqBody(tt.body)
			require.NoError(t, err)

			resp, err := http.Post(init.Server.URL+tt.reqURL, "application/json", reqBody)
			require.NoError(t, err)

			defer resp.Body.Close()

			if resp.StatusCode != tt.wantCode {
				t.Errorf("status is not correct: %d, want code: %d", resp.StatusCode, tt.wantCode)
			}
		})
	}
}

func TestIntegration_ClearCart(t *testing.T) {
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
			body: controller.AddItemRequest{
				UserID: 1,
				SKUID:  1001,
				Count:  9,
			},
			reqURL:   AddItemHttpReqURL,
			wantCode: http.StatusOK,
		},
		{
			name: TestSuccessName,
			body: controller.UserIDRequest{
				UserID: 1,
			},
			reqURL:   ClearCartHttpReqURL,
			wantCode: http.StatusOK,
		},
		{
			name: TesNotFoundName,
			body: controller.UserIDRequest{
				UserID: 2,
			},
			reqURL:   ClearCartHttpReqURL,
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, err := createReqBody(tt.body)
			require.NoError(t, err)

			resp, err := http.Post(init.Server.URL+tt.reqURL, "application/json", reqBody)
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
