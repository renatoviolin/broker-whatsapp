package util

import (
	"broker/dto"
	"broker/infra/logger"
	"encoding/json"
)

func GetPayload1() (payload dto.WebhookInput) {
	err := json.Unmarshal([]byte(payload1), &payload)
	if err != nil {
		logger.Error("util", "get-payload", err.Error())
	}
	return
}

func GetPayload2() (payload dto.WebhookInput) {
	err := json.Unmarshal([]byte(payload2), &payload)
	if err != nil {
		logger.Error("util", "get-payload", err.Error())
	}
	return
}

func GetPayload3() (payload dto.WebhookInput) {
	err := json.Unmarshal([]byte(payload3), &payload)
	if err != nil {
		logger.Error("util", "get-payload", err.Error())
	}
	return
}

func GetPayloadStatus() (payload dto.WebhookInput) {
	err := json.Unmarshal([]byte(payloadStatus), &payload)
	if err != nil {
		logger.Error("util", "get-payload", err.Error())
	}
	return
}

func GetPayloadResponseButton() (payload dto.WebhookInput) {
	err := json.Unmarshal([]byte(payloadResponseButton), &payload)
	if err != nil {
		logger.Error("util", "get-payload", err.Error())
	}
	return
}

func GetPayloadInvalid() (payload dto.WebhookInput) {
	err := json.Unmarshal([]byte(payloadInvalid), &payload)
	if err != nil {
		logger.Error("util", "get-payload-invalid", err.Error())
	}
	return
}

func GetPayloadInvalidClientID() (payload dto.WebhookInput) {
	err := json.Unmarshal([]byte(payloadInvalidClientID), &payload)
	if err != nil {
		logger.Error("util", "get-payload-invalid", err.Error())
	}
	return
}
