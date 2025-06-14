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
}

type StockUsecase struct {
	tx repository.PgTxManagerInterface
}

func NewStockUsecase(pgTx repository.PgTxManager) *StockUsecase {
	return &StockUsecase{tx: &pgTx}
}

func (u *StockUsecase) AddSkuUsecase(ctx context.Context, item models.SKU) error {
	err := u.tx.WithTx(ctx, func(sri repository.StockRepoInterface) error {
		repoModel, err := sri.GetSkuBySkuIdRepo(ctx, item.SkuId)
		if err != nil {
			err = sri.AddSkuRepo(ctx, item)
			if err != nil {
				return err
			}
		} else {
			if repoModel.UserId != item.UserId {
				return errors.New("user_id mismatch")
			}

			//update
			err := sri.UpdateSkuBySkuIdRepo(ctx, item)
			if err != nil {
				return err
			}
		}

		return nil
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
