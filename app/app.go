package app

import (
	"log/slog"
	"os"

	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/config"
	"github.com/urfave/cli/v2"
)

type KhedraApp struct {
	Cli        *cli.App
	config     *config.Config
	fileLogger *slog.Logger
	progLogger *slog.Logger
}

func NewKhedraApp() *KhedraApp {
	cfg := config.MustLoadConfig("config.yaml")
	fileLogger, progLogger := config.NewLoggers(cfg.Logging)
	cli := initializeCli()

	k := &KhedraApp{
		config:     cfg,
		fileLogger: fileLogger,
		progLogger: progLogger,
		Cli:        cli,
	}

	return k
}

// Run runs the Khedra cli
func (k *KhedraApp) Run() error {
	return k.Cli.Run(os.Args)
}
