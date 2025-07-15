package producer

import (
	"stocks/internal/models"
	"time"
)

type ProducerMessageDTO struct {
	Type      string
	Service   string
	Timestamp time.Time
	SKU       models.SKUID
	Count     uint16
	Price     uint32
}
