package config

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"
	sdkerr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	configv20220802 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/config/v20220802"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
)

// go test ./tencentcloud/services/config/ -run "TestConfigListAggregateCompliancePacksDataSource" -v -count=1 -gcflags="all=-l"

func ptrStrAggCompliancePacks(s string) *string {
	return &s
}

func ptrUint64AggCompliancePacks(u uint64) *uint64 {
	return &u
}

func newConfigListAggregateCompliancePacksDataSource(t *testing.T) *ConfigListAggregateCompliancePacksDataSource {
	t.Helper()
	d := NewConfigListAggregateCompliancePacksDataSource()
	ds, ok := d.(*ConfigListAggregateCompliancePacksDataSource)
	if !ok {
		t.Fatalf("factory returned %T, want *ConfigListAggregateCompliancePacksDataSource", d)
	}
	return ds
}

func configListAggregateCompliancePacksSchema(t *testing.T) schema.Schema {
	t.Helper()
	d := newConfigListAggregateCompliancePacksDataSource(t)
	schemaResp := &schema.SchemaResponse{}
	d.Schema(context.Background(), schema.SchemaRequest{}, schemaResp)
	if schemaResp.Diagnostics.HasError() {
		t.Fatalf("failed to build data source schema: %v", schemaResp.Diagnostics)
	}
	return schemaResp.Schema
}

// newConfigListAggregateCompliancePacksReadRequest builds a ReadRequest whose
// config carries the given account_group_id and optional filters.
func newConfigListAggregateCompliancePacksReadRequest(t *testing.T, filters map[string]tftypes.Value) datasource.ReadRequest {
	t.Helper()
	ctx := context.Background()
	s := configListAggregateCompliancePacksSchema(t)

	attrTypes := map[string]tftypes.Type{}
	valueMap := map[string]tftypes.Value{}
	for name, attr := range s.Attributes {
		tt := attr.Type().TerraformType(ctx)
		attrTypes[name] = tt
		valueMap[name] = tftypes.NewValue(tt, nil)
	}
	for name, v := range filters {
		valueMap[name] = v
	}

	raw := tftypes.NewValue(tftypes.Object{AttributeTypes: attrTypes}, valueMap)
	req := datasource.ReadRequest{}
	req.Config.Raw = raw
	req.Config.Schema = s
	return req
}

func setConfigListAggregateCompliancePacksClient(d *ConfigListAggregateCompliancePacksDataSource, client *connectivity.TencentCloudClient) {
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := datasource.ConfigureRequest{ProviderData: meta}
	configureResp := &datasource.ConfigureResponse{}
	d.Configure(context.Background(), configureReq, configureResp)
}

func mockAggCompliancePacksResponse(total uint64, items []*configv20220802.ConfigCompliancePack) *configv20220802.ListAggregateCompliancePacksResponse {
	resp := configv20220802.NewListAggregateCompliancePacksResponse()
	resp.Response = &configv20220802.ListAggregateCompliancePacksResponseParams{
		Total:     &total,
		Items:     items,
		RequestId: ptrStrAggCompliancePacks("fake-request-id"),
	}
	return resp
}

// TestConfigListAggregateCompliancePacksDataSource_Metadata validates the type name.
func TestConfigListAggregateCompliancePacksDataSource_Metadata(t *testing.T) {
	d := newConfigListAggregateCompliancePacksDataSource(t)
	resp := &datasource.MetadataResponse{}
	d.Metadata(context.Background(), datasource.MetadataRequest{}, resp)
	if resp.TypeName != "tencentcloud_config_list_aggregate_compliance_packs" {
		t.Errorf("Metadata.TypeName = %q, want %q", resp.TypeName, "tencentcloud_config_list_aggregate_compliance_packs")
	}
}

// TestConfigListAggregateCompliancePacksDataSource_Schema validates
// required/optional/computed flags on every attribute.
func TestConfigListAggregateCompliancePacksDataSource_Schema(t *testing.T) {
	s := configListAggregateCompliancePacksSchema(t)

	cases := map[string]struct{ required, optional, computed bool }{
		"account_group_id":                    {required: true},
		"compliance_pack_name":                {optional: true},
		"risk_level":                          {optional: true},
		"status":                              {optional: true},
		"compliance_result":                   {optional: true},
		"order_type":                          {optional: true},
		"include_compliance_pack_rule_result": {optional: true},
		"total":                               {computed: true},
		"items":                               {computed: true},
	}

	for name, want := range cases {
		attr, ok := s.Attributes[name]
		if !ok {
			t.Errorf("attribute %q missing from schema", name)
			continue
		}
		if attr.IsRequired() != want.required || attr.IsOptional() != want.optional || attr.IsComputed() != want.computed {
			t.Errorf("attribute %q flags (required=%v, optional=%v, computed=%v), want (required=%v, optional=%v, computed=%v)",
				name, attr.IsRequired(), attr.IsOptional(), attr.IsComputed(), want.required, want.optional, want.computed)
		}
	}
}

