package config

import (
	"testing"

	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestConfigMustLoad(t *testing.T) {
	var configFile string
	defer types.SetupTestOld(t, &configFile, types.GetConfigFn, types.EstablishConfig)()
	assert.FileExists(t, configFile)
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
