package install

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/TrueBlocks/trueblocks-khedra/v5/pkg/types"
)

func TestDetectRecentCorruption(t *testing.T) {
	// Create a fake corrupt file in temp config directory
	fn := types.GetConfigFnNoCreate()
	dir := filepath.Dir(fn)
	corrupt := filepath.Join(dir, "config.draft.json.corrupt-9999999999")
	if err := os.WriteFile(corrupt, []byte("{}"), 0o600); err != nil {
		t.Fatalf("write corrupt file: %v", err)
	}
	t.Cleanup(func() { _ = os.Remove(corrupt) })
	if !DetectRecentCorruption() {
		t.Fatalf("expected corruption detection true")
	}
}
