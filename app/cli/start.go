package cli

import (
	"context"
	"fiber-starter/config"
	"fiber-starter/db"
	"fiber-starter/pkg/utils"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/urfave/cli/v2"
)

type CliApp struct {
	Config    *config.Config
	DB        *pgxpool.Pool
	Fiber     *fiber.App
	Validator *config.Validator
	Redis     *config.Redis
	RabbitMQ  *config.RabbitMQ
}

func New(c *config.Config) *CliApp {
	fiberConfig := fiber.Config{
		AppName: c.App.Name,
	}

	redis, err := config.SetupRedis()
	if err != nil {
		log.Fatalln(err)
	}

	rabbitMQ, err := config.SetupRabbitMQ()
	if err != nil {
		log.Fatalln(err)
	}

	return &CliApp{
		Config:    c,
		DB:        db.Init(c),
		Fiber:     fiber.New(fiberConfig),
		Validator: config.SetupValidator(&c.App),
		Redis:     redis,
		RabbitMQ:  rabbitMQ,
	}
}

func (cliApp *CliApp) Start(c *cli.Context) error {
	switch queue := c.String("queue"); queue {
	case utils.CMDQueueSendMail:
		return cliApp.QueueSendMailHandler(c)
	}

	// Create channel for idle connections.
	sng := make(chan os.Signal, 1)
	signal.Notify(sng, os.Interrupt) // Catch OS signals.
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	go func() {
		for range sng {
			fmt.Println("Shutting down..")
			cliApp.RabbitMQ.RabbitMQConnection.Close()

			// Received an interrupt signal, shutdown.
			if err := cliApp.Fiber.Shutdown(); err != nil {
				// Error from closing listeners, or context timeout:
				log.Printf("Oops... Server is not shutting down! Reason: %v", err)
			}

			select {
			case <-time.After(21 * time.Second):
				fmt.Println("Not all connections done")
			case <-ctx.Done():

			}
		}
	}()

	// Run server.
	err := cliApp.Fiber.Listen(fmt.Sprintf(":%s", cliApp.Config.App.Port))
	if err != nil {
		log.Printf("Oops... Server is not running! Reason: %v", err)
	}

	return err
}

func (cliApp *CliApp) Flags() []cli.Flag {
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:  "queue",
			Value: "",
			Usage: "Run queue subsrciber in this server",
		},
	}

	return flags
}
