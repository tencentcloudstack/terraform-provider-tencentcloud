package mongodb_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	mongodb_sdk "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/mongodb/v20190725"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/mongodb"
)

type mockMetaForCpuScale struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForCpuScale) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForCpuScale{}

func newMockMetaForCpuScale() *mockMetaForCpuScale {
	return &mockMetaForCpuScale{client: &connectivity.TencentCloudClient{}}
}

func ptrStringCpu(s string) *string { return &s }
func ptrInt64Cpu(v int64) *int64    { return &v }

// Test: Create scenario (enable CPU elastic scaling)
func TestTencentCloudMongodbInstanceCpuScaleConfig_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mongodbClient := &mongodb_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForCpuScale().client, "UseMongodbClient", mongodbClient)

	// Mock ScaleUpDBInstanceCpu
	patches.ApplyMethodFunc(mongodbClient, "ScaleUpDBInstanceCpu", func(request *mongodb_sdk.ScaleUpDBInstanceCpuRequest) (*mongodb_sdk.ScaleUpDBInstanceCpuResponse, error) {
		resp := mongodb_sdk.NewScaleUpDBInstanceCpuResponse()
		resp.Response = &mongodb_sdk.ScaleUpDBInstanceCpuResponseParams{
			FlowId:    ptrInt64Cpu(12345),
			RequestId: ptrStringCpu("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeDBInstances (read after create)
	patches.ApplyMethodFunc(mongodbClient, "DescribeDBInstances", func(request *mongodb_sdk.DescribeDBInstancesRequest) (*mongodb_sdk.DescribeDBInstancesResponse, error) {
		resp := mongodb_sdk.NewDescribeDBInstancesResponse()
		resp.Response = &mongodb_sdk.DescribeDBInstancesResponseParams{
			InstanceDetails: []*mongodb_sdk.InstanceDetail{
				{
					InstanceId:   ptrStringCpu("cmgo-test1234"),
					InstanceName: ptrStringCpu("test-instance"),
					Status:       ptrInt64Cpu(2),
				},
			},
			RequestId: ptrStringCpu("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeAsyncRequestInfo
	patches.ApplyMethodFunc(mongodbClient, "DescribeAsyncRequestInfo", func(request *mongodb_sdk.DescribeAsyncRequestInfoRequest) (*mongodb_sdk.DescribeAsyncRequestInfoResponse, error) {
		resp := mongodb_sdk.NewDescribeAsyncRequestInfoResponse()
		resp.Response = &mongodb_sdk.DescribeAsyncRequestInfoResponseParams{
			Status:    ptrStringCpu("success"),
			RequestId: ptrStringCpu("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForCpuScale()
	res := mongodb.ResourceTencentCloudMongodbInstanceCpuScaleConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "cmgo-test1234",
		"extra_cpu":   2,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "cmgo-test1234", d.Id())
}

// Test: Read scenario (verify instance exists)
func TestTencentCloudMongodbInstanceCpuScaleConfig_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mongodbClient := &mongodb_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForCpuScale().client, "UseMongodbClient", mongodbClient)

	patches.ApplyMethodFunc(mongodbClient, "DescribeDBInstances", func(request *mongodb_sdk.DescribeDBInstancesRequest) (*mongodb_sdk.DescribeDBInstancesResponse, error) {
		resp := mongodb_sdk.NewDescribeDBInstancesResponse()
		resp.Response = &mongodb_sdk.DescribeDBInstancesResponseParams{
			InstanceDetails: []*mongodb_sdk.InstanceDetail{
				{
					InstanceId:   ptrStringCpu("cmgo-test1234"),
					InstanceName: ptrStringCpu("test-instance"),
					Status:       ptrInt64Cpu(2),
				},
			},
			RequestId: ptrStringCpu("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForCpuScale()
	res := mongodb.ResourceTencentCloudMongodbInstanceCpuScaleConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "cmgo-test1234",
		"extra_cpu":   2,
	})
	d.SetId("cmgo-test1234")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "cmgo-test1234", d.Get("instance_id"))
}

// Test: Read scenario (instance not found)
func TestTencentCloudMongodbInstanceCpuScaleConfig_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mongodbClient := &mongodb_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForCpuScale().client, "UseMongodbClient", mongodbClient)

	patches.ApplyMethodFunc(mongodbClient, "DescribeDBInstances", func(request *mongodb_sdk.DescribeDBInstancesRequest) (*mongodb_sdk.DescribeDBInstancesResponse, error) {
		resp := mongodb_sdk.NewDescribeDBInstancesResponse()
		resp.Response = &mongodb_sdk.DescribeDBInstancesResponseParams{
			InstanceDetails: []*mongodb_sdk.InstanceDetail{},
			RequestId:       ptrStringCpu("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForCpuScale()
	res := mongodb.ResourceTencentCloudMongodbInstanceCpuScaleConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "cmgo-test1234",
	})
	d.SetId("cmgo-test1234")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// Test: Update scenario (adjust extra_cpu)
func TestTencentCloudMongodbInstanceCpuScaleConfig_Update_ScaleUp(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mongodbClient := &mongodb_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForCpuScale().client, "UseMongodbClient", mongodbClient)

	// Mock ScaleUpDBInstanceCpu
	patches.ApplyMethodFunc(mongodbClient, "ScaleUpDBInstanceCpu", func(request *mongodb_sdk.ScaleUpDBInstanceCpuRequest) (*mongodb_sdk.ScaleUpDBInstanceCpuResponse, error) {
		resp := mongodb_sdk.NewScaleUpDBInstanceCpuResponse()
		resp.Response = &mongodb_sdk.ScaleUpDBInstanceCpuResponseParams{
			FlowId:    ptrInt64Cpu(12346),
			RequestId: ptrStringCpu("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeDBInstances
	patches.ApplyMethodFunc(mongodbClient, "DescribeDBInstances", func(request *mongodb_sdk.DescribeDBInstancesRequest) (*mongodb_sdk.DescribeDBInstancesResponse, error) {
		resp := mongodb_sdk.NewDescribeDBInstancesResponse()
		resp.Response = &mongodb_sdk.DescribeDBInstancesResponseParams{
			InstanceDetails: []*mongodb_sdk.InstanceDetail{
				{
					InstanceId:   ptrStringCpu("cmgo-test1234"),
					InstanceName: ptrStringCpu("test-instance"),
					Status:       ptrInt64Cpu(2),
				},
			},
			RequestId: ptrStringCpu("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeAsyncRequestInfo
	patches.ApplyMethodFunc(mongodbClient, "DescribeAsyncRequestInfo", func(request *mongodb_sdk.DescribeAsyncRequestInfoRequest) (*mongodb_sdk.DescribeAsyncRequestInfoResponse, error) {
		resp := mongodb_sdk.NewDescribeAsyncRequestInfoResponse()
		resp.Response = &mongodb_sdk.DescribeAsyncRequestInfoResponseParams{
			Status:    ptrStringCpu("success"),
			RequestId: ptrStringCpu("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForCpuScale()
	res := mongodb.ResourceTencentCloudMongodbInstanceCpuScaleConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "cmgo-test1234",
		"extra_cpu":   4,
	})
	d.SetId("cmgo-test1234")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

// Test: Update scenario (remove extra_cpu → scale down)
func TestTencentCloudMongodbInstanceCpuScaleConfig_Update_ScaleDown(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mongodbClient := &mongodb_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForCpuScale().client, "UseMongodbClient", mongodbClient)

	// Mock ScaleDownDBInstanceCpu
	patches.ApplyMethodFunc(mongodbClient, "ScaleDownDBInstanceCpu", func(request *mongodb_sdk.ScaleDownDBInstanceCpuRequest) (*mongodb_sdk.ScaleDownDBInstanceCpuResponse, error) {
		resp := mongodb_sdk.NewScaleDownDBInstanceCpuResponse()
		resp.Response = &mongodb_sdk.ScaleDownDBInstanceCpuResponseParams{
			FlowId:    ptrInt64Cpu(12347),
			RequestId: ptrStringCpu("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeDBInstances
	patches.ApplyMethodFunc(mongodbClient, "DescribeDBInstances", func(request *mongodb_sdk.DescribeDBInstancesRequest) (*mongodb_sdk.DescribeDBInstancesResponse, error) {
		resp := mongodb_sdk.NewDescribeDBInstancesResponse()
		resp.Response = &mongodb_sdk.DescribeDBInstancesResponseParams{
			InstanceDetails: []*mongodb_sdk.InstanceDetail{
				{
					InstanceId:   ptrStringCpu("cmgo-test1234"),
					InstanceName: ptrStringCpu("test-instance"),
					Status:       ptrInt64Cpu(2),
				},
			},
			RequestId: ptrStringCpu("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeAsyncRequestInfo
	patches.ApplyMethodFunc(mongodbClient, "DescribeAsyncRequestInfo", func(request *mongodb_sdk.DescribeAsyncRequestInfoRequest) (*mongodb_sdk.DescribeAsyncRequestInfoResponse, error) {
		resp := mongodb_sdk.NewDescribeAsyncRequestInfoResponse()
		resp.Response = &mongodb_sdk.DescribeAsyncRequestInfoResponseParams{
			Status:    ptrStringCpu("success"),
			RequestId: ptrStringCpu("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForCpuScale()
	res := mongodb.ResourceTencentCloudMongodbInstanceCpuScaleConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "cmgo-test1234",
		"extra_cpu":   0,
	})
	d.SetId("cmgo-test1234")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

// Test: Delete scenario (close elastic scaling)
func TestTencentCloudMongodbInstanceCpuScaleConfig_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	mongodbClient := &mongodb_sdk.Client{}
	patches.ApplyMethodReturn(newMockMetaForCpuScale().client, "UseMongodbClient", mongodbClient)

	// Mock ScaleDownDBInstanceCpu
	patches.ApplyMethodFunc(mongodbClient, "ScaleDownDBInstanceCpu", func(request *mongodb_sdk.ScaleDownDBInstanceCpuRequest) (*mongodb_sdk.ScaleDownDBInstanceCpuResponse, error) {
		resp := mongodb_sdk.NewScaleDownDBInstanceCpuResponse()
		resp.Response = &mongodb_sdk.ScaleDownDBInstanceCpuResponseParams{
			FlowId:    ptrInt64Cpu(12348),
			RequestId: ptrStringCpu("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeAsyncRequestInfo
	patches.ApplyMethodFunc(mongodbClient, "DescribeAsyncRequestInfo", func(request *mongodb_sdk.DescribeAsyncRequestInfoRequest) (*mongodb_sdk.DescribeAsyncRequestInfoResponse, error) {
		resp := mongodb_sdk.NewDescribeAsyncRequestInfoResponse()
		resp.Response = &mongodb_sdk.DescribeAsyncRequestInfoResponseParams{
			Status:    ptrStringCpu("success"),
			RequestId: ptrStringCpu("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForCpuScale()
	res := mongodb.ResourceTencentCloudMongodbInstanceCpuScaleConfig()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"instance_id": "cmgo-test1234",
		"extra_cpu":   2,
	})
	d.SetId("cmgo-test1234")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}
