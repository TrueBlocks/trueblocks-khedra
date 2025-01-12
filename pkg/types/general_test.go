package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v2"
)

// Testing status: reviewed

// TestGeneralNew tests the initialization of the General type to ensure it is
// created correctly with valid default or input values.
func TestNewGeneral(t *testing.T) {
	g := NewGeneral()
	assert.Equal(t, "~/.khedra/data", g.DataFolder)
}

// TestGeneralValidation validates the functionality of the General type to ensure
// that invalid data is caught and proper validation rules are applied.
func TestGeneralValidation(t *testing.T) {
	defer SetTestEnv([]string{"TEST_MODE=true"})()
	tests := []struct {
		name    string
		general General
		wantErr bool
	}{
		{
			name: "Valid General struct",
			general: General{
				DataFolder: createTempDir(t, true),
			},
			wantErr: false,
		},
		{
			name: "Non-existent DataFolder",
			general: General{
				DataFolder: "/non/existent/path",
			},
			wantErr: false,
		},
		{
			name: "Non-writable DataFolder",
			general: General{
				DataFolder: createTempDir(t, false),
			},
			wantErr: false,
		},
		{
			name: "Empty DataFolder",
			general: General{
				DataFolder: "",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate.Struct(tt.general)
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
	err := yaml.Unmarshal([]byte(content), &g)
	assert.NoError(t, err)
	assert.Equal(t, "expected/folder/name", g.DataFolder)

	out, err := yaml.Marshal(&g)
	assert.NoError(t, err)
	assert.Contains(t, string(out), "dataFolder: expected/folder/name")
}
