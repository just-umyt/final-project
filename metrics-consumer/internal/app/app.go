package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"metrics-consumer/internal/config"
	"metrics-consumer/internal/consumer"
	"metrics-consumer/internal/handler"
)

var (
	ErrLoadEnv = "error loading .env file: %v"
)

func RunApp(env string) error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	if err := config.LoadConfig(env); err != nil {
		err = fmt.Errorf(ErrLoadEnv, err)
		return err
	}

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
