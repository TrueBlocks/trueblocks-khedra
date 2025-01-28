package app

import (
	"fmt"
	"log/slog"
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
			slog.Info("Connected to", "chain", chain.Name, "rpc", rpc)
		}
	}

	var activeServices []services.Servicer
	chains := strings.Split(strings.ReplaceAll(k.config.ChainList(false /* enabledOnly */), " ", ""), ",")
	scraperSvc := services.NewScrapeService(
		k.logger,
		"all",
		chains,
		k.config.Services["scraper"].Sleep,
		k.config.Services["scraper"].BatchSize,
	)
	monitorSvc := services.NewMonitorService(nil)
	apiSvc := services.NewApiService(k.logger)
	ipfsSvc := services.NewIpfsService(k.logger)
	controlService := services.NewControlService(k.logger)
	activeServices = append(activeServices, controlService)
	activeServices = append(activeServices, scraperSvc)
	activeServices = append(activeServices, monitorSvc)
	activeServices = append(activeServices, apiSvc)
	activeServices = append(activeServices, ipfsSvc)
	slog.Info("Starting khedra daemon", "services", len(activeServices))
	serviceManager := services.NewServiceManager(activeServices, k.logger)
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

/*
SHOW ALL THE TB_KHEDRA_ VARIABLES FOUND
SHOW THE DATA FOLDER
THERE USED TO BE A DIFFERENCE BETWEEN THE INDEXED CHAINS AND THE CHAINS REQUIRING RPCS
TELL THE USER WHICH CHAINS ARE BEING PROCESSED
TELL THE USER WHICH SERVICES ARE BEING STARTED

// EstablishConfig either reads an existing configuration file or creates it if it doesn't exist.
func (a *App) EstablishConfig() error {
	var ok bool
	var err error
	if a.Config.ConfigPath, ok = os.LookupEnv("TB_NODE_DATADIR"); !ok {
		return errors.New("environment variable `TB_NODE_DATADIR` is required but not found")
	} else {
		if a.Config.ConfigPath, err = cleanDataPath(a.Config.ConfigPath); err != nil {
			return err
		}
	}
	a.Logger.Info("data directory", "dataDir", a.Config.ConfigPath)

	var targets string
	chainStr, ok := os.LookupEnv("TB_NODE_CHAINS")
	if !ok {
		chainStr, targets = "mainnet", "mainnet"
	} else {
		if chainStr, targets, err = cleanChainString(chainStr); err != nil {
			return err
		}
	}
	a.Logger.Info("configured chains", "chainStr", chainStr, "targets", targets)
	a.Config.Targets = strings.Split(targets, ",")

	chains := strings.Split(chainStr, ",")
	for _, chain := range chains {
		key := "TB_NODE_" + strings.ToUpper(chain) + "RPC"
		if providerUrl, ok := os.LookupEnv(key); !ok {
			msg := fmt.Sprintf("environment variable `%s` is required but not found (implied by TB_NODE_CHAINS=%s)", key, chainStr)
			return errors.New(msg)
		} else {
			providerUrl = strings.Trim(providerUrl, "/")
			if !isValidURL(providerUrl) {
				return fmt.Errorf("invalid URL for %s: %s", key, providerUrl)
			}
			a.Config.ProviderMap[chain] = providerUrl
		}
	}

	// // Set the environment trueblocks-core needs
	os.Setenv("XDG_CONFIG_HOME", a.Config.ConfigPath)
	os.Setenv("TB_SETTINGS_DEFAULTCHAIN", "mainnet")
	os.Setenv("TB_SETTINGS_INDEXPATH", a.Config.IndexPath())
	os.Setenv("TB_SETTINGS_CACHEPATH", a.Config.CachePath())
	for chain, providerUrl := range a.Config.ProviderMap {
		envKey := "TB_CHAINS_" + strings.ToUpper(chain) + "_RPCPROVIDER"
		os.Setenv(envKey, providerUrl)
	}

	for _, env := range os.Environ() {
		if (strings.HasPrefix(env, "TB_") || strings.HasPrefix(env, "XDG_")) && strings.Contains(env, "=") {
			parts := strings.Split(env, "=")
			if len(parts) > 1 {
				a.Logger.Info("environment", parts[0], parts[1])
			} else {
				a.Logger.Info("environment", parts[0], "<empty>")
			}
		}
	}

	for _, chain := range chains {
		providerUrl := a.Config.ProviderMap[chain]
		if err := a.tryConnect(chain, providerUrl, 5); err != nil {
			return err
		} else {
			a.Logger.Info("test connection", "result", "okay", "chain", chain, "providerUrl", providerUrl)
		}
	}

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



	handleService := func(i int, feature Feature) (int, error) {
		if hasValue(i) {
			if mode, err := validateOnOff(os.Args[i+1]); err == nil {
				switch feature {
				case Scrape:
					a.Scrape = mode
					if a.IsOn(Scrape) {
						scrapeSvc := services.NewScrapeService(
							a.Logger,
							string(a.InitMode),
							a.Config.Targets,
							a.Sleep,
							a.BlockCnt,
						)
						activeServices = append(activeServices, scrapeSvc)
					}
				case Api:
					a.Api = mode
					if a.IsOn(Api) {
						apiSvc := services.NewApiService(a.Logger)
						activeServices = append(activeServices, apiSvc)
					}
				case Ipfs:
					a.Ipfs = mode
					if a.IsOn(Ipfs) {
						ipfsSvc := services.NewIpfsService(a.Logger)
						activeServices = append(activeServices, ipfsSvc)
					}
				case Monitor:
					a.Monitor = mode
					if a.IsOn(Monitor) {
						monSvc := services.NewMonitorService(a.Logger)
						activeServices = append(activeServices, monSvc)
					}
				}
				return i + 1, nil
			} else {
				return i, fmt.Errorf("parsing --%s: %w", feature.String(), err)
			}
		}
		return i, fmt.Errorf("%w for --%s", ErrMissingArgument, feature.String())
	}
	controlService := services.NewControlService(a.Logger)
	activeServices = append([]services.Servicer{controlService}, activeServices...)
}




FOR RUNNING CORE
	os.Setenv("XDG_CONFIG_HOME", a.Config.ConfigPath)
	os.Setenv("TB_ SETTINGS_DEFAULTCHAIN", "mainnet")
	os.Setenv("TB_ SETTINGS_INDEXPATH", a.Config.IndexPath())
	os.Setenv("TB_ SETTINGS_CACHEPATH", a.Config.CachePath())
	for chain, providerUrl := range a.Config.ProviderMap {
		envKey := "TB_CHAINS_" + strings.ToUpper(chain) + "_RPCPROVIDER"
		os.Setenv(envKey, providerUrl)
	}
	for _, env := range os.Environ() {
		if (strings.HasPrefix(env, "TB_") || strings.HasPrefix(env, "XDG_")) && strings.Contains(env, "=") {
			parts := strings.Split(env, "=")
			if len(parts) > 1 {
				a.Logger.Info("environment", parts[0], parts[1])
			} else {
				a.Logger.Info("environment", parts[0], "<empty>")
			}
		}
	}
	for _, chain := range chains {
		providerUrl := a.Config.ProviderMap[chain]
		if err := a.tryConnect(chain, providerUrl, 5); err != nil {
			return err
		} else {
			a.Logger.Info("test connection", "result", "okay", "chain", chain, "providerUrl", providerUrl)
		}
	}
USED TO CHECK THAT IF THE USER SPECIFIED A CHAIN IN THE ENV, THEN IT HAD TO EXISTING IN TRUEBLOCKS.TOML
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
USED TO ESTABLISH FOLDERS
	if err := file.Establish Folder(a.Config.ConfigPath); err != nil {
		return err
	}
	for _, chain := range chains {
		chainConfig := filepath.Join(a.Config.ConfigPath, "config", chain)
		if err := file.Establish Folder(chainConfig); err != nil {
			return err
		}
	}
WOULD CREATE A MINIMAL TRUEBLOCKS.TOML IF NOT FOUND

*/
