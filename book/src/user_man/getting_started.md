# Getting Started

## Overview

The `config` package in Khedra manages application configuration through the `config.yaml` file. This file allows you to specify key parameters for running Khedra, including logging, blockchain chains, and services. Additionally, environment variables can override specific configuration options. This document outlines the configuration file structure, validation rules, default values, and environment variable usage.

---

## Quick Start

1. **Copy the example configuration file**:

   ```bash
   cp config.yaml.example config.yaml
   ```

   Modify the `config.yaml` file according to your requirements.

2. **Location of `config.yaml`**:
   - By default, the `config.yaml` file must be in the current directory or in the `~/.khedra` folder. If the file is not present, Khedra creates a default configuration file in `~/.khedra`.

3. **Using Environment Variables**:
   - Environment variables starting with `TB_KHEDRA_` can override specific values in `config.yaml`. For example:

     ```bash
     export TB_KHEDRA_GENERAL_DATADIR="/path/override"
     ```

---

## Configuration File Format

The `config.yaml` file is structured as follows:

```yaml
general:
  data_dir: "~/.khedra/data"     # Path to the data directory (must exist and be writable)
  log_level: "info"              # Log level: debug, info, warn, error

chains:
  - name: "mainnet"              # Blockchain name
    rpcs:                        # List of RPC endpoints (at least one is required)
      - "rpc_endpoint_for_mainnet"
    enabled: true                # Whether this chain is enabled

  - name: "sepolia"
    rpcs:
      - "rpc_endpoint_for_sepolia"
    enabled: true

  - name: "gnosis"
    rpcs:
      - "rpc_endpoint_for_gnosis"
    enabled: false

  - name: "optimism"
    rpcs:
      - "rpc_endpoint_for_optimism"
    enabled: false

services:
  - name: "api"                  # Service name (api, scraper, monitor, ipfs)
    enabled: true                # Whether this service is enabled
    port: 8080                   # Port number for the service

  - name: "scraper"
    enabled: true
    sleep: 60                    # Time (in seconds) between scraping operations
    batch_size: 500              # Number of blocks processed in each batch (50-10,000)

  - name: "monitor"
    enabled: true
    sleep: 60                    # Time (in seconds) between updates
    batch_size: 500              # Number of blocks processed in each batch (50-10,000)

  - name: "ipfs"
    enabled: true
    port: 5001                   # Port number for the service

logging:
  folder: "~/.khedra/logs"       # Path to log directory
  filename: "khedra.log"         # Log file name
  max_size_mb: 10                # Max log file size in MB
  max_backups: 5                 # Number of backup log files to keep
  max_age_days: 30               # Number of days to retain old logs
  compress: true                 # Whether to compress backup logs
```

---

## Using Environment Variables

Khedra allows configuration values to be overridden using environment variables. The environment variable naming convention is:

`TB_KHEDRA_<section>_<key>`

For example:

- To override the `general.data_dir` value:

  ```bash
  export TB_KHEDRA_GENERAL_DATADIR="/path/override"
  ```

- To set the `log_level`:

  ```bash
  export TB_KHEDRA_GENERAL_LOGLEVEL="debug"
  ```

### Precedence Rules

1. Default values are loaded first.
2. Values from `config.yaml` override the defaults.
3. Environment variables (e.g., `TB_KHEDRA_<section>_<key>`) take precedence over both the defaults and the file.

Environment variables must conform to the same validation rules as the configuration file.

---

## Configuration Sections

### General Settings

- **`data_dir`**: Path to the data directory. This directory must exist and be writable.
- **`log_level`**: Logging level. Possible values: `debug`, `info`, `warn`, `error`.

### Chains (Blockchains)

Defines the blockchain networks to interact with. Each chain must have:

- **`name`**: Chain name (e.g., `mainnet`).
- **`rpcs`**: List of RPC endpoints. At least one valid and reachable endpoint is required.
- **`enabled`**: Whether the chain is active.

### Services (API, Scraper, Monitor, IPFS)

Defines various services provided by Khedra. Supported services:

- **API**:
  - Requires `port` to be specified.
- **Scraper** and **Monitor**:
  - **`sleep`**: Duration (seconds) between operations.
  - **`batch_size`**: Number of blocks to process in each operation (50-10,000).
- **IPFS**:
  - Requires `port` to be specified.

### Logging Configuration

Controls the application's logging behavior:

- **`folder`**: Directory for storing logs.
- **`filename`**: Name of the log file.
- **`max_size_mb`**: Maximum log file size before rotation.
- **`max_backups`**: Number of old log files to retain.
- **`max_age_days`**: Retention period for old logs.
- **`compress`**: Whether to compress rotated logs.

---

## Validation Rules

The configuration file and environment variables are validated on load with the following rules:

### General

- `data_dir`: Must be a valid, existing directory and writable.
- `log_level`: Must be one of `debug`, `info`, `warn`, `error`.

### Chains

- `name`: Required and non-empty.
- `rpcs`: Must include at least one valid and reachable RPC URL.
- `enabled`: Defaults to `false` if not specified.

### Services

- `name`: Required and non-empty. Must be one of `api`, `scraper`, `monitor`, `ipfs`.
- `enabled`: Defaults to `false` if not specified.
- `port`: For API and IPFS services, must be between 1024 and 65535.
- `sleep`: Must be non-negative.
- `batch_size`: Must be between 50 and 10,000.

### Logging

- `folder`: Must exist and be writable.
- `filename`: Must end with `.log`.
- `max_size_mb`: Minimum value of 5.
- `max_backups`: Minimum value of 1.
- `max_age_days`: Minimum value of 1.

---

## Default Values

If the configuration file is not found or incomplete, Khedra uses the following defaults:

- **Data directory**: `~/.khedra/data`
- **Log level**: `info`
- **Logging configuration**:
  - Folder: `~/.khedra/logs`
  - Filename: `khedra.log`
  - Max size: 10 MB
  - Max backups: 3
  - Max age: 10 days
  - Compression: Enabled
- **Chains**: Only `mainnet` and `sepolia` enabled by default.
- **Services**: All services (`api`, `scraper`, `monitor`, `ipfs`) enabled with default configurations.

---

## Common Commands

1. **Validate Configuration**:
   Khedra validates the `config.yaml` file and environment variables automatically on startup.

2. **Run Khedra**:

   ```bash
   ./khedra --version
   ```

   Ensure that your `config.yaml` file is properly set up.

3. **Override Configuration with Environment Variables**:

   Use environment variables to override specific configurations:

   ```bash
   export TB_KHEDRA_GENERAL_DATADIR="/new/path"
   ./khedra
   ```

For additional details, see the technical specification.
