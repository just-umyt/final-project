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
	"testing"
)

func TestAddStock(t *testing.T) {
	usecaseMock := mock.NewIStockUsecaseMock(t)
	t.Cleanup(func() {
		usecaseMock.MinimockFinish()
	})

	usecaseMock.AddStockMock.Set(func(ctx context.Context, stock usecase.AddStockDTO) error {
		return idCheck(stock.SKUID)
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
		{
			name: "sql error",
			body: AddStockRequest{
				SKUID:    20000,
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
			if err != nil {
				t.Error(err)
			}

			controller.AddStock(w, req)

			if w.Result().StatusCode != tt.wantCode {
				t.Errorf("failed test with code :%d", w.Result().StatusCode)
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
		return idCheck(delStock.SKUID)
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
		{
			name: "sql error",
			body: DeleteStockRequest{
				UserID: 1,
				SKUID:  20000,
			},
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, req, err := generateWriterRequest(tt.body)
			if err != nil {
				t.Error(err)
			}

			controller.DeleteStockBySKU(w, req)

			if w.Result().StatusCode != tt.wantCode {
				t.Errorf("failed test with code :%d", w.Result().StatusCode)
			}
		})
	}
}

func TestGetItemsByLocation(t *testing.T) {
	usecaseMock := mock.NewIStockUsecaseMock(t)

	t.Cleanup(func() {
		usecaseMock.MinimockFinish()
	})

	usecaseMock.GetStocksByLocationMock.Set(func(ctx context.Context, param usecase.GetItemByLocDTO) (usecase.ItemsByLocDTO, error) {
		if param.UserID != 1 {
			return usecase.ItemsByLocDTO{}, errors.New("sql err")
		}

		return usecase.ItemsByLocDTO{Stocks: []usecase.StockDTO{{Count: 1}}}, nil
	})

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
		{
			name: "sql err",
			body: GetItemsByLocRequest{
				UserID:      2,
				Location:    "AG",
				PageSize:    1,
				CurrentPage: 1,
			},
			wantCode: http.StatusInternalServerError,
		},
	}

	controller := NewStockController(usecaseMock)

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, req, err := generateWriterRequest(tt.body)
			if err != nil {
				t.Error(err)
			}

			controller.GetItemsByLocation(w, req)

			if w.Result().StatusCode != tt.wantCode {
				t.Errorf("failed test with code :%d", w.Result().StatusCode)
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
		if err := idCheck(sku); err != nil {
			return usecase.StockDTO{}, err
		}

		return usecase.StockDTO{}, nil
	})

	controller := NewStockController(usecaseMock)

	tests := []struct {
		name     string
		body     any
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
			body:     `{}`,
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "not found",
			body:     GetItemBySKURequest{SKU: 1},
			wantCode: http.StatusNotFound,
		},
		{
			name: "sql err",
			body: GetItemBySKURequest{
				SKU: 20000,
			},
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		reqBody, err := json.Marshal(tt.body)
		if err != nil {
			t.Error(err)
		}

		w := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/stocks/item/get", bytes.NewReader(reqBody))

		controller.GetItemBySKU(w, req)

		if w.Result().StatusCode != tt.wantCode {
			t.Errorf("failed test with code :%d", w.Result().StatusCode)
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

	if i > 10101 {
		return errors.New("sql error")
	}

	return nil
}
