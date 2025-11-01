package app

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/v5/pkg/colors"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/v5/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/v5/pkg/utils"
	"github.com/TrueBlocks/trueblocks-khedra/v6/pkg/types"
	"github.com/TrueBlocks/trueblocks-sdk/v6/services"
	"github.com/urfave/cli/v2"
)

type KhedraApp struct {
	cli            *cli.App
	config         *types.Config
	logger         *types.CustomLogger
	controlSvc     *services.ControlService
	serviceManager *services.ServiceManager
}

// RestartAllServices restarts all services except the control service directly via service manager.
func (k *KhedraApp) RestartAllServices() error {
	if k.serviceManager == nil {
		return fmt.Errorf("service manager not initialized")
	}

	k.logger.Info("Restarting all services (except control) directly via service manager")

	// Get all services that can be restarted (this excludes control service automatically)
	results, err := k.serviceManager.Restart("all")
	if err != nil {
		k.logger.Error("Failed to restart services", "error", err)
		return err
	}

	for _, result := range results {
		serviceName := result["name"]
		status := result["status"]
		k.logger.Info("Service restart result", "service", serviceName, "status", status)
	}

	k.logger.Info("All restartable services restarted successfully")
	return nil
}

func NewKhedraApp() *KhedraApp {
	k := KhedraApp{}
	if k.isRunning() {
		logger.Panic(colors.BrightBlue + "khedra is already running - cannot run..." + colors.Off)
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
