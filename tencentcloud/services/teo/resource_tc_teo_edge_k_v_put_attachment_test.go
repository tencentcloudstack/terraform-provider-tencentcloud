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

// mockMetaEdgeKVPut implements tccommon.ProviderMeta
type mockMetaEdgeKVPut struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaEdgeKVPut) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaEdgeKVPut{}

func newMockMetaEdgeKVPut() *mockMetaEdgeKVPut {
	return &mockMetaEdgeKVPut{client: &connectivity.TencentCloudClient{}}
}

func ptrStringEdgeKVPut(s string) *string {
	return &s
}

// go test ./tencentcloud/services/teo/ -run "TestTeoEdgeKVPut" -v -count=1 -gcflags="all=-l"

// TestTeoEdgeKVPut_Create tests Create operation
func TestTeoEdgeKVPut_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaEdgeKVPut().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "EdgeKVPutWithContext", func(_ context.Context, request *teov20220901.EdgeKVPutRequest) (*teov20220901.EdgeKVPutResponse, error) {
		assert.NotNil(t, request.ZoneId)
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.NotNil(t, request.Namespace)
		assert.Equal(t, "test-namespace", *request.Namespace)
		assert.NotNil(t, request.Key)
		assert.Equal(t, "test-key", *request.Key)
		assert.NotNil(t, request.Value)
		assert.Equal(t, "test-value", *request.Value)

		resp := teov20220901.NewEdgeKVPutResponse()
		resp.Response = &teov20220901.EdgeKVPutResponseParams{
			RequestId: ptrStringEdgeKVPut("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "EdgeKVGetWithContext", func(_ context.Context, request *teov20220901.EdgeKVGetRequest) (*teov20220901.EdgeKVGetResponse, error) {
		resp := teov20220901.NewEdgeKVGetResponse()
		resp.Response = &teov20220901.EdgeKVGetResponseParams{
			Data: []*teov20220901.KeyValuePair{
				{
					Key:   ptrStringEdgeKVPut("test-key"),
					Value: ptrStringEdgeKVPut("test-value"),
				},
			},
			RequestId: ptrStringEdgeKVPut("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaEdgeKVPut()
	res := teo.ResourceTencentCloudTeoEdgeKVPut()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":   "zone-12345678",
		"namespace": "test-namespace",
		"key":       "test-key",
		"value":     "test-value",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-12345678"+tccommon.FILED_SP+"test-namespace"+tccommon.FILED_SP+"test-key", d.Id())
}

// TestTeoEdgeKVPut_CreateWithExpiration tests Create with expiration parameters
func TestTeoEdgeKVPut_CreateWithExpiration(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaEdgeKVPut().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "EdgeKVPutWithContext", func(_ context.Context, request *teov20220901.EdgeKVPutRequest) (*teov20220901.EdgeKVPutResponse, error) {
		assert.NotNil(t, request.ZoneId)
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.NotNil(t, request.ExpirationTTL)
		assert.Equal(t, int64(3600), *request.ExpirationTTL)

		resp := teov20220901.NewEdgeKVPutResponse()
		resp.Response = &teov20220901.EdgeKVPutResponseParams{
			RequestId: ptrStringEdgeKVPut("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "EdgeKVGetWithContext", func(_ context.Context, request *teov20220901.EdgeKVGetRequest) (*teov20220901.EdgeKVGetResponse, error) {
		resp := teov20220901.NewEdgeKVGetResponse()
		resp.Response = &teov20220901.EdgeKVGetResponseParams{
			Data: []*teov20220901.KeyValuePair{
				{
					Key:   ptrStringEdgeKVPut("test-key"),
					Value: ptrStringEdgeKVPut("test-value"),
				},
			},
			RequestId: ptrStringEdgeKVPut("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaEdgeKVPut()
	res := teo.ResourceTencentCloudTeoEdgeKVPut()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":        "zone-12345678",
		"namespace":      "test-namespace",
		"key":            "test-key",
		"value":          "test-value",
		"expiration_ttl": 3600,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-12345678"+tccommon.FILED_SP+"test-namespace"+tccommon.FILED_SP+"test-key", d.Id())
}

// TestTeoEdgeKVPut_Read tests Read operation
func TestTeoEdgeKVPut_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaEdgeKVPut().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "EdgeKVGetWithContext", func(_ context.Context, request *teov20220901.EdgeKVGetRequest) (*teov20220901.EdgeKVGetResponse, error) {
		assert.NotNil(t, request.ZoneId)
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.NotNil(t, request.Namespace)
		assert.Equal(t, "test-namespace", *request.Namespace)
		assert.Equal(t, 1, len(request.Keys))
		assert.Equal(t, "test-key", *request.Keys[0])

		resp := teov20220901.NewEdgeKVGetResponse()
		resp.Response = &teov20220901.EdgeKVGetResponseParams{
			Data: []*teov20220901.KeyValuePair{
				{
					Key:   ptrStringEdgeKVPut("test-key"),
					Value: ptrStringEdgeKVPut("test-value-read"),
				},
			},
			RequestId: ptrStringEdgeKVPut("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaEdgeKVPut()
	res := teo.ResourceTencentCloudTeoEdgeKVPut()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":   "zone-12345678",
		"namespace": "test-namespace",
		"key":       "test-key",
		"value":     "test-value",
	})
	d.SetId("zone-12345678" + tccommon.FILED_SP + "test-namespace" + tccommon.FILED_SP + "test-key")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "test-value-read", d.Get("value"))
	assert.Equal(t, "zone-12345678", d.Get("zone_id"))
	assert.Equal(t, "test-namespace", d.Get("namespace"))
	assert.Equal(t, "test-key", d.Get("key"))
}

// TestTeoEdgeKVPut_ReadNotFound tests Read when key does not exist
func TestTeoEdgeKVPut_ReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaEdgeKVPut().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "EdgeKVGetWithContext", func(_ context.Context, request *teov20220901.EdgeKVGetRequest) (*teov20220901.EdgeKVGetResponse, error) {
		resp := teov20220901.NewEdgeKVGetResponse()
		resp.Response = &teov20220901.EdgeKVGetResponseParams{
			Data:      []*teov20220901.KeyValuePair{},
			RequestId: ptrStringEdgeKVPut("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaEdgeKVPut()
	res := teo.ResourceTencentCloudTeoEdgeKVPut()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":   "zone-12345678",
		"namespace": "test-namespace",
		"key":       "test-key",
		"value":     "test-value",
	})
	d.SetId("zone-12345678" + tccommon.FILED_SP + "test-namespace" + tccommon.FILED_SP + "test-key")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestTeoEdgeKVPut_Delete tests Delete operation
func TestTeoEdgeKVPut_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaEdgeKVPut().client, "UseTeoV20220901Client", teoClient)

	deleteCalled := false
	patches.ApplyMethodFunc(teoClient, "EdgeKVDeleteWithContext", func(_ context.Context, request *teov20220901.EdgeKVDeleteRequest) (*teov20220901.EdgeKVDeleteResponse, error) {
		assert.NotNil(t, request.ZoneId)
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.NotNil(t, request.Namespace)
		assert.Equal(t, "test-namespace", *request.Namespace)
		assert.Equal(t, 1, len(request.Keys))
		assert.Equal(t, "test-key", *request.Keys[0])
		deleteCalled = true

		resp := teov20220901.NewEdgeKVDeleteResponse()
		resp.Response = &teov20220901.EdgeKVDeleteResponseParams{
			RequestId: ptrStringEdgeKVPut("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaEdgeKVPut()
	res := teo.ResourceTencentCloudTeoEdgeKVPut()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":   "zone-12345678",
		"namespace": "test-namespace",
		"key":       "test-key",
		"value":     "test-value",
	})
	d.SetId("zone-12345678" + tccommon.FILED_SP + "test-namespace" + tccommon.FILED_SP + "test-key")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
	assert.True(t, deleteCalled)
}
