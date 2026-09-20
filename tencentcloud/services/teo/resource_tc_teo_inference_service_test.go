package teo_test

import (
	"context"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teo "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

type mockMetaTeoInferenceService struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaTeoInferenceService) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaTeoInferenceService{}

func newMockMetaTeoInferenceService() *mockMetaTeoInferenceService {
	return &mockMetaTeoInferenceService{client: &connectivity.TencentCloudClient{}}
}

func ptrStringTeoIs(s string) *string    { return &s }
func ptrInt64TeoIs(v int64) *int64       { return &v }
func ptrFloat64TeoIs(v float64) *float64 { return &v }

// buildInferenceServiceResponse builds a canned DescribeInferenceServices response
// containing a single InferenceService that matches the given id.
func buildInferenceServiceResponse(id string) *teo.DescribeInferenceServicesResponse {
	resp := teo.NewDescribeInferenceServicesResponse()
	resp.Response = &teo.DescribeInferenceServicesResponseParams{
		TotalCount: ptrInt64TeoIs(1),
		Services: []*teo.InferenceService{
			{
				ServiceId:            ptrStringTeoIs(id),
				Name:                 ptrStringTeoIs("tf-example-inference-service"),
				Description:          ptrStringTeoIs("tf example inference service"),
				ListenPort:           ptrInt64TeoIs(8500),
				RequestPaths:         []*string{ptrStringTeoIs("/v1/predict")},
				Status:               ptrStringTeoIs("Running"),
				ScalingStatus:        ptrStringTeoIs("Normal"),
				CurrentInstanceCount: ptrInt64TeoIs(1),
				InferenceURL:         ptrStringTeoIs("https://example.inference.edgeone.ai"),
				CreateTime:           ptrStringTeoIs("2026-09-20T10:00:00Z"),
				UpdateTime:           ptrStringTeoIs("2026-09-20T10:00:00Z"),
			},
		},
		RequestId: ptrStringTeoIs("fake-request-id"),
	}
	return resp
}

