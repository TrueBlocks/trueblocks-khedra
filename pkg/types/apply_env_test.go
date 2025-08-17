package types

import (
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Testing status: reviewed

// TestApplyEnv ensures that applyEnv correctly applies environment variable values to a Config struct,
// covering various scenarios such as valid/invalid booleans, integers, empty values, unknown keys,
// and handling multiple chains/services.
func TestApplyEnv(t *testing.T) {
	applyChainSettings := func() {
		defer setEnv(map[string]string{
			"TB_KHEDRA_GENERAL_DATAFOLDER":     "/env/data",
			"TB_KHEDRA_CHAINS_MAINNET_RPCS":    "http://rpc1.mainnet,http://rpc2.mainnet",
			"TB_KHEDRA_CHAINS_MAINNET_ENABLED": "true",
		})()

		cfg := Config{
			General: General{
				DataFolder: "/default/data",
				Strategy:   "download",
				Detail:     "index",
			},
			Chains: map[string]Chain{
				"mainnet": {
					RPCs:    []string{"http://default.rpc"},
					Enabled: false,
				},
			},
			Services: map[string]Service{},
		}

		keys := getEnvironmentKeys(cfg, InEnv)
		err := applyEnv(keys, &cfg)
		assert.NoError(t, err)

		expected := Config{
			General: General{
				DataFolder: "/env/data",
				Strategy:   "download",
				Detail:     "index",
			},
			Chains: map[string]Chain{
				"mainnet": {
					RPCs:    []string{"http://rpc1.mainnet", "http://rpc2.mainnet"},
					Enabled: true,
				},
			},
			Services: map[string]Service{},
		}

		assert.Equal(t, expected, cfg)
	}
	t.Run("Apply Chain Settings", func(t *testing.T) { applyChainSettings() })

	applyServiceSettings := func() {
		defer setEnv(map[string]string{
			"TB_KHEDRA_SERVICES_API_ENABLED": "false",
			"TB_KHEDRA_SERVICES_API_PORT":    "9090",
		})()

		cfg := Config{
			Chains: map[string]Chain{},
			Services: map[string]Service{
				"api": {
					Name:    "api",
					Enabled: true,
					Port:    8080,
				},
			},
		}

		keys := getEnvironmentKeys(cfg, InEnv)
		err := applyEnv(keys, &cfg)
		assert.NoError(t, err)

		expected := Config{
			Chains: map[string]Chain{},
			Services: map[string]Service{
				"api": {
					Name:    "api",
					Enabled: false,
					Port:    9090,
				},
			},
		}

		assert.Equal(t, expected, cfg)
	}
	t.Run("Apply Service Settings", func(t *testing.T) { applyServiceSettings() })

	partialLoggingUpdate := func() {
		defer setEnv(map[string]string{
			"TB_KHEDRA_LOGGING_FOLDER":   "/env/logs",
			"TB_KHEDRA_LOGGING_TOFILE":   "true",
			"TB_KHEDRA_LOGGING_COMPRESS": "true",
		})()

		cfg := Config{
			Chains:   map[string]Chain{},
			Services: map[string]Service{},
			Logging: Logging{
				Folder:   "/default/logs",
				Compress: false,
			},
		}

		keys := getEnvironmentKeys(cfg, InEnv)
		err := applyEnv(keys, &cfg)
		assert.NoError(t, err)

		expected := Config{
			Chains:   map[string]Chain{},
			Services: map[string]Service{},
			Logging: Logging{
				Folder:   "/env/logs",
				ToFile:   true,
				Compress: true,
			},
		}

		assert.Equal(t, expected, cfg)
	}
	t.Run("Partial Logging Update", func(t *testing.T) { partialLoggingUpdate() })

	invalidBoolean := func() {
		defer setEnv(map[string]string{
			"TB_KHEDRA_CHAINS_MAINNET_ENABLED": "not_a_bool",
		})()

		cfg := Config{
			Chains: map[string]Chain{
				"mainnet": {
					Enabled: false,
				},
			},
		}

		keys := getEnvironmentKeys(cfg, InEnv)
		err := applyEnv(keys, &cfg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "environment variable has an invalid value")
	}
	t.Run("Invalid Boolean", func(t *testing.T) { invalidBoolean() })

	invalidInteger := func() {
		defer setEnv(map[string]string{
			"TB_KHEDRA_LOGGING_MAXSIZE": "not_a_number",
		})()

		cfg := Config{
			Logging: Logging{
				MaxSize: 50,
			},
		}

		keys := getEnvironmentKeys(cfg, InEnv)
		err := applyEnv(keys, &cfg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "environment variable has an invalid value")
	}
	t.Run("InvalidInteger", func(t *testing.T) { invalidInteger() })

	emptyValues := func() {
		cases := []struct {
			name   string
			envVar string
		}{
			{"Empty RPC List", "TB_KHEDRA_CHAINS_MAINNET_RPCS"},
			{"Empty Data Folder", "TB_KHEDRA_GENERAL_DATAFOLDER"},
		}

		for _, c := range cases {
			t.Run(c.name, func(t *testing.T) {
				defer setEnv(map[string]string{
					c.envVar: "",
				})()
				cfg := Config{
					General: General{
						DataFolder: "/default/data",
						Strategy:   "download",
						Detail:     "index",
					},
					Chains: map[string]Chain{
						"mainnet": {
							RPCs: []string{"http://default.rpc"},
						},
					},
				}

				keys := getEnvironmentKeys(cfg, InEnv)
				err := applyEnv(keys, &cfg)
				assert.Error(t, err)
				assert.Contains(t, err.Error(), "environment variable is empty:")
			})
		}
	}
	t.Run("Empty Values", func(t *testing.T) { emptyValues() })

	unknownKeys := func() {
		defer setEnv(map[string]string{
			"TB_KHEDRA_CHAINS_MAINNET_FOO": "some_value",
		})()

		cfg := Config{
			Chains: map[string]Chain{
				"mainnet": {
					RPCs:    []string{"http://default.rpc"},
					Enabled: true,
				},
			},
		}

		keys := getEnvironmentKeys(cfg, InEnv)
		err := applyEnv(keys, &cfg)

		assert.NoError(t, err, "applying unknown key should not produce an error")

		expected := Config{
			Chains: map[string]Chain{
				"mainnet": {
					RPCs:    []string{"http://default.rpc"},
					Enabled: true,
				},
			},
		}
		assert.Equal(t, expected, cfg, "config should remain unchanged for unknown keys")
	}
	t.Run("Unknown Keys", func(t *testing.T) { unknownKeys() })

	incompleteKeys := func() {
		defer setEnv(map[string]string{
			"TB_KHEDRA_CHAINS_MAINNET": "some_value",
		})()

		cfg := Config{
			Chains: map[string]Chain{
				"mainnet": {
					RPCs:    []string{"http://default.rpc"},
					Enabled: false,
				},
			},
		}
		keys := getEnvironmentKeys(cfg, InEnv)
		err := applyEnv(keys, &cfg)
		assert.NoError(t, err)
		assert.Equal(t, []string{"http://default.rpc"}, cfg.Chains["mainnet"].RPCs)
		assert.False(t, cfg.Chains["mainnet"].Enabled)
	}
	t.Run("Incomplete Keys", func(t *testing.T) { incompleteKeys() })

	multipleChainsServices := func() {
		defer setEnv(map[string]string{
			"TB_KHEDRA_CHAINS_MAINNET_RPCS":     "https://mainnet-rpc-1,https://mainnet-rpc-2",
			"TB_KHEDRA_CHAINS_MAINNET_ENABLED":  "true",
			"TB_KHEDRA_CHAINS_TESTNET_RPCS":     "https://testnet-rpc",
			"TB_KHEDRA_CHAINS_TESTNET_ENABLED":  "false",
			"TB_KHEDRA_SERVICES_API_ENABLED":    "true",
			"TB_KHEDRA_SERVICES_API_PORT":       "8080",
			"TB_KHEDRA_SERVICES_WORKER_ENABLED": "false",
			"TB_KHEDRA_SERVICES_WORKER_PORT":    "9090",
		})()

		cfg := Config{
			Chains: map[string]Chain{
				"mainnet": {
					RPCs:    []string{"http://default-mainnet"},
					Enabled: false,
				},
				"testnet": {
					RPCs:    []string{"http://default-testnet"},
					Enabled: true,
				},
			},
			Services: map[string]Service{
				"api": {
					Name:    "api",
					Enabled: false,
					Port:    7000,
				},
				"worker": {
					Name:    "worker",
					Enabled: true,
					Port:    6000,
				},
			},
		}

		keys := getEnvironmentKeys(cfg, InEnv)
		err := applyEnv(keys, &cfg)
		assert.NoError(t, err)

		assert.Equal(t, []string{"https://mainnet-rpc-1", "https://mainnet-rpc-2"}, cfg.Chains["mainnet"].RPCs)
		assert.True(t, cfg.Chains["mainnet"].Enabled)
		assert.Equal(t, []string{"https://testnet-rpc"}, cfg.Chains["testnet"].RPCs)
		assert.False(t, cfg.Chains["testnet"].Enabled)
		assert.True(t, cfg.Services["api"].Enabled)
		assert.Equal(t, 8080, cfg.Services["api"].Port)
		assert.False(t, cfg.Services["worker"].Enabled)
		assert.Equal(t, 9090, cfg.Services["worker"].Port)
	}
	t.Run("Mulitple Chains Services", func(t *testing.T) { multipleChainsServices() })

	largeRpcList := func() {
		rpcList := make([]string, 100)
		for i := 0; i < 100; i++ {
			rpcList[i] = fmt.Sprintf("http://rpc%d.example.com", i)
		}
		defer setEnv(map[string]string{
			"TB_KHEDRA_CHAINS_MAINNET_RPCS": strings.Join(rpcList, ","),
		})()

		cfg := Config{
			Chains: map[string]Chain{
				"mainnet": {
					RPCs:    []string{},
					Enabled: true,
				},
			},
		}

		keys := getEnvironmentKeys(cfg, InEnv)
		err := applyEnv(keys, &cfg)
		assert.NoError(t, err)
		assert.Len(t, cfg.Chains["mainnet"].RPCs, 100)
	}
	t.Run("LargeRpcList", func(t *testing.T) { largeRpcList() })
}

// Focused test 1 (from design item 1): apply general Strategy and Detail via environment.
func TestApplyEnv_GeneralStrategyDetail(t *testing.T) {
	defer setEnv(map[string]string{
		KeyStrategy: "fastsync",
		KeyDetail:   "full",
	})()
	cfg := NewConfig()
	// Sanity preconditions
	if cfg.General.Strategy == "fastsync" || cfg.General.Detail == "full" {
		t.Fatalf("precondition failed: defaults already match test values")
	}
	keys := []string{KeyStrategy, KeyDetail}
	if err := applyEnv(keys, &cfg); err != nil {
		t.Fatalf("applyEnv returned error: %v", err)
	}
	if cfg.General.Strategy != "fastsync" {
		t.Fatalf("expected Strategy=fastsync got=%s", cfg.General.Strategy)
	}
	if cfg.General.Detail != "full" {
		t.Fatalf("expected Detail=full got=%s", cfg.General.Detail)
	}
}

// Focused test 2: apply remaining logging fields (Filename, MaxBackups, MaxAge) and success path.
func TestApplyEnv_LoggingRemaining(t *testing.T) {
	defer setEnv(map[string]string{
		KeyLoggingFilename:   "custom.log",
		KeyLoggingMaxBackups: "7",
		KeyLoggingMaxAge:     "21",
	})()
	cfg := NewConfig()
	// Change defaults to ensure values actually update
	cfg.Logging.Filename = "old.log"
	cfg.Logging.MaxBackups = 3
	cfg.Logging.MaxAge = 10
	keys := []string{KeyLoggingFilename, KeyLoggingMaxBackups, KeyLoggingMaxAge}
	if err := applyEnv(keys, &cfg); err != nil {
		t.Fatalf("applyEnv error: %v", err)
	}
	if cfg.Logging.Filename != "custom.log" {
		t.Fatalf("Filename not updated: %s", cfg.Logging.Filename)
	}
	if cfg.Logging.MaxBackups != 7 {
		t.Fatalf("MaxBackups expected 7 got %d", cfg.Logging.MaxBackups)
	}
	if cfg.Logging.MaxAge != 21 {
		t.Fatalf("MaxAge expected 21 got %d", cfg.Logging.MaxAge)
	}
}

// Focused test 3: logging parse errors for bool and int fields (ToFile, Compress, MaxSize, MaxBackups, MaxAge)
func TestApplyEnv_LoggingParseErrors(t *testing.T) {
	cases := []struct{ key, val string }{
		{KeyLoggingToFile, "notabool"},
		{KeyLoggingCompress, "??"},
		{KeyLoggingMaxSize, "size"},
		{KeyLoggingMaxBackups, "five"},
		{KeyLoggingMaxAge, "age"},
	}
	for _, c := range cases {
		t.Run(c.key, func(t *testing.T) {
			defer setEnv(map[string]string{c.key: c.val})()
			cfg := NewConfig()
			if err := applyEnv([]string{c.key}, &cfg); err == nil {
				t.Fatalf("expected parse error for key %s value %s", c.key, c.val)
			}
		})
	}
}

// Focused test 4: chainID valid & invalid
func TestApplyEnv_ChainID(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		defer setEnv(map[string]string{"TB_KHEDRA_CHAINS_MAINNET_CHAINID": "999"})()
		cfg := NewConfig() // default mainnet ChainID=1
		if err := applyEnv([]string{"TB_KHEDRA_CHAINS_MAINNET_CHAINID"}, &cfg); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Chains["mainnet"].ChainID != 999 {
			t.Fatalf("expected ChainID 999 got %d", cfg.Chains["mainnet"].ChainID)
		}
	})
	t.Run("Invalid", func(t *testing.T) {
		defer setEnv(map[string]string{"TB_KHEDRA_CHAINS_MAINNET_CHAINID": "abc"})()
		cfg := NewConfig()
		if err := applyEnv([]string{"TB_KHEDRA_CHAINS_MAINNET_CHAINID"}, &cfg); err == nil {
			t.Fatalf("expected parse error for invalid chainid")
		}
	})
}

