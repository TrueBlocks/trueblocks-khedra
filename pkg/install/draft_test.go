package install

import (
	"os"
	"testing"
)

func TestLoadDraft_CorruptDetection(t *testing.T) {
	draftPath, _, _ := draftPaths()
	tmpFile := draftPath + ".testcorrupt"
	os.WriteFile(tmpFile, []byte("{not:valid,json}"), 0644)
	defer os.Remove(tmpFile)

	// Simulate corrupt file detection
	os.Rename(tmpFile, draftPath)
	defer os.Remove(draftPath)

	_, err := LoadDraft()
	if err == nil {
		t.Errorf("Expected error for corrupt draft, got nil")
	}
}

func TestSaveDraftAtomicAndLoadDraft_RoundTrip(t *testing.T) {
	draft := NewDraftFromConfig("testsession")
	err := SaveDraftAtomic(draft)
	if err != nil {
		t.Fatalf("SaveDraftAtomic failed: %v", err)
	}

	loaded, err := LoadDraft()
	if err != nil {
		t.Fatalf("LoadDraft failed: %v", err)
	}
	if loaded.Meta.Schema != DraftSchema || loaded.Meta.Session != "testsession" {
		t.Errorf("Loaded draft does not match saved draft")
	}
	_ = RemoveDraft()
}
