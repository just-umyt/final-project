package usecase

import (
	"cart/internal/models"
	"cart/internal/repository"
	"cart/internal/services"
	"context"
	"errors"
)

type Repository struct {
	repository.ICartRepo
	repository.IPgTxManager
}

type ICartUsecase interface {
	AddItem(ctx context.Context, addItem AddItemDTO) error
	DeleteItem(ctx context.Context, delItem DeleteItemDTO) error
	GetItemsByUserID(ctx context.Context, userID models.UserID) (ListItemsDTO, error)
	ClearCartByUserID(ctx context.Context, userID models.UserID) error
}

type CartUsecase struct {
	Repository
	skuService services.IStockService
}

var (
	ErrNotFound       error = errors.New("not found")
	ErrNotEnoughStock error = errors.New("not enough stock")
)

func NewCartUsecase(repo Repository, service services.IStockService) *CartUsecase {
	return &CartUsecase{Repository: repo, skuService: service}
}

func (u *CartUsecase) AddItem(ctx context.Context, addItem AddItemDTO) error {
	item, err := u.skuService.GetItemInfo(ctx, addItem.SKUID)
	if err != nil {
		return err
	}

	if item.Count < addItem.Count {
		return ErrNotEnoughStock
	}

	return u.IPgTxManager.WithTx(ctx, func(repo repository.ICartRepo) error {
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
			err := repo.UpdateItemByUserID(ctx, cart)
			if errors.Is(err, repository.ErrNotFound) {
				return ErrNotFound
			}

			return err
		}

		err = repo.AddItem(ctx, cart)
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}

		return err
	})
}

func (u *CartUsecase) DeleteItem(ctx context.Context, delItem DeleteItemDTO) error {
	err := u.ICartRepo.DeleteItem(ctx, delItem.UserID, delItem.SKUID)
	if errors.Is(err, repository.ErrNotFound) {
		return ErrNotFound
	}

	return err
}

func (u *CartUsecase) GetItemsByUserID(ctx context.Context, userID models.UserID) (ListItemsDTO, error) {
	var list ListItemsDTO

	cart, err := u.ICartRepo.GetCartByUserID(ctx, userID)
	if err != nil {
		return list, err
	}

	for id, c := range cart {
		sku, err := u.skuService.GetItemInfo(ctx, id)
		if err != nil {
			return ListItemsDTO{}, err
		}

		if c <= sku.Count {
			list.Items = append(list.Items, sku)
			list.TotalPrice += sku.Price
		}
	}

	return list, nil
}

func (u *CartUsecase) ClearCartByUserID(ctx context.Context, userID models.UserID) error {
	return u.IPgTxManager.WithTx(ctx, func(repo repository.ICartRepo) error {
		err := repo.ClearCartByUserID(ctx, userID)
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}

		return err
	})
}
