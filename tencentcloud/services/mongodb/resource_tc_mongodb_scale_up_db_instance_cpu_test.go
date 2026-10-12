package mongodb_test

import (
	"fmt"
	"testing"

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

	var scaleUpCallCount int
	patches.ApplyMethodFunc(mongodbClient, "ScaleUpDBInstanceCpu", func(request *mongodb_sdk.ScaleUpDBInstanceCpuRequest) (*mongodb_sdk.ScaleUpDBInstanceCpuResponse, error) {
		scaleUpCallCount++
		resp := mongodb_sdk.NewScaleUpDBInstanceCpuResponse()
		resp.Response = &mongodb_sdk.ScaleUpDBInstanceCpuResponseParams{
			FlowId:    ptrInt64ScaleUp(12345),
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
	assert.Equal(t, 12345, d.Get("flow_id"))
	assert.Equal(t, true, d.Get("is_scaled"))
	assert.Equal(t, int64(4), d.Get("extra_cpu"))
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
	assert.Equal(t, int64(4), d.Get("extra_cpu"))
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

func TestMongodbScaleUpDbInstanceCpu_Update_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mongodbClient := &mongodb_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForScaleUpCpu().client, "UseMongodbClient", mongodbClient)

	var scaleUpCallCount int
	patches.ApplyMethodFunc(mongodbClient, "ScaleUpDBInstanceCpu", func(request *mongodb_sdk.ScaleUpDBInstanceCpuRequest) (*mongodb_sdk.ScaleUpDBInstanceCpuResponse, error) {
		scaleUpCallCount++
		resp := mongodb_sdk.NewScaleUpDBInstanceCpuResponse()
		resp.Response = &mongodb_sdk.ScaleUpDBInstanceCpuResponseParams{
			FlowId:    ptrInt64ScaleUp(67890),
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
			ExtraCpu:    ptrInt64ScaleUp(8),
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
		"extra_cpu":   8,
	})
	d.SetId("cmgo-test1234")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, 1, scaleUpCallCount)
	assert.Equal(t, 67890, d.Get("flow_id"))
	assert.Equal(t, int64(8), d.Get("extra_cpu"))
}

func TestMongodbScaleUpDbInstanceCpu_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mongodbClient := &mongodb_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForScaleUpCpu().client, "UseMongodbClient", mongodbClient)

	var scaleDownCallCount int
	patches.ApplyMethodFunc(mongodbClient, "ScaleDownDBInstanceCpu", func(request *mongodb_sdk.ScaleDownDBInstanceCpuRequest) (*mongodb_sdk.ScaleDownDBInstanceCpuResponse, error) {
		scaleDownCallCount++
		resp := mongodb_sdk.NewScaleDownDBInstanceCpuResponse()
		resp.Response = &mongodb_sdk.ScaleDownDBInstanceCpuResponseParams{
			FlowId:    ptrInt64ScaleUp(99999),
			RequestId: ptrStrScaleUp("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(mongodbClient, "CheckDBInstanceElasticCpuScalable", func(request *mongodb_sdk.CheckDBInstanceElasticCpuScalableRequest) (*mongodb_sdk.CheckDBInstanceElasticCpuScalableResponse, error) {
		resp := mongodb_sdk.NewCheckDBInstanceElasticCpuScalableResponse()
		resp.Response = &mongodb_sdk.CheckDBInstanceElasticCpuScalableResponseParams{
			Scalable:    ptrBoolScaleUp(true),
			IsScaled:    ptrBoolScaleUp(false),
			IsLocked:    ptrBoolScaleUp(false),
			Reason:      ptrStrScaleUp(""),
			MaxExtraCpu: ptrInt64ScaleUp(24),
			ExtraCpu:    ptrInt64ScaleUp(0),
			TriggerType: ptrInt64ScaleUp(0),
			ScaleUpTime: ptrStrScaleUp(""),
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

	err := res.Delete(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, 1, scaleDownCallCount)
}

func TestMongodbScaleUpDbInstanceCpu_Delete_NotScaledTolerance(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mongodbClient := &mongodb_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForScaleUpCpu().client, "UseMongodbClient", mongodbClient)

	patches.ApplyMethodFunc(mongodbClient, "ScaleDownDBInstanceCpu", func(request *mongodb_sdk.ScaleDownDBInstanceCpuRequest) (*mongodb_sdk.ScaleDownDBInstanceCpuResponse, error) {
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
