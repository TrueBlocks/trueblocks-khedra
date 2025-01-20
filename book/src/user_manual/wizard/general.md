# General Configuration Screen

## Function

`getGeneralScreen(cfg *types.Config) wizard.Screen`

## Purpose

- Allows users to configure high-level application settings.

## Key Features

- Define global settings that apply across the application.

## Example Usage

```go
screen := getGeneralScreen(cfg)
wizard.AddScreen(screen)
```
