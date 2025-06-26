package usecase

import (
	"cart/internal/models"
	"cart/internal/repository"
	"cart/internal/services"
	"context"
	"errors"
)

type CartUsecaseInterface interface {
	AddItem(ctx context.Context, addItem AddItemDTO) error
	DeleteItem(ctx context.Context, delItem DeleteItemDTO) error
	GetItemsByUserID(ctx context.Context, userID models.UserID) (ListItemsDTO, error)
	ClearCartByUserID(ctx context.Context, userID models.UserID) error
}

type CartUsecase struct {
	tx         repository.PgTxManagerInterface
	skuService services.StockService
}

var (
	ErrNotFound       error = errors.New("not found")
	ErrNotEnoughStock error = errors.New("not enough stock")
)

func NewCartUsecase(pgTx repository.PgTxManager, service services.StockService) *CartUsecase {
	return &CartUsecase{tx: &pgTx, skuService: service}
}

func (u *CartUsecase) AddItem(ctx context.Context, addItem AddItemDTO) error {
	item, err := u.skuService.GetItemInfo(ctx, addItem.SKUID)
	if err != nil {
		return err
	}

	if item.Count < addItem.Count {
		return ErrNotEnoughStock
	}

	return u.tx.WithTx(ctx, func(repo repository.CartRepoInterface) error {
		cartID, err := repo.GetCartIDByUserID(ctx, addItem.UserID, addItem.SKUID)
		if err != nil {
			return err
		}

		cart := models.Cart{
			UserID: addItem.UserID,
			SKUID:  addItem.SKUID,
			Count:  addItem.Count,
		}

		if cartID > 0 {
			return repo.UpdateItemByUserID(ctx, cart)
		}

		return repo.AddItem(ctx, cart)
	})
}

func (u *CartUsecase) DeleteItem(ctx context.Context, delItem DeleteItemDTO) error {
	return u.tx.WithTx(ctx, func(repo repository.CartRepoInterface) error {
		rowsAffect, err := repo.DeleteItem(ctx, delItem.UserID, delItem.SKUID)
		if err != nil {
			return err
		}

		if rowsAffect < 1 {
			return ErrNotFound
		}

		return nil
	})
}

func (u *CartUsecase) GetItemsByUserID(ctx context.Context, userID models.UserID) (ListItemsDTO, error) {
	var skuIDs []models.SKUID

	err := u.tx.WithTx(ctx, func(repo repository.CartRepoInterface) error {
		var err error

		skuIDs, err = repo.GetSKUIDsByUserID(ctx, userID)
		if err != nil {
			return err
		}

		return nil
	})

	var list ListItemsDTO

	for _, id := range skuIDs {
		sku, err := u.skuService.GetItemInfo(ctx, id)
		if err != nil {
			return ListItemsDTO{}, err
		}

		list.Items = append(list.Items, sku)
		list.TotalPrice += sku.Price
	}

	return list, err
}

func (u *CartUsecase) ClearCartByUserID(ctx context.Context, userID models.UserID) error {
	return u.tx.WithTx(ctx, func(repo repository.CartRepoInterface) error {
		rowsAffect, err := repo.ClearCartByUserID(ctx, userID)
		if err != nil {
			return err
		}

		if rowsAffect < 1 {
			return ErrNotFound
		}

		return nil
	})
}
