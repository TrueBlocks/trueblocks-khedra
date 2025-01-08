package app

import (
	"fmt"
	"os"
	"strconv"
	"testing"

	coreFile "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/types"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

// TODO: FIX TEST
// func TestChainEnvOverrides(t *testing.T) {
// 	defer types.SetupTest([]string{
// 		"TB_KHEDRA_CHAINS_MAINNET_RPCS=http://rpc1.mainnet,http://rpc2.mainnet",
// 		"TB_KHEDRA_CHAINS_MAINNET_ENABLED=false",
// 		"TB_KHEDRA_CHAINS_SEPOLIA_ENABLED=true",
// 	})()

// 	if cfg, err := LoadConfig(); err != nil {
// 		t.Error(err)
// 	} else {
// 		assert.Equal(t, []string{"http://rpc1.mainnet", "http://rpc2.mainnet"}, cfg.Chains["mainnet"].RPCs, "RPCs for mainnet should be overridden by environment variable")
// 		assert.False(t, cfg.Chains["mainnet"].Enabled, "Enabled flag for mainnet should be overridden by environment variable")
// 		assert.True(t, cfg.Chains["sepolia"].Enabled, "Enabled flag for sepolia should be overridden by environment variable")
// 	}
// }

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

func TestServiceInvalidPort(t *testing.T) {
	defer types.SetupTest([]string{
		"TB_KHEDRA_SERVICES_API_PORT=invalid_port",
	})()

	if cfg, err := LoadConfig(); err != nil {
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
	if cfg, err = LoadConfig(); err != nil {
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

	if cfg, err := LoadConfig(); err != nil {
		assert.Error(t, err, "An error should occur if an unknown chain is defined in the environment but not in the configuration file")
	} else {
		t.Error("loadConfig should return an error for invalid chain", cfg.Chains["unknown"])
	}
}

func TestChainEmptyRPCs(t *testing.T) {
	defer types.SetupTest([]string{
		"TB_KHEDRA_CHAINS_MAINNET_RPCS=",
	})()

	if cfg, err := LoadConfig(); err != nil {
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

	if cfg, err := LoadConfig(); err != nil {
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

func TestLoadFileConfig_ValidFile(t *testing.T) {
	defer types.SetupTest([]string{})()

	cfg := types.NewConfig()
	chain := cfg.Chains["mainnet"]
	chain.RPCs = []string{"http://localhost:8545", "http://localhost:8546"}
	cfg.Chains["mainnet"] = chain
	bytes, _ := yaml.Marshal(cfg)
	coreFile.StringToAsciiFile(types.GetConfigFn(), string(bytes))

	result, err := loadFileConfig()
	assert.NoError(t, err)
	assert.Equal(t, cfg, result)
}

func TestLoadFileConfig_InvalidFile(t *testing.T) {
	defer types.SetupTest([]string{})()
	coreFile.StringToAsciiFile(types.GetConfigFn(), "invalid: [:::]")

	_, err := loadFileConfig()
	assert.Error(t, err)
}

func TestLoadFileConfig_MissingFile(t *testing.T) {
	defer types.SetupTest([]string{})()
	os.Remove(types.GetConfigFn())

	_, err := loadFileConfig()
	assert.Error(t, err)
}

func TestLoadEnvConfig_ValidEnvVars(t *testing.T) {
	defer types.SetupTest([]string{
		"TB_KHEDRA_CHAINS_MAINNET_RPCS=http://rpc1.mainnet,http://rpc2.mainnet",
		"TB_KHEDRA_CHAINS_MAINNET_ENABLED=true",
		"TB_KHEDRA_SERVICES_API_PORT=9090",
		"TB_KHEDRA_SERVICES_API_ENABLED=false",
	})()

	result, err := loadEnvConfig()
	assert.NoError(t, err)
	assert.Equal(t, []string{"http://rpc1.mainnet", "http://rpc2.mainnet"}, result.Chains["mainnet"].RPCs)
	assert.True(t, result.Chains["mainnet"].Enabled)
	assert.Equal(t, 9090, result.Services["api"].Port)
	assert.False(t, result.Services["api"].Enabled)
}

func TestLoadEnvConfig_InvalidBoolean(t *testing.T) {
	defer types.SetupTest([]string{
		"TB_KHEDRA_CHAINS_MAINNET_ENABLED=not_a_bool",
	})()

	_, err := loadEnvConfig()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal environment config")
}

func TestLoadEnvConfig_NoEnvVars(t *testing.T) {
	defer types.SetupTest([]string{})()

	result, err := loadEnvConfig()
	assert.NoError(t, err)
	assert.Empty(t, result.Chains)
	assert.Empty(t, result.Services)
}

func TestLoadEnvConfig_InvalidInteger(t *testing.T) {
	defer types.SetupTest([]string{
		"TB_KHEDRA_SERVICES_API_PORT=not_an_int",
	})()

	_, err := loadEnvConfig()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to unmarshal environment config")
}

func TestLoadEnvConfig_SliceParsingError(t *testing.T) {
	defer types.SetupTest([]string{
		"TB_KHEDRA_CHAINS_MAINNET_RPCS=http://valid.rpc,http://invalid\\rpc",
	})()

	result, err := loadEnvConfig()
	assert.NoError(t, err)
	assert.Equal(t, []string{"http://valid.rpc", "http://invalid\\rpc"}, result.Chains["mainnet"].RPCs)
}

func TestLoadEnvConfig_CaseInsensitiveKeys(t *testing.T) {
	defer types.SetupTest([]string{
		"tb_khedra_chains_mainnet_rpcs=http://rpc1.mainnet,http://rpc2.mainnet",
		"TB_KHEDRA_CHAINS_MAINNET_ENABLED=true",
	})()

	result, err := loadEnvConfig()
	assert.NoError(t, err)
	assert.Empty(t, result.Chains["mainnet"].RPCs, "Expected mainnet RPCs to be empty due to incorrect case")
	assert.True(t, result.Chains["mainnet"].Enabled)
}

func TestLoadEnvConfig_AllVariables(t *testing.T) {
	defer types.SetupTest([]string{
		"TB_KHEDRA_GENERAL_DATAFOLDER=/data/khedra",

		"TB_KHEDRA_CHAINS_MAINNET_RPCS=http://rpc1.mainnet,http://rpc2.mainnet",
		"TB_KHEDRA_CHAINS_MAINNET_ENABLED=false",
		"TB_KHEDRA_CHAINS_SEPOLIA_RPCS=http://rpc1.sepolia,http://rpc2.sepolia",
		"TB_KHEDRA_CHAINS_SEPOLIA_ENABLED=false",

		"TB_KHEDRA_SERVICES_SCRAPER_ENABLED=false",
		"TB_KHEDRA_SERVICES_SCRAPER_PORT=500",
		"TB_KHEDRA_SERVICES_SCRAPER_SLEEP=500",
		"TB_KHEDRA_SERVICES_SCRAPER_BATCHSIZE=500",
		"TB_KHEDRA_SERVICES_MONITOR_ENABLED=false",
		"TB_KHEDRA_SERVICES_MONITOR_PORT=500",
		"TB_KHEDRA_SERVICES_MONITOR_SLEEP=500",
		"TB_KHEDRA_SERVICES_MONITOR_BATCHSIZE=500",
		"TB_KHEDRA_SERVICES_API_ENABLED=false",
		"TB_KHEDRA_SERVICES_API_PORT=500",
		"TB_KHEDRA_SERVICES_API_SLEEP=500",
		"TB_KHEDRA_SERVICES_API_BATCHSIZE=500",
		"TB_KHEDRA_SERVICES_IPFS_ENABLED=false",
		"TB_KHEDRA_SERVICES_IPFS_PORT=500",
		"TB_KHEDRA_SERVICES_IPFS_SLEEP=500",
		"TB_KHEDRA_SERVICES_IPFS_BATCHSIZE=500",

		"TB_KHEDRA_LOGGING_FOLDER=/var/log/khedra",
		"TB_KHEDRA_LOGGING_FILENAME=khedra.log",
		"TB_KHEDRA_LOGGING_MAXSIZE=500",
		"TB_KHEDRA_LOGGING_MAXBACKUPS=500",
		"TB_KHEDRA_LOGGING_MAXAGE=500",
		"TB_KHEDRA_LOGGING_COMPRESS=false",
		"TB_KHEDRA_LOGGING_LEVEL=error",
	})()

	result, err := loadEnvConfig()
	assert.NoError(t, err)
	assert.Equal(t, "/data/khedra", result.General.DataFolder)

	assert.Equal(t, []string{"http://rpc1.mainnet", "http://rpc2.mainnet"}, result.Chains["mainnet"].RPCs)
	assert.False(t, result.Chains["mainnet"].Enabled)
	assert.Equal(t, []string{"http://rpc1.sepolia", "http://rpc2.sepolia"}, result.Chains["sepolia"].RPCs)
	assert.False(t, result.Chains["sepolia"].Enabled)

	assert.False(t, result.Services["scraper"].Enabled)
	assert.Equal(t, 500, result.Services["scraper"].Port)
	assert.Equal(t, 500, result.Services["scraper"].Sleep)
	assert.Equal(t, 500, result.Services["scraper"].BatchSize)

	assert.False(t, result.Services["monitor"].Enabled)
	assert.Equal(t, 500, result.Services["monitor"].Port)
	assert.Equal(t, 500, result.Services["monitor"].Sleep)
	assert.Equal(t, 500, result.Services["monitor"].BatchSize)

	assert.False(t, result.Services["api"].Enabled)
	assert.Equal(t, 500, result.Services["api"].Port)
	assert.Equal(t, 500, result.Services["api"].Sleep)
	assert.Equal(t, 500, result.Services["api"].BatchSize)

	assert.False(t, result.Services["ipfs"].Enabled)
	assert.Equal(t, 500, result.Services["ipfs"].Port)
	assert.Equal(t, 500, result.Services["ipfs"].Sleep)
	assert.Equal(t, 500, result.Services["ipfs"].BatchSize)

	assert.Equal(t, "/var/log/khedra", result.Logging.Folder)
	assert.Equal(t, "khedra.log", result.Logging.Filename)
	assert.Equal(t, 500, result.Logging.MaxSize)
	assert.Equal(t, 500, result.Logging.MaxBackups)
	assert.Equal(t, 500, result.Logging.MaxAge)
	assert.False(t, result.Logging.Compress)
	assert.Equal(t, "error", result.Logging.Level)
}

func TestMergeConfigs_ChainsMerge(t *testing.T) {
	fileCfg := types.Config{
		Chains: map[string]types.Chain{
			"mainnet": {RPCs: []string{"http://file.rpc"}, Enabled: false},
		},
	}
	envCfg := types.Config{
		Chains: map[string]types.Chain{
			"mainnet": {RPCs: []string{"http://env.rpc"}, Enabled: true},
		},
	}

	result, err := mergeConfigs(fileCfg, envCfg)
	assert.NoError(t, err)
	assert.Equal(t, []string{"http://env.rpc"}, result.Chains["mainnet"].RPCs)
	assert.True(t, result.Chains["mainnet"].Enabled)
}

func TestMergeConfigs_ChainsPreserveFile(t *testing.T) {
	fileCfg := types.Config{
		Chains: map[string]types.Chain{
			"mainnet": {RPCs: []string{"http://file.rpc"}, Enabled: true},
		},
	}
	envCfg := types.Config{
		Chains: map[string]types.Chain{
			"mainnet": {RPCs: []string{}, Enabled: false},
		},
	}

	result, err := mergeConfigs(fileCfg, envCfg)
	assert.NoError(t, err)
	assert.Equal(t, []string{"http://file.rpc"}, result.Chains["mainnet"].RPCs)
	assert.True(t, result.Chains["mainnet"].Enabled)
}

func TestMergeConfigs_ServicesMerge(t *testing.T) {
	fileCfg := types.Config{
		Services: map[string]types.Service{
			"api": {BatchSize: 100, Port: 8080, Enabled: false},
		},
	}
	envCfg := types.Config{
		Services: map[string]types.Service{
			"api": {BatchSize: 200, Port: 9090, Enabled: true},
		},
	}

	result, err := mergeConfigs(fileCfg, envCfg)
	assert.NoError(t, err)
	assert.Equal(t, 200, result.Services["api"].BatchSize)
	assert.Equal(t, 9090, result.Services["api"].Port)
	assert.True(t, result.Services["api"].Enabled)
}

func TestMergeConfigs_ServicesPreserveFile(t *testing.T) {
	fileCfg := types.Config{
		Services: map[string]types.Service{
			"api": {BatchSize: 100, Port: 8080, Enabled: true},
		},
	}
	envCfg := types.Config{
		Services: map[string]types.Service{
			"api": {BatchSize: 0, Port: 0, Enabled: false},
		},
	}

	result, err := mergeConfigs(fileCfg, envCfg)
	assert.NoError(t, err)
	assert.Equal(t, 100, result.Services["api"].BatchSize)
	assert.Equal(t, 8080, result.Services["api"].Port)
	assert.True(t, result.Services["api"].Enabled)
}

func TestMergeConfigs_LoggingMerge(t *testing.T) {
	fileCfg := types.Config{
		Logging: types.Logging{
			Folder:     "/var/log/file",
			Filename:   "file.log",
			MaxSize:    50,
			MaxBackups: 5,
			MaxAge:     7,
			Compress:   false,
		},
	}
	envCfg := types.Config{
		Logging: types.Logging{
			Folder:     "/var/log/env",
			Filename:   "env.log",
			MaxSize:    100,
			MaxBackups: 10,
			MaxAge:     30,
			Compress:   true,
		},
	}

	result, err := mergeConfigs(fileCfg, envCfg)
	assert.NoError(t, err)
	assert.Equal(t, "/var/log/env", result.Logging.Folder)
	assert.Equal(t, "env.log", result.Logging.Filename)
	assert.Equal(t, 100, result.Logging.MaxSize)
	assert.Equal(t, 10, result.Logging.MaxBackups)
	assert.Equal(t, 30, result.Logging.MaxAge)
	assert.True(t, result.Logging.Compress)
}

func TestMergeConfigs_LoggingPreserveFile(t *testing.T) {
	fileCfg := types.Config{
		Logging: types.Logging{
			Folder:     "/var/log/file",
			Filename:   "file.log",
			MaxSize:    50,
			MaxBackups: 5,
			MaxAge:     7,
			Compress:   false,
		},
	}
	envCfg := types.Config{
		Logging: types.Logging{
			Folder:     "",
			Filename:   "",
			MaxSize:    0,
			MaxBackups: 0,
			MaxAge:     0,
			Compress:   false,
		},
	}

	result, err := mergeConfigs(fileCfg, envCfg)
	assert.NoError(t, err)
	assert.Equal(t, "/var/log/file", result.Logging.Folder)
	assert.Equal(t, "file.log", result.Logging.Filename)
	assert.Equal(t, 50, result.Logging.MaxSize)
	assert.Equal(t, 5, result.Logging.MaxBackups)
	assert.Equal(t, 7, result.Logging.MaxAge)
	assert.False(t, result.Logging.Compress)
}

func TestMergeConfigs_GeneralMerge(t *testing.T) {
	fileCfg := types.Config{
		General: types.General{DataFolder: "/data/file"},
	}
	envCfg := types.Config{
		General: types.General{DataFolder: "/data/env"},
	}

	result, err := mergeConfigs(fileCfg, envCfg)
	assert.NoError(t, err)
	assert.Equal(t, "/data/env", result.General.DataFolder)
}

func TestMergeConfigs_GeneralPreserveFile(t *testing.T) {
	fileCfg := types.Config{
		General: types.General{DataFolder: "/data/file"},
	}
	envCfg := types.Config{
		General: types.General{DataFolder: ""},
	}

	result, err := mergeConfigs(fileCfg, envCfg)
	assert.NoError(t, err)
	assert.Equal(t, "/data/file", result.General.DataFolder)
}

func TestValidateConfig_ValidConfig(t *testing.T) {
	cfg := types.Config{
		Chains: map[string]types.Chain{
			"mainnet": {RPCs: []string{"http://rpc1.mainnet"}, Enabled: true},
		},
		Services: map[string]types.Service{
			"api": {Port: 8080, BatchSize: 100, Enabled: true},
		},
		Logging: types.Logging{
			Folder:   "~/.khedra",
			Filename: "khedra.log",
			MaxSize:  50,
		},
		General: types.General{
			DataFolder: "~/.khedra/data",
		},
	}

	err := validateConfig(cfg)
	assert.NoError(t, err)
}

func TestValidateConfig_MissingRPCs(t *testing.T) {
	cfg := types.Config{
		Chains: map[string]types.Chain{
			"mainnet": {RPCs: []string{}, Enabled: true},
		},
	}

	err := validateConfig(cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "chain mainnet has no RPCs defined")
}

// TODO: FIX TEST
// func TestValidateConfig_InvalidServicePort(t *testing.T) {
// 	cfg := types.Config{
// 		Services: map[string]types.Service{
// 			"api": {Port: 80, BatchSize: 100, Enabled: true},
// 		},
// 	}

// 	err := validateConfig(cfg)
// 	assert.Error(t, err)
// 	assert.Contains(t, err.Error(), "service api has an invalid port")
// }

// TODO: FIX TEST
// func TestValidateConfig_InvalidBatchSize(t *testing.T) {
// 	cfg := types.Config{
// 		Services: map[string]types.Service{
// 			"api": {Port: 8080, BatchSize: 0, Enabled: true},
// 		},
// 	}

// 	err := validateConfig(cfg)
// 	assert.Error(t, err)
// 	assert.Contains(t, err.Error(), "service api has an invalid batch size")
// }

func TestValidateConfig_MissingLoggingConfig(t *testing.T) {
	cfg := types.Config{
		Logging: types.Logging{
			Folder:   "",
			Filename: "",
		},
	}

	err := validateConfig(cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "logging folder is not defined")
}

func TestValidateConfig_MissingGeneralConfig(t *testing.T) {
	cfg := types.Config{
		General: types.General{
			DataFolder: "",
		},
	}

	err := validateConfig(cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "logging folder is not defined")
}

func TestInitializeFolders_AllFoldersExist(t *testing.T) {
	cfg := types.Config{
		Logging: types.Logging{
			Folder: "/tmp/test-logging-folder",
		},
		General: types.General{
			DataFolder: "/tmp/test-data-folder",
		},
	}

	// Ensure folders exist before running the test
	os.MkdirAll(cfg.Logging.Folder, os.ModePerm)
	os.MkdirAll(cfg.General.DataFolder, os.ModePerm)

	err := initializeFolders(cfg)
	assert.NoError(t, err)

	// Clean up
	os.RemoveAll(cfg.Logging.Folder)
	os.RemoveAll(cfg.General.DataFolder)
}

func TestInitializeFolders_CreateMissingFolders(t *testing.T) {
	cfg := types.Config{
		Logging: types.Logging{
			Folder: "/tmp/test-missing-logging-folder",
		},
		General: types.General{
			DataFolder: "/tmp/test-missing-data-folder",
		},
	}

	// Ensure folders do not exist before running the test
	os.RemoveAll(cfg.Logging.Folder)
	os.RemoveAll(cfg.General.DataFolder)

	err := initializeFolders(cfg)
	assert.NoError(t, err)

	// Verify that folders were created
	_, err = os.Stat(cfg.Logging.Folder)
	assert.NoError(t, err)
	_, err = os.Stat(cfg.General.DataFolder)
	assert.NoError(t, err)

	// Clean up
	os.RemoveAll(cfg.Logging.Folder)
	os.RemoveAll(cfg.General.DataFolder)
}

func TestInitializeFolders_ErrorOnInvalidPath(t *testing.T) {
	cfg := types.Config{
		Logging: types.Logging{
			Folder: "/invalid-folder-path/\\0",
		},
		General: types.General{
			DataFolder: "/tmp/test-data-folder",
		},
	}

	err := initializeFolders(cfg)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to create folder")

	// Clean up
	os.RemoveAll(cfg.General.DataFolder)
}

// ------------------------------------------------------------
// TODO: FIX TEST
// func TestLoadConfig_ValidConfig(t *testing.T) {
// 	defer types.SetupTest([]string{})()

// 	// Set up a valid configuration file
// 	cfg := types.Config{
// 		Chains: map[string]types.Chain{
// 			"mainnet": {RPCs: []string{"http://rpc1.mainnet"}, Enabled: true},
// 		},
// 		Services: map[string]types.Service{
// 			"api": {Port: 8080, BatchSize: 100, Enabled: true},
// 		},
// 		Logging: types.Logging{
// 			Folder:   "/tmp/test-logging-folder",
// 			Filename: "test.log",
// 			MaxSize:  50,
// 		},
// 		General: types.General{
// 			DataFolder: "/tmp/test-data-folder",
// 		},
// 	}

// 	// Save the config to the file
// 	bytes, _ := yaml.Marshal(cfg)
// 	coreFile.StringToAsciiFile(types.GetConfigFn(), string(bytes))

// 	// Run LoadConfig
// 	result, err := LoadConfig()
// 	assert.NoError(t, err)
// 	assert.Equal(t, cfg, result)

// 	// Clean up
// 	os.RemoveAll(cfg.Logging.Folder)
// 	os.RemoveAll(cfg.General.DataFolder)
// 	os.Remove(types.GetConfigFn())
// }

func TestLoadConfig_InvalidFileConfig(t *testing.T) {
	defer types.SetupTest([]string{})()

	// Write an invalid configuration file
	os.WriteFile(types.GetConfigFn(), []byte("invalid_yaml"), 0644)

	// Run LoadConfig
	_, err := LoadConfig()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "failed to load file configuration")

	// Clean up
	os.Remove(types.GetConfigFn())
}

// TODO: FIX TEST
// func TestLoadConfig_EnvOverrides(t *testing.T) {
// 	defer types.SetupTest([]string{
// 		"TB_KHEDRA_CHAINS_MAINNET_RPCS=http://env.rpc1.mainnet,http://env.rpc2.mainnet",
// 		"TB_KHEDRA_SERVICES_API_PORT=9090",
// 		"TB_KHEDRA_GENERAL_DATADIR=/tmp/env-data-folder",
// 	})()

// 	// Set up a base configuration file
// 	cfg := types.Config{
// 		Chains: map[string]types.Chain{
// 			"mainnet": {RPCs: []string{"http://rpc1.mainnet"}, Enabled: true},
// 		},
// 		Services: map[string]types.Service{
// 			"api": {Port: 8080, BatchSize: 100, Enabled: true},
// 		},
// 		Logging: types.Logging{
// 			Folder:   "/tmp/test-logging-folder",
// 			Filename: "test.log",
// 			MaxSize:  50,
// 		},
// 		General: types.General{
// 			DataFolder: "/tmp/test-data-folder",
// 		},
// 	}

// 	bytes, _ := yaml.Marshal(cfg)
// 	coreFile.StringToAsciiFile(types.GetConfigFn(), string(bytes))

// 	// Run LoadConfig
// 	result, err := LoadConfig()
// 	assert.NoError(t, err)
// 	assert.Equal(t, []string{"http://env.rpc1.mainnet", "http://env.rpc2.mainnet"}, result.Chains["mainnet"].RPCs)
// 	assert.Equal(t, 9090, result.Services["api"].Port)
// 	assert.Equal(t, "/tmp/env-data-folder", result.General.DataFolder)

// 	// Clean up
// 	os.RemoveAll(cfg.Logging.Folder)
// 	os.RemoveAll(result.General.DataFolder)
// 	os.Remove(types.GetConfigFn())
// }

func TestLoadConfig_ValidationFailure(t *testing.T) {
	defer types.SetupTest([]string{})()

	// Set up an invalid configuration file
	cfg := types.Config{
		Chains: map[string]types.Chain{
			"mainnet": {RPCs: []string{}, Enabled: true},
		},
		Services: map[string]types.Service{
			"api": {Port: 8080, BatchSize: 0, Enabled: true},
		},
		Logging: types.Logging{
			Folder:   "",
			Filename: "",
			MaxSize:  0,
		},
		General: types.General{
			DataFolder: "",
		},
	}

	// Save the config to the file
	bytes, _ := yaml.Marshal(cfg)
	coreFile.StringToAsciiFile(types.GetConfigFn(), string(bytes))

	// Run LoadConfig
	_, err := LoadConfig()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "validation failed")

	// Clean up
	os.Remove(types.GetConfigFn())
}
