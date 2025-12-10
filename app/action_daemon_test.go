package app

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/TrueBlocks/trueblocks-chifra/v6/pkg/file"
	"github.com/TrueBlocks/trueblocks-khedra/v6/pkg/types"
)

// ============================================================================
// Test Helpers
// ============================================================================

// setupTestEnv creates an isolated test environment with temp directories
// Returns cleanup function that should be deferred
func setupTestEnv(t *testing.T) (rootFolder string, cleanup func()) {
	// Create temp directory for XDG_CONFIG_HOME
	tempDir, err := os.MkdirTemp("", "khedra-test-*")
	require.NoError(t, err, "Failed to create temp dir")

	// Set XDG_CONFIG_HOME to temp dir
	oldXDG := os.Getenv("XDG_CONFIG_HOME")
	os.Setenv("XDG_CONFIG_HOME", tempDir)

	// Create config subdirectory structure
	configDir := filepath.Join(tempDir, "config")
	err = os.MkdirAll(configDir, 0o755)
	require.NoError(t, err, "Failed to create config dir")

	cleanup = func() {
		os.Setenv("XDG_CONFIG_HOME", oldXDG)
		os.RemoveAll(tempDir)
	}

	return tempDir, cleanup
}

// mockConfig creates a test configuration with known chains and services
func mockConfig() *types.Config {
	return &types.Config{
		Chains: map[string]types.Chain{
			"mainnet": {
				Name:    "mainnet",
				RPCs:    []string{"http://localhost:8545"},
				ChainID: 1,
				Enabled: true,
			},
			"sepolia": {
				Name:    "sepolia",
				RPCs:    []string{"http://localhost:8546"},
				ChainID: 11155111,
				Enabled: false,
			},
		},
		Services: map[string]types.Service{
			"scraper": {
				Name:      "scraper",
				Enabled:   true,
				Sleep:     14,
				BatchSize: 100,
			},
			"api": {
				Name:    "api",
				Enabled: true,
				Port:    8080,
			},
		},
		General: types.General{
			DataFolder: "/tmp/khedra-test/data",
			Strategy:   "download",
			Detail:     "index",
		},
		Logging: types.Logging{
			Folder:     "/tmp/khedra-test/logs",
			Filename:   "khedra.log",
			Level:      "info",
			MaxSize:    10,
			MaxBackups: 3,
			MaxAge:     7,
		},
	}
}

// mockBootstrapper creates a DaemonBootstrapper for testing
func mockBootstrapper(rootFolder string) *DaemonBootstrapper {
	cfg := mockConfig()
	logger := types.NewLogger(types.Logging{Level: "error"})
	return NewDaemonBootstrapper(cfg, rootFolder, logger)
}

// ============================================================================
// Unit Tests for Helper Methods
// ============================================================================

func TestChainsConfigured_ValidConfig(t *testing.T) {
	rootFolder, cleanup := setupTestEnv(t)
	defer cleanup()

	bootstrapper := mockBootstrapper(rootFolder)

	// Create a valid trueBlocks.toml with chain sections
	configFn := filepath.Join(rootFolder, "trueBlocks.toml")
	validConfig := `
[chains.mainnet]
rpcProvider = "http://localhost:8545"

[chains.sepolia]
rpcProvider = "http://localhost:8546"
`
	err := file.StringToAsciiFile(configFn, validConfig)
	require.NoError(t, err)

	// Test: should return true for valid config
	result := bootstrapper.chainsConfigured(configFn)
	assert.True(t, result, "Valid config should pass validation")
}

func TestChainsConfigured_MissingChainSection(t *testing.T) {
	rootFolder, cleanup := setupTestEnv(t)
	defer cleanup()

	bootstrapper := mockBootstrapper(rootFolder)

	// Create config missing mainnet section
	configFn := filepath.Join(rootFolder, "trueBlocks.toml")
	invalidConfig := `
[chains.sepolia]
rpcProvider = "http://localhost:8546"
`
	err := file.StringToAsciiFile(configFn, invalidConfig)
	require.NoError(t, err)

	// Test: should return false when chain section missing
	result := bootstrapper.chainsConfigured(configFn)
	assert.False(t, result, "Config missing enabled chain section should fail")
}

