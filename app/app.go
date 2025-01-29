package app

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/types"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/utils"
	"github.com/urfave/cli/v2"
)

type KhedraApp struct {
	cli       *cli.App
	config    *types.Config
	logger    *types.CustomLogger
	chainList *types.ChainList
}

func NewKhedraApp() *KhedraApp {
	var err error
	k := KhedraApp{}

	// If khedra is already running, one of these ports is serving the
	// control API. We need to make sure it's not running and fail if
	// it is.
	cntlSvcPorts := []string{"8338", "8337", "8336", "8335"}
	for _, port := range cntlSvcPorts {
		if utils.PingServer("http://localhost:" + port) {
			msg := fmt.Sprintf("Error: Khedra is already running (control service port :%s is in use). Quitting...", port)
			fmt.Println(colors.Red+msg, colors.Off)
			os.Exit(1)
		}
	}

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
	k.logger = types.NewLogger(cfg.Logging)
	slog.SetDefault(k.logger.GetLogger())

	return cfg, nil
}
