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

func ptrStringEdgeKVNamespace(s string) *string {
	return &s
}

func helper_int64(v int64) *int64 {
	return &v
}

// go test ./tencentcloud/services/teo/ -run "TestTeoEdgeKVNamespace" -v -count=1 -gcflags="all=-l"

// TestTeoEdgeKVNamespace_Create tests Create operation
func TestTeoEdgeKVNamespace_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaEdgeKVNamespace().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateEdgeKVNamespaceWithContext", func(_ context.Context, request *teov20220901.CreateEdgeKVNamespaceRequest) (*teov20220901.CreateEdgeKVNamespaceResponse, error) {
		assert.NotNil(t, request.ZoneId)
		assert.Equal(t, "zone-2o3h21ed2t68", *request.ZoneId)
		assert.NotNil(t, request.Namespace)
		assert.Equal(t, "my-namespace", *request.Namespace)
		assert.NotNil(t, request.Remark)
		assert.Equal(t, "test remark", *request.Remark)

		resp := teov20220901.NewCreateEdgeKVNamespaceResponse()
		resp.Response = &teov20220901.CreateEdgeKVNamespaceResponseParams{
			RequestId: ptrStringEdgeKVNamespace("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeEdgeKVNamespacesWithContext", func(_ context.Context, request *teov20220901.DescribeEdgeKVNamespacesRequest) (*teov20220901.DescribeEdgeKVNamespacesResponse, error) {
		assert.NotNil(t, request.ZoneId)
		assert.Equal(t, "zone-2o3h21ed2t68", *request.ZoneId)

		resp := teov20220901.NewDescribeEdgeKVNamespacesResponse()
		resp.Response = &teov20220901.DescribeEdgeKVNamespacesResponseParams{
			TotalCount: helper_int64(1),
			KVNamespaces: []*teov20220901.KVNamespace{
				{
					Namespace: ptrStringEdgeKVNamespace("my-namespace"),
					Remark:    ptrStringEdgeKVNamespace("test remark"),
				},
			},
			RequestId: ptrStringEdgeKVNamespace("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaEdgeKVNamespace()
	res := teo.ResourceTencentCloudTeoEdgeKVNamespace()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":   "zone-2o3h21ed2t68",
		"namespace": "my-namespace",
		"remark":    "test remark",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-2o3h21ed2t68#my-namespace", d.Id())
	assert.Equal(t, "test remark", d.Get("remark"))
}

// TestTeoEdgeKVNamespace_Read tests Read operation
func TestTeoEdgeKVNamespace_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaEdgeKVNamespace().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeEdgeKVNamespacesWithContext", func(_ context.Context, request *teov20220901.DescribeEdgeKVNamespacesRequest) (*teov20220901.DescribeEdgeKVNamespacesResponse, error) {
		assert.NotNil(t, request.ZoneId)
		assert.Equal(t, "zone-2o3h21ed2t68", *request.ZoneId)

		resp := teov20220901.NewDescribeEdgeKVNamespacesResponse()
		resp.Response = &teov20220901.DescribeEdgeKVNamespacesResponseParams{
			TotalCount: helper_int64(1),
			KVNamespaces: []*teov20220901.KVNamespace{
				{
					Namespace: ptrStringEdgeKVNamespace("my-namespace"),
					Remark:    ptrStringEdgeKVNamespace("test remark"),
				},
			},
			RequestId: ptrStringEdgeKVNamespace("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaEdgeKVNamespace()
	res := teo.ResourceTencentCloudTeoEdgeKVNamespace()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":   "zone-2o3h21ed2t68",
		"namespace": "my-namespace",
	})
	d.SetId("zone-2o3h21ed2t68#my-namespace")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-2o3h21ed2t68", d.Get("zone_id"))
	assert.Equal(t, "my-namespace", d.Get("namespace"))
	assert.Equal(t, "test remark", d.Get("remark"))
}

// TestTeoEdgeKVNamespace_ReadNotFound tests Read when resource is not found
func TestTeoEdgeKVNamespace_ReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaEdgeKVNamespace().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeEdgeKVNamespacesWithContext", func(_ context.Context, request *teov20220901.DescribeEdgeKVNamespacesRequest) (*teov20220901.DescribeEdgeKVNamespacesResponse, error) {
		resp := teov20220901.NewDescribeEdgeKVNamespacesResponse()
		resp.Response = &teov20220901.DescribeEdgeKVNamespacesResponseParams{
			TotalCount:   helper_int64(0),
			KVNamespaces: []*teov20220901.KVNamespace{},
			RequestId:    ptrStringEdgeKVNamespace("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaEdgeKVNamespace()
	res := teo.ResourceTencentCloudTeoEdgeKVNamespace()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":   "zone-2o3h21ed2t68",
		"namespace": "my-namespace",
	})
	d.SetId("zone-2o3h21ed2t68#my-namespace")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestTeoEdgeKVNamespace_Update tests Update operation
func TestTeoEdgeKVNamespace_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaEdgeKVNamespace().client, "UseTeoV20220901Client", teoClient)

	modifyCalled := false
	patches.ApplyMethodFunc(teoClient, "ModifyEdgeKVNamespaceWithContext", func(_ context.Context, request *teov20220901.ModifyEdgeKVNamespaceRequest) (*teov20220901.ModifyEdgeKVNamespaceResponse, error) {
		assert.NotNil(t, request.ZoneId)
		assert.Equal(t, "zone-2o3h21ed2t68", *request.ZoneId)
		assert.NotNil(t, request.Namespace)
		assert.Equal(t, "my-namespace", *request.Namespace)
		assert.NotNil(t, request.Remark)
		assert.Equal(t, "updated remark", *request.Remark)
		modifyCalled = true

		resp := teov20220901.NewModifyEdgeKVNamespaceResponse()
		resp.Response = &teov20220901.ModifyEdgeKVNamespaceResponseParams{
			RequestId: ptrStringEdgeKVNamespace("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeEdgeKVNamespacesWithContext", func(_ context.Context, request *teov20220901.DescribeEdgeKVNamespacesRequest) (*teov20220901.DescribeEdgeKVNamespacesResponse, error) {
		resp := teov20220901.NewDescribeEdgeKVNamespacesResponse()
		resp.Response = &teov20220901.DescribeEdgeKVNamespacesResponseParams{
			TotalCount: helper_int64(1),
			KVNamespaces: []*teov20220901.KVNamespace{
				{
					Namespace: ptrStringEdgeKVNamespace("my-namespace"),
					Remark:    ptrStringEdgeKVNamespace("updated remark"),
				},
			},
			RequestId: ptrStringEdgeKVNamespace("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaEdgeKVNamespace()
	res := teo.ResourceTencentCloudTeoEdgeKVNamespace()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":   "zone-2o3h21ed2t68",
		"namespace": "my-namespace",
		"remark":    "updated remark",
	})
	d.SetId("zone-2o3h21ed2t68#my-namespace")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.True(t, modifyCalled)
	assert.Equal(t, "updated remark", d.Get("remark"))
}

// TestTeoEdgeKVNamespace_Delete tests Delete operation
func TestTeoEdgeKVNamespace_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaEdgeKVNamespace().client, "UseTeoV20220901Client", teoClient)

	deleteCalled := false
	patches.ApplyMethodFunc(teoClient, "DeleteEdgeKVNamespaceWithContext", func(_ context.Context, request *teov20220901.DeleteEdgeKVNamespaceRequest) (*teov20220901.DeleteEdgeKVNamespaceResponse, error) {
		assert.NotNil(t, request.ZoneId)
		assert.Equal(t, "zone-2o3h21ed2t68", *request.ZoneId)
		assert.NotNil(t, request.Namespace)
		assert.Equal(t, "my-namespace", *request.Namespace)
		deleteCalled = true

		resp := teov20220901.NewDeleteEdgeKVNamespaceResponse()
		resp.Response = &teov20220901.DeleteEdgeKVNamespaceResponseParams{
			RequestId: ptrStringEdgeKVNamespace("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaEdgeKVNamespace()
	res := teo.ResourceTencentCloudTeoEdgeKVNamespace()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":   "zone-2o3h21ed2t68",
		"namespace": "my-namespace",
		"remark":    "test remark",
	})
	d.SetId("zone-2o3h21ed2t68#my-namespace")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
	assert.True(t, deleteCalled)
}
