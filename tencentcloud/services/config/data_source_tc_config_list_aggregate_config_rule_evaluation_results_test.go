package config_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	configv20220802 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/config"
)

func ptrStringAggEval(s string) *string { return &s }
func ptrUint64AggEval(v uint64) *uint64 { return &v }

// aggAnnotationTftypes returns the tftypes.Object describing an annotation element.
func aggAnnotationTftypes() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"configuration": tftypes.String,
			"desired_value": tftypes.String,
			"operator":      tftypes.String,
			"property":      tftypes.String,
		},
	}
}

// aggEvalItemTftypes returns the tftypes.Object describing an items element.
func aggEvalItemTftypes() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"resource_id":                 tftypes.String,
			"resource_type":              tftypes.String,
			"resource_region":           tftypes.String,
			"resource_name":              tftypes.String,
			"config_rule_id":            tftypes.String,
			"config_rule_name":         tftypes.String,
			"compliance_pack_id":      tftypes.String,
			"risk_level":              tftypes.Number,
			"annotation":               tftypes.List{ElementType: aggAnnotationTftypes()},
			"compliance_type":         tftypes.String,
			"invoking_event_message_type": tftypes.String,
			"config_rule_invoked_time":    tftypes.String,
			"result_recorded_time":       tftypes.String,
			"resource_owner_id":          tftypes.Number,
			"resource_owner_name":       tftypes.String,
		},
	}
}

func aggEvalSchema(t *testing.T) schema.Schema {
	ds := config.NewConfigListAggregateConfigRuleEvaluationResultsDataSource()
	schemaResp := &datasource.SchemaResponse{}
	ds.Schema(context.Background(), datasource.SchemaRequest{}, schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("failed to build data source schema: %v", schemaResp.Diagnostics)
	}
	return schemaResp.Schema
}

// newAggEvalReadRequest builds a datasource.ReadRequest with the given raw
// config values. configRuleId and accountGroupId are required; passing ""
// for the optional list/string params produces a null value.
func newAggEvalReadRequest(
	t *testing.T,
	configRuleId, accountGroupId string,
	resourceTypes, complianceTypes []string,
) datasource.ReadRequest {
	s := aggEvalSchema(t)

	var resourceTypesVal tftypes.Value
	if resourceTypes == nil {
		resourceTypesVal = tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, nil)
	} else {
		elems := make([]tftypes.Value, 0, len(resourceTypes))
		for _, rt := range resourceTypes {
			elems = append(elems, tftypes.NewValue(tftypes.String, rt))
		}
		resourceTypesVal = tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, elems)
	}

	var complianceTypesVal tftypes.Value
	if complianceTypes == nil {
		complianceTypesVal = tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, nil)
	} else {
		elems := make([]tftypes.Value, 0, len(complianceTypes))
		for _, ct := range complianceTypes {
			elems = append(elems, tftypes.NewValue(tftypes.String, ct))
		}
		complianceTypesVal = tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, elems)
	}

	raw := tftypes.NewValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"config_rule_id":     tftypes.String,
			"account_group_id":   tftypes.String,
			"resource_type":      tftypes.List{ElementType: tftypes.String},
			"compliance_type":    tftypes.List{ElementType: tftypes.String},
			"resource_owner_id":   tftypes.Number,
			"resource_owner_ids": tftypes.List{ElementType: tftypes.Number},
			"total":              tftypes.Number,
			"items":              tftypes.List{ElementType: aggEvalItemTftypes()},
		},
	}, map[string]tftypes.Value{
		"config_rule_id":     tftypes.NewValue(tftypes.String, configRuleId),
		"account_group_id":   tftypes.NewValue(tftypes.String, accountGroupId),
		"resource_type":      resourceTypesVal,
		"compliance_type":    complianceTypesVal,
		"resource_owner_id":  tftypes.NewValue(tftypes.Number, nil),
		"resource_owner_ids": tftypes.NewValue(tftypes.List{ElementType: tftypes.Number}, nil),
		"total":             tftypes.NewValue(tftypes.Number, nil),
		"items":            tftypes.NewValue(tftypes.List{ElementType: aggEvalItemTftypes()}, nil),
	})

	req := datasource.ReadRequest{}
	req.Config.Raw = raw
	req.Config.Schema = s
	return req
}

