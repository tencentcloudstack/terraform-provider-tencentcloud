package ags_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	agsv20250920 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ags/v20250920"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svcags "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ags"
)

// go test ./tencentcloud/services/ags/ -run "TestAgsSandboxTool" -v -count=1 -gcflags="all=-l"

type mockMetaSandboxTool struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaSandboxTool) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaSandboxTool{}

func newMockMetaSandboxTool() *mockMetaSandboxTool {
	return &mockMetaSandboxTool{client: &connectivity.TencentCloudClient{}}
}

func ptrStr(s string) *string {
	return &s
}

// TestAgsSandboxToolResource_Create tests resource creation
func TestAgsSandboxToolResource_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	agsClient := &agsv20250920.Client{}
	patches.ApplyMethodReturn(newMockMetaSandboxTool().client, "UseAgsV20250920Client", agsClient)

	patches.ApplyMethodFunc(agsClient, "CreateSandboxToolWithContext", func(ctx interface{}, request *agsv20250920.CreateSandboxToolRequest) (*agsv20250920.CreateSandboxToolResponse, error) {
		resp := agsv20250920.NewCreateSandboxToolResponse()
		resp.Response = &agsv20250920.CreateSandboxToolResponseParams{
			ToolId:    ptrStr("sbt-test12345"),
			RequestId: ptrStr("fake-request-id"),
		}
		return resp, nil
	})

	// Mock the Read after create
	patches.ApplyMethodFunc(agsClient, "DescribeSandboxToolList", func(request *agsv20250920.DescribeSandboxToolListRequest) (*agsv20250920.DescribeSandboxToolListResponse, error) {
		resp := agsv20250920.NewDescribeSandboxToolListResponse()
		resp.Response = &agsv20250920.DescribeSandboxToolListResponseParams{
			SandboxToolSet: []*agsv20250920.SandboxTool{
				{
					ToolId:                ptrStr("sbt-test12345"),
					ToolName:              ptrStr("test-tool"),
					ToolType:              ptrStr("custom"),
					Status:                ptrStr("ACTIVE"),
					Description:           ptrStr("test description"),
					DefaultTimeoutSeconds: func() *uint64 { v := uint64(600); return &v }(),
					CreateTime:            ptrStr("2025-01-01T00:00:00Z"),
					UpdateTime:            ptrStr("2025-01-01T00:00:00Z"),
				},
			},
			TotalCount: func() *int64 { v := int64(1); return &v }(),
			RequestId:  ptrStr("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaSandboxTool()
	res := svcags.ResourceTencentCloudAgsSandboxTool()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"tool_name":       "test-tool",
		"tool_type":       "custom",
		"description":     "test description",
		"default_timeout": "10m",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "sbt-test12345", d.Id())
	assert.Equal(t, "ACTIVE", d.Get("status"))
	assert.Equal(t, "test-tool", d.Get("tool_name"))
}

// TestAgsSandboxToolResource_Read tests resource read
func TestAgsSandboxToolResource_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	agsClient := &agsv20250920.Client{}
	patches.ApplyMethodReturn(newMockMetaSandboxTool().client, "UseAgsV20250920Client", agsClient)

	patches.ApplyMethodFunc(agsClient, "DescribeSandboxToolList", func(request *agsv20250920.DescribeSandboxToolListRequest) (*agsv20250920.DescribeSandboxToolListResponse, error) {
		resp := agsv20250920.NewDescribeSandboxToolListResponse()
		resp.Response = &agsv20250920.DescribeSandboxToolListResponseParams{
			SandboxToolSet: []*agsv20250920.SandboxTool{
				{
					ToolId:                ptrStr("sbt-test12345"),
					ToolName:              ptrStr("test-tool"),
					ToolType:              ptrStr("custom"),
					Status:                ptrStr("ACTIVE"),
					Description:           ptrStr("test description"),
					DefaultTimeoutSeconds: func() *uint64 { v := uint64(600); return &v }(),
					CreateTime:            ptrStr("2025-01-01T00:00:00Z"),
					UpdateTime:            ptrStr("2025-01-01T00:00:00Z"),
					StatusReason:          ptrStr(""),
				},
			},
			TotalCount: func() *int64 { v := int64(1); return &v }(),
			RequestId:  ptrStr("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaSandboxTool()
	res := svcags.ResourceTencentCloudAgsSandboxTool()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("sbt-test12345")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "sbt-test12345", d.Id())
	assert.Equal(t, "test-tool", d.Get("tool_name"))
	assert.Equal(t, "ACTIVE", d.Get("status"))
	assert.Equal(t, "test description", d.Get("description"))
	assert.Equal(t, "2025-01-01T00:00:00Z", d.Get("create_time"))
}

// TestAgsSandboxToolResource_ReadEmpty tests resource read when empty
func TestAgsSandboxToolResource_ReadEmpty(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	agsClient := &agsv20250920.Client{}
	patches.ApplyMethodReturn(newMockMetaSandboxTool().client, "UseAgsV20250920Client", agsClient)

	patches.ApplyMethodFunc(agsClient, "DescribeSandboxToolList", func(request *agsv20250920.DescribeSandboxToolListRequest) (*agsv20250920.DescribeSandboxToolListResponse, error) {
		resp := agsv20250920.NewDescribeSandboxToolListResponse()
		resp.Response = &agsv20250920.DescribeSandboxToolListResponseParams{
			SandboxToolSet: []*agsv20250920.SandboxTool{},
			TotalCount:     func() *int64 { v := int64(0); return &v }(),
			RequestId:      ptrStr("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaSandboxTool()
	res := svcags.ResourceTencentCloudAgsSandboxTool()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("sbt-not-exists")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestAgsSandboxToolResource_Update tests resource update
func TestAgsSandboxToolResource_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	agsClient := &agsv20250920.Client{}
	patches.ApplyMethodReturn(newMockMetaSandboxTool().client, "UseAgsV20250920Client", agsClient)

	updateCalled := false
	patches.ApplyMethodFunc(agsClient, "UpdateSandboxToolWithContext", func(ctx interface{}, request *agsv20250920.UpdateSandboxToolRequest) (*agsv20250920.UpdateSandboxToolResponse, error) {
		updateCalled = true
		resp := agsv20250920.NewUpdateSandboxToolResponse()
		resp.Response = &agsv20250920.UpdateSandboxToolResponseParams{
			RequestId: ptrStr("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(agsClient, "DescribeSandboxToolList", func(request *agsv20250920.DescribeSandboxToolListRequest) (*agsv20250920.DescribeSandboxToolListResponse, error) {
		resp := agsv20250920.NewDescribeSandboxToolListResponse()
		resp.Response = &agsv20250920.DescribeSandboxToolListResponseParams{
			SandboxToolSet: []*agsv20250920.SandboxTool{
				{
					ToolId:                ptrStr("sbt-test12345"),
					ToolName:              ptrStr("test-tool"),
					ToolType:              ptrStr("custom"),
					Status:                ptrStr("ACTIVE"),
					Description:           ptrStr("updated description"),
					DefaultTimeoutSeconds: func() *uint64 { v := uint64(600); return &v }(),
				},
			},
			TotalCount: func() *int64 { v := int64(1); return &v }(),
			RequestId:  ptrStr("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaSandboxTool()
	res := svcags.ResourceTencentCloudAgsSandboxTool()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"tool_name":       "test-tool",
		"tool_type":       "custom",
		"description":     "updated description",
		"default_timeout": "10m",
	})
	d.SetId("sbt-test12345")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.True(t, updateCalled)
	assert.Equal(t, "updated description", d.Get("description"))
}

// TestAgsSandboxToolResource_Delete tests resource delete
func TestAgsSandboxToolResource_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	agsClient := &agsv20250920.Client{}
	patches.ApplyMethodReturn(newMockMetaSandboxTool().client, "UseAgsV20250920Client", agsClient)

	deleteCalled := false
	patches.ApplyMethodFunc(agsClient, "DeleteSandboxToolWithContext", func(ctx interface{}, request *agsv20250920.DeleteSandboxToolRequest) (*agsv20250920.DeleteSandboxToolResponse, error) {
		deleteCalled = true
		resp := agsv20250920.NewDeleteSandboxToolResponse()
		resp.Response = &agsv20250920.DeleteSandboxToolResponseParams{
			RequestId: ptrStr("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaSandboxTool()
	res := svcags.ResourceTencentCloudAgsSandboxTool()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("sbt-test12345")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
	assert.True(t, deleteCalled)
}

// TestAgsSandboxToolResource_CreateNilResponse tests nil response handling
func TestAgsSandboxToolResource_CreateNilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	agsClient := &agsv20250920.Client{}
	patches.ApplyMethodReturn(newMockMetaSandboxTool().client, "UseAgsV20250920Client", agsClient)

	patches.ApplyMethodFunc(agsClient, "CreateSandboxToolWithContext", func(ctx interface{}, request *agsv20250920.CreateSandboxToolRequest) (*agsv20250920.CreateSandboxToolResponse, error) {
		// Return nil response
		return nil, nil
	})

	meta := newMockMetaSandboxTool()
	res := svcags.ResourceTencentCloudAgsSandboxTool()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"tool_name": "test-tool",
		"tool_type": "custom",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
}

// TestAgsSandboxToolResource_Import tests import
func TestAgsSandboxToolResource_Import(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	agsClient := &agsv20250920.Client{}
	patches.ApplyMethodReturn(newMockMetaSandboxTool().client, "UseAgsV20250920Client", agsClient)

	patches.ApplyMethodFunc(agsClient, "DescribeSandboxToolList", func(request *agsv20250920.DescribeSandboxToolListRequest) (*agsv20250920.DescribeSandboxToolListResponse, error) {
		resp := agsv20250920.NewDescribeSandboxToolListResponse()
		resp.Response = &agsv20250920.DescribeSandboxToolListResponseParams{
			SandboxToolSet: []*agsv20250920.SandboxTool{
				{
					ToolId:   ptrStr("sbt-imported"),
					ToolName: ptrStr("imported-tool"),
					ToolType: ptrStr("custom"),
					Status:   ptrStr("ACTIVE"),
				},
			},
			TotalCount: func() *int64 { v := int64(1); return &v }(),
			RequestId:  ptrStr("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaSandboxTool()
	res := svcags.ResourceTencentCloudAgsSandboxTool()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	// Simulate import
	importData, err := res.Importer.State(d, meta)
	assert.NoError(t, err)
	assert.Len(t, importData, 1)
	assert.Equal(t, "sbt-imported", importData[0].Id())
}
