package usecase

import (
	"cart/internal/models"
	"cart/internal/repository"
	repoMock "cart/internal/repository/mock"
	"cart/internal/services"
	"cart/internal/usecase/mock"
	"context"
	"errors"

	"testing"
)

const (
	testSuccesName   = "Succes"
	testNotFoundName = "NotFound"
)

func TestAddItem(t *testing.T) {
	t.Parallel()

	serviceMock := mock.NewIStockServiceMock(t)
	repoMock := repoMock.NewICartRepoMock(t)
	trxMock := mock.NewIPgTxManagerMock(t)

	t.Cleanup(func() {
		repoMock.MinimockFinish()
		trxMock.MinimockFinish()
		serviceMock.MinimockFinish()
	})

	serviceMock.GetItemInfoMock.Set(func(ctx context.Context, skuID models.SKUID) (services.ItemDTO, error) {
		if skuID < 1001 {
			return services.ItemDTO{}, errors.New("not found")
		} else if skuID > 1001 {
			return services.ItemDTO{Count: 1}, nil
		}

		return services.ItemDTO{Count: 10}, nil
	})

	repoMock.UpdateItemByUserIDMock.Set(func(ctx context.Context, cart models.Cart) (err error) {
		if cart.UserID != 1 {
			return repository.ErrNotFound
		}

		return nil
	})

	repoMock.AddItemMock.Return(nil)

	trxMock.WithTxMock.Set(func(ctx context.Context, fn func(repository.ICartRepo) error) (err error) {
		return fn(repoMock)
	})

	cartUsecase := NewCartUsecase(repoMock, trxMock, serviceMock)

	tests := []struct {
		name    string
		body    AddItemDTO
		wantErr bool
	}{
		{
			name: testSuccesName,
			body: AddItemDTO{
				UserID: 1,
				SKUID:  1001,
				Count:  5,
			},
			wantErr: false,
		},
		{
			name: "ErrorGetInfoCheck",
			body: AddItemDTO{
				UserID: 1,
				SKUID:  1000,
				Count:  5,
			},
			wantErr: true,
		},
		{
			name: "ErrorNotEnoughStock",
			body: AddItemDTO{
				UserID: 1,
				SKUID:  1002,
				Count:  5,
			},
			wantErr: true,
		},
		{
			name: "ErrorUpdateNotfound",
			body: AddItemDTO{
				UserID: 2,
				SKUID:  1001,
				Count:  5,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cartUsecase.AddItem(t.Context(), tt.body)
			if (err != nil) != tt.wantErr {
				t.Error(err)
			}
		})
	}
}

func TestDeleteItem(t *testing.T) {
	t.Parallel()

	serviceMock := mock.NewIStockServiceMock(t)
	repoMock := repoMock.NewICartRepoMock(t)
	trxMock := mock.NewIPgTxManagerMock(t)

	t.Cleanup(func() {
		repoMock.MinimockFinish()
		trxMock.MinimockFinish()
		serviceMock.MinimockFinish()
	})

	repoMock.DeleteItemMock.Set(func(ctx context.Context, userID models.UserID, skuID models.SKUID) (err error) {
		if skuID != 1001 {
			return repository.ErrNotFound
		}

		return nil
	})

	cartUsecase := NewCartUsecase(repoMock, trxMock, serviceMock)

	tests := []struct {
		name    string
		body    DeleteItemDTO
		wantErr bool
	}{
		{
			name: testSuccesName,
			body: DeleteItemDTO{
				UserID: 1,
				SKUID:  1001,
			},
			wantErr: false,
		},
		{
			name: testNotFoundName,
			body: DeleteItemDTO{
				UserID: 1,
				SKUID:  1,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cartUsecase.DeleteItem(t.Context(), tt.body)
			if (err != nil) != tt.wantErr {
				t.Error(err)
			}
		})
	}
}

func TestGetItemsByUserID(t *testing.T) {
	t.Parallel()

	serviceMock := mock.NewIStockServiceMock(t)
	repoMock := repoMock.NewICartRepoMock(t)
	trxMock := mock.NewIPgTxManagerMock(t)

	t.Cleanup(func() {
		repoMock.MinimockFinish()
		trxMock.MinimockFinish()
		serviceMock.MinimockFinish()
	})

	repoMock.GetCartByUserIDMock.Set(func(ctx context.Context, userID models.UserID) (ca1 []models.CartItem, err error) {
		if userID > 1 {
			return []models.CartItem{}, errors.New("sql error")
		}

		return []models.CartItem{{SKUID: models.SKUID(1001), Count: 10}}, nil
	})

	serviceMock.GetItemInfoMock.Return(services.ItemDTO{}, nil)

	cartUsecase := NewCartUsecase(repoMock, trxMock, serviceMock)

	tests := []struct {
		name    string
		body    models.UserID
		wantErr bool
	}{
		{
			name:    testNotFoundName,
			body:    1,
			wantErr: false,
		},
		{
			name:    "SqlError",
			body:    2,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := cartUsecase.GetItemsByUserID(t.Context(), tt.body)
			if (err != nil) != tt.wantErr {
				t.Error(err)
			}
		})
	}
}

func TestClearCartByUserID(t *testing.T) {
	t.Parallel()

	serviceMock := mock.NewIStockServiceMock(t)
	repoMock := repoMock.NewICartRepoMock(t)
	trxMock := mock.NewIPgTxManagerMock(t)

	t.Cleanup(func() {
		repoMock.MinimockFinish()
		trxMock.MinimockFinish()
		serviceMock.MinimockFinish()
	})

	repoMock.ClearCartByUserIDMock.Set(func(ctx context.Context, userID models.UserID) (err error) {
		if userID != 1 {
			return repository.ErrNotFound
		}

		return nil
	})

	trxMock.WithTxMock.Set(func(ctx context.Context, fn func(repository.ICartRepo) error) (err error) { return fn(repoMock) })

	cartUsecase := NewCartUsecase(repoMock, trxMock, serviceMock)

	tests := []struct {
		name    string
		body    models.UserID
		wantErr bool
	}{
		{
			name:    testSuccesName,
			body:    1,
			wantErr: false,
		},
		{
			name:    testNotFoundName,
			body:    2,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cartUsecase.ClearCartByUserID(t.Context(), tt.body)
			if (err != nil) != tt.wantErr {
				t.Error(err)
			}
		})
	}
}
