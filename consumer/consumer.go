package main

import (
	"fmt"

	messagebroker "github.com/RusseLHuang/push-message-example/message_broker"
	pushregistry "github.com/RusseLHuang/push-message-example/push_registry"
	"github.com/spf13/viper"
)

func main() {

	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	pushRegistry := pushregistry.NewPushRegistry()
	messageBrokerClient := messagebroker.NewMessageBrokerClient()

	consumerHandler := NewConsumerHandler(pushRegistry, messageBrokerClient)

	messageBrokerClient.Receive("movie_recommendation", consumerHandler.Handle)
}
