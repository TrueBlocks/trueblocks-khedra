package install

// KnownChain carries metadata for selectable chains in the install wizard.
type KnownChain struct {
	Name    string
	ChainID int
	// Future: recommended RPCs, explorer, etc.
}

// KnownChains returns a static ordered list of commonly used chains.
func KnownChains() []KnownChain {
	return []KnownChain{
		{Name: "mainnet", ChainID: 1},
		{Name: "gnosis", ChainID: 100},
		{Name: "sepolia", ChainID: 11155111},
		{Name: "holesky", ChainID: 17000},
	}
}

// IsKnownChain returns true if name matches a known chain.
func IsKnownChain(name string) bool {
	for _, kc := range KnownChains() {
		if kc.Name == name {
			return true
		}
	}
	return false
}
