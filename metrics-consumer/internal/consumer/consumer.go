package consumer

import (
	"context"
	"log"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

const (
	sessionTimeoutMs = 7000
	readTimeoutMs    = 50000
)

type IHandler interface {
	HandleMessage(message []byte, offset kafka.Offset, partition int32)
}

type Consumer struct {
	handler  IHandler
	consumer *kafka.Consumer
}

func NewConsumer(handler IHandler, address string, topic, consumerGroup string) (*Consumer, error) {
	config := &kafka.ConfigMap{
		"bootstrap.servers":        address,
		"group.id":                 consumerGroup,
		"session.timeout.ms":       sessionTimeoutMs,
		"enable.auto.offset.store": false,
		"enable.partition.eof":     false,
		"enable.auto.commit":       false,
		"auto.commit.interval.ms":  5000,
		"auto.offset.reset":        "earliest",
	}

	consumer, err := kafka.NewConsumer(config)
	if err != nil {
		return nil, err
	}

	if err = consumer.Subscribe(topic, nil); err != nil {
		return nil, err
	}

	return &Consumer{
		handler:  handler,
		consumer: consumer,
	}, nil
}

func (c *Consumer) Start(ctx context.Context) {
	for ctx.Err() == nil {

		kafkaMsg, err := c.consumer.ReadMessage(readTimeoutMs)
		if err != nil {
			if err.Error() == kafka.ErrTimedOut.String() {
				continue
			}
			log.Printf("error kafka read message: %v", err)
		}

		if kafkaMsg == nil {
			continue
		}

		c.handler.HandleMessage(kafkaMsg.Value, kafkaMsg.TopicPartition.Offset, kafkaMsg.TopicPartition.Partition)

		if _, err := c.consumer.StoreMessage(kafkaMsg); err != nil {
			log.Printf("error kafka store message: %v", err)
			continue
		}
	}
}

func (c *Consumer) Stop() error {

	if _, err := c.consumer.Commit(); err != nil {
		return err
	}

	return c.consumer.Close()
}
