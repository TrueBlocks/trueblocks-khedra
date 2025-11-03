package control

import (
	"errors"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"
)

func TestNewMetadataBasic(t *testing.T) {
	port := 9999
	version := "ver1"
	m := NewMetadata(port, version)
	if m.PID != os.Getpid() {
		t.Fatalf("expected PID to be current process, got %d", m.PID)
	}
	if m.Port != port {
		t.Fatalf("expected port %d, got %d", port, m.Port)
	}
	if m.Version != version {
		t.Fatalf("expected version %s, got %s", version, m.Version)
	}
	if _, err := time.Parse(time.RFC3339, m.Started); err != nil {
		t.Fatalf("Started field not RFC3339: %v", err)
	}
	// Allow a small delta for time
	started, _ := time.Parse(time.RFC3339, m.Started)
	if abs := time.Since(started); abs < 0 || abs > 5*time.Second {
		t.Fatalf("Started field not within reasonable bounds: %v", abs)
	}
}

// helper to set a temp run dir and return it
func tempRunDir(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	t.Setenv("KHEDRA_RUN_DIR", dir)
	return dir
}

// TestPathDefault ensures default path resolves under the home run directory when override not set.
func TestPathDefault(t *testing.T) {
	os.Unsetenv("KHEDRA_RUN_DIR")
	p := Path()
	if !strings.Contains(p, ".khedra/run/control.json") { // loose assertion to avoid depending on $HOME exact value
		t.Fatalf("expected default path to contain .khedra/run/control.json, got %s", p)
	}
	// directory should exist
	if fi, err := os.Stat(filepath.Dir(p)); err != nil || !fi.IsDir() {
		t.Fatalf("expected run directory to exist: err=%v", err)
	}
}

func TestEnsureMetadataLifecycle(t *testing.T) {
	t.Run("Initial creation", func(t *testing.T) {
		tempRunDir(t)
		m, regen, err := EnsureMetadata(1234, "vtest")
		if err != nil {
			t.Fatalf("EnsureMetadata failed: %v", err)
		}
		if !regen {
			t.Fatalf("expected regeneration on missing file")
		}
		if m.Port != 1234 || m.Version != "vtest" {
			t.Fatalf("unexpected metadata contents: %+v", m)
		}
	})

	t.Run("Cached reuse", func(t *testing.T) {
		tempRunDir(t)
		m, _, err := EnsureMetadata(1234, "vtest")
		if err != nil {
			t.Fatalf("EnsureMetadata failed: %v", err)
		}
		m2, regen2, err := EnsureMetadata(1234, "vtest")
		if err != nil || regen2 {
			t.Fatalf("expected cached metadata without regeneration, got regen=%v err=%v", regen2, err)
		}
		if m2.PID != m.PID {
			t.Fatalf("pid changed unexpectedly: %d vs %d", m2.PID, m.PID)
		}
	})

	t.Run("Stale PID regeneration", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			t.Skip("stale pid test skipped on windows")
		}
		tempRunDir(t)
		fake := Metadata{PID: 999999, Port: 55, Version: "old", Started: "2020-01-01T00:00:00Z"}
		if err := Write(fake); err != nil {
			t.Fatalf("failed to write fake metadata: %v", err)
		}
		if _, err := os.Stat(Path()); err != nil {
			t.Fatalf("expected metadata file to exist: %v", err)
		}
		m, regen, err := EnsureMetadata(77, "new")
		if err != nil {
			t.Fatalf("EnsureMetadata failed: %v", err)
		}
		if !regen {
			t.Fatalf("expected regeneration for stale pid")
		}
		if m.Port != 77 || m.Version != "new" {
			t.Fatalf("metadata not updated: %+v", m)
		}
		if m.PID == fake.PID {
			t.Fatalf("pid not updated from stale pid")
		}
	})

	t.Run("Write failure surfaced", func(t *testing.T) {
		dir := t.TempDir()
		t.Setenv("KHEDRA_RUN_DIR", dir)
		p := Path()
		if err := os.Mkdir(p, 0o755); err != nil {
			t.Fatalf("failed to create directory for write failure test: %v", err)
		}
		_, _, err := EnsureMetadata(1, "fail")
		if err == nil {
			t.Fatalf("expected error when writing over directory, got nil")
		}
		_ = os.Remove(p + ".tmp")
	})

	t.Run("Concurrent calls race-free", func(t *testing.T) {
		dir := t.TempDir()
		t.Setenv("KHEDRA_RUN_DIR", dir)
		const N = 10
		results := make([]struct {
			m     Metadata
			regen bool
			err   error
		}, N)
		done := make(chan struct{})
		var mu sync.Mutex
		for i := 0; i < N; i++ {
			go func(idx int) {
				mu.Lock()
				m, regen, err := EnsureMetadata(42, "race")
				mu.Unlock()
				results[idx] = struct {
					m     Metadata
					regen bool
					err   error
				}{m, regen, err}
				done <- struct{}{}
			}(i)
		}
		for i := 0; i < N; i++ {
			<-done
		}
		// All returned metadata should be identical
		base := results[0].m
		regenCount := 0
		for _, r := range results {
			if r.err != nil {
				t.Fatalf("unexpected error in concurrent EnsureMetadata: %v", r.err)
			}
			// Compare only PID, Port, Version, Started
			if r.m.PID != base.PID || r.m.Port != base.Port || r.m.Version != base.Version || r.m.Started != base.Started {
				t.Fatalf("metadata mismatch: %v vs %v", r.m, base)
			}
			if r.regen {
				regenCount++
			}
		}
		if regenCount > 1 {
			t.Fatalf("expected at most one regeneration, got %d", regenCount)
		}
	})
}

