# Changes Needed in Core and SDK for Scraper-Monitor Coordination

## Overview

This document describes the changes needed in trueblocks-core and trueblocks-sdk to support channel-based coordination between scraper and monitor services.

## Core Changes (trueblocks-chifra)

### 1. Add Channel Support to ScrapeOptions

**File:** `chifra/internal/scrape/options.go`

```go
type ScrapeOptions struct {
    // ... existing fields ...
    
    // CompletionChan is an optional channel for scrape completion notifications.
    // If non-nil, an event is sent when each scrape run completes.
    // This enables in-process coordination with monitors or other services.
    CompletionChan chan<- ScrapeCompletedEvent
}

// ScrapeCompletedEvent is sent when a scrape run completes
type ScrapeCompletedEvent struct {
    Chain      string
    FirstBlock uint64
    LastBlock  uint64
    Meta       *types.MetaData
}
```

### 2. Send Events When Scrape Completes

**File:** `chifra/internal/scrape/scrape_batch.go` (or wherever ScrapeRunOnce finishes)

```go
func (opts *ScrapeOptions) ScrapeRunOnce() ([]types.Message, *types.MetaData, error) {
    // ... existing scrape logic ...
    
    // Send completion notification if channel provided
    if opts.CompletionChan != nil && meta != nil {
        event := ScrapeCompletedEvent{
            Chain:      opts.Globals.Chain,
            FirstBlock: meta.FirstBlock,
            LastBlock:  meta.LastBlock,
            Meta:       meta,
        }
        // Non-blocking send to prevent hanging scraper if channel full
        select {
        case opts.CompletionChan <- event:
        default:
            logger.Warn("Completion channel full, skipping notification")
        }
    }
    
    return msg, meta, nil
}
```

### 3. Deprecate Notify Package

**File:** `chifra/pkg/notify/notification.go`

Add deprecation notice at top:
```go
// Package notify provides HTTP-based notifications for scraper events.
//
// DEPRECATED: This package is deprecated in favor of channel-based
// coordination. Instead of HTTP notifications, pass a channel to
// ScrapeOptions.CompletionChan. This notify package will be removed
// in a future version.
package notify
```

**File:** `chifra/internal/scrape/options.go`

Deprecate the Notify flag:
```go
type ScrapeOptions struct {
    // ... other fields ...
    
    // Notify enables HTTP-based notifications (DEPRECATED)
    // DEPRECATED: Use CompletionChan instead for in-process coordination.
    // This flag will be removed in a future version.
    Notify bool `json:"notify,omitempty"`
}
```

## SDK Changes (trueblocks-sdk)

### 1. Add Channel Support to ScrapeService

**File:** `services/service_scraper.go`

```go
type ScrapeService struct {
    paused        bool
    logger        *slog.Logger
    initMode      string
    configTargets []string
    sleep         int
    blockCnt      int
    ctx           context.Context
    cancel        context.CancelFunc
    
    // completionChan enables coordination with monitors
    // When set, scraper sends events when it completes processing blocks
    completionChan chan<- ScrapeCompletedEvent
}

// ScrapeCompletedEvent is sent when scraper completes a batch
// This mirrors the type in chifra/internal/scrape
type ScrapeCompletedEvent struct {
    Chain      string
    FirstBlock uint64
    LastBlock  uint64
    Meta       interface{} // *types.MetaData from core
}

// SetCompletionChannel configures the scraper to send completion events
// Pass a channel to receive notifications when scrape batches complete
func (s *ScrapeService) SetCompletionChannel(ch chan<- ScrapeCompletedEvent) {
    s.completionChan = ch
}
```

### 2. Pass Channel Through to Core

**File:** `services/service_scraper.go`

Update `scrapeOneChain` to pass channel through:

```go
func (s *ScrapeService) scrapeOneChain(chain string) (*scraperReport, error) {
    // ... existing setup code ...
    
    opts := sdk.ScrapeOptions{
        BlockCnt: uint64(s.blockCnt),
        Globals: sdk.Globals{
            Chain: chain,
        },
    }
    
    // Pass completion channel through to core if configured
    if s.completionChan != nil {
        opts.CompletionChan = s.completionChan
    }
    
    // ... rest of existing code ...
}
```

### 3. Add ProcessBlockRange to MonitorService

**File:** `services/monitor/monitor_service.go` (or similar)

```go
// ProcessBlockRange triggers the monitor to process a specific block range
// This is called by the coordinator when the scraper completes a batch
func (m *MonitorService) ProcessBlockRange(chain string, firstBlock, lastBlock uint64) error {
    // Implementation will depend on monitor service structure
    // Should trigger monitor to export/process blocks for watched addresses
    return nil
}
```

## Testing Plan

1. **Core Tests**: Test that ScrapeRunOnce sends events when channel provided
2. **SDK Tests**: Test that SetCompletionChannel properly configures service
3. **Integration Tests**: Test khedra coordinator receives events and triggers monitors
4. **Backward Compatibility**: Test that existing notify-based code still works

## Migration Path

1. **Phase 1** (Current): Implement channel support alongside existing notify
2. **Phase 2** (Next version): Default to channels, warn if using notify
3. **Phase 3** (Future version): Remove notify package entirely

## Benefits

- **Zero Config**: No notify URL needed
- **Type Safe**: Compile-time guarantees
- **Zero Overhead**: In-process, no HTTP/JSON overhead  
- **Backpressure**: Buffered channel handles fast scraper
- **Clean Shutdown**: Close channel to stop coordination
- **Testable**: Easy to test with mock channels
