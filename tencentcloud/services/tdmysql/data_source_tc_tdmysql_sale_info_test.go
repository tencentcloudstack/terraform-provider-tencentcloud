package tdmysql_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"

	tdmysqlv20211122 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tdmysql/v20211122"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	svtdmysql "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/tdmysql"
)

func tdmysqlSaleInfoStr(s string) *string { return &s }
func tdmysqlSaleInfoInt64(v int64) *int64 { return &v }
func tdmysqlSaleInfoBool(v bool) *bool    { return &v }

// tdmysqlSaleInfoSchema returns the framework data source schema used to
// build the raw config/state value.
func tdmysqlSaleInfoSchema(t *testing.T) schema.Schema {
	d := svtdmysql.NewTdmysqlSaleInfoDataSource()
	resp := &datasource.SchemaResponse{}
	d.Schema(context.Background(), datasource.SchemaRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("failed to build data source schema: %v", resp.Diagnostics)
	}

	return resp.Schema
}

// tdmysqlSaleInfoRootType builds the tftypes.Object type for the data source
// root, mirroring every schema attribute.
func tdmysqlSaleInfoRootType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"id":            tftypes.String,
			"src_region":    tftypes.String,
			"instance_id":   tftypes.String,
			"instance_type": tftypes.String,
			"region_list":   tftypes.List{ElementType: tdmysqlSaleInfoRegionItemType()},
		},
	}
}

func tdmysqlSaleInfoRegionItemType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"region":                   tftypes.String,
			"zone_list":                tftypes.List{ElementType: tdmysqlSaleInfoZoneItemType()},
			"region_name":              tftypes.String,
			"enable":                   tftypes.Number,
			"available_zone_num":       tftypes.Number,
			"is_supported_log_replica": tftypes.Bool,
			"zone_group":               tftypes.List{ElementType: tdmysqlSaleInfoZoneGroupItemType()},
			"is_support_serverless":    tftypes.Bool,
		},
	}
}

func tdmysqlSaleInfoZoneItemType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"zone":                  tftypes.String,
			"zone_name":             tftypes.String,
			"enable":                tftypes.Number,
			"is_default_master":     tftypes.Number,
			"available_disk_types":  tftypes.List{ElementType: tftypes.String},
			"support_types":         tftypes.List{ElementType: tftypes.String},
			"is_support_serverless": tftypes.Bool,
		},
	}
}

func tdmysqlSaleInfoZoneGroupItemType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"zone_num":              tftypes.Number,
			"zones":                 tftypes.List{ElementType: tftypes.String},
			"support_types":         tftypes.List{ElementType: tftypes.String},
			"available_disk_types":  tftypes.List{ElementType: tftypes.String},
			"is_support_serverless": tftypes.Bool,
		},
	}
}

// buildTdmysqlSaleInfoRawValue builds the raw config value with the given
// query inputs; every other attribute is unset (nil).
func buildTdmysqlSaleInfoRawValue(srcRegion, instanceId, instanceType string) tftypes.Value {
	attrs := map[string]tftypes.Value{
		"id":            tftypes.NewValue(tftypes.String, nil),
		"src_region":    tftypes.NewValue(tftypes.String, nilIfEmpty(srcRegion)),
		"instance_id":   tftypes.NewValue(tftypes.String, nilIfEmpty(instanceId)),
		"instance_type": tftypes.NewValue(tftypes.String, nilIfEmpty(instanceType)),
		"region_list":   tftypes.NewValue(tftypes.List{ElementType: tdmysqlSaleInfoRegionItemType()}, nil),
	}

	return tftypes.NewValue(tdmysqlSaleInfoRootType(), attrs)
}

func nilIfEmpty(s string) interface{} {
	if s == "" {
		return nil
	}

	return s
}

// setTdmysqlSaleInfoClient injects a mock connectivity client into the data
// source instance via the Configure method.
func setTdmysqlSaleInfoClient(d *svtdmysql.TdmysqlSaleInfoDataSource, client *connectivity.TencentCloudClient) {
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := datasource.ConfigureRequest{ProviderData: meta}
	configureResp := &datasource.ConfigureResponse{}
	d.Configure(context.Background(), configureReq, configureResp)
}

func newTdmysqlSaleInfoReadRequest(sch schema.Schema, raw tftypes.Value) datasource.ReadRequest {
	req := datasource.ReadRequest{}
	req.Config.Raw = raw
	req.Config.Schema = sch

	return req
}

func newTdmysqlSaleInfoReadResponse(sch schema.Schema) *datasource.ReadResponse {
	resp := &datasource.ReadResponse{}
	resp.State.Schema = sch

	return resp
}

