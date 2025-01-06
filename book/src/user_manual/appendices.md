# Appendices

## Glossary of Terms

- **EVM**: Ethereum Virtual Machine, the runtime environment for smart contracts in Ethereum and similar blockchains.
- **RPC**: Remote Procedure Call, a protocol allowing the application to communicate with blockchain nodes.
- **Indexing**: The process of organizing blockchain data for fast and efficient retrieval.
- **IPFS**: InterPlanetary File System, a decentralized storage system for sharing and retrieving data.

## Frequently Asked Questions (FAQ)

### 1. What chains are supported by Khedra?

Khedra supports Ethereum mainnet and other EVM-compatible chains such as Sepolia and Gnosis. Additional chains can be added by configuring the `TB_NODE_CHAINS` environment variable.

### 2. Do I need an RPC endpoint for every chain?

Yes, each chain you want to index or interact with requires a valid RPC endpoint specified in the `.env` file.

### 3. Can I run Khedra without IPFS?

Yes, IPFS integration is optional and can be enabled or disabled using the `--ipfs` command-line option.

## References and Further Reading

- [TrueBlocks GitHub Repository](https://github.com/TrueBlocks/trueblocks-khedra)
- [TrueBlocks Official Website](https://trueblocks.io)
- [Ethereum Developer Documentation](https://ethereum.org/en/developers/)
- [IPFS Documentation](https://docs.ipfs.io)

## Index

- Address Monitoring: Chapter 4, Section "Monitoring Addresses"
- Advanced Operations: Chapter 5
- API Access: Chapter 4, Section "Accessing the REST API"
- Blockchain Indexing: Chapter 4, Section "Indexing Blockchains"
- Chains: Chapter 3, Section "Terminology and Concepts"
- Configuration Management: Chapter 4, Section "Managing Configurations"
- Glossary: Chapter 7, Section "Glossary of Terms"
- IPFS Integration: Chapter 5, Section "Integrating with IPFS"
- Logging and Debugging: Chapter 6, Section "Log Files and Debugging"
- RPC Endpoints: Chapter 2, Section "Initial Configuration"
- Troubleshooting: Chapter 6
