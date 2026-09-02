package cat_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	cat "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/cat/v20180409"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	svccat "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cat"
)

type mockMetaForCatInstantTasksDS struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForCatInstantTasksDS) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForCatInstantTasksDS{}

func newMockMetaForCatInstantTasksDS() *mockMetaForCatInstantTasksDS {
	return &mockMetaForCatInstantTasksDS{client: &connectivity.TencentCloudClient{}}
}

func ptrStringCITDS(s string) *string { return &s }

func ptrUint64CITDS(v uint64) *uint64 { return &v }

func ptrFloat64CITDS(v float64) *float64 { return &v }

func buildSingleInstantTask(taskId, targetAddress, status string, taskType, probeTime, nodeCount, taskCategory uint64, successRate float64) *cat.SingleInstantTask {
	return &cat.SingleInstantTask{
		TaskId:        ptrStringCITDS(taskId),
		TargetAddress: ptrStringCITDS(targetAddress),
		TaskType:      ptrUint64CITDS(taskType),
		ProbeTime:     ptrUint64CITDS(probeTime),
		Status:        ptrStringCITDS(status),
		SuccessRate:   ptrFloat64CITDS(successRate),
		NodeCount:     ptrUint64CITDS(nodeCount),
		TaskCategory:  ptrUint64CITDS(taskCategory),
	}
}

func TestCatInstantTasksDS_ReadBasic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	catClient := &cat.Client{}
	patches.ApplyMethodReturn(newMockMetaForCatInstantTasksDS().client, "UseCatClient", catClient)

	patches.ApplyMethodFunc(catClient, "DescribeInstantTasksWithContext", func(ctx context.Context, request *cat.DescribeInstantTasksRequest) (*cat.DescribeInstantTasksResponse, error) {
		resp := cat.NewDescribeInstantTasksResponse()
		resp.Response = &cat.DescribeInstantTasksResponseParams{
			Tasks: []*cat.SingleInstantTask{
				buildSingleInstantTask("task-001", "http://www.example.com", "success", 1, 1667923200, 3, 1, 99.5),
				buildSingleInstantTask("task-002", "http://www.test.com", "failed", 2, 1667923300, 5, 2, 80.0),
			},
			Total:     ptrUint64CITDS(2),
			RequestId: ptrStringCITDS("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForCatInstantTasksDS()
	res := svccat.DataSourceTencentCloudCatInstantTasks()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	tasks := d.Get("tasks").([]interface{})
	assert.Len(t, tasks, 2)

	task0 := tasks[0].(map[string]interface{})
	assert.Equal(t, "task-001", task0["task_id"].(string))
	assert.Equal(t, "http://www.example.com", task0["target_address"].(string))
	assert.Equal(t, 1, task0["task_type"].(int))
	assert.Equal(t, 1667923200, task0["probe_time"].(int))
	assert.Equal(t, "success", task0["status"].(string))
	assert.Equal(t, 99.5, task0["success_rate"].(float64))
	assert.Equal(t, 3, task0["node_count"].(int))
	assert.Equal(t, 1, task0["task_category"].(int))

	task1 := tasks[1].(map[string]interface{})
	assert.Equal(t, "task-002", task1["task_id"].(string))
	assert.Equal(t, "http://www.test.com", task1["target_address"].(string))

	total := d.Get("total").(int)
	assert.Equal(t, 2, total)
}

func TestCatInstantTasksDS_ReadEmpty(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	catClient := &cat.Client{}
	patches.ApplyMethodReturn(newMockMetaForCatInstantTasksDS().client, "UseCatClient", catClient)

	patches.ApplyMethodFunc(catClient, "DescribeInstantTasksWithContext", func(ctx context.Context, request *cat.DescribeInstantTasksRequest) (*cat.DescribeInstantTasksResponse, error) {
		resp := cat.NewDescribeInstantTasksResponse()
		resp.Response = &cat.DescribeInstantTasksResponseParams{
			Tasks:     []*cat.SingleInstantTask{},
			Total:     ptrUint64CITDS(0),
			RequestId: ptrStringCITDS("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForCatInstantTasksDS()
	res := svccat.DataSourceTencentCloudCatInstantTasks()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	err := res.Read(d, meta)
	assert.NoError(t, err)

	tasks := d.Get("tasks").([]interface{})
	assert.Len(t, tasks, 0)

	total := d.Get("total").(int)
	assert.Equal(t, 0, total)
}

func TestCatInstantTasksDS_ReadError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	catClient := &cat.Client{}
	patches.ApplyMethodReturn(newMockMetaForCatInstantTasksDS().client, "UseCatClient", catClient)

	patches.ApplyMethodFunc(catClient, "DescribeInstantTasksWithContext", func(ctx context.Context, request *cat.DescribeInstantTasksRequest) (*cat.DescribeInstantTasksResponse, error) {
		return nil, assert.AnError
	})

	meta := newMockMetaForCatInstantTasksDS()
	res := svccat.DataSourceTencentCloudCatInstantTasks()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	err := res.Read(d, meta)
	assert.Error(t, err)
}

func TestCatInstantTasksDS_Schema(t *testing.T) {
	res := svccat.DataSourceTencentCloudCatInstantTasks()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "tasks")
	assert.Contains(t, res.Schema, "total")
	assert.Contains(t, res.Schema, "result_output_file")

	tasksSchema := res.Schema["tasks"]
	assert.Equal(t, schema.TypeList, tasksSchema.Type)
	assert.True(t, tasksSchema.Computed)

	elemRes := tasksSchema.Elem.(*schema.Resource)
	assert.Contains(t, elemRes.Schema, "task_id")
	assert.Contains(t, elemRes.Schema, "target_address")
	assert.Contains(t, elemRes.Schema, "task_type")
	assert.Contains(t, elemRes.Schema, "probe_time")
	assert.Contains(t, elemRes.Schema, "status")
	assert.Contains(t, elemRes.Schema, "success_rate")
	assert.Contains(t, elemRes.Schema, "node_count")
	assert.Contains(t, elemRes.Schema, "task_category")

	totalSchema := res.Schema["total"]
	assert.Equal(t, schema.TypeInt, totalSchema.Type)
	assert.True(t, totalSchema.Computed)

	outputSchema := res.Schema["result_output_file"]
	assert.Equal(t, schema.TypeString, outputSchema.Type)
	assert.True(t, outputSchema.Optional)
}
