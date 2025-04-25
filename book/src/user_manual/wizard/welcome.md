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
│ "h" or "help" to get more information, or "e" to edit the file directly.     │
│                                                                              │
│ Press enter to continue.                                                     │
│                                                                              │
│ Keyboard: [h] Help [q] Quit [b] Back [enter] Continue                        │
└──────────────────────────────────────────────────────────────────────────────┘
```

## Purpose

- Introduces the wizard to the user
- Orients the user to the configuration process
- Provides clear navigation instructions

## Navigation Options

- **Enter**: Proceed to the next screen
- **h/help**: Open browser with documentation
- **q/quit**: Exit the wizard
- **b/back**: Return to previous screen
- **e/edit**: Edit configuration file directly

You may directly edit the configuration from any screen by typing `e` or `edit`. This will open the configuration file in the user's preferred text editor (defined by the EDITOR environment variable).

The welcome screen serves as the entry point to the configuration process, designed to be approachable while providing clear direction on how to proceed.
