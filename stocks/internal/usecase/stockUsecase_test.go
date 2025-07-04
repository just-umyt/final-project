package usecase

import (
	"context"
	"errors"
	"stocks/internal/models"
	"stocks/internal/repository"
	"stocks/internal/repository/mock"
	txMock "stocks/internal/usecase/mock"

	"testing"
)

func TestAddStock(t *testing.T) {
	repoMock := mock.NewIStockRepoMock(t)
	trxMock := txMock.NewIPgTxManagerMock(t)

	t.Cleanup(func() {
		repoMock.MinimockFinish()
		trxMock.MinimockFinish()
	})

	repoMock.GetItemBySKUMock.Set(func(ctx context.Context, skuID models.SKUID) (i1 models.Item, err error) {
		switch skuID {
		case 0:
			return models.Item{}, repository.ErrNotFound
		case 1001:
			return models.Item{SKU: models.SKU{ID: 1001}, Stock: models.Stock{UserID: 0}}, nil
		case 2020:
			return models.Item{SKU: models.SKU{ID: 2020}, Stock: models.Stock{UserID: 1}}, nil
		}

		return models.Item{Stock: models.Stock{ID: 3033}}, errors.New("sql err")
	})

	repoMock.AddStockMock.Set(func(ctx context.Context, stock models.Stock) (err error) {
		if stock.Count > 1 {
			return repository.ErrNotFound
		}

		return nil
	})

	repoMock.UpdateStockMock.Set(func(ctx context.Context, stock models.Stock) (err error) {
		if stock.Count > 1 {
			return repository.ErrNotFound
		}

		return nil
	})

	trxMock.WithTxMock.Set(func(ctx context.Context, fn func(repository.IStockRepo) error) (err error) {
		return fn(repoMock)
	})

	usecase := NewStockUsecase(repoMock, trxMock)

	tests := []struct {
		name      string
		stock     AddStockDTO
		wantError bool
	}{
		{
			name: "add",
			stock: AddStockDTO{
				SKUID:  1001,
				UserID: 0,
			},
			wantError: false,
		},
		{
			name: "add sql err",
			stock: AddStockDTO{
				SKUID:  1001,
				UserID: 0,
				Count:  2,
			},
			wantError: true,
		},
		{
			name: "update",
			stock: AddStockDTO{
				SKUID:  2020,
				UserID: 1,
			},
			wantError: false,
		},
		{
			name: "update sql err",
			stock: AddStockDTO{
				SKUID:  2020,
				UserID: 1,
				Count:  2,
			},
			wantError: true,
		},
		{
			name: "error user id",
			stock: AddStockDTO{
				SKUID:  2020,
				UserID: 3,
			},
			wantError: true,
		},
		{
			name:      "get item err",
			stock:     AddStockDTO{},
			wantError: true,
		},
		{
			name: "get item sql err",
			stock: AddStockDTO{
				SKUID: 3033,
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := usecase.AddStock(t.Context(), tt.stock)

			if (err != nil) != tt.wantError {
				t.Error(err)
			}

		})
	}
}

func TestDeleteStockBySKU(t *testing.T) {
	repoMock := mock.NewIStockRepoMock(t)
	trxMock := txMock.NewIPgTxManagerMock(t)

	t.Cleanup(func() {
		repoMock.MinimockFinish()
		trxMock.MinimockFinish()
	})

	repoMock.DeleteStockMock.Set(func(ctx context.Context, skuID models.SKUID, userID models.UserID) (err error) {
		if userID > 1 {
			return repository.ErrNotFound
		}

		return nil
	})

	usecase := NewStockUsecase(repoMock, trxMock)

	tests := []struct {
		name      string
		body      DeleteStockDTO
		wantError bool
	}{
		{
			name: "valid",
			body: DeleteStockDTO{
				UserID: 1,
				SKUID:  1001,
			},
			wantError: false,
		},
		{
			name: "sql err",
			body: DeleteStockDTO{
				UserID: 2,
				SKUID:  1001,
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := usecase.DeleteStockBySKU(t.Context(), tt.body)

			if (err != nil) != tt.wantError {
				t.Error(err)
			}
		})
	}
}

func TestGetStockByLocation(t *testing.T) {
	repoMock := mock.NewIStockRepoMock(t)
	trxMock := txMock.NewIPgTxManagerMock(t)

	t.Cleanup(func() {
		repoMock.MinimockFinish()
		trxMock.MinimockFinish()
	})

	repoMock.GetItemsByLocationMock.Set(func(ctx context.Context, param repository.GetStockByLocation) ([]models.Item, error) {
		if param.UserID > 1 {
			return []models.Item{}, errors.New("sql error")
		}

		return []models.Item{
			{
				SKU:   models.SKU{},
				Stock: models.Stock{},
			},
		}, nil
	})

	trxMock.WithTxMock.Set(func(ctx context.Context, fn func(repository.IStockRepo) error) (err error) {
		return fn(repoMock)
	})

	usecase := NewStockUsecase(repoMock, trxMock)

	tests := []struct {
		name      string
		body      GetItemByLocDTO
		wantError bool
	}{
		{
			name: "valid",
			body: GetItemByLocDTO{
				UserID: 1,
			},
			wantError: false,
		},
		{
			name: "sql err",
			body: GetItemByLocDTO{
				UserID: 2,
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := usecase.GetStocksByLocation(t.Context(), tt.body)
			if (err != nil) != tt.wantError {
				t.Error(err)
			}
		})
	}

}

func TestGetItemBySKU(t *testing.T) {
	repoMock := mock.NewIStockRepoMock(t)
	trxMock := txMock.NewIPgTxManagerMock(t)

	t.Cleanup(func() {
		repoMock.MinimockFinish()
		trxMock.MinimockFinish()
	})

	repoMock.GetItemBySKUMock.Set(func(ctx context.Context, skuID models.SKUID) (models.Item, error) {
		switch skuID {
		case 1001:
			return models.Item{}, nil
		case 2020:
			return models.Item{SKU: models.SKU{ID: skuID}}, errors.New("not found")
		default:
			return models.Item{}, errors.New("sql error")
		}
	})

	trxMock.WithTxMock.Set(func(ctx context.Context, fn func(repository.IStockRepo) error) (err error) {
		return fn(repoMock)
	})

	usecase := NewStockUsecase(repoMock, trxMock)

	tests := []struct {
		name    string
		body    models.SKUID
		wantErr bool
	}{
		{
			name:    "valid",
			body:    1001,
			wantErr: false,
		},
		{
			name:    "not found",
			body:    2020,
			wantErr: true,
		},
		{
			name:    "sql err",
			body:    3003,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := usecase.GetItemBySKU(t.Context(), tt.body)
			if (err != nil) != tt.wantErr {
				t.Error(err)
			}
		})
	}
}
