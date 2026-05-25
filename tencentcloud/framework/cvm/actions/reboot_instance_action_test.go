package cvmactions

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/action"
)

// TestRebootInstanceAction_Metadata verifies that the type name equals
// "<provider>_reboot_instance".
func TestRebootInstanceAction_Metadata(t *testing.T) {
	a := NewRebootInstanceAction()

	var resp action.MetadataResponse
	a.Metadata(context.Background(), action.MetadataRequest{ProviderTypeName: "tencentcloud"}, &resp)

	if got, want := resp.TypeName, "tencentcloud_reboot_instance"; got != want {
		t.Fatalf("Metadata.TypeName = %q, want %q", got, want)
	}
}

// TestRebootInstanceAction_Schema verifies that the schema's attribute set
// matches the proposal: instance_id (Required) + force (Optional).
func TestRebootInstanceAction_Schema(t *testing.T) {
	a := NewRebootInstanceAction()

	var resp action.SchemaResponse
	a.Schema(context.Background(), action.SchemaRequest{}, &resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema returned errors: %v", resp.Diagnostics)
	}

	wantAttrs := map[string]struct{}{
		"instance_id": {},
		"force":       {},
	}
	for name := range resp.Schema.Attributes {
		if _, ok := wantAttrs[name]; !ok {
			t.Errorf("unexpected attribute %q in schema", name)
			continue
		}
		delete(wantAttrs, name)
	}
	if len(wantAttrs) > 0 {
		missing := make([]string, 0, len(wantAttrs))
		for k := range wantAttrs {
			missing = append(missing, k)
		}
		t.Errorf("missing attributes in schema: %v", missing)
	}
}

// TestInstanceIDPattern_AcceptsValidIDs exercises the instance_id regex
// against both legal and illegal ids to confirm the validation behaviour.
//
// We test the regex directly (rather than driving the framework runtime
// chain via Invoke) to avoid the boilerplate of constructing a tfsdk.Config.
// End-to-end coverage of Invoke is provided by acceptance tests, which are
// out of scope for this change.
func TestInstanceIDPattern_AcceptsValidIDs(t *testing.T) {
	t.Parallel()

	cases := []struct {
		id    string
		valid bool
	}{
		{"ins-abc", true},
		{"ins-123", true},
		{"ins-abc123def", true},
		{"ins-", false},     // missing body
		{"ins-ABC", false},  // contains uppercase
		{"INS-abc", false},  // uppercase prefix
		{"ins-abc!", false}, // special character
		{"i-abc", false},    // wrong prefix (VPC style)
		{"", false},         // empty string
	}

	for _, tc := range cases {
		tc := tc
		t.Run(tc.id, func(t *testing.T) {
			t.Parallel()
			got := instanceIDPattern.MatchString(tc.id)
			if got != tc.valid {
				t.Errorf("instanceIDPattern.MatchString(%q) = %v, want %v", tc.id, got, tc.valid)
			}
		})
	}
}
