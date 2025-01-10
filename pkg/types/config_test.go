package types

import (
	"path/filepath"
	"testing"

	coreFile "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestConfigNew(t *testing.T) {
	cfg := NewConfig()

	assert.NotNil(t, cfg.General)
	assert.Equal(t, "~/.khedra/data", cfg.General.DataFolder)

	assert.NotNil(t, cfg.Chains)
	assert.Equal(t, 1, len(cfg.Chains))
	assert.NotNil(t, cfg.Chains["mainnet"])
	assert.NotNil(t, cfg.Chains["sepolia"])

	service := cfg.Chains["mainnet"]
	assert.Equal(t, "mainnet", service.Name)
	assert.Equal(t, "http://localhost:8545", service.RPCs[0])
	assert.True(t, service.Enabled)

	assert.NotNil(t, cfg.Services)
	assert.Equal(t, 4, len(cfg.Services))
	assert.NotNil(t, cfg.Services["scraper"])
	assert.NotNil(t, cfg.Services["monitor"])
	assert.NotNil(t, cfg.Services["api"])
	assert.NotNil(t, cfg.Services["cmd"])

	svc := cfg.Services["scraper"]
	assert.Equal(t, "scraper", svc.Name)
	assert.False(t, svc.Enabled)
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
	assert.False(t, svc.Enabled)
	assert.Equal(t, 8080, svc.Port)
	assert.Equal(t, 0, svc.Sleep)
	assert.Equal(t, 0, svc.BatchSize)

	svc = cfg.Services["ipfs"]
	assert.Equal(t, "ipfs", svc.Name)
	assert.False(t, svc.Enabled)
	assert.Equal(t, 5001, svc.Port)
	assert.Equal(t, 0, svc.Sleep)
	assert.Equal(t, 0, svc.BatchSize)

	assert.NotNil(t, cfg.Logging)
	assert.Equal(t, "~/.khedra/logs", cfg.Logging.Folder)
	assert.Equal(t, "khedra.log", cfg.Logging.Filename)
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
	bytes, _ := yaml.Marshal(cfg)
	coreFile.StringToAsciiFile(configFile, string(bytes))

	assert.FileExists(t, configFile)
}
