package http

type AddStockRequest struct {
	SKUID    uint32 `json:"sku"`
	UserID   int64  `json:"userId"`
	Count    uint16 `json:"count"`
	Price    uint32 `json:"price"`
	Location string `json:"location"`
}

type DeleteStockRequest struct {
	UserID int64  `json:"userId"`
	SKUID  uint32 `json:"sku"`
}

type GetItemsByLocRequest struct {
	UserID      int64  `json:"userId"`
	Location    string `json:"location"`
	PageSize    int64  `json:"pageSize"`
	CurrentPage int64  `json:"currentPage"`
}

type GetItemBySKURequest struct {
	SKU uint32 `json:"sku"`
}
