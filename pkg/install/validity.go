package install

import (
	"strings"

	coreFile "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/logger"
	"github.com/TrueBlocks/trueblocks-khedra/v5/pkg/types"
	yamlv2 "gopkg.in/yaml.v2"
)

// Configured returns true only if a config file exists AND contains an enabled mainnet
// chain with at least one RPC URL present. Deep RPC reachability probing is deferred
// (performed during daemon startup); here we only perform a shallow structural check
// to decide whether to present the installation wizard.
func Configured() bool {
	fn := types.GetConfigFnNoCreate()
	if !coreFile.FileExists(fn) {
		return false
	}

	cfg := types.NewConfig()
	b := coreFile.AsciiFileToString(fn)
	if len(b) == 0 {
		logger.Warn("config file empty", "file", fn)
		return false
	}
	if err := yamlv2.Unmarshal([]byte(b), &cfg); err != nil {
		logger.Warn("config yaml unmarshal failed", "file", fn, "err", err)
		return false
	}
	// Must have a mainnet key
	if main, ok := cfg.Chains["mainnet"]; ok {
		if !main.Enabled {
			logger.Info("mainnet chain present but disabled; treating as not configured")
			return false
		}
		if len(main.RPCs) == 0 || strings.TrimSpace(main.RPCs[0]) == "" {
			logger.Info("mainnet enabled but no RPC present")
			return false
		}
		if main.ChainID == 0 {
			logger.Info("mainnet chainId zero", "cfgFile", fn, "chain", main, "rawLen", len(b))
			return false
		}
		return true
	}
	logger.Info("no suitable mainnet chain found", "file", fn, "chains", len(cfg.Chains))
	return false
}
