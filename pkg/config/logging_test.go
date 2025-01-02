package config

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLogging(t *testing.T) {
	logging := NewLogging()

	// Check default values
	assert.Equal(t, "~/.khedra/logs", logging.Folder)
	assert.Equal(t, "khedra.log", logging.Filename)
	assert.Equal(t, 10, logging.MaxSizeMb)
	assert.Equal(t, 3, logging.MaxBackups)
	assert.Equal(t, 10, logging.MaxAgeDays)
	assert.True(t, logging.Compress)
}

func TestLogLevel(t *testing.T) {
	logging := NewLogging()
	assert.Equal(t, "info", logging.LogLevel)
}
