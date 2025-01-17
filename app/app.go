package app

import (
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

func (k *KhedraApp) ConfigMaker() (types.Config, error) {
	cfg, err := LoadConfig()
	if err != nil {
		return types.Config{}, err
	}
	k.config = &cfg
	k.fileLogger, k.progLogger = types.NewLoggers(cfg.Logging)
	return cfg, nil
}
