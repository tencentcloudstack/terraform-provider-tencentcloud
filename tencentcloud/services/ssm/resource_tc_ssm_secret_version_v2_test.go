package ssm

import (
	"context"
	"testing"

	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/types"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
)

// TestSsmSecretVersionV2Resource_Metadata validates the resource type name.
func TestSsmSecretVersionV2Resource_Metadata(t *testing.T) {
	r := &SsmSecretVersionV2Resource{}
	req := resource.MetadataRequest{}
	resp := &resource.MetadataResponse{}

	r.Metadata(context.Background(), req, resp)

	if resp.TypeName != "tencentcloud_ssm_secret_version_v2" {
		t.Errorf("Metadata.TypeName = %q, want %q", resp.TypeName, "tencentcloud_ssm_secret_version_v2")
	}
}

// TestSsmSecretVersionV2Resource_Schema validates required / optional / computed
// flags on every attribute, including the WriteOnly flag on secret_string_wo.
func TestSsmSecretVersionV2Resource_Schema(t *testing.T) {
	r := &SsmSecretVersionV2Resource{}
	req := resource.SchemaRequest{}
	resp := &resource.SchemaResponse{}

	r.Schema(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema returned diagnostics: %v", resp.Diagnostics)
	}

	type expect struct {
		required, optional, computed, sensitive, writeOnly bool
	}
	cases := map[string]expect{
		"id":               {computed: true},
		"secret_name":      {required: true},
		"version_id":       {required: true},
		"secret_binary":    {optional: true, sensitive: true},
		"secret_string":    {optional: true, sensitive: true},
		"secret_string_wo": {optional: true, sensitive: true, writeOnly: true},
	}

	for name, want := range cases {
		attr, ok := resp.Schema.Attributes[name].(schema.StringAttribute)
		if !ok {
			t.Errorf("attribute %q missing or not a StringAttribute", name)
			continue
		}
		if attr.Required != want.required {
			t.Errorf("attribute %q Required=%v, want %v", name, attr.Required, want.required)
		}
		if attr.Optional != want.optional {
			t.Errorf("attribute %q Optional=%v, want %v", name, attr.Optional, want.optional)
		}
		if attr.Computed != want.computed {
			t.Errorf("attribute %q Computed=%v, want %v", name, attr.Computed, want.computed)
		}
		if attr.Sensitive != want.sensitive {
			t.Errorf("attribute %q Sensitive=%v, want %v", name, attr.Sensitive, want.sensitive)
		}
		if attr.WriteOnly != want.writeOnly {
			t.Errorf("attribute %q WriteOnly=%v, want %v", name, attr.WriteOnly, want.writeOnly)
		}
	}
}

// TestSsmSecretVersionV2Resource_ImplementsInterfaces verifies the resource
// satisfies the framework interfaces it advertises in package-level _ asserts.
func TestSsmSecretVersionV2Resource_ImplementsInterfaces(t *testing.T) {
	var _ resource.Resource = &SsmSecretVersionV2Resource{}
	var _ resource.ResourceWithConfigure = &SsmSecretVersionV2Resource{}
	var _ resource.ResourceWithImportState = &SsmSecretVersionV2Resource{}
}

// TestSsmSecretVersionV2Resource_Factory validates the registry factory.
func TestSsmSecretVersionV2Resource_Factory(t *testing.T) {
	got := NewSsmSecretVersionV2Resource()
	if got == nil {
		t.Fatal("NewSsmSecretVersionV2Resource returned nil")
	}
	if _, ok := got.(*SsmSecretVersionV2Resource); !ok {
		t.Fatalf("factory returned %T, want *SsmSecretVersionV2Resource", got)
	}
}

// TestSsmSecretVersionV2Resource_ConfigureNilProviderData ensures the
// pre-Configure invocation pattern (nil ProviderData) is silently no-op.
func TestSsmSecretVersionV2Resource_ConfigureNilProviderData(t *testing.T) {
	r := &SsmSecretVersionV2Resource{}
	req := resource.ConfigureRequest{}
	resp := &resource.ConfigureResponse{}

	r.Configure(context.Background(), req, resp)

	if resp.Diagnostics.HasError() {
		t.Fatalf("Configure with nil ProviderData returned diagnostics: %v", resp.Diagnostics)
	}
}

// TestSsmSecretVersionV2Resource_ConfigureWrongType ensures Configure adds an
// error diagnostic when ProviderData is of an unexpected type.
func TestSsmSecretVersionV2Resource_ConfigureWrongType(t *testing.T) {
	r := &SsmSecretVersionV2Resource{}
	req := resource.ConfigureRequest{ProviderData: "not-a-meta"}
	resp := &resource.ConfigureResponse{}

	r.Configure(context.Background(), req, resp)

	if !resp.Diagnostics.HasError() {
		t.Error("expected error diagnostic for wrong ProviderData type")
	}
}

// TestSplitSsmVersionV2Id validates parsing of the `<secret_name>#<version_id>` id.
func TestSplitSsmVersionV2Id(t *testing.T) {
	good := "my-secret" + tccommon.FILED_SP + "v1"
	name, ver, ok := splitSsmVersionV2Id(good)
	if !ok || name != "my-secret" || ver != "v1" {
		t.Errorf("splitSsmVersionV2Id(%q) = (%q, %q, %v); want (\"my-secret\", \"v1\", true)", good, name, ver, ok)
	}

	bad := []string{
		"",
		"only-name",
		tccommon.FILED_SP + "v1",
		"name" + tccommon.FILED_SP,
		"a" + tccommon.FILED_SP + "b" + tccommon.FILED_SP + "c",
	}
	for _, in := range bad {
		if _, _, ok := splitSsmVersionV2Id(in); ok {
			t.Errorf("splitSsmVersionV2Id(%q) ok=true, want false", in)
		}
	}
}

