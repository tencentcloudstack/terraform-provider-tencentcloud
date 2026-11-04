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

func ptrStrDSAggregate(s string) *string    { return &s }
func ptrUint64DSAggregate(v uint64) *uint64 { return &v }

// configListAggregateCompliancePacksSchema returns the framework data source
// schema used to build the raw config/state value.
func configListAggregateCompliancePacksSchema(t *testing.T) schema.Schema {
	d := svcconfig.NewConfigListAggregateCompliancePacksDataSource()
	resp := &datasource.SchemaResponse{}
	d.Schema(context.Background(), datasource.SchemaRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("failed to build data source schema: %v", resp.Diagnostics)
	}
	return resp.Schema
}

// configListAggregateCompliancePacksRootType builds the tftypes.Object type for
// the data source root, mirroring every schema attribute.
func configListAggregateCompliancePacksRootType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"account_group_id":                    tftypes.String,
			"compliance_pack_name":                tftypes.String,
			"risk_level":                          tftypes.List{ElementType: tftypes.Number},
			"status":                              tftypes.String,
			"compliance_result":                   tftypes.List{ElementType: tftypes.String},
			"order_type":                          tftypes.String,
			"include_compliance_pack_rule_result": tftypes.String,
			"total":                               tftypes.Number,
			"items":                               tftypes.List{ElementType: configAggregateCompliancePackItemType()},
		},
	}
}

// configAggregateCompliancePackItemType builds the tftypes.Object type for an
// items element.
func configAggregateCompliancePackItemType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"status":               tftypes.String,
			"compliance_pack_id":   tftypes.String,
			"description":          tftypes.String,
			"compliance_pack_name": tftypes.String,
			"risk_level":           tftypes.Number,
			"compliance_result":    tftypes.String,
			"create_time":          tftypes.String,
			"no_compliant_names":   tftypes.List{ElementType: tftypes.String},
			"rule_count":           tftypes.Number,
		},
	}
}

// buildAggregateCompliancePacksRawValue builds the raw config value with the
// given query inputs (only the required account_group_id must be set).
func buildAggregateCompliancePacksRawValue(accountGroupId string, riskLevels []float64, status, orderType, complianceResult string) tftypes.Value {
	attrs := map[string]tftypes.Value{
		"account_group_id":                    tftypes.NewValue(tftypes.String, accountGroupId),
		"compliance_pack_name":                tftypes.NewValue(tftypes.String, nil),
		"status":                              tftypes.NewValue(tftypes.String, nil),
		"order_type":                          tftypes.NewValue(tftypes.String, nil),
		"include_compliance_pack_rule_result": tftypes.NewValue(tftypes.String, nil),
		"total":                               tftypes.NewValue(tftypes.Number, nil),
	}

	// risk_level is a list of int64 (tftypes.Number for the protocol layer)
	riskLevelVals := make([]tftypes.Value, 0, len(riskLevels))
	for _, rl := range riskLevels {
		riskLevelVals = append(riskLevelVals, tftypes.NewValue(tftypes.Number, rl))
	}
	riskLevelType := tftypes.List{ElementType: tftypes.Number}
	if len(riskLevelVals) == 0 {
		attrs["risk_level"] = tftypes.NewValue(riskLevelType, nil)
	} else {
		attrs["risk_level"] = tftypes.NewValue(riskLevelType, riskLevelVals)
	}

	// compliance_result is a list of string
	if complianceResult == "" {
		attrs["compliance_result"] = tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, nil)
	} else {
		attrs["compliance_result"] = tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, []tftypes.Value{
			tftypes.NewValue(tftypes.String, complianceResult),
		})
	}

	// Build empty items list to avoid leaving state unparsed.
	itemsType := tftypes.List{ElementType: configAggregateCompliancePackItemType()}
	attrs["items"] = tftypes.NewValue(itemsType, nil)

	if status != "" {
		attrs["status"] = tftypes.NewValue(tftypes.String, status)
	}
	if orderType != "" {
		attrs["order_type"] = tftypes.NewValue(tftypes.String, orderType)
	}

	return tftypes.NewValue(configListAggregateCompliancePacksRootType(), attrs)
}

// setAggregateCompliancePacksClient injects a mock connectivity client into the
// data source instance via the Configure method.
func setAggregateCompliancePacksClient(d *svcconfig.ConfigListAggregateCompliancePacksDataSource, client *connectivity.TencentCloudClient) {
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := datasource.ConfigureRequest{ProviderData: meta}
	configureResp := &datasource.ConfigureResponse{}
	d.Configure(context.Background(), configureReq, configureResp)
}

