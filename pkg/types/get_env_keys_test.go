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
			"TB_KHEDRA_GENERAL_DATAFOLDER",
			"TB_KHEDRA_GENERAL_STRATEGY",
			"TB_KHEDRA_GENERAL_DETAIL",
			"TB_KHEDRA_CHAINS_MAINNET_ENABLED",
			"TB_KHEDRA_CHAINS_MAINNET_RPCS",
			"TB_KHEDRA_CHAINS_MAINNET_CHAINID",
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

	// New test: skip list enforcement (no suffix matches)
	t.Run("SkipList", func(t *testing.T) {
		cfg := NewConfig()
		keys := getEnvironmentKeys(cfg, InStruct)
		skipSuffixes := []string{"_NAME", "_API_BATCHSIZE", "_API_SLEEP", "_IPFS_BATCHSIZE", "_IPFS_SLEEP", "_MONITOR_PORT", "_SCRAPER_PORT"}
		for _, k := range keys {
			for _, suf := range skipSuffixes {
				if len(k) >= len(suf) && k[len(k)-len(suf):] == suf {
					t.Fatalf("key %s should have been skipped (suffix %s)", k, suf)
				}
			}
		}
	})

	// New test: dynamic addition of extra chain adds its keys (services are limited to fixed set)
	t.Run("DynamicMaps", func(t *testing.T) {
		cfg := NewConfig()
		cfg.Chains["goerli"] = NewChain("goerli", 5)
		keys := getEnvironmentKeys(cfg, InStruct)
		set := map[string]bool{}
		for _, k := range keys {
			set[k] = true
		}
		want := []string{
			"TB_KHEDRA_CHAINS_GOERLI_ENABLED",
			"TB_KHEDRA_CHAINS_GOERLI_RPCS",
			"TB_KHEDRA_CHAINS_GOERLI_CHAINID",
		}
		for _, w := range want {
			if !set[w] {
				t.Fatalf("expected dynamic key %s not found", w)
			}
		}
	})

	// New test: pointer vs value equivalence
	t.Run("PointerInput", func(t *testing.T) {
		cfg := NewConfig()
		valKeys := getEnvironmentKeys(cfg, InStruct)
		ptrKeys := getEnvironmentKeys(&cfg, InStruct)
		if len(valKeys) != len(ptrKeys) {
			t.Fatalf("length mismatch val=%d ptr=%d", len(valKeys), len(ptrKeys))
		}
		m := map[string]int{}
		for _, k := range valKeys {
			m[k]++
		}
		for _, k := range ptrKeys {
			m[k]--
		}
		for k, v := range m {
			if v != 0 {
				t.Fatalf("key mismatch: %s count %d", k, v)
			}
		}
	})

	// New test: InEnv single var filtering
	t.Run("InEnvSingleVar", func(t *testing.T) {
		defer setEnv(map[string]string{"TB_KHEDRA_LOGGING_LEVEL": "debug"})()
		cfg := NewConfig()
		keys := getEnvironmentKeys(cfg, InEnv)
		if len(keys) != 1 || keys[0] != "TB_KHEDRA_LOGGING_LEVEL" {
			t.Fatalf("expected only LOGGING_LEVEL got %v", keys)
		}
	})
}
