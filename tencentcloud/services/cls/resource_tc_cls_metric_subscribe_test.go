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
	localcls "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/cls"
)

type mockMetaForClsMetricSubscribe struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForClsMetricSubscribe) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForClsMetricSubscribe{}

func newMockMetaForClsMetricSubscribe() *mockMetaForClsMetricSubscribe {
	return &mockMetaForClsMetricSubscribe{client: &connectivity.TencentCloudClient{}}
}

func ptrStringMS(s string) *string { return &s }
func ptrUint64MS(v uint64) *uint64 { return &v }
func ptrInt64MS(v int64) *int64    { return &v }

func buildMetricSubscribeInfo() *clsv20201016.MetricSubscribeInfo {
	return &clsv20201016.MetricSubscribeInfo{
		TaskId:    ptrStringMS("task-test-123"),
		TopicId:   ptrStringMS("topic-test-123"),
		Name:      ptrStringMS("tf-example-metric-subscribe"),
		Namespace: ptrStringMS("QCE/CVM"),
		Metrics: []*clsv20201016.MetricConfig{
			{
				MetricName: ptrStringMS("cpu_usage"),
				Periods:    []*uint64{ptrUint64MS(60), ptrUint64MS(300)},
				MetricLabels: []*clsv20201016.MetricLabel{
					{
						Key:   ptrStringMS("label_key"),
						Value: ptrStringMS("label_value"),
					},
				},
			},
		},
		InstanceInfo: &clsv20201016.InstanceConfig{
			InstanceDimension: []*string{ptrStringMS("InstanceId")},
			Instances: []*clsv20201016.Instance{
				{
					Values: []*string{ptrStringMS("ins-xxxxxxxx")},
				},
			},
		},
		Enable:     ptrUint64MS(2),
		Status:     ptrUint64MS(2),
		CreateTime: ptrInt64MS(1700000000),
		UpdateTime: ptrInt64MS(1700000100),
	}
}

func metricSubscribeRawInput() map[string]interface{} {
	return map[string]interface{}{
		"name":      "tf-example-metric-subscribe",
		"topic_id":  "topic-test-123",
		"namespace": "QCE/CVM",
		"enable":    2,
		"metrics": []interface{}{
			map[string]interface{}{
				"metric_name": "cpu_usage",
				"periods":     []interface{}{60, 300},
				"metric_labels": []interface{}{
					map[string]interface{}{
						"key":   "label_key",
						"value": "label_value",
					},
				},
			},
		},
		"instance_info": []interface{}{
			map[string]interface{}{
				"instance_dimension": []interface{}{"InstanceId"},
				"instances": []interface{}{
					map[string]interface{}{
						"values": []interface{}{"ins-xxxxxxxx"},
					},
				},
			},
		},
	}
}

// TestClsMetricSubscribe_Create_Success tests the Create flow maps fields to the
// CreateMetricSubscribe request and sets the composite id topicId#taskId.
func TestClsMetricSubscribe_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newMockMetaForClsMetricSubscribe().client, "UseClsV20201016Client", clsClient)

	var capturedRequest *clsv20201016.CreateMetricSubscribeRequest
	patches.ApplyMethodFunc(clsClient, "CreateMetricSubscribeWithContext", func(_ context.Context, request *clsv20201016.CreateMetricSubscribeRequest) (*clsv20201016.CreateMetricSubscribeResponse, error) {
		capturedRequest = request
		resp := clsv20201016.NewCreateMetricSubscribeResponse()
		resp.Response = &clsv20201016.CreateMetricSubscribeResponseParams{
			TaskId:    ptrStringMS("task-test-123"),
			RequestId: ptrStringMS("fake-request-id"),
		}
		return resp, nil
	})

	info := buildMetricSubscribeInfo()
	patches.ApplyMethodFunc(&localcls.ClsService{}, "DescribeClsMetricSubscribeById", func(_ context.Context, _ string, _ string) (*clsv20201016.MetricSubscribeInfo, error) {
		return info, nil
	})

	meta := newMockMetaForClsMetricSubscribe()
	res := localcls.ResourceTencentCloudClsMetricSubscribe()
	d := schema.TestResourceDataRaw(t, res.Schema, metricSubscribeRawInput())

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "topic-test-123"+tccommon.FILED_SP+"task-test-123", d.Id())

	// Verify request mapping
	assert.NotNil(t, capturedRequest)
	assert.Equal(t, "tf-example-metric-subscribe", *capturedRequest.Name)
	assert.Equal(t, "topic-test-123", *capturedRequest.TopicId)
	assert.Equal(t, "QCE/CVM", *capturedRequest.Namespace)
	assert.NotNil(t, capturedRequest.Metrics)
	assert.Equal(t, 1, len(capturedRequest.Metrics))
	assert.Equal(t, "cpu_usage", *capturedRequest.Metrics[0].MetricName)
	assert.Equal(t, 2, len(capturedRequest.Metrics[0].Periods))
	assert.Equal(t, uint64(60), *capturedRequest.Metrics[0].Periods[0])
	assert.Equal(t, uint64(300), *capturedRequest.Metrics[0].Periods[1])
	assert.Equal(t, 1, len(capturedRequest.Metrics[0].MetricLabels))
	assert.Equal(t, "label_key", *capturedRequest.Metrics[0].MetricLabels[0].Key)
	assert.Equal(t, "label_value", *capturedRequest.Metrics[0].MetricLabels[0].Value)
	assert.NotNil(t, capturedRequest.InstanceInfo)
	assert.Equal(t, 1, len(capturedRequest.InstanceInfo.InstanceDimension))
	assert.Equal(t, "InstanceId", *capturedRequest.InstanceInfo.InstanceDimension[0])
	assert.Equal(t, 1, len(capturedRequest.InstanceInfo.Instances))
	assert.Equal(t, 1, len(capturedRequest.InstanceInfo.Instances[0].Values))
	assert.Equal(t, "ins-xxxxxxxx", *capturedRequest.InstanceInfo.Instances[0].Values[0])

	// Verify computed fields read back into state
	assert.Equal(t, "task-test-123", d.Get("task_id"))
	assert.Equal(t, 2, d.Get("status"))
	assert.Equal(t, 1700000000, d.Get("create_time"))
	assert.Equal(t, 1700000100, d.Get("update_time"))
}

