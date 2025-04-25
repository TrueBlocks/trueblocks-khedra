package app

import (
	"fmt"
	"strconv"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-khedra/v5/pkg/types"
	"github.com/TrueBlocks/trueblocks-khedra/v5/pkg/wizard"
)

// getServicePortsScreen creates a configuration screen for service ports
func getServicePortsScreen() wizard.Screen {
	title := `Service Port Configuration`
	subtitle := ``
	instructions := `Enter port numbers and press enter.`
	body := `
Each enabled service needs a unique port to listen on. The ports must not
conflict with other applications on your system.

Default ports:
- API: 8080
- Scraper: 8081
- Monitor: 8082
- IPFS: 8083

Ports will be validated to ensure they are available.
`
	replacements := []wizard.Replacement{
		{Color: colors.Yellow, Values: []string{title}},
		{Color: colors.BrightBlue, Values: []string{
			"API:", "Scraper:", "Monitor:", "IPFS:",
		}},
	}

	// Create questions for each service port configuration
	var questions []wizard.Questioner

	// We'll add questions dynamically based on which services are enabled
	apiPortQuestion := createPortQuestion("api")
	scraperPortQuestion := createPortQuestion("scraper")
	monitorPortQuestion := createPortQuestion("monitor")
	ipfsPortQuestion := createPortQuestion("ipfs")

	questions = append(questions,
		&apiPortQuestion,
		&scraperPortQuestion,
		&monitorPortQuestion,
		&ipfsPortQuestion)

	style := wizard.NewStyle()

	return wizard.Screen{
		Title:        title,
		Subtitle:     subtitle,
		Instructions: instructions,
		Body:         body,
		Replacements: replacements,
		Questions:    questions,
		Style:        style,
	}
}

// createPortQuestion creates a question for configuring a service port
func createPortQuestion(serviceName string) wizard.Question {
	return wizard.Question{
		Question: fmt.Sprintf(`Enter port for the "%s" service:`, serviceName),
		Hint: fmt.Sprintf(`The port must be a number between 1024 and 65535|that is not in use. Default port for %s is %s.`,
			serviceName, getDefaultPort(serviceName)),
		ValidationType: "port", // For real-time port validation
		PrepareFn: func(input string, q *wizard.Question) (string, error) {
			if cfg, ok := q.Screen.Wizard.Backing.(*types.Config); ok {
				service, exists := cfg.Services[serviceName]
				if !exists || !service.Enabled {
					// Skip this question if service is disabled
					return "", validSkipNext()
				}

				// Return current port or default
				return strconv.Itoa(service.Port), nil
			}
			return getDefaultPort(serviceName), nil
		},
		Validate: func(input string, q *wizard.Question) (string, error) {
			if cfg, ok := q.Screen.Wizard.Backing.(*types.Config); ok {
				// Try to parse port number
				port, err := strconv.Atoi(input)
				if err != nil {
					return input, fmt.Errorf("port must be a number %w", wizard.ErrValidate)
				}

				// Validate port range
				if port < 1024 || port > 65535 {
					return input, fmt.Errorf("port must be between 1024 and 65535 %w", wizard.ErrValidate)
				}

				// Update service configuration
				service := cfg.Services[serviceName]
				service.Port = port
				cfg.Services[serviceName] = service

				// Save configuration
				if err := cfg.WriteToFile(types.GetConfigFnNoCreate()); err != nil {
					return input, fmt.Errorf("failed to save configuration: %s %w", err.Error(), wizard.ErrValidate)
				}

				return "", validOk("Port configuration saved successfully", "")
			}
			return input, fmt.Errorf("failed to access configuration %w", wizard.ErrValidate)
		},
	}
}

// getDefaultPort returns the default port for a service
func getDefaultPort(serviceName string) string {
	switch serviceName {
	case "api":
		return "8080"
	case "scraper":
		return "8081"
	case "monitor":
		return "8082"
	case "ipfs":
		return "8083"
	default:
		return "8080"
	}
}
