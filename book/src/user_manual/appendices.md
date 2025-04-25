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

- Address Monitoring:
  - Documentation: Chapter 4, Section "Monitoring Addresses"
  - Implementation: [`app/action_daemon.go`](/Users/jrush/Development/trueblocks-core/khedra/app/action_daemon.go) (Monitor service initialization and `MonitorsOptions` struct)

- API Access:
  - Documentation: Chapter 4, Section "Accessing the REST API"
  - Implementation: [`app/action_daemon.go`](/Users/jrush/Development/trueblocks-core/khedra/app/action_daemon.go) (API service initialization)

- Blockchain Indexing:
  - Documentation: Chapter 4, Section "Indexing Blockchains"
  - Implementation: [`app/action_daemon.go`](/Users/jrush/Development/trueblocks-core/khedra/app/action_daemon.go) (Scraper service initialization)

- Chains Configuration:
  - Documentation: Chapter 3, Section "Terminology and Concepts"
  - Implementation: 
    - [`app/action_init_chains.go`](/Users/jrush/Development/trueblocks-core/khedra/app/action_init_chains.go) (Chain wizard implementation)
    - [`pkg/types/chain.go`](/Users/jrush/Development/trueblocks-core/khedra/pkg/types/chain.go) (Chain struct definition and validation)

- Configuration Management:
  - Documentation: Chapter 4, Section "Managing Configurations"
  - Implementation: 
    - [`app/config.go`](/Users/jrush/Development/trueblocks-core/khedra/app/config.go) (Configuration loading and initialization)
    - [`app/action_config_show.go`](/Users/jrush/Development/trueblocks-core/khedra/app/action_config_show.go) (Show config command)
    - [`app/action_config_edit.go`](/Users/jrush/Development/trueblocks-core/khedra/app/action_config_edit.go) (Edit config command)

- Glossary: Chapter 7, Section "Glossary of Terms"

- IPFS Integration:
  - Documentation: Chapter 5, Section "Integrating with IPFS"
  - Implementation: 
    - [`app/action_daemon.go`](/Users/jrush/Development/trueblocks-core/khedra/app/action_daemon.go) (IPFS service initialization)
    - [`pkg/types/service.go`](/Users/jrush/Development/trueblocks-core/khedra/pkg/types/service.go) (IPFS service definition)

- Logging and Debugging:
  - Documentation: Chapter 6, Section "Log Files and Debugging"
  - Implementation: 
    - [`app/action_init_logging.go`](/Users/jrush/Development/trueblocks-core/khedra/app/action_init_logging.go) (Logging configuration)
    - [`pkg/types/general.go`](/Users/jrush/Development/trueblocks-core/khedra/pkg/types/general.go) (Logging struct definition)

- RPC Endpoints:
  - Documentation: Chapter 2, Section "Initial Configuration"
  - Implementation: 
    - [`pkg/validate/try_connect.go`](/Users/jrush/Development/trueblocks-core/khedra/pkg/validate/try_connect.go) (RPC connection validation)
    - [`app/has_valid_rpc.go`](/Users/jrush/Development/trueblocks-core/khedra/app/has_valid_rpc.go) (RPC validation logic)

- Service Configuration:
  - Documentation: Chapter 2, Section "Configuration File Format"
  - Implementation: 
    - [`app/action_init_services.go`](/Users/jrush/Development/trueblocks-core/khedra/app/action_init_services.go) (Services wizard implementation)
    - [`pkg/types/service.go`](/Users/jrush/Development/trueblocks-core/khedra/pkg/types/service.go) (Service struct definition and validation)

- Troubleshooting:
  - Documentation: Chapter 6, Section "Troubleshooting"
  - Implementation: Error handling throughout the codebase, especially in:
    - [`app/action_daemon.go`](/Users/jrush/Development/trueblocks-core/khedra/app/action_daemon.go) (Service error handling)
    - [`app/config.go`](/Users/jrush/Development/trueblocks-core/khedra/app/config.go) (Configuration error handling)

- Wizard Interface:
  - Documentation: Chapter 6, Section "Installation Wizard"
  - Implementation: 
    - [`pkg/wizard/`](/Users/jrush/Development/trueblocks-core/khedra/pkg/wizard/) directory (Wizard framework)
    - [`app/action_init.go`](/Users/jrush/Development/trueblocks-core/khedra/app/action_init.go) (Wizard initialization)
