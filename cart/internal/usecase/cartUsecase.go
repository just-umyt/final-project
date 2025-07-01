package usecase

import (
	"cart/internal/models"
	"cart/internal/repository"
	"cart/internal/services"
	"context"
	"errors"
	"log"
)

type ICartRepo interface {
	GetCartInfoByUserID(ctx context.Context, userID models.UserID, skuID models.SKUID) (models.CartID, uint16, error)
	UpdateItemByUserID(ctx context.Context, cart models.Cart) error
	AddItem(ctx context.Context, cart models.Cart) error
	DeleteItem(ctx context.Context, userID models.UserID, skuID models.SKUID) error
	GetCartByUserID(ctx context.Context, userID models.UserID) ([]models.CartItem, error)
	ClearCartByUserID(ctx context.Context, userID models.UserID) error
}

type IPgTxManager interface {
	WithTx(ctx context.Context, fn func(ICartRepo) error) error
}

type IStockService interface {
	GetItemInfo(ctx context.Context, skuID models.SKUID) (services.ItemDTO, error)
}

type CartUsecase struct {
	skuService IStockService
	cartRepo   ICartRepo
	trManager  IPgTxManager
}

var (
	ErrNotFound       error = errors.New("not found")
	ErrNotEnoughStock error = errors.New("not enough stock")
)

func NewCartUsecase(cartRepo ICartRepo, trManager IPgTxManager, service IStockService) *CartUsecase {
	return &CartUsecase{cartRepo: cartRepo, trManager: trManager, skuService: service}
}

func (u *CartUsecase) AddItem(ctx context.Context, addItem AddItemDTO) error {
	item, err := u.skuService.GetItemInfo(ctx, addItem.SKUID)
	if err != nil {
		return err
	}

	if item.Count < addItem.Count {
		return ErrNotEnoughStock
	}

	return u.trManager.WithTx(ctx, func(repo ICartRepo) error {
		cartID, count, err := repo.GetCartInfoByUserID(ctx, addItem.UserID, addItem.SKUID)
		if err != nil {
			return err
		}

		cart := models.Cart{
			UserID: addItem.UserID,
			SKUID:  addItem.SKUID,
			Count:  count + addItem.Count,
		}

		if cartID > 0 {
			err := repo.UpdateItemByUserID(ctx, cart)
			if errors.Is(err, ErrNotFound) {
				return ErrNotFound
			}

			return err
		}

		err = repo.AddItem(ctx, cart)
		if errors.Is(err, ErrNotFound) {
			return ErrNotFound
		}

		return err
	})
}

func (u *CartUsecase) DeleteItem(ctx context.Context, delItem DeleteItemDTO) error {
	err := u.cartRepo.DeleteItem(ctx, delItem.UserID, delItem.SKUID)
	if errors.Is(err, ErrNotFound) {
		return ErrNotFound
	}

	return err
}

func (u *CartUsecase) GetItemsByUserID(ctx context.Context, userID models.UserID) (ListItemsDTO, error) {
	var list ListItemsDTO

	carts, err := u.cartRepo.GetCartByUserID(ctx, userID)
	if err != nil {
		return list, err
	}

	for _, cart := range carts {
		sku, err := u.skuService.GetItemInfo(ctx, cart.SKUID)
		if err != nil {
			return ListItemsDTO{}, err
		}

		realCount := cart.Count

		if cart.Count > sku.Count {
			log.Printf("Warning: user requested %d of SKU %d, but only %d in stock. Adjusting.",
				cart.Count, cart.SKUID, sku.Count)
			realCount = sku.Count
		}
		sku.Count = realCount
		list.Items = append(list.Items, sku)
		list.TotalPrice += uint32(realCount) * sku.Price
	}

	return list, nil
}

func (u *CartUsecase) ClearCartByUserID(ctx context.Context, userID models.UserID) error {
	return u.trManager.WithTx(ctx, func(repo ICartRepo) error {
		err := repo.ClearCartByUserID(ctx, userID)
		if errors.Is(err, repository.ErrNotFound) {
			return ErrNotFound
		}

		return err
	})
}
