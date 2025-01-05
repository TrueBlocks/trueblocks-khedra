package types

import (
	"os"
	"testing"

	coreFile "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/knadh/koanf/maps"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	yamlv2 "gopkg.in/yaml.v2"
)

func TestChainNewChain(t *testing.T) {
	c := NewChain("TestChain")

	expectedName := "TestChain"
	if c.Name != expectedName {
		t.Errorf("Expected Name to be '%s', got '%s'", expectedName, c.Name)
	}

	expectedRPCs := []string{"http://localhost:8545"}
	if len(c.RPCs) != len(expectedRPCs) || c.RPCs[0] != expectedRPCs[0] {
		t.Errorf("Expected RPCs to be '%v', got '%v'", expectedRPCs, c.RPCs)
	}

	if !c.Enabled {
		t.Errorf("Expected Enabled to be true, got %v", c.Enabled)
	}
}

func TestChainValidation(t *testing.T) {
	os.Setenv("TEST_MODE", "true")
	defer os.Unsetenv("TEST_MODE")

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
			err := Validate.Struct(tt.chain)
			checkValidationErrors(t, tt.name, err, tt.wantErr)
		})
	}
}

func TestChainValidation2(t *testing.T) {
	os.Setenv("TEST_MODE", "true")
	invalidChain := Chain{}

	err := Validate.Struct(invalidChain)
	if err == nil {
		t.Errorf("Expected validation error for missing Name and RPCs, got nil")
	}

	validChain := Chain{
		Name:    "ValidChain",
		RPCs:    []string{"http://localhost:8545"},
		Enabled: true,
	}

	err = Validate.Struct(validChain)
	if err != nil {
		t.Errorf("Expected no validation error, but got: %s", err)
	}
}

func TestChainReadAndWrite(t *testing.T) {
	tempFilePath := "temp_chain_config.yaml"
	defer os.Remove(tempFilePath)

	content := `
name: "TestChain"
rpcs:
  - "http://localhost:8545"
enabled: true
`
	err := coreFile.StringToAsciiFile(tempFilePath, content)
	if err != nil {
		t.Fatalf("Failed to write temporary file: %s", err)
	}

	k := koanf.New(".")
	err = k.Load(file.Provider(tempFilePath), yaml.Parser())
	if err != nil {
		t.Fatalf("Failed to load configuration using koanf: %s", err)
	}

	name := k.String("name")
	if name != "TestChain" {
		t.Errorf("Expected name to be 'TestChain', got '%s'", name)
	}

	rpcs := k.Strings("rpcs")
	if len(rpcs) != 1 || rpcs[0] != "http://localhost:8545" {
		t.Errorf("Expected RPCs to contain 'http://localhost:8545', got '%v'", rpcs)
	}

	enabled := k.Bool("enabled")
	if !enabled {
		t.Errorf("Expected enabled to be true, got %v", enabled)
	}

	output := maps.Unflatten(k.All(), ".")
	outputFilePath := "output_chain_config.yaml"
	defer os.Remove(outputFilePath)

	yamlContent, err := yamlv2.Marshal(output)
	if err != nil {
		t.Fatalf("Failed to marshal config to YAML: %s", err)
	}

	err = coreFile.StringToAsciiFile(outputFilePath, string(yamlContent))
	if err != nil {
		t.Fatalf("Failed to write output file: %s", err)
	}
}
