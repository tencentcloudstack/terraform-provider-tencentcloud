package teo_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// go test ./tencentcloud/services/teo/ -run "TestSecurityApiService" -v -count=1 -gcflags="all=-l"

func ptrStringVal(s string) *string {
	return &s
}

// TestSecurityApiService_Create_Success tests Create calls API and sets ID
func TestSecurityApiService_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateSecurityAPIServiceWithContext", func(ctx context.Context, request *teov20220901.CreateSecurityAPIServiceRequest) (*teov20220901.CreateSecurityAPIServiceResponse, error) {
		resp := teov20220901.NewCreateSecurityAPIServiceResponse()
		resp.Response = &teov20220901.CreateSecurityAPIServiceResponseParams{
			APIServiceIds: []*string{ptrStringVal("svc-12345"), ptrStringVal("svc-67890")},
			RequestId:     ptrStringVal("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeSecurityAPIService", func(request *teov20220901.DescribeSecurityAPIServiceRequest) (*teov20220901.DescribeSecurityAPIServiceResponse, error) {
		resp := teov20220901.NewDescribeSecurityAPIServiceResponse()
		resp.Response = &teov20220901.DescribeSecurityAPIServiceResponseParams{
			TotalCount: ptrInt64(2),
			APIServices: []*teov20220901.APIService{
				{
					Id:       ptrStringVal("svc-12345"),
					Name:     ptrStringVal("my-api-service-1"),
					BasePath: ptrStringVal("/api/v1"),
				},
				{
					Id:       ptrStringVal("svc-67890"),
					Name:     ptrStringVal("my-api-service-2"),
					BasePath: ptrStringVal("/api/v2"),
				},
			},
			RequestId: ptrStringVal("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeSecurityAPIResource", func(request *teov20220901.DescribeSecurityAPIResourceRequest) (*teov20220901.DescribeSecurityAPIResourceResponse, error) {
		resp := teov20220901.NewDescribeSecurityAPIResourceResponse()
		resp.Response = &teov20220901.DescribeSecurityAPIResourceResponseParams{
			TotalCount:   ptrInt64(0),
			APIResources: []*teov20220901.APIResource{},
			RequestId:    ptrStringVal("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoSecurityApiService()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-2qtuhspy7cr6",
		"api_services": []interface{}{
			map[string]interface{}{
				"name":      "my-api-service-1",
				"base_path": "/api/v1",
			},
			map[string]interface{}{
				"name":      "my-api-service-2",
				"base_path": "/api/v2",
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-2qtuhspy7cr6#svc-12345,svc-67890", d.Id())
}

// TestSecurityApiService_Create_APIError tests Create handles API error
func TestSecurityApiService_Create_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateSecurityAPIServiceWithContext", func(ctx context.Context, request *teov20220901.CreateSecurityAPIServiceRequest) (*teov20220901.CreateSecurityAPIServiceResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid zone_id")
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoSecurityApiService()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-invalid",
		"api_services": []interface{}{
			map[string]interface{}{
				"name":      "test",
				"base_path": "/api",
			},
		},
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestSecurityApiService_Read_Success tests Read retrieves API service data
func TestSecurityApiService_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeSecurityAPIService", func(request *teov20220901.DescribeSecurityAPIServiceRequest) (*teov20220901.DescribeSecurityAPIServiceResponse, error) {
		resp := teov20220901.NewDescribeSecurityAPIServiceResponse()
		resp.Response = &teov20220901.DescribeSecurityAPIServiceResponseParams{
			TotalCount: ptrInt64(1),
			APIServices: []*teov20220901.APIService{
				{
					Id:       ptrStringVal("svc-12345"),
					Name:     ptrStringVal("my-api-service"),
					BasePath: ptrStringVal("/api/v1"),
				},
			},
			RequestId: ptrStringVal("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeSecurityAPIResource", func(request *teov20220901.DescribeSecurityAPIResourceRequest) (*teov20220901.DescribeSecurityAPIResourceResponse, error) {
		resp := teov20220901.NewDescribeSecurityAPIResourceResponse()
		resp.Response = &teov20220901.DescribeSecurityAPIResourceResponseParams{
			TotalCount:   ptrInt64(0),
			APIResources: []*teov20220901.APIResource{},
			RequestId:    ptrStringVal("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoSecurityApiService()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-2qtuhspy7cr6",
	})
	d.SetId("zone-2qtuhspy7cr6#svc-12345")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-2qtuhspy7cr6", d.Get("zone_id"))

	apiServices := d.Get("api_services").([]interface{})
	assert.Equal(t, 1, len(apiServices))
	svcMap := apiServices[0].(map[string]interface{})
	assert.Equal(t, "my-api-service", svcMap["name"])
	assert.Equal(t, "/api/v1", svcMap["base_path"])

	apiServiceIds := d.Get("api_service_ids").([]interface{})
	assert.Equal(t, 1, len(apiServiceIds))
	assert.Equal(t, "svc-12345", apiServiceIds[0])
}

// TestSecurityApiService_Read_NotFound tests Read handles resource not found
func TestSecurityApiService_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeSecurityAPIService", func(request *teov20220901.DescribeSecurityAPIServiceRequest) (*teov20220901.DescribeSecurityAPIServiceResponse, error) {
		resp := teov20220901.NewDescribeSecurityAPIServiceResponse()
		resp.Response = &teov20220901.DescribeSecurityAPIServiceResponseParams{
			TotalCount:  ptrInt64(0),
			APIServices: []*teov20220901.APIService{},
			RequestId:   ptrStringVal("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoSecurityApiService()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-2qtuhspy7cr6",
	})
	d.SetId("zone-2qtuhspy7cr6#svc-notfound")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestSecurityApiService_Update_ApiResources tests Update with api_resources changes
func TestSecurityApiService_Update_ApiResources(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifySecurityAPIResourceWithContext", func(ctx context.Context, request *teov20220901.ModifySecurityAPIResourceRequest) (*teov20220901.ModifySecurityAPIResourceResponse, error) {
		resp := teov20220901.NewModifySecurityAPIResourceResponse()
		resp.Response = &teov20220901.ModifySecurityAPIResourceResponseParams{
			RequestId: ptrStringVal("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeSecurityAPIService", func(request *teov20220901.DescribeSecurityAPIServiceRequest) (*teov20220901.DescribeSecurityAPIServiceResponse, error) {
		resp := teov20220901.NewDescribeSecurityAPIServiceResponse()
		resp.Response = &teov20220901.DescribeSecurityAPIServiceResponseParams{
			TotalCount: ptrInt64(1),
			APIServices: []*teov20220901.APIService{
				{
					Id:       ptrStringVal("svc-12345"),
					Name:     ptrStringVal("my-api-service"),
					BasePath: ptrStringVal("/api/v1"),
				},
			},
			RequestId: ptrStringVal("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeSecurityAPIResource", func(request *teov20220901.DescribeSecurityAPIResourceRequest) (*teov20220901.DescribeSecurityAPIResourceResponse, error) {
		resp := teov20220901.NewDescribeSecurityAPIResourceResponse()
		resp.Response = &teov20220901.DescribeSecurityAPIResourceResponseParams{
			TotalCount:   ptrInt64(0),
			APIResources: []*teov20220901.APIResource{},
			RequestId:    ptrStringVal("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoSecurityApiService()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-2qtuhspy7cr6",
		"api_services": []interface{}{
			map[string]interface{}{
				"name":      "my-api-service",
				"base_path": "/api/v1",
			},
		},
		"api_service_ids": []interface{}{"svc-12345"},
		"api_resources": []interface{}{
			map[string]interface{}{
				"name":            "my-resource",
				"api_service_ids": []interface{}{"svc-12345"},
				"path":            "/api/v1/users",
				"methods":         []interface{}{"GET", "POST"},
			},
		},
	})
	d.SetId("zone-2qtuhspy7cr6#svc-12345")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

// TestSecurityApiService_Delete_Success tests Delete removes API service
func TestSecurityApiService_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteSecurityAPIServiceWithContext", func(ctx context.Context, request *teov20220901.DeleteSecurityAPIServiceRequest) (*teov20220901.DeleteSecurityAPIServiceResponse, error) {
		resp := teov20220901.NewDeleteSecurityAPIServiceResponse()
		resp.Response = &teov20220901.DeleteSecurityAPIServiceResponseParams{
			RequestId: ptrStringVal("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoSecurityApiService()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-2qtuhspy7cr6",
		"api_services": []interface{}{
			map[string]interface{}{
				"name":      "my-api-service",
				"base_path": "/api/v1",
			},
		},
		"api_service_ids": []interface{}{"svc-12345"},
	})
	d.SetId("zone-2qtuhspy7cr6#svc-12345")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestSecurityApiService_Delete_APIError tests Delete handles API error
func TestSecurityApiService_Delete_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMeta().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteSecurityAPIServiceWithContext", func(ctx context.Context, request *teov20220901.DeleteSecurityAPIServiceRequest) (*teov20220901.DeleteSecurityAPIServiceResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=Service not found")
	})

	meta := newMockMeta()
	res := teo.ResourceTencentCloudTeoSecurityApiService()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-2qtuhspy7cr6",
		"api_services": []interface{}{
			map[string]interface{}{
				"name":      "my-api-service",
				"base_path": "/api/v1",
			},
		},
		"api_service_ids": []interface{}{"svc-12345"},
	})
	d.SetId("zone-2qtuhspy7cr6#svc-12345")

	err := res.Delete(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
}

// TestSecurityApiService_Schema validates schema definition
func TestSecurityApiService_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoSecurityApiService()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.NotNil(t, res.Update)
	assert.NotNil(t, res.Delete)
	assert.NotNil(t, res.Importer)

	// Check required fields
	assert.Contains(t, res.Schema, "zone_id")
	zoneId := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneId.Type)
	assert.True(t, zoneId.Required)
	assert.True(t, zoneId.ForceNew)

	assert.Contains(t, res.Schema, "api_services")
	apiServices := res.Schema["api_services"]
	assert.Equal(t, schema.TypeList, apiServices.Type)
	assert.True(t, apiServices.Required)
	assert.True(t, apiServices.ForceNew)

	// Check computed fields
	assert.Contains(t, res.Schema, "api_service_ids")
	apiServiceIds := res.Schema["api_service_ids"]
	assert.Equal(t, schema.TypeList, apiServiceIds.Type)
	assert.True(t, apiServiceIds.Computed)

	// Check optional fields
	assert.Contains(t, res.Schema, "api_resources")
	apiResources := res.Schema["api_resources"]
	assert.Equal(t, schema.TypeList, apiResources.Type)
	assert.True(t, apiResources.Optional)
}
