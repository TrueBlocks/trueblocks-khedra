# Command-Line Interface

Khedra provides a simple, focused command-line interface (CLI) for managing the system. The CLI is designed around a core workflow of initialization, daemon management, and service control via REST API.

## CLI Architecture

The CLI is built using the `urfave/cli` library, providing a lightweight and user-friendly interface focused on essential operations.

### Design Principles

1. **Simplicity**: Minimal command set focused on core functionality
2. **Clarity**: Each command has a clear, single purpose
3. **Automation-Friendly**: Services are controlled via REST API for scriptability
4. **Self-Documenting**: Built-in help for all commands

### Implementation Structure

```go
func initCli(k *KhedraApp) *cli.App {
    app := &cli.App{
        Name:     "khedra",
        Usage:    "A tool to index, monitor, serve, and share blockchain data",
        Version:  version.Version,
        Commands: []*cli.Command{
            // Core commands only - no complex service management via CLI
        },
    }
    return app
}
```

## Command Structure

Khedra's CLI is organized into logical command groups.

### Root Command

The base `khedra` command serves as the entry point and provides global options:

```
khedra [global options] command [command options] [arguments...]
```

Global options include:
- `--config`: Specify an alternate configuration file path
- `--format`: Control output format (text, json, csv)
- `--verbose`: Enable verbose output
- `--quiet`: Suppress non-error output
- `--chain`: Specify the target blockchain (defaults to "mainnet")

### Service Management Commands

Commands for controlling Khedra's services:

#### Service Lifecycle Commands

- `khedra daemon`: Start Khedra in daemon mode with all configured services
  ```bash
  # Start with default configuration
  khedra daemon
  
  # Start with specific log level
  TB_KHEDRA_LOGGING_LEVEL=debug khedra daemon
  
  # Start with custom configuration file
  khedra daemon --config=/path/to/config.yaml
  ```

- `khedra start [service...]`: Start specific services (if supported)
  ```bash
  # Start all services
  khedra start
  
  # Start specific services
  khedra start scraper api
  
  # Start in foreground mode
  khedra start --foreground
  ```

- `khedra stop [service...]`: Stop specific services (if supported)
  ```bash
  # Stop all services
  khedra stop
  
  # Stop specific services  
  khedra stop scraper monitor
  
  # Force stop (ungraceful shutdown)
  khedra stop --force
  ```

- `khedra restart [service...]`: Restart specific services (if supported)
  ```bash
  # Restart all services
  khedra restart
  
  # Restart specific services
  khedra restart api ipfs
  
  # Restart with configuration reload
  khedra restart --reload-config
  ```

#### Service Status Commands

- `khedra status`: Show status of all services
  ```bash
  # Basic status
  khedra status
  
  # Detailed status with metrics
  khedra status --verbose
  
  # Status for specific services
  khedra status scraper monitor
  
  # JSON output for scripting
  khedra status --format=json
  ```

#### Service Control via API

All service management operations are also available via REST API:

```bash
# Check service status
curl http://localhost:8080/api/v1/services

# Start a service
curl -X POST http://localhost:8080/api/v1/services/scraper/start

# Stop a service  
curl -X POST http://localhost:8080/api/v1/services/monitor/stop

# Get detailed service metrics
curl http://localhost:8080/api/v1/services/api?verbose=true
```

### Index Management Commands

Commands for managing the Unchained Index:

#### Index Status and Information

- `khedra index status [chain]`: Show index status for specified chain
  ```bash
  # Status for default chain
  khedra index status
  
  # Status for specific chain
  khedra index status --chain=mainnet
  
  # Show gaps in the index
  khedra index status --show-gaps
  
  # Show detailed analytics
  khedra index status --analytics
  
  # Export status to file
  khedra index status --output=/path/to/status.json
  ```

#### Index Maintenance Commands

- `khedra index rebuild [options]`: Rebuild portions of the index
  ```bash
  # Rebuild specific block range
  khedra index rebuild --start=18000000 --end=18001000
  
  # Rebuild from specific block to latest
  khedra index rebuild --from=18000000
  
  # Rebuild with specific batch size
  khedra index rebuild --batch-size=1000 --start=18000000 --end=18001000
  
  # Rebuild for specific chain
  khedra index rebuild --chain=sepolia --start=4000000 --end=4001000
  ```

- `khedra index verify [options]`: Verify index integrity
  ```bash
  # Verify entire index
  khedra index verify
  
  # Verify specific block range
  khedra index verify --start=18000000 --end=18001000
  
  # Verify and attempt to repair issues
  khedra index verify --repair
  
  # Verify for specific chain
  khedra index verify --chain=mainnet
  ```

