package config

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"

	configv20220802 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802"
	sdkerr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/helper"
)

// --- Schema / Metadata / Factory ---

// TestConfigAggregateCompliancePackResource_Metadata validates the type name.
func TestConfigAggregateCompliancePackResource_Metadata(t *testing.T) {
	r := &ConfigAggregateCompliancePackResource{}
	req := resource.MetadataRequest{}
	resp := &resource.MetadataResponse{}

	r.Metadata(context.Background(), req, resp)

	if resp.TypeName != "tencentcloud_config_aggregate_compliance_pack" {
		t.Errorf("Metadata.TypeName = %q, want %q", resp.TypeName, "tencentcloud_config_aggregate_compliance_pack")
	}
}

// TestConfigAggregateCompliancePackResource_Schema validates required/optional/computed
// flags on every top-level attribute.
func TestConfigAggregateCompliancePackResource_Schema(t *testing.T) {
	r := &ConfigAggregateCompliancePackResource{}
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
		"id":                  {computed: true},
		"compliance_pack_name": {required: true},
		"account_group_id":    {required: true},
		"risk_level":          {required: true},
		"description":         {optional: true},
		"config_rules":        {required: true},
		"status":              {computed: true},
		"compliance_pack_id": {computed: true},
		"create_time":         {computed: true},
	}

	for name, want := range cases {
		attr, ok := resp.Schema.Attributes[name]
		if !ok {
			t.Errorf("attribute %q missing from schema", name)
			continue
		}
		switch a := attr.(type) {
		case schema.StringAttribute:
			if a.Required != want.required {
				t.Errorf("attribute %q Required=%v, want %v", name, a.Required, want.required)
			}
			if a.Optional != want.optional {
				t.Errorf("attribute %q Optional=%v, want %v", name, a.Optional, want.optional)
			}
			if a.Computed != want.computed {
				t.Errorf("attribute %q Computed=%v, want %v", name, a.Computed, want.computed)
			}
		case schema.Int64Attribute:
			if a.Required != want.required {
				t.Errorf("attribute %q Required=%v, want %v", name, a.Required, want.required)
			}
			if a.Optional != want.optional {
				t.Errorf("attribute %q Optional=%v, want %v", name, a.Optional, want.optional)
			}
			if a.Computed != want.computed {
				t.Errorf("attribute %q Computed=%v, want %v", name, a.Computed, want.computed)
			}
		case schema.ListNestedAttribute:
			if a.Required != want.required {
				t.Errorf("attribute %q Required=%v, want %v", name, a.Required, want.required)
			}
			if a.Optional != want.optional {
				t.Errorf("attribute %q Optional=%v, want %v", name, a.Optional, want.optional)
			}
			if a.Computed != want.computed {
				t.Errorf("attribute %q Computed=%v, want %v", name, a.Computed, want.computed)
			}
		default:
			t.Errorf("attribute %q unexpected type %T", name, attr)
		}
	}

	// account_group_id must force replacement.
	cr, ok := resp.Schema.Attributes["account_group_id"].(schema.StringAttribute)
	if !ok {
		t.Fatal("account_group_id not a StringAttribute")
	}
	if len(cr.PlanModifiers) == 0 {
		t.Error("account_group_id should have a RequiresReplace plan modifier")
	}
}

// TestConfigAggregateCompliancePackResource_ImplementsInterfaces verifies the
// resource satisfies the framework interfaces advertised in the _ asserts.
func TestConfigAggregateCompliancePackResource_ImplementsInterfaces(t *testing.T) {
	var _ resource.ResourceWithConfigure = &ConfigAggregateCompliancePackResource{}
	var _ resource.ResourceWithImportState = &ConfigAggregateCompliancePackResource{}
}

