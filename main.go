package main

import (
	"fmt"
	"os"
	"path/filepath"

	coreFile "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/v5/pkg/file"
	"github.com/TrueBlocks/trueblocks-khedra/v6/app"
	"github.com/joho/godotenv"
)

func main() {
	k := app.NewKhedraApp()
	k.Run()
}

func init() {
	if pwd, err := os.Getwd(); err == nil {
		if coreFile.FileExists(filepath.Join(pwd, ".env")) {
			if err = godotenv.Load(filepath.Join(pwd, ".env")); err != nil {
				fmt.Fprintf(os.Stderr, "Found .env, but could not read it\n")
			}
		}
	}
}
