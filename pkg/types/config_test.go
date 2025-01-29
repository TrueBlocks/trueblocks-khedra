package types

import (
	"os"
	"path/filepath"
	"testing"

	coreFile "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/stretchr/testify/assert"
	yamlv2 "gopkg.in/yaml.v2"
)

// Testing status: reviewed

func TestConfigNew(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	assert.NoError(t, err)

	expectedDataFolder := filepath.Join(homeDir, ".khedra", "data")
	expectedLogsFolder := filepath.Join(homeDir, ".khedra", "logs")

	cfg := NewConfig()

	assert.NotNil(t, cfg.General)
	assert.Equal(t, expectedDataFolder, cfg.General.DataFolder)
	assert.Equal(t, "download", cfg.General.Strategy)
	assert.Equal(t, "index", cfg.General.Detail)

	assert.NotNil(t, cfg.Chains)
	assert.Equal(t, 1, len(cfg.Chains))
	assert.NotNil(t, cfg.Chains["mainnet"])
	assert.Equal(t, 1, cfg.Chains["mainnet"].ChainId)

	chain := cfg.Chains["mainnet"]
	assert.Equal(t, "mainnet", chain.Name)
	assert.Equal(t, "http://localhost:8545", chain.RPCs[0])
	assert.True(t, chain.Enabled)

	assert.NotNil(t, cfg.Services)
	assert.Equal(t, 4, len(cfg.Services))
	assert.NotNil(t, cfg.Services["scraper"])
	assert.NotNil(t, cfg.Services["monitor"])
	assert.NotNil(t, cfg.Services["api"])
	assert.NotNil(t, cfg.Services["cmd"])

	svc := cfg.Services["scraper"]
	assert.Equal(t, "scraper", svc.Name)
	assert.True(t, svc.Enabled)
	assert.Equal(t, 0, svc.Port)
	assert.Equal(t, 10, svc.Sleep)
	assert.Equal(t, 500, svc.BatchSize)

	svc = cfg.Services["monitor"]
	assert.Equal(t, "monitor", svc.Name)
	assert.False(t, svc.Enabled)
	assert.Equal(t, 0, svc.Port)
	assert.Equal(t, 12, svc.Sleep)
	assert.Equal(t, 500, svc.BatchSize)

	svc = cfg.Services["api"]
	assert.Equal(t, "api", svc.Name)
	assert.True(t, svc.Enabled)
	assert.Equal(t, 8080, svc.Port)
	assert.Equal(t, 0, svc.Sleep)
	assert.Equal(t, 0, svc.BatchSize)

	svc = cfg.Services["ipfs"]
	assert.Equal(t, "ipfs", svc.Name)
	assert.True(t, svc.Enabled)
	assert.Equal(t, 5001, svc.Port)
	assert.Equal(t, 0, svc.Sleep)
	assert.Equal(t, 0, svc.BatchSize)

	assert.NotNil(t, cfg.Logging)
	assert.Equal(t, expectedLogsFolder, cfg.Logging.Folder)
	assert.Equal(t, "khedra.log", cfg.Logging.Filename)
	assert.False(t, cfg.Logging.ToFile)
	assert.Equal(t, 10, cfg.Logging.MaxSize)
	assert.Equal(t, 3, cfg.Logging.MaxBackups)
	assert.Equal(t, 10, cfg.Logging.MaxAge)
	assert.True(t, cfg.Logging.Compress)
	assert.Equal(t, "info", cfg.Logging.Level)
}

func TestConfigEstablish(t *testing.T) {
	tmpDir := (t.TempDir())
	configFile := filepath.Join(tmpDir, "config.yaml")

	cfg := NewConfig()
	bytes, _ := yamlv2.Marshal(cfg)
	coreFile.StringToAsciiFile(configFile, string(bytes))

	assert.FileExists(t, configFile)
	os.Remove(configFile)
}