// TestTeoInferenceService_Create verifies the Create flow sets the composite id from the mock response.
func TestTeoInferenceService_Create(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teo.Client{}
	patches.ApplyMethodReturn(newMockMetaTeoInferenceService().client, "UseTeoV20220901Client", teoClient)

	var capturedRequest *teo.CreateInferenceServiceRequest
	patches.ApplyMethodFunc(teoClient, "CreateInferenceServiceWithContext", func(_ context.Context, request *teo.CreateInferenceServiceRequest) (*teo.CreateInferenceServiceResponse, error) {
		capturedRequest = request
		assert.Equal(t, "zone-3fkff38fyw8s", *request.ZoneId)
		assert.Equal(t, "tf-example-inference-service", *request.Name)
		assert.Equal(t, int64(8500), *request.ListenPort)

		resp := teo.NewCreateInferenceServiceResponse()
		resp.Response = &teo.CreateInferenceServiceResponseParams{
			ServiceId: ptrStringTeoIs("is-abc123"),
			RequestId: ptrStringTeoIs("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeInferenceServicesWithContext", func(_ context.Context, request *teo.DescribeInferenceServicesRequest) (*teo.DescribeInferenceServicesResponse, error) {
		return buildInferenceServiceResponse("is-abc123"), nil
	})

	meta := newMockMetaTeoInferenceService()
	res := teo.ResourceTencentCloudTeoInferenceService()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-3fkff38fyw8s",
		"name":        "tf-example-inference-service",
		"listen_port": 8500,
		"description": "tf example inference service",
		"containers": []interface{}{
			map[string]interface{}{
				"image_type": "TCR",
			},
		},
		"resource_config": []interface{}{
			map[string]interface{}{
				"scaling_mode": "Auto",
			},
		},
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-3fkff38fyw8s#is-abc123", d.Id())
	assert.NotNil(t, capturedRequest)
}

// TestTeoInferenceService_Create_EmptyServiceId verifies Create rejects an empty ServiceId.
func TestTeoInferenceService_Create_EmptyServiceId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teo.Client{}
	patches.ApplyMethodReturn(newMockMetaTeoInferenceService().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateInferenceServiceWithContext", func(_ context.Context, request *teo.CreateInferenceServiceRequest) (*teo.CreateInferenceServiceResponse, error) {
		resp := teo.NewCreateInferenceServiceResponse()
		resp.Response = &teo.CreateInferenceServiceResponseParams{
			ServiceId: ptrStringTeoIs(""),
			RequestId: ptrStringTeoIs("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaTeoInferenceService()
	res := teo.ResourceTencentCloudTeoInferenceService()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-3fkff38fyw8s",
		"name":        "tf-example-inference-service",
		"listen_port": 8500,
		"containers": []interface{}{
			map[string]interface{}{
				"image_type": "TCR",
			},
		},
		"resource_config": []interface{}{
			map[string]interface{}{
				"scaling_mode": "Auto",
			},
		},
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ServiceId")
}

// TestTeoInferenceService_Read verifies Read populates computed fields from the response.
func TestTeoInferenceService_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teo.Client{}
	patches.ApplyMethodReturn(newMockMetaTeoInferenceService().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeInferenceServicesWithContext", func(_ context.Context, request *teo.DescribeInferenceServicesRequest) (*teo.DescribeInferenceServicesResponse, error) {
		return buildInferenceServiceResponse("is-read123"), nil
	})

	meta := newMockMetaTeoInferenceService()
	res := teo.ResourceTencentCloudTeoInferenceService()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-3fkff38fyw8s",
		"name":        "tf-example-inference-service",
		"listen_port": 8500,
	})
	d.SetId("zone-3fkff38fyw8s#is-read123")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-3fkff38fyw8s#is-read123", d.Id())

	assert.Equal(t, "tf-example-inference-service", d.Get("name"))
	assert.Equal(t, "tf example inference service", d.Get("description"))
	assert.Equal(t, int(8500), d.Get("listen_port"))
	assert.Equal(t, "is-read123", d.Get("service_id"))
	assert.Equal(t, "Running", d.Get("status"))
	assert.Equal(t, "Normal", d.Get("scaling_status"))
	assert.Equal(t, int(1), d.Get("current_instance_count"))
}

// TestTeoInferenceService_Read_NotFound verifies Read clears state when the describe returns empty.
func TestTeoInferenceService_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teo.Client{}
	patches.ApplyMethodReturn(newMockMetaTeoInferenceService().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeInferenceServicesWithContext", func(_ context.Context, request *teo.DescribeInferenceServicesRequest) (*teo.DescribeInferenceServicesResponse, error) {
		resp := teo.NewDescribeInferenceServicesResponse()
		resp.Response = &teo.DescribeInferenceServicesResponseParams{
			TotalCount: ptrInt64TeoIs(0),
			Services:   []*teo.InferenceService{},
			RequestId:  ptrStringTeoIs("fake-request-id-empty"),
		}
		return resp, nil
	})

	meta := newMockMetaTeoInferenceService()
	res := teo.ResourceTencentCloudTeoInferenceService()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-3fkff38fyw8s",
	})
	d.SetId("zone-3fkff38fyw8s#is-missing")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestTeoInferenceService_Update verifies Update triggers ModifyInferenceService with the changed description.
