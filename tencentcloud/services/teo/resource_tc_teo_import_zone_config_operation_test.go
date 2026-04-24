package teo_test

import (
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// mockMetaImportZoneConfig implements tccommon.ProviderMeta
type mockMetaImportZoneConfig struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaImportZoneConfig) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaImportZoneConfig{}

func newMockMetaImportZoneConfig() *mockMetaImportZoneConfig {
	return &mockMetaImportZoneConfig{client: &connectivity.TencentCloudClient{}}
}

// go test ./tencentcloud/services/teo/ -run "TestImportZoneConfigOperation" -v -count=1 -gcflags="all=-l"

// TestImportZoneConfigOperation_Success tests Create with successful import
func TestImportZoneConfigOperation_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaImportZoneConfig().client, "UseTeoV20220901Client", teoClient)

	// Patch ImportZoneConfig to return TaskId
	patches.ApplyMethodFunc(teoClient, "ImportZoneConfig", func(request *teov20220901.ImportZoneConfigRequest) (*teov20220901.ImportZoneConfigResponse, error) {
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.Contains(t, *request.Content, "GlobalAccelerate")
		resp := teov20220901.NewImportZoneConfigResponse()
		resp.Response = &teov20220901.ImportZoneConfigResponseParams{
			TaskId:    ptrStringImportZoneConfig("task-87654321"),
			RequestId: ptrStringImportZoneConfig("fake-request-id"),
		}
		return resp, nil
	})

	// Patch DescribeZoneConfigImportResult to return success
	patches.ApplyMethodFunc(teoClient, "DescribeZoneConfigImportResult", func(request *teov20220901.DescribeZoneConfigImportResultRequest) (*teov20220901.DescribeZoneConfigImportResultResponse, error) {
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.Equal(t, "task-87654321", *request.TaskId)
		resp := teov20220901.NewDescribeZoneConfigImportResultResponse()
		resp.Response = &teov20220901.DescribeZoneConfigImportResultResponseParams{
			Status:     ptrStringImportZoneConfig("success"),
			Message:    ptrStringImportZoneConfig(""),
			ImportTime: ptrStringImportZoneConfig("2025-01-01T00:00:00Z"),
			FinishTime: ptrStringImportZoneConfig("2025-01-01T00:01:00Z"),
			RequestId:  ptrStringImportZoneConfig("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaImportZoneConfig()
	res := teo.ResourceTencentCloudTeoImportZoneConfigOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"content": `{"GlobalAccelerate":{}}`,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-12345678", d.Id())
	assert.Equal(t, "task-87654321", d.Get("task_id").(string))
	assert.Equal(t, "success", d.Get("status").(string))
	assert.Equal(t, "", d.Get("message").(string))
	assert.Equal(t, "2025-01-01T00:00:00Z", d.Get("import_time").(string))
	assert.Equal(t, "2025-01-01T00:01:00Z", d.Get("finish_time").(string))
}

// TestImportZoneConfigOperation_Failure tests Create when import task fails
func TestImportZoneConfigOperation_Failure(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaImportZoneConfig().client, "UseTeoV20220901Client", teoClient)

	// Patch ImportZoneConfig to return TaskId
	patches.ApplyMethodFunc(teoClient, "ImportZoneConfig", func(request *teov20220901.ImportZoneConfigRequest) (*teov20220901.ImportZoneConfigResponse, error) {
		resp := teov20220901.NewImportZoneConfigResponse()
		resp.Response = &teov20220901.ImportZoneConfigResponseParams{
			TaskId:    ptrStringImportZoneConfig("task-failed-123"),
			RequestId: ptrStringImportZoneConfig("fake-request-id"),
		}
		return resp, nil
	})

	// Patch DescribeZoneConfigImportResult to return failure
	patches.ApplyMethodFunc(teoClient, "DescribeZoneConfigImportResult", func(request *teov20220901.DescribeZoneConfigImportResultRequest) (*teov20220901.DescribeZoneConfigImportResultResponse, error) {
		resp := teov20220901.NewDescribeZoneConfigImportResultResponse()
		resp.Response = &teov20220901.DescribeZoneConfigImportResultResponseParams{
			Status:    ptrStringImportZoneConfig("failure"),
			Message:   ptrStringImportZoneConfig("invalid configuration content"),
			RequestId: ptrStringImportZoneConfig("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaImportZoneConfig()
	res := teo.ResourceTencentCloudTeoImportZoneConfigOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"content": `{"invalid": true}`,
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "invalid configuration content")
}

// TestImportZoneConfigOperation_APIError tests Create handles ImportZoneConfig API error
func TestImportZoneConfigOperation_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaImportZoneConfig().client, "UseTeoV20220901Client", teoClient)

	// Patch ImportZoneConfig to return error
	patches.ApplyMethodFunc(teoClient, "ImportZoneConfig", func(request *teov20220901.ImportZoneConfigRequest) (*teov20220901.ImportZoneConfigResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=Zone not found")
	})

	meta := newMockMetaImportZoneConfig()
	res := teo.ResourceTencentCloudTeoImportZoneConfigOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-invalid",
		"content": `{"test": true}`,
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
}

// TestImportZoneConfigOperation_Read tests Read is no-op
func TestImportZoneConfigOperation_Read(t *testing.T) {
	res := teo.ResourceTencentCloudTeoImportZoneConfigOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"content": `{"test": true}`,
	})
	d.SetId("zone-12345678")

	err := res.Read(d, nil)
	assert.NoError(t, err)
}

// TestImportZoneConfigOperation_Delete tests Delete is no-op
func TestImportZoneConfigOperation_Delete(t *testing.T) {
	res := teo.ResourceTencentCloudTeoImportZoneConfigOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"content": `{"test": true}`,
	})
	d.SetId("zone-12345678")

	err := res.Delete(d, nil)
	assert.NoError(t, err)
}

// TestImportZoneConfigOperation_Schema validates schema definition
func TestImportZoneConfigOperation_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoImportZoneConfigOperation()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.Nil(t, res.Update)
	assert.NotNil(t, res.Delete)

	assert.Contains(t, res.Schema, "zone_id")
	assert.Contains(t, res.Schema, "content")
	assert.Contains(t, res.Schema, "task_id")
	assert.Contains(t, res.Schema, "status")
	assert.Contains(t, res.Schema, "message")
	assert.Contains(t, res.Schema, "import_time")
	assert.Contains(t, res.Schema, "finish_time")

	zoneId := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneId.Type)
	assert.True(t, zoneId.Required)
	assert.True(t, zoneId.ForceNew)

	content := res.Schema["content"]
	assert.Equal(t, schema.TypeString, content.Type)
	assert.True(t, content.Required)
	assert.True(t, content.ForceNew)

	taskId := res.Schema["task_id"]
	assert.Equal(t, schema.TypeString, taskId.Type)
	assert.True(t, taskId.Computed)

	status := res.Schema["status"]
	assert.Equal(t, schema.TypeString, status.Type)
	assert.True(t, status.Computed)

	message := res.Schema["message"]
	assert.Equal(t, schema.TypeString, message.Type)
	assert.True(t, message.Computed)

	importTime := res.Schema["import_time"]
	assert.Equal(t, schema.TypeString, importTime.Type)
	assert.True(t, importTime.Computed)

	finishTime := res.Schema["finish_time"]
	assert.Equal(t, schema.TypeString, finishTime.Type)
	assert.True(t, finishTime.Computed)
}

func ptrStringImportZoneConfig(s string) *string {
	return &s
}
