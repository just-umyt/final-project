package usecase

import (
	"context"
	"errors"
	"stocks/internal/models"
	"stocks/internal/repository"
)

type StockUsecaseInterface interface {
	AddStock(ctx context.Context, stock AddStockDTO) error
	DeleteStockBySKU(ctx context.Context, delStock DeleteStockDTO) error
	GetStocksByLocation(ctx context.Context, param GetItemByLocDTO) (ItemsByLocDTO, error)
	GetItemBySKU(ctx context.Context, sku models.SKUID) (StockDTO, error)
}

type StockUsecase struct {
	tx repository.PgTxManagerInterface
}

var (
	ErrNotFound error = errors.New("not found")
	ErrUserID   error = errors.New("user id is not matched")
)

func NewStockUsecase(pgTx repository.PgTxManager) *StockUsecase {
	return &StockUsecase{tx: &pgTx}
}

func (u *StockUsecase) AddStock(ctx context.Context, stock AddStockDTO) error {
	return u.tx.WithTx(ctx, func(repo repository.StockRepoInterface) error {
		item, err := repo.GetItemBySKU(ctx, stock.SKUID)
		if err != nil {
			if item.SKU.ID == 0 {
				return ErrNotFound
			} else {
				return err
			}
		}

		newItem := models.Stock{
			Count:    stock.Count,
			Price:    stock.Price,
			Location: stock.Location,
			UserID:   stock.UserID,
			SKUID:    stock.SKUID,
		}

		switch item.Stock.UserID {
		case 0:
			return repo.AddStock(ctx, newItem)
		case stock.UserID:
			return repo.UpdateStock(ctx, newItem)
		default:
			return ErrUserID
		}
	})
}

func (u *StockUsecase) DeleteStockBySKU(ctx context.Context, delStock DeleteStockDTO) error {
	return u.tx.WithTx(ctx, func(repo repository.StockRepoInterface) error {
		rows, err := repo.DeleteStock(ctx, delStock.SKUID, delStock.UserID)
		if err != nil {
			return err
		}

		if rows == 0 {
			return ErrNotFound
		}

		return nil
	})
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

	err := u.tx.WithTx(ctx, func(repo repository.StockRepoInterface) error {
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
	err := u.tx.WithTx(ctx, func(repo repository.StockRepoInterface) error {
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
