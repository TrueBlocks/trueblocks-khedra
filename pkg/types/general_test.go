package types

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/TrueBlocks/trueblocks-khedra/v2/pkg/validate"
	"github.com/stretchr/testify/assert"
	yamlv2 "gopkg.in/yaml.v2"
)

// Testing status: reviewed

// TestGeneralNew tests the initialization of the General type to ensure it is
// created correctly with valid default or input values.
func TestNewGeneral(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	assert.NoError(t, err)

	expectedPath := filepath.Join(homeDir, ".khedra", "data")
	g := NewGeneral()
	assert.Equal(t, expectedPath, g.DataFolder)
}

// TestGeneralValidation validates the functionality of the General type to ensure
// that invalid data is caught and proper validation rules are applied.
func TestGeneralValidation(t *testing.T) {
	defer SetTestEnv([]string{})()
	tests := []struct {
		name    string
		general General
		wantErr bool
	}{
		{
			name: "Valid General struct with all fields",
			general: General{
				DataFolder: createTempDir(t, true),
				Strategy:   "scratch",
				Detail:     "bloomFilters",
			},
			wantErr: false,
		},
		{
			name: "Non-existent DataFolder with valid strategy and detail",
			general: General{
				DataFolder: "/non/existent/path",
				Strategy:   "download",
				Detail:     "entireIndex",
			},
			wantErr: false,
		},
		{
			name: "Non-writable DataFolder with valid strategy and detail",
			general: General{
				DataFolder: createTempDir(t, false),
				Strategy:   "scratch",
				Detail:     "bloomFilters",
			},
			wantErr: false,
		},
		{
			name: "Empty DataFolder with valid strategy and detail",
			general: General{
				DataFolder: "",
				Strategy:   "download",
				Detail:     "entireIndex",
			},
			wantErr: true,
		},
		{
			name: "Invalid Strategy",
			general: General{
				DataFolder: createTempDir(t, true),
				Strategy:   "invalid_strategy",
				Detail:     "bloomFilters",
			},
			wantErr: true,
		},
		{
			name: "Invalid Detail",
			general: General{
				DataFolder: createTempDir(t, true),
				Strategy:   "scratch",
				Detail:     "invalid_detail",
			},
			wantErr: true,
		},
		{
			name: "Empty Strategy",
			general: General{
				DataFolder: createTempDir(t, true),
				Strategy:   "",
				Detail:     "bloomFilters",
			},
			wantErr: true,
		},
		{
			name: "Empty Detail",
			general: General{
				DataFolder: createTempDir(t, true),
				Strategy:   "scratch",
				Detail:     "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Validate(&tt.general)
			if tt.wantErr {
				assert.Error(t, err, "Expected error for test case '%s'", tt.name)
			} else {
				assert.NoError(t, err, "Did not expect error for test case '%s'", tt.name)
			}
		})
	}
}

// TestGeneralReadAndWrite tests the reading and writing operations of the General type
// to confirm accurate data handling and storage.
func TestGeneralSerialization(t *testing.T) {
	content := `
dataFolder: "expected/folder/name"
`
	var g General
	err := yamlv2.Unmarshal([]byte(content), &g)
	assert.NoError(t, err)
	assert.Equal(t, "expected/folder/name", g.DataFolder)

	out, err := yamlv2.Marshal(&g)
	assert.NoError(t, err)
	assert.Contains(t, string(out), "dataFolder: expected/folder/name")
}

func TestInvalidYAMLInput(t *testing.T) {
	content := `
dataFolder: "expected/folder/name"
strategy: download
detail: [invalid_array]
`
	var g General
	err := yamlv2.Unmarshal([]byte(content), &g)
	assert.Error(t, err, "Expected error for invalid YAML input")
}

func TestUnsupportedCharactersInFields(t *testing.T) {
	content := `
dataFolder: "expected/folder/\x00name"
strategy: "\tinvalid_strategy"
detail: "entireIndex"
`
	var g General
	err := yamlv2.Unmarshal([]byte(content), &g)
	assert.NoError(t, err, "Unexpected error during YAML parsing")
	err = validate.Validate(&g)
	assert.Error(t, err, "Expected validation error for unsupported characters")
}
