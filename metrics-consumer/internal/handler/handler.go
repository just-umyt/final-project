package handler

import (
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) HandleMessage(message []byte, offset kafka.Offset, partition int32) {
	log.Printf("Message: [%s], Offset: [%d], Partition: [%d]", message, offset, partition)
}
