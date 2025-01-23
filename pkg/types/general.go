package types

import (
	"os"
	"path/filepath"
)

// General represents configuration for data storage, ensuring the data folder is specified,
// validated for existence, and serialized for YAML-based configuration management.
type General struct {
	DataFolder string `koanf:"dataFolder" yaml:"dataFolder" validate:"required,folder_exists"`
}

func NewGeneral() General {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic("could not determine user home directory")
	}
	return General{
		DataFolder: filepath.Join(homeDir, ".khedra", "data"),
	}
}
