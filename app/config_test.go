package app

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	yamlv2 "gopkg.in/yaml.v2"

	coreFile "github.com/TrueBlocks/trueblocks-chifra/v6/pkg/file"
	"github.com/TrueBlocks/trueblocks-khedra/v6/pkg/types"
)

// Testing status: not_reviewed

func TestConfigMustLoad(t *testing.T) {
	defer types.SetupTest([]string{})()
	assert.FileExists(t, types.GetConfigFn())
}

func TestLoadConfig_ValidConfig(t *testing.T) {
	defer types.SetupTest([]string{})()

	cfg := types.Config{
		Chains: map[string]types.Chain{
			"mainnet": {Name: "mainnet", RPCs: []string{"http://rpc1.mainnet"}, Enabled: true, ChainID: 1},
		},
		Services: map[string]types.Service{
			"scraper": types.NewService("scraper"),
			"monitor": types.NewService("monitor"),
			"api":     types.NewService("api"),
			"ipfs":    types.NewService("ipfs"),
		},
		Logging: types.NewLogging(),
		General: types.General{
			DataFolder: "/tmp/test-data-folder",
			Strategy:   "scratch",
			Detail:     "index",
		},
	}

	bytes, _ := yamlv2.Marshal(cfg)
	_ = coreFile.StringToAsciiFile(types.GetConfigFn(), string(bytes))

	result, err := LoadConfig()
	assert.NoError(t, err, t.Name())
	assert.Equal(t, cfg, result)

	os.RemoveAll(cfg.Logging.Folder)
	os.RemoveAll(cfg.General.DataFolder)
	os.Remove(types.GetConfigFn())
}

func TestLoadConfig_InvalidFileConfig(t *testing.T) {
	defer types.SetupTest([]string{})()

	_ = os.WriteFile(types.GetConfigFn(), []byte("invalid_yaml"), 0o644)

	_, err := LoadConfig()
	assert.Error(t, err, t.Name())
	if err != nil {
		assert.Contains(t, err.Error(), "failed to load file configuration")
	}

	os.Remove(types.GetConfigFn())
}

func TestLoadConfig_EnvOverrides(t *testing.T) {
	defer types.SetupTest([]string{
		"TB_KHEDRA_CHAINS_MAINNET_RPCS=http://env.rpc1.mainnet,http://env.rpc2.mainnet",
		"TB_KHEDRA_SERVICES_API_PORT=9090",
		"TB_KHEDRA_GENERAL_DATAFOLDER=/tmp/env-data-folder",
		"TB_KHEDRA_GENERAL_DETAIL=index",
	})()

	cfg := types.Config{
		Chains: map[string]types.Chain{
			"mainnet": {
				RPCs:    []string{"http://rpc1.mainnet"},
				Enabled: true,
				ChainID: 1,
			},
		},
		Services: map[string]types.Service{
			"api": {
				Port:      8080,
				BatchSize: 100,
				Enabled:   true,
			},
		},
		Logging: types.Logging{
			Folder:     "/tmp/test-logging-folder",
			Filename:   "test.log",
			MaxSize:    50,
			MaxBackups: 1,
			MaxAge:     1,
			Level:      "info",
		},
		General: types.General{
			DataFolder: "/tmp/test-data-folder",
			Strategy:   "download",
			Detail:     "bloom",
		},
	}

	bytes, _ := yamlv2.Marshal(cfg)
	_ = coreFile.StringToAsciiFile(types.GetConfigFn(), string(bytes))

	result, err := LoadConfig()
	assert.NoError(t, err, t.Name())
	assert.Equal(t, []string{"http://env.rpc1.mainnet", "http://env.rpc2.mainnet"}, result.Chains["mainnet"].RPCs)
	assert.Equal(t, 9090, result.Services["api"].Port)
	assert.Equal(t, "/tmp/env-data-folder", result.General.DataFolder)
	assert.Equal(t, "index", result.General.Detail)

	os.RemoveAll(cfg.Logging.Folder)
	os.RemoveAll(result.General.DataFolder)
	os.Remove(types.GetConfigFn())
}

