package ags

import (
	"context"
	"errors"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-framework/attr"
	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/types"

	ags "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920"
	sdkerr "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
)

func TestAgsSandboxToolResourceInterfaces(t *testing.T) {
	r := &AgsSandboxToolResource{}
	var _ resource.Resource = r
	var _ resource.ResourceWithConfigure = r
	var _ resource.ResourceWithImportState = r
}

func setupMockClient() (*AgsSandboxToolResource, *gomonkey.Patches) {
	r := &AgsSandboxToolResource{}
	patches := gomonkey.NewPatches()
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyFunc((*AgsSandboxToolResource).Client, func(_ *AgsSandboxToolResource) *connectivity.TencentCloudClient {
		return mockClient
	})
	patches.ApplyFunc((*connectivity.TencentCloudClient).UseAgsClient,
		func(_ *connectivity.TencentCloudClient) *ags.Client {
			return &ags.Client{}
		})
	return r, patches
}

func TestAgsSandboxToolRead_NilResponse(t *testing.T) {
	r, patches := setupMockClient()
	defer patches.Reset()

	patches.ApplyFunc((*ags.Client).DescribeSandboxToolListWithContext,
		func(_ *ags.Client, _ context.Context, _ *ags.DescribeSandboxToolListRequest) (*ags.DescribeSandboxToolListResponse, error) {
			return nil, nil
		})

	tool, err := r.read(context.Background(), "tool-123")
	if err == nil {
		t.Fatal("expected error for nil response")
	}
	if tool != nil {
		t.Fatal("expected nil tool for nil response")
	}
}

func TestAgsSandboxToolRead_EmptyResponse(t *testing.T) {
	r, patches := setupMockClient()
	defer patches.Reset()

	patches.ApplyFunc((*ags.Client).DescribeSandboxToolListWithContext,
		func(_ *ags.Client, _ context.Context, _ *ags.DescribeSandboxToolListRequest) (*ags.DescribeSandboxToolListResponse, error) {
			return &ags.DescribeSandboxToolListResponse{
				Response: nil,
			}, nil
		})

	tool, err := r.read(context.Background(), "tool-123")
	if err == nil {
		t.Fatal("expected error for nil Response")
	}
	if tool != nil {
		t.Fatal("expected nil tool for nil Response")
	}
}

func TestAgsSandboxToolRead_NotFound(t *testing.T) {
	r, patches := setupMockClient()
	defer patches.Reset()

	patches.ApplyFunc((*ags.Client).DescribeSandboxToolListWithContext,
		func(_ *ags.Client, _ context.Context, _ *ags.DescribeSandboxToolListRequest) (*ags.DescribeSandboxToolListResponse, error) {
			return &ags.DescribeSandboxToolListResponse{
				Response: &ags.DescribeSandboxToolListResponseParams{
					SandboxToolSet: []*ags.SandboxTool{},
				},
			}, nil
		})

	tool, err := r.read(context.Background(), "tool-123")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tool != nil {
		t.Fatal("expected nil tool for empty result set")
	}
}

func TestAgsSandboxToolRead_Success(t *testing.T) {
	r, patches := setupMockClient()
	defer patches.Reset()

	toolId := "tool-abc123"
	toolName := "test-tool"
	toolType := "custom"
	status := "ACTIVE"
	desc := "test description"
	defaultTimeout := uint64(300)
	createTime := "2024-01-01T00:00:00Z"
	updateTime := "2024-01-02T00:00:00Z"

	patches.ApplyFunc((*ags.Client).DescribeSandboxToolListWithContext,
		func(_ *ags.Client, _ context.Context, _ *ags.DescribeSandboxToolListRequest) (*ags.DescribeSandboxToolListResponse, error) {
			return &ags.DescribeSandboxToolListResponse{
				Response: &ags.DescribeSandboxToolListResponseParams{
					SandboxToolSet: []*ags.SandboxTool{
						{
							ToolId:                &toolId,
							ToolName:              &toolName,
							ToolType:              &toolType,
							Status:                &status,
							Description:           &desc,
							DefaultTimeoutSeconds: &defaultTimeout,
							CreateTime:            &createTime,
							UpdateTime:            &updateTime,
						},
					},
				},
			}, nil
		})

	tool, err := r.read(context.Background(), toolId)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if tool == nil {
		t.Fatal("expected non-nil tool")
	}
	if *tool.ToolId != toolId {
		t.Fatalf("expected toolId %q, got %q", toolId, *tool.ToolId)
	}
	if *tool.ToolName != toolName {
		t.Fatalf("expected toolName %q, got %q", toolName, *tool.ToolName)
	}
}

func TestAgsSandboxToolRead_Error(t *testing.T) {
	r, patches := setupMockClient()
	defer patches.Reset()

	patches.ApplyFunc((*ags.Client).DescribeSandboxToolListWithContext,
		func(_ *ags.Client, _ context.Context, _ *ags.DescribeSandboxToolListRequest) (*ags.DescribeSandboxToolListResponse, error) {
			return nil, errors.New("connection refused")
		})

	tool, err := r.read(context.Background(), "tool-123")
	if err == nil {
		t.Fatal("expected error")
	}
	if tool != nil {
		t.Fatal("expected nil tool on error")
	}
}

