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
	svcteo "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// mockMetaMultiPathGateway implements tccommon.ProviderMeta
type mockMetaMultiPathGateway struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaMultiPathGateway) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaMultiPathGateway{}

func newMockMetaMultiPathGateway() *mockMetaMultiPathGateway {
	return &mockMetaMultiPathGateway{client: &connectivity.TencentCloudClient{}}
}

func ptrStringMPG(s string) *string {
	return &s
}

func ptrInt64MPG(n int64) *int64 {
	return &n
}

// go test ./tencentcloud/services/teo/ -run "TestTeoMultiPathGateway" -v -count=1 -gcflags="all=-l"

// TestTeoMultiPathGateway_Create_CloudType tests creating a cloud type gateway
func TestTeoMultiPathGateway_Create_CloudType(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaMultiPathGateway().client, "UseTeoV20220901Client", teoClient)

	// Mock CreateMultiPathGatewayWithContext
	patches.ApplyMethodFunc(teoClient, "CreateMultiPathGatewayWithContext", func(_ context.Context, _ *teov20220901.CreateMultiPathGatewayRequest) (*teov20220901.CreateMultiPathGatewayResponse, error) {
		resp := teov20220901.NewCreateMultiPathGatewayResponse()
		resp.Response = &teov20220901.CreateMultiPathGatewayResponseParams{
			GatewayId: ptrStringMPG("gw-12345678"),
			RequestId: ptrStringMPG("fake-request-id"),
		}
		return resp, nil
	})

	// Mock TeoService.DescribeTeoMultiPathGatewayById for StateChangeConf and Read
	patches.ApplyMethodFunc(&svcteo.TeoService{}, "DescribeTeoMultiPathGatewayById", func(_ context.Context, zoneId string, gatewayId string) (*teov20220901.MultiPathGateway, error) {
		return &teov20220901.MultiPathGateway{
			GatewayId:   ptrStringMPG("gw-12345678"),
			GatewayName: ptrStringMPG("test-cloud-gw"),
			GatewayType: ptrStringMPG("cloud"),
			GatewayPort: ptrInt64MPG(8080),
			Status:      ptrStringMPG("online"),
			RegionId:    ptrStringMPG("ap-guangzhou"),
			NeedConfirm: ptrStringMPG("false"),
		}, nil
	})

	meta := newMockMetaMultiPathGateway()
	res := svcteo.ResourceTencentCloudTeoMultiPathGateway()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-test1234",
		"gateway_type": "cloud",
		"gateway_name": "test-cloud-gw",
		"gateway_port": 8080,
		"region_id":    "ap-guangzhou",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)

	// Verify composite ID
	assert.Equal(t, "zone-test1234#gw-12345678", d.Id())

	// Verify gateway_id computed field
	assert.Equal(t, "gw-12345678", d.Get("gateway_id").(string))

	// Verify status
	assert.Equal(t, "online", d.Get("status").(string))
}

// TestTeoMultiPathGateway_Create_PrivateType tests creating a private type gateway
func TestTeoMultiPathGateway_Create_PrivateType(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaMultiPathGateway().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateMultiPathGatewayWithContext", func(_ context.Context, _ *teov20220901.CreateMultiPathGatewayRequest) (*teov20220901.CreateMultiPathGatewayResponse, error) {
		resp := teov20220901.NewCreateMultiPathGatewayResponse()
		resp.Response = &teov20220901.CreateMultiPathGatewayResponseParams{
			GatewayId: ptrStringMPG("gw-87654321"),
			RequestId: ptrStringMPG("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(&svcteo.TeoService{}, "DescribeTeoMultiPathGatewayById", func(_ context.Context, zoneId string, gatewayId string) (*teov20220901.MultiPathGateway, error) {
		return &teov20220901.MultiPathGateway{
			GatewayId:   ptrStringMPG("gw-87654321"),
			GatewayName: ptrStringMPG("test-private-gw"),
			GatewayType: ptrStringMPG("private"),
			GatewayPort: ptrInt64MPG(9090),
			GatewayIP:   ptrStringMPG("10.0.0.1"),
			Status:      ptrStringMPG("online"),
			NeedConfirm: ptrStringMPG("false"),
		}, nil
	})

	meta := newMockMetaMultiPathGateway()
	res := svcteo.ResourceTencentCloudTeoMultiPathGateway()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-test1234",
		"gateway_type": "private",
		"gateway_name": "test-private-gw",
		"gateway_port": 9090,
		"gateway_ip":   "10.0.0.1",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)

	assert.Equal(t, "zone-test1234#gw-87654321", d.Id())
	assert.Equal(t, "gw-87654321", d.Get("gateway_id").(string))
	assert.Equal(t, "private", d.Get("gateway_type").(string))
	assert.Equal(t, "10.0.0.1", d.Get("gateway_ip").(string))
}

