package config_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	configv20220802 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	svcconfig "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/config"
)

func ptrStrLACR(s string) *string { return &s }

func ptrUint64LACR(v uint64) *uint64 { return &v }

func ptrInt64LACR(v int64) *int64 { return &v }

func ptrBoolLACR(v bool) *bool { return &v }

func newMockProviderMetaLACR() *sharedmeta.ProviderMeta {
	return &sharedmeta.ProviderMeta{Client: &connectivity.TencentCloudClient{}}
}

func setConfigListAggregateConfigRulesClient(d *svcconfig.ConfigListAggregateConfigRulesDataSource, client *connectivity.TencentCloudClient) {
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := datasource.ConfigureRequest{ProviderData: meta}
	configureResp := &datasource.ConfigureResponse{}
	d.Configure(context.Background(), configureReq, configureResp)
}

// configListAggregateConfigRulesSchema returns the data source schema used to
// build the raw config value for ReadRequest.
func configListAggregateConfigRulesSchema(t *testing.T) schema.Schema {
	d := svcconfig.NewConfigListAggregateConfigRulesDataSource().(*svcconfig.ConfigListAggregateConfigRulesDataSource)
	schemaResp := &datasource.SchemaResponse{}
	d.Schema(context.Background(), datasource.SchemaRequest{}, schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("failed to build data source schema: %v", schemaResp.Diagnostics)
	}
	return schemaResp.Schema
}

// newConfigListAggregateConfigRulesReadRequest builds a ReadRequest carrying
// the given account_group_id as a raw tftypes value. All other attributes are
// null because account_group_id is the only required input exercised here.
func newConfigListAggregateConfigRulesReadRequest(t *testing.T, accountGroupId string) datasource.ReadRequest {
	s := configListAggregateConfigRulesSchema(t)

	schemaType := s.Type()
	tfType := schemaType.TerraformType(context.Background())

	objType, ok := tfType.(tftypes.Object)
	if !ok {
		t.Fatalf("schema type is %T, want tftypes.Object", tfType)
	}

	values := make(map[string]tftypes.Value, len(objType.AttributeTypes))
	for name, typ := range objType.AttributeTypes {
		// null for optional inputs and computed outputs
		values[name] = tftypes.NewValue(typ, nil)
	}
	values["account_group_id"] = tftypes.NewValue(tftypes.String, accountGroupId)

	raw := tftypes.NewValue(objType, values)
	req := datasource.ReadRequest{}
	req.Config.Raw = raw
	req.Config.Schema = s
	return req
}

func buildAggregateConfigRule() *configv20220802.ConfigRule {
	return &configv20220802.ConfigRule{
		Identifier:            ptrStrLACR("cam-user-mfa-check"),
		RuleName:              ptrStrLACR("cam-user-mfa-check"),
		RiskLevel:             ptrInt64LACR(1),
		ServiceFunction:       ptrStrLACR("cam-user-mfa-check"),
		CreateTime:            ptrStrLACR("2024-01-01 00:00:00"),
		Description:           ptrStrLACR("check cam user mfa"),
		Status:                ptrStrLACR("ACTIVE"),
		ComplianceResult:      ptrStrLACR("COMPLIANT"),
		ConfigRuleInvokedTime: ptrStrLACR("2024-01-02 00:00:00"),
		ConfigRuleId:          ptrStrLACR("cr-abcdefg"),
		IdentifierType:        ptrStrLACR("SYSTEM"),
		CompliancePackId:      ptrStrLACR(""),
		CompliancePackName:    ptrStrLACR(""),
		AccountGroupId:        ptrStrLACR("ca-aggregator-abcdefg"),
		AccountGroupName:      ptrStrLACR("tf-aggregator"),
		RuleOwnerId:           ptrUint64LACR(100000000000),
		InputParameter: []*configv20220802.InputParameter{
			{ParameterKey: ptrStrLACR("maxIdle"), Type: ptrStrLACR("Optional"), Value: ptrStrLACR("30")},
		},
		SourceCondition: []*configv20220802.SourceConditionForManage{
			{EmptyAs: ptrStrLACR("NON_COMPLIANT"), SelectPath: ptrStrLACR("$.maxIdle"), Operator: ptrStrLACR("lt"), Required: ptrBoolLACR(true), DesiredValue: ptrStrLACR("30")},
		},
		ResourceType:            []*string{ptrStrLACR("Qcloud::CAM::User")},
		Labels:                  []*string{ptrStrLACR("security")},
		Annotation:              &configv20220802.Annotation{Configuration: ptrStrLACR("idle=60"), DesiredValue: ptrStrLACR("30"), Operator: ptrStrLACR("lt"), Property: ptrStrLACR("$.maxIdle")},
		TriggerType:             []*configv20220802.TriggerType{{MessageType: ptrStrLACR("ScheduledNotification"), MaximumExecutionFrequency: ptrStrLACR("One_Hour")}},
		ManageInputParameter:    []*configv20220802.InputParameterForManage{{ValueType: ptrStrLACR("Integer"), ParameterKey: ptrStrLACR("maxIdle"), Type: ptrStrLACR("Optional"), DefaultValue: ptrStrLACR("30"), Description: ptrStrLACR("max idle days")}},
		RegionsScope:            []*string{ptrStrLACR("ap-guangzhou")},
		TagsScope:               []*configv20220802.Tag{{TagKey: ptrStrLACR("env"), TagValue: ptrStrLACR("prod")}},
		ExcludeResourceIdsScope: []*string{ptrStrLACR("usr-123")},
		ManageTriggerType:       []*string{ptrStrLACR("ScheduledNotification")},
	}
}

