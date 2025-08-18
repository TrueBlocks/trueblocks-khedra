package install

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"sync"
	"time"

	coreFile "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-khedra/v5/pkg/types"
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
	if err := fsyncDir(dir); err != nil { // best effort
		// continue; not fatal on some platforms
	}
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
