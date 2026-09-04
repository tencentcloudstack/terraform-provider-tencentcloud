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

type mockMetaForCatProbeTasksDS struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForCatProbeTasksDS) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForCatProbeTasksDS{}

func newMockMetaForCatProbeTasksDS() *mockMetaForCatProbeTasksDS {
	return &mockMetaForCatProbeTasksDS{client: &connectivity.TencentCloudClient{}}
}

func ptrStringCPTDS(s string) *string { return &s }

func ptrInt64CPTDS(v int64) *int64 { return &v }

func buildProbeTask(name, taskId, targetAddress, parameters, createdAt, cron string, taskType, nodeIpType, interval, status, payMode, orderState, taskCategory, cronState, subSyncFlag int64) *cat.ProbeTask {
	return &cat.ProbeTask{
		Name:          ptrStringCPTDS(name),
		TaskId:        ptrStringCPTDS(taskId),
		TaskType:      ptrInt64CPTDS(taskType),
		Nodes:         []*string{ptrStringCPTDS("node-1"), ptrStringCPTDS("node-2")},
		NodeIpType:    ptrInt64CPTDS(nodeIpType),
		Interval:      ptrInt64CPTDS(interval),
		Parameters:    ptrStringCPTDS(parameters),
		Status:        ptrInt64CPTDS(status),
		TargetAddress: ptrStringCPTDS(targetAddress),
		PayMode:       ptrInt64CPTDS(payMode),
		OrderState:    ptrInt64CPTDS(orderState),
		TaskCategory:  ptrInt64CPTDS(taskCategory),
		CreatedAt:     ptrStringCPTDS(createdAt),
		Cron:          ptrStringCPTDS(cron),
		CronState:     ptrInt64CPTDS(cronState),
		TagInfoList: []*cat.KeyValuePair{
			{
				Key:   ptrStringCPTDS("Environment"),
				Value: ptrStringCPTDS("Production"),
			},
		},
		SubSyncFlag: ptrInt64CPTDS(subSyncFlag),
	}
}

// go test ./tencentcloud/services/cat/ -run "TestCatProbeTasksDS" -v -count=1 -gcflags="all=-l"

