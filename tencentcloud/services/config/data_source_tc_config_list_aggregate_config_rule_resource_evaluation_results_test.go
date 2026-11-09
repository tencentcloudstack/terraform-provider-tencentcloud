package config_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	configv20220802 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	svcconfig "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/config"
)

func ptrStringAggEval(s string) *string { return &s }

func ptrInt64AggEval(v int64) *int64 { return &v }

// aggEvalResultSchema returns the framework data source schema.
func aggEvalResultSchema(t *testing.T) datasource.Schema {
	d := svcconfig.NewConfigListAggregateConfigRuleResourceEvaluationResultsDataSource()
	schemaResp := &datasource.SchemaResponse{}
	d.Schema(context.Background(), datasource.SchemaRequest{}, schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("failed to build data source schema: %v", schemaResp.Diagnostics)
	}
	return schemaResp.Schema
}

// newAggEvalReadRequest builds a datasource.ReadRequest carrying the
// given account_group_id value in the raw config.
func newAggEvalReadRequest(t *testing.T, accountGroupId string) datasource.ReadRequest {
	s := aggEvalResultSchema(t)

	var accountGroupIdVal tftypes.Value
	if accountGroupId == "" {
		accountGroupIdVal = tftypes.NewValue(tftypes.String, nil)
	} else {
		accountGroupIdVal = tftypes.NewValue(tftypes.String, accountGroupId)
	}

	// Build the root object type from the schema.
	rootType := s.Type().TerraformType(context.Background())
	objType, ok := rootType.(tftypes.Object)
	if !ok {
		t.Fatalf("schema root type is not tftypes.Object: %T", rootType)
	}

	// items is computed; supply a null list so the config object is well-formed.
	itemsType, itemsOk := objType.AttributeTypes["items"]
	if !itemsOk {
		t.Fatalf("schema does not contain items attribute")
	}

	raw := tftypes.NewValue(objType, map[string]tftypes.Value{
		"account_group_id": accountGroupIdVal,
		"items":            tftypes.NewValue(itemsType, nil),
	})

	req := datasource.ReadRequest{}
	req.Config.Raw = raw
	req.Config.Schema = s
	return req
}

// setAggEvalDataSourceClient injects a mock connectivity client into the
// data source instance via the Configure method.
func setAggEvalDataSourceClient(d *svcconfig.ConfigListAggregateConfigRuleResourceEvaluationResultsDataSource, client *connectivity.TencentCloudClient) {
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := datasource.ConfigureRequest{ProviderData: meta}
	configureResp := &datasource.ConfigureResponse{}
	d.Configure(context.Background(), configureReq, configureResp)
}

// buildAggEvalResult assembles a ConfigRuleResourceEvaluationResult with all
// leaf fields populated, including tags and annotation.
func buildAggEvalResult() *configv20220802.ConfigRuleResourceEvaluationResult {
	return &configv20220802.ConfigRuleResourceEvaluationResult{
		RuleId:          ptrStringAggEval("cr-rule-001"),
		RuleName:        ptrStringAggEval("rule-name"),
		RuleIdentifier:  ptrStringAggEval("rule-identifier"),
		RuleDescription: ptrStringAggEval("rule-description"),
		RuleOwnerId:     ptrInt64AggEval(100012345678),
		RuleRiskLevel:   ptrInt64AggEval(1),
		ResourceId:      ptrStringAggEval("ins-resource-001"),
		ResourceType:    ptrStringAggEval("QCS::CVM::Instance"),
		ResourceName:    ptrStringAggEval("resource-name"),
		ResourceRegion:  ptrStringAggEval("ap-guangzhou"),
		ResourceTags: []*configv20220802.Tag{
			{
				TagKey:   ptrStringAggEval("env"),
				TagValue: ptrStringAggEval("prod"),
			},
		},
		Annotation: &configv20220802.Annotation{
			Configuration: ptrStringAggEval("actual-config"),
			DesiredValue:  ptrStringAggEval("desired-config"),
			Operator:      ptrStringAggEval("EQUAL"),
			Property:      ptrStringAggEval("$.property"),
		},
	}
}

