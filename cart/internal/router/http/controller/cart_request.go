package controller

type AddItemRequest struct {
	UserID int64  `json:"userId" validate:"required,min=1, "`
	SKUID  uint32 `json:"sku" validate:"required, min=4"`
	Count  uint16 `json:"count" validate:"required, min=1"`
}

type DeleteItemRequest struct {
	UserID int64  `json:"userId"`
	SKUID  uint32 `json:"sku"`
}

type UserIDRequest struct {
	UserID int64 `json:"userId"`
}
