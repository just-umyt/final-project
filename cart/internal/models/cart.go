package models

type Cart struct {
	ID     CartID
	UserID UserID
	SKUID  SKUID
	Count  uint16
}

type CartItem struct {
	SKUID SKUID
	Count uint16
}
