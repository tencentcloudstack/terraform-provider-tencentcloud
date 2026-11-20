package monitor_test

import (
	"context"
	"errors"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/stretchr/testify/assert"

	monitorv20230616 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20230616"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	svmonitor "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/monitor"
)

func ptrStrMonitor(s string) *string { return &s }
func ptrInt64Monitor(v int64) *int64 { return &v }

// monitorDispenseRegionSchema returns the framework data source schema used
// to build the raw config/state value.
func monitorDispenseRegionSchema(t *testing.T) schema.Schema {
	d := svmonitor.NewMonitorDispenseRegionDataSource()
	resp := &datasource.SchemaResponse{}
	d.Schema(context.Background(), datasource.SchemaRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("failed to build data source schema: %v", resp.Diagnostics)
	}
	return resp.Schema
}

// monitorDispenseRegionItemType builds the tftypes.Object type for a
// region_list element.
func monitorDispenseRegionItemType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"region":         tftypes.String,
			"region_cn_name": tftypes.String,
			"region_en_name": tftypes.String,
			"rule_number":    tftypes.Number,
		},
	}
}

// monitorDispenseRegionRootType builds the tftypes.Object type for the data
// source root, mirroring every schema attribute.
func monitorDispenseRegionRootType() tftypes.Object {
	return tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"region_list": tftypes.List{ElementType: monitorDispenseRegionItemType()},
		},
	}
}

// buildMonitorDispenseRegionRawValue builds the raw config value. The data
// source has no configurable inputs; all attributes are computed outputs
// initialized as null.
func buildMonitorDispenseRegionRawValue() tftypes.Value {
	attrs := map[string]tftypes.Value{
		"region_list": tftypes.NewValue(tftypes.List{ElementType: monitorDispenseRegionItemType()}, nil),
	}
	return tftypes.NewValue(monitorDispenseRegionRootType(), attrs)
}

// setMonitorDispenseRegionClient injects a mock connectivity client into the
// data source instance via the Configure method.
func setMonitorDispenseRegionClient(d *svmonitor.MonitorDispenseRegionDataSource, client *connectivity.TencentCloudClient) {
	meta := &sharedmeta.ProviderMeta{Client: client}
	configureReq := datasource.ConfigureRequest{ProviderData: meta}
	configureResp := &datasource.ConfigureResponse{}
	d.Configure(context.Background(), configureReq, configureResp)
}

// newMonitorDispenseRegionReadRequest builds a ReadRequest with the given
// config value and schema.
func newMonitorDispenseRegionReadRequest(sch schema.Schema, raw tftypes.Value) datasource.ReadRequest {
	req := datasource.ReadRequest{}
	req.Config.Raw = raw
	req.Config.Schema = sch
	return req
}

// newMonitorDispenseRegionReadResponse builds a ReadResponse whose State
// carries the schema, so resp.State.Set can reflect the model into a raw
// value.
func newMonitorDispenseRegionReadResponse(sch schema.Schema) *datasource.ReadResponse {
	resp := &datasource.ReadResponse{}
	resp.State.Schema = sch
	return resp
}

// go test ./tencentcloud/services/monitor/ -run "TestMonitorDispenseRegion" -v -count=1 -gcflags="all=-l"

// TestMonitorDispenseRegion_MetadataAndSchema validates metadata and schema
// definition.
func TestMonitorDispenseRegion_MetadataAndSchema(t *testing.T) {
	d := svmonitor.NewMonitorDispenseRegionDataSource()

	metaResp := &datasource.MetadataResponse{}
	d.Metadata(context.Background(), datasource.MetadataRequest{}, metaResp)
	assert.Equal(t, "tencentcloud_monitor_dispense_region", metaResp.TypeName)

	sch := monitorDispenseRegionSchema(t)
	attrs := sch.Attributes

	regionList, ok := attrs["region_list"].(schema.ListNestedAttribute)
	assert.True(t, ok, "attribute region_list should be ListNestedAttribute")
	assert.True(t, regionList.Computed)
	assert.False(t, regionList.Optional)
	assert.False(t, regionList.Required)

	nested := regionList.NestedObject.Attributes
	assert.True(t, nested["region"].(schema.StringAttribute).Computed)
	assert.True(t, nested["region_cn_name"].(schema.StringAttribute).Computed)
	assert.True(t, nested["region_en_name"].(schema.StringAttribute).Computed)
	assert.True(t, nested["rule_number"].(schema.Int64Attribute).Computed)

	assert.Len(t, attrs, 1)
	assert.Len(t, nested, 4)
}

