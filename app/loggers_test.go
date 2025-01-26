package app

import (
	"os"
	"os/exec"
	"testing"
)

// Testing status: reviewed

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
