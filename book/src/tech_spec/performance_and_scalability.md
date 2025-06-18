# Performance and Scalability

This section details Khedra's performance characteristics, optimization strategies, and scalability considerations.

## Performance Tuning Guide

### Optimal Configuration for Different Use Cases

#### Light Usage (Personal Development/Testing)
```yaml
services:
  scraper:
    enabled: true
    sleep: 30        # Longer sleep for less aggressive indexing
    batchSize: 100   # Smaller batches to reduce memory usage
  monitor:
    enabled: false   # Disable if not needed
    sleep: 60
    batchSize: 50
  api:
    enabled: true
    port: 8080       # Standard configuration
  ipfs:
    enabled: false   # Disable to reduce resource usage
```

**Expected Performance:** 5-10 blocks/sec indexing, minimal resource usage

#### Standard Usage (Regular Development/Analysis)
```yaml
services:
  scraper:
    enabled: true
    sleep: 12        # Balanced sleep interval
    batchSize: 500   # Default batch size
  monitor:
    enabled: true    # Enable for address tracking
    sleep: 12
    batchSize: 100
  api:
    enabled: true
    port: 8080
  ipfs:
    enabled: true    # Enable for collaboration
    port: 8083
```

**Expected Performance:** 15-25 blocks/sec indexing, moderate resource usage

#### High-Performance Usage (Production/Heavy Analysis)
```yaml
services:
  scraper:
    enabled: true
    sleep: 5         # Aggressive indexing
    batchSize: 2000  # Large batches for efficiency
  monitor:
    enabled: true
    sleep: 5         # Fast monitoring
    batchSize: 500
  api:
    enabled: true
    port: 8080
  ipfs:
    enabled: true
    port: 8083
```

**Expected Performance:** 25-40 blocks/sec indexing, high resource usage

### Batch Size Optimization Guidelines

#### Factors Affecting Optimal Batch Size

1. **RPC Endpoint Performance:**
   - Fast/unlimited RPC: 1000-5000 blocks
   - Standard RPC: 500-1000 blocks  
   - Rate-limited RPC: 50-200 blocks

2. **Available Memory:**
   - 8GB+ RAM: 1000-2000 blocks
   - 4-8GB RAM: 500-1000 blocks
   - <4GB RAM: 100-500 blocks

3. **Network Latency:**
   - Local RPC node: 2000-5000 blocks
   - Same region: 1000-2000 blocks
   - Remote/high latency: 100-500 blocks

#### Batch Size Tuning Process

1. **Start with defaults** (500 blocks)
2. **Monitor performance** metrics
3. **Adjust based on bottlenecks:**
   - If RPC timeouts: decrease batch size
   - If memory issues: decrease batch size
   - If CPU idle time: increase batch size
   - If slow overall progress: increase batch size

### Sleep Interval Recommendations

#### Based on System Resources

**High-End Systems (8+ cores, 16GB+ RAM):**
- Scraper: 5-10 seconds
- Monitor: 5-10 seconds

**Mid-Range Systems (4-8 cores, 8-16GB RAM):**
- Scraper: 10-15 seconds  
- Monitor: 10-15 seconds

**Resource-Constrained Systems (<4 cores, <8GB RAM):**
- Scraper: 20-30 seconds
- Monitor: 30-60 seconds

#### Based on RPC Provider

**Unlimited/Premium RPC:**
- Scraper: 5-10 seconds
- Monitor: 5-10 seconds

**Standard RPC with rate limits:**
- Scraper: 15-30 seconds
- Monitor: 30-60 seconds

**Free/heavily limited RPC:**
- Scraper: 60-120 seconds
- Monitor: 120-300 seconds

### RPC Endpoint Optimization

#### Choosing RPC Providers

**Recommended for High Performance:**
1. Local RPC node (best performance)
2. Premium providers (Alchemy, Infura Pro)
3. Archive nodes with trace support

**Configuration for Multiple RPC Endpoints:**
```yaml
chains:
  mainnet:
    rpcs:
      - "https://eth-mainnet.alchemyapi.io/v2/YOUR_KEY"
      - "https://mainnet.infura.io/v3/YOUR_KEY"  
      - "https://rpc.ankr.com/eth"
    enabled: true
```

#### RPC Performance Testing

```bash
# Test RPC response time
time curl -X POST -H "Content-Type: application/json" \
  --data '{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}' \
  https://your-rpc-endpoint

# Test batch request performance  
time curl -X POST -H "Content-Type: application/json" \
  --data '[{"jsonrpc":"2.0","method":"eth_getBlockByNumber","params":["0x1000000",false],"id":1}]' \
  https://your-rpc-endpoint
```

