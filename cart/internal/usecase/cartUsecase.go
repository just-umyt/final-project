package usecase

import (
	"cart/internal/models"
	"cart/internal/repository"
	"cart/internal/services"
	"context"
	"errors"
)

type CartUsecaseInterface interface {
	CartAddItem(ctx context.Context, cartDto CartAddItemDto) error
	CartDeleteItem(ctx context.Context, item DeleteItemDto) error
	CartListByUserId(ctx context.Context, userId models.UserID) (ListItemsDto, error)
	CartClearByUserId(ctx context.Context, userId models.UserID) error
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

func (u *CartUsecase) CartAddItem(ctx context.Context, cartDto CartAddItemDto) error {
	sku, err := u.skuService.GetItemInfo(ctx, cartDto.SkuId)
	if err != nil {
		return err
	}

	if sku.Count < cartDto.Count {
		return ErrNotEnoughStock
	}

	return u.tx.WithTx(ctx, func(cri repository.CartRepoInterface) error {
		cartId, err := cri.GetItemIdByUserId(ctx, cartDto.UserId, cartDto.SkuId)
		if err != nil {
			return err
		}

		cart := models.Cart{
			UserId: cartDto.UserId,
			SKUId:  cartDto.SkuId,
			Count:  cartDto.Count,
		}

		if cartId > 0 {
			return cri.UpdateItemByUserId(ctx, cart)
		}

		return cri.AddItem(ctx, cart)
	})
}

func (u *CartUsecase) CartDeleteItem(ctx context.Context, item DeleteItemDto) error {
	return u.tx.WithTx(ctx, func(cri repository.CartRepoInterface) error {
		rowsAffect, err := cri.DeleteItem(ctx, item.UserId, item.SkuId)
		if err != nil {
			return err
		}

		if rowsAffect < 1 {
			return ErrNotFound
		}

		return nil
	})
}

func (u *CartUsecase) CartListByUserId(ctx context.Context, userId models.UserID) (ListItemsDto, error) {
	var skuIds []models.SKUID

	err := u.tx.WithTx(ctx, func(cri repository.CartRepoInterface) error {
		var err error

		skuIds, err = cri.GetItemsSkuIdByCartId(ctx, userId)
		if err != nil {
			return err
		}

		return nil
	})

	var list ListItemsDto

	for _, e := range skuIds {
		sku, err := u.skuService.GetItemInfo(ctx, e)
		if err != nil {
			return ListItemsDto{}, err
		}

		list.TotalPrice += sku.Price
		list.SKUs = append(list.SKUs, sku)
	}

	return list, err
}

func (u *CartUsecase) CartClearByUserId(ctx context.Context, userId models.UserID) error {
	return u.tx.WithTx(ctx, func(cri repository.CartRepoInterface) error {
		rowsAffect, err := cri.ClearCartByUserId(ctx, userId)
		if err != nil {
			return err
		}

		if rowsAffect < 1 {
			return ErrNotFound
		}

		return nil
	})
}
