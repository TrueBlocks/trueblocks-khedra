# Maintenance and Troubleshooting

This chapter covers routine maintenance tasks and solutions to common issues you might encounter when using Khedra.

## Routine Maintenance

### Regular Updates

To keep Khedra running smoothly, periodically check for and install updates:

```bash
# Check current version
khedra version

# Update to the latest version
go get -u github.com/TrueBlocks/trueblocks-khedra/v5

# Rebuild and install
cd <path_for_khedra_github_repo>
git pull --recurse-submodules
go build -o bin/khedra main.go
./bin/khedra version
```

### Log Rotation

Khedra automatically rotates logs based on your configuration, but you should periodically check log usage:

```bash
# Check log directory size
du -sh ~/.khedra/logs

# List log files
ls -la ~/.khedra/logs
```

If logs are consuming too much space, adjust your logging configuration:

```yaml
logging:
  maxSize: 10      # Maximum size in MB before rotation
  maxBackups: 5    # Number of rotated files to keep
  maxAge: 30       # Days to keep rotated logs
  compress: true   # Compress rotated logs
```

### Index Verification

Periodically verify the integrity of your Unchained Index:

```bash
chifra chunks index --check --chain <chain_name>
```

This checks for any gaps or inconsistencies in the index and reports issues.

### Cache Management

You may check on the cache size and prune old caches (by hand) to free up space:

```bash
# Check cache size
chifra status --verbose
```

## Troubleshooting

### Common Issues and Solutions

#### Service Won't Start

**Symptoms:** A service fails to start or immediately stops.

**Solutions:**

1. Check the logs for error messages:

   ```bash
   tail -n 100 ~/.khedra/logs/khedra.log
   ```

2. Verify the service's port isn't in use by another application:

   ```bash
   lsof -i :<port_number>
   ```

3. Ensure the RPC endpoints are accessible:

   ```bash
   chifra status
   ```

4. Try starting with verbose logging:

   ```bash
   TB_KHEDRA_LOGGING_LEVEL=debug TB_KHEDRA_LOGGING_TOFILE=true khedra start
   ```

#### Slow Indexing

**Symptoms:** Indexing is progressing much slower than expected.

**Solutions:**

1. Check RPC endpoint performance:

   ```bash
   chifra status --diagnose
   ```

2. Increase or lower batch size in configuration:

   ```yaml
   services:
     scraper:
       batchSize: 1000  # Increase from default
   ```

3. Monitor system resources to identify bottlenecks:

   ```bash
   top -c -p $(pgrep khedra)
   ```

4. Consider using a faster RPC endpoint or running your own node.

#### Index Gaps

**Symptoms:** The index status shows gaps in block coverage.

**Solutions:**

1. Identify the missing ranges:

   ```bash
   chifra chunks index --check --chain <chain_name>
   ```

2. In very rare cases, you may truncate the index and it will rebuild from that spot. BE CAREFUL -- THIS IS NOT A RECOMMENDED SOLUTION.

   ```bash
   chifra chunks index --truncate <block_number> --chain <chain_name>
   ```

#### API Connection Issues

**Symptoms:** Unable to connect to Khedra's API.

**Solutions:**

1. Verify the API service is running:

   ```bash
   curl http://localhost:8080/status
   ```

2. Check if the configured port is accessible:

   ```bash
   lsof -i :8080
   ```

3. Look for firewall or permission issues:

   ```bash
   sudo lsof -i :8080
   ```

#### IPFS Connectivity Problems

**Symptoms:** Unable to publish or fetch via IPFS.

**Solutions:**

1. Check IPFS service status:

   ```bash
   ps -ef | grep ipfs
   ```

2. Restart khedra

### Log Analysis

Khedra's logs are your best resource for troubleshooting. Here's how to use them effectively:

```bash
# View recent log entries
tail -f ~/.khedra/logs/khedra.log

# Search for error messages
grep -i error ~/.khedra/logs/khedra.log

# Find logs related to a specific service
grep "scraper" ~/.khedra/logs/khedra.log

# Find logs related to a specific address
grep "0x742d35Cc6634C0532925a3b844Bc454e4438f44e" ~/.khedra/logs/khedra.log
```

### Getting Help

If you encounter issues you can't resolve:

1. Check the [Khedra GitHub repository](https://github.com/TrueBlocks/trueblocks-khedra) for known issues
2. Search the [discussions forum](https://github.com/TrueBlocks/trueblocks-khedra/discussions) for similar problems
3. Submit a detailed issue report including:
   - Khedra version (`khedra version`)
   - Relevant log extracts
   - Steps to reproduce the problem
   - Your configuration (with sensitive data redacted)

Regular maintenance and prompt troubleshooting will keep your Khedra installation running smoothly and efficiently.

## Implementation Details

The maintenance and troubleshooting procedures described in this document are implemented in several key files:

### Service Management

- **Service Lifecycle Management**: [`app/action_daemon.go`](/Users/jrush/Development/trueblocks-core/khedra/app/action_daemon.go) - Contains the core service management code that starts, stops, and monitors services
- **Service Health Checks**: Service status monitoring is implemented in the daemon action function

### RPC Connection Management

- **RPC Endpoint Testing**: [`pkg/validate/try_connect.go`](/Users/jrush/Development/trueblocks-core/khedra/pkg/validate/try_connect.go) - Contains the `TestRpcEndpoint` function used to verify endpoints are functioning correctly
- **RPC Validation**: [`app/has_valid_rpc.go`](/Users/jrush/Development/trueblocks-core/khedra/app/has_valid_rpc.go) - Implements validation logic for RPC endpoints

### Logging System

- **Log Configuration**: Defined in the `Logging` struct in [`pkg/types/general.go`](/Users/jrush/Development/trueblocks-core/khedra/pkg/types/general.go) which handles log rotation and management
- **Logger Implementation**: Custom logger in [`pkg/types/custom_logger.go`](/Users/jrush/Development/trueblocks-core/khedra/pkg/types/custom_logger.go) that provides structured logging capabilities

### Error Recovery

The troubleshooting techniques described are supported by robust error handling throughout the codebase, especially in:

- **Service error handling**: Found in the daemon action function
- **Validation error reporting**: Implemented in the validation framework
- **Index management functions**: For identifying and fixing gaps in the index
