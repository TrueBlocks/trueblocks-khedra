package install

import (
	"os"
	"path/filepath"
	"testing"
)

// writeConfig writes a config YAML to a temp file and points
// types.GetConfigFnNoCreate at it via KHEDRA_TEST_CONFIG_FN.
func writeConfig(t *testing.T, yaml string) func() {
	t.Helper()
	dir := t.TempDir()
	fn := filepath.Join(dir, "config.yaml")
	if err := os.WriteFile(fn, []byte(yaml), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}
	prev, hadPrev := os.LookupEnv("KHEDRA_TEST_CONFIG_FN")
	if err := os.Setenv("KHEDRA_TEST_CONFIG_FN", fn); err != nil {
		t.Fatalf("setenv: %v", err)
	}
	return func() {
		if hadPrev {
			_ = os.Setenv("KHEDRA_TEST_CONFIG_FN", prev)
		} else {
			_ = os.Unsetenv("KHEDRA_TEST_CONFIG_FN")
		}
	}
}

// resetAccessibilityCache clears the in-memory cache so test runs don't
// observe stale probe results from prior tests.
func resetAccessibilityCache() {
	accessibilityCacheMu.Lock()
	defer accessibilityCacheMu.Unlock()
	accessibilityCache = map[string]accessibilityCacheEntry{}
}

// When general.skipMainnetProbe is true the daemon must be able to come up
// without a reachable mainnet RPC — see issue #22 (kurtosis / air-gapped).
func TestConfigured_SkipMainnetProbe_SkipsReachabilityProbe(t *testing.T) {
	resetAccessibilityCache()
	// Use an obviously unreachable endpoint: a TEST-NET-1 address on a closed
	// port. If the probe runs, Configured() will return false.
	yaml := `general:
  dataFolder: "/tmp/khedra-test"
  strategy: "download"
  detail: "index"
  skipMainnetProbe: true
chains:
  mainnet:
    rpcs:
      - "http://192.0.2.1:1/unreachable"
    chainId: 1
    enabled: false
  kurtosis:
    rpcs:
      - "http://localhost:8545"
    chainId: 3151908
    enabled: true
`
	restore := writeConfig(t, yaml)
	defer restore()

	if !Configured() {
		t.Fatalf("Configured() = false; expected true when skipMainnetProbe=true even with unreachable mainnet RPC")
	}
}

// Default behavior: skipMainnetProbe omitted (== false) → probe must run.
func TestConfigured_SkipMainnetProbeOmitted_RequiresReachableProbe(t *testing.T) {
	resetAccessibilityCache()
	yaml := `general:
  dataFolder: "/tmp/khedra-test"
  strategy: "download"
  detail: "index"
chains:
  mainnet:
    rpcs:
      - "http://192.0.2.1:1/unreachable"
    chainId: 1
    enabled: true
`
	restore := writeConfig(t, yaml)
	defer restore()

	if Configured() {
		t.Fatalf("Configured() = true; expected false when skipMainnetProbe is omitted and mainnet RPC is unreachable")
	}
}

// skipMainnetProbe explicitly false → probe must run.
func TestConfigured_SkipMainnetProbeFalse_RequiresReachableProbe(t *testing.T) {
	resetAccessibilityCache()
	yaml := `general:
  dataFolder: "/tmp/khedra-test"
  strategy: "download"
  detail: "index"
  skipMainnetProbe: false
chains:
  mainnet:
    rpcs:
      - "http://192.0.2.1:1/unreachable"
    chainId: 1
    enabled: true
`
	restore := writeConfig(t, yaml)
	defer restore()

	if Configured() {
		t.Fatalf("Configured() = true; expected false when skipMainnetProbe=false and mainnet RPC is unreachable")
	}
}

// skipMainnetProbe=true with no RPCs at all is still invalid — downstream code
// reads main.RPCs[0] (env wiring in action_daemon.go), so the structural check
// must remain.
func TestConfigured_SkipMainnetProbe_StillRequiresRPC(t *testing.T) {
	resetAccessibilityCache()
	yaml := `general:
  dataFolder: "/tmp/khedra-test"
  strategy: "download"
  detail: "index"
  skipMainnetProbe: true
chains:
  mainnet:
    rpcs: []
    chainId: 1
    enabled: false
`
	restore := writeConfig(t, yaml)
	defer restore()

	if Configured() {
		t.Fatalf("Configured() = true; expected false when mainnet has no RPCs even with skipMainnetProbe=true")
	}
}

// skipMainnetProbe=true with chainId == 0 is still rejected — malformed YAML.
func TestConfigured_SkipMainnetProbe_StillRequiresChainID(t *testing.T) {
	resetAccessibilityCache()
	yaml := `general:
  dataFolder: "/tmp/khedra-test"
  strategy: "download"
  detail: "index"
  skipMainnetProbe: true
chains:
  mainnet:
    rpcs:
      - "http://localhost:8545"
    chainId: 0
    enabled: false
`
	restore := writeConfig(t, yaml)
	defer restore()

	if Configured() {
		t.Fatalf("Configured() = true; expected false when mainnet has chainId=0 even with skipMainnetProbe=true")
	}
}
