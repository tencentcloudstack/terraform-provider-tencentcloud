package mongodb_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	mongodb_sdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mongodb/v20190725"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/mongodb"
)

type mockMetaForScaleUpCpu struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForScaleUpCpu) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForScaleUpCpu{}

func newMockMetaForScaleUpCpu() *mockMetaForScaleUpCpu {
	return &mockMetaForScaleUpCpu{client: &connectivity.TencentCloudClient{}}
}

func ptrStrScaleUp(s string) *string { return &s }
func ptrBoolScaleUp(b bool) *bool    { return &b }
func ptrInt64ScaleUp(v int64) *int64 { return &v }

func TestMongodbScaleUpDbInstanceCpu_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mongodbClient := &mongodb_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForScaleUpCpu().client, "UseMongodbClient", mongodbClient)

	var (
		scaleUpCallCount       int
		describeAsyncCallCount int
		waitedFlowId           string
	)
	patches.ApplyMethodFunc(mongodbClient, "ScaleUpDBInstanceCpuWithContext", func(ctx context.Context, request *mongodb_sdk.ScaleUpDBInstanceCpuRequest) (*mongodb_sdk.ScaleUpDBInstanceCpuResponse, error) {
		scaleUpCallCount++
		resp := mongodb_sdk.NewScaleUpDBInstanceCpuResponse()
		resp.Response = &mongodb_sdk.ScaleUpDBInstanceCpuResponseParams{
			FlowId:    ptrInt64ScaleUp(12345),
			RequestId: ptrStrScaleUp("fake-request-id"),
		}
		return resp, nil
	})

	// Create waits on the flow ID returned by ScaleUpDBInstanceCpu.
	patches.ApplyMethodFunc(mongodbClient, "DescribeAsyncRequestInfo", func(request *mongodb_sdk.DescribeAsyncRequestInfoRequest) (*mongodb_sdk.DescribeAsyncRequestInfoResponse, error) {
		describeAsyncCallCount++
		waitedFlowId = *request.AsyncRequestId
		resp := mongodb_sdk.NewDescribeAsyncRequestInfoResponse()
		resp.Response = &mongodb_sdk.DescribeAsyncRequestInfoResponseParams{
			Status:    ptrStrScaleUp(mongodb.MONGODB_TASK_SUCCESS),
			RequestId: ptrStrScaleUp("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(mongodbClient, "CheckDBInstanceElasticCpuScalable", func(request *mongodb_sdk.CheckDBInstanceElasticCpuScalableRequest) (*mongodb_sdk.CheckDBInstanceElasticCpuScalableResponse, error) {
		resp := mongodb_sdk.NewCheckDBInstanceElasticCpuScalableResponse()
		resp.Response = &mongodb_sdk.CheckDBInstanceElasticCpuScalableResponseParams{
			Scalable:    ptrBoolScaleUp(true),
			IsScaled:    ptrBoolScaleUp(true),
			IsLocked:    ptrBoolScaleUp(false),
			Reason:      ptrStrScaleUp(""),
			MaxExtraCpu: ptrInt64ScaleUp(24),
			ExtraCpu:    ptrInt64ScaleUp(4),
			TriggerType: ptrInt64ScaleUp(1),
			ScaleUpTime: ptrStrScaleUp("2024-01-01 00:00:00"),
			RequestId:   ptrStrScaleUp("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForScaleUpCpu()
	res := mongodb.ResourceTencentCloudMongodbScaleUpDbInstanceCpu()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "cmgo-test1234",
		"extra_cpu":   4,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "cmgo-test1234", d.Id())
	assert.Equal(t, 1, scaleUpCallCount)
	assert.Equal(t, 1, describeAsyncCallCount)
	assert.Equal(t, "12345", waitedFlowId)
	assert.Equal(t, true, d.Get("is_scaled"))
	assert.Equal(t, 4, d.Get("extra_cpu"))
}

func TestMongodbScaleUpDbInstanceCpu_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mongodbClient := &mongodb_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForScaleUpCpu().client, "UseMongodbClient", mongodbClient)

	patches.ApplyMethodFunc(mongodbClient, "CheckDBInstanceElasticCpuScalable", func(request *mongodb_sdk.CheckDBInstanceElasticCpuScalableRequest) (*mongodb_sdk.CheckDBInstanceElasticCpuScalableResponse, error) {
		resp := mongodb_sdk.NewCheckDBInstanceElasticCpuScalableResponse()
		resp.Response = &mongodb_sdk.CheckDBInstanceElasticCpuScalableResponseParams{
			Scalable:    ptrBoolScaleUp(true),
			IsScaled:    ptrBoolScaleUp(true),
			IsLocked:    ptrBoolScaleUp(false),
			Reason:      ptrStrScaleUp(""),
			MaxExtraCpu: ptrInt64ScaleUp(24),
			ExtraCpu:    ptrInt64ScaleUp(4),
			TriggerType: ptrInt64ScaleUp(1),
			ScaleUpTime: ptrStrScaleUp("2024-01-01 00:00:00"),
			RequestId:   ptrStrScaleUp("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForScaleUpCpu()
	res := mongodb.ResourceTencentCloudMongodbScaleUpDbInstanceCpu()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "cmgo-test1234",
	})
	d.SetId("cmgo-test1234")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "cmgo-test1234", d.Id())
	assert.Equal(t, "cmgo-test1234", d.Get("instance_id"))
	assert.Equal(t, true, d.Get("scalable"))
	assert.Equal(t, true, d.Get("is_scaled"))
	assert.Equal(t, false, d.Get("is_locked"))
	assert.Equal(t, 24, d.Get("max_extra_cpu"))
	assert.Equal(t, 4, d.Get("extra_cpu"))
	assert.Equal(t, 1, d.Get("trigger_type"))
	assert.Equal(t, "2024-01-01 00:00:00", d.Get("scale_up_time"))
}

func TestMongodbScaleUpDbInstanceCpu_Read_EmptyResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mongodbClient := &mongodb_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForScaleUpCpu().client, "UseMongodbClient", mongodbClient)

	patches.ApplyMethodFunc(mongodbClient, "CheckDBInstanceElasticCpuScalable", func(request *mongodb_sdk.CheckDBInstanceElasticCpuScalableRequest) (*mongodb_sdk.CheckDBInstanceElasticCpuScalableResponse, error) {
		resp := mongodb_sdk.NewCheckDBInstanceElasticCpuScalableResponse()
		resp.Response = nil
		return resp, nil
	})

	meta := newMockMetaForScaleUpCpu()
	res := mongodb.ResourceTencentCloudMongodbScaleUpDbInstanceCpu()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "cmgo-test1234",
	})
	d.SetId("cmgo-test1234")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestMongodbScaleUpDbInstanceCpu_SchemaImmutable(t *testing.T) {
	res := mongodb.ResourceTencentCloudMongodbScaleUpDbInstanceCpu()

	// Both fields are mandatory inputs, and extra_cpu can only be set on create:
	// changing it forces a replacement, and the resource deliberately exposes no
	// Update handler.
	assert.True(t, res.Schema["extra_cpu"].Required)
	assert.True(t, res.Schema["instance_id"].Required)
	assert.False(t, res.Schema["extra_cpu"].Computed)
	assert.True(t, res.Schema["extra_cpu"].ForceNew)
	assert.True(t, res.Schema["instance_id"].ForceNew)
	assert.Nil(t, res.Update)

	// Create and Delete both wait for the scale flow to finish, so the timeouts
	// must stay declared.
	assert.NotNil(t, res.Timeouts)
	assert.NotNil(t, res.Timeouts.Create)
	assert.Equal(t, 30*time.Minute, *res.Timeouts.Create)
	assert.NotNil(t, res.Timeouts.Delete)
	assert.Equal(t, 30*time.Minute, *res.Timeouts.Delete)
}