// TestTdmysqlSaleInfo_MetadataAndSchema validates metadata and schema
// definition.
func TestTdmysqlSaleInfo_MetadataAndSchema(t *testing.T) {
	d := svtdmysql.NewTdmysqlSaleInfoDataSource()

	metaResp := &datasource.MetadataResponse{}
	d.Metadata(context.Background(), datasource.MetadataRequest{}, metaResp)
	assert.Equal(t, "tencentcloud_tdmysql_sale_info", metaResp.TypeName)

	sch := tdmysqlSaleInfoSchema(t)
	attrs := sch.Attributes

	assert.True(t, attrs["src_region"].(schema.StringAttribute).Optional)
	assert.True(t, attrs["instance_id"].(schema.StringAttribute).Optional)
	assert.True(t, attrs["instance_type"].(schema.StringAttribute).Optional)
	assert.True(t, attrs["region_list"].(schema.ListNestedAttribute).Computed)
	assert.True(t, attrs["id"].(schema.StringAttribute).Computed)

	// Paging knobs and total-count aggregates are NEVER exposed.
	assert.NotContains(t, attrs, "limit")
	assert.NotContains(t, attrs, "offset")
	assert.NotContains(t, attrs, "total_count")
}

// TestTdmysqlSaleInfo_ReadBasic tests a successful read with one region and
// asserts every leaf field is read back.
func TestTdmysqlSaleInfo_ReadBasic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	tdmysqlClient := &tdmysqlv20211122.Client{}
	patches.ApplyMethodReturn(client, "UseTdmysqlV20211122Client", tdmysqlClient)
	patches.ApplyMethodFunc(tdmysqlClient, "DescribeSaleInfo", func(request *tdmysqlv20211122.DescribeSaleInfoRequest) (*tdmysqlv20211122.DescribeSaleInfoResponse, error) {
		if request.InstanceType == nil || *request.InstanceType != "serverless" {
			t.Errorf("expected instance_type to be passed to the API request")
		}

		return &tdmysqlv20211122.DescribeSaleInfoResponse{Response: &tdmysqlv20211122.DescribeSaleInfoResponseParams{
			RegionList: []*tdmysqlv20211122.DescribeSaleRegionInfo{
				{
					Region:                tdmysqlSaleInfoStr("ap-guangzhou"),
					RegionName:            tdmysqlSaleInfoStr("华南地区(广州)"),
					Enable:                tdmysqlSaleInfoInt64(1),
					AvailableZoneNum:      tdmysqlSaleInfoInt64(3),
					IsSupportedLogReplica: tdmysqlSaleInfoBool(true),
					IsSupportServerless:   tdmysqlSaleInfoBool(true),
					ZoneList: []*tdmysqlv20211122.DescribeSaleZonesInfo{
						{
							Zone:                tdmysqlSaleInfoStr("ap-guangzhou-3"),
							ZoneName:            tdmysqlSaleInfoStr("广州三区"),
							Enable:              tdmysqlSaleInfoInt64(1),
							IsDefaultMaster:     tdmysqlSaleInfoInt64(1),
							AvailableDiskTypes:  []*string{tdmysqlSaleInfoStr("CLOUD_HSSD")},
							SupportTypes:        []*string{tdmysqlSaleInfoStr("exclusive")},
							IsSupportServerless: tdmysqlSaleInfoBool(true),
						},
					},
					ZoneGroup: []*tdmysqlv20211122.DescribeSaleZonesGroup{
						{
							ZoneNum:             tdmysqlSaleInfoInt64(2),
							Zones:               []*string{tdmysqlSaleInfoStr("ap-guangzhou-3"), tdmysqlSaleInfoStr("ap-guangzhou-4")},
							SupportTypes:        []*string{tdmysqlSaleInfoStr("exclusive")},
							AvailableDiskTypes:  []*string{tdmysqlSaleInfoStr("CLOUD_TCS")},
							IsSupportServerless: tdmysqlSaleInfoBool(false),
						},
					},
				},
			},
			RequestId: tdmysqlSaleInfoStr("fake-request-id"),
		}}, nil
	})

	ds := svtdmysql.NewTdmysqlSaleInfoDataSource().(*svtdmysql.TdmysqlSaleInfoDataSource)
	setTdmysqlSaleInfoClient(ds, client)

	sch := tdmysqlSaleInfoSchema(t)
	raw := buildTdmysqlSaleInfoRawValue("", "", "serverless")
	req := newTdmysqlSaleInfoReadRequest(sch, raw)
	resp := newTdmysqlSaleInfoReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	var result svtdmysql.TdmysqlSaleInfoDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &result)...)
	assert.False(t, resp.Diagnostics.HasError())

	// The state id is a random token.
	assert.NotEmpty(t, result.Id.ValueString())

	assert.False(t, result.RegionList.IsNull())

	var regions []svtdmysql.TdmysqlSaleRegionItemModel
	resp.Diagnostics.Append(result.RegionList.ElementsAs(context.Background(), &regions, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, regions, 1)

	assert.Equal(t, "ap-guangzhou", regions[0].Region.ValueString())
	assert.Equal(t, "华南地区(广州)", regions[0].RegionName.ValueString())
	assert.Equal(t, int64(1), regions[0].Enable.ValueInt64())
	assert.Equal(t, int64(3), regions[0].AvailableZoneNum.ValueInt64())
	assert.True(t, regions[0].IsSupportedLogReplica.ValueBool())
	assert.True(t, regions[0].IsSupportServerless.ValueBool())

	assert.False(t, regions[0].ZoneList.IsNull())
	var zones []svtdmysql.TdmysqlSaleZoneItemModel
	resp.Diagnostics.Append(regions[0].ZoneList.ElementsAs(context.Background(), &zones, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, zones, 1)
	assert.Equal(t, "ap-guangzhou-3", zones[0].Zone.ValueString())
	assert.Equal(t, "广州三区", zones[0].ZoneName.ValueString())
	assert.Equal(t, int64(1), zones[0].Enable.ValueInt64())
	assert.Equal(t, int64(1), zones[0].IsDefaultMaster.ValueInt64())
	assert.True(t, zones[0].IsSupportServerless.ValueBool())
	assert.False(t, zones[0].AvailableDiskTypes.IsNull())
	var diskTypes []string
	resp.Diagnostics.Append(zones[0].AvailableDiskTypes.ElementsAs(context.Background(), &diskTypes, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Equal(t, []string{"CLOUD_HSSD"}, diskTypes)

	assert.False(t, regions[0].ZoneGroup.IsNull())
	var groups []svtdmysql.TdmysqlSaleZoneGroupItemModel
	resp.Diagnostics.Append(regions[0].ZoneGroup.ElementsAs(context.Background(), &groups, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, groups, 1)
	assert.Equal(t, int64(2), groups[0].ZoneNum.ValueInt64())
	assert.False(t, groups[0].IsSupportServerless.ValueBool())
	assert.False(t, groups[0].Zones.IsNull())
	var groupZones []string
	resp.Diagnostics.Append(groups[0].Zones.ElementsAs(context.Background(), &groupZones, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Equal(t, []string{"ap-guangzhou-3", "ap-guangzhou-4"}, groupZones)
}

// TestTdmysqlSaleInfo_ReadEmpty covers an empty result: a legitimate answer,
// not a failure — the data source must succeed with an empty list.
func TestTdmysqlSaleInfo_ReadEmpty(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	tdmysqlClient := &tdmysqlv20211122.Client{}
	patches.ApplyMethodReturn(client, "UseTdmysqlV20211122Client", tdmysqlClient)
	patches.ApplyMethodFunc(tdmysqlClient, "DescribeSaleInfo", func(request *tdmysqlv20211122.DescribeSaleInfoRequest) (*tdmysqlv20211122.DescribeSaleInfoResponse, error) {
		return &tdmysqlv20211122.DescribeSaleInfoResponse{Response: &tdmysqlv20211122.DescribeSaleInfoResponseParams{}}, nil
	})

	ds := svtdmysql.NewTdmysqlSaleInfoDataSource().(*svtdmysql.TdmysqlSaleInfoDataSource)
	setTdmysqlSaleInfoClient(ds, client)

	sch := tdmysqlSaleInfoSchema(t)
	raw := buildTdmysqlSaleInfoRawValue("", "", "serverless")
	req := newTdmysqlSaleInfoReadRequest(sch, raw)
	resp := newTdmysqlSaleInfoReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	var result svtdmysql.TdmysqlSaleInfoDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &result)...)
	assert.False(t, resp.Diagnostics.HasError())

	assert.False(t, result.RegionList.IsNull())

	var regions []svtdmysql.TdmysqlSaleRegionItemModel
	resp.Diagnostics.Append(result.RegionList.ElementsAs(context.Background(), &regions, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, regions, 0)
}

// TestTdmysqlSaleInfo_ReadApiError covers an API-level failure: the error
// must surface through diagnostics instead of silently persisting empty
// state.
func TestTdmysqlSaleInfo_ReadApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	tdmysqlClient := &tdmysqlv20211122.Client{}
	patches.ApplyMethodReturn(client, "UseTdmysqlV20211122Client", tdmysqlClient)
	patches.ApplyMethodFunc(tdmysqlClient, "DescribeSaleInfo", func(request *tdmysqlv20211122.DescribeSaleInfoRequest) (*tdmysqlv20211122.DescribeSaleInfoResponse, error) {
		return nil, fmt.Errorf("mocked API error")
	})

	ds := svtdmysql.NewTdmysqlSaleInfoDataSource().(*svtdmysql.TdmysqlSaleInfoDataSource)
	setTdmysqlSaleInfoClient(ds, client)

	sch := tdmysqlSaleInfoSchema(t)
	raw := buildTdmysqlSaleInfoRawValue("", "", "serverless")
	req := newTdmysqlSaleInfoReadRequest(sch, raw)
	resp := newTdmysqlSaleInfoReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError(), "expected error diagnostics for API failure")
}
