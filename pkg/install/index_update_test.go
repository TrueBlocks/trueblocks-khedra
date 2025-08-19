package install

import "testing"

func TestUpdateIndexStrategy(t *testing.T) {
	d := NewDraftFromConfig("test-session")

	// case 1: download + blooms only (detail 'bloom')
	UpdateIndexStrategy(d, "download", "bloom")
	if d.Config.General.Strategy != "download" || d.Config.General.Detail != "bloom" {
		t.Fatalf("strategy/detail not set: %+v", d.Config.General)
	}
	if d.Meta.EstDiskGB != 6 || d.Meta.EstHours != 1 { // expectation from new EstimateIndex
		t.Fatalf("unexpected estimates (download,bloom): disk=%d hours=%d", d.Meta.EstDiskGB, d.Meta.EstHours)
	}

	// case 2: scratch + full index (detail 'index')
	UpdateIndexStrategy(d, "scratch", "index")
	if d.Meta.EstDiskGB != 168 || d.Meta.EstHours != 100 {
		t.Fatalf("unexpected estimates (scratch,index): disk=%d hours=%d", d.Meta.EstDiskGB, d.Meta.EstHours)
	}
}
