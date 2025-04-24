# System Architecture

## Architectural Overview

Khedra employs a modular, service-oriented architecture designed for flexibility, resilience, and extensibility. The system is structured around a central application core that coordinates multiple specialized services, each with distinct responsibilities.

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

The main application container that initializes, configures, and manages the lifecycle of all services. It provides:

- Service registration and coordination
- Application startup and shutdown sequences
- Signal handling for graceful termination
- Global state management
- Cross-service communication

Implementation: `app/khedra.go`

### 2. Service Framework

Khedra implements five primary services:

#### 2.1 Control Service

- Provides management endpoints for other services
- Handles service health monitoring
- Enables runtime reconfiguration
- Serves as the primary management interface

Implementation: `pkg/services/control/service.go`

#### 2.2 Scraper Service

- Processes blockchain data to build the Unchained Index
- Extracts address appearances from transactions, logs, and traces
- Manages indexing state and progress tracking
- Handles retry logic for failed operations
- Implements batch processing with configurable parameters

Implementation: `pkg/services/scraper/service.go`

#### 2.3 Monitor Service

- Tracks specified addresses for on-chain activity
- Maintains focused indices for monitored addresses
- Processes real-time blocks for quick notifications
- Supports flexible notification configurations
- Manages monitor definitions and states

Implementation: `pkg/services/monitor/service.go`

#### 2.4 API Service

- Exposes RESTful endpoints for data access
- Implements query interfaces for the index and monitors
- Handles authentication and rate limiting
- Provides structured data responses in multiple formats
- Includes Swagger documentation for API endpoints

Implementation: `pkg/services/api/service.go`

#### 2.5 IPFS Service

- Manages distributed sharing of index data
- Handles chunking of index data for efficient distribution
- Implements publishing and retrieval mechanisms
- Provides peer discovery and connection management
- Integrates with the IPFS network protocol

Implementation: `pkg/services/ipfs/service.go`

### 3. Configuration Manager

A centralized system for managing application settings, including:

- Configuration file parsing and validation
- Environment variable integration
- Runtime configuration updates
- Defaults management
- Chain-specific configuration handling

Implementation: `pkg/config/config.go`

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

The interface layer between Khedra and blockchain nodes:

- RPC client implementations
- Connection pooling and management
- Request rate limiting and backoff strategies
- Error handling and resilience patterns
- Chain-specific adaptations

Implementation: `pkg/rpc/client.go`

## Communication Patterns

Khedra employs several communication patterns between components:

1. **Service-to-Service Communication**: Structured message passing between services using channels
2. **RPC Communication**: JSON-RPC communication with blockchain nodes
3. **REST API**: HTTP-based communication for external interfaces
4. **File-Based Storage**: Persistent data storage using structured files

## Deployment Architecture

Khedra supports multiple deployment models:

1. **Standalone Application**: Single-process deployment for individual users
2. **Docker Container**: Containerized deployment for managed environments
3. **Distributed Deployment**: Multiple instances sharing index data via IPFS

## Security Architecture

Security considerations in Khedra's architecture include:

1. **Local-First Processing**: Minimizes exposure of query data
2. **API Authentication**: Optional key-based authentication for API access
3. **Configuration Protection**: Secure handling of RPC credentials
4. **Update Verification**: Integrity checks for application updates
5. **Resource Isolation**: Service-level resource constraints

The modular design of Khedra allows for individual components to be extended, replaced, or enhanced without affecting the entire system, providing a solid foundation for future development and integration.
