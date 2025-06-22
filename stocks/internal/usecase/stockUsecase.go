package usecase

import (
	"context"
	"errors"
	"stocks/internal/models"
	"stocks/internal/repository"
)

type StockUsecaseInterface interface {
	AddStockUsecase(ctx context.Context, stockDto AddStockDto) error
	DeleteStockBySkuIdUsecase(ctx context.Context, deleteDto DeleteStockDto) error
	GetStocksByLocationUsecase(ctx context.Context, paginationByLoc GetSkuByLocationParamsDto) (StockByLocDto, error)
	GetSkuStocksBySkuIdUsecase(ctx context.Context, skuId models.SKUID) (StockDto, error)
}

type StockUsecase struct {
	tx repository.PgTxManagerInterface
}

var (
	ErrNotFound error = errors.New("not found")
	ErrUserId   error = errors.New("user id is not matched")
)

// const (
// 	NotFoundError = "not found"
// 	UserIdError   = "user id is not matched"
// )

func NewStockUsecase(pgTx repository.PgTxManager) *StockUsecase {
	return &StockUsecase{tx: &pgTx}
}

func (u *StockUsecase) AddStockUsecase(ctx context.Context, stockDto AddStockDto) error {
	return u.tx.WithTx(ctx, func(sri repository.StockRepoInterface) error {
		sku, stock, err := sri.GetSkuStockBySkuId(ctx, stockDto.SkuId)
		if err != nil {
			if sku.SkuID == 0 {
				return ErrNotFound
			} else {
				return err
			}
		}

		newItem := models.Stock{
			Count:    stockDto.Count,
			Price:    stockDto.Price,
			Location: stockDto.Location,
			UserId:   stockDto.UserId,
			SkuId:    stockDto.SkuId,
		}

		switch stock.UserId {
		case 0:
			return sri.AddStock(ctx, newItem)
		case stockDto.UserId:
			return sri.UpdateStock(ctx, newItem)
		default:
			return ErrUserId
		}
	})
}

func (u *StockUsecase) DeleteStockBySkuIdUsecase(ctx context.Context, deleteDto DeleteStockDto) error {
	return u.tx.WithTx(ctx, func(sri repository.StockRepoInterface) error {
		rows, err := sri.DeleteStock(ctx, deleteDto.SkuId, deleteDto.UserId)
		if err != nil {
			return err
		}

		if rows == 0 {
			return ErrNotFound
		}

		return nil
	})
}

func (u *StockUsecase) GetStocksByLocationUsecase(ctx context.Context, paginationByLoc GetSkuByLocationParamsDto) (StockByLocDto, error) {
	var items StockByLocDto

	limit := paginationByLoc.PageSize
	offset := limit * (paginationByLoc.CurrentPage - 1)

	params := repository.GetSkuByLocationParameter{
		User_id:  paginationByLoc.User_id,
		Location: paginationByLoc.Location,
		Limit:    limit,
		Offset:   offset,
	}

	err := u.tx.WithTx(ctx, func(sri repository.StockRepoInterface) error {
		stocksFromRepo, err := sri.GetStocksByLocation(ctx, params)
		if err != nil {
			return err
		}

		for _, repoStock := range stocksFromRepo {
			item := StockDto{
				SkuDto: SkuDto{
					SkuId: repoStock.SkuID,
					Name:  repoStock.Name,
					Type:  repoStock.Type,
				},
				Price:    repoStock.Price,
				Count:    repoStock.Count,
				Location: repoStock.Location,
				UserId:   repoStock.UserId,
			}

			items.Stocks = append(items.Stocks, item)
		}

		return nil
	})

	items.TotalCount = len(items.Stocks)
	items.PageNumber = paginationByLoc.CurrentPage

	return items, err
}

func (u *StockUsecase) GetSkuStocksBySkuIdUsecase(ctx context.Context, skuId models.SKUID) (StockDto, error) {
	var stockDto StockDto
	err := u.tx.WithTx(ctx, func(sri repository.StockRepoInterface) error {
		sku, stock, err := sri.GetSkuStockBySkuId(ctx, skuId)
		if err != nil {
			if sku.SkuID == 0 {
				return ErrNotFound
			} else {
				return err
			}
		}

		stockDto = StockDto{
			SkuDto: SkuDto{
				SkuId: sku.SkuID,
				Name:  sku.Name,
				Type:  sku.Type,
			},
			Price:    stock.Price,
			Count:    stock.Count,
			Location: stock.Location,
			UserId:   stock.UserId,
		}

		return nil
	})

	return stockDto, err
}
