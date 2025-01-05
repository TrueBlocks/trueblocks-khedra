package types

import (
	"fmt"
	"os"
	"testing"

	"github.com/go-playground/validator"
)

func TestServiceValidation(t *testing.T) {
	tests := []struct {
		name    string
		service Service
		wantErr bool
	}{
		{
			name: "Valid Service with Sleep positive",
			service: Service{
				Name:  "api",
				Port:  8080,
				Sleep: 5, // Valid Sleep
			},
			wantErr: false,
		},
		{
			name: "Valid Service with BatchSize at min",
			service: Service{
				Name:      "scraper",
				BatchSize: 50, // Minimum valid BatchSize
				Sleep:     1,
			},
			wantErr: false,
		},
		{
			name: "Valid Service with BatchSize at max",
			service: Service{
				Name:      "scraper",
				BatchSize: 10000, // Maximum valid BatchSize
				Sleep:     1,
			},
			wantErr: false,
		},
		{
			name: "Valid Service with Sleep unset (0)",
			service: Service{
				Name:  "api",
				Port:  8080,
				Sleep: 0, // Optional, no validation
			},
			wantErr: false,
		},
		{
			name: "Valid Service with all optional values unset (zero)",
			service: Service{
				Name: "api",
			},
			wantErr: true,
		},
		{
			name: "Valid Service with Port within range",
			service: Service{
				Name: "api",
				Port: 8080, // Valid Port
			},
			wantErr: false,
		},
		{
			name: "Invalid Service with Port below 1024",
			service: Service{
				Name: "api",
				Port: 100, // Invalid Port
			},
			wantErr: true,
		},
		{
			name: "Invalid Service with Port above 65535",
			service: Service{
				Name: "api",
				Port: 70000, // Invalid Port
			},
			wantErr: true,
		},
		{
			name: "Invalid Service with BatchSize below min",
			service: Service{
				Name:      "scraper",
				BatchSize: 40, // Invalid BatchSize
			},
			wantErr: true,
		},
		{
			name: "Invalid Service with BatchSize above max",
			service: Service{
				Name:      "scraper",
				BatchSize: 20000, // Invalid BatchSize
			},
			wantErr: true,
		},
		{
			name: "Valid Service with all fields set to valid values",
			service: Service{
				Name:      "api",
				Enabled:   true,
				Port:      8080,
				Sleep:     5,
				BatchSize: 500,
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate.Struct(tt.service)
			checkValidationErrors(t, tt.name, err, tt.wantErr)
		})
	}
}

func TestAPIServiceValidation(t *testing.T) {
	tests := []struct {
		name    string
		service Service
		wantErr bool
	}{
		{
			name: "Valid API service with Port",
			service: Service{
				Name:    "api",
				Enabled: true,
				Port:    8080,
			},
			wantErr: false,
		},
		{
			name: "Invalid API service without Port",
			service: Service{
				Name:    "api",
				Enabled: true,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate.Struct(tt.service)
			checkValidationErrors(t, tt.name, err, tt.wantErr)
		})
	}
}

func TestScraperServiceValidation(t *testing.T) {
	tests := []struct {
		name    string
		service Service
		wantErr bool
	}{
		{
			name: "Valid Scraper service with required fields",
			service: Service{
				Name:      "scraper",
				Enabled:   true,
				Sleep:     60,
				BatchSize: 500,
			},
			wantErr: false,
		},
		{
			name: "Invalid Scraper service without Sleep",
			service: Service{
				Name:      "scraper",
				Enabled:   true,
				BatchSize: 500,
			},
			wantErr: true,
		},
		{
			name: "Invalid Scraper service without BatchSize",
			service: Service{
				Name:    "scraper",
				Enabled: true,
				Sleep:   60,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate.Struct(tt.service)
			checkValidationErrors(t, tt.name, err, tt.wantErr)
		})
	}
}

func TestMonitorServiceValidation(t *testing.T) {
	tests := []struct {
		name    string
		service Service
		wantErr bool
	}{
		{
			name: "Valid Monitor service with required fields",
			service: Service{
				Name:      "monitor",
				Enabled:   true,
				Sleep:     60,
				BatchSize: 500,
			},
			wantErr: false,
		},
		{
			name: "Invalid Monitor service without Sleep",
			service: Service{
				Name:      "monitor",
				Enabled:   true,
				BatchSize: 500,
			},
			wantErr: true,
		},
		{
			name: "Invalid Monitor service without BatchSize",
			service: Service{
				Name:    "monitor",
				Enabled: true,
				Sleep:   60,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate.Struct(tt.service)
			checkValidationErrors(t, tt.name, err, tt.wantErr)
		})
	}
}

func TestIPFSServiceValidation(t *testing.T) {
	tests := []struct {
		name    string
		service Service
		wantErr bool
	}{
		{
			name: "Valid IPFS service with Port",
			service: Service{
				Name:    "ipfs",
				Enabled: true,
				Port:    5001,
			},
			wantErr: false,
		},
		{
			name: "Invalid IPFS service without Port",
			service: Service{
				Name:    "ipfs",
				Enabled: true,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate.Struct(tt.service)
			checkValidationErrors(t, tt.name, err, tt.wantErr)
		})
	}
}

func TestServiceListValidation(t *testing.T) {
	services := []Service{
		{
			Name:    "api",
			Enabled: true,
			Port:    8080,
		},
		{
			Name:      "scraper",
			Enabled:   true,
			Sleep:     60,
			BatchSize: 500,
		},
		{
			Name:      "monitor",
			Enabled:   true,
			Sleep:     60,
			BatchSize: 500,
		},
		{
			Name:    "ipfs",
			Enabled: true,
			Port:    5001,
		},
	}

	for i, service := range services {
		t.Run(fmt.Sprintf("Service %d Validation", i+1), func(t *testing.T) {
			err := Validate.Struct(service)
			if err != nil {
				t.Errorf("Validation failed for service %d: %v", i+1, err)
			}
		})
	}
}

// createTempDir creates a temporary directory for testing.
// If writable is false, it makes the directory non-writable.
func createTempDir(t *testing.T, writable bool) string {
	dir, err := os.MkdirTemp("", "test_general")
	if err != nil {
		t.Fatalf("Failed to create temp directory: %v", err)
	}

	if !writable {
		err := os.Chmod(dir, 0500) // Read and execute permissions only
		if err != nil {
			t.Fatalf("Failed to make directory non-writable: %v", err)
		}
	}

	return dir
}

func checkValidationErrors(t *testing.T, name string, err error, wantErr bool) {
	t.Helper() // Marks this function as a helper, so the line numbers in errors refer to the caller.

	if (err != nil) != wantErr {
		if err != nil {
			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				for _, fieldErr := range validationErrors {
					t.Errorf("Validation error in test '%s' on field '%s': tag='%s', param='%s', value='%v'",
						name,
						fieldErr.Field(),
						fieldErr.Tag(),
						fieldErr.Param(),
						fieldErr.Value(),
					)
				}
			} else {
				t.Errorf("Unexpected error in test '%s': %v", name, err)
			}
		} else {
			t.Errorf("Test '%s': expected error = %v, got error = %v", name, wantErr, err != nil)
		}
	}
}