// Focused test 5: services Sleep & BatchSize valid & invalid
func TestApplyEnv_ServiceSleepBatchSize(t *testing.T) {
	t.Run("Valid", func(t *testing.T) {
		defer setEnv(map[string]string{
			"TB_KHEDRA_SERVICES_API_SLEEP":     "15",
			"TB_KHEDRA_SERVICES_API_BATCHSIZE": "250",
		})()
		cfg := NewConfig()
		// ensure defaults differ
		if cfg.Services["api"].Sleep == 15 || cfg.Services["api"].BatchSize == 250 {
			t.Fatalf("precondition mismatch: defaults already match")
		}
		keys := []string{"TB_KHEDRA_SERVICES_API_SLEEP", "TB_KHEDRA_SERVICES_API_BATCHSIZE"}
		if err := applyEnv(keys, &cfg); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if cfg.Services["api"].Sleep != 15 {
			t.Fatalf("expected Sleep=15 got %d", cfg.Services["api"].Sleep)
		}
		if cfg.Services["api"].BatchSize != 250 {
			t.Fatalf("expected BatchSize=250 got %d", cfg.Services["api"].BatchSize)
		}
	})
	t.Run("Invalid", func(t *testing.T) {
		cases := []string{"TB_KHEDRA_SERVICES_API_SLEEP", "TB_KHEDRA_SERVICES_API_BATCHSIZE"}
		vals := []string{"nope", "-"}
		for i, k := range cases {
			t.Run(k, func(t *testing.T) {
				defer setEnv(map[string]string{k: vals[i]})()
				cfg := NewConfig()
				if err := applyEnv([]string{k}, &cfg); err == nil {
					t.Fatalf("expected parse error for %s", k)
				}
			})
		}
	})
}