// TestStringValueOrNull verifies "" maps to Null and non-empty strings round-trip.
func TestStringValueOrNull(t *testing.T) {
	if v := stringValueOrNull(""); !v.IsNull() {
		t.Errorf("stringValueOrNull(\"\") = %v, want Null", v)
	}
	if v := stringValueOrNull("abc"); v.IsNull() || v.ValueString() != "abc" {
		t.Errorf("stringValueOrNull(\"abc\") = %v, want \"abc\"", v)
	}
}

// TestStringRequiresReplace validates the custom ForceNew-equivalent plan
// modifier across create / destroy / unchanged / changed scenarios.
func TestStringRequiresReplace(t *testing.T) {
	type tc struct {
		name      string
		state     types.String
		plan      types.String
		wantForce bool
	}
	cases := []tc{
		{"create:state-null", types.StringNull(), types.StringValue("v1"), false},
		{"destroy:plan-null", types.StringValue("v1"), types.StringNull(), false},
		{"unchanged", types.StringValue("v1"), types.StringValue("v1"), false},
		{"changed", types.StringValue("v1"), types.StringValue("v2"), true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			req := planmodifier.StringRequest{StateValue: c.state, PlanValue: c.plan}
			resp := &planmodifier.StringResponse{PlanValue: c.plan}
			stringRequiresReplace{}.PlanModifyString(context.Background(), req, resp)
			if resp.RequiresReplace != c.wantForce {
				t.Errorf("RequiresReplace = %v, want %v", resp.RequiresReplace, c.wantForce)
			}
		})
	}
}

// TestStringUseStateForUnknown validates that an unknown plan value is
// replaced by the prior state value, while known / null state is left alone.
func TestStringUseStateForUnknown(t *testing.T) {
	t.Run("state-null:no-op", func(t *testing.T) {
		req := planmodifier.StringRequest{StateValue: types.StringNull(), PlanValue: types.StringUnknown()}
		resp := &planmodifier.StringResponse{PlanValue: req.PlanValue}
		stringUseStateForUnknown{}.PlanModifyString(context.Background(), req, resp)
		if !resp.PlanValue.IsUnknown() {
			t.Errorf("PlanValue = %v, want Unknown", resp.PlanValue)
		}
	})
	t.Run("plan-known:no-op", func(t *testing.T) {
		req := planmodifier.StringRequest{StateValue: types.StringValue("prev"), PlanValue: types.StringValue("new")}
		resp := &planmodifier.StringResponse{PlanValue: req.PlanValue}
		stringUseStateForUnknown{}.PlanModifyString(context.Background(), req, resp)
		if resp.PlanValue.ValueString() != "new" {
			t.Errorf("PlanValue = %q, want \"new\"", resp.PlanValue.ValueString())
		}
	})
	t.Run("unknown-plan:copied-from-state", func(t *testing.T) {
		req := planmodifier.StringRequest{StateValue: types.StringValue("prev"), PlanValue: types.StringUnknown()}
		resp := &planmodifier.StringResponse{PlanValue: req.PlanValue}
		stringUseStateForUnknown{}.PlanModifyString(context.Background(), req, resp)
		if resp.PlanValue.ValueString() != "prev" {
			t.Errorf("PlanValue = %v, want \"prev\"", resp.PlanValue)
		}
	})
}

// TestSsmSecretVersionV2Resource_BuildSecretValueParam covers the
// "exactly one of" rule and the wire-mapping of secret_string_wo.
func TestSsmSecretVersionV2Resource_BuildSecretValueParam(t *testing.T) {
	r := &SsmSecretVersionV2Resource{}

	// none set -> error
	if _, err := r.buildSecretValueParam("n", "v",
		types.StringNull(), types.StringNull(), types.StringNull()); err == nil {
		t.Error("expected error when no secret value field is set")
	}

	// two set -> error
	if _, err := r.buildSecretValueParam("n", "v",
		types.StringValue("b"), types.StringValue("s"), types.StringNull()); err == nil {
		t.Error("expected error when two secret value fields are set")
	}

	// secret_binary only
	p, err := r.buildSecretValueParam("n", "v",
		types.StringValue("YmluCg=="), types.StringNull(), types.StringNull())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p["secret_binary"] != "YmluCg==" || p["secret_name"] != "n" || p["version_id"] != "v" {
		t.Errorf("secret_binary param mismatch: %+v", p)
	}
	if _, ok := p["secret_string"]; ok {
		t.Errorf("secret_string should not be present, got %+v", p)
	}

	// secret_string only
	p, err = r.buildSecretValueParam("n", "v",
		types.StringNull(), types.StringValue("plain"), types.StringNull())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p["secret_string"] != "plain" {
		t.Errorf("secret_string param mismatch: %+v", p)
	}

	// secret_string_wo only -> mapped onto secret_string on the wire
	p, err = r.buildSecretValueParam("n", "v",
		types.StringNull(), types.StringNull(), types.StringValue("wo"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p["secret_string"] != "wo" {
		t.Errorf("secret_string_wo should map to secret_string, got %+v", p)
	}
	if _, ok := p["secret_string_wo"]; ok {
		t.Errorf("secret_string_wo should not be in param, got %+v", p)
	}
}
