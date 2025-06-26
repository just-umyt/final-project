package http

type StockByLocResponse struct {
	Items      []ItemResponse `json:"stocks"`
	TotalCount int            `json:"total_count"`
	PageNumber int64          `json:"page_number"`
}

type ItemResponse struct {
	SKU      uint32 `json:"sku"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Count    uint16 `json:"count,omitempty"`
	Price    uint32 `json:"price,omitempty"`
	Location string `json:"location,omitempty"`
	UserID   int64  `json:"user_id,omitempty"`
}
