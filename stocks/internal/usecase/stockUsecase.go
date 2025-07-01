package usecase

import (
	"context"
	"errors"
	"stocks/internal/models"
	"stocks/internal/repository"
)

type IStockRepo interface {
	GetItemBySKU(ctx context.Context, skuID models.SKUID) (models.Item, error)
	AddStock(ctx context.Context, stock models.Stock) error
	UpdateStock(ctx context.Context, stock models.Stock) error
	DeleteStock(ctx context.Context, skuID models.SKUID, userID models.UserID) error
	GetItemsByLocation(ctx context.Context, param repository.GetStockByLocation) ([]models.Item, error)
}

type IPgTxManager interface {
	WithTx(ctx context.Context, fn func(IStockRepo) error) error
}

type StockUsecase struct {
	stockRepo IStockRepo
	trManager IPgTxManager
}

var (
	ErrNotFound error = errors.New("not found")
	ErrUserID   error = errors.New("user id is not matched")
)

func NewStockUsecase(repo IStockRepo, trManager IPgTxManager) *StockUsecase {
	return &StockUsecase{stockRepo: repo, trManager: trManager}
}

func (u *StockUsecase) AddStock(ctx context.Context, stock AddStockDTO) error {
	return u.trManager.WithTx(ctx, func(repo IStockRepo) error {
		item, err := repo.GetItemBySKU(ctx, stock.SKUID)
		if err != nil {
			if item.Stock.ID == 0 {
				return ErrNotFound
			}

			return err
		}

		newItem := models.Stock{
			Count:    item.Stock.Count + stock.Count,
			Price:    stock.Price,
			Location: stock.Location,
			UserID:   stock.UserID,
			SKUID:    stock.SKUID,
		}

		switch item.Stock.UserID {
		case 0:
			err := repo.AddStock(ctx, newItem)
			if errors.Is(err, repository.ErrNotFound) {
				return ErrNotFound
			}

			return err
		case stock.UserID:
			err := repo.UpdateStock(ctx, newItem)
			if errors.Is(err, repository.ErrNotFound) {
				return ErrNotFound
			}

			return err
		default:
			return ErrUserID
		}
	})
}

func (u *StockUsecase) DeleteStockBySKU(ctx context.Context, delStock DeleteStockDTO) error {
	err := u.stockRepo.DeleteStock(ctx, delStock.SKUID, delStock.UserID)
	if errors.Is(err, repository.ErrNotFound) {
		return ErrNotFound
	}

	return err
}

func (u *StockUsecase) GetStocksByLocation(ctx context.Context, param GetItemByLocDTO) (ItemsByLocDTO, error) {
	var items ItemsByLocDTO

	limit := param.PageSize
	offset := limit * (param.CurrentPage - 1)

	params := repository.GetStockByLocation{
		UserID:   param.UserID,
		Location: param.Location,
		Limit:    limit,
		Offset:   offset,
	}

	err := u.trManager.WithTx(ctx, func(repo IStockRepo) error {
		stocksFromRepo, err := repo.GetItemsByLocation(ctx, params)
		if err != nil {
			return err
		}

		for _, s := range stocksFromRepo {
			item := StockDTO{
				SKU: SKUDTO{
					SKUID: s.SKU.ID,
					Name:  s.SKU.Name,
					Type:  s.SKU.Type,
				},
				Price:    s.Stock.Price,
				Count:    s.Stock.Count,
				Location: s.Stock.Location,
				UserID:   s.Stock.UserID,
			}

			items.Stocks = append(items.Stocks, item)
		}

		return nil
	})

	items.TotalCount = len(items.Stocks)
	items.PageNumber = param.CurrentPage

	return items, err
}

func (u *StockUsecase) GetItemBySKU(ctx context.Context, sku models.SKUID) (StockDTO, error) {
	var stockDTO StockDTO
	err := u.trManager.WithTx(ctx, func(repo IStockRepo) error {
		item, err := repo.GetItemBySKU(ctx, sku)
		if err != nil {
			if item.SKU.ID == 0 {
				return ErrNotFound
			} else {
				return err
			}
		}

		stockDTO = StockDTO{
			SKU: SKUDTO{
				SKUID: item.SKU.ID,
				Name:  item.SKU.Name,
				Type:  item.SKU.Type,
			},
			Price:    item.Stock.Price,
			Count:    item.Stock.Count,
			Location: item.Stock.Location,
			UserID:   item.Stock.UserID,
		}

		return nil
	})

	return stockDTO, err
}
