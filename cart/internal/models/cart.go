package models

type Cart struct {
	Id     CartID
	UserId UserID
	SKUId  SKUID
	Count  uint16
}