// TestConfigListAggregateCompliancePacksDataSource_Read_Success tests a
// successful read with filters and verifies all leaf fields are written to state.
func TestConfigListAggregateCompliancePacksDataSource_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	configClient := &configv20220802.Client{}
	patches.ApplyMethodReturn(client, "UseConfigV20220802Client", configClient)

	var gotRequest *configv20220802.ListAggregateCompliancePacksRequest
	patches.ApplyMethodFunc(configClient, "ListAggregateCompliancePacksWithContext", func(_ context.Context, request *configv20220802.ListAggregateCompliancePacksRequest) (*configv20220802.ListAggregateCompliancePacksResponse, error) {
		gotRequest = request
		items := []*configv20220802.ConfigCompliancePack{
			{
				Status:             ptrStrAggCompliancePacks("ACTIVE"),
				CompliancePackId:   ptrStrAggCompliancePacks("cp-1"),
				Description:        ptrStrAggCompliancePacks("pack one"),
				CompliancePackName: ptrStrAggCompliancePacks("Best-Practices"),
				RiskLevel:          ptrUint64AggCompliancePacks(1),
				ComplianceResult:   ptrStrAggCompliancePacks("NON_COMPLIANT"),
				CreateTime:         ptrStrAggCompliancePacks("2025-01-01 00:00:00"),
				NoCompliantNames:   []*string{ptrStrAggCompliancePacks("rule-a"), ptrStrAggCompliancePacks("rule-b")},
				RuleCount:          ptrUint64AggCompliancePacks(3),
			},
			{
				Status:             ptrStrAggCompliancePacks("NO_ACTIVE"),
				CompliancePackId:   ptrStrAggCompliancePacks("cp-2"),
				CompliancePackName: ptrStrAggCompliancePacks("Empty-Pack"),
				RiskLevel:          ptrUint64AggCompliancePacks(2),
			},
		}
		return mockAggCompliancePacksResponse(2, items), nil
	})

	d := newConfigListAggregateCompliancePacksDataSource(t)
	setConfigListAggregateCompliancePacksClient(d, client)

	req := newConfigListAggregateCompliancePacksReadRequest(t, map[string]tftypes.Value{
		"account_group_id":     tftypes.NewValue(tftypes.String, "aggr-123"),
		"compliance_pack_name": tftypes.NewValue(tftypes.String, "Best-Practices"),
		"risk_level": tftypes.NewValue(tftypes.Set{ElementType: tftypes.Number}, []tftypes.Value{
			tftypes.NewValue(tftypes.Number, 1),
		}),
		"status":                              tftypes.NewValue(tftypes.String, "ACTIVE"),
		"order_type":                          tftypes.NewValue(tftypes.String, "desc"),
		"include_compliance_pack_rule_result": tftypes.NewValue(tftypes.String, "NO"),
	})
	resp := &datasource.ReadResponse{}
	d.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	assert.NotNil(t, gotRequest)
	assert.Equal(t, "aggr-123", *gotRequest.AccountGroupId)
	assert.Equal(t, "Best-Practices", *gotRequest.CompliancePackName)
	assert.Len(t, gotRequest.RiskLevel, 1)
	assert.Equal(t, uint64(1), *gotRequest.RiskLevel[0])
	assert.Equal(t, "ACTIVE", *gotRequest.Status)
	assert.Equal(t, "desc", *gotRequest.OrderType)
	assert.Equal(t, "NO", *gotRequest.IncludeCompliancePackRuleResult)
	assert.Nil(t, gotRequest.ComplianceResult)

	var out ConfigListAggregateCompliancePacksDataSourceModel
	if diags := resp.State.Get(context.Background(), &out); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags)
	}
	assert.Equal(t, "aggr-123", out.AccountGroupId.ValueString())
	assert.Equal(t, int64(2), out.Total.ValueInt64())

	var packs []configListAggregateCompliancePacksItemModel
	if diags := out.Items.ElementsAs(context.Background(), &packs, false); diags.HasError() {
		t.Fatalf("failed to read items: %v", diags)
	}
	assert.Len(t, packs, 2)

	p0 := packs[0]
	assert.Equal(t, "ACTIVE", p0.Status.ValueString())
	assert.Equal(t, "cp-1", p0.CompliancePackId.ValueString())
	assert.Equal(t, "pack one", p0.Description.ValueString())
	assert.Equal(t, "Best-Practices", p0.CompliancePackName.ValueString())
	assert.Equal(t, int64(1), p0.RiskLevel.ValueInt64())
	assert.Equal(t, "NON_COMPLIANT", p0.ComplianceResult.ValueString())
	assert.Equal(t, "2025-01-01 00:00:00", p0.CreateTime.ValueString())
	assert.Equal(t, int64(3), p0.RuleCount.ValueInt64())

	var names []string
	if diags := p0.NoCompliantNames.ElementsAs(context.Background(), &names, false); diags.HasError() {
		t.Fatalf("failed to read no_compliant_names: %v", diags)
	}
	assert.Equal(t, []string{"rule-a", "rule-b"}, names)

	// Second pack has no description / result / create_time / names / count.
	p1 := packs[1]
	assert.Equal(t, "NO_ACTIVE", p1.Status.ValueString())
	assert.Equal(t, "cp-2", p1.CompliancePackId.ValueString())
	assert.True(t, p1.Description.IsNull(), "description should be null")
	assert.True(t, p1.ComplianceResult.IsNull(), "compliance_result should be null")
	assert.True(t, p1.CreateTime.IsNull(), "create_time should be null")
	assert.True(t, p1.RuleCount.IsNull(), "rule_count should be null")
	assert.True(t, p1.NoCompliantNames.IsNull(), "no_compliant_names should be null")
}

