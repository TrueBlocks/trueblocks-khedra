package app

import (
	"log"
	"log/slog"
	"os"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
	"github.com/TrueBlocks/trueblocks-khedra/v5/pkg/types"
	"github.com/urfave/cli/v2"
)

type KhedraApp struct {
	cli    *cli.App
	config *types.Config
	logger *types.CustomLogger
}

func NewKhedraApp() *KhedraApp {
	k := KhedraApp{}
	if k.isRunning() {
		log.Fatal(colors.BrightBlue + "khedra is already running - cannot run..." + colors.Off)
	}
	k.cli = initCli(&k)
	return &k
}

func (k *KhedraApp) Run() {
	_ = k.cli.Run(os.Args)
}

func (k *KhedraApp) isRunning() bool {
	okArgs := map[string]bool{
		"help":      true,
		"-h":        true,
		"--help":    true,
		"version":   true,
		"-v":        true,
		"--version": true,
	}

	if len(os.Args) < 2 || len(os.Args) == 2 && os.Args[1] == "config" {
		return false
	}

	for i, arg := range os.Args {
		if okArgs[arg] {
			return false
		} else if arg == "config" && i < len(os.Args)-1 && os.Args[i+1] == "show" {
			return false
		}
	}

	ports := []string{"8338", "8337", "8336", "8335"}
	for _, port := range ports {
		if utils.PingServer("http://localhost:" + port) {
			return true
		}
	}

	return false
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
