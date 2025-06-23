package services

import "cart/internal/models"

type GetSkuRequest struct {
	SkuId models.SKUID `json:"sku"`
}

type Response struct {
	Message StockResponse `json:"message"`
}

type StockResponse struct {
	SkuId    uint32 `json:"sku"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Count    uint16 `json:"count,omitempty"`
	Price    uint32 `json:"price,omitempty"`
	Location string `json:"location,omitempty"`
	UserId   int64  `json:"user_id,omitempty"`
}