// TestConfigListAggregateCompliancePacksDataSource_Read_EmptyResult tests
// that an empty list is a valid result, not an error.
func TestConfigListAggregateCompliancePacksDataSource_Read_EmptyResult(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	configClient := &configv20220802.Client{}
	patches.ApplyMethodReturn(client, "UseConfigV20220802Client", configClient)
	patches.ApplyMethodFunc(configClient, "ListAggregateCompliancePacksWithContext", func(_ context.Context, _ *configv20220802.ListAggregateCompliancePacksRequest) (*configv20220802.ListAggregateCompliancePacksResponse, error) {
		return mockAggCompliancePacksResponse(0, nil), nil
	})

	d := newConfigListAggregateCompliancePacksDataSource(t)
	setConfigListAggregateCompliancePacksClient(d, client)

	req := newConfigListAggregateCompliancePacksReadRequest(t, map[string]tftypes.Value{
		"account_group_id": tftypes.NewValue(tftypes.String, "aggr-empty"),
	})
	resp := &datasource.ReadResponse{}
	d.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	var out ConfigListAggregateCompliancePacksDataSourceModel
	if diags := resp.State.Get(context.Background(), &out); diags.HasError() {
		t.Fatalf("failed to read back state: %v", diags)
	}
	assert.Equal(t, int64(0), out.Total.ValueInt64())
	var packs []configListAggregateCompliancePacksItemModel
	if diags := out.Items.ElementsAs(context.Background(), &packs, false); diags.HasError() {
		t.Fatalf("failed to read items: %v", diags)
	}
	assert.Len(t, packs, 0)
}

// TestConfigListAggregateCompliancePacksDataSource_Read_APIError tests API
// error propagation.
func TestConfigListAggregateCompliancePacksDataSource_Read_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	configClient := &configv20220802.Client{}
	patches.ApplyMethodReturn(client, "UseConfigV20220802Client", configClient)
	patches.ApplyMethodFunc(configClient, "ListAggregateCompliancePacksWithContext", func(_ context.Context, _ *configv20220802.ListAggregateCompliancePacksRequest) (*configv20220802.ListAggregateCompliancePacksResponse, error) {
		return nil, sdkerr.New("InvalidParameter", "invalid account group id")
	})

	d := newConfigListAggregateCompliancePacksDataSource(t)
	setConfigListAggregateCompliancePacksClient(d, client)

	req := newConfigListAggregateCompliancePacksReadRequest(t, map[string]tftypes.Value{
		"account_group_id": tftypes.NewValue(tftypes.String, "aggr-bad"),
	})
	resp := &datasource.ReadResponse{}
	d.Read(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError(), "expected error diagnostics for API failure")
}

// TestConfigListAggregateCompliancePacksDataSource_Read_NilResponse tests
// that a nil Response body is an API failure, not an empty result.
func TestConfigListAggregateCompliancePacksDataSource_Read_NilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	configClient := &configv20220802.Client{}
	patches.ApplyMethodReturn(client, "UseConfigV20220802Client", configClient)
	patches.ApplyMethodFunc(configClient, "ListAggregateCompliancePacksWithContext", func(_ context.Context, _ *configv20220802.ListAggregateCompliancePacksRequest) (*configv20220802.ListAggregateCompliancePacksResponse, error) {
		resp := configv20220802.NewListAggregateCompliancePacksResponse()
		resp.Response = nil
		return resp, nil
	})

	d := newConfigListAggregateCompliancePacksDataSource(t)
	setConfigListAggregateCompliancePacksClient(d, client)

	req := newConfigListAggregateCompliancePacksReadRequest(t, map[string]tftypes.Value{
		"account_group_id": tftypes.NewValue(tftypes.String, "aggr-123"),
	})
	resp := &datasource.ReadResponse{}
	d.Read(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError(), "expected error diagnostics for nil response")
}

