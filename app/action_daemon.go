package app

import (
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	_ "github.com/TrueBlocks/trueblocks-khedra/v2/pkg/env"
	"github.com/TrueBlocks/trueblocks-sdk/v4/services"
	"github.com/urfave/cli/v2"
)

func (k *KhedraApp) daemonAction(c *cli.Context) error {
	_ = c // linter
	if err := k.loadConfigIfInitialized(); err != nil {
		return err
	}
	k.logger.Info("Starting khedra daemon...config loaded...")

	for _, ch := range k.config.Chains {
		if ch.Enabled {
			if !ch.HasValidRpc(4) {
				return fmt.Errorf("chain %s has no valid RPC", ch.Name)
			}
			Progress("Connected to", "chain", ch.Name)
		}
	}
	k.logger.Info("Processing chains...", "chainList", k.config.EnabledChains())
	os.Exit(0)

	os.Setenv("XDG_CONFIG_HOME", k.config.General.DataFolder)
	os.Setenv("TB_SETTINGS_DEFAULTCHAIN", "mainnet")
	os.Setenv("TB_SETTINGS_INDEXPATH", k.config.IndexPath())
	os.Setenv("TB_SETTINGS_CACHEPATH", k.config.CachePath())
	for key, ch := range k.config.Chains {
		if ch.Enabled {
			envKey := "TB_CHAINS_" + strings.ToUpper(key) + "_RPCPROVIDER"
			os.Setenv(envKey, ch.RPCs[0])
		}
	}

	for _, env := range os.Environ() {
		if (strings.HasPrefix(env, "TB_") || strings.HasPrefix(env, "XDG_")) && strings.Contains(env, "=") {
			parts := strings.Split(env, "=")
			if len(parts) > 1 {
				k.logger.Info("environment", parts[0], parts[1])
			} else {
				k.logger.Info("environment", parts[0], "<empty>")
			}
		}
	}

	fmt.Println("Services...", k.config.ServiceList(true /* enabledOnly */))

	for key, ch := range k.config.Chains {
		chain := k.chainList.ChainsMap[ch.ChainId]
		fmt.Println(colors.Blue, "Chain", key, ch, chain, colors.Off)
		fmt.Println("Sleeping 1...")
		time.Sleep(3 * time.Second)
	}

	configFn := filepath.Join(k.config.General.DataFolder, "trueBlocks.toml")
	if file.FileExists(configFn) {
		fmt.Println("Config file found", configFn)
	} else {
		fmt.Println("Config file not found", configFn)
	}

	fmt.Println("Sleeping...")
	time.Sleep(2 * time.Second)
	os.Exit(0)

	var activeServices []services.Servicer
	chains := strings.Split(strings.ReplaceAll(k.config.EnabledChains(), " ", ""), ",")
	scraperSvc := services.NewScrapeService(
		k.logger.GetLogger(),
		"all",
		chains,
		k.config.Services["scraper"].Sleep,
		k.config.Services["scraper"].BatchSize,
	)
	monitorSvc := services.NewMonitorService(nil)
	apiSvc := services.NewApiService(k.logger.GetLogger())
	ipfsSvc := services.NewIpfsService(k.logger.GetLogger())
	controlService := services.NewControlService(k.logger.GetLogger())
	activeServices = append(activeServices, controlService)
	activeServices = append(activeServices, scraperSvc)
	activeServices = append(activeServices, monitorSvc)
	activeServices = append(activeServices, apiSvc)
	activeServices = append(activeServices, ipfsSvc)
	slog.Info("Starting khedra daemon", "services", len(activeServices))
	serviceManager := services.NewServiceManager(activeServices, k.logger.GetLogger())
	for _, svc := range activeServices {
		if controlSvc, ok := svc.(*services.ControlService); ok {
			controlSvc.AttachServiceManager(serviceManager)
		}
	}
	if err := serviceManager.StartAllServices(); err != nil {
		k.Fatal(err.Error())
	}
	serviceManager.HandleSignals()
	if true {
		select {}
	}

	return nil
}

// If trueBlocks.io file exists, check that it contains records for each enabled chain

// If trueBlocks.io does not exist create it with a template

/*
	configFn := filepath.Join(a.Config.ConfigPath, "trueBlocks.toml")
	if file.FileExists(configFn) {
		a.Logger.Info("config loaded", "configFile", configFn, "nChains", len(a.Config.ProviderMap))
		// check to make sure the config file has all the chains
		contents := file.AsciiFileToString(configFn)
		for chain := range a.Config.ProviderMap {
			search := "[chains." + chain + "]"
			if !strings.Contains(contents, search) {
				msg := fmt.Sprintf("config file {%s} does not contain {%s}", configFn, search)
				msg = colors.ColoredWith(msg, colors.Red)
				return errors.New(msg)
			}
		}
		return nil
	}

	if err := file.EstablishFolder(a.Config.ConfigPath); err != nil {
		return err
	}
	for _, chain := range chains {
		chainConfig := filepath.Join(a.Config.ConfigPath, "config", chain)
		if err := file.EstablishFolder(chainConfig); err != nil {
			return err
		}
	}

	tmpl, err := template.New("tmpl").Parse(configTmpl)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err = tmpl.Execute(&buf, &a.Config); err != nil {
		return err
	}

	_ = file.StringToAsciiFile(configFn, buf.String())
	a.Logger.Info("Created config file", "configFile", configFn, "nChains", len(a.Config.ProviderMap))

	return nil
}

var configTmpl string = `[version]
  current = "v4.0.0"

[settings]
  cachePath = "{{.CachePath}}"
  defaultChain = "mainnet"
  indexPath = "{{.IndexPath}}"

[keys]
  [keys.etherscan]
    apiKey = ""

[chains]{{.ChainDescriptors}}
`

*/