// TestConfigAggregateCompliancePackResource_Factory validates the registry factory.
func TestConfigAggregateCompliancePackResource_Factory(t *testing.T) {
	got := NewConfigAggregateCompliancePackResource()
	if got == nil {
		t.Fatal("NewConfigAggregateCompliancePackResource returned nil")
	}
	if _, ok := got.(*ConfigAggregateCompliancePackResource); !ok {
		t.Fatalf("factory returned %T, want *ConfigAggregateCompliancePackResource", got)
	}
}

// TestConfigAggregateCompliancePackResource_ConfigureNilProviderData ensures the
// pre-Configure invocation pattern (nil ProviderData) is silently no-op.
func TestConfigAggregateCompliancePackResource_ConfigureNilProviderData(t *testing.T) {
	r := &ConfigAggregateCompliancePackResource{}
	req := resource.ConfigureRequest{}
	resp := &resource.ConfigureResponse{}

	r.Configure(context.Background(), req, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("Configure with nil ProviderData returned diagnostics: %v", resp.Diagnostics)
	}
}

// --- ID split ---

// TestSplitAggregateCompliancePackId validates parsing of the composite id.
func TestSplitAggregateCompliancePackId(t *testing.T) {
	good := "ca-xxxxxxxx" + tccommon.FILED_SP + "cp-yyyyyyyy"
	ag, cp, ok := splitAggregateCompliancePackId(good)
	if !ok || ag != "ca-xxxxxxxx" || cp != "cp-yyyyyyyy" {
		t.Errorf("splitAggregateCompliancePackId(%q) = (%q, %q, %v); want (\"ca-xxxxxxxx\", \"cp-yyyyyyyy\", true)", good, ag, cp, ok)
	}

	bad := []string{
		"",
		"only-group",
		tccommon.FILED_SP + "cp",
		"ca" + tccommon.FILED_SP,
		"a" + tccommon.FILED_SP + "b" + tccommon.FILED_SP + "c",
	}
	for _, in := range bad {
		if _, _, ok := splitAggregateCompliancePackId(in); ok {
			t.Errorf("splitAggregateCompliancePackId(%q) ok=true, want false", in)
		}
	}
}

// --- expand / flatten ---

// buildRulesList constructs a types.List of config_rules for tests.
func buildRulesList(ctx context.Context, t *testing.T) types.List {
	t.Helper()
	rule := aggregateCompliancePackRuleModel{
		RuleName:              types.StringValue("rule-1"),
		RiskLevel:             types.Int64Value(1),
		Identifier:            types.StringValue("conformance-pack-1"),
		Description:           types.StringValue("rule description"),
		ManagedRuleIdentifier: types.StringValue("qcloud::cvm::check-instance-vpc-subnet"),
		ConfigRuleId:          types.StringValue("cr-aaaa"),
		CompliancePackId:      types.StringValue("cp-yyyyyyyy"),
	}
	param := aggregateCompliancePackInputParameterModel{
		ParameterKey: types.StringValue("param1"),
		Type:          types.StringValue("Require"),
		Value:         types.StringValue("value1"),
	}
	paramList, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: inputParameterAttrTypes()}, []aggregateCompliancePackInputParameterModel{param})
	if diags.HasError() {
		t.Fatalf("build input_parameter list failed: %v", diags)
	}
	rule.InputParameter = paramList

	list, diags := types.ListValueFrom(ctx, types.ObjectType{AttrTypes: configRulesAttrTypes()}, []aggregateCompliancePackRuleModel{rule})
	if diags.HasError() {
		t.Fatalf("build config_rules list failed: %v", diags)
	}
	return list
}

