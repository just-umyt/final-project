package http

type StockByLocResponse struct {
	Stocks     []StockResponse `json:"stocks"`
	TotalCount int             `json:"total_count"`
	PageNumber int64           `json:"page_number"`
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
