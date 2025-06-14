package models

type SKU struct {
	SkuId    SKUID
	Name     string
	Count    uint16
	Type     string
	Price    uint32
	Location string
	UserId   UserID
}
