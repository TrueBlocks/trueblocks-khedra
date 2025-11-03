package types

import (
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strings"
	"time"

	"github.com/TrueBlocks/trueblocks-core/src/apps/chifra/v6/pkg/rpc"
)

// RpcTestResult contains the result of an RPC endpoint connectivity test
type RpcTestResult struct {
	Reachable    bool
	ResponseTime time.Duration
	BlockNumber  string
	ChainID      string
	ErrorMessage string
}

// TestRpcEndpoint performs comprehensive testing of an RPC endpoint
// This includes checking if it's reachable, measuring response time,
// and verifying it responds to basic Ethereum JSON-RPC methods
func TestRpcEndpoint(endpoint string) RpcTestResult {
	result := RpcTestResult{
		Reachable: false,
	}

	// Skip testing for non-HTTP endpoints
	if !strings.HasPrefix(endpoint, "http://") && !strings.HasPrefix(endpoint, "https://") {
		if strings.HasPrefix(endpoint, "ws://") || strings.HasPrefix(endpoint, "wss://") {
			result.Reachable = true
			result.ErrorMessage = "WebSocket endpoints cannot be fully tested"
		} else {
			result.ErrorMessage = "Invalid endpoint URL format"
		}
		return result
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	// Test 1: eth_blockNumber - checks basic connectivity and gets current block
	startTime := time.Now()
	blockNumberPayload := strings.NewReader(`{"jsonrpc":"2.0","method":"eth_blockNumber","params":[],"id":1}`)
	resp, err := client.Post(endpoint, "application/json", blockNumberPayload)

	if err != nil {
		result.ErrorMessage = fmt.Sprintf("Connection failed: %s", err.Error())
		return result
	}

	responseTime := time.Since(startTime)
	result.ResponseTime = responseTime

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		result.ErrorMessage = fmt.Sprintf("Connection returned status code: %d", resp.StatusCode)
		return result
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		result.ErrorMessage = fmt.Sprintf("Failed to read response: %s", err.Error())
		return result
	}

	// Parse the response
	var blockNumberResponse struct {
		Result string `json:"result"`
		Error  *struct {
			Message string `json:"message"`
		} `json:"error"`
	}

	if err := json.Unmarshal(body, &blockNumberResponse); err != nil {
		result.ErrorMessage = fmt.Sprintf("Failed to parse response: %s", err.Error())
		return result
	}

	if blockNumberResponse.Error != nil {
		result.ErrorMessage = fmt.Sprintf("RPC error: %s", blockNumberResponse.Error.Message)
		return result
	}

	result.BlockNumber = blockNumberResponse.Result

	// Test 2: eth_chainId - gets the chain ID to verify correct network
	chainIDPayload := strings.NewReader(`{"jsonrpc":"2.0","method":"eth_chainId","params":[],"id":1}`)
	chainResp, err := client.Post(endpoint, "application/json", chainIDPayload)

	if err == nil && chainResp.StatusCode == http.StatusOK {
		defer chainResp.Body.Close()
		chainBody, err := io.ReadAll(chainResp.Body)

		if err == nil {
			var chainIDResponse struct {
				Result string `json:"result"`
			}

			if json.Unmarshal(chainBody, &chainIDResponse) == nil {
				result.ChainID = chainIDResponse.Result
			}
		}
	}

	// If we got this far, the endpoint is reachable
	result.Reachable = true
	return result
}

// FormatRpcTestResult formats an RPC test result as a user-friendly string
func FormatRpcTestResult(result RpcTestResult) string {
	if !result.Reachable {
		return fmt.Sprintf("❌ RPC endpoint unreachable: %s", result.ErrorMessage)
	}

	responseTimeMs := result.ResponseTime.Milliseconds()

	output := fmt.Sprintf("✅ RPC endpoint is reachable (response time: %dms)\n", responseTimeMs)

	if result.BlockNumber != "" {
		blockNum, _ := ParseHexNumber(result.BlockNumber)
		output += fmt.Sprintf("📦 Current block: %s\n", blockNum)
	}

	if result.ChainID != "" {
		chainID, _ := ParseHexNumber(result.ChainID)
		output += fmt.Sprintf("🔗 Chain ID: %s", chainID)
	}

	return output
}

// ParseHexNumber converts a hex string to a decimal string
func ParseHexNumber(hexString string) (string, error) {
	// Remove "0x" prefix if present
	hexString = strings.TrimPrefix(hexString, "0x")

	// Parse hex string
	var number int64
	_, err := fmt.Sscanf(hexString, "%x", &number)
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%d", number), nil
}

func TryConnect(chain, providerUrl string, maxAttempts int) error {
	for i := 1; i <= maxAttempts; i++ {
		_, err := rpc.PingRpc(providerUrl)
		if err == nil {
			return nil
		} else {
			slog.Warn("retrying RPC", "chain", chain, "provider", providerUrl)
			if i < maxAttempts {
				time.Sleep(1 * time.Second)
			}
		}
	}

	fv := NewFieldValidator("ping_rpc", "Chain", "rpc", fmt.Sprintf("[%s]", chain))
	return Failed(fv, fmt.Sprintf("cannot connect to RPC (%s-%s) after %d attempts", chain, providerUrl, maxAttempts), "")
}
