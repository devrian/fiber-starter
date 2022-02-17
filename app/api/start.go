package api

import (
	"context"
	"fiber-starter/config"
	"fiber-starter/db"
	"fmt"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/urfave/cli/v2"
)

type ApiApp struct {
	Config    *config.Config
	DB        *pgxpool.Pool
	Fiber     *fiber.App
	Validator *config.Validator
	Redis     *config.Redis
	RabbitMQ  *config.RabbitMQ
}

func New(c *config.Config) *ApiApp {
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

	return &ApiApp{
		Config:    c,
		DB:        db.Init(c),
		Fiber:     fiber.New(fiberConfig),
		Validator: config.SetupValidator(&c.App),
		Redis:     redis,
		RabbitMQ:  rabbitMQ,
	}
}

func (app *ApiApp) Start(c *cli.Context) error {
	// Create channel for idle connections.
	sng := make(chan os.Signal, 1)
	signal.Notify(sng, os.Interrupt) // Catch OS signals.
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	go func() {
		for range sng {
			fmt.Println("Shutting down..")

			// Received an interrupt signal, shutdown.
			if err := app.Fiber.Shutdown(); err != nil {
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
	err := app.Fiber.Listen(fmt.Sprintf(":%s", app.Config.App.Port))
	if err != nil {
		log.Printf("Oops... Server is not running! Reason: %v", err)
	}

	return err
}

func (app *ApiApp) Flags() []cli.Flag {
	flags := []cli.Flag{
		&cli.StringFlag{
			Name:  "host",
			Value: fmt.Sprintf(`%s:%s`, app.Config.App.Host, app.Config.App.Port),
			Usage: "Run API service with custom host",
		},
	}

	return flags
}
