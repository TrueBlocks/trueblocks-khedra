package app

import (
	"fmt"

	coreFile "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/types"
)

func LoadConfig() (types.Config, error) {
	cfg, err := loadFileConfig()
	if err != nil {
		return types.Config{}, fmt.Errorf("failed to load file configuration: %w", err)
	}
	keys := types.GetEnvironmentKeys(cfg, types.InEnv)
	if err := types.ApplyEnv(keys, &cfg); err != nil {
		return types.Config{}, fmt.Errorf("failed to apply environment configuration: %w", err)
	}

	if err := finalCleanup(&cfg); err != nil {
		return types.Config{}, fmt.Errorf("failed to finalize configuration: %w", err)
	}

	if err := validateConfig(cfg); err != nil {
		return types.Config{}, fmt.Errorf("validation error: %w", err)
	}

	if err := initializeFolders(cfg); err != nil {
		return types.Config{}, fmt.Errorf("failed to initialize folders: %w", err)
	}

	return cfg, nil
}

func (k *KhedraApp) loadConfigIfInitialized() error {
	fn := types.GetConfigFnNoCreate()
	if !coreFile.FileExists(fn) {
		return fmt.Errorf("not initialized you must run `khedra init` first")
	}

	if _, err := k.ConfigMaker(); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	return nil
}