- `khedra index optimize`: Optimize index storage and performance
  ```bash
  # Optimize entire index
  khedra index optimize
  
  # Optimize specific chunks
  khedra index optimize --chunks=1800,1801,1802
  
  # Optimize and compress
  khedra index optimize --compress
  ```

### Monitor Commands

Commands for managing address monitors:

#### Monitor Management

- `khedra monitor add ADDRESS [ADDRESS...]`: Add addresses to monitor
  ```bash
  # Add single address
  khedra monitor add 0x742d35Cc6634C0532925a3b844Bc454e4438f44e
  
  # Add multiple addresses
  khedra monitor add 0x742d35Cc... 0x1234567... 0xabcdef...
  
  # Add with custom name
  khedra monitor add 0x742d35Cc... --name="Vitalik Buterin"
  
  # Add with notification settings
  khedra monitor add 0x742d35Cc... --notifications=webhook,email
  
  # Add from file
  khedra monitor add --file=/path/to/addresses.txt
  ```

- `khedra monitor remove ADDRESS [ADDRESS...]`: Remove monitored addresses
  ```bash
  # Remove specific addresses
  khedra monitor remove 0x742d35Cc6634C0532925a3b844Bc454e4438f44e
  
  # Remove without confirmation prompt
  khedra monitor remove 0x742d35Cc... --force
  
  # Remove all monitors (with confirmation)
  khedra monitor remove --all
  
  # Remove monitors matching pattern
  khedra monitor remove --pattern="test_*"
  ```

#### Monitor Information

- `khedra monitor list`: List all monitored addresses
  ```bash
  # List all monitors
  khedra monitor list
  
  # List with detailed information
  khedra monitor list --details
  
  # List only active monitors
  khedra monitor list --active-only
  
  # Export to CSV
  khedra monitor list --format=csv --output=monitors.csv
  
  # Filter by pattern
  khedra monitor list --pattern="*buterin*"
  ```

- `khedra monitor activity ADDRESS`: Show activity for monitored address
  ```bash
  # Show recent activity
  khedra monitor activity 0x742d35Cc6634C0532925a3b844Bc454e4438f44e
  
  # Show activity in block range
  khedra monitor activity 0x742d35Cc... --from=18000000 --to=18001000
  
  # Limit number of results
  khedra monitor activity 0x742d35Cc... --limit=100
  
  # Show activity in specific time range
  khedra monitor activity 0x742d35Cc... --since="2023-01-01" --until="2023-12-31"
  
  # Export activity to file
  khedra monitor activity 0x742d35Cc... --output=/path/to/activity.json
  ```

#### Monitor Configuration

- `khedra monitor config ADDRESS`: Configure monitor settings
  ```bash
  # Update monitor name
  khedra monitor config 0x742d35Cc... --name="Updated Name"
  
  # Configure notifications
  khedra monitor config 0x742d35Cc... --notifications=webhook --webhook-url=https://...
  
  # Set monitoring thresholds
  khedra monitor config 0x742d35Cc... --min-value=1.0 --currency=ETH
  
  # Disable/enable monitor
  khedra monitor config 0x742d35Cc... --enabled=false
  ```

### Chain Management Commands

Commands for managing blockchain connections:

#### Chain Configuration

- `khedra chains list`: List configured blockchain networks
  ```bash
  # List all configured chains
  khedra chains list
  
  # List only enabled chains
  khedra chains list --enabled-only
  
  # Show detailed chain information
  khedra chains list --details
  
  # Export chain configuration
  khedra chains list --format=json --output=chains.json
  ```

- `khedra chains add NAME URL`: Add new blockchain network
  ```bash
  # Add new chain with single RPC
  khedra chains add polygon https://polygon-rpc.com
  
  # Add chain and enable immediately
  khedra chains add arbitrum https://arb1.arbitrum.io/rpc --enable
  
  # Add chain with multiple RPC endpoints
  khedra chains add optimism https://mainnet.optimism.io,https://opt-mainnet.g.alchemy.com/v2/KEY
  
  # Add chain with custom configuration
  khedra chains add base https://mainnet.base.org --chain-id=8453 --symbol=ETH
  ```

- `khedra chains remove NAME`: Remove blockchain network
  ```bash
  # Remove specific chain
  khedra chains remove polygon
  
  # Remove without confirmation
  khedra chains remove arbitrum --force
  
  # Remove and clean up data
  khedra chains remove optimism --cleanup-data
  ```

#### Chain Testing and Validation

- `khedra chains test NAME`: Test connection to blockchain network
  ```bash
  # Test specific chain connectivity
  khedra chains test mainnet
  
  # Test with verbose output
  khedra chains test polygon --verbose
  
  # Test all RPC endpoints for a chain
  khedra chains test mainnet --test-all-rpcs
  
  # Test and benchmark performance
  khedra chains test mainnet --benchmark
  ```

