package usecase

import (
	"context"
	"errors"
	"stocks/internal/dto"
	"stocks/internal/models"
	"stocks/internal/repository"
)

type StockUsecaseInterface interface {
	AddSkuUsecase(ctx context.Context, item models.SKU) error
	GetSkuBySkuIdUsecase(ctx context.Context, id models.SKUID) (*dto.GetSkuDto, error)
	DeleteSkuBySkuIdUsecase(ctx context.Context, deleteDtp dto.DeleteSkuDto) error
	GetSkuByLocationUsecase(ctx context.Context, locationDto dto.GetSkuByLocationParamsDto) ([]dto.GetSkuDto, error)
}

type StockUsecase struct {
	tx repository.PgTxManagerInterface
}

func NewStockUsecase(pgTx repository.PgTxManager) *StockUsecase {
	return &StockUsecase{tx: &pgTx}
}

func (u *StockUsecase) AddSkuUsecase(ctx context.Context, item models.SKU) error {
	err := u.tx.WithTx(ctx, func(sri repository.StockRepoInterface) error {
		userId, err := sri.GetUserIdBySkuIdRepo(ctx, item.SkuId)
		if err != nil {
			return sri.AddSkuRepo(ctx, item)
		}

		if *userId != item.UserId {
			return errors.New("user_id mismatch")
		}

		//update
		return sri.UpdateSkuBySkuIdRepo(ctx, item)
	})

	if err != nil {
		return err
	}

	return nil
}

func (u *StockUsecase) GetSkuBySkuIdUsecase(ctx context.Context, id models.SKUID) (*dto.GetSkuDto, error) {
	newDto := new(dto.GetSkuDto)

	err := u.tx.WithTx(ctx, func(sri repository.StockRepoInterface) error {
		repoModel, err := sri.GetSkuBySkuIdRepo(ctx, id)
		if err != nil {
			return err
		}

		newDto.Count = repoModel.Count
		newDto.Name = repoModel.Name
		newDto.Price = repoModel.Price
		newDto.Type = repoModel.Type

		return nil
	})
	if err != nil {
		return nil, err
	}

	return newDto, nil
}

func (u *StockUsecase) DeleteSkuBySkuIdUsecase(ctx context.Context, deleteDto dto.DeleteSkuDto) error {
	return u.tx.WithTx(ctx, func(sri repository.StockRepoInterface) error {
		return sri.DeleteSkuBySkuIdRepo(ctx, deleteDto.SkuId, deleteDto.UserId)
	})
}

func (u *StockUsecase) GetSkuByLocationUsecase(ctx context.Context, paginationByLoc dto.GetSkuByLocationParamsDto) ([]dto.GetSkuDto, error) {
	var items []dto.GetSkuDto
	var err error

	limit := paginationByLoc.PageSize
	offset := limit * (paginationByLoc.CurrentPage - 1)

	params := repository.GetSkuByLocationParameter{
		User_id:  paginationByLoc.User_id,
		Location: paginationByLoc.Location,
		Limit:    limit,
		Offset:   offset,
	}

	err = u.tx.WithTx(ctx, func(sri repository.StockRepoInterface) error {
		itemsFromRepo, err := sri.GetSkusByLocationRepo(ctx, params)

		for _, itemRepo := range itemsFromRepo {
			item := dto.GetSkuDto{
				Sku:      itemRepo.Sku,
				Name:     itemRepo.Name,
				Type:     itemRepo.Type,
				Price:    itemRepo.Price,
				Count:    itemRepo.Count,
				Location: itemRepo.Location,
			}

			items = append(items, item)
		}

		return err
	})

	return items, err
}
