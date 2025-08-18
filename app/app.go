package app

import (
	"log"
	"log/slog"
	"os"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
	"github.com/TrueBlocks/trueblocks-khedra/v5/pkg/types"
	"github.com/TrueBlocks/trueblocks-sdk/v5/services"
	"github.com/urfave/cli/v2"
)

type KhedraApp struct {
	cli            *cli.App
	config         *types.Config
	logger         *types.CustomLogger
	controlSvc     *services.ControlService
	serviceManager *services.ServiceManager
}

// ReloadConfigAndServices reloads the finalized config from disk and restarts all services.
func (k *KhedraApp) ReloadConfigAndServices() error {
	cfg, err := LoadConfig()
	if err != nil {
		k.logger.Error("Failed to reload config", "error", err)
		return err
	}
	k.config = &cfg
	k.logger = types.NewLogger(cfg.Logging)
	slog.SetDefault(k.logger.GetLogger())

	// Re-initialize control service and service manager
	if err := k.initializeControlSvc(); err != nil {
		k.logger.Error("Failed to re-initialize control service", "error", err)
		return err
	}
	// Restart all services
	if err := k.serviceManager.StartAllServices(); err != nil {
		k.logger.Error("Failed to restart services", "error", err)
		return err
	}
	k.logger.Info("Config and services reloaded successfully")
	return nil
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
		"pause":     true,
		"unpause":   true,
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
