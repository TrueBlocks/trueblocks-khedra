# Performance and Scalability

This chapter documents the tunable parameters that exist in the current codebase.

## 1. Tunable Parameters

The runtime tuning levers exposed in `config.yaml`:

- `services.scraper.sleep` (seconds between batch cycles)
- `services.scraper.batchSize` (approximate blocks per batch)
- `services.monitor.sleep` (monitor disabled by default)
- `services.monitor.batchSize`
- `enabled` flags per service

API & IPFS have `port` values; scraper/monitor do not expose listeners here.

## 2. Practical Tuning Guidance

Start with defaults (sleep ≈ 10–12s, batchSize 500). Adjust slowly:

- High idle CPU & stable RAM → increase batchSize moderately (e.g. 750 → 1000)
- RPC errors / rate limits → decrease batchSize or increase sleep
- Memory pressure → decrease batchSize first; then increase sleep if needed
- Slow catch-up → cautiously increase batchSize (ensure RPC can handle)

### Heuristic Table

| Situation                       | Action                              |
|---------------------------------|-------------------------------------|
| Idle resources                  | Increase batchSize                  |
| RPC timeouts / 429s             | Decrease batchSize or add sleep     |
| Rising memory usage             | Decrease batchSize                  |
| Far behind chain tip            | Increase batchSize if safe          |

## 3. Measuring

Use OS tools:
- CPU / Mem: `top`, `ps`
- Disk growth: `du -sh ~/.khedra`
- RPC latency: ad hoc `curl` against `eth_blockNumber`

Use external observation tools; there is no embedded metrics endpoint.

## 4. Scaling Modes

Single-process only. You may run separate processes per chain (with isolated config directories) if desired. Keep combined RPC load in mind.

## 5. Monitor Service

Monitor is disabled by default and supports pause/unpause.

## 6. Summary

Tuning today = batchSize + sleep + enable/disable services. Measure effects externally. Accuracy now takes precedence over aspirational detail.

_Concise by design._

