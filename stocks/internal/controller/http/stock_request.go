package http

type AddStockRequest struct {
	SKUID    uint32 `json:"sku"`
	UserID   int64  `json:"user_id"`
	Count    uint16 `json:"count"`
	Price    uint32 `json:"price"`
	Location string `json:"location"`
}

type DeleteStockRequest struct {
	UserID int64  `json:"user_id"`
	SKUID  uint32 `json:"sku"`
}

type GetItemsByLocRequest struct {
	UserID      int64  `json:"user_id"`
	Location    string `json:"location"`
	PageSize    int64  `json:"page_size"`
	CurrentPage int64  `json:"current_page"`
}

type GetItemBySKURequest struct {
	SKU uint32 `json:"sku"`
}
