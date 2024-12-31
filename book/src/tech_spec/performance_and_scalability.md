# Performance and Scalability

## Performance Benchmarks

Khedra is designed to handle high-throughput blockchain data. Typical performance benchmarks include:

- Processing speed: ~500 blocks per second (depending on RPC response time).
- REST API response time: <50ms for standard queries.

## Strategies for Handling Large-Scale Data

1. Use high-performance RPC endpoints with low latency.
2. Increase local storage capacity to handle large blockchain data.
3. Scale horizontally by running multiple instances of Khedra for different chains.

## Resource Optimization Guidelines

- Limit the number of chains processed simultaneously to reduce system load.
- Configure `--sleep` duration to balance processing speed with system resource usage.