// newAggregateCompliancePacksReadRequest builds a ReadRequest with the given
// config value and schema.
func newAggregateCompliancePacksReadRequest(sch schema.Schema, raw tftypes.Value) datasource.ReadRequest {
	req := datasource.ReadRequest{}
	req.Config.Raw = raw
	req.Config.Schema = sch
	return req
}

// newAggregateCompliancePacksReadResponse builds a ReadResponse whose State
// carries the schema, so resp.State.Set can reflect the model into a raw value.
func newAggregateCompliancePacksReadResponse(sch schema.Schema) *datasource.ReadResponse {
	resp := &datasource.ReadResponse{}
	resp.State.Schema = sch
	return resp
}

// go test ./tencentcloud/services/config/ -run "TestConfigListAggregateCompliancePacks" -v -count=1 -gcflags="all=-l"

// TestConfigListAggregateCompliancePacks_MetadataAndSchema validates metadata
// and schema definition.
func TestConfigListAggregateCompliancePacks_MetadataAndSchema(t *testing.T) {
	d := svcconfig.NewConfigListAggregateCompliancePacksDataSource()

	metaResp := &datasource.MetadataResponse{}
	d.Metadata(context.Background(), datasource.MetadataRequest{}, metaResp)
	assert.Equal(t, "tencentcloud_config_list_aggregate_compliance_packs", metaResp.TypeName)

	sch := configListAggregateCompliancePacksSchema(t)
	attrs := sch.Attributes

	// Required inputs
	assert.True(t, attrs["account_group_id"].(schema.StringAttribute).Required)

	// Optional inputs
	for _, opt := range []string{"compliance_pack_name", "status", "order_type", "include_compliance_pack_rule_result"} {
		attr, ok := attrs[opt].(schema.StringAttribute)
		assert.True(t, ok, "attribute %s should be StringAttribute", opt)
		assert.True(t, attr.Optional, "attribute %s should be Optional", opt)
	}

	// List optional inputs
	assert.True(t, attrs["risk_level"].(schema.ListAttribute).Optional)
	assert.True(t, attrs["compliance_result"].(schema.ListAttribute).Optional)

	// Computed outputs
	assert.True(t, attrs["total"].(schema.Int64Attribute).Computed)
	assert.True(t, attrs["items"].(schema.ListNestedAttribute).Computed)

	// Nested item attributes
	nested := attrs["items"].(schema.ListNestedAttribute).NestedObject.Attributes
	for _, leaf := range []string{"status", "compliance_pack_id", "description", "compliance_pack_name", "compliance_result", "create_time"} {
		nestedAttr, ok := nested[leaf].(schema.StringAttribute)
		assert.True(t, ok, "nested attribute %s should be StringAttribute", leaf)
		assert.True(t, nestedAttr.Computed, "nested attribute %s should be Computed", leaf)
	}
	assert.True(t, nested["risk_level"].(schema.Int64Attribute).Computed)
	assert.True(t, nested["rule_count"].(schema.Int64Attribute).Computed)
	assert.True(t, nested["no_compliant_names"].(schema.ListAttribute).Computed)
}

