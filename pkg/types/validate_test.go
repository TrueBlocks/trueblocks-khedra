package types

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// Testing status: not_reviewed
// Integration tests for the validation system using real config scenarios

// ============================================================================
// Test 1: Chain Validation
// ============================================================================

func TestValidateChain_DisabledChain(t *testing.T) {
	defer SetupTest([]string{})()
	// Disabled chain should pass even with empty RPCs
	chain := Chain{
		Name:    "testnet",
		RPCs:    []string{},
		ChainID: 5,
		Enabled: false,
	}
	err := Validate(&chain)
	assert.NoError(t, err, "Disabled chain should pass validation even without RPCs")
}

func TestValidateChain_EnabledWithValidRPC(t *testing.T) {
	defer SetupTest([]string{})()
	chain := Chain{
		Name:    "mainnet",
		RPCs:    []string{"http://localhost:8545"},
		ChainID: 1,
		Enabled: true,
	}
	err := Validate(&chain)
	assert.NoError(t, err, "Enabled chain with valid RPC should pass")
}

func TestValidateChain_EnabledWithoutName(t *testing.T) {
	defer SetupTest([]string{})()
	chain := Chain{
		Name:    "",
		RPCs:    []string{"http://localhost:8545"},
		ChainID: 1,
		Enabled: true,
	}
	err := Validate(&chain)
	assert.Error(t, err, "Enabled chain without name should fail")
	assert.Contains(t, err.Error(), "Name", "Error should mention Name field")
}

func TestValidateChain_EnabledWithoutRPC(t *testing.T) {
	defer SetupTest([]string{})()
	chain := Chain{
		Name:    "mainnet",
		RPCs:    []string{},
		ChainID: 1,
		Enabled: true,
	}
	err := Validate(&chain)
	assert.Error(t, err, "Enabled chain without RPC should fail")
	assert.Contains(t, err.Error(), "RPCs", "Error should mention RPCs field")
}

func TestValidateChain_EnabledWithInvalidURL(t *testing.T) {
	defer SetupTest([]string{})()
	chain := Chain{
		Name:    "mainnet",
		RPCs:    []string{"not a valid url"},
		ChainID: 1,
		Enabled: true,
	}
	err := Validate(&chain)
	assert.Error(t, err, "Enabled chain with invalid RPC URL should fail")
}

func TestValidateChain_ChainIDZero(t *testing.T) {
	defer SetupTest([]string{})()
	chain := Chain{
		Name:    "mainnet",
		RPCs:    []string{"http://localhost:8545"},
		ChainID: 0,
		Enabled: true,
	}
	err := Validate(&chain)
	assert.Error(t, err, "Chain with zero ChainID should fail")
	assert.Contains(t, err.Error(), "ChainID", "Error should mention ChainID field")
}

func TestValidateChain_ChainIDPositive(t *testing.T) {
	defer SetupTest([]string{})()
	chain := Chain{
		Name:    "mainnet",
		RPCs:    []string{"http://localhost:8545"},
		ChainID: 1,
		Enabled: true,
	}
	err := Validate(&chain)
	assert.NoError(t, err, "Chain with positive ChainID should pass")
}

func TestValidateChain_DisabledChainIgnoresReqIfEnabled(t *testing.T) {
	defer SetupTest([]string{})()
	chain := Chain{
		Name:    "",         // Would fail if enabled
		RPCs:    []string{}, // Would fail if enabled
		ChainID: 1,
		Enabled: false,
	}
	err := Validate(&chain)
	// Disabled chain should skip req_if_enabled validators
	assert.NoError(t, err, "Disabled chain should skip req_if_enabled validators")
}

// ============================================================================
// Test 2: Service Validation
// ============================================================================

func TestValidateService_APIServiceWithValidPort(t *testing.T) {
	defer SetupTest([]string{})()
	svc := Service{
		Name:    "api",
		Enabled: true,
		Port:    8080,
	}
	// Test with Config root (as it would be during config validation)
	cfg := NewConfig()
	cfg.Services["api"] = svc
	err := Validate(&cfg)
	assert.NoError(t, err, "API service with valid port should pass")
}

func TestValidateService_APIServicePortTooLow(t *testing.T) {
	defer SetupTest([]string{})()
	svc := Service{
		Name:    "api",
		Enabled: true,
		Port:    80, // Below 1024
	}
	cfg := NewConfig()
	cfg.Services["api"] = svc
	err := Validate(&cfg)
	assert.Error(t, err, "API service with port < 1024 should fail")
	assert.Contains(t, err.Error(), "Port", "Error should mention Port")
}

