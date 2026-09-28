package cls_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	cls "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	localcls "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cls"
)

type mockMetaForRemoteWriteTask struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForRemoteWriteTask) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForRemoteWriteTask{}

func newMockMetaForRemoteWriteTask() *mockMetaForRemoteWriteTask {
	return &mockMetaForRemoteWriteTask{client: &connectivity.TencentCloudClient{}}
}

func ptrStrRemoteWrite(s string) *string {
	return &s
}

func ptrUint64RemoteWrite(v uint64) *uint64 {
	return &v
}

func ptrInt64RemoteWrite(v int64) *int64 {
	return &v
}

// TestRemoteWriteTask_Create_Basic tests Create with basic parameters
func TestRemoteWriteTask_Create_Basic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &cls.Client{}
	patches.ApplyMethodReturn(newMockMetaForRemoteWriteTask().client, "UseClsClient", clsClient)

	var capturedRequest *cls.CreateRemoteWriteTaskRequest
	patches.ApplyMethodFunc(clsClient, "CreateRemoteWriteTaskWithContext", func(_ context.Context, request *cls.CreateRemoteWriteTaskRequest) (*cls.CreateRemoteWriteTaskResponse, error) {
		capturedRequest = request
		resp := cls.NewCreateRemoteWriteTaskResponse()
		resp.Response = &cls.CreateRemoteWriteTaskResponseParams{
			TaskId:    ptrStrRemoteWrite("task-test-123"),
			RequestId: ptrStrRemoteWrite("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(clsClient, "DescribeRemoteWriteTasksWithContext", func(_ context.Context, request *cls.DescribeRemoteWriteTasksRequest) (*cls.DescribeRemoteWriteTasksResponse, error) {
		resp := cls.NewDescribeRemoteWriteTasksResponse()
		resp.Response = &cls.DescribeRemoteWriteTasksResponseParams{
			Infos: []*cls.RemoteWriteInfo{
				{
					TaskId:         ptrStrRemoteWrite("task-test-123"),
					TopicId:        ptrStrRemoteWrite("topic-test-123"),
					Name:           ptrStrRemoteWrite("test-remote-write"),
					Target:         ptrStrRemoteWrite("prometheus"),
					RemoteWriteURL: ptrStrRemoteWrite("http://prometheus.example.com/api/v1/write"),
					AuthType:       ptrUint64RemoteWrite(0),
					NetType:        ptrUint64RemoteWrite(1),
					Status:         ptrInt64RemoteWrite(1),
					Enable:         ptrUint64RemoteWrite(1),
				},
			},
			TotalCount: ptrUint64RemoteWrite(1),
			RequestId:  ptrStrRemoteWrite("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForRemoteWriteTask()
	res := localcls.ResourceTencentCloudClsRemoteWriteTask()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"topic_id":         "topic-test-123",
		"name":             "test-remote-write",
		"target":           "prometheus",
		"remote_write_url": "http://prometheus.example.com/api/v1/write",
		"auth_type":        0,
		"net_type":         1,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "task-test-123#topic-test-123", d.Id())
	assert.NotNil(t, capturedRequest)
	assert.Equal(t, "topic-test-123", *capturedRequest.TopicId)
	assert.Equal(t, "test-remote-write", *capturedRequest.Name)
	assert.Equal(t, "prometheus", *capturedRequest.Target)
	assert.Equal(t, uint64(0), *capturedRequest.AuthType)
	assert.Equal(t, uint64(1), *capturedRequest.NetType)
}

// TestRemoteWriteTask_Create_WithAuthInfo tests Create with auth_info block
func TestRemoteWriteTask_Create_WithAuthInfo(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &cls.Client{}
	patches.ApplyMethodReturn(newMockMetaForRemoteWriteTask().client, "UseClsClient", clsClient)

	var capturedRequest *cls.CreateRemoteWriteTaskRequest
	patches.ApplyMethodFunc(clsClient, "CreateRemoteWriteTaskWithContext", func(_ context.Context, request *cls.CreateRemoteWriteTaskRequest) (*cls.CreateRemoteWriteTaskResponse, error) {
		capturedRequest = request
		resp := cls.NewCreateRemoteWriteTaskResponse()
		resp.Response = &cls.CreateRemoteWriteTaskResponseParams{
			TaskId:    ptrStrRemoteWrite("task-test-456"),
			RequestId: ptrStrRemoteWrite("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(clsClient, "DescribeRemoteWriteTasksWithContext", func(_ context.Context, request *cls.DescribeRemoteWriteTasksRequest) (*cls.DescribeRemoteWriteTasksResponse, error) {
		resp := cls.NewDescribeRemoteWriteTasksResponse()
		resp.Response = &cls.DescribeRemoteWriteTasksResponseParams{
			Infos: []*cls.RemoteWriteInfo{
				{
					TaskId:         ptrStrRemoteWrite("task-test-456"),
					TopicId:        ptrStrRemoteWrite("topic-test-123"),
					Name:           ptrStrRemoteWrite("test-remote-write"),
					Target:         ptrStrRemoteWrite("prometheus"),
					RemoteWriteURL: ptrStrRemoteWrite("https://prometheus.example.com/api/v1/write"),
					AuthType:       ptrUint64RemoteWrite(1),
					NetType:        ptrUint64RemoteWrite(2),
					AuthInfo: &cls.RemoteWriteAuthInfo{
						Username: ptrStrRemoteWrite("admin"),
						Password: ptrStrRemoteWrite("my-password"),
					},
				},
			},
			TotalCount: ptrUint64RemoteWrite(1),
			RequestId:  ptrStrRemoteWrite("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForRemoteWriteTask()
	res := localcls.ResourceTencentCloudClsRemoteWriteTask()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"topic_id":         "topic-test-123",
		"name":             "test-remote-write",
		"target":           "prometheus",
		"remote_write_url": "https://prometheus.example.com/api/v1/write",
		"auth_type":        1,
		"net_type":         2,
		"auth_info": []interface{}{
			map[string]interface{}{
				"username": "admin",
				"password": "my-password",
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "task-test-456#topic-test-123", d.Id())
	assert.NotNil(t, capturedRequest.AuthInfo)
	assert.Equal(t, "admin", *capturedRequest.AuthInfo.Username)
	assert.Equal(t, "my-password", *capturedRequest.AuthInfo.Password)
}

// TestRemoteWriteTask_Read_Basic tests Read populates fields from API response
func TestRemoteWriteTask_Read_Basic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &cls.Client{}
	patches.ApplyMethodReturn(newMockMetaForRemoteWriteTask().client, "UseClsClient", clsClient)

	patches.ApplyMethodFunc(clsClient, "DescribeRemoteWriteTasksWithContext", func(_ context.Context, request *cls.DescribeRemoteWriteTasksRequest) (*cls.DescribeRemoteWriteTasksResponse, error) {
		resp := cls.NewDescribeRemoteWriteTasksResponse()
		resp.Response = &cls.DescribeRemoteWriteTasksResponseParams{
			Infos: []*cls.RemoteWriteInfo{
				{
					TaskId:         ptrStrRemoteWrite("task-test-789"),
					TopicId:        ptrStrRemoteWrite("topic-test-123"),
					Name:           ptrStrRemoteWrite("test-remote-write"),
					Target:         ptrStrRemoteWrite("prometheus"),
					RemoteWriteURL: ptrStrRemoteWrite("http://prometheus.example.com/api/v1/write"),
					AuthType:       ptrUint64RemoteWrite(0),
					NetType:        ptrUint64RemoteWrite(1),
					Status:         ptrInt64RemoteWrite(1),
					Enable:         ptrUint64RemoteWrite(1),
					CreateTime:     ptrStrRemoteWrite("2024-01-01 00:00:00"),
					UpdateTime:     ptrStrRemoteWrite("2024-01-02 00:00:00"),
					LogsetId:       ptrStrRemoteWrite("logset-test-123"),
				},
			},
			TotalCount: ptrUint64RemoteWrite(1),
			RequestId:  ptrStrRemoteWrite("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForRemoteWriteTask()
	res := localcls.ResourceTencentCloudClsRemoteWriteTask()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("task-test-789#topic-test-123")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "task-test-789#topic-test-123", d.Id())
	assert.Equal(t, "topic-test-123", d.Get("topic_id"))
	assert.Equal(t, "test-remote-write", d.Get("name"))
	assert.Equal(t, "prometheus", d.Get("target"))
	assert.Equal(t, "http://prometheus.example.com/api/v1/write", d.Get("remote_write_url"))
	assert.Equal(t, 0, d.Get("auth_type"))
	assert.Equal(t, 1, d.Get("net_type"))
	assert.Equal(t, 1, d.Get("status"))
	assert.Equal(t, "2024-01-01 00:00:00", d.Get("create_time"))
	assert.Equal(t, "logset-test-123", d.Get("logset_id"))
}

// TestRemoteWriteTask_Read_NilAuthInfo tests Read handles nil AuthInfo
func TestRemoteWriteTask_Read_NilAuthInfo(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &cls.Client{}
	patches.ApplyMethodReturn(newMockMetaForRemoteWriteTask().client, "UseClsClient", clsClient)

	patches.ApplyMethodFunc(clsClient, "DescribeRemoteWriteTasksWithContext", func(_ context.Context, request *cls.DescribeRemoteWriteTasksRequest) (*cls.DescribeRemoteWriteTasksResponse, error) {
		resp := cls.NewDescribeRemoteWriteTasksResponse()
		resp.Response = &cls.DescribeRemoteWriteTasksResponseParams{
			Infos: []*cls.RemoteWriteInfo{
				{
					TaskId:         ptrStrRemoteWrite("task-test-789"),
					TopicId:        ptrStrRemoteWrite("topic-test-123"),
					Name:           ptrStrRemoteWrite("test-remote-write"),
					Target:         ptrStrRemoteWrite("prometheus"),
					RemoteWriteURL: ptrStrRemoteWrite("http://prometheus.example.com/api/v1/write"),
					AuthType:       ptrUint64RemoteWrite(0),
					NetType:        ptrUint64RemoteWrite(1),
					AuthInfo:       nil,
				},
			},
			TotalCount: ptrUint64RemoteWrite(1),
			RequestId:  ptrStrRemoteWrite("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForRemoteWriteTask()
	res := localcls.ResourceTencentCloudClsRemoteWriteTask()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("task-test-789#topic-test-123")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "task-test-789#topic-test-123", d.Id())
	authInfo := d.Get("auth_info").([]interface{})
	assert.Equal(t, 0, len(authInfo))
}

// TestRemoteWriteTask_Read_EmptyResult tests Read clears ID when result is empty
func TestRemoteWriteTask_Read_EmptyResult(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &cls.Client{}
	patches.ApplyMethodReturn(newMockMetaForRemoteWriteTask().client, "UseClsClient", clsClient)

	patches.ApplyMethodFunc(clsClient, "DescribeRemoteWriteTasksWithContext", func(_ context.Context, request *cls.DescribeRemoteWriteTasksRequest) (*cls.DescribeRemoteWriteTasksResponse, error) {
		resp := cls.NewDescribeRemoteWriteTasksResponse()
		resp.Response = &cls.DescribeRemoteWriteTasksResponseParams{
			Infos:      []*cls.RemoteWriteInfo{},
			TotalCount: ptrUint64RemoteWrite(0),
			RequestId:  ptrStrRemoteWrite("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForRemoteWriteTask()
	res := localcls.ResourceTencentCloudClsRemoteWriteTask()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("task-not-exist#topic-not-exist")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestRemoteWriteTask_Update_Basic tests Update calls ModifyRemoteWriteTask
func TestRemoteWriteTask_Update_Basic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &cls.Client{}
	patches.ApplyMethodReturn(newMockMetaForRemoteWriteTask().client, "UseClsClient", clsClient)

	var capturedModifyRequest *cls.ModifyRemoteWriteTaskRequest
	patches.ApplyMethodFunc(clsClient, "ModifyRemoteWriteTaskWithContext", func(_ context.Context, request *cls.ModifyRemoteWriteTaskRequest) (*cls.ModifyRemoteWriteTaskResponse, error) {
		capturedModifyRequest = request
		resp := cls.NewModifyRemoteWriteTaskResponse()
		resp.Response = &cls.ModifyRemoteWriteTaskResponseParams{
			RequestId: ptrStrRemoteWrite("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(clsClient, "DescribeRemoteWriteTasksWithContext", func(_ context.Context, request *cls.DescribeRemoteWriteTasksRequest) (*cls.DescribeRemoteWriteTasksResponse, error) {
		resp := cls.NewDescribeRemoteWriteTasksResponse()
		resp.Response = &cls.DescribeRemoteWriteTasksResponseParams{
			Infos: []*cls.RemoteWriteInfo{
				{
					TaskId:         ptrStrRemoteWrite("task-test-update"),
					TopicId:        ptrStrRemoteWrite("topic-test-123"),
					Name:           ptrStrRemoteWrite("updated-name"),
					Target:         ptrStrRemoteWrite("prometheus"),
					RemoteWriteURL: ptrStrRemoteWrite("http://prometheus.example.com/api/v1/write"),
					AuthType:       ptrUint64RemoteWrite(0),
					NetType:        ptrUint64RemoteWrite(1),
					Enable:         ptrUint64RemoteWrite(0),
				},
			},
			TotalCount: ptrUint64RemoteWrite(1),
			RequestId:  ptrStrRemoteWrite("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForRemoteWriteTask()
	res := localcls.ResourceTencentCloudClsRemoteWriteTask()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"topic_id":         "topic-test-123",
		"name":             "updated-name",
		"target":           "prometheus",
		"remote_write_url": "http://prometheus.example.com/api/v1/write",
		"auth_type":        0,
		"net_type":         1,
		"enable":           0,
	})
	d.SetId("task-test-update#topic-test-123")

	// Simulate HasChange by setting old value first
	_ = d.Set("enable", 1)

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.NotNil(t, capturedModifyRequest)
	assert.Equal(t, "task-test-update", *capturedModifyRequest.TaskId)
	assert.Equal(t, "topic-test-123", *capturedModifyRequest.TopicId)
}

// TestRemoteWriteTask_Delete_Basic tests Delete calls DeleteRemoteWriteTask
func TestRemoteWriteTask_Delete_Basic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &cls.Client{}
	patches.ApplyMethodReturn(newMockMetaForRemoteWriteTask().client, "UseClsClient", clsClient)

	var capturedDeleteRequest *cls.DeleteRemoteWriteTaskRequest
	patches.ApplyMethodFunc(clsClient, "DeleteRemoteWriteTaskWithContext", func(_ context.Context, request *cls.DeleteRemoteWriteTaskRequest) (*cls.DeleteRemoteWriteTaskResponse, error) {
		capturedDeleteRequest = request
		resp := cls.NewDeleteRemoteWriteTaskResponse()
		resp.Response = &cls.DeleteRemoteWriteTaskResponseParams{
			RequestId: ptrStrRemoteWrite("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForRemoteWriteTask()
	res := localcls.ResourceTencentCloudClsRemoteWriteTask()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"topic_id": "topic-test-123",
	})
	d.SetId("task-test-delete#topic-test-123")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
	assert.NotNil(t, capturedDeleteRequest)
	assert.Equal(t, "task-test-delete", *capturedDeleteRequest.TaskId)
	assert.Equal(t, "topic-test-123", *capturedDeleteRequest.TopicId)
}

// TestRemoteWriteTask_Schema tests the schema definition
func TestRemoteWriteTask_Schema(t *testing.T) {
	res := localcls.ResourceTencentCloudClsRemoteWriteTask()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "topic_id")
	assert.Contains(t, res.Schema, "name")
	assert.Contains(t, res.Schema, "target")
	assert.Contains(t, res.Schema, "remote_write_url")
	assert.Contains(t, res.Schema, "auth_type")
	assert.Contains(t, res.Schema, "net_type")
	assert.Contains(t, res.Schema, "vpc_id")
	assert.Contains(t, res.Schema, "virtual_gateway_type")
	assert.Contains(t, res.Schema, "enable")
	assert.Contains(t, res.Schema, "auth_info")
	assert.Contains(t, res.Schema, "task_id")
	assert.Contains(t, res.Schema, "status")
	assert.Contains(t, res.Schema, "create_time")
	assert.Contains(t, res.Schema, "update_time")
	assert.Contains(t, res.Schema, "logset_id")

	assert.True(t, res.Schema["topic_id"].Required)
	assert.True(t, res.Schema["name"].Required)
	assert.True(t, res.Schema["target"].Required)
	assert.True(t, res.Schema["remote_write_url"].Required)
	assert.True(t, res.Schema["auth_type"].Required)
	assert.True(t, res.Schema["net_type"].Required)

	assert.True(t, res.Schema["vpc_id"].Optional)
	assert.True(t, res.Schema["virtual_gateway_type"].Optional)
	assert.True(t, res.Schema["enable"].Optional)
	assert.True(t, res.Schema["auth_info"].Optional)

	assert.True(t, res.Schema["task_id"].Computed)
	assert.True(t, res.Schema["status"].Computed)
	assert.True(t, res.Schema["create_time"].Computed)
	assert.True(t, res.Schema["update_time"].Computed)
	assert.True(t, res.Schema["logset_id"].Computed)
}
