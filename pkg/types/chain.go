package types

import "github.com/TrueBlocks/trueblocks-khedra/v2/pkg/validate"

type Chain struct {
	Name    string   `koanf:"name" json:"name,omitempty" validate:"req_if_enabled"` // Must be non-empty
	RPCs    []string `koanf:"rpcs" validate:"req_if_enabled,dive,strict_url"`       // Must have at least one reachable RPC URL
	ChainId int      `koanf:"chainId" json:"chainId,omitempty" validate:"non_zero"` // Must be non-zero
	Enabled bool     `koanf:"enabled"`                                              // Defaults to false if not specified
}

func NewChain(chain string, chainId int) Chain {
	return Chain{
		Name:    chain,
		RPCs:    []string{"http://localhost:8545"},
		ChainId: chainId,
		Enabled: true,
	}
}

func (ch *Chain) IsEnabled() bool {
	return ch.Enabled
}

func (ch *Chain) HasValidRpc() bool {
	for _, rpc := range ch.RPCs {
		if err := validate.TryConnect(ch.Name, rpc, 2); err == nil {
			return true
		}
	}
	return false
}
