package wizard

import (
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
)

// KeyboardShortcut represents a keyboard shortcut with a key and description
type KeyboardShortcut struct {
	Key         string
	Description string
}

// GetShortcutBar returns a formatted string of keyboard shortcuts
func GetShortcutBar(shortcuts []KeyboardShortcut) string {
	var output strings.Builder

	output.WriteString("\n")
	// output.WriteString(colors.Black + "Keyboard shortcuts: " + colors.Off)
	output.WriteString(" ")

	for i, shortcut := range shortcuts {
		if i > 0 {
			output.WriteString(" | ")
		}
		output.WriteString(colors.BrightBlue + shortcut.Key + colors.Off + ": " + shortcut.Description)
	}

	return output.String()
}

// GetDefaultShortcuts returns the default set of keyboard shortcuts
func GetDefaultShortcuts() []KeyboardShortcut {
	return []KeyboardShortcut{
		{"Enter", "Continue"},
		{"b", "Back"},
		{"h", "Help"},
		{"q", "Quit"},
	}
}

// GetShortcutBarForScreen returns a shortcut bar appropriate for the current screen
func GetShortcutBarForScreen(screenTitle string, w *Wizard) string {
	_ = w           // delint
	_ = screenTitle // delint
	shortcuts := GetDefaultShortcuts()

	// Add context-specific shortcuts based on screen title
	// title := strings.ToLower(screenTitle)

	// if strings.Contains(title, "summary") {
	// 	// Add summary-specific shortcuts - but remove template functionality
	// 	// No more "save" shortcut for templates
	// }

	return GetShortcutBar(shortcuts)
}
