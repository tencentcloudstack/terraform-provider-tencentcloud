package monitor_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	monitorv20230616 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/monitor/v20230616"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/monitor"
)

func newTestMonitorExtMetricDataSource(t *testing.T, client *connectivity.TencentCloudClient) *monitor.MonitorExtMetricDataSource {
	t.Helper()
	ds := monitor.NewMonitorExtMetricDataSource().(*monitor.MonitorExtMetricDataSource)
	ds.Configure(context.Background(), datasource.ConfigureRequest{
		ProviderData: &sharedmeta.ProviderMeta{Client: client},
	}, &datasource.ConfigureResponse{})
	return ds
}

func testMonitorExtMetricSchema(t *testing.T) schema.Schema {
	t.Helper()
	resp := &datasource.SchemaResponse{}
	ds := monitor.NewMonitorExtMetricDataSource()
	ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("failed to build data source schema: %v", resp.Diagnostics)
	}
	return resp.Schema
}

func monitorExtMetricAttrTypes() map[string]tftypes.Type {
	return map[string]tftypes.Type{
		"ext_namespace": tftypes.String,
		"ext_metric_list": tftypes.List{ElementType: tftypes.Object{
			AttributeTypes: map[string]tftypes.Type{
				"metric_name":    tftypes.String,
				"metric_c_name":  tftypes.String,
				"cn_meaning":     tftypes.String,
				"en_meaning":     tftypes.String,
				"unit":           tftypes.String,
				"dimension_flag": tftypes.Bool,
			},
		}},
	}
}

func newMonitorExtMetricReadRequest(t *testing.T, extNamespace string) datasource.ReadRequest {
	t.Helper()
	raw := tftypes.NewValue(tftypes.Object{
		AttributeTypes: monitorExtMetricAttrTypes(),
	}, map[string]tftypes.Value{
		"ext_namespace": tftypes.NewValue(tftypes.String, extNamespace),
		"ext_metric_list": tftypes.NewValue(tftypes.List{ElementType: tftypes.Object{AttributeTypes: map[string]tftypes.Type{
			"metric_name":    tftypes.String,
			"metric_c_name":  tftypes.String,
			"cn_meaning":     tftypes.String,
			"en_meaning":     tftypes.String,
			"unit":           tftypes.String,
			"dimension_flag": tftypes.Bool,
		}}}, nil),
	})

	req := datasource.ReadRequest{}
	req.Config.Raw = raw
	req.Config.Schema = testMonitorExtMetricSchema(t)
	return req
}

func runMonitorExtMetricRead(t *testing.T, ds *monitor.MonitorExtMetricDataSource, req datasource.ReadRequest) *datasource.ReadResponse {
	t.Helper()
	resp := &datasource.ReadResponse{
		State: tfsdk.State{
			Schema: testMonitorExtMetricSchema(t),
		},
	}
	ds.Read(context.Background(), req, resp)
	return resp
}

func testMonitorExtMetricResponse(metrics ...*monitorv20230616.ExtMetric) *monitorv20230616.DescribeExtMetricResponse {
	return &monitorv20230616.DescribeExtMetricResponse{
		Response: &monitorv20230616.DescribeExtMetricResponseParams{
			ExtMetricList: metrics,
		},
	}
}

