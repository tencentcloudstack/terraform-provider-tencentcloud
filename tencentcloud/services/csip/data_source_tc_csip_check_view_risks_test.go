package csip_test

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

	csipv20221121 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/csip/v20221121"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	svcsip "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/csip"
)

func csipCheckViewRisksStr(s string) *string    { return &s }
func csipCheckViewRisksUint64(v uint64) *uint64 { return &v }
func csipCheckViewRisksInt64(v int64) *int64    { return &v }

// csipCheckViewRisksSchema returns the framework data source schema used to
// build the raw config/state value.
func csipCheckViewRisksSchema(t *testing.T) schema.Schema {
	d := svcsip.NewCsipCheckViewRisksDataSource()
	resp := &datasource.SchemaResponse{}
	d.Schema(context.Background(), datasource.SchemaRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("failed to build data source schema: %v", resp.Diagnostics)
	}

	return resp.Schema
}

// csipCheckViewRisksFilterType builds the tftypes.Object type for a filters
// element.
func csipCheckViewRisksFilterType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"name":        tftypes.String,
			"values":      tftypes.List{ElementType: tftypes.String},
			"exact_match": tftypes.String,
		},
	}
}

// csipCheckViewRisksStandardTermType builds the tftypes.Object type for a
// standard_terms element.
func csipCheckViewRisksStandardTermType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"tag":   tftypes.String,
			"terms": tftypes.List{ElementType: tftypes.String},
		},
	}
}

// csipCheckViewRiskItemType builds the tftypes.Object type for a
// check_view_risk_list element.
func csipCheckViewRiskItemType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"risk_rule_id":        tftypes.String,
			"risk_title":          tftypes.String,
			"check_type":          tftypes.String,
			"severity":            tftypes.String,
			"risk_desc":           tftypes.String,
			"create_time":         tftypes.String,
			"update_time":         tftypes.String,
			"provider":            tftypes.String,
			"risk_status":         tftypes.Number,
			"asset_count":         tftypes.Number,
			"risk_count":          tftypes.Number,
			"asset_type":          tftypes.String,
			"event_type":          tftypes.String,
			"classify":            tftypes.String,
			"standard_terms":      tftypes.List{ElementType: csipCheckViewRisksStandardTermType()},
			"asset_type_icon_url": tftypes.String,
		},
	}
}

// csipCheckViewRisksStandardNameItemType builds the tftypes.Object type for a
// standard_name_list element.
func csipCheckViewRisksStandardNameItemType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"id":   tftypes.Number,
			"name": tftypes.String,
		},
	}
}

// csipCheckViewRisksAttributeOptionItemType builds the tftypes.Object type for
// an asset_type_list / provider_list element.
func csipCheckViewRisksAttributeOptionItemType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"text":  tftypes.String,
			"value": tftypes.String,
		},
	}
}

// csipCheckViewRisksRootType builds the tftypes.Object type for the data
// source root, mirroring every schema attribute.
func csipCheckViewRisksRootType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"id":                   tftypes.String,
			"member_id":            tftypes.List{ElementType: tftypes.String},
			"filters":              tftypes.List{ElementType: csipCheckViewRisksFilterType()},
			"order":                tftypes.String,
			"by":                   tftypes.String,
			"check_view_risk_list": tftypes.List{ElementType: csipCheckViewRiskItemType()},
			"standard_name_list":   tftypes.List{ElementType: csipCheckViewRisksStandardNameItemType()},
			"asset_type_list":      tftypes.List{ElementType: csipCheckViewRisksAttributeOptionItemType()},
			"provider_list":        tftypes.List{ElementType: csipCheckViewRisksAttributeOptionItemType()},
		},
	}
}

