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

// mockMeta implements tccommon.ProviderMeta
type mockMetaForPurgeTask struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForPurgeTask) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForPurgeTask{}

func newMockMetaForPurgeTask() *mockMetaForPurgeTask {
	return &mockMetaForPurgeTask{client: &connectivity.TencentCloudClient{}}
}

func ptrStringForPurge(s string) *string {
	return &s
}

// go test ./tencentcloud/services/teo/ -run "TestPurgeTaskOperation" -v -count=1 -gcflags="all=-l"

// TestPurgeTaskOperation_PurgeURLSuccess tests creating a purge_url type task successfully
func TestPurgeTaskOperation_PurgeURLSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForPurgeTask().client, "UseTeoV20220901Client", teoClient)

	// Patch CreatePurgeTask to return success
	patches.ApplyMethodFunc(teoClient, "CreatePurgeTask", func(request *teov20220901.CreatePurgeTaskRequest) (*teov20220901.CreatePurgeTaskResponse, error) {
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.Equal(t, "purge_url", *request.Type)

		resp := teov20220901.NewCreatePurgeTaskResponse()
		resp.Response = &teov20220901.CreatePurgeTaskResponseParams{
			JobId:     ptrStringForPurge("job-abc123"),
			RequestId: ptrStringForPurge("fake-request-id"),
		}
		return resp, nil
	})

	// Patch DescribePurgeTasks to return success status
	patches.ApplyMethodFunc(teoClient, "DescribePurgeTasks", func(request *teov20220901.DescribePurgeTasksRequest) (*teov20220901.DescribePurgeTasksResponse, error) {
		assert.Equal(t, "zone-12345678", *request.ZoneId)

		resp := teov20220901.NewDescribePurgeTasksResponse()
		resp.Response = &teov20220901.DescribePurgeTasksResponseParams{
			TotalCount: ptrUint64(1),
			Tasks: []*teov20220901.Task{
				{
					JobId:      ptrStringForPurge("job-abc123"),
					Target:     ptrStringForPurge("https://example.com/path1"),
					Type:       ptrStringForPurge("purge_url"),
					Method:     ptrStringForPurge("invalidate"),
					Status:     ptrStringForPurge("success"),
					CreateTime: ptrStringForPurge("2024-01-01T00:00:00Z"),
					UpdateTime: ptrStringForPurge("2024-01-01T00:01:00Z"),
				},
			},
			RequestId: ptrStringForPurge("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForPurgeTask()
	res := teo.ResourceTencentCloudTeoPurgeTaskOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"type":    "purge_url",
		"targets": []interface{}{"https://example.com/path1", "https://example.com/path2"},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())
	assert.Equal(t, "job-abc123", d.Get("job_id").(string))

	tasks := d.Get("tasks").([]interface{})
	assert.Len(t, tasks, 1)
	taskMap := tasks[0].(map[string]interface{})
	assert.Equal(t, "job-abc123", taskMap["job_id"])
	assert.Equal(t, "https://example.com/path1", taskMap["target"])
	assert.Equal(t, "purge_url", taskMap["type"])
	assert.Equal(t, "success", taskMap["status"])
}

// TestPurgeTaskOperation_PurgeCacheTagSuccess tests creating a purge_cache_tag type task successfully
func TestPurgeTaskOperation_PurgeCacheTagSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForPurgeTask().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreatePurgeTask", func(request *teov20220901.CreatePurgeTaskRequest) (*teov20220901.CreatePurgeTaskResponse, error) {
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.Equal(t, "purge_cache_tag", *request.Type)
		assert.NotNil(t, request.CacheTag)
		assert.Len(t, request.CacheTag.Domains, 2)

		resp := teov20220901.NewCreatePurgeTaskResponse()
		resp.Response = &teov20220901.CreatePurgeTaskResponseParams{
			JobId:     ptrStringForPurge("job-cachetag-456"),
			RequestId: ptrStringForPurge("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribePurgeTasks", func(request *teov20220901.DescribePurgeTasksRequest) (*teov20220901.DescribePurgeTasksResponse, error) {
		resp := teov20220901.NewDescribePurgeTasksResponse()
		resp.Response = &teov20220901.DescribePurgeTasksResponseParams{
			TotalCount: ptrUint64(1),
			Tasks: []*teov20220901.Task{
				{
					JobId:      ptrStringForPurge("job-cachetag-456"),
					Target:     ptrStringForPurge("tag1"),
					Type:       ptrStringForPurge("purge_cache_tag"),
					Status:     ptrStringForPurge("success"),
					CreateTime: ptrStringForPurge("2024-01-01T00:00:00Z"),
					UpdateTime: ptrStringForPurge("2024-01-01T00:01:00Z"),
				},
			},
			RequestId: ptrStringForPurge("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForPurgeTask()
	res := teo.ResourceTencentCloudTeoPurgeTaskOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"type":    "purge_cache_tag",
		"cache_tag": []interface{}{
			map[string]interface{}{
				"domains": []interface{}{"example.com", "www.example.com"},
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())
	assert.Equal(t, "job-cachetag-456", d.Get("job_id").(string))
}

// TestPurgeTaskOperation_APIError tests Create handles API error
func TestPurgeTaskOperation_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForPurgeTask().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreatePurgeTask", func(request *teov20220901.CreatePurgeTaskRequest) (*teov20220901.CreatePurgeTaskResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid zone_id")
	})

	meta := newMockMetaForPurgeTask()
	res := teo.ResourceTencentCloudTeoPurgeTaskOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-invalid",
		"type":    "purge_url",
		"targets": []interface{}{"https://example.com/path1"},
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestPurgeTaskOperation_Read tests Read is no-op
func TestPurgeTaskOperation_Read(t *testing.T) {
	res := teo.ResourceTencentCloudTeoPurgeTaskOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"type":    "purge_url",
	})
	d.SetId("test-id")

	err := res.Read(d, nil)
	assert.NoError(t, err)
}

// TestPurgeTaskOperation_Delete tests Delete is no-op
func TestPurgeTaskOperation_Delete(t *testing.T) {
	res := teo.ResourceTencentCloudTeoPurgeTaskOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"type":    "purge_url",
	})
	d.SetId("test-id")

	err := res.Delete(d, nil)
	assert.NoError(t, err)
}