// TestConfigListAggregateConfigRulesDS_MetadataType validates the registry type name.
func TestConfigListAggregateConfigRulesDS_MetadataType(t *testing.T) {
	d := svcconfig.NewConfigListAggregateConfigRulesDataSource().(*svcconfig.ConfigListAggregateConfigRulesDataSource)
	req := datasource.MetadataRequest{}
	resp := &datasource.MetadataResponse{}

	d.Metadata(context.Background(), req, resp)

	assert.Equal(t, "tencentcloud_config_list_aggregate_config_rules", resp.TypeName)
}

// TestConfigListAggregateConfigRulesDS_Schema validates input/computed attribute flags.
func TestConfigListAggregateConfigRulesDS_Schema(t *testing.T) {
	d := svcconfig.NewConfigListAggregateConfigRulesDataSource().(*svcconfig.ConfigListAggregateConfigRulesDataSource)
	req := datasource.SchemaRequest{}
	resp := &datasource.SchemaResponse{}

	d.Schema(context.Background(), req, resp)
	assert.False(t, resp.Diagnostics.HasError(), "schema returned diagnostics: %v", resp.Diagnostics)

	attr, ok := resp.Schema.Attributes["account_group_id"].(schema.StringAttribute)
	assert.True(t, ok, "account_group_id must be a StringAttribute")
	assert.True(t, attr.Required, "account_group_id must be Required")

	for _, name := range []string{"order_type", "state", "rule_name"} {
		a, ok := resp.Schema.Attributes[name].(schema.StringAttribute)
		assert.True(t, ok, "%s must be a StringAttribute", name)
		assert.True(t, a.Optional, "%s must be Optional", name)
	}

	totalAttr, ok := resp.Schema.Attributes["total"].(schema.Int64Attribute)
	assert.True(t, ok, "total must be an Int64Attribute")
	assert.True(t, totalAttr.Computed, "total must be Computed")

	itemsAttr, ok := resp.Schema.Attributes["items"].(schema.ListNestedAttribute)
	assert.True(t, ok, "items must be a ListNestedAttribute")
	assert.True(t, itemsAttr.Computed, "items must be Computed")
	assert.Contains(t, itemsAttr.NestedObject.Attributes, "config_rule_id")
	assert.Contains(t, itemsAttr.NestedObject.Attributes, "rule_name")
}

