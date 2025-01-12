package app

import (
	"fmt"

	_ "github.com/TrueBlocks/trueblocks-khedra/v2/pkg/env"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

func (k *KhedraApp) configShowAction(c *cli.Context) error {
	_ = c // liinter
	k.ConfigMaker()
	cfg, err := LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}
	bytes, err := yaml.Marshal(&cfg)
	if err != nil {
		return fmt.Errorf("failed to unmarshal config: %w", err)
	}
	fmt.Println(string(bytes))
	return nil
}
