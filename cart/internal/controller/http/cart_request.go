package http

type AddItemRequest struct {
	UserID int64  `json:"userId"`
	SKUID  uint32 `json:"sku"`
	Count  uint16 `json:"count"`
}

type DeleteItemRequest struct {
	UserID int64  `json:"userId"`
	SKUID  uint32 `json:"sku"`
}

type UserIDRequest struct {
	UserID int64 `json:"userId"`
}
