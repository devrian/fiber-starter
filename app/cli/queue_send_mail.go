package cli

import (
	"encoding/json"
	"fiber-starter/app/service"
	"fiber-starter/pkg/utils"
	"log"
	"os"

	"github.com/urfave/cli/v2"
)

func (cliApp *CliApp) QueueSendMailHandler(c *cli.Context) error {
	conn := cliApp.RabbitMQ.RabbitMQConnection
	ch := cliApp.RabbitMQ.RabbitMQChannel
	defer conn.Close()
	defer ch.Close()

	queue, err := ch.QueueDeclare(utils.CMDQueueSendMail, true, false, false, false, nil)
	if err != nil {
		log.Fatalf("%s: %v", "[CLI - QueueSendMailHandler] Could not declare `add` queue", err)
		return err
	}

	err = ch.Qos(1, 0, false)
	if err != nil {
		log.Fatalf("%s: %v", "[CLI - QueueSendMailHandler] Could not configure QoS", err)
		return err
	}

	msgCh, err := ch.Consume(queue.Name, "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("%s: %s", "Could not register consumer", err)
		return err
	}

	stopChan := make(chan bool)

	go func() {
		log.Printf("Consumer ready, PID: %d", os.Getpid())
		for d := range msgCh {
			var qData utils.MailData
			err := json.Unmarshal(d.Body, &qData)
			if err != nil {
				log.Printf("Error decoding JSON: %s", err)
			}

			mailService := service.NewMailService()
			err = mailService.MailSender(qData.Receivers, qData.Usage, qData.Data)
			if err != nil {
				log.Printf("Error Send Mail: %v", err)
			}

			if err := d.Ack(false); err != nil {
				log.Printf("Error acknowledging message : %v", err)
			} else {
				log.Printf("Acknowledged message")
			}
		}
	}()

	// Stop for program termination
	<-stopChan

	return nil
}