// Focused test 6: new chain and service entries created when absent
func TestApplyEnv_NewEntriesCreated(t *testing.T) {
	defer setEnv(map[string]string{
		"TB_KHEDRA_CHAINS_FRESH_CHAINID":    "123",
		"TB_KHEDRA_CHAINS_FRESH_ENABLED":    "true",
		"TB_KHEDRA_SERVICES_NEWSVC_ENABLED": "true",
		"TB_KHEDRA_SERVICES_NEWSVC_PORT":    "5555",
	})()
	cfg := NewConfig()
	// Remove mainnet & api to emphasize creation (optional)
	delete(cfg.Chains, "mainnet")
	delete(cfg.Services, "api")
	keys := []string{"TB_KHEDRA_CHAINS_FRESH_CHAINID", "TB_KHEDRA_CHAINS_FRESH_ENABLED", "TB_KHEDRA_SERVICES_NEWSVC_ENABLED", "TB_KHEDRA_SERVICES_NEWSVC_PORT"}
	if err := applyEnv(keys, &cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	fresh, ok := cfg.Chains["fresh"]
	if !ok {
		t.Fatalf("expected new chain 'fresh' created")
	}
	if fresh.ChainID != 123 || !fresh.Enabled {
		t.Fatalf("unexpected fresh chain values: %+v", fresh)
	}
	svc, ok := cfg.Services["newsvc"]
	if !ok {
		t.Fatalf("expected new service 'newsvc' created")
	}
	if !svc.Enabled || svc.Port != 5555 {
		t.Fatalf("unexpected newsvc values: %+v", svc)
	}
}

// Focused test 7: unknown service sub-key ignored (e.g., _FOO)
func TestApplyEnv_ServiceUnknownSubKeyIgnored(t *testing.T) {
	defer setEnv(map[string]string{"TB_KHEDRA_SERVICES_API_FOO": "bar"})()
	cfg := NewConfig()
	before := cfg.Services["api"]
	// get recognized keys only (unknown key not included)
	keys := getEnvironmentKeys(cfg, InEnv)
	if err := applyEnv(keys, &cfg); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	after := cfg.Services["api"]
	if before != after {
		t.Fatalf("service changed despite unknown sub-key: before=%+v after=%+v", before, after)
	}
}

// setEnv sets the specified environment variables, then returns a cleanup function that restores
// previous environment variable values. This is useful for scoped testing where environment variables
// must be set temporarily.
func setEnv(envVars map[string]string) func() {
	originalEnv := make(map[string]string)
	for key := range envVars {
		if val, exists := os.LookupEnv(key); exists {
			originalEnv[key] = val
		}
	}

	for key, val := range envVars {
		os.Setenv(key, val)
	}

	return func() {
		for key := range envVars {
			if val, exists := originalEnv[key]; exists {
				os.Setenv(key, val)
			} else {
				os.Unsetenv(key)
			}
		}
	}
}
