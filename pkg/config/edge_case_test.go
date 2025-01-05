// edge_case_tests.go
package config

import (
	"fmt"
	"path/filepath"
	"strconv"
	"testing"

	coreFile "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/types"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestServiceInvalidPort(t *testing.T) {
	defer types.SetTempEnv("TB_KHEDRA_SERVICES_API_PORT", "invalid_port")()
	defer types.SetTempEnv("TEST_MODE", "true")()

	tmpDir := t.TempDir()
	configFile := filepath.Join(tmpDir, "config.yaml")

	types.EstablishConfig(configFile)

	_, err := loadConfig()
	assert.Error(t, err, "loadConfig should return an error for invalid port value")
	assert.Contains(t, err.Error(), "invalid_port", "Error message should indicate invalid port")
}

func TestChainLargeNumberOf(t *testing.T) {
	var configFile string
	defer types.SetupTest(t, &configFile, types.GetConfigFn, types.EstablishConfig)()

	cfg := types.NewConfig()
	cfg.Chains = make(map[string]types.Chain)
	nChains := 1000
	for i := 0; i < nChains; i++ {
		chainName := "chain" + strconv.Itoa(i)
		// fmt.Println(chainName)
		cfg.Chains[chainName] = types.Chain{
			Name:    chainName,
			RPCs:    []string{fmt.Sprintf("http://%s.rpc", chainName)},
			Enabled: true,
		}
	}

	bytes, _ := yaml.Marshal(cfg)
	coreFile.StringToAsciiFile(configFile, string(bytes))

	// Load the configuration and verify all chains are present (two are there from defaults)
	cfg = MustLoadConfig(configFile)
	assert.Equal(t, nChains+2, len(cfg.Chains), "All chains should be loaded correctly")
}

func TestChainMissingInConfig(t *testing.T) {
	defer types.SetTempEnv("TB_KHEDRA_CHAINS_UNKNOWN_NAME", "unknown")()
	defer types.SetTempEnv("TB_KHEDRA_CHAINS_UNKNOWN_RPCS", "http://unknown.rpc")()
	defer types.SetTempEnv("TB_KHEDRA_CHAINS_UNKNOWN_ENABLED", "true")()
	defer types.SetupTest(t, nil, types.GetConfigFn, types.EstablishConfig)()

	_, err := loadConfig()
	assert.Error(t, err, "An error should occur if an unknown chain is defined in the environment but not in the configuration file")
}

func TestChainEmptyRPCs(t *testing.T) {
	var configFile string

	defer types.SetTempEnv("TB_KHEDRA_CHAINS_MAINNET_RPCS", "")()
	defer types.SetupTest(t, &configFile, types.GetConfigFn, types.EstablishConfig)()

	cfg := MustLoadConfig(configFile)
	assert.NotEmpty(t, cfg.Chains["mainnet"].RPCs, "Mainnet RPCs should not be empty in the final configuration")
}
