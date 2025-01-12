package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"testing"
)

// Testing status: reviewed

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
		// {"ProgWithNewline", func(msg string, v ...any) { k.Prog(msg + "\n") }, "prog message\n", &bufferProg, true},
		{"ProgWithArgs", func(msg string, v ...any) { k.Prog(msg, "arg1", 42) }, "prog message", &bufferProg, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			bufferFile.Reset()
			bufferProg.Reset()

			tt.logFunc(tt.message)

			if tt.expectInLog {
				logOutput := tt.buffer.String()
				fmt.Printf("DEBUG: Log Output: %s\n", logOutput)
				var logEntry map[string]any
				err := json.Unmarshal([]byte(logOutput), &logEntry)
				if err != nil {
					t.Fatalf("failed to parse log output: %v", err)
				}
				if logEntry["msg"] != tt.message {
					t.Errorf("expected message %q, got %q", tt.message, logEntry["msg"])
				}

				// Check for arguments in log output (if any)
				if tt.name == "ProgWithArgs" {
					if !bytes.Contains([]byte(logOutput), []byte("arg1")) || !bytes.Contains([]byte(logOutput), []byte("42")) {
						t.Fatalf("expected arguments not found in log: %s", logOutput)
					}
				}
			}
		})
	}
}

func TestFatal(t *testing.T) {
	if os.Getenv("TEST_FATAL") == "1" {
		k := &KhedraApp{}
		k.Fatal("fatal message")
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestFatal")
	cmd.Env = append(os.Environ(), "TEST_FATAL=1")
	err := cmd.Run()
	if err == nil || err.Error() != "exit status 1" {
		t.Fatalf("expected Fatal to exit with status 1, got %v", err)
	}
}