// TestTeoMultiPathGateway_Create_NilGatewayId tests Create when response GatewayId is nil
func TestTeoMultiPathGateway_Create_NilGatewayId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaMultiPathGateway().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateMultiPathGatewayWithContext", func(_ context.Context, _ *teov20220901.CreateMultiPathGatewayRequest) (*teov20220901.CreateMultiPathGatewayResponse, error) {
		resp := teov20220901.NewCreateMultiPathGatewayResponse()
		resp.Response = &teov20220901.CreateMultiPathGatewayResponseParams{
			GatewayId: nil,
			RequestId: ptrStringMPG("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaMultiPathGateway()
	res := svcteo.ResourceTencentCloudTeoMultiPathGateway()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-test1234",
		"gateway_type": "cloud",
		"gateway_name": "test-gw",
		"gateway_port": 8080,
		"region_id":    "ap-guangzhou",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "GatewayId is nil")
}

// TestTeoMultiPathGateway_Create_APIError tests Create handles API error
func TestTeoMultiPathGateway_Create_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaMultiPathGateway().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateMultiPathGatewayWithContext", func(_ context.Context, _ *teov20220901.CreateMultiPathGatewayRequest) (*teov20220901.CreateMultiPathGatewayResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid zone_id")
	})

	meta := newMockMetaMultiPathGateway()
	res := svcteo.ResourceTencentCloudTeoMultiPathGateway()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-invalid",
		"gateway_type": "cloud",
		"gateway_name": "test-gw",
		"gateway_port": 8080,
		"region_id":    "ap-guangzhou",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestTeoMultiPathGateway_Read tests Read populates all fields correctly
func TestTeoMultiPathGateway_Read(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyMethodFunc(&svcteo.TeoService{}, "DescribeTeoMultiPathGatewayById", func(_ context.Context, zoneId string, gatewayId string) (*teov20220901.MultiPathGateway, error) {
		assert.Equal(t, "zone-test1234", zoneId)
		assert.Equal(t, "gw-12345678", gatewayId)
		return &teov20220901.MultiPathGateway{
			GatewayId:   ptrStringMPG("gw-12345678"),
			GatewayName: ptrStringMPG("test-gw"),
			GatewayType: ptrStringMPG("cloud"),
			GatewayPort: ptrInt64MPG(8080),
			Status:      ptrStringMPG("online"),
			RegionId:    ptrStringMPG("ap-guangzhou"),
			NeedConfirm: ptrStringMPG("false"),
		}, nil
	})

	meta := newMockMetaMultiPathGateway()
	res := svcteo.ResourceTencentCloudTeoMultiPathGateway()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-test1234",
		"gateway_type": "cloud",
		"gateway_name": "test-gw",
		"gateway_port": 8080,
	})
	d.SetId("zone-test1234#gw-12345678")

	err := res.Read(d, meta)
	assert.NoError(t, err)

	assert.Equal(t, "gw-12345678", d.Get("gateway_id").(string))
	assert.Equal(t, "cloud", d.Get("gateway_type").(string))
	assert.Equal(t, "test-gw", d.Get("gateway_name").(string))
	assert.Equal(t, 8080, d.Get("gateway_port").(int))
	assert.Equal(t, "online", d.Get("status").(string))
	assert.Equal(t, "ap-guangzhou", d.Get("region_id").(string))
	assert.Equal(t, "false", d.Get("need_confirm").(string))
}

// TestTeoMultiPathGateway_Read_NotFound tests Read when gateway is not found
func TestTeoMultiPathGateway_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyMethodFunc(&svcteo.TeoService{}, "DescribeTeoMultiPathGatewayById", func(_ context.Context, zoneId string, gatewayId string) (*teov20220901.MultiPathGateway, error) {
		return nil, nil
	})

	meta := newMockMetaMultiPathGateway()
	res := svcteo.ResourceTencentCloudTeoMultiPathGateway()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-test1234",
		"gateway_type": "cloud",
		"gateway_name": "test-gw",
		"gateway_port": 8080,
	})
	d.SetId("zone-test1234#gw-nonexistent")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestTeoMultiPathGateway_Read_NilFields tests Read with nil fields in response