// buildCsipCheckViewRisksRawValue builds the raw config value with the given
// query inputs; every other attribute is unset (nil).
func buildCsipCheckViewRisksRawValue(order, by string, memberIds []string, withFilter bool) tftypes.Value {
	attrs := map[string]tftypes.Value{
		"id":                   tftypes.NewValue(tftypes.String, nil),
		"order":                tftypes.NewValue(tftypes.String, order),
		"by":                   tftypes.NewValue(tftypes.String, by),
		"member_id":            tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, nil),
		"filters":              tftypes.NewValue(tftypes.List{ElementType: csipCheckViewRisksFilterType()}, nil),
		"check_view_risk_list": tftypes.NewValue(tftypes.List{ElementType: csipCheckViewRiskItemType()}, nil),
		"standard_name_list":   tftypes.NewValue(tftypes.List{ElementType: csipCheckViewRisksStandardNameItemType()}, nil),
		"asset_type_list":      tftypes.NewValue(tftypes.List{ElementType: csipCheckViewRisksAttributeOptionItemType()}, nil),
		"provider_list":        tftypes.NewValue(tftypes.List{ElementType: csipCheckViewRisksAttributeOptionItemType()}, nil),
	}

	if len(memberIds) > 0 {
		memberIdVals := make([]tftypes.Value, 0, len(memberIds))
		for _, id := range memberIds {
			memberIdVals = append(memberIdVals, tftypes.NewValue(tftypes.String, id))
		}
		attrs["member_id"] = tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, memberIdVals)
	}

	if withFilter {
		filter := tftypes.NewValue(csipCheckViewRisksFilterType(), map[string]tftypes.Value{
			"name":        tftypes.NewValue(tftypes.String, "Status"),
			"values":      tftypes.NewValue(tftypes.List{ElementType: tftypes.String}, []tftypes.Value{tftypes.NewValue(tftypes.String, "CHECKING")}),
			"exact_match": tftypes.NewValue(tftypes.String, "1"),
		})
		attrs["filters"] = tftypes.NewValue(tftypes.List{ElementType: csipCheckViewRisksFilterType()}, []tftypes.Value{filter})
	}

	return tftypes.NewValue(csipCheckViewRisksRootType(), attrs)
}

// setCsipCheckViewRisksClient injects a mock connectivity client into the data
// source instance via the Configure method.
func setCsipCheckViewRisksClient(d *svcsip.CsipCheckViewRisksDataSource, client *connectivity.TencentCloudClient) {
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := datasource.ConfigureRequest{ProviderData: meta}
	configureResp := &datasource.ConfigureResponse{}
	d.Configure(context.Background(), configureReq, configureResp)
}

// newCsipCheckViewRisksReadRequest builds a ReadRequest with the given config
// value and schema.
func newCsipCheckViewRisksReadRequest(sch schema.Schema, raw tftypes.Value) datasource.ReadRequest {
	req := datasource.ReadRequest{}
	req.Config.Raw = raw
	req.Config.Schema = sch

	return req
}

// newCsipCheckViewRisksReadResponse builds a ReadResponse whose State carries
// the schema, so resp.State.Set can reflect the model into a raw value.
func newCsipCheckViewRisksReadResponse(sch schema.Schema) *datasource.ReadResponse {
	resp := &datasource.ReadResponse{}
	resp.State.Schema = sch

	return resp
}

// go test ./tencentcloud/services/csip/ -run "TestCsipCheckViewRisks" -v -count=1 -gcflags="all=-l"

// TestCsipCheckViewRisks_MetadataAndSchema validates metadata and schema
// definition.
func TestCsipCheckViewRisks_MetadataAndSchema(t *testing.T) {
	d := svcsip.NewCsipCheckViewRisksDataSource()

	metaResp := &datasource.MetadataResponse{}
	d.Metadata(context.Background(), datasource.MetadataRequest{}, metaResp)
	assert.Equal(t, "tencentcloud_csip_check_view_risks", metaResp.TypeName)

	sch := csipCheckViewRisksSchema(t)
	attrs := sch.Attributes

	// Optional inputs
	assert.True(t, attrs["member_id"].(schema.ListAttribute).Optional)
	assert.True(t, attrs["filters"].(schema.ListNestedAttribute).Optional)
	assert.True(t, attrs["order"].(schema.StringAttribute).Optional)
	assert.True(t, attrs["by"].(schema.StringAttribute).Optional)

	// Computed outputs
	assert.True(t, attrs["id"].(schema.StringAttribute).Computed)
	assert.True(t, attrs["check_view_risk_list"].(schema.ListNestedAttribute).Computed)
	assert.True(t, attrs["standard_name_list"].(schema.ListNestedAttribute).Computed)
	assert.True(t, attrs["asset_type_list"].(schema.ListNestedAttribute).Computed)
	assert.True(t, attrs["provider_list"].(schema.ListNestedAttribute).Computed)

	// Nested filter attributes
	filterNested := attrs["filters"].(schema.ListNestedAttribute).NestedObject.Attributes
	assert.True(t, filterNested["name"].(schema.StringAttribute).Required)
	assert.True(t, filterNested["values"].(schema.ListAttribute).Required)
	assert.True(t, filterNested["exact_match"].(schema.StringAttribute).Optional)

	// Nested risk item attributes
	riskNested := attrs["check_view_risk_list"].(schema.ListNestedAttribute).NestedObject.Attributes
	for _, leaf := range []string{"risk_rule_id", "risk_title", "check_type", "severity", "risk_desc", "create_time", "update_time", "provider", "asset_type", "event_type", "classify", "asset_type_icon_url"} {
		nestedAttr, ok := riskNested[leaf].(schema.StringAttribute)
		assert.True(t, ok, "nested attribute %s should be StringAttribute", leaf)
		assert.True(t, nestedAttr.Computed, "nested attribute %s should be Computed", leaf)
	}
	for _, leaf := range []string{"risk_status", "asset_count", "risk_count"} {
		nestedAttr, ok := riskNested[leaf].(schema.Int64Attribute)
		assert.True(t, ok, "nested attribute %s should be Int64Attribute", leaf)
		assert.True(t, nestedAttr.Computed, "nested attribute %s should be Computed", leaf)
	}
	assert.True(t, riskNested["standard_terms"].(schema.ListNestedAttribute).Computed)

	// Paging knobs and total-count aggregates are NEVER exposed.
	assert.NotContains(t, attrs, "limit")
	assert.NotContains(t, attrs, "offset")
	assert.NotContains(t, attrs, "total_count")
}

