package coordinator

import (
	"log/slog"

	"github.com/TrueBlocks/trueblocks-sdk/v6/services"
	"github.com/TrueBlocks/trueblocks-sdk/v6/services/monitor"
)

// ScraperCompletedEvent is sent when the scraper completes processing a batch
// This is a "wake up" signal - monitors track their own block state
type ScraperCompletedEvent = services.ScrapeCompletedEvent

// ScraperMonitorCoordinator coordinates between scraper and monitor services
// When the scraper completes a batch, it invokes HandleEvent callback
// The coordinator then triggers the appropriate monitor to process those blocks
type ScraperMonitorCoordinator struct {
	monitors map[string]*monitor.MonitorService // keyed by chain
	logger   *slog.Logger
}

// NewScraperMonitorCoordinator creates a new coordinator
func NewScraperMonitorCoordinator(logger *slog.Logger) *ScraperMonitorCoordinator {
	return &ScraperMonitorCoordinator{
		monitors: make(map[string]*monitor.MonitorService),
		logger:   logger,
	}
}

// RegisterMonitor registers a monitor service for a specific chain
func (c *ScraperMonitorCoordinator) RegisterMonitor(chain string, mon *monitor.MonitorService) {
	c.monitors[chain] = mon
	c.logger.Info("Monitor registered for chain", "chain", chain)
}

// HandleEvent processes a scraper completion event (used as callback)
func (c *ScraperMonitorCoordinator) HandleEvent(event ScraperCompletedEvent) {
	mon, ok := c.monitors[event.Chain]
	if !ok {
		c.logger.Warn("Monitor NOT triggered - no monitor registered for chain", "chain", event.Chain, "registeredChains", len(c.monitors))
		return
	}

	c.logger.Info("Scraper completed, triggering monitor", "chain", event.Chain)

	// Trigger monitor to process newly scraped blocks (monitor determines range from its state)
	if err := mon.ProcessMonitors(event.Chain); err != nil {
		c.logger.Error("Failed to trigger monitor", "chain", event.Chain, "error", err)
	}
}