// TestExpandAggregateCompliancePackRules verifies the plan list -> SDK struct
// conversion covers all leaf fields.
func TestExpandAggregateCompliancePackRules(t *testing.T) {
	ctx := context.Background()
	list := buildRulesList(ctx, t)

	rules, diags := expandAggregateCompliancePackRules(ctx, list)
	if diags.HasError() {
		t.Fatalf("expandAggregateCompliancePackRules returned diagnostics: %v", diags)
	}
	if len(rules) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(rules))
	}
	r := rules[0]
	if r.RuleName == nil || *r.RuleName != "rule-1" {
		t.Errorf("RuleName = %v, want rule-1", r.RuleName)
	}
	if r.RiskLevel == nil || *r.RiskLevel != 1 {
		t.Errorf("RiskLevel = %v, want 1", r.RiskLevel)
	}
	if r.Identifier == nil || *r.Identifier != "conformance-pack-1" {
		t.Errorf("Identifier = %v, want conformance-pack-1", r.Identifier)
	}
	if r.Description == nil || *r.Description != "rule description" {
		t.Errorf("Description = %v, want rule description", r.Description)
	}
	if r.ManagedRuleIdentifier == nil || *r.ManagedRuleIdentifier != "qcloud::cvm::check-instance-vpc-subnet" {
		t.Errorf("ManagedRuleIdentifier = %v, want preset id", r.ManagedRuleIdentifier)
	}
	if r.ConfigRuleId == nil || *r.ConfigRuleId != "cr-aaaa" {
		t.Errorf("ConfigRuleId = %v, want cr-aaaa", r.ConfigRuleId)
	}
	if r.CompliancePackId == nil || *r.CompliancePackId != "cp-yyyyyyyy" {
		t.Errorf("CompliancePackId = %v, want cp-yyyyyyyy", r.CompliancePackId)
	}
	if len(r.InputParameter) != 1 {
		t.Fatalf("expected 1 input parameter, got %d", len(r.InputParameter))
	}
	if r.InputParameter[0].ParameterKey == nil || *r.InputParameter[0].ParameterKey != "param1" {
		t.Errorf("ParameterKey = %v, want param1", r.InputParameter[0].ParameterKey)
	}
	if r.InputParameter[0].Type == nil || *r.InputParameter[0].Type != "Require" {
		t.Errorf("Type = %v, want Require", r.InputParameter[0].Type)
	}
	if r.InputParameter[0].Value == nil || *r.InputParameter[0].Value != "value1" {
		t.Errorf("Value = %v, want value1", r.InputParameter[0].Value)
	}
}

// TestExpandAggregateCompliancePackRules_Null verifies a null list yields no rules.
func TestExpandAggregateCompliancePackRules_Null(t *testing.T) {
	rules, diags := expandAggregateCompliancePackRules(context.Background(), types.ListNull(types.ObjectType{AttrTypes: configRulesAttrTypes()}))
	if diags.HasError() {
		t.Fatalf("unexpected diagnostics: %v", diags)
	}
	if rules != nil {
		t.Errorf("expected nil rules for null list, got %v", rules)
	}
}

