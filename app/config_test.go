package app

import (
	"fmt"
	"strconv"
	"testing"

	coreFile "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/types"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestChainEnvOverrides(t *testing.T) {
	defer types.SetupTest([]string{
		"TB_KHEDRA_CHAINS_MAINNET_RPCS=http://rpc1.mainnet,http://rpc2.mainnet",
		"TB_KHEDRA_CHAINS_MAINNET_ENABLED=false",
		"TB_KHEDRA_CHAINS_SEPOLIA_ENABLED=true",
	})()

	if cfg, _, err := LoadConfig(); err != nil {
		t.Error(err)
	} else {
		assert.Equal(t, []string{"http://rpc1.mainnet", "http://rpc2.mainnet"}, cfg.Chains["mainnet"].RPCs, "RPCs for mainnet should be overridden by environment variable")
		assert.False(t, cfg.Chains["mainnet"].Enabled, "Enabled flag for mainnet should be overridden by environment variable")
		assert.True(t, cfg.Chains["sepolia"].Enabled, "Enabled flag for sepolia should be overridden by environment variable")
	}
}

func TestChainInvalidBooleanValue(t *testing.T) {
	defer types.SetupTest([]string{
		"TB_KHEDRA_CHAINS_MAINNET_ENABLED=not_a_bool",
	})()

	if cfg, _, err := LoadConfig(); err != nil {
		assert.Error(t, err, "loadConfig should return an error for invalid boolean value")
		assert.Contains(t, err.Error(), "cannot parse", "Error message should indicate the inability to parse the boolean value")
		assert.Contains(t, err.Error(), "chains[mainnet].enabled", "Error message should point to the problematic field")
	} else {
		t.Error("loadConfig should return an error for invalid boolean value", cfg.Chains["mainnet"].Enabled)
	}
}

func TestServiceEnvironmentVariableOverrides(t *testing.T) {
	defer types.SetupTest([]string{
		"TB_KHEDRA_SERVICES_API_ENABLED=false",
		"TB_KHEDRA_SERVICES_API_PORT=9090",
	})()

	if cfg, _, err := LoadConfig(); err != nil {
		t.Error(err)
	} else {
		apiService, exists := cfg.Services["api"]
		assert.True(t, exists, "API service should exist in the configuration")
		assert.False(t, apiService.Enabled, "Enabled flag for API service should be overridden by environment variable")
		assert.Equal(t, 9090, apiService.Port, "Port for API service should be overridden by environment variable")
	}
}

func TestServiceMultipleEnvironmentOverrides(t *testing.T) {
	defer types.SetupTest([]string{
		"TB_KHEDRA_SERVICES_API_ENABLED=false",
		"TB_KHEDRA_SERVICES_SCRAPER_ENABLED=true",
		"TB_KHEDRA_SERVICES_SCRAPER_PORT=8081",
	})()

	if cfg, _, err := LoadConfig(); err != nil {
		t.Error(err)
	} else {
		apiService, apiExists := cfg.Services["api"]
		scraperService, scraperExists := cfg.Services["scraper"]
		assert.True(t, apiExists, "API service should exist in the configuration")
		assert.True(t, scraperExists, "Scraper service should exist in the configuration")
		assert.False(t, apiService.Enabled, "Enabled flag for API service should be overridden by environment variable")
		assert.True(t, scraperService.Enabled, "Enabled flag for Scraper service should be overridden by environment variable")
		assert.Equal(t, 8081, scraperService.Port, "Port for Scraper service should be overridden by environment variable")
	}
}

func TestEnvNoVariables(t *testing.T) {
	defer types.SetupTest([]string{})()

	if cfg, _, err := LoadConfig(); err != nil {
		t.Error(err)
	} else {
		mainnet := cfg.Chains["mainnet"]
		assert.NotEqual(t, mainnet, nil, "mainnet chain should exist in the configuration")
		assert.Equal(t, []string{"http://localhost:8545"}, cfg.Chains["mainnet"].RPCs, "RPCs for mainnet should remain as default")
		assert.True(t, mainnet.Enabled, "Enabled flag for mainnet should remain as default")
	}
}

func TestServiceInvalidPort(t *testing.T) {
	defer types.SetupTest([]string{
		"TB_KHEDRA_SERVICES_API_PORT=invalid_port",
	})()

	if cfg, _, err := LoadConfig(); err != nil {
		assert.Error(t, err, "loadConfig should return an error for invalid port value")
		assert.Contains(t, err.Error(), "invalid_port", "Error message should indicate invalid port")
	} else {
		t.Error("loadConfig should return an error for invalid port", cfg.Services["api"])
	}
}

func TestChainLargeNumberOfChains(t *testing.T) {
	defer types.SetupTest([]string{})()

	nChains := 1000
	cfg := types.NewConfig()
	cfg.Chains = make(map[string]types.Chain)
	for i := 0; i < nChains; i++ {
		chainName := "chain" + strconv.Itoa(i)
		cfg.Chains[chainName] = types.Chain{
			Name:    chainName,
			RPCs:    []string{fmt.Sprintf("http://%s.rpc", chainName)},
			Enabled: true,
		}
	}

	bytes, _ := yaml.Marshal(cfg)
	coreFile.StringToAsciiFile(types.GetConfigFn(), string(bytes))

	var err error
	if cfg, _, err = LoadConfig(); err != nil {
		t.Error(err)
	} else {
		assert.Equal(t, nChains+2, len(cfg.Chains), "All chains should be loaded correctly")
	}
}

func TestChainMissingInConfig(t *testing.T) {
	defer types.SetupTest([]string{
		"TB_KHEDRA_CHAINS_UNKNOWN_NAME=unknown",
		"TB_KHEDRA_CHAINS_UNKNOWN_RPCS=http://unknown.rpc",
		"TB_KHEDRA_CHAINS_UNKNOWN_ENABLED=true",
	})()

	if cfg, _, err := LoadConfig(); err != nil {
		assert.Error(t, err, "An error should occur if an unknown chain is defined in the environment but not in the configuration file")
	} else {
		t.Error("loadConfig should return an error for invalid chain", cfg.Chains["unknown"])
	}
}

func TestChainEmptyRPCs(t *testing.T) {
	defer types.SetupTest([]string{
		"TB_KHEDRA_CHAINS_MAINNET_RPCS=",
	})()

	if cfg, _, err := LoadConfig(); err != nil {
		t.Error(err)
	} else {
		assert.NotEmpty(t, cfg.Chains["mainnet"].RPCs, "Mainnet RPCs should not be empty in the final configuration")
	}
}

func TestConfigMustLoad(t *testing.T) {
	defer types.SetupTest([]string{})()
	assert.FileExists(t, types.GetConfigFn())
}

func TestConfigMustLoadDefaults(t *testing.T) {
	defer types.SetupTest([]string{})()

	if cfg, _, err := LoadConfig(); err != nil {
		t.Error(err)
	} else {
		for name, service := range cfg.Services {
			switch name {
			case "scraper", "monitor":
				assert.GreaterOrEqual(t, service.BatchSize, 50)
				assert.LessOrEqual(t, service.BatchSize, 10000)
				assert.GreaterOrEqual(t, service.Sleep, 0)
			case "api", "ipfs":
				assert.GreaterOrEqual(t, service.Port, 1024)
				assert.LessOrEqual(t, service.Port, 65535)
			}
		}
	}
}