func TestIsAgsResourceNotFound(t *testing.T) {
	if isAgsResourceNotFound(nil) {
		t.Fatal("nil error should not be notFound")
	}
	if isAgsResourceNotFound(errors.New("generic error")) {
		t.Fatal("generic error should not be notFound")
	}
	sdkNotFound := &sdkerr.TencentCloudSDKError{
		Code:    "ResourceNotFound",
		Message: "tool not found",
	}
	if !isAgsResourceNotFound(sdkNotFound) {
		t.Fatal("ResourceNotFound should be notFound")
	}
	sdkOther := &sdkerr.TencentCloudSDKError{
		Code:    "InvalidParameter",
		Message: "bad param",
	}
	if isAgsResourceNotFound(sdkOther) {
		t.Fatal("InvalidParameter should not be notFound")
	}
}

func TestBuildNetworkConfigSDK(t *testing.T) {
	nc := &agsNetworkConfigModel{
		NetworkMode: types.StringValue("PUBLIC"),
	}
	sdk := buildNetworkConfigSDK(nc)
	if sdk == nil {
		t.Fatal("expected non-nil NetworkConfiguration")
	}
	if sdk.NetworkMode == nil || *sdk.NetworkMode != "PUBLIC" {
		t.Fatalf("expected NetworkMode PUBLIC, got %v", sdk.NetworkMode)
	}
	if sdk.VpcConfig != nil {
		t.Fatal("expected nil VpcConfig when not configured")
	}
}

func TestBuildNetworkConfigSDK_WithVpc(t *testing.T) {
	subnetIds, _ := types.SetValueFrom(context.Background(), types.StringType, []string{"subnet-1", "subnet-2"})
	sgIds, _ := types.SetValueFrom(context.Background(), types.StringType, []string{"sg-1"})
	nc := &agsNetworkConfigModel{
		NetworkMode: types.StringValue("VPC"),
		VpcConfig: &agsVpcConfigModel{
			SubnetIds:        subnetIds,
			SecurityGroupIds: sgIds,
		},
	}
	sdk := buildNetworkConfigSDK(nc)
	if sdk == nil {
		t.Fatal("expected non-nil NetworkConfiguration")
	}
	if sdk.VpcConfig == nil {
		t.Fatal("expected non-nil VpcConfig")
	}
	if len(sdk.VpcConfig.SubnetIds) != 2 {
		t.Fatalf("expected 2 subnet IDs, got %d", len(sdk.VpcConfig.SubnetIds))
	}
	if len(sdk.VpcConfig.SecurityGroupIds) != 1 {
		t.Fatalf("expected 1 security group ID, got %d", len(sdk.VpcConfig.SecurityGroupIds))
	}
}

func TestBuildNetworkConfigSDK_Nil(t *testing.T) {
	sdk := buildNetworkConfigSDK(nil)
	if sdk != nil {
		t.Fatal("expected nil for nil input")
	}
}

func TestBuildNetworkConfigState(t *testing.T) {
	networkMode := "SANDBOX"
	sdk := &ags.NetworkConfiguration{
		NetworkMode: &networkMode,
	}
	state := buildNetworkConfigState(sdk)
	if state == nil {
		t.Fatal("expected non-nil state")
	}
	if state.NetworkMode.ValueString() != "SANDBOX" {
		t.Fatalf("expected SANDBOX, got %s", state.NetworkMode.ValueString())
	}
}

func TestBuildNetworkConfigState_Nil(t *testing.T) {
	state := buildNetworkConfigState(nil)
	if state != nil {
		t.Fatal("expected nil for nil input")
	}
}

func TestBuildTagsSDK(t *testing.T) {
	ctx := context.Background()
	tags := []agsTagModel{
		{Key: types.StringValue("env"), Value: types.StringValue("prod")},
		{Key: types.StringValue("team"), Value: types.StringValue("backend")},
	}
	list, listDiags := types.ListValueFrom(ctx,
		types.ObjectType{AttrTypes: map[string]attr.Type{"key": types.StringType, "value": types.StringType}},
		tags)
	if listDiags.HasError() {
		t.Fatalf("failed to create list: %v", listDiags)
	}
	var d diag.Diagnostics
	sdkTags := buildTagsSDK(ctx, list, &d)
	if d.HasError() {
		t.Fatalf("unexpected diags: %v", d)
	}
	if len(sdkTags) != 2 {
		t.Fatalf("expected 2 tags, got %d", len(sdkTags))
	}
	if *sdkTags[0].Key != "env" {
		t.Fatalf("expected first tag key 'env', got %q", *sdkTags[0].Key)
	}
	if *sdkTags[0].Value != "prod" {
		t.Fatalf("expected first tag value 'prod', got %q", *sdkTags[0].Value)
	}
}