// TestFlattenAggregateComplianceConfigRules verifies the SDK struct -> list
// conversion reads back leaf fields and preserves managed_rule_identifier /
// compliance_pack_id from the previous state.
func TestFlattenAggregateComplianceConfigRules(t *testing.T) {
	ctx := context.Background()
	prevList := buildRulesList(ctx, t)

	riskLevel := uint64(1)
	rules := []*configv20220802.ComplianceConfigRule{
		{
			RuleName:       helper.String("rule-1"),
			RiskLevel:      &riskLevel,
			Identifier:      helper.String("conformance-pack-1"),
			Description:    helper.String("rule description"),
			ConfigRuleId:   helper.String("cr-aaaa"),
			InputParameter: []*configv20220802.InputParameter{
				{
					ParameterKey: helper.String("param1"),
					Type:          helper.String("Require"),
					Value:         helper.String("value1"),
				},
			},
		},
	}

	list := flattenAggregateComplianceConfigRules(ctx, rules, prevList)
	if list.IsNull() {
		t.Fatal("flattened list is null")
	}

	var out []aggregateCompliancePackRuleModel
	if diags := list.ElementsAs(ctx, &out, false); diags.HasError() {
		t.Fatalf("elements as failed: %v", diags)
	}
	if len(out) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(out))
	}
	r := out[0]
	if r.RuleName.ValueString() != "rule-1" {
		t.Errorf("RuleName = %q, want rule-1", r.RuleName.ValueString())
	}
	if r.RiskLevel.ValueInt64() != 1 {
		t.Errorf("RiskLevel = %d, want 1", r.RiskLevel.ValueInt64())
	}
	if r.Identifier.ValueString() != "conformance-pack-1" {
		t.Errorf("Identifier = %q, want conformance-pack-1", r.Identifier.ValueString())
	}
	if r.Description.ValueString() != "rule description" {
		t.Errorf("Description = %q, want rule description", r.Description.ValueString())
	}
	if r.ConfigRuleId.ValueString() != "cr-aaaa" {
		t.Errorf("ConfigRuleId = %q, want cr-aaaa", r.ConfigRuleId.ValueString())
	}
	// managed_rule_identifier / compliance_pack_id are not echoed by the cloud;
	// they must be carried over from the previous state.
	if r.ManagedRuleIdentifier.ValueString() != "qcloud::cvm::check-instance-vpc-subnet" {
		t.Errorf("ManagedRuleIdentifier = %q, want preserved value", r.ManagedRuleIdentifier.ValueString())
	}
	if r.CompliancePackId.ValueString() != "cp-yyyyyyyy" {
		t.Errorf("CompliancePackId = %q, want preserved value", r.CompliancePackId.ValueString())
	}
	if r.InputParameter.IsNull() {
		t.Fatal("InputParameter is null")
	}
	var params []aggregateCompliancePackInputParameterModel
	if diags := r.InputParameter.ElementsAs(ctx, &params, false); diags.HasError() {
		t.Fatalf("input parameter elements as failed: %v", diags)
	}
	if len(params) != 1 {
		t.Fatalf("expected 1 input parameter, got %d", len(params))
	}
	if params[0].ParameterKey.ValueString() != "param1" {
		t.Errorf("ParameterKey = %q, want param1", params[0].ParameterKey.ValueString())
	}
	if params[0].Type.ValueString() != "Require" {
		t.Errorf("Type = %q, want Require", params[0].Type.ValueString())
	}
	if params[0].Value.ValueString() != "value1" {
		t.Errorf("Value = %q, want value1", params[0].Value.ValueString())
	}
}

// TestFlattenAggregateComplianceConfigRules_Empty verifies an empty cloud list
// yields a null list (not an error).
func TestFlattenAggregateComplianceConfigRules_Empty(t *testing.T) {
	list := flattenAggregateComplianceConfigRules(context.Background(), nil, types.ListNull(types.ObjectType{AttrTypes: configRulesAttrTypes()}))
	if !list.IsNull() {
		t.Errorf("expected null list for empty rules, got %v", list)
	}
}

