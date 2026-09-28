package cls_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	clsv20201016 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cls/v20201016"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cls"
)

type mockMetaForClsClsDeliverTask struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForClsClsDeliverTask) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForClsClsDeliverTask{}

func newMockMetaForClsClsDeliverTask() *mockMetaForClsClsDeliverTask {
	return &mockMetaForClsClsDeliverTask{client: &connectivity.TencentCloudClient{}}
}

func ptrStringCDT(s string) *string { return &s }
func ptrUint64CDT(v uint64) *uint64 { return &v }

func buildClsDeliverTaskInfo(taskId string) *clsv20201016.CLSDeliverTaskInfo {
	return &clsv20201016.CLSDeliverTaskInfo{
		TaskId:   ptrStringCDT(taskId),
		TaskName: ptrStringCDT("tf-example-deliver-task"),
		Uin:      ptrUint64CDT(1000000000),
		SourceTopicConfig: &clsv20201016.SourceTopicConfig{
			TopicFilterType: ptrUint64CDT(1),
			LogsetId:        ptrStringCDT("source-logset-id"),
			Topics: []*clsv20201016.SourceTopicInfo{
				{
					TopicId: ptrStringCDT("source-topic-id"),
				},
			},
		},
		TargetTopicConfig: &clsv20201016.TargetTopicConfig{
			AccountType: ptrUint64CDT(1),
			Region:      ptrStringCDT("ap-guangzhou"),
			LogsetId:    ptrStringCDT("target-logset-id"),
			TopicId:     ptrStringCDT("target-topic-id"),
		},
		DeliverRule: &clsv20201016.DeliverRule{
			DataScope: ptrUint64CDT(3),
		},
		Compliance:     ptrUint64CDT(1),
		Status:         ptrUint64CDT(0),
		Enable:         ptrUint64CDT(0),
		Progress:       ptrUint64CDT(100),
		HasServicesLog: ptrUint64CDT(2),
		CreateTime:     ptrUint64CDT(1696000000),
		UpdateTime:     ptrUint64CDT(1696000000),
	}
}

func buildClsDeliverTaskReadResp(taskId string) *clsv20201016.DescribeCLSDeliverTasksResponse {
	resp := clsv20201016.NewDescribeCLSDeliverTasksResponse()
	resp.Response = &clsv20201016.DescribeCLSDeliverTasksResponseParams{
		Infos: []*clsv20201016.CLSDeliverTaskInfo{buildClsDeliverTaskInfo(taskId)},
		Total: ptrUint64CDT(1),
	}
	return resp
}

func buildClsDeliverTaskEmptyResp() *clsv20201016.DescribeCLSDeliverTasksResponse {
	resp := clsv20201016.NewDescribeCLSDeliverTasksResponse()
	resp.Response = &clsv20201016.DescribeCLSDeliverTasksResponseParams{
		Infos: []*clsv20201016.CLSDeliverTaskInfo{},
		Total: ptrUint64CDT(0),
	}
	return resp
}

// TestClsDeliverTask_Create_Basic verifies that Create sets the id and the
// create request carries the schema fields.
func TestClsDeliverTask_Create_Basic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newMockMetaForClsClsDeliverTask().client, "UseClsClient", clsClient)

	var capturedRequest *clsv20201016.CreateCLSDeliverTaskRequest
	patches.ApplyMethodFunc(clsClient, "CreateCLSDeliverTaskWithContext", func(_ context.Context, request *clsv20201016.CreateCLSDeliverTaskRequest) (*clsv20201016.CreateCLSDeliverTaskResponse, error) {
		capturedRequest = request
		resp := clsv20201016.NewCreateCLSDeliverTaskResponse()
		resp.Response = &clsv20201016.CreateCLSDeliverTaskResponseParams{
			TaskId: ptrStringCDT("deliver-task-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(clsClient, "DescribeCLSDeliverTasksWithContext", func(_ context.Context, request *clsv20201016.DescribeCLSDeliverTasksRequest) (*clsv20201016.DescribeCLSDeliverTasksResponse, error) {
		return buildClsDeliverTaskReadResp("deliver-task-id"), nil
	})

	meta := newMockMetaForClsClsDeliverTask()
	res := cls.ResourceTencentCloudClsClsDeliverTask()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"task_name":  "tf-example-deliver-task",
		"compliance": 1,
		"source_topic_config": []interface{}{
			map[string]interface{}{
				"topic_filter_type": 1,
				"logset_id":         "source-logset-id",
				"topics": []interface{}{
					map[string]interface{}{
						"topic_id": "source-topic-id",
					},
				},
			},
		},
		"target_topic_config": []interface{}{
			map[string]interface{}{
				"account_type": 1,
				"region":       "ap-guangzhou",
				"logset_id":    "target-logset-id",
				"topic_id":     "target-topic-id",
			},
		},
		"deliver_rule": []interface{}{
			map[string]interface{}{
				"data_scope": 3,
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "deliver-task-id", d.Id())

	// Verify the request carries the schema fields.
	assert.NotNil(t, capturedRequest)
	assert.Equal(t, "tf-example-deliver-task", *capturedRequest.TaskName)
	assert.NotNil(t, capturedRequest.SourceTopicConfig)
	assert.Equal(t, "source-logset-id", *capturedRequest.SourceTopicConfig.LogsetId)
	assert.NotNil(t, capturedRequest.TargetTopicConfig)
	assert.Equal(t, "ap-guangzhou", *capturedRequest.TargetTopicConfig.Region)
	assert.NotNil(t, capturedRequest.DeliverRule)
	assert.Equal(t, uint64(3), *capturedRequest.DeliverRule.DataScope)
	assert.Equal(t, uint64(1), *capturedRequest.Compliance)
}

// TestClsDeliverTask_Read_Basic verifies that Read populates the state from
// the API response.
func TestClsDeliverTask_Read_Basic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newMockMetaForClsClsDeliverTask().client, "UseClsClient", clsClient)

	patches.ApplyMethodFunc(clsClient, "DescribeCLSDeliverTasksWithContext", func(_ context.Context, request *clsv20201016.DescribeCLSDeliverTasksRequest) (*clsv20201016.DescribeCLSDeliverTasksResponse, error) {
		return buildClsDeliverTaskReadResp("deliver-task-id"), nil
	})

	meta := newMockMetaForClsClsDeliverTask()
	res := cls.ResourceTencentCloudClsClsDeliverTask()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("deliver-task-id")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "deliver-task-id", d.Id())
	assert.Equal(t, "tf-example-deliver-task", d.Get("task_name").(string))
	assert.Equal(t, 1, d.Get("compliance").(int))
}

