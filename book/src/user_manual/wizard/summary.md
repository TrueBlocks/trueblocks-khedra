# Summary Screen

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
│                                                                              │
│ Keyboard: [h] Help [q] Quit [b] Back [e] Edit [enter] Finish                 │
└──────────────────────────────────────────────────────────────────────────────┘
```

## Purpose

- Provides a review of all configured settings
- Offers a final chance to make adjustments before saving
- Summarizes the configuration in a clear, readable format

## Configuration Summary Display

The summary screen presents the configuration organized by section:

1. **General Settings**
   - Data folder location
   - Download strategy
   - Logging configuration

2. **Services Configuration**
   - Enabled/disabled status for each service
   - Port numbers and key parameters
   - Resource allocations

3. **Chain Settings**
   - Configured blockchains
   - RPC endpoints
   - Chain-specific settings

## Final Options

From the summary screen, you can:

1. **Finish**: Accept the configuration and write it to the config file
2. **Edit**: Open the configuration in a text editor for manual changes
3. **Back**: Return to previous screens to make adjustments
4. **Help**: Access documentation about configuration options
5. **Quit**: Exit without saving changes

When the user chooses to finish, the wizard writes the configuration to `~/.khedra/config.yaml` by default, or to an alternative location if specified during the process.

If the user chooses to edit the file directly, the wizard will invoke the system's default editor (or the editor specified in the EDITOR environment variable) and then reload the configuration after editing.
