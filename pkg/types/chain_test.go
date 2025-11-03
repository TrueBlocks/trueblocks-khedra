package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Testing status: reviewed

func TestNewChain(t *testing.T) {
	c := NewChain("TestChain", -1)

	assert.Equal(t, "TestChain", c.Name)
	assert.Equal(t, []string{"http://localhost:8545"}, c.RPCs)
	assert.True(t, c.Enabled)
}

func TestChainValidation(t *testing.T) {
	defer SetupTest([]string{})()
	tests := []struct {
		name    string
		chain   Chain
		wantErr bool
	}{
		{
			name: "Valid Chain with one valid RPC",
			chain: Chain{
				Name:    "mainnet",
				RPCs:    []string{"https://mainnet.infura.io/v3/YOUR_PROJECT_ID"},
				Enabled: true,
				ChainID: 1,
			},
			wantErr: false,
		},
		{
			name: "Valid Chain with multiple valid RPCs",
			chain: Chain{
				Name:    "sepolia",
				RPCs:    []string{"https://sepolia.infura.io/v3/YOUR_PROJECT_ID", "https://another.valid.rpc"},
				Enabled: false,
				ChainID: 1,
			},
			wantErr: false,
		},
		{
			name: "Invalid Chain with missing Name",
			chain: Chain{
				Name:    "",
				RPCs:    []string{"https://mainnet.infura.io/v3/YOUR_PROJECT_ID"},
				Enabled: true,
				ChainID: 1,
			},
			wantErr: true,
		},
		{
			name: "Invalid Chain with empty RPCs",
			chain: Chain{
				Name:    "mainnet",
				RPCs:    []string{},
				Enabled: true,
				ChainID: 1,
			},
			wantErr: true,
		},
		{
			name: "Invalid Chain with an invalid RPC URL",
			chain: Chain{
				Name:    "mainnet",
				RPCs:    []string{"invalid-url"},
				Enabled: true,
				ChainID: 1,
			},
			wantErr: true,
		},
		{
			name: "Valid Chain with missing Enabled field",
			chain: Chain{
				Name:    "mainnet",
				RPCs:    []string{"https://mainnet.infura.io/v3/YOUR_PROJECT_ID"},
				ChainID: 1,
			},
			wantErr: false, // Enabled defaults to false, which is valid
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(&tt.chain)
			if tt.wantErr {
				assert.Error(t, err, "Expected error for test case '%s'", tt.name)
			} else {
				assert.NoError(t, err, "Did not expect error for test case '%s'", tt.name)
			}
		})
	}
}

func TestChainReadAndWrite(t *testing.T) {
	tempFilePath := "temp_chain_config.yaml"
	content := `
name: "TestChain"
rpcs:
  - "http://localhost:8545"
enabled: true
chainId: 0
`

	assertions := func(t *testing.T, chain *Chain) {
		assert.Equal(t, "TestChain", chain.Name, "Expected name to be 'TestChain', got '%s'", chain.Name)
		assert.Equal(t, []string{"http://localhost:8545"}, chain.RPCs, "Expected RPCs to contain 'http://localhost:8545', got '%v'", chain.RPCs)
		assert.True(t, chain.Enabled, "Expected enabled to be true, got %v", chain.Enabled)
	}

	ReadAndWriteWithAssertions(t, tempFilePath, content, assertions)
}

// Added targeted tests per ai/TestDesign_chain.go.md
func TestChainValidation_NonZeroChainID(t *testing.T) {
	defer SetupTest([]string{})()
	ch := Chain{
		Name:    "dev",
		RPCs:    []string{"http://localhost:8545"},
		ChainID: 0,
		Enabled: true,
	}
	err := Validate(&ch)
	assert.Error(t, err, "expected error for zero ChainID when enabled")
}

func TestChainValidation_DisabledVariants(t *testing.T) {
	defer SetupTest([]string{})()
	tests := []struct {
		name    string
		chain   Chain
		wantErr bool
	}{
		{
			name: "Disabled chain missing name and rpcs with valid ChainID passes",
			chain: Chain{
				Enabled: false,
				ChainID: 123,
			},
			wantErr: false,
		},
		{
			name: "Disabled chain missing name/rpcs but zero ChainID still errors",
			chain: Chain{
				Enabled: false,
				ChainID: 0,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(&tt.chain)
			if tt.wantErr {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestChain_IsEnabled(t *testing.T) {
	ch := NewChain("sample", 999)
	assert.True(t, ch.IsEnabled())
	ch.Enabled = false
	assert.False(t, ch.IsEnabled())
}

func TestChain_MetadataHelpersReturnUnknown(t *testing.T) {
	// Use an improbable chain ID; no chain list item expected
	ch := Chain{
		Name:    "mystery",
		RPCs:    []string{"http://localhost:8545"},
		ChainID: 987654321,
		Enabled: true,
	}
	assert.Equal(t, "Unknown", ch.Symbol())
	assert.Equal(t, "Unknown", ch.RemoteExplorer())
}
