package usecase

import "stocks/internal/models"

type AddStockDto struct {
	SkuId    models.SKUID
	UserId   models.UserID
	Count    uint16
	Price    uint32
	Location string
}

type DeleteStockDto struct {
	UserId models.UserID
	SkuId  models.SKUID
}

type GetSkuByLocationParamsDto struct {
	User_id     models.UserID
	Location    string
	PageSize    int64
	CurrentPage int64
}

type SkuDto struct {
	SkuId models.SKUID
	Name  string
	Type  string
}

type StockDto struct {
	SkuDto
	Count    uint16
	Price    uint32
	Location string
	UserId   models.UserID
}

type StockByLocDto struct {
	Stocks     []StockDto
	TotalCount int
	PageNumber int64
}