func TestMongodbScaleUpDbInstanceCpu_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mongodbClient := &mongodb_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForScaleUpCpu().client, "UseMongodbClient", mongodbClient)

	var (
		scaleDownCallCount     int
		describeAsyncCallCount int
		waitedFlowId           string
	)
	patches.ApplyMethodFunc(mongodbClient, "ScaleDownDBInstanceCpuWithContext", func(ctx context.Context, request *mongodb_sdk.ScaleDownDBInstanceCpuRequest) (*mongodb_sdk.ScaleDownDBInstanceCpuResponse, error) {
		scaleDownCallCount++
		resp := mongodb_sdk.NewScaleDownDBInstanceCpuResponse()
		resp.Response = &mongodb_sdk.ScaleDownDBInstanceCpuResponseParams{
			FlowId:    ptrInt64ScaleUp(99999),
			RequestId: ptrStrScaleUp("fake-request-id"),
		}
		return resp, nil
	})

	// Delete waits on the flow ID returned by ScaleDownDBInstanceCpu.
	patches.ApplyMethodFunc(mongodbClient, "DescribeAsyncRequestInfo", func(request *mongodb_sdk.DescribeAsyncRequestInfoRequest) (*mongodb_sdk.DescribeAsyncRequestInfoResponse, error) {
		describeAsyncCallCount++
		waitedFlowId = *request.AsyncRequestId
		resp := mongodb_sdk.NewDescribeAsyncRequestInfoResponse()
		resp.Response = &mongodb_sdk.DescribeAsyncRequestInfoResponseParams{
			Status:    ptrStrScaleUp(mongodb.MONGODB_TASK_SUCCESS),
			RequestId: ptrStrScaleUp("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForScaleUpCpu()
	res := mongodb.ResourceTencentCloudMongodbScaleUpDbInstanceCpu()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "cmgo-test1234",
	})
	d.SetId("cmgo-test1234")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, 1, scaleDownCallCount)
	assert.Equal(t, 1, describeAsyncCallCount)
	assert.Equal(t, "99999", waitedFlowId)
}

func TestMongodbScaleUpDbInstanceCpu_Delete_NotScaledTolerance(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mongodbClient := &mongodb_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForScaleUpCpu().client, "UseMongodbClient", mongodbClient)

	patches.ApplyMethodFunc(mongodbClient, "ScaleDownDBInstanceCpuWithContext", func(ctx context.Context, request *mongodb_sdk.ScaleDownDBInstanceCpuRequest) (*mongodb_sdk.ScaleDownDBInstanceCpuResponse, error) {
		return nil, fmt.Errorf("instance is not in scaled state")
	})

	meta := newMockMetaForScaleUpCpu()
	res := mongodb.ResourceTencentCloudMongodbScaleUpDbInstanceCpu()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "cmgo-test1234",
	})
	d.SetId("cmgo-test1234")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}
