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

type mockMetaPrefetchTask struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaPrefetchTask) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaPrefetchTask{}

func newMockMetaPrefetchTask() *mockMetaPrefetchTask {
	return &mockMetaPrefetchTask{client: &connectivity.TencentCloudClient{}}
}

// go test ./tencentcloud/services/teo/ -run "TestPrefetchTaskOperation" -v -count=1 -gcflags="all=-l"

// TestPrefetchTaskOperation_Success tests Create with successful prefetch task
func TestPrefetchTaskOperation_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaPrefetchTask().client, "UseTeoV20220901Client", teoClient)

	// Patch CreatePrefetchTask to return JobId
	patches.ApplyMethodFunc(teoClient, "CreatePrefetchTask", func(request *teov20220901.CreatePrefetchTaskRequest) (*teov20220901.CreatePrefetchTaskResponse, error) {
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		resp := teov20220901.NewCreatePrefetchTaskResponse()
		resp.Response = &teov20220901.CreatePrefetchTaskResponseParams{
			JobId:     ptrStringPrefetchTask("job-87654321"),
			RequestId: ptrStringPrefetchTask("fake-request-id"),
		}
		return resp, nil
	})

	// Patch DescribePrefetchTasks to return success
	patches.ApplyMethodFunc(teoClient, "DescribePrefetchTasks", func(request *teov20220901.DescribePrefetchTasksRequest) (*teov20220901.DescribePrefetchTasksResponse, error) {
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		resp := teov20220901.NewDescribePrefetchTasksResponse()
		resp.Response = &teov20220901.DescribePrefetchTasksResponseParams{
			TotalCount: ptrUint64PrefetchTask(1),
			Tasks: []*teov20220901.Task{
				{
					JobId:      ptrStringPrefetchTask("job-87654321"),
					Target:     ptrStringPrefetchTask("http://www.example.com/example.txt"),
					Type:       ptrStringPrefetchTask("prefetch"),
					Method:     ptrStringPrefetchTask("default"),
					Status:     ptrStringPrefetchTask("success"),
					CreateTime: ptrStringPrefetchTask("2025-01-01T00:00:00Z"),
					UpdateTime: ptrStringPrefetchTask("2025-01-01T00:01:00Z"),
				},
			},
			RequestId: ptrStringPrefetchTask("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaPrefetchTask()
	res := teo.ResourceTencentCloudTeoPrefetchTaskOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"targets": []interface{}{"http://www.example.com/example.txt"},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-12345678#job-87654321", d.Id())
	assert.Equal(t, "job-87654321", d.Get("job_id").(string))
}

// TestPrefetchTaskOperation_Failure tests Create when prefetch task fails
func TestPrefetchTaskOperation_Failure(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaPrefetchTask().client, "UseTeoV20220901Client", teoClient)

	// Patch CreatePrefetchTask to return JobId
	patches.ApplyMethodFunc(teoClient, "CreatePrefetchTask", func(request *teov20220901.CreatePrefetchTaskRequest) (*teov20220901.CreatePrefetchTaskResponse, error) {
		resp := teov20220901.NewCreatePrefetchTaskResponse()
		resp.Response = &teov20220901.CreatePrefetchTaskResponseParams{
			JobId:     ptrStringPrefetchTask("job-failed-123"),
			RequestId: ptrStringPrefetchTask("fake-request-id"),
		}
		return resp, nil
	})

	// Patch DescribePrefetchTasks to return failure
	patches.ApplyMethodFunc(teoClient, "DescribePrefetchTasks", func(request *teov20220901.DescribePrefetchTasksRequest) (*teov20220901.DescribePrefetchTasksResponse, error) {
		resp := teov20220901.NewDescribePrefetchTasksResponse()
		resp.Response = &teov20220901.DescribePrefetchTasksResponseParams{
			TotalCount: ptrUint64PrefetchTask(1),
			Tasks: []*teov20220901.Task{
				{
					JobId:       ptrStringPrefetchTask("job-failed-123"),
					Target:      ptrStringPrefetchTask("http://www.example.com/example.txt"),
					Status:      ptrStringPrefetchTask("failed"),
					FailType:    ptrStringPrefetchTask("originPullFailed"),
					FailMessage: ptrStringPrefetchTask("origin server returned 404"),
				},
			},
			RequestId: ptrStringPrefetchTask("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaPrefetchTask()
	res := teo.ResourceTencentCloudTeoPrefetchTaskOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"targets": []interface{}{"http://www.example.com/example.txt"},
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "originPullFailed")
}

// TestPrefetchTaskOperation_APIError tests Create handles CreatePrefetchTask API error
func TestPrefetchTaskOperation_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaPrefetchTask().client, "UseTeoV20220901Client", teoClient)

	// Patch CreatePrefetchTask to return error
	patches.ApplyMethodFunc(teoClient, "CreatePrefetchTask", func(request *teov20220901.CreatePrefetchTaskRequest) (*teov20220901.CreatePrefetchTaskResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=Zone not found")
	})

	meta := newMockMetaPrefetchTask()
	res := teo.ResourceTencentCloudTeoPrefetchTaskOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-invalid",
		"targets": []interface{}{"http://www.example.com/example.txt"},
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
}

// TestPrefetchTaskOperation_Read tests Read is no-op
func TestPrefetchTaskOperation_Read(t *testing.T) {
	res := teo.ResourceTencentCloudTeoPrefetchTaskOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"targets": []interface{}{"http://www.example.com/example.txt"},
	})
	d.SetId("zone-12345678:tcc:job-87654321")

	err := res.Read(d, nil)
	assert.NoError(t, err)
}

// TestPrefetchTaskOperation_Delete tests Delete is no-op
func TestPrefetchTaskOperation_Delete(t *testing.T) {
	res := teo.ResourceTencentCloudTeoPrefetchTaskOperation()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
		"targets": []interface{}{"http://www.example.com/example.txt"},
	})
	d.SetId("zone-12345678:tcc:job-87654321")

	err := res.Delete(d, nil)
	assert.NoError(t, err)
}

