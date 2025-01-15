package types

import (
	"testing"

	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/validate"
	"github.com/stretchr/testify/assert"
)

// Testing status: reviewed

func TestNewChain(t *testing.T) {
	c := NewChain("TestChain")

	assert.Equal(t, "TestChain", c.Name)
	assert.Equal(t, []string{"http://localhost:8545"}, c.RPCs)
	assert.True(t, c.Enabled)
}

func TestChainValidation(t *testing.T) {
	SetupTest([]string{})
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
			},
			wantErr: false,
		},
		{
			name: "Valid Chain with multiple valid RPCs",
			chain: Chain{
				Name:    "sepolia",
				RPCs:    []string{"https://sepolia.infura.io/v3/YOUR_PROJECT_ID", "https://another.valid.rpc"},
				Enabled: false,
			},
			wantErr: false,
		},
		{
			name: "Invalid Chain with missing Name",
			chain: Chain{
				Name:    "",
				RPCs:    []string{"https://mainnet.infura.io/v3/YOUR_PROJECT_ID"},
				Enabled: true,
			},
			wantErr: true,
		},
		{
			name: "Invalid Chain with empty RPCs",
			chain: Chain{
				Name:    "mainnet",
				RPCs:    []string{},
				Enabled: true,
			},
			wantErr: true,
		},
		{
			name: "Invalid Chain with an invalid RPC URL",
			chain: Chain{
				Name:    "mainnet",
				RPCs:    []string{"invalid-url"},
				Enabled: true,
			},
			wantErr: true,
		},
		{
			name: "Valid Chain with missing Enabled field",
			chain: Chain{
				Name: "mainnet",
				RPCs: []string{"https://mainnet.infura.io/v3/YOUR_PROJECT_ID"},
			},
			wantErr: false, // Enabled defaults to false, which is valid
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Validate(&tt.chain)
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
`

	assertions := func(t *testing.T, chain *Chain) {
		assert.Equal(t, "TestChain", chain.Name, "Expected name to be 'TestChain', got '%s'", chain.Name)
		assert.Equal(t, []string{"http://localhost:8545"}, chain.RPCs, "Expected RPCs to contain 'http://localhost:8545', got '%v'", chain.RPCs)
		assert.True(t, chain.Enabled, "Expected enabled to be true, got %v", chain.Enabled)
	}

	ReadAndWriteWithAssertions[Chain](t, tempFilePath, content, assertions)
}
