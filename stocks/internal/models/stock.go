package models

type SKU struct {
	ID   SKUID
	Name string
	Type string
}

type Stock struct {
	ID       StockID
	SKUID    SKUID
	Count    uint16
	Price    uint32
	Location string
	UserID   UserID
}

type Item struct {
	SKU   SKU
	Stock Stock
}
