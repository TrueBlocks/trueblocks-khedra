package app

import (
	"strings"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-khedra/v5/pkg/wizard"
)

// registerHelpHandler registers the context-sensitive help handler
func registerHelpHandler() {
	wizard.SetHelpHandler(func(screen *wizard.Screen, question *wizard.Question) string {
		screenTitle := strings.ToLower(screen.Title)
		questionText := ""
		if question != nil {
			questionText = strings.ToLower(question.Question)
		}

		helpContent := strings.Builder{}
		helpContent.WriteString(colors.Yellow + "📚 Help Information" + colors.Off + "\n\n")

		switch {
		case strings.Contains(screenTitle, "welcome"):
			helpContent.WriteString(getWelcomeHelp())
		case strings.Contains(screenTitle, "configuration templates"):
			helpContent.WriteString(getTemplatesHelp())
		case strings.Contains(screenTitle, "configuration mode"):
			helpContent.WriteString(getModeHelp())
		case strings.Contains(screenTitle, "general settings"):
			helpContent.WriteString(getGeneralHelp(questionText))
		case strings.Contains(screenTitle, "services"):
			helpContent.WriteString(getServicesHelp(questionText))
		case strings.Contains(screenTitle, "service port"):
			helpContent.WriteString(getPortsHelp(questionText))
		case strings.Contains(screenTitle, "chains"):
			helpContent.WriteString(getChainsHelp(questionText))
		case strings.Contains(screenTitle, "logging"):
			helpContent.WriteString(getLoggingHelp(questionText))
		case strings.Contains(screenTitle, "summary"):
			helpContent.WriteString(getSummaryHelp())
		default:
			helpContent.WriteString("No specific help is available for this screen.\n\n")
		}

		// Add navigation help for all screens
		helpContent.WriteString(colors.BrightBlue + "Navigation:" + colors.Off + "\n")
		helpContent.WriteString("- Press 'Enter' to proceed to the next item\n")
		helpContent.WriteString("- Type 'b' or 'back' to go to the previous screen\n")
		helpContent.WriteString("- Press 'Ctrl+C' to exit the wizard at any time\n")

		// Simply return the help content string - don't try to format it here
		return helpContent.String()
	})
}

// Helper functions for each screen

func getWelcomeHelp() string {
	return `Welcome to Khedra, the TrueBlocks configuration wizard!

This wizard will guide you through the process of setting up your 
TrueBlocks installation. It will help you configure:

- General settings such as data location
- Services you wish to enable
- Blockchain networks you want to interact with
- Logging configuration

You can use the help function at any time by typing 'h' or 'help'.

`
}

func getTemplatesHelp() string {
	return `Configuration Templates

Templates allow you to save and reuse configuration settings:

- Select a template number to load a previously saved template
- Press Enter to start fresh without a template
- Type "save" to save your current configuration as a template
- Type "delete" followed by a template number to remove it

Templates are useful for managing different environments (development, 
production) or for quickly setting up TrueBlocks on multiple machines.

`
}

func getModeHelp() string {
	return `Configuration Modes

TrueBlocks offers different configuration modes:

- Basic: Sets up minimal configuration for getting started
- Standard: Recommended for most users, enables common features
- Advanced: Full control over all settings
- Expert: Complete access to all configuration options

Choose the mode that best fits your needs and experience level.

`
}

func getGeneralHelp(_ string) string {
	// Provide more specific help based on the question
	return `General Settings

This section configures the core settings for TrueBlocks:

- Data Folder: Where TrueBlocks stores all of its data
- Strategy: How TrueBlocks processes the blockchain data
- Detail: The level of detail to index from the blockchain

These settings affect performance and disk usage.

`
}

func getServicesHelp(questionText string) string {
	// Service-specific help
	if strings.Contains(questionText, "api") {
		return `API Service

The API service provides a REST API for accessing TrueBlocks data.
Enable this if you want to:
- Access TrueBlocks data from other applications
- Use the TrueBlocks explorer interface
- Integrate with external tools

The API runs on port 8080 by default.

`
	}

	if strings.Contains(questionText, "scraper") {
		return `Scraper Service

The scraper continuously indexes the blockchain to provide fast access to:
- Transaction histories for accounts
- Traces and logs
- Other blockchain data

This service should be enabled if you want to use TrueBlocks for
historical analysis or monitoring.

`
	}

	// Default services help
	return `Services Configuration

TrueBlocks offers several services that can be enabled:

- API: REST API for accessing TrueBlocks data
- Scraper: Indexes the blockchain for fast data access
- Monitor: Watches accounts for new transactions
- IPFS: Provides decentralized storage for shared index data

Enable only the services you need to conserve resources.

`
}

func getPortsHelp(questionText string) string {
	_ = questionText
	return `Service Ports Configuration

Each enabled service needs a unique port to listen on:

- API: Default port 8080
- Scraper: Default port 8081
- Monitor: Default port 8082
- IPFS: Default port 8083

Requirements:
- Ports must be between 1024-65535
- Ports must not be in use by other applications
- Ensure your firewall allows traffic if needed

`
}

func getChainsHelp(_ string) string {
	// Chain-specific help
	return `Blockchain Networks Configuration

Configure which blockchain networks you want to use with TrueBlocks:

- Enable/disable specific networks
- Configure RPC endpoints for each network
- Set chain-specific parameters

TrueBlocks supports multiple Ethereum-compatible networks including
Ethereum mainnet, testnets, and L2 solutions.

`
}

func getLoggingHelp(questionText string) string {
	_ = questionText
	return `Logging Configuration

Configure how TrueBlocks logs information:

- Log to File: Whether to save logs to disk
- Log Folder: Where to store log files
- Log Level: How verbose the logging should be
  * Error: Only errors
  * Warning: Errors and warnings
  * Info: General information (recommended)
  * Debug: Detailed information for troubleshooting

Logs are useful for troubleshooting issues with TrueBlocks.

`
}

func getSummaryHelp() string {
	return `Configuration Summary

This screen provides a complete overview of your configuration:

- Review all settings before finalizing
- Use "b" or "back" to return to previous screens to make changes
- Type "save" to save this configuration as a template
- Press Enter to save the configuration and exit the wizard

After saving, you can start using TrueBlocks with your new configuration.

`
}
