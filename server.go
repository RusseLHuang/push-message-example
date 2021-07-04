package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/RusseLHuang/push-message-example/message"
	messagebroker "github.com/RusseLHuang/push-message-example/message_broker"
	"github.com/RusseLHuang/push-message-example/node"
	pushregistry "github.com/RusseLHuang/push-message-example/push_registry"
	wsconnection "github.com/RusseLHuang/push-message-example/ws_connection"
	"github.com/gorilla/mux"
	"github.com/spf13/viper"
)

func main() {
	flag.Parse()
	log.SetFlags(0)

	viper.AddConfigPath(".")
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	err := viper.ReadInConfig() // Find and read the config file
	if err != nil {             // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	nodeService := node.NewNode()
	pushRegistry := pushregistry.NewPushRegistry()
	wsConnectionManager := wsconnection.NewWSConnectionManager(pushRegistry, nodeService)
	messageBrokerClient := messagebroker.NewMessageBrokerClient()

	wsConnectionController := wsconnection.NewWSConnectionController(wsConnectionManager, messageBrokerClient)
	messageController := message.NewMessageController(wsConnectionManager)

	router := mux.NewRouter()

	router.HandleFunc("/connect/{clientID}", wsConnectionController.Connect)
	router.HandleFunc("/message/client/{clientID}", messageController.Send)

	log.Println("Starting Server")
	src := &http.Server{
		Handler: router,
	}
	src.ListenAndServe()
}