// TestConfigListAggregateConfigRulesDS_ReadBasic covers a single matched rule and
// asserts every leaf field round-trips into state.
func TestConfigListAggregateConfigRulesDS_ReadBasic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	configClient := &configv20220802.Client{}
	patches.ApplyMethodReturn(newMockProviderMetaLACR().Client, "UseConfigV20220802Client", configClient)

	patches.ApplyMethodFunc(configClient, "ListAggregateConfigRulesWithContext", func(_ context.Context, request *configv20220802.ListAggregateConfigRulesRequest) (*configv20220802.ListAggregateConfigRulesResponse, error) {
		assert.Equal(t, "ca-aggregator-abcdefg", *request.AccountGroupId)
		resp := configv20220802.NewListAggregateConfigRulesResponse()
		resp.Response = &configv20220802.ListAggregateConfigRulesResponseParams{
			Total:     ptrUint64LACR(1),
			Items:     []*configv20220802.ConfigRule{buildAggregateConfigRule()},
			RequestId: ptrStrLACR("fake-request-id"),
		}
		return resp, nil
	})

	d := svcconfig.NewConfigListAggregateConfigRulesDataSource().(*svcconfig.ConfigListAggregateConfigRulesDataSource)
	setConfigListAggregateConfigRulesClient(d, newMockProviderMetaLACR().Client)

	rreq := newConfigListAggregateConfigRulesReadRequest(t, "ca-aggregator-abcdefg")
	rresp := &datasource.ReadResponse{}
	d.Read(context.Background(), rreq, rresp)

	assert.False(t, rresp.Diagnostics.HasError(), "Read returned diagnostics: %v", rresp.Diagnostics)

	var model svcconfig.ConfigListAggregateConfigRulesDataSourceModel
	if diags := rresp.State.Get(context.Background(), &model); diags.HasError() {
		t.Fatalf("State.Get returned diagnostics: %v", diags)
	}

	assert.Equal(t, int64(1), model.Total.ValueInt64())

	items := make([]svcconfig.AggregateConfigRuleItemModelExport, 0)
	if diags := model.Items.ElementsAs(context.Background(), &items, false); diags.HasError() {
		t.Fatalf("Items.ElementsAs returned diagnostics: %v", diags)
	}
	assert.Len(t, items, 1)

	item := items[0]
	assert.Equal(t, "cam-user-mfa-check", item.RuleName.ValueString())
	assert.Equal(t, "cam-user-mfa-check", item.Identifier.ValueString())
	assert.Equal(t, "cr-abcdefg", item.ConfigRuleId.ValueString())
	assert.Equal(t, int64(1), item.RiskLevel.ValueInt64())
	assert.Equal(t, int64(100000000000), item.RuleOwnerId.ValueInt64())
	assert.Equal(t, "ca-aggregator-abcdefg", item.AccountGroupId.ValueString())
	assert.Equal(t, "tf-aggregator", item.AccountGroupName.ValueString())
	assert.Equal(t, "ACTIVE", item.Status.ValueString())
	assert.Equal(t, "COMPLIANT", item.ComplianceResult.ValueString())
	assert.Equal(t, "SYSTEM", item.IdentifierType.ValueString())
	assert.Equal(t, "2024-01-01 00:00:00", item.CreateTime.ValueString())
	assert.Equal(t, "check cam user mfa", item.Description.ValueString())
	assert.Equal(t, "2024-01-02 00:00:00", item.ConfigRuleInvokedTime.ValueString())

	// input_parameter
	ipItems := make([]svcconfig.AggregateInputParameterModelExport, 0)
	if diags := item.InputParameter.ElementsAs(context.Background(), &ipItems, false); diags.HasError() {
		t.Fatalf("input_parameter ElementsAs diagnostics: %v", diags)
	}
	assert.Len(t, ipItems, 1)
	assert.Equal(t, "maxIdle", ipItems[0].ParameterKey.ValueString())
	assert.Equal(t, "Optional", ipItems[0].Type.ValueString())
	assert.Equal(t, "30", ipItems[0].Value.ValueString())

	// source_condition
	scItems := make([]svcconfig.AggregateSourceConditionModelExport, 0)
	if diags := item.SourceCondition.ElementsAs(context.Background(), &scItems, false); diags.HasError() {
		t.Fatalf("source_condition ElementsAs diagnostics: %v", diags)
	}
	assert.Len(t, scItems, 1)
	assert.Equal(t, "NON_COMPLIANT", scItems[0].EmptyAs.ValueString())
	assert.Equal(t, "$.maxIdle", scItems[0].SelectPath.ValueString())
	assert.Equal(t, "lt", scItems[0].Operator.ValueString())
	assert.True(t, scItems[0].Required.ValueBool())
	assert.Equal(t, "30", scItems[0].DesiredValue.ValueString())

	// resource_type string list
	rtItems := make([]types.String, 0)
	if diags := item.ResourceType.ElementsAs(context.Background(), &rtItems, false); diags.HasError() {
		t.Fatalf("resource_type ElementsAs diagnostics: %v", diags)
	}
	assert.Len(t, rtItems, 1)
	assert.Equal(t, "Qcloud::CAM::User", rtItems[0].ValueString())

	// labels string list
	lbItems := make([]types.String, 0)
	if diags := item.Labels.ElementsAs(context.Background(), &lbItems, false); diags.HasError() {
		t.Fatalf("labels ElementsAs diagnostics: %v", diags)
	}
	assert.Len(t, lbItems, 1)
	assert.Equal(t, "security", lbItems[0].ValueString())

	// annotation
	anItems := make([]svcconfig.AggregateAnnotationModelExport, 0)
	if diags := item.Annotation.ElementsAs(context.Background(), &anItems, false); diags.HasError() {
		t.Fatalf("annotation ElementsAs diagnostics: %v", diags)
	}
	assert.Len(t, anItems, 1)
	assert.Equal(t, "idle=60", anItems[0].Configuration.ValueString())
	assert.Equal(t, "30", anItems[0].DesiredValue.ValueString())
	assert.Equal(t, "lt", anItems[0].Operator.ValueString())
	assert.Equal(t, "$.maxIdle", anItems[0].Property.ValueString())

	// trigger_type
	ttItems := make([]svcconfig.AggregateTriggerTypeModelExport, 0)
	if diags := item.TriggerType.ElementsAs(context.Background(), &ttItems, false); diags.HasError() {
		t.Fatalf("trigger_type ElementsAs diagnostics: %v", diags)
	}
	assert.Len(t, ttItems, 1)
	assert.Equal(t, "ScheduledNotification", ttItems[0].MessageType.ValueString())
	assert.Equal(t, "One_Hour", ttItems[0].MaximumExecutionFrequency.ValueString())

	// manage_input_parameter
	mipItems := make([]svcconfig.AggregateManageInputParameterModelExport, 0)
	if diags := item.ManageInputParameter.ElementsAs(context.Background(), &mipItems, false); diags.HasError() {
		t.Fatalf("manage_input_parameter ElementsAs diagnostics: %v", diags)
	}
	assert.Len(t, mipItems, 1)
	assert.Equal(t, "Integer", mipItems[0].ValueType.ValueString())
	assert.Equal(t, "maxIdle", mipItems[0].ParameterKey.ValueString())
	assert.Equal(t, "Optional", mipItems[0].Type.ValueString())
	assert.Equal(t, "30", mipItems[0].DefaultValue.ValueString())
	assert.Equal(t, "max idle days", mipItems[0].Description.ValueString())

	// tags_scope
	tsItems := make([]svcconfig.AggregateTagsScopeModelExport, 0)
	if diags := item.TagsScope.ElementsAs(context.Background(), &tsItems, false); diags.HasError() {
		t.Fatalf("tags_scope ElementsAs diagnostics: %v", diags)
	}
	assert.Len(t, tsItems, 1)
	assert.Equal(t, "env", tsItems[0].TagKey.ValueString())
	assert.Equal(t, "prod", tsItems[0].TagValue.ValueString())

	// regions_scope
	rsItems := make([]types.String, 0)
	if diags := item.RegionsScope.ElementsAs(context.Background(), &rsItems, false); diags.HasError() {
		t.Fatalf("regions_scope ElementsAs diagnostics: %v", diags)
	}
	assert.Len(t, rsItems, 1)
	assert.Equal(t, "ap-guangzhou", rsItems[0].ValueString())

	// exclude_resource_ids_scope
	exItems := make([]types.String, 0)
	if diags := item.ExcludeResourceIdsScope.ElementsAs(context.Background(), &exItems, false); diags.HasError() {
		t.Fatalf("exclude_resource_ids_scope ElementsAs diagnostics: %v", diags)
	}
	assert.Len(t, exItems, 1)
	assert.Equal(t, "usr-123", exItems[0].ValueString())

	// manage_trigger_type
	mtItems := make([]types.String, 0)
	if diags := item.ManageTriggerType.ElementsAs(context.Background(), &mtItems, false); diags.HasError() {
		t.Fatalf("manage_trigger_type ElementsAs diagnostics: %v", diags)
	}
	assert.Len(t, mtItems, 1)
	assert.Equal(t, "ScheduledNotification", mtItems[0].ValueString())
}

