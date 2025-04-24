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

    If no `config.yaml` file is found, **khedra** creates a default configuration in its default location.

5. **Using Environment Variables**:

   You may override configuration options using environment variables, each of which must take the form `TB_KHEDRA_<section>_<key>`.

   For example, the following overrides the `general.dataFolder` value.

     ```bash
     export TB_KHEDRA_GENERAL_DATAFOLDER="/path/override"
     ```

    You'll notice that underbars (`_`) in the `<key>` names are not needed.

---

## Configuration File Format

The `config.yaml` file (shown here with default values) is structured as follows:

```yaml
# Khedra Configuration File
# Version: 2.0

general:
  dataFolder: "~/.khedra/data"  # See note 1
  strategy: "download"          # How to build the Unchained Index [download* | scrape]
  detail: "index"               # How detailed to log the processes [index* | blooms]

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
      - "rpc_endpoint_for_gnosis" # must be a reachable URL if the chain is enabled
    enabled: false                # in this example, this chain is disabled
  optimism:
    rpcs:
      - "rpc_endpoint_for_optimism"
    enabled: false

services:                          # See note 5
  scraper:                         # Required. (One of: api, scraper, monitor, ipfs, control)
    enabled: true                  # `true` if the service is enabled
    sleep: 12                      # Seconds between scraping batches (see note 6)
    batchSize: 500                 # Number of blocks to process in a batch (range: 50-10000)

  monitor:
    enabled: true
    sleep: 12                      # Seconds between scraping batches (see note 6)
    batchSize: 500                 # Number of blocks processed in a batch (range: 50-10000)

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
  toFile: false                    # If true, will write to above file. Screen only otherwise
  level: "info"                    # One of: debug, info, warn, error
  maxSize: 10                      # Max log file size in MB
  maxBackups: 5                    # Number of backup log files to keep
  maxAge: 30                       # Number of days to retain old logs
  compress: true                   # Whether to compress backup logs
```

**Notes:**

1. The `dataFolder` value must be a valid, existing directory that is writable. You may wish to change this value to a location with suitable disc space. Depending on configuration, the Unchained Index and binary caches may get large (> 200GB in some cases).

2. The `chains` section is required. At least one chain must be enabled. An RPC for `mainnet` is required even if `mainnet` is disabled. The software reads `mainnet` smart contracts (such as the *Unchained Index* and *UniSwap*) during normal operation.

