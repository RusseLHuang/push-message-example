package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	messagebroker "github.com/RusseLHuang/push-message-example/message_broker"
	pushregistry "github.com/RusseLHuang/push-message-example/push_registry"
)

type ConsumerHandler struct {
	PushRegistry        *pushregistry.PushRegistry
	MessageBrokerClient *messagebroker.MessageBrokerClient
}

func NewConsumerHandler(
	pushRegistry *pushregistry.PushRegistry,
	messageBrokerClient *messagebroker.MessageBrokerClient,
) *ConsumerHandler {
	return &ConsumerHandler{
		PushRegistry:        pushRegistry,
		MessageBrokerClient: messageBrokerClient,
	}
}

func (ch *ConsumerHandler) Handle(key string) {

	connectionEndpoint := ch.PushRegistry.GetPersistentConnectionID(key)
	if connectionEndpoint != "" {
		afterTime := 10 * time.Second

		ticker := time.After(afterTime)
		<-ticker

		endpoint := fmt.Sprintf("http://%s:80/message/client/%s", connectionEndpoint, key)
		log.Println("Endpoint: %s", endpoint)

		resp, err := http.Get(endpoint)
		if err != nil {
			log.Fatal("Fail to call server endpoint")
		}

		fmt.Println(resp)

		defer resp.Body.Close()

		ch.MessageBrokerClient.Send(key, "movie_recommendation")
	}

}