// TestConfigListAggregateConfigRulesDS_ReadEmpty covers an account group that has
// no rules yet. An empty result is a legitimate answer, not a failure.
func TestConfigListAggregateConfigRulesDS_ReadEmpty(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	configClient := &configv20220802.Client{}
	patches.ApplyMethodReturn(newMockProviderMetaLACR().Client, "UseConfigV20220802Client", configClient)

	patches.ApplyMethodFunc(configClient, "ListAggregateConfigRulesWithContext", func(_ context.Context, request *configv20220802.ListAggregateConfigRulesRequest) (*configv20220802.ListAggregateConfigRulesResponse, error) {
		resp := configv20220802.NewListAggregateConfigRulesResponse()
		resp.Response = &configv20220802.ListAggregateConfigRulesResponseParams{
			Total:     ptrUint64LACR(0),
			Items:     []*configv20220802.ConfigRule{},
			RequestId: ptrStrLACR("fake-request-id"),
		}
		return resp, nil
	})

	d := svcconfig.NewConfigListAggregateConfigRulesDataSource().(*svcconfig.ConfigListAggregateConfigRulesDataSource)
	setConfigListAggregateConfigRulesClient(d, newMockProviderMetaLACR().Client)

	rreq := newConfigListAggregateConfigRulesReadRequest(t, "ca-aggregator-empty")
	rresp := &datasource.ReadResponse{}
	d.Read(context.Background(), rreq, rresp)

	assert.False(t, rresp.Diagnostics.HasError(), "an account group without rules must not be reported as an error: %v", rresp.Diagnostics)

	var model svcconfig.ConfigListAggregateConfigRulesDataSourceModel
	if diags := rresp.State.Get(context.Background(), &model); diags.HasError() {
		t.Fatalf("State.Get returned diagnostics: %v", diags)
	}
	assert.Equal(t, int64(0), model.Total.ValueInt64())
	assert.False(t, model.Items.IsNull(), "items must be an empty list, not null")
	items := make([]svcconfig.AggregateConfigRuleItemModelExport, 0)
	_ = model.Items.ElementsAs(context.Background(), &items, false)
	assert.Empty(t, items)
}

