package types

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	coreFile "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	_ "github.com/TrueBlocks/trueblocks-khedra/v2/pkg/env"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/utils"
	"gopkg.in/yaml.v2"
)

type Config struct {
	General  General            `koanf:"general" validate:"dive"`
	Chains   map[string]Chain   `koanf:"chains" validate:"dive"`
	Services map[string]Service `koanf:"services" validate:"dive"`
	Logging  Logging            `koanf:"logging" validate:"dive"`
}

func NewConfig() Config {
	chains := map[string]Chain{
		"mainnet": NewChain("mainnet"),
	}
	services := map[string]Service{
		"scraper": NewService("scraper"),
		"monitor": NewService("monitor"),
		"api":     NewService("api"),
		"ipfs":    NewService("ipfs"),
	}

	return Config{
		General:  NewGeneral(),
		Chains:   chains,
		Services: services,
		Logging:  NewLogging(),
	}
}

func GetConfigFnNoCreate() string {
	if os.Getenv("TEST_MODE") == "true" {
		tmpDir := os.TempDir()
		return filepath.Join(tmpDir, "config.yaml")
	}

	// current folder
	fn := utils.ResolvePath("config.yaml")
	if coreFile.FileExists(fn) {
		return fn
	}

	// expanded default config folder
	return utils.ResolvePath(filepath.Join(mustGetConfigPath(), "config.yaml"))
}

// GetConfigFn returns the path to the config file which must
// be either in the current folder or in the default location. If
// there is no such file, establish it
func GetConfigFn() string {
	if os.Getenv("TEST_MODE") == "true" {
		tmpDir := os.TempDir()
		return filepath.Join(tmpDir, "config.yaml")
	}

	fn := GetConfigFnNoCreate()
	if coreFile.FileExists(fn) {
		return fn
	}

	cfg := NewConfig()
	bytes, _ := yaml.Marshal(cfg)
	coreFile.StringToAsciiFile(fn, string(bytes))

	return fn
}

func mustGetConfigPath() string {
	var err error
	cfgDir := utils.ResolvePath("~/.khedra")

	if !coreFile.FolderExists(cfgDir) {
		if err = coreFile.EstablishFolder(cfgDir); err != nil {
			log.Fatalf("error establishing log folder %s: %v", cfgDir, err)
		}
	}

	if writable := isWritable(cfgDir); !writable {
		log.Fatalf("log directory %s is not writable: %v", cfgDir, err)
	}

	return cfgDir
}

// isWritable checks to see if a folder is writable
func isWritable(path string) bool {
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
