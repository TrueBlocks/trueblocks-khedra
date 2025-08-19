# Introduction

## What is Khedra?

Khedra (pronounced *kɛd-ɾɑ*) is TrueBlocks' service management system that provides blockchain indexing, monitoring, and data serving capabilities for EVM-compatible blockchains. It creates and maintains the Unchained Index - a comprehensive, permissionless index of address appearances across blockchain data.

At its core, Khedra creates and maintains the Unchained Index - a permissionless index of address appearances across blockchain data, including transactions, event logs, execution traces, and more. This detailed indexing enables powerful monitoring capabilities for any address on any supported chain.

## Key Features

### 1. Comprehensive Indexing

The Scraper service indexes address appearances from multiple sources:
- Transaction senders and recipients
- Event log topics and data fields  
- Internal calls from execution traces
- Block rewards and consensus activities

This detailed indexing enables fast lookups of any address's complete on-chain history.

### 2. Address Monitoring

The Monitor service provides real-time tracking of specific addresses:
- Detects new transactions involving monitored addresses
- Captures relevant events and interactions
- Supports monitoring multiple addresses simultaneously

### 3. Service Management

Khedra operates through five core services with runtime control:

- **Scraper**: Builds and maintains the Unchained Index *(pausable)*
- **Monitor**: Tracks specific addresses *(pausable)*
- **API**: Provides REST endpoints for data access
- **Control**: HTTP interface for service management
- **IPFS**: Distributed data sharing (optional)

Services marked as *pausable* can be stopped and resumed without restarting the entire system.

### 4. Multi-Chain Support

While Ethereum mainnet is the primary focus, Khedra supports any EVM-compatible blockchain:
- Test networks (Sepolia, Goerli)
- Layer 2 solutions (Optimism, Arbitrum, Polygon)
- Alternative EVMs (Gnosis Chain, Base)

Each chain requires only a valid RPC endpoint to begin indexing.

### 5. Privacy-Preserving Design

Khedra runs entirely on your local machine:
- No data sent to third-party servers
- Complete control over your queries and data
- Local-first architecture for maximum privacy

## Architecture

### Service Communication

Services communicate through:
- Shared configuration system
- HTTP APIs for control operations
- Local file system for data storage
- Optional IPFS for distributed sharing

### Runtime Control

The Control service provides HTTP endpoints for:
- Pausing/unpausing indexing and monitoring
- Checking service status
- Managing service lifecycle

This enables automation and integration with other systems.

## Use Cases

Khedra excels in various blockchain data scenarios:

- **Account Monitoring**: Track transactions and interactions for specific addresses
- **Index Building**: Create comprehensive local blockchain indices
- **Data Analysis**: Extract on-chain patterns and insights  
- **Custom Applications**: Build specialized tools using the REST API
- **Research**: Analyze blockchain data with complete privacy

## Getting Started

The following sections will guide you through:

1. Installing and configuring Khedra
2. Understanding service management
3. Using pause/unpause functionality
4. Working with the REST API
5. Maintenance and troubleshooting

Khedra provides the foundation for building powerful blockchain data applications while maintaining complete control over your data and privacy.
