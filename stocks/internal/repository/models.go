package repository

import "stocks/internal/models"

type GetSKU struct {
	Name   string        `db:"name"`
	Price  uint32        `db:"price"`
	Count  int           `db:"count"`
	Type   string        `db:"type"`
	UserId models.UserID `db:"user_id"`
}

type GetSkuByLocation struct {
	Sku      models.SKUID `db:"sku_id"`
	Name     string       `db:"name"`
	Type     string       `db:"type"`
	Price    uint32       `db:"price"`
	Count    int          `db:"count"`
	Location string       `db:"location"`
}

type GetSkuByLocationParameter struct {
	User_id  models.UserID
	Location string
	Limit    int64
	Offset   int64
}
