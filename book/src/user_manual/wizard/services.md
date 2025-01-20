# Services Configuration Screen

## Function

`getServicesScreen(cfg *types.Config) wizard.Screen`

## Purpose

- Enables users to select and configure services.

## Key Features

- Choose relevant services.
- Configure each selected service.

## Example Usage

```go
screen := getServicesScreen(cfg)
wizard.AddScreen(screen)
```
