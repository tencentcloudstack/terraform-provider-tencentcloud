package cam

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-framework/types/basetypes"

	camv20190116 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cam/v20190116"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

// --- Schema / Metadata / Factory ---

// TestCamUpdatePasswordRuleConfigResource_Metadata validates the type name.
func TestCamUpdatePasswordRuleConfigResource_Metadata(t *testing.T) {
	r := &CamUpdatePasswordRuleConfigResource{}
	req := resource.MetadataRequest{}
	resp := &resource.MetadataResponse{}

	r.Metadata(context.Background(), req, resp)

	if resp.TypeName != "tencentcloud_cam_update_password_rule_config" {
		t.Errorf("Metadata.TypeName = %q, want %q", resp.TypeName, "tencentcloud_cam_update_password_rule_config")
	}
}

// TestCamUpdatePasswordRuleConfigResource_Schema validates required/optional/computed
// flags on every top-level attribute and on the nested rules attributes.
func TestCamUpdatePasswordRuleConfigResource_Schema(t *testing.T) {
	r := &CamUpdatePasswordRuleConfigResource{}
	req := resource.SchemaRequest{}
	resp := &resource.SchemaResponse{}

	r.Schema(context.Background(), req, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Schema returned diagnostics: %v", resp.Diagnostics)
	}

	type expect struct {
		required, optional, computed bool
	}
	cases := map[string]expect{
		"id":          {computed: true},
		"rules":       {required: true},
		"update_time": {computed: true},
		"modifier":    {computed: true},
	}
	for name, want := range cases {
		attr, ok := resp.Schema.Attributes[name]
		if !ok {
			t.Errorf("attribute %q missing from schema", name)
			continue
		}
		switch a := attr.(type) {
		case schema.StringAttribute:
			if a.Required != want.required || a.Optional != want.optional || a.Computed != want.computed {
				t.Errorf("attribute %q (Required=%v, Optional=%v, Computed=%v), want %v", name, a.Required, a.Optional, a.Computed, want)
			}
		case schema.SingleNestedAttribute:
			if a.Required != want.required || a.Optional != want.optional || a.Computed != want.computed {
				t.Errorf("attribute %q (Required=%v, Optional=%v, Computed=%v), want %v", name, a.Required, a.Optional, a.Computed, want)
			}
		default:
			t.Errorf("attribute %q unexpected type %T", name, attr)
		}
	}

	rules, ok := resp.Schema.Attributes["rules"].(schema.SingleNestedAttribute)
	if !ok {
		t.Fatal("rules not a SingleNestedAttribute")
	}
	ruleCases := map[string]expect{
		"minimum_length":                   {optional: true, computed: true},
		"must_contain":                     {optional: true, computed: true},
		"force_password_change":            {optional: true, computed: true},
		"reuse_password_limit":             {optional: true, computed: true},
		"retry_password_limit":             {optional: true, computed: true},
		"password_expiration_invalidation": {optional: true, computed: true},
	}
	for name, want := range ruleCases {
		attr, ok := rules.Attributes[name]
		if !ok {
			t.Errorf("nested attribute %q missing from rules schema", name)
			continue
		}
		switch a := attr.(type) {
		case schema.Int64Attribute:
			if a.Optional != want.optional || a.Computed != want.computed {
				t.Errorf("nested attribute %q (Optional=%v, Computed=%v), want %v", name, a.Optional, a.Computed, want)
			}
		case schema.StringAttribute:
			if a.Optional != want.optional || a.Computed != want.computed {
				t.Errorf("nested attribute %q (Optional=%v, Computed=%v), want %v", name, a.Optional, a.Computed, want)
			}
		default:
			t.Errorf("nested attribute %q unexpected type %T", name, attr)
		}
	}
}

