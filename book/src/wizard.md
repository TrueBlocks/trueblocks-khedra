# Wizard Screen Documentation

## Introduction

Khedra's configuration wizard provides a streamlined, interactive way to set up your installation. Rather than manually editing the `config.yaml` file, the wizard walks you through each configuration section with clear explanations and validation.

### User Interface Features

The wizard provides several helpful features:

- **Keyboard Navigation**: Use arrow keys and shortcuts to navigate
- **Contextual Help**: Press 'h' on any screen for detailed documentation
- **Editor Integration**: Press 'e' to directly edit configuration files
- **Validation**: Input is checked for correctness before proceeding
- **Visual Cues**: Consistent layout with clear indicators for navigation options

## Using the Wizard

Start the Wizard with:

```bash
khedra init
```

## Implementation Details

The configuration wizard described in this document is implemented through a package of Go files in the `pkg/wizard` directory:

### Core Wizard Framework

- **Main Wizard Structure**: [`pkg/wizard/wizard.go`](/Users/jrush/Development/trueblocks-core/khedra/pkg/wizard/wizard.go) - Defines the `Wizard` struct and methods for managing the wizard state, navigation between screens, and execution flow

- **Screen Component**: [`pkg/wizard/screen.go`](/Users/jrush/Development/trueblocks-core/khedra/pkg/wizard/screen.go) - Implements the `Screen` struct representing individual wizard pages with questions and display logic

- **Question Framework**: [`pkg/wizard/question.go`](/Users/jrush/Development/trueblocks-core/khedra/pkg/wizard/question.go) - Provides the `Question` struct and interface for gathering and validating user input

- **User Interface**: 
  - [`pkg/wizard/display.go`](/Users/jrush/Development/trueblocks-core/khedra/pkg/wizard/display.go) - Handles rendering screens and questions
  - [`pkg/wizard/style.go`](/Users/jrush/Development/trueblocks-core/khedra/pkg/wizard/style.go) - Controls visual styling of the wizard interface
  - [`pkg/boxes/boxes.go`](/Users/jrush/Development/trueblocks-core/khedra/pkg/boxes/boxes.go) - Implements the ASCII box drawing for wizard screens

- **Navigation**: 
  - [`pkg/wizard/navigation.go`](/Users/jrush/Development/trueblocks-core/khedra/pkg/wizard/navigation.go) - Implements the navigation bar and controls
  - [`pkg/wizard/shortcuts.go`](/Users/jrush/Development/trueblocks-core/khedra/pkg/wizard/shortcuts.go) - Handles keyboard shortcuts

### Wizard Screen Implementations

The specific wizard screens visible in the user interface are implemented in these files:

- **Welcome Screen**: [`app/action_init_welcome.go`](/Users/jrush/Development/trueblocks-core/khedra/app/action_init_welcome.go)
- **General Settings**: [`app/action_init_general.go`](/Users/jrush/Development/trueblocks-core/khedra/app/action_init_general.go)
- **Services Config**: [`app/action_init_services.go`](/Users/jrush/Development/trueblocks-core/khedra/app/action_init_services.go)
- **Chain Config**: [`app/action_init_chains.go`](/Users/jrush/Development/trueblocks-core/khedra/app/action_init_chains.go)
- **Summary Screen**: [`app/action_init_summary.go`](/Users/jrush/Development/trueblocks-core/khedra/app/action_init_summary.go)

### Integration with Configuration System

The wizard integrates with the configuration system through:

- **Configuration Loading**: In the `ReloaderFn` function passed to the wizard
- **Configuration Validation**: Through the validation functions for each input field
- **Configuration Saving**: In the final step of the wizard workflow

The wizard framework uses a screen-based approach with forward/backward navigation, real-time validation, and contextual help, exactly as described in this document.
