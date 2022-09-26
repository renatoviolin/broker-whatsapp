package kafka

import (
	"broker/infra/logger"
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
)

type Consumer struct {
	handler  sarama.ConsumerGroupHandler
	consumer sarama.ConsumerGroup
	topic    []string
}

func NewConsumer(servers []string, topic []string, consumerGroup string, handler sarama.ConsumerGroupHandler) *Consumer {
	logger.Info("consumer", "new-consumer", "connecting to: "+servers[0])
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Net.DialTimeout = 10 * time.Second
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Consumer.Group.Rebalance.Timeout = 1 * time.Second
	consumer, err := sarama.NewConsumerGroup(servers, consumerGroup, config)
	if err != nil {
		logger.Fatal("consumer", "new-consumer", err.Error())
	}
	logger.Info("consumer", "new-consumer", "connected: ready to consume topic: "+topic[0])
	return &Consumer{
		consumer: consumer,
		handler:  handler,
		topic:    topic,
	}
}

func (h *Consumer) StartConsumer() {
	logger.Info("consumer", "start-consumer", "starting consumer in separeted goroutine....")
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		for {
			err := h.consumer.Consume(context.Background(), h.topic, h.handler)
			if err != nil {
				logger.Error("consumer", "consumer-goroutine", err.Error())
			}
		}
	}()

	<-sigterm
	logger.Info("consumer", "start-consumer", "terminating: via signal")
}
