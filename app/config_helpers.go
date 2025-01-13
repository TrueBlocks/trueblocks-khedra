package app

import (
	"fmt"

	coreFile "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/types"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/utils"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// loadFileConfig loads configuration from a YAML file using the Koanf library.
// It unmarshals the file content into a types.Config object, setting default names
// for chains and services based on their keys. If the file cannot be read, parsed,
// or unmarshaled, an error is returned.
func loadFileConfig() (types.Config, error) {
	fileK := koanf.New(".")
	fn := types.GetConfigFn()
	if coreFile.FileSize(fn) == 0 {
		return types.Config{}, fmt.Errorf("config file is empty: %s", fn)
	}

	if err := fileK.Load(file.Provider(fn), yaml.Parser()); err != nil {
		return types.Config{}, fmt.Errorf("failed to load file config %s: %w", fn, err)
	}

	fileCfg := types.NewConfig()
	if err := fileK.Unmarshal("", &fileCfg); err != nil {
		return types.Config{}, fmt.Errorf("failed to unmarshal file config: %w", err)
	}

	for key, chain := range fileCfg.Chains {
		chain.Name = key
		fileCfg.Chains[key] = chain
	}

	for key, service := range fileCfg.Services {
		service.Name = key
		fileCfg.Services[key] = service
	}

	return fileCfg, nil
}

func validateConfig(cfg types.Config) error {
	for key, chain := range cfg.Chains {
		if len(chain.RPCs) == 0 {
			return fmt.Errorf("chain %s has no RPCs defined", key)
		}
	}

	if cfg.Logging.Folder == "" {
		return fmt.Errorf("logging folder is not defined")
	}
	if cfg.Logging.Filename == "" {
		return fmt.Errorf("logging filename is not defined")
	}
	if cfg.General.DataFolder == "" {
		return fmt.Errorf("general data directory is not defined")
	}

	return nil
}

// initializeFolders ensures that the specified folders in the configuration exist.
// It resolves each folder path, attempts to create missing folders, and returns
// an error if any folder cannot be created.
func initializeFolders(cfg types.Config) error {
	folders := []string{
		utils.ResolvePath(cfg.General.DataFolder),
		utils.ResolvePath(cfg.Logging.Folder),
	}

	for _, folder := range folders {
		if err := coreFile.EstablishFolder(folder); err != nil {
			return fmt.Errorf("failed to create folder %s: %v", folder, err)
		}
	}

	return nil
}
