package types

import (
	"bufio"
	"fmt"
	"html/template"
	"os"
	"path/filepath"
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/base"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	coreFile "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
	sdk "github.com/TrueBlocks/trueblocks-sdk/v5"
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
	if testFn, ok := os.LookupEnv("KHEDRA_TEST_CONFIG_FN"); ok && testFn != "" {
		return testFn
	}
	if base.IsTestMode() {
		return filepath.Join(os.TempDir(), "config.yaml")
	}
	fn := utils.ResolvePath("config.yaml")
	if coreFile.FileExists(fn) {
		return fn
	}
	return utils.ResolvePath(filepath.Join(mustGetConfigPath(), "config.yaml"))
}

func GetConfigFn() string {
	if testFn, ok := os.LookupEnv("KHEDRA_TEST_CONFIG_FN"); ok && testFn != "" {
		if coreFile.FileExists(testFn) {
			return testFn
		}
		cfg := NewConfig()
		if err := cfg.WriteToFile(testFn); err != nil {
			fmt.Println(colors.Red+"error writing config file: %v", err, colors.Off)
		}
		return testFn
	}
	if base.IsTestMode() {
		return filepath.Join(os.TempDir(), "config.yaml")
	}
	fn := GetConfigFnNoCreate()
	if coreFile.FileExists(fn) {
		return fn
	}
	cfg := NewConfig()
	if err := cfg.WriteToFile(fn); err != nil {
		fmt.Println(colors.Red+"error writing config file: %v", err, colors.Off)
	}
	return fn
}

func mustGetConfigPath() string {
	var err error
	cfgDir := utils.ResolvePath("~/.khedra")
	if !coreFile.FolderExists(cfgDir) {
		if err = coreFile.EstablishFolder(cfgDir); err != nil {
			logger.Panicf("error establishing log folder %s: %v", cfgDir, err)
		}
	}
	if !isWritable(cfgDir) {
		logger.Panicf("log directory %s is not writable: %v", cfgDir, err)
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
		// Try to clean up test file, but don't fail if cleanup fails
		os.Remove(tmpFile) // Ignore error - file may already be gone or have permission issues
	}
	return true
}

func (c *Config) EnabledChains() string {
	var ret []string
	for k, ch := range c.Chains {
		if ch.Enabled {
			ret = append(ret, k)
		}
	}
	return strings.Join(ret, ",")
}

func (c *Config) ServiceList(enabledOnly bool) string {
	var ret []string
	for k, svc := range c.Services {
		if !enabledOnly || svc.Enabled {
			ret = append(ret, k)
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
	var b strings.Builder
	sc := bufio.NewScanner(strings.NewReader(input))
	for sc.Scan() {
		line := sc.Text()
		if !strings.HasSuffix(line, ": 0") {
			b.WriteString(line + "\n")
		}
	}
	return b.String()
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

// ConfigTemplate returns the YAML template for the config file.
func ConfigTemplate() string {
	return tmpl
}
