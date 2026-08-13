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

type mockMetaComponentBinding struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaComponentBinding) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaComponentBinding{}

func newMockMetaComponentBinding() *mockMetaComponentBinding {
	return &mockMetaComponentBinding{client: &connectivity.TencentCloudClient{}}
}

func ptrStringCB(s string) *string {
	return &s
}

func ptrInt64CB(i int64) *int64 {
	return &i
}

func TestTeoFunctionComponentBinding_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaComponentBinding().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyFunctionComponentBindingsWithContext", func(_ context.Context, request *teov20220901.ModifyFunctionComponentBindingsRequest) (*teov20220901.ModifyFunctionComponentBindingsResponse, error) {
		assert.NotNil(t, request.ZoneId)
		assert.Equal(t, "zone-test123", *request.ZoneId)
		assert.NotNil(t, request.FunctionId)
		assert.Equal(t, "ef-test456", *request.FunctionId)
		assert.NotNil(t, request.Operation)
		assert.Equal(t, "rebind", *request.Operation)
		assert.Equal(t, 1, len(request.FunctionComponentBindings))
		assert.Equal(t, "kv_namespace", *request.FunctionComponentBindings[0].Type)
		assert.Equal(t, "MY_KV", *request.FunctionComponentBindings[0].VariableName)
		assert.NotNil(t, request.FunctionComponentBindings[0].KVNamespaceParameters)
		assert.Equal(t, "zone-test123", *request.FunctionComponentBindings[0].KVNamespaceParameters.ZoneId)
		assert.Equal(t, "my-namespace", *request.FunctionComponentBindings[0].KVNamespaceParameters.Namespace)

		resp := teov20220901.NewModifyFunctionComponentBindingsResponse()
		resp.Response = &teov20220901.ModifyFunctionComponentBindingsResponseParams{
			RequestId: ptrStringCB("fake-request-id-create"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionComponentBindingsWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionComponentBindingsRequest) (*teov20220901.DescribeFunctionComponentBindingsResponse, error) {
		resp := teov20220901.NewDescribeFunctionComponentBindingsResponse()
		resp.Response = &teov20220901.DescribeFunctionComponentBindingsResponseParams{
			TotalCount: ptrInt64CB(1),
			FunctionComponentBindings: []*teov20220901.FunctionComponentBinding{
				{
					Type:         ptrStringCB("kv_namespace"),
					VariableName: ptrStringCB("MY_KV"),
					KVNamespaceParameters: &teov20220901.KVNamespaceParameters{
						ZoneId:    ptrStringCB("zone-test123"),
						Namespace: ptrStringCB("my-namespace"),
					},
				},
			},
			RequestId: ptrStringCB("fake-request-id-read"),
		}
		return resp, nil
	})

	meta := newMockMetaComponentBinding()
	res := teo.ResourceTencentCloudTeoFunctionComponentBinding()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-test123",
		"function_id": "ef-test456",
		"function_component_bindings": []interface{}{
			map[string]interface{}{
				"type":          "kv_namespace",
				"variable_name": "MY_KV",
				"kv_namespace_parameters": []interface{}{
					map[string]interface{}{
						"zone_id":   "zone-test123",
						"namespace": "my-namespace",
					},
				},
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-test123#ef-test456", d.Id())
}

func TestTeoFunctionComponentBinding_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaComponentBinding().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionComponentBindingsWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionComponentBindingsRequest) (*teov20220901.DescribeFunctionComponentBindingsResponse, error) {
		assert.NotNil(t, request.ZoneId)
		assert.Equal(t, "zone-test123", *request.ZoneId)
		assert.NotNil(t, request.FunctionId)
		assert.Equal(t, "ef-test456", *request.FunctionId)

		resp := teov20220901.NewDescribeFunctionComponentBindingsResponse()
		resp.Response = &teov20220901.DescribeFunctionComponentBindingsResponseParams{
			TotalCount: ptrInt64CB(1),
			FunctionComponentBindings: []*teov20220901.FunctionComponentBinding{
				{
					Type:         ptrStringCB("kv_namespace"),
					VariableName: ptrStringCB("MY_KV"),
					KVNamespaceParameters: &teov20220901.KVNamespaceParameters{
						ZoneId:    ptrStringCB("zone-test123"),
						Namespace: ptrStringCB("my-namespace"),
					},
				},
			},
			RequestId: ptrStringCB("fake-request-id-read"),
		}
		return resp, nil
	})

	meta := newMockMetaComponentBinding()
	res := teo.ResourceTencentCloudTeoFunctionComponentBinding()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":                     "zone-test123",
		"function_id":                 "ef-test456",
		"function_component_bindings": []interface{}{},
	})
	d.SetId("zone-test123#ef-test456")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-test123#ef-test456", d.Id())
	assert.Equal(t, "zone-test123", d.Get("zone_id"))
	assert.Equal(t, "ef-test456", d.Get("function_id"))
}

