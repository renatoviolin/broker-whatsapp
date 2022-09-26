package main

import (
	"broker/app/services"
	"broker/controller"
	"broker/infra/http_client"
	"broker/infra/kafka"
	"broker/infra/logger"
	"broker/infra/redis"
	"broker/util"
	"os"

	"github.com/Shopify/sarama"
)

func init() {
	logger.InitLog()
	err := util.LoadVars()
	if err != nil {
		logger.Error("main", "load-vars", err.Error())
	}
}

func main() {
	servers := []string{os.Getenv("KAFKA_SERVER")}
	topic := os.Getenv("KAFKA_TOPIC")
	producer := kafka.NewProducer(servers)
	redis := redis.NewRedisClient(os.Getenv("REDIS_CONNECTION"))
	httpClient := http_client.NewClient()
	sender := services.NewSendService(httpClient)
	routerService := services.NewRouterService(producer, redis, sender)
	controller := controller.NewRouterController(routerService)
	handler := NewConsumerHandler(controller)

	kafkaConsumer := kafka.NewConsumer(servers, []string{topic}, "cg-1", handler)
	kafkaConsumer.StartConsumer()
}

type ConsumerHandler struct {
	routerController controller.RouterController
}

func NewConsumerHandler(routerController controller.RouterController) *ConsumerHandler {
	return &ConsumerHandler{routerController: routerController}
}

func (h *ConsumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for {
		select {
		case message := <-claim.Messages():
			err := h.routerController.HandleMessageFromKafka(message.Value, message.Offset)
			if err == nil {
				session.MarkMessage(message, "")
			} else {
				logger.Error("consume-claim", "handle-from-Kafka", err.Error())
			}
		case <-session.Context().Done():
			return nil
		}
	}
}

func (h *ConsumerHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}
