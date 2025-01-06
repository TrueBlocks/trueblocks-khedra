package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestGeneralNew tests the initialization of the General type to ensure it is
// created correctly with valid default or input values.
func TestGeneralNew(t *testing.T) {
	defer SetTestEnv([]string{"TEST_MODE=true"})()
	g := NewGeneral()
	expected := "~/.khedra/data"
	assert.Equal(t, expected, g.DataDir, "Expected DataDir to be '%s', got '%s'", expected, g.DataDir)
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
				DataDir: createTempDir(t, true),
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
				DataDir: createTempDir(t, false),
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
			if tt.wantErr {
				assert.Error(t, err, "Expected error for test case '%s'", tt.name)
			} else {
				assert.NoError(t, err, "Did not expect error for test case '%s'", tt.name)
			}
			fmt.Println()
		})
	}

	invalidGeneral := General{}
	err := Validate.Struct(invalidGeneral)
	assert.Error(t, err, "Expected validation error for missing DataDir, got nil")

	validGeneral := General{
		DataDir: "/valid/path",
	}
	err = Validate.Struct(validGeneral)
	assert.NoError(t, err, "Expected no validation error, but got: %s", err)
}

// TestGeneralReadAndWrite tests the reading and writing operations of the General type
// to confirm accurate data handling and storage.
func TestGeneralReadAndWrite(t *testing.T) {
	tempFilePath := "temp_config.yaml"
	content := `
data_dir: "expected/folder/name"
`

	assertions := func(t *testing.T, general *General) {
		assert.Equal(t, "expected/folder/name", general.DataDir, "Expected data_dir to be 'expected/folder/name', got '%s'", general.DataDir)
	}

	ReadAndWriteTest[General](t, tempFilePath, content, assertions)
}
