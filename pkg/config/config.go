package config

import (
	"encoding/json"

	coreFile "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"gopkg.in/yaml.v2"
)

type General struct {
	DataPath string `koanf:"data_dir" validate:"required,path_exists,is_writable"`
	LogLevel string `koanf:"log_level" validate:"oneof=debug info warn error"`
}

type Chain struct {
	Name    string   `koanf:"name" validate:"required"`                                // Must be non-empty
	RPCs    []string `koanf:"rpcs" validate:"required,min=1,dive,strict_url,ping_one"` // Must have at least one reachable RPC URL
	Enabled bool     `koanf:"enabled"`                                                 // Defaults to false if not specified
}

type Service struct {
	Name       string `koanf:"name" validate:"required,oneof=api scraper monitor ipfs"`  // Must be non-empty
	Enabled    bool   `koanf:"enabled"`                                                  // Defaults to false if not specified
	Port       int    `koanf:"port,omitempty" validate:"opt_min=1024,opt_max=65535"`     // Must be between 1024 and 65535
	Sleep      int    `koanf:"sleep,omitempty"`                                          // Must be non-negative
	BatchSize  int    `koanf:"batch_size,omitempty" validate:"opt_min=50,opt_max=10000"` // Must be between 50 and 10000
	RetryCnt   int    `koanf:"retry_cnt,omitempty"`                                      // Must be at least 1
	RetryDelay int    `koanf:"retry_delay,omitempty"`                                    // Must be at least 1
}

type Logging struct {
	Folder     string `koanf:"folder" validate:"required,dirpath"`
	Filename   string `koanf:"filename" validate:"required,endswith=.log"`
	MaxSizeMb  int    `koanf:"max_size_mb" validate:"required,min=5"`
	MaxBackups int    `koanf:"max_backups" validate:"required,min=1"`
	MaxAgeDays int    `koanf:"max_age_days" validate:"required,min=1"`
	Compress   bool   `koanf:"compress"`
}

type Config struct {
	General  General   `koanf:"general" validate:"required"`           // Validate General struct
	Chains   []Chain   `koanf:"chains" validate:"required,min=1,dive"` // Validate each Chain struct
	Services []Service `koanf:"services" validate:"required,dive"`     // Validate each Service struct
	Logging  Logging   `koanf:"logging" validate:"required"`           // Validate Logging struct
}

func NewConfig() Config {
	return Config{
		Logging: NewLogging(),
	}
}

func (c *Config) String() string {
	bytes, _ := json.Marshal(c)
	return string(bytes)
}

func NewLogging() Logging {
	return Logging{
		Folder:     "~/.khedra/logs",
		Filename:   "khedra.log",
		MaxSizeMb:  10,
		MaxBackups: 3,
		MaxAgeDays: 10,
		Compress:   true,
	}
}

func establishConfig(fn string) bool {
	cfg := NewConfig()
	bytes, _ := yaml.Marshal(cfg)
	coreFile.StringToAsciiFile(fn, string(bytes))
	return coreFile.FileExists(fn)
}
