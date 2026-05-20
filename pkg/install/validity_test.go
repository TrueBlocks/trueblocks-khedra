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

func TestMainnetExplicitlyDisabled(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		yaml string
		want bool
	}{
		{
			name: "explicit false",
			yaml: "chains:\n  mainnet:\n    enabled: false\n    chainId: 1\n",
			want: true,
		},
		{
			name: "omitted",
			yaml: "chains:\n  mainnet:\n    chainId: 1\n",
			want: false,
		},
		{
			name: "explicit true",
			yaml: "chains:\n  mainnet:\n    enabled: true\n    chainId: 1\n",
			want: false,
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := mainnetExplicitlyDisabled([]byte(tc.yaml)); got != tc.want {
				t.Fatalf("mainnetExplicitlyDisabled() = %v, want %v", got, tc.want)
			}
		})
	}
}

// When mainnet is disabled the daemon must be able to come up without a
// reachable mainnet RPC — see issue #22 (kurtosis / air-gapped deployments).
func TestConfigured_MainnetDisabled_SkipsReachabilityProbe(t *testing.T) {
	resetAccessibilityCache()
	// Use an obviously unreachable endpoint: a reserved-block address on a
	// closed port. If the probe runs, Configured() will return false.
	yaml := `chains:
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
		t.Fatalf("Configured() = false; expected true when mainnet.enabled=false even with unreachable mainnet RPC")
	}
}

// When enabled is omitted, treat as "must verify RPC" (legacy configs).
func TestConfigured_MainnetOmittedEnabled_RequiresReachableProbe(t *testing.T) {
	resetAccessibilityCache()
	yaml := `chains:
  mainnet:
    rpcs:
      - "http://192.0.2.1:1/unreachable"
    chainId: 1
  kurtosis:
    rpcs:
      - "http://localhost:8545"
    chainId: 3151908
    enabled: true
`
	restore := writeConfig(t, yaml)
	defer restore()

	if Configured() {
		t.Fatalf("Configured() = true; expected false when mainnet.enabled is omitted and RPC is unreachable")
	}
}

// When mainnet is enabled the existing strict probe must still gate
// Configured() — an unreachable mainnet RPC means the wizard should run.
func TestConfigured_MainnetEnabled_RequiresReachableProbe(t *testing.T) {
	resetAccessibilityCache()
	yaml := `chains:
  mainnet:
    rpcs:
      - "http://192.0.2.1:1/unreachable"
    chainId: 1
    enabled: true
`
	restore := writeConfig(t, yaml)
	defer restore()

	if Configured() {
		t.Fatalf("Configured() = true; expected false when mainnet.enabled=true with unreachable RPC")
	}
}

// Disabled mainnet with no RPCs at all is still an invalid config — we don't
// want to silently let through configs that would break code paths that read
// main.RPCs[0] later (env wiring, etc.).
func TestConfigured_MainnetDisabled_StillRequiresRPC(t *testing.T) {
	resetAccessibilityCache()
	yaml := `chains:
  mainnet:
    rpcs: []
    chainId: 1
    enabled: false
`
	restore := writeConfig(t, yaml)
	defer restore()

	if Configured() {
		t.Fatalf("Configured() = true; expected false when mainnet has no RPCs even if disabled")
	}
}

// Disabled mainnet with chainId == 0 is still rejected — the YAML is malformed
// and the wizard should fix it.
func TestConfigured_MainnetDisabled_StillRequiresChainID(t *testing.T) {
	resetAccessibilityCache()
	yaml := `chains:
  mainnet:
    rpcs:
      - "http://localhost:8545"
    chainId: 0
    enabled: false
`
	restore := writeConfig(t, yaml)
	defer restore()

	if Configured() {
		t.Fatalf("Configured() = true; expected false when mainnet has chainId=0 even if disabled")
	}
}
