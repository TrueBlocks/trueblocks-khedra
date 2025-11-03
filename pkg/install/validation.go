package install

import (
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/TrueBlocks/trueblocks-khedra/v6/pkg/types"
)

// FieldError identifies a specific invalid field; empty Field means global error.
type FieldError struct {
	Field   string `json:"field"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ValidationError aggregates errors (placeholder for future expansion).
type ValidationError struct {
	Errors []FieldError `json:"errors"`
}

// normalizeDraft trims and applies legacy conversions prior to validation.
func normalizeDraft(d *Draft) {
	if d == nil {
		return
	}
	g := d.Config.General
	g.DataFolder = strings.TrimSpace(g.DataFolder)
	g.Strategy = strings.TrimSpace(g.Strategy)
	g.Detail = strings.TrimSpace(g.Detail)
	if g.Detail == "blooms" { // legacy
		g.Detail = "bloom"
	}
	d.Config.General = g
	// Trim first RPC of each chain (others if present)
	for name, ch := range d.Config.Chains {
		for i, rpc := range ch.RPCs {
			ch.RPCs[i] = strings.TrimSpace(rpc)
		}
		d.Config.Chains[name] = ch
	}
	// Trim logging folder/filename
	lg := d.Config.Logging
	lg.Folder = strings.TrimSpace(lg.Folder)
	lg.Filename = strings.TrimSpace(lg.Filename)
	d.Config.Logging = lg
}

// Individual pure validators -------------------------------------------------

func ValidatePaths(d *Draft) []FieldError {
	var out []FieldError
	if d == nil || strings.TrimSpace(d.Config.General.DataFolder) == "" {
		out = append(out, FieldError{Field: "general.dataFolder", Code: "required", Message: "data folder is required"})
	}
	return out
}

func ValidateIndex(d *Draft) []FieldError {
	var out []FieldError
	strategy := d.Config.General.Strategy
	if strategy == "" {
		out = append(out, FieldError{Field: "general.strategy", Code: "required", Message: "strategy is required"})
	} else if strategy != "download" && strategy != "scratch" {
		out = append(out, FieldError{Field: "general.strategy", Code: "invalid_strategy", Message: fmt.Sprintf("unknown strategy '%s'", strategy)})
	}
	det := d.Config.General.Detail
	if det == "" { // allow empty -> will default later, no error
	} else if det != "index" && det != "bloom" {
		out = append(out, FieldError{Field: "general.detail", Code: "invalid_detail", Message: fmt.Sprintf("unknown detail '%s'", det)})
	}
	return out
}

var chainNameRe = regexp.MustCompile(`^[a-zA-Z0-9-_]+$`)

func ValidateChains(d *Draft) []FieldError {
	var out []FieldError
	if d == nil {
		return []FieldError{{Field: "", Code: "internal", Message: "nil draft"}}
	}
	hasMainnet := false
	for name, ch := range d.Config.Chains {
		// name validation (always)
		if !chainNameRe.MatchString(name) {
			out = append(out, FieldError{Field: fmt.Sprintf("chains.%s.name", name), Code: "invalid_name", Message: fmt.Sprintf("invalid chain name '%s'", name)})
		}
		if name == "mainnet" || ch.ChainID == 1 {
			hasMainnet = true
			if ch.Enabled {
				if len(ch.RPCs) == 0 || strings.TrimSpace(firstRPC(ch.RPCs)) == "" {
					out = append(out, FieldError{Field: "chains.mainnet.rpc", Code: "mainnet_missing_rpc", Message: "mainnet enabled but RPC missing"})
				} else if !validRPCScheme(firstRPC(ch.RPCs)) {
					out = append(out, FieldError{Field: "chains.mainnet.rpc", Code: "invalid_rpc_scheme", Message: "mainnet RPC must start with http(s)://"})
				}
			} else {
				out = append(out, FieldError{Field: "chains.mainnet.rpc", Code: "require_mainnet", Message: "mainnet must be enabled"})
			}
		} else if ch.Enabled { // non-mainnet enabled
			if len(ch.RPCs) == 0 || strings.TrimSpace(firstRPC(ch.RPCs)) == "" {
				out = append(out, FieldError{Field: fmt.Sprintf("chains.%s.rpc", name), Code: "missing_rpc", Message: fmt.Sprintf("chain %s enabled but RPC missing", name)})
			} else if !validRPCScheme(firstRPC(ch.RPCs)) {
				out = append(out, FieldError{Field: fmt.Sprintf("chains.%s.rpc", name), Code: "invalid_rpc_scheme", Message: fmt.Sprintf("chain %s RPC must start with http(s)://", name)})
			}
		}
	}
	if !hasMainnet { // mainnet chain block absent entirely
		out = append(out, FieldError{Field: "chains.mainnet.rpc", Code: "require_mainnet", Message: "mainnet chain definition required"})
	}
	return out
}

func firstRPC(rpcs []string) string {
	if len(rpcs) == 0 {
		return ""
	}
	return rpcs[0]
}

func validRPCScheme(u string) bool {
	return strings.HasPrefix(u, "http://") || strings.HasPrefix(u, "https://")
}

func ValidateServices(d *Draft) []FieldError {
	var out []FieldError

	// Check that at least one service is enabled
	hasEnabledService := false
	for _, svc := range d.Config.Services {
		if svc.Enabled {
			hasEnabledService = true
			break
		}
	}
	if !hasEnabledService {
		out = append(out, FieldError{Field: "services", Code: "no_services_enabled", Message: "At least one service must be enabled"})
	}

	// Check port conflicts for enabled services
	seen := map[int]string{}
	for name, svc := range d.Config.Services {
		if !svc.Enabled || svc.Port == 0 { // only ports for tcp services (api/ipfs)
			continue
		}
		if svc.Port < 1024 || svc.Port > 65535 {
			out = append(out, FieldError{Field: fmt.Sprintf("services.%s.port", name), Code: "port_out_of_range", Message: fmt.Sprintf("service %s port %d out of range (1024-65535)", name, svc.Port)})
			continue
		}
		if other, ok := seen[svc.Port]; ok {
			out = append(out, FieldError{Field: fmt.Sprintf("services.%s.port", name), Code: "port_conflict", Message: fmt.Sprintf("port %d reused by %s and %s", svc.Port, other, name)})
		} else {
			seen[svc.Port] = name
		}
	}
	return out
}

func ValidateLogging(d *Draft) []FieldError {
	var out []FieldError
	lg := d.Config.Logging
	lvl := lg.Level
	if lvl != "" && lvl != "info" && lvl != "debug" && lvl != "warn" && lvl != "error" {
		out = append(out, FieldError{Field: "logging.level", Code: "invalid_level", Message: fmt.Sprintf("unknown logging level '%s'", lvl)})
	}
	if lg.ToFile {
		if strings.TrimSpace(lg.Folder) == "" {
			out = append(out, FieldError{Field: "logging.folder", Code: "log_folder_required", Message: "folder required when file logging enabled"})
		}
		if strings.TrimSpace(lg.Filename) == "" {
			out = append(out, FieldError{Field: "logging.filename", Code: "log_filename_required", Message: "filename required when file logging enabled"})
		}
	}
	return out
}

// ValidateDraftPhase validates a draft for a given phase (step:* or final) returning FieldErrors.
func ValidateDraftPhase(d *Draft, phase string) []FieldError {
	normalizeDraft(d)
	var all []FieldError
	// Simple validation - only validate current screen when moving forward
	switch phase {
	case "step:paths":
		all = append(all, ValidatePaths(d)...)
	case "step:chains":
		all = append(all, ValidateChains(d)...)
	case "step:index":
		all = append(all, ValidateIndex(d)...)
	case "step:services":
		all = append(all, ValidateServices(d)...)
	case "step:logging":
		all = append(all, ValidateLogging(d)...)
	case "final":
		// Final validation - check EVERYTHING before writing real config
		all = append(all, ValidatePaths(d)...)
		all = append(all, ValidateChains(d)...)
		all = append(all, ValidateIndex(d)...)
		all = append(all, ValidateServices(d)...)
		all = append(all, ValidateLogging(d)...)
	default:
		// treat unknown as final
		all = append(all, ValidatePaths(d)...)
		all = append(all, ValidateIndex(d)...)
		all = append(all, ValidateChains(d)...)
		all = append(all, ValidateServices(d)...)
		all = append(all, ValidateLogging(d)...)
	}
	return all
}

// EnsureDataFolder attempts to create the data folder if it does not exist.
func EnsureDataFolder(path string) error {
	if strings.TrimSpace(path) == "" {
		return errors.New("empty data folder path")
	}
	if fi, err := os.Stat(path); err == nil {
		if !fi.IsDir() {
			return fmt.Errorf("data folder path exists and is not a directory: %s", path)
		}
		return nil
	}
	if err := os.MkdirAll(path, 0o700); err != nil { // create with 0700
		return err
	}
	return nil
}

// ApplyDraft validates the draft then writes it as the final config (config.yaml). It backs up
// any existing final config and removes the draft on success.
func ApplyDraft() error {
	d, err := LoadDraft()
	if err != nil {
		return fmt.Errorf("cannot load draft: %w", err)
	}
	ferrs := ValidateDraftPhase(d, "final")
	if len(ferrs) > 0 {
		msgs := make([]string, 0, len(ferrs))
		for _, fe := range ferrs {
			msgs = append(msgs, fe.Message)
		}
		return fmt.Errorf("draft invalid: %s", strings.Join(msgs, "; "))
	}
	if err := EnsureDataFolder(d.Config.General.DataFolder); err != nil {
		return err
	}
	finalPath := types.GetConfigFnNoCreate()
	finalDir := filepath.Dir(finalPath)
	if fi, err := os.Stat(finalDir); err != nil || !fi.IsDir() {
		if err := os.MkdirAll(finalDir, 0o700); err != nil {
			return err
		}
	}
	original := []byte{}
	if data, err := os.ReadFile(finalPath); err == nil {
		original = data
	}
	_, _ = BackupFinalConfig()
	tmp := finalPath + ".tmp-new"
	if err := d.Config.WriteToFile(tmp); err != nil {
		_ = os.Remove(tmp)
		return err
	}
	if f, err := os.Open(tmp); err == nil {
		_ = f.Sync()
		_ = f.Close()
	}
	if err := os.Rename(tmp, finalPath); err != nil {
		_ = os.Remove(tmp)
		if len(original) > 0 { // rollback attempt
			_ = os.WriteFile(finalPath, original, 0o600)
		}
		return err
	}
	if dirF, err := os.Open(filepath.Dir(finalPath)); err == nil { // Attempt dir sync (best effort)
		_ = dirF.Sync()
		_ = dirF.Close()
	}
	if fi, err := os.Stat(finalPath); err != nil || fi.Size() == 0 { // Verify non-empty final file
		if len(original) > 0 {
			_ = os.WriteFile(finalPath, original, 0o600)
		}
		return io.ErrUnexpectedEOF
	}
	if err := RemoveDraft(); err != nil { // Remove draft
		return err
	}
	return nil
}
