package usecase

import (
	"cart/internal/models"
	"cart/internal/services"
)

type CartAddItemDto struct {
	UserId models.UserID
	SkuId  models.SKUID
	Count  uint16
}

type DeleteItemDto struct {
	UserId models.UserID
	SkuId  models.SKUID
}

type ListItemsDto struct {
	SKUs       []services.Item
	TotalPrice uint32
}