// TestClsDeliverTask_Read_EmptyResult verifies that when the API returns an
// empty list, the resource id is cleared.
func TestClsDeliverTask_Read_EmptyResult(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newMockMetaForClsClsDeliverTask().client, "UseClsClient", clsClient)

	patches.ApplyMethodFunc(clsClient, "DescribeCLSDeliverTasksWithContext", func(_ context.Context, request *clsv20201016.DescribeCLSDeliverTasksRequest) (*clsv20201016.DescribeCLSDeliverTasksResponse, error) {
		return buildClsDeliverTaskEmptyResp(), nil
	})

	meta := newMockMetaForClsClsDeliverTask()
	res := cls.ResourceTencentCloudClsClsDeliverTask()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("delivered-task-id")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Empty(t, d.Id())
}

// TestClsDeliverTask_Update_Basic verifies that when mutable fields change,
// ModifyCLSDeliverTask is called with the task id and changed fields.
func TestClsDeliverTask_Update_Basic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newMockMetaForClsClsDeliverTask().client, "UseClsClient", clsClient)

	var capturedRequest *clsv20201016.ModifyCLSDeliverTaskRequest
	patches.ApplyMethodFunc(clsClient, "ModifyCLSDeliverTaskWithContext", func(_ context.Context, request *clsv20201016.ModifyCLSDeliverTaskRequest) (*clsv20201016.ModifyCLSDeliverTaskResponse, error) {
		capturedRequest = request
		resp := clsv20201016.NewModifyCLSDeliverTaskResponse()
		resp.Response = &clsv20201016.ModifyCLSDeliverTaskResponseParams{}
		return resp, nil
	})

	patches.ApplyMethodFunc(clsClient, "DescribeCLSDeliverTasksWithContext", func(_ context.Context, request *clsv20201016.DescribeCLSDeliverTasksRequest) (*clsv20201016.DescribeCLSDeliverTasksResponse, error) {
		return buildClsDeliverTaskReadResp("deliver-task-id"), nil
	})

	meta := newMockMetaForClsClsDeliverTask()
	res := cls.ResourceTencentCloudClsClsDeliverTask()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"task_name":  "tf-example-deliver-task-updated",
		"compliance": 1,
		"source_topic_config": []interface{}{
			map[string]interface{}{
				"topic_filter_type": 1,
				"logset_id":         "source-logset-id",
				"topics": []interface{}{
					map[string]interface{}{
						"topic_id": "source-topic-id",
					},
				},
			},
		},
		"target_topic_config": []interface{}{
			map[string]interface{}{
				"account_type": 1,
				"region":       "ap-guangzhou",
				"logset_id":    "target-logset-id",
				"topic_id":     "target-topic-id",
			},
		},
		"deliver_rule": []interface{}{
			map[string]interface{}{
				"data_scope": 3,
			},
		},
		"enable": 0,
	})
	d.SetId("deliver-task-id")

	// Force only `task_name` to be detected as changed so the modify path runs.
	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return key == "task_name"
	})

	err := res.Update(d, meta)
	assert.NoError(t, err)

	assert.NotNil(t, capturedRequest)
	assert.Equal(t, "deliver-task-id", *capturedRequest.TaskId)
	assert.Equal(t, "tf-example-deliver-task-updated", *capturedRequest.TaskName)
	// compliance should NOT be sent in update (no Compliance field on ModifyCLSDeliverTaskRequest struct).
	assert.Equal(t, "deliver-task-id", d.Id())
}

// TestClsDeliverTask_Delete_Basic verifies that Delete calls
// DeleteCLSDeliverTask with the task id.
func TestClsDeliverTask_Delete_Basic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newMockMetaForClsClsDeliverTask().client, "UseClsClient", clsClient)

	var capturedRequest *clsv20201016.DeleteCLSDeliverTaskRequest
	patches.ApplyMethodFunc(clsClient, "DeleteCLSDeliverTaskWithContext", func(_ context.Context, request *clsv20201016.DeleteCLSDeliverTaskRequest) (*clsv20201016.DeleteCLSDeliverTaskResponse, error) {
		capturedRequest = request
		resp := clsv20201016.NewDeleteCLSDeliverTaskResponse()
		resp.Response = &clsv20201016.DeleteCLSDeliverTaskResponseParams{}
		return resp, nil
	})

	meta := newMockMetaForClsClsDeliverTask()
	res := cls.ResourceTencentCloudClsClsDeliverTask()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})
	d.SetId("deliver-task-id")

	err := res.Delete(d, meta)
	assert.NoError(t, err)

	assert.NotNil(t, capturedRequest)
	assert.Equal(t, "deliver-task-id", *capturedRequest.TaskId)
}
