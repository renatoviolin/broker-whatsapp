package services

import (
	"broker/app/entity"
	"broker/dto"
	"broker/infra/kafka"
	"broker/infra/logger"
	"broker/infra/redis"
	"encoding/json"
	"errors"
	"fmt"
	"os"
)

const (
	MESSAGE_FROM_USER   = "message-from-user"
	MESSAGE_FROM_STATUS = "message-from-status"
	ROUTER_TOPIC        = "router-topic"
	ROUTER_ID           = "router"
	CLIENT_A            = "client_a"
	CLIENT_B            = "client_b"
	CLIENT_C            = "client_c"
)

var Id2Topic = map[string]string{
	"client_a": "client_A_in",
	"client_b": "client_B_in",
	"client_c": "client_C_in",
}

type RouterService struct {
	redis      *redis.RedisClient
	producer   *kafka.Producer
	sender     Sender
	errorTopic string
}

func NewRouterService(producer *kafka.Producer, redis *redis.RedisClient, sender Sender) RouterService {
	return RouterService{
		producer:   producer,
		redis:      redis,
		errorTopic: os.Getenv("KAFKA_ERROR_TOPIC"),
		sender:     sender,
	}
}

func (h *RouterService) ExecuteRouter(input dto.WebhookInput) error {
	inputBytes, err := json.Marshal(input)
	if err != nil {
		return err
	}

	// 1. read the ID whatsapp number
	waID, origin, err := h.GetWhatsappMetadata(input)
	if err != nil {
		return err
	}

	text, _ := h.GetWhatsappResponseText(input)
	if text == "sair" {
		if err = h.UpdateStatus(waID, ROUTER_ID, origin); err != nil {
			return err
		}
		if err := h.sender.SendText("Saindo....", waID); err != nil {
			logger.Error("router", "generate-menu", err.Error())
			return err
		}
		if err := h.GenerateMenu(waID); err != nil {
			logger.Error("router", "generate-menu", err.Error())
			return err
		}
		return nil
	}

	// 2. check in redis if the ID is associated with any client
	// if expired or not found, means the first contact from the client
	clientID, err := h.redis.Get(waID)
	if err != nil {
		if origin == MESSAGE_FROM_USER {
			if err := h.GenerateMenu(waID); err != nil {
				logger.Error("router", "generate-menu", err.Error())
				return err
			}
		}
		if err = h.UpdateStatus(waID, ROUTER_ID, origin); err != nil {
			return err
		}
		if _, _, err = h.producer.Produce(inputBytes, ROUTER_TOPIC); err != nil {
			return err
		}
		return nil
	}

	// 3. check if the message is set as ROUTER owner
	if clientID == ROUTER_ID {
		desiredClientID, err := h.ReadMenu(input)
		if err != nil {
			if origin == MESSAGE_FROM_USER {
				if err := h.GenerateMenu(waID); err != nil {
					logger.Error("router", "generate-menu", err.Error())
					return err
				}
			}
			if err = h.UpdateStatus(waID, ROUTER_ID, origin); err != nil {
				return err
			}
			if _, _, err = h.producer.Produce(inputBytes, ROUTER_TOPIC); err != nil {
				return err
			}
			return nil
		}

		// has a valid desiredClientID, set to that desired clientID
		if err = h.UpdateStatus(waID, desiredClientID, origin); err == nil {
			if origin == MESSAGE_FROM_USER {
				if err := h.sender.SendText(fmt.Sprintf("Ok, estou te transferindo para: %s", desiredClientID), waID); err != nil {
					logger.Error("router", "generate-menu", err.Error())
					return err
				}
			}
			_, _, err = h.producer.Produce(inputBytes, ROUTER_ID)
		}
		return err
	}

	// 4. has a valid client_id, send to specific topic
	if err = h.UpdateStatus(waID, clientID, origin); err != nil {
		return err
	}
	_, _, err = h.producer.Produce(inputBytes, Id2Topic[clientID])
	return err
}

func (h *RouterService) GenerateMenu(waID string) error {
	actions := []entity.Row{
		{ID: "client_a", Title: "Cliente A"},
		{ID: "client_b", Title: "Cliente B"},
		{ID: "client_c", Title: "Cliente C"},
	}

	return h.sender.SendList("Escolha o sistema:", actions, waID)
}

func (h *RouterService) ReadMenu(input dto.WebhookInput) (string, error) {
	clientID, err := h.GetWhatsappClientID(input)
	if err != nil {
		return "", err
	}
	switch clientID {
	case CLIENT_A:
		return CLIENT_A, nil
	case CLIENT_B:
		return CLIENT_B, nil
	case CLIENT_C:
		return CLIENT_C, nil
	default:
		return "", errors.New("client not found in menu response")
	}
}

func (h *RouterService) UpdateStatus(whatsappID string, clientID string, origin string) (err error) {
	if origin == MESSAGE_FROM_STATUS {
		return nil
	}
	return h.redis.Save(whatsappID, clientID)
}

func (h *RouterService) GetWhatsappMetadata(input dto.WebhookInput) (waID string, origin string, err error) {
	var tempWaID string
	if len(input.Entry) > 0 {
		if len(input.Entry[0].Changes) > 0 {
			if len(input.Entry[0].Changes[0].Value.Contacts) > 0 {
				tempWaID = input.Entry[0].Changes[0].Value.Contacts[0].WaID
				if tempWaID != "" {
					waID = tempWaID
				}
				origin = MESSAGE_FROM_USER
			} else if len(input.Entry[0].Changes[0].Value.Statuses) > 0 {
				tempWaID = input.Entry[0].Changes[0].Value.Statuses[0].RecipientID
				if tempWaID != "" {
					waID = tempWaID
				}
				origin = MESSAGE_FROM_STATUS
			}
		}
	}
	if waID == "" || origin == "" {
		return "", "", errors.New("unable to find whatsapp_number/origin in payload")
	}
	return waID, origin, nil
}

func (h *RouterService) GetWhatsappResponseText(input dto.WebhookInput) (text string, err error) {
	if len(input.Entry) > 0 {
		if len(input.Entry[0].Changes) > 0 {
			if len(input.Entry[0].Changes[0].Value.Messages) > 0 {
				text = input.Entry[0].Changes[0].Value.Messages[0].Text.Body
			}
		}
	}
	if text == "" {
		return "", errors.New("unable to find response text")
	}
	return text, nil
}

func (h *RouterService) GetWhatsappClientID(input dto.WebhookInput) (clientID string, err error) {
	if len(input.Entry) > 0 {
		if len(input.Entry[0].Changes) > 0 {
			if len(input.Entry[0].Changes[0].Value.Messages) > 0 {
				clientID = input.Entry[0].Changes[0].Value.Messages[0].Interactive.ListReply.ID
			}
		}
	}
	if clientID == "" {
		return "", errors.New("unable to find clientID")
	}
	return clientID, nil
}