// TestConfigListAggregateConfigRulesDS_ReadNilResponse covers the case where the
// API returns a nil response. This is a failure, not an empty result.
func TestConfigListAggregateConfigRulesDS_ReadNilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	configClient := &configv20220802.Client{}
	patches.ApplyMethodReturn(newMockProviderMetaLACR().Client, "UseConfigV20220802Client", configClient)

	patches.ApplyMethodFunc(configClient, "ListAggregateConfigRulesWithContext", func(_ context.Context, request *configv20220802.ListAggregateConfigRulesRequest) (*configv20220802.ListAggregateConfigRulesResponse, error) {
		return nil, fmt.Errorf("internal server error")
	})

	d := svcconfig.NewConfigListAggregateConfigRulesDataSource().(*svcconfig.ConfigListAggregateConfigRulesDataSource)
	setConfigListAggregateConfigRulesClient(d, newMockProviderMetaLACR().Client)

	rreq := newConfigListAggregateConfigRulesReadRequest(t, "ca-aggregator-abcdefg")
	rresp := &datasource.ReadResponse{}
	d.Read(context.Background(), rreq, rresp)

	assert.True(t, rresp.Diagnostics.HasError(), "an API error must surface as diagnostics")
	assert.Contains(t, rresp.Diagnostics.Errors()[0].Summary(), "querying Config aggregate config rules")
}

// TestConfigListAggregateConfigRulesDS_ConfigureNilProviderData ensures the
// pre-Configure invocation pattern (nil ProviderData) is silently no-op.
func TestConfigListAggregateConfigRulesDS_ConfigureNilProviderData(t *testing.T) {
	d := svcconfig.NewConfigListAggregateConfigRulesDataSource().(*svcconfig.ConfigListAggregateConfigRulesDataSource)
	req := datasource.ConfigureRequest{}
	resp := &datasource.ConfigureResponse{}

	d.Configure(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "Configure with nil ProviderData returned diagnostics: %v", resp.Diagnostics)
}

// TestConfigListAggregateConfigRulesDS_ImplementsInterface is a compile-time guard
// ensuring the data source implements the framework interfaces it advertises.
func TestConfigListAggregateConfigRulesDS_ImplementsInterface(t *testing.T) {
	var _ datasource.DataSourceWithConfigure = &svcconfig.ConfigListAggregateConfigRulesDataSource{}
}
