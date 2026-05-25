package metalists

import "testing"

// TestRegionEntries_NotEmptyAndPopulated checks that the placeholder
// region list contains at least 5 rows (matching the proposal) and that
// every row has a non-empty ID and Name.
//
// This is the minimal regression for the L0 placeholder: it ensures the
// stub file is not silently emptied or broken by future edits.
func TestRegionEntries_NotEmptyAndPopulated(t *testing.T) {
	got := RegionEntries()

	if len(got) < 5 {
		t.Fatalf("RegionEntries() returned %d entries, want at least 5", len(got))
	}
	for i, entry := range got {
		if entry.ID == "" {
			t.Errorf("entry[%d].ID is empty", i)
		}
		if entry.Name == "" {
			t.Errorf("entry[%d].Name is empty", i)
		}
	}
}

// TestRegionEntries_ReturnsDefensiveCopy verifies that RegionEntries
// returns a copy rather than the underlying slice; external mutation
// must not affect subsequent calls.
func TestRegionEntries_ReturnsDefensiveCopy(t *testing.T) {
	first := RegionEntries()
	if len(first) == 0 {
		t.Skip("no entries to mutate")
	}

	// Tamper with the returned value.
	first[0].ID = "TAMPERED"

	second := RegionEntries()
	if second[0].ID == "TAMPERED" {
		t.Fatalf("RegionEntries returned shared slice; expected defensive copy")
	}
}
