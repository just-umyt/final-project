package models

type SKU struct {
	SkuID SKUID
	Name  string
	Type  string
}

type Stock struct {
	Id       StockID
	SkuId    SKUID
	Count    uint16
	Price    uint32
	Location string
	UserId   UserID
}

type FullStock struct {
	SKU
	Stock
}
