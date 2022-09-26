package services

import (
	"broker/app/entity"
	"broker/infra/http_client"
	"encoding/json"
	"errors"
)

var (
	ErrHttpRequest = errors.New("error while request whatsapp URL")
)

type Sender struct {
	httpClient http_client.Client
}

func NewSendService(httpClient http_client.Client) Sender {
	return Sender{
		httpClient: httpClient,
	}
}

var URL = "https://graph.facebook.com/v12.0/102261592636695/messages?access_token=EAAH5GKHVwZAgBALF4jZC7rJyZB3WkWFdZB7DZB6nApUamPRWMMtFYEaJZCr8nnYsCCun0ZAOBpvmeZCCHZAqJcESdlUXHI32mkWw77NHhZAxI3Pfyu9atWv1LyQ9gpLGlr5hYxuuUGW0xmZC2zbCXfwjWZAY3pmTfqZAjwgZBFdUsJBlZCUB39CBrOpiToMMCmR7TbnHQrIsmGNB1CxggZDZD"

func (h *Sender) SendText(text string, to string) error {
	payload, err := entity.NewTextPayload(text, to)
	if err != nil {
		return err
	}

	bytePayload, _ := json.Marshal(payload)
	_, statusCode, err := h.httpClient.Post(URL, bytePayload)
	if err != nil {
		return err
	}
	if statusCode >= 400 {
		return ErrHttpRequest
	}
	return nil
}

func (h *Sender) SendList(body string, rows []entity.Row, to string) error {
	payload, err := entity.NewListPayload(body, rows, to)
	if err != nil {
		return err
	}
	bytePayload, _ := json.Marshal(payload)
	_, statusCode, err := h.httpClient.Post(URL, bytePayload)
	if err != nil {
		return err
	}
	if statusCode >= 400 {
		return ErrHttpRequest
	}
	return nil
}
