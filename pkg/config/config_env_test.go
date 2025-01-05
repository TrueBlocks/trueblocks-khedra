package config

import (
	"testing"

	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestChainEnvOverrides(t *testing.T) {
	defer types.SetupTest([]string{
		"TB_KHEDRA_CHAINS_MAINNET_RPCS=http://rpc1.mainnet,http://rpc2.mainnet",
		"TB_KHEDRA_CHAINS_MAINNET_ENABLED=false",
		"TB_KHEDRA_CHAINS_SEPOLIA_ENABLED=true",
	})()

	if cfg, err := LoadConfig(); err != nil {
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

	if cfg, err := LoadConfig(); err != nil {
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

	if cfg, err := LoadConfig(); err != nil {
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

	if cfg, err := LoadConfig(); err != nil {
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

	if cfg, err := LoadConfig(); err != nil {
		t.Error(err)
	} else {
		mainnet := cfg.Chains["mainnet"]
		assert.NotEqual(t, mainnet, nil, "mainnet chain should exist in the configuration")
		assert.Equal(t, []string{"http://localhost:8545"}, cfg.Chains["mainnet"].RPCs, "RPCs for mainnet should remain as default")
		assert.True(t, mainnet.Enabled, "Enabled flag for mainnet should remain as default")
	}
}
