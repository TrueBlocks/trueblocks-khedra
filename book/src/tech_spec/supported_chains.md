# Supported Chains

## List of Supported Blockchains

Khedra supports Ethereum mainnet and other EVM-compatible chains like:

- Sepolia
- Gnosis
- Optimism

## Requirements for RPC Endpoints

Each chain requires a valid RPC endpoint. For example:

- `TB_NODE_MAINNETRPC`: Mainnet RPC URL.
- `TB_NODE_SEPOLIARPC`: Sepolia RPC URL.

## Handling Multiple Chains

To enable multiple chains, set `TB_NODE_CHAINS` in the `.env` file:

```env
TB_NODE_CHAINS="mainnet,sepolia,gnosis"
```

Ensure each chain has a corresponding RPC endpoint.
