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

// mockMetaForSecurityApiResource implements tccommon.ProviderMeta
type mockMetaForSecurityApiResource struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForSecurityApiResource) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForSecurityApiResource{}

func newMockMetaForSecurityApiResource() *mockMetaForSecurityApiResource {
	return &mockMetaForSecurityApiResource{client: &connectivity.TencentCloudClient{}}
}

func ptrStringForSecurityApiResource(s string) *string {
	return &s
}

func ptrInt64ForSecurityApiResource(i int64) *int64 {
	return &i
}

// go test ./tencentcloud/services/teo/ -run "TestSecurityApiResource" -v -count=1 -gcflags="all=-l"

// TestSecurityApiResource_Create_Success tests Create calls API and sets ID
func TestSecurityApiResource_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForSecurityApiResource().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateSecurityAPIResourceWithContext", func(ctx context.Context, request *teov20220901.CreateSecurityAPIResourceRequest) (*teov20220901.CreateSecurityAPIResourceResponse, error) {
		resp := teov20220901.NewCreateSecurityAPIResourceResponse()
		resp.Response = &teov20220901.CreateSecurityAPIResourceResponseParams{
			APIResourceIds: []*string{
				ptrStringForSecurityApiResource("res-001"),
				ptrStringForSecurityApiResource("res-002"),
			},
			RequestId: ptrStringForSecurityApiResource("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeSecurityAPIResourceWithContext", func(ctx context.Context, request *teov20220901.DescribeSecurityAPIResourceRequest) (*teov20220901.DescribeSecurityAPIResourceResponse, error) {
		resp := teov20220901.NewDescribeSecurityAPIResourceResponse()
		resp.Response = &teov20220901.DescribeSecurityAPIResourceResponseParams{
			TotalCount: ptrInt64ForSecurityApiResource(2),
			APIResources: []*teov20220901.APIResource{
				{
					Id:   ptrStringForSecurityApiResource("res-001"),
					Name: ptrStringForSecurityApiResource("test-api-1"),
					Path: ptrStringForSecurityApiResource("/api/v1/test"),
					Methods: []*string{
						ptrStringForSecurityApiResource("GET"),
						ptrStringForSecurityApiResource("POST"),
					},
				},
				{
					Id:   ptrStringForSecurityApiResource("res-002"),
					Name: ptrStringForSecurityApiResource("test-api-2"),
					Path: ptrStringForSecurityApiResource("/api/v2/test"),
					Methods: []*string{
						ptrStringForSecurityApiResource("GET"),
					},
				},
			},
			RequestId: ptrStringForSecurityApiResource("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForSecurityApiResource()
	res := teo.ResourceTencentCloudTeoSecurityApiResource()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"api_resources": []interface{}{
			map[string]interface{}{
				"name":    "test-api-1",
				"path":    "/api/v1/test",
				"methods": []interface{}{"GET", "POST"},
			},
			map[string]interface{}{
				"name":    "test-api-2",
				"path":    "/api/v2/test",
				"methods": []interface{}{"GET"},
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-1234567890", d.Id())
}

// TestSecurityApiResource_Create_APIError tests Create handles API error
func TestSecurityApiResource_Create_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForSecurityApiResource().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateSecurityAPIResourceWithContext", func(ctx context.Context, request *teov20220901.CreateSecurityAPIResourceRequest) (*teov20220901.CreateSecurityAPIResourceResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid zone_id")
	})

	meta := newMockMetaForSecurityApiResource()
	res := teo.ResourceTencentCloudTeoSecurityApiResource()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-invalid",
		"api_resources": []interface{}{
			map[string]interface{}{
				"name": "test-api",
			},
		},
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestSecurityApiResource_Read_Success tests Read retrieves API resource data
func TestSecurityApiResource_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForSecurityApiResource().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeSecurityAPIResourceWithContext", func(ctx context.Context, request *teov20220901.DescribeSecurityAPIResourceRequest) (*teov20220901.DescribeSecurityAPIResourceResponse, error) {
		resp := teov20220901.NewDescribeSecurityAPIResourceResponse()
		resp.Response = &teov20220901.DescribeSecurityAPIResourceResponseParams{
			TotalCount: ptrInt64ForSecurityApiResource(1),
			APIResources: []*teov20220901.APIResource{
				{
					Id:   ptrStringForSecurityApiResource("res-001"),
					Name: ptrStringForSecurityApiResource("test-api-1"),
					Path: ptrStringForSecurityApiResource("/api/v1/test"),
					Methods: []*string{
						ptrStringForSecurityApiResource("GET"),
						ptrStringForSecurityApiResource("POST"),
					},
					APIServiceIds: []*string{
						ptrStringForSecurityApiResource("svc-001"),
					},
					RequestConstraint: ptrStringForSecurityApiResource(`{"key":"value"}`),
				},
			},
			RequestId: ptrStringForSecurityApiResource("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForSecurityApiResource()
	res := teo.ResourceTencentCloudTeoSecurityApiResource()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"api_resources": []interface{}{
			map[string]interface{}{
				"name": "test-api-1",
			},
		},
	})
	d.SetId("zone-1234567890")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-1234567890", d.Id())
}

// TestSecurityApiResource_Read_NotFound tests Read handles resource not found
func TestSecurityApiResource_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForSecurityApiResource().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeSecurityAPIResourceWithContext", func(ctx context.Context, request *teov20220901.DescribeSecurityAPIResourceRequest) (*teov20220901.DescribeSecurityAPIResourceResponse, error) {
		resp := teov20220901.NewDescribeSecurityAPIResourceResponse()
		resp.Response = &teov20220901.DescribeSecurityAPIResourceResponseParams{
			TotalCount:   ptrInt64ForSecurityApiResource(0),
			APIResources: []*teov20220901.APIResource{},
			RequestId:    ptrStringForSecurityApiResource("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForSecurityApiResource()
	res := teo.ResourceTencentCloudTeoSecurityApiResource()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"api_resources": []interface{}{
			map[string]interface{}{
				"name": "test-api-1",
			},
		},
	})
	d.SetId("zone-1234567890")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestSecurityApiResource_Update_Success tests Update modifies API resources
func TestSecurityApiResource_Update_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForSecurityApiResource().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifySecurityAPIResourceWithContext", func(ctx context.Context, request *teov20220901.ModifySecurityAPIResourceRequest) (*teov20220901.ModifySecurityAPIResourceResponse, error) {
		resp := teov20220901.NewModifySecurityAPIResourceResponse()
		resp.Response = &teov20220901.ModifySecurityAPIResourceResponseParams{
			RequestId: ptrStringForSecurityApiResource("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeSecurityAPIResourceWithContext", func(ctx context.Context, request *teov20220901.DescribeSecurityAPIResourceRequest) (*teov20220901.DescribeSecurityAPIResourceResponse, error) {
		resp := teov20220901.NewDescribeSecurityAPIResourceResponse()
		resp.Response = &teov20220901.DescribeSecurityAPIResourceResponseParams{
			TotalCount: ptrInt64ForSecurityApiResource(1),
			APIResources: []*teov20220901.APIResource{
				{
					Id:   ptrStringForSecurityApiResource("res-001"),
					Name: ptrStringForSecurityApiResource("test-api-1-updated"),
					Path: ptrStringForSecurityApiResource("/api/v1/updated"),
					Methods: []*string{
						ptrStringForSecurityApiResource("GET"),
					},
				},
			},
			RequestId: ptrStringForSecurityApiResource("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForSecurityApiResource()
	res := teo.ResourceTencentCloudTeoSecurityApiResource()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"api_resources": []interface{}{
			map[string]interface{}{
				"id":      "res-001",
				"name":    "test-api-1-updated",
				"path":    "/api/v1/updated",
				"methods": []interface{}{"GET"},
			},
		},
	})
	d.SetId("zone-1234567890")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

// TestSecurityApiResource_Update_APIError tests Update handles API error
func TestSecurityApiResource_Update_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForSecurityApiResource().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifySecurityAPIResourceWithContext", func(ctx context.Context, request *teov20220901.ModifySecurityAPIResourceRequest) (*teov20220901.ModifySecurityAPIResourceResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid parameter")
	})

	meta := newMockMetaForSecurityApiResource()
	res := teo.ResourceTencentCloudTeoSecurityApiResource()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"api_resources": []interface{}{
			map[string]interface{}{
				"id":   "res-001",
				"name": "test-api-1",
			},
		},
	})
	d.SetId("zone-1234567890")

	err := res.Update(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestSecurityApiResource_Delete_Success tests Delete removes API resources
func TestSecurityApiResource_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForSecurityApiResource().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteSecurityAPIResourceWithContext", func(ctx context.Context, request *teov20220901.DeleteSecurityAPIResourceRequest) (*teov20220901.DeleteSecurityAPIResourceResponse, error) {
		resp := teov20220901.NewDeleteSecurityAPIResourceResponse()
		resp.Response = &teov20220901.DeleteSecurityAPIResourceResponseParams{
			RequestId: ptrStringForSecurityApiResource("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForSecurityApiResource()
	res := teo.ResourceTencentCloudTeoSecurityApiResource()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"api_resources": []interface{}{
			map[string]interface{}{
				"name": "test-api-1",
			},
		},
		"api_resource_ids": []interface{}{"res-001", "res-002"},
	})
	d.SetId("zone-1234567890")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestSecurityApiResource_Delete_APIError tests Delete handles API error
func TestSecurityApiResource_Delete_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForSecurityApiResource().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteSecurityAPIResourceWithContext", func(ctx context.Context, request *teov20220901.DeleteSecurityAPIResourceRequest) (*teov20220901.DeleteSecurityAPIResourceResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=API resource not found")
	})

	meta := newMockMetaForSecurityApiResource()
	res := teo.ResourceTencentCloudTeoSecurityApiResource()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"api_resources": []interface{}{
			map[string]interface{}{
				"name": "test-api-1",
			},
		},
		"api_resource_ids": []interface{}{"res-001"},
	})
	d.SetId("zone-1234567890")

	err := res.Delete(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
}

// TestSecurityApiResource_Schema validates schema definition
func TestSecurityApiResource_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoSecurityApiResource()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.NotNil(t, res.Update)
	assert.NotNil(t, res.Delete)
	assert.NotNil(t, res.Importer)

	// Check zone_id field
	assert.Contains(t, res.Schema, "zone_id")
	zoneId := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneId.Type)
	assert.True(t, zoneId.Required)
	assert.True(t, zoneId.ForceNew)

	// Check api_resources field
	assert.Contains(t, res.Schema, "api_resources")
	apiResources := res.Schema["api_resources"]
	assert.Equal(t, schema.TypeList, apiResources.Type)
	assert.True(t, apiResources.Required)

	// Check api_resource_ids field
	assert.Contains(t, res.Schema, "api_resource_ids")
	apiResourceIds := res.Schema["api_resource_ids"]
	assert.Equal(t, schema.TypeList, apiResourceIds.Type)
	assert.True(t, apiResourceIds.Computed)
}
