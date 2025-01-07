package app

import (
	"bytes"
	"log/slog"
	"testing"
)

func TestLoggingMethods(t *testing.T) {
	var bufferFile bytes.Buffer
	var bufferProg bytes.Buffer

	fileHandler := slog.NewJSONHandler(&bufferFile, &slog.HandlerOptions{Level: slog.LevelDebug})
	progHandler := slog.NewJSONHandler(&bufferProg, &slog.HandlerOptions{Level: slog.LevelDebug})

	fileLogger := slog.New(fileHandler)
	progLogger := slog.New(progHandler)

	k := &KhedraApp{
		fileLogger: fileLogger,
		progLogger: progLogger,
	}

	tests := []struct {
		name        string
		logFunc     func(string, ...any)
		message     string
		buffer      *bytes.Buffer
		expectInLog bool
	}{
		{"Debug", k.Debug, "debug message", &bufferFile, true},
		{"Info", k.Info, "info message", &bufferProg, true},
		{"Warn", k.Warn, "warn message", &bufferProg, true},
		{"Error", k.Error, "error message", &bufferProg, true},
		{"ProgNoNewline", func(msg string, v ...any) { k.Prog(msg) }, "prog message", &bufferProg, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bufferFile.Reset()
			bufferProg.Reset()

			tt.logFunc(tt.message)

			if tt.expectInLog {
				logOutput := tt.buffer.String()
				expected := `"msg":"` + tt.message + `"`
				if !bytes.Contains([]byte(logOutput), []byte(expected)) {
					t.Fatalf("expected log message not found in: %s", logOutput)
				}
			}
		})
	}
}
