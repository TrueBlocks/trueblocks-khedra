package main

import (
	"github.com/TrueBlocks/trueblocks-khedra/v2/app"
)

func main() {
	k := app.NewKhedraApp()
	k.Run()
}
