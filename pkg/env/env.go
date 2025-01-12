package env

import (
	"fmt"
	"os"
	"path/filepath"

	coreFile "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/joho/godotenv"
)

func init() {
	if pwd, err := os.Getwd(); err == nil {
		if coreFile.FileExists(filepath.Join(pwd, ".env")) {
			if err = godotenv.Load(filepath.Join(pwd, ".env")); err != nil {
				fmt.Fprintf(os.Stderr, "Found .env, but could not read it\n")
			}
		}
	}
}
