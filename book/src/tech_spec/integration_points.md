# Integration Points

## Integration with External APIs

Khedra exposes data through a REST API, making it compatible with external applications. Example use cases:

- Fetching transaction details for a given address.
- Retrieving block information for analysis.

## Interfacing with IPFS

If the IPFS service is enabled in `config.yaml` it will be started with the daemon. There is **no** `--ipfs on` CLI flag; previous documentation using that syntax was incorrect.

## Customizing for Specific Use Cases

Users can tailor the configuration by:

- Editing `config.yaml` (wizard or `khedra config edit`) to set chains, RPC endpoints, and service enablement.
- Using shell scripts to automate pause/unpause via the control endpoints.