// TestCamUpdatePasswordRuleConfigResource_ImplementsInterfaces verifies the
// resource satisfies the framework interfaces advertised by the var _ asserts.
func TestCamUpdatePasswordRuleConfigResource_ImplementsInterfaces(t *testing.T) {
	var _ resource.ResourceWithConfigure = &CamUpdatePasswordRuleConfigResource{}
	var _ resource.ResourceWithImportState = &CamUpdatePasswordRuleConfigResource{}
}

// TestCamUpdatePasswordRuleConfigResource_Factory validates the registry factory.
func TestCamUpdatePasswordRuleConfigResource_Factory(t *testing.T) {
	got := NewCamUpdatePasswordRuleConfigResource()
	if got == nil {
		t.Fatal("NewCamUpdatePasswordRuleConfigResource returned nil")
	}
	if _, ok := got.(*CamUpdatePasswordRuleConfigResource); !ok {
		t.Fatalf("factory returned %T, want *CamUpdatePasswordRuleConfigResource", got)
	}
}

// TestCamUpdatePasswordRuleConfigResource_ConfigureNilProviderData ensures the
// pre-Configure invocation pattern (nil ProviderData) is a silent no-op.
func TestCamUpdatePasswordRuleConfigResource_ConfigureNilProviderData(t *testing.T) {
	r := &CamUpdatePasswordRuleConfigResource{}
	req := resource.ConfigureRequest{}
	resp := &resource.ConfigureResponse{}

	r.Configure(context.Background(), req, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Configure with nil ProviderData returned diagnostics: %v", resp.Diagnostics)
	}
}

// --- passwordRulesTouched ---

// TestPasswordRulesTouched verifies that untouched (all-zero) defaults are not
// treated as an existing configuration.
func TestPasswordRulesTouched(t *testing.T) {
	cases := []struct {
		name   string
		params *camv20190116.GetPasswordRulesResponseParams
		want   bool
	}{
		{"nil params", nil, false},
		{"nil update time", &camv20190116.GetPasswordRulesResponseParams{Rules: &camv20190116.PassWordRule{MinimumLength: helper.Int64(8)}}, false},
		{
			"empty update time",
			&camv20190116.GetPasswordRulesResponseParams{
				Rules:      &camv20190116.PassWordRule{MinimumLength: helper.Int64(8), MustContain: helper.String("Aa1!")},
				UpdateTime: helper.String(""),
			},
			false,
		},
		{
			// A live unmodified account reports non-zero rule defaults with the
			// zero-sentinel timestamp; this must NOT be treated as touched.
			"zero sentinel update time with non-zero defaults",
			&camv20190116.GetPasswordRulesResponseParams{
				Rules: &camv20190116.PassWordRule{
					MinimumLength:       helper.Int64(8),
					MustContain:         helper.String("Aa1!"),
					RetryPasswordLimit:  helper.Int64(10),
					ReusePasswordLimit:  helper.Int64(1),
					ForcePasswordChange: helper.Int64(0),
				},
				UpdateTime: helper.String("0000-00-00 00:00:00"),
			},
			false,
		},
		{
			"real update time",
			&camv20190116.GetPasswordRulesResponseParams{
				Rules:      &camv20190116.PassWordRule{MinimumLength: helper.Int64(8)},
				UpdateTime: helper.String("2026-11-10 15:37:23"),
				Modifier:   helper.String("100053131847"),
			},
			true,
		},
		{
			"real update time with zero rules",
			&camv20190116.GetPasswordRulesResponseParams{
				Rules:      &camv20190116.PassWordRule{},
				UpdateTime: helper.String("2026-11-10 15:37:23"),
			},
			true,
		},
	}
	for _, c := range cases {
		c := c
		t.Run(c.name, func(t *testing.T) {
			if got := passwordRulesTouched(c.params); got != c.want {
				t.Fatalf("passwordRulesTouched = %v, want %v", got, c.want)
			}
		})
	}
}

// --- expand / flatten ---

