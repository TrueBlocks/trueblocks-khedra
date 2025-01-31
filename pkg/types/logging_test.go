package types

import (
	"log/slog"
	"os"
	"path/filepath"
	"testing"

	"github.com/alecthomas/assert/v2"
)

func TestLoggingNew(t *testing.T) {
	homeDir, err := os.UserHomeDir()
	assert.NoError(t, err)

	expectedFolder := filepath.Join(homeDir, ".khedra", "logs")

	logging := NewLogging()
	assert.Equal(t, expectedFolder, logging.Folder)
	assert.Equal(t, "khedra.log", logging.Filename)
	assert.False(t, logging.ToFile)
	assert.Equal(t, 10, logging.MaxSize)
	assert.Equal(t, 3, logging.MaxBackups)
	assert.Equal(t, 10, logging.MaxAge)
	assert.Equal(t, "info", logging.Level)
	assert.True(t, logging.Compress)
}

// func TestLoggingValidation(t *testing.T) {
// 	tempDir := createTempDir(t, true)

// 	tests := []struct {
// 		name    string
// 		logging Logging
// 		wantErr bool
// 	}{
// 		{
// 			name: "Valid Logging struct",
// 			logging: Logging{
// 				Folder:     tempDir,
// 				Filename:   "app.log",
// 				ToFile:     true,
// 				MaxSize:    10,
// 				MaxBackups: 3,
// 				MaxAge:     7,
// 				Compress:   true,
// 				Level:      "info",
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "Valid Logging level warn",
// 			logging: Logging{
// 				Folder:     tempDir,
// 				Filename:   "app.log",
// 				ToFile:     false,
// 				MaxSize:    10,
// 				MaxBackups: 3,
// 				MaxAge:     7,
// 				Compress:   true,
// 				Level:      "warn",
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "Invalid Logging level",
// 			logging: Logging{
// 				Folder:     tempDir,
// 				Filename:   "app.log",
// 				ToFile:     false,
// 				MaxSize:    10,
// 				MaxBackups: 3,
// 				MaxAge:     7,
// 				Compress:   true,
// 				Level:      "bogus",
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "Missing Folder",
// 			logging: Logging{
// 				Filename:   "app.log",
// 				ToFile:     false,
// 				MaxSize:    10,
// 				MaxBackups: 3,
// 				MaxAge:     7,
// 				Compress:   true,
// 				Level:      "info",
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "Non-existent Folder",
// 			logging: Logging{
// 				Folder:     "/non/existent/path",
// 				Filename:   "app.log",
// 				ToFile:     true,
// 				MaxSize:    10,
// 				MaxBackups: 3,
// 				MaxAge:     7,
// 				Compress:   true,
// 				Level:      "info",
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "Missing Filename",
// 			logging: Logging{
// 				Folder:     tempDir,
// 				MaxSize:    10,
// 				ToFile:     false,
// 				MaxBackups: 3,
// 				MaxAge:     7,
// 				Compress:   true,
// 				Level:      "info",
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "Filename without .log extension",
// 			logging: Logging{
// 				Folder:     tempDir,
// 				Filename:   "app.txt",
// 				ToFile:     false,
// 				MaxSize:    10,
// 				MaxBackups: 3,
// 				MaxAge:     7,
// 				Compress:   true,
// 				Level:      "info",
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "MaxSize is zero",
// 			logging: Logging{
// 				Folder:     tempDir,
// 				Filename:   "app.log",
// 				ToFile:     false,
// 				MaxSize:    0,
// 				MaxBackups: 3,
// 				MaxAge:     7,
// 				Compress:   true,
// 				Level:      "info",
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "MaxBackups is negative",
// 			logging: Logging{
// 				Folder:     tempDir,
// 				Filename:   "app.log",
// 				ToFile:     false,
// 				MaxSize:    10,
// 				MaxBackups: -1,
// 				MaxAge:     7,
// 				Compress:   true,
// 				Level:      "info",
// 			},
// 			wantErr: true,
// 		},
// 		{
// 			name: "MaxAge is negative",
// 			logging: Logging{
// 				Folder:     tempDir,
// 				Filename:   "app.log",
// 				ToFile:     true,
// 				MaxSize:    10,
// 				MaxBackups: 3,
// 				MaxAge:     -1,
// 				Compress:   true,
// 				Level:      "info",
// 			},
// 			wantErr: true,
// 		},
// 	}

// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			err := validate.Validate(&tt.logging)
// 			if tt.wantErr {
// 				assert.Error(t, err, "Expected error for test case '%s'", tt.name)
// 			} else {
// 				assert.NoError(t, err, "Did not expect error for test case '%s'", tt.name)
// 			}
// 		})
// 	}
// }

func TestLoggingReadAndWrite(t *testing.T) {
	tempFilePath := "temp_config.yaml"
	content := `
      folder: ~/.khedra/logs
      filename: khedra.log
      toFile: false
      maxSize: 10
      maxBackups: 3
      maxAge: 10
      compress: true
      level: debug
    `

	assertions := func(t *testing.T, logging *Logging) {
		assert.Equal(t, "~/.khedra/logs", logging.Folder)
		assert.Equal(t, "khedra.log", logging.Filename)
		assert.False(t, logging.ToFile)
		assert.Equal(t, "debug", logging.Level)
		assert.Equal(t, 10, logging.MaxSize)
		assert.Equal(t, 3, logging.MaxBackups)
		assert.Equal(t, 10, logging.MaxAge)
		assert.True(t, logging.Compress)
	}

	ReadAndWriteWithAssertions[Logging](t, tempFilePath, content, assertions)
}

func TestConvertLevelUnsupported(t *testing.T) {
	level := convertLevel("unsupported")
	if level != slog.LevelInfo {
		t.Errorf("expected fallback to DefaultLevel, got %v", level)
	}
}
