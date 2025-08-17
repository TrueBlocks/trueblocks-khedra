package types

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	coreFile "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/stretchr/testify/assert"
	yamlv2 "gopkg.in/yaml.v2"
)

// Testing status: reviewed

func TestConfigNew(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	assert.NoError(t, err)

	expectedDataFolder := filepath.Join(homeDir, ".khedra", "data")
	expectedLogsFolder := filepath.Join(homeDir, ".khedra", "logs")

	cfg := NewConfig()

	assert.NotNil(t, cfg.General)
	assert.Equal(t, expectedDataFolder, cfg.General.DataFolder)
	assert.Equal(t, "download", cfg.General.Strategy)
	assert.Equal(t, "index", cfg.General.Detail)

	assert.NotNil(t, cfg.Chains)
	assert.Equal(t, 1, len(cfg.Chains))
	assert.NotNil(t, cfg.Chains["mainnet"])
	assert.Equal(t, 1, cfg.Chains["mainnet"].ChainID)

	chain := cfg.Chains["mainnet"]
	assert.Equal(t, "mainnet", chain.Name)
	assert.Equal(t, "http://localhost:8545", chain.RPCs[0])
	assert.True(t, chain.Enabled)

	assert.NotNil(t, cfg.Services)
	assert.Equal(t, 4, len(cfg.Services))
	assert.NotNil(t, cfg.Services["scraper"])
	assert.NotNil(t, cfg.Services["monitor"])
	assert.NotNil(t, cfg.Services["api"])
	assert.NotNil(t, cfg.Services["ipfs"])
	if _, ok := cfg.Services["cmd"]; ok {
		t.Fatalf("unexpected legacy 'cmd' service present in default config")
	}

	svc := cfg.Services["scraper"]
	assert.Equal(t, "scraper", svc.Name)
	assert.True(t, svc.Enabled)
	assert.Equal(t, 0, svc.Port)
	assert.Equal(t, 10, svc.Sleep)
	assert.Equal(t, 500, svc.BatchSize)

	svc = cfg.Services["monitor"]
	assert.Equal(t, "monitor", svc.Name)
	assert.False(t, svc.Enabled)
	assert.Equal(t, 0, svc.Port)
	assert.Equal(t, 12, svc.Sleep)
	assert.Equal(t, 500, svc.BatchSize)

	svc = cfg.Services["api"]
	assert.Equal(t, "api", svc.Name)
	assert.True(t, svc.Enabled)
	assert.Equal(t, 8080, svc.Port)
	assert.Equal(t, 0, svc.Sleep)
	assert.Equal(t, 0, svc.BatchSize)

	svc = cfg.Services["ipfs"]
	assert.Equal(t, "ipfs", svc.Name)
	assert.True(t, svc.Enabled)
	assert.Equal(t, 5001, svc.Port)
	assert.Equal(t, 0, svc.Sleep)
	assert.Equal(t, 0, svc.BatchSize)

	assert.NotNil(t, cfg.Logging)
	assert.Equal(t, expectedLogsFolder, cfg.Logging.Folder)
	assert.Equal(t, "khedra.log", cfg.Logging.Filename)
	assert.False(t, cfg.Logging.ToFile)
	assert.Equal(t, 10, cfg.Logging.MaxSize)
	assert.Equal(t, 3, cfg.Logging.MaxBackups)
	assert.Equal(t, 10, cfg.Logging.MaxAge)
	assert.True(t, cfg.Logging.Compress)
	assert.Equal(t, "info", cfg.Logging.Level)
}

func TestConfigEstablish(t *testing.T) {
	tmpDir := (t.TempDir())
	configFile := filepath.Join(tmpDir, "config.yaml")

	cfg := NewConfig()
	bytes, _ := yamlv2.Marshal(cfg)
	_ = coreFile.StringToAsciiFile(configFile, string(bytes))

	assert.FileExists(t, configFile)
	os.Remove(configFile)
}

// Step 1: Paths and simple helpers
func TestConfig_Paths(t *testing.T) {
	cfg := NewConfig()
	tempDir := t.TempDir()
	cfg.General.DataFolder = tempDir
	if got, want := cfg.IndexPath(), filepath.Join(tempDir, "unchained"); got != want {
		t.Fatalf("IndexPath mismatch got=%s want=%s", got, want)
	}
	if got, want := cfg.CachePath(), filepath.Join(tempDir, "cache"); got != want {
		t.Fatalf("CachePath mismatch got=%s want=%s", got, want)
	}
}