// buildRulesObject constructs a types.Object rules value for tests.
func buildRulesObject(ctx context.Context, t *testing.T) types.Object {
	t.Helper()
	model := CamUpdatePasswordRuleConfigRuleModel{
		MinimumLength:                  types.Int64Value(8),
		MustContain:                    types.StringValue("char&num"),
		ForcePasswordChange:            types.Int64Value(90),
		ReusePasswordLimit:             types.Int64Value(5),
		RetryPasswordLimit:             types.Int64Value(10),
		PasswordExpirationInvalidation: types.Int64Value(1),
	}
	obj, diags := types.ObjectValueFrom(ctx, passwordRulesAttrTypes(), model)
	if diags.HasError() {
		t.Fatalf("build rules object failed: %v", diags)
	}
	return obj
}

// TestExpandPasswordRules verifies the rules object -> SDK struct conversion
// covers all leaf fields.
func TestExpandPasswordRules(t *testing.T) {
	ctx := context.Background()
	obj := buildRulesObject(ctx, t)

	var diags diag.Diagnostics
	rule := expandPasswordRules(ctx, obj, &diags)
	if diags.HasError() {
		t.Fatalf("expandPasswordRules returned diagnostics: %v", diags)
	}
	if rule.MinimumLength == nil || *rule.MinimumLength != 8 {
		t.Errorf("MinimumLength = %v, want 8", rule.MinimumLength)
	}
	if rule.MustContain == nil || *rule.MustContain != "char&num" {
		t.Errorf("MustContain = %v, want char&num", rule.MustContain)
	}
	if rule.ForcePasswordChange == nil || *rule.ForcePasswordChange != 90 {
		t.Errorf("ForcePasswordChange = %v, want 90", rule.ForcePasswordChange)
	}
	if rule.ReusePasswordLimit == nil || *rule.ReusePasswordLimit != 5 {
		t.Errorf("ReusePasswordLimit = %v, want 5", rule.ReusePasswordLimit)
	}
	if rule.RetryPasswordLimit == nil || *rule.RetryPasswordLimit != 10 {
		t.Errorf("RetryPasswordLimit = %v, want 10", rule.RetryPasswordLimit)
	}
	if rule.PasswordExpirationInvalidation == nil || *rule.PasswordExpirationInvalidation != 1 {
		t.Errorf("PasswordExpirationInvalidation = %v, want 1", rule.PasswordExpirationInvalidation)
	}
}

// TestExpandPasswordRules_Null verifies a null object yields an empty rule
// struct without diagnostics.
func TestExpandPasswordRules_Null(t *testing.T) {
	var diags diag.Diagnostics
	rule := expandPasswordRules(context.Background(), types.ObjectNull(passwordRulesAttrTypes()), &diags)
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if rule == nil {
		t.Fatal("expected non-nil empty rule for null object")
	}
	if rule.MinimumLength != nil || rule.MustContain != nil {
		t.Errorf("expected empty rule, got %+v", rule)
	}
}

// TestFlattenPasswordRules verifies the SDK struct -> rules object conversion
// reads back all leaf fields.
func TestFlattenPasswordRules(t *testing.T) {
	ctx := context.Background()
	prev := types.ObjectNull(passwordRulesAttrTypes())

	rules := &camv20190116.PassWordRule{
		MinimumLength:                  helper.Int64(10),
		MustContain:                    helper.String("char"),
		ForcePasswordChange:            helper.Int64(180),
		ReusePasswordLimit:             helper.Int64(3),
		RetryPasswordLimit:             helper.Int64(5),
		PasswordExpirationInvalidation: helper.Int64(2),
	}

	obj := flattenPasswordRules(ctx, rules, prev, &diag.Diagnostics{})
	var out CamUpdatePasswordRuleConfigRuleModel
	if diags := obj.As(ctx, &out, basetypes.ObjectAsOptions{}); diags.HasError() {
		t.Fatalf("object as failed: %v", diags)
	}
	if out.MinimumLength.ValueInt64() != 10 {
		t.Errorf("MinimumLength = %v, want 10", out.MinimumLength)
	}
	if out.MustContain.ValueString() != "char" {
		t.Errorf("MustContain = %q, want char", out.MustContain.ValueString())
	}
	if out.ForcePasswordChange.ValueInt64() != 180 {
		t.Errorf("ForcePasswordChange = %v, want 180", out.ForcePasswordChange)
	}
	if out.ReusePasswordLimit.ValueInt64() != 3 {
		t.Errorf("ReusePasswordLimit = %v, want 3", out.ReusePasswordLimit)
	}
	if out.RetryPasswordLimit.ValueInt64() != 5 {
		t.Errorf("RetryPasswordLimit = %v, want 5", out.RetryPasswordLimit)
	}
	if out.PasswordExpirationInvalidation.ValueInt64() != 2 {
		t.Errorf("PasswordExpirationInvalidation = %v, want 2", out.PasswordExpirationInvalidation)
	}
}

