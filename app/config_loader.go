package app

import (
	"fmt"
	"path/filepath"

	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"

	"github.com/TrueBlocks/trueblocks-chifra/v6/pkg/base"
	coreFile "github.com/TrueBlocks/trueblocks-chifra/v6/pkg/file"
	"github.com/TrueBlocks/trueblocks-chifra/v6/pkg/utils"
	"github.com/TrueBlocks/trueblocks-khedra/v6/pkg/types"
)

type ConfigLoader struct{}

func NewConfigLoader() *ConfigLoader {
	return &ConfigLoader{}
}

func (cl *ConfigLoader) Load() (types.Config, error) {
	cfg, err := cl.loadFromFile()
	if err != nil {
		return types.Config{}, fmt.Errorf("failed to load file configuration: %w", err)
	}

	if err := cl.applyEnvironment(&cfg); err != nil {
		return types.Config{}, fmt.Errorf("failed to apply environment configuration: %w", err)
	}

	if err := cl.cleanup(&cfg); err != nil {
		return types.Config{}, fmt.Errorf("failed to finalize configuration: %w", err)
	}

	if err := cl.validate(cfg); err != nil {
		return types.Config{}, fmt.Errorf("validation error: %w", err)
	}

	if err := cl.initializeFolders(cfg); err != nil {
		return types.Config{}, fmt.Errorf("failed to initialize folders: %w", err)
	}

	return cfg, nil
}

func (cl *ConfigLoader) loadFromFile() (types.Config, error) {
	fileK := koanf.New(".")
	fn := types.GetConfigFn()
	if coreFile.FileSize(fn) == 0 || len(coreFile.AsciiFileToString(fn)) == 0 {
		return types.Config{}, fmt.Errorf("config file is empty: %s", fn)
	}

	if err := fileK.Load(file.Provider(fn), MyParser()); err != nil {
		return types.Config{}, fmt.Errorf("failed to load file config %s: %w", fn, err)
	}

	fileCfg := types.NewConfig()
	if err := fileK.Unmarshal("", &fileCfg); err != nil {
		return types.Config{}, fmt.Errorf("failed to unmarshal file config: %w", err)
	}

	for key, chain := range fileCfg.Chains {
		chain.Name = key
		fileCfg.Chains[key] = chain
	}

	for key, service := range fileCfg.Services {
		service.Name = key
		fileCfg.Services[key] = service
	}

	return fileCfg, nil
}

func (cl *ConfigLoader) applyEnvironment(cfg *types.Config) error {
	keys := types.GetEnvironmentKeys(*cfg, types.InEnv)
	return types.ApplyEnv(keys, cfg)
}

func (cl *ConfigLoader) cleanup(cfg *types.Config) error {
	cfg.General.DataFolder = filepath.Clean(cfg.General.DataFolder)
	cfg.Logging.Folder = filepath.Clean(cfg.Logging.Folder)
	return nil
}

func (cl *ConfigLoader) validate(cfg types.Config) error {
	if err := types.Validate(&cfg); err != nil {
		return err
	}

	if base.IsTestMode() {
		return nil
	}

	svcList := cfg.ServiceList(true)
	if len(svcList) == 0 {
		return fmt.Errorf("at least one service must be enabled")
	}
	chList := cfg.EnabledChains()
	if len(chList) == 0 {
		return fmt.Errorf("at least one chain must be enabled")
	}
	if ch, ok := cfg.Chains["mainnet"]; !ok || len(ch.RPCs) == 0 {
		return fmt.Errorf("mainnet RPC must be provided")
	}

	return nil
}

func (cl *ConfigLoader) initializeFolders(cfg types.Config) error {
	folders := []string{
		utils.ResolvePath(cfg.General.DataFolder),
		utils.ResolvePath(cfg.Logging.Folder),
	}

	for _, folder := range folders {
		if err := coreFile.EstablishFolder(folder); err != nil {
			return fmt.Errorf("failed to create folder %s: %v", folder, err)
		}
	}

	return nil
}