func TestTeoMultiPathGateway_Read_NilFields(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyMethodFunc(&svcteo.TeoService{}, "DescribeTeoMultiPathGatewayById", func(_ context.Context, zoneId string, gatewayId string) (*teov20220901.MultiPathGateway, error) {
		return &teov20220901.MultiPathGateway{
			GatewayId:   ptrStringMPG("gw-12345678"),
			GatewayName: ptrStringMPG("test-gw"),
			GatewayType: ptrStringMPG("cloud"),
			GatewayPort: ptrInt64MPG(8080),
			Status:      nil,
			RegionId:    nil,
			NeedConfirm: nil,
			GatewayIP:   nil,
		}, nil
	})

	meta := newMockMetaMultiPathGateway()
	res := svcteo.ResourceTencentCloudTeoMultiPathGateway()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-test1234",
		"gateway_type": "cloud",
		"gateway_name": "test-gw",
		"gateway_port": 8080,
	})
	d.SetId("zone-test1234#gw-12345678")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-test1234#gw-12345678", d.Id())
}

// TestTeoMultiPathGateway_Update tests Update with gateway_name change
func TestTeoMultiPathGateway_Update(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaMultiPathGateway().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyMultiPathGatewayWithContext", func(_ context.Context, request *teov20220901.ModifyMultiPathGatewayRequest) (*teov20220901.ModifyMultiPathGatewayResponse, error) {
		assert.Equal(t, "zone-test1234", *request.ZoneId)
		assert.Equal(t, "gw-12345678", *request.GatewayId)
		assert.Equal(t, "updated-gw-name", *request.GatewayName)
		resp := teov20220901.NewModifyMultiPathGatewayResponse()
		resp.Response = &teov20220901.ModifyMultiPathGatewayResponseParams{
			RequestId: ptrStringMPG("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(&svcteo.TeoService{}, "DescribeTeoMultiPathGatewayById", func(_ context.Context, zoneId string, gatewayId string) (*teov20220901.MultiPathGateway, error) {
		return &teov20220901.MultiPathGateway{
			GatewayId:   ptrStringMPG("gw-12345678"),
			GatewayName: ptrStringMPG("updated-gw-name"),
			GatewayType: ptrStringMPG("cloud"),
			GatewayPort: ptrInt64MPG(8080),
			Status:      ptrStringMPG("online"),
			RegionId:    ptrStringMPG("ap-guangzhou"),
			NeedConfirm: ptrStringMPG("false"),
		}, nil
	})

	meta := newMockMetaMultiPathGateway()
	res := svcteo.ResourceTencentCloudTeoMultiPathGateway()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-test1234",
		"gateway_type": "cloud",
		"gateway_name": "updated-gw-name",
		"gateway_port": 8080,
		"region_id":    "ap-guangzhou",
	})
	d.SetId("zone-test1234#gw-12345678")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

// TestTeoMultiPathGateway_Update_APIError tests Update handles API error
func TestTeoMultiPathGateway_Update_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaMultiPathGateway().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyMultiPathGatewayWithContext", func(_ context.Context, _ *teov20220901.ModifyMultiPathGatewayRequest) (*teov20220901.ModifyMultiPathGatewayResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid gateway_name")
	})

	meta := newMockMetaMultiPathGateway()
	res := svcteo.ResourceTencentCloudTeoMultiPathGateway()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-test1234",
		"gateway_type": "cloud",
		"gateway_name": "updated-gw-name",
		"gateway_port": 8080,
		"region_id":    "ap-guangzhou",
	})
	d.SetId("zone-test1234#gw-12345678")

	err := res.Update(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestTeoMultiPathGateway_Delete tests successful deletion
func TestTeoMultiPathGateway_Delete(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaMultiPathGateway().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(&svcteo.TeoService{}, "DescribeTeoMultiPathGatewayById", func(_ context.Context, zoneId string, gatewayId string) (*teov20220901.MultiPathGateway, error) {
		return &teov20220901.MultiPathGateway{
			GatewayId: ptrStringMPG("gw-12345678"),
			Status:    ptrStringMPG("online"),
		}, nil
	})

	patches.ApplyMethodFunc(teoClient, "DeleteMultiPathGatewayWithContext", func(_ context.Context, request *teov20220901.DeleteMultiPathGatewayRequest) (*teov20220901.DeleteMultiPathGatewayResponse, error) {
		assert.Equal(t, "zone-test1234", *request.ZoneId)
		assert.Equal(t, "gw-12345678", *request.GatewayId)
		resp := teov20220901.NewDeleteMultiPathGatewayResponse()
		resp.Response = &teov20220901.DeleteMultiPathGatewayResponseParams{
			RequestId: ptrStringMPG("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaMultiPathGateway()
	res := svcteo.ResourceTencentCloudTeoMultiPathGateway()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-test1234",
		"gateway_type": "cloud",
		"gateway_name": "test-gw",
		"gateway_port": 8080,
	})
	d.SetId("zone-test1234#gw-12345678")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestTeoMultiPathGateway_Delete_NotFound tests deletion when gateway is already gone
func TestTeoMultiPathGateway_Delete_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	patches.ApplyMethodFunc(&svcteo.TeoService{}, "DescribeTeoMultiPathGatewayById", func(_ context.Context, zoneId string, gatewayId string) (*teov20220901.MultiPathGateway, error) {
		return nil, nil
	})

	meta := newMockMetaMultiPathGateway()
	res := svcteo.ResourceTencentCloudTeoMultiPathGateway()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-test1234",
		"gateway_type": "cloud",
		"gateway_name": "test-gw",
		"gateway_port": 8080,
	})
	d.SetId("zone-test1234#gw-nonexistent")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestTeoMultiPathGateway_Delete_APIError tests deletion handles API error
func TestTeoMultiPathGateway_Delete_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaMultiPathGateway().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(&svcteo.TeoService{}, "DescribeTeoMultiPathGatewayById", func(_ context.Context, zoneId string, gatewayId string) (*teov20220901.MultiPathGateway, error) {
		return &teov20220901.MultiPathGateway{
			GatewayId: ptrStringMPG("gw-12345678"),
			Status:    ptrStringMPG("online"),
		}, nil
	})

	patches.ApplyMethodFunc(teoClient, "DeleteMultiPathGatewayWithContext", func(_ context.Context, _ *teov20220901.DeleteMultiPathGatewayRequest) (*teov20220901.DeleteMultiPathGatewayResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=Gateway not found")
	})

	meta := newMockMetaMultiPathGateway()
	res := svcteo.ResourceTencentCloudTeoMultiPathGateway()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-test1234",
		"gateway_type": "cloud",
		"gateway_name": "test-gw",
		"gateway_port": 8080,
	})
	d.SetId("zone-test1234#gw-12345678")

	err := res.Delete(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
}

