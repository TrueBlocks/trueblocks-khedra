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

### Service-Specific Troubleshooting

#### Scraper Service Issues

**Symptoms:** Scraper service fails to start, stops unexpectedly, or indexes slowly.

**Common Issues and Solutions:**

1. **RPC Connection Failures:**
   ```bash
   # Test RPC connectivity
   curl -X POST -H "Content-Type: application/json" \
     --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
     http://your-rpc-endpoint
   
   # Check RPC provider limits
   grep -i "rate limit\|too many requests" ~/.khedra/logs/khedra.log
   ```

2. **Batch Size Optimization:**
   ```yaml
   # For fast RPC endpoints
   services:
     scraper:
       batchSize: 2000
       sleep: 5
   
   # For slower/limited RPC endpoints  
   services:
     scraper:
       batchSize: 100
       sleep: 30
   ```

3. **Memory Issues:**
   ```bash
   # Monitor scraper memory usage
   ps -o pid,vsz,rss,comm -p $(pgrep -f "scraper")
   
   # Reduce batch size if memory usage is high
   ```

4. **Scraper-Specific Log Analysis:**
   ```bash
   # Filter scraper logs
   grep "scraper" ~/.khedra/logs/khedra.log | tail -50
   
   # Look for specific errors
   grep -E "error|failed|timeout" ~/.khedra/logs/khedra.log | grep scraper
   ```

#### Monitor Service Issues

**Symptoms:** Monitor service doesn't detect address activity or sends duplicate notifications.

**Common Issues and Solutions:**

1. **No Monitored Addresses:**
   ```bash
   # Check if addresses are properly configured
   chifra list --monitors
   
   # Add addresses to monitor
   chifra monitors --addrs 0x742d35Cc6634C0532925a3b844Bc454e4438f44e
   ```

2. **Monitor Service Dependencies:**
   ```bash
   # Ensure scraper is running for real-time monitoring
   curl http://localhost:8080/api/v1/services/scraper
   
   # Check if index is up to date
   chifra status --index
   ```

3. **Monitor Configuration Issues:**
   ```yaml
   services:
     monitor:
       enabled: true
       sleep: 12        # Check every 12 seconds
       batchSize: 100   # Process 100 addresses at once
   ```

4. **Monitor-Specific Logs:**
   ```bash
   # Filter monitor logs
   grep "monitor" ~/.khedra/logs/khedra.log | tail -50
   
   # Check for address activity detection
   grep -i "activity\|appearance" ~/.khedra/logs/khedra.log
   ```

#### API Service Issues

**Symptoms:** API service returns errors, timeouts, or incorrect data.

**Common Issues and Solutions:**

1. **Port Conflicts:**
   ```bash
   # Check if API port is available
   lsof -i :8080
   
   # Change API port if needed
   export TB_KHEDRA_SERVICES_API_PORT=8081
   ```

2. **API Performance Issues:**
   ```bash
   # Test API response time
   time curl http://localhost:8080/status
   
   # Check for slow queries
   grep -E "slow|timeout" ~/.khedra/logs/khedra.log | grep api
   ```

3. **API Authentication Issues:**
   ```bash
   # Verify API is accessible
   curl -v http://localhost:8080/api/v1/services
   
   # Check for auth-related errors
   grep -i "auth\|unauthorized" ~/.khedra/logs/khedra.log
   ```

4. **Data Consistency Issues:**
   ```bash
   # Compare API data with direct index queries
   chifra list 0x742d35Cc6634C0532925a3b844Bc454e4438f44e
   curl http://localhost:8080/api/v1/list/0x742d35Cc6634C0532925a3b844Bc454e4438f44e
   ```

#### IPFS Service Issues

**Symptoms:** IPFS service fails to start, can't connect to network, or sharing fails.

**Common Issues and Solutions:**

1. **IPFS Daemon Issues:**
   ```bash
   # Check IPFS daemon status
   ps aux | grep ipfs
   
   # Restart IPFS if needed
   curl -X POST http://localhost:8080/api/v1/services/ipfs/restart
   ```

2. **IPFS Port Conflicts:**
   ```bash
   # Check IPFS ports
   lsof -i :5001  # IPFS API port
   lsof -i :4001  # IPFS swarm port
   
   # Configure different IPFS port
   export TB_KHEDRA_SERVICES_IPFS_PORT=5002
   ```

3. **IPFS Network Connectivity:**
   ```bash
   # Test IPFS connectivity
   curl http://localhost:5001/api/v0/id
   
   # Check peer connections
   curl http://localhost:5001/api/v0/swarm/peers
   ```

4. **Index Sharing Issues:**
   ```bash
   # Check IPFS pinning status
   curl http://localhost:5001/api/v0/pin/ls
   
   # Verify index chunk uploads
   grep -i "ipfs\|pin" ~/.khedra/logs/khedra.log
   ```

#### Control Service Issues

**Symptoms:** Cannot manage other services via API or CLI commands fail.

**Common Issues and Solutions:**

1. **Control Service Availability:**
   ```bash
   # Verify control service is running
   curl http://localhost:8080/api/v1/services
   
   # Check control service logs
   grep "control" ~/.khedra/logs/khedra.log
   ```

2. **Service Management Failures:**
   ```bash
   # Test individual service control
   curl -X POST http://localhost:8080/api/v1/services/scraper/status
   
   # Check for permission issues
   grep -i "permission\|access denied" ~/.khedra/logs/khedra.log
   ```

3. **Configuration Issues:**
   ```bash
   # Verify control service configuration
   khedra config show | grep -A5 -B5 control
   
   # Test configuration validation
   khedra config validate
   ```

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
