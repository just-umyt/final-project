package usecase

import (
	"cart/internal/dto"
	"cart/internal/models"
	"cart/internal/repository"
	"cart/internal/services"
	"context"
	"errors"
)

type CartUsecaseInterface interface {
	CartAddItemUsecase(ctx context.Context, cartDto dto.CartAddItemDto) error
	CartDeleteItemUsecase(ctx context.Context, item dto.DeleteItemDto) error
	CartListByUserIdUsecase(ctx context.Context, userId models.UserID) (dto.ListDto, error)
	CartClearByUserIdUsecase(ctx context.Context, userId models.UserID) error
}

type CartUsecase struct {
	tx         repository.PgTxManagerInterface
	skuService services.SkuGetServiceInterface
}

func NewCartUsecase(pgTx repository.PgTxManager, service services.SkuGetServiceInterface) *CartUsecase {
	return &CartUsecase{tx: &pgTx, skuService: service}
}

func (u *CartUsecase) CartAddItemUsecase(ctx context.Context, cartDto dto.CartAddItemDto) error {
	//get_sku
	sku, err := u.skuService.GetSku(ctx, cartDto.SkuId)
	if err != nil {
		return err
	}

	//validate counts
	if sku.Count < 1 || sku.Count < cartDto.Count {
		return errors.New("not enough stock")
	}

	return u.tx.WithTx(ctx, func(cri repository.CartRepoInterface) error {
		cartId, err := cri.GetItemIdByUserIdRepo(ctx, cartDto.UserId, cartDto.SkuId)
		if err != nil {
			return err
		}

		cart := models.Cart{
			UserId: cartDto.UserId,
			SKUId:  cartDto.SkuId,
			Count:  cartDto.Count,
		}

		if cartId > 0 {
			return cri.UpdateItemByUserIdRepo(ctx, cart)
		}

		return cri.AddItemRepo(ctx, cart)
	})
}

func (u *CartUsecase) CartDeleteItemUsecase(ctx context.Context, item dto.DeleteItemDto) error {
	return u.tx.WithTx(ctx, func(cri repository.CartRepoInterface) error {
		rowsAffect, err := cri.DeleteItemRepo(ctx, item.UserId, item.SkuId)
		if err != nil {
			return err
		}

		if rowsAffect < 1 {
			return errors.New("not found")
		}

		return nil
	})
}

func (u *CartUsecase) CartListByUserIdUsecase(ctx context.Context, userId models.UserID) (dto.ListDto, error) {
	var skuIds []models.SKUID

	err := u.tx.WithTx(ctx, func(cri repository.CartRepoInterface) error {
		var err error
		skuIds, err = cri.GetItemsSkuIdByCartId(ctx, userId)
		if err != nil {
			return err
		}

		return nil
	})

	var list dto.ListDto
	var skus []dto.SKU

	for _, e := range skuIds {
		sku, err := u.skuService.GetSku(ctx, e)
		if err != nil {
			return dto.ListDto{}, err
		}

		list.TotalPrice += sku.Price
		skus = append(skus, sku)
	}

	list.SKUs = skus

	return list, err
}

func (u *CartUsecase) CartClearByUserIdUsecase(ctx context.Context, userId models.UserID) error {
	return u.tx.WithTx(ctx, func(cri repository.CartRepoInterface) error {
		rowsAffect, err := cri.ClearCartByUserIdRepo(ctx, userId)
		if err != nil {
			return nil
		}

		if rowsAffect < 1 {
			return errors.New("not found")
		}

		return nil
	})
}
