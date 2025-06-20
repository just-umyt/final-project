package http

type AddStockRequest struct {
	SkuId    uint32 `json:"sku"`
	UserId   int64  `json:"user_id"`
	Count    uint16 `json:"count"`
	Price    uint32 `json:"price"`
	Location string `json:"location"`
}

type DeleteStockRequest struct {
	UserId int64  `json:"user_id"`
	SkuId  uint32 `json:"sku"`
}

type GetSkuByLocationParamsRequest struct {
	User_id     int64  `json:"user_id"`
	Location    string `json:"location"`
	PageSize    int64  `json:"page_size"`
	CurrentPage int64  `json:"current_page"`
}

type GetSkuBySkuIdRequest struct {
	SkuId uint32 `json:"sku"`
}
