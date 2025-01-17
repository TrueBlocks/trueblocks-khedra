package types

type Chain struct {
	Name    string   `koanf:"name" validate:"req_if_enabled"`                 // Must be non-empty
	RPCs    []string `koanf:"rpcs" validate:"req_if_enabled,dive,strict_url"` // Must have at least one reachable RPC URL
	Enabled bool     `koanf:"enabled"`                                        // Defaults to false if not specified
}

func NewChain(chain string) Chain {
	return Chain{
		Name:    chain,
		RPCs:    []string{"http://localhost:8545"},
		Enabled: true,
	}
}

func (c *Chain) IsEnabled() bool {
	return c.Enabled
}