func TestBuildStorageSourceSDK(t *testing.T) {
	src := &agsStorageSourceModel{
		Cos: &agsCosStorageModel{
			Endpoint:   types.StringValue("cos.ap-guangzhou.myqcloud.com"),
			BucketName: types.StringValue("my-bucket"),
			BucketPath: types.StringValue("/data"),
		},
	}
	sdk := buildStorageSourceSDK(src)
	if sdk == nil {
		t.Fatal("expected non-nil StorageSource")
	}
	if sdk.Cos == nil {
		t.Fatal("expected non-nil Cos")
	}
	if *sdk.Cos.BucketName != "my-bucket" {
		t.Fatalf("expected bucket my-bucket, got %q", *sdk.Cos.BucketName)
	}
	if sdk.Image != nil {
		t.Fatal("expected nil Image when not configured")
	}
}

func TestBuildStorageSourceSDK_Nil(t *testing.T) {
	sdk := buildStorageSourceSDK(nil)
	if sdk != nil {
		t.Fatal("expected nil for nil input")
	}
}

func TestBuildCustomConfigState(t *testing.T) {
	image := "nginx:latest"
	regType := "personal"
	cfg := &ags.CustomConfigurationDetail{
		Image:             &image,
		ImageRegistryType: &regType,
	}
	state := buildCustomConfigState(cfg)
	if state == nil {
		t.Fatal("expected non-nil state")
	}
	if state.Image.ValueString() != "nginx:latest" {
		t.Fatalf("expected nginx:latest, got %s", state.Image.ValueString())
	}
}

func TestBuildCustomConfigState_Nil(t *testing.T) {
	state := buildCustomConfigState(nil)
	if state != nil {
		t.Fatal("expected nil for nil input")
	}
}

func TestBuildLogConfigState(t *testing.T) {
	topicId := "topic-12345"
	cfg := &ags.LogConfiguration{
		CLSConfig: &ags.CLSConfig{
			TopicId: &topicId,
		},
	}
	state := buildLogConfigState(cfg)
	if state == nil {
		t.Fatal("expected non-nil state")
	}
	if state.CLSConfig == nil {
		t.Fatal("expected non-nil CLSConfig")
	}
	if state.CLSConfig.TopicId.ValueString() != "topic-12345" {
		t.Fatalf("expected topic-12345, got %s", state.CLSConfig.TopicId.ValueString())
	}
}

func TestBuildLogConfigState_Nil(t *testing.T) {
	state := buildLogConfigState(nil)
	if state != nil {
		t.Fatal("expected nil for nil input")
	}
}

func TestBuildComputerConfigObjState(t *testing.T) {
	imageId := "img-abc"
	cfg := &ags.ComputerConfiguration{
		WAAConfiguration: &ags.WAAConfiguration{
			ImageId: &imageId,
		},
	}
	state := buildComputerConfigObjState(cfg)
	if state == nil {
		t.Fatal("expected non-nil state")
	}
	if state.WAAConfiguration == nil {
		t.Fatal("expected non-nil WAAConfiguration")
	}
	if state.WAAConfiguration.ImageId.ValueString() != "img-abc" {
		t.Fatalf("expected img-abc, got %s", state.WAAConfiguration.ImageId.ValueString())
	}
}

func TestRefreshModel(t *testing.T) {
	m := &AgsSandboxToolResourceModel{}
	toolId := "tool-xyz"
	toolName := "my-tool"
	toolType := "browser"
	status := "ACTIVE"
	desc := "a test tool"
	createTime := "2024-01-01T00:00:00Z"
	updateTime := "2024-01-02T00:00:00Z"
	defaultTimeout := uint64(600)

	tool := &ags.SandboxTool{
		ToolId:                &toolId,
		ToolName:              &toolName,
		ToolType:              &toolType,
		Status:                &status,
		Description:           &desc,
		DefaultTimeoutSeconds: &defaultTimeout,
		CreateTime:            &createTime,
		UpdateTime:            &updateTime,
		NetworkConfiguration:  &ags.NetworkConfiguration{},
	}

	var d diag.Diagnostics
	refreshModel(context.Background(), m, tool, &d)
	if d.HasError() {
		t.Fatalf("unexpected errors: %v", d)
	}
	if m.ToolId.ValueString() != toolId {
		t.Fatalf("expected tool_id %q, got %s", toolId, m.ToolId.ValueString())
	}
	if m.ToolName.ValueString() != toolName {
		t.Fatalf("expected tool_name %q, got %s", toolName, m.ToolName.ValueString())
	}
	if m.DefaultTimeout.IsNull() {
		t.Fatal("expected DefaultTimeout to be set from DefaultTimeoutSeconds")
	}
	if m.DefaultTimeout.ValueString() != "600s" {
		t.Fatalf("expected DefaultTimeout '600s', got %s", m.DefaultTimeout.ValueString())
	}
}
