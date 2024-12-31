# Integration Points

## Integration with External APIs

Khedra exposes data through a REST API, making it compatible with external applications. Example use cases:

- Fetching transaction details for a given address.
- Retrieving block information for analysis.

## Interfacing with IPFS

Data indexed by Khedra can be pinned to IPFS for decentralized storage:

```bash
./khedra --ipfs on
```

## Customizing for Specific Use Cases

Users can tailor the configuration by:

- Adjusting `.env` variables to include specific chains and RPC endpoints.
- Writing custom scripts to query the REST API and process the data.
