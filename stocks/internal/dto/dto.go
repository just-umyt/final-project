package dto

import "stocks/internal/models"

type AddSkuDto struct {
	UserId   models.UserID `json:"user_id"`
	SkuId    models.SKUID  `json:"sku"`
	Count    uint16        `json:"count"`
	Name     string        `json:"name"`
	Type     string        `json:"type"`
	Price    uint32        `json:"price"`
	Location string        `json:"location"`
}

type GetSkuDto struct {
	Sku      models.SKUID `json:"sku_id"`
	Name     string       `json:"name"`
	Type     string       `json:"type"`
	Price    uint32       `json:"price"`
	Count    int          `json:"count"`
	Location string       `json:"location"`
}

type GetSkuBySkuIdDto struct {
	SkuId models.SKUID `json:"sku"`
}

type DeleteSkuDto struct {
	UserId models.UserID `json:"user_id"`
	SkuId  models.SKUID  `json:"sku"`
}

type GetSkuByLocationParamsDto struct {
	User_id     models.UserID `json:"user_id"`
	Location    string        `json:"location"`
	PageSize    int64         `json:"page_size"`
	CurrentPage int64         `json:"current_page"`
}

type PaginationResonse struct {
	Items []GetSkuDto `json:"items"`
	Err   error       `json:"error"`
}
