package producer

import (
	"cart/internal/models"
	"time"
)

type ProducerMessageDTO struct {
	Type      string
	Service   string
	Timestamp time.Time
	CartID    models.CartID
	SKU       models.SKUID
	Count     uint16
	Status    string
	Reason    string
}
