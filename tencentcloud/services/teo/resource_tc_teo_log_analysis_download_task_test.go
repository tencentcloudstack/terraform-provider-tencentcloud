package teo_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svcteo "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// go test ./tencentcloud/services/teo/ -run "TestTeoLogAnalysisDownloadTask" -v -count=1 -gcflags="all=-l"

type mockMetaForLogAnalysisDownloadTask struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForLogAnalysisDownloadTask) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForLogAnalysisDownloadTask{}

func newMockMetaForLogAnalysisDownloadTask() *mockMetaForLogAnalysisDownloadTask {
	return &mockMetaForLogAnalysisDownloadTask{client: &connectivity.TencentCloudClient{}}
}

func ptrStringLogAnalysisTask(s string) *string {
	return &s
}

func ptrInt64LogAnalysisTask(n int64) *int64 {
	return &n
}

// strVal reads a schema value that may have been stored as a string or a *string pointer.
func strVal(v interface{}) string {
	switch t := v.(type) {
	case string:
		return t
	case *string:
		if t != nil {
			return *t
		}
	}
	return ""
}

// buildLogAnalysisDownloadTaskResp builds a fake DescribeLogAnalysisDownloadTasks response containing one task.
func buildLogAnalysisDownloadTaskResp(zoneId, area, taskId string) *teov20220901.DescribeLogAnalysisDownloadTasksResponse {
	resp := teov20220901.NewDescribeLogAnalysisDownloadTasksResponse()
	resp.Response = &teov20220901.DescribeLogAnalysisDownloadTasksResponseParams{
		TotalCount: ptrInt64LogAnalysisTask(1),
		Tasks: []*teov20220901.LogAnalysisDownloadTask{
			{
				TaskId:     ptrStringLogAnalysisTask(taskId),
				ZoneId:     ptrStringLogAnalysisTask(zoneId),
				Area:       ptrStringLogAnalysisTask(area),
				StartTime:  ptrStringLogAnalysisTask("2020-04-29T00:00:00Z"),
				EndTime:    ptrStringLogAnalysisTask("2020-04-30T00:00:00Z"),
				LogType:    ptrStringLogAnalysisTask("l7-access-logs"),
				Condition:  ptrStringLogAnalysisTask(""),
				Format:     ptrStringLogAnalysisTask("csv"),
				Sort:       ptrStringLogAnalysisTask("desc"),
				Status:     ptrStringLogAnalysisTask("completed"),
				CreateTime: ptrStringLogAnalysisTask("2020-04-29T01:00:00Z"),
				Url:        ptrStringLogAnalysisTask("https://example.com/log.csv"),
				ExpireTime: ptrStringLogAnalysisTask("2020-05-02T01:00:00Z"),
			},
		},
		RequestId: ptrStringLogAnalysisTask("fake-request-id"),
	}
	return resp
}

// TestTeoLogAnalysisDownloadTask_Create tests creating a log analysis download task successfully
func TestTeoLogAnalysisDownloadTask_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForLogAnalysisDownloadTask().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateLogAnalysisDownloadTaskWithContext", func(_ context.Context, request *teov20220901.CreateLogAnalysisDownloadTaskRequest) (*teov20220901.CreateLogAnalysisDownloadTaskResponse, error) {
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.Equal(t, "mainland", *request.Area)
		assert.Equal(t, "2020-04-29T00:00:00Z", *request.StartTime)
		assert.Equal(t, "2020-04-30T00:00:00Z", *request.EndTime)

		resp := teov20220901.NewCreateLogAnalysisDownloadTaskResponse()
		resp.Response = &teov20220901.CreateLogAnalysisDownloadTaskResponseParams{
			TaskId:    ptrStringLogAnalysisTask("task-abc123"),
			RequestId: ptrStringLogAnalysisTask("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeLogAnalysisDownloadTasks", func(request *teov20220901.DescribeLogAnalysisDownloadTasksRequest) (*teov20220901.DescribeLogAnalysisDownloadTasksResponse, error) {
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.Equal(t, "mainland", *request.Area)
		return buildLogAnalysisDownloadTaskResp("zone-12345678", "mainland", "task-abc123"), nil
	})

	meta := newMockMetaForLogAnalysisDownloadTask()
	res := svcteo.ResourceTencentCloudTeoLogAnalysisDownloadTask()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-12345678",
		"area":       "mainland",
		"start_time": "2020-04-29T00:00:00Z",
		"end_time":   "2020-04-30T00:00:00Z",
		"log_type":   "l7-access-logs",
		"format":     "csv",
		"sort":       "desc",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())
	assert.Equal(t, "zone-12345678#mainland#task-abc123", d.Id())
	assert.Equal(t, "task-abc123", strVal(d.Get("task_id")))
	assert.Equal(t, "completed", strVal(d.Get("status")))
	assert.Equal(t, "https://example.com/log.csv", strVal(d.Get("url")))
}

