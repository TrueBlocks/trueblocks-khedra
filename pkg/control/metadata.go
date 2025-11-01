package control

import (
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/v6/pkg/utils"
)

const MetadataSchema = 1

type Metadata struct {
	Schema  int    `json:"schema"`
	PID     int    `json:"pid"`
	Port    int    `json:"port"`
	Version string `json:"version"`
	Started string `json:"started"`
}

// Path returns the location of the control metadata file. It may be overridden
// in tests (or power users) by setting KHEDRA_RUN_DIR to an alternate directory.
func Path() string {
	if custom := os.Getenv("KHEDRA_RUN_DIR"); custom != "" {
		_ = os.MkdirAll(custom, 0o755)
		return filepath.Join(custom, "control.json")
	}
	runDir := utils.ResolvePath("~/.khedra/run")
	_ = os.MkdirAll(runDir, 0o755)
	return filepath.Join(runDir, "control.json")
}

func Write(meta Metadata) error {
	meta.Schema = MetadataSchema
	b, err := json.MarshalIndent(meta, "", "  ")
	if err != nil {
		return err
	}
	fn := Path()
	tmp := fn + ".tmp"
	if err := os.WriteFile(tmp, b, 0o644); err != nil {
		return err
	}
	if err := os.Rename(tmp, fn); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	return nil
}

func Read() (Metadata, error) {
	fn := Path()
	b, err := os.ReadFile(fn)
	if err != nil {
		return Metadata{}, err
	}
	var m Metadata
	if err := json.Unmarshal(b, &m); err != nil {
		return Metadata{}, err
	}
	if m.Schema == 0 {
		return Metadata{}, errors.New("invalid metadata schema")
	}
	return m, nil
}

func NewMetadata(port int, version string) Metadata {
	return Metadata{
		PID:     os.Getpid(),
		Port:    port,
		Version: version,
		Started: time.Now().UTC().Format(time.RFC3339),
	}
}

// isProcessAlive attempts a best-effort check for a live process. On Unix we use
// a signal 0; on other platforms we conservatively return true (to avoid
// excessive churn from false negatives). A PID matching the current process is
// considered alive.
func isProcessAlive(pid int) bool {
	if pid <= 0 {
		return false
	}
	if pid == os.Getpid() {
		return true
	}
	if runtime.GOOS != "windows" { // best effort for unix-y systems
		err := syscall.Kill(pid, 0)
		if err == nil {
			return true
		}
		if errno, ok := err.(syscall.Errno); ok {
			if errno == syscall.EPERM { // exists but no permission
				return true
			}
		}
		return false
	}
	// On Windows (or unknown), assume alive so we don't thrash metadata.
	return true
}

// EnsureMetadata guarantees a control metadata file exists and is "fresh". A
// metadata file is considered stale iff the recorded PID does not represent a
// currently running process (and is not the current process). When stale, a new
// metadata object is written using the supplied port & version. The returned
// boolean indicates whether regeneration occurred.
func EnsureMetadata(port int, version string) (Metadata, bool, error) {
	m, err := Read()
	if err != nil {
		fresh := NewMetadata(port, version)
		if werr := Write(fresh); werr != nil {
			return fresh, true, werr
		}
		return fresh, true, nil
	}
	if m.PID != os.Getpid() && !isProcessAlive(m.PID) {
		fresh := NewMetadata(port, version)
		if werr := Write(fresh); werr != nil {
			return fresh, true, werr
		}
		return fresh, true, nil
	}
	return m, false, nil
}
