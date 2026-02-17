package main

import (
	"context"
	"os"

	"github.com/pdcgo/shared/pkg/cloud_logging"
	"github.com/urfave/cli/v3"
)

type App *cli.Command

func NewApp() App {
	return &cli.Command{
		Name: "tracking service",
		Commands: []*cli.Command{
			{
				Name: "check_order",
			},
		},
		Action: func(ctx context.Context, c *cli.Command) error {

			return nil
		},
	}
}

func main() {
	cloud_logging.SetCloudLoggingDefault()
	app, err := InitializeApp()
	if err != nil {
		panic(err)
	}

	var command *cli.Command = app

	err = command.Run(context.Background(), os.Args)
	if err != nil {
		panic(err)
	}
}
