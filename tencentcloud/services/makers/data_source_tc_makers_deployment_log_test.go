package makers_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	makerssdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/makers/v20261113"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/internal/sharedmeta"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/makers"
)

func newTestMakersDeploymentLogDataSource(t *testing.T, client *connectivity.TencentCloudClient) *makers.MakersDeploymentLogDataSource {
	t.Helper()
	ds := makers.NewMakersDeploymentLogDataSource().(*makers.MakersDeploymentLogDataSource)
	ds.Configure(context.Background(), datasource.ConfigureRequest{
		ProviderData: &sharedmeta.ProviderMeta{Client: client},
	}, &datasource.ConfigureResponse{})
	return ds
}

func testMakersDeploymentLogSchema(t *testing.T) schema.Schema {
	t.Helper()
	resp := &datasource.SchemaResponse{}
	ds := makers.NewMakersDeploymentLogDataSource()
	ds.Schema(context.Background(), datasource.SchemaRequest{}, resp)
	if resp.Diagnostics.HasError() {
		t.Fatalf("failed to build data source schema: %v", resp.Diagnostics)
	}
	return resp.Schema
}

func newMakersDeploymentLogReadRequest(t *testing.T, zoneId, projectId, deploymentId string) datasource.ReadRequest {
	t.Helper()
	raw := tftypes.NewValue(tftypes.Object{
		AttributeTypes: map[string]tftypes.Type{
			"zone_id":       tftypes.String,
			"project_id":    tftypes.String,
			"deployment_id": tftypes.String,
			// Computed attributes are present as null in the config object,
			// exactly like Terraform provides it at read time.
			"log_url": tftypes.String,
		},
	}, map[string]tftypes.Value{
		"zone_id":       tftypes.NewValue(tftypes.String, zoneId),
		"project_id":    tftypes.NewValue(tftypes.String, projectId),
		"deployment_id": tftypes.NewValue(tftypes.String, deploymentId),
		"log_url":       tftypes.NewValue(tftypes.String, nil),
	})

	req := datasource.ReadRequest{}
	req.Config.Raw = raw
	req.Config.Schema = testMakersDeploymentLogSchema(t)
	return req
}

func runMakersDeploymentLogRead(t *testing.T, ds *makers.MakersDeploymentLogDataSource, req datasource.ReadRequest) *datasource.ReadResponse {
	t.Helper()
	resp := &datasource.ReadResponse{
		State: tfsdk.State{
			Schema: testMakersDeploymentLogSchema(t),
		},
	}
	ds.Read(context.Background(), req, resp)
	return resp
}

func TestMakersDeploymentLog_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	service := &makers.MakersService{}
	patches.ApplyMethodFunc(service, "DescribeDeploymentLog", func(ctx context.Context, request *makerssdk.DescribeDeploymentLogRequest) (*string, error) {
		if request.ZoneId == nil || *request.ZoneId != "zone-1" {
			t.Errorf("unexpected ZoneId: %v", request.ZoneId)
		}
		if request.ProjectId == nil || *request.ProjectId != "project-1" {
			t.Errorf("unexpected ProjectId: %v", request.ProjectId)
		}
		if request.DeploymentId == nil || *request.DeploymentId != "deploy-1" {
			t.Errorf("unexpected DeploymentId: %v", request.DeploymentId)
		}
		return ptrStrMakersDeploymentLog("https://log.example.com/build-1.log"), nil
	})

	ds := newTestMakersDeploymentLogDataSource(t, client)
	req := newMakersDeploymentLogReadRequest(t, "zone-1", "project-1", "deploy-1")
	resp := runMakersDeploymentLogRead(t, ds, req)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error diagnostics: %v", resp.Diagnostics)
	}

	var state makers.MakersDeploymentLogDataSourceModel
	if err := resp.State.Get(context.Background(), &state); err != nil {
		t.Fatalf("failed to get state: %v", err)
	}
	if state.ZoneId.ValueString() != "zone-1" {
		t.Errorf("ZoneId = %q, want zone-1", state.ZoneId.ValueString())
	}
	if state.ProjectId.ValueString() != "project-1" {
		t.Errorf("ProjectId = %q, want project-1", state.ProjectId.ValueString())
	}
	if state.DeploymentId.ValueString() != "deploy-1" {
		t.Errorf("DeploymentId = %q, want deploy-1", state.DeploymentId.ValueString())
	}
	if state.LogUrl.ValueString() != "https://log.example.com/build-1.log" {
		t.Errorf("LogUrl = %q, want https://log.example.com/build-1.log", state.LogUrl.ValueString())
	}
}

func TestMakersDeploymentLog_Read_NullLogUrl(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	service := &makers.MakersService{}
	patches.ApplyMethodFunc(service, "DescribeDeploymentLog", func(ctx context.Context, request *makerssdk.DescribeDeploymentLogRequest) (*string, error) {
		return nil, nil
	})

	ds := newTestMakersDeploymentLogDataSource(t, client)
	req := newMakersDeploymentLogReadRequest(t, "zone-2", "project-2", "deploy-2")
	resp := runMakersDeploymentLogRead(t, ds, req)

	if resp.Diagnostics.HasError() {
		t.Fatalf("unexpected error diagnostics: %v", resp.Diagnostics)
	}

	var state makers.MakersDeploymentLogDataSourceModel
	if err := resp.State.Get(context.Background(), &state); err != nil {
		t.Fatalf("failed to get state: %v", err)
	}
	if !state.LogUrl.IsNull() {
		t.Errorf("LogUrl should be null, got %q", state.LogUrl.ValueString())
	}
}

func TestMakersDeploymentLog_Read_ApiError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	service := &makers.MakersService{}
	patches.ApplyMethodFunc(service, "DescribeDeploymentLog", func(ctx context.Context, request *makerssdk.DescribeDeploymentLogRequest) (*string, error) {
		return nil, fmt.Errorf("internal error")
	})

	ds := newTestMakersDeploymentLogDataSource(t, client)
	req := newMakersDeploymentLogReadRequest(t, "zone-3", "project-3", "deploy-3")
	resp := runMakersDeploymentLogRead(t, ds, req)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error diagnostics on API failure")
	}
}

func TestMakersDeploymentLog_Read_EmptyResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	client := &connectivity.TencentCloudClient{}
	service := &makers.MakersService{}
	patches.ApplyMethodFunc(service, "DescribeDeploymentLog", func(ctx context.Context, request *makerssdk.DescribeDeploymentLogRequest) (*string, error) {
		return nil, fmt.Errorf("DescribeDeploymentLog returned an empty response")
	})

	ds := newTestMakersDeploymentLogDataSource(t, client)
	req := newMakersDeploymentLogReadRequest(t, "zone-4", "project-4", "deploy-4")
	resp := runMakersDeploymentLogRead(t, ds, req)

	if !resp.Diagnostics.HasError() {
		t.Fatal("expected error diagnostics on empty response")
	}
}

func ptrStrMakersDeploymentLog(s string) *string {
	return &s
}
