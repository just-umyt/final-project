package usecase

import "stocks/internal/models"

type AddStockDTO struct {
	SKUID    models.SKUID
	UserID   models.UserID
	Count    uint16
	Price    uint32
	Location string
}

type DeleteStockDTO struct {
	UserID models.UserID
	SKUID  models.SKUID
}

type GetItemByLocDTO struct {
	UserID      models.UserID
	Location    string
	PageSize    int64
	CurrentPage int64
}

type SKUDTO struct {
	SKUID models.SKUID
	Name  string
	Type  string
}

type StockDTO struct {
	SKU      SKUDTO
	Count    uint16
	Price    uint32
	Location string
	UserID   models.UserID
}

type ItemsByLocDTO struct {
	Stocks     []StockDTO
	TotalCount int
	PageNumber int64
}
