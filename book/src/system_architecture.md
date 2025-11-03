# System Architecture

## Architectural Overview

Khedra employs a modular, service-oriented architecture. A central application core wires up a set of specialized services, each with a narrow responsibility.

### High-Level Architecture Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                      Khedra Application                          │
├─────────┬─────────┬─────────┬─────────┬─────────────────────────┤
│ Control │ Scraper │ Monitor │   API   │         IPFS            │
│ Service │ Service │ Service │ Service │        Service          │
├─────────┴─────────┴─────────┴─────────┴─────────────────────────┤
│                      Configuration Manager                       │
├─────────────────────────────────────────────────────────────────┤
│                          Data Layer                              │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────────────┐ │
│  │ Unchained│  │  Binary  │  │ Monitor  │  │ Chain-Specific   │ │
│  │   Index  │  │  Caches  │  │   Data   │  │     Data         │ │
│  └──────────┘  └──────────┘  └──────────┘  └──────────────────┘ │
├─────────────────────────────────────────────────────────────────┤
│                      Blockchain Connectors                       │
└─────────────────────────────────────────────────────────────────┘
             ▲                    ▲                     ▲
             │                    │                     │
 ┌───────────┴──────────┐ ┌──────┴───────┐  ┌──────────┴──────────┐
 │  Ethereum Mainnet    │ │   Testnets   │  │   Other EVM Chains  │
 └──────────────────────┘ └──────────────┘  └─────────────────────┘
```

## Core Components

### 1. Khedra Application

The main application container initializes configuration and wires up the enabled services. Today it provides:

- Basic service instantiation (no dynamic registration at runtime)
- One–time startup (no hot restart orchestration)
- OS signal handling for shutdown via the underlying service manager

There is no cross‑service message bus, restart policy, or runtime dependency graph.

Implementation: `app/app.go`, `app/action_daemon.go`

### 2. Service Framework

Khedra implements five primary services:

#### 2.1 Control Service

Current (implemented) responsibilities:

- Exposes a minimal HTTP interface for pausing / unpausing pausable services
- Reports simple paused / running / not‑pausable status via `/isPaused`
- Listens on the first available port in the range 8338–8335

Not implemented: start/stop/restart of individual services, runtime configuration mutation, health or metrics aggregation, dependency ordering logic, automatic restarts, or a generalized management API surface.

Implementation entry: constructed in `app/action_daemon.go` (via `services.NewControlService`).

#### 2.2 Scraper Service

Intended role (some functionality provided by the upstream SDK library):

- Processes blockchain data in batches (batch size & sleep interval configurable in `config.yaml`)
- Capable of being paused / unpaused through the Control Service endpoints

Detailed index storage, retry semantics, and appearance extraction logic live in the shared TrueBlocks SDK (not in this repository) and are therefore abstracted from this codebase. Paths: created through `services.NewScrapeService` in `app/action_daemon.go`.

#### 2.3 Monitor Service

Current state:

- Instantiated when enabled but disabled by default
- Supports pause / unpause
- Advanced notification / registry features are not implemented here.

Implementation entry: created via `services.NewMonitorService` in `app/action_daemon.go`.

#### 2.4 API Service

When enabled it exposes HTTP endpoints (details provided by the SDK). This repository does not implement authentication, rate limiting, Swagger generation, or multi‑format response logic.

Implementation entry: `services.NewApiService` in `app/action_daemon.go`.

#### 2.5 IPFS Service

Optional. Created only if enabled. Within this codebase we only instantiate via `services.NewIpfsService`.

### 3. Configuration Manager

Implemented as a YAML backed configuration (`~/.khedra/config.yaml` by default) created / edited through the init wizard or `khedra config edit`. Runtime (hot) reconfiguration is **not** supported; changes require a daemon restart.

Implementation: `pkg/types/config.go` and related helpers in `app/`.

### 4. Data Layer

The persistent storage infrastructure for Khedra:

#### 4.1 Unchained Index

- Core data structure mapping addresses to appearances
- Optimized for fast lookups and efficient storage
- Implements chunking for distributed sharing
- Includes versioning for format compatibility

Implementation: `pkg/index/index.go`

#### 4.2 Binary Caches

- Stores raw blockchain data for efficient reprocessing
- Implements cache invalidation and management
- Optimizes storage space usage with compression
- Supports pruning and maintenance operations

Implementation: `pkg/cache/cache.go`

#### 4.3 Monitor Data

- Stores monitor definitions and state
- Tracks monitored address appearances
- Maintains notification history
- Implements efficient storage for frequent updates

Implementation: `pkg/monitor/data.go`

#### 4.4 Chain-Specific Data

- Segregates data by blockchain
- Stores chain metadata and state
- Manages chain-specific configurations
- Handles chain reorganizations

Implementation: `pkg/chains/data.go`

### 5. Blockchain Connectors

Low‑level RPC client logic is handled in upstream TrueBlocks components; this repository primarily validates configured RPC endpoints (see `HasValidRpc` usage in `app/action_daemon.go`).

## Communication Patterns

Khedra employs several communication patterns between components:

1. **RPC Communication**: JSON-RPC communication with blockchain nodes (through upstream SDK)
2. **Minimal Control HTTP**: `/isPaused`, `/pause`, `/unpause` endpoints for operational control
3. **File-Based Storage**: Index / cache paths determined by config (actual index logic external)

## Deployment Architecture

Khedra supports multiple deployment models:

1. **Standalone Application**: Single-process deployment
2. **(Removed)** Prior Docker support has been removed (see project README)

## Security Notes (Current Scope)

Current implementation is local‑first and depends on the operator to secure the host machine. Features such as authenticated API access, update integrity verification, and formal resource isolation are not implemented in this repository.

---