// TestMonitorDispenseRegion_ReadBasic tests a successful read with two
// region records and asserts every leaf field is read back.
func TestMonitorDispenseRegion_ReadBasic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	service := &svmonitor.MonitorService{}
	patches.ApplyMethodFunc(service, "DescribeDispenseRegions", func(_ context.Context) ([]*monitorv20230616.DispenseRegion, error) {
		return []*monitorv20230616.DispenseRegion{
			{
				Region:       ptrStrMonitor("ap-guangzhou"),
				RegionCnName: ptrStrMonitor("华南地区(广州)"),
				RegionEnName: ptrStrMonitor("South China (Guangzhou)"),
				RuleNumber:   ptrInt64Monitor(3),
			},
			{
				Region:       ptrStrMonitor("ap-shanghai"),
				RegionCnName: ptrStrMonitor("华东地区(上海)"),
				RegionEnName: ptrStrMonitor("East China (Shanghai)"),
				RuleNumber:   ptrInt64Monitor(0),
			},
		}, nil
	})

	ds := svmonitor.NewMonitorDispenseRegionDataSource().(*svmonitor.MonitorDispenseRegionDataSource)
	setMonitorDispenseRegionClient(ds, client)

	sch := monitorDispenseRegionSchema(t)
	raw := buildMonitorDispenseRegionRawValue()
	req := newMonitorDispenseRegionReadRequest(sch, raw)
	resp := newMonitorDispenseRegionReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	var result svmonitor.MonitorDispenseRegionDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &result)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.False(t, result.RegionList.IsNull())

	var items []svmonitor.MonitorDispenseRegionItemModel
	resp.Diagnostics.Append(result.RegionList.ElementsAs(context.Background(), &items, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, items, 2)

	// assert every leaf field of both elements
	assert.Equal(t, "ap-guangzhou", items[0].Region.ValueString())
	assert.Equal(t, "华南地区(广州)", items[0].RegionCnName.ValueString())
	assert.Equal(t, "South China (Guangzhou)", items[0].RegionEnName.ValueString())
	assert.Equal(t, int64(3), items[0].RuleNumber.ValueInt64())
	assert.Equal(t, "ap-shanghai", items[1].Region.ValueString())
	assert.Equal(t, "华东地区(上海)", items[1].RegionCnName.ValueString())
	assert.Equal(t, "East China (Shanghai)", items[1].RegionEnName.ValueString())
	assert.Equal(t, int64(0), items[1].RuleNumber.ValueInt64())
}

// TestMonitorDispenseRegion_ReadEmpty covers an empty region list. An empty
// result is a legitimate answer, not a failure: the data source must succeed
// with an empty region_list.
func TestMonitorDispenseRegion_ReadEmpty(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	service := &svmonitor.MonitorService{}
	patches.ApplyMethodFunc(service, "DescribeDispenseRegions", func(_ context.Context) ([]*monitorv20230616.DispenseRegion, error) {
		return []*monitorv20230616.DispenseRegion{}, nil
	})

	ds := svmonitor.NewMonitorDispenseRegionDataSource().(*svmonitor.MonitorDispenseRegionDataSource)
	setMonitorDispenseRegionClient(ds, client)

	sch := monitorDispenseRegionSchema(t)
	raw := buildMonitorDispenseRegionRawValue()
	req := newMonitorDispenseRegionReadRequest(sch, raw)
	resp := newMonitorDispenseRegionReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.False(t, resp.Diagnostics.HasError(), "expected no error diagnostics, got: %v", resp.Diagnostics)

	var result svmonitor.MonitorDispenseRegionDataSourceModel
	resp.Diagnostics.Append(resp.State.Get(context.Background(), &result)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.False(t, result.RegionList.IsNull())

	var items []svmonitor.MonitorDispenseRegionItemModel
	resp.Diagnostics.Append(result.RegionList.ElementsAs(context.Background(), &items, false)...)
	assert.False(t, resp.Diagnostics.HasError())
	assert.Len(t, items, 0)
}

// TestMonitorDispenseRegion_ReadApiError covers an API-level failure: the
// data source must surface the error through diagnostics instead of silently
// persisting empty state.
func TestMonitorDispenseRegion_ReadApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	service := &svmonitor.MonitorService{}
	patches.ApplyMethodFunc(service, "DescribeDispenseRegions", func(_ context.Context) ([]*monitorv20230616.DispenseRegion, error) {
		return nil, errors.New("DescribeDispenseRegion failed")
	})

	ds := svmonitor.NewMonitorDispenseRegionDataSource().(*svmonitor.MonitorDispenseRegionDataSource)
	setMonitorDispenseRegionClient(ds, client)

	sch := monitorDispenseRegionSchema(t)
	raw := buildMonitorDispenseRegionRawValue()
	req := newMonitorDispenseRegionReadRequest(sch, raw)
	resp := newMonitorDispenseRegionReadResponse(sch)

	ds.Read(context.Background(), req, resp)

	assert.True(t, resp.Diagnostics.HasError(), "expected error diagnostics for API failure")
}
