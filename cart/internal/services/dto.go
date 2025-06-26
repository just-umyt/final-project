package services

import "cart/internal/models"

type ItemDTO struct {
	SKUID    models.SKUID
	Name     string
	Type     string
	Count    uint16
	Price    uint32
	Location string
	UserID   models.UserID
}