// TestFlattenPasswordRules_NilRules verifies a nil SDK struct yields an all-null
// rules object (no diagnostics).
func TestFlattenPasswordRules_NilRules(t *testing.T) {
	ctx := context.Background()
	obj := flattenPasswordRules(ctx, nil, types.ObjectNull(passwordRulesAttrTypes()), &diag.Diagnostics{})
	if obj.IsNull() {
		t.Fatal("expected a non-null object with null attributes")
	}
	var out CamUpdatePasswordRuleConfigRuleModel
	if diags := obj.As(ctx, &out, basetypes.ObjectAsOptions{}); diags.HasError() {
		t.Fatalf("object as failed: %v", diags)
	}
	if !out.MinimumLength.IsNull() || !out.MustContain.IsNull() {
		t.Errorf("expected all-null attributes, got %+v", out)
	}
}

// TestFlattenPasswordRules_PreservesPrevious verifies that a field the cloud
// omits keeps the previously configured value.
func TestFlattenPasswordRules_PreservesPrevious(t *testing.T) {
	ctx := context.Background()
	prev := buildRulesObject(ctx, t)

	// The cloud only echoes back minimum_length.
	rules := &camv20190116.PassWordRule{MinimumLength: helper.Int64(12)}

	obj := flattenPasswordRules(ctx, rules, prev, &diag.Diagnostics{})
	var out CamUpdatePasswordRuleConfigRuleModel
	if diags := obj.As(ctx, &out, basetypes.ObjectAsOptions{}); diags.HasError() {
		t.Fatalf("object as failed: %v", diags)
	}
	if out.MinimumLength.ValueInt64() != 12 {
		t.Errorf("MinimumLength = %v, want 12 (cloud value)", out.MinimumLength)
	}
	if out.MustContain.ValueString() != "char&num" {
		t.Errorf("MustContain = %q, want char&num (previous value)", out.MustContain.ValueString())
	}
	if out.ForcePasswordChange.ValueInt64() != 90 {
		t.Errorf("ForcePasswordChange = %v, want 90 (previous value)", out.ForcePasswordChange)
	}
}

// --- service layer ---

// mockPasswordRulesClient patches the CAM client accessor so service calls hit
// an empty SDK client whose API methods can be stubbed per-test.
func mockPasswordRulesClient(t *testing.T) (*gomonkey.Patches, CamService, *camv20190116.Client) {
	t.Helper()
	patches := gomonkey.NewPatches()
	client := &connectivity.TencentCloudClient{Region: "ap-guangzhou"}
	camClient := &camv20190116.Client{}
	patches.ApplyMethodReturn(client, "UseCamClient", camClient)
	service := NewCamService(client)
	return patches, service, camClient
}

