package services

import (
	"broker/app/entity"
	"broker/infra/http_client"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_SendText(t *testing.T) {
	httpClient := http_client.NewClient()
	service := NewSendService(httpClient)
	err := service.SendText("enviando a partir do golang", "5516993259256")
	require.NoError(t, err)
}

func Test_InvalidNumber(t *testing.T) {
	httpClient := http_client.NewClient()
	service := NewSendService(httpClient)
	err := service.SendText("enviando a partir do golang", "5516993259255")
	require.Error(t, err)
}

func Test_SendList(t *testing.T) {
	httpClient := http_client.NewClient()
	service := NewSendService(httpClient)

	actions := []entity.Row{
		{ID: "client_a", Title: "Cliente A"},
		{ID: "client_b", Title: "Cliente B"},
		{ID: "client_c", Title: "Cliente C"},
	}

	err := service.SendList("Escolha o Sistema", actions, "5516993259256")
	require.NoError(t, err)
}
