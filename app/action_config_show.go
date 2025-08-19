package app

import (
	"fmt"

	coreFile "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-khedra/v5/pkg/types"
	"github.com/urfave/cli/v2"
	yamlv2 "gopkg.in/yaml.v2"
)

func (k *KhedraApp) configShowAction(c *cli.Context) error {
	_ = c // linter
	fn := types.GetConfigFnNoCreate()
	if !coreFile.FileExists(fn) {
		return fmt.Errorf("not initialized you must run `khedra init` first")
	}

	cfg, err := LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	bytes, err := yamlv2.Marshal(&cfg)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}
	fmt.Println(string(bytes))
	return nil
}
