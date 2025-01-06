# Getting Started

## Overview

**Khedra** runs primarily from a configuration file called `config.yaml`. This file lives at `~/.khedra/config.yaml` by default. If the file is not found, **Khedra** creates a default configuration in this location.

The config file allows you to specify key parameters for running **khedra**, including which chains to index/monitor, which services to enable, how detailed to log the processes, and where and how to publish (that is, share) the results.

You may use environment variables to override specific options. This document outlines the configuration file structure, validation rules, default values, and environment variable usage.

---

## Quick Start

1. **Download, build, and test khedra**:

   ```bash
   git clone https://github.com/TrueBlocks/trueblocks-khedra.git
   cd trueblocks-khedra
   go build -o khedra main.go
   ./khedra version
   ```

   You should get something similar to `khedra v4.0.0-release`.

2. **Establish the config file and edit values for your system**:

   ```bash
   mkdir -p ~/.khedra
   cp config.yaml.example ~/.khedra/config.yaml
   ./khedra config edit
   ```

   Modify the file according to your requirements (see below).

   The minimal configuration needed is to provide a valid RPC to Ethereum mainnet. (All configurations require access to Ethereum `mainnet`.)

   You may configure as many other EVM-compatible chains (each with its own RPC) as you like.

3. **Location of the configuration file**:

    By default, the config file resides at `~/.khedra/config.yaml`. (The folder and the file will be created if it does not exist).

    You may, however, place a `config.yaml` file in the current working folder (the folder from which you run **khedra**). If found locally, this configuration file will dominate. This allows for running multiple instances of the software concurrently.

    If no `config.yaml` file is found, **khedra** creates a default configuration in its default location.

4. **Using Environment Variables**:

   You may override configuration options using environment variables, each of which must take the form `TB_KHEDRA_<section>_<key>`.

   For example, the following overrides the `general.data_dir` value.

     ```bash
     export TB_KHEDRA_GENERAL_DATADIR="/path/override"
     ```

    You'll notice that underbars (`_`) in the `<key>` names are not needed.

---

## Configuration File Format

The `config.yaml` file (shown here with default values) is structured as follows:

```yaml
# Khedra Configuration File
# Version: 2.0

general:
  data_dir: "~/.khedra/data"     # See note 1

chains:
  mainnet:                       # Blockchain name (see notes 2, 3, and 4)
    rpcs:                        # A list of RPC endpoints (at least one is required)
      - "rpc_endpoint_for_mainnet"
    enabled: true                # `true` if this chain is enabled
  sepolia:
    rpcs:
      - "rpc_endpoint_for_sepolia"
    enabled: true
  gnosis:                         # Add as many chains as your machine can handle
    rpcs:
      - "rpc_endpoint_for_gnosis" # must be a reachable, valid URL if the chain is enabled
    enabled: false                # this chain is disabled
  optimism:
    rpcs:
      - "rpc_endpoint_for_optimism"
    enabled: false

services:                          # See note 5
  scraper:               # Required. (One of: api, scraper, monitor, ipfs, control)
    enabled: true                  # `true` if the service is enabled
    sleep: 12                      # Seconds between scraping batches (see note 6)
    batch_size: 500                # Number of blocks to process in a batch (range: 50-10000)

  monitor:
    enabled: true
    sleep: 12                      # Seconds between scraping batches (see note 6)
    batch_size: 500                # Number of blocks processed in a batch (range: 50-10000)

  api:
    enabled: true
    port: 8080                     # Port number for API service (the port must be available)

  ipfs:
    enabled: true
    port: 5001                     # Port number for IPFS service (the port must be available)

  control:
    enabled: true                  # Always enabled - false values are invalid
    port: 5001                     # Port number for IPFS service (the port must be available)

logging:
  folder: "~/.khedra/logs"         # Path to log directory (must exist and be writable)
  filename: "khedra.log"           # Log file name (must end with .log)
  log_level: "info"                # One of: debug, info, warn, error
  max_size_mb: 10                  # Max log file size in MB
  max_backups: 5                   # Number of backup log files to keep
  max_age_days: 30                 # Number of days to retain old logs
  compress: true                   # Whether to compress backup logs
```

**Notes:**

1. The `data_dir` value must be a valid, existing directory that is writable. You may wish to change this value to a location with suitable disc scape. Depending on configuration, the Unchained Index and binary caches may approach 200GB.

2. The `chains` section is required. At least one chain must be enabled.

3. If chains other than Ethereum `mainnet` are configured, you must also configure Ethereum `mainnet`. The software reads `mainnet` smart contracts (such as the *Unchained Index* and *UniSwap*) during normal operation.

