package config

import (
	coreFile "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"gopkg.in/yaml.v2"
)

type Config struct {
	General  General            `koanf:"general" validate:"required"`
	Chains   map[string]Chain   `koanf:"chains" validate:"required,min=1,dive"`
	Services map[string]Service `koanf:"services" validate:"required,min=1,dive"`
	Logging  Logging            `koanf:"logging" validate:"required"`
}

func NewConfig() Config {
	chains := map[string]Chain{
		"mainnet": NewChain("mainnet"),
		"sepolia": NewChain("sepolia"),
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

func establishConfig(fn string) bool {
	cfg := NewConfig()

	// Ensure all required fields have valid defaults
	for name, service := range cfg.Services {
		switch service.Name {
		case "scraper", "monitor":
			if service.BatchSize == 0 {
				service.BatchSize = 500 // Default BatchSize
			}
			if service.Sleep == 0 {
				service.Sleep = 10 // Default Sleep
			}
		case "api", "ipfs":
			if service.Port == 0 {
				service.Port = 8080 // Default Port
			}
		}
		cfg.Services[name] = service
	}

	return writeConfig(&cfg, fn)
}

func writeConfig(cfg *Config, fn string) bool {
	bytes, _ := yaml.Marshal(cfg)
	coreFile.StringToAsciiFile(fn, string(bytes))
	return coreFile.FileExists(fn)
}
