package teo_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

type mockMetaFunctionComponentBinding struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaFunctionComponentBinding) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaFunctionComponentBinding{}

func newMockMetaFunctionComponentBinding() *mockMetaFunctionComponentBinding {
	return &mockMetaFunctionComponentBinding{client: &connectivity.TencentCloudClient{}}
}

func ptrStringFCB(s string) *string {
	return &s
}

func ptrInt64FCB(n int64) *int64 {
	return &n
}

// go test ./tencentcloud/services/teo/ -run "TestFunctionComponentBinding" -v -count=1 -gcflags="all=-l"

// TestFunctionComponentBinding_Create_Success tests Create calls API and sets composite ID
func TestFunctionComponentBinding_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaFunctionComponentBinding().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyFunctionComponentBindingsWithContext", func(_ context.Context, request *teov20220901.ModifyFunctionComponentBindingsRequest) (*teov20220901.ModifyFunctionComponentBindingsResponse, error) {
		resp := teov20220901.NewModifyFunctionComponentBindingsResponse()
		resp.Response = &teov20220901.ModifyFunctionComponentBindingsResponseParams{
			RequestId: ptrStringFCB("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionComponentBindingsWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionComponentBindingsRequest) (*teov20220901.DescribeFunctionComponentBindingsResponse, error) {
		resp := teov20220901.NewDescribeFunctionComponentBindingsResponse()
		resp.Response = &teov20220901.DescribeFunctionComponentBindingsResponseParams{
			TotalCount: ptrInt64FCB(1),
			FunctionComponentBindings: []*teov20220901.FunctionComponentBinding{
				{
					Type:         ptrStringFCB("kv_namespace"),
					VariableName: ptrStringFCB("MY_KV"),
					KVNamespaceParameters: &teov20220901.KVNamespaceParameters{
						ZoneId:    ptrStringFCB("zone-abc123"),
						Namespace: ptrStringFCB("my-namespace"),
					},
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaFunctionComponentBinding()
	res := teo.ResourceTencentCloudTeoFunctionComponentBinding()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-abc123",
		"function_id": "func-xyz789",
		"function_component_bindings": []interface{}{
			map[string]interface{}{
				"type":          "kv_namespace",
				"variable_name": "MY_KV",
				"kv_namespace_parameters": []interface{}{
					map[string]interface{}{
						"zone_id":   "zone-abc123",
						"namespace": "my-namespace",
					},
				},
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-abc123#func-xyz789", d.Id())
}

// TestFunctionComponentBinding_Create_APIError tests Create handles API error
func TestFunctionComponentBinding_Create_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaFunctionComponentBinding().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyFunctionComponentBindingsWithContext", func(_ context.Context, request *teov20220901.ModifyFunctionComponentBindingsRequest) (*teov20220901.ModifyFunctionComponentBindingsResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid zone_id")
	})

	meta := newMockMetaFunctionComponentBinding()
	res := teo.ResourceTencentCloudTeoFunctionComponentBinding()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-invalid",
		"function_id": "func-xyz789",
		"function_component_bindings": []interface{}{
			map[string]interface{}{
				"type":          "kv_namespace",
				"variable_name": "MY_KV",
				"kv_namespace_parameters": []interface{}{
					map[string]interface{}{
						"zone_id":   "zone-invalid",
						"namespace": "my-namespace",
					},
				},
			},
		},
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestFunctionComponentBinding_Read_Success tests Read retrieves binding data
func TestFunctionComponentBinding_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaFunctionComponentBinding().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionComponentBindingsWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionComponentBindingsRequest) (*teov20220901.DescribeFunctionComponentBindingsResponse, error) {
		resp := teov20220901.NewDescribeFunctionComponentBindingsResponse()
		resp.Response = &teov20220901.DescribeFunctionComponentBindingsResponseParams{
			TotalCount: ptrInt64FCB(1),
			FunctionComponentBindings: []*teov20220901.FunctionComponentBinding{
				{
					Type:         ptrStringFCB("kv_namespace"),
					VariableName: ptrStringFCB("MY_KV"),
					KVNamespaceParameters: &teov20220901.KVNamespaceParameters{
						ZoneId:    ptrStringFCB("zone-abc123"),
						Namespace: ptrStringFCB("my-namespace"),
					},
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaFunctionComponentBinding()
	res := teo.ResourceTencentCloudTeoFunctionComponentBinding()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":                     "zone-abc123",
		"function_id":                 "func-xyz789",
		"function_component_bindings": []interface{}{},
	})
	d.SetId("zone-abc123#func-xyz789")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-abc123", d.Get("zone_id"))
	assert.Equal(t, "func-xyz789", d.Get("function_id"))
	bindings := d.Get("function_component_bindings").([]interface{})
	assert.Equal(t, 1, len(bindings))
}

// TestFunctionComponentBinding_Read_EmptyBindings tests Read with empty bindings
func TestFunctionComponentBinding_Read_EmptyBindings(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaFunctionComponentBinding().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionComponentBindingsWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionComponentBindingsRequest) (*teov20220901.DescribeFunctionComponentBindingsResponse, error) {
		resp := teov20220901.NewDescribeFunctionComponentBindingsResponse()
		resp.Response = &teov20220901.DescribeFunctionComponentBindingsResponseParams{
			TotalCount:                ptrInt64FCB(0),
			FunctionComponentBindings: []*teov20220901.FunctionComponentBinding{},
		}
		return resp, nil
	})

	meta := newMockMetaFunctionComponentBinding()
	res := teo.ResourceTencentCloudTeoFunctionComponentBinding()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":                     "zone-abc123",
		"function_id":                 "func-xyz789",
		"function_component_bindings": []interface{}{},
	})
	d.SetId("zone-abc123#func-xyz789")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	bindings := d.Get("function_component_bindings").([]interface{})
	assert.Equal(t, 0, len(bindings))
}

