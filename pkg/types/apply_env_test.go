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