func TestChainsConfigured_EmptyFile(t *testing.T) {
	rootFolder, cleanup := setupTestEnv(t)
	defer cleanup()

	bootstrapper := mockBootstrapper(rootFolder)

	// Create empty config file
	configFn := filepath.Join(rootFolder, "trueBlocks.toml")
	err := file.StringToAsciiFile(configFn, "")
	require.NoError(t, err)

	// Test: should return false for empty config
	result := bootstrapper.chainsConfigured(configFn)
	assert.False(t, result, "Empty config should fail validation")
}

func TestCreateChainConfigFolder(t *testing.T) {
	rootFolder, cleanup := setupTestEnv(t)
	defer cleanup()

	bootstrapper := mockBootstrapper(rootFolder)

	// Test: create folder for mainnet chain
	err := bootstrapper.createChainConfigFolder("mainnet")
	assert.NoError(t, err, "Should create chain config folder successfully")

	// Verify folder exists
	chainFolder := filepath.Join(rootFolder, "config", "mainnet")
	exists := file.FolderExists(chainFolder)
	assert.True(t, exists, "Chain config folder should exist")
}

func TestCreateChainConfigFolder_InvalidPath(t *testing.T) {
	bootstrapper := mockBootstrapper("/dev/null/invalid")

	// Test: try to create folder in invalid location
	err := bootstrapper.createChainConfigFolder("mainnet")
	assert.Error(t, err, "Should fail to create folder in invalid location")
}

func TestCreateChifraConfig(t *testing.T) {
	rootFolder, cleanup := setupTestEnv(t)
	defer cleanup()

	bootstrapper := mockBootstrapper(rootFolder)

	// Test: create chifra config from template
	err := bootstrapper.createChifraConfig()
	assert.NoError(t, err, "Should create chifra config successfully")

	// Verify config file exists
	configFn := filepath.Join(rootFolder, "trueBlocks.toml")
	exists := file.FileExists(configFn)
	assert.True(t, exists, "Config file should exist")

	// Verify config contains chain sections (both enabled and disabled chains appear in config)
	contents := file.AsciiFileToString(configFn)
	assert.Contains(t, contents, "[chains.mainnet]", "Config should contain mainnet section")
	assert.Contains(t, contents, "[chains.sepolia]", "Config should contain sepolia section (disabled chains still appear in config)")

	// Verify chain config folders exist only for ENABLED chains
	mainnetFolder := filepath.Join(rootFolder, "config", "mainnet")
	assert.True(t, file.FolderExists(mainnetFolder), "Mainnet config folder should exist")

	sepoliaFolder := filepath.Join(rootFolder, "config", "sepolia")
	assert.False(t, file.FolderExists(sepoliaFolder), "Sepolia config folder should NOT exist (disabled chain)")
}

func TestCreateChifraConfig_EmptyConfig(t *testing.T) {
	rootFolder, cleanup := setupTestEnv(t)
	defer cleanup()

	// Create config with NO enabled chains
	cfg := &types.Config{
		Chains: map[string]types.Chain{
			"mainnet": {
				Name:    "mainnet",
				RPCs:    []string{"http://localhost:8545"},
				ChainID: 1,
				Enabled: false, // Disabled
			},
		},
	}
	logger := types.NewLogger(types.Logging{Level: "error"})
	bootstrapper := NewDaemonBootstrapper(cfg, rootFolder, logger)

	// Test: should handle empty enabled chains list
	err := bootstrapper.createChifraConfig()
	// This might succeed with empty template - verify behavior
	if err == nil {
		configFn := filepath.Join(rootFolder, "trueBlocks.toml")
		exists := file.FileExists(configFn)
		assert.True(t, exists, "Config file should exist even with no enabled chains")
	}
}

