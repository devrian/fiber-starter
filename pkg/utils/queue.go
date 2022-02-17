package utils

import (
	"encoding/json"
	"fiber-starter/config"
	"log"

	"github.com/streadway/amqp"
)

const (
	CMDQueueSendMail = "cmd_queue_send_mail"
)

type QueueRabbitMQ struct {
	RabbitMQ *config.RabbitMQ
	Data     interface{}
}

func PushQueueToRabbitMQ(q QueueRabbitMQ, qName string) error {
	conn, ch, err := config.Connect(q.RabbitMQ.Host)
	if err != nil {
		log.Fatalf("[%s] %s: %v", qName, "Error when connect to queue", err)
	}
	defer conn.Close()
	defer ch.Close()

	queue, err := ch.QueueDeclare(qName, true, false, false, false, nil)
	if err != nil {
		log.Fatalf("[%s] %s: %v", qName, "Could not declare queue", err)
	}

	qDataJSON, err := json.Marshal(q.Data)
	if err != nil {
		log.Fatalf("[%s] %s: %v", qName, "Error encoding JSON", err)
	}

	err = ch.Publish("", queue.Name, false, false, amqp.Publishing{
		DeliveryMode: amqp.Persistent,
		ContentType:  "text/plain",
		Body:         qDataJSON,
	})
	if err != nil {
		log.Fatalf("[%s] %s: %v", qName, "Something error when publish the messages", err)
	}

	return err
}
