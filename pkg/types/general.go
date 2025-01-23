package types

import (
	"os"
	"path/filepath"
)

// General represents configuration for data storage, ensuring the data folder is specified,
// validated for existence, and serialized for YAML-based configuration management.
type General struct {
	DataFolder       string `koanf:"dataFolder" yaml:"dataFolder" validate:"required,folder_exists"`
	DownloadStrategy string `koanf:"downloadStrategy" yaml:"downloadStrategy"`
	DownloadDetail   string `koanf:"downloadDetail" yaml:"downloadDetail"`
}

func NewGeneral() General {
	return General{
		DataFolder:       getDefaultDataFolder(),
		DownloadStrategy: getDefaultDownloadStrategy(),
		DownloadDetail:   getDefaultDownloadDetail(),
	}
}

func getDefaultDataFolder() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		panic("could not determine user home directory")
	}
	return filepath.Join(homeDir, ".khedra", "data")
}

func getDefaultDownloadStrategy() string {
	return "download"
}

func getDefaultDownloadDetail() string {
	return "entire index"
}
