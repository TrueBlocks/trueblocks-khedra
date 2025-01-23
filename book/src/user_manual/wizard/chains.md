# Chains Configuration Screen

## Function

```ascii
┌──────────────────────────────────────────────────────────────────────────────┐
│ Chain Settings                                                               │
├──────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│ Khedra will index any number of EVM chains, however it requires an           │
│ RPC endpoint for each to do so. Fast, dedicated local endpoints are          │
│ preferred. Likely, you will get rate limited if you point to a remote        │
│ endpoing, but if you do, you may use the Sleep option to slow down           │
│ operation. See "help".                                                       │
│                                                                              │
│ You may add chains to the list by typing the chain's name. Remove chains     │
│ with "remove <chain>". Or, an easier way is to edit the configuration        │
│ file directly by typing "edit". The mainnet chain is required.               │
│                                                                              │
│ Press enter to continue.                                                     │
└──────────────────────────────────────────────────────────────────────────────┘
```

## Purpose

- Configure blockchain networks (chains) for monitoring or interaction.

## Key Features

- Select relevant chains.
- Set chain-specific parameters such as endpoints and authentication.

## Example Usage

```go
screen := getChainsScreen(cfg)
wizard.AddScreen(screen)
```