// setAggEvalClient injects a mock connectivity client into the data source
// instance via the Configure method.
func setAggEvalClient(ds *config.ConfigListAggregateConfigRuleEvaluationResultsDataSource, client *connectivity.TencentCloudClient) {
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := datasource.ConfigureRequest{ProviderData: meta}
	configureResp := &datasource.ConfigureResponse{}
	ds.Configure(context.Background(), configureReq, configureResp)
}

func buildAggEvaluationResult() *configv20220802.AggregateEvaluationResult {
	riskLevel := int64(1)
	resourceOwnerId := uint64(100000000000)
	return &configv20220802.AggregateEvaluationResult{
		ResourceId:                ptrStringAggEval("qcs::cvm::instance:ap-guangzhou:uin/100000000000:instance/ins-abc"),
		ResourceType:              ptrStringAggEval("QCS::CVM::Instance"),
		ResourceRegion:            ptrStringAggEval("ap-guangzhou"),
		ResourceName:             ptrStringAggEval("tf-instance"),
		ConfigRuleId:            ptrStringAggEval("cr-abcdefg1234"),
		ConfigRuleName:         ptrStringAggEval("tf-rule"),
		CompliancePackId:      ptrStringAggEval("cp-abcdefg1234"),
		RiskLevel:              &riskLevel,
		Annotation: &configv20220802.Annotation{
			Configuration: ptrStringAggEval("{\"disk\":50}"),
			DesiredValue:  ptrStringAggEval("{\"disk\":100}"),
			Operator:      ptrStringAggEval("less_than"),
			Property:      ptrStringAggEval("$.disk"),
		},
		ComplianceType:          ptrStringAggEval("NON_COMPLIANT"),
		InvokingEventMessageType: ptrStringAggEval("ScheduledNotification"),
		ConfigRuleInvokedTime:    ptrStringAggEval("2024-01-01 00:00:00"),
		ResultRecordedTime:       ptrStringAggEval("2024-01-01 00:00:01"),
		ResourceOwnerId:          &resourceOwnerId,
		ResourceOwnerName:       ptrStringAggEval("owner-name"),
	}
}

func newAggEvalResponse(items []*configv20220802.AggregateEvaluationResult) *configv20220802.ListAggregateConfigRuleEvaluationResultsResponse {
	resp := configv20220802.NewListAggregateConfigRuleEvaluationResultsResponse()
	resp.Response = &configv20220802.ListAggregateConfigRuleEvaluationResultsResponseParams{
		Total:    ptrUint64AggEval(uint64(len(items))),
		Items:    items,
		RequestId: ptrStringAggEval("fake-request-id"),
	}
	return resp
}

