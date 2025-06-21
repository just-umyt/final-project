package usecase

import (
	"cart/internal/models"
	"cart/internal/repository"
	"cart/internal/services"
	"context"
	"errors"
)

type CartUsecaseInterface interface {
	CartAddItemUsecase(ctx context.Context, cartDto CartAddItemDto) error
	CartDeleteItemUsecase(ctx context.Context, item DeleteItemDto) error
	CartListByUserIdUsecase(ctx context.Context, userId models.UserID) (ListDto, error)
	CartClearByUserIdUsecase(ctx context.Context, userId models.UserID) error
}

type CartUsecase struct {
	tx         repository.PgTxManagerInterface
	skuService services.SkuGetService
}

const (
	NotFoundError  = "not found"
	NotEnoughStock = "not enough stock"
)

func NewCartUsecase(pgTx repository.PgTxManager, service services.SkuGetService) *CartUsecase {
	return &CartUsecase{tx: &pgTx, skuService: service}
}

func (u *CartUsecase) CartAddItemUsecase(ctx context.Context, cartDto CartAddItemDto) error {
	sku, err := u.skuService.GetItemInfo(ctx, cartDto.SkuId)
	if err != nil {
		return err
	}

	if sku.Count < 1 || sku.Count < cartDto.Count {
		return errors.New(NotEnoughStock)
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

func (u *CartUsecase) CartDeleteItemUsecase(ctx context.Context, item DeleteItemDto) error {
	return u.tx.WithTx(ctx, func(cri repository.CartRepoInterface) error {
		rowsAffect, err := cri.DeleteItem(ctx, item.UserId, item.SkuId)
		if err != nil {
			return err
		}

		if rowsAffect < 1 {
			return errors.New(NotFoundError)
		}

		return nil
	})
}

func (u *CartUsecase) CartListByUserIdUsecase(ctx context.Context, userId models.UserID) (ListDto, error) {
	var skuIds []models.SKUID

	err := u.tx.WithTx(ctx, func(cri repository.CartRepoInterface) error {
		var err error

		skuIds, err = cri.GetItemsSkuIdByCartId(ctx, userId)
		if err != nil {
			return err
		}

		return nil
	})

	var list ListDto

	for _, e := range skuIds {
		sku, err := u.skuService.GetItemInfo(ctx, e)
		if err != nil {
			return ListDto{}, err
		}

		list.TotalPrice += sku.Price
		list.SKUs = append(list.SKUs, sku)
	}

	return list, err
}

func (u *CartUsecase) CartClearByUserIdUsecase(ctx context.Context, userId models.UserID) error {
	return u.tx.WithTx(ctx, func(cri repository.CartRepoInterface) error {
		rowsAffect, err := cri.ClearCartByUserId(ctx, userId)
		if err != nil {
			return err
		}

		if rowsAffect < 1 {
			return errors.New(NotFoundError)
		}

		return nil
	})
}