### System Resource Monitoring

#### Key Metrics to Monitor

1. **CPU Usage:**
   ```bash
   top -p $(pgrep khedra)
   ```

2. **Memory Usage:**
   ```bash
   ps -o pid,vsz,rss,comm -p $(pgrep khedra)
   ```

3. **Disk I/O:**
   ```bash
   iotop -p $(pgrep khedra)
   ```

4. **Network Usage:**
   ```bash
   nethogs -p $(pgrep khedra)
   ```

#### Performance Thresholds

**CPU Usage:**
- <50%: Can increase batch size or decrease sleep
- 50-80%: Optimal range
- >80%: Decrease batch size or increase sleep

**Memory Usage:**
- <2GB: Can increase batch size
- 2-4GB: Monitor for memory leaks
- >4GB: Decrease batch size

**Disk I/O:**
- High read: Index queries are efficient
- High write: Indexing in progress (normal)
- Very high write: May need to reduce batch size

### Scaling Considerations

#### Horizontal Scaling Strategies

1. **Chain Separation:**
   - Run separate Khedra instances per blockchain
   - Distribute chains across multiple servers
   - Use load balancer for API access

2. **Service Separation:**
   - Run API service on separate instances
   - Dedicated IPFS nodes for data sharing
   - Centralized monitoring service

3. **Geographic Distribution:**
   - Deploy close to RPC providers
   - Regional API instances for lower latency
   - IPFS network for global data sharing

#### Vertical Scaling Guidelines

**Memory Scaling:**
- 8GB: Single chain, moderate usage
- 16GB: Multiple chains or heavy usage
- 32GB+: High-performance production usage

**CPU Scaling:**
- 4 cores: Basic usage
- 8 cores: Standard production
- 16+ cores: High-performance or multiple chains

**Storage Scaling:**
- SSD required for optimal performance
- 100GB per chain per year (estimate)
- Consider compression and archival strategies

## Service Metrics and Monitoring

### Available Performance Metrics

Each Khedra service exposes performance metrics that can be accessed through the Control Service API. These metrics provide insight into service health, performance, and resource utilization.

#### Control Service Metrics

**Service Status Metrics:**
- `uptime`: Service runtime duration since last start
- `state`: Current service state (running, paused, stopped, etc.)
- `last_started`: Timestamp of last service start
- `restart_count`: Number of times service has been restarted
- `health_score`: Overall service health indicator (0-100)

**System Resource Metrics:**
- `memory_usage_bytes`: Current memory consumption
- `cpu_usage_percent`: Current CPU utilization
- `goroutines_count`: Number of active goroutines
- `gc_cycles`: Garbage collection statistics

#### Scraper Service Metrics

**Indexing Performance:**
- `blocks_processed_total`: Total number of blocks indexed
- `blocks_per_second`: Current indexing throughput
- `batch_size_current`: Current batch size setting
- `batch_processing_time_ms`: Average time per batch
- `index_chunks_created`: Number of index chunks generated
- `appearances_extracted_total`: Total address appearances found

**RPC Performance:**
- `rpc_requests_total`: Total RPC requests made
- `rpc_requests_failed`: Number of failed RPC requests
- `rpc_response_time_ms`: Average RPC response time
- `rpc_rate_limit_hits`: Number of rate limit encounters
- `rpc_endpoint_health`: Status of each configured RPC endpoint

**Processing State:**
- `current_block_number`: Latest block being processed
- `target_block_number`: Target block (chain tip)
- `blocks_behind`: Number of blocks behind chain tip
- `indexing_progress_percent`: Overall indexing completion percentage

#### Monitor Service Metrics

**Monitoring Performance:**
- `addresses_monitored`: Number of addresses being tracked
- `monitoring_checks_total`: Total monitoring checks performed
- `activity_detected_total`: Number of activities detected
- `notifications_sent_total`: Number of notifications dispatched
- `false_positives`: Number of false positive detections

**Detection Metrics:**
- `detection_latency_ms`: Time from block to activity detection
- `monitoring_batch_size`: Current batch size for monitoring
- `monitoring_frequency_seconds`: Current monitoring interval

#### API Service Metrics

**Request Performance:**
- `api_requests_total`: Total API requests served
- `api_requests_per_second`: Current request throughput
- `api_response_time_ms`: Average response time
- `api_errors_total`: Number of API errors
- `api_cache_hits`: Number of cache hits
- `api_cache_misses`: Number of cache misses

**Endpoint Metrics:**
- `status_endpoint_calls`: Calls to status endpoints
- `index_endpoint_calls`: Calls to index query endpoints
- `monitor_endpoint_calls`: Calls to monitor endpoints
- `admin_endpoint_calls`: Calls to admin endpoints

