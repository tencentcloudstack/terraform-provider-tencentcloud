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
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

type mockMetaEdgeKVNamespace struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaEdgeKVNamespace) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaEdgeKVNamespace{}

func newMockMetaEdgeKVNamespace() *mockMetaEdgeKVNamespace {
	return &mockMetaEdgeKVNamespace{client: &connectivity.TencentCloudClient{}}
}

func ptrStringEKVN(s string) *string {
	return &s
}

func ptrInt64EKVN(i int64) *int64 {
	return &i
}

func ptrBoolEKVN(b bool) *bool {
	return &b
}

func TestTeoEdgeKVNamespace_Create_Normal(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaEdgeKVNamespace().client, "UseTeoV20220901Client", teoClient)

	// Mock Create API
	patches.ApplyMethodFunc(teoClient, "CreateEdgeKVNamespaceWithContext", func(_ context.Context, request *teov20220901.CreateEdgeKVNamespaceRequest) (*teov20220901.CreateEdgeKVNamespaceResponse, error) {
		assert.NotNil(t, request.ZoneId)
		assert.Equal(t, "zone-test123", *request.ZoneId)
		assert.NotNil(t, request.Namespace)
		assert.Equal(t, "test-ns", *request.Namespace)
		assert.NotNil(t, request.Remark)
		assert.Equal(t, "test remark", *request.Remark)

		resp := teov20220901.NewCreateEdgeKVNamespaceResponse()
		resp.Response = &teov20220901.CreateEdgeKVNamespaceResponseParams{
			RequestId: ptrStringEKVN("fake-request-id-create"),
		}
		return resp, nil
	})

	// Mock Read API (called after Create)
	patches.ApplyMethodFunc(teoClient, "DescribeEdgeKVNamespacesWithContext", func(_ context.Context, request *teov20220901.DescribeEdgeKVNamespacesRequest) (*teov20220901.DescribeEdgeKVNamespacesResponse, error) {
		resp := teov20220901.NewDescribeEdgeKVNamespacesResponse()
		resp.Response = &teov20220901.DescribeEdgeKVNamespacesResponseParams{
			TotalCount: ptrInt64EKVN(1),
			KVNamespaces: []*teov20220901.KVNamespace{
				{
					Namespace:    ptrStringEKVN("test-ns"),
					Remark:       ptrStringEKVN("test remark"),
					Capacity:     ptrInt64EKVN(1073741824),
					CapacityUsed: ptrInt64EKVN(0),
					CreatedOn:    ptrStringEKVN("2024-01-01T00:00:00Z"),
					ModifiedOn:   ptrStringEKVN("2024-01-01T00:00:00Z"),
				},
			},
			RequestId: ptrStringEKVN("fake-request-id-read"),
		}
		return resp, nil
	})

	meta := newMockMetaEdgeKVNamespace()
	res := teo.ResourceTencentCloudTeoEdgeKVNamespace()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":   "zone-test123",
		"namespace": "test-ns",
		"remark":    "test remark",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-test123#test-ns", d.Id())
	assert.Equal(t, "test remark", d.Get("remark"))
	assert.Equal(t, 1073741824, d.Get("capacity"))
	assert.Equal(t, 0, d.Get("capacity_used"))
	assert.Equal(t, "2024-01-01T00:00:00Z", d.Get("created_on"))
	assert.Equal(t, "2024-01-01T00:00:00Z", d.Get("modified_on"))
}

