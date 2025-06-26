package usecase

import (
	"cart/internal/models"
	"cart/internal/services"
)

type AddItemDTO struct {
	UserID models.UserID
	SKUID  models.SKUID
	Count  uint16
}

type DeleteItemDTO struct {
	UserID models.UserID
	SKUID  models.SKUID
}

type ListItemsDTO struct {
	Items      []services.ItemDTO
	TotalPrice uint32
}