// TestConfigListAggregateConfigRuleEvaluationResultsDS_ReadBasic covers the
// happy path: one non-compliant evaluation result with a populated annotation.
func TestConfigListAggregateConfigRuleEvaluationResultsDS_ReadBasic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	configClient := &configv20220802.Client{}
	patches.ApplyMethodReturn(client, "UseConfigV20220802Client", configClient)

	patches.ApplyMethodFunc(configClient, "ListAggregateConfigRuleEvaluationResultsWithContext", func(_ context.Context, request *configv20220802.ListAggregateConfigRuleEvaluationResultsRequest) (*configv20220802.ListAggregateConfigRuleEvaluationResultsResponse, error) {
		// Required filters must reach the request.
		assert.Equal(t, "cr-abcdefg1234", *request.ConfigRuleId)
		assert.Equal(t, "ca-aggregator-abcdefg", *request.AccountGroupId)
		return newAggEvalResponse([]*configv20220802.AggregateEvaluationResult{buildAggEvaluationResult()}), nil
	})

	ds := config.NewConfigListAggregateConfigRuleEvaluationResultsDataSource().(*config.ConfigListAggregateConfigRuleEvaluationResultsDataSource)
	setAggEvalClient(ds, client)

	req := newAggEvalReadRequest(t, "cr-abcdefg1234", "ca-aggregator-abcdefg", nil, nil)
	resp := &datasource.ReadResponse{}
	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	var model config.ConfigListAggregateConfigRuleEvaluationResultsDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &model)...)
	assert.False(t, resp.Diagnostics.HasError(), "state read failed: %v", resp.Diagnostics)

	assert.Equal(t, int64(1), model.Total.ValueInt64())
	assert.False(t, model.Items.IsNull())
	assert.Equal(t, 1, len(model.Items.Elements()))

	items := make([]config.AggregateEvaluationResultItemModel, 0, 1)
	resp.Diagnostics.Append(model.Items.ElementsAs(context.Background(), &items, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	item := items[0]
	assert.Equal(t, "qcs::cvm::instance:ap-guangzhou:uin/100000000000:instance/ins-abc", item.ResourceId.ValueString())
	assert.Equal(t, "QCS::CVM::Instance", item.ResourceType.ValueString())
	assert.Equal(t, "ap-guangzhou", item.ResourceRegion.ValueString())
	assert.Equal(t, "tf-instance", item.ResourceName.ValueString())
	assert.Equal(t, "cr-abcdefg1234", item.ConfigRuleId.ValueString())
	assert.Equal(t, "tf-rule", item.ConfigRuleName.ValueString())
	assert.Equal(t, "cp-abcdefg1234", item.CompliancePackId.ValueString())
	assert.Equal(t, int64(1), item.RiskLevel.ValueInt64())
	assert.Equal(t, "NON_COMPLIANT", item.ComplianceType.ValueString())
	assert.Equal(t, "ScheduledNotification", item.InvokingEventMessageType.ValueString())
	assert.Equal(t, "2024-01-01 00:00:00", item.ConfigRuleInvokedTime.ValueString())
	assert.Equal(t, "2024-01-01 00:00:01", item.ResultRecordedTime.ValueString())
	assert.Equal(t, int64(100000000000), item.ResourceOwnerId.ValueInt64())
	assert.Equal(t, "owner-name", item.ResourceOwnerName.ValueString())

	// annotation is a single-element list
	assert.False(t, item.Annotation.IsNull())
	anns := make([]config.AggregateAnnotationItemModel, 0, 1)
	resp.Diagnostics.Append(item.Annotation.ElementsAs(context.Background(), &anns, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, anns, 1)
	assert.Equal(t, "{\"disk\":50}", anns[0].Configuration.ValueString())
	assert.Equal(t, "{\"disk\":100}", anns[0].DesiredValue.ValueString())
	assert.Equal(t, "less_than", anns[0].Operator.ValueString())
	assert.Equal(t, "$.disk", anns[0].Property.ValueString())
}

// TestConfigListAggregateConfigRuleEvaluationResultsDS_ReadEmpty covers an
// account group that has no evaluation results yet. An empty result is a
// legitimate answer, not a failure: the data source must succeed and expose
// total=0 with an empty items list.
func TestConfigListAggregateConfigRuleEvaluationResultsDS_ReadEmpty(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	configClient := &configv20220802.Client{}
	patches.ApplyMethodReturn(client, "UseConfigV20220802Client", configClient)

	patches.ApplyMethodFunc(configClient, "ListAggregateConfigRuleEvaluationResultsWithContext", func(_ context.Context, request *configv20220802.ListAggregateConfigRuleEvaluationResultsRequest) (*configv20220802.ListAggregateConfigRuleEvaluationResultsResponse, error) {
		return newAggEvalResponse([]*configv20220802.AggregateEvaluationResult{}), nil
	})

	ds := config.NewConfigListAggregateConfigRuleEvaluationResultsDataSource().(*config.ConfigListAggregateConfigRuleEvaluationResultsDataSource)
	setAggEvalClient(ds, client)

	req := newAggEvalReadRequest(t, "cr-abcdefg1234", "ca-aggregator-abcdefg", nil, nil)
	resp := &datasource.ReadResponse{}
	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "an account group without results must not be reported as an error")

	var model config.ConfigListAggregateConfigRuleEvaluationResultsDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &model)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Equal(t, int64(0), model.Total.ValueInt64())
	assert.False(t, model.Items.IsNull())
	assert.Equal(t, 0, len(model.Items.Elements()))
}

