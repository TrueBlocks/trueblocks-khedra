# Getting Started

## Quick Start

Get Khedra running in 3 simple steps:

### 1. Initialize Khedra

Run the configuration wizard to set up your blockchain connections and services:

```bash
khedra init
```

This interactive wizard will guide you through:
- Setting up blockchain RPC connections (Ethereum, Polygon, etc.)
- Configuring which services to enable (scraper, monitor, API, IPFS)
- Setting up logging and data storage paths

### 2. Start Khedra

Start all configured services:

```bash
khedra daemon
```

This starts the daemon with all enabled services. The Control Service runs automatically and manages other services.

### 3. Control Services via API

Once running, manage services through the REST API:

```bash
# Check service status
curl http://localhost:8080/api/v1/services

# Start/stop individual services
curl -X POST http://localhost:8080/api/v1/services/scraper/start
curl -X POST http://localhost:8080/api/v1/services/monitor/stop

# Get system status
curl http://localhost:8080/api/v1/status
```

That's it! Your Khedra instance is now indexing blockchain data and ready for queries.

---

## Detailed Configuration

**Khedra** runs primarily from a configuration file called `config.yaml`. This file lives at `~/.khedra/config.yaml` by default. If the file is not found, **Khedra** creates a default configuration in this location.

The config file allows you to specify key parameters for running **khedra**, including which chains to index/monitor, which services to enable, how detailed to log the processes, and where and how to publish (that is, share) the results.

You may use environment variables to override specific options. This document outlines the configuration file structure, validation rules, default values, and environment variable usage.

## Installation

1. **Download, build, and test khedra**:

   ```bash
   git clone https://github.com/TrueBlocks/trueblocks-khedra.git
   cd trueblocks-khedra
   go build -o khedra main.go
   ./khedra version
   ```

   You should get something similar to `khedra v4.0.0-release`.

2. **You may edit the config file with**:

   ```bash
   ./khedra config edit
   ```

   Modify the file according to your requirements (see below).

   The minimal configuration needed is to provide a valid RPC to Ethereum mainnet. (All configurations require access to Ethereum `mainnet`.)

   You may configure as many other EVM-compatible chains (each with its own RPC) as you like.

3. **Use the Wizard**:

    You may also use the **khedra** wizard to create a configuration file. The wizard will prompt you for the required information and generate a `config.yaml` file.
  
    ```bash
    ./khedra init
    ```

4. **Location of the configuration file**:

    By default, the config file resides at `~/.khedra/config.yaml`. (The folder and the file will be created if it does not exist).

    You may, however, place a `config.yaml` file in the current working folder (the folder from which you run **khedra**). If found locally, this configuration file will dominate. This allows for running multiple instances of the software concurrently.

---

## Advanced Configuration Examples

### Production Deployment Configuration

For production environments with high availability and performance requirements:

```yaml
general:
  indexPath: "/var/lib/khedra/index"     # Fast SSD storage
  cachePath: "/var/lib/khedra/cache"     # Local SSD cache
  dataDir: "/var/lib/khedra"             # Dedicated data directory

chains:
  mainnet:
    rpcs:
      - "https://eth-mainnet.alchemyapi.io/v2/YOUR_PREMIUM_KEY"
      - "https://mainnet.infura.io/v3/YOUR_PROJECT_ID"
      - "https://rpc.ankr.com/eth"              # Fallback
      - "https://ethereum.publicnode.com"       # Additional fallback
    enabled: true

  polygon:
    rpcs:
      - "https://polygon-mainnet.g.alchemy.com/v2/YOUR_KEY"
      - "https://polygon-rpc.com"
    enabled: true

  arbitrum:
    rpcs:
      - "https://arb-mainnet.g.alchemy.com/v2/YOUR_KEY"
      - "https://arb1.arbitrum.io/rpc"
    enabled: true

services:
  scraper:
    enabled: true
    sleep: 5                     # Aggressive indexing
    batchSize: 2000             # Large batches for efficiency

  monitor:
    enabled: true
    sleep: 5                    # Fast monitoring
    batchSize: 500

  api:
    enabled: true
    port: 8080

  ipfs:
    enabled: true
    port: 8083

logging:
  folder: "/var/log/khedra"     # System log directory
  filename: "khedra.log"
  toFile: true                  # Always log to file in production
  level: "info"                 # Balanced logging
  maxSize: 100                  # Larger log files
  maxBackups: 10               # More backup files
  maxAge: 90                   # Longer retention
  compress: true               # Compress old logs
```

