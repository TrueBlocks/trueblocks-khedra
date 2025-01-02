// edge_case_tests.go
package config

import (
	"fmt"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMissingChainInConfig(t *testing.T) {
	// Set environment variables for a chain not in the config file
	defer setTempEnvVar("TB_KHEDRA_CHAINS_UNKNOWN_RPCS", "http://unknown.rpc")()
	defer setTempEnvVar("TB_KHEDRA_CHAINS_UNKNOWN_ENABLED", "true")()
	defer setTempEnvVar("TEST_MODE", "true")()

	// Use a temporary directory to simulate missing config
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yaml")

	originalGetConfigFn := getConfigFn
	getConfigFn = func() string { return configFile }
	defer func() { getConfigFn = originalGetConfigFn }()

	establishConfig(configFile)

	cfg := MustLoadConfig(configFile)

	// Verify the chain is not added to the configuration
	chainIndex := -1
	for i, chain := range cfg.Chains {
		if chain.Name == "unknown" {
			chainIndex = i
			break
		}
	}

	assert.Equal(t, -1, chainIndex, "Unknown chain should not be added to the configuration")
}

// func TestInvalidPortForService(t *testing.T) {
// 	// Set an invalid port for the API service
// 	defer setTempEnvVar("TB_KHEDRA_SERVICES_API_PORT", "invalid_port")()
// 	defer setTempEnvVar("TEST_MODE", "true")()

// 	// Use a temporary directory to simulate missing config
// 	tmpDir := t.TempDir()
// 	configFile := filepath.Join(tmpDir, "config.yaml")

// 	originalGetConfigFn := getConfigFn
// 	getConfigFn = func() string { return configFile }
// 	defer func() { getConfigFn = originalGetConfigFn }()

// 	establishConfig(configFile)

// 	// Load the configuration and expect an error
// 	_, err := loadConfig()
// 	assert.Error(t, err, "loadConfig should return an error for invalid port value")
// 	assert.Contains(t, err.Error(), "invalid port value", "Error message should indicate invalid port")
// }

// func TestEmptyRPCsForChain(t *testing.T) {
// 	// Set RPCs for the mainnet chain to an invalid empty value
// 	defer setTempEnvVar("TB_KHEDRA_CHAINS_MAINNET_RPCS", "")()
// 	defer setTempEnvVar("TEST_MODE", "true")()

// 	// Use a temporary directory to simulate missing config
// 	tmpDir := t.TempDir()
// 	configFile := filepath.Join(tmpDir, "config.yaml")

// 	originalGetConfigFn := getConfigFn
// 	getConfigFn = func() string { return configFile }
// 	defer func() { getConfigFn = originalGetConfigFn }()

// 	establishConfig(configFile)

// 	// Load the configuration and expect a validation error
// 	_, err := loadConfig()
// 	assert.Error(t, err, "loadConfig should return an error for empty RPCs")
// 	assert.Contains(t, err.Error(), "strict_url", "Error message should indicate strict_url validation failure")
// }

func TestLargeNumberOfChains(t *testing.T) {
	// Set a large number of chains in the configuration
	defer setTempEnvVar("TEST_MODE", "true")()

	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yaml")

	originalGetConfigFn := getConfigFn
	getConfigFn = func() string { return configFile }
	defer func() { getConfigFn = originalGetConfigFn }()

	establishConfig(configFile)

	cfg := NewConfig()
	cfg.Chains = []Chain{}
	nChains := 1000
	for i := 0; i < nChains; i++ {
		chainName := "chain" + strconv.Itoa(i)
		cfg.Chains = append(cfg.Chains, Chain{
			Name:    chainName,
			RPCs:    []string{fmt.Sprintf("http://%s.rpc", chainName)},
			Enabled: true,
		})
	}

	// Write the large config to the file
	writeConfig(&cfg, configFile)

	// Load the configuration and verify all chains are present
	cfg = *MustLoadConfig(configFile)
	assert.Equal(t, nChains, len(cfg.Chains), "All chains should be loaded correctly")
}
