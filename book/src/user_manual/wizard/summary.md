# Summary Screen

## Function

```ascii
┌──────────────────────────────────────────────────────────────────────────────┐
│ Summary                                                                      │
├──────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│ Question: Would you like to edit the config by hand?                         │
│ Current:  no                                                                 │
│                                                                              │
│                                                                              │
│                                                                              │
│                                                                              │
│                                                                              │
│                                                                              │
│                                                                              │
│                                                                              │
│ Press enter to finish the wizard. ("b"=back, "h"=help)                       │
└──────────────────────────────────────────────────────────────────────────────┘
```

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