### Multi-Chain Development Environment

For developers working with multiple blockchain networks:

```yaml
general:
  indexPath: "~/.khedra/dev/index"
  cachePath: "~/.khedra/dev/cache"

chains:
  mainnet:
    rpcs:
      - "https://eth-mainnet.alchemyapi.io/v2/YOUR_DEV_KEY"
    enabled: true

  sepolia:
    rpcs:
      - "https://eth-sepolia.g.alchemy.com/v2/YOUR_KEY"
      - "https://sepolia.infura.io/v3/YOUR_PROJECT_ID"
    enabled: true

  polygon:
    rpcs:
      - "https://polygon-mumbai.g.alchemy.com/v2/YOUR_KEY"
    enabled: true

  optimism:
    rpcs:
      - "https://opt-goerli.g.alchemy.com/v2/YOUR_KEY"
    enabled: true

  arbitrum:
    rpcs:
      - "https://arb-goerli.g.alchemy.com/v2/YOUR_KEY"
    enabled: true

  base:
    rpcs:
      - "https://base-goerli.g.alchemy.com/v2/YOUR_KEY"
    enabled: true

services:
  scraper:
    enabled: true
    sleep: 15                   # Moderate speed for development
    batchSize: 500

  monitor:
    enabled: true               # Enable for testing monitoring features
    sleep: 30
    batchSize: 100

  api:
    enabled: true
    port: 8080

  ipfs:
    enabled: false              # Disable to reduce resource usage

logging:
  folder: "~/.khedra/dev/logs"
  filename: "khedra-dev.log"
  toFile: true
  level: "debug"               # Verbose logging for development
  maxSize: 10
  maxBackups: 5
  maxAge: 7                    # Shorter retention for dev
  compress: false              # No compression for easier reading
```

### High-Availability Load-Balanced Setup

Configuration for running multiple Khedra instances behind a load balancer:

```yaml
# Instance 1: Primary indexing node
general:
  indexPath: "/shared/khedra/index"      # Shared storage
  cachePath: "/local/khedra/cache1"      # Local cache per instance

chains:
  mainnet:
    rpcs:
      - "https://eth-mainnet-primary.alchemyapi.io/v2/KEY1"
      - "https://eth-mainnet-backup.infura.io/v3/PROJECT1"
    enabled: true

services:
  scraper:
    enabled: true              # Primary indexer
    sleep: 5
    batchSize: 2000

  monitor:
    enabled: false             # Disabled on indexing nodes

  api:
    enabled: false             # Dedicated API nodes

  ipfs:
    enabled: true              # IPFS on indexing nodes
    port: 8083

logging:
  folder: "/var/log/khedra"
  filename: "khedra-indexer-1.log"
  toFile: true
  level: "info"
---

# Instance 2: API-only node
general:
  indexPath: "/shared/khedra/index"      # Same shared storage
  cachePath: "/local/khedra/cache2"      # Different local cache

chains:
  mainnet:
    rpcs:
      - "https://eth-mainnet-api.alchemyapi.io/v2/KEY2"
    enabled: true

services:
  scraper:
    enabled: false             # No indexing on API nodes

  monitor:
    enabled: true              # Monitoring on API nodes
    sleep: 10
    batchSize: 200

  api:
    enabled: true              # Primary function
    port: 8080

  ipfs:
    enabled: false             # Not needed on API nodes

logging:
  folder: "/var/log/khedra"
  filename: "khedra-api-2.log"
  toFile: true
  level: "warn"               # Less verbose for API nodes
```

