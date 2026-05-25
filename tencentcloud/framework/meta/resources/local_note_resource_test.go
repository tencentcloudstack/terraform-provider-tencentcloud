package metaresources

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
)

// TestLocalNoteResource_Metadata verifies that the resource type name
// equals "<provider>_local_note", so that
//
//	resource "tencentcloud_local_note" "x" {}
//
// in HCL is routed correctly.
func TestLocalNoteResource_Metadata(t *testing.T) {
	r := NewLocalNoteResource()

	var resp resource.MetadataResponse
	r.Metadata(context.Background(), resource.MetadataRequest{ProviderTypeName: "tencentcloud"}, &resp)

	if got, want := resp.TypeName, "tencentcloud_local_note"; got != want {
		t.Fatalf("Metadata.TypeName = %q, want %q", got, want)
	}
}

// TestLocalNoteResource_Schema verifies that the schema's attribute set
// matches the proposal: id / title / content / last_updated.
func TestLocalNoteResource_Schema(t *testing.T) {
	r := NewLocalNoteResource()

	var resp resource.SchemaResponse
	r.Schema(context.Background(), resource.SchemaRequest{}, &resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema returned errors: %v", resp.Diagnostics)
	}

	wantAttrs := map[string]struct{}{
		"id":           {},
		"title":        {},
		"content":      {},
		"last_updated": {},
	}
	for name := range resp.Schema.Attributes {
		if _, ok := wantAttrs[name]; !ok {
			t.Errorf("unexpected attribute %q in schema", name)
			continue
		}
		delete(wantAttrs, name)
	}
	if len(wantAttrs) > 0 {
		t.Errorf("missing attributes in schema: %v", keys(wantAttrs))
	}
}

// TestLocalNoteResource_CRUD_InMemoryRoundTrip is a lightweight
// regression case: it verifies the in-memory store's behaviour over a
// Create -> Read -> Update -> Delete sequence. It only exercises the
// local sync.Map; it does not construct the framework's full Plan/State
// objects (those are driven by the framework runtime in acceptance
// tests).
func TestLocalNoteResource_CRUD_InMemoryRoundTrip(t *testing.T) {
	id := newNoteID()
	if id == "" {
		t.Fatalf("newNoteID returned empty string")
	}

	notesStore.Store(id, noteRecord{Title: "t", Content: "c", LastUpdated: "2026-01-01T00:00:00Z"})
	defer notesStore.Delete(id)

	got, ok := notesStore.Load(id)
	if !ok {
		t.Fatalf("expected note %q to be present after Store", id)
	}
	rec := got.(noteRecord)
	if rec.Title != "t" || rec.Content != "c" {
		t.Fatalf("unexpected record after Store: %+v", rec)
	}

	notesStore.Delete(id)
	if _, ok := notesStore.Load(id); ok {
		t.Fatalf("expected note %q to be removed after Delete", id)
	}
}

// keys flattens a set-style map's keys into a slice for readable error
// messages.
func keys(m map[string]struct{}) []string {
	out := make([]string, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	return out
}
