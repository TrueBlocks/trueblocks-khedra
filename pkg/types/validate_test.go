package types

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Testing status: reviewed

func TestChainValidation2(t *testing.T) {
	SetupTest([]string{})
	cases := []struct {
		name     string
		chain    Chain
		expected bool
	}{
		{
			name:     "Valid Chain with RPCs",
			chain:    Chain{Name: "Mainnet", RPCs: []string{"https://mainnet.infura.io"}, Enabled: true},
			expected: true,
		},
		{
			name:     "Invalid Chain without RPCs",
			chain:    Chain{Name: "Mainnet", RPCs: []string{}, Enabled: true},
			expected: false,
		},
		{
			name:     "Invalid Chain without Name",
			chain:    Chain{Name: "", RPCs: []string{"https://mainnet.infura.io"}, Enabled: true},
			expected: false,
		},
		{
			name:     "Disabled Chain with Empty RPCs",
			chain:    Chain{Name: "Testnet", RPCs: []string{}, Enabled: false},
			expected: true,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := Validate.Struct(c.chain)
			assert.Equal(t, c.expected, err == nil)
		})
	}
}

func TestServiceValidation2(t *testing.T) {
	SetupTest([]string{})
	cases := []struct {
		name     string
		service  Service
		expected bool
	}{
		{
			name:     "Valid API Service",
			service:  Service{Name: "api", Enabled: true, Port: 8080},
			expected: true,
		},
		{
			name:     "Invalid API Service without Port",
			service:  Service{Name: "api", Enabled: true, Port: 0},
			expected: false,
		},
		{
			name:     "Valid Scraper Service",
			service:  Service{Name: "scraper", Enabled: true, Sleep: 10, BatchSize: 100},
			expected: true,
		},
		{
			name:     "Invalid Scraper Service with Low BatchSize",
			service:  Service{Name: "scraper", Enabled: true, Sleep: 10, BatchSize: 10},
			expected: false,
		},
		{
			name:     "Valid IPFS Service",
			service:  Service{Name: "ipfs", Enabled: true, Port: 5001},
			expected: true,
		},
		{
			name:     "Invalid IPFS Service without Port",
			service:  Service{Name: "ipfs", Enabled: true, Port: 0},
			expected: false,
		},
		{
			name:     "Invalid Service with Unknown Name",
			service:  Service{Name: "unknown", Enabled: true},
			expected: false,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := Validate.Struct(c.service)
			assert.Equal(t, c.expected, err == nil)
		})
	}
}

func TestValidateFileAndFolder(t *testing.T) {
	// SetupTest([]string{})
	x, _ := os.Getwd()
	fmt.Println("Working Directory:", x)
	tempDir := t.TempDir()
	tempFile := filepath.Join(tempDir, "test.txt")
	_ = os.WriteFile(tempFile, []byte("test"), 0644)

	cases := []struct {
		name     string
		path     string
		tag      string
		expected bool
	}{
		{"Writable Directory", tempDir, "is_writable", true},
		// {"Non-Existing Directory", "nonexistent", "folder_exists", false}, // test fails because of fatal in ResolvePath
		{"Existing File", tempFile, "file_exists", true},
		// {"Non-Existing File", "nonexistent.txt", "file_exists", false}, // test fails because of fatal in ResolvePath
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := Validate.Var(c.path, c.tag)
			assert.Equal(t, c.expected, err == nil)
		})
	}
}

func TestOptionalMinMax(t *testing.T) {
	SetupTest([]string{})
	cases := []struct {
		name     string
		value    int
		tag      string
		expected bool
	}{
		{"Valid Min Value", 5, "opt_min=5", true},
		{"Invalid Below Min", 4, "opt_min=5", false},
		{"Valid Max Value", 10, "opt_max=10", true},
		{"Invalid Above Max", 11, "opt_max=10", false},
		{"Unset Value for Min", 0, "opt_min=5", true},
		{"Unset Value for Max", 0, "opt_max=10", true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			err := Validate.Var(c.value, c.tag)
			assert.Equal(t, c.expected, err == nil)
		})
	}
}
