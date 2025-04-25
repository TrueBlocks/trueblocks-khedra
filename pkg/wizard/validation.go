package wizard

import (
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
)

// ValidationFeedback represents the result of real-time validation
type ValidationFeedback struct {
	IsValid  bool
	Message  string
	Severity string // "error", "warning", "info"
}

// ValidationFunc is a function type for real-time validation
type ValidationFunc func(input string) ValidationFeedback

// Global map of registered validation functions
var validationFunctions = make(map[string]ValidationFunc)

// RegisterValidationFunc registers a validation function for a specific question type
func RegisterValidationFunc(questionType string, validationFunc ValidationFunc) {
	validationFunctions[questionType] = validationFunc
}

// GetValidationFunc retrieves a validation function for a question type
func GetValidationFunc(questionType string) ValidationFunc {
	if fn, ok := validationFunctions[questionType]; ok {
		return fn
	}
	return nil
}

// FormatValidationFeedback formats validation feedback with appropriate colors
func FormatValidationFeedback(feedback ValidationFeedback) string {
	if !feedback.IsValid {
		return colors.Red + "✗ " + feedback.Message + colors.Off
	}

	switch feedback.Severity {
	case "warning":
		return colors.BrightYellow + "⚠ " + feedback.Message + colors.Off
	case "info":
		return colors.BrightBlue + "ℹ " + feedback.Message + colors.Off
	default:
		return colors.Green + "✓ " + feedback.Message + colors.Off
	}
}
