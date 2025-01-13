package utils

import (
	"log"
	"os"
	"path/filepath"
	"strings"
)

// ResolvePath returns an absolute path expanded for ~, $HOME or other env variables
func ResolvePath(path string) string {
	if path == "" {
		log.Fatalf("path cannot be empty")
	}

	if strings.HasPrefix(path, "~") {
		if path == "~" || strings.HasPrefix(path, "~/") {
			home, err := os.UserHomeDir()
			if err != nil {
				log.Fatalf("failed to resolve home directory: %v", err)
			}
			path = filepath.Join(home, strings.TrimPrefix(path, "~"))
		} else {
			log.Fatalf("unsupported path format: %s", path)
		}
	}

	for _, part := range strings.Split(path, "/") {
		if strings.HasPrefix(part, "$") {
			envVar := strings.TrimPrefix(part, "$")
			if os.Getenv(envVar) == "" {
				log.Fatalf("path contains unset environment variable: %s", part)
			}
		}
	}

	path = os.ExpandEnv(path)

	absolutePath, err := filepath.Abs(path)
	if err != nil {
		log.Fatalf("failed to resolve absolute path: %v", err)
	}

	return absolutePath
}