// TestClsMetricSubscribe_Create_EmptyTaskId tests that an empty TaskId in the
// response returns an error (NonRetryableError) instead of writing an empty id.
func TestClsMetricSubscribe_Create_EmptyTaskId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newMockMetaForClsMetricSubscribe().client, "UseClsV20201016Client", clsClient)

	patches.ApplyMethodFunc(clsClient, "CreateMetricSubscribeWithContext", func(_ context.Context, request *clsv20201016.CreateMetricSubscribeRequest) (*clsv20201016.CreateMetricSubscribeResponse, error) {
		resp := clsv20201016.NewCreateMetricSubscribeResponse()
		resp.Response = &clsv20201016.CreateMetricSubscribeResponseParams{
			TaskId:    ptrStringMS(""),
			RequestId: ptrStringMS("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForClsMetricSubscribe()
	res := localcls.ResourceTencentCloudClsMetricSubscribe()
	d := schema.TestResourceDataRaw(t, res.Schema, metricSubscribeRawInput())

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Equal(t, "", d.Id())
}

// TestClsMetricSubscribe_Read_Success tests Read populates state from the
// DescribeClsMetricSubscribeById service method response.
func TestClsMetricSubscribe_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	info := buildMetricSubscribeInfo()
	patches.ApplyMethodFunc(&localcls.ClsService{}, "DescribeClsMetricSubscribeById", func(_ context.Context, _ string, _ string) (*clsv20201016.MetricSubscribeInfo, error) {
		return info, nil
	})

	meta := newMockMetaForClsMetricSubscribe()
	res := localcls.ResourceTencentCloudClsMetricSubscribe()
	d := schema.TestResourceDataRaw(t, res.Schema, metricSubscribeRawInput())
	d.SetId("topic-test-123" + tccommon.FILED_SP + "task-test-123")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "topic-test-123"+tccommon.FILED_SP+"task-test-123", d.Id())
	assert.Equal(t, "tf-example-metric-subscribe", d.Get("name"))
	assert.Equal(t, "topic-test-123", d.Get("topic_id"))
	assert.Equal(t, "QCE/CVM", d.Get("namespace"))
	assert.Equal(t, "task-test-123", d.Get("task_id"))
	assert.Equal(t, 2, d.Get("status"))
	assert.Equal(t, 1700000000, d.Get("create_time"))
	assert.Equal(t, 1700000100, d.Get("update_time"))
}

// TestClsMetricSubscribe_Read_NotFound tests that Read clears state (SetId(""))
// when the service method returns nil (resource not found).
func TestClsMetricSubscribe_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyMethodFunc(&localcls.ClsService{}, "DescribeClsMetricSubscribeById", func(_ context.Context, _ string, _ string) (*clsv20201016.MetricSubscribeInfo, error) {
		return nil, nil
	})

	meta := newMockMetaForClsMetricSubscribe()
	res := localcls.ResourceTencentCloudClsMetricSubscribe()
	d := schema.TestResourceDataRaw(t, res.Schema, metricSubscribeRawInput())
	d.SetId("topic-test-123" + tccommon.FILED_SP + "task-test-123")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestClsMetricSubscribe_Update_Success tests the Update flow builds a
// ModifyMetricSubscribe request with changed fields and calls Read.
func TestClsMetricSubscribe_Update_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newMockMetaForClsMetricSubscribe().client, "UseClsV20201016Client", clsClient)

	var capturedModifyRequest *clsv20201016.ModifyMetricSubscribeRequest
	patches.ApplyMethodFunc(clsClient, "ModifyMetricSubscribeWithContext", func(_ context.Context, request *clsv20201016.ModifyMetricSubscribeRequest) (*clsv20201016.ModifyMetricSubscribeResponse, error) {
		capturedModifyRequest = request
		resp := clsv20201016.NewModifyMetricSubscribeResponse()
		resp.Response = &clsv20201016.ModifyMetricSubscribeResponseParams{
			RequestId: ptrStringMS("fake-request-id"),
		}
		return resp, nil
	})

	info := buildMetricSubscribeInfo()
	patches.ApplyMethodFunc(&localcls.ClsService{}, "DescribeClsMetricSubscribeById", func(_ context.Context, _ string, _ string) (*clsv20201016.MetricSubscribeInfo, error) {
		return info, nil
	})

	meta := newMockMetaForClsMetricSubscribe()
	res := localcls.ResourceTencentCloudClsMetricSubscribe()
	d := schema.TestResourceDataRaw(t, res.Schema, metricSubscribeRawInput())
	d.SetId("topic-test-123" + tccommon.FILED_SP + "task-test-123")

	// Simulate a change on "name" to trigger the Modify call
	_ = d.Set("name", "tf-example-metric-subscribe-updated")

	err := res.Update(d, meta)
	assert.NoError(t, err)

	// Verify the modify request carried TopicId/TaskId and the updated name
	assert.NotNil(t, capturedModifyRequest)
	assert.Equal(t, "topic-test-123", *capturedModifyRequest.TopicId)
	assert.Equal(t, "task-test-123", *capturedModifyRequest.TaskId)
	assert.Equal(t, "tf-example-metric-subscribe-updated", *capturedModifyRequest.Name)
}

