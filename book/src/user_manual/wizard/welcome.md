# Welcome Screen

## Function

```ascii
┌──────────────────────────────────────────────────────────────────────────────┐
│ ╔═══════════════════════════════════════════════════════╗                    │
│ ║                     KHEDRA WIZARD                     ║                    │
│ ║                                                       ║                    │
│ ║   Index, monitor, serve, and share blockchain data.   ║                    │
│ ╚═══════════════════════════════════════════════════════╝                    │
│                                                                              │
│ Welcome to Khedra, the world's only local-first indexer/monitor for          │
│ EVM blockchains. This wizard will help you configure Khedra. There are       │
│ three groups of settings: General, Services, and Chains.                     │
│                                                                              │
│ Type "q" or "quit" to quit, "b" or "back" to return to a previous screen,    │
│ or "help" to get more information.                                           │
│                                                                              │
│ Press enter to continue.                                                     │
└──────────────────────────────────────────────────────────────────────────────┘
```

## Purpose

- Introduces the wizard to the user.
- Provides navigation instructions.

## Key Features

- Display introduction message.
- Outline navigation commands (e.g., "quit," "back").

## Example Usage

```go
screen := getWelcomeScreen(cfg)
wizard.AddScreen(screen)
```