func TestTeoInferenceService_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teo.Client{}
	patches.ApplyMethodReturn(newMockMetaTeoInferenceService().client, "UseTeoV20220901Client", teoClient)

	var capturedModifyRequest *teo.ModifyInferenceServiceRequest
	patches.ApplyMethodFunc(teoClient, "ModifyInferenceServiceWithContext", func(_ context.Context, request *teo.ModifyInferenceServiceRequest) (*teo.ModifyInferenceServiceResponse, error) {
		capturedModifyRequest = request
		assert.Equal(t, "zone-3fkff38fyw8s", *request.ZoneId)
		assert.Equal(t, "is-update123", *request.ServiceId)

		resp := teo.NewModifyInferenceServiceResponse()
		resp.Response = &teo.ModifyInferenceServiceResponseParams{
			RequestId: ptrStringTeoIs("fake-request-id-modify"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeInferenceServicesWithContext", func(_ context.Context, request *teo.DescribeInferenceServicesRequest) (*teo.DescribeInferenceServicesResponse, error) {
		return buildInferenceServiceResponse("is-update123"), nil
	})

	meta := newMockMetaTeoInferenceService()
	res := teo.ResourceTencentCloudTeoInferenceService()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-3fkff38fyw8s",
		"name":        "tf-example-inference-service",
		"listen_port": 8500,
		"description": "updated description",
	})
	d.SetId("zone-3fkff38fyw8s#is-update123")

	patches.ApplyMethodFunc(d, "HasChange", func(key string) bool {
		return key == "description"
	})

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.NotNil(t, capturedModifyRequest)
}

// TestTeoInferenceService_Delete verifies Delete issues OperateInferenceService with Operation=Delete.
func TestTeoInferenceService_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teo.Client{}
	patches.ApplyMethodReturn(newMockMetaTeoInferenceService().client, "UseTeoV20220901Client", teoClient)

	var capturedRequest *teo.OperateInferenceServiceRequest
	patches.ApplyMethodFunc(teoClient, "OperateInferenceServiceWithContext", func(_ context.Context, request *teo.OperateInferenceServiceRequest) (*teo.OperateInferenceServiceResponse, error) {
		capturedRequest = request
		assert.Equal(t, "zone-3fkff38fyw8s", *request.ZoneId)
		assert.Equal(t, "is-del123", *request.ServiceId)
		assert.Equal(t, "Delete", *request.Operation)

		resp := teo.NewOperateInferenceServiceResponse()
		resp.Response = &teo.OperateInferenceServiceResponseParams{
			RequestId: ptrStringTeoIs("fake-request-id-delete"),
		}
		return resp, nil
	})

	meta := newMockMetaTeoInferenceService()
	res := teo.ResourceTencentCloudTeoInferenceService()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":     "zone-3fkff38fyw8s",
		"name":        "tf-example-inference-service",
		"listen_port": 8500,
	})
	d.SetId("zone-3fkff38fyw8s#is-del123")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
	assert.NotNil(t, capturedRequest)
}

// TestTeoInferenceService_Schema validates the schema definition.
func TestTeoInferenceService_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoInferenceService()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.NotNil(t, res.Update)
	assert.NotNil(t, res.Delete)
	assert.NotNil(t, res.Importer)

	assert.Contains(t, res.Schema, "zone_id")
	assert.Contains(t, res.Schema, "name")
	assert.Contains(t, res.Schema, "listen_port")
	assert.Contains(t, res.Schema, "containers")
	assert.Contains(t, res.Schema, "resource_config")
	assert.Contains(t, res.Schema, "affinity_config")
	assert.Contains(t, res.Schema, "request_paths")
	assert.Contains(t, res.Schema, "description")

	assert.Contains(t, res.Schema, "service_id")
	assert.Contains(t, res.Schema, "status")
	assert.Contains(t, res.Schema, "scaling_status")
	assert.Contains(t, res.Schema, "current_instance_count")
	assert.Contains(t, res.Schema, "inference_url")
	assert.Contains(t, res.Schema, "create_time")
	assert.Contains(t, res.Schema, "update_time")

	zoneId := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneId.Type)
	assert.True(t, zoneId.Required)
	assert.True(t, zoneId.ForceNew)

	name := res.Schema["name"]
	assert.Equal(t, schema.TypeString, name.Type)
	assert.True(t, name.Required)
	assert.True(t, name.ForceNew)

	serviceId := res.Schema["service_id"]
	assert.Equal(t, schema.TypeString, serviceId.Type)
	assert.True(t, serviceId.Computed)
}