// TestClsMetricSubscribe_Delete_Success tests the Delete flow builds a
// DeleteMetricSubscribe request with TaskId and TopicId from the composite id.
func TestClsMetricSubscribe_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	clsClient := &clsv20201016.Client{}
	patches.ApplyMethodReturn(newMockMetaForClsMetricSubscribe().client, "UseClsV20201016Client", clsClient)

	var capturedRequest *clsv20201016.DeleteMetricSubscribeRequest
	patches.ApplyMethodFunc(clsClient, "DeleteMetricSubscribeWithContext", func(_ context.Context, request *clsv20201016.DeleteMetricSubscribeRequest) (*clsv20201016.DeleteMetricSubscribeResponse, error) {
		capturedRequest = request
		resp := clsv20201016.NewDeleteMetricSubscribeResponse()
		resp.Response = &clsv20201016.DeleteMetricSubscribeResponseParams{
			RequestId: ptrStringMS("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForClsMetricSubscribe()
	res := localcls.ResourceTencentCloudClsMetricSubscribe()
	d := schema.TestResourceDataRaw(t, res.Schema, metricSubscribeRawInput())
	d.SetId("topic-test-123" + tccommon.FILED_SP + "task-test-123")

	err := res.Delete(d, meta)
	assert.NoError(t, err)

	// Verify the delete request carried TaskId and TopicId from the composite id
	assert.NotNil(t, capturedRequest)
	assert.Equal(t, "task-test-123", *capturedRequest.TaskId)
	assert.Equal(t, "topic-test-123", *capturedRequest.TopicId)
}

// TestClsMetricSubscribe_Schema tests key schema definitions.
func TestClsMetricSubscribe_Schema(t *testing.T) {
	res := localcls.ResourceTencentCloudClsMetricSubscribe()

	assert.NotNil(t, res)

	// ForceNew fields
	topicIdSchema := res.Schema["topic_id"]
	assert.True(t, topicIdSchema.ForceNew)
	assert.True(t, topicIdSchema.Required)

	// Computed fields
	for _, field := range []string{"task_id", "status", "create_time", "update_time"} {
		s := res.Schema[field]
		assert.True(t, s.Computed, "field %s should be Computed", field)
	}

	// instance_info MaxItems=1
	instanceInfoSchema := res.Schema["instance_info"]
	assert.Equal(t, 1, instanceInfoSchema.MaxItems)

	// Importer registered
	assert.NotNil(t, res.Importer)
}
