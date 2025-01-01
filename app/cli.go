package app

import (
	"log/slog"

	"github.com/urfave/cli/v2"
)

func initializeCli() *cli.App {
	return &cli.App{
		Name:  "khedra",
		Usage: "A tool to index, monitor, serve, and share blockchain data",
		Commands: []*cli.Command{
			{
				Name:  "init",
				Usage: "Initializes Khedra",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "mode",
						Usage: "Initialization mode (all, blooms, none)",
					},
				},
				Action: func(c *cli.Context) error {
					slog.Info("command calls: init", "mode", c.String("mode"))
					return nil
				},
			},
			{
				Name:  "scrape",
				Usage: "Controls the blockchain scraper",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "state",
						Usage: "Scraper state (on, off)",
					},
				},
				Action: func(c *cli.Context) error {
					slog.Info("command calls: scrape", "state", c.String("state"))
					return nil
				},
			},
			{
				Name:  "api",
				Usage: "Starts or stops the API server",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:  "state",
						Usage: "API state (on, off)",
					},
				},
				Action: func(c *cli.Context) error {
					slog.Info("command calls: api", "state", c.String("state"))
					return nil
				},
			},
			{
				Name:  "sleep",
				Usage: "Sets the duration between updates",
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:  "duration",
						Usage: "Sleep duration in seconds",
					},
				},
				Action: func(c *cli.Context) error {
					slog.Info("command calls: sleep", "duration", c.Int("duration"))
					return nil
				},
			},
		},
	}
}