- `khedra chains validate`: Validate all chain configurations
  ```bash
  # Validate all chains
  khedra chains validate
  
  # Validate specific chain
  khedra chains validate --chain=mainnet
  
  # Validate and show detailed results
  khedra chains validate --verbose
  
  # Validate and attempt to fix issues
  khedra chains validate --auto-fix
  ```

### Configuration Commands

Commands for managing Khedra's configuration:

#### Configuration Display and Editing

- `khedra config show`: Display current configuration
  ```bash
  # Show complete configuration
  khedra config show
  
  # Show specific section
  khedra config show --section=services
  khedra config show --section=chains
  khedra config show --section=logging
  
  # Hide sensitive information
  khedra config show --redact
  
  # Export configuration
  khedra config show --format=yaml --output=config.yaml
  khedra config show --format=json --output=config.json
  ```

- `khedra config edit`: Open configuration in editor
  ```bash
  # Edit with default editor
  khedra config edit
  
  # Edit with specific editor
  khedra config edit --editor=vim
  khedra config edit --editor=code
  
  # Edit specific section
  khedra config edit --section=services
  
  # Edit and validate before saving
  khedra config edit --validate
  ```

#### Configuration Management

- `khedra config wizard`: Run interactive configuration wizard
  ```bash
  # Run full configuration wizard
  khedra config wizard
  
  # Run simplified wizard
  khedra config wizard --simple
  
  # Run wizard for specific section
  khedra config wizard --section=services
  khedra config wizard --section=chains
  
  # Run wizard and save to custom location
  khedra config wizard --output=/path/to/config.yaml
  ```

- `khedra config validate`: Validate configuration
  ```bash
  # Validate current configuration
  khedra config validate
  
  # Validate specific file
  khedra config validate --file=/path/to/config.yaml
  
  # Validate and show detailed errors
  khedra config validate --verbose
  
  # Validate against specific schema version
  khedra config validate --schema-version=v1.0
  ```

- `khedra config reset`: Reset configuration to defaults
  ```bash
  # Reset entire configuration (with confirmation)
  khedra config reset
  
  # Reset specific section
  khedra config reset --section=services
  
  # Reset without confirmation prompt
  khedra config reset --force
  
  # Reset and backup current config
  khedra config reset --backup=/path/to/backup.yaml
  ```

### Utility Commands

Commands for various utility operations:

#### Version and Information

- `khedra version`: Display version information
  ```bash
  # Show version
  khedra version
  
  # Show detailed build information
  khedra version --build-info
  
  # Show version in JSON format
  khedra version --json
  
  # Check for updates
  khedra version --check-updates
  ```

#### Initialization and Setup

- `khedra init`: Initialize Khedra configuration and data directories
  ```bash
  # Initialize with wizard
  khedra init
  
  # Initialize with minimal setup
  khedra init --minimal
  
  # Initialize in specific directory
  khedra init --data-dir=/custom/path
  
  # Initialize with custom configuration
  khedra init --config=/path/to/config.yaml
  
  # Reinitialize existing setup
  khedra init --force
  ```

#### Health and Diagnostics

- `khedra health`: Check overall system health
  ```bash
  # Basic health check
  khedra health
  
  # Detailed health report
  khedra health --detailed
  
  # Health check with remediation suggestions
  khedra health --suggest-fixes
  
  # Export health report
  khedra health --output=/path/to/health-report.json
  ```

- `khedra doctor`: Run comprehensive system diagnostics
  ```bash
  # Run all diagnostic checks
  khedra doctor
  
  # Run specific diagnostic categories
  khedra doctor --checks=network,storage,performance
  
  # Run diagnostics and attempt auto-fixes
  khedra doctor --auto-fix
  
  # Run diagnostics for specific chain
  khedra doctor --chain=mainnet
  ```

#### Log and Debug Commands

- `khedra logs`: Display and manage log files
  ```bash
  # Show recent logs
  khedra logs
  
  # Follow logs in real-time
  khedra logs --follow
  
  # Show logs for specific service
  khedra logs --service=scraper
  khedra logs --service=monitor
  
  # Filter logs by level
  khedra logs --level=error
  khedra logs --level=debug
  
  # Show logs in specific time range
  khedra logs --since="2023-01-01" --until="2023-12-31"
  
  # Export logs to file
  khedra logs --output=/path/to/logs.txt
  ```

#### Data Management

- `khedra cleanup`: Clean up temporary files and optimize storage
  ```bash
  # Clean temporary files
  khedra cleanup
  
  # Clean and optimize index
  khedra cleanup --optimize-index
  
  # Clean old log files
  khedra cleanup --logs --older-than=30d
  
  # Dry run to see what would be cleaned
  khedra cleanup --dry-run
  
  # Clean specific data types
  khedra cleanup --cache --temp --logs
  ```
