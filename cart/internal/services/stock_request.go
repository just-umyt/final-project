package services

import "cart/internal/models"

type getSKUIDRequest struct {
	SKUID models.SKUID `json:"sku"`
}

type httpResponse struct {
	Message stock `json:"message"`
}

type stock struct {
	SKUID    uint32 `json:"sku"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Count    uint16 `json:"count,omitempty"`
	Price    uint32 `json:"price,omitempty"`
	Location string `json:"location,omitempty"`
	UserID   int64  `json:"user_id,omitempty"`
}
