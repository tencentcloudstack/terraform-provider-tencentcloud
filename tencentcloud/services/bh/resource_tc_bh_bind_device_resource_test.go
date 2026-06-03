package bh_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	bhv20230418 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/bh/v20230418"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/bh"
)

type mockMetaBindDeviceResource struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaBindDeviceResource) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaBindDeviceResource{}

func newMockMetaBindDeviceResource() *mockMetaBindDeviceResource {
	return &mockMetaBindDeviceResource{client: &connectivity.TencentCloudClient{}}
}

func ptrStringBDR(s string) *string {
	return &s
}

func ptrUint64BDR(i uint64) *uint64 {
	return &i
}

func TestBhBindDeviceResource_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	bhClient := &bhv20230418.Client{}
	patches.ApplyMethodReturn(newMockMetaBindDeviceResource().client, "UseBhV20230418Client", bhClient)

	patches.ApplyMethodFunc(bhClient, "BindDeviceResourceWithContext", func(_ context.Context, request *bhv20230418.BindDeviceResourceRequest) (*bhv20230418.BindDeviceResourceResponse, error) {
		assert.NotNil(t, request.DeviceIdSet)
		assert.Equal(t, 2, len(request.DeviceIdSet))
		assert.Equal(t, uint64(123), *request.DeviceIdSet[0])
		assert.Equal(t, uint64(456), *request.DeviceIdSet[1])
		assert.NotNil(t, request.ResourceId)
		assert.Equal(t, "bh-saas-abc123", *request.ResourceId)

		resp := bhv20230418.NewBindDeviceResourceResponse()
		resp.Response = &bhv20230418.BindDeviceResourceResponseParams{
			RequestId: ptrStringBDR("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(bhClient, "DescribeDevicesWithContext", func(_ context.Context, request *bhv20230418.DescribeDevicesRequest) (*bhv20230418.DescribeDevicesResponse, error) {
		resp := bhv20230418.NewDescribeDevicesResponse()
		resp.Response = &bhv20230418.DescribeDevicesResponseParams{
			TotalCount: ptrUint64BDR(1),
			DeviceSet: []*bhv20230418.Device{
				{
					Id:       ptrUint64BDR(123),
					DomainId: ptrStringBDR("dm-domain01"),
					Resource: &bhv20230418.Resource{
						ResourceId: ptrStringBDR("bh-saas-abc123"),
					},
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaBindDeviceResource()
	res := bh.ResourceTencentCloudBhBindDeviceResource()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"device_id_set": []interface{}{123, 456},
		"resource_id":   "bh-saas-abc123",
		"domain_id":     "dm-domain01",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "123,456#bh-saas-abc123", d.Id())
}

func TestBhBindDeviceResource_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	bhClient := &bhv20230418.Client{}
	patches.ApplyMethodReturn(newMockMetaBindDeviceResource().client, "UseBhV20230418Client", bhClient)

	patches.ApplyMethodFunc(bhClient, "DescribeDevicesWithContext", func(_ context.Context, request *bhv20230418.DescribeDevicesRequest) (*bhv20230418.DescribeDevicesResponse, error) {
		assert.NotNil(t, request.IdSet)
		assert.Equal(t, uint64(123), *request.IdSet[0])

		resp := bhv20230418.NewDescribeDevicesResponse()
		resp.Response = &bhv20230418.DescribeDevicesResponseParams{
			TotalCount: ptrUint64BDR(1),
			DeviceSet: []*bhv20230418.Device{
				{
					Id:              ptrUint64BDR(123),
					DomainId:        ptrStringBDR("dm-domain01"),
					ManageDimension: ptrUint64BDR(1),
					ManageAccountId: ptrUint64BDR(100),
					Namespace:       ptrStringBDR("default"),
					Workload:        ptrStringBDR("deployment/nginx"),
					Resource: &bhv20230418.Resource{
						ResourceId: ptrStringBDR("bh-saas-abc123"),
					},
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaBindDeviceResource()
	res := bh.ResourceTencentCloudBhBindDeviceResource()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"device_id_set": []interface{}{123, 456},
		"resource_id":   "bh-saas-abc123",
	})
	d.SetId("123,456#bh-saas-abc123")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "dm-domain01", d.Get("domain_id"))
	assert.Equal(t, 1, d.Get("manage_dimension"))
	assert.Equal(t, 100, d.Get("manage_account_id"))
	assert.Equal(t, "default", d.Get("namespace"))
	assert.Equal(t, "deployment/nginx", d.Get("workload"))
}

func TestBhBindDeviceResource_ReadDeviceNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	bhClient := &bhv20230418.Client{}
	patches.ApplyMethodReturn(newMockMetaBindDeviceResource().client, "UseBhV20230418Client", bhClient)

	patches.ApplyMethodFunc(bhClient, "DescribeDevicesWithContext", func(_ context.Context, request *bhv20230418.DescribeDevicesRequest) (*bhv20230418.DescribeDevicesResponse, error) {
		resp := bhv20230418.NewDescribeDevicesResponse()
		resp.Response = &bhv20230418.DescribeDevicesResponseParams{
			TotalCount: ptrUint64BDR(0),
			DeviceSet:  []*bhv20230418.Device{},
		}
		return resp, nil
	})

	meta := newMockMetaBindDeviceResource()
	res := bh.ResourceTencentCloudBhBindDeviceResource()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"device_id_set": []interface{}{123},
		"resource_id":   "bh-saas-abc123",
	})
	d.SetId("123#bh-saas-abc123")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestBhBindDeviceResource_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	bhClient := &bhv20230418.Client{}
	patches.ApplyMethodReturn(newMockMetaBindDeviceResource().client, "UseBhV20230418Client", bhClient)

	patches.ApplyMethodFunc(bhClient, "BindDeviceResourceWithContext", func(_ context.Context, request *bhv20230418.BindDeviceResourceRequest) (*bhv20230418.BindDeviceResourceResponse, error) {
		assert.NotNil(t, request.DeviceIdSet)
		assert.Equal(t, 1, len(request.DeviceIdSet))
		assert.Equal(t, uint64(123), *request.DeviceIdSet[0])
		assert.NotNil(t, request.ResourceId)
		assert.Equal(t, "bh-saas-new456", *request.ResourceId)

		resp := bhv20230418.NewBindDeviceResourceResponse()
		resp.Response = &bhv20230418.BindDeviceResourceResponseParams{
			RequestId: ptrStringBDR("fake-request-id-update"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(bhClient, "DescribeDevicesWithContext", func(_ context.Context, request *bhv20230418.DescribeDevicesRequest) (*bhv20230418.DescribeDevicesResponse, error) {
		resp := bhv20230418.NewDescribeDevicesResponse()
		resp.Response = &bhv20230418.DescribeDevicesResponseParams{
			TotalCount: ptrUint64BDR(1),
			DeviceSet: []*bhv20230418.Device{
				{
					Id:       ptrUint64BDR(123),
					DomainId: ptrStringBDR("dm-domain01"),
					Resource: &bhv20230418.Resource{
						ResourceId: ptrStringBDR("bh-saas-new456"),
					},
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaBindDeviceResource()
	res := bh.ResourceTencentCloudBhBindDeviceResource()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"device_id_set": []interface{}{123},
		"resource_id":   "bh-saas-new456",
		"domain_id":     "dm-domain01",
	})
	d.SetId("123#bh-saas-abc123")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

func TestBhBindDeviceResource_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	bhClient := &bhv20230418.Client{}
	patches.ApplyMethodReturn(newMockMetaBindDeviceResource().client, "UseBhV20230418Client", bhClient)

	patches.ApplyMethodFunc(bhClient, "BindDeviceResourceWithContext", func(_ context.Context, request *bhv20230418.BindDeviceResourceRequest) (*bhv20230418.BindDeviceResourceResponse, error) {
		assert.NotNil(t, request.DeviceIdSet)
		assert.Equal(t, 2, len(request.DeviceIdSet))
		assert.Equal(t, uint64(123), *request.DeviceIdSet[0])
		assert.Equal(t, uint64(456), *request.DeviceIdSet[1])
		assert.NotNil(t, request.ResourceId)
		assert.Equal(t, "", *request.ResourceId)

		resp := bhv20230418.NewBindDeviceResourceResponse()
		resp.Response = &bhv20230418.BindDeviceResourceResponseParams{
			RequestId: ptrStringBDR("fake-request-id-delete"),
		}
		return resp, nil
	})

	meta := newMockMetaBindDeviceResource()
	res := bh.ResourceTencentCloudBhBindDeviceResource()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"device_id_set": []interface{}{123, 456},
		"resource_id":   "bh-saas-abc123",
	})
	d.SetId("123,456#bh-saas-abc123")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}
