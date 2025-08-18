package types

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
)

type Chain struct {
	Name    string   `koanf:"name" yaml:"name" json:"name,omitempty" validate:"req_if_enabled"`                 // Must be non-empty
	RPCs    []string `koanf:"rpcs" yaml:"rpcs" json:"rpcs,omitempty" validate:"req_if_enabled,dive,strict_url"` // Must have at least one reachable RPC URL
	ChainID int      `koanf:"chainId" yaml:"chainId" json:"chainId,omitempty" validate:"non_zero"`              // Must be non-zero
	Enabled bool     `koanf:"enabled" yaml:"enabled" json:"enabled,omitempty"`                                  // Defaults to false if not specified
}

func NewChain(chain string, chainId int) Chain {
	return Chain{
		Name:    chain,
		RPCs:    []string{"http://localhost:8545"},
		ChainID: chainId,
		Enabled: true,
	}
}

func (ch Chain) IsEnabled() bool {
	return ch.Enabled
}

func (cc Chain) Symbol() string {
	if item := utils.GetChainListItem("~/.khedra", cc.ChainID); item != nil {
		return item.NativeCurrency.Symbol
	}
	return "Unknown"
}

func (cc Chain) RemoteExplorer() string {
	if item := utils.GetChainListItem("~/.khedra", cc.ChainID); item != nil {
		if len(item.Explorers) > 0 {
			return item.Explorers[0].URL
		}
	}
	return "Unknown"
}
