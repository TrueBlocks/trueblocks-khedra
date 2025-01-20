# Welcome Screen

## Function

`getWelcomeScreen(cfg *types.Config) wizard.Screen`

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