func TestTeoEdgeKVNamespace_Create_WithoutRemark(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaEdgeKVNamespace().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateEdgeKVNamespaceWithContext", func(_ context.Context, request *teov20220901.CreateEdgeKVNamespaceRequest) (*teov20220901.CreateEdgeKVNamespaceResponse, error) {
		assert.NotNil(t, request.ZoneId)
		assert.Equal(t, "zone-test123", *request.ZoneId)
		assert.NotNil(t, request.Namespace)
		assert.Equal(t, "test-ns-no-remark", *request.Namespace)

		resp := teov20220901.NewCreateEdgeKVNamespaceResponse()
		resp.Response = &teov20220901.CreateEdgeKVNamespaceResponseParams{
			RequestId: ptrStringEKVN("fake-request-id-create"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeEdgeKVNamespacesWithContext", func(_ context.Context, request *teov20220901.DescribeEdgeKVNamespacesRequest) (*teov20220901.DescribeEdgeKVNamespacesResponse, error) {
		resp := teov20220901.NewDescribeEdgeKVNamespacesResponse()
		resp.Response = &teov20220901.DescribeEdgeKVNamespacesResponseParams{
			TotalCount: ptrInt64EKVN(1),
			KVNamespaces: []*teov20220901.KVNamespace{
				{
					Namespace:    ptrStringEKVN("test-ns-no-remark"),
					Remark:       ptrStringEKVN(""),
					Capacity:     ptrInt64EKVN(1073741824),
					CapacityUsed: ptrInt64EKVN(0),
					CreatedOn:    ptrStringEKVN("2024-01-01T00:00:00Z"),
					ModifiedOn:   ptrStringEKVN("2024-01-01T00:00:00Z"),
				},
			},
			RequestId: ptrStringEKVN("fake-request-id-read"),
		}
		return resp, nil
	})

	meta := newMockMetaEdgeKVNamespace()
	res := teo.ResourceTencentCloudTeoEdgeKVNamespace()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":   "zone-test123",
		"namespace": "test-ns-no-remark",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-test123#test-ns-no-remark", d.Id())
}

func TestTeoEdgeKVNamespace_Read_Normal(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaEdgeKVNamespace().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeEdgeKVNamespacesWithContext", func(_ context.Context, request *teov20220901.DescribeEdgeKVNamespacesRequest) (*teov20220901.DescribeEdgeKVNamespacesResponse, error) {
		assert.NotNil(t, request.ZoneId)
		assert.Equal(t, "zone-test123", *request.ZoneId)
		assert.NotNil(t, request.Limit)
		assert.Equal(t, int64(1000), *request.Limit)
		assert.NotNil(t, request.Filters)
		assert.Equal(t, 1, len(request.Filters))
		assert.Equal(t, "namespace", *request.Filters[0].Name)
		assert.Equal(t, "test-ns", *request.Filters[0].Values[0])
		assert.Equal(t, false, *request.Filters[0].Fuzzy)

		resp := teov20220901.NewDescribeEdgeKVNamespacesResponse()
		resp.Response = &teov20220901.DescribeEdgeKVNamespacesResponseParams{
			TotalCount: ptrInt64EKVN(1),
			KVNamespaces: []*teov20220901.KVNamespace{
				{
					Namespace:    ptrStringEKVN("test-ns"),
					Remark:       ptrStringEKVN("test remark"),
					Capacity:     ptrInt64EKVN(1073741824),
					CapacityUsed: ptrInt64EKVN(512),
					CreatedOn:    ptrStringEKVN("2024-01-01T00:00:00Z"),
					ModifiedOn:   ptrStringEKVN("2024-06-15T12:00:00Z"),
				},
			},
			RequestId: ptrStringEKVN("fake-request-id-read"),
		}
		return resp, nil
	})

	meta := newMockMetaEdgeKVNamespace()
	res := teo.ResourceTencentCloudTeoEdgeKVNamespace()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":   "zone-test123",
		"namespace": "test-ns",
		"remark":    "test remark",
	})
	d.SetId("zone-test123#test-ns")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-test123#test-ns", d.Id())
	assert.Equal(t, "zone-test123", d.Get("zone_id"))
	assert.Equal(t, "test-ns", d.Get("namespace"))
	assert.Equal(t, "test remark", d.Get("remark"))
	assert.Equal(t, 1073741824, d.Get("capacity"))
	assert.Equal(t, 512, d.Get("capacity_used"))
	assert.Equal(t, "2024-01-01T00:00:00Z", d.Get("created_on"))
	assert.Equal(t, "2024-06-15T12:00:00Z", d.Get("modified_on"))
}

