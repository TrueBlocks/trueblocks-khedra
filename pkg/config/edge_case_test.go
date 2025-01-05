// edge_case_tests.go
package config

import (
	"fmt"
	"strconv"
	"testing"

	coreFile "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/types"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

func TestServiceInvalidPort(t *testing.T) {
	defer types.SetupTest([]string{
		"TB_KHEDRA_SERVICES_API_PORT=invalid_port",
	})()

	if cfg, err := LoadConfig(); err != nil {
		assert.Error(t, err, "loadConfig should return an error for invalid port value")
		assert.Contains(t, err.Error(), "invalid_port", "Error message should indicate invalid port")
	} else {
		t.Error("loadConfig should return an error for invalid port", cfg.Services["api"])
	}
}

func TestChainLargeNumberOfChains(t *testing.T) {
	defer types.SetupTest([]string{})()

	nChains := 1000
	cfg := types.NewConfig()
	cfg.Chains = make(map[string]types.Chain)
	for i := 0; i < nChains; i++ {
		chainName := "chain" + strconv.Itoa(i)
		cfg.Chains[chainName] = types.Chain{
			Name:    chainName,
			RPCs:    []string{fmt.Sprintf("http://%s.rpc", chainName)},
			Enabled: true,
		}
	}

	bytes, _ := yaml.Marshal(cfg)
	coreFile.StringToAsciiFile(types.GetConfigFn(), string(bytes))

	var err error
	if cfg, err = LoadConfig(); err != nil {
		t.Error(err)
	} else {
		assert.Equal(t, nChains+2, len(cfg.Chains), "All chains should be loaded correctly")
	}
}

func TestChainMissingInConfig(t *testing.T) {
	defer types.SetupTest([]string{
		"TB_KHEDRA_CHAINS_UNKNOWN_NAME=unknown",
		"TB_KHEDRA_CHAINS_UNKNOWN_RPCS=http://unknown.rpc",
		"TB_KHEDRA_CHAINS_UNKNOWN_ENABLED=true",
	})()

	if cfg, err := LoadConfig(); err != nil {
		assert.Error(t, err, "An error should occur if an unknown chain is defined in the environment but not in the configuration file")
	} else {
		t.Error("loadConfig should return an error for invalid chain", cfg.Chains["unknown"])
	}
}

func TestChainEmptyRPCs(t *testing.T) {
	defer types.SetupTest([]string{
		"TB_KHEDRA_CHAINS_MAINNET_RPCS=",
	})()

	if cfg, err := LoadConfig(); err != nil {
		t.Error(err)
	} else {
		assert.NotEmpty(t, cfg.Chains["mainnet"].RPCs, "Mainnet RPCs should not be empty in the final configuration")
	}
}
