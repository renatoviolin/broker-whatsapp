package services

import (
	"broker/infra/http_client"
	"broker/infra/kafka"
	"broker/infra/redis"
	"broker/util"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// Route a Status message (do not update redis and send to ROUTER topic)
func Test_Route_RouterStatusMessage(t *testing.T) {
	payload := util.GetPayloadStatus()
	require.NotNil(t, payload)

	redis := redis.NewRedisClient(os.Getenv("REDIS_CONNECTION"))
	redis.ClearAll()
	producer := kafka.NewProducer([]string{"localhost:9093"})
	httpClient := http_client.NewClient()
	sender := NewSendService(httpClient)
	service := NewRouterService(producer, redis, sender)
	err := service.ExecuteRouter(payload)
	require.NoError(t, err)
	_, err = redis.Get("111")
	require.Error(t, err)
}

// Route first message from client to broker
// - Message must be set to ROUTER_ID
// - Message must be published in ROUTER_TOPIC
// - Menu must be sent to whatsapp
func Test_Route_RouterMessage_1(t *testing.T) {
	payload := util.GetPayload1()
	require.NotNil(t, payload)

	redis := redis.NewRedisClient(os.Getenv("REDIS_CONNECTION"))
	redis.ClearAll()
	producer := kafka.NewProducer([]string{"localhost:9093"})
	httpClient := http_client.NewClient()
	sender := NewSendService(httpClient)
	service := NewRouterService(producer, redis, sender)
	err := service.ExecuteRouter(payload)
	require.NoError(t, err)
	clientID, err := redis.Get("5516993259256")
	require.NoError(t, err)
	require.Equal(t, ROUTER_ID, clientID)
}

// Route response message with invalid clientID
// - Message must be kept to ROUTER_ID
// - Message must be published in ROUTER_TOPIC
// - Menu must be sent to whatsapp again
func Test_Route_RouterMessage_ChooseInvalidClientID(t *testing.T) {
	payload := util.GetPayloadInvalidClientID()
	require.NotNil(t, payload)

	redis := redis.NewRedisClient(os.Getenv("REDIS_CONNECTION"))
	producer := kafka.NewProducer([]string{"localhost:9093"})
	httpClient := http_client.NewClient()
	sender := NewSendService(httpClient)
	service := NewRouterService(producer, redis, sender)
	err := service.ExecuteRouter(payload)
	require.NoError(t, err)
	clientID, err := redis.Get("5516993259256")
	require.NoError(t, err)
	require.Equal(t, "router", clientID)
}

// Route message from client to ROUTER_TOPIC
// - Message must be set to ROUTER_ID
// - Message must be published in ROUTER_TOPIC
// - Must send the menu to user again
func Test_Route_RouterMessage_2(t *testing.T) {
	payload := util.GetPayload2()
	require.NotNil(t, payload)

	redis := redis.NewRedisClient(os.Getenv("REDIS_CONNECTION"))
	producer := kafka.NewProducer([]string{"localhost:9093"})
	httpClient := http_client.NewClient()
	sender := NewSendService(httpClient)
	service := NewRouterService(producer, redis, sender)
	err := service.ExecuteRouter(payload)
	require.NoError(t, err)
	clientID, err := redis.Get("5516993259256")
	require.NoError(t, err)
	require.Equal(t, ROUTER_ID, clientID)
}

// Route message from client to ROUTER_TOPIC
// - Message must be set to CLIENT_A
// - Message must be published in ROUTER_TOPIC
func Test_Route_RouterMessage_2a(t *testing.T) {
	payload := util.GetPayload2()
	require.NotNil(t, payload)

	redis := redis.NewRedisClient(os.Getenv("REDIS_CONNECTION"))
	producer := kafka.NewProducer([]string{"localhost:9093"})
	httpClient := http_client.NewClient()
	sender := NewSendService(httpClient)
	service := NewRouterService(producer, redis, sender)
	err := service.ExecuteRouter(payload)
	require.NoError(t, err)
	clientID, err := redis.Get("5516993259256")
	require.NoError(t, err)
	require.Equal(t, CLIENT_A, clientID)
}

// Route seconds message from client to broker
// - Message already routed to to CLIENT_A
// - Message must be published in CLIENT_A_IN
// - Must update the Expires of the waID
func Test_Route_RouterMessage_3(t *testing.T) {
	payload := util.GetPayload3()
	require.NotNil(t, payload)

	redis := redis.NewRedisClient(os.Getenv("REDIS_CONNECTION"))
	producer := kafka.NewProducer([]string{"localhost:9093"})
	httpClient := http_client.NewClient()
	sender := NewSendService(httpClient)
	service := NewRouterService(producer, redis, sender)
	err := service.ExecuteRouter(payload)
	require.NoError(t, err)
	clientID, err := redis.Get("5516993259256")
	require.NoError(t, err)
	require.Equal(t, CLIENT_A, clientID)
}

func Test_Route_RouterMessage_4_updateTTL(t *testing.T) {
	time.Sleep(3 * time.Second)
	payload := util.GetPayload3()
	require.NotNil(t, payload)

	redis := redis.NewRedisClient(os.Getenv("REDIS_CONNECTION"))
	ttlPrev, err := redis.GetTTL("5516993259256")
	require.NoError(t, err)
	clientID, err := redis.Get("5516993259256")
	require.NoError(t, err)
	require.Equal(t, CLIENT_A, clientID)

	producer := kafka.NewProducer([]string{"localhost:9093"})
	httpClient := http_client.NewClient()
	sender := NewSendService(httpClient)
	service := NewRouterService(producer, redis, sender)
	err = service.ExecuteRouter(payload)
	require.NoError(t, err)
	clientID2, err := redis.Get("5516993259256")
	require.NoError(t, err)
	require.Equal(t, CLIENT_A, clientID2)
	ttlCurrent, err := redis.GetTTL("5516993259256")
	require.NoError(t, err)
	require.Greater(t, ttlCurrent, ttlPrev)
}

// func Test_Route_Client3(t *testing.T) {
// 	payload := util.GetPayload3()
// 	require.NotNil(t, payload)

// 	redis := redis.NewRedisClient(os.Getenv("REDIS_CONNECTION"))
// 	redis.Save("5516993259256", CLIENT_C)

// 	producer := kafka.NewProducer([]string{"localhost:9093"})
// 	httpClient := http_client.NewClient()
// 	sender := NewSendService(httpClient)
// 	service := NewRouterService(producer, redis, sender)
// 	service.ExecuteRouter(payload)
// 	clientID, err := redis.Get("5516993259256")
// 	require.NoError(t, err)
// 	require.Equal(t, CLIENT_C, clientID)
// }

// func Test_Route_ExpiredMessageRedirectToRouters(t *testing.T) {
// 	payload := util.GetPayload3()
// 	require.NotNil(t, payload)

// 	redis := redis.NewRedisClient(os.Getenv("REDIS_CONNECTION"))
// 	redis.Client.Expire(context.Background(), "333", -2)

// 	producer := kafka.NewProducer([]string{"localhost:9093"})
// 	httpClient := http_client.NewClient()
// 	sender := NewSendService(httpClient)
// 	service := NewRouterService(producer, redis, sender)
// 	service.ExecuteRouter(payload)
// 	clientID, err := redis.Get("333")
// 	require.NoError(t, err)
// 	require.Equal(t, ROUTER_ID, clientID)
// }

// // func Test_Route_SendToRouter(t *testing.T) {
// // 	payload := util.GetPayload1()
// // 	require.NotNil(t, payload)

// // 	redis := redis.NewRedisClient(os.Getenv("REDIS_CONNECTION"))
// // 	producer := kafka.NewProducer([]string{"localhost:9093"})
// // 	service := NewRouterService(producer, redis)
// // 	service.ExecuteRouter(payload)
// // }

func Test_GetWhatsappNumber_OK(t *testing.T) {
	payload := util.GetPayloadStatus()
	require.NotNil(t, payload)

	redis := redis.NewRedisClient(os.Getenv("REDIS_CONNECTION"))
	producer := kafka.NewProducer([]string{"localhost:9093"})
	httpClient := http_client.NewClient()
	sender := NewSendService(httpClient)
	service := NewRouterService(producer, redis, sender)

	waID, origin, err := service.GetWhatsappMetadata(payload)
	require.NoError(t, err)
	require.Equal(t, "111", waID)
	require.Equal(t, MESSAGE_FROM_STATUS, origin)

	payload = util.GetPayload2()
	waID, origin, err = service.GetWhatsappMetadata(payload)
	require.NoError(t, err)
	require.Equal(t, "5516993259256", waID)
	require.Equal(t, MESSAGE_FROM_USER, origin)

	payload = util.GetPayload3()
	waID, origin, err = service.GetWhatsappMetadata(payload)
	require.NoError(t, err)
	require.Equal(t, "5516993259256", waID)
	require.Equal(t, MESSAGE_FROM_USER, origin)
}

func Test_GetWhatsappClientID(t *testing.T) {
	payload := util.GetPayloadResponseButton()
	require.NotNil(t, payload)

	redis := redis.NewRedisClient(os.Getenv("REDIS_CONNECTION"))
	producer := kafka.NewProducer([]string{"localhost:9093"})
	httpClient := http_client.NewClient()
	sender := NewSendService(httpClient)
	service := NewRouterService(producer, redis, sender)

	clientID, err := service.GetWhatsappClientID(payload)
	require.NoError(t, err)
	require.Equal(t, "id1", clientID)

	payload = util.GetPayload1()
	clientID, err = service.GetWhatsappClientID(payload)
	require.Error(t, err)
	require.Equal(t, "", clientID)

	payload = util.GetPayload2()
	clientID, err = service.GetWhatsappClientID(payload)
	require.Error(t, err)
	require.Equal(t, "", clientID)

	payload = util.GetPayload3()
	clientID, err = service.GetWhatsappClientID(payload)
	require.Error(t, err)
	require.Equal(t, "", clientID)
}

// func Test_GetWhatsappNumber_NotFound(t *testing.T) {
// 	payload := util.GetPayloadInvalid()
// 	require.NotNil(t, payload)

// 	redis := redis.NewRedisClient(os.Getenv("REDIS_CONNECTION"))
// 	producer := kafka.NewProducer([]string{"localhost:9093"})
// 	httpClient := http_client.NewClient()
// 	sender := NewSendService(httpClient)
// 	service := NewRouterService(producer, redis, sender)

// 	waID, origin, err := service.GetWhatsappMetadata(payload)
// 	require.Error(t, err)
// 	require.Equal(t, "", waID)
// 	require.Equal(t, "", origin)
// }

// func Test_GetWhatsappResponseText(t *testing.T) {
// 	payload := util.GetPayload2()
// 	require.NotNil(t, payload)

// 	redis := redis.NewRedisClient(os.Getenv("REDIS_CONNECTION"))
// 	producer := kafka.NewProducer([]string{"localhost:9093"})
// 	httpClient := http_client.NewClient()
// 	sender := NewSendService(httpClient)
// 	service := NewRouterService(producer, redis, sender)

// 	text, err := service.GetWhatsappResponseText(payload)
// 	require.NoError(t, err)
// 	require.Equal(t, "Resposta do renato", text)

// 	payload = util.GetPayload3()
// 	text, err = service.GetWhatsappResponseText(payload)
// 	require.NoError(t, err)
// 	require.Equal(t, "outra resposta do usuário", text)

// 	payload = util.GetPayload1()
// 	text, err = service.GetWhatsappResponseText(payload)
// 	require.Error(t, err)
// 	require.Equal(t, "", text)
// }