// TestConfigListAggregateCompliancePacks_ReadBasic tests a successful read with
// one compliance pack and asserts every leaf field is read back.
func TestConfigListAggregateCompliancePacks_ReadBasic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	configClient := &configv20220802.Client{}
	patches.ApplyMethodReturn(client, "UseConfigV20220802Client", configClient)

	patches.ApplyMethodFunc(configClient, "ListAggregateCompliancePacksWithContext", func(_ context.Context, request *configv20220802.ListAggregateCompliancePacksRequest) (*configv20220802.ListAggregateCompliancePacksResponse, error) {
		resp := configv20220802.NewListAggregateCompliancePacksResponse()
		resp.Response = &configv20220802.ListAggregateCompliancePacksResponseParams{
			Total: ptrUint64DSAggregate(1),
			Items: []*configv20220802.ConfigCompliancePack{
				{
					Status:             ptrStrDSAggregate("ACTIVE"),
					CompliancePackId:   ptrStrDSAggregate("cp-abcdefg"),
					Description:        ptrStrDSAggregate("tf test pack"),
					CompliancePackName: ptrStrDSAggregate("tf-aggregate-pack"),
					RiskLevel:          ptrUint64DSAggregate(1),
					ComplianceResult:   ptrStrDSAggregate("NON_COMPLIANT"),
					CreateTime:         ptrStrDSAggregate("2024-01-01 00:00:00"),
					NoCompliantNames:   []*string{ptrStrDSAggregate("rule-a"), ptrStrDSAggregate("rule-b")},
					RuleCount:          ptrUint64DSAggregate(2),
				},
			},
			RequestId: ptrStrDSAggregate("fake-request-id"),
		}
		return resp, nil
	})

	ds := svcconfig.NewConfigListAggregateCompliancePacksDataSource().(*svcconfig.ConfigListAggregateCompliancePacksDataSource)
	setAggregateCompliancePacksClient(ds, client)

	sch := configListAggregateCompliancePacksSchema(t)
	raw := buildAggregateCompliancePacksRawValue("ca-aggregator-abcdefg", nil, "", "", "")
	req := newAggregateCompliancePacksReadRequest(sch, raw)
	resp := newAggregateCompliancePacksReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	var result svcconfig.ConfigListAggregateCompliancePacksDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &result)...)
	assert.False(t, resp.Diagnostics.HasError())

	assert.Equal(t, int64(1), result.Total.ValueInt64())
	assert.False(t, result.Items.IsNull())

	var items []svcconfig.ConfigAggregateCompliancePackItemModel
	resp.Diagnostics.Append(result.Items.ElementsAs(context.Background(), &items, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, items, 1)

	// assert every leaf field
	assert.Equal(t, "ACTIVE", items[0].Status.ValueString())
	assert.Equal(t, "cp-abcdefg", items[0].CompliancePackId.ValueString())
	assert.Equal(t, "tf test pack", items[0].Description.ValueString())
	assert.Equal(t, "tf-aggregate-pack", items[0].CompliancePackName.ValueString())
	assert.Equal(t, int64(1), items[0].RiskLevel.ValueInt64())
	assert.Equal(t, "NON_COMPLIANT", items[0].ComplianceResult.ValueString())
	assert.Equal(t, "2024-01-01 00:00:00", items[0].CreateTime.ValueString())
	assert.Equal(t, int64(2), items[0].RuleCount.ValueInt64())

	// no_compliant_names set is read back
	assert.False(t, items[0].NoCompliantNames.IsNull())
	var ruleNames []types.String
	resp.Diagnostics.Append(items[0].NoCompliantNames.ElementsAs(context.Background(), &ruleNames, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, ruleNames, 2)
	assert.Equal(t, "rule-a", ruleNames[0].ValueString())
	assert.Equal(t, "rule-b", ruleNames[1].ValueString())
}

// TestConfigListAggregateCompliancePacks_ReadEmpty covers an account group that
// has no compliance packs. An empty result is a legitimate answer, not a
// failure: the data source must succeed with total = 0 and an empty items list.
func TestConfigListAggregateCompliancePacks_ReadEmpty(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	configClient := &configv20220802.Client{}
	patches.ApplyMethodReturn(client, "UseConfigV20220802Client", configClient)

	patches.ApplyMethodFunc(configClient, "ListAggregateCompliancePacksWithContext", func(_ context.Context, request *configv20220802.ListAggregateCompliancePacksRequest) (*configv20220802.ListAggregateCompliancePacksResponse, error) {
		resp := configv20220802.NewListAggregateCompliancePacksResponse()
		resp.Response = &configv20220802.ListAggregateCompliancePacksResponseParams{
			Total:     ptrUint64DSAggregate(0),
			Items:     []*configv20220802.ConfigCompliancePack{},
			RequestId: ptrStrDSAggregate("fake-request-id"),
		}
		return resp, nil
	})

	ds := svcconfig.NewConfigListAggregateCompliancePacksDataSource().(*svcconfig.ConfigListAggregateCompliancePacksDataSource)
	setAggregateCompliancePacksClient(ds, client)

	sch := configListAggregateCompliancePacksSchema(t)
	raw := buildAggregateCompliancePacksRawValue("ca-aggregator-empty", nil, "", "", "")
	req := newAggregateCompliancePacksReadRequest(sch, raw)
	resp := newAggregateCompliancePacksReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "an account group without compliance packs must not be reported as an error")

	var result svcconfig.ConfigListAggregateCompliancePacksDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &result)...)
	assert.False(t, resp.Diagnostics.HasError())

	assert.Equal(t, int64(0), result.Total.ValueInt64())
	assert.False(t, result.Items.IsNull())
	assert.Equal(t, 0, len(result.Items.Elements()))
}

