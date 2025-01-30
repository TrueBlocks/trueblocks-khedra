package types

import (
	"bufio"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	coreFile "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	_ "github.com/TrueBlocks/trueblocks-khedra/v2/pkg/env"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/utils"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v4"
)

type Config struct {
	General  General            `koanf:"general" validate:"dive"`
	Chains   map[string]Chain   `koanf:"chains" validate:"dive"`
	Services map[string]Service `koanf:"services" validate:"dive"`
	Logging  Logging            `koanf:"logging" validate:"dive"`
}

func NewConfig() Config {
	chains := map[string]Chain{
		"mainnet": NewChain("mainnet", 1),
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

func (c *Config) Version() string {
	v := sdk.Version()
	v = strings.Replace(v, "GHC-TrueBlocks//", "", 1)
	v = strings.Replace(v, "-release", "", 1)
	return v
}

func (c *Config) IndexPath() string {
	return filepath.Join(c.General.DataFolder, "unchained")
}

func (c *Config) CachePath() string {
	return filepath.Join(c.General.DataFolder, "cache")
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
	err := cfg.WriteToFile(fn)
	if err != nil {
		fmt.Println(colors.Red+"error writing config file: %v", err, colors.Off)
	}

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

func (c *Config) EnabledChains() string {
	var ret []string
	for key, ch := range c.Chains {
		if ch.Enabled {
			ret = append(ret, key)
		}
	}
	return strings.Join(ret, ",")
}

func (c *Config) ServiceList(enabledOnly bool) string {
	var ret []string
	for key, svc := range c.Services {
		if !enabledOnly || svc.Enabled {
			ret = append(ret, key)
		}
	}
	return strings.Join(ret, ",")
}

// WriteToFile writes the Config struct to a file using a YAML template with comments.
func (c *Config) WriteToFile(fn string) error {
	c.General.DataFolder = filepath.Clean(c.General.DataFolder)
	c.Logging.Folder = filepath.Clean(c.Logging.Folder)

	t, err := template.New("config").Parse(strings.TrimSpace(tmpl) + "\n")
	if err != nil {
		return err
	}

	var builder strings.Builder
	if err := t.Execute(&builder, c); err != nil {
		return err
	}
	processed := RemoveZeroLines(builder.String())

	f, err := os.Create(fn)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(processed)
	return err
}

func RemoveZeroLines(input string) string {
	var builder strings.Builder
	scanner := bufio.NewScanner(strings.NewReader(input))

	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasSuffix(line, ": 0") {
			builder.WriteString(line + "\n")
		}
	}

	return builder.String()
}

const tmpl = `
# Khedra Configuration File
#
# For more information see the users manual at https://khedra.trueblocks.io
#
# You may easily edit this file with "khedra config edit" or by typing
# "edit" on the "khedra init" command line.
#
# You may add as many chains as you wish. The Mainnet RPC is required even
# though you may disable it from processing. Additional chains require
# a working RPC endpoint. The file will be validated when loaded.
#
# Note that any comments you write in this file will be overwritten.

general:
  dataFolder: "{{ .General.DataFolder }}"
  strategy: "{{ .General.Strategy }}"
  detail: "{{ .General.Detail }}"

chains:
{{- range $key, $value := .Chains }}
  {{ $key }}:
    name: "{{ $value.Name }}"
    rpcs: 
{{- range $rpc := $value.RPCs }}
      - "{{ $rpc }}"
{{- end }}
    enabled: {{ $value.Enabled }}
    chainId: {{ $value.ChainID }}
{{- end }}

services:
{{- range $key, $value := .Services }}
  {{ $key }}:
    name: "{{ $value.Name }}"
    enabled: {{ $value.Enabled }}
    port: {{ $value.Port }}
    sleep: {{ $value.Sleep }}
    batchSize: {{ $value.BatchSize }}
{{- end }}

logging:
  folder: "{{ .Logging.Folder }}"
  filename: "{{ .Logging.Filename }}"
  toFile: {{ .Logging.ToFile }}
  maxSize: {{ .Logging.MaxSize }}
  maxBackups: {{ .Logging.MaxBackups }}
  maxAge: {{ .Logging.MaxAge }}
  compress: {{ .Logging.Compress }}
  level: "{{ .Logging.Level }}"
`
