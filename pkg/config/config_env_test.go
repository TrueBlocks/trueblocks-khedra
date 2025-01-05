package config

import (
	"testing"

	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestChainMultipleEnvironmentOverrides(t *testing.T) {
	env := []string{
		"TB_KHEDRA_CHAINS_MAINNET_RPCS=http://rpc1.mainnet,http://rpc2.mainnet",
		"TB_KHEDRA_CHAINS_SEPOLIA_ENABLED=true",
	}

	var configFile string
	defer types.SetupTest2(t, env, &configFile)()

	cfg := MustLoadConfig(configFile)
	assert.Equal(t, []string{"http://rpc1.mainnet", "http://rpc2.mainnet"}, cfg.Chains["mainnet"].RPCs, "RPCs for mainnet should be overridden by environment variable")
	assert.True(t, cfg.Chains["sepolia"].Enabled, "Enabled flag for sepolia should be overridden by environment variable")
}

func TestChainEnvironmentVariableOverrides(t *testing.T) {
	env := []string{
		"TB_KHEDRA_CHAINS_MAINNET_RPCS=http://rpc1.mainnet,http://rpc2.mainnet",
		"TB_KHEDRA_CHAINS_MAINNET_ENABLED=false",
	}

	var configFile string
	defer types.SetupTest2(t, env, &configFile)()

	cfg := MustLoadConfig(configFile)
	assert.Equal(t, []string{"http://rpc1.mainnet", "http://rpc2.mainnet"}, cfg.Chains["mainnet"].RPCs, "RPCs for mainnet should be overridden by environment variable")
	assert.False(t, cfg.Chains["mainnet"].Enabled, "Enabled flag for mainnet should be overridden by environment variable")
}

func TestChainInvalidBooleanValue(t *testing.T) {
	defer types.SetTempEnv("TB_KHEDRA_CHAINS_MAINNET_ENABLED", "not_a_bool")()
	defer types.SetupTest(t, nil, types.GetConfigFn, types.EstablishConfig)()

	_, err := loadConfig()
	assert.Error(t, err, "loadConfig should return an error for invalid boolean value")
	assert.Contains(t, err.Error(), "cannot parse", "Error message should indicate the inability to parse the boolean value")
	assert.Contains(t, err.Error(), "chains[mainnet].enabled", "Error message should point to the problematic field")
}

func TestEnvironmentVariableOverridesForServices(t *testing.T) {
	env := []string{
		"TB_KHEDRA_SERVICES_API_ENABLED=false",
		"TB_KHEDRA_SERVICES_API_PORT=9090",
	}

	var configFile string
	defer types.SetupTest2(t, env, &configFile)()

	cfg := MustLoadConfig(configFile)
	apiService, exists := cfg.Services["api"]
	assert.True(t, exists, "API service should exist in the configuration")
	assert.False(t, apiService.Enabled, "Enabled flag for API service should be overridden by environment variable")
	assert.Equal(t, 9090, apiService.Port, "Port for API service should be overridden by environment variable")
}

func TestMultipleServicesEnvironmentOverrides(t *testing.T) {
	env := []string{
		"TB_KHEDRA_SERVICES_API_ENABLED=false",
		"TB_KHEDRA_SERVICES_SCRAPER_ENABLED=true",
		"TB_KHEDRA_SERVICES_SCRAPER_PORT=8081",
	}

	var configFile string
	defer types.SetupTest2(t, env, &configFile)()

	cfg := MustLoadConfig(configFile)
	apiService, apiExists := cfg.Services["api"]
	scraperService, scraperExists := cfg.Services["scraper"]
	assert.True(t, apiExists, "API service should exist in the configuration")
	assert.True(t, scraperExists, "Scraper service should exist in the configuration")
	assert.False(t, apiService.Enabled, "Enabled flag for API service should be overridden by environment variable")
	assert.True(t, scraperService.Enabled, "Enabled flag for Scraper service should be overridden by environment variable")
	assert.Equal(t, 8081, scraperService.Port, "Port for Scraper service should be overridden by environment variable")
}

// func TestNoEnvironmentVariables(t *testing.T) {
// 	var configFile string
// 	defer types.SetupTest2(t, []string{}, &configFile)() // types.GetConfigFn, types.EstablishConfig)()

// 	s := file.AsciiFileToString(configFile)
// 	fmt.Println(s)
// 	fmt.Println(configFile)

// 	cfg := MustLoadConfig(configFile)
// 	mainnet := cfg.Chains["mainnet"]
// 	assert.NotEqual(t, mainnet, nil, "mainnet chain should exist in the configuration")
// 	for _, chain := range cfg.Chains {
// 		fmt.Println(chain)
// 		fmt.Println(chain.RPCs)
// 	}
// 	// assert.Equal(t, []string{"http://localhost:8545"}, cfg.Chains["mainnet"].RPCs, "RPCs for mainnet should remain as default")
// 	// fmt.Println(mainnet)
// 	// fmt.Println(mainnet.Enabled)
// 	assert.True(t, mainnet.Enabled, "Enabled flag for mainnet should remain as default")
// }
