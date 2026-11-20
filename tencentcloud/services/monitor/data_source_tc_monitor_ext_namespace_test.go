package monitor_test

import (
	"context"
	"errors"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/monitor"
)

func newTestMonitorExtNamespaceDataSource(t *testing.T, client *connectivity.TencentCloudClient) *monitor.MonitorExtNamespaceDataSource {
	t.Helper()
	ds := monitor.NewMonitorExtNamespaceDataSource().(*monitor.MonitorExtNamespaceDataSource)
	ds.Configure(context.Background(), datasource.ConfigureRequest{
		ProviderData: &sharedmeta.ProviderMeta{Client: client},
	}, &datasource.ConfigureResponse{})
	return ds
}

func testMonitorExtNamespaceSchema(t *testing.T) schema.Schema {
	t.Helper()
	resp := &datasource.SchemaResponse{}
	ds := monitor.NewMonitorExtNamespaceDataSource()
	ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("failed to build data source schema: %v", resp.Diagnostics)
	}
	return resp.Schema
}

func newMonitorExtNamespaceReadRequest(t *testing.T) datasource.ReadRequest {
	t.Helper()
	raw := tftypes.NewValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"ext_namespace_list": tftypes.Set{ElementType: tftypes.String},
		},
	}, map[string]tftypes.Value{
		// Computed attributes are present as null in the config object,
		// exactly like Terraform provides it at read time.
		"ext_namespace_list": tftypes.NewValue(tftypes.Set{ElementType: tftypes.String}, nil),
	})

	req := datasource.ReadRequest{}
	req.Config.Raw = raw
	req.Config.Schema = testMonitorExtNamespaceSchema(t)
	return req
}

func runMonitorExtNamespaceRead(t *testing.T, ds *monitor.MonitorExtNamespaceDataSource, req datasource.ReadRequest) *datasource.ReadResponse {
	t.Helper()
	resp := &datasource.ReadResponse{
		State: tfsdk.State{
			Schema: testMonitorExtNamespaceSchema(t),
		},
	}
	ds.Read(context.Background(), req, resp)
	return resp
}

func ptrStrMonitorExtNamespace(s string) *string {
	return &s
}

func TestMonitorExtNamespace_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	service := &monitor.MonitorService{}
	patches.ApplyMethodFunc(service, "DescribeExtNamespace", func(ctx context.Context) ([]*string, error) {
		return []*string{
			ptrStrMonitorExtNamespace("ns-one"),
			ptrStrMonitorExtNamespace("ns-two"),
		}, nil
	})

	ds := newTestMonitorExtNamespaceDataSource(t, client)
	req := newMonitorExtNamespaceReadRequest(t)
	resp := runMonitorExtNamespaceRead(t, ds, req)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error diagnostics: %v", resp.Diagnostics)
	}

	var state monitor.MonitorExtNamespaceDataSourceModel
	if err := resp.State.Get(context.Background(), &state); err != nil {
		t.Fatalf("failed to get state: %v", err)
	}
	if state.ExtNamespaceList.IsNull() || state.ExtNamespaceList.IsUnknown() {
		t.Fatal("ext_namespace_list should not be null or unknown")
	}
	length := len(state.ExtNamespaceList.Elements())
	if length != 2 {
		t.Fatalf("ext_namespace_list length = %d, want 2", length)
	}
	values := make([]string, 0, length)
	for _, element := range state.ExtNamespaceList.Elements() {
		s, ok := element.(interface{ ValueString() string })
		if !ok {
			t.Fatalf("unexpected element type: %T", element)
		}
		values = append(values, s.ValueString())
	}
	if values[0] != "ns-one" || values[1] != "ns-two" {
		t.Errorf("ext_namespace_list = %v, want [ns-one ns-two]", values)
	}
}

func TestMonitorExtNamespace_Read_EmptyList(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	service := &monitor.MonitorService{}
	patches.ApplyMethodFunc(service, "DescribeExtNamespace", func(ctx context.Context) ([]*string, error) {
		return []*string{}, nil
	})

	ds := newTestMonitorExtNamespaceDataSource(t, client)
	req := newMonitorExtNamespaceReadRequest(t)
	resp := runMonitorExtNamespaceRead(t, ds, req)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error diagnostics: %v", resp.Diagnostics)
	}

	var state monitor.MonitorExtNamespaceDataSourceModel
	if err := resp.State.Get(context.Background(), &state); err != nil {
		t.Fatalf("failed to get state: %v", err)
	}
	if state.ExtNamespaceList.IsNull() || state.ExtNamespaceList.IsUnknown() {
		t.Fatal("ext_namespace_list should be an empty set, not null or unknown")
	}
	if len(state.ExtNamespaceList.Elements()) != 0 {
		t.Errorf("ext_namespace_list length = %d, want 0", len(state.ExtNamespaceList.Elements()))
	}
}

func TestMonitorExtNamespace_Read_ApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	service := &monitor.MonitorService{}
	patches.ApplyMethodFunc(service, "DescribeExtNamespace", func(ctx context.Context) ([]*string, error) {
		return nil, errors.New("internal error")
	})

	ds := newTestMonitorExtNamespaceDataSource(t, client)
	req := newMonitorExtNamespaceReadRequest(t)
	resp := runMonitorExtNamespaceRead(t, ds, req)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error diagnostics on API failure")
	}
}
