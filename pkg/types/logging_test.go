package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewLogging(t *testing.T) {
	logging := NewLogging()

	// Check default values
	assert.Equal(t, "~/.khedra/logs", logging.Folder)
	assert.Equal(t, "khedra.log", logging.Filename)
	assert.Equal(t, 10, logging.MaxSizeMb)
	assert.Equal(t, 3, logging.MaxBackups)
	assert.Equal(t, 10, logging.MaxAgeDays)
	assert.Equal(t, "info", logging.LogLevel)
	assert.True(t, logging.Compress)
}

func TestLoggingValidation(t *testing.T) {
	tempDir := createTempDir(t, true) // Helper function to create a temp writable directory

	tests := []struct {
		name    string
		logging Logging
		wantErr bool
	}{
		{
			name: "Valid Logging struct",
			logging: Logging{
				Folder:     tempDir,
				Filename:   "app.log",
				MaxSizeMb:  10,
				MaxBackups: 3,
				MaxAgeDays: 7,
				Compress:   true,
				LogLevel:   "info",
			},
			wantErr: false,
		},
		{
			name: "Valid Logging level warn",
			logging: Logging{
				Folder:     tempDir,
				Filename:   "app.log",
				MaxSizeMb:  10,
				MaxBackups: 3,
				MaxAgeDays: 7,
				Compress:   true,
				LogLevel:   "warn",
			},
			wantErr: false,
		},
		{
			name: "Invalid Logging level",
			logging: Logging{
				Folder:     tempDir,
				Filename:   "app.log",
				MaxSizeMb:  10,
				MaxBackups: 3,
				MaxAgeDays: 7,
				Compress:   true,
				LogLevel:   "bogus",
			},
			wantErr: true,
		},
		{
			name: "Missing Folder",
			logging: Logging{
				Filename:   "app.log",
				MaxSizeMb:  10,
				MaxBackups: 3,
				MaxAgeDays: 7,
				Compress:   true,
				LogLevel:   "info",
			},
			wantErr: true,
		},
		{
			name: "Non-existent Folder",
			logging: Logging{
				Folder:     "/non/existent/path",
				Filename:   "app.log",
				MaxSizeMb:  10,
				MaxBackups: 3,
				MaxAgeDays: 7,
				Compress:   true,
				LogLevel:   "info",
			},
			wantErr: true,
		},
		{
			name: "Missing Filename",
			logging: Logging{
				Folder:     tempDir,
				MaxSizeMb:  10,
				MaxBackups: 3,
				MaxAgeDays: 7,
				Compress:   true,
				LogLevel:   "info",
			},
			wantErr: true,
		},
		{
			name: "Filename without .log extension",
			logging: Logging{
				Folder:     tempDir,
				Filename:   "app.txt",
				MaxSizeMb:  10,
				MaxBackups: 3,
				MaxAgeDays: 7,
				Compress:   true,
				LogLevel:   "info",
			},
			wantErr: true,
		},
		{
			name: "MaxSizeMb is zero",
			logging: Logging{
				Folder:     tempDir,
				Filename:   "app.log",
				MaxSizeMb:  0,
				MaxBackups: 3,
				MaxAgeDays: 7,
				Compress:   true,
				LogLevel:   "info",
			},
			wantErr: true,
		},
		{
			name: "MaxBackups is negative",
			logging: Logging{
				Folder:     tempDir,
				Filename:   "app.log",
				MaxSizeMb:  10,
				MaxBackups: -1,
				MaxAgeDays: 7,
				Compress:   true,
				LogLevel:   "info",
			},
			wantErr: true,
		},
		{
			name: "MaxAgeDays is negative",
			logging: Logging{
				Folder:     tempDir,
				Filename:   "app.log",
				MaxSizeMb:  10,
				MaxBackups: 3,
				MaxAgeDays: -1,
				Compress:   true,
				LogLevel:   "info",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate.Struct(tt.logging) // or any struct being validated
			checkValidationErrors(t, tt.name, err, tt.wantErr)
		})
	}
}
