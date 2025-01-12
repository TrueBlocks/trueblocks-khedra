package app

import (
	"fmt"
	"os"

	_ "github.com/TrueBlocks/trueblocks-khedra/v2/pkg/env"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/types"
	"github.com/urfave/cli/v2"
)

func (k *KhedraApp) configEditAction(c *cli.Context) error {
	_ = c // liinter
	k.ConfigMaker()
	editor := os.Getenv("EDITOR")
	if editor == "" {
		return fmt.Errorf("EDITOR environment variable not set")
	} else if editor == "testing" {
		fmt.Println("Would have edited:")
		return nil
	}
	configPath := types.GetConfigFn()
	cmd := execCommand(editor, configPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to open config for editing: %w", err)
	}
	return nil
}
