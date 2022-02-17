package config

import (
	"log"
	"os"

	"github.com/streadway/amqp"
)

type RabbitMQ struct {
	Host               string
	RabbitMQConnection *amqp.Connection
	RabbitMQChannel    *amqp.Channel
}

func SetupRabbitMQ() (*RabbitMQ, error) {
	host := os.Getenv("RABBITMQ_HOST")
	conn, ch, err := Connect(host)

	rabbitMQ := RabbitMQ{
		Host:               host,
		RabbitMQConnection: conn,
		RabbitMQChannel:    ch,
	}

	return &rabbitMQ, err
}

func handleError(err error, msg string) {
	if err != nil {
		log.Fatalf("%s: %s", msg, err)
	}
}

// Connect ...
func Connect(host string) (*amqp.Connection, *amqp.Channel, error) {
	var (
		err         error
		conn        *amqp.Connection
		amqpChannel *amqp.Channel
	)

	conn, err = amqp.Dial(host)
	handleError(err, "Can't connect to AMQP")

	amqpChannel, err = conn.Channel()
	handleError(err, "Can't create a amqpChannel")

	return conn, amqpChannel, err
}