// TestFunctionComponentBinding_Update_Success tests Update calls API and then Read
func TestFunctionComponentBinding_Update_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaFunctionComponentBinding().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyFunctionComponentBindingsWithContext", func(_ context.Context, request *teov20220901.ModifyFunctionComponentBindingsRequest) (*teov20220901.ModifyFunctionComponentBindingsResponse, error) {
		resp := teov20220901.NewModifyFunctionComponentBindingsResponse()
		resp.Response = &teov20220901.ModifyFunctionComponentBindingsResponseParams{
			RequestId: ptrStringFCB("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionComponentBindingsWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionComponentBindingsRequest) (*teov20220901.DescribeFunctionComponentBindingsResponse, error) {
		resp := teov20220901.NewDescribeFunctionComponentBindingsResponse()
		resp.Response = &teov20220901.DescribeFunctionComponentBindingsResponseParams{
			TotalCount: ptrInt64FCB(1),
			FunctionComponentBindings: []*teov20220901.FunctionComponentBinding{
				{
					Type:         ptrStringFCB("kv_namespace"),
					VariableName: ptrStringFCB("MY_KV_UPDATED"),
					KVNamespaceParameters: &teov20220901.KVNamespaceParameters{
						ZoneId:    ptrStringFCB("zone-abc123"),
						Namespace: ptrStringFCB("my-namespace-updated"),
					},
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaFunctionComponentBinding()
	res := teo.ResourceTencentCloudTeoFunctionComponentBinding()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-abc123",
		"function_id": "func-xyz789",
		"function_component_bindings": []interface{}{
			map[string]interface{}{
				"type":          "kv_namespace",
				"variable_name": "MY_KV_UPDATED",
				"kv_namespace_parameters": []interface{}{
					map[string]interface{}{
						"zone_id":   "zone-abc123",
						"namespace": "my-namespace-updated",
					},
				},
			},
		},
	})
	d.SetId("zone-abc123#func-xyz789")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

// TestFunctionComponentBinding_Delete_Success tests Delete calls API with empty bindings
func TestFunctionComponentBinding_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaFunctionComponentBinding().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyFunctionComponentBindingsWithContext", func(_ context.Context, request *teov20220901.ModifyFunctionComponentBindingsRequest) (*teov20220901.ModifyFunctionComponentBindingsResponse, error) {
		assert.Equal(t, "rebind", *request.Operation)
		assert.Equal(t, 0, len(request.FunctionComponentBindings))
		resp := teov20220901.NewModifyFunctionComponentBindingsResponse()
		resp.Response = &teov20220901.ModifyFunctionComponentBindingsResponseParams{
			RequestId: ptrStringFCB("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaFunctionComponentBinding()
	res := teo.ResourceTencentCloudTeoFunctionComponentBinding()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":                     "zone-abc123",
		"function_id":                 "func-xyz789",
		"function_component_bindings": []interface{}{},
	})
	d.SetId("zone-abc123#func-xyz789")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestFunctionComponentBinding_Delete_APIError tests Delete handles API error
func TestFunctionComponentBinding_Delete_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaFunctionComponentBinding().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyFunctionComponentBindingsWithContext", func(_ context.Context, request *teov20220901.ModifyFunctionComponentBindingsRequest) (*teov20220901.ModifyFunctionComponentBindingsResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=Function not found")
	})

	meta := newMockMetaFunctionComponentBinding()
	res := teo.ResourceTencentCloudTeoFunctionComponentBinding()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":                     "zone-abc123",
		"function_id":                 "func-xyz789",
		"function_component_bindings": []interface{}{},
	})
	d.SetId("zone-abc123#func-xyz789")

	err := res.Delete(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
}

// TestFunctionComponentBinding_Schema validates schema definition
func TestFunctionComponentBinding_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoFunctionComponentBinding()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.NotNil(t, res.Update)
	assert.NotNil(t, res.Delete)
	assert.NotNil(t, res.Importer)

	assert.Contains(t, res.Schema, "zone_id")
	zoneId := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneId.Type)
	assert.True(t, zoneId.Required)
	assert.True(t, zoneId.ForceNew)

	assert.Contains(t, res.Schema, "function_id")
	functionId := res.Schema["function_id"]
	assert.Equal(t, schema.TypeString, functionId.Type)
	assert.True(t, functionId.Required)
	assert.True(t, functionId.ForceNew)

	assert.Contains(t, res.Schema, "function_component_bindings")
	bindings := res.Schema["function_component_bindings"]
	assert.Equal(t, schema.TypeList, bindings.Type)
	assert.True(t, bindings.Required)
	assert.False(t, bindings.ForceNew)
}
