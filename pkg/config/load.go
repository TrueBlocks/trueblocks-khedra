package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	coreFile "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
)

// Declare a variable for the function that retrieves the config file path
var getConfigFn = mustGetConfigFn

func MustLoadConfig(fn string) *Config {
	var err error
	var cfg Config
	if cfg, err = loadConfig(); err != nil {
		log.Fatalf("error loading config: %v", err)
	}
	return &cfg
}

func loadConfig() (Config, error) {
	fn := getConfigFn()
	if err := k.Load(file.Provider(fn), yaml.Parser()); err != nil {
		return Config{}, fmt.Errorf("koanf.Load failed: %v", err)
	}

	cfg := NewConfig()
	if err := k.Unmarshal("", &cfg); err != nil {
		return cfg, fmt.Errorf("koanf.Unmarshal failed: %v", err)
	}
	cfg.Logging.Folder = expandPath(cfg.Logging.Folder)
	coreFile.EstablishFolder(cfg.Logging.Folder)

	return cfg, nil
}

// mustGetConfigFn returns the path to the config file which must
// be either in the current folder or in the default location. If
// there is no such file, establish it
func mustGetConfigFn() string {
	// current folder
	fn := expandPath("config.yaml")
	if coreFile.FileExists(fn) {
		return fn
	}

	// expanded default config folder
	fn = expandPath(filepath.Join(mustGetConfigDir(), "config.yaml"))
	if coreFile.FileExists(fn) {
		return fn
	}

	_ = establishConfig(fn)
	return fn
}

func mustGetConfigDir() string {
	var err error
	cfgDir := expandPath("~/.khedra")

	if !coreFile.FolderExists(cfgDir) {
		if err = coreFile.EstablishFolder(cfgDir); err != nil {
			log.Fatalf("error establishing log folder %s: %v", cfgDir, err)
		}

		if writable := IsWritable(cfgDir); !writable {
			log.Fatalf("log directory %s is not writable: %v", cfgDir, err)
		}
	}

	return cfgDir
}

// expandPath returns an absolute path expanded for ~, $HOME or other env variables
func expandPath(path string) string {
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

// IsWritable checks to see if a folder is writable
func IsWritable(path string) bool {
	tmpFile := filepath.Join(path, ".test")

	if fil, err := os.Create(tmpFile); err != nil {
		fmt.Println(fmt.Errorf("folder %s is not writable: %v", path, err))
		return false
	} else {
		fil.Close()
		if err := os.Remove(tmpFile); err != nil {
			fmt.Println(fmt.Errorf("error cleaning up test file in %s: %v", path, err))
			return false
		}
	}

	return true
}
