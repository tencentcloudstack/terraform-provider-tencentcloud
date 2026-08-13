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

// mockMetaFunctionReplica implements tccommon.ProviderMeta
type mockMetaFunctionReplica struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaFunctionReplica) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaFunctionReplica{}

func newMockMetaFunctionReplica() *mockMetaFunctionReplica {
	return &mockMetaFunctionReplica{client: &connectivity.TencentCloudClient{}}
}

func ptrStringFunctionReplica(s string) *string {
	return &s
}

func ptrInt64FunctionReplica(i int64) *int64 {
	return &i
}

// go test ./tencentcloud/services/teo/ -run "TestTeoFunctionReplica_Create" -v -count=1 -gcflags="all=-l"
func TestTeoFunctionReplica_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaFunctionReplica().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateFunctionReplicaWithContext", func(_ context.Context, request *teov20220901.CreateFunctionReplicaRequest) (*teov20220901.CreateFunctionReplicaResponse, error) {
		assert.Equal(t, "zone-test123", *request.ZoneId)
		assert.Equal(t, "ef-test456", *request.FunctionId)
		assert.Equal(t, "test-replica", *request.ReplicaName)
		assert.Equal(t, "console.log(123)", *request.Content)
		assert.Equal(t, "test remark", *request.Remark)

		resp := teov20220901.NewCreateFunctionReplicaResponse()
		resp.Response = &teov20220901.CreateFunctionReplicaResponseParams{
			RequestId: ptrStringFunctionReplica("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionReplicasWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
		resp := teov20220901.NewDescribeFunctionReplicasResponse()
		resp.Response = &teov20220901.DescribeFunctionReplicasResponseParams{
			TotalCount: ptrInt64FunctionReplica(1),
			FunctionReplicas: []*teov20220901.FunctionReplica{
				{
					FunctionId:  ptrStringFunctionReplica("ef-test456"),
					ReplicaName: ptrStringFunctionReplica("test-replica"),
					Content:     ptrStringFunctionReplica("console.log(123)"),
					Remark:      ptrStringFunctionReplica("test remark"),
					CreatedOn:   ptrStringFunctionReplica("2024-01-01T00:00:00Z"),
					ModifiedOn:  ptrStringFunctionReplica("2024-01-01T00:00:00Z"),
				},
			},
			RequestId: ptrStringFunctionReplica("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaFunctionReplica()
	res := teo.ResourceTencentCloudTeoFunctionReplica()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-test123",
		"function_id":  "ef-test456",
		"replica_name": "test-replica",
		"content":      "console.log(123)",
		"remark":       "test remark",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-test123#ef-test456#test-replica", d.Id())
}

// go test ./tencentcloud/services/teo/ -run "TestTeoFunctionReplica_Read" -v -count=1 -gcflags="all=-l"
func TestTeoFunctionReplica_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaFunctionReplica().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionReplicasWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
		assert.Equal(t, "zone-test123", *request.ZoneId)
		assert.Equal(t, "ef-test456", *request.FunctionId)
		assert.Equal(t, int64(200), *request.Limit)

		resp := teov20220901.NewDescribeFunctionReplicasResponse()
		resp.Response = &teov20220901.DescribeFunctionReplicasResponseParams{
			TotalCount: ptrInt64FunctionReplica(1),
			FunctionReplicas: []*teov20220901.FunctionReplica{
				{
					FunctionId:  ptrStringFunctionReplica("ef-test456"),
					ReplicaName: ptrStringFunctionReplica("test-replica"),
					Content:     ptrStringFunctionReplica("console.log(123)"),
					Remark:      ptrStringFunctionReplica("test remark"),
					CreatedOn:   ptrStringFunctionReplica("2024-01-01T00:00:00Z"),
					ModifiedOn:  ptrStringFunctionReplica("2024-01-01T00:00:00Z"),
				},
			},
			RequestId: ptrStringFunctionReplica("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaFunctionReplica()
	res := teo.ResourceTencentCloudTeoFunctionReplica()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-test123",
		"function_id":  "ef-test456",
		"replica_name": "test-replica",
		"content":      "console.log(123)",
		"remark":       "test remark",
	})
	d.SetId("zone-test123#ef-test456#test-replica")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "console.log(123)", d.Get("content"))
	assert.Equal(t, "test remark", d.Get("remark"))
}

// go test ./tencentcloud/services/teo/ -run "TestTeoFunctionReplica_ReadNotFound" -v -count=1 -gcflags="all=-l"
func TestTeoFunctionReplica_ReadNotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaFunctionReplica().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionReplicasWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
		resp := teov20220901.NewDescribeFunctionReplicasResponse()
		resp.Response = &teov20220901.DescribeFunctionReplicasResponseParams{
			TotalCount:       ptrInt64FunctionReplica(0),
			FunctionReplicas: []*teov20220901.FunctionReplica{},
			RequestId:        ptrStringFunctionReplica("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaFunctionReplica()
	res := teo.ResourceTencentCloudTeoFunctionReplica()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-test123",
		"function_id":  "ef-test456",
		"replica_name": "test-replica",
		"content":      "console.log(123)",
		"remark":       "test remark",
	})
	d.SetId("zone-test123#ef-test456#test-replica")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// go test ./tencentcloud/services/teo/ -run "TestTeoFunctionReplica_Update" -v -count=1 -gcflags="all=-l"
func TestTeoFunctionReplica_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaFunctionReplica().client, "UseTeoV20220901Client", teoClient)

	modifyCalled := false
	patches.ApplyMethodFunc(teoClient, "ModifyFunctionReplicaWithContext", func(_ context.Context, request *teov20220901.ModifyFunctionReplicaRequest) (*teov20220901.ModifyFunctionReplicaResponse, error) {
		modifyCalled = true
		assert.Equal(t, "zone-test123", *request.ZoneId)
		assert.Equal(t, "ef-test456", *request.FunctionId)
		assert.Equal(t, "test-replica", *request.ReplicaName)
		assert.Equal(t, "console.log(456)", *request.Content)
		assert.Equal(t, "updated remark", *request.Remark)

		resp := teov20220901.NewModifyFunctionReplicaResponse()
		resp.Response = &teov20220901.ModifyFunctionReplicaResponseParams{
			RequestId: ptrStringFunctionReplica("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeFunctionReplicasWithContext", func(_ context.Context, request *teov20220901.DescribeFunctionReplicasRequest) (*teov20220901.DescribeFunctionReplicasResponse, error) {
		resp := teov20220901.NewDescribeFunctionReplicasResponse()
		resp.Response = &teov20220901.DescribeFunctionReplicasResponseParams{
			TotalCount: ptrInt64FunctionReplica(1),
			FunctionReplicas: []*teov20220901.FunctionReplica{
				{
					FunctionId:  ptrStringFunctionReplica("ef-test456"),
					ReplicaName: ptrStringFunctionReplica("test-replica"),
					Content:     ptrStringFunctionReplica("console.log(456)"),
					Remark:      ptrStringFunctionReplica("updated remark"),
					CreatedOn:   ptrStringFunctionReplica("2024-01-01T00:00:00Z"),
					ModifiedOn:  ptrStringFunctionReplica("2024-01-02T00:00:00Z"),
				},
			},
			RequestId: ptrStringFunctionReplica("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaFunctionReplica()
	res := teo.ResourceTencentCloudTeoFunctionReplica()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-test123",
		"function_id":  "ef-test456",
		"replica_name": "test-replica",
		"content":      "console.log(456)",
		"remark":       "updated remark",
	})
	d.SetId("zone-test123#ef-test456#test-replica")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	_ = modifyCalled
}

// go test ./tencentcloud/services/teo/ -run "TestTeoFunctionReplica_Delete" -v -count=1 -gcflags="all=-l"
func TestTeoFunctionReplica_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaFunctionReplica().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteFunctionReplicaWithContext", func(_ context.Context, request *teov20220901.DeleteFunctionReplicaRequest) (*teov20220901.DeleteFunctionReplicaResponse, error) {
		assert.Equal(t, "zone-test123", *request.ZoneId)
		assert.Equal(t, "ef-test456", *request.FunctionId)
		assert.Equal(t, 1, len(request.ReplicaNames))
		assert.Equal(t, "test-replica", *request.ReplicaNames[0])

		resp := teov20220901.NewDeleteFunctionReplicaResponse()
		resp.Response = &teov20220901.DeleteFunctionReplicaResponseParams{
			RequestId: ptrStringFunctionReplica("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaFunctionReplica()
	res := teo.ResourceTencentCloudTeoFunctionReplica()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-test123",
		"function_id":  "ef-test456",
		"replica_name": "test-replica",
		"content":      "console.log(123)",
		"remark":       "test remark",
	})
	d.SetId("zone-test123#ef-test456#test-replica")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}