// TestRefreshAggregateCompliancePackModel verifies cloud-owned top-level fields
// are written into the model and identity parts are reconstructed.
func TestRefreshAggregateCompliancePackModel(t *testing.T) {
	ctx := context.Background()
	m := &ConfigAggregateCompliancePackResourceModel{}
	riskLevel := uint64(2)
	item := &configv20220802.DescribeAggregateCompliancePackResponseParams{
		CompliancePackName: helper.String("pack-name"),
		Description:        helper.String("desc"),
		CreateTime:         helper.String("2024-01-01 00:00:00"),
		RiskLevel:          &riskLevel,
		CompliancePackId:   helper.String("cp-zzz"),
		Status:             helper.String("ACTIVE"),
		ConfigRules:        nil,
	}

	refreshAggregateCompliancePackModel(ctx, m, "ca-aaa", "cp-zzz", item)

	if m.CompliancePackId.ValueString() != "cp-zzz" {
		t.Errorf("CompliancePackId = %q, want cp-zzz", m.CompliancePackId.ValueString())
	}
	if m.AccountGroupId.ValueString() != "ca-aaa" {
		t.Errorf("AccountGroupId = %q, want ca-aaa", m.AccountGroupId.ValueString())
	}
	wantId := "ca-aaa" + tccommon.FILED_SP + "cp-zzz"
	if m.Id.ValueString() != wantId {
		t.Errorf("Id = %q, want %q", m.Id.ValueString(), wantId)
	}
	if m.CompliancePackName.ValueString() != "pack-name" {
		t.Errorf("CompliancePackName = %q, want pack-name", m.CompliancePackName.ValueString())
	}
	if m.Description.ValueString() != "desc" {
		t.Errorf("Description = %q, want desc", m.Description.ValueString())
	}
	if m.CreateTime.ValueString() != "2024-01-01 00:00:00" {
		t.Errorf("CreateTime = %q, want 2024-01-01 00:00:00", m.CreateTime.ValueString())
	}
	if m.Status.ValueString() != "ACTIVE" {
		t.Errorf("Status = %q, want ACTIVE", m.Status.ValueString())
	}
	if m.RiskLevel.ValueInt64() != 2 {
		t.Errorf("RiskLevel = %d, want 2", m.RiskLevel.ValueInt64())
	}
	if !m.ConfigRules.IsNull() {
		t.Errorf("ConfigRules should be null for empty cloud rules, got %v", m.ConfigRules)
	}
}

// --- Service layer with mocked cloud API ---

// mockAggregateClient wires a ConfigService whose underlying client returns a
// patched config SDK client for the duration of the test.
func mockAggregateClient(t *testing.T) (*gomonkey.Patches, ConfigService, *configv20220802.Client) {
	t.Helper()
	patches := gomonkey.NewPatches()
	client := &connectivity.TencentCloudClient{Region: "ap-guangzhou"}
	configClient := &configv20220802.Client{}
	patches.ApplyMethodReturn(client, "UseConfigV20220802Client", configClient)
	service := NewConfigService(client)
	return patches, service, configClient
}

// TestServiceAddAggregateCompliancePack covers the create service path with a
// mocked cloud API, asserting the returned compliance pack id and request echo.
func TestServiceAddAggregateCompliancePack(t *testing.T) {
	patches, service, configClient := mockAggregateClient(t)
	defer patches.Reset()

	var captured *configv20220802.AddAggregateCompliancePackRequest
	patches.ApplyMethodFunc(configClient, "AddAggregateCompliancePackWithContext",
		func(ctx context.Context, request *configv20220802.AddAggregateCompliancePackRequest) (*configv20220802.AddAggregateCompliancePackResponse, error) {
			captured = request
			resp := configv20220802.NewAddAggregateCompliancePackResponse()
			resp.Response = &configv20220802.AddAggregateCompliancePackResponseParams{
				CompliancePackId: helper.String("cp-new"),
			}
			return resp, nil
		})

	req := configv20220802.NewAddAggregateCompliancePackRequest()
	req.CompliancePackName = helper.String("pack-name")
	req.AccountGroupId = helper.String("ca-aaa")
	req.RiskLevel = helper.Uint64(1)
	req.ConfigRules = []*configv20220802.CompliancePackRule{
		{RuleName: helper.String("rule-1"), RiskLevel: helper.Uint64(1)},
	}

	id, err := service.AddAggregateCompliancePack(context.Background(), req)
	if err != nil {
		t.Fatalf("AddAggregateCompliancePack failed: %v", err)
	}
	if id != "cp-new" {
		t.Errorf("id = %q, want cp-new", id)
	}
	if captured == nil {
		t.Fatal("AddAggregateCompliancePackWithContext was not called")
	}
	if captured.CompliancePackName == nil || *captured.CompliancePackName != "pack-name" {
		t.Errorf("captured CompliancePackName = %v, want pack-name", captured.CompliancePackName)
	}
}

