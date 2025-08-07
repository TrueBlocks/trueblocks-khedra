# Using Khedra

This chapter covers the practical aspects of working with Khedra for blockchain data indexing and monitoring.

## Command Overview

Khedra provides five essential commands:

```text
NAME:
   khedra - A tool to index, monitor, serve, and share blockchain data

USAGE:
   khedra [global options] command [command options]

COMMANDS:
   init     Initializes Khedra configuration
   daemon   Runs Khedra's services  
   config   Manages Khedra configuration
   pause    Pause services (scraper, monitor, all)
   unpause  Unpause services (scraper, monitor, all)
   help, h  Shows help for commands

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

## Getting Started

### 1. Initialize Configuration

Set up Khedra's configuration interactively:

```bash
khedra init
```

The initialization wizard configures:
- **General Settings**: Data directories, logging preferences
- **Chain Configuration**: RPC endpoints for blockchain networks
- **Service Settings**: Which services to enable (scraper, monitor, API, IPFS)
- **Port Configuration**: Network ports for HTTP services

### 2. Start Services

Launch all configured services:

```bash
khedra daemon
```

This starts:
- **Scraper**: Blockchain indexing service
- **Monitor**: Address monitoring service
- **API**: REST endpoints (if enabled)
- **Control**: Service management interface
- **IPFS**: Distributed sharing (if enabled)

The daemon runs until interrupted (Ctrl+C) or receives SIGTERM.

### 3. Manage Configuration

View or edit configuration:

```bash
# Display current configuration
khedra config show

# Edit configuration in default editor
khedra config edit
```

Changes require restarting the daemon to take effect.

## Service Management

Control individual services at runtime without stopping the daemon:

### Pause Services

Temporarily stop service operations:

```bash
# Pause specific services
khedra pause scraper    # Stop blockchain indexing
khedra pause monitor    # Stop address monitoring

# Pause all pausable services
khedra pause all
```

### Resume Services

Restart paused services:

```bash
# Resume specific services
khedra unpause scraper
khedra unpause monitor

# Resume all paused services  
khedra unpause all
```

### Service Types

**Pausable Services**: 
- `scraper`: Can be paused to stop indexing
- `monitor`: Can be paused to stop address monitoring

**Always-On Services**:
- `control`: Provides service management API
- `api`: Serves data queries (cannot be paused)
- `ipfs`: Handles distributed sharing (cannot be paused)

## REST API Control

The Control service (port 8338) provides HTTP endpoints for automation:

### Check Service Status

```bash
# All service status
curl "http://localhost:8338/isPaused"

# Specific service status
curl "http://localhost:8338/isPaused?name=scraper"
```

Response format:
```json
[
  {"name": "scraper", "status": "running"},
  {"name": "monitor", "status": "paused"},
  {"name": "control", "status": "not pausable"}
]
```

### Control Operations

```bash
# Pause services
curl -X POST "http://localhost:8338/pause?name=scraper"
curl -X POST "http://localhost:8338/pause?name=all"

# Resume services
curl -X POST "http://localhost:8338/unpause?name=scraper"
curl -X POST "http://localhost:8338/unpause?name=all"
```

## Common Workflows

### Initial Setup

1. **Install**: Build or install Khedra binary
2. **Initialize**: Run `khedra init` to configure
3. **Start**: Run `khedra daemon` to begin indexing
4. **Monitor**: Use pause/unpause for operational control

### Operational Management

```bash
# Check what's running
curl "http://localhost:8338/isPaused"

# Pause indexing during maintenance
khedra pause scraper

# Resume normal operations
khedra unpause scraper

# Pause everything for system maintenance
khedra pause all
khedra unpause all
```

### Configuration Updates

```bash
# View current settings
khedra config show

# Edit configuration
khedra config edit

# Restart to apply changes
# (Stop daemon with Ctrl+C, then restart)
khedra daemon
```

## Environment Variables

Control behavior with environment variables:

- `TB_KHEDRA_WAIT_FOR_NODE`: Wait for specific node process (e.g., `erigon`, `geth`)
- `TB_KHEDRA_WAIT_SECONDS`: Seconds to wait after node detection (default: 30)
- `TB_KHEDRA_LOGGING_LEVEL`: Log verbosity (`debug`, `info`, `warn`, `error`)
- `EDITOR`: Editor for `config edit` command

Example:
```bash
TB_KHEDRA_LOGGING_LEVEL=debug khedra daemon
```

## Troubleshooting

### Common Issues

**Configuration not found**: Run `khedra init` to create initial configuration

**Port conflicts**: Control service automatically finds available ports (8338, 8337, 8336, 8335)

**Service not pausable**: Only `scraper` and `monitor` services can be paused

**Control API unavailable**: Ensure daemon is running and control service is enabled

### Getting Help

```bash
# Command-specific help
khedra init --help
khedra daemon --help
khedra pause --help

# General help
khedra --help

# Version information
khedra --version
```

### Debug Information

Enable verbose logging:
```bash
TB_KHEDRA_LOGGING_LEVEL=debug khedra daemon
```

Check service status via API:
```bash
curl "http://localhost:8338/isPaused" | jq
```

Monitor log output for service-specific issues and configuration problems.
