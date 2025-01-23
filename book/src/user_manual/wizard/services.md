# Services Configuration Screen

## Function

```ascii
┌──────────────────────────────────────────────────────────────────────────────┐
│ Services Settings                                                            │
├──────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│ Khedra provides five services. The first, "control," exposes endpoints to    │
│ control the other four: "scrape", "monitor", "api", and "ipfs".              │
│                                                                              │
│ You may disable/enable any combination of services, but at least one must    │
│ be enabled.                                                                  │
│                                                                              │
│ The next few screens will allow you to configure each service.               │
│                                                                              │
│                                                                              │
│                                                                              │
│ Press enter to continue.                                                     │
└──────────────────────────────────────────────────────────────────────────────┘
```

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
