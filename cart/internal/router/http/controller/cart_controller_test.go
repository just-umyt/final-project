package controller

import (
	"bytes"
	"cart/internal/models"
	"cart/internal/router/http/controller/mock"
	"cart/internal/usecase"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	testValidRequestName        = "ValidRequest"
	testNotFoundName            = "NotFound"
	testInternalServerErrorName = "ErrorInternalServer"
)

var (
	errInternalServer error = errors.New("nternal server error")
)

func TestAddItem(t *testing.T) {
	t.Parallel()

	usecaseMock := mock.NewICartUsecaseMock(t)
	t.Cleanup(func() { usecaseMock.MinimockFinish() })

	usecaseMock.AddItemMock.Set(func(ctx context.Context, addItem usecase.AddItemDTO) (err error) {
		if addItem.Count > 5 {
			return usecase.ErrNotEnoughStock
		} else if addItem.Count < 5 {
			return errInternalServer
		}

		return nil
	})

	cartController := NewCartController(usecaseMock)

	tests := []struct {
		name     string
		body     any
		wantCode int
	}{
		{
			name: testValidRequestName,
			body: AddItemRequest{
				UserID: 1,
				SKUID:  1001,
				Count:  5,
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
			body:     AddItemRequest{},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "NotEnough",
			body: AddItemRequest{
				UserID: 1,
				SKUID:  100,
				Count:  6,
			},
			wantCode: http.StatusPreconditionFailed,
		},
		{
			name: testInternalServerErrorName,
			body: AddItemRequest{
				UserID: 1,
				SKUID:  100,
				Count:  4,
			},
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, req, err := generateWriterRequest(tt.body)
			require.NoError(t, err)

			cartController.AddItem(w, req)

			if w.Result().StatusCode != tt.wantCode {
				t.Errorf("failed test with code :%d", w.Result().StatusCode)
			}
		})
	}
}

func TestCartClear(t *testing.T) {
	t.Parallel()

	usecaseMock := mock.NewICartUsecaseMock(t)
	t.Cleanup(func() { usecaseMock.MinimockFinish() })

	usecaseMock.ClearCartByUserIDMock.Set(func(ctx context.Context, userID models.UserID) (err error) {
		if userID != 1 {
			return usecase.ErrNotFound
		}

		return nil
	})

	cartController := NewCartController(usecaseMock)

	tests := []struct {
		name     string
		body     any
		wantCode int
	}{
		{
			name:     testValidRequestName,
			body:     UserIDRequest{UserID: 1},
			wantCode: http.StatusOK,
		},
		{
			name:     testNotFoundName,
			body:     UserIDRequest{UserID: 0},
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, req, err := generateWriterRequest(tt.body)
			require.NoError(t, err)

			cartController.CartClear(w, req)

			if w.Result().StatusCode != tt.wantCode {
				t.Errorf("failed test with code :%d", w.Result().StatusCode)
			}
		})
	}
}

func TestDeleteItem(t *testing.T) {
	t.Parallel()

	usecaseMock := mock.NewICartUsecaseMock(t)

	t.Cleanup(func() { usecaseMock.MinimockFinish() })

	usecaseMock.DeleteItemMock.Set(func(ctx context.Context, delItem usecase.DeleteItemDTO) error {
		if delItem.SKUID != 1001 {
			return usecase.ErrNotFound
		}

		return nil
	})

	cartController := NewCartController(usecaseMock)

	tests := []struct {
		name     string
		body     any
		wantCode int
	}{
		{
			name: testValidRequestName,
			body: DeleteItemRequest{
				UserID: 1,
				SKUID:  1001,
			},
			wantCode: http.StatusOK,
		},
		{
			name: testNotFoundName,
			body: DeleteItemRequest{
				UserID: 1,
				SKUID:  2020,
			},
			wantCode: http.StatusNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, req, err := generateWriterRequest(tt.body)
			require.NoError(t, err)

			cartController.DeleteItem(w, req)

			if w.Result().StatusCode != tt.wantCode {
				t.Errorf("failed test with code :%d", w.Result().StatusCode)
			}
		})
	}
}

func TestCartList(t *testing.T) {
	t.Parallel()

	usecaseMock := mock.NewICartUsecaseMock(t)

	t.Cleanup(func() { usecaseMock.MinimockFinish() })

	usecaseMock.GetItemsByUserIDMock.Set(func(ctx context.Context, userID models.UserID) (usecase.ListItemsDTO, error) {
		if userID > 1 {
			return usecase.ListItemsDTO{}, errInternalServer
		}

		return usecase.ListItemsDTO{}, nil
	})

	cartController := NewCartController(usecaseMock)

	tests := []struct {
		name     string
		body     any
		wantCode int
	}{
		{
			name:     testValidRequestName,
			body:     UserIDRequest{UserID: 1},
			wantCode: http.StatusOK,
		},
		{
			name:     testInternalServerErrorName,
			body:     UserIDRequest{UserID: 2},
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w, req, err := generateWriterRequest(tt.body)
			require.NoError(t, err)

			cartController.CartList(w, req)

			if w.Result().StatusCode != tt.wantCode {
				t.Errorf("failed test with code :%d", w.Result().StatusCode)
			}
		})
	}
}

func generateWriterRequest(body any) (*httptest.ResponseRecorder, *http.Request, error) {
	reqBody, err := json.Marshal(body)
	if err != nil {
		return nil, nil, errors.New("failed converting ANY to BYTE")
	}

	w := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/", bytes.NewReader(reqBody))

	return w, req, nil
}
