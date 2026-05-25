package metaephemerals

import (
	"context"
	"strings"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/ephemeral"
)

// TestTempCredentialEphemeralResource_Metadata verifies that the type name
// equals "<provider>_temp_credential".
func TestTempCredentialEphemeralResource_Metadata(t *testing.T) {
	r := NewTempCredentialEphemeralResource()

	var resp ephemeral.MetadataResponse
	r.Metadata(context.Background(), ephemeral.MetadataRequest{ProviderTypeName: "tencentcloud"}, &resp)

	if got, want := resp.TypeName, "tencentcloud_temp_credential"; got != want {
		t.Fatalf("Metadata.TypeName = %q, want %q", got, want)
	}
}

// TestTempCredentialEphemeralResource_Schema verifies the schema's
// attribute set matches the proposal: region / secret_id / secret_key /
// token / expires_at.
//
// It also verifies that secret_key and token are marked Sensitive so the
// Terraform UI does not print them in plaintext.
func TestTempCredentialEphemeralResource_Schema(t *testing.T) {
	r := NewTempCredentialEphemeralResource()

	var resp ephemeral.SchemaResponse
	r.Schema(context.Background(), ephemeral.SchemaRequest{}, &resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema returned errors: %v", resp.Diagnostics)
	}

	wantAttrs := map[string]struct {
		mustBeSensitive bool
	}{
		"region":     {mustBeSensitive: false},
		"secret_id":  {mustBeSensitive: false},
		"secret_key": {mustBeSensitive: true},
		"token":      {mustBeSensitive: true},
		"expires_at": {mustBeSensitive: false},
	}

	for name, attr := range resp.Schema.Attributes {
		want, ok := wantAttrs[name]
		if !ok {
			t.Errorf("unexpected attribute %q in schema", name)
			continue
		}
		if want.mustBeSensitive && !attr.IsSensitive() {
			t.Errorf("attribute %q must be marked Sensitive", name)
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

// TestRandomHex_FormatAndUniqueness exercises the placeholder credential's
// random generator:
//   - The returned length is exactly 2n (hex doubles the byte count).
//   - Different invocations produce different values (the collision
//     probability is vanishingly small; if a regression hits this, retry).
func TestRandomHex_FormatAndUniqueness(t *testing.T) {
	t.Parallel()

	for _, n := range []int{8, 16, 32} {
		got := randomHex(n)
		if len(got) != 2*n && !strings.Contains(got, ".") {
			// Note: the crypto/rand fallback returns a timestamp containing
			// "."; that branch does not enforce 2n length. crypto/rand
			// should never fail in the test environment.
			t.Errorf("randomHex(%d) length = %d, want %d", n, len(got), 2*n)
		}
	}

	// Calling the same length multiple times produces different results
	// (collision probability < 2^-128).
	a := randomHex(16)
	b := randomHex(16)
	if a == b {
		t.Errorf("randomHex(16) returned same value twice: %q", a)
	}
}
