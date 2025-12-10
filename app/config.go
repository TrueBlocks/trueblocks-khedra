package app

import (
	"fmt"

	"github.com/TrueBlocks/trueblocks-khedra/v6/pkg/install"
	"github.com/TrueBlocks/trueblocks-khedra/v6/pkg/types"
)

func LoadConfig() (types.Config, error) {
	loader := NewConfigLoader()
	return loader.Load()
}

func (k *KhedraApp) loadConfigIfInitialized() error {
	if !install.Configured() {
		return fmt.Errorf("not initialized you must run `khedra init` first")
	}

	if _, err := k.ConfigMaker(); err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	return nil
}
