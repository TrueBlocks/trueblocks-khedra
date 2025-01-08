package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGeneralNew tests the initialization of the General type to ensure it is
// created correctly with valid default or input values.
func TestGeneralNew(t *testing.T) {
	defer SetTestEnv([]string{"TEST_MODE=true"})()
	g := NewGeneral()
	expected := "~/.khedra/data"
	assert.Equal(t, expected, g.DataFolder, "Expected DataFolder to be '%s', got '%s'", expected, g.DataFolder)
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

	invalidGeneral := General{}
	err := Validate.Struct(invalidGeneral)
	assert.Error(t, err, "Expected validation error for missing DataFolder, got nil")

	validGeneral := General{
		DataFolder: "/valid/path",
	}
	err = Validate.Struct(validGeneral)
	assert.NoError(t, err, "Expected no validation error, but got: %s", err)
}

// TestGeneralReadAndWrite tests the reading and writing operations of the General type
// to confirm accurate data handling and storage.
func TestGeneralReadAndWrite(t *testing.T) {
	tempFilePath := "temp_config.yaml"
	content := `
dataFolder: "expected/folder/name"
`

	assertions := func(t *testing.T, general *General) {
		assert.Equal(t, "expected/folder/name", general.DataFolder, "Expected dataFolder to be 'expected/folder/name', got '%s'", general.DataFolder)
	}

	ReadAndWriteWithAssertions[General](t, tempFilePath, content, assertions)
}
