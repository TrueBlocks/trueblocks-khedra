package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// func TestEnvironmentVariableOverridesForChains(t *testing.T) {
// 	defer setTempEnvVar("TEST_MODE", "true")()
// 	defer setTempEnvVar("TB_KHEDRA_CHAINS_MAINNET_RPCS", "http://rpc1.mainnet,http://rpc2.mainnet")()
// 	defer setTempEnvVar("TB_KHEDRA_CHAINS_MAINNET_ENABLED", "false")()

// 	// Use a temporary directory to simulate missing config
// 	tmpDir := t.TempDir()
// 	configFile := filepath.Join(tmpDir, "config.yaml")

// 	// Mock getConfigFn to return the temporary config path
// 	originalGetConfigFn := getConfigFn
// 	getConfigFn = func() string { return configFile }
// 	defer func() { getConfigFn = originalGetConfigFn }()

// 	// Establish the config file if it doesn't exist
// 	establishConfig(configFile)

// 	// Load the configuration
// 	cfg := MustLoadConfig(configFile)

// 	// Validate the overrides
// 	mainnetIndex := -1
// 	for i, chain := range cfg.Chains {
// 		if chain.Name == "mainnet" {
// 			mainnetIndex = i
// 			break
// 		}
// 	}

// 	assert.NotEqual(t, -1, mainnetIndex, "mainnet chain should exist in the configuration")
// 	assert.Equal(t, []string{"http://rpc1.mainnet", "http://rpc2.mainnet"}, cfg.Chains[mainnetIndex].RPCs, "RPCs for mainnet should be overridden by environment variable")
// 	assert.False(t, cfg.Chains[mainnetIndex].Enabled, "Enabled flag for mainnet should be overridden by environment variable")
// }

// func TestInvalidBooleanValueForChains(t *testing.T) {
// 	defer setTempEnvVar("TEST_MODE", "true")()
// 	defer setTempEnvVar("TB_KHEDRA_CHAINS_MAINNET_ENABLED", "not_a_bool")()

// 	// Use a temporary directory to simulate missing config
// 	tmpDir := t.TempDir()
// 	configFile := filepath.Join(tmpDir, "config.yaml")

// 	// Mock getConfigFn to return the temporary config path
// 	originalGetConfigFn := getConfigFn
// 	getConfigFn = func() string { return configFile }
// 	defer func() { getConfigFn = originalGetConfigFn }()

// 	// Establish the config file if it doesn't exist
// 	establishConfig(configFile)

// 	// Attempt to load the configuration and expect an error
// 	_, err := loadConfig()
// 	assert.Error(t, err, "loadConfig should return an error for invalid boolean value")
// 	assert.Contains(t, err.Error(), "invalid boolean value", "Error message should indicate invalid boolean")
// }

func TestMissingEnvironmentVariables(t *testing.T) {
	defer setTempEnvVar("TEST_MODE", "true")()

	// Use a temporary directory to simulate missing config
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yaml")

	// Mock getConfigFn to return the temporary config path
	originalGetConfigFn := getConfigFn
	getConfigFn = func() string { return configFile }
	defer func() { getConfigFn = originalGetConfigFn }()

	// Establish the config file if it doesn't exist
	establishConfig(configFile)

	// Load the configuration
	cfg := MustLoadConfig(configFile)

	// Validate that default values are used when no environment variables are set
	mainnetIndex := -1
	for i, chain := range cfg.Chains {
		if chain.Name == "mainnet" {
			mainnetIndex = i
			break
		}
	}

	assert.NotEqual(t, -1, mainnetIndex, "mainnet chain should exist in the configuration")
	assert.Equal(t, []string{"http://localhost:8545"}, cfg.Chains[mainnetIndex].RPCs, "RPCs for mainnet should remain as default")
	assert.True(t, cfg.Chains[mainnetIndex].Enabled, "Enabled flag for mainnet should remain as default")
}

// func TestMultipleChainsEnvironmentOverrides(t *testing.T) {
// 	defer setTempEnvVar("TEST_MODE", "true")()
// 	defer setTempEnvVar("TB_KHEDRA_CHAINS_MAINNET_RPCS", "http://rpc1.mainnet,http://rpc2.mainnet")()
// 	defer setTempEnvVar("TB_KHEDRA_CHAINS_SEPOLIA_ENABLED", "true")()

// 	// Use a temporary directory to simulate missing config
// 	tmpDir := t.TempDir()
// 	configFile := filepath.Join(tmpDir, "config.yaml")