func TestTeoFunctionComponentBinding_Read_Empty(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaComponentBinding().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionComponentBindingsWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionComponentBindingsRequest) (*teov20220901.DescribeFunctionComponentBindingsResponse, error) {
		resp := teov20220901.NewDescribeFunctionComponentBindingsResponse()
		resp.Response = &teov20220901.DescribeFunctionComponentBindingsResponseParams{
			TotalCount:                ptrInt64CB(0),
			FunctionComponentBindings: []*teov20220901.FunctionComponentBinding{},
			RequestId:                 ptrStringCB("fake-request-id-empty"),
		}
		return resp, nil
	})

	meta := newMockMetaComponentBinding()
	res := teo.ResourceTencentCloudTeoFunctionComponentBinding()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":                     "zone-test123",
		"function_id":                 "ef-test456",
		"function_component_bindings": []interface{}{},
	})
	d.SetId("zone-test123#ef-test456")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-test123#ef-test456", d.Id())
}

func TestTeoFunctionComponentBinding_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaComponentBinding().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyFunctionComponentBindingsWithContext", func(_ context.Context, request *teov20220901.ModifyFunctionComponentBindingsRequest) (*teov20220901.ModifyFunctionComponentBindingsResponse, error) {
		assert.NotNil(t, request.ZoneId)
		assert.Equal(t, "zone-test123", *request.ZoneId)
		assert.NotNil(t, request.FunctionId)
		assert.Equal(t, "ef-test456", *request.FunctionId)
		assert.NotNil(t, request.Operation)
		assert.Equal(t, "rebind", *request.Operation)

		resp := teov20220901.NewModifyFunctionComponentBindingsResponse()
		resp.Response = &teov20220901.ModifyFunctionComponentBindingsResponseParams{
			RequestId: ptrStringCB("fake-request-id-update"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionComponentBindingsWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionComponentBindingsRequest) (*teov20220901.DescribeFunctionComponentBindingsResponse, error) {
		resp := teov20220901.NewDescribeFunctionComponentBindingsResponse()
		resp.Response = &teov20220901.DescribeFunctionComponentBindingsResponseParams{
			TotalCount: ptrInt64CB(1),
			FunctionComponentBindings: []*teov20220901.FunctionComponentBinding{
				{
					Type:         ptrStringCB("kv_namespace"),
					VariableName: ptrStringCB("MY_KV_UPDATED"),
					KVNamespaceParameters: &teov20220901.KVNamespaceParameters{
						ZoneId:    ptrStringCB("zone-test123"),
						Namespace: ptrStringCB("my-namespace-updated"),
					},
				},
			},
			RequestId: ptrStringCB("fake-request-id-read"),
		}
		return resp, nil
	})

	meta := newMockMetaComponentBinding()
	res := teo.ResourceTencentCloudTeoFunctionComponentBinding()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-test123",
		"function_id": "ef-test456",
		"function_component_bindings": []interface{}{
			map[string]interface{}{
				"type":          "kv_namespace",
				"variable_name": "MY_KV_UPDATED",
				"kv_namespace_parameters": []interface{}{
					map[string]interface{}{
						"zone_id":   "zone-test123",
						"namespace": "my-namespace-updated",
					},
				},
			},
		},
	})
	d.SetId("zone-test123#ef-test456")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

func TestTeoFunctionComponentBinding_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaComponentBinding().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyFunctionComponentBindingsWithContext", func(_ context.Context, request *teov20220901.ModifyFunctionComponentBindingsRequest) (*teov20220901.ModifyFunctionComponentBindingsResponse, error) {
		assert.NotNil(t, request.ZoneId)
		assert.Equal(t, "zone-test123", *request.ZoneId)
		assert.NotNil(t, request.FunctionId)
		assert.Equal(t, "ef-test456", *request.FunctionId)
		assert.NotNil(t, request.Operation)
		assert.Equal(t, "rebind", *request.Operation)
		assert.Equal(t, 0, len(request.FunctionComponentBindings))

		resp := teov20220901.NewModifyFunctionComponentBindingsResponse()
		resp.Response = &teov20220901.ModifyFunctionComponentBindingsResponseParams{
			RequestId: ptrStringCB("fake-request-id-delete"),
		}
		return resp, nil
	})

	meta := newMockMetaComponentBinding()
	res := teo.ResourceTencentCloudTeoFunctionComponentBinding()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":                     "zone-test123",
		"function_id":                 "ef-test456",
		"function_component_bindings": []interface{}{},
	})
	d.SetId("zone-test123#ef-test456")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}
