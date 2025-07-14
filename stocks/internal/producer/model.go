package producer

type Payload struct {
	SKU   uint32 `json:"sku"`
	Count uint16 `json:"count"`
	Price uint32 `json:"price"`
}

type Message struct {
	Type      string  `json:"type"`
	Service   string  `json:"service"`
	Timestamp string  `json:"timestamp"`
	Payload   Payload `json:"payload"`
}
