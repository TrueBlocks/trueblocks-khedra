# Introduction

## What is Khedra?

Khedra (pronounced *kɛd-ɾɑ*) is an all-in-one blockchain indexing and monitoring solution for EVM-compatible blockchains. It provides a comprehensive suite of tools to index, monitor, serve, and share blockchain data in a local-first, privacy-preserving manner.

At its core, Khedra creates and maintains the Unchained Index - a permissionless index of address appearances across blockchain data, including transactions, event logs, execution traces, and more. This detailed indexing enables powerful monitoring capabilities for any address on any supported chain.

## Key Features

### 1. Comprehensive Indexing

Khedra indexes address appearances from multiple sources:
- Transactions (senders and recipients)
- Event logs (topics and data fields)
- Execution traces (internal calls)
- Smart contract state changes
- Block rewards and staking activities
- Genesis allocations

The resulting index allows for lightning-fast lookups of any address's complete on-chain history.

### 2. Multi-Chain Support

While Ethereum mainnet is required, Khedra works with any EVM-compatible blockchain, including:
- Test networks (Sepolia, etc.)
- Layer 2 solutions (Optimism, Arbitrum)
- Alternative EVMs (Gnosis Chain, etc.)

Each chain requires only a valid RPC endpoint to begin indexing.

### 3. Modular Service Architecture

Khedra operates through five interconnected services:
- **Control Service**: Central management API
- **Scraper Service**: Builds and maintains the Unchained Index
- **Monitor Service**: Tracks specific addresses of interest
- **API Service**: Provides data access via REST endpoints
- **IPFS Service**: Enables distributed sharing of index data

These services can be enabled or disabled independently to suit your needs.

### 4. Privacy-Preserving Design

Unlike traditional blockchain explorers that track user behavior, Khedra:
- Runs entirely on your local machine
- Never sends queries to third-party servers
- Doesn't track or log your address lookups
- Gives you complete control over your data

### 5. Distributed Index Sharing

The Unchained Index can be optionally shared and downloaded via IPFS, creating a collaborative network where:
- Users can contribute to building different parts of the index
- New users can download existing index portions instead of rebuilding
- The index becomes more resilient through distribution

## Use Cases

Khedra excels in numerous blockchain data scenarios:

- **Account History**: Track complete transaction and interaction history for any address
- **Balance Tracking**: Monitor native and ERC-20 token balances over time
- **Smart Contract Monitoring**: Watch for specific events or interactions with contracts
- **Auditing and Accounting**: Export complete financial histories for tax or business purposes
- **Custom Indexing**: Build specialized indices for specific protocols or applications
- **Data Analysis**: Extract patterns and insights from comprehensive on-chain data

## Getting Started

The following chapters will guide you through:

1. Installing and configuring Khedra
2. Understanding the core concepts and architecture
3. Using the various components and services
4. Advanced operations and customization
5. Maintenance and troubleshooting

Whether you're a developer, researcher, trader, or blockchain enthusiast, Khedra provides the tools you need to extract maximum value from blockchain data while maintaining your privacy and autonomy.

## Implementation Details

The core features of Khedra described in this introduction are implemented in the following Go files:

- **Main Application Structure**: The primary application entry point is defined in [`main.go`](/Users/jrush/Development/trueblocks-core/khedra/main.go) which initializes the `KhedraApp` struct defined in [`app/app.go`](/Users/jrush/Development/trueblocks-core/khedra/app/app.go).

- **Service Architecture Implementation**: 
  - The service framework is initialized in [`app/action_daemon.go`](/Users/jrush/Development/trueblocks-core/khedra/app/action_daemon.go)
  - Service definitions and validation are in [`pkg/types/service.go`](/Users/jrush/Development/trueblocks-core/khedra/pkg/types/service.go)
  - Service initialization happens in the `daemonAction` function

- **Configuration System**: 
  - Configuration loading and validation: [`app/config.go`](/Users/jrush/Development/trueblocks-core/khedra/app/config.go)
  - Environment variable processing: [`pkg/types/apply_env.go`](/Users/jrush/Development/trueblocks-core/khedra/pkg/types/apply_env.go)

- **Multi-Chain Support**: 
  - Chain configuration and validation: [`app/action_init_chains.go`](/Users/jrush/Development/trueblocks-core/khedra/app/action_init_chains.go)
  - RPC connection validation: [`pkg/validate/try_connect.go`](/Users/jrush/Development/trueblocks-core/khedra/pkg/validate/try_connect.go)