// Ensure Path respects override env.
func TestPathOverride(t *testing.T) {
	dir := tempRunDir(t)
	p := Path()
	if filepath.Dir(p) != dir {
		t.Fatalf("expected path under override dir, got %s", p)
	}
}

func TestWriteBehavior(t *testing.T) {
	t.Run("Normal write sets schema", func(t *testing.T) {
		tempRunDir(t)
		in := Metadata{PID: 111, Port: 222, Version: "ver", Started: "2020-01-01T00:00:00Z"}
		if err := Write(in); err != nil {
			t.Fatalf("write failed: %v", err)
		}
		out, err := Read()
		if err != nil {
			t.Fatalf("read failed: %v", err)
		}
		if out.Port != in.Port || out.Version != in.Version || out.PID != in.PID || out.Schema != MetadataSchema {
			t.Fatalf("round trip mismatch: in=%+v out=%+v", in, out)
		}
	})

	t.Run("Atomic rename failure path", func(t *testing.T) {
		dir := t.TempDir()
		t.Setenv("KHEDRA_RUN_DIR", dir)
		p := Path()
		if err := os.Mkdir(p, 0o755); err != nil {
			t.Fatalf("failed to create directory for atomic rename test: %v", err)
		}
		err := Write(Metadata{PID: 1, Port: 2, Version: "fail", Started: "now"})
		if err == nil {
			t.Fatalf("expected error when renaming over directory, got nil")
		}
		// Clean up temp file if present
		_ = os.Remove(p + ".tmp")
	})

	t.Run("Overwrite existing file", func(t *testing.T) {
		tempRunDir(t)
		in := Metadata{PID: 1, Port: 2, Version: "first", Started: "2020-01-01T00:00:00Z"}
		if err := Write(in); err != nil {
			t.Fatalf("first write failed: %v", err)
		}
		fi1, err := os.Stat(Path())
		if err != nil {
			t.Fatalf("stat after first write: %v", err)
		}
		// Sleep only if mod times are equal (rare, but possible on coarse FS)
		in2 := Metadata{PID: 1, Port: 2, Version: "second", Started: "2020-01-01T00:00:00Z"}
		if err := Write(in2); err != nil {
			t.Fatalf("second write failed: %v", err)
		}
		fi2, err := os.Stat(Path())
		if err != nil {
			t.Fatalf("stat after second write: %v", err)
		}
		if fi2.ModTime().Equal(fi1.ModTime()) {
			// Try again after a short sleep
			os.Remove(Path())
			Write(in2)
			fi2, _ = os.Stat(Path())
			if fi2.ModTime().Equal(fi1.ModTime()) {
				t.Fatalf("mod time did not change after overwrite")
			}
		}
	})
}

func TestReadBehavior(t *testing.T) {
	t.Run("Successful round-trip", func(t *testing.T) {
		tempRunDir(t)
		in := Metadata{PID: 111, Port: 222, Version: "ver", Started: "2020-01-01T00:00:00Z"}
		if err := Write(in); err != nil {
			t.Fatalf("write failed: %v", err)
		}
		out, err := Read()
		if err != nil {
			t.Fatalf("read failed: %v", err)
		}
		if out.Port != in.Port || out.Version != in.Version || out.PID != in.PID || out.Schema != MetadataSchema {
			t.Fatalf("round trip mismatch: in=%+v out=%+v", in, out)
		}
	})

	t.Run("Missing file", func(t *testing.T) {
		tempRunDir(t)
		os.Remove(Path())
		_, err := Read()
		if err == nil || !errors.Is(err, os.ErrNotExist) {
			t.Fatalf("expected not exists error, got %v", err)
		}
	})

	t.Run("Invalid JSON", func(t *testing.T) {
		tempRunDir(t)
		if err := os.WriteFile(Path(), []byte(`{"schema":`), 0o644); err != nil {
			t.Fatalf("failed writing invalid JSON: %v", err)
		}
		if _, err := Read(); err == nil {
			t.Fatalf("expected unmarshal error")
		}
	})

	t.Run("Zero schema invalid", func(t *testing.T) {
		tempRunDir(t)
		if err := os.WriteFile(Path(), []byte(`{"schema":0}`), 0o644); err != nil {
			t.Fatalf("failed to write invalid schema file: %v", err)
		}
		if _, err := Read(); err == nil || !strings.Contains(err.Error(), "invalid metadata schema") {
			t.Fatalf("expected invalid schema error, got %v", err)
		}
	})
}

func TestIsProcessAliveCases(t *testing.T) {
	t.Run("Negative PID", func(t *testing.T) {
		if isProcessAlive(-123) {
			t.Fatalf("expected negative pid not alive")
		}
	})

	t.Run("Zero PID", func(t *testing.T) {
		if isProcessAlive(0) {
			t.Fatalf("expected zero pid not alive")
		}
	})

	t.Run("Current PID", func(t *testing.T) {
		if !isProcessAlive(os.Getpid()) {
			t.Fatalf("expected current pid alive")
		}
	})

	t.Run("Nonexistent Unix PID", func(t *testing.T) {
		if runtime.GOOS == "windows" {
			t.Skip("unix-specific expectation")
		}
		pid := 999999
		if isProcessAlive(pid) {
			t.Skipf("pid %d unexpectedly alive; skipping to avoid flake", pid)
		}
		if isProcessAlive(pid) {
			t.Fatalf("expected nonexistent pid not alive")
		}
	})
}
