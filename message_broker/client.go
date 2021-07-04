package messagebroker

import (
	"log"

	"github.com/streadway/amqp"
)

type MessageBrokerClient struct {
	Connection *amqp.Connection
}

func NewMessageBrokerClient() *MessageBrokerClient {
	connection, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
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

func (m *MessageBrokerClient) Receive() {
	m.Connection.Close()
}

func (m *MessageBrokerClient) Close() {
	m.Connection.Close()
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}
