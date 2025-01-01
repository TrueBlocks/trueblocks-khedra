package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMustGetConfigFn(t *testing.T) {
	// Test fallback to create config if missing
	tmpDir := t.TempDir()

	// Mock mustGetConfigDir to point to temp directory
	oldMustGetConfigDir := mustGetConfigDir
	getConfigFn = func() string { return tmpDir }
	defer func() { getConfigFn = oldMustGetConfigDir }()

	configPath := mustGetConfigFn()

	// Verify the file is created
	assert.FileExists(t, configPath)
}

func TestMustLoadConfig_Defaults(t *testing.T) {
	// Use a temporary directory to simulate missing config
	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yaml")

	// Mock getConfigFn to return the temporary config path
	originalGetConfigFn := getConfigFn
	getConfigFn = func() string { return configFile }
	defer func() { getConfigFn = originalGetConfigFn }()

	// Establish the config file if it doesn't exist
	establishConfig(configFile)

	os.Setenv("TEST_MODE", "true")
	defer os.Unsetenv("TEST_MODE")
	cfg := MustLoadConfig(configFile)

	// Expand the expected path
	expectedFolder := expandPath("~/.khedra/logs")

	// Verify defaults
	assert.NotNil(t, cfg)
	assert.Equal(t, expectedFolder, cfg.Logging.Folder)
	assert.Equal(t, "khedra.log", cfg.Logging.Filename)
}
