package services

import "cart/internal/models"

type Item struct {
	SkuId    models.SKUID
	Name     string
	Type     string
	Count    uint16
	Price    uint32
	Location string
	UserId   models.UserID
}
