package app

import (
	"fmt"

	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
	"github.com/urfave/cli/v2"
)

func (k *KhedraApp) versionAction(c *cli.Context) error {
	_ = c // linter
	fmt.Println("khedra version " + sdk.Version())
	return nil
}
