package app

import (
	"log"
	"log/slog"
	"os"

	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/types"
	"github.com/urfave/cli/v2"
)

type KhedraApp struct {
	cli        *cli.App
	config     *types.Config
	fileLogger *slog.Logger
	progLogger *slog.Logger
}

func (k *KhedraApp) Run() {
	k.cli = initCli(k)
	k.cli.Run(os.Args)
}

func (k *KhedraApp) ConfigMaker() {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}
	k.config = &cfg
	k.fileLogger, k.progLogger = types.NewLoggers(cfg.Logging)
}
