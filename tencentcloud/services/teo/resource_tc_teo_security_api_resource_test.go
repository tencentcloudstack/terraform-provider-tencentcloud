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

// mockMetaForSecurityAPIResource implements tccommon.ProviderMeta
type mockMetaForSecurityAPIResource struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForSecurityAPIResource) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForSecurityAPIResource{}

func newMockMetaForSecurityAPIResource() *mockMetaForSecurityAPIResource {
	return &mockMetaForSecurityAPIResource{client: &connectivity.TencentCloudClient{}}
}

func ptrStrForSecurityAPIResource(s string) *string {
	return &s
}

func ptrInt64ForSecurityAPIResource(i int64) *int64 {
	return &i
}

// go test ./tencentcloud/services/teo/ -run "TestSecurityAPIResource" -v -count=1 -gcflags="all=-l"

// TestSecurityAPIResource_Create_Success tests Create calls API and sets composite ID
func TestSecurityAPIResource_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForSecurityAPIResource().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateSecurityAPIResourceWithContext", func(ctx context.Context, request *teov20220901.CreateSecurityAPIResourceRequest) (*teov20220901.CreateSecurityAPIResourceResponse, error) {
		resp := teov20220901.NewCreateSecurityAPIResourceResponse()
		resp.Response = &teov20220901.CreateSecurityAPIResourceResponseParams{
			APIResourceIds: []*string{
				ptrStrForSecurityAPIResource("apires-001"),
			},
			RequestId: ptrStrForSecurityAPIResource("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeSecurityAPIResourceWithContext", func(ctx context.Context, request *teov20220901.DescribeSecurityAPIResourceRequest) (*teov20220901.DescribeSecurityAPIResourceResponse, error) {
		resp := teov20220901.NewDescribeSecurityAPIResourceResponse()
		resp.Response = &teov20220901.DescribeSecurityAPIResourceResponseParams{
			TotalCount: ptrInt64ForSecurityAPIResource(1),
			APIResources: []*teov20220901.APIResource{
				{
					Id:   ptrStrForSecurityAPIResource("apires-001"),
					Name: ptrStrForSecurityAPIResource("test-api-1"),
					Path: ptrStrForSecurityAPIResource("/api/v1/orders"),
					Methods: []*string{
						ptrStrForSecurityAPIResource("GET"),
						ptrStrForSecurityAPIResource("POST"),
					},
				},
			},
			RequestId: ptrStrForSecurityAPIResource("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForSecurityAPIResource()
	res := teo.ResourceTencentCloudTeoSecurityAPIResource()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"api_resources": []interface{}{
			map[string]interface{}{
				"name":    "test-api-1",
				"path":    "/api/v1/orders",
				"methods": []interface{}{"GET", "POST"},
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-1234567890#apires-001", d.Id())
}

// TestSecurityAPIResource_Create_APIError tests Create handles API error
func TestSecurityAPIResource_Create_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForSecurityAPIResource().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateSecurityAPIResourceWithContext", func(ctx context.Context, request *teov20220901.CreateSecurityAPIResourceRequest) (*teov20220901.CreateSecurityAPIResourceResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid zone_id")
	})

	meta := newMockMetaForSecurityAPIResource()
	res := teo.ResourceTencentCloudTeoSecurityAPIResource()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-invalid",
		"api_resources": []interface{}{
			map[string]interface{}{
				"name": "test-api",
				"path": "/api/v1/test",
			},
		},
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestSecurityAPIResource_Read_Success tests Read retrieves API resource data
func TestSecurityAPIResource_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForSecurityAPIResource().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeSecurityAPIResourceWithContext", func(ctx context.Context, request *teov20220901.DescribeSecurityAPIResourceRequest) (*teov20220901.DescribeSecurityAPIResourceResponse, error) {
		resp := teov20220901.NewDescribeSecurityAPIResourceResponse()
		resp.Response = &teov20220901.DescribeSecurityAPIResourceResponseParams{
			TotalCount: ptrInt64ForSecurityAPIResource(1),
			APIResources: []*teov20220901.APIResource{
				{
					Id:   ptrStrForSecurityAPIResource("apires-001"),
					Name: ptrStrForSecurityAPIResource("test-api-1"),
					Path: ptrStrForSecurityAPIResource("/api/v1/orders"),
					Methods: []*string{
						ptrStrForSecurityAPIResource("GET"),
						ptrStrForSecurityAPIResource("POST"),
					},
					APIServiceIds: []*string{
						ptrStrForSecurityAPIResource("svc-001"),
					},
					RequestConstraint: ptrStrForSecurityAPIResource(`{"key":"value"}`),
				},
			},
			RequestId: ptrStrForSecurityAPIResource("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForSecurityAPIResource()
	res := teo.ResourceTencentCloudTeoSecurityAPIResource()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"api_resources": []interface{}{
			map[string]interface{}{
				"name": "test-api-1",
				"path": "/api/v1/orders",
			},
		},
	})
	d.SetId("zone-1234567890#apires-001")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-1234567890#apires-001", d.Id())
}

// TestSecurityAPIResource_Read_NotFound tests Read handles resource not found
func TestSecurityAPIResource_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForSecurityAPIResource().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeSecurityAPIResourceWithContext", func(ctx context.Context, request *teov20220901.DescribeSecurityAPIResourceRequest) (*teov20220901.DescribeSecurityAPIResourceResponse, error) {
		resp := teov20220901.NewDescribeSecurityAPIResourceResponse()
		resp.Response = &teov20220901.DescribeSecurityAPIResourceResponseParams{
			TotalCount:   ptrInt64ForSecurityAPIResource(0),
			APIResources: []*teov20220901.APIResource{},
			RequestId:    ptrStrForSecurityAPIResource("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForSecurityAPIResource()
	res := teo.ResourceTencentCloudTeoSecurityAPIResource()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"api_resources": []interface{}{
			map[string]interface{}{
				"name": "test-api-1",
				"path": "/api/v1/orders",
			},
		},
	})
	d.SetId("zone-1234567890#apires-nonexistent")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestSecurityAPIResource_Update_Success tests Update modifies API resources
func TestSecurityAPIResource_Update_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForSecurityAPIResource().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifySecurityAPIResourceWithContext", func(ctx context.Context, request *teov20220901.ModifySecurityAPIResourceRequest) (*teov20220901.ModifySecurityAPIResourceResponse, error) {
		resp := teov20220901.NewModifySecurityAPIResourceResponse()
		resp.Response = &teov20220901.ModifySecurityAPIResourceResponseParams{
			RequestId: ptrStrForSecurityAPIResource("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeSecurityAPIResourceWithContext", func(ctx context.Context, request *teov20220901.DescribeSecurityAPIResourceRequest) (*teov20220901.DescribeSecurityAPIResourceResponse, error) {
		resp := teov20220901.NewDescribeSecurityAPIResourceResponse()
		resp.Response = &teov20220901.DescribeSecurityAPIResourceResponseParams{
			TotalCount: ptrInt64ForSecurityAPIResource(1),
			APIResources: []*teov20220901.APIResource{
				{
					Id:   ptrStrForSecurityAPIResource("apires-001"),
					Name: ptrStrForSecurityAPIResource("test-api-1-updated"),
					Path: ptrStrForSecurityAPIResource("/api/v1/updated"),
					Methods: []*string{
						ptrStrForSecurityAPIResource("GET"),
					},
				},
			},
			RequestId: ptrStrForSecurityAPIResource("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForSecurityAPIResource()
	res := teo.ResourceTencentCloudTeoSecurityAPIResource()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"api_resources": []interface{}{
			map[string]interface{}{
				"name":    "test-api-1-updated",
				"path":    "/api/v1/updated",
				"methods": []interface{}{"GET"},
			},
		},
	})
	d.SetId("zone-1234567890#apires-001")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

// TestSecurityAPIResource_Update_APIError tests Update handles API error
func TestSecurityAPIResource_Update_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForSecurityAPIResource().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifySecurityAPIResourceWithContext", func(ctx context.Context, request *teov20220901.ModifySecurityAPIResourceRequest) (*teov20220901.ModifySecurityAPIResourceResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid parameter")
	})

	meta := newMockMetaForSecurityAPIResource()
	res := teo.ResourceTencentCloudTeoSecurityAPIResource()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"api_resources": []interface{}{
			map[string]interface{}{
				"name": "test-api-1",
				"path": "/api/v1/test",
			},
		},
	})
	d.SetId("zone-1234567890#apires-001")

	err := res.Update(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestSecurityAPIResource_Delete_Success tests Delete removes API resources
func TestSecurityAPIResource_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForSecurityAPIResource().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteSecurityAPIResourceWithContext", func(ctx context.Context, request *teov20220901.DeleteSecurityAPIResourceRequest) (*teov20220901.DeleteSecurityAPIResourceResponse, error) {
		resp := teov20220901.NewDeleteSecurityAPIResourceResponse()
		resp.Response = &teov20220901.DeleteSecurityAPIResourceResponseParams{
			RequestId: ptrStrForSecurityAPIResource("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForSecurityAPIResource()
	res := teo.ResourceTencentCloudTeoSecurityAPIResource()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"api_resources": []interface{}{
			map[string]interface{}{
				"name": "test-api-1",
				"path": "/api/v1/orders",
			},
		},
	})
	d.SetId("zone-1234567890#apires-001")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestSecurityAPIResource_Delete_APIError tests Delete handles API error
func TestSecurityAPIResource_Delete_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForSecurityAPIResource().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteSecurityAPIResourceWithContext", func(ctx context.Context, request *teov20220901.DeleteSecurityAPIResourceRequest) (*teov20220901.DeleteSecurityAPIResourceResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=API resource not found")
	})

	meta := newMockMetaForSecurityAPIResource()
	res := teo.ResourceTencentCloudTeoSecurityAPIResource()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-1234567890",
		"api_resources": []interface{}{
			map[string]interface{}{
				"name": "test-api-1",
				"path": "/api/v1/orders",
			},
		},
	})
	d.SetId("zone-1234567890#apires-001")

	err := res.Delete(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
}

// TestSecurityAPIResource_Schema validates schema definition
func TestSecurityAPIResource_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoSecurityAPIResource()

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
	assert.Equal(t, 1, apiResources.MaxItems)

	// Check nested schema has path as Required
	elemRes := apiResources.Elem.(*schema.Resource)
	assert.Contains(t, elemRes.Schema, "path")
	assert.True(t, elemRes.Schema["path"].Required)

	// Check nested schema has id as Computed
	assert.Contains(t, elemRes.Schema, "id")
	assert.True(t, elemRes.Schema["id"].Computed)
}
