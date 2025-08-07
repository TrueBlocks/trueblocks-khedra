# TrueBlocks Khedra

## Overview

Khedra is TrueBlocks' service management system that provides specialized blockchain data indexing, monitoring, and serving capabilities. It operates as a collection of microservices managed through a unified CLI and REST API interface.

## Key Features

- **Blockchain Indexing**: High-performance scraping and indexing of Ethereum transaction data
- **Address Monitoring**: Real-time monitoring of specific Ethereum addresses
- **REST API**: HTTP interface for querying indexed data
- **Service Management**: Runtime control over individual services (pause/unpause)
- **Multi-Chain Support**: Configurable support for various Ethereum-compatible networks
- **IPFS Integration**: Distributed data sharing capabilities

## Quick Start

### Prerequisites

- [TrueBlocks Core](https://github.com/TrueBlocks/trueblocks-core) installed
- Go 1.23+ for building from source (Go Version)
- Running Ethereum node (Erigon, Geth, etc.) or RPC endpoint access

### Installation

```bash
# Clone and build
git clone https://github.com/TrueBlocks/trueblocks-core.git
cd trueblocks-core
make

# Initialize configuration
./bin/khedra init

# Start services
./bin/khedra daemon
```

## Usage

### Basic Commands

```bash
# Interactive configuration setup
khedra init

# Start all enabled services
khedra daemon

# View current configuration
khedra config show

# Edit configuration
khedra config edit
```

### Service Management

Control individual services at runtime:

```bash
# Pause services
khedra pause scraper    # Pause blockchain indexing
khedra pause monitor    # Pause address monitoring
khedra pause all        # Pause all pausable services

# Resume services
khedra unpause scraper
khedra unpause monitor
khedra unpause all
```

**Pausable Services**: `scraper`, `monitor`  
**Always-On Services**: `control`, `api`, `ipfs`

### REST API Control

Service management via HTTP (Control Service on port 8338 or an available value):

```bash
# Check service status
curl "http://localhost:8338/isPaused"

# Pause/unpause services
curl -X POST "http://localhost:8338/pause?name=scraper"
curl -X POST "http://localhost:8338/unpause?name=scraper"
curl -X POST "http://localhost:8338/pause?name=all"
```

## Configuration

Khedra uses YAML configuration managed through the `init` wizard or direct editing.

### Environment Variables

- `TB_KHEDRA_WAIT_FOR_NODE`: Node process to wait for before starting (e.g., `erigon`, `geth`)
- `TB_KHEDRA_WAIT_SECONDS`: Seconds to wait after node detection -- allows node to initialize (default: 30)
- `TB_KHEDRA_LOGGING_LEVEL`: Log level (`debug`, `info`, `warn`, `error`)

### Configuration Sections

- **General**: Data directories, logging preferences
- **Chains**: RPC endpoints, indexing settings per blockchain
- **Services**: Which services to enable and their specific settings
- **Ports**: Network ports for API, control, and IPFS services

## Architecture

Khedra consists of five core services:

1. **Scraper**: Indexes blockchain transactions and builds the Unchained Index
2. **Monitor**: Tracks specific addresses and detects relevant transactions
3. **API**: Provides REST endpoints for querying indexed data
4. **Control**: HTTP interface for service management (pause/unpause)
5. **IPFS**: Distributed data sharing and chunk distribution

Services communicate through shared data structures and can be independently controlled.

## Configuration

Before using `khedra`, you may need to configure it to point at the TrueBlocks indexing data or specify custom indexing rules:

- **Config File**: By default, `khedra` may look for a configuration file at `~/.trueblocks/trueblocks-khedra.conf`.  
- **Environment Variables**:  
  - `KHEDRA_DATA_DIR`: Path to where you want `khedra` to store or read data.  
  - `KHEDRA_LOG_LEVEL`: Adjusts the verbosity of logs (`DEBUG`, `INFO`, `WARN`, `ERROR`).
  - `TB_KHEDRA_WAIT_FOR_NODE`: (Optional) Name of the node process to wait for before starting (e.g., `erigon`, `geth`). If not set, Khedra starts immediately.
  - `TB_KHEDRA_WAIT_SECONDS`: (Optional) Number of seconds to wait for node stabilization after detection (default: 30). Only used when `TB_KHEDRA_WAIT_FOR_NODE` is set.

Refer to the sample configuration file (`.conf.example`) in this repo for a template of possible settings.

## Usage

Once khedra is built and configured, you may use these commands:

### Basic Commands

```bash
# Initialize khedra configuration
khedra init

# Start the daemon with all enabled services
khedra daemon

# View or edit configuration
khedra config show
khedra config edit
```

### Service Management

Khedra provides runtime control over individual services:

```bash
# Pause a specific service (scraper or monitor)
khedra pause scraper
khedra pause monitor

# Pause all pausable services
khedra pause all

# Unpause services
khedra unpause scraper
khedra unpause monitor
khedra unpause all
```

**Note**: Only the `scraper` and `monitor` services can be paused/unpaused. The `control`, `api`, and `ipfs` services cannot be paused. The service must be enabled in configuration and currently running to be paused.

## Docker Support

As of the latest version, we've removed docker support for this tool.

## Documentation

Complete documentation is available in the [Khedra Book](./book/), including:

- **User Manual**: Installation, configuration, and usage guides
- **Technical Specification**: Architecture, APIs, and implementation details
- **Command Reference**: Complete CLI and REST API documentation

## Contributing

Khedra is part of the [TrueBlocks](https://github.com/TrueBlocks/trueblocks-core) ecosystem. Please refer to the main TrueBlocks repository for:

- Contributing guidelines
- Development setup
- Issue reporting
- Community guidelines

## License

This project uses the same license as TrueBlocks Core. See the [LICENSE](LICENSE) file for details.

---

*TrueBlocks Khedra - Blockchain indexing and monitoring for the decentralized web*
