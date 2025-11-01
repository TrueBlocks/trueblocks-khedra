package install

import (
	"testing"

	"github.com/TrueBlocks/trueblocks-khedra/v6/pkg/types"
)

func newDraft() *Draft {
	return &Draft{Config: types.NewConfig(), Meta: DraftMeta{Schema: DraftSchema}}
}

// ApplyDraft tests - removed due to config directory permission issues
// TODO: Re-add these tests with proper test environment setup

// ValidateIndex tests
func TestValidateIndex_InvalidStrategy(t *testing.T) {
	d := newDraft()
	d.Config.General.Strategy = "weird"
	ferrs := ValidateDraftPhase(d, "step:index")
	if !hasCode(ferrs, "invalid_strategy") {
		t.Fatalf("expected invalid_strategy, got %+v", ferrs)
	}
}

func TestValidateIndex_InvalidDetail(t *testing.T) {
	d := newDraft()
	d.Config.General.Detail = "other"
	ferrs := ValidateDraftPhase(d, "step:index")
	if !hasCode(ferrs, "invalid_detail") {
		t.Fatalf("expected invalid_detail, got %+v", ferrs)
	}
}

// ValidateChains tests
func TestValidateChains_MainnetMissing(t *testing.T) {
	d := newDraft()
	delete(d.Config.Chains, "mainnet")
	ferrs := ValidateDraftPhase(d, "step:chains")
	if !hasCode(ferrs, "require_mainnet") {
		t.Fatalf("expected require_mainnet, got %+v", ferrs)
	}
}

func TestValidateChains_MainnetEnabledNoRPC(t *testing.T) {
	d := newDraft()
	ch := d.Config.Chains["mainnet"]
	ch.RPCs = []string{""}
	d.Config.Chains["mainnet"] = ch
	ferrs := ValidateDraftPhase(d, "step:chains")
	if !hasCode(ferrs, "mainnet_missing_rpc") {
		t.Fatalf("expected mainnet_missing_rpc, got %+v", ferrs)
	}
}

func TestValidateChains_InvalidRPCScheme(t *testing.T) {
	d := newDraft()
	ch := d.Config.Chains["mainnet"]
	ch.RPCs = []string{"ws://localhost:8545"}
	d.Config.Chains["mainnet"] = ch
	ferrs := ValidateDraftPhase(d, "step:chains")
	if !hasCode(ferrs, "invalid_rpc_scheme") {
		t.Fatalf("expected invalid_rpc_scheme, got %+v", ferrs)
	}
}

// ValidateServices tests
func TestValidateServices_PortConflict(t *testing.T) {
	d := newDraft()
	api := d.Config.Services["api"]
	ipfs := d.Config.Services["ipfs"]
	ipfs.Port = api.Port // conflict
	d.Config.Services["ipfs"] = ipfs
	ferrs := ValidateDraftPhase(d, "step:services")
	if !hasCode(ferrs, "port_conflict") {
		t.Fatalf("expected port_conflict, got %+v", ferrs)
	}
}

func TestValidateServices_NoEnabledOK(t *testing.T) {
	d := newDraft()
	for k, svc := range d.Config.Services {
		svc.Enabled = false
		d.Config.Services[k] = svc
	}
	ferrs := ValidateDraftPhase(d, "step:services")
	if len(filterCodes(ferrs, "port_conflict")) > 0 {
		t.Fatalf("unexpected port_conflict: %+v", ferrs)
	}
}

// ValidateLogging tests
func TestValidateLogging_FileLoggingNoFolder(t *testing.T) {
	d := newDraft()
	d.Config.Logging.ToFile = true
	d.Config.Logging.Filename = "k.log"
	d.Config.Logging.Folder = "" // missing
	ferrs := ValidateDraftPhase(d, "step:logging")
	if !hasCode(ferrs, "log_folder_required") {
		t.Fatalf("expected log_folder_required, got %+v", ferrs)
	}
}

// Helper functions
func hasCode(ferrs []FieldError, code string) bool { return len(filterCodes(ferrs, code)) > 0 }
func filterCodes(ferrs []FieldError, code string) []FieldError {
	var out []FieldError
	for _, fe := range ferrs {
		if fe.Code == code {
			out = append(out, fe)
		}
	}
	return out
}
