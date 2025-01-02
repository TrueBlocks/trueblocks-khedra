package config

import (
	coreFile "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"gopkg.in/yaml.v2"
)

type Config struct {
	General  General   `koanf:"general" validate:"required"`           // Validate General struct
	Chains   []Chain   `koanf:"chains" validate:"required,min=1,dive"` // Validate each Chain struct
	Services []Service `koanf:"services" validate:"required,dive"`     // Validate each Service struct
	Logging  Logging   `koanf:"logging" validate:"required"`           // Validate Logging struct
}

func NewConfig() Config {
	chains := []Chain{NewChain("mainnet"), NewChain("sepolia")}
	services := []Service{
		NewService("scraper"),
		NewService("monitor"),
		NewService("api"),
		NewService("ipfs"),
	}

	return Config{
		General:  NewGeneral(),
		Chains:   chains,
		Services: services,
		Logging:  NewLogging(),
	}
}

func establishConfig(fn string) bool {
	cfg := NewConfig()
	return writeConfig(&cfg, fn)
}

func writeConfig(cfg *Config, fn string) bool {
	bytes, _ := yaml.Marshal(cfg)
	coreFile.StringToAsciiFile(fn, string(bytes))
	return coreFile.FileExists(fn)
}