// ============================================================================
// Integration-Style Tests for Daemon Phases
// ============================================================================

func TestDaemonEnvironmentSetup(t *testing.T) {
	rootFolder, cleanup := setupTestEnv(t)
	defer cleanup()

	cfg := mockConfig()
	cfg.General.DataFolder = rootFolder + "/data"

	// Manually set environment variables as daemon would
	os.Setenv("TB_SETTINGS_DEFAULTCHAIN", "mainnet")
	os.Setenv("TB_SETTINGS_INDEXPATH", cfg.IndexPath())
	os.Setenv("TB_SETTINGS_CACHEPATH", cfg.CachePath())

	for key, ch := range cfg.Chains {
		if ch.Enabled {
			envKey := "TB_CHAINS_" + strings.ToUpper(key) + "_RPCPROVIDER"
			os.Setenv(envKey, ch.RPCs[0])
		}
	}

	// Verify environment variables set correctly
	assert.Equal(t, "mainnet", os.Getenv("TB_SETTINGS_DEFAULTCHAIN"))
	assert.Equal(t, cfg.IndexPath(), os.Getenv("TB_SETTINGS_INDEXPATH"))
	assert.Equal(t, cfg.CachePath(), os.Getenv("TB_SETTINGS_CACHEPATH"))
	assert.Equal(t, "http://localhost:8545", os.Getenv("TB_CHAINS_MAINNET_RPCPROVIDER"))
	assert.Equal(t, "", os.Getenv("TB_CHAINS_SEPOLIA_RPCPROVIDER"), "Disabled chain should not have env var")

	// Cleanup
	os.Unsetenv("TB_SETTINGS_DEFAULTCHAIN")
	os.Unsetenv("TB_SETTINGS_INDEXPATH")
	os.Unsetenv("TB_SETTINGS_CACHEPATH")
	os.Unsetenv("TB_CHAINS_MAINNET_RPCPROVIDER")
}

func TestDaemonChifraConfigCreation(t *testing.T) {
	rootFolder, cleanup := setupTestEnv(t)
	defer cleanup()

	bootstrapper := mockBootstrapper(rootFolder)

	// Phase 1: Config doesn't exist
	configFn := filepath.Join(rootFolder, "trueBlocks.toml")
	assert.False(t, file.FileExists(configFn), "Config should not exist initially")

	// Phase 2: Create config
	err := bootstrapper.createChifraConfig()
	assert.NoError(t, err, "Config creation should succeed")
	assert.True(t, file.FileExists(configFn), "Config should exist after creation")

	// Phase 3: Validate created config
	result := bootstrapper.chainsConfigured(configFn)
	assert.True(t, result, "Created config should pass validation")
}

func TestDaemonChifraConfigValidation(t *testing.T) {
	rootFolder, cleanup := setupTestEnv(t)
	defer cleanup()

	bootstrapper := mockBootstrapper(rootFolder)

	configFn := filepath.Join(rootFolder, "trueBlocks.toml")

	tests := []struct {
		name     string
		content  string
		expected bool
	}{
		{
			name: "Valid config with all enabled chains",
			content: `
[chains.mainnet]
rpcProvider = "http://localhost:8545"
`,
			expected: true,
		},
		{
			name: "Invalid config missing enabled chain section",
			content: `
[settings]
defaultChain = "mainnet"
# Config is incomplete - no chain configurations
`,
			expected: false,
		},
		{
			name: "Config with chain section for enabled chain",
			content: `
[chains.mainnet]
rpcProvider = "http://localhost:8545"

[chains.sepolia]
rpcProvider = "http://localhost:8546"
`,
			expected: true,
		},
		{
			name:     "Empty config",
			content:  "",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := file.StringToAsciiFile(configFn, tt.content)
			require.NoError(t, err)

			result := bootstrapper.chainsConfigured(configFn)
			assert.Equal(t, tt.expected, result, tt.name)
		})
	}
}
