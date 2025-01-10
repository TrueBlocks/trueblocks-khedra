package main

import (
	"os"

	"github.com/TrueBlocks/trueblocks-khedra/v2/app"
)

func main() {
	app.
		NewKhedraApp().
		Run(os.Args)
}
