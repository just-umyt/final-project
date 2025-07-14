package producer

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const (
	acks = 0

	//if we want acks -1 uncomment config below
	// retries = 3
	// retryBackoffMs    = 500
	// enableIdempotence = true

	flushTimeout = 5000

	partitionID = 0

	ErrCreateProducer = "error creating kafka producer: %v"
	ErrSendMsg        = "error sending message to kafka: %v"
	ErrMarshallMsg    = "error marshaling message to json: %v"
)

var (
	ErrUnknownType = errors.New("err unknown event type")
)

type Producer struct {
	producer *kafka.Producer
}

func NewProducer(address string) (*Producer, error) {
	config := &kafka.ConfigMap{
		"bootstrap.servers": address,
		"acks":              acks,
		//if we want acks -1 uncomment config below
		// "retries":           retries,
		// "retry.backoff.ms": retryBackoffMs,
		// "enable.idempotence": enableIdempotence,
	}

	prod, err := kafka.NewProducer(config)
	if err != nil {
		return nil, fmt.Errorf(ErrCreateProducer, err)
	}

	return &Producer{producer: prod}, nil
}

func (p *Producer) Produce(dto ProducerMessageDTO, topic string, partionID int32, t time.Time) error {
	message := Message{
		Type:      dto.Type,
		Service:   dto.Service,
		Timestamp: dto.Timestamp.Format(time.RFC3339),
		Payload: Payload{
			SKU:   uint32(dto.SKU),
			Price: dto.Price,
			Count: dto.Count,
		},
	}

	jsonMsg, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf(ErrMarshallMsg, err)
	}

	kafkaMessage := &kafka.Message{
		TopicPartition: kafka.TopicPartition{
			Topic:     &topic,
			Partition: partitionID,
		},
		Key:       nil,
		Value:     jsonMsg,
		Timestamp: t,
	}

	kafkaChan := make(chan kafka.Event)
	if err := p.producer.Produce(kafkaMessage, kafkaChan); err != nil {
		return fmt.Errorf(ErrSendMsg, err)
	}

	event := <-kafkaChan
	switch e := event.(type) {
	case *kafka.Message:
		return nil
	case kafka.Error:
		return e
	default:
		return ErrUnknownType
	}
}

func (p *Producer) Close() {
	p.producer.Flush(flushTimeout)
	p.producer.Close()
}
