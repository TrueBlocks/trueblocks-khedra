package config

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
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
	// Use a temporary directory to test config creation
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yaml")

	created := establishConfig(configFile)

	// Verify the file is created and exists
	assert.True(t, created)
	assert.FileExists(t, configFile)
}
