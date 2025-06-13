package models

type SKU struct {
	ID       SKUID
	Name     string
	Count    uint16
	Type     string
	Price    uint32
	Location string
	UserId   UserID
}
