package main

import (
	"os"

	"github.com/TrueBlocks/trueblocks-khedra/v2/app"
)

func main() {
	// Create a new Khedra app...
	k := app.NewKhedraApp()
	k.Debug("Khedra started.")
	defer k.Debug("Khedra stopped.")

	// ...and run it
	if err := k.Run(); err != nil {
		k.Error(err.Error())
		os.Exit(1)
	}
}
