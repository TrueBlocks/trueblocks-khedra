package app

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"strconv"
	"time"
)

// handleWaitForNode waits for a specified node process to start before proceeding.
// It uses two environment variables:
// - TB_KHEDRA_WAIT_FOR_NODE: Name of the node process to wait for (e.g., "erigon", "geth")
// - TB_KHEDRA_WAIT_SECONDS: Number of seconds to wait for stabilization (default: 30)
//
// The function polls every 2 seconds using pgrep to check if the process is running.
// Once found, it waits for the specified stabilization period, counting down in
// 3-second intervals to provide user feedback.
func (k *KhedraApp) handleWaitForNode() error {
	nodeName := os.Getenv("TB_KHEDRA_WAIT_FOR_NODE")
	if nodeName == "" {
		return nil
	}

	k.logger.Info(fmt.Sprintf("Waiting for node process '%s' to start...\n", nodeName))

	// Wait for the node process to start
	for {
		cmd := exec.Command("pgrep", "-f", nodeName)
		err := cmd.Run()
		if err == nil {
			// Process found
			break
		}
		log.Print(".")
		time.Sleep(2 * time.Second)
	}

	// Get wait time from environment variable, default to 30 seconds
	waitSeconds := 30
	if waitSecondsEnv := os.Getenv("TB_KHEDRA_WAIT_SECONDS"); waitSecondsEnv != "" {
		if parsed, err := strconv.Atoi(waitSecondsEnv); err == nil && parsed > 0 {
			waitSeconds = parsed
		}
	}

	spins := waitSeconds / 3
	wait := 3
	k.logger.Info(fmt.Sprintf("\nNode '%s' detected. Waiting %d seconds for stabilization...\n", nodeName, spins*wait))
	for i := spins; i > 0; i-- {
		k.logger.Info(fmt.Sprintf("Stabilizing... %2d seconds remaining", i*wait))
		time.Sleep(time.Duration(wait) * time.Second)
	}
	k.logger.Info("Ready to proceed.")

	return nil
}