// TestConfigListAggregateCompliancePacks_ReadNilItems covers the case where
// the API returns a nil items slice. The flatten logic must not panic.
func TestConfigListAggregateCompliancePacks_ReadNilItems(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	configClient := &configv20220802.Client{}
	patches.ApplyMethodReturn(client, "UseConfigV20220802Client", configClient)

	patches.ApplyMethodFunc(configClient, "ListAggregateCompliancePacksWithContext", func(_ context.Context, request *configv20220802.ListAggregateCompliancePacksRequest) (*configv20220802.ListAggregateCompliancePacksResponse, error) {
		resp := configv20220802.NewListAggregateCompliancePacksResponse()
		resp.Response = &configv20220802.ListAggregateCompliancePacksResponseParams{
			Total:     ptrUint64DSAggregate(0),
			Items:     nil,
			RequestId: ptrStrDSAggregate("fake-request-id"),
		}
		return resp, nil
	})

	ds := svcconfig.NewConfigListAggregateCompliancePacksDataSource().(*svcconfig.ConfigListAggregateCompliancePacksDataSource)
	setAggregateCompliancePacksClient(ds, client)

	sch := configListAggregateCompliancePacksSchema(t)
	raw := buildAggregateCompliancePacksRawValue("ca-aggregator-nil", nil, "", "", "")
	req := newAggregateCompliancePacksReadRequest(sch, raw)
	resp := newAggregateCompliancePacksReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError())

	var result svcconfig.ConfigListAggregateCompliancePacksDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &result)...)
	assert.False(t, resp.Diagnostics.HasError())

	assert.Equal(t, int64(0), result.Total.ValueInt64())
	assert.False(t, result.Items.IsNull())
	assert.Equal(t, 0, len(result.Items.Elements()))
}

// TestConfigListAggregateCompliancePacks_ReadNullFields covers a compliance pack
// whose nullable fields (Description, ComplianceResult, NoCompliantNames, RuleCount)
// are returned as nil by the API. They must be read back as null, not zero values.
func TestConfigListAggregateCompliancePacks_ReadNullFields(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	configClient := &configv20220802.Client{}
	patches.ApplyMethodReturn(client, "UseConfigV20220802Client", configClient)

	patches.ApplyMethodFunc(configClient, "ListAggregateCompliancePacksWithContext", func(_ context.Context, request *configv20220802.ListAggregateCompliancePacksRequest) (*configv20220802.ListAggregateCompliancePacksResponse, error) {
		resp := configv20220802.NewListAggregateCompliancePacksResponse()
		resp.Response = &configv20220802.ListAggregateCompliancePacksResponseParams{
			Total: ptrUint64DSAggregate(1),
			Items: []*configv20220802.ConfigCompliancePack{
				{
					Status:           ptrStrDSAggregate("ACTIVE"),
					CompliancePackId: ptrStrDSAggregate("cp-nullfields"),
					// Description nil
					CompliancePackName: ptrStrDSAggregate("tf-null-pack"),
					// RiskLevel nil
					// ComplianceResult nil
					CreateTime: ptrStrDSAggregate("2024-02-02 00:00:00"),
					// NoCompliantNames nil
					// RuleCount nil
				},
			},
			RequestId: ptrStrDSAggregate("fake-request-id"),
		}
		return resp, nil
	})

	ds := svcconfig.NewConfigListAggregateCompliancePacksDataSource().(*svcconfig.ConfigListAggregateCompliancePacksDataSource)
	setAggregateCompliancePacksClient(ds, client)

	sch := configListAggregateCompliancePacksSchema(t)
	raw := buildAggregateCompliancePacksRawValue("ca-aggregator-nullfields", nil, "", "", "")
	req := newAggregateCompliancePacksReadRequest(sch, raw)
	resp := newAggregateCompliancePacksReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError())

	var result svcconfig.ConfigListAggregateCompliancePacksDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &result)...)
	assert.False(t, resp.Diagnostics.HasError())

	assert.Equal(t, int64(1), result.Total.ValueInt64())

	var items []svcconfig.ConfigAggregateCompliancePackItemModel
	resp.Diagnostics.Append(result.Items.ElementsAs(context.Background(), &items, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, items, 1)

	// non-nil fields carry their values
	assert.Equal(t, "ACTIVE", items[0].Status.ValueString())
	assert.Equal(t, "cp-nullfields", items[0].CompliancePackId.ValueString())
	assert.Equal(t, "tf-null-pack", items[0].CompliancePackName.ValueString())
	assert.Equal(t, "2024-02-02 00:00:00", items[0].CreateTime.ValueString())

	// nullable fields are null, not zero
	assert.True(t, items[0].Description.IsNull(), "description should be null when API returns nil")
	assert.True(t, items[0].ComplianceResult.IsNull(), "compliance_result should be null when API returns nil")
	assert.True(t, items[0].RiskLevel.IsNull(), "risk_level should be null when API returns nil")
	assert.True(t, items[0].RuleCount.IsNull(), "rule_count should be null when API returns nil")
	assert.True(t, items[0].NoCompliantNames.IsNull(), "no_compliant_names should be null when API returns nil")
}

