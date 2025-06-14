package usecase

import (
	"context"
	"stocks/internal/models"
	"stocks/internal/repository"
)

type StockUsecase struct {
	tx repository.PgTxManagerInterface
}

func NewStockUsecase(pgTx repository.PgTxManager) *StockUsecase {
	return &StockUsecase{tx: &pgTx}
}

func (u *StockUsecase) AddItemUsecase(ctx context.Context, item models.SKU) error {
	err := u.tx.WithTx(ctx, func(sri repository.StockRepoInterface) error {
		_, err := sri.GetBySku(ctx, item.ID)
		if err != nil {
			return err
		}

		err = sri.AddItem(ctx, item)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		return err
	}

	return nil
}

func (u *StockUsecase) GetBySku(ctx context.Context, id models.SKUID) (*GetSKU, error) {
	var model *GetSKU

	err := u.tx.WithTx(ctx, func(sri repository.StockRepoInterface) error {
		resp, err := sri.GetBySku(ctx, id)
		if err != nil {
			return err
		}

		model.Count = resp.Count
		model.Name = resp.Name
		model.Price = resp.Price

		return nil
	})
	if err != nil {
		return nil, err
	}

	return model, nil
}
