package usecase

import "stocks/internal/models"

type GetSKU struct {
	Name   string
	Price  float64
	Count  int
	Type   string
	UserId models.UserID
}
