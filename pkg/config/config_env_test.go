package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMultipleChainsEnvironmentOverrides(t *testing.T) {
	defer setTempEnvVar("TEST_MODE", "true")()
	defer setTempEnvVar("TB_KHEDRA_CHAINS_MAINNET_RPCS", "http://rpc1.mainnet,http://rpc2.mainnet")()
	defer setTempEnvVar("TB_KHEDRA_CHAINS_SEPOLIA_ENABLED", "true")()

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

	assert.Equal(t, []string{"http://rpc1.mainnet", "http://rpc2.mainnet"}, cfg.Chains["mainnet"].RPCs, "RPCs for mainnet should be overridden by environment variable")
	assert.True(t, cfg.Chains["sepolia"].Enabled, "Enabled flag for sepolia should be overridden by environment variable")
}

func TestEnvironmentVariableOverridesForServices(t *testing.T) {
	// Set environment variables to override the configuration
	defer setTempEnvVar("TB_KHEDRA_SERVICES_API_ENABLED", "false")()
	defer setTempEnvVar("TB_KHEDRA_SERVICES_API_PORT", "9090")()
	defer setTempEnvVar("TEST_MODE", "true")()

	// Create a temporary directory for the config file
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yaml")

	// Mock getConfigFn to return the temporary config path
	originalGetConfigFn := getConfigFn
	getConfigFn = func() string { return configFile }
	defer func() { getConfigFn = originalGetConfigFn }()

	// Establish the configuration file
	establishConfig(configFile)

	// Load the configuration
	cfg := MustLoadConfig(configFile)

	// Check if the API service exists in the configuration
	apiService, exists := cfg.Services["api"]
	assert.True(t, exists, "API service should exist in the configuration")

	// Validate that the overrides were applied correctly
	assert.False(t, apiService.Enabled, "Enabled flag for API service should be overridden by environment variable")
	assert.Equal(t, 9090, apiService.Port, "Port for API service should be overridden by environment variable")
}

func TestMultipleServicesEnvironmentOverrides(t *testing.T) {
	// Set environment variables to override the configuration
	defer setTempEnvVar("TB_KHEDRA_SERVICES_API_ENABLED", "false")()
	defer setTempEnvVar("TB_KHEDRA_SERVICES_SCRAPER_ENABLED", "true")()
	defer setTempEnvVar("TB_KHEDRA_SERVICES_SCRAPER_PORT", "8081")()
	defer setTempEnvVar("TEST_MODE", "true")()

	// Create a temporary directory for the config file
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yaml")

	// Mock getConfigFn to return the temporary config path
	originalGetConfigFn := getConfigFn
	getConfigFn = func() string { return configFile }
	defer func() { getConfigFn = originalGetConfigFn }()

	// Establish the configuration file
	establishConfig(configFile)

	// Load the configuration
	cfg := MustLoadConfig(configFile)

	// Check if the API and Scraper services exist in the configuration
	apiService, apiExists := cfg.Services["api"]
	scraperService, scraperExists := cfg.Services["scraper"]

	assert.True(t, apiExists, "API service should exist in the configuration")
	assert.True(t, scraperExists, "Scraper service should exist in the configuration")

	// Validate that the overrides were applied correctly
	assert.False(t, apiService.Enabled, "Enabled flag for API service should be overridden by environment variable")
	assert.True(t, scraperService.Enabled, "Enabled flag for Scraper service should be overridden by environment variable")
	assert.Equal(t, 8081, scraperService.Port, "Port for Scraper service should be overridden by environment variable")
}

func TestNoEnvironmentVariables(t *testing.T) {
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

	assert.NotEqual(t, cfg.Chains["mainnet"], nil, "mainnet chain should exist in the configuration")
	assert.Equal(t, []string{"http://localhost:8545"}, cfg.Chains["mainnet"].RPCs, "RPCs for mainnet should remain as default")
	assert.True(t, cfg.Chains["mainnet"].Enabled, "Enabled flag for mainnet should remain as default")
}

func TestEnvironmentVariableOverridesForChains(t *testing.T) {
	defer setTempEnvVar("TEST_MODE", "true")()
	defer setTempEnvVar("TB_KHEDRA_CHAINS_MAINNET_RPCS", "http://rpc1.mainnet,http://rpc2.mainnet")()
	defer setTempEnvVar("TB_KHEDRA_CHAINS_MAINNET_ENABLED", "false")()

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

	assert.Equal(t, []string{"http://rpc1.mainnet", "http://rpc2.mainnet"}, cfg.Chains["mainnet"].RPCs, "RPCs for mainnet should be overridden by environment variable")
	assert.False(t, cfg.Chains["mainnet"].Enabled, "Enabled flag for mainnet should be overridden by environment variable")
}

func TestInvalidBooleanValueForChains(t *testing.T) {
	defer setTempEnvVar("TEST_MODE", "true")()
	defer setTempEnvVar("TB_KHEDRA_CHAINS_MAINNET_ENABLED", "not_a_bool")()

	// Use a temporary directory to simulate missing config
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yaml")

	// Mock getConfigFn to return the temporary config path
	originalGetConfigFn := getConfigFn
	getConfigFn = func() string { return configFile }
	defer func() { getConfigFn = originalGetConfigFn }()

	// Establish the config file if it doesn't exist
	establishConfig(configFile)

	// Attempt to load the configuration and expect an error
	_, err := loadConfig()
	assert.Error(t, err, "loadConfig should return an error for invalid boolean value")
	assert.Contains(t, err.Error(), "cannot parse", "Error message should indicate the inability to parse the boolean value")
	assert.Contains(t, err.Error(), "chains[mainnet].enabled", "Error message should point to the problematic field")
}

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
