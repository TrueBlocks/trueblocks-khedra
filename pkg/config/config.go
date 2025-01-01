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
	return Config{
		General: NewGeneral(),
		Chains:  []Chain{NewChain()},
		Services: []Service{
			NewService("scraper"),
			NewService("monitor"),
			NewService("api"),
			NewService("ipfs"),
		},
		Logging: NewLogging(),
	}
}

func establishConfig(fn string) bool {
	cfg := NewConfig()
	bytes, _ := yaml.Marshal(cfg)
	coreFile.StringToAsciiFile(fn, string(bytes))
	return coreFile.FileExists(fn)
}
