package types

import (
	"fmt"
	"os"
	"testing"

	coreFile "github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/file"
	"github.com/knadh/koanf/maps"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	yamlv2 "gopkg.in/yaml.v2"
)

func TestNewGeneral(t *testing.T) {
	g := NewGeneral()

	expectedDataDir := "~/.khedra/data"
	if g.DataDir != expectedDataDir {
		t.Errorf("Expected DataDir to be '%s', got '%s'", expectedDataDir, g.DataDir)
	}
}

func TestGeneralValidation(t *testing.T) {
	// Test cases for validation
	tests := []struct {
		name    string
		general General
		wantErr bool
	}{
		{
			name: "Valid General struct",
			general: General{
				DataDir: createTempDir(t, true), // Create a writable temp directory
			},
			wantErr: false,
		},
		{
			name: "Non-existent DataDir",
			general: General{
				DataDir: "/non/existent/path",
			},
			wantErr: false,
		},
		{
			name: "Non-writable DataDir",
			general: General{
				DataDir: createTempDir(t, false), // Create a non-writable temp directory
			},
			wantErr: false,
		},
		{
			name: "Empty DataDir",
			general: General{
				DataDir: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate.Struct(tt.general)
			checkValidationErrors(t, tt.name, err, tt.wantErr)
			fmt.Println()
		})
	}
}

func TestGeneralValidation2(t *testing.T) {
	invalidGeneral := General{}

	err := Validate.Struct(invalidGeneral)
	if err == nil {
		t.Errorf("Expected validation error for missing DataDir, got nil")
	}

	validGeneral := General{
		DataDir: "/valid/path",
	}

	err = Validate.Struct(validGeneral)
	if err != nil {
		t.Errorf("Expected no validation error, but got: %s", err)
	}
}

func TestReadAndWrite(t *testing.T) {
	tempFilePath := "temp_config.yaml"
	defer os.Remove(tempFilePath)

	content := `
data_dir: "/tmp/khedra/data"
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

	dataDir := k.String("data_dir")
	if dataDir != "/tmp/khedra/data" {
		t.Errorf("Expected data_dir to be '/tmp/khedra/data', got '%s'", dataDir)
	}

	output := maps.Unflatten(k.All(), ".")
	outputFilePath := "output_config.yaml"
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
