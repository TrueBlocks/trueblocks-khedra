package app

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/colors"
	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/pkg/utils"
	"github.com/urfave/cli/v2"
)

// pauseAction handles the pause command
func (k *KhedraApp) pauseAction(c *cli.Context) error {
	if c.NArg() < 1 {
		return fmt.Errorf("service name is required")
	}

	serviceName := c.Args().First()
	if !isValidServiceName(serviceName) {
		return fmt.Errorf("invalid service name '%s'. Valid services are: %s", serviceName, getValidServiceNames())
	}

	// Find the running khedra control service
	controlURL, err := findControlServiceURL()
	if err != nil {
		return err
	}

	// Handle "all" service - now supported directly by the API
	if serviceName == "all" {
		fmt.Printf("Pausing all pausable services...\n")
		result, err := callControlEndpoint(controlURL, "pause", "all")
		if err != nil {
			return fmt.Errorf("failed to pause all services: %w", err)
		}
		displayServiceControlResult("pause", result)
		return nil
	}

	// Call the pause endpoint for single service
	fmt.Printf("Pausing service '%s'...\n", serviceName)
	result, err := callControlEndpoint(controlURL, "pause", serviceName)
	if err != nil {
		return fmt.Errorf("failed to pause service: %w", err)
	}

	// Display the result
	displayServiceControlResult("pause", result)
	return nil
}

// unpauseAction handles the unpause command
func (k *KhedraApp) unpauseAction(c *cli.Context) error {
	if c.NArg() < 1 {
		return fmt.Errorf("service name is required")
	}

	serviceName := c.Args().First()
	if !isValidServiceName(serviceName) {
		return fmt.Errorf("invalid service name '%s'. Valid services are: %s", serviceName, getValidServiceNames())
	}

	// Find the running khedra control service
	controlURL, err := findControlServiceURL()
	if err != nil {
		return err
	}

	// Handle "all" service - now supported directly by the API
	if serviceName == "all" {
		fmt.Printf("Unpausing all pausable services...\n")
		result, err := callControlEndpoint(controlURL, "unpause", "all")
		if err != nil {
			return fmt.Errorf("failed to unpause all services: %w", err)
		}
		displayServiceControlResult("unpause", result)
		return nil
	}

	// Call the unpause endpoint for single service
	fmt.Printf("Unpausing service '%s'...\n", serviceName)
	result, err := callControlEndpoint(controlURL, "unpause", serviceName)
	if err != nil {
		return fmt.Errorf("failed to unpause service: %w", err)
	}

	// Display the result
	displayServiceControlResult("unpause", result)
	return nil
}

// findControlServiceURL finds the URL of the running control service
func findControlServiceURL() (string, error) {
	ports := []string{"8338", "8337", "8336", "8335"}
	for _, port := range ports {
		url := "http://localhost:" + port
		if utils.PingServer(url) {
			return url, nil
		}
	}
	return "", fmt.Errorf("khedra daemon is not running. Start it with 'khedra daemon'")
}

// callControlEndpoint makes an HTTP request to the control service
func callControlEndpoint(baseURL, endpoint, serviceName string) ([]map[string]string, error) {
	// Build the URL with service name parameter
	u, err := url.Parse(baseURL + "/" + endpoint)
	if err != nil {
		return nil, err
	}

	query := u.Query()
	query.Set("name", serviceName)
	u.RawQuery = query.Encode()

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Make the request
	resp, err := client.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("failed to connect to control service: %w", err)
	}
	defer resp.Body.Close()

	// Read the response
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("control service returned error: %s (status %d)", string(body), resp.StatusCode)
	}

	// Parse the JSON response
	var result []map[string]string
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result, nil
}

// isValidServiceName checks if the given service name is valid and pausable
func isValidServiceName(serviceName string) bool {
	validServices := strings.Split(strings.ReplaceAll(getValidServiceNames(), " ", ""), ",")
	for _, valid := range validServices {
		if serviceName == valid {
			return true
		}
	}
	return false
}

// getValidServiceNames returns a comma-separated list of valid service names
func getValidServiceNames() string {
	return "scraper, monitor, all"
}

// displayServiceControlResult displays the result of a pause/unpause operation
func displayServiceControlResult(action string, results []map[string]string) {
	if len(results) == 0 {
		fmt.Printf("No results returned for %s operation\n", action)
		return
	}

	for _, result := range results {
		name := result["name"]
		status := result["status"]

		// Color code the status
		var coloredStatus string
		if strings.Contains(status, "paused") || strings.Contains(status, "unpaused") {
			coloredStatus = colors.BrightGreen + status + colors.Off
		} else if strings.Contains(status, "already") {
			coloredStatus = colors.BrightBlue + status + colors.Off
		} else {
			coloredStatus = colors.BrightRed + status + colors.Off
		}

		fmt.Printf("Service '%s': %s\n", name, coloredStatus)
	}
}