// TestServiceAddAggregateCompliancePack_NilResponse asserts a nil response is
// treated as an API failure (not an empty result).
func TestServiceAddAggregateCompliancePack_NilResponse(t *testing.T) {
	patches, service, configClient := mockAggregateClient(t)
	defer patches.Reset()

	patches.ApplyMethodReturn(configClient, "AddAggregateCompliancePackWithContext", (*configv20220802.AddAggregateCompliancePackResponse)(nil), nil)

	_, err := service.AddAggregateCompliancePack(context.Background(), configv20220802.NewAddAggregateCompliancePackRequest())
	if err == nil {
		t.Fatal("expected error for nil response, got nil")
	}
}

// TestServiceDescribeAggregateCompliancePackById covers the read path and the
// "not found" (nil item) signal.
func TestServiceDescribeAggregateCompliancePackById(t *testing.T) {
	patches, service, configClient := mockAggregateClient(t)
	defer patches.Reset()

	riskLevel := uint64(1)
	patches.ApplyMethodFunc(configClient, "DescribeAggregateCompliancePackWithContext",
		func(ctx context.Context, request *configv20220802.DescribeAggregateCompliancePackRequest) (*configv20220802.DescribeAggregateCompliancePackResponse, error) {
			if request.AccountGroupId == nil || *request.AccountGroupId != "ca-aaa" {
				t.Errorf("captured AccountGroupId = %v, want ca-aaa", request.AccountGroupId)
			}
			resp := configv20220802.NewDescribeAggregateCompliancePackResponse()
			resp.Response = &configv20220802.DescribeAggregateCompliancePackResponseParams{
				CompliancePackName: helper.String("pack-name"),
				RiskLevel:          &riskLevel,
				CompliancePackId:   helper.String("cp-bbb"),
			}
			return resp, nil
		})

	item, err := service.DescribeAggregateCompliancePackById(context.Background(), "ca-aaa", "cp-bbb")
	if err != nil {
		t.Fatalf("DescribeAggregateCompliancePackById failed: %v", err)
	}
	if item == nil || item.CompliancePackId == nil || *item.CompliancePackId != "cp-bbb" {
		t.Errorf("item = %+v, want CompliancePackId cp-bbb", item)
	}
}

// TestServiceDescribeAggregateCompliancePackById_NotFound asserts a response
// with a nil CompliancePackId yields a non-nil params whose id is nil (the
// resource Read interprets this as "not found").
func TestServiceDescribeAggregateCompliancePackById_NotFound(t *testing.T) {
	patches, service, configClient := mockAggregateClient(t)
	defer patches.Reset()

	patches.ApplyMethodFunc(configClient, "DescribeAggregateCompliancePackWithContext",
		func(ctx context.Context, request *configv20220802.DescribeAggregateCompliancePackRequest) (*configv20220802.DescribeAggregateCompliancePackResponse, error) {
			resp := configv20220802.NewDescribeAggregateCompliancePackResponse()
			resp.Response = &configv20220802.DescribeAggregateCompliancePackResponseParams{}
			return resp, nil
		})

	item, err := service.DescribeAggregateCompliancePackById(context.Background(), "ca-aaa", "cp-bbb")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item == nil || item.CompliancePackId != nil {
		t.Errorf("expected item with nil CompliancePackId, got %+v", item)
	}
}

