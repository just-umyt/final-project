package controller

type AddItemRequest struct {
	UserID int64  `json:"userId" validate:"required"`
	SKUID  uint32 `json:"sku" validate:"required"`
	Count  uint16 `json:"count" validate:"required"`
}

type DeleteItemRequest struct {
	UserID int64  `json:"userId" validate:"required"`
	SKUID  uint32 `json:"sku" validate:"required"`
}

type UserIDRequest struct {
	UserID int64 `json:"userId"`
}
