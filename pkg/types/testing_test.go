package types

import (
	"os"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Added tests per ai/TestDesign_testing.go.md

func TestSetupTest_RestoreEnv(t *testing.T) {
	// Ensure a clean starting state
	os.Unsetenv("FOO_FOR_TEST")
	os.Unsetenv("TEST_MODE")
	_, hadFoo := os.LookupEnv("FOO_FOR_TEST")
	_, hadTestMode := os.LookupEnv("TEST_MODE")

	cleanup := SetupTest([]string{"FOO_FOR_TEST=bar"})

	// During test
	if v, ok := os.LookupEnv("FOO_FOR_TEST"); assert.True(t, ok) {
		assert.Equal(t, "bar", v)
	}
	assert.Equal(t, "true", os.Getenv("TEST_MODE"))
	cfgFn := os.Getenv("KHEDRA_TEST_CONFIG_FN")
	assert.NotEmpty(t, cfgFn)
	if cfgFn != "" {
		_, err := os.Stat(cfgFn)
		assert.NoError(t, err, "expected temp config file to exist")
	}

	// Perform cleanup
	cleanup()

	// After cleanup: restoration
	_, fooStill := os.LookupEnv("FOO_FOR_TEST")
	if !hadFoo {
		assert.False(t, fooStill, "FOO_FOR_TEST should be unset after cleanup")
	}
	if !hadTestMode { // only assert if we purposefully unset earlier
		_, tmStill := os.LookupEnv("TEST_MODE")
		assert.False(t, tmStill, "TEST_MODE should be unset after cleanup if not originally present")
	}
}

func TestSetupTest_CreatesTempConfig(t *testing.T) {
	cleanup := SetupTest([]string{})
	cfgFn := os.Getenv("KHEDRA_TEST_CONFIG_FN")
	if cfgFn == "" {
		t.Fatal("KHEDRA_TEST_CONFIG_FN not set")
	}
	// File should exist now
	if _, err := os.Stat(cfgFn); err != nil {
		t.Fatalf("expected config file to exist: %v", err)
	}
	// Parent dir should exist
	dir := filepathDir(cfgFn)
	if _, err := os.Stat(dir); err != nil {
		t.Fatalf("expected config dir to exist: %v", err)
	}
	cleanup()
	// After cleanup file and directory should be gone
	if _, err := os.Stat(cfgFn); !os.IsNotExist(err) {
		t.Fatalf("expected config file removed, got err=%v", err)
	}
}

// filepathDir is a tiny wrapper to allow testing with reduced imports (avoid pulling full filepath here unnecessarily)
func filepathDir(p string) string {
	// simple last separator split (os specific)
	sep := string(os.PathSeparator)
	last := -1
	for i := len(p) - 1; i >= 0; i-- {
		if string(p[i]) == sep { // this works because PathSeparator is one byte on supported systems
			last = i
			break
		}
	}
	if last <= 0 {
		return "."
	}
	return p[:last]
}

func TestReadAndWriteWithAssertions_RoundTrip(t *testing.T) {
	type Simple struct {
		Name  string `koanf:"name" yaml:"name"`
		Value int    `koanf:"value" yaml:"value"`
	}

	tempFile, err := os.CreateTemp("", "roundtrip_*.yaml")
	if err != nil {
		t.Fatalf("failed to create temp file: %v", err)
	}
	tempPath := tempFile.Name()
	tempFile.Close()

	content := "name: tester\nvalue: 42\n"
	ReadAndWriteWithAssertions(t, tempPath, content, func(t *testing.T, s *Simple) {
		assert.Equal(t, "tester", s.Name)
		assert.Equal(t, 42, s.Value)
	})
}

func TestCreateTempDir_NonWritable(t *testing.T) {
	if runtime.GOOS == "windows" { // permission bits unreliable on Windows
		t.Skip("skipping non-writable dir test on windows")
	}
	dir := createTempDir(t, false)
	defer os.RemoveAll(dir)
	info, err := os.Stat(dir)
	if err != nil {
		t.Fatalf("stat failed: %v", err)
	}
	mode := info.Mode().Perm()
	assert.Equal(t, os.FileMode(0500), mode, "expected directory to have 0500 permissions")
	// Ensure owner write bit not set
	assert.Zero(t, mode&0200, "directory unexpectedly writable")
}
