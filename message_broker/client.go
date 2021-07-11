package messagebroker

import (
	"fmt"
	"log"

	"github.com/spf13/viper"
	"github.com/streadway/amqp"
)

type MessageBrokerClient struct {
	Connection *amqp.Connection
}

func NewMessageBrokerClient() *MessageBrokerClient {
	brokerHost := viper.Get("brokerHost")
	brokerPort := viper.Get("brokerPort")
	uri := fmt.Sprintf("amqp://guest:guest@%s:%v/", brokerHost, brokerPort)
	connection, err := amqp.Dial(uri)
	failOnError(err, "Failed to connect to RabbitMQ")

	return &MessageBrokerClient{
		Connection: connection,
	}
}

func (m *MessageBrokerClient) Send(
	message string,
	queueName string,
) {
	channel, err := m.Connection.Channel()
	failOnError(err, "Failed to open channel")

	defer channel.Close()

	queue, err := channel.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	failOnError(err, "Failed to declare a queue")

	err = channel.Publish(
		"",         // exchange
		queue.Name, // routing key
		false,      // mandatory
		false,      // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        []byte(message),
		})

	failOnError(err, "Failed to publish a message")
}

func (m *MessageBrokerClient) Receive(queueName string, handlerFunc func(messageBody string)) {
	channel, err := m.Connection.Channel()
	failOnError(err, "Failed to open channel")

	defer channel.Close()

	queue, err := channel.QueueDeclare(
		queueName, // name
		false,     // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)

	failOnError(err, "Failed to declare a queue")

	msg, err := channel.Consume(
		queue.Name, // queue
		"",         // consumer
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	failOnError(err, "Failed to register a consumer")

	forever := make(chan bool)

	go func() {
		for d := range msg {
			log.Println("Receive a message from queue")
			handlerFunc(string(d.Body))
		}
	}()

	<-forever
}

func (m *MessageBrokerClient) Close() {
	m.Connection.Close()
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
