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

// mockMetaForFunctionComponentBinding implements tccommon.ProviderMeta
type mockMetaForFunctionComponentBinding struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForFunctionComponentBinding) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForFunctionComponentBinding{}

func newMockMetaForFunctionComponentBinding() *mockMetaForFunctionComponentBinding {
	return &mockMetaForFunctionComponentBinding{client: &connectivity.TencentCloudClient{}}
}

func ptrStrFCB(s string) *string {
	return &s
}

func ptrInt64FCB(i int64) *int64 {
	return &i
}

// go test ./tencentcloud/services/teo/ -run "TestFunctionComponentBinding" -v -count=1 -gcflags="all=-l"

// TestFunctionComponentBinding_Read_Success tests Read populates state from API
func TestFunctionComponentBinding_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForFunctionComponentBinding().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionComponentBindingsWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionComponentBindingsRequest) (*teov20220901.DescribeFunctionComponentBindingsResponse, error) {
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.Equal(t, "ef-abcdefgh", *request.FunctionId)
		assert.Equal(t, int64(0), *request.Offset)
		assert.Equal(t, int64(1000), *request.Limit)
		resp := teov20220901.NewDescribeFunctionComponentBindingsResponse()
		resp.Response = &teov20220901.DescribeFunctionComponentBindingsResponseParams{
			TotalCount: ptrInt64FCB(2),
			FunctionComponentBindings: []*teov20220901.FunctionComponentBinding{
				{
					Type:         ptrStrFCB("kv_namespace"),
					VariableName: ptrStrFCB("MY_KV"),
					KVNamespaceParameters: &teov20220901.KVNamespaceParameters{
						ZoneId:    ptrStrFCB("zone-12345678"),
						Namespace: ptrStrFCB("my-namespace"),
					},
				},
				{
					Type:         ptrStrFCB("kv_namespace"),
					VariableName: ptrStrFCB("MY_KV_2"),
					KVNamespaceParameters: &teov20220901.KVNamespaceParameters{
						ZoneId:    ptrStrFCB("zone-12345678"),
						Namespace: ptrStrFCB("my-namespace-2"),
					},
				},
			},
			RequestId: ptrStrFCB("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForFunctionComponentBinding()
	res := teo.ResourceTencentCloudTeoFunctionComponentBinding()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":                     "zone-12345678",
		"function_id":                 "ef-abcdefgh",
		"function_component_bindings": []interface{}{},
	})
	d.SetId("zone-12345678" + tccommon.FILED_SP + "ef-abcdefgh")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-12345678", d.Get("zone_id"))
	assert.Equal(t, "ef-abcdefgh", d.Get("function_id"))

	bindings := d.Get("function_component_bindings").([]interface{})
	assert.Equal(t, 2, len(bindings))

	binding0 := bindings[0].(map[string]interface{})
	assert.Equal(t, "kv_namespace", binding0["type"])
	assert.Equal(t, "MY_KV", binding0["variable_name"])
	kvParams0 := binding0["kv_namespace_parameters"].([]interface{})
	assert.Equal(t, 1, len(kvParams0))
	kvParamsMap0 := kvParams0[0].(map[string]interface{})
	assert.Equal(t, "zone-12345678", kvParamsMap0["zone_id"])
	assert.Equal(t, "my-namespace", kvParamsMap0["namespace"])

	binding1 := bindings[1].(map[string]interface{})
	assert.Equal(t, "kv_namespace", binding1["type"])
	assert.Equal(t, "MY_KV_2", binding1["variable_name"])
}

// TestFunctionComponentBinding_Read_EmptyBindings tests Read with empty binding list
func TestFunctionComponentBinding_Read_EmptyBindings(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForFunctionComponentBinding().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionComponentBindingsWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionComponentBindingsRequest) (*teov20220901.DescribeFunctionComponentBindingsResponse, error) {
		resp := teov20220901.NewDescribeFunctionComponentBindingsResponse()
		resp.Response = &teov20220901.DescribeFunctionComponentBindingsResponseParams{
			TotalCount:                ptrInt64FCB(0),
			FunctionComponentBindings: []*teov20220901.FunctionComponentBinding{},
			RequestId:                 ptrStrFCB("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForFunctionComponentBinding()
	res := teo.ResourceTencentCloudTeoFunctionComponentBinding()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":                     "zone-12345678",
		"function_id":                 "ef-abcdefgh",
		"function_component_bindings": []interface{}{},
	})
	d.SetId("zone-12345678" + tccommon.FILED_SP + "ef-abcdefgh")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	// CONFIG type resource should not clear ID on empty bindings
	assert.NotEqual(t, "", d.Id())

	bindings := d.Get("function_component_bindings").([]interface{})
	assert.Equal(t, 0, len(bindings))
}

