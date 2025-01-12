package utils

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

// ResolvePath returns an absolute path expanded for ~, $HOME or other env variables
func ResolvePath(path string) string {
	if strings.HasPrefix(path, "~/") || path == "~" {
		home, err := os.UserHomeDir()
		if err != nil {
			log.Fatalf("error getting home directory: %v", err)
		}
		path = filepath.Join(home, strings.TrimPrefix(path, "~"))
	}

	var err error
	path, err = filepath.Abs(os.ExpandEnv(path))
	if err != nil {
		log.Fatalf("error making path absolute: %v", err)
	}

	return path
}
