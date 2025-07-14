package app

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"metrics-consumer/internal/consumer"
	"metrics-consumer/internal/handler"
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

	// sigChan := make(chan os.Signal, 1)
	// signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-ctx.Done()

	return cons.Stop()
}
