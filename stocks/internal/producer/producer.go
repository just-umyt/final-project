package producer

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const (
	acks         = 0
	flushTimeout = 5000
	partitionID  = 1

	ErrCreateProducer = "error creating kafka producer: %v"
	ErrSendMsg        = "error sending message to kafka: %v"
	ErrMarshallMsg    = "error marshaling message to json: %v"
	ErrKafkaRespond   = "error kafka respond: %v"
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
	}

	prod, err := kafka.NewProducer(config)
	if err != nil {
		return nil, fmt.Errorf(ErrCreateProducer, err)
	}

	return &Producer{producer: prod}, nil
}

func (p *Producer) Produce(dto ProducerMessageDTO, topic string, t time.Time) error {
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
		log.Printf(ErrMarshallMsg, err)
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
		log.Printf(ErrSendMsg, err)
	}

	go func() {
		event := <-kafkaChan

		switch e := event.(type) {
		case *kafka.Message:
			err = nil
		case kafka.Error:
			err = fmt.Errorf(ErrKafkaRespond, e)
		default:
			err = ErrUnknownType
		}
	}()

	return err
}

func (p *Producer) Close() {
	p.producer.Flush(flushTimeout)
	p.producer.Close()
}
