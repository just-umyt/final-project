package integration

type AddStockRequest struct {
	SKUID    uint32 `json:"sku" validate:"required"`
	UserID   int64  `json:"userId" validate:"required"`
	Count    uint16 `json:"count" validate:"required"`
	Price    uint32 `json:"price" validate:"required"`
	Location string `json:"location" validate:"required"`
}

type DeleteStockRequest struct {
	UserID int64  `json:"userId" validate:"required"`
	SKUID  uint32 `json:"sku" validate:"required"`
}

type GetItemsByLocRequest struct {
	UserID      int64  `json:"userId" validate:"required"`
	Location    string `json:"location" validate:"required"`
	PageSize    int64  `json:"pageSize" validate:"required"`
	CurrentPage int64  `json:"currentPage" validate:"required"`
}

type GetItemBySKURequest struct {
	SKU uint32 `json:"sku"`
}