// 	// Mock getConfigFn to return the temporary config path
// 	originalGetConfigFn := getConfigFn
// 	getConfigFn = func() string { return configFile }
// 	defer func() { getConfigFn = originalGetConfigFn }()

// 	// Establish the config file if it doesn't exist
// 	establishConfig(configFile)

// 	// Load the configuration
// 	cfg := MustLoadConfig(configFile)

// 	// Validate overrides for mainnet
// 	mainnetIndex := -1
// 	sepoliaIndex := -1
// 	for i, chain := range cfg.Chains {
// 		fmt.Println("Looking for chain:", chain)
// 		if chain.Name == "mainnet" {
// 			fmt.Println("found mainnet")
// 			mainnetIndex = i
// 		} else if chain.Name == "sepolia" {
// 			fmt.Println("found sepolia")
// 			sepoliaIndex = i
// 		}
// 	}

// 	assert.NotEqual(t, -1, mainnetIndex, "mainnet chain should exist in the configuration")
// 	assert.Equal(t, []string{"http://rpc1.mainnet", "http://rpc2.mainnet"}, cfg.Chains[mainnetIndex].RPCs, "RPCs for mainnet should be overridden by environment variable")

// 	assert.NotEqual(t, -1, sepoliaIndex, "sepolia chain should exist in the configuration")
// 	assert.True(t, cfg.Chains[sepoliaIndex].Enabled, "Enabled flag for sepolia should be overridden by environment variable")
// }

// func TestEnvironmentVariableOverridesForServices(t *testing.T) {
// 	defer setTempEnvVar("TB_KHEDRA_SERVICES_API_ENABLED", "false")()
// 	defer setTempEnvVar("TB_KHEDRA_SERVICES_API_PORT", "9090")()
// 	defer setTempEnvVar("TEST_MODE", "true")()

// 	tmpDir := t.TempDir()
// 	configFile := filepath.Join(tmpDir, "config.yaml")

// 	originalGetConfigFn := getConfigFn
// 	getConfigFn = func() string { return configFile }
// 	defer func() { getConfigFn = originalGetConfigFn }()

// 	establishConfig(configFile)

// 	cfg := MustLoadConfig(configFile)

// 	apiIndex := -1
// 	for i, service := range cfg.Services {
// 		if service.Name == "api" {
// 			apiIndex = i
// 			break
// 		}
// 	}

// 	assert.NotEqual(t, -1, apiIndex, "API service should exist in the configuration")
// 	assert.False(t, cfg.Services[apiIndex].Enabled, "Enabled flag for API service should be overridden by environment variable")
// 	assert.Equal(t, 9090, cfg.Services[apiIndex].Port, "Port for API service should be overridden by environment variable")
// }

// func TestMultipleServicesEnvironmentOverrides(t *testing.T) {
// 	defer setTempEnvVar("TB_KHEDRA_SERVICES_API_ENABLED", "false")()
// 	defer setTempEnvVar("TB_KHEDRA_SERVICES_SCRAPER_ENABLED", "true")()
// 	defer setTempEnvVar("TB_KHEDRA_SERVICES_SCRAPER_PORT", "8081")()
// 	defer setTempEnvVar("TEST_MODE", "true")()

// 	tmpDir := t.TempDir()
// 	configFile := filepath.Join(tmpDir, "config.yaml")

// 	originalGetConfigFn := getConfigFn
// 	getConfigFn = func() string { return configFile }
// 	defer func() { getConfigFn = originalGetConfigFn }()

// 	establishConfig(configFile)

// 	cfg := MustLoadConfig(configFile)

// 	apiIndex := -1
// 	scraperIndex := -1
// 	for i, service := range cfg.Services {
// 		if service.Name == "api" {
// 			apiIndex = i
// 		}
// 		if service.Name == "scraper" {
// 			scraperIndex = i
// 		}
// 	}

// 	assert.NotEqual(t, -1, apiIndex, "API service should exist in the configuration")
// 	assert.NotEqual(t, -1, scraperIndex, "Scraper service should exist in the configuration")

// 	assert.False(t, cfg.Services[apiIndex].Enabled, "Enabled flag for API service should be overridden by environment variable")
// 	assert.True(t, cfg.Services[scraperIndex].Enabled, "Enabled flag for Scraper service should be overridden by environment variable")
// 	assert.Equal(t, 8081, cfg.Services[scraperIndex].Port, "Port for Scraper service should be overridden by environment variable")
// }

func setTempEnvVar(key, value string) func() {
	originalValue, exists := os.LookupEnv(key)
	os.Setenv(key, value)
	return func() {
		if exists {
			os.Setenv(key, originalValue)
		} else {
			os.Unsetenv(key)
		}
	}
}