3. [This repository](https://github.com/ethereum-lists/chains) is used to identify chain names. Using consistent chain names aides in sharing indexes. Use these values in your configuration if you wish to fully participate in sharing the *Unchained Index*.

4. The `services` section is required. At least one service must be enabled.

5. When a `scraper` or `monitor` is "catching up" to a chain, the `sleep` value is ignored.

---

## Using Environment Variables

**Khedra** allows configuration values to be overridden at runtime using environment variables. The value of an environment variable takes precedence over the defaults and the configuration file.

### Naming Evirnment Variables

The environment variable naming convention is:

`TB_KHEDRA_<section>_<key>`

For example:

- To override the `general.dataFolder` value:

  ```bash
  export TB_KHEDRA_GENERAL_DATAFOLDER="/path/override"
  ```

- To override `logging.level`:

  ```bash
  export TB_KHEDRA_LOGGING_LEVEL="debug"
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

- **`dataFolder`**: The location where **khedra** stores all of its data. This directory must exist and be writable.
- **`strategy`**: The strategy used to initialize the *Unchained Index*. With `download` (the default), the Unchained Index smart contract will be consulted the index will be downloaded from IPFS. With `scrape` the entire index will be created from scratch. The former takes a lot less time, but relies on values created by a third party. The later (`scrape`) uses only the RPC as a source which means it takes significantly longer, but is most secure as no third-party trust is required.
- **`detail`**: The detail level of the dowloaded or scraped index. With `index` both the Bloom filters and the Index chunks are either downloaded or build (depending on `strategy`). With `blooms`, only the Bloom filters are retained. Index chunks are downloaded on an as needed basis through `chifra export`. The former is much larger and takes much longer to `download` (if `strategy` is `scrape` no time savings is seen). The later is much smaller and faster to `download`. Downloading or creating the full `index` is the default.

### Chains (Blockchains)

Defines the blockchain networks to interact with. Each chain must have:

- **`name`**: Chain name (e.g., `mainnet`).
- **`rpcs`**: List of RPC endpoints. At least one valid and reachable endpoint is required. `mainnet` RPC is required, but you are not required to index it.
- **`enabled`**: Whether the chain is being actively indexed.

#### Behavior for Empty RPCs

- If the `RPCs` field is empty in the environment, it is ignored and the configuration file's value is preserved.
- If the `RPCs` field is empty in the final configuration (after merging), the chain is treated as it would be if it were disabled.

---

### Services (API, Scraper, Monitor, IPFS)

Defines various services provided by **Khedra**. Supported services:

- **API**:
  - An API server for the `chifra` command line interface. See [API Documentation](https://trueblocks.io/api/) for details.
  - Requires `port` to be specified in the configuration.
- **Scraper** and **Monitor**:
  - These two services are used to scrape and monitor the blockchain data respectively. Each runs "periodically" to keep the index or monitor data up to date.
  - **`sleep`**: Duration (seconds) between operations.
  - **`batchSize`**: Number of blocks to process in each operation (50-10,000).
- **IPFS**:
  - A service for interacting with IPFS (InterPlanetary File System). This service starts an internal IPFS daemon if it's not already running. The scraper service may use IPFS to pin and share the index if so configured.
  - Requires `port` to be specified.

### Logging Configuration

Controls the application's logging behavior:

- **`folder`**: Directory for storing logs.
- **`filename`**: Name of the log file.
- **`toFile`**: If `true`, logs are written to the specified file. If `false`, logs are only printed to the console.
- **`level`**: Logging level. Possible values: `debug`, `info`, `warn`, `error`.
- **`maxSize`**: Maximum log file size before rotation.
- **`maxBackups`**: Number of old log files to retain.
- **`maxAge`**: Retention period for old logs.
- **`compress`**: Whether to compress rotated logs.

---

## Validation Rules

The configuration file and environment variables are validated when the program starts with the following rules:

### General

- `dataFolder`: Must be a valid, existing directory and writable.
- `strategy`: Must be either `download` or `scrape`.
- `detail`: Must be either `index` or `blooms`.

### Chains

- `name`: Required and non-empty.
- `rpcs`: Must include at least one valid and reachable RPC URL.
- **Empty RPC Behavior**: Ignored from the environment, but required in the final configuration.
- `enabled`: Defaults to `false` if not specified.

#### Notes on chains section

1. The `mainnet` RPC is required even if indexing the chain is disabled. The software reads `mainnet` smart contracts (such as the *Unchained Index* and *UniSwap*) during normal operation.
2. It is always best to have a dedicated RPC endpoint. If you are using a public RPC endpoint, be sure to check the rate limits and usage policies of the provider and set the `sleep` and `batchSize` values for the services appropriately. Some providers (all providers?) will block or throttle requests if they exceed certain limits.

### Services

- `name`: Required and non-empty. Must be one of `api`, `scraper`, `monitor`, `ipfs`.
- `enabled`: Defaults to `false` if not specified.
- `port`: For API and IPFS services, must be between 1024 and 65535. Ignored for other services.
- `sleep`: Must be non-negative. Ignored by API and IPFS services.
- `batchSize`: Must be between 50 and 10,000. Ignored by API and IPFS services.

### Logging

- `folder`: Must exist and be writable.
- `filename`: Must end with `.log`.
- `toFile`: Must be `true` or `false`.
- `level`: Must be one of `debug`, `info`, `warn`, `error`.
- `maxSize`: Minimum value of 5.
- `maxBackups`: Minimum value of 1.
- `maxAge`: Minimum value of 1.

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
- **Chains**: Only `mainnet` and `gnosis` enabled by default.
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
   export TB_KHEDRA_GENERAL_DATAFOLDER="/new/path"
   ./khedra
   ```

For additional details, see the technical specification.

## Implementation Details

The configuration system and initialization described in this section are implemented in these Go files:

- **Configuration Loading**: [`app/config.go`](/Users/jrush/Development/trueblocks-core/khedra/app/config.go) - Contains the `LoadConfig()` function that loads, merges, and validates configuration from files and environment variables

- **Configuration Validation**: 
  - [`pkg/types/general.go`](/Users/jrush/Development/trueblocks-core/khedra/pkg/types/general.go) - Validates general settings
  - [`pkg/types/chain.go`](/Users/jrush/Development/trueblocks-core/khedra/pkg/types/chain.go) - Validates chain settings
  - [`pkg/types/service.go`](/Users/jrush/Development/trueblocks-core/khedra/pkg/types/service.go) - Validates service settings

- **Environment Variables**: [`pkg/types/apply_env.go`](/Users/jrush/Development/trueblocks-core/khedra/pkg/types/apply_env.go) - Contains functions for applying environment variables to the configuration

- **Initialization Command**: [`app/action_init.go`](/Users/jrush/Development/trueblocks-core/khedra/app/action_init.go) - Implements the `init` command to set up the initial configuration

- **Folder and Path Management**: Found in the `initializeFolders()` function in [`app/config.go`](/Users/jrush/Development/trueblocks-core/khedra/app/config.go) which ensures required directories exist
