package models

type Cart struct {
	ID     CartID
	UserID UserID
	SKUID  SKUID
	Count  uint16
}
