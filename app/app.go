package app

import (
	"fmt"
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
	chainList  *types.ChainList
}

func NewKhedraApp() *KhedraApp {
	var err error
	k := KhedraApp{}
	if k.chainList, err = types.UpdateChainList(); err != nil {
		fmt.Println(err.Error())
	}
	k.cli = initCli(&k)
	return &k
}

func (k *KhedraApp) Run() {
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
