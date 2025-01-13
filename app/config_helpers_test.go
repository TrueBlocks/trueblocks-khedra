package app

import (
	"fmt"
	"os"
	"testing"

	coreFile "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/types"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

// Testing status: not_reviewed

// ---------------------------------------------------------
func TestLoadFileConfig(t *testing.T) {
	invalidFile := func() {
		defer types.SetupTest([]string{})()
		coreFile.StringToAsciiFile(types.GetConfigFn(), "invalid: [:::]")

		_, err := loadFileConfig()
		assert.Error(t, err)
	}
	t.Run("Invalid File", func(t *testing.T) { invalidFile() })

	validFile := func() {
		defer types.SetupTest([]string{})()

		cfg := types.NewConfig()
		chain := cfg.Chains["mainnet"]
		chain.RPCs = []string{"http://localhost:8545", "http://localhost:8546"}
		cfg.Chains["mainnet"] = chain
		bytes, _ := yaml.Marshal(cfg)
		coreFile.StringToAsciiFile(types.GetConfigFn(), string(bytes))
		// fmt.Println(string(bytes))

		result, err := loadFileConfig()
		assert.NoError(t, err)
		assert.Equal(t, cfg, result)
	}
	t.Run("Valid File", func(t *testing.T) { validFile() })

	missingFile := func() {
		defer types.SetupTest([]string{})()
		os.Remove(types.GetConfigFn())

		_, err := loadFileConfig()
		assert.Error(t, err)
	}
	t.Run("Missing File", func(t *testing.T) { missingFile() })

	emptyFile := func() {
		defer types.SetupTest([]string{})()
		coreFile.StringToAsciiFile(types.GetConfigFn(), "")
		result, err := loadFileConfig()
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "config file is empty")
		fmt.Println(result)
	}
	t.Run("Empty File", func(t *testing.T) { emptyFile() })
}

// ---------------------------------------------------------
func TestValidateConfig(t *testing.T) {
	validConfig := func() {
		cfg := types.Config{
			Chains: map[string]types.Chain{
				"mainnet": {RPCs: []string{"http://rpc1.mainnet"}, Enabled: true},
			},
			Services: map[string]types.Service{
				"api": {Port: 8080, BatchSize: 100, Enabled: true},
			},
			Logging: types.Logging{
				Folder:   "~/.khedra",
				Filename: "khedra.log",
				MaxSize:  50,
			},
			General: types.General{
				DataFolder: "~/.khedra/data",
			},
		}

		err := validateConfig(cfg)
		assert.NoError(t, err)
	}
	t.Run("Valid Config", func(t *testing.T) { validConfig() })

	missingRPCs := func() {
		cfg := types.Config{
			Chains: map[string]types.Chain{
				"mainnet": {RPCs: []string{}, Enabled: true},
			},
		}

		err := validateConfig(cfg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "chain mainnet has no RPCs defined")
	}
	t.Run("Missing RPCs", func(t *testing.T) { missingRPCs() })

	invalidLoggingFolder := func() {
		cfg := types.NewConfig()
		cfg.Logging.Folder = ""

		err := validateConfig(cfg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "logging folder is not defined")
	}
	t.Run("Invalid Logging Folder", func(t *testing.T) { invalidLoggingFolder() })

	missingLoggingConfig := func() {
		cfg := types.Config{
			Logging: types.Logging{
				Folder:   "",
				Filename: "",
			},
		}

		err := validateConfig(cfg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "logging folder is not defined")
	}
	t.Run("Missing Logging Config", func(t *testing.T) { missingLoggingConfig() })

	missingGeneralConfig := func() {
		cfg := types.Config{
			General: types.General{
				DataFolder: "",
			},
		}

		err := validateConfig(cfg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "logging folder is not defined")
	}
	t.Run("Missing General Config", func(t *testing.T) { missingGeneralConfig() })
}

// ---------------------------------------------------------
func TestInitializeFolders(t *testing.T) {
	cleanup := func(cfg types.Config) {
		os.RemoveAll(cfg.Logging.Folder)
		os.RemoveAll(cfg.General.DataFolder)
	}

	allFoldersExist := func() {
		cfg := types.Config{
			Logging: types.Logging{
				Folder: "/tmp/test-logging-folder",
			},
			General: types.General{
				DataFolder: "/tmp/test-data-folder",
			},
		}

		os.MkdirAll(cfg.Logging.Folder, os.ModePerm)
		os.MkdirAll(cfg.General.DataFolder, os.ModePerm)

		err := initializeFolders(cfg)
		assert.NoError(t, err)

		cleanup(cfg)
	}
	t.Run("All Folders Exist", func(t *testing.T) { allFoldersExist() })

	createMissingFolders := func() {
		cfg := types.Config{
			Logging: types.Logging{
				Folder: "/tmp/test-missing-logging-folder",
			},
			General: types.General{
				DataFolder: "/tmp/test-missing-data-folder",
			},
		}

		cleanup(cfg)

		err := initializeFolders(cfg)
		assert.NoError(t, err)

		_, err = os.Stat(cfg.Logging.Folder)
		assert.NoError(t, err)
		_, err = os.Stat(cfg.General.DataFolder)
		assert.NoError(t, err)

		cleanup(cfg)
	}
	t.Run("Create Missing Folders", func(t *testing.T) { createMissingFolders() })

	errorOnInvalidPath := func() {
		cfg := types.Config{
			Logging: types.Logging{
				Folder: "/invalid-folder-path/\\0",
			},
			General: types.General{
				DataFolder: "/tmp/test-data-folder",
			},
		}

		err := initializeFolders(cfg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "failed to create folder")

		cleanup(cfg)
	}
	t.Run("Error On Invalid Path", func(t *testing.T) { errorOnInvalidPath() })
}
