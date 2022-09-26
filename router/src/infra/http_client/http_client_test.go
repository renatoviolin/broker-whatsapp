package http_client

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestHttpClientGet(t *testing.T) {
	url := "https://reqres.in/api/users"
	client := NewClient()
	data := []byte(`{"name": "morpheus","job": "leader"}`)
	res, status, err := client.Post(url, data)
	require.Equal(t, 201, status)
	require.NoError(t, err)
	require.NotEmpty(t, res)
}
