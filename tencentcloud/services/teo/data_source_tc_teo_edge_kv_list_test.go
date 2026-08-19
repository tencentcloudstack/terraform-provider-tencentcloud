package teo_test

import (
	"context"
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// go test ./tencentcloud/services/teo/ -run "TestTeoEdgeKvListDataSource" -v -count=1 -gcflags="all=-l"

// TestTeoEdgeKvListDataSource_ReadSuccess tests successful read with a single page of keys and an empty cursor
func TestTeoEdgeKvListDataSource_ReadSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "EdgeKVListWithContext", func(_ context.Context, request *teov20220901.EdgeKVListRequest) (*teov20220901.EdgeKVListResponse, error) {
		assert.Equal(t, "zone-2qtuhspy7cr6", *request.ZoneId)
		assert.Equal(t, "default", *request.Namespace)
		assert.Equal(t, int64(1000), *request.Limit)
		resp := teov20220901.NewEdgeKVListResponse()
		resp.Response = &teov20220901.EdgeKVListResponseParams{
			Keys: []*string{
				ptrString("key1"),
				ptrString("key2"),
				ptrString("key3"),
			},
			Cursor:    ptrString(""),
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.DataSourceTencentCloudTeoEdgeKvList()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":   "zone-2qtuhspy7cr6",
		"namespace": "default",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	keys := d.Get("keys").([]interface{})
	assert.Len(t, keys, 3)
	assert.Equal(t, "key1", keys[0].(string))
	assert.Equal(t, "key2", keys[1].(string))
	assert.Equal(t, "key3", keys[2].(string))
}

// TestTeoEdgeKvListDataSource_Paginated tests that keys from multiple pages are accumulated
func TestTeoEdgeKvListDataSource_Paginated(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoClient", teoClient)

	var callCount int32
	patches.ApplyMethodFunc(teoClient, "EdgeKVListWithContext", func(_ context.Context, request *teov20220901.EdgeKVListRequest) (*teov20220901.EdgeKVListResponse, error) {
		n := atomic.AddInt32(&callCount, 1)
		if n == 1 {
			resp := teov20220901.NewEdgeKVListResponse()
			resp.Response = &teov20220901.EdgeKVListResponseParams{
				Keys: []*string{
					ptrString("page1-key1"),
					ptrString("page1-key2"),
				},
				Cursor:    ptrString("next-page-cursor"),
				RequestId: ptrString("fake-request-id-1"),
			}
			return resp, nil
		}
		// Second page returns remaining keys and an empty cursor to terminate the loop
		resp := teov20220901.NewEdgeKVListResponse()
		resp.Response = &teov20220901.EdgeKVListResponseParams{
			Keys: []*string{
				ptrString("page2-key1"),
			},
			Cursor:    ptrString(""),
			RequestId: ptrString("fake-request-id-2"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.DataSourceTencentCloudTeoEdgeKvList()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":   "zone-2qtuhspy7cr6",
		"namespace": "default",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	keys := d.Get("keys").([]interface{})
	assert.Len(t, keys, 3)
	assert.Equal(t, "page1-key1", keys[0].(string))
	assert.Equal(t, "page1-key2", keys[1].(string))
	assert.Equal(t, "page2-key1", keys[2].(string))
	assert.Equal(t, int32(2), atomic.LoadInt32(&callCount))
}

// TestTeoEdgeKvListDataSource_PrefixFilter tests that request.Prefix is set to the configured value
func TestTeoEdgeKvListDataSource_PrefixFilter(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "EdgeKVListWithContext", func(_ context.Context, request *teov20220901.EdgeKVListRequest) (*teov20220901.EdgeKVListResponse, error) {
		assert.NotNil(t, request.Prefix)
		assert.Equal(t, "user_", *request.Prefix)
		resp := teov20220901.NewEdgeKVListResponse()
		resp.Response = &teov20220901.EdgeKVListResponseParams{
			Keys: []*string{
				ptrString("user_1"),
				ptrString("user_2"),
			},
			Cursor:    ptrString(""),
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.DataSourceTencentCloudTeoEdgeKvList()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":   "zone-2qtuhspy7cr6",
		"namespace": "default",
		"prefix":    "user_",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	keys := d.Get("keys").([]interface{})
	assert.Len(t, keys, 2)
	assert.Equal(t, "user_1", keys[0].(string))
	assert.Equal(t, "user_2", keys[1].(string))
}

// TestTeoEdgeKvListDataSource_NilResponse tests that a nil response returns an error and does not clear state id
func TestTeoEdgeKvListDataSource_NilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "EdgeKVListWithContext", func(_ context.Context, request *teov20220901.EdgeKVListRequest) (*teov20220901.EdgeKVListResponse, error) {
		return nil, nil
	})

	meta := newMockMeta()
	res := teo.DataSourceTencentCloudTeoEdgeKvList()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":   "zone-2qtuhspy7cr6",
		"namespace": "default",
	})
	d.SetId("existing-state-id")

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "teo_edge_kv_list")
	// state id should NOT be cleared by the read; the NonRetryableError preserves it
	assert.Equal(t, "existing-state-id", d.Id())
}

// TestTeoEdgeKvListDataSource_APIError tests that an API error is propagated
func TestTeoEdgeKvListDataSource_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "EdgeKVListWithContext", func(_ context.Context, request *teov20220901.EdgeKVListRequest) (*teov20220901.EdgeKVListResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceUnavailable.NamespaceNotFound, Message=namespace not found")
	})

	meta := newMockMeta()
	res := teo.DataSourceTencentCloudTeoEdgeKvList()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":   "zone-2qtuhspy7cr6",
		"namespace": "missing",
	})

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "NamespaceNotFound")
}

// TestTeoEdgeKvListDataSource_Schema validates schema definition
func TestTeoEdgeKvListDataSource_Schema(t *testing.T) {
	res := teo.DataSourceTencentCloudTeoEdgeKvList()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Read)

	assert.Contains(t, res.Schema, "zone_id")
	assert.Contains(t, res.Schema, "namespace")
	assert.Contains(t, res.Schema, "prefix")
	assert.Contains(t, res.Schema, "cursor")
	assert.Contains(t, res.Schema, "keys")
	assert.Contains(t, res.Schema, "result_output_file")

	// zone_id is Required TypeString
	zoneId := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneId.Type)
	assert.True(t, zoneId.Required)

	// namespace is Required TypeString
	namespace := res.Schema["namespace"]
	assert.Equal(t, schema.TypeString, namespace.Type)
	assert.True(t, namespace.Required)

	// prefix is Optional TypeString
	prefix := res.Schema["prefix"]
	assert.Equal(t, schema.TypeString, prefix.Type)
	assert.True(t, prefix.Optional)

	// cursor is Optional and Computed TypeString
	cursor := res.Schema["cursor"]
	assert.Equal(t, schema.TypeString, cursor.Type)
	assert.True(t, cursor.Optional)
	assert.True(t, cursor.Computed)

	// keys is Computed TypeList of TypeString
	keys := res.Schema["keys"]
	assert.Equal(t, schema.TypeList, keys.Type)
	assert.True(t, keys.Computed)

	// result_output_file is Optional TypeString
	outputFile := res.Schema["result_output_file"]
	assert.Equal(t, schema.TypeString, outputFile.Type)
	assert.True(t, outputFile.Optional)
}
