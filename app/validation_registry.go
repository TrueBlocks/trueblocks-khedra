package app

import (
	"strings"

	"github.com/TrueBlocks/trueblocks-khedra/v5/pkg/types"
	"github.com/TrueBlocks/trueblocks-khedra/v5/pkg/wizard"
)

// registerValidationFunctions registers all the real-time validation functions
func registerValidationFunctions() {
	// Register validation functions for context-based validation
	wizard.RegisterValidationFunc("rpc", func(input string) wizard.ValidationFeedback {
		// Create an adapter between the old and new validation function types
		if _, err := types.ValidateRpcEndpointRT(input); err != nil {
			return wizard.ValidationFeedback{IsValid: false, Message: err.Error(), Severity: "error"}
		}
		result := types.TestRpcEndpoint(input)
		if strings.Contains(strings.ToLower(result.ErrorMessage), "websocket endpoints cannot be fully tested") {
			return wizard.ValidationFeedback{IsValid: false, Message: "WebSocket RPC endpoints (ws://, wss://) are not supported; use http:// or https://", Severity: "error"}
		}
		if !result.Reachable {
			return wizard.ValidationFeedback{IsValid: false, Message: result.ErrorMessage, Severity: "error"}
		}
		formattedResult := types.FormatRpcTestResult(result)
		return wizard.ValidationFeedback{IsValid: true, Message: formattedResult, Severity: "info"}
	})

	wizard.RegisterValidationFunc("chainid", func(input string) wizard.ValidationFeedback {
		isValid, message, severity := types.ValidateWithFeedback(input, "chainId")
		return wizard.ValidationFeedback{
			IsValid:  isValid,
			Message:  message,
			Severity: severity,
		}
	})

	wizard.RegisterValidationFunc("port", func(input string) wizard.ValidationFeedback {
		isValid, message, severity := types.ValidateWithFeedback(input, "port")
		return wizard.ValidationFeedback{
			IsValid:  isValid,
			Message:  message,
			Severity: severity,
		}
	})

	wizard.RegisterValidationFunc("folder", func(input string) wizard.ValidationFeedback {
		isValid, errorMsg, warningMsg := types.ValidateFolder(input)
		message := errorMsg
		severity := "error"

		if isValid {
			message = warningMsg
			if message != "" {
				severity = "warning"
			} else {
				message = "Folder is valid"
				severity = "info"
			}
		}

		return wizard.ValidationFeedback{
			IsValid:  isValid,
			Message:  message,
			Severity: severity,
		}
	})

	wizard.RegisterValidationFunc("loglevel", func(input string) wizard.ValidationFeedback {
		// Simple validation for log levels
		validLevels := []string{"debug", "info", "warn", "error"}
		isValid := false

		for _, level := range validLevels {
			if input == level {
				isValid = true
				break
			}
		}

		if isValid {
			return wizard.ValidationFeedback{
				IsValid:  true,
				Message:  "Log level is valid",
				Severity: "info",
			}
		}

		return wizard.ValidationFeedback{
			IsValid:  false,
			Message:  "Log level must be one of: debug, info, warn, error",
			Severity: "error",
		}
	})
}
