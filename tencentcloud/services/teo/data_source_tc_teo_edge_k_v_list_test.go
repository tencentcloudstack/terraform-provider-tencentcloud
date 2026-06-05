package teo_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

type mockMetaForEdgeKVList struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForEdgeKVList) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForEdgeKVList{}

func newMockMetaForEdgeKVList() *mockMetaForEdgeKVList {
	return &mockMetaForEdgeKVList{client: &connectivity.TencentCloudClient{}}
}

func ptrStrEdgeKVList(s string) *string {
	return &s
}

// go test ./tencentcloud/services/teo/ -run "TestEdgeKVList" -v -count=1 -gcflags="all=-l"

// TestEdgeKVList_Read_Success tests successful read of edge kv list
func TestEdgeKVList_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForEdgeKVList().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "EdgeKVList", func(request *teov20220901.EdgeKVListRequest) (*teov20220901.EdgeKVListResponse, error) {
		resp := teov20220901.NewEdgeKVListResponse()
		resp.Response = &teov20220901.EdgeKVListResponseParams{
			Keys: []*string{
				ptrStrEdgeKVList("key1"),
				ptrStrEdgeKVList("key2"),
				ptrStrEdgeKVList("key3"),
			},
			Cursor:    ptrStrEdgeKVList("next_cursor_value"),
			RequestId: ptrStrEdgeKVList("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForEdgeKVList()
	res := teo.DataSourceTencentCloudTeoEdgeKVList()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":   "zone-test123",
		"namespace": "test_namespace",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	keys := d.Get("keys").([]interface{})
	assert.Equal(t, 3, len(keys))
	assert.Equal(t, "key1", keys[0].(string))
	assert.Equal(t, "key2", keys[1].(string))
	assert.Equal(t, "key3", keys[2].(string))
	assert.Equal(t, "next_cursor_value", d.Get("cursor").(string))
}

// TestEdgeKVList_Read_WithPrefix tests read with prefix filter
func TestEdgeKVList_Read_WithPrefix(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForEdgeKVList().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "EdgeKVList", func(request *teov20220901.EdgeKVListRequest) (*teov20220901.EdgeKVListResponse, error) {
		assert.NotNil(t, request.Prefix)
		assert.Equal(t, "config_", *request.Prefix)
		assert.NotNil(t, request.Limit)
		assert.Equal(t, int64(1000), *request.Limit)

		resp := teov20220901.NewEdgeKVListResponse()
		resp.Response = &teov20220901.EdgeKVListResponseParams{
			Keys: []*string{
				ptrStrEdgeKVList("config_a"),
				ptrStrEdgeKVList("config_b"),
			},
			Cursor:    ptrStrEdgeKVList(""),
			RequestId: ptrStrEdgeKVList("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForEdgeKVList()
	res := teo.DataSourceTencentCloudTeoEdgeKVList()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":   "zone-test123",
		"namespace": "test_namespace",
		"prefix":    "config_",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	keys := d.Get("keys").([]interface{})
	assert.Equal(t, 2, len(keys))
	assert.Equal(t, "config_a", keys[0].(string))
	assert.Equal(t, "config_b", keys[1].(string))
}

// TestEdgeKVList_Read_EmptyResult tests read with empty result
func TestEdgeKVList_Read_EmptyResult(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForEdgeKVList().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "EdgeKVList", func(request *teov20220901.EdgeKVListRequest) (*teov20220901.EdgeKVListResponse, error) {
		resp := teov20220901.NewEdgeKVListResponse()
		resp.Response = &teov20220901.EdgeKVListResponseParams{
			Keys:      []*string{},
			Cursor:    ptrStrEdgeKVList(""),
			RequestId: ptrStrEdgeKVList("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForEdgeKVList()
	res := teo.DataSourceTencentCloudTeoEdgeKVList()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":   "zone-test123",
		"namespace": "test_namespace",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	keys := d.Get("keys").([]interface{})
	assert.Equal(t, 0, len(keys))
}

// TestEdgeKVList_Read_WithCursor tests read with cursor for pagination
func TestEdgeKVList_Read_WithCursor(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForEdgeKVList().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "EdgeKVList", func(request *teov20220901.EdgeKVListRequest) (*teov20220901.EdgeKVListResponse, error) {
		assert.NotNil(t, request.Cursor)
		assert.Equal(t, "prev_cursor_value", *request.Cursor)

		resp := teov20220901.NewEdgeKVListResponse()
		resp.Response = &teov20220901.EdgeKVListResponseParams{
			Keys: []*string{
				ptrStrEdgeKVList("key4"),
				ptrStrEdgeKVList("key5"),
			},
			Cursor:    ptrStrEdgeKVList("next_cursor_value_2"),
			RequestId: ptrStrEdgeKVList("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForEdgeKVList()
	res := teo.DataSourceTencentCloudTeoEdgeKVList()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":   "zone-test123",
		"namespace": "test_namespace",
		"cursor":    "prev_cursor_value",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	keys := d.Get("keys").([]interface{})
	assert.Equal(t, 2, len(keys))
	assert.Equal(t, "key4", keys[0].(string))
	assert.Equal(t, "key5", keys[1].(string))
	assert.Equal(t, "next_cursor_value_2", d.Get("cursor").(string))
}
