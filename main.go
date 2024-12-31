package main

import (
	"log"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/config"
	"github.com/urfave/cli/v2"
	"gopkg.in/natefinch/lumberjack.v2"
)

type Khedra struct {
	Config *config.Config
	Logger *slog.Logger
}

func main() {
	cfg := config.MustLoadConfig("config.yaml")
	slog.Info("Logging to", "filename", filepath.Join(cfg.Logging.Folder, cfg.Logging.Filename))

	fileLogger := &lumberjack.Logger{
		Filename:   filepath.Join(cfg.Logging.Folder, cfg.Logging.Filename),
		MaxSize:    cfg.Logging.MaxSizeMb,
		MaxBackups: cfg.Logging.MaxBackups,
		MaxAge:     cfg.Logging.MaxAgeDays,
		Compress:   cfg.Logging.Compress,
	}

	fileHandler := slog.NewJSONHandler(fileLogger, nil)
	logger := slog.New(fileHandler)
	// slog.SetDefault(logger)

	logger.Info("Starting Khedra CLI")
	logger.Info("Logging to", "filename", filepath.Join(cfg.Logging.Folder, cfg.Logging.Filename))
	logger.Info(cfg.String())

	app := &cli.App{
		Name:  "khedra",
		Usage: "A CLI tool for Khedra",
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

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}