func TestLoadConfig_ValidationFailure(t *testing.T) {
	defer types.SetupTest([]string{})()

	cfg := types.Config{
		Chains: map[string]types.Chain{
			"mainnet": {
				RPCs:    []string{},
				Enabled: true,
				ChainID: 1,
			},
		},
		Services: map[string]types.Service{
			"api": {
				Port:      8080,
				BatchSize: 0,
				Enabled:   true,
			},
		},
		Logging: types.Logging{
			Folder:   "",
			Filename: "",
			MaxSize:  0,
		},
		General: types.General{
			DataFolder: "",
			Strategy:   "download",
			Detail:     "index",
		},
	}

	bytes, _ := yamlv2.Marshal(cfg)
	_ = coreFile.StringToAsciiFile(types.GetConfigFn(), string(bytes))

	_, err := LoadConfig()
	assert.Error(t, err, t.Name())
	assert.Contains(t, err.Error(), "is required")

	os.Remove(types.GetConfigFn())
}

func TestChainEnvOverrides(t *testing.T) {
	defer types.SetupTest([]string{
		"TB_KHEDRA_CHAINS_MAINNET_RPCS=http://rpc1.mainnet,http://rpc2.mainnet",
		"TB_KHEDRA_CHAINS_MAINNET_ENABLED=false",
		"TB_KHEDRA_CHAINS_MAINNET_CHAINID=2",
	})()

	cfg, err := LoadConfig()
	assert.NoError(t, err, t.Name())
	assert.Equal(t, []string{"http://rpc1.mainnet", "http://rpc2.mainnet"}, cfg.Chains["mainnet"].RPCs, "RPCs for mainnet should be overridden by environment variable")
	assert.Equal(t, 2, cfg.Chains["mainnet"].ChainID, "ChainID for mainnet should be overridden by environment variable")
}

func TestChainInvalidBooleanValue(t *testing.T) {
	defer types.SetupTest([]string{
		"TB_KHEDRA_CHAINS_MAINNET_ENABLED=not_a_bool",
	})()

	if cfg, err := LoadConfig(); err != nil {
		assert.Error(t, err, "loadConfig should return an error for invalid boolean value", t.Name())
		assert.Contains(t, err.Error(), "failed to apply environment configuration: environment variable has an invalid value", "Error message should indicate the inability to parse the boolean value")
		assert.Contains(t, err.Error(), "key=[enabled], value=[]", "Error message should point to the problematic field")
	} else {
		t.Error("loadConfig should return an error for invalid boolean value", cfg.Chains["mainnet"].Enabled, t.Name())
	}
}

func TestServiceEnvironmentVariableOverrides(t *testing.T) {
	defer types.SetupTest([]string{
		"TB_KHEDRA_SERVICES_API_ENABLED=false",
		"TB_KHEDRA_SERVICES_API_PORT=9090",
	})()

	if cfg, err := LoadConfig(); err != nil {
		t.Error(err, t.Name())
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
		t.Error(err, t.Name())
	} else {
		apiService, apiExists := cfg.Services["api"]
		scraperService, scraperExists := cfg.Services["scraper"]
		assert.True(t, apiExists, "API service should exist in the configuration")
		assert.True(t, scraperExists, "Scraper service should exist in the configuration")
		assert.False(t, apiService.Enabled, "Enabled flag for API service should be overridden by environment variable")
		assert.True(t, scraperService.Enabled, "Enabled flag for Scraper service should be overridden by environment variable")
		assert.Equal(t, 0, scraperService.Port, "Port for Scraper service should NOT be overridden by environment variable")
	}
}

func TestEnvNoVariables(t *testing.T) {
	defer types.SetupTest([]string{})()

	if cfg, err := LoadConfig(); err != nil {
		t.Error(err, t.Name())
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

	cfg, err := LoadConfig()
	if err != nil {
		assert.Error(t, err, "loadConfig should return an error for invalid port value", t.Name())
		assert.Contains(t, err.Error(), "failed to apply environment configuration: environment variable has an invalid value: key=[port], value=[]", "Error message should indicate invalid port")
	} else {
		t.Error("loadConfig should return an error for invalid port", cfg.Services["api"], t.Name())
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
			ChainID: i + 1,
		}
	}

	bytes, _ := yamlv2.Marshal(cfg)
	_ = coreFile.StringToAsciiFile(types.GetConfigFn(), string(bytes))

	var err error
	if cfg, err = LoadConfig(); err != nil {
		t.Error(err, t.Name())
	} else {
		assert.Equal(t, nChains+1, len(cfg.Chains), "All chains should be loaded correctly")
	}
}

