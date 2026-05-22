package install

import (
	"strings"
	"sync"
	"time"

	coreFile "github.com/TrueBlocks/trueblocks-chifra/v6/pkg/file"
	"github.com/TrueBlocks/trueblocks-chifra/v6/pkg/logger"
	"github.com/TrueBlocks/trueblocks-chifra/v6/pkg/rpc"
	yamlv2 "gopkg.in/yaml.v2"

	"github.com/TrueBlocks/trueblocks-khedra/v6/pkg/types"
)

// Cache for mainnet accessibility checks
var (
	accessibilityCache    = map[string]accessibilityCacheEntry{}
	accessibilityCacheMu  sync.Mutex
	accessibilityCacheTTL = 30 * time.Second
)

type accessibilityCacheEntry struct {
	accessible bool
	checkedAt  time.Time
}

// checkMainnetAccessible checks if mainnet RPC is accessible and returns chainId == 1
// Uses an in-memory cache with 30 second TTL to avoid hitting RPC repeatedly
func checkMainnetAccessible(rpcUrl string) bool {
	accessibilityCacheMu.Lock()
	defer accessibilityCacheMu.Unlock()

	// Check cache first
	if entry, ok := accessibilityCache[rpcUrl]; ok {
		if time.Since(entry.checkedAt) < accessibilityCacheTTL {
			logger.Info("mainnet accessibility cache hit", "url", rpcUrl, "accessible", entry.accessible)
			return entry.accessible
		}
		// Cache expired, delete old entry
		delete(accessibilityCache, rpcUrl)
	}

	// Perform the actual RPC check
	probe, err := rpc.PingRpc(rpcUrl)
	if err != nil {
		logger.Info("mainnet RPC ping failed", "url", rpcUrl, "error", err)
	}
	accessible := probe.OK && (probe.ChainID == "0x1" || probe.ChainID == "1")

	// Cache the result
	accessibilityCache[rpcUrl] = accessibilityCacheEntry{
		accessible: accessible,
		checkedAt:  time.Now(),
	}

	if !accessible {
		if !probe.OK {
			logger.Info("mainnet RPC unreachable", "url", rpcUrl, "error", probe.Error)
		} else {
			logger.Info("mainnet RPC returns wrong chainId", "url", rpcUrl, "chainId", probe.ChainID, "expected", "1")
		}
	}

	return accessible
}

// Configured returns true only if a config file exists AND contains a mainnet
// configuration. By default mainnet RPC must be reachable and return
// chainId == 1. When general.skipMainnetProbe is true, the reachability probe
// is skipped so the daemon can run in hermetic / air-gapped environments
// (CI, containers, isolated devnets) where no mainnet node is available;
// structural checks (RPC non-empty, chainId non-zero) still apply.
// Deep RPC reachability probing is performed here but cached for performance;
// this determines whether to present the installation wizard.
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
		if len(main.RPCs) == 0 || strings.TrimSpace(main.RPCs[0]) == "" {
			logger.Info("mainnet RPC missing but is required")
			return false
		}

		if main.ChainID == 0 {
			logger.Info("mainnet chainId zero", "cfgFile", fn, "chain", main, "rawLen", len(b))
			return false
		}

		if cfg.General.SkipMainnetProbe {
			return true
		}

		// Verify mainnet RPC is reachable and returns chainId == 1 (with caching)
		if !checkMainnetAccessible(main.RPCs[0]) {
			return false
		}
		return true
	}
	logger.Info("no suitable mainnet chain found", "file", fn, "chains", len(cfg.Chains))
	return false
}