// TestConfigListAggregateConfigRuleResourceEvaluationResultsDS_ReadBasic
// verifies that a single fully populated result is read back into state with
// every leaf field, nested tags and annotation intact.
func TestConfigListAggregateConfigRuleResourceEvaluationResultsDS_ReadBasic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyMethodFunc(new(svcconfig.ConfigService), "DescribeConfigListAggregateConfigRuleResourceEvaluationResultsByFilter",
		func(_ context.Context, _ map[string]interface{}) ([]*configv20220802.ConfigRuleResourceEvaluationResult, error) {
			return []*configv20220802.ConfigRuleResourceEvaluationResult{buildAggEvalResult()}, nil
		})

	ds := svcconfig.NewConfigListAggregateConfigRuleResourceEvaluationResultsDataSource().(*svcconfig.ConfigListAggregateConfigRuleResourceEvaluationResultsDataSource)
	setAggEvalDataSourceClient(ds, &connectivity.TencentCloudClient{})

	req := newAggEvalReadRequest(t, "ca-aggregator-abcdefg")
	resp := &datasource.ReadResponse{}
	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	var state svcconfig.ConfigListAggregateConfigRuleResourceEvaluationResultsModel
	resp.State.Get(context.Background(), &state)

	assert.Equal(t, "ca-aggregator-abcdefg", state.AccountGroupId.ValueString())
	assert.False(t, state.Items.IsNull())

	var items []svcconfig.AggregateConfigRuleResourceEvaluationResultItemModel
	diags := state.Items.ElementsAs(context.Background(), &items, false)
	assert.False(t, diags.HasError(), "expected no error decoding items, got: %v", diags)
	assert.Len(t, items, 1)

	item := items[0]
	assert.Equal(t, "cr-rule-001", item.RuleId.ValueString())
	assert.Equal(t, "rule-name", item.RuleName.ValueString())
	assert.Equal(t, "rule-identifier", item.RuleIdentifier.ValueString())
	assert.Equal(t, "rule-description", item.RuleDescription.ValueString())
	assert.Equal(t, int64(100012345678), item.RuleOwnerId.ValueInt64())
	assert.Equal(t, int64(1), item.RuleRiskLevel.ValueInt64())
	assert.Equal(t, "ins-resource-001", item.ResourceId.ValueString())
	assert.Equal(t, "QCS::CVM::Instance", item.ResourceType.ValueString())
	assert.Equal(t, "resource-name", item.ResourceName.ValueString())
	assert.Equal(t, "ap-guangzhou", item.ResourceRegion.ValueString())

	// resource_tags
	assert.False(t, item.ResourceTags.IsNull())
	var tags []svcconfig.ResourceTagItemModel
	tagDiags := item.ResourceTags.ElementsAs(context.Background(), &tags, false)
	assert.False(t, tagDiags.HasError())
	assert.Len(t, tags, 1)
	assert.Equal(t, "env", tags[0].TagKey.ValueString())
	assert.Equal(t, "prod", tags[0].TagValue.ValueString())

	// annotation
	assert.False(t, item.Annotation.IsNull())
	var annotations []svcconfig.AnnotationItemModel
	annDiags := item.Annotation.ElementsAs(context.Background(), &annotations, false)
	assert.False(t, annDiags.HasError())
	assert.Len(t, annotations, 1)
	assert.Equal(t, "actual-config", annotations[0].Configuration.ValueString())
	assert.Equal(t, "desired-config", annotations[0].DesiredValue.ValueString())
	assert.Equal(t, "EQUAL", annotations[0].Operator.ValueString())
	assert.Equal(t, "$.property", annotations[0].Property.ValueString())
}

// TestConfigListAggregateConfigRuleResourceEvaluationResultsDS_ReadEmpty
// covers an account with no evaluation results: the data source must succeed
// and expose an empty items list.
func TestConfigListAggregateConfigRuleResourceEvaluationResultsDS_ReadEmpty(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyMethodFunc(new(svcconfig.ConfigService), "DescribeConfigListAggregateConfigRuleResourceEvaluationResultsByFilter",
		func(_ context.Context, _ map[string]interface{}) ([]*configv20220802.ConfigRuleResourceEvaluationResult, error) {
			return []*configv20220802.ConfigRuleResourceEvaluationResult{}, nil
		})

	ds := svcconfig.NewConfigListAggregateConfigRuleResourceEvaluationResultsDataSource().(*svcconfig.ConfigListAggregateConfigRuleResourceEvaluationResultsDataSource)
	setAggEvalDataSourceClient(ds, &connectivity.TencentCloudClient{})

	req := newAggEvalReadRequest(t, "")
	resp := &datasource.ReadResponse{}
	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "an account without results must not be reported as an error")

	var state svcconfig.ConfigListAggregateConfigRuleResourceEvaluationResultsModel
	resp.State.Get(context.Background(), &state)
	assert.True(t, state.AccountGroupId.IsNull())

	// An empty list should be represented as a non-null, zero-length list.
	assert.False(t, state.Items.IsNull())
	var items []svcconfig.AggregateConfigRuleResourceEvaluationResultItemModel
	diags := state.Items.ElementsAs(context.Background(), &items, false)
	assert.False(t, diags.HasError())
	assert.Empty(t, items)
}

