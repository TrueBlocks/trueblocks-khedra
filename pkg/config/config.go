package config

import (
	"encoding/json"

	coreFile "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/knadh/koanf/v2"
	"gopkg.in/yaml.v2"
)

var k = koanf.New(".")

type Config struct {
	Logging Logging `koanf:"logging"`
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

func establishConfig(fn string) bool {
	cfg := NewConfig()
	bytes, _ := yaml.Marshal(cfg)
	coreFile.StringToAsciiFile(fn, string(bytes))
	return coreFile.FileExists(fn)
}