// TestConfigListAggregateCompliancePacks_ReadAPIError tests API error handling.
func TestConfigListAggregateCompliancePacks_ReadAPIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	configClient := &configv20220802.Client{}
	patches.ApplyMethodReturn(client, "UseConfigV20220802Client", configClient)

	patches.ApplyMethodFunc(configClient, "ListAggregateCompliancePacksWithContext", func(_ context.Context, request *configv20220802.ListAggregateCompliancePacksRequest) (*configv20220802.ListAggregateCompliancePacksResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InternalError, Message=db error")
	})

	ds := svcconfig.NewConfigListAggregateCompliancePacksDataSource().(*svcconfig.ConfigListAggregateCompliancePacksDataSource)
	setAggregateCompliancePacksClient(ds, client)

	sch := configListAggregateCompliancePacksSchema(t)
	raw := buildAggregateCompliancePacksRawValue("ca-aggregator-err", nil, "", "", "")
	req := newAggregateCompliancePacksReadRequest(sch, raw)
	resp := newAggregateCompliancePacksReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError())
	assert.Contains(t, resp.Diagnostics.Errors()[0].Summary(), "Error querying Config aggregate compliance packs")
	assert.Contains(t, resp.Diagnostics.Errors()[0].Detail(), "InternalError")
}

// TestConfigListAggregateCompliancePacks_ConfigureNilProviderData validates
// that Configure handles nil ProviderData gracefully.
func TestConfigListAggregateCompliancePacks_ConfigureNilProviderData(t *testing.T) {
	d := &svcconfig.ConfigListAggregateCompliancePacksDataSource{}
	req := datasource.ConfigureRequest{}
	resp := &datasource.ConfigureResponse{}

	d.Configure(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError())
}

// TestConfigListAggregateCompliancePacks_Factory validates that the registry
// factory returns a non-nil DataSource.
func TestConfigListAggregateCompliancePacks_Factory(t *testing.T) {
	got := svcconfig.NewConfigListAggregateCompliancePacksDataSource()
	assert.NotNil(t, got)
}

// TestConfigListAggregateCompliancePacks_QueriesNilFields tests that the
// ReadRequest reads back the null query inputs allocated at config time.
func TestConfigListAggregateCompliancePacks_ReadQueryStruct(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	configClient := &configv20220802.Client{}
	patches.ApplyMethodReturn(client, "UseConfigV20220802Client", configClient)

	var gotRequest *configv20220802.ListAggregateCompliancePacksRequest
	patches.ApplyMethodFunc(configClient, "ListAggregateCompliancePacksWithContext", func(_ context.Context, request *configv20220802.ListAggregateCompliancePacksRequest) (*configv20220802.ListAggregateCompliancePacksResponse, error) {
		gotRequest = request
		resp := configv20220802.NewListAggregateCompliancePacksResponse()
		resp.Response = &configv20220802.ListAggregateCompliancePacksResponseParams{
			Total:     ptrUint64DSAggregate(0),
			Items:     []*configv20220802.ConfigCompliancePack{},
			RequestId: ptrStrDSAggregate("fake-request-id"),
		}
		return resp, nil
	})

	ds := svcconfig.NewConfigListAggregateCompliancePacksDataSource().(*svcconfig.ConfigListAggregateCompliancePacksDataSource)
	setAggregateCompliancePacksClient(ds, client)

	// Build config with risk_level and status set
	raw := buildAggregateCompliancePacksRawValue("ca-aggregator-struct", []float64{1, 2}, "ACTIVE", "desc", "NON_COMPLIANT")
	sch := configListAggregateCompliancePacksSchema(t)
	req := newAggregateCompliancePacksReadRequest(sch, raw)
	resp := newAggregateCompliancePacksReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError())
	assert.NotNil(t, gotRequest)
	assert.Equal(t, "ca-aggregator-struct", *gotRequest.AccountGroupId)
	assert.Equal(t, "ACTIVE", *gotRequest.Status)
	assert.Equal(t, "desc", *gotRequest.OrderType)
	assert.Len(t, gotRequest.RiskLevel, 2)
	assert.Equal(t, uint64(1), *gotRequest.RiskLevel[0])
	assert.Equal(t, uint64(2), *gotRequest.RiskLevel[1])
	assert.Len(t, gotRequest.ComplianceResult, 1)
	assert.Equal(t, "NON_COMPLIANT", *gotRequest.ComplianceResult[0])
}
