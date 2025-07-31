package usecase

import (
	"context"
	"errors"
	"stocks/internal/models"
	logMock "stocks/internal/observability/log/mock"
	"stocks/internal/repository"
	repositoryMock "stocks/internal/repository/mock"
	"stocks/internal/usecase/mock"

	"testing"
)

const (
	testSuccesName   = "Succes"
	testSqlErrorName = "ErrorSqlGetItem"
)

var (
	errSql = errors.New("sql error")
)

func TestAddStock(t *testing.T) {
	repoMock := repositoryMock.NewIStockRepoMock(t)
	trxMock := mock.NewIPgTxManagerMock(t)
	kafkaMock := mock.NewIProducerMock(t)
	logger := logMock.NewLoggerMock(t)

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

		return models.Item{Stock: models.Stock{ID: 3033}}, errSql
	})

	repoMock.AddStockMock.Return(nil)

	repoMock.UpdateStockMock.Return(nil)

	kafkaMock.ProduceMock.Return(nil)

	trxMock.WithTxMock.Set(func(ctx context.Context, fn func(repository.IStockRepo) error) (err error) {
		return fn(repoMock)
	})

	logger.InfoMock.Return()

	usecase := NewStockUsecase(repoMock, trxMock, kafkaMock, logger)

	tests := []struct {
		name    string
		stock   AddStockDTO
		wantErr error
	}{
		{
			name:    "ErrorGetItem",
			stock:   AddStockDTO{},
			wantErr: ErrNotFound,
		},
		{
			name: "ErrorSqlGetItem",
			stock: AddStockDTO{
				SKUID: 3033,
			},
			wantErr: errSql,
		},
		{
			name: "Add",
			stock: AddStockDTO{
				SKUID:  1001,
				UserID: 0,
			},
			wantErr: nil,
		},
		{
			name: "Update",
			stock: AddStockDTO{
				SKUID:  2020,
				UserID: 1,
			},
			wantErr: nil,
		},
		{
			name: "ErrorUserId",
			stock: AddStockDTO{
				SKUID:  2020,
				UserID: 3,
			},
			wantErr: ErrUserID,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := usecase.AddStock(t.Context(), tt.stock)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("wanted: %v, respond: %v", tt.wantErr.Error(), err)
			}
		})
	}
}

func TestDeleteStockBySKU(t *testing.T) {
	repoMock := repositoryMock.NewIStockRepoMock(t)
	trxMock := mock.NewIPgTxManagerMock(t)
	kafkaMock := mock.NewIProducerMock(t)
	logger := logMock.NewLoggerMock(t)

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

	usecase := NewStockUsecase(repoMock, trxMock, kafkaMock, logger)

	tests := []struct {
		name    string
		body    DeleteStockDTO
		wantErr error
	}{
		{
			name: testSuccesName,
			body: DeleteStockDTO{
				UserID: 1,
				SKUID:  1001,
			},
			wantErr: nil,
		},
		{
			name: testSqlErrorName,
			body: DeleteStockDTO{
				UserID: 2,
				SKUID:  1001,
			},
			wantErr: ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := usecase.DeleteStockBySKU(t.Context(), tt.body)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("wanted: %v, respond: %v", tt.wantErr.Error(), err)
			}
		})
	}
}

func TestGetStockByLocation(t *testing.T) {
	repoMock := repositoryMock.NewIStockRepoMock(t)
	trxMock := mock.NewIPgTxManagerMock(t)
	kafkaMock := mock.NewIProducerMock(t)
	logger := logMock.NewLoggerMock(t)

	t.Cleanup(func() {
		repoMock.MinimockFinish()
		trxMock.MinimockFinish()
	})

	repoMock.GetItemsByLocationMock.Set(func(ctx context.Context, param repository.GetStockByLocation) ([]models.Item, error) {
		if param.UserID > 1 {
			return []models.Item{}, errSql
		}

		return []models.Item{
			{
				SKU: models.SKU{
					ID: 1001,
				},
				Stock: models.Stock{
					Location: "AG",
				},
			},
		}, nil
	})

	trxMock.WithTxMock.Set(func(ctx context.Context, fn func(repository.IStockRepo) error) (err error) {
		return fn(repoMock)
	})

	usecase := NewStockUsecase(repoMock, trxMock, kafkaMock, logger)

	tests := []struct {
		name    string
		body    GetItemByLocDTO
		want    ItemsByLocDTO
		wantErr error
	}{
		{
			name: testSuccesName,
			body: GetItemByLocDTO{
				UserID: 1,
			},
			want: ItemsByLocDTO{
				Stocks:     []StockDTO{},
				TotalCount: 1,
				PageNumber: 1,
			},
			wantErr: nil,
		},
		{
			name: testSqlErrorName,
			body: GetItemByLocDTO{
				UserID: 2,
			},
			want: ItemsByLocDTO{
				Stocks:     []StockDTO{{}},
				TotalCount: 0,
				PageNumber: 1,
			},
			wantErr: errSql,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			items, err := usecase.GetStocksByLocation(t.Context(), tt.body)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("wanted: %v, respond: %v", tt.wantErr.Error(), err)
			}

			if items.TotalCount != tt.want.TotalCount {
				t.Errorf("wanted total count:%d,  respond: %d", items.TotalCount, tt.want.TotalCount)
			}

			t.Log(err)
		})
	}
}

func TestGetItemBySKU(t *testing.T) {
	repoMock := repositoryMock.NewIStockRepoMock(t)
	trxMock := mock.NewIPgTxManagerMock(t)
	kafkaMock := mock.NewIProducerMock(t)
	logger := logMock.NewLoggerMock(t)

	t.Cleanup(func() {
		repoMock.MinimockFinish()
		trxMock.MinimockFinish()
	})

	repoMock.GetItemBySKUMock.Set(func(ctx context.Context, skuID models.SKUID) (models.Item, error) {
		if skuID != 1001 {
			return models.Item{SKU: models.SKU{}}, errors.New("not found")
		}

		return models.Item{SKU: models.SKU{ID: 1001}}, nil
	})

	trxMock.WithTxMock.Set(func(ctx context.Context, fn func(repository.IStockRepo) error) (err error) {
		return fn(repoMock)
	})

	usecase := NewStockUsecase(repoMock, trxMock, kafkaMock, logger)

	tests := []struct {
		name    string
		body    models.SKUID
		want    StockDTO
		wantErr error
	}{
		{
			name: testSuccesName,
			body: 1001,
			want: StockDTO{
				SKU: SKUDTO{
					SKUID: 1001,
				},
			},
			wantErr: nil,
		},
		{
			name:    "NotFound",
			body:    2020,
			want:    StockDTO{},
			wantErr: ErrNotFound,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item, err := usecase.GetItemBySKU(t.Context(), tt.body)
			if !errors.Is(err, tt.wantErr) {
				t.Errorf("wanted: %v, respond: %v", tt.wantErr.Error(), err)
			}

			if item != tt.want {
				t.Errorf("wanted: %v, respond: %v", tt.want, item)
			}

			t.Log(err)
		})
	}
}
