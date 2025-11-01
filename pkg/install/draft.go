package install

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	coreFile "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/v5/pkg/file"
	"github.com/TrueBlocks/trueblocks-khedra/v6/pkg/types"
	yamlv2 "gopkg.in/yaml.v2"
)

const (
	DraftSchema    = 1
	draftFileName  = "config.draft.json"
	backupFileName = "config.prev.yaml" // single rolling backup of final yaml
)

var (
	draftMu sync.Mutex
)

// DraftMeta carries metadata about the draft file.
type DraftMeta struct {
	Schema    int       `json:"schema"`
	Updated   time.Time `json:"updated"`
	Session   string    `json:"session"`
	EstDiskGB int       `json:"estDiskGb,omitempty"`
	EstHours  int       `json:"estHours,omitempty"`
}

// Draft holds an in-progress configuration plus metadata.
type Draft struct {
	Meta   DraftMeta    `json:"meta"`
	Config types.Config `json:"config"`
}

// draftPaths derives related paths from the final config path.
func draftPaths() (draftPath string, finalPath string, dir string) {
	finalPath = types.GetConfigFnNoCreate()
	dir = filepath.Dir(finalPath)
	draftPath = filepath.Join(dir, draftFileName)
	return
}

// DraftFilePath returns the absolute path to the draft config file (whether or not it exists).
func DraftFilePath() string {
	d, _, _ := draftPaths()
	return d
}

// LoadDraft loads the draft if it exists. Returns (nil, os.ErrNotExist) if absent.
// On corruption (unmarshal error) the file is archived with a timestamp suffix and error returned.
func LoadDraft() (*Draft, error) {
	draftMu.Lock()
	defer draftMu.Unlock()
	draftPath, _, _ := draftPaths()
	if !coreFile.FileExists(draftPath) {
		return nil, os.ErrNotExist
	}
	raw := coreFile.AsciiFileToString(draftPath)
	if raw == "" { // treat empty as corrupt
		archiveCorrupt(draftPath, []byte(raw))
		SetCorruptionFlag()
		return nil, errors.New("empty draft (archived)")
	}
	var d Draft
	if err := json.Unmarshal([]byte(raw), &d); err != nil {
		archiveCorrupt(draftPath, []byte(raw))
		SetCorruptionFlag()
		return nil, fmt.Errorf("corrupt draft archived: %w", err)
	}
	if d.Meta.Schema != DraftSchema { // schema mismatch -> archive and restart
		archiveCorrupt(draftPath, []byte(raw))
		SetCorruptionFlag()
		return nil, fmt.Errorf("draft schema mismatch (expected %d got %d)", DraftSchema, d.Meta.Schema)
	}
	return &d, nil
}

// SaveDraftAtomic writes the draft atomically.
func SaveDraftAtomic(d *Draft) error {
	draftMu.Lock()
	defer draftMu.Unlock()
	draftPath, _, dir := draftPaths()
	if d == nil {
		return errors.New("nil draft")
	}
	// Backfill missing chain IDs using known chain list (lightweight) before persisting.
	for key, ch := range d.Config.Chains {
		if ch.ChainID == 0 {
			for _, kc := range KnownChains() {
				if kc.Name == key {
					ch.ChainID = kc.ChainID
					break
				}
			}
			// Ensure mainnet defaults to 1 even if not in KnownChains for some reason
			if ch.ChainID == 0 && key == "mainnet" {
				ch.ChainID = 1
			}
			d.Config.Chains[key] = ch
		}
	}
	d.Meta.Updated = time.Now().UTC()
	d.Meta.Schema = DraftSchema
	data, err := json.MarshalIndent(d, "", "  ")
	if err != nil {
		return err
	}
	tmp := fmt.Sprintf("%s.tmp-%d-%d", draftPath, os.Getpid(), rand.Int63())
	if err := os.WriteFile(tmp, data, 0o600); err != nil {
		return err
	}
	if err := fsyncFile(tmp); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	// if err := fsyncDir(dir); err != nil { // best effort
	// 	// continue; not fatal on some platforms
	// }
	if err := os.Rename(tmp, draftPath); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	_ = fsyncDir(dir)
	return nil
}

// NewDraftFromConfig seeds a new draft from an existing final config (or default if missing).
func NewDraftFromConfig(session string) *Draft {
	cfg := types.NewConfig()
	finalPath := types.GetConfigFnNoCreate()
	if coreFile.FileExists(finalPath) {
		if raw := coreFile.AsciiFileToString(finalPath); len(raw) > 0 {
			// Best-effort parse YAML into cfg (ignore error; fallback default config)
			_ = yamlv2.Unmarshal([]byte(raw), &cfg)
		}
	}
	return &Draft{Meta: DraftMeta{Schema: DraftSchema, Updated: time.Now().UTC(), Session: session}, Config: cfg}
}

// BackupFinalConfig creates/overwrites the single rolling backup of the final config file.
func BackupFinalConfig() (string, error) {
	_, finalPath, dir := draftPaths()
	if !coreFile.FileExists(finalPath) {
		return "", os.ErrNotExist
	}
	backupPath := filepath.Join(dir, backupFileName)
	data := coreFile.AsciiFileToString(finalPath)
	tmp := backupPath + ".tmp"
	if err := os.WriteFile(tmp, []byte(data), 0o600); err != nil {
		return "", err
	}
	if err := fsyncFile(tmp); err != nil {
		_ = os.Remove(tmp)
		return "", err
	}
	if err := os.Rename(tmp, backupPath); err != nil {
		_ = os.Remove(tmp)
		return "", err
	}
	_ = fsyncDir(dir)
	return backupPath, nil
}

