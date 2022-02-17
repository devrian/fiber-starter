package main

import (
	"fiber-starter/app/api"
	"fiber-starter/app/api/routes"
	"fiber-starter/app/cli"
	"fiber-starter/config"
	"fmt"
	"log"
	"os"
	"runtime"

	command "github.com/urfave/cli/v2"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	c := config.New()
	apiApp := api.New(c)
	cliApp := cli.New(c)

	routes.Configure(apiApp)

	cmd := &command.App{
		Name:  c.App.Name,
		Usage: fmt.Sprintf(`%s, BE App`, c.App.Name),
		Commands: []*command.Command{
			{
				Name:   "api",
				Usage:  "sample API service",
				Flags:  apiApp.Flags(),
				Action: apiApp.Start,
			},
			{
				Name:   "cmd",
				Usage:  "sample Command service",
				Flags:  cliApp.Flags(),
				Action: cliApp.Start,
			},
		},
		Action: func(cli *command.Context) error {
			fmt.Printf("%s version:%s\n", cli.App.Name, c.App.Version)
			return nil
		},
	}

	err := cmd.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