// TestConfigListAggregateConfigRuleResourceEvaluationResultsDS_ReadNilAnnotation
// covers a result whose annotation is nil: the annotation field must be an
// empty list, not an error.
func TestConfigListAggregateConfigRuleResourceEvaluationResultsDS_ReadNilAnnotation(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	result := &configv20220802.ConfigRuleResourceEvaluationResult{
		RuleId:       ptrStringAggEval("cr-rule-002"),
		ResourceId:   ptrStringAggEval("ins-resource-002"),
		ResourceTags: nil,
		Annotation:   nil,
	}

	patches.ApplyMethodFunc(new(svcconfig.ConfigService), "DescribeConfigListAggregateConfigRuleResourceEvaluationResultsByFilter",
		func(_ context.Context, _ map[string]interface{}) ([]*configv20220802.ConfigRuleResourceEvaluationResult, error) {
			return []*configv20220802.ConfigRuleResourceEvaluationResult{result}, nil
		})

	ds := svcconfig.NewConfigListAggregateConfigRuleResourceEvaluationResultsDataSource().(*svcconfig.ConfigListAggregateConfigRuleResourceEvaluationResultsDataSource)
	setAggEvalDataSourceClient(ds, &connectivity.TencentCloudClient{})

	req := newAggEvalReadRequest(t, "ca-aggregator-abcdefg")
	resp := &datasource.ReadResponse{}
	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	var state svcconfig.ConfigListAggregateConfigRuleResourceEvaluationResultsModel
	resp.State.Get(context.Background(), &state)

	var items []svcconfig.AggregateConfigRuleResourceEvaluationResultItemModel
	diags := state.Items.ElementsAs(context.Background(), &items, false)
	assert.False(t, diags.HasError())
	assert.Len(t, items, 1)

	// Both nested lists must be empty (non-null, zero-length) for a nil response field.
	assert.False(t, items[0].ResourceTags.IsNull())
	var tags []svcconfig.ResourceTagItemModel
	assert.False(t, items[0].ResourceTags.ElementsAs(context.Background(), &tags, false).HasError())
	assert.Empty(t, tags)

	assert.False(t, items[0].Annotation.IsNull())
	var annotations []svcconfig.AnnotationItemModel
	assert.False(t, items[0].Annotation.ElementsAs(context.Background(), &annotations, false).HasError())
	assert.Empty(t, annotations)
}

// TestConfigListAggregateConfigRuleResourceEvaluationResultsDS_ReadAPIError
// verifies that an underlying API failure surfaces as an error diagnostic
// rather than being silently swallowed.
func TestConfigListAggregateConfigRuleResourceEvaluationResultsDS_ReadAPIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyMethodFunc(new(svcconfig.ConfigService), "DescribeConfigListAggregateConfigRuleResourceEvaluationResultsByFilter",
		func(_ context.Context, _ map[string]interface{}) ([]*configv20220802.ConfigRuleResourceEvaluationResult, error) {
			return nil, aggEvalAPIError("InternalError")
		})

	ds := svcconfig.NewConfigListAggregateConfigRuleResourceEvaluationResultsDataSource().(*svcconfig.ConfigListAggregateConfigRuleResourceEvaluationResultsDataSource)
	setAggEvalDataSourceClient(ds, &connectivity.TencentCloudClient{})

	req := newAggEvalReadRequest(t, "ca-aggregator-abcdefg")
	resp := &datasource.ReadResponse{}
	ds.Read(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError(), "an API failure must produce error diagnostics")
}

// TestConfigListAggregateConfigRuleResourceEvaluationResultsDS_Schema checks
// that the schema declares the expected attributes with the right optionality.
func TestConfigListAggregateConfigRuleResourceEvaluationResultsDS_Schema(t *testing.T) {
	s := aggEvalResultSchema(t)

	assert.Contains(t, s.Attributes, "account_group_id")
	assert.Contains(t, s.Attributes, "items")

	accountGroupAttr := s.Attributes["account_group_id"]
	assert.True(t, accountGroupAttr.IsOptional())

	itemsAttr := s.Attributes["items"]
	assert.True(t, itemsAttr.IsComputed())
}

// aggEvalAPIError is a minimal error type used to simulate an API failure.
type aggEvalAPIError string

func (e aggEvalAPIError) Error() string { return string(e) }
