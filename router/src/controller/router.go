package controller

import (
	"broker/app/services"
	"broker/dto"
	"broker/infra/logger"
	"encoding/json"
)

type RouterController struct {
	routerService services.RouterService
}

func NewRouterController(routerService services.RouterService) RouterController {
	return RouterController{routerService: routerService}
}

func (h *RouterController) HandleMessageFromKafka(input []byte, offset int64) error {
	var payload dto.WebhookInput
	err := json.Unmarshal(input, &payload)
	if err != nil {
		logger.Error("controller", "handle-from-kafka-1", err.Error())
		return err
	}

	err = h.routerService.ExecuteRouter(payload)
	if err != nil {
		logger.Error("controller", "handle-from-kafka-2", err.Error())
		return err
	}

	return nil
}