// Step 2: EnabledChains (treat as set)
func TestConfig_EnabledChains(t *testing.T) {
	cfg := NewConfig()
	cfg.Chains["alt"] = NewChain("alt", 99)
	c := cfg.Chains["mainnet"]
	c.Enabled = false
	cfg.Chains["mainnet"] = c
	list := strings.Split(cfg.EnabledChains(), ",")
	have := map[string]bool{}
	for _, v := range list {
		if v != "" {
			have[v] = true
		}
	}
	if !have["alt"] {
		t.Fatalf("expected 'alt' in enabled chains: %v", list)
	}
	if have["mainnet"] {
		t.Fatalf("did not expect 'mainnet' enabled: %v", list)
	}
}

// Step 3: ServiceList variants
func TestConfig_ServiceList(t *testing.T) {
	cfg := NewConfig()
	all := strings.Split(cfg.ServiceList(false), ",")
	wantSet := map[string]bool{"api": true, "scraper": true, "monitor": true, "ipfs": true}
	for _, s := range all {
		delete(wantSet, s)
	}
	if len(wantSet) != 0 {
		t.Fatalf("missing services in full list: %v", wantSet)
	}
	enabledOnly := strings.Split(cfg.ServiceList(true), ",")
	for _, s := range enabledOnly {
		if s == "monitor" {
			t.Fatalf("monitor should not appear in enabled-only list: %v", enabledOnly)
		}
	}
}

// Step 4: RemoveZeroLines direct
func TestRemoveZeroLines(t *testing.T) {
	input := "a: 1\nzeroA: 0\nkeep: 42\nzeroB: 0\ntrailing: 5\n"
	out := RemoveZeroLines(input)
	if strings.Contains(out, "zeroA: 0") || strings.Contains(out, "zeroB: 0") {
		t.Fatalf("zero lines not removed: %s", out)
	}
	if !strings.Contains(out, "a: 1") || !strings.Contains(out, "keep: 42") || !strings.Contains(out, "trailing: 5") {
		t.Fatalf("expected retained lines missing: %s", out)
	}
}

// Step 5: WriteToFile removes ': 0' lines
func TestConfig_WriteToFile_RemoveZeroLines(t *testing.T) {
	cfg := NewConfig()
	cfg.Logging.MaxAge = 0
	cfg.Logging.MaxBackups = 0
	tempDir := t.TempDir()
	fn := filepath.Join(tempDir, "out.yaml")
	if err := cfg.WriteToFile(fn); err != nil {
		t.Fatalf("WriteToFile error: %v", err)
	}
	bytes, err := os.ReadFile(fn)
	if err != nil {
		t.Fatalf("read error: %v", err)
	}
	content := string(bytes)
	if strings.Contains(content, ": 0\n") {
		t.Fatalf("found unexpected ': 0' line in output:\n%s", content)
	}
}

func TestConfig_Version(t *testing.T) {
	cfg := NewConfig()
	v := cfg.Version()
	if strings.HasSuffix(v, "-") {
		t.Fatalf("version ends with trailing dash (from sdk.Version()): %q", v)
	}
	if strings.Contains(v, "GHC-TrueBlocks//") {
		t.Fatalf("version contains build prefix not stripped: %q", v)
	}
	if strings.Contains(v, "-release") {
		t.Fatalf("version contains '-release' suffix not stripped: %q", v)
	}
	if v == "" {
		t.Fatalf("version unexpectedly empty")
	}
}

// Optional path cleaning test: ensure trailing slashes are cleaned in output file
func TestConfig_WriteToFile_PathCleaning(t *testing.T) {
	cfg := NewConfig()
	base := t.TempDir()
	// Intentionally add trailing slashes
	cfg.General.DataFolder = filepath.Join(base, "data///")
	cfg.Logging.Folder = filepath.Join(base, "logs///")
	outFn := filepath.Join(base, "cfg.yaml")
	if err := cfg.WriteToFile(outFn); err != nil {
		t.Fatalf("WriteToFile error: %v", err)
	}
	bytes, err := os.ReadFile(outFn)
	if err != nil {
		t.Fatalf("read error: %v", err)
	}
	content := string(bytes)
	if strings.Contains(content, "data///") || strings.Contains(content, "logs///") {
		t.Fatalf("found uncleaned path with trailing slashes in output:\n%s", content)
	}
	cleanedData := filepath.Clean(filepath.Join(base, "data"))
	cleanedLogs := filepath.Clean(filepath.Join(base, "logs"))
	if !strings.Contains(content, cleanedData) {
		t.Fatalf("expected cleaned data folder path %s in output:\n%s", cleanedData, content)
	}
	if !strings.Contains(content, cleanedLogs) {
		t.Fatalf("expected cleaned logs folder path %s in output:\n%s", cleanedLogs, content)
	}
}
