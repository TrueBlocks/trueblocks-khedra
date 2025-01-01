package main

import (
	"os"

	"github.com/TrueBlocks/trueblocks-khedra/v2/app"
)

func main() {
	// Create a new Khedra app...
	k := app.NewKhedraApp()

	k.FileLogger.Info("Khedra started.")
	defer k.FileLogger.Info("Khedra stopped.")

	// ...and run it
	if err := k.Run(); err != nil {
		k.ProgLogger.Error(err.Error())
		os.Exit(1)
	}
}
