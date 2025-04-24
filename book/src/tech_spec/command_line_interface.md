# Command-Line Interface

Khedra provides a comprehensive command-line interface (CLI) for interacting with the system. This section details the CLI's design, implementation, and available commands.

## CLI Architecture

The CLI is built using a hierarchical command structure implemented with the Cobra library, providing a consistent and user-friendly interface.

### Design Principles

1. **Consistency**: Uniform command structure and option naming
2. **Discoverability**: Self-documenting with built-in help
3. **Composability**: Commands can be combined in pipelines
4. **Feedback**: Clear status and progress information
5. **Automation-Friendly**: Structured output for scripting

### Implementation Structure

```go
func NewRootCommand() *cobra.Command {
    root := &cobra.Command{
        Use:   "khedra",
        Short: "Khedra is a blockchain indexing and monitoring tool",
        Long:  `A comprehensive tool for indexing, monitoring, and querying EVM blockchains`,
    }
    
    // Add global flags
    root.PersistentFlags().StringVar(&cfgFile, "config", "", "config file path")
    root.PersistentFlags().StringVar(&format, "format", "text", "output format (text, json, csv)")
    
    // Add commands
    root.AddCommand(NewStartCommand())
    root.AddCommand(NewStopCommand())
    root.AddCommand(NewStatusCommand())
    // ... additional commands
    
    return root
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

- `khedra start`: Start all or specified services
  - `--services=service1,service2`: Specify services to start
  - `--foreground`: Run in the foreground (don't daemonize)

- `khedra stop`: Stop all or specified services
  - `--services=service1,service2`: Specify services to stop

- `khedra restart`: Restart all or specified services
  - `--services=service1,service2`: Specify services to restart

- `khedra status`: Show status of all or specified services
  - `--services=service1,service2`: Specify services to check
  - `--verbose`: Show detailed status information

### Index Management Commands

Commands for managing the Unchained Index:

- `khedra index status`: Show index status
  - `--show-gaps`: Display gaps in the index
  - `--analytics`: Show index analytics

- `khedra index rebuild`: Rebuild portions of the index
  - `--start=X`: Starting block number
  - `--end=Y`: Ending block number

- `khedra index verify`: Verify index integrity
  - `--repair`: Attempt to repair issues

### Monitor Commands

Commands for managing address monitors:

- `khedra monitor add ADDRESS [ADDRESS...]`: Add addresses to monitor
  - `--name=NAME`: Assign a name to the monitor
  - `--notifications=webhook,email`: Configure notification methods

- `khedra monitor remove ADDRESS [ADDRESS...]`: Remove monitored addresses
  - `--force`: Remove without confirmation

- `khedra monitor list`: List all monitored addresses
  - `--details`: Show detailed information

- `khedra monitor activity ADDRESS`: Show activity for a monitored address
  - `--from=X`: Starting block number
  - `--to=Y`: Ending block number
  - `--limit=N`: Limit number of results

### Chain Management Commands

Commands for managing blockchain connections:

- `khedra chains list`: List configured chains
  - `--enabled-only`: Show only enabled chains

- `khedra chains add NAME URL`: Add a new chain configuration
  - `--enable`: Enable the chain after adding

- `khedra chains test NAME`: Test connection to a chain
  - `--verbose`: Show detailed test results

### Configuration Commands

Commands for managing Khedra's configuration:

- `khedra config show`: Display current configuration
  - `--redact`: Hide sensitive information
  - `--section=SECTION`: Show only specified section

- `khedra config edit`: Open configuration in an editor

- `khedra config wizard`: Run interactive configuration wizard
  - `--simple`: Run simplified wizard with fewer options

### Utility Commands

Additional utility commands:

- `khedra version`: Show version information
  - `--check-update`: Check for updates

- `khedra cache prune`: Prune old cache data
  - `--older-than=30d`: Prune data older than specified period

- `khedra export`: Export data for external use
  - `--address=ADDR`: Export data for specific address
  - `--format=csv`: Export format
  - `--output=FILE`: Output file path

## Implementation Details

### Command Execution Flow

1. **Parsing**: The CLI parses command-line arguments and flags
2. **Validation**: Command options are validated for correctness
3. **Configuration**: The application configuration is loaded
4. **Execution**: The command is executed with provided options
5. **Output**: Results are formatted according to the specified format

### Output Formatting

The CLI supports multiple output formats:

1. **Text**: Human-readable formatted text (default)
2. **JSON**: Structured JSON for programmatic processing
3. **CSV**: Comma-separated values for spreadsheet import

Output formatting is implemented through a formatter interface:

```go
type OutputFormatter interface {
    Format(data interface{}) ([]byte, error)
}

// Implementations for different formats
type TextFormatter struct{}
type JSONFormatter struct{}
type CSVFormatter struct{}
```

### Error Handling

CLI commands follow consistent error handling patterns:

1. **User Errors**: Clear messages for incorrect usage
2. **System Errors**: Detailed information for system-level issues
3. **Exit Codes**: Specific exit codes for different error types

Example error handling in a command:

```go
func runStatusCommand(cmd *cobra.Command, args []string) error {
    services, err := getServicesToCheck(cmd)
    if err != nil {
        return fmt.Errorf("invalid service selection: %w", err)
    }
    
    status, err := app.GetServiceStatus(services)
    if err != nil {
        return fmt.Errorf("failed to get status: %w", err)
    }
    
    formatter := getFormatter(cmd)
    output, err := formatter.Format(status)
    if err != nil {
        return fmt.Errorf("failed to format output: %w", err)
    }
    
    fmt.Fprintln(cmd.OutOrStdout(), string(output))
    return nil
}
```

### Command Autocompletion

The CLI generates shell completion scripts for popular shells:

- `khedra completion bash`: Generate Bash completion script
- `khedra completion zsh`: Generate Zsh completion script
- `khedra completion fish`: Generate Fish completion script
- `khedra completion powershell`: Generate PowerShell completion script

### Command Documentation

All commands include detailed help information accessible via:

- `khedra --help`: General help
- `khedra command --help`: Command-specific help
- `khedra command subcommand --help`: Subcommand-specific help

This comprehensive CLI design provides users with a powerful and flexible interface for interacting with Khedra's functionality through the command line.
