# Monitor Service

The Monitor Service watches specified Ethereum addresses across configured chains and executes custom commands when new blocks are detected.

## Overview

The monitor service provides automated address monitoring with:
- Per-chain watchlist management
- Custom command execution with template variables
- Batch processing for efficiency
- Parallel command execution
- Automatic chain head tracking

## Configuration

### Enable Monitor Service

Edit `~/.khedra/config.yaml`:

```yaml
general:
  monitorsFolder: "~/.khedra/monitors"  # Directory for monitor files

services:
  monitor:
    enabled: true      # Enable the monitor service
    sleep: 12          # Seconds between iterations when caught up
    batchSize: 8       # Addresses to freshen per batch
    concurrency: 4     # Parallel command workers
```

### Directory Structure

The monitor service uses a flat directory structure:

```
~/.khedra/monitors/
  watchlist-mainnet.txt      # Mainnet addresses
  commands-mainnet.yaml      # Mainnet commands
  watchlist-gnosis.txt       # Gnosis chain addresses
  commands-gnosis.yaml       # Gnosis chain commands
```

Files follow the pattern: `watchlist-{chain}.txt` and `commands-{chain}.yaml`

## Watchlist Files

### Format

Each line contains an address with an optional starting block:

```
address[,starting_block]
```

### Example (`watchlist-mainnet.txt`)

```
# Mainnet addresses to monitor
0x1234567890123456789012345678901234567890
0xabcdefabcdefabcdefabcdefabcdefabcdefabcd,15000000
```

### Rules

- One address per line
- Address can be:
  - Ethereum address (0x + 40 hex chars)
  - ENS name (e.g., trueblocks.eth)
- Optional starting block after comma
- Lines starting with `#` are comments
- Empty lines are ignored

## Commands Files

### Format

YAML format with command definitions:

```yaml
commands:
  - id: command_identifier
    command: executable_name
    arguments:
      - arg1
      - arg2
    output: /path/to/output
```

### Template Variables

Commands support dynamic template variables:

| Variable | Description | Example |
|----------|-------------|---------|
| `{address}` | Monitored address | `0x1234...` |
| `{chain}` | Chain name | `mainnet` |
| `{first_block}` | First block in iteration | `15000000` |
| `{last_block}` | Last block in iteration | `15000100` |
| `{block_count}` | Number of blocks | `101` |

### Example (`commands-mainnet.yaml`)

```yaml
commands:
  # Export transactions for monitored address
  - id: export_transactions
    command: chifra
    arguments:
      - export
      - "{address}"
      - --chain
      - "{chain}"
      - --first_block
      - "{first_block}"
      - --last_block
      - "{last_block}"
    output: /dev/null

  # Export logs
  - id: export_logs
    command: chifra
    arguments:
      - export
      - --logs
      - "{address}"
      - --chain
      - "{chain}"
    output: /dev/null
```

## Service Behavior

### Initialization

On startup, the monitor service:
1. Creates `monitorsFolder` if it doesn't exist
2. Loads watchlists for each enabled chain
3. Parses command definitions
4. Reports loaded monitors and commands

### Processing Loop

For each iteration:
1. Batch freshen addresses using `chifra export --freshen`
2. Get current chain head
3. For each address:
   - Calculate block range to process
   - Expand template variables
   - Execute configured commands in parallel
4. Sleep if caught up, otherwise continue immediately

### Silent Operation

The monitor service will silently sleep if:
- No watchlist files exist
- Watchlist files are empty
- No commands are configured
- No enabled chains are found

### Error Handling

- **Watchlist errors**: Service logs error and skips that chain
- **Command failures**: Logged but don't stop processing
- **Fail-early**: If >50% of monitors fail, service stops with error

## Getting Started

### 1. Copy Example Files

```bash
# Copy examples to monitors directory
cp examples/monitors/watchlist-mainnet.txt ~/.khedra/monitors/
cp examples/monitors/commands-mainnet.yaml ~/.khedra/monitors/
```

### 2. Add Addresses

Edit the watchlist file:

```bash
# Add your monitored addresses
echo "0xYourAddressHere" >> ~/.khedra/monitors/watchlist-mainnet.txt
```

### 3. Customize Commands

Edit `~/.khedra/monitors/commands-mainnet.yaml` to define what happens when new blocks appear for monitored addresses.

### 4. Enable Service

Update `~/.khedra/config.yaml`:

```yaml
services:
  monitor:
    enabled: true
```

### 5. Start Khedra

```bash
khedra daemon
```

## Advanced Usage

### Multiple Chains

Create separate watchlist and commands files for each chain:

```bash
# Add gnosis monitoring
cp examples/monitors/watchlist-gnosis.txt ~/.khedra/monitors/
cp examples/monitors/commands-gnosis.yaml ~/.khedra/monitors/

# Edit gnosis watchlist
nano ~/.khedra/monitors/watchlist-gnosis.txt
```

Ensure the chain is enabled in `config.yaml`.

### Performance Tuning

Adjust service parameters based on your needs:

```yaml
services:
  monitor:
    batchSize: 8       # Increase for more addresses
    concurrency: 4     # Increase for more parallel processing
    sleep: 12          # Decrease for more frequent checks
```

### Custom Commands

Create any command that uses the template variables:

```yaml
commands:
  - id: custom_export
    command: /path/to/script.sh
    arguments:
      - "{address}"
      - "{chain}"
      - "{first_block}"
      - "{last_block}"
    output: /tmp/monitor.log
```

## Monitoring and Control

### View Service Status

```bash
# Check if monitor is running and paused state
curl http://localhost:8338/status

# Get detailed status
curl http://localhost:8338/isPaused?name=monitor
```

### Pause/Unpause

```bash
# Pause monitor service
curl -X POST http://localhost:8338/pause?name=monitor

# Resume monitor service
curl -X POST http://localhost:8338/unpause?name=monitor
```

### Restart

```bash
# Restart monitor service
curl -X POST http://localhost:8338/restart?name=monitor
```

## Troubleshooting

### Service Not Starting

1. Check logs: `~/.khedra/logs/khedra.log`
2. Verify watchlist files exist and are readable
3. Validate address format (Ethereum address or ENS name)
4. Check commands YAML syntax

### No Processing Activity

1. Verify service is enabled: `enabled: true`
2. Check if service is paused: `curl http://localhost:8338/isPaused?name=monitor`
3. Ensure watchlist files contain valid addresses
4. Verify chains are enabled in config

### Command Failures

1. Check command executable exists and is in PATH
2. Verify template variables expand correctly
3. Test command manually with actual values
4. Review logs for error messages

## Examples

See `examples/monitors/` directory for complete examples:
- `watchlist-mainnet.txt` - Mainnet watchlist example
- `commands-mainnet.yaml` - Mainnet commands example
- `watchlist-gnosis.txt` - Gnosis chain watchlist example
- `commands-gnosis.yaml` - Gnosis chain commands example
- `README.md` - Detailed setup instructions

## See Also

- [Configuration Guide](../config.yaml.example)
- [Service Management](../README.md#service-management)
- [chifra export Documentation](https://chifra.trueblocks.io/docs/chifra/accounts/#chifra-export)
