package usecase

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"stocks/internal/dto"
	"stocks/internal/models"
	"stocks/internal/repository"
)

type StockUsecaseInterface interface {
	AddStockUsecase(ctx context.Context, stockDto dto.AddStockDto) dto.ErrorResponse
	DeleteStockBySkuIdUsecase(ctx context.Context, deleteDto dto.DeleteStockDto) dto.ErrorResponse
	GetStocksByLocationUsecase(ctx context.Context, paginationByLoc dto.GetSkuByLocationParamsDto) (dto.StockByLocDto, dto.ErrorResponse)
	GetSkuStocksBySkuIdUsecase(ctx context.Context, skuId models.SKUID) (dto.StockDto, dto.ErrorResponse)
}

type StockUsecase struct {
	tx repository.PgTxManagerInterface
}

func NewStockUsecase(pgTx repository.PgTxManager) *StockUsecase {
	return &StockUsecase{tx: &pgTx}
}

func (u *StockUsecase) AddStockUsecase(ctx context.Context, stockDto dto.AddStockDto) dto.ErrorResponse {
	var errorRes dto.ErrorResponse
	errorRes.Message = u.tx.WithTx(ctx, func(sri repository.StockRepoInterface) error {
		sku, stock, err := sri.GetSkuStockBySkuIdRepo(ctx, stockDto.SkuId)
		if err != nil {
			if sku.SkuID == 0 {
				errorRes.Code = http.StatusNotFound
				return fmt.Errorf("not found: error is: %v", err.Error())
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
			return sri.AddStockRepo(ctx, newItem)
		case stockDto.UserId:
			return sri.UpdateStockRepo(ctx, newItem)
		default:
			return errors.New("user id is not matched")
		}
	})

	if errorRes.Message != nil && errorRes.Code == 0 {
		errorRes.Code = http.StatusInternalServerError
	}

	return errorRes
}

func (u *StockUsecase) DeleteStockBySkuIdUsecase(ctx context.Context, deleteDto dto.DeleteStockDto) dto.ErrorResponse {
	var errorRes dto.ErrorResponse
	errorRes.Message = u.tx.WithTx(ctx, func(sri repository.StockRepoInterface) error {
		rows, err := sri.DeleteStockRepo(ctx, deleteDto.SkuId, deleteDto.UserId)
		if err != nil {
			return err
		}

		if rows == 0 {
			errorRes.Code = http.StatusNotFound
			return errors.New("not found")
		}

		return nil
	})

	if errorRes.Message != nil && errorRes.Code == 0 {
		errorRes.Code = http.StatusInternalServerError
	}

	return errorRes
}

func (u *StockUsecase) GetStocksByLocationUsecase(ctx context.Context, paginationByLoc dto.GetSkuByLocationParamsDto) (dto.StockByLocDto, dto.ErrorResponse) {
	var items dto.StockByLocDto
	var errorRes dto.ErrorResponse

	limit := paginationByLoc.PageSize
	offset := limit * (paginationByLoc.CurrentPage - 1)

	params := repository.GetSkuByLocationParameter{
		User_id:  paginationByLoc.User_id,
		Location: paginationByLoc.Location,
		Limit:    limit,
		Offset:   offset,
	}

	errorRes.Message = u.tx.WithTx(ctx, func(sri repository.StockRepoInterface) error {
		stocksFromRepo, err := sri.GetStocksByLocationRepo(ctx, params)
		if err != nil {
			return err
		}

		for _, repoStock := range stocksFromRepo {
			item := dto.StockDto{
				SkuDto: dto.SkuDto{
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

	if errorRes.Message != nil && errorRes.Code == 0 {
		errorRes.Code = http.StatusInternalServerError
	}

	items.TotalCount = len(items.Stocks)
	items.PageNumber = paginationByLoc.CurrentPage

	return items, errorRes
}

func (u *StockUsecase) GetSkuStocksBySkuIdUsecase(ctx context.Context, skuId models.SKUID) (dto.StockDto, dto.ErrorResponse) {
	var stockDto dto.StockDto
	var errorRes dto.ErrorResponse

	errorRes.Message = u.tx.WithTx(ctx, func(sri repository.StockRepoInterface) error {
		sku, stock, err := sri.GetSkuStockBySkuIdRepo(ctx, skuId)
		if err != nil {
			if sku.SkuID == 0 {
				errorRes.Code = http.StatusNotFound
				return fmt.Errorf("not found: error is: %v", err.Error())
			} else {
				return err
			}
		}

		stockDto = dto.StockDto{
			SkuDto: dto.SkuDto{
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

	if errorRes.Message != nil && errorRes.Code == 0 {
		errorRes.Code = http.StatusInternalServerError
	}

	return stockDto, errorRes
}
