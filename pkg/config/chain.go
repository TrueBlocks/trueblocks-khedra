package config

type Chain struct {
	Name    string   `koanf:"name" validate:"required"`                                // Must be non-empty
	RPCs    []string `koanf:"rpcs" validate:"required,min=1,dive,strict_url,ping_one"` // Must have at least one reachable RPC URL
	Enabled bool     `koanf:"enabled"`                                                 // Defaults to false if not specified
}

func NewChain(chain string) Chain {
	return Chain{
		Name:    chain,
		RPCs:    []string{"http://localhost:8545"},
		Enabled: true,
	}
}
