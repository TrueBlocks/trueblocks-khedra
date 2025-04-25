# General Configuration Screen

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
│                                                                              │
│ Keyboard: [h] Help [q] Quit [b] Back [enter] Continue                        │
└──────────────────────────────────────────────────────────────────────────────┘
```

## Purpose

- Allows users to configure high-level application settings
- Sets up crucial file paths for data storage
- Configures logging behavior

## Key Features

- Define the main data folder location with path expansion support
- Configure index download and update strategies
- Set up logging preferences for troubleshooting
- Options for path expansion (supporting $HOME and ~/ notation)
- Disk space requirement warnings
- Input validation for directory existence and write permissions

## Configuration Options

The General Settings screen presents these key configuration options:

1. **Data Folder**: Where Khedra stores all index and cache data
   - Default: `~/.khedra/data`
   - Must be a writable location with sufficient disk space

2. **Index Download Strategy**:
   - IPFS-first: Prioritize downloading from the distributed network
   - Local-first: Prioritize building the index locally
   - Hybrid: Balance between downloading and local building