4. We've used [this repository](https://github.com/ethereum-lists/chains) to identify chain names. Using consistent chain names aides in sharing indexes. Use these values in your configuration if you wish to fully participate in sharing the *Unchained Index*.

5. The `services` section is required. At least one service must be enabled.

6. When a `scraper` or `monitor` is "catching up" to a chain, the `sleep` value is ignored.

---

## Using Environment Variables

**Khedra** allows configuration values to be overridden at runtime using environment variables. The value of an environment variable takes precedence over the defaults and the configuration file.

The environment variable naming convention is:

`TB_KHEDRA_<section>_<key>`

For example:

- To override the `general.data_dir` value:

  ```bash
  export TB_KHEDRA_GENERAL_DATADIR="/path/override"
  ```

- To override `logging.log_level`:

  ```bash
  export TB_KHEDRA_LOGGING_LOGLEVEL="debug"
  ```

- To override `services[0].batch_size`:

  ```bash
  export TB_KHEDRA_GENERAL_LOGLEVEL="debug"
  ```

Underbars (`_`) in `<key>` names are not used and should be omitted.

### Overriding Chains and Services

Environment variables can also be used to override values for chains and services settings. The naming convention for these sections is as follows:

`TB_KHEDRA_<section>_<name>_<key>`

Where:

- `<section>` is either `CHAINS` or `SERVICES`.
- `<name>` is the name of the chain or service (converted to uppercase).
- `<key>` is the specific field to override.

#### Examples

To override the RPC endpoints for the `mainnet` chain:

```bash
export TB_KHEDRA_CHAINS_MAINNET_RPCS="http://rpc1.mainnet,http://rpc2.mainnet"
```

You may list mulitple RPC endpoints by separating them with commas.

To disable the `mainnet` chain:

```bash
export TB_KHEDRA_CHAINS_MAINNET_ENABLED="false"
```

To enable the `api` service:

```bash
export TB_KHEDRA_SERVICES_API_ENABLED="true"
```

To set the port for the `api` service:

```bash
export TB_KHEDRA_SERVICES_API_PORT="8088"
```

### Precedence Rules

1. Default values are loaded first,
2. Values from `config.yaml` override the defaults,
3. Environment variables take precedence over both the defaults and the file.

The values set by environment variables must conform to the same validation rules as the configuration file.

---

## Configuration Sections

### General Settings

- **`data_dir`**: The location where **khedra** stores all of its data. This directory must exist and be writable.

### Chains (Blockchains)

Defines the blockchain networks to interact with. Each chain must have:

- **`name`**: Chain name (e.g., `mainnet`).
- **`rpcs`**: List of RPC endpoints. At least one valid and reachable endpoint is required.
- **`enabled`**: Whether the chain is active.

#### Behavior for Empty RPCs
- If the `RPCs` field is empty in the environment, it is ignored and the configuration file's value is preserved.
- If the `RPCs` field is empty in the final configuration (after merging), the configuration will be rejected.

---

### Services (API, Scraper, Monitor, IPFS)

Defines various services provided by **Khedra**. Supported services:

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
- **`log_level`**: Logging level. Possible values: `debug`, `info`, `warn`, `error`.
- **`max_size_mb`**: Maximum log file size before rotation.
- **`max_backups`**: Number of old log files to retain.
- **`max_age_days`**: Retention period for old logs.
- **`compress`**: Whether to compress rotated logs.

---

## Validation Rules

The configuration file and environment variables are validated on load with the following rules:

### General

- `data_dir`: Must be a valid, existing directory and writable.

### Chains

- `name`: Required and non-empty.
- `rpcs`: Must include at least one valid and reachable RPC URL.
- **Empty RPC Behavior**: Ignored from the environment, but required in the final configuration.
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
- `log_level`: Must be one of `debug`, `info`, `warn`, `error`.
- `max_size_mb`: Minimum value of 5.
- `max_backups`: Minimum value of 1.
- `max_age_days`: Minimum value of 1.

---

## Default Values

If the configuration file is not found or incomplete, **Khedra** uses the following defaults:

- **Data directory**: `~/.khedra/data`
- **Logging configuration**:
  - Folder: `~/.khedra/logs`
  - Filename: `khedra.log`
  - Max size: 10 MB
  - Max backups: 3
  - Max age: 10 days
  - Compression: Enabled
  - Log level: `info`
- **Chains**: Only `mainnet` and `sepolia` enabled by default.
- **Services**: All services (`api`, `scraper`, `monitor`, `ipfs`) enabled with default configurations.

---

## Common Commands

1. **Validate Configuration**:
   **Khedra** validates the `config.yaml` file and environment variables automatically on startup.

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