// TestPrefetchTaskOperation_Schema validates schema definition
func TestPrefetchTaskOperation_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoPrefetchTaskOperation()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.Nil(t, res.Update)
	assert.NotNil(t, res.Delete)

	assert.Contains(t, res.Schema, "zone_id")
	assert.Contains(t, res.Schema, "targets")
	assert.Contains(t, res.Schema, "mode")
	assert.Contains(t, res.Schema, "headers")
	assert.Contains(t, res.Schema, "prefetch_media_segments")
	assert.Contains(t, res.Schema, "job_id")
	assert.Contains(t, res.Schema, "tasks")

	zoneId := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneId.Type)
	assert.True(t, zoneId.Required)
	assert.True(t, zoneId.ForceNew)

	targets := res.Schema["targets"]
	assert.Equal(t, schema.TypeList, targets.Type)
	assert.True(t, targets.Required)
	assert.True(t, targets.ForceNew)

	mode := res.Schema["mode"]
	assert.Equal(t, schema.TypeString, mode.Type)
	assert.True(t, mode.Optional)
	assert.True(t, mode.ForceNew)

	headers := res.Schema["headers"]
	assert.Equal(t, schema.TypeList, headers.Type)
	assert.True(t, headers.Optional)
	assert.True(t, headers.ForceNew)

	prefetchMediaSegments := res.Schema["prefetch_media_segments"]
	assert.Equal(t, schema.TypeString, prefetchMediaSegments.Type)
	assert.True(t, prefetchMediaSegments.Optional)
	assert.True(t, prefetchMediaSegments.ForceNew)

	jobId := res.Schema["job_id"]
	assert.Equal(t, schema.TypeString, jobId.Type)
	assert.True(t, jobId.Computed)

	tasks := res.Schema["tasks"]
	assert.Equal(t, schema.TypeList, tasks.Type)
	assert.True(t, tasks.Computed)
}

func ptrStringPrefetchTask(s string) *string {
	return &s
}

func ptrUint64PrefetchTask(v uint64) *uint64 {
	return &v
}
