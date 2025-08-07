# Command-Line Interface

Khedra provides a simple, focused command-line interface (CLI) for managing the system. The CLI is designed around essential operations: initialization, daemon management, configuration, and service control.

## CLI Architecture

The CLI is built using the `urfave/cli` library, providing a lightweight and user-friendly interface focused on core functionality.

### Design Principles

1. **Simplicity**: Minimal command set focused on essential operations
2. **Clarity**: Each command has a clear, single purpose
3. **REST API Integration**: Service control via HTTP API for automation
4. **Self-Documenting**: Built-in help for all commands

## Command Overview

Khedra implements five core commands:

### Essential Commands

#### `khedra init`
Initialize Khedra configuration interactively.

```bash
khedra init
```

Launches an interactive wizard that configures:
- General settings (data directories, logging)
- Chain configurations (RPC endpoints, indexing preferences)
- Service settings (which services to enable)
- Service ports (API, control, IPFS services)

#### `khedra daemon`
Start Khedra daemon with all configured services.

```bash
khedra daemon
```

Starts all enabled services:
- **Scraper**: Blockchain indexing service
- **Monitor**: Address monitoring service  
- **API**: REST API service (if enabled)
- **Control**: Service management HTTP interface
- **IPFS**: Distributed data sharing (if enabled)

The daemon runs until interrupted (Ctrl+C) or receives a termination signal.

#### `khedra config`
Manage Khedra configuration.

```bash
# Display current configuration
khedra config show

# Edit configuration in default editor
khedra config edit
```

Configuration management:
- `show`: Display current configuration in readable format
- `edit`: Open configuration file in system editor (respects `$EDITOR` environment variable)

#### `khedra pause <service>`
Pause running services.

```bash
# Pause specific services
khedra pause scraper
khedra pause monitor

# Pause all pausable services
khedra pause all
```

**Supported Services**:
- `scraper`: Blockchain indexing service
- `monitor`: Address monitoring service
- `all`: All pausable services

**Non-Pausable Services**: `control`, `api`, `ipfs` (these provide critical system functionality)

#### `khedra unpause <service>`
Resume paused services.

```bash
# Resume specific services
khedra unpause scraper
khedra unpause monitor  

# Resume all paused services
khedra unpause all
```

Same service support as pause command. Services must be paused to be unpaused.

### Control Service API

All pause/unpause operations are also available via REST API on the Control Service (default port 8338):

#### Status Queries
```bash
# Check all service status
curl "http://localhost:8338/isPaused"

# Check specific service
curl "http://localhost:8338/isPaused?name=scraper"
curl "http://localhost:8338/isPaused?name=monitor"
```

#### Pause Operations
```bash
# Pause specific service
curl -X POST "http://localhost:8338/pause?name=scraper"
curl -X POST "http://localhost:8338/pause?name=monitor"

# Pause all pausable services
curl -X POST "http://localhost:8338/pause?name=all"
curl -X POST "http://localhost:8338/pause"    # alternative
```

#### Unpause Operations
```bash
# Unpause specific service
curl -X POST "http://localhost:8338/unpause?name=scraper" 
curl -X POST "http://localhost:8338/unpause?name=monitor"

# Unpause all services
curl -X POST "http://localhost:8338/unpause?name=all"
curl -X POST "http://localhost:8338/unpause"   # alternative
```

#### API Responses

Status queries return JSON arrays:
```json
[
  {"name": "scraper", "status": "running"},
  {"name": "monitor", "status": "paused"}, 
  {"name": "control", "status": "not pausable"},
  {"name": "ipfs", "status": "not pausable"}
]
```

Control operations return operation results:
```json
[
  {"name": "scraper", "status": "paused"}
]
```

#### Error Handling

Invalid service names return HTTP 400:
```json
{"error": "service 'invalid' not found or is not pausable"}
```

## Usage Examples

### Complete Startup Workflow

```bash
# 1. Initialize configuration (first time only)
khedra init

# 2. Start daemon
khedra daemon
```

### Service Management During Operation

```bash
# Check what's running
curl "http://localhost:8338/isPaused"

# Pause indexing temporarily  
khedra pause scraper

# Resume when ready
khedra unpause scraper

# Pause everything for maintenance
khedra pause all

# Resume normal operations
khedra unpause all
```

### Configuration Management

```bash
# View current settings
khedra config show

# Modify configuration
khedra config edit

# Restart daemon to apply changes
# (stop with Ctrl+C, then restart)
khedra daemon
```

## Environment Variables

Khedra respects these environment variables:

- `TB_KHEDRA_WAIT_FOR_NODE`: Node process name to wait for before starting (e.g., `erigon`, `geth`)
- `TB_KHEDRA_WAIT_SECONDS`: Seconds to wait after node detection (default: 30)
- `TB_KHEDRA_LOGGING_LEVEL`: Log level (`debug`, `info`, `warn`, `error`)
- `EDITOR`: Editor for `config edit` command

## Error Handling

### Common Issues

**Service not found**: Ensure service name is correct (`scraper`, `monitor`, or `all`)

**Control service unavailable**: Verify daemon is running and control service is enabled

**Permission denied**: Ensure proper file permissions for configuration and data directories

**Port conflicts**: Control service automatically finds available ports (8338, 8337, 8336, 8335)

### Debugging

Enable debug logging:
```bash
TB_KHEDRA_LOGGING_LEVEL=debug khedra daemon
```

Check service status:
```bash
curl "http://localhost:8338/isPaused"
```

View configuration:
```bash
khedra config show
```