func TestValidateService_APIServicePortTooHigh(t *testing.T) {
	defer SetupTest([]string{})()
	svc := Service{
		Name:    "api",
		Enabled: true,
		Port:    70000, // Above 65535
	}
	cfg := NewConfig()
	cfg.Services["api"] = svc
	err := Validate(&cfg)
	assert.Error(t, err, "API service with port > 65535 should fail")
	assert.Contains(t, err.Error(), "Port", "Error should mention Port")
}

func TestValidateService_ScraperWithValidSleep(t *testing.T) {
	defer SetupTest([]string{})()
	svc := Service{
		Name:      "scraper",
		Enabled:   true,
		Sleep:     10,
		BatchSize: 500,
	}
	cfg := NewConfig()
	cfg.Services["scraper"] = svc
	err := Validate(&cfg)
	assert.NoError(t, err, "Scraper with valid sleep and batch size should pass")
}

func TestValidateService_ScraperWithInvalidSleep(t *testing.T) {
	defer SetupTest([]string{})()
	svc := Service{
		Name:      "scraper",
		Enabled:   true,
		Sleep:     0, // Must be positive
		BatchSize: 500,
	}
	cfg := NewConfig()
	cfg.Services["scraper"] = svc
	err := Validate(&cfg)
	assert.Error(t, err, "Scraper with sleep <= 0 should fail")
	assert.Contains(t, err.Error(), "Sleep", "Error should mention Sleep")
}

func TestValidateService_ScraperWithBatchSizeTooSmall(t *testing.T) {
	defer SetupTest([]string{})()
	svc := Service{
		Name:      "scraper",
		Enabled:   true,
		Sleep:     10,
		BatchSize: 25, // Below 50
	}
	cfg := NewConfig()
	cfg.Services["scraper"] = svc
	err := Validate(&cfg)
	assert.Error(t, err, "Scraper with batch size < 50 should fail")
	assert.Contains(t, err.Error(), "BatchSize", "Error should mention BatchSize")
}

func TestValidateService_ScraperWithBatchSizeTooLarge(t *testing.T) {
	defer SetupTest([]string{})()
	svc := Service{
		Name:      "scraper",
		Enabled:   true,
		Sleep:     10,
		BatchSize: 15000, // Above 10000
	}
	cfg := NewConfig()
	cfg.Services["scraper"] = svc
	err := Validate(&cfg)
	assert.Error(t, err, "Scraper with batch size > 10000 should fail")
	assert.Contains(t, err.Error(), "BatchSize", "Error should mention BatchSize")
}

func TestValidateService_DisabledServiceIgnoresFieldValidation(t *testing.T) {
	defer SetupTest([]string{})()
	svc := Service{
		Name:      "scraper",
		Enabled:   false,
		Sleep:     -5, // Invalid but should be ignored
		BatchSize: 1,  // Invalid but should be ignored
	}
	cfg := NewConfig()
	cfg.Services["scraper"] = svc
	err := Validate(&cfg)
	assert.NoError(t, err, "Disabled service should not validate its fields")
}

func TestValidateService_InvalidServiceName(t *testing.T) {
	defer SetupTest([]string{})()
	svc := Service{
		Name:    "invalid",
		Enabled: true,
		Port:    8080,
	}
	err := Validate(&svc)
	assert.Error(t, err, "Service with invalid name should fail")
	assert.Contains(t, err.Error(), "invalid", "Error should mention the invalid service name")
}

// ============================================================================
// Test 3: Config Validation
// ============================================================================

func TestValidateConfig_ValidDefaultConfig(t *testing.T) {
	defer SetupTest([]string{})()
	cfg := NewConfig()
	err := Validate(&cfg)
	assert.NoError(t, err, "Default config should be valid")
}

func TestValidateConfig_NoEnabledServices(t *testing.T) {
	defer SetupTest([]string{})()
	cfg := NewConfig()
	// Disable all services
	for name, svc := range cfg.Services {
		svc.Enabled = false
		cfg.Services[name] = svc
	}
	err := Validate(&cfg)
	assert.NoError(t, err, "Validation framework should pass; config.go helper validates service count")
}

func TestValidateConfig_NoEnabledChains(t *testing.T) {
	defer SetupTest([]string{})()
	cfg := NewConfig()
	// Disable all chains
	for name, chain := range cfg.Chains {
		chain.Enabled = false
		cfg.Chains[name] = chain
	}
	err := Validate(&cfg)
	assert.NoError(t, err, "Validation framework should pass; config.go helper validates chain count")
}

