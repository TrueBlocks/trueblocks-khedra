package app

import (
	"fmt"

	coreFile "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/types"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/utils"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/validate"
	"github.com/goccy/go-yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// loadFileConfig loads configuration from a YAML file using the Koanf library.
// It unmarshals the file content into a types.Config object, setting default names
// for chains and services based on their keys. If the file cannot be read, parsed,
// or unmarshaled, an error is returned.
func loadFileConfig() (types.Config, error) {
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

func validateConfig(cfg types.Config) error {
	if err := validate.Validate(&cfg); err != nil {
		return err
	}
	return nil
}

// initializeFolders ensures that the specified folders in the configuration exist.
// It resolves each folder path, attempts to create missing folders, and returns
// an error if any folder cannot be created.
func initializeFolders(cfg types.Config) error {
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

type YamlComments struct{}

func MyParser() *YamlComments {
	return &YamlComments{}
}

func (p *YamlComments) Unmarshal(b []byte) (map[string]interface{}, error) {
	var out map[string]interface{}
	if err := yaml.Unmarshal(b, &out); err != nil {
		return nil, err
	}
	if out["general"] == nil {
		return out, fmt.Errorf("invalid config file: general key not found")
	}

	return out, nil
}

func (p *YamlComments) Marshal(o map[string]interface{}) ([]byte, error) {
	comments := []*yaml.Comment{{Texts: []string{"This is a file-level comment"}}}
	cm := yaml.CommentMap{
		"x": comments,
	}
	data, err := yaml.MarshalWithOptions(o, yaml.WithComment(cm))
	if err != nil {
		return nil, fmt.Errorf("failed to marshal with comment: %w", err)
	}
	return data, nil
}
