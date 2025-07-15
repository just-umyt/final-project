package producer

type Payload struct {
	CartID uint32 `json:"cartId"`
	SKU    uint32 `json:"sku"`
	Count  uint16 `json:"count"`
	Status string `json:"status"`
	Reason string `json:"reason,omitempty"`
}

type Message struct {
	Type      string  `json:"type"`
	Service   string  `json:"service"`
	Timestamp string  `json:"timestamp"`
	Payload   Payload `json:"payload"`
}
