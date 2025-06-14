package dto

type AddSkuDto struct {
	UserId   int64  `json:"user_id"`
	SkuId    uint32 `json:"sku"`
	Count    uint16 `json:"count"`
	Name     string `json:"name"`
	Type     string `json:"type"`
	Price    uint32 `json:"price"`
	Location string `json:"location"`
}

type GetSkuDto struct {
	Name  string  `json:"name"`
	Price float64 `json:"price"`
	Count int     `json:"count"`
	Type  string  `json:"type"`
}

type GetSkuBySkuId struct {
	SkuId uint32 `json:"sku"`
}