// TestPurgeTaskOperation_Schema validates schema definition
func TestPurgeTaskOperation_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoPurgeTaskOperation()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.Nil(t, res.Update)
	assert.NotNil(t, res.Delete)

	assert.Contains(t, res.Schema, "zone_id")
	assert.Contains(t, res.Schema, "type")
	assert.Contains(t, res.Schema, "method")
	assert.Contains(t, res.Schema, "targets")
	assert.Contains(t, res.Schema, "cache_tag")
	assert.Contains(t, res.Schema, "job_id")
	assert.Contains(t, res.Schema, "tasks")

	zoneId := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneId.Type)
	assert.True(t, zoneId.Required)
	assert.True(t, zoneId.ForceNew)

	purgeType := res.Schema["type"]
	assert.Equal(t, schema.TypeString, purgeType.Type)
	assert.True(t, purgeType.Required)
	assert.True(t, purgeType.ForceNew)

	method := res.Schema["method"]
	assert.Equal(t, schema.TypeString, method.Type)
	assert.True(t, method.Optional)
	assert.True(t, method.ForceNew)

	targets := res.Schema["targets"]
	assert.Equal(t, schema.TypeList, targets.Type)
	assert.True(t, targets.Optional)
	assert.True(t, targets.ForceNew)

	cacheTag := res.Schema["cache_tag"]
	assert.Equal(t, schema.TypeList, cacheTag.Type)
	assert.True(t, cacheTag.Optional)
	assert.True(t, cacheTag.ForceNew)
	assert.Equal(t, 1, cacheTag.MaxItems)

	jobId := res.Schema["job_id"]
	assert.Equal(t, schema.TypeString, jobId.Type)
	assert.True(t, jobId.Computed)

	tasks := res.Schema["tasks"]
	assert.Equal(t, schema.TypeList, tasks.Type)
	assert.True(t, tasks.Computed)
}
