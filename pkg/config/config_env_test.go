package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestMultipleChainsEnvironmentOverrides(t *testing.T) {
	var configFile string

	defer setTempEnvVar("TB_KHEDRA_CHAINS_MAINNET_RPCS", "http://rpc1.mainnet,http://rpc2.mainnet")()
	defer setTempEnvVar("TB_KHEDRA_CHAINS_SEPOLIA_ENABLED", "true")()
	defer setupTest(t, &configFile)()

	// Load the configuration
	cfg := MustLoadConfig(configFile)

	assert.Equal(t, []string{"http://rpc1.mainnet", "http://rpc2.mainnet"}, cfg.Chains["mainnet"].RPCs, "RPCs for mainnet should be overridden by environment variable")
	assert.True(t, cfg.Chains["sepolia"].Enabled, "Enabled flag for sepolia should be overridden by environment variable")
}

func TestEnvironmentVariableOverridesForServices(t *testing.T) {
	var configFile string

	defer setTempEnvVar("TB_KHEDRA_SERVICES_API_ENABLED", "false")()
	defer setTempEnvVar("TB_KHEDRA_SERVICES_API_PORT", "9090")()
	defer setTempEnvVar("TEST_MODE", "true")()
	defer setupTest(t, &configFile)()

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
	var configFile string

	defer setTempEnvVar("TB_KHEDRA_SERVICES_API_ENABLED", "false")()
	defer setTempEnvVar("TB_KHEDRA_SERVICES_SCRAPER_ENABLED", "true")()
	defer setTempEnvVar("TB_KHEDRA_SERVICES_SCRAPER_PORT", "8081")()
	defer setupTest(t, &configFile)()

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
	var configFile string
	defer setupTest(t, &configFile)()

	cfg := MustLoadConfig(configFile)

	assert.NotEqual(t, cfg.Chains["mainnet"], nil, "mainnet chain should exist in the configuration")
	assert.Equal(t, []string{"http://localhost:8545"}, cfg.Chains["mainnet"].RPCs, "RPCs for mainnet should remain as default")
	assert.True(t, cfg.Chains["mainnet"].Enabled, "Enabled flag for mainnet should remain as default")
}

func TestEnvironmentVariableOverridesForChains(t *testing.T) {
	var configFile string

	defer setTempEnvVar("TB_KHEDRA_CHAINS_MAINNET_RPCS", "http://rpc1.mainnet,http://rpc2.mainnet")()
	defer setTempEnvVar("TB_KHEDRA_CHAINS_MAINNET_ENABLED", "false")()
	defer setupTest(t, &configFile)()

	cfg := MustLoadConfig(configFile)

	assert.Equal(t, []string{"http://rpc1.mainnet", "http://rpc2.mainnet"}, cfg.Chains["mainnet"].RPCs, "RPCs for mainnet should be overridden by environment variable")
	assert.False(t, cfg.Chains["mainnet"].Enabled, "Enabled flag for mainnet should be overridden by environment variable")
}

func TestInvalidBooleanValueForChains(t *testing.T) {
	defer setTempEnvVar("TB_KHEDRA_CHAINS_MAINNET_ENABLED", "not_a_bool")()
	defer setupTest(t, nil)()

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

// setupTest sets up a temporary folder, updates the getConfigFn pointer, calls establishConfig,
// and assigns the config file path to the provided string pointer if it is not nil.
// Returns a cleanup function to restore the original state.
func setupTest(t *testing.T, configFile *string) func() {
	os.Setenv("TEST_MODE", "true")

	// Use a temporary directory to simulate the config environment
	tmpDir := t.TempDir()
	tempConfigFile := filepath.Join(tmpDir, "config.yaml")

	// If configFile pointer is not nil, assign the path to it
	if configFile != nil {
		*configFile = tempConfigFile
	}

	// Mock getConfigFn to return the temporary config path
	originalGetConfigFn := getConfigFn
	getConfigFn = func() string { return tempConfigFile }

	// Establish the config file
	types.EstablishConfig(tempConfigFile)

	// Return a cleanup function to restore the original state
	return func() {
		os.Unsetenv("TEST_MODE")
		getConfigFn = originalGetConfigFn
	}
}