// TestFunctionComponentBinding_Update_Success tests Update calls ModifyFunctionComponentBindings API
func TestFunctionComponentBinding_Update_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForFunctionComponentBinding().client, "UseTeoV20220901Client", teoClient)

	modifyCalled := false
	patches.ApplyMethodFunc(teoClient, "ModifyFunctionComponentBindingsWithContext", func(_ context.Context, request *teov20220901.ModifyFunctionComponentBindingsRequest) (*teov20220901.ModifyFunctionComponentBindingsResponse, error) {
		modifyCalled = true
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.Equal(t, "ef-abcdefgh", *request.FunctionId)
		assert.Equal(t, "rebind", *request.Operation)
		assert.Equal(t, 1, len(request.FunctionComponentBindings))
		assert.Equal(t, "kv_namespace", *request.FunctionComponentBindings[0].Type)
		assert.Equal(t, "MY_KV", *request.FunctionComponentBindings[0].VariableName)
		assert.Equal(t, "zone-12345678", *request.FunctionComponentBindings[0].KVNamespaceParameters.ZoneId)
		assert.Equal(t, "my-namespace", *request.FunctionComponentBindings[0].KVNamespaceParameters.Namespace)
		resp := teov20220901.NewModifyFunctionComponentBindingsResponse()
		resp.Response = &teov20220901.ModifyFunctionComponentBindingsResponseParams{
			RequestId: ptrStrFCB("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionComponentBindingsWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionComponentBindingsRequest) (*teov20220901.DescribeFunctionComponentBindingsResponse, error) {
		resp := teov20220901.NewDescribeFunctionComponentBindingsResponse()
		resp.Response = &teov20220901.DescribeFunctionComponentBindingsResponseParams{
			TotalCount: ptrInt64FCB(1),
			FunctionComponentBindings: []*teov20220901.FunctionComponentBinding{
				{
					Type:         ptrStrFCB("kv_namespace"),
					VariableName: ptrStrFCB("MY_KV"),
					KVNamespaceParameters: &teov20220901.KVNamespaceParameters{
						ZoneId:    ptrStrFCB("zone-12345678"),
						Namespace: ptrStrFCB("my-namespace"),
					},
				},
			},
			RequestId: ptrStrFCB("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForFunctionComponentBinding()
	res := teo.ResourceTencentCloudTeoFunctionComponentBinding()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-12345678",
		"function_id": "ef-abcdefgh",
		"function_component_bindings": []interface{}{
			map[string]interface{}{
				"type":          "kv_namespace",
				"variable_name": "MY_KV",
				"kv_namespace_parameters": []interface{}{
					map[string]interface{}{
						"zone_id":   "zone-12345678",
						"namespace": "my-namespace",
					},
				},
			},
		},
	})
	d.SetId("zone-12345678" + tccommon.FILED_SP + "ef-abcdefgh")

	// Simulate HasChange by marking the resource as new (Create calls Update)
	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.True(t, modifyCalled)
}

// TestFunctionComponentBinding_Create_SetsCompositeID tests Create sets the composite ID
func TestFunctionComponentBinding_Create_SetsCompositeID(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForFunctionComponentBinding().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyFunctionComponentBindingsWithContext", func(_ context.Context, request *teov20220901.ModifyFunctionComponentBindingsRequest) (*teov20220901.ModifyFunctionComponentBindingsResponse, error) {
		resp := teov20220901.NewModifyFunctionComponentBindingsResponse()
		resp.Response = &teov20220901.ModifyFunctionComponentBindingsResponseParams{
			RequestId: ptrStrFCB("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionComponentBindingsWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionComponentBindingsRequest) (*teov20220901.DescribeFunctionComponentBindingsResponse, error) {
		resp := teov20220901.NewDescribeFunctionComponentBindingsResponse()
		resp.Response = &teov20220901.DescribeFunctionComponentBindingsResponseParams{
			TotalCount:                ptrInt64FCB(0),
			FunctionComponentBindings: []*teov20220901.FunctionComponentBinding{},
			RequestId:                 ptrStrFCB("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForFunctionComponentBinding()
	res := teo.ResourceTencentCloudTeoFunctionComponentBinding()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":                     "zone-aabbccdd",
		"function_id":                 "ef-11223344",
		"function_component_bindings": []interface{}{},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-aabbccdd"+tccommon.FILED_SP+"ef-11223344", d.Id())
}
