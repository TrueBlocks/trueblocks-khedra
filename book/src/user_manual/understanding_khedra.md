# Understanding Khedra

## Core Concepts

### The Unchained Index

The foundation of Khedra is the Unchained Index - a specialized data structure that maps blockchain addresses to their appearances in blockchain data. Think of it as a reverse index: while a blockchain explorer lets you look up a transaction and see which addresses were involved, the Unchained Index lets you look up an address and see all transactions where it appears.

The index captures appearances from multiple sources:

- **External Transactions**: Direct sends and receives
- **Internal Transactions**: Contract-to-contract calls (from traces)
- **Event Logs**: Events emitted by smart contracts
- **State Changes**: Modifications to contract storage
- **Special Appearances**: Block rewards, validators, etc.

What makes this particularly powerful is that the index includes trace-derived appearances - meaning it captures internal contract interactions that normal blockchain explorers miss.

### Address Appearances

An "appearance" in Khedra means any instance where an address is referenced in blockchain data. Each appearance record contains:

- The address that appeared
- The block number where it appeared
- The transaction index within that block
- Additional metadata about the appearance type

These compact records allow Khedra to quickly answer the fundamental question: "Where does this address appear in the blockchain?"

### Local-First Architecture

Khedra operates as a "local-first" application, meaning:

1. All data processing happens on your local machine
2. Your queries never leave your computer 
3. You maintain complete ownership of your data
4. The application continues to work without internet access

This approach maximizes privacy and resilience while minimizing dependency on external services.

### Distributed Collaboration

While Khedra is local-first, it also embraces distributed collaboration through IPFS integration:

- The Unchained Index can be shared and downloaded in chunks
- Users can contribute to different parts of the index
- New users can bootstrap quickly by downloading existing index portions
- The system becomes more resilient as more people participate

This creates a hybrid model that preserves privacy while enabling community benefits.

## System Architecture

### Service Components

Khedra is organized into five core services:

1. **Control Service**
   - Central management interface
   - Exposes API endpoints for service control
   - Handles configuration and coordinating other services

2. **Scraper Service**
   - Processes blockchain data to build the Unchained Index
   - Extracts address appearances from blocks, transactions, and traces
   - Works in configurable batches with adjustable sleep intervals

3. **Monitor Service**
   - Tracks specific addresses of interest
   - Provides notifications for address activities
   - Maintains focused indices for monitored addresses

4. **API Service**
   - Exposes data through REST endpoints (defined here: [API Docs](https://trueblocks.io/api/))
   - Provides query interfaces for the index and monitors
   - Enables integration with other tools and services

5. **IPFS Service**
   - Facilitates distributed sharing of index data
   - Handles publishing and retrieving chunks via IPFS
   - Enables collaborative index building

### Data Flow

Here's how data flows through the Khedra system:

1. The Scraper retrieves blockchain data from configured RPC endpoints
2. Address appearances are extracted and added to the Unchained Index
3. The Monitor service checks new blocks for appearances of watched addresses
4. The API service provides query access to the indexed data
5. Optionally, index chunks are shared via the IPFS service

### Directory Structure

Khedra organizes its data with this structure:

```bash
~/.khedra/
├── config.yaml       # Main configuration file
├── data/             # Main data directory
│   ├── mainnet/      # Chain-specific data
│   │   ├── cache/    # Binary caches
│   │   ├── monitors/ # Address monitor data
│   │   └── index/    # Unchained Index chunks
│   └── [other-chains]/
└── logs/             # Application logs
```

The above structure may vary depending on your version and configuration. Each chain has its own subdirectory, allowing Khedra to manage multiple chains simultaneously.

## Terminology

To help navigate Khedra effectively, here are key terms you'll encounter:

- **Appearance**: Any reference to an address in blockchain data
- **Chunk**: A portion of the Unchained Index covering a range of blocks
- **Finalized**: Blocks that have reached consensus and won't be reorganized
- **Monitor**: A configuration to track specific addresses of interest
- **RPC**: Remote Procedure Call - the method for communicating with blockchain nodes
- **Trace**: Detailed execution record of a transaction, including internal calls

Understanding these core concepts provides the foundation for effectively using Khedra's capabilities, which we'll explore in the following chapters.

## Implementation Details

The core concepts and system architecture described in this chapter are implemented in the following Go files:

### The Unchained Index

The Unchained Index implementation is handled primarily by the TrueBlocks-core library, with Khedra providing the service framework. The primary code files for index interactions are:

- **Index Management**: The `scraper` service, implemented in the service framework initialized in [`app/action_daemon.go`](/Users/jrush/Development/trueblocks-core/khedra/app/action_daemon.go)

### Address Monitoring

The monitoring system for tracking address appearances is implemented in:

- **Monitor Service**: Initialized in [`app/action_daemon.go`](/Users/jrush/Development/trueblocks-core/khedra/app/action_daemon.go) with references to the `monitor` service
- **Monitor Configuration**: Service settings defined in [`pkg/types/service.go`](/Users/jrush/Development/trueblocks-core/khedra/pkg/types/service.go)
- **Monitor Options**: Monitor-specific options defined and processed in [`app/action_daemon.go`](/Users/jrush/Development/trueblocks-core/khedra/app/action_daemon.go) with the `MonitorsOptions` struct

### Service Components

The five core services are defined and initialized in these files:

- **Service Definitions**: [`pkg/types/service.go`](/Users/jrush/Development/trueblocks-core/khedra/pkg/types/service.go) defines the `Service` struct and validation rules
- **Service Initialization**: [`app/action_daemon.go`](/Users/jrush/Development/trueblocks-core/khedra/app/action_daemon.go) in the `daemonAction` function initializes each service based on configuration
- **Service Manager**: The `ServiceManager` is created in the `daemonAction` function to coordinate all services

### Directory Structure

The directory structure described in this chapter is established by:

- **Folder Initialization**: The `initializeFolders` function in [`app/config.go`](/Users/jrush/Development/trueblocks-core/khedra/app/config.go)
- **Path Resolution**: Path management and expansion functions throughout the codebase handle the directory structure