// TestCsipCheckViewRisks_ReadBasic tests a successful read with one record
// and asserts every leaf field is read back.
func TestCsipCheckViewRisks_ReadBasic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	csipClient := &csipv20221121.Client{}
	patches.ApplyMethodReturn(client, "UseCsipClient", csipClient)
	patches.ApplyMethodFunc(csipClient, "DescribeCheckViewRisks", func(request *csipv20221121.DescribeCheckViewRisksRequest) (*csipv20221121.DescribeCheckViewRisksResponse, error) {
		if request.Order == nil || *request.Order != "desc" {
			t.Errorf("expected order to be passed to the API request")
		}
		if request.By == nil || *request.By != "update_time" {
			t.Errorf("expected by to be passed to the API request")
		}
		if len(request.MemberId) != 1 || *request.MemberId[0] != "mem-6wfo0fzks3" {
			t.Errorf("expected member_id to be passed to the API request, got %v", request.MemberId)
		}
		if len(request.Filters) != 1 {
			t.Fatalf("expected 1 filter to be passed to the API request, got %d", len(request.Filters))
		}
		if request.Filters[0].Name == nil || *request.Filters[0].Name != "Status" {
			t.Errorf("expected filter name to be passed to the API request")
		}
		if len(request.Filters[0].Values) != 1 || *request.Filters[0].Values[0] != "CHECKING" {
			t.Errorf("expected filter values to be passed to the API request")
		}
		if request.Filters[0].ExactMatch == nil || *request.Filters[0].ExactMatch != "1" {
			t.Errorf("expected filter exact_match to be passed to the API request")
		}

		return &csipv20221121.DescribeCheckViewRisksResponse{Response: &csipv20221121.DescribeCheckViewRisksResponseParams{
			CheckViewRiskList: []*csipv20221121.CheckViewRiskItem{
				{
					RiskRuleId: csipCheckViewRisksStr("tc-001"),
					RiskTitle:  csipCheckViewRisksStr("tf-test-risk-title"),
					CheckType:  csipCheckViewRisksStr("网络安全"),
					Severity:   csipCheckViewRisksStr("high"),
					RiskDesc:   csipCheckViewRisksStr("tf-test-risk-desc"),
					CreateTime: csipCheckViewRisksStr("2025-05-05 12:00:00"),
					UpdateTime: csipCheckViewRisksStr("2025-05-06 12:00:00"),
					Provider:   csipCheckViewRisksStr("tencent"),
					RiskStatus: csipCheckViewRisksUint64(0),
					AssetCount: csipCheckViewRisksUint64(1),
					RiskCount:  csipCheckViewRisksUint64(2),
					AssetType:  csipCheckViewRisksStr("安全组"),
					EventType:  csipCheckViewRisksStr("不合理规则"),
					Classify:   csipCheckViewRisksStr("expose"),
					StandardTerms: []*csipv20221121.StandardTerm{
						{
							Tag:   csipCheckViewRisksStr("等保三级"),
							Terms: []*string{csipCheckViewRisksStr("term-a"), csipCheckViewRisksStr("term-b")},
						},
					},
					AssetTypeIconURL: csipCheckViewRisksStr("https://example.com/icon.png"),
				},
			},
			StandardNameList: []*csipv20221121.StandardItem{
				{
					ID:   csipCheckViewRisksUint64(1),
					Name: csipCheckViewRisksStr("等保三级"),
				},
			},
			AssetTypeList: []*csipv20221121.AttributeOptionSet{
				{
					Text:  csipCheckViewRisksStr("cvm_instance"),
					Value: csipCheckViewRisksStr("CVM"),
				},
			},
			ProviderList: []*csipv20221121.AttributeOptionSet{
				{
					Text:  csipCheckViewRisksStr("tencent_cloud"),
					Value: csipCheckViewRisksStr("Tencent Cloud"),
				},
			},
			RequestId: csipCheckViewRisksStr("fake-request-id"),
		}}, nil
	})

	ds := svcsip.NewCsipCheckViewRisksDataSource().(*svcsip.CsipCheckViewRisksDataSource)
	setCsipCheckViewRisksClient(ds, client)

	sch := csipCheckViewRisksSchema(t)
	raw := buildCsipCheckViewRisksRawValue("desc", "update_time", []string{"mem-6wfo0fzks3"}, true)
	req := newCsipCheckViewRisksReadRequest(sch, raw)
	resp := newCsipCheckViewRisksReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	var result svcsip.CsipCheckViewRisksDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &result)...)
	assert.False(t, resp.Diagnostics.HasError())

	// The state id is a random token.
	assert.NotEmpty(t, result.Id.ValueString())

	// check_view_risk_list
	assert.False(t, result.CheckViewRiskList.IsNull())
	var riskItems []svcsip.CsipCheckViewRiskItemModel
	resp.Diagnostics.Append(result.CheckViewRiskList.ElementsAs(context.Background(), &riskItems, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, riskItems, 1)

	assert.Equal(t, "tc-001", riskItems[0].RiskRuleId.ValueString())
	assert.Equal(t, "tf-test-risk-title", riskItems[0].RiskTitle.ValueString())
	assert.Equal(t, "网络安全", riskItems[0].CheckType.ValueString())
	assert.Equal(t, "high", riskItems[0].Severity.ValueString())
	assert.Equal(t, "tf-test-risk-desc", riskItems[0].RiskDesc.ValueString())
	assert.Equal(t, "2025-05-05 12:00:00", riskItems[0].CreateTime.ValueString())
	assert.Equal(t, "2025-05-06 12:00:00", riskItems[0].UpdateTime.ValueString())
	assert.Equal(t, "tencent", riskItems[0].Provider.ValueString())
	assert.Equal(t, int64(0), riskItems[0].RiskStatus.ValueInt64())
	assert.Equal(t, int64(1), riskItems[0].AssetCount.ValueInt64())
	assert.Equal(t, int64(2), riskItems[0].RiskCount.ValueInt64())
	assert.Equal(t, "安全组", riskItems[0].AssetType.ValueString())
	assert.Equal(t, "不合理规则", riskItems[0].EventType.ValueString())
	assert.Equal(t, "expose", riskItems[0].Classify.ValueString())
	assert.Equal(t, "https://example.com/icon.png", riskItems[0].AssetTypeIconUrl.ValueString())

	// standard_terms nested list
	assert.False(t, riskItems[0].StandardTerms.IsNull())
	var standardTerms []svcsip.CsipCheckViewRiskStandardTermModel
	resp.Diagnostics.Append(riskItems[0].StandardTerms.ElementsAs(context.Background(), &standardTerms, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, standardTerms, 1)
	assert.Equal(t, "等保三级", standardTerms[0].Tag.ValueString())
	assert.False(t, standardTerms[0].Terms.IsNull())
	var terms []types.String
	resp.Diagnostics.Append(standardTerms[0].Terms.ElementsAs(context.Background(), &terms, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, terms, 2)
	assert.Equal(t, "term-a", terms[0].ValueString())
	assert.Equal(t, "term-b", terms[1].ValueString())

	// standard_name_list
	assert.False(t, result.StandardNameList.IsNull())
	var standardNameItems []svcsip.CsipCheckViewStandardNameItemModel
	resp.Diagnostics.Append(result.StandardNameList.ElementsAs(context.Background(), &standardNameItems, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, standardNameItems, 1)
	assert.Equal(t, int64(1), standardNameItems[0].Id.ValueInt64())
	assert.Equal(t, "等保三级", standardNameItems[0].Name.ValueString())

	// asset_type_list
	assert.False(t, result.AssetTypeList.IsNull())
	var assetTypeItems []svcsip.CsipCheckViewAttributeOptionItemModel
	resp.Diagnostics.Append(result.AssetTypeList.ElementsAs(context.Background(), &assetTypeItems, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, assetTypeItems, 1)
	assert.Equal(t, "cvm_instance", assetTypeItems[0].Text.ValueString())
	assert.Equal(t, "CVM", assetTypeItems[0].Value.ValueString())

	// provider_list
	assert.False(t, result.ProviderList.IsNull())
	var providerItems []svcsip.CsipCheckViewAttributeOptionItemModel
	resp.Diagnostics.Append(result.ProviderList.ElementsAs(context.Background(), &providerItems, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, providerItems, 1)
	assert.Equal(t, "tencent_cloud", providerItems[0].Text.ValueString())
	assert.Equal(t, "Tencent Cloud", providerItems[0].Value.ValueString())
}

// TestCsipCheckViewRisks_ReadEmpty covers an empty result: a legitimate
// answer, not a failure — the data source must succeed with empty lists.
func TestCsipCheckViewRisks_ReadEmpty(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	csipClient := &csipv20221121.Client{}
	patches.ApplyMethodReturn(client, "UseCsipClient", csipClient)
	patches.ApplyMethodFunc(csipClient, "DescribeCheckViewRisks", func(request *csipv20221121.DescribeCheckViewRisksRequest) (*csipv20221121.DescribeCheckViewRisksResponse, error) {
		return &csipv20221121.DescribeCheckViewRisksResponse{Response: &csipv20221121.DescribeCheckViewRisksResponseParams{
			CheckViewRiskList: []*csipv20221121.CheckViewRiskItem{},
			RequestId:         csipCheckViewRisksStr("fake-request-id"),
		}}, nil
	})

	ds := svcsip.NewCsipCheckViewRisksDataSource().(*svcsip.CsipCheckViewRisksDataSource)
	setCsipCheckViewRisksClient(ds, client)

	sch := csipCheckViewRisksSchema(t)
	raw := buildCsipCheckViewRisksRawValue("desc", "update_time", nil, false)
	req := newCsipCheckViewRisksReadRequest(sch, raw)
	resp := newCsipCheckViewRisksReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	var result svcsip.CsipCheckViewRisksDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &result)...)
	assert.False(t, resp.Diagnostics.HasError())

	assert.False(t, result.CheckViewRiskList.IsNull())
	assert.Equal(t, 0, len(result.CheckViewRiskList.Elements()))
	assert.False(t, result.StandardNameList.IsNull())
	assert.Equal(t, 0, len(result.StandardNameList.Elements()))
	assert.False(t, result.AssetTypeList.IsNull())
	assert.Equal(t, 0, len(result.AssetTypeList.Elements()))
	assert.False(t, result.ProviderList.IsNull())
	assert.Equal(t, 0, len(result.ProviderList.Elements()))
	assert.NotEmpty(t, result.Id.ValueString())
}

// TestCsipCheckViewRisks_ReadApiError covers an API-level failure: the error
// must surface through diagnostics instead of silently persisting empty
// state.
func TestCsipCheckViewRisks_ReadApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	csipClient := &csipv20221121.Client{}
	patches.ApplyMethodReturn(client, "UseCsipClient", csipClient)
	patches.ApplyMethodFunc(csipClient, "DescribeCheckViewRisks", func(request *csipv20221121.DescribeCheckViewRisksRequest) (*csipv20221121.DescribeCheckViewRisksResponse, error) {
		return nil, fmt.Errorf("mocked API error")
	})

	ds := svcsip.NewCsipCheckViewRisksDataSource().(*svcsip.CsipCheckViewRisksDataSource)
	setCsipCheckViewRisksClient(ds, client)

	sch := csipCheckViewRisksSchema(t)
	raw := buildCsipCheckViewRisksRawValue("desc", "update_time", nil, false)
	req := newCsipCheckViewRisksReadRequest(sch, raw)
	resp := newCsipCheckViewRisksReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError(), "expected error diagnostics for API failure")
}

// TestCsipCheckViewRisksService_ReadAllPages verifies that the service layer
// loops through every page with the fixed page size and returns the complete
// set (the total count stays internal).
func TestCsipCheckViewRisksService_ReadAllPages(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	csipClient := &csipv20221121.Client{}
	patches.ApplyMethodReturn(client, "UseCsipClient", csipClient)

	callCount := 0
	patches.ApplyMethodFunc(csipClient, "DescribeCheckViewRisks", func(request *csipv20221121.DescribeCheckViewRisksRequest) (*csipv20221121.DescribeCheckViewRisksResponse, error) {
		callCount++
		// The service must drive pagination itself with the page size 100.
		if request.Limit == nil || *request.Limit != 100 {
			t.Errorf("expected Limit 100 on call %d, got %v", callCount, request.Limit)
		}
		var offset uint64
		if request.Offset != nil {
			offset = *request.Offset
		}
		switch callCount {
		case 1:
			if offset != 0 {
				t.Errorf("expected offset 0 on first call, got %d", offset)
			}
			riskList := make([]*csipv20221121.CheckViewRiskItem, 0, 100)
			for i := 0; i < 100; i++ {
				riskList = append(riskList, &csipv20221121.CheckViewRiskItem{
					RiskRuleId: csipCheckViewRisksStr(fmt.Sprintf("tc-page1-%d", i)),
					Severity:   csipCheckViewRisksStr("high"),
				})
			}

			return &csipv20221121.DescribeCheckViewRisksResponse{Response: &csipv20221121.DescribeCheckViewRisksResponseParams{
				TotalCount:        csipCheckViewRisksInt64(101),
				CheckViewRiskList: riskList,
				StandardNameList: []*csipv20221121.StandardItem{
					{ID: csipCheckViewRisksUint64(1), Name: csipCheckViewRisksStr("等保三级")},
				},
				RequestId: csipCheckViewRisksStr("fake-request-id"),
			}}, nil
		case 2:
			if offset != 100 {
				t.Errorf("expected offset 100 on second call, got %d", offset)
			}

			return &csipv20221121.DescribeCheckViewRisksResponse{Response: &csipv20221121.DescribeCheckViewRisksResponseParams{
				TotalCount: csipCheckViewRisksInt64(101),
				CheckViewRiskList: []*csipv20221121.CheckViewRiskItem{
					{RiskRuleId: csipCheckViewRisksStr("tc-page2-0"), Severity: csipCheckViewRisksStr("low")},
				},
				RequestId: csipCheckViewRisksStr("fake-request-id"),
			}}, nil
		default:
			t.Errorf("unexpected extra API call %d", callCount)

			return &csipv20221121.DescribeCheckViewRisksResponse{}, nil
		}
	})

	service := svcsip.NewCsipService(client)
	describeCheckViewRisksRequest := csipv20221121.NewDescribeCheckViewRisksRequest()
	describeCheckViewRisksRequest.Order = csipCheckViewRisksStr("desc")

	checkViewRiskList, standardNameList, assetTypeList, providerList, err := service.DescribeCheckViewRisks(context.Background(), describeCheckViewRisksRequest)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if callCount != 2 {
		t.Errorf("expected 2 API calls, got %d", callCount)
	}
	if len(checkViewRiskList) != 101 {
		t.Fatalf("expected 101 risks across pages, got %d", len(checkViewRiskList))
	}
	if checkViewRiskList[0].RiskRuleId == nil || *checkViewRiskList[0].RiskRuleId != "tc-page1-0" {
		t.Errorf("unexpected first risk id: %v", checkViewRiskList[0].RiskRuleId)
	}
	if checkViewRiskList[100].RiskRuleId == nil || *checkViewRiskList[100].RiskRuleId != "tc-page2-0" {
		t.Errorf("unexpected last risk id: %v", checkViewRiskList[100].RiskRuleId)
	}
	if len(standardNameList) != 1 {
		t.Errorf("expected 1 standard name, got %d", len(standardNameList))
	}
	if len(assetTypeList) != 0 || len(providerList) != 0 {
		t.Errorf("expected empty asset type / provider lists, got %d / %d", len(assetTypeList), len(providerList))
	}
}