### Resource-Constrained Environment

Configuration for running Khedra on limited hardware (VPS, Raspberry Pi, etc.):

```yaml
general:
  indexPath: "~/.khedra/index"
  cachePath: "~/.khedra/cache"

chains:
  mainnet:
    rpcs:
      - "https://ethereum.publicnode.com"    # Free RPC
      - "https://rpc.ankr.com/eth"           # Backup free RPC
    enabled: true

  # Only enable additional chains if needed
  sepolia:
    rpcs:
      - "https://ethereum-sepolia.publicnode.com"
    enabled: false               # Disabled to save resources

services:
  scraper:
    enabled: true
    sleep: 60                   # Very conservative indexing
    batchSize: 50               # Small batches

  monitor:
    enabled: false              # Disable to save resources
    sleep: 300
    batchSize: 10

  api:
    enabled: true
    port: 8080

  ipfs:
    enabled: false              # Disable to save bandwidth/storage

logging:
  folder: "~/.khedra/logs"
  filename: "khedra.log"
  toFile: false                 # Console only to save disk space
  level: "warn"                # Minimal logging
  maxSize: 5                   # Small log files
  maxBackups: 2                # Minimal retention
  maxAge: 7
  compress: true
```

### Security-Focused Configuration

Configuration with enhanced security for sensitive environments:

```yaml
general:
  indexPath: "/encrypted/khedra/index"   # Encrypted storage
  cachePath: "/encrypted/khedra/cache"

chains:
  mainnet:
    rpcs:
      - "https://your-private-node.internal:8545"  # Private RPC node
    enabled: true

services:
  scraper:
    enabled: true
    sleep: 10
    batchSize: 1000

  monitor:
    enabled: true
    sleep: 15
    batchSize: 100

  api:
    enabled: true
    port: 8080                  # Consider using reverse proxy with TLS

  ipfs:
    enabled: false              # Disable external data sharing

logging:
  folder: "/secure/logs/khedra"
  filename: "khedra.log"
  toFile: true
  level: "info"
  maxSize: 50
  maxBackups: 20              # Extended retention for audit
  maxAge: 365                 # Long retention for compliance
  compress: true

# Environment variables for sensitive data:
# TB_KHEDRA_CHAINS_MAINNET_RPCS="https://user:pass@private-node:8545"
# TB_KHEDRA_API_AUTH_TOKEN="your-secure-api-token"
# TB_KHEDRA_WAIT_FOR_NODE="erigon"          # (Optional) Wait for node process before starting
# TB_KHEDRA_WAIT_SECONDS="60"              # (Optional) Wait time for node stabilization (default: 30)
```

### Testing and CI/CD Configuration

Configuration optimized for automated testing environments:

```yaml
general:
  indexPath: "./test-data/index"
  cachePath: "./test-data/cache"

chains:
  sepolia:                      # Use testnet for testing
    rpcs:
      - "https://ethereum-sepolia.publicnode.com"
    enabled: true

  mainnet:
    rpcs:
      - "https://ethereum.publicnode.com"
    enabled: false              # Disabled for testing

services:
  scraper:
    enabled: true
    sleep: 30                   # Conservative for CI resources
    batchSize: 100

  monitor:
    enabled: true               # Test monitoring functionality
    sleep: 60
    batchSize: 50

  api:
    enabled: true
    port: 8080

  ipfs:
    enabled: false              # Not needed for testing

logging:
  folder: "./test-logs"
  filename: "khedra-test.log"
  toFile: true
  level: "debug"               # Verbose for troubleshooting tests
  maxSize: 10
  maxBackups: 3
  maxAge: 1                    # Short retention for CI
  compress: false              # Easier to read in CI logs
```
