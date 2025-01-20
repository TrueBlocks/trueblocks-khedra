# Chains Configuration Screen

## Function

`getChainsScreen(cfg *types.Config) wizard.Screen`

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
