package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestMustGetConfigFn(t *testing.T) {
	var configFile string
	defer types.SetupTest(t, &configFile, types.GetConfigFn, types.EstablishConfig)()
	assert.FileExists(t, configFile)
}

func TestMustLoadConfig_Defaults(t *testing.T) {
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yaml")

	// Establish the config file if it doesn't exist
	types.EstablishConfig(configFile)

	os.Setenv("TEST_MODE", "true")
	defer os.Unsetenv("TEST_MODE")
	cfg := MustLoadConfig(configFile)

	// Verify Services
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