func TestTeoEdgeKVNamespace_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaEdgeKVNamespace().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeEdgeKVNamespacesWithContext", func(_ context.Context, request *teov20220901.DescribeEdgeKVNamespacesRequest) (*teov20220901.DescribeEdgeKVNamespacesResponse, error) {
		resp := teov20220901.NewDescribeEdgeKVNamespacesResponse()
		resp.Response = &teov20220901.DescribeEdgeKVNamespacesResponseParams{
			TotalCount:   ptrInt64EKVN(0),
			KVNamespaces: []*teov20220901.KVNamespace{},
			RequestId:    ptrStringEKVN("fake-request-id-read"),
		}
		return resp, nil
	})

	meta := newMockMetaEdgeKVNamespace()
	res := teo.ResourceTencentCloudTeoEdgeKVNamespace()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":   "zone-test123",
		"namespace": "test-ns-notexist",
	})
	d.SetId("zone-test123#test-ns-notexist")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

func TestTeoEdgeKVNamespace_Update_Remark(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaEdgeKVNamespace().client, "UseTeoV20220901Client", teoClient)

	// Mock Modify API
	patches.ApplyMethodFunc(teoClient, "ModifyEdgeKVNamespaceWithContext", func(_ context.Context, request *teov20220901.ModifyEdgeKVNamespaceRequest) (*teov20220901.ModifyEdgeKVNamespaceResponse, error) {
		assert.NotNil(t, request.ZoneId)
		assert.Equal(t, "zone-test123", *request.ZoneId)
		assert.NotNil(t, request.Namespace)
		assert.Equal(t, "test-ns", *request.Namespace)
		assert.NotNil(t, request.Remark)
		assert.Equal(t, "updated remark", *request.Remark)

		resp := teov20220901.NewModifyEdgeKVNamespaceResponse()
		resp.Response = &teov20220901.ModifyEdgeKVNamespaceResponseParams{
			RequestId: ptrStringEKVN("fake-request-id-update"),
		}
		return resp, nil
	})

	// Mock Read API (called after Update)
	patches.ApplyMethodFunc(teoClient, "DescribeEdgeKVNamespacesWithContext", func(_ context.Context, request *teov20220901.DescribeEdgeKVNamespacesRequest) (*teov20220901.DescribeEdgeKVNamespacesResponse, error) {
		resp := teov20220901.NewDescribeEdgeKVNamespacesResponse()
		resp.Response = &teov20220901.DescribeEdgeKVNamespacesResponseParams{
			TotalCount: ptrInt64EKVN(1),
			KVNamespaces: []*teov20220901.KVNamespace{
				{
					Namespace:    ptrStringEKVN("test-ns"),
					Remark:       ptrStringEKVN("updated remark"),
					Capacity:     ptrInt64EKVN(1073741824),
					CapacityUsed: ptrInt64EKVN(0),
					CreatedOn:    ptrStringEKVN("2024-01-01T00:00:00Z"),
					ModifiedOn:   ptrStringEKVN("2024-06-15T12:00:00Z"),
				},
			},
			RequestId: ptrStringEKVN("fake-request-id-read"),
		}
		return resp, nil
	})

	meta := newMockMetaEdgeKVNamespace()
	res := teo.ResourceTencentCloudTeoEdgeKVNamespace()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":   "zone-test123",
		"namespace": "test-ns",
		"remark":    "updated remark",
	})
	d.SetId("zone-test123#test-ns")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-test123#test-ns", d.Id())
	assert.Equal(t, "updated remark", d.Get("remark"))
}

func TestTeoEdgeKVNamespace_Delete_Normal(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaEdgeKVNamespace().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteEdgeKVNamespaceWithContext", func(_ context.Context, request *teov20220901.DeleteEdgeKVNamespaceRequest) (*teov20220901.DeleteEdgeKVNamespaceResponse, error) {
		assert.NotNil(t, request.ZoneId)
		assert.Equal(t, "zone-test123", *request.ZoneId)
		assert.NotNil(t, request.Namespace)
		assert.Equal(t, "test-ns", *request.Namespace)

		resp := teov20220901.NewDeleteEdgeKVNamespaceResponse()
		resp.Response = &teov20220901.DeleteEdgeKVNamespaceResponseParams{
			RequestId: ptrStringEKVN("fake-request-id-delete"),
		}
		return resp, nil
	})

	meta := newMockMetaEdgeKVNamespace()
	res := teo.ResourceTencentCloudTeoEdgeKVNamespace()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":   "zone-test123",
		"namespace": "test-ns",
		"remark":    "test remark",
	})
	d.SetId("zone-test123#test-ns")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}
