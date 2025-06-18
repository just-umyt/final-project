package dto

import "stocks/internal/models"

type SkuDto struct {
	SkuId models.SKUID `json:"sku"`
	Name  string       `json:"name"`
	Type  string       `json:"type"`
}

type StockDto struct {
	SkuDto
	Count    uint16        `json:"count,omitempty"`
	Price    uint32        `json:"price,omitempty"`
	Location string        `json:"location,omitempty"`
	UserId   models.UserID `json:"user_id,omitempty"`
}

type StockByLocDto struct {
	Stocks     []StockDto `json:"stocks"`
	TotalCount int        `json:"total_count"`
	PageNumber int64      `json:"page_number"`
}

type AddStockDto struct {
	SkuId    models.SKUID  `json:"sku"`
	UserId   models.UserID `json:"user_id"`
	Count    uint16        `json:"count"`
	Price    uint32        `json:"price"`
	Location string        `json:"location"`
}

type DeleteStockDto struct {
	UserId models.UserID `json:"user_id"`
	SkuId  models.SKUID  `json:"sku"`
}

type ErrorResponse struct {
	Message error `json:"message"`
	Code    int   `json:"code"`
}

type GetSkuByLocationParamsDto struct {
	User_id     models.UserID `json:"user_id"`
	Location    string        `json:"location"`
	PageSize    int64         `json:"page_size"`
	CurrentPage int64         `json:"current_page"`
}

type GetSkuBySkuIdDto struct {
	SkuId models.SKUID `json:"sku"`
}
