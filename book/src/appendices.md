# Appendices

## Glossary of Terms

## Glossary of Technical Terms

- **EVM**: Ethereum Virtual Machine, the runtime environment for smart contracts in Ethereum and similar blockchains.
- **RPC**: Remote Procedure Call, a protocol for interacting with blockchain nodes.
- **Indexing**: The process of organizing blockchain data for fast and efficient retrieval.
- **IPFS**: InterPlanetary File System, a decentralized storage solution for sharing and retrieving data.

## Frequently Asked Questions (FAQ)

### 1. What chains are supported by Khedra?

Khedra supports Ethereum mainnet and other EVM-compatible chains such as Sepolia and Gnosis. Additional chains can be added by configuring the `TB_NODE_CHAINS` environment variable.

### 2. Do I need an RPC endpoint for every chain?

Yes, each chain you want to index or interact with requires a valid RPC endpoint specified in the `.env` file.

### 3. Can I run Khedra without IPFS?

Yes, IPFS integration is optional and can be enabled or disabled using the `--ipfs` command-line option.

## References and Further Reading


## Additional Technical References and Resources

- [TrueBlocks GitHub Repository](https://github.com/TrueBlocks/trueblocks-khedra)
- [TrueBlocks Official Website](https://trueblocks.io)
- [Ethereum Developer Documentation](https://ethereum.org/en/developers/)
- [IPFS Documentation](https://docs.ipfs.io)
## Index

- Address Monitoring:
  - Implementation: [`app/action_daemon.go`](/Users/jrush/Development/trueblocks-core/khedra/app/action_daemon.go) (Monitor service initialization and `MonitorsOptions` struct)

- API Access:
  - Implementation: [`app/action_daemon.go`](/Users/jrush/Development/trueblocks-core/khedra/app/action_daemon.go) (API service initialization)

- Blockchain Indexing:
  - Implementation: [`app/action_daemon.go`](/Users/jrush/Development/trueblocks-core/khedra/app/action_daemon.go) (Scraper service initialization)

- Chains Configuration:
  - Implementation: 
    - [`app/action_init_chains.go`](/Users/jrush/Development/trueblocks-core/khedra/app/action_init_chains.go) (Chain wizard implementation)
    - [`pkg/types/chain.go`](/Users/jrush/Development/trueblocks-core/khedra/pkg/types/chain.go) (Chain struct definition and validation)

- Configuration Management:
  - Implementation: 
    - [`app/config.go`](/Users/jrush/Development/trueblocks-core/khedra/app/config.go) (Configuration loading and initialization)
    - [`app/action_config_show.go`](/Users/jrush/Development/trueblocks-core/khedra/app/action_config_show.go) (Show config command)
    - [`app/action_config_edit.go`](/Users/jrush/Development/trueblocks-core/khedra/app/action_config_edit.go) (Edit config command)

- Glossary: Chapter 7, Section "Glossary of Terms"
- IPFS Integration:
  - Implementation: 
    - [`app/action_daemon.go`](/Users/jrush/Development/trueblocks-core/khedra/app/action_daemon.go) (IPFS service initialization)
    - [`pkg/types/service.go`](/Users/jrush/Development/trueblocks-core/khedra/pkg/types/service.go) (IPFS service definition)

- Logging and Debugging:
  - Implementation: 
    - [`app/action_init_logging.go`](/Users/jrush/Development/trueblocks-core/khedra/app/action_init_logging.go) (Logging configuration)
    - [`pkg/types/general.go`](/Users/jrush/Development/trueblocks-core/khedra/pkg/types/general.go) (Logging struct definition)

- RPC Endpoints:
  - Implementation: 
    - [`pkg/validate/try_connect.go`](/Users/jrush/Development/trueblocks-core/khedra/pkg/validate/try_connect.go) (RPC connection validation)
    - [`app/has_valid_rpc.go`](/Users/jrush/Development/trueblocks-core/khedra/app/has_valid_rpc.go) (RPC validation logic)

- Service Configuration:
  - Implementation: 
    - [`app/action_init_services.go`](/Users/jrush/Development/trueblocks-core/khedra/app/action_init_services.go) (Services wizard implementation)
    - [`pkg/types/service.go`](/Users/jrush/Development/trueblocks-core/khedra/pkg/types/service.go) (Service struct definition and validation)

- Troubleshooting:
  - Implementation: Error handling throughout the codebase, especially in:
    - [`app/action_daemon.go`](/Users/jrush/Development/trueblocks-core/khedra/app/action_daemon.go) (Service error handling)
    - [`app/config.go`](/Users/jrush/Development/trueblocks-core/khedra/app/config.go) (Configuration error handling)

- Wizard Interface:
  - Implementation: 
    - [`pkg/wizard/`](/Users/jrush/Development/trueblocks-core/khedra/pkg/wizard/) directory (Wizard framework)
    - [`app/action_init.go`](/Users/jrush/Development/trueblocks-core/khedra/app/action_init.go) (Wizard initialization)


---

## Technical Index (from Technical Appendices)

- **Address Monitoring**: Section 3, Core Functionalities
- **API Access**: Section 3, Core Functionalities
- **Architecture Overview**: Section 2, System Architecture
- **Blockchain Indexing**: Section 3, Core Functionalities
- **Configuration Files**: Section 4, Technical Design
- **Data Flow**: Section 4, Technical Design
- **Error Handling**: Section 4, Technical Design
- **Integration Points**: Section 8, Integration Points
- **IPFS Integration**: Section 3, Core Functionalities; Section 8, Integration Points
- **Logging**: Section 4, Technical Design
- **Performance Tuning**: Section 7, Performance and Scalability (benchmarks removed; only tuning guidance retained)
- **REST API**: Section 3, Core Functionalities; Section 8, Integration Points
- **RPC Requirements**: Section 5, Supported Chains
- **Scalability Strategies**: Section 7, Performance and Scalability
- **System Components**: Section 2, System Architecture
- **Testing Guidelines**: Section 9, Testing and Validation