func TestValidateConfig_MainnetMissingRPC(t *testing.T) {
	defer SetupTest([]string{})()
	cfg := NewConfig()
	chain := cfg.Chains["mainnet"]
	chain.RPCs = []string{}
	cfg.Chains["mainnet"] = chain
	err := Validate(&cfg)
	assert.Error(t, err, "Mainnet without RPC should fail validation")
}

func TestValidateConfig_MultipleChains(t *testing.T) {
	defer SetupTest([]string{})()
	cfg := NewConfig()
	cfg.Chains["sepolia"] = Chain{
		Name:    "sepolia",
		RPCs:    []string{"http://localhost:8545"},
		ChainID: 11155111,
		Enabled: true,
	}
	err := Validate(&cfg)
	assert.NoError(t, err, "Config with multiple chains should validate")
}

func TestValidateConfig_EnabledChainWithInvalidRPC(t *testing.T) {
	defer SetupTest([]string{})()
	cfg := NewConfig()
	chain := cfg.Chains["mainnet"]
	chain.RPCs = []string{"not a url"}
	cfg.Chains["mainnet"] = chain
	err := Validate(&cfg)
	assert.Error(t, err, "Enabled chain with invalid RPC should fail")
}

// ============================================================================
// Test 4: req_if_enabled Behavior
// ============================================================================

func TestReqIfEnabled_DisabledChainMissingName(t *testing.T) {
	defer SetupTest([]string{})()
	chain := Chain{
		Name:    "",
		RPCs:    []string{"http://localhost:8545"},
		ChainID: 1,
		Enabled: false,
	}
	err := Validate(&chain)
	assert.NoError(t, err, "Disabled chain with missing name should pass due to req_if_enabled")
}

func TestReqIfEnabled_EnabledChainMissingName(t *testing.T) {
	defer SetupTest([]string{})()
	chain := Chain{
		Name:    "",
		RPCs:    []string{"http://localhost:8545"},
		ChainID: 1,
		Enabled: true,
	}
	err := Validate(&chain)
	assert.Error(t, err, "Enabled chain with missing name should fail")
}

func TestReqIfEnabled_DisabledServiceMissingPort(t *testing.T) {
	defer SetupTest([]string{})()
	svc := Service{
		Name:    "api",
		Enabled: false,
		Port:    0, // Invalid but ignored when disabled
	}
	cfg := NewConfig()
	cfg.Services["api"] = svc
	err := Validate(&cfg)
	assert.NoError(t, err, "Disabled service should not require port validation")
}

func TestReqIfEnabled_EnabledServiceMissingPort(t *testing.T) {
	defer SetupTest([]string{})()
	svc := Service{
		Name:    "api",
		Enabled: true,
		Port:    0, // Must be in valid range
	}
	cfg := NewConfig()
	cfg.Services["api"] = svc
	err := Validate(&cfg)
	assert.Error(t, err, "Enabled API service should require valid port")
}

// ============================================================================
// Test 5: Error Messages and Context
// ============================================================================

func TestValidationError_MultipleErrors(t *testing.T) {
	defer SetupTest([]string{})()
	cfg := NewConfig()
	// Break multiple things
	chain := cfg.Chains["mainnet"]
	chain.Name = ""         // Missing name (when enabled)
	chain.RPCs = []string{} // Missing RPC (when enabled)
	chain.ChainID = 0       // Missing ChainID
	cfg.Chains["mainnet"] = chain

	err := Validate(&cfg)
	assert.Error(t, err, "Multiple validation errors should be aggregated")
	errStr := err.Error()
	// All errors should be in the message
	assert.Contains(t, errStr, "Name", "Should report Name error")
	assert.Contains(t, errStr, "RPCs", "Should report RPCs error")
	assert.Contains(t, errStr, "ChainID", "Should report ChainID error")
}

func TestValidationError_PreservesContext(t *testing.T) {
	defer SetupTest([]string{})()
	cfg := NewConfig()
	// Create invalid service in map
	cfg.Services["api"] = Service{
		Name:    "api",
		Enabled: true,
		Port:    500, // Too low
	}
	err := Validate(&cfg)
	assert.Error(t, err)
	errStr := err.Error()
	// Should mention the field and context
	assert.Contains(t, errStr, "Port", "Error should mention Port field")
	assert.Contains(t, errStr, "api", "Error should mention api service context")
}
