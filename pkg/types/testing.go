package types

import (
	"os"
	"sort"
	"strings"
	"testing"

	coreFile "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/utils"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	yamlv2 "gopkg.in/yaml.v2"
)

// SetTestEnv sets an environment variable to a temporary value and returns a function to restore the original value.
func SetTestEnv(env []string) func() {
	os.Setenv("TEST_MODE", "true")
	for i := 0; i < len(env); i++ {
		parts := strings.Split(env[i], "=")
		if len(parts) == 2 {
			os.Setenv(parts[0], parts[1])
		}
	}

	return func() {
		os.Unsetenv("TEST_MODE")
		for i := 0; i < len(env); i++ {
			parts := strings.Split(env[i], "=")
			if len(parts) == 2 {
				os.Unsetenv(parts[0])
			}
		}
	}
}

// SetupTest sets up the test environment by setting environment variables and establishing
// a temporary configuration file. It takes a slice of environment variable strings in the
// format "KEY=VALUE". Returns a cleanup function which may be deferred to remove the
// temporary configuration file and unset the environment variables.
func SetupTest(env []string) func() {
	originalArgs := os.Args
	envCleanup := SetTestEnv(env)
	tempConfigFile := GetConfigFn()
	cfg := NewConfig()
	bytes, _ := yamlv2.Marshal(cfg)
	coreFile.StringToAsciiFile(tempConfigFile, string(bytes))
	return func() {
		envCleanup()
		os.Remove(tempConfigFile)
		os.Args = originalArgs
	}
}

// ReadAndWriteWithAssertions writes the provided content to a temporary file, reads it
// back using koanf, performs assertions on the loaded configuration, and
// writes the output to another file. It takes the temporary file path,
// content to write, a function for assertions, and a testing object.
func ReadAndWriteWithAssertions[T any](t *testing.T, tempFilePath string, content string, assertions func(*testing.T, *T)) {
	defer os.Remove(tempFilePath)

	err := coreFile.StringToAsciiFile(tempFilePath, content)
	if err != nil {
		t.Fatalf("Failed to write temporary file: %s", err)
	}

	k := koanf.New(".")
	err = k.Load(file.Provider(tempFilePath), yaml.Parser())
	if err != nil {
		t.Fatalf("Failed to load configuration using koanf: %s", err)
	}

	var instance T
	err = k.Unmarshal("", &instance)
	if err != nil {
		t.Fatalf("Failed to unmarshal data into type: %s", err)
	}

	assertions(t, &instance)

	marshaledContent, err := yamlv2.Marshal(instance)
	if err != nil {
		t.Fatalf("Failed to marshal instance back to YAML: %s", err)
	}

	// cleanAndSort should take a string, remove all whitespace and then return the
	// sort with all of its characters sorted alphabetically (make comparing possible)
	cleanAndSort := func(s string) string {
		s = utils.RemoveAny(s, "\n\r\t\"` ")
		chars := strings.Split(s, "")
		sort.Strings(chars)
		return strings.Join(chars, "")
	}

	marshalCleaned := cleanAndSort(strings.ToLower(string(marshaledContent)))
	contentCleaned := cleanAndSort(strings.ToLower(content))
	if marshalCleaned != contentCleaned {
		t.Errorf("Mismatch between marshaled content and input content.\nMarshaled:\n%s\nInput:\n%s", marshalCleaned, contentCleaned)
	}
}

// createTempDir creates a temporary directory for testing.
// If writable is false, it makes the directory non-writable.
func createTempDir(t *testing.T, writable bool) string {
	t.Helper()
	dir, err := os.MkdirTemp("", "test_general")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	if !writable {
		err := os.Chmod(dir, 0500) // Read and execute permissions only
		if err != nil {
			t.Fatalf("Failed to make directory non-writable: %v", err)
		}
	}

	return dir
}
