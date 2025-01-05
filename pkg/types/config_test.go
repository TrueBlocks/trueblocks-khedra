package types

import (
	"path/filepath"
	"testing"

	coreFile "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestNewConfig(t *testing.T) {
	cfg := NewConfig()

	// Check default values
	assert.NotNil(t, cfg.Logging)
	assert.Equal(t, "~/.khedra/logs", cfg.Logging.Folder)
	assert.Equal(t, "khedra.log", cfg.Logging.Filename)
	assert.Equal(t, 10, cfg.Logging.MaxSizeMb)
	assert.Equal(t, 3, cfg.Logging.MaxBackups)
	assert.Equal(t, 10, cfg.Logging.MaxAgeDays)
	assert.True(t, cfg.Logging.Compress)
}

func TestEstablishConfig(t *testing.T) {
	tmpDir := (t.TempDir())
	configFile := filepath.Join(tmpDir, "config.yaml")

	cfg := NewConfig()
	bytes, _ := yaml.Marshal(cfg)
	coreFile.StringToAsciiFile(configFile, string(bytes))
	created := coreFile.FileExists(configFile)

	// Verify the file is created and exists
	assert.True(t, created)
	assert.FileExists(t, configFile)
}
