# General Configuration Screen

## Function

```ascii
┌──────────────────────────────────────────────────────────────────────────────┐
│ General Settings                                                             │
├──────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│ The General group of options controls where Khedra stores the Unchained      │
│ Index and its caches. It also helps you choose a download strategy for       │
│ the index and helps you set up Khedra's logging options.                     │
│                                                                              │
│ Choose your folders carefully. The index and logs can get quite large        │
│ depending on the configuration. As always, type "help" to get more           │
│ information.                                                                 │
│                                                                              │
│ You may use $HOME or ~/ in your paths to refer to your home directory.       │
│                                                                              │
│ Press enter to continue.                                                     │
└──────────────────────────────────────────────────────────────────────────────┘
```

## Purpose

- Allows users to configure high-level application settings.

## Key Features

- Define global settings that apply across the application.

## Example Usage

```go
screen := getGeneralScreen(cfg)
wizard.AddScreen(screen)
```
