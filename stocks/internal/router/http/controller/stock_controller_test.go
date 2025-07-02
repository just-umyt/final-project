package controller

import (
	// myHttp "stocks/internal/router/http"

	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"stocks/internal/models"
	"stocks/internal/router/http/controller/mock"
	"stocks/internal/usecase"
	"testing"
)

func TestAddStock(t *testing.T) {
	usecaseMock := mock.NewIStockUsecaseMock(t)
	t.Cleanup(func() {
		usecaseMock.MinimockFinish()
	})

	usecaseMock.AddStockMock.Set(func(ctx context.Context, stock usecase.AddStockDTO) error {
		if err := notFoundCheck(stock.SKUID); err != nil {
			return err
		}

		return nil
	})

	controller := NewStockController(usecaseMock)

	tests := []struct {
		name     string
		body     any
		wantCode int
	}{
		{
			name: "valid request",
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
			name:     "bad request",
			body:     `{}`,
			wantCode: http.StatusBadRequest,
		},
		{
			name: "validation",
			body: AddStockRequest{
				SKUID:    2020,
				Count:    10,
				Price:    100,
				Location: "AG",
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "not found",
			body: AddStockRequest{
				SKUID:    100,
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
			reqBody, err := json.Marshal(tt.body)
			if err != nil {
				t.Errorf("failed converting STRUCT to BYTE")
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "//stocks/item/add", bytes.NewReader(reqBody))

			controller.AddStock(w, req)

			if w.Result().StatusCode != tt.wantCode {
				t.Errorf("failed test with code :%d", w.Result().StatusCode)

				return
			}
		})

	}
}

func TestDeleteStockBySKU(t *testing.T) {
	usecaseMock := mock.NewIStockUsecaseMock(t)

	t.Cleanup(func() {
		usecaseMock.MinimockFinish()
	})

	usecaseMock.DeleteStockBySKUMock.Set(func(ctx context.Context, delStock usecase.DeleteStockDTO) (err error) {
		if err := notFoundCheck(delStock.SKUID); err != nil {
			return err
		}

		return nil
	})

	controller := NewStockController(usecaseMock)

	tests := []struct {
		name     string
		body     any
		wantCode int
	}{
		{
			name: "valid test",
			body: DeleteStockRequest{
				UserID: 1,
				SKUID:  1001,
			},
			wantCode: http.StatusOK,
		},
		{
			name:     "bad request",
			body:     `{}`,
			wantCode: http.StatusBadRequest,
		},
		{
			name: "validation",
			body: DeleteStockRequest{
				UserID: 1,
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "not found",
			body: DeleteStockRequest{
				UserID: 1,
				SKUID:  1,
			},
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, err := json.Marshal(tt.body)
			if err != nil {
				t.Errorf("failed converting STRUCT to BYTE")
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/stocks/item/delete", bytes.NewReader(reqBody))

			controller.DeleteStockBySKU(w, req)

			if w.Result().StatusCode != tt.wantCode {
				t.Errorf("failed test with code :%d", w.Result().StatusCode)

				return
			}
		})
	}
}

func TestGetItemsByLocation(t *testing.T) {
	usecaseMock := mock.NewIStockUsecaseMock(t)

	t.Cleanup(func() {
		usecaseMock.MinimockFinish()
	})

	usecaseMock.GetStocksByLocationMock.Return(usecase.ItemsByLocDTO{Stocks: []usecase.StockDTO{{Count: 1}}}, nil)

	tests := []struct {
		name     string
		body     any
		wantCode int
	}{
		{
			name: "valid",
			body: GetItemsByLocRequest{
				UserID:      1,
				Location:    "AG",
				PageSize:    1,
				CurrentPage: 1,
			},
			wantCode: http.StatusOK,
		},
		{
			name:     "bad request",
			body:     `{}`,
			wantCode: http.StatusBadRequest,
		},
		{
			name: "validation",
			body: GetItemsByLocRequest{
				UserID: 1,
			},
			wantCode: http.StatusBadRequest,
		},
	}

	controller := NewStockController(usecaseMock)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			reqBody, err := json.Marshal(tt.body)
			if err != nil {
				t.Errorf("failed converting STRUCT to BYTE")
			}

			w := httptest.NewRecorder()
			req := httptest.NewRequest(http.MethodPost, "/stocks/list/location", bytes.NewReader(reqBody))

			controller.GetItemsByLocation(w, req)

			if w.Result().StatusCode != tt.wantCode {
				t.Errorf("failed test with code :%d", w.Result().StatusCode)

				return
			}
		})
	}
}

func TestGetItemBySKU(t *testing.T) {
	usecaseMock := mock.NewIStockUsecaseMock(t)
	t.Cleanup(func() {
		usecaseMock.MinimockFinish()
	})

	usecaseMock.GetItemBySKUMock.Set(func(ctx context.Context, sku models.SKUID) (s1 usecase.StockDTO, err error) {
		if err := notFoundCheck(sku); err != nil {
			return usecase.StockDTO{}, err
		}

		return usecase.StockDTO{}, nil
	})

	controller := NewStockController(usecaseMock)

	tests := []struct {
		name     string
		body     GetItemBySKURequest
		wantCode int
	}{
		{
			name: "valid",
			body: GetItemBySKURequest{
				SKU: 1001,
			},
			wantCode: http.StatusOK,
		},
		{
			name:     "bad request",
			body:     GetItemBySKURequest{},
			wantCode: 400,
		},
		{
			name:     "not found",
			body:     GetItemBySKURequest{SKU: 1},
			wantCode: 404,
		},
	}

	for _, tt := range tests {
		reqBody, err := json.Marshal(tt.body)
		if err != nil {
			t.Errorf("failed converting STRUCT to BYTE")
		}

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/stocks/item/get", bytes.NewReader(reqBody))

		controller.GetItemBySKU(w, req)

		if w.Result().StatusCode != tt.wantCode {
			t.Errorf("failed test with code :%d", w.Result().StatusCode)

			return
		}
	}

}

func notFoundCheck(i models.SKUID) error {
	if i < 1001 || i > 10101 {
		return usecase.ErrNotFound
	}

	return nil
}