func TestMonitorExtMetric_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	service := &monitor.MonitorService{}
	patches.ApplyMethodFunc(service, "DescribeExtMetric", func(ctx context.Context, request *monitorv20230616.DescribeExtMetricRequest) (*monitorv20230616.DescribeExtMetricResponse, error) {
		if request.ExtNamespace == nil || *request.ExtNamespace != "qcs::monitor: namespace-1" {
			t.Errorf("unexpected ExtNamespace: %v", request.ExtNamespace)
		}
		return testMonitorExtMetricResponse(
			&monitorv20230616.ExtMetric{
				MetricName:    strPtrMonitorExtMetric("cpu_usage"),
				MetricCName:   strPtrMonitorExtMetric("CPU使用率"),
				CNMeaning:     strPtrMonitorExtMetric("CPU使用率"),
				EnMeaning:     strPtrMonitorExtMetric("CPU usage"),
				Unit:          strPtrMonitorExtMetric("%"),
				DimensionFlag: boolPtrMonitorExtMetric(true),
			},
			&monitorv20230616.ExtMetric{
				MetricName:    strPtrMonitorExtMetric("mem_usage"),
				MetricCName:   nil,
				CNMeaning:     nil,
				EnMeaning:     nil,
				Unit:          nil,
				DimensionFlag: boolPtrMonitorExtMetric(false),
			},
		), nil
	})

	ds := newTestMonitorExtMetricDataSource(t, client)
	req := newMonitorExtMetricReadRequest(t, "qcs::monitor: namespace-1")
	resp := runMonitorExtMetricRead(t, ds, req)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error diagnostics: %v", resp.Diagnostics)
	}

	var state monitor.MonitorExtMetricDataSourceModel
	if err := resp.State.Get(context.Background(), &state); err != nil {
		t.Fatalf("failed to get state: %v", err)
	}
	if state.ExtNamespace.ValueString() != "qcs::monitor: namespace-1" {
		t.Errorf("ExtNamespace = %q, want qcs::monitor: namespace-1", state.ExtNamespace.ValueString())
	}
	if len(state.ExtMetricList.Elements()) != 2 {
		t.Fatalf("ExtMetricList length = %d, want 2", len(state.ExtMetricList.Elements()))
	}

	var items []monitor.MonitorExtMetricItemModel
	if diags := state.ExtMetricList.ElementsAs(context.Background(), &items, false); diags.HasError() {
		t.Fatalf("failed to decode ext_metric_list: %v", diags)
	}

	if items[0].MetricName.ValueString() != "cpu_usage" {
		t.Errorf("MetricName = %q, want cpu_usage", items[0].MetricName.ValueString())
	}
	if items[0].MetricCName.ValueString() != "CPU使用率" {
		t.Errorf("MetricCName = %q, want CPU使用率", items[0].MetricCName.ValueString())
	}
	if items[0].CNMeaning.ValueString() != "CPU使用率" {
		t.Errorf("CNMeaning = %q, want CPU使用率", items[0].CNMeaning.ValueString())
	}
	if items[0].EnMeaning.ValueString() != "CPU usage" {
		t.Errorf("EnMeaning = %q, want CPU usage", items[0].EnMeaning.ValueString())
	}
	if items[0].Unit.ValueString() != "%" {
		t.Errorf("Unit = %q, want %%", items[0].Unit.ValueString())
	}
	if !items[0].DimensionFlag.ValueBool() {
		t.Errorf("DimensionFlag = %v, want true", items[0].DimensionFlag.ValueBool())
	}

	if items[1].MetricName.ValueString() != "mem_usage" {
		t.Errorf("MetricName = %q, want mem_usage", items[1].MetricName.ValueString())
	}
	for name, v := range map[string]bool{
		"MetricCName":   items[1].MetricCName.IsNull(),
		"CNMeaning":     items[1].CNMeaning.IsNull(),
		"EnMeaning":     items[1].EnMeaning.IsNull(),
		"Unit":          items[1].Unit.IsNull(),
		"DimensionFlag": items[1].DimensionFlag.ValueBool() == false,
	} {
		if !v {
			t.Errorf("second metric field %s should be null/false", name)
		}
	}
}

func TestMonitorExtMetric_Read_EmptyList(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	service := &monitor.MonitorService{}
	patches.ApplyMethodFunc(service, "DescribeExtMetric", func(ctx context.Context, request *monitorv20230616.DescribeExtMetricRequest) (*monitorv20230616.DescribeExtMetricResponse, error) {
		return testMonitorExtMetricResponse(), nil
	})

	ds := newTestMonitorExtMetricDataSource(t, client)
	req := newMonitorExtMetricReadRequest(t, "namespace-empty")
	resp := runMonitorExtMetricRead(t, ds, req)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error diagnostics: %v", resp.Diagnostics)
	}

	var state monitor.MonitorExtMetricDataSourceModel
	if err := resp.State.Get(context.Background(), &state); err != nil {
		t.Fatalf("failed to get state: %v", err)
	}
	if len(state.ExtMetricList.Elements()) != 0 {
		t.Errorf("ExtMetricList length = %d, want 0", len(state.ExtMetricList.Elements()))
	}
}

func TestMonitorExtMetric_Read_ApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	service := &monitor.MonitorService{}
	patches.ApplyMethodFunc(service, "DescribeExtMetric", func(ctx context.Context, request *monitorv20230616.DescribeExtMetricRequest) (*monitorv20230616.DescribeExtMetricResponse, error) {
		return nil, fmt.Errorf("internal error")
	})

	ds := newTestMonitorExtMetricDataSource(t, client)
	req := newMonitorExtMetricReadRequest(t, "namespace-err")
	resp := runMonitorExtMetricRead(t, ds, req)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error diagnostics on API failure")
	}
}

func TestMonitorExtMetric_Read_NilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	service := &monitor.MonitorService{}
	patches.ApplyMethodFunc(service, "DescribeExtMetric", func(ctx context.Context, request *monitorv20230616.DescribeExtMetricRequest) (*monitorv20230616.DescribeExtMetricResponse, error) {
		return nil, nil
	})

	ds := newTestMonitorExtMetricDataSource(t, client)
	req := newMonitorExtMetricReadRequest(t, "namespace-nil")
	resp := runMonitorExtMetricRead(t, ds, req)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error diagnostics on nil response")
	}
}

func strPtrMonitorExtMetric(s string) *string {
	return &s
}

func boolPtrMonitorExtMetric(b bool) *bool {
	return &b
}