// TestServiceGetPasswordRules covers the read path, echoing the response params.
func TestServiceGetPasswordRules(t *testing.T) {
	patches, service, camClient := mockPasswordRulesClient(t)
	defer patches.Reset()

	patches.ApplyMethodFunc(camClient, "GetPasswordRulesWithContext",
		func(ctx context.Context, request *camv20190116.GetPasswordRulesRequest) (*camv20190116.GetPasswordRulesResponse, error) {
			resp := camv20190116.NewGetPasswordRulesResponse()
			resp.Response = &camv20190116.GetPasswordRulesResponseParams{
				Rules: &camv20190116.PassWordRule{
					MinimumLength: helper.Int64(8),
					MustContain:   helper.String("char&num"),
				},
				UpdateTime: helper.String("2024-01-01 00:00:00"),
				Modifier:   helper.String("100000000001"),
			}
			return resp, nil
		})

	params, err := service.GetPasswordRules(context.Background())
	if err != nil {
		t.Fatalf("GetPasswordRules failed: %v", err)
	}
	if params == nil || params.Rules == nil || params.Rules.MinimumLength == nil || *params.Rules.MinimumLength != 8 {
		t.Errorf("params = %+v, want Rules.MinimumLength 8", params)
	}
}

// TestServiceGetPasswordRules_NoRules asserts a response without rules yields a
// distinguishable error (errPasswordRulesNotFound).
func TestServiceGetPasswordRules_NoRules(t *testing.T) {
	patches, service, camClient := mockPasswordRulesClient(t)
	defer patches.Reset()

	patches.ApplyMethodFunc(camClient, "GetPasswordRulesWithContext",
		func(ctx context.Context, request *camv20190116.GetPasswordRulesRequest) (*camv20190116.GetPasswordRulesResponse, error) {
			resp := camv20190116.NewGetPasswordRulesResponse()
			resp.Response = &camv20190116.GetPasswordRulesResponseParams{}
			return resp, nil
		})

	params, err := service.GetPasswordRules(context.Background())
	if err == nil {
		t.Fatal("expected error for a response without rules")
	}
	if params != nil {
		t.Errorf("expected nil params, got %+v", params)
	}
}

// TestServiceUpdatePasswordRules covers the write path and the request echo.
func TestServiceUpdatePasswordRules(t *testing.T) {
	patches, service, camClient := mockPasswordRulesClient(t)
	defer patches.Reset()

	var captured *camv20190116.UpdatePasswordRulesRequest
	patches.ApplyMethodFunc(camClient, "UpdatePasswordRulesWithContext",
		func(ctx context.Context, request *camv20190116.UpdatePasswordRulesRequest) (*camv20190116.UpdatePasswordRulesResponse, error) {
			captured = request
			resp := camv20190116.NewUpdatePasswordRulesResponse()
			resp.Response = &camv20190116.UpdatePasswordRulesResponseParams{}
			return resp, nil
		})

	rules := &camv20190116.PassWordRule{MinimumLength: helper.Int64(8)}
	if err := service.UpdatePasswordRules(context.Background(), rules); err != nil {
		t.Fatalf("UpdatePasswordRules failed: %v", err)
	}
	if captured == nil {
		t.Fatal("UpdatePasswordRulesWithContext was not called")
	}
	if captured.Rules == nil || captured.Rules.MinimumLength == nil || *captured.Rules.MinimumLength != 8 {
		t.Errorf("captured Rules = %+v, want MinimumLength 8", captured.Rules)
	}
}

// TestServiceUpdatePasswordRules_NilResponse asserts a nil response is treated
// as an API failure.
func TestServiceUpdatePasswordRules_NilResponse(t *testing.T) {
	patches, service, camClient := mockPasswordRulesClient(t)
	defer patches.Reset()

	patches.ApplyMethodReturn(camClient, "UpdatePasswordRulesWithContext", (*camv20190116.UpdatePasswordRulesResponse)(nil), nil)

	err := service.UpdatePasswordRules(context.Background(), &camv20190116.PassWordRule{})
	if err == nil {
		t.Fatal("expected error for nil response, got nil")
	}
}