// TestConfigListAggregateCompliancePacksDataSource_Read_NoClient tests the
// guard when Configure has not injected a client.
func TestConfigListAggregateCompliancePacksDataSource_Read_NoClient(t *testing.T) {
	d := newConfigListAggregateCompliancePacksDataSource(t)

	req := newConfigListAggregateCompliancePacksReadRequest(t, map[string]tftypes.Value{
		"account_group_id": tftypes.NewValue(tftypes.String, "aggr-123"),
	})
	resp := &datasource.ReadResponse{}
	d.Read(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError(), "expected error diagnostics when client is not configured")
}

// TestDescribeConfigListAggregateCompliancePacksByFilter_Pagination tests the
// pagination loop: two pages of 100 items each are collected into one slice.
func TestDescribeConfigListAggregateCompliancePacksByFilter_Pagination(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	configClient := &configv20220802.Client{}
	patches.ApplyMethodReturn(client, "UseConfigV20220802Client", configClient)

	callCount := 0
	patches.ApplyMethodFunc(configClient, "ListAggregateCompliancePacksWithContext", func(_ context.Context, request *configv20220802.ListAggregateCompliancePacksRequest) (*configv20220802.ListAggregateCompliancePacksResponse, error) {
		callCount++
		items := make([]*configv20220802.ConfigCompliancePack, 0, 100)
		for i := 0; i < 100; i++ {
			id := fmt.Sprintf("cp-%d-%d", callCount, i)
			items = append(items, &configv20220802.ConfigCompliancePack{
				CompliancePackId: ptrStrAggCompliancePacks(id),
			})
		}
		return mockAggCompliancePacksResponse(150, items), nil
	})

	service := NewConfigService(client)
	paramMap := map[string]interface{}{"AccountGroupId": "aggr-123"}
	items, total, err := service.DescribeConfigListAggregateCompliancePacksByFilter(context.Background(), paramMap)

	assert.NoError(t, err)
	assert.Equal(t, uint64(150), total)
	assert.Equal(t, 2, callCount)
	assert.Len(t, items, 150)
}

// TestDescribeConfigListAggregateCompliancePacksByFilter_SinglePage tests the
// loop stops after one page when fewer items than the page size are returned.
func TestDescribeConfigListAggregateCompliancePacksByFilter_SinglePage(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	configClient := &configv20220802.Client{}
	patches.ApplyMethodReturn(client, "UseConfigV20220802Client", configClient)

	callCount := 0
	patches.ApplyMethodFunc(configClient, "ListAggregateCompliancePacksWithContext", func(_ context.Context, _ *configv20220802.ListAggregateCompliancePacksRequest) (*configv20220802.ListAggregateCompliancePacksResponse, error) {
		callCount++
		items := []*configv20220802.ConfigCompliancePack{
			{CompliancePackId: ptrStrAggCompliancePacks("cp-1")},
		}
		return mockAggCompliancePacksResponse(1, items), nil
	})

	service := NewConfigService(client)
	paramMap := map[string]interface{}{
		"AccountGroupId":                  "aggr-123",
		"CompliancePackName":              "Best-Practices",
		"Status":                          "ACTIVE",
		"OrderType":                       "asc",
		"IncludeCompliancePackRuleResult": "NO",
	}
	items, total, err := service.DescribeConfigListAggregateCompliancePacksByFilter(context.Background(), paramMap)

	assert.NoError(t, err)
	assert.Equal(t, uint64(1), total)
	assert.Equal(t, 1, callCount)
	assert.Len(t, items, 1)
}

// TestDescribeConfigListAggregateCompliancePacksByFilter_Error tests the
// error path of the service method.
func TestDescribeConfigListAggregateCompliancePacksByFilter_Error(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	configClient := &configv20220802.Client{}
	patches.ApplyMethodReturn(client, "UseConfigV20220802Client", configClient)
	patches.ApplyMethodFunc(configClient, "ListAggregateCompliancePacksWithContext", func(_ context.Context, _ *configv20220802.ListAggregateCompliancePacksRequest) (*configv20220802.ListAggregateCompliancePacksResponse, error) {
		return nil, sdkerr.New("AuthFailure", "auth failed")
	})

	service := NewConfigService(client)
	paramMap := map[string]interface{}{"AccountGroupId": "aggr-123"}
	items, total, err := service.DescribeConfigListAggregateCompliancePacksByFilter(context.Background(), paramMap)

	assert.Error(t, err)
	assert.Equal(t, uint64(0), total)
	assert.Len(t, items, 0)
}