// TestTeoMultiPathGateway_Schema tests schema definition
func TestTeoMultiPathGateway_Schema(t *testing.T) {
	res := svcteo.ResourceTencentCloudTeoMultiPathGateway()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.NotNil(t, res.Update)
	assert.NotNil(t, res.Delete)
	assert.NotNil(t, res.Importer)

	// Required fields
	assert.Contains(t, res.Schema, "zone_id")
	zoneIdSchema := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneIdSchema.Type)
	assert.True(t, zoneIdSchema.Required)
	assert.True(t, zoneIdSchema.ForceNew)

	assert.Contains(t, res.Schema, "gateway_type")
	gatewayTypeSchema := res.Schema["gateway_type"]
	assert.Equal(t, schema.TypeString, gatewayTypeSchema.Type)
	assert.True(t, gatewayTypeSchema.Required)
	assert.True(t, gatewayTypeSchema.ForceNew)

	assert.Contains(t, res.Schema, "gateway_name")
	gatewayNameSchema := res.Schema["gateway_name"]
	assert.Equal(t, schema.TypeString, gatewayNameSchema.Type)
	assert.True(t, gatewayNameSchema.Required)
	assert.False(t, gatewayNameSchema.ForceNew)

	assert.Contains(t, res.Schema, "gateway_port")
	gatewayPortSchema := res.Schema["gateway_port"]
	assert.Equal(t, schema.TypeInt, gatewayPortSchema.Type)
	assert.True(t, gatewayPortSchema.Required)
	assert.False(t, gatewayPortSchema.ForceNew)

	// Optional fields
	assert.Contains(t, res.Schema, "region_id")
	regionIdSchema := res.Schema["region_id"]
	assert.Equal(t, schema.TypeString, regionIdSchema.Type)
	assert.True(t, regionIdSchema.Optional)
	assert.True(t, regionIdSchema.ForceNew)

	assert.Contains(t, res.Schema, "gateway_ip")
	gatewayIpSchema := res.Schema["gateway_ip"]
	assert.Equal(t, schema.TypeString, gatewayIpSchema.Type)
	assert.True(t, gatewayIpSchema.Optional)

	// Computed fields
	assert.Contains(t, res.Schema, "gateway_id")
	gatewayIdSchema := res.Schema["gateway_id"]
	assert.Equal(t, schema.TypeString, gatewayIdSchema.Type)
	assert.True(t, gatewayIdSchema.Computed)

	assert.Contains(t, res.Schema, "status")
	statusSchema := res.Schema["status"]
	assert.Equal(t, schema.TypeString, statusSchema.Type)
	assert.True(t, statusSchema.Computed)

	assert.Contains(t, res.Schema, "need_confirm")
	needConfirmSchema := res.Schema["need_confirm"]
	assert.Equal(t, schema.TypeString, needConfirmSchema.Type)
	assert.True(t, needConfirmSchema.Computed)
}
