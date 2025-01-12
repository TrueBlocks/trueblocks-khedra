package app

import (
	"fmt"

	sdk "github.com/TrueBlocks/trueblocks-sdk/v4"
	"github.com/urfave/cli/v2"
)

func (k *KhedraApp) versionAction(c *cli.Context) error {
	_ = c // liinter
	fmt.Println("khedra version " + sdk.Version())
	return nil
}
