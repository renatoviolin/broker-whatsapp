package http_client

import (
	"broker/infra/logger"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"time"
)

type Client struct {
	client http.Client
}

func NewClient() Client {
	return Client{
		client: http.Client{Timeout: 5 * time.Second, Transport: &http.Transport{}},
	}
}

func (h *Client) Post(url string, bodyRequest []byte) (body []byte, statusCode int, err error) {
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(bodyRequest))
	req.Header.Add("Content-Type", "application/json")
	req.Close = true
	if err != nil {
		logger.Error("http_client", "post", fmt.Sprintf("could not create request: %s", err))
		return []byte{}, 500, err
	}

	return h.processRequest(req)
}

func (h *Client) processRequest(req *http.Request) (body []byte, statusCode int, err error) {
	res, err := h.client.Do(req)
	if err != nil {
		return body, 500, err
	}
	if res != nil {
		body, err = io.ReadAll(res.Body)
		if err != nil {
			return body, res.StatusCode, err
		}
	}

	defer res.Body.Close()
	if res.StatusCode >= 400 {
		return body, res.StatusCode, fmt.Errorf("error while requesting: %s", string(body))
	}
	return body, res.StatusCode, nil
}