func TestCatProbeTasksDS_ReadBasic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	catClient := &cat.Client{}
	patches.ApplyMethodReturn(newMockMetaForCatProbeTasksDS().client, "UseCatClient", catClient)

	patches.ApplyMethodFunc(catClient, "DescribeProbeTasksWithContext", func(ctx context.Context, request *cat.DescribeProbeTasksRequest) (*cat.DescribeProbeTasksResponse, error) {
		resp := cat.NewDescribeProbeTasksResponse()
		resp.Response = &cat.DescribeProbeTasksResponseParams{
			TaskSet: []*cat.ProbeTask{
				buildProbeTask("task-001", "probe-task-001", "http://www.example.com", "params", "2024-01-01 00:00:00", "0 0 * * *", 1, 1, 5, 2, 2, 1, 1, 1, 0),
				buildProbeTask("task-002", "probe-task-002", "http://www.test.com", "params2", "2024-01-02 00:00:00", "0 0 * * *", 5, 0, 10, 6, 1, 1, 2, 2, 1),
			},
			Total:     ptrInt64CPTDS(2),
			RequestId: ptrStringCPTDS("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForCatProbeTasksDS()
	res := svccat.DataSourceTencentCloudCatProbeTasks()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	taskSet := d.Get("task_set").([]interface{})
	assert.Len(t, taskSet, 2)

	task0 := taskSet[0].(map[string]interface{})
	assert.Equal(t, "task-001", task0["name"].(string))
	assert.Equal(t, "probe-task-001", task0["task_id"].(string))
	assert.Equal(t, 1, task0["task_type"].(int))
	assert.Equal(t, "http://www.example.com", task0["target_address"].(string))
	assert.Equal(t, 2, task0["status"].(int))
	assert.Equal(t, 2, task0["pay_mode"].(int))
	assert.Equal(t, "params", task0["parameters"].(string))

	nodes := task0["nodes"].([]interface{})
	assert.Len(t, nodes, 2)
	assert.Equal(t, "node-1", nodes[0].(string))

	tagInfoList := task0["tag_info_list"].([]interface{})
	assert.Len(t, tagInfoList, 1)
	tagInfo := tagInfoList[0].(map[string]interface{})
	assert.Equal(t, "Environment", tagInfo["key"].(string))
	assert.Equal(t, "Production", tagInfo["value"].(string))

	total := d.Get("total").(int)
	assert.Equal(t, 2, total)
}

func TestCatProbeTasksDS_ReadWithFilters(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	catClient := &cat.Client{}
	patches.ApplyMethodReturn(newMockMetaForCatProbeTasksDS().client, "UseCatClient", catClient)

	patches.ApplyMethodFunc(catClient, "DescribeProbeTasksWithContext", func(ctx context.Context, request *cat.DescribeProbeTasksRequest) (*cat.DescribeProbeTasksResponse, error) {
		assert.Equal(t, "my-probe-task", *request.TaskName)
		assert.Equal(t, "http://www.example.com", *request.TargetAddress)
		assert.Equal(t, int64(2), *request.PayMode)
		assert.Equal(t, 1, len(request.TaskStatus))
		assert.Equal(t, int64(2), *request.TaskStatus[0])
		assert.Equal(t, 1, len(request.TagFilters))
		assert.Equal(t, "Environment", *request.TagFilters[0].Key)
		assert.Equal(t, "Production", *request.TagFilters[0].Value)

		resp := cat.NewDescribeProbeTasksResponse()
		resp.Response = &cat.DescribeProbeTasksResponseParams{
			TaskSet: []*cat.ProbeTask{
				buildProbeTask("my-probe-task", "probe-task-001", "http://www.example.com", "params", "2024-01-01 00:00:00", "", 1, 1, 5, 2, 2, 1, 1, 0, 0),
			},
			Total:     ptrInt64CPTDS(1),
			RequestId: ptrStringCPTDS("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForCatProbeTasksDS()
	res := svccat.DataSourceTencentCloudCatProbeTasks()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"task_name":      "my-probe-task",
		"target_address": "http://www.example.com",
		"pay_mode":       2,
		"task_status":    []interface{}{2},
		"tag_filters": []interface{}{
			map[string]interface{}{
				"key":   "Environment",
				"value": "Production",
			},
		},
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	taskSet := d.Get("task_set").([]interface{})
	assert.Len(t, taskSet, 1)

	total := d.Get("total").(int)
	assert.Equal(t, 1, total)
}

func TestCatProbeTasksDS_ReadEmpty(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	catClient := &cat.Client{}
	patches.ApplyMethodReturn(newMockMetaForCatProbeTasksDS().client, "UseCatClient", catClient)

	patches.ApplyMethodFunc(catClient, "DescribeProbeTasksWithContext", func(ctx context.Context, request *cat.DescribeProbeTasksRequest) (*cat.DescribeProbeTasksResponse, error) {
		resp := cat.NewDescribeProbeTasksResponse()
		resp.Response = &cat.DescribeProbeTasksResponseParams{
			TaskSet:   []*cat.ProbeTask{},
			Total:     ptrInt64CPTDS(0),
			RequestId: ptrStringCPTDS("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForCatProbeTasksDS()
	res := svccat.DataSourceTencentCloudCatProbeTasks()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	taskSet := d.Get("task_set").([]interface{})
	assert.Len(t, taskSet, 0)

	total := d.Get("total").(int)
	assert.Equal(t, 0, total)
}

func TestCatProbeTasksDS_ReadError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	catClient := &cat.Client{}
	patches.ApplyMethodReturn(newMockMetaForCatProbeTasksDS().client, "UseCatClient", catClient)

	patches.ApplyMethodFunc(catClient, "DescribeProbeTasksWithContext", func(ctx context.Context, request *cat.DescribeProbeTasksRequest) (*cat.DescribeProbeTasksResponse, error) {
		return nil, assert.AnError
	})

	meta := newMockMetaForCatProbeTasksDS()
	res := svccat.DataSourceTencentCloudCatProbeTasks()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	err := res.Read(d, meta)
	assert.Error(t, err)
}

func TestCatProbeTasksDS_Schema(t *testing.T) {
	res := svccat.DataSourceTencentCloudCatProbeTasks()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "task_ids")
	assert.Contains(t, res.Schema, "task_name")
	assert.Contains(t, res.Schema, "target_address")
	assert.Contains(t, res.Schema, "task_status")
	assert.Contains(t, res.Schema, "pay_mode")
	assert.Contains(t, res.Schema, "order_state")
	assert.Contains(t, res.Schema, "task_type")
	assert.Contains(t, res.Schema, "task_category")
	assert.Contains(t, res.Schema, "order_by")
	assert.Contains(t, res.Schema, "ascend")
	assert.Contains(t, res.Schema, "tag_filters")
	assert.Contains(t, res.Schema, "task_set")
	assert.Contains(t, res.Schema, "total")
	assert.Contains(t, res.Schema, "result_output_file")

	taskIDsSchema := res.Schema["task_ids"]
	assert.Equal(t, schema.TypeList, taskIDsSchema.Type)
	assert.True(t, taskIDsSchema.Optional)

	taskNameSchema := res.Schema["task_name"]
	assert.Equal(t, schema.TypeString, taskNameSchema.Type)
	assert.True(t, taskNameSchema.Optional)

	taskStatusSchema := res.Schema["task_status"]
	assert.Equal(t, schema.TypeList, taskStatusSchema.Type)
	assert.True(t, taskStatusSchema.Optional)

	payModeSchema := res.Schema["pay_mode"]
	assert.Equal(t, schema.TypeInt, payModeSchema.Type)
	assert.True(t, payModeSchema.Optional)

	ascendSchema := res.Schema["ascend"]
	assert.Equal(t, schema.TypeBool, ascendSchema.Type)
	assert.True(t, ascendSchema.Optional)

	tagFiltersSchema := res.Schema["tag_filters"]
	assert.Equal(t, schema.TypeList, tagFiltersSchema.Type)
	assert.True(t, tagFiltersSchema.Optional)
	tagFiltersElemRes := tagFiltersSchema.Elem.(*schema.Resource)
	assert.Contains(t, tagFiltersElemRes.Schema, "key")
	assert.Contains(t, tagFiltersElemRes.Schema, "value")
	assert.True(t, tagFiltersElemRes.Schema["key"].Required)
	assert.True(t, tagFiltersElemRes.Schema["value"].Required)

	taskSetSchema := res.Schema["task_set"]
	assert.Equal(t, schema.TypeList, taskSetSchema.Type)
	assert.True(t, taskSetSchema.Computed)
	taskSetElemRes := taskSetSchema.Elem.(*schema.Resource)
	assert.Contains(t, taskSetElemRes.Schema, "name")
	assert.Contains(t, taskSetElemRes.Schema, "task_id")
	assert.Contains(t, taskSetElemRes.Schema, "task_type")
	assert.Contains(t, taskSetElemRes.Schema, "nodes")
	assert.Contains(t, taskSetElemRes.Schema, "node_ip_type")
	assert.Contains(t, taskSetElemRes.Schema, "interval")
	assert.Contains(t, taskSetElemRes.Schema, "parameters")
	assert.Contains(t, taskSetElemRes.Schema, "status")
	assert.Contains(t, taskSetElemRes.Schema, "target_address")
	assert.Contains(t, taskSetElemRes.Schema, "pay_mode")
	assert.Contains(t, taskSetElemRes.Schema, "order_state")
	assert.Contains(t, taskSetElemRes.Schema, "task_category")
	assert.Contains(t, taskSetElemRes.Schema, "created_at")
	assert.Contains(t, taskSetElemRes.Schema, "cron")
	assert.Contains(t, taskSetElemRes.Schema, "cron_state")
	assert.Contains(t, taskSetElemRes.Schema, "tag_info_list")
	assert.Contains(t, taskSetElemRes.Schema, "sub_sync_flag")

	totalSchema := res.Schema["total"]
	assert.Equal(t, schema.TypeInt, totalSchema.Type)
	assert.True(t, totalSchema.Computed)

	outputSchema := res.Schema["result_output_file"]
	assert.Equal(t, schema.TypeString, outputSchema.Type)
	assert.True(t, outputSchema.Optional)
}
