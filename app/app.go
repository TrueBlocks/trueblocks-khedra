package app

import (
	"log/slog"
	"os"

	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/config"
	"github.com/urfave/cli/v2"
)

type KhedraApp struct {
	Cli        *cli.App
	Config     *config.Config
	FileLogger *slog.Logger
	ProgLogger *slog.Logger
}

func NewKhedraApp() *KhedraApp {
	cfg := config.MustLoadConfig("config.yaml")
	fileLogger, progLogger := config.NewLoggers(cfg.Logging)
	cli := initializeCli()

	k := &KhedraApp{
		Config:     cfg,
		FileLogger: fileLogger,
		ProgLogger: progLogger,
		Cli:        cli,
	}

	return k
}

// Run runs the Khedra cli
func (k *KhedraApp) Run() error {
	return k.Cli.Run(os.Args)
}
