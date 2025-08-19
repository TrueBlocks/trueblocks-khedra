# Supported Chains

This section details the blockchain networks supported by Khedra, the technical requirements for each, and the implementation approaches for multi-chain support.

## Chain Support Architecture

Khedra implements a flexible architecture for supporting multiple EVM-compatible blockchains simultaneously.

### Chain Abstraction Layer

At the core of Khedra's multi-chain support is a chain abstraction layer that:

1. Normalizes differences between chain implementations
2. Provides a uniform interface for blockchain interactions
3. Manages chain-specific configurations and behaviors
4. Isolates chain-specific code from the core application logic

```go
// Simplified Chain interface
type Chain interface {
    // Return the chain name
    Name() string
    
    // Return the chain ID
    ChainID() uint64
    
    // Get RPC client for this chain
    Client() rpc.Client
    
    // Get path to chain-specific data directory
    DataDir() string
    
    // Check if this chain requires special handling for a feature
    SupportsFeature(feature string) bool
    
    // Get chain-specific configuration
    Config() ChainConfig
}
```

## Core Chain Requirements

For a blockchain to be fully supported by Khedra, it must meet these technical requirements:

### RPC Support

The chain must provide an Ethereum-compatible JSON-RPC API with these essential methods:

1. **Basic Methods**:
   - `eth_blockNumber`: Get the latest block number
   - `eth_getBlockByNumber`: Retrieve block data
   - `eth_getTransactionReceipt`: Get transaction receipts with logs
   - `eth_chainId`: Return the chain identifier

2. **Trace Support**:
   - Either `debug_traceTransaction` or `trace_transaction`: Retrieve execution traces
   - Alternatively: `trace_block` or `debug_traceBlockByNumber`: Get all traces in a block

### Data Structures

The chain must use compatible data structures:

1. **Addresses**: 20-byte Ethereum-compatible addresses
2. **Transactions**: Compatible transaction format with standard fields
3. **Logs**: EVM-compatible event logs
4. **Traces**: Call traces in a format compatible with Khedra's processors

### Consensus and Finality

The chain should have:

1. **Deterministic Finality**: Clear rules for when blocks are considered final
2. **Manageable Reorgs**: Limited blockchain reorganizations
3. **Block Time Consistency**: Relatively consistent block production times

## Ethereum Mainnet

Ethereum mainnet is the primary supported chain and is required even when indexing other chains.

### Special Considerations

1. **Block Range**: Support for full historical range from genesis
2. **Archive Node**: Full archive node required for historical traces
3. **Trace Support**: Must support either Geth or Parity trace methods
4. **Size Considerations**: Largest data volume among supported chains

### Implementation Details

```go
// Ethereum mainnet-specific configuration
type EthereumMainnetChain struct {
    BaseChain
    traceMethod string  // "geth" or "parity" style traces
}

func (c *EthereumMainnetChain) ProcessTraces(traces []interface{}) []Appearance {
    // Mainnet-specific trace processing logic
    // ...
}
```

## EVM-Compatible Chains

Khedra supports a variety of EVM-compatible chains with minimal configuration.

### Officially Supported Chains

These chains are officially supported with tested implementations:

1. **Ethereum Testnets**:
   - Sepolia
   - Goerli (legacy support)

2. **Layer 2 Networks**:
   - Optimism
   - Arbitrum
   - Polygon

3. **EVM Sidechains**:
   - Gnosis Chain (formerly xDai)
   - Avalanche C-Chain
   - Binance Smart Chain

### Chain Configuration

Each chain is configured with these parameters:

```yaml
chains:
  mainnet:  # Chain identifier
    rpcs:   # List of RPC endpoints
      - "https://ethereum-rpc.example.com"
    enabled: true  # Whether the chain is active
    trace_support: "geth"  # Trace API style
    # Chain-specific overrides
    scraper:
      batch_size: 500
```

### Chain-Specific Adaptations

Some chains require special handling:

1. **Optimism/Arbitrum**: Modified trace processing for rollup architecture
2. **Polygon**: Adjusted finality assumptions for PoS consensus
3. **BSC/Avalanche**: Faster block times requiring different batch sizing

## Chain Detection and Validation

Khedra implements robust chain detection and validation:

### Auto-Detection

When connecting to an RPC endpoint:

1. Query `eth_chainId` to determine the actual chain
2. Verify against the configured chain identifier
3. Detect trace method support through API probing
4. Identify chain-specific capabilities

### Connection Management

For each configured chain:

1. **Connection Pool**: Maintain multiple connections for parallel operations
2. **Failover Support**: Try alternative endpoints when primary fails
3. **Health Monitoring**: Track endpoint reliability and performance
4. **Rate Limiting**: Respect provider-specific rate limits

## Data Isolation

Khedra maintains strict data isolation between chains:

1. **Chain-Specific Directories**: Separate storage locations for each chain
2. **Independent Indices**: Each chain has its own Unchained Index
3. **Configuration Isolation**: Chain-specific settings don't affect other chains
4. **Parallel Processing**: Chains can be processed concurrently

## Adding New Chain Support

For adding support for a new EVM-compatible chain:

1. **Configuration**: Add the chain definition to `config.yaml`
2. **Custom Handling**: Implement any chain-specific processors if needed
3. **Testing**: Verify compatibility with the new chain
4. **Documentation**: Update supported chains documentation

### Example: Adding a New Chain

```go
// Register a new chain type
func RegisterNewChain() {
    registry.RegisterChain("new-chain", func(config ChainConfig) (Chain, error) {
        return &NewChain{
            BaseChain: NewBaseChain(config),
            // Chain-specific initialization
        }, nil
    })
}

// Implement chain-specific behavior
type NewChain struct {
    BaseChain
    // Chain-specific fields
}

func (c *NewChain) SupportsFeature(feature string) bool {
    // Chain-specific feature support
    switch feature {
    case "trace":
        return true
    case "state_diff":
        return false
    default:
        return c.BaseChain.SupportsFeature(feature)
    }
}
```

Khedra's flexible chain support architecture allows it to adapt to the evolving ecosystem of EVM-compatible blockchains while maintaining consistent indexing and monitoring capabilities across all supported networks.
