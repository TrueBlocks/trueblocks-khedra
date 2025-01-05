package types

import (
	"os"
	"strings"
	"testing"
)

// SetTempEnv sets an environment variable to a temporary value and returns a function to restore the original value.
func SetTempEnv(key, value string) func() {
	originalValue, exists := os.LookupEnv(key)
	os.Setenv(key, value)
	return func() {
		if exists {
			os.Setenv(key, originalValue)
		} else {
			os.Unsetenv(key)
		}
	}
}

// SetupTest sets up a temporary folder, updates the getConfigFn pointer, calls establishConfig,
// and assigns the config file path to the provided string pointer if it is not nil.
// Returns a cleanup function to restore the original state.
func SetupTest(
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

func SetupTest2(t *testing.T, env []string, configFile *string) func() {
	os.Setenv("TEST_MODE", "true")
	for i := 0; i < len(env); i++ {
		parts := strings.Split(env[i], "=")
		if len(parts) == 2 {
			os.Setenv(parts[0], parts[1])
		}
	}

	tempConfigFile := GetConfigFn()
	if configFile != nil {
		*configFile = tempConfigFile
	}
	EstablishConfig(tempConfigFile)

	return func() {
		os.Unsetenv("TEST_MODE")
	}
}
