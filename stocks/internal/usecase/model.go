package usecase

import "stocks/internal/models"

type GetSKU struct {
	Name   string
	Price  uint32
	Count  int
	Type   string
	UserId models.UserID
}
