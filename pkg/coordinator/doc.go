// Package coordinator provides coordination between services
//
// The scraper-monitor coordinator implements an event-driven model where
// the scraper service notifies monitors when new blocks have been scraped.
// This replaces the previous independent polling model where monitors ran
// on their own schedule regardless of scraper activity.
//
// Architecture:
//
//	Scraper -> Channel -> Coordinator -> Monitor
//
// The channel is passed from khedra down through the SDK to the core scraper,
// enabling type-safe, zero-overhead in-process communication.
package coordinator
