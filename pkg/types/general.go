package types

import (
	"os"
	"path/filepath"

	"github.com/TrueBlocks/trueblocks-chifra/v6/pkg/logger"
)

// General represents configuration for data storage, ensuring the data folder is specified,
// validated for existence, and serialized for YAML-based configuration management.
type General struct {
	DataFolder     string `koanf:"dataFolder" yaml:"dataFolder" json:"dataFolder,omitempty" validate:"required,folder_exists"`
	MonitorsFolder string `koanf:"monitorsFolder" yaml:"monitorsFolder" json:"monitorsFolder,omitempty" validate:"required"`
	Strategy       string `koanf:"strategy" yaml:"strategy" json:"strategy,omitempty" validate:"oneof=download scratch"`
	Detail         string `koanf:"detail" yaml:"detail" json:"detail,omitempty" validate:"oneof=index bloom"`
}

func NewGeneral() General {
	return General{
		DataFolder:     getDefaultDataFolder(),
		MonitorsFolder: getDefaultMonitorsFolder(),
		Strategy:       getDefaultStrategy(),
		Detail:         getDefaultDetail(),
	}
}

func getDefaultDataFolder() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		logger.Panic("could not determine user home directory")
	}
	return filepath.Join(homeDir, ".khedra", "data")
}

func getDefaultMonitorsFolder() string {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		logger.Panic("could not determine user home directory")
	}
	return filepath.Join(homeDir, ".khedra", "monitors")
}

func getDefaultStrategy() string {
	return "download"
}

func getDefaultDetail() string {
	return "index"
}
