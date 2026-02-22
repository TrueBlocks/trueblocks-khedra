# Monitor Service Example Files

These example files demonstrate how to configure the khedra monitor service.

## Installation

1. **Create the monitors directory** (if it doesn't exist):
   ```bash
   mkdir -p ~/.khedra/monitors
   ```

2. **Copy the example files** for the chains you want to monitor:
   ```bash
   # For mainnet
   cp watchlist-mainnet.txt ~/.khedra/monitors/
   cp commands-mainnet.yaml ~/.khedra/monitors/
   
   # For gnosis chain
   cp watchlist-gnosis.txt ~/.khedra/monitors/
   cp commands-gnosis.yaml ~/.khedra/monitors/
   ```

3. **Edit the watchlist files** to add your monitored addresses:
   ```bash
   # Add addresses to the watchlist
   echo "0xYourAddressHere" >> ~/.khedra/monitors/watchlist-mainnet.txt
   ```

4. **Enable the monitor service** in your `~/.khedra/config.yaml`:
   ```yaml
   services:
     monitor:
       enabled: true
       sleep: 12
       batchSize: 8
       concurrency: 4
   ```

5. **Start khedra**:
   ```bash
   khedra daemon
   ```

## File Formats

### Watchlist Files

Format: `address[,starting_block]`

Example:
```
0x1234567890123456789012345678901234567890
0xabcdefabcdefabcdefabcdefabcdefabcdefabcd,15000000
```

### Commands Files

YAML format defining commands to execute for each monitored address.

Template variables available:
- `{address}` - The monitored address
- `{chain}` - The chain name
- `{first_block}` - First block in this iteration
- `{last_block}` - Last block in this iteration
- `{block_count}` - Number of blocks being processed

Example:
```yaml
commands:
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
```

## Behavior

- The monitor service will **silently sleep** if:
  - No watchlist files exist for enabled chains
  - Watchlist files are empty
  - No commands are defined
- Khedra will automatically create the `~/.khedra/monitors` directory if it doesn't exist
- The service respects the `enabled` flag in the configuration
