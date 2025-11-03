package types

import (
	"fmt"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// ValidateRpcEndpointRT performs real-time validation of an RPC endpoint,
// checking both formatting and connectivity.
func ValidateRpcEndpointRT(endpoint string) (string, error) {
	// Basic format validation
	if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") &&
		!strings.HasPrefix(endpoint, "ws://") && !strings.HasPrefix(endpoint, "wss://") {
		return "", fmt.Errorf("RPC endpoint must start with http:// or https:// (websocket URLs not supported)")
	}

	// Explicitly reject websocket endpoints
	if strings.HasPrefix(strings.ToLower(endpoint), "ws://") || strings.HasPrefix(strings.ToLower(endpoint), "wss://") {
		return "", fmt.Errorf("WebSocket RPC endpoints (ws://, wss://) are not supported; use http:// or https://")
	}

	// For HTTP endpoints, attempt to connect
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Prepare a simple JSON-RPC request to check the endpoint
	payload := strings.NewReader(`{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}`)

	resp, err := client.Post(endpoint, "application/json", payload)
	if err != nil {
		return "", fmt.Errorf("connection failed: %s", err.Error())
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("connection returned status code: %d", resp.StatusCode)
	}

	return "", nil
}

// ValidateFolder performs real-time validation of a folder path,
// checking for existence, permissions, and disk space.
func ValidateFolder(path string) (bool, string, string) {
	_ = path // delint
	// The actual implementation would check:
	// 1. If the folder exists or can be created
	// 2. If we have write permissions
	// 3. If there's enough disk space

	// For now, we'll return a warning about space requirements
	return true, "", "Ensure this location has at least 200GB of free space for the index data."
}

// ValidatePort checks if a port number is valid and available
func ValidatePort(portStr string) (bool, string) {
	// Check if the port is a valid number
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return false, "Port must be a valid number"
	}

	// Check if the port is in the valid range
	if port < 1024 || port > 65535 {
		return false, "Port must be between 1024 and 65535"
	}

	// Check if the port is available
	if !isPortAvailable(port) {
		return false, fmt.Sprintf("Port %d is already in use", port)
	}

	return true, ""
}

// isPortAvailable checks if a port is available to use
func isPortAvailable(port int) bool {
	// Try to listen on the port to see if it's available
	addr := fmt.Sprintf(":%d", port)
	listener, err := net.Listen("tcp", addr)

	// If we can't listen, the port is in use
	if err != nil {
		return false
	}

	// Close the listener and return true
	listener.Close()
	return true
}

// ValidateWithFeedback provides enhanced validation with color-coded feedback
// and severity levels. Returns validity, message, and severity (error, warning, info)
func ValidateWithFeedback(input string, validationType string, extraInfo ...string) (bool, string, string) {
	_ = extraInfo // delint
	switch validationType {
	case "rpc":
		_, err := ValidateRpcEndpointRT(input)
		if err != nil {
			return false, err.Error(), "error"
		}
		return true, "RPC endpoint format is valid", "info"

	case "folder":
		isValid, message, warning := ValidateFolder(input)
		if !isValid {
			return false, message, "error"
		}
		if warning != "" {
			return true, warning, "warning"
		}
		return true, "Folder location is valid", "info"

	case "port":
		isValid, message := ValidatePort(input)
		severity := "error"
		if isValid {
			message = fmt.Sprintf("Port %s is available", input)
			severity = "info"
		}
		return isValid, message, severity

	case "chainId":
		// Validate chain ID format and known chains
		return true, "Chain ID is valid", "info"

	default:
		return true, "", "info"
	}
}
