# Command-Line Interface

Khedra provides a deliberately small command-line interface (CLI) focused on what actually exists today: initialization, daemon startup, configuration viewing/editing, and pausing / unpausing certain services.

## CLI Architecture

The CLI is built using the `urfave/cli` library. There are no hidden subcommands beyond those listed below, and no status / metrics / restart commands at present.

### Design Principles

1. **Simplicity**: Minimal command set focused on essential operations
2. **Clarity**: Each command has a clear, single purpose
3. **REST API Integration**: Service control via HTTP API for automation
4. **Self-Documenting**: Built-in help for all commands

## Command Overview

Khedra implements these core commands (current implementation):

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

Starts enabled services in a simple order: Control first, then the remaining enabled services in whatever order the configuration map iteration yields (not guaranteed / currently unordered). Services:
- **Scraper** (pausable)
- **Monitor** (pausable, disabled by default; functionality limited)
- **API** (if enabled)
- **IPFS** (if enabled)
- **Control** (always started)

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

**Non-Pausable Services**: `control`, `api`, `ipfs`.

#### `khedra unpause <service>`
Resume paused services.

```bash
# Resume specific services
khedra unpause scraper
khedra unpause monitor  

# Resume all paused services
khedra unpause all
```

Same service support as pause command. A service must be paused to unpause it. Only `scraper` and `monitor` are recognized plus the alias `all`.

### Control Service API

Pause/unpause operations are available via a minimal HTTP interface on the Control Service (first available of ports 8338, 8337, 8336, 8335). Mutating operations use HTTP GET.

#### Status Queries
```bash
# Check all service status
curl "http://localhost:8338/isPaused"

# Check specific service
curl "http://localhost:8338/isPaused?name=scraper"
curl "http://localhost:8338/isPaused?name=monitor"
```

#### Pause Operations (implemented as HTTP GET)
```bash
# Pause specific service
curl "http://localhost:8338/pause?name=scraper"
curl "http://localhost:8338/pause?name=monitor"

# Pause all pausable services
curl "http://localhost:8338/pause?name=all"
curl "http://localhost:8338/pause"    # alternative
```

#### Unpause Operations (implemented as HTTP GET)
```bash
# Unpause specific service
curl "http://localhost:8338/unpause?name=scraper"
curl "http://localhost:8338/unpause?name=monitor"

# Unpause all services
curl "http://localhost:8338/unpause?name=all"
curl "http://localhost:8338/unpause"   # alternative
```

#### API Responses

Status queries return simple JSON arrays like:
```json
[
  {"name": "scraper", "status": "running"},
  {"name": "monitor", "status": "paused"}, 
  {"name": "control", "status": "not pausable"},
  {"name": "ipfs", "status": "not pausable"}
]
```

Control operations return result arrays. Example:
```json
[
  {"name": "scraper", "status": "paused"}
]
```

#### Error Handling

Invalid service names return an error JSON body with 400.
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

## Environment Variables (current)

- `TB_KHEDRA_WAIT_FOR_NODE` (optional): process name to block on before starting
- `TB_KHEDRA_WAIT_SECONDS` (default 30 if waiting): post-detect delay
- `TB_KHEDRA_LOGGING_LEVEL`: one of `debug|info|warn|error`
- `EDITOR`: used by `khedra config edit`

## Error Handling

### Common Issues

**Service not found**: Ensure service name is correct (`scraper`, `monitor`, or `all`)

**Control service unavailable**: Verify daemon is running and control service is enabled

**Permission denied**: Ensure proper file permissions for configuration and data directories

**Port conflicts**: Control service scans 8338 → 8335 and uses the first open port

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
