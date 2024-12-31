# Technical Design

## Configuration Files and Environment Variables

Khedra uses a `.env` file for configuration. Key variables include:

- `TB_NODE_DATADIR`: Directory for storing data.
- `TB_NODE_MAINNETRPC`: RPC endpoint for Ethereum mainnet.
- `TB_NODE_CHAINS`: List of chains to index.

## Initialization Process

1. Validate `.env` configuration.
2. Connect to RPC endpoints for the specified chains.
3. Initialize the blockchain index if necessary.

## Data Flow and Processing

- **Input**: Blockchain data retrieved via RPC.
- **Processing**: Indexing, storing, and optionally pinning data to IPFS.
- **Output**: Indexed data accessible through the REST API.

## Error Handling and Logging

Logs are written to the console with adjustable levels (`Debug`, `Info`, `Warn`, `Error`). Errors during initialization or RPC interactions are logged and reported.
