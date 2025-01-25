package app

import (
	"fmt"
	"strings"

	coreFile "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	_ "github.com/TrueBlocks/trueblocks-khedra/v2/pkg/env"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/types"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/validate"
	"github.com/TrueBlocks/trueblocks-sdk/v4/services"
	"github.com/urfave/cli/v2"
)

func (k *KhedraApp) daemonAction(c *cli.Context) error {
	_ = c // linter
	fn := types.GetConfigFnNoCreate()
	if !coreFile.FileExists(fn) {
		return fmt.Errorf("not initialized you must run `khedra init` first")
	}

	_, err := k.ConfigMaker()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	for _, chain := range k.config.Chains {
		for _, rpc := range chain.RPCs {
			if err := validate.TryConnect(chain.Name, rpc, 5); err != nil {
				return err
			}
			k.Info("Connected to", "chain", chain.Name, "rpc", rpc)
		}
	}

	var activeServices []services.Servicer
	chains := strings.Split(strings.ReplaceAll(k.config.ChainList(), " ", ""), ",")
	scraperSvc := services.NewScrapeService(
		k.progLogger,
		"all",
		chains,
		k.config.Services["scraper"].Sleep,
		k.config.Services["scraper"].BatchSize,
	)
	monitorSvc := services.NewMonitorService(nil)
	apiSvc := services.NewApiService(k.progLogger)
	ipfsSvc := services.NewIpfsService(k.progLogger)
	controlService := services.NewControlService(k.progLogger)
	activeServices = append(activeServices, controlService)
	activeServices = append(activeServices, scraperSvc)
	activeServices = append(activeServices, monitorSvc)
	activeServices = append(activeServices, apiSvc)
	activeServices = append(activeServices, ipfsSvc)
	k.Info("Starting khedra daemon", "services", len(activeServices))
	serviceManager := services.NewServiceManager(activeServices, k.progLogger)
	for _, svc := range activeServices {
		if controlSvc, ok := svc.(*services.ControlService); ok {
			controlSvc.AttachServiceManager(serviceManager)
		}
	}
	if err := serviceManager.StartAllServices(); err != nil {
		k.Fatal(err.Error())
	}
	serviceManager.HandleSignals()
	select {}

	return nil
}
