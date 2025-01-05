package types

import (
	"os"
	"strings"
	"testing"
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
	cleanupFn := SetTestEnv(env)
	tempConfigFile := GetConfigFn()
	EstablishConfig(tempConfigFile)
	return func() {
		cleanupFn()
		os.Remove(tempConfigFile)
	}
}

// SetupTestOld sets up a temporary folder, updates the getConfigFn pointer, calls establishConfig,
// and assigns the config file path to the provided string pointer if it is not nil.
// Returns a cleanup function to restore the original state.
func SetupTestOld(
	t *testing.T,
	configFile *string,
	getConfigFnPtr func() string,
	establishConfigFn func(configPath string) bool,
) func() {
	os.Setenv("TEST_MODE", "true")

	tempConfigFile := getConfigFnPtr()
	if configFile != nil {
		*configFile = tempConfigFile
	}
	establishConfigFn(tempConfigFile)

	return func() {
		os.Unsetenv("TEST_MODE")
	}
}
