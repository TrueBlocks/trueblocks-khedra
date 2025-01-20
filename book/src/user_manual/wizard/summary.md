# Summary Screen

## Function

`getSummaryScreen(cfg *types.Config) wizard.Screen`

## Purpose

- Summarizes the configurations selected in previous steps.

## Key Features

- Display all settings for review.
- Allow users to confirm or modify selections.

## Example Usage

```go
screen := getSummaryScreen(cfg)
wizard.AddScreen(screen)
```
