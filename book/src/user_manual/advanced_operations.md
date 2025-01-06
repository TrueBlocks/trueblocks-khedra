# Advanced Operations

## Integrating with IPFS

Enable IPFS support with:

```bash
./khedra --ipfs on
```

This will pin indexed blockchain data to IPFS, ensuring decentralized storage and retrieval.

## Customizing Chain Indexing

Specify additional chains by updating the `TB_NODE_CHAINS` environment variable. Example:

```env
TB_NODE_CHAINS="mainnet,sepolia,gnosis"
```

Ensure each chain has a valid RPC endpoint configured.

## Utilizing Command-Line Options

Key options include:

- `--init [all|blooms|none]`: Specify the type of index initialization.
- `--scrape [on|off]`: Enable or disable the scraper.
- `--api [on|off]`: Enable or disable the API.
- `--sleep [int]`: Set the sleep duration between updates in seconds.
