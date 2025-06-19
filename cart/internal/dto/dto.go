package dto

import "cart/internal/models"

type CartAddItemDto struct {
	UserId models.UserID
	SkuId  models.SKUID
	Count  uint16
}

type SKU struct {
	SkuId    models.SKUID
	Name     string
	Type     string
	Count    uint16
	Price    uint32
	Location string
	UserId   models.UserID
}

type DeleteItemDto struct {
	UserId models.UserID
	SkuId  models.SKUID
}

type ListDto struct {
	SKUs       []SKU
	TotalPrice uint32
}
