package app

import (
	"fmt"

	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/types"
)

func LoadConfig() (types.Config, error) {
	cfg, err := loadFileConfig()
	if err != nil {
		return types.Config{}, fmt.Errorf("failed to load file configuration: %w", err)
	}
	keys := types.GetEnvironmentKeys(cfg, types.InEnv)
	if err := types.ApplyEnv(keys, &cfg); err != nil {
		return types.Config{}, fmt.Errorf("failed to apply environment configuration: %w", err)
	}

	if err := validateConfig(cfg); err != nil {
		return types.Config{}, fmt.Errorf("validation error: %w", err)
	}

	if err := initializeFolders(cfg); err != nil {
		return types.Config{}, fmt.Errorf("failed to initialize folders: %w", err)
	}

	return cfg, nil
}

/*
SHOW ALL THE TB_KHEDRA_ VARIABLES FOUND
SHOW THE DATA FOLDER
THERE USED TO BE A DIFFERENCE BETWEEN THE INDEXED CHAINS AND THE CHAINS REQUIRING RPCS
TELL THE USER WHICH CHAINS ARE BEING PROCESSED
TELL THE USER WHICH SERVICES ARE BEING STARTED
RUN INIT EVERY TIME WE START?
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