func TestChainMissingInConfig(t *testing.T) {
	defer types.SetupTest([]string{
		"TB_KHEDRA_CHAINS_UNKNOWN_NAME=unknown",
		"TB_KHEDRA_CHAINS_UNKNOWN_RPCS=http://unknown.rpc",
		"TB_KHEDRA_CHAINS_UNKNOWN_ENABLED=true",
	})()

	if cfg, err := LoadConfig(); err != nil {
		assert.Error(t, err, "An error should occur if an unknown chain is defined in the environment but not in the configuration file", t.Name())
	} else {
		assert.Equal(t, 1, len(cfg.Chains), "Only the mainnet chain should be loaded")
		assert.Equal(t, "mainnet", cfg.Chains["mainnet"].Name, "Only the mainnet chain should be loaded")
	}
}

func TestChainEmptyRPCs(t *testing.T) {
	defer types.SetupTest([]string{
		"TB_KHEDRA_CHAINS_MAINNET_RPCS=",
	})()

	if cfg, err := LoadConfig(); err != nil {
		assert.Error(t, err, t.Name())
	} else {
		assert.NotEmpty(t, cfg.Chains["mainnet"].RPCs, "Mainnet RPCs should not be empty in the final configuration")
	}
}

func TestConfigMustLoadDefaults(t *testing.T) {
	defer types.SetupTest([]string{})()

	if cfg, err := LoadConfig(); err != nil {
		t.Error(err, t.Name())
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

// Phase 4a: Config pipeline tests

func TestLoadFileConfig_EmptyFile(t *testing.T) {
	defer types.SetupTest([]string{})()
	_ = os.WriteFile(types.GetConfigFn(), []byte(""), 0o644)

	loader := NewConfigLoader()
	_, err := loader.loadFromFile()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "empty")
}

func TestLoadFileConfig_SetsChainNames(t *testing.T) {
	defer types.SetupTest([]string{})()

	cfg := types.Config{
		Chains: map[string]types.Chain{
			"mainnet": {RPCs: []string{"http://rpc1"}, Enabled: true, ChainID: 1},
			"sepolia": {RPCs: []string{"http://rpc2"}, Enabled: true, ChainID: 11155111},
		},
		Services: map[string]types.Service{
			"scraper": types.NewService("scraper"),
		},
		Logging: types.NewLogging(),
		General: types.General{DataFolder: "/tmp/data", Strategy: "scratch", Detail: "index"},
	}

	bytes, _ := yamlv2.Marshal(cfg)
	_ = coreFile.StringToAsciiFile(types.GetConfigFn(), string(bytes))

	loader := NewConfigLoader()
	result, err := loader.loadFromFile()
	assert.NoError(t, err)
	assert.Equal(t, "mainnet", result.Chains["mainnet"].Name)
	assert.Equal(t, "sepolia", result.Chains["sepolia"].Name)
	assert.Equal(t, "scraper", result.Services["scraper"].Name)
}

func TestFinalCleanup_CleansPaths(t *testing.T) {
	cfg := types.Config{
		General: types.General{DataFolder: "/tmp//data/../data/"},
		Logging: types.Logging{Folder: "/tmp/logs/./"},
	}

	loader := NewConfigLoader()
	err := loader.cleanup(&cfg)
	assert.NoError(t, err)
	assert.Equal(t, "/tmp/data", cfg.General.DataFolder)
	assert.Equal(t, "/tmp/logs", cfg.Logging.Folder)
}

func TestInitializeFolders_CreatesDirectories(t *testing.T) {
	tempDir := t.TempDir()
	cfg := types.Config{
		General: types.General{DataFolder: filepath.Join(tempDir, "data")},
		Logging: types.Logging{Folder: filepath.Join(tempDir, "logs")},
	}

	loader := NewConfigLoader()
	err := loader.initializeFolders(cfg)
	assert.NoError(t, err)

	_, err = os.Stat(filepath.Join(tempDir, "data"))
	assert.NoError(t, err)
	_, err = os.Stat(filepath.Join(tempDir, "logs"))
	assert.NoError(t, err)
}