// archiveCorrupt renames a corrupt draft to a timestamped file for diagnostics.
func archiveCorrupt(path string, contents []byte) {
	ts := time.Now().Unix()
	newName := fmt.Sprintf("%s.corrupt-%d", path, ts)
	_ = os.Rename(path, newName)
	if !coreFile.FileExists(newName) {
		_ = os.WriteFile(newName, contents, 0o600)
	}
}

// fsync helpers (best effort, ignore on unsupported platforms)
func fsyncFile(fn string) error {
	f, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer f.Close()
	if err := f.Sync(); err != nil {
		return err
	}
	return nil
}

func fsyncDir(dir string) error {
	df, err := os.Open(dir)
	if err != nil {
		return err
	}
	defer df.Close()
	return df.Sync()
}

// RemoveDraft deletes the draft file (used after successful apply).
func RemoveDraft() error {
	draftPath, _, _ := draftPaths()
	if coreFile.FileExists(draftPath) {
		return os.Remove(draftPath)
	}
	return nil
}

// SaveDraft is an alias for SaveDraftAtomic for consistency
func SaveDraft(d *Draft) error {
	return SaveDraftAtomic(d)
}

// ApplyFormToDraft updates a draft config with form values
func ApplyFormToDraft(d *Draft, form map[string][]string) {
	if d == nil {
		return
	}

	// Helper to get single form value
	getValue := func(key string) string {
		if vals, ok := form[key]; ok && len(vals) > 0 {
			return strings.TrimSpace(vals[0])
		}
		return ""
	}

	// Helper to check if form key is present (for checkboxes)
	hasKey := func(key string) bool {
		_, exists := form[key]
		return exists
	}

	// Update General settings
	if v := getValue("dataFolder"); v != "" {
		d.Config.General.DataFolder = v
	}

	// Handle strategy and detail changes with estimation updates
	strategy := getValue("strategy")
	detail := getValue("detail")
	if strategy != "" || detail != "" {
		// Use current values if not provided in form
		if strategy == "" {
			strategy = d.Config.General.Strategy
		}
		if detail == "" {
			detail = d.Config.General.Detail
		}
		// This updates both the config and the estimates
		UpdateIndexStrategy(d, strategy, detail)
	}

	// Update Chains - handle enable/disable checkboxes and RPC URLs
	for name, chain := range d.Config.Chains {
		enableKey := name + "_enabled"
		rpcKey := "chain_rpc_" + name

		// Update enabled state if checkbox is present in form
		if hasKey(enableKey) {
			val := getValue(enableKey)
			switch val {
			case "1":
				chain.Enabled = true
			case "0":
				chain.Enabled = false
			}
			// If value is neither "1" nor "0", do nothing
		}

		// Update RPC URL if present
		if rpcUrl := getValue(rpcKey); rpcUrl != "" {
			if len(chain.RPCs) == 0 {
				chain.RPCs = []string{rpcUrl}
			} else {
				chain.RPCs[0] = rpcUrl
			}
		}

		d.Config.Chains[name] = chain
	}

	// Update Services
	if d.Config.Services == nil {
		d.Config.Services = make(map[string]types.Service)
	}
	for _, serviceName := range []string{"scraper", "monitor", "api", "ipfs"} {
		key := serviceName + "_enabled"
		if hasKey(key) {
			// Get existing service or create new one
			service := d.Config.Services[serviceName]
			service.Name = serviceName
			val := getValue(key)
			switch val {
			case "1":
				service.Enabled = true
			case "0":
				service.Enabled = false
			}
			// If value is neither "1" nor "0", do nothing
			d.Config.Services[serviceName] = service
		}
	}

	// Update Logging
	if v := getValue("level"); v != "" {
		d.Config.Logging.Level = v
	}
	if v := getValue("folder"); v != "" {
		d.Config.Logging.Folder = v
	}
	if v := getValue("filename"); v != "" {
		d.Config.Logging.Filename = v
	}
	if v := getValue("maxSize"); v != "" {
		if size, err := strconv.Atoi(v); err == nil {
			d.Config.Logging.MaxSize = size
		}
	}
	if v := getValue("maxBackups"); v != "" {
		if backups, err := strconv.Atoi(v); err == nil {
			d.Config.Logging.MaxBackups = backups
		}
	}
	if v := getValue("maxAge"); v != "" {
		if age, err := strconv.Atoi(v); err == nil {
			d.Config.Logging.MaxAge = age
		}
	}

	// Update logging checkboxes
	if hasKey("toFile") {
		val := getValue("toFile")
		switch val {
		case "1":
			d.Config.Logging.ToFile = true
		case "0":
			d.Config.Logging.ToFile = false
		}
		// If value is neither "1" nor "0", do nothing
	}
	if hasKey("compress") {
		val := getValue("compress")
		switch val {
		case "1":
			d.Config.Logging.Compress = true
		case "0":
			d.Config.Logging.Compress = false
		}
		// If value is neither "1" nor "0", do nothing
	}
}
