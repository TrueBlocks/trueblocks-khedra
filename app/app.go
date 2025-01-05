package app

import (
	"log"
	"log/slog"
	"os"

	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/config"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/types"
	"github.com/urfave/cli/v2"
)

type KhedraApp struct {
	Cli        *cli.App
	config     *types.Config
	fileLogger *slog.Logger
	progLogger *slog.Logger
}

func NewKhedraApp() *KhedraApp {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	fileLogger, progLogger := types.NewLoggers(cfg.Logging)
	cli := initializeCli()

	k := &KhedraApp{
		config:     &cfg,
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
