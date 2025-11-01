package app

import (
	"fmt"
	"os"
	"os/exec"

	coreFile "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/v5/pkg/file"
	"github.com/TrueBlocks/trueblocks-khedra/v6/pkg/types"
	"github.com/urfave/cli/v2"
)

func (k *KhedraApp) configEditAction(c *cli.Context) error {
	_ = c // linter
	fn := types.GetConfigFnNoCreate()
	if !coreFile.FileExists(fn) {
		return fmt.Errorf("not initialized you must run `khedra init` first")
	}

	editor := os.Getenv("EDITOR")
	switch editor {
	case "":
		return fmt.Errorf("EDITOR environment variable not set")
	case "testing":
		fmt.Println("Would have edited:")
		return nil
	}
	configPath := types.GetConfigFn()
	cmd := exec.Command(editor, configPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to open config for editing: %w", err)
	}
	return nil
}