// TestTeoLogAnalysisDownloadTask_Read tests reading a task and populating state
func TestTeoLogAnalysisDownloadTask_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForLogAnalysisDownloadTask().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeLogAnalysisDownloadTasks", func(request *teov20220901.DescribeLogAnalysisDownloadTasksRequest) (*teov20220901.DescribeLogAnalysisDownloadTasksResponse, error) {
		return buildLogAnalysisDownloadTaskResp("zone-read-test", "overseas", "task-read-123"), nil
	})

	meta := newMockMetaForLogAnalysisDownloadTask()
	res := svcteo.ResourceTencentCloudTeoLogAnalysisDownloadTask()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-read-test",
		"area":    "overseas",
	})
	d.SetId("zone-read-test#overseas#task-read-123")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-read-test#overseas#task-read-123", d.Id())
	assert.Equal(t, "task-read-123", strVal(d.Get("task_id")))
	assert.Equal(t, "completed", strVal(d.Get("status")))
	assert.Equal(t, "2020-04-29T01:00:00Z", strVal(d.Get("create_time")))
	assert.Equal(t, "https://example.com/log.csv", strVal(d.Get("url")))
	assert.Equal(t, "2020-05-02T01:00:00Z", strVal(d.Get("expire_time")))
	assert.Equal(t, "l7-access-logs", strVal(d.Get("log_type")))
	assert.Equal(t, "csv", strVal(d.Get("format")))
	assert.Equal(t, "desc", strVal(d.Get("sort")))
}

// TestTeoLogAnalysisDownloadTask_ReadNotFound tests Read clears state when the task is absent
func TestTeoLogAnalysisDownloadTask_ReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForLogAnalysisDownloadTask().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeLogAnalysisDownloadTasks", func(request *teov20220901.DescribeLogAnalysisDownloadTasksRequest) (*teov20220901.DescribeLogAnalysisDownloadTasksResponse, error) {
		resp := teov20220901.NewDescribeLogAnalysisDownloadTasksResponse()
		resp.Response = &teov20220901.DescribeLogAnalysisDownloadTasksResponseParams{
			TotalCount: ptrInt64LogAnalysisTask(0),
			Tasks:      nil,
		}
		return resp, nil
	})

	meta := newMockMetaForLogAnalysisDownloadTask()
	res := svcteo.ResourceTencentCloudTeoLogAnalysisDownloadTask()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-gone",
		"area":    "mainland",
	})
	d.SetId("zone-gone#mainland#task-gone")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Empty(t, d.Id())
}

// TestTeoLogAnalysisDownloadTask_UpdateImmutable tests that Update returns an error when an immutable field changes.
// schema.TestResourceDataRaw does not track diffs, so HasChange always reports false.
// We patch (*schema.ResourceData).HasChange to simulate a change on the `end_time` field.
func TestTeoLogAnalysisDownloadTask_UpdateImmutable(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForLogAnalysisDownloadTask().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(new(schema.ResourceData), "HasChange", func(d *schema.ResourceData, key string) bool {
		return key == "end_time"
	})

	meta := newMockMetaForLogAnalysisDownloadTask()
	res := svcteo.ResourceTencentCloudTeoLogAnalysisDownloadTask()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-update-test",
		"area":       "mainland",
		"start_time": "2020-04-29T00:00:00Z",
		"end_time":   "2020-04-30T00:00:00Z",
	})
	d.SetId("zone-update-test#mainland#task-update-123")

	err := res.Update(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "end_time")
	assert.Contains(t, err.Error(), "immutable")
}

// TestTeoLogAnalysisDownloadTask_UpdateNoChange tests that Update falls through to Read when no field changes
func TestTeoLogAnalysisDownloadTask_UpdateNoChange(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForLogAnalysisDownloadTask().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeLogAnalysisDownloadTasks", func(request *teov20220901.DescribeLogAnalysisDownloadTasksRequest) (*teov20220901.DescribeLogAnalysisDownloadTasksResponse, error) {
		return buildLogAnalysisDownloadTaskResp("zone-update-test", "mainland", "task-update-123"), nil
	})

	meta := newMockMetaForLogAnalysisDownloadTask()
	res := svcteo.ResourceTencentCloudTeoLogAnalysisDownloadTask()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-update-test",
		"area":       "mainland",
		"start_time": "2020-04-29T00:00:00Z",
		"end_time":   "2020-04-30T00:00:00Z",
	})
	d.SetId("zone-update-test#mainland#task-update-123")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

// TestTeoLogAnalysisDownloadTask_Delete tests Delete is a no-op that does not call any cloud API
func TestTeoLogAnalysisDownloadTask_Delete(t *testing.T) {
	res := svcteo.ResourceTencentCloudTeoLogAnalysisDownloadTask()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":    "zone-delete-test",
		"area":       "mainland",
		"start_time": "2020-04-29T00:00:00Z",
		"end_time":   "2020-04-30T00:00:00Z",
	})
	d.SetId("zone-delete-test#mainland#task-delete-123")

	err := res.Delete(d, nil)
	assert.NoError(t, err)
	assert.Equal(t, "zone-delete-test#mainland#task-delete-123", d.Id())
}

// TestTeoLogAnalysisDownloadTask_Schema validates the schema definition
func TestTeoLogAnalysisDownloadTask_Schema(t *testing.T) {
	res := svcteo.ResourceTencentCloudTeoLogAnalysisDownloadTask()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.NotNil(t, res.Update)
	assert.NotNil(t, res.Delete)
	assert.NotNil(t, res.Importer)

	inputFields := []string{"zone_id", "area", "start_time", "end_time", "log_type", "condition", "format", "sort"}
	for _, f := range inputFields {
		s, ok := res.Schema[f]
		assert.True(t, ok, f)
		assert.Equal(t, schema.TypeString, s.Type, f)
		assert.True(t, s.ForceNew, f)
	}
	for _, f := range []string{"zone_id", "area", "start_time", "end_time"} {
		assert.True(t, res.Schema[f].Required, f)
	}
	for _, f := range []string{"log_type", "condition", "format", "sort"} {
		assert.True(t, res.Schema[f].Optional, f)
	}
	for _, f := range []string{"task_id", "status", "create_time", "url", "expire_time"} {
		s, ok := res.Schema[f]
		assert.True(t, ok, f)
		assert.True(t, s.Computed, f)
	}
}
