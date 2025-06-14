package repository

import "stocks/internal/models"

type GetSKU struct {
	Name   string        `db:"name"`
	Price  float64       `db:"price"`
	Count  int           `db:"count"`
	Type   string        `db:"type"`
	UserId models.UserID `db:"user_id"`
}