// TestServiceUpdateAggregateCompliancePack covers the update path.
func TestServiceUpdateAggregateCompliancePack(t *testing.T) {
	patches, service, configClient := mockAggregateClient(t)
	defer patches.Reset()

	var captured *configv20220802.UpdateAggregateCompliancePackRequest
	patches.ApplyMethodFunc(configClient, "UpdateAggregateCompliancePackWithContext",
		func(ctx context.Context, request *configv20220802.UpdateAggregateCompliancePackRequest) (*configv20220802.UpdateAggregateCompliancePackResponse, error) {
			captured = request
			resp := configv20220802.NewUpdateAggregateCompliancePackResponse()
			resp.Response = &configv20220802.UpdateAggregateCompliancePackResponseParams{}
			return resp, nil
		})

	req := configv20220802.NewUpdateAggregateCompliancePackRequest()
	req.AccountGroupId = helper.String("ca-aaa")
	req.CompliancePackId = helper.String("cp-bbb")
	req.CompliancePackName = helper.String("pack-name-updated")

	if err := service.UpdateAggregateCompliancePack(context.Background(), req); err != nil {
		t.Fatalf("UpdateAggregateCompliancePack failed: %v", err)
	}
	if captured == nil {
		t.Fatal("UpdateAggregateCompliancePackWithContext was not called")
	}
	if captured.CompliancePackName == nil || *captured.CompliancePackName != "pack-name-updated" {
		t.Errorf("captured CompliancePackName = %v, want pack-name-updated", captured.CompliancePackName)
	}
}

// TestServiceDeleteAggregateCompliancePack covers the delete path.
func TestServiceDeleteAggregateCompliancePack(t *testing.T) {
	patches, service, configClient := mockAggregateClient(t)
	defer patches.Reset()

	var captured *configv20220802.DeleteAggregateCompliancePackRequest
	patches.ApplyMethodFunc(configClient, "DeleteAggregateCompliancePackWithContext",
		func(ctx context.Context, request *configv20220802.DeleteAggregateCompliancePackRequest) (*configv20220802.DeleteAggregateCompliancePackResponse, error) {
			captured = request
			resp := configv20220802.NewDeleteAggregateCompliancePackResponse()
			resp.Response = &configv20220802.DeleteAggregateCompliancePackResponseParams{}
			return resp, nil
		})

	req := configv20220802.NewDeleteAggregateCompliancePackRequest()
	req.AccountGroupId = helper.String("ca-aaa")
	req.CompliancePackId = helper.String("cp-bbb")

	if err := service.DeleteAggregateCompliancePack(context.Background(), req); err != nil {
		t.Fatalf("DeleteAggregateCompliancePack failed: %v", err)
	}
	if captured == nil {
		t.Fatal("DeleteAggregateCompliancePackWithContext was not called")
	}
	if captured.AccountGroupId == nil || *captured.AccountGroupId != "ca-aaa" {
		t.Errorf("captured AccountGroupId = %v, want ca-aaa", captured.AccountGroupId)
	}
	if captured.CompliancePackId == nil || *captured.CompliancePackId != "cp-bbb" {
		t.Errorf("captured CompliancePackId = %v, want cp-bbb", captured.CompliancePackId)
	}
}

// TestServiceDeleteAggregateCompliancePack_APIError asserts a NotFound-class
// SDK error is tolerated by the resource Delete (the resource checks
// IsExpectError on a real *TencentCloudSDKError).
func TestServiceDeleteAggregateCompliancePack_APIError(t *testing.T) {
	patches, service, configClient := mockAggregateClient(t)
	defer patches.Reset()

	// A real NotFound-class SDK error (long code like "ResourceNotFound.CompliancePack")
	// must be recognized by IsExpectError via its short-code form.
	notFoundErr := sdkerr.NewTencentCloudSDKError("ResourceNotFound.CompliancePack", "not found", "")
	if !tccommon.IsExpectError(notFoundErr, []string{"ResourceNotFound"}) {
		t.Fatal("IsExpectError should recognize ResourceNotFound.CompliancePack via short code")
	}

	// The service surfaces any error; the resource Delete tolerates
	// NotFound-class errors via IsExpectError.
	patches.ApplyMethodReturn(configClient, "DeleteAggregateCompliancePackWithContext",
		(*configv20220802.DeleteAggregateCompliancePackResponse)(nil),
		fmt.Errorf("ResourceNotFound: compliance pack not found"))

	if err := service.DeleteAggregateCompliancePack(context.Background(), configv20220802.NewDeleteAggregateCompliancePackRequest()); err == nil {
		t.Fatal("expected error from delete, got nil")
	}
}
