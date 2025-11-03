package install

import (
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/TrueBlocks/trueblocks-khedra/v6/pkg/types"
)

// DetectRecentCorruption scans for archived corrupt draft files created within the
// last 24 hours. Filenames follow pattern: config.draft.json.corrupt-<unixTs>
// Returns true if at least one such file exists.
func DetectRecentCorruption() bool {
	finalPath := types.GetConfigFnNoCreate()
	dir := filepath.Dir(finalPath)
	entries, err := os.ReadDir(dir)
	if err != nil {
		return false
	}
	cutoff := time.Now().Add(-24 * time.Hour)
	prefix := "config.draft.json.corrupt-"
	for _, e := range entries {
		name := e.Name()
		if !strings.HasPrefix(name, prefix) {
			continue
		}
		info, err := e.Info()
		if err != nil {
			continue
		}
		if info.ModTime().After(cutoff) {
			return true
		}
	}
	return false
}

// corruptionFlagFile is a sentinel file created when a corruption event is detected and
// a replacement draft is successfully loaded. Its presence (recent) triggers a one-time banner.
func corruptionFlagFile() string {
	fn := types.GetConfigFnNoCreate()
	return filepath.Join(filepath.Dir(fn), "config.draft.json.corruption.flag")
}

// SetCorruptionFlag creates/updates the sentinel file timestamp.
func SetCorruptionFlag() {
	_ = os.WriteFile(corruptionFlagFile(), []byte(time.Now().UTC().Format(time.RFC3339)), 0o600)
}

// ClearCorruptionFlag removes the sentinel flag file.
func ClearCorruptionFlag() {
	_ = os.Remove(corruptionFlagFile())
}

// ConsumeCorruptionFlag returns true once if the sentinel exists (and is recent) then clears it.
func ConsumeCorruptionFlag(maxAge time.Duration) bool {
	fn := corruptionFlagFile()
	st, err := os.Stat(fn)
	if err != nil {
		return false
	}
	if time.Since(st.ModTime()) > maxAge {
		_ = os.Remove(fn)
		return false
	}
	_ = os.Remove(fn)
	return true
}
