package repository

import "stocks/internal/models"

type GetSkuByLocationParameter struct {
	User_id  models.UserID
	Location string
	Limit    int64
	Offset   int64
}

type Stock struct {
	Id       *models.StockID `db:"id"`
	SkuId    *models.SKUID   `db:"sku_id"`
	Count    *uint16         `db:"count"`
	Price    *uint32         `db:"price"`
	Location *string         `db:"location"`
	UserId   *models.UserID  `db:"user_id"`
}

type SKU struct {
	SkuId models.SKUID `db:"sku_id"`
	Name  string       `db:"name"`
	Type  string       `db:"type"`
}