#### IPFS Service Metrics

**Network Performance:**
- `ipfs_peers_connected`: Number of connected IPFS peers
- `ipfs_data_uploaded_bytes`: Total data uploaded to IPFS
- `ipfs_data_downloaded_bytes`: Total data downloaded from IPFS
- `ipfs_pin_operations`: Number of pin operations performed
- `ipfs_chunks_shared`: Number of index chunks shared

**Synchronization Metrics:**
- `ipfs_sync_operations`: Number of sync operations
- `ipfs_sync_latency_ms`: Average sync operation time
- `ipfs_failed_retrievals`: Number of failed chunk retrievals

### Accessing Service Metrics

#### REST API Access

Metrics are available through the Control Service API:

```bash
# Get metrics for all services
curl http://localhost:8080/api/v1/metrics

# Get metrics for specific service
curl http://localhost:8080/api/v1/services/scraper/metrics

# Get detailed metrics with verbose output
curl http://localhost:8080/api/v1/services/scraper?verbose=true&include=metrics

# Get metrics in different formats
curl http://localhost:8080/api/v1/metrics?format=json
curl http://localhost:8080/api/v1/metrics?format=prometheus
```

#### CLI Access

```bash
# Show basic service status with key metrics
khedra status --metrics

# Show detailed metrics for all services
khedra metrics

# Show metrics for specific service
khedra metrics --service=scraper

# Export metrics to file
khedra metrics --output=/path/to/metrics.json

# Watch metrics in real-time
khedra metrics --watch --interval=5s
```

#### Programmatic Access

```go
// Example: Getting service metrics programmatically
import "github.com/TrueBlocks/trueblocks-khedra/v5/pkg/client"

client := client.NewKhedraClient("http://localhost:8080")
metrics, err := client.GetServiceMetrics("scraper")
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Blocks per second: %f\n", metrics["blocks_per_second"])
fmt.Printf("Memory usage: %d bytes\n", metrics["memory_usage_bytes"])
```

### Interpreting Metrics

#### Performance Health Indicators

**Scraper Service Health:**
- **Healthy**: `blocks_per_second > 10`, `rpc_response_time_ms < 500`, `memory_usage_bytes` stable
- **Warning**: `blocks_per_second < 5`, `rpc_response_time_ms > 1000`, `blocks_behind > 1000`
- **Critical**: `blocks_per_second < 1`, `rpc_requests_failed > 10%`, `memory_usage_bytes` increasing rapidly

**Monitor Service Health:**
- **Healthy**: `detection_latency_ms < 30000`, `false_positives < 5%`, all monitored addresses active
- **Warning**: `detection_latency_ms > 60000`, `false_positives > 10%`, some addresses not responding
- **Critical**: `detection_latency_ms > 300000`, `false_positives > 25%`, monitoring completely behind

**API Service Health:**
- **Healthy**: `api_response_time_ms < 100`, `api_errors_total < 1%`, `api_cache_hits > 80%`
- **Warning**: `api_response_time_ms > 500`, `api_errors_total > 5%`, `api_cache_hits < 60%`
- **Critical**: `api_response_time_ms > 2000`, `api_errors_total > 15%`, service unresponsive

#### Resource Utilization Thresholds

**Memory Usage:**
- **Normal**: < 2GB per service
- **High**: 2-4GB per service (monitor for leaks)
- **Critical**: > 4GB per service (immediate attention required)

**CPU Usage:**
- **Normal**: < 50% average
- **High**: 50-80% average (acceptable under load)
- **Critical**: > 80% sustained (performance degradation likely)

### Metrics-Based Troubleshooting

#### High Resource Usage

**High Memory Usage:**
```bash
# Check memory metrics
curl http://localhost:8080/api/v1/services/scraper/metrics | jq '.memory_usage_bytes'

# If memory usage is high:
# 1. Reduce batch size
# 2. Increase sleep interval
# 3. Check for memory leaks in logs
# 4. Restart service if memory continues growing
```

**High CPU Usage:**
```bash
# Check CPU metrics and goroutine count
curl http://localhost:8080/api/v1/metrics | jq '.cpu_usage_percent, .goroutines_count'

# If CPU usage is high:
# 1. Reduce batch size
# 2. Increase sleep interval
# 3. Check for infinite loops in logs
# 4. Verify RPC endpoint performance
```

#### Performance Degradation

**Slow Indexing:**
```bash
# Check indexing performance
curl http://localhost:8080/api/v1/services/scraper/metrics | jq '.blocks_per_second, .rpc_response_time_ms'

# Troubleshooting steps:
# 1. Check RPC response times
# 2. Verify network connectivity
# 3. Adjust batch size based on performance
# 4. Check for rate limiting
```

**API Response Delays:**
```bash
# Check API performance
curl http://localhost:8080/api/v1/services/api/metrics | jq '.api_response_time_ms, .api_cache_hits'

# Troubleshooting steps:
# 1. Check cache hit ratio
# 2. Verify index integrity
# 3. Monitor concurrent request load
# 4. Check for slow database queries
```

#### Service Failures

**RPC Connection Issues:**
```bash
# Check RPC health metrics
curl http://localhost:8080/api/v1/services/scraper/metrics | jq '.rpc_requests_failed, .rpc_rate_limit_hits'

# Troubleshooting steps:
# 1. Test RPC endpoints directly
# 2. Increase sleep intervals if rate limited
# 3. Switch to backup RPC endpoints
# 4. Check network connectivity
```

### Alerting and Monitoring Setup

#### Prometheus Integration

Khedra can export metrics in Prometheus format for integration with monitoring systems:

```yaml
# prometheus.yml configuration
scrape_configs:
  - job_name: 'khedra'
    static_configs:
      - targets: ['localhost:8080']
    metrics_path: '/api/v1/metrics'
    params:
      format: ['prometheus']
    scrape_interval: 30s
```

#### Grafana Dashboard

Key metrics to monitor in Grafana:

**Performance Dashboard:**
- Blocks per second (Scraper)
- API response times
- Memory and CPU usage
- RPC response times

**Health Dashboard:**
- Service uptime
- Error rates
- Detection latency
- System resource utilization

#### Alerting Rules

Example alerting rules for common issues:

```yaml
# Slow indexing alert
- alert: SlowIndexing
  expr: khedra_blocks_per_second < 5
  for: 5m
  labels:
    severity: warning
  annotations:
    summary: "Khedra indexing is slow"
    description: "Indexing rate is {{ $value }} blocks/sec, below threshold"

# High memory usage alert
- alert: HighMemoryUsage
  expr: khedra_memory_usage_bytes > 4000000000
  for: 2m
  labels:
    severity: critical
  annotations:
    summary: "High memory usage detected"
    description: "Memory usage is {{ $value }} bytes, above 4GB threshold"

# API response time alert
- alert: SlowAPI
  expr: khedra_api_response_time_ms > 1000
  for: 3m
  labels:
    severity: warning
  annotations:
    summary: "API responses are slow"
    description: "Average response time is {{ $value }}ms"
```

#### Custom Monitoring Scripts

```bash
#!/bin/bash
# Simple monitoring script
METRICS_URL="http://localhost:8080/api/v1/metrics"

# Check blocks per second
BPS=$(curl -s $METRICS_URL | jq -r '.scraper.blocks_per_second // 0')
if (( $(echo "$BPS < 5" | bc -l) )); then
    echo "WARNING: Slow indexing detected: $BPS blocks/sec"
fi

# Check memory usage
MEMORY=$(curl -s $METRICS_URL | jq -r '.scraper.memory_usage_bytes // 0')
if (( MEMORY > 4000000000 )); then
    echo "CRITICAL: High memory usage: $((MEMORY/1024/1024))MB"
fi

# Check API health
API_TIME=$(curl -s $METRICS_URL | jq -r '.api.api_response_time_ms // 0')
if (( $(echo "$API_TIME > 1000" | bc -l) )); then
    echo "WARNING: Slow API responses: ${API_TIME}ms"
fi
```

### Best Practices for Metrics Monitoring

#### Regular Monitoring

1. **Establish Baselines**: Monitor metrics during normal operation to establish performance baselines
2. **Set Appropriate Thresholds**: Configure alerts based on your specific environment and requirements
3. **Monitor Trends**: Look for gradual degradation over time, not just immediate issues
4. **Correlate Metrics**: Use multiple metrics together to diagnose issues accurately

#### Performance Optimization

1. **Use Metrics for Tuning**: Adjust batch sizes and sleep intervals based on actual performance metrics
2. **Monitor Resource Efficiency**: Track resource usage to optimize system utilization
3. **Identify Bottlenecks**: Use metrics to identify which component is limiting performance
4. **Validate Changes**: Use metrics to verify that configuration changes improve performance

#### Operational Excellence

1. **Automate Monitoring**: Set up automated alerts for critical metrics
2. **Create Dashboards**: Visualize key metrics for easier monitoring
3. **Document Thresholds**: Maintain documentation of what constitutes healthy vs. problematic metrics
4. **Regular Reviews**: Periodically review and adjust monitoring thresholds based on operational experience
