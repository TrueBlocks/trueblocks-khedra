package app

import (
	"bytes"
	"fmt"
	"html/template"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

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
			k.logger.Progress("Connected to", "chain", ch.Name)
		}
	}
	k.logger.Info("Processing chains...", "chainList", k.config.EnabledChains())

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
				k.logger.Progress("environment", parts[0], parts[1])
			} else {
				k.logger.Progress("environment", parts[0], "<empty>")
			}
		}
	}

	k.logger.Progress("Starting services", "services", k.config.ServiceList(true /* enabledOnly */))

	configFn := filepath.Join(k.config.General.DataFolder, "trueBlocks.toml")
	if file.FileExists(configFn) {
		k.logger.Info("Config file found", "fn", configFn)
		if !k.chainsConfigured(configFn) {
			k.logger.Error("Config file not configured", "fn", configFn)
			return fmt.Errorf("config file not configured")
		}
	} else {
		k.logger.Warn("Config file not found", "fn", configFn)
		if err := k.createChifraConfig(); err != nil {
			k.logger.Error("Error creating config file", "error", err)
			return err
		}
	}

	var activeServices []services.Servicer
	controlService := services.NewControlService(k.logger.GetLogger())
	activeServices = append(activeServices, controlService)
	for _, svc := range k.config.Services {
		switch svc.Name {
		case "scraper":
			chains := strings.Split(strings.ReplaceAll(k.config.EnabledChains(), " ", ""), ",")
			scraperSvc := services.NewScrapeService(
				k.logger.GetLogger(),
				"all",
				chains,
				k.config.Services["scraper"].Sleep,
				k.config.Services["scraper"].BatchSize,
			)
			activeServices = append(activeServices, scraperSvc)
			if !svc.Enabled {
				scraperSvc.Pause()
			}
		case "monitor":
			monitorSvc := services.NewMonitorService(nil)
			activeServices = append(activeServices, monitorSvc)
			if !svc.Enabled {
				monitorSvc.Pause()
			}
		case "api":
			if svc.Enabled {
				apiSvc := services.NewApiService(k.logger.GetLogger())
				activeServices = append(activeServices, apiSvc)
			}
		case "ipfs":
			if svc.Enabled {
				ipfsSvc := services.NewIpfsService(k.logger.GetLogger())
				activeServices = append(activeServices, ipfsSvc)
			}
		}
	}

	slog.Info("Starting khedra daemon", "services", len(activeServices))
	serviceManager := services.NewServiceManager(activeServices, k.logger.GetLogger())
	for _, svc := range activeServices {
		if controlSvc, ok := svc.(*services.ControlService); ok {
			controlSvc.AttachServiceManager(serviceManager)
		}
	}
	if err := serviceManager.StartAllServices(); err != nil {
		k.logger.Fatal(err.Error())
	}

	serviceManager.HandleSignals()

	if true {
		select {}
	}

	return nil
}

func (k *KhedraApp) chainsConfigured(configFn string) bool {
	chainStr := k.config.EnabledChains()
	chains := strings.Split(chainStr, ",")

	k.logger.Info("chifra config loaded")
	k.logger.Info("checking", "configFile", configFn, "nChains", len(chains))

	contents := file.AsciiFileToString(configFn)
	for _, chain := range chains {
		search := "[chains." + chain + "]"
		if !strings.Contains(contents, search) {
			msg := fmt.Sprintf("config file {%s} does not contain {%s}", configFn, search)
			k.logger.Error(msg)
			return false
		}
	}
	return true
}

func (k *KhedraApp) createChifraConfig() error {
	if err := file.EstablishFolder(k.config.General.DataFolder); err != nil {
		return err
	}

	chainStr := k.config.EnabledChains()

	chains := strings.Split(chainStr, ",")
	for _, chain := range chains {
		if err := k.createChainConfig(chain); err != nil {
			return err
		}
	}

	tmpl, err := template.New("tmpl").Parse(configTmpl)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err = tmpl.Execute(&buf, &k.config); err != nil {
		return err
	}
	if len(buf.String()) == 0 {
		return fmt.Errorf("empty config file")
	}

	configFn := filepath.Join(k.config.General.DataFolder, "trueBlocks.toml")
	err = file.StringToAsciiFile(configFn, buf.String())
	if err != nil {
		return err
	}
	k.logger.Info("Created config file", "configFile", configFn, "nChains", len(chains))
	return nil
}

func (k *KhedraApp) createChainConfig(chain string) error {
	chainConfig := filepath.Join(k.config.General.DataFolder, "config", chain)
	if err := file.EstablishFolder(chainConfig); err != nil {
		return fmt.Errorf("failed to create folder %s: %w", chainConfig, err)
	}

	k.logger.Progress("Creating chain config", "chainConfig", chainConfig)

	// baseURL := "https://raw.githubusercontent.com/TrueBlocks/trueblocks-core/refs/heads/master/src/other/install/per-chain/"
	// url, err := url.JoinPath(baseURL, chain, "allocs.csv")
	// if err != nil {
	// 	return err
	// }
	// allocFn := filepath.Join(chainConfig, "allocs.csv")
	// dur := 100 * 365 * 24 * time.Hour // 100 years
	// if _, err := utils.DownloadAndStore(url, allocFn, dur); err != nil {
	// 	return fmt.Errorf("failed to download and store allocs.csv for chain %s: %w", chain, err)
	// }

	return nil
}

var configTmpl string = `[version]
  current = "{{.Version}}"

[settings]
  cachePath = "{{.CachePath}}"
  defaultChain = "mainnet"
  indexPath = "{{.IndexPath}}"

[keys]
  [keys.etherscan]
    apiKey = ""

[chains]
{{- range .Chains}}
  [chains.{{.Name}}]
    chain = "{{.Name}}"
    chainId = "{{.ChainID}}"
    remoteExplorer = "{{.RemoteExplorer}}"
    rpcProvider = "{{ index .RPCs 0 }}"
    symbol = "{{.Symbol}}"
{{end -}}
`
