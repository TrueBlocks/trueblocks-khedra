// edge_case_tests.go
package config

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/types"
	"github.com/stretchr/testify/assert"
)

func TestInvalidPortForService(t *testing.T) {
	defer setTempEnvVar("TB_KHEDRA_SERVICES_API_PORT", "invalid_port")()
	defer setTempEnvVar("TEST_MODE", "true")()
	defer setupTest(t, nil)()

	_, err := loadConfig()
	assert.Error(t, err, "loadConfig should return an error for invalid port value")
	assert.Contains(t, err.Error(), "invalid_port", "Error message should indicate invalid port")
}

func TestLargeNumberOfChains(t *testing.T) {
	var configFile string
	defer setupTest(t, &configFile)()

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

	// Write the large config to the file
	types.WriteConfig(&cfg, configFile)

	// Load the configuration and verify all chains are present
	cfg = MustLoadConfig(configFile)
	assert.Equal(t, nChains+2, len(cfg.Chains), "All chains should be loaded correctly")
}

func TestMissingChainInConfig(t *testing.T) {
	defer setTempEnvVar("TB_KHEDRA_CHAINS_UNKNOWN_NAME", "unknown")()
	defer setTempEnvVar("TB_KHEDRA_CHAINS_UNKNOWN_RPCS", "http://unknown.rpc")()
	defer setTempEnvVar("TB_KHEDRA_CHAINS_UNKNOWN_ENABLED", "true")()
	defer setupTest(t, nil)()

	_, err := loadConfig()
	assert.Error(t, err, "An error should occur if an unknown chain is defined in the environment but not in the configuration file")
}

func TestEmptyRPCsForChain(t *testing.T) {
	var configFile string

	defer setTempEnvVar("TB_KHEDRA_CHAINS_MAINNET_RPCS", "")()
	defer setupTest(t, &configFile)()

	cfg := MustLoadConfig(configFile)
	assert.NotEmpty(t, cfg.Chains["mainnet"].RPCs, "Mainnet RPCs should not be empty in the final configuration")
}
