# Using Khedra

This chapter covers the practical aspects of working with Khedra once it's installed and configured.

## Understanding Khedra's Command Structure

Khedra provides a streamlined set of commands designed to index, monitor, serve, and share blockchain data:

```text
NAME:
   khedra - A tool to index, monitor, serve, and share blockchain data

USAGE:
   khedra [global options] command [command options]

VERSION:
   v5.1.0

COMMANDS:
   init     Initializes Khedra
   daemon   Runs Khedras services
   config   Manages Khedra configuration
   help, h  Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --help, -h     show help
   --version, -v  print the version
```

## Getting Started with Khedra

### Initializing Khedra

Before using Khedra, you need to initialize it. This sets up the necessary data structures and configurations:

```bash
khedra init
```

During initialization, Khedra will:

- Set up its directory structure
- Configure initial settings
- Prepare the system for indexing blockchain data

### Managing Configuration

To view or modify Khedra's configuration:

```text
khedra config [show | edit]
```

The configuration command allows you to:

- View current settings
- Update connection parameters
- Adjust service behaviors
- Configure chain connections

### Running Khedra's Services

To start Khedra's daemon services:

```bash
khedra daemon
```

This command:

- Starts the indexing service
- Enables the API server if configured
- Processes monitored addresses
- Handles data serving capabilities

You can use various options with the daemon command to customize its behavior. For detailed options:

```text
khedra daemon --help
```

## Common Workflows

### Basic Setup

1. Install Khedra using the installation instructions
2. Initialize the system:

   ```text
   khedra init
   ```

3. Configure as needed:

   ```text
   khedra config edit
   ```

4. Start the daemon services:

   ```text
   khedra daemon
   ```

### Checking System Status

You can view the current status of Khedra by examining the daemon process:

```text
curl http://localhost:8338/status | jq
```

- **Note:** The port for the above command defaults to one of 8338, 8337, 8336 or 8335 in that order whichever one is first available. If none of those ports is available, the daemon will not start.

### Accessing the Data API

If so configured, when the daemon is running, it provides API endpoints for accessing blockchain data. The default configuration typically serves on:

```curl
curl http://localhost:8080/status
```

See the [API documentation](https://trueblocks.io/api/) for more details on available endpoints and their usage.

## Getting Help

Each command provides detailed help information. To access help for any command:

```bash
khedra [command] --help
```

For general help:

```bash
khedra --help
```

### Version Information

To check which version of Khedra you're running:

```bash
khedra --version
```

## Advanced Usage

For more detailed information about advanced operations and configurations, please refer to the documentation for each specific command:

```bash
khedra init --help
khedra daemon --help
khedra config --help
```

The next chapter covers advanced operations for users who want to maximize Khedra's capabilities.

## Implementation Details

The command structure and functionality described in this section are implemented in these Go files:

### Core Command Structure

- **CLI Framework**: [`app/cli.go`](/Users/jrush/Development/trueblocks-core/khedra/app/cli.go) - Defines the top-level command structure using the `urfave/cli` package

### Command Implementations

- **Init Command**: [`app/action_init.go`](/Users/jrush/Development/trueblocks-core/khedra/app/action_init.go) - Handles initialization of Khedra
  - **Welcome Screen**: [`app/action_init_welcome.go`](/Users/jrush/Development/trueblocks-core/khedra/app/action_init_welcome.go)
  - **General Settings**: [`app/action_init_general.go`](/Users/jrush/Development/trueblocks-core/khedra/app/action_init_general.go)
  - **Services Config**: [`app/action_init_services.go`](/Users/jrush/Development/trueblocks-core/khedra/app/action_init_services.go)
  - **Chain Config**: [`app/action_init_chains.go`](/Users/jrush/Development/trueblocks-core/khedra/app/action_init_chains.go)
  - **Summary Screen**: [`app/action_init_summary.go`](/Users/jrush/Development/trueblocks-core/khedra/app/action_init_summary.go)

- **Daemon Command**: [`app/action_daemon.go`](/Users/jrush/Development/trueblocks-core/khedra/app/action_daemon.go) - Implements the daemon service that runs the various Khedra services

- **Config Commands**: 
  - [`app/action_config_show.go`](/Users/jrush/Development/trueblocks-core/khedra/app/action_config_show.go) - Displays current configuration
  - [`app/action_config_edit.go`](/Users/jrush/Development/trueblocks-core/khedra/app/action_config_edit.go) - Opens configuration in editor

- **Version Command**: [`app/action_version.go`](/Users/jrush/Development/trueblocks-core/khedra/app/action_version.go) - Shows version information

### Helper Functions

- **Command Line Arguments**: [`app/args.go`](/Users/jrush/Development/trueblocks-core/khedra/app/args.go) - Processes command line arguments
- **Help System**: [`app/help_system.go`](/Users/jrush/Development/trueblocks-core/khedra/app/help_system.go) - Provides help text for commands
