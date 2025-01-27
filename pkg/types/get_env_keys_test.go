package types

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Testing status: not_reviewed

func TestGetEnvironmentKeys(t *testing.T) {
	testGetEnvKeys := func(t *testing.T) {
		cfg := NewConfig()
		keys := getEnvironmentKeys(cfg, InStruct)
		sort.Strings(keys)
		assert.ElementsMatch(t, []string{
			"TB_KHEDRA_CHAINS_MAINNET_ENABLED",
			"TB_KHEDRA_CHAINS_MAINNET_RPCS",
			"TB_KHEDRA_GENERAL_DATAFOLDER",
			"TB_KHEDRA_GENERAL_STRATEGY",
			"TB_KHEDRA_GENERAL_DETAIL",
			"TB_KHEDRA_LOGGING_COMPRESS",
			"TB_KHEDRA_LOGGING_FILENAME",
			"TB_KHEDRA_LOGGING_FOLDER",
			"TB_KHEDRA_LOGGING_TOFILE",
			"TB_KHEDRA_LOGGING_LEVEL",
			"TB_KHEDRA_LOGGING_MAXAGE",
			"TB_KHEDRA_LOGGING_MAXBACKUPS",
			"TB_KHEDRA_LOGGING_MAXSIZE",
			"TB_KHEDRA_SERVICES_API_ENABLED",
			"TB_KHEDRA_SERVICES_API_PORT",
			"TB_KHEDRA_SERVICES_IPFS_ENABLED",
			"TB_KHEDRA_SERVICES_IPFS_PORT",
			"TB_KHEDRA_SERVICES_MONITOR_BATCHSIZE",
			"TB_KHEDRA_SERVICES_MONITOR_ENABLED",
			"TB_KHEDRA_SERVICES_MONITOR_SLEEP",
			"TB_KHEDRA_SERVICES_SCRAPER_BATCHSIZE",
			"TB_KHEDRA_SERVICES_SCRAPER_ENABLED",
			"TB_KHEDRA_SERVICES_SCRAPER_SLEEP",
		}, keys)
	}
	t.Run("Test getEnv", testGetEnvKeys)

	testGetEnvKeysInEnv := func(t *testing.T) {
		defer setEnv(map[string]string{
			"TB_KHEDRA_CHAINS_MAINNET_ENABLED": "false",
			"TB_KHEDRA_LOGGING_FILENAME":       "\"A filename\"",
			"TB_KHEDRA_GENERAL_STRATEGY":       "scratch",
		})()
		cfg := NewConfig()
		keys := getEnvironmentKeys(cfg, InEnv)
		sort.Strings(keys)
		assert.ElementsMatch(t, []string{
			"TB_KHEDRA_CHAINS_MAINNET_ENABLED",
			"TB_KHEDRA_LOGGING_FILENAME",
			"TB_KHEDRA_GENERAL_STRATEGY",
		}, keys)
	}
	t.Run("Test getEnv", testGetEnvKeysInEnv)
}
