package utils

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExpandPath(t *testing.T) {
	homeDir, _ := os.UserHomeDir()
	currentDir, _ := os.Getwd()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"Expand home directory", "~/test", filepath.Join(homeDir, "test")},
		{"Expand env variable", "$HOME/test", filepath.Join(homeDir, "test")},
		{"Relative path", "./test", filepath.Join(currentDir, "test")},
		{"Absolute path", "/usr/local/test", "/usr/local/test"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExpandPath(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}
