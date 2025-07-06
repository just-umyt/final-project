package controller

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"stocks/internal/models"
	"stocks/internal/router/http/controller/mock"
	"stocks/internal/usecase"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	testValidRequestName = "ValidRequest"
	testNotFoundName     = "NotFound"
)

func TestAddStock(t *testing.T) {
	t.Parallel()

	usecaseMock := mock.NewIStockUsecaseMock(t)
	t.Cleanup(func() {
		usecaseMock.MinimockFinish()
	})

	usecaseMock.AddStockMock.Set(func(ctx context.Context, stock usecase.AddStockDTO) error {
		return idCheck(stock.SKUID)
	})

	stockController := NewStockController(usecaseMock)

	tests := []struct {
		name     string
		body     any
		wantCode int
	}{
		{
			name: testValidRequestName,
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
			name:     "BadRequest",
			body:     `{}`,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "Validation",
			body:     AddStockRequest{},
			wantCode: http.StatusBadRequest,
		},
		{
			name: testNotFoundName,
			body: AddStockRequest{
				SKUID:    1000,
				UserID:   1,
				Count:    10,
				Price:    100,
				Location: "AG",
			},
			wantCode: http.StatusNotFound,
		},
		{
			name: "SqlError",
			body: AddStockRequest{
				SKUID:    1002,
				UserID:   1,
				Count:    10,
				Price:    100,
				Location: "AG",
			},
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, req, err := generateWriterRequest(tt.body)
			require.NoError(t, err)

			stockController.AddStock(w, req)

			if w.Result().StatusCode != tt.wantCode {
				t.Errorf("failed test with code :%d", w.Result().StatusCode)
			}
		})
	}
}

func TestDeleteStockBySKU(t *testing.T) {
	t.Parallel()

	usecaseMock := mock.NewIStockUsecaseMock(t)

	t.Cleanup(func() {
		usecaseMock.MinimockFinish()
	})

	usecaseMock.DeleteStockBySKUMock.Return(nil)

	stockController := NewStockController(usecaseMock)

	tests := []struct {
		name     string
		body     any
		wantCode int
	}{
		{
			name: testValidRequestName,
			body: DeleteStockRequest{
				UserID: 1,
				SKUID:  1001,
			},
			wantCode: http.StatusOK,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, req, err := generateWriterRequest(tt.body)
			require.NoError(t, err)

			stockController.DeleteStockBySKU(w, req)

			if w.Result().StatusCode != tt.wantCode {
				t.Errorf("failed test with code :%d", w.Result().StatusCode)
			}
		})
	}
}

func TestGetItemsByLocation(t *testing.T) {
	t.Parallel()

	usecaseMock := mock.NewIStockUsecaseMock(t)

	t.Cleanup(func() {
		usecaseMock.MinimockFinish()
	})

	usecaseMock.GetStocksByLocationMock.Return(usecase.ItemsByLocDTO{Stocks: []usecase.StockDTO{{Count: 1}}}, nil)

	tests := []struct {
		name     string
		body     any
		want     string
		wantCode int
	}{
		{
			name: testValidRequestName,
			body: GetItemsByLocRequest{
				UserID:      1,
				Location:    "AG",
				PageSize:    1,
				CurrentPage: 1,
			},
			want:     `{"message":{"stocks":[{"sku":0,"name":"","type":"","count":1}],"totalCount":0,"pageNumber":0}}`,
			wantCode: http.StatusOK,
		},
	}

	stockController := NewStockController(usecaseMock)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, req, err := generateWriterRequest(tt.body)
			require.NoError(t, err)

			stockController.GetItemsByLocation(w, req)

			if w.Result().StatusCode != tt.wantCode {
				t.Errorf("failed test with code :%d", w.Result().StatusCode)
			}

			resp := strings.TrimSpace(w.Body.String())

			if resp != tt.want {
				t.Errorf("respond: %s, wanted: %s", resp, tt.want)
			}
		})
	}
}

func TestGetItemBySKU(t *testing.T) {
	t.Parallel()

	usecaseMock := mock.NewIStockUsecaseMock(t)
	t.Cleanup(func() {
		usecaseMock.MinimockFinish()
	})

	usecaseMock.GetItemBySKUMock.Set(func(ctx context.Context, sku models.SKUID) (s1 usecase.StockDTO, err error) {
		if err := idCheck(sku); err != nil {
			return usecase.StockDTO{}, err
		}

		return usecase.StockDTO{SKU: usecase.SKUDTO{SKUID: 1001}}, nil
	})

	stockController := NewStockController(usecaseMock)

	tests := []struct {
		name     string
		body     any
		want     string
		wantCode int
	}{
		{
			name:     testValidRequestName,
			body:     GetItemBySKURequest{SKU: 1001},
			want:     `{"message":{"sku":1001,"name":"","type":""}}`,
			wantCode: http.StatusOK,
		},
		{
			name:     testNotFoundName,
			body:     GetItemBySKURequest{SKU: 1000},
			want:     `{"error":"not found"}`,
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		w, req, err := generateWriterRequest(tt.body)
		require.NoError(t, err)

		stockController.GetItemBySKU(w, req)

		if w.Result().StatusCode != tt.wantCode {
			t.Errorf("failed test with code :%d", w.Result().StatusCode)
		}

		resp := strings.TrimSpace(w.Body.String())

		if resp != tt.want {
			t.Errorf("respond: %s, wanted: %s", resp, tt.want)
		}
	}
}

func generateWriterRequest(body any) (*httptest.ResponseRecorder, *http.Request, error) {
	reqBody, err := json.Marshal(body)
	if err != nil {
		return nil, nil, errors.New("failed converting STRUCT to BYTE")
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(reqBody))

	return w, req, nil
}

func idCheck(i models.SKUID) error {
	if i < 1001 {
		return usecase.ErrNotFound
	}

	if i > 1001 {
		return errors.New("sql error")
	}

	return nil
}
