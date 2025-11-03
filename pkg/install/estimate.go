package install

// EstimateIndex returns rough (order-of-magnitude) estimates for disk (GB) and initial hours
// for a given acquisition strategy and detail level. These are intentionally conservative
// and should be refined with telemetry later.
func EstimateIndex(strategy, detail string) (diskGB int, hours int) {
	// New semantics:
	// detail = index  -> full index + blooms (heavy)
	// detail = bloom     -> blooms only (light)
	// Sizes provided by user: index=168GB, bloom-only=6GB.
	if detail == "blooms" { // backward compatibility normalization
		detail = "bloom"
	}
	full := detail == "index"
	switch strategy {
	case "download":
		if full {
			return 168, 3 // heavy download incl verify (rough)
		}
		return 6, 1 // light download
	case "scratch":
		if full {
			return 168, 100 // building full index locally (very rough placeholder)
		}
		return 6, 12 // building blooms only takes time but far less disk
	default:
		return 0, 0
	}
}
