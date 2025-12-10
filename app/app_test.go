package app

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/TrueBlocks/trueblocks-khedra/v6/pkg/types"
)

// ============================================================================
// Service Initialization Tests
// ============================================================================

func TestInitializeControlSvc_CreatesServiceManager(t *testing.T) {
	cfg := types.NewConfig()
	// Disable services that bind to ports to avoid conflicts
	api := cfg.Services["api"]
	api.Enabled = false
	cfg.Services["api"] = api
	ipfs := cfg.Services["ipfs"]
	ipfs.Enabled = false
	cfg.Services["ipfs"] = ipfs
	app := &KhedraApp{config: &cfg}

	err := app.initializeControlSvc()
	require.NoError(t, err, "Should initialize control service successfully")

	assert.NotNil(t, app.controlSvc, "Control service should be initialized")
	assert.NotNil(t, app.serviceManager, "Service manager should be initialized")
	assert.NotNil(t, app.logger, "Logger should be initialized")
	assert.NotNil(t, app.config, "Config should be initialized")
}

func TestInitializeControlSvc_Idempotent(t *testing.T) {
	cfg := types.NewConfig()
	// Disable services that bind to ports to avoid conflicts
	api := cfg.Services["api"]
	api.Enabled = false
	cfg.Services["api"] = api
	ipfs := cfg.Services["ipfs"]
	ipfs.Enabled = false
	cfg.Services["ipfs"] = ipfs
	app := &KhedraApp{config: &cfg}

	// First call
	err := app.initializeControlSvc()
	require.NoError(t, err)
	firstControlSvc := app.controlSvc
	firstServiceManager := app.serviceManager

	// Second call should return early
	err = app.initializeControlSvc()
	require.NoError(t, err)

	assert.Same(t, firstControlSvc, app.controlSvc, "Control service should be same instance")
	assert.Same(t, firstServiceManager, app.serviceManager, "Service manager should be same instance")
}

func TestInitializeControlSvc_WithExistingLogger(t *testing.T) {
	customLogger := types.NewLogger(types.Logging{Level: "debug"})
	app := &KhedraApp{
		logger: customLogger,
	}

	err := app.initializeControlSvc()
	require.NoError(t, err)

	assert.Same(t, customLogger, app.logger, "Should preserve existing logger")
}

func TestInitializeControlSvc_WithExistingConfig(t *testing.T) {
	customConfig := &types.Config{
		General: types.General{
			DataFolder: "/custom/path",
		},
		Services: map[string]types.Service{
			"scraper": {
				Name:      "scraper",
				Enabled:   true,
				Sleep:     20,
				BatchSize: 50,
			},
			"api": {
				Name:    "api",
				Enabled: true,
				Port:    0, // ephemeral port
			},
		},
		Chains: map[string]types.Chain{
			"mainnet": {
				Name:    "mainnet",
				Enabled: true,
				RPCs:    []string{"http://localhost:8545"},
			},
		},
	}
	app := &KhedraApp{
		config: customConfig,
	}

	err := app.initializeControlSvc()
	require.NoError(t, err)

	assert.Same(t, customConfig, app.config, "Should preserve existing config")
}

// ============================================================================
// RestartAllServices Tests
// ============================================================================

func TestRestartAllServices_NilServiceManager(t *testing.T) {
	app := &KhedraApp{
		logger: types.NewLogger(types.Logging{Level: "error"}),
	}

	err := app.RestartAllServices()
	assert.Error(t, err, "Should fail when service manager not initialized")
	assert.Contains(t, err.Error(), "service manager not initialized")
}

func TestRestartAllServices_WithServiceManager(t *testing.T) {
	cfg := types.NewConfig()
	// Disable services that bind to ports to avoid conflicts
	api := cfg.Services["api"]
	api.Enabled = false
	cfg.Services["api"] = api
	ipfs := cfg.Services["ipfs"]
	ipfs.Enabled = false
	cfg.Services["ipfs"] = ipfs
	app := &KhedraApp{config: &cfg}

	// Initialize control service which creates service manager
	err := app.initializeControlSvc()
	require.NoError(t, err)

	// Restart should not error (may have no restartable services, but that's OK)
	err = app.RestartAllServices()
	assert.NoError(t, err, "Restart should succeed even with minimal services")
}

// ============================================================================
// Service Manager Integration Tests
// ============================================================================

func TestServiceManager_ReceivesAllConfiguredServices(t *testing.T) {
	customConfig := &types.Config{
		Services: map[string]types.Service{
			"scraper": {Name: "scraper", Enabled: true, Sleep: 14, BatchSize: 100},
			"monitor": {Name: "monitor", Enabled: false},
			"api":     {Name: "api", Enabled: true, Port: 0}, // ephemeral port
			"ipfs":    {Name: "ipfs", Enabled: false},
		},
		Chains: map[string]types.Chain{
			"mainnet": {Name: "mainnet", Enabled: true, RPCs: []string{"http://localhost:8545"}},
		},
	}

	app := &KhedraApp{config: customConfig}
	err := app.initializeControlSvc()
	require.NoError(t, err)

	// Service manager should be created with services
	assert.NotNil(t, app.serviceManager, "Service manager should exist")

	// Verify restart works (exercises service manager functionality)
	err = app.RestartAllServices()
	assert.NoError(t, err, "Should be able to restart services")
}

func TestControlService_AttachedToServiceManager(t *testing.T) {
	cfg := types.NewConfig()
	// Disable services that bind to ports to avoid conflicts
	api := cfg.Services["api"]
	api.Enabled = false
	cfg.Services["api"] = api
	ipfs := cfg.Services["ipfs"]
	ipfs.Enabled = false
	cfg.Services["ipfs"] = ipfs
	app := &KhedraApp{config: &cfg}
	err := app.initializeControlSvc()
	require.NoError(t, err)

	// The control service should be attached to service manager
	// We can verify by checking restart doesn't panic
	assert.NotPanics(t, func() {
		_ = app.RestartAllServices()
	}, "Control service operations should not panic")
}