// TestConfigListAggregateConfigRuleEvaluationResultsDS_ReadNilItems covers
// the case where the API reports a total of 0 but returns a nil Items slice.
// The flatten logic must not panic and the data source must still succeed.
func TestConfigListAggregateConfigRuleEvaluationResultsDS_ReadNilItems(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	configClient := &configv20220802.Client{}
	patches.ApplyMethodReturn(client, "UseConfigV20220802Client", configClient)

	patches.ApplyMethodFunc(configClient, "ListAggregateConfigRuleEvaluationResultsWithContext", func(_ context.Context, request *configv20220802.ListAggregateConfigRuleEvaluationResultsRequest) (*configv20220802.ListAggregateConfigRuleEvaluationResultsResponse, error) {
		resp := configv20220802.NewListAggregateConfigRuleEvaluationResultsResponse()
		resp.Response = &configv20220802.ListAggregateConfigRuleEvaluationResultsResponseParams{
			Total:    ptrUint64AggEval(0),
			Items:    nil,
			RequestId: ptrStringAggEval("fake-request-id"),
		}
		return resp, nil
	})

	ds := config.NewConfigListAggregateConfigRuleEvaluationResultsDataSource().(*config.ConfigListAggregateConfigRuleEvaluationResultsDataSource)
	setAggEvalClient(ds, client)

	req := newAggEvalReadRequest(t, "cr-abcdefg1234", "ca-aggregator-abcdefg", nil, nil)
	resp := &datasource.ReadResponse{}
	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError())

	var model config.ConfigListAggregateConfigRuleEvaluationResultsDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &model)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Equal(t, int64(0), model.Total.ValueInt64())
	assert.False(t, model.Items.IsNull())
	assert.Equal(t, 0, len(model.Items.Elements()))
}

// TestConfigListAggregateConfigRuleEvaluationResultsDS_ReadNilAnnotation
// covers a result whose Annotation pointer is nil. The data source must emit
// a null annotation list rather than panic.
func TestConfigListAggregateConfigRuleEvaluationResultsDS_ReadNilAnnotation(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	configClient := &configv20220802.Client{}
	patches.ApplyMethodReturn(client, "UseConfigV20220802Client", configClient)

	result := buildAggEvaluationResult()
	result.Annotation = nil
	result.CompliancePackId = nil
	result.ComplianceType = nil
	result.ResourceOwnerId = nil
	result.ResourceOwnerName = nil

	patches.ApplyMethodFunc(configClient, "ListAggregateConfigRuleEvaluationResultsWithContext", func(_ context.Context, request *configv20220802.ListAggregateConfigRuleEvaluationResultsRequest) (*configv20220802.ListAggregateConfigRuleEvaluationResultsResponse, error) {
		return newAggEvalResponse([]*configv20220802.AggregateEvaluationResult{result}), nil
	})

	ds := config.NewConfigListAggregateConfigRuleEvaluationResultsDataSource().(*config.ConfigListAggregateConfigRuleEvaluationResultsDataSource)
	setAggEvalClient(ds, client)

	req := newAggEvalReadRequest(t, "cr-abcdefg1234", "ca-aggregator-abcdefg", nil, nil)
	resp := &datasource.ReadResponse{}
	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError())

	var model config.ConfigListAggregateConfigRuleEvaluationResultsDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &model)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Equal(t, int64(1), model.Total.ValueInt64())

	items := make([]config.AggregateEvaluationResultItemModel, 0, 1)
	resp.Diagnostics.Append(model.Items.ElementsAs(context.Background(), &items, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, items, 1)
	// nil Annotation pointer on the cloud side must become a null list, not a panic.
	assert.True(t, items[0].Annotation.IsNull())
	// nil nullable top-level fields must become null.
	assert.True(t, items[0].CompliancePackId.IsNull())
	assert.True(t, items[0].ComplianceType.IsNull())
	assert.True(t, items[0].ResourceOwnerId.IsNull())
	assert.True(t, items[0].ResourceOwnerName.IsNull())
}

// TestConfigListAggregateConfigRuleEvaluationResultsDS_ReadAPIError covers an
// API call failure: the data source must surface the error via diagnostics
// instead of returning an empty state.
func TestConfigListAggregateConfigRuleEvaluationResultsDS_ReadAPIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	configClient := &configv20220802.Client{}
	patches.ApplyMethodReturn(client, "UseConfigV20220802Client", configClient)

	patches.ApplyMethodFunc(configClient, "ListAggregateConfigRuleEvaluationResultsWithContext", func(_ context.Context, request *configv20220802.ListAggregateConfigRuleEvaluationResultsRequest) (*configv20220802.ListAggregateConfigRuleEvaluationResultsResponse, error) {
		return nil, assert.AnError
	})

	ds := config.NewConfigListAggregateConfigRuleEvaluationResultsDataSource().(*config.ConfigListAggregateConfigRuleEvaluationResultsDataSource)
	setAggEvalClient(ds, client)

	req := newAggEvalReadRequest(t, "cr-abcdefg1234", "ca-aggregator-abcdefg", nil, nil)
	resp := &datasource.ReadResponse{}
	ds.Read(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError(), "an API failure must be surfaced as a diagnostic error")
}