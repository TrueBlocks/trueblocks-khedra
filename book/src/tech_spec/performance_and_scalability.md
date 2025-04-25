# Performance and Scalability

This section details Khedra's performance characteristics, optimization strategies, and scalability considerations.

## Performance Benchmarks

### Indexing Performance

Typical indexing performance metrics on reference hardware (8-core CPU, 16GB RAM, SSD storage):

| Chain     | Block Processing Rate | Trace Processing Rate | Disk Space per 1M Blocks |
| --------- | --------------------- | --------------------- | ------------------------ |
| Mainnet   | 15-25 blocks/sec      | 5-15 blocks/sec       | 1.5-2.5 GB               |
| Testnets  | 20-30 blocks/sec      | 8-20 blocks/sec       | 0.5-1.5 GB               |
| L2 Chains | 30-50 blocks/sec      | 10-25 blocks/sec      | 1.0-2.0 GB               |

Factors affecting indexing performance:

1. **RPC Endpoint Performance**: Quality and latency of the blockchain RPC connection
2. **Trace Availability**: Whether traces are locally available or require RPC calls
3. **Block Complexity**: Number of transactions and traces in each block
4. **Hardware Specifications**: CPU cores, memory, and disk I/O capacity
5. **Network Conditions**: Bandwidth and latency for remote RPC endpoints

### Query Performance

Performance metrics for common queries:

| Query Type                | Response Time (cold) | Response Time (warm) |
| ------------------------- | -------------------- | -------------------- |
| Address Appearance Lookup | 50-200ms             | 10-50ms              |
| Block Range Scan          | 100-500ms            | 20-100ms             |
| Monitor Status Check      | 20-100ms             | 5-20ms               |
| API Status Endpoints      | 5-20ms               | 1-5ms                |

Factors affecting query performance:

1. **Index Structure**: Organization and optimization of the Unchained Index
2. **Memory Cache**: Availability of data in memory versus disk access
3. **Query Complexity**: Number of addresses and block range size
4. **Hardware Specifications**: Particularly memory and disk speed
5. **Concurrent Load**: Number of simultaneous queries being processed

## Performance Optimization Strategies

### Memory Management

Khedra implements several memory optimization techniques:

1. **Bloom Filters**: Space-efficient probabilistic data structures to quickly determine if an address might appear in a block
2. **LRU Caching**: Least Recently Used caching for frequently accessed data
3. **Memory Pooling**: Reuse of allocated memory for similar operations
4. **Batch Processing**: Processing multiple items in batches to amortize overhead
5. **Incremental GC**: Tuned garbage collection to minimize pause times

Implementation example:

```go
// Bloom filter implementation for quick address lookups
type AppearanceBloomFilter struct {
    filter     *bloom.BloomFilter
    capacity   uint
    errorRate  float64
}

func NewAppearanceBloomFilter(expectedItems uint) *AppearanceBloomFilter {
    return &AppearanceBloomFilter{
        filter:    bloom.NewWithEstimates(uint(expectedItems), 0.01),
        capacity:  expectedItems,
        errorRate: 0.01,
    }
}

func (bf *AppearanceBloomFilter) Add(address []byte) {
    bf.filter.Add(address)
}

func (bf *AppearanceBloomFilter) MayContain(address []byte) bool {
    return bf.filter.Test(address)
}
```

### Disk I/O Optimization

Strategies for optimizing disk operations:

1. **Sequential Writes**: Organize write patterns for sequential access where possible
2. **Write Batching**: Combine multiple small writes into larger operations
3. **Read-Ahead Buffering**: Anticipate and pre-load data likely to be needed
4. **Cache Warming**: Proactively load frequently accessed data into memory
5. **Compression**: Reduce storage requirements and I/O bandwidth

Example implementation:

```go
// Batched write implementation
type BatchedWriter struct {
    buffer     []byte
    maxSize    int
    flushThreshold int
    target     io.Writer
    mutex      sync.Mutex
}

func (w *BatchedWriter) Write(p []byte) (n int, err error) {
    w.mutex.Lock()
    defer w.mutex.Unlock()
    
    // Add to buffer
    w.buffer = append(w.buffer, p...)
    
    // Flush if threshold reached
    if len(w.buffer) >= w.flushThreshold {
        return w.Flush()
    }
    
    return len(p), nil
}

func (w *BatchedWriter) Flush() (n int, err error) {
    if len(w.buffer) == 0 {
        return 0, nil
    }
    
    n, err = w.target.Write(w.buffer)
    w.buffer = w.buffer[:0] // Clear buffer
    return n, err
}
```

### Concurrency Management

Techniques for efficient parallel processing:

1. **Worker Pools**: Fixed-size pools of worker goroutines for controlled parallelism
2. **Pipeline Processing**: Multi-stage processing with each stage running concurrently
3. **Batched Distribution**: Group work items for efficient parallelization
4. **Backpressure Mechanisms**: Prevent resource exhaustion during high load
5. **Adaptive Parallelism**: Adjust concurrency based on system load and resources

Example worker pool implementation:

```go
// Worker pool for parallel block processing
type BlockWorkerPool struct {
    workers     int
    queue       chan BlockTask
    results     chan BlockResult
    wg          sync.WaitGroup
    ctx         context.Context
    cancel      context.CancelFunc
}

func NewBlockWorkerPool(workers int) *BlockWorkerPool {
    ctx, cancel := context.WithCancel(context.Background())
    pool := &BlockWorkerPool{
        workers: workers,
        queue:   make(chan BlockTask, workers*2),
        results: make(chan BlockResult, workers*2),
        ctx:     ctx,
        cancel:  cancel,
    }
    
    // Start worker goroutines
    pool.wg.Add(workers)
    for i = 0; i < workers; i++ {
        go pool.worker(i)
    }
    
    return pool
}

func (p *BlockWorkerPool) worker(id int) {
    defer p.wg.Done()
    
    for {
        select {
        case <-p.ctx.Done():
            return
        case task, ok := <-p.queue:
            if !ok {
                return
            }
            
            result := processBlock(task)
            p.results <- result
        }
    }
}
```

## Scalability Considerations

### Vertical Scaling

Khedra is designed to efficiently utilize additional resources when available:

1. **CPU Utilization**: Automatic adjustment to use available CPU cores
2. **Memory Utilization**: Configurable memory limits for caching and processing
3. **Storage Scaling**: Support for high-performance storage devices
4. **I/O Optimization**: Tuning based on available I/O capacity

Configuration parameters for vertical scaling:

```yaml
services:
  scraper:
    concurrency: 8         # Number of parallel workers
    memory_limit: "4GB"    # Maximum memory usage
    batch_size: 1000       # Items per processing batch
```

### Horizontal Scaling

While Khedra runs as a single process, it supports distributed operation through:

1. **Multi-Instance Deployment**: Running multiple instances focusing on different chains
2. **Shared Index via IPFS**: Collaborative building and sharing of the index
3. **Split Processing**: Dividing block ranges between instances
4. **API Load Balancing**: Distributing API queries across instances

Example multi-instance deployment:

```
Instance 1: Mainnet indexing
Instance 2: L2 chains indexing
Instance 3: API service
Instance 4: Monitor service
```

### Data Volume Management

Strategies for handling large data volumes:

1. **Selective Indexing**: Configure which data types to index (transactions, logs, traces)
2. **Retention Policies**: Automatically prune older cache data while preserving the index
3. **Compression**: Reduce storage requirements through data compression
4. **Tiered Storage**: Move less frequently accessed data to lower-cost storage

Example retention configuration:

```yaml
cache:
  retention:
    blocks: "30d"      # Keep block data for 30 days
    traces: "15d"      # Keep trace data for 15 days
    receipts: "60d"    # Keep receipt data for 60 days
  compression: true    # Enable data compression
```

## Performance Monitoring

Khedra includes built-in performance monitoring capabilities:

1. **Metrics Collection**: Runtime statistics for key operations
2. **Performance Logging**: Timing information for critical paths
3. **Resource Monitoring**: Tracking of CPU, memory, and disk usage
4. **Bottleneck Detection**: Identification of performance limitations

Example metrics available:

```json
{
  "scraper": {
    "blocks_processed": 1520489,
    "blocks_per_second": 18.5,
    "current_block": 18245367,
    "last_processed": "2023-06-15T14:23:45Z",
    "memory_usage_mb": 2458,
    "rpc_calls": 247896,
    "processing_latency_ms": 54
  }
}
```

## Resource Requirements

Recommended system specifications based on usage patterns:

### Minimum Requirements

- **CPU**: 4 cores
- **Memory**: 8GB RAM
- **Storage**: 250GB SSD
- **Network**: 10Mbps stable connection
- **Supported Workload**: Monitoring a few addresses, single chain, limited API usage

### Recommended Configuration

- **CPU**: 8 cores
- **Memory**: 16GB RAM
- **Storage**: 1TB NVMe SSD
- **Network**: 100Mbps stable connection
- **Supported Workload**: Full indexing of mainnet, multiple monitored addresses, moderate API usage

### High-Performance Configuration

- **CPU**: 16+ cores
- **Memory**: 32GB+ RAM
- **Storage**: 2TB+ NVMe SSD with high IOPS
- **Network**: 1Gbps+ connection
- **Supported Workload**: Multiple chains, extensive monitoring, heavy API usage, IPFS participation

These performance optimizations and scalability considerations enable Khedra to handle the demands of blockchain data processing efficiently across a wide range of hardware configurations and usage scenarios.
