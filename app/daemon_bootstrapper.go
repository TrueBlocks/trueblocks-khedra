package app

import (
	"bytes"
	"fmt"
	"html/template"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/TrueBlocks/trueblocks-chifra/v6/pkg/file"
	"github.com/TrueBlocks/trueblocks-khedra/v6/pkg/types"
)

// DaemonBootstrapper handles daemon initialization phases
type DaemonBootstrapper struct {
	config     *types.Config
	rootFolder string
	logger     *types.CustomLogger
}

// NewDaemonBootstrapper creates a new bootstrapper instance
func NewDaemonBootstrapper(config *types.Config, rootFolder string, logger *types.CustomLogger) *DaemonBootstrapper {
	return &DaemonBootstrapper{
		config:     config,
		rootFolder: rootFolder,
		logger:     logger,
	}
}

// ValidateEnvironment checks that required environment variables are set
func (db *DaemonBootstrapper) ValidateEnvironment() error {
	required := map[string]string{
		"XDG_CONFIG_HOME":          db.rootFolder,
		"TB_SETTINGS_DEFAULTCHAIN": "mainnet",
		"TB_SETTINGS_INDEXPATH":    db.config.IndexPath(),
		"TB_SETTINGS_CACHEPATH":    db.config.CachePath(),
	}

	for key, expectedValue := range required {
		if value := strings.TrimSpace(strings.ToLower(expectedValue)); value == "" {
			return fmt.Errorf("required environment variable %s has empty expected value", key)
		}
	}

	return nil
}

// EnsureConfig creates or validates the chifra config file
func (db *DaemonBootstrapper) EnsureConfig() error {
	configFn := filepath.Join(db.rootFolder, "trueBlocks.toml")

	if file.FileExists(configFn) {
		db.logger.Info("Config file found", "fn", configFn)
		if !db.chainsConfigured(configFn) {
			db.logger.Error("Config file not configured", "fn", configFn)
			return fmt.Errorf("config file not configured")
		}
		return nil
	}

	db.logger.Warn("Config file not found", "fn", configFn)
	if err := db.createChifraConfig(); err != nil {
		db.logger.Error("Error creating config file", "error", err)
		return err
	}

	return nil
}

// chainsConfigured validates that all enabled chains have config sections
func (db *DaemonBootstrapper) chainsConfigured(configFn string) bool {
	chainStr := db.config.EnabledChains()
	chains := strings.Split(chainStr, ",")

	db.logger.Info("chifra config loaded")
	db.logger.Info("checking", "configFile", configFn, "nChains", len(chains))

	contents := file.AsciiFileToString(configFn)
	for _, chain := range chains {
		search := "[chains." + chain + "]"
		if !strings.Contains(contents, search) {
			msg := fmt.Sprintf("config file {%s} does not contain {%s}", configFn, search)
			db.logger.Error(msg)
			return false
		}
	}
	return true
}

// createChifraConfig generates a new trueBlocks.toml config file
func (db *DaemonBootstrapper) createChifraConfig() error {
	if err := file.EstablishFolder(db.rootFolder); err != nil {
		return err
	}

	chainStr := db.config.EnabledChains()
	chains := strings.Split(chainStr, ",")
	for _, chain := range chains {
		if err := db.createChainConfigFolder(chain); err != nil {
			return err
		}
	}

	tmpl, err := template.New("tmpl").Parse(configTmpl)
	if err != nil {
		return err
	}

	var buf bytes.Buffer
	if err = tmpl.Execute(&buf, db.config); err != nil {
		return err
	}
	if len(buf.String()) == 0 {
		return fmt.Errorf("empty config file")
	}

	configFn := filepath.Join(db.rootFolder, "trueBlocks.toml")
	err = file.StringToAsciiFile(configFn, buf.String())
	if err != nil {
		return err
	}
	db.logger.Info("Created config file", "configFile", configFn, "nChains", len(chains))
	return nil
}

// createChainConfigFolder creates the chain-specific config folder and downloads allocs.csv
func (db *DaemonBootstrapper) createChainConfigFolder(chain string) error {
	chainConfig := filepath.Join(db.rootFolder, "config", chain)
	if err := file.EstablishFolder(chainConfig); err != nil {
		return fmt.Errorf("failed to create folder %s: %w", chainConfig, err)
	}

	baseURL := "https://raw.githubusercontent.com/TrueBlocks/trueblocks-core/refs/heads/master/src/other/install/per-chain"
	allocURL, err := url.JoinPath(baseURL, chain, "allocs.csv")
	if err != nil {
		return err
	}
	allocFn := filepath.Join(chainConfig, "allocs.csv")
	dur := 100 * 365 * 24 * time.Hour // 100 years
	if _, err := downloadAndStore(allocURL, allocFn, dur); err != nil {
		db.logger.Warn(fmt.Errorf("failed to download and store allocs.csv for chain %s: %w", chain, err).Error())
		// It's not an error to not have an allocation file. IsArchiveNode assumes archive if not present.
		return nil
	}
	db.logger.Progress("Creating chain config", "chainConfig", allocFn)
	db.logger.Progress("Creating chain config", "source", allocURL)

	return nil
}
