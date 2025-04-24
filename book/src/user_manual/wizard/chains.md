# Chain Settings Screen

```ascii
┌──────────────────────────────────────────────────────────────────────────────┐
│ Chain Settings                                                               │
├──────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│ Khedra works with the Ethereum mainnet chain and any EVM-compatible          │
│ blockchain. Each chain requires at least one RPC endpoint URL and a          │
│ chain name.                                                                  │
│                                                                              │
│ Ethereum mainnet must be configured even if other chains are enabled.        │
│ The format of an RPC endpoint is protocol://host:port. For example:          │
│ http://localhost:8545 or https://mainnet.infura.io/v3/YOUR-PROJECT-ID.       │
│                                                                              │
│ The next few screens will help you configure your chains.                    │
│                                                                              │
│ Press enter to continue.                                                     │
│                                                                              │
│ Keyboard: [h] Help [q] Quit [b] Back [enter] Continue                        │
└──────────────────────────────────────────────────────────────────────────────┘
```

## Purpose

- Configures blockchain connections for indexing and monitoring
- Ensures proper RPC endpoint setup for each chain
- Explains the requirement for Ethereum mainnet

## Key Features

- Multiple chain support with standardized naming
- RPC endpoint configuration and validation
- Clear explanation of requirements and format

## Chain Configuration

The chains configuration screen guides you through setting up:

1. **Ethereum Mainnet (Required)**
   - At least one valid RPC endpoint
   - Used for core functionality and the Unchained Index

2. **Additional EVM Chains (Optional)**
   - Sepolia, Gnosis, Optimism, and other EVM-compatible chains
   - Each requires at least one RPC endpoint
   - Enable/disable option for each chain

## RPC Endpoint Requirements

For each chain, you must provide:

- A valid RPC URL in the format `protocol://host:port`
- Proper authentication details if required (e.g., Infura project ID)
- Endpoints with sufficient capabilities for indexing (archive nodes recommended)

THIS TEXT NEEDS TO BE REVIEWED.
## Validation Checks

The wizard performs these validations on each RPC endpoint:

- URL format validation
- Connection test to verify the endpoint is reachable
- Chain ID verification to ensure the endpoint matches the selected chain
- API method support check for required JSON-RPC methods
THIS TEXT NEEDS TO BE REVIEWED.

## Implementation

The chain configuration uses the Screen struct with specialized validation for RPC endpoints. The wizard prioritizes setting up Ethereum mainnet first, then offers options to configure additional chains as needed.

For each chain, the wizard walks through enabling the chain, configuring RPC endpoints, and validating the connection before proceeding to the next chain.
