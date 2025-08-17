package types

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

// Testing status: reviewed

func TestNewService(t *testing.T) {
	tests := []struct {
		name        string
		serviceType string
		expected    Service
		shouldPanic bool
	}{
		{
			name:        "Create scraper service",
			serviceType: "scraper",
			expected: Service{
				Name:      "scraper",
				Enabled:   true,
				Sleep:     10,
				BatchSize: 500,
			},
		},
		{
			name:        "Create monitor service",
			serviceType: "monitor",
			expected: Service{
				Name:      "monitor",
				Enabled:   false,
				Sleep:     12,
				BatchSize: 500,
			},
		},
		{
			name:        "Create API service",
			serviceType: "api",
			expected: Service{
				Name:    "api",
				Enabled: true,
				Port:    8080,
			},
		},
		{
			name:        "Create IPFS service",
			serviceType: "ipfs",
			expected: Service{
				Name:    "ipfs",
				Enabled: true,
				Port:    5001,
			},
		},
		{
			name:        "Unknown service type",
			serviceType: "unknown",
			shouldPanic: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.shouldPanic {
				assert.Panics(t, func() {
					NewService(tt.serviceType)
				}, "Expected panic for unknown service type")
			} else {
				service := NewService(tt.serviceType)
				assert.Equal(t, tt.expected, service)

				// Validate the returned service
				err := Validate(&service)
				assert.NoError(t, err, "Validation failed for service: %v", service)
			}
		})
	}
}

func TestServiceValidationUnified(t *testing.T) {
	tests := []struct {
		name    string
		service Service
		wantErr bool
	}{
		{
			name: "Valid API service with Port",
			service: Service{
				Name:    "api",
				Port:    8080,
				Enabled: true,
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
			name: "Invalid Monitor service without BatchSize",
			service: Service{
				Name:    "monitor",
				Enabled: true,
				Sleep:   60,
			},
			wantErr: true,
		},
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
		{
			name: "Valid Scraper with BatchSize at min",
			service: Service{
				Name:      "scraper",
				BatchSize: 50, // Minimum valid BatchSize
				Sleep:     1,
				Enabled:   true,
			},
			wantErr: false,
		},
		{
			name: "Invalid Scraper with BatchSize below min",
			service: Service{
				Name:      "scraper",
				BatchSize: 40, // Invalid BatchSize
				Enabled:   true,
			},
			wantErr: true,
		},
		{
			name: "Invalid Scraper with BatchSize below min (disabled)",
			service: Service{
				Name:      "scraper",
				BatchSize: 40, // Invalid BatchSize
				Enabled:   false,
			},
			wantErr: false,
		},
		{
			name: "Valid Scraper with BatchSize at max",
			service: Service{
				Name:      "scraper",
				BatchSize: 10000, // Maximum valid BatchSize
				Sleep:     1,
				Enabled:   true,
			},
			wantErr: false,
		},
		{
			name: "Invalid API service with Port above 65535",
			service: Service{
				Name:    "api",
				Port:    70000, // Invalid Port
				Enabled: true,
			},
			wantErr: true,
		},
		{
			name: "Invalid API service with Port above 65535 (disabled)",
			service: Service{
				Name:    "api",
				Port:    70000, // Invalid Port
				Enabled: false,
			},
			wantErr: false,
		},
		{
			name: "Valid Service with Sleep positive",
			service: Service{
				Name:    "api",
				Port:    8080,
				Sleep:   5, // Valid Sleep
				Enabled: true,
			},
			wantErr: false,
		},
		{
			name: "Valid Service with Sleep unset (0)",
			service: Service{
				Name:    "api",
				Port:    8080,
				Sleep:   0, // Optional, no validation
				Enabled: true,
			},
			wantErr: false,
		},
		{
			name: "Valid Service with all optional values unset (zero)",
			service: Service{
				Name:    "api",
				Enabled: true,
			},
			wantErr: true,
		},
		{
			name: "Valid Service with Port within range",
			service: Service{
				Name:    "api",
				Port:    8080, // Valid Port
				Enabled: true,
			},
			wantErr: false,
		},
		{
			name: "Invalid Service with Port below 1024",
			service: Service{
				Name:    "api",
				Port:    100, // Invalid Port
				Enabled: true,
			},
			wantErr: true,
		},
		{
			name: "Invalid Service with BatchSize above max",
			service: Service{
				Name:      "scraper",
				BatchSize: 20000, // Invalid BatchSize
				Enabled:   true,
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
			err := Validate(&tt.service)
			if tt.wantErr {
				assert.Error(t, err, "Expected an error in test '%s', but got none", tt.name)
			} else {
				assert.NoError(t, err, "Unexpected error in test '%s': %v", tt.name, err)
			}
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
			err := Validate(&service)
			assert.NoError(t, err, "Validation failed for service %d: %v", i+1, err)
		})
	}
}

// Added minimal missing coverage per ai/TestDesign_service.go.md
func TestService_IsEnabled(t *testing.T) {
	api := NewService("api")
	assert.True(t, api.IsEnabled(), "api service should start enabled")
	api.Enabled = false
	assert.False(t, api.IsEnabled(), "api service should reflect disabled state")

	monitor := NewService("monitor") // constructor sets Enabled false
	assert.False(t, monitor.IsEnabled(), "monitor service should start disabled")
	monitor.Enabled = true
	assert.True(t, monitor.IsEnabled(), "monitor service should reflect enabled state after change")
}

func TestServiceValidation_UnknownNameOneOf(t *testing.T) {
	// Directly construct a Service with an invalid Name to trigger oneof failure
	s := Service{Name: "random", Enabled: true}
	err := Validate(&s)
	assert.Error(t, err, "expected validation error for unknown service name")
	if err != nil {
		assert.Contains(t, err.Error(), "service_field", "should include custom validator tag in error")
		// The struct tag includes both 'required' and 'oneof', ensure at least a hint of the name issue
		// We expect multiple service_field validator failures referencing unknown service name
		assert.Contains(t, err.Error(), "unknown service name", "should note unknown service name in custom validator failures")
	}
}
