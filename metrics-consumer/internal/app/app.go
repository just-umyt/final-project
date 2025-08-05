package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"metrics-consumer/internal/consumer"
	"metrics-consumer/internal/handler"
)

const (
	ErrLoadEnv = "error loading .env file: %v"
)

func RunApp() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	address := os.Getenv("KAFKA_BROKERS")

	topic := os.Getenv("KAFKA_TOPIC")

	consumerGroup := os.Getenv("KAFKA_CONSUMER_GROUP")

	hand := handler.NewHandler()

	cons, err := consumer.NewConsumer(hand, address, topic, consumerGroup)
	if err != nil {
		return err
	}

	go func() {
		cons.Start(ctx)
	}()

	<-ctx.Done()

	return cons.Stop()
}
