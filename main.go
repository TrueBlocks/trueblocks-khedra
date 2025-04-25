package main

import (
	"github.com/TrueBlocks/trueblocks-khedra/v5/app"
)

func main() {
	k := app.NewKhedraApp()
	k.Run()
}
