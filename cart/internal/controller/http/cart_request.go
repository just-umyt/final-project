package http

type CartAddItemRequest struct {
	UserId int64  `json:"user_id"`
	SkuId  uint32 `json:"sku"`
	Count  uint16 `json:"count"`
}

type DeleteItemRequest struct {
	UserId int64  `json:"user_id"`
	SkuId  uint32 `json:"sku"`
}

type UserIdRequest struct {
	UserId int64 `json:"user_id"`
}
