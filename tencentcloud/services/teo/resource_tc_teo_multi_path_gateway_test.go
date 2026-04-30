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

// mockMetaForMultiPathGateway implements tccommon.ProviderMeta
type mockMetaForMultiPathGateway struct {
	client *connectivity.TencentCloudClient
}

func (m *mockMetaForMultiPathGateway) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mockMetaForMultiPathGateway{}

func newMockMetaForMultiPathGateway() *mockMetaForMultiPathGateway {
	return &mockMetaForMultiPathGateway{client: &connectivity.TencentCloudClient{}}
}

func ptrStrMPG(s string) *string {
	return &s
}

func ptrInt64MPG(i int64) *int64 {
	return &i
}

// go test ./tencentcloud/services/teo/ -run "TestMultiPathGateway" -v -count=1 -gcflags="all=-l"

// TestMultiPathGateway_Create_Success tests Create calls API and sets composite ID
func TestMultiPathGateway_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForMultiPathGateway().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateMultiPathGatewayWithContext", func(_ context.Context, request *teov20220901.CreateMultiPathGatewayRequest) (*teov20220901.CreateMultiPathGatewayResponse, error) {
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.Equal(t, "cloud", *request.GatewayType)
		assert.Equal(t, "test-gateway", *request.GatewayName)
		assert.Equal(t, "ap-guangzhou", *request.RegionId)
		resp := teov20220901.NewCreateMultiPathGatewayResponse()
		resp.Response = &teov20220901.CreateMultiPathGatewayResponseParams{
			GatewayId: ptrStrMPG("mpgw-abcdefgh"),
			RequestId: ptrStrMPG("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeMultiPathGateways", func(request *teov20220901.DescribeMultiPathGatewaysRequest) (*teov20220901.DescribeMultiPathGatewaysResponse, error) {
		resp := teov20220901.NewDescribeMultiPathGatewaysResponse()
		resp.Response = &teov20220901.DescribeMultiPathGatewaysResponseParams{
			Gateways: []*teov20220901.MultiPathGateway{
				{
					GatewayId:   ptrStrMPG("mpgw-abcdefgh"),
					GatewayName: ptrStrMPG("test-gateway"),
					GatewayType: ptrStrMPG("cloud"),
					Status:      ptrStrMPG("online"),
					RegionId:    ptrStrMPG("ap-guangzhou"),
				},
			},
			TotalCount: ptrInt64MPG(1),
			RequestId:  ptrStrMPG("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForMultiPathGateway()
	res := teo.ResourceTencentCloudTeoMultiPathGateway()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-12345678",
		"gateway_type": "cloud",
		"gateway_name": "test-gateway",
		"region_id":    "ap-guangzhou",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-12345678#mpgw-abcdefgh", d.Id())
}

// TestMultiPathGateway_Create_Private tests Create with private gateway type
func TestMultiPathGateway_Create_Private(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForMultiPathGateway().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateMultiPathGatewayWithContext", func(_ context.Context, request *teov20220901.CreateMultiPathGatewayRequest) (*teov20220901.CreateMultiPathGatewayResponse, error) {
		assert.Equal(t, "private", *request.GatewayType)
		assert.Equal(t, "1.2.3.4", *request.GatewayIP)
		assert.Equal(t, int64(8080), *request.GatewayPort)
		resp := teov20220901.NewCreateMultiPathGatewayResponse()
		resp.Response = &teov20220901.CreateMultiPathGatewayResponseParams{
			GatewayId: ptrStrMPG("mpgw-private123"),
			RequestId: ptrStrMPG("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeMultiPathGateways", func(request *teov20220901.DescribeMultiPathGatewaysRequest) (*teov20220901.DescribeMultiPathGatewaysResponse, error) {
		resp := teov20220901.NewDescribeMultiPathGatewaysResponse()
		resp.Response = &teov20220901.DescribeMultiPathGatewaysResponseParams{
			Gateways: []*teov20220901.MultiPathGateway{
				{
					GatewayId:   ptrStrMPG("mpgw-private123"),
					GatewayName: ptrStrMPG("test-private-gw"),
					GatewayType: ptrStrMPG("private"),
					GatewayPort: ptrInt64MPG(8080),
					GatewayIP:   ptrStrMPG("1.2.3.4"),
					Status:      ptrStrMPG("online"),
				},
			},
			TotalCount: ptrInt64MPG(1),
			RequestId:  ptrStrMPG("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForMultiPathGateway()
	res := teo.ResourceTencentCloudTeoMultiPathGateway()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-12345678",
		"gateway_type": "private",
		"gateway_name": "test-private-gw",
		"gateway_ip":   "1.2.3.4",
		"gateway_port": 8080,
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-12345678#mpgw-private123", d.Id())
}

// TestMultiPathGateway_Create_APIError tests Create handles API error
func TestMultiPathGateway_Create_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForMultiPathGateway().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateMultiPathGatewayWithContext", func(_ context.Context, request *teov20220901.CreateMultiPathGatewayRequest) (*teov20220901.CreateMultiPathGatewayResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid zone_id")
	})

	meta := newMockMetaForMultiPathGateway()
	res := teo.ResourceTencentCloudTeoMultiPathGateway()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-invalid",
		"gateway_type": "cloud",
		"gateway_name": "test-gateway",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestMultiPathGateway_Read_Success tests Read populates state from API
func TestMultiPathGateway_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForMultiPathGateway().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeMultiPathGateways", func(request *teov20220901.DescribeMultiPathGatewaysRequest) (*teov20220901.DescribeMultiPathGatewaysResponse, error) {
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		resp := teov20220901.NewDescribeMultiPathGatewaysResponse()
		resp.Response = &teov20220901.DescribeMultiPathGatewaysResponseParams{
			Gateways: []*teov20220901.MultiPathGateway{
				{
					GatewayId:   ptrStrMPG("mpgw-abcdefgh"),
					GatewayName: ptrStrMPG("test-gateway"),
					GatewayType: ptrStrMPG("cloud"),
					GatewayPort: ptrInt64MPG(443),
					Status:      ptrStrMPG("online"),
					GatewayIP:   ptrStrMPG(""),
					RegionId:    ptrStrMPG("ap-guangzhou"),
					NeedConfirm: ptrStrMPG("no"),
				},
			},
			TotalCount: ptrInt64MPG(1),
			RequestId:  ptrStrMPG("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForMultiPathGateway()
	res := teo.ResourceTencentCloudTeoMultiPathGateway()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-12345678",
		"gateway_type": "cloud",
		"gateway_name": "test-gateway",
	})
	d.SetId("zone-12345678#mpgw-abcdefgh")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "mpgw-abcdefgh", d.Get("gateway_id"))
	assert.Equal(t, "test-gateway", d.Get("gateway_name"))
	assert.Equal(t, "cloud", d.Get("gateway_type"))
	assert.Equal(t, "online", d.Get("status"))
	assert.Equal(t, "ap-guangzhou", d.Get("region_id"))
	assert.Equal(t, "no", d.Get("need_confirm"))
}

// TestMultiPathGateway_Read_NotFound tests Read handles resource not found
func TestMultiPathGateway_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForMultiPathGateway().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeMultiPathGateways", func(request *teov20220901.DescribeMultiPathGatewaysRequest) (*teov20220901.DescribeMultiPathGatewaysResponse, error) {
		resp := teov20220901.NewDescribeMultiPathGatewaysResponse()
		resp.Response = &teov20220901.DescribeMultiPathGatewaysResponseParams{
			Gateways:   []*teov20220901.MultiPathGateway{},
			TotalCount: ptrInt64MPG(0),
			RequestId:  ptrStrMPG("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForMultiPathGateway()
	res := teo.ResourceTencentCloudTeoMultiPathGateway()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-12345678",
		"gateway_type": "cloud",
		"gateway_name": "test-gateway",
	})
	d.SetId("zone-12345678#mpgw-notfound")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestMultiPathGateway_Update_Success tests Update calls Modify API
func TestMultiPathGateway_Update_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForMultiPathGateway().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "ModifyMultiPathGatewayWithContext", func(_ context.Context, request *teov20220901.ModifyMultiPathGatewayRequest) (*teov20220901.ModifyMultiPathGatewayResponse, error) {
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.Equal(t, "mpgw-abcdefgh", *request.GatewayId)
		assert.Equal(t, "updated-name", *request.GatewayName)
		resp := teov20220901.NewModifyMultiPathGatewayResponse()
		resp.Response = &teov20220901.ModifyMultiPathGatewayResponseParams{
			RequestId: ptrStrMPG("fake-request-id"),
		}
		return resp, nil
	})

	patches.ApplyMethodFunc(teoClient, "DescribeMultiPathGateways", func(request *teov20220901.DescribeMultiPathGatewaysRequest) (*teov20220901.DescribeMultiPathGatewaysResponse, error) {
		resp := teov20220901.NewDescribeMultiPathGatewaysResponse()
		resp.Response = &teov20220901.DescribeMultiPathGatewaysResponseParams{
			Gateways: []*teov20220901.MultiPathGateway{
				{
					GatewayId:   ptrStrMPG("mpgw-abcdefgh"),
					GatewayName: ptrStrMPG("updated-name"),
					GatewayType: ptrStrMPG("cloud"),
					Status:      ptrStrMPG("online"),
					RegionId:    ptrStrMPG("ap-guangzhou"),
				},
			},
			TotalCount: ptrInt64MPG(1),
			RequestId:  ptrStrMPG("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForMultiPathGateway()
	res := teo.ResourceTencentCloudTeoMultiPathGateway()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-12345678",
		"gateway_type": "cloud",
		"gateway_name": "updated-name",
	})
	d.SetId("zone-12345678#mpgw-abcdefgh")

	err := res.Update(d, meta)
	assert.NoError(t, err)
}

// TestMultiPathGateway_Delete_Success tests Delete calls API successfully
func TestMultiPathGateway_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForMultiPathGateway().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteMultiPathGatewayWithContext", func(_ context.Context, request *teov20220901.DeleteMultiPathGatewayRequest) (*teov20220901.DeleteMultiPathGatewayResponse, error) {
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		assert.Equal(t, "mpgw-abcdefgh", *request.GatewayId)
		resp := teov20220901.NewDeleteMultiPathGatewayResponse()
		resp.Response = &teov20220901.DeleteMultiPathGatewayResponseParams{
			RequestId: ptrStrMPG("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForMultiPathGateway()
	res := teo.ResourceTencentCloudTeoMultiPathGateway()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-12345678",
		"gateway_type": "cloud",
		"gateway_name": "test-gateway",
	})
	d.SetId("zone-12345678#mpgw-abcdefgh")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestMultiPathGateway_Delete_APIError tests Delete handles API error
func TestMultiPathGateway_Delete_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForMultiPathGateway().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DeleteMultiPathGatewayWithContext", func(_ context.Context, request *teov20220901.DeleteMultiPathGatewayRequest) (*teov20220901.DeleteMultiPathGatewayResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=Gateway not found")
	})

	meta := newMockMetaForMultiPathGateway()
	res := teo.ResourceTencentCloudTeoMultiPathGateway()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-12345678",
		"gateway_type": "cloud",
		"gateway_name": "test-gateway",
	})
	d.SetId("zone-12345678#mpgw-abcdefgh")

	err := res.Delete(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
}

// TestMultiPathGateway_Schema validates schema definition
func TestMultiPathGateway_Schema(t *testing.T) {
	res := teo.ResourceTencentCloudTeoMultiPathGateway()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Create)
	assert.NotNil(t, res.Read)
	assert.NotNil(t, res.Update)
	assert.NotNil(t, res.Delete)
	assert.NotNil(t, res.Importer)

	assert.Contains(t, res.Schema, "zone_id")
	assert.Contains(t, res.Schema, "gateway_type")
	assert.Contains(t, res.Schema, "gateway_name")
	assert.Contains(t, res.Schema, "gateway_port")
	assert.Contains(t, res.Schema, "region_id")
	assert.Contains(t, res.Schema, "gateway_ip")
	assert.Contains(t, res.Schema, "gateway_id")
	assert.Contains(t, res.Schema, "status")
	assert.Contains(t, res.Schema, "need_confirm")

	zoneId := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneId.Type)
	assert.True(t, zoneId.Required)
	assert.True(t, zoneId.ForceNew)

	gatewayType := res.Schema["gateway_type"]
	assert.Equal(t, schema.TypeString, gatewayType.Type)
	assert.True(t, gatewayType.Required)
	assert.True(t, gatewayType.ForceNew)

	gatewayName := res.Schema["gateway_name"]
	assert.Equal(t, schema.TypeString, gatewayName.Type)
	assert.True(t, gatewayName.Required)
	assert.False(t, gatewayName.ForceNew)

	gatewayPort := res.Schema["gateway_port"]
	assert.Equal(t, schema.TypeInt, gatewayPort.Type)
	assert.True(t, gatewayPort.Optional)
	assert.True(t, gatewayPort.Computed)

	regionId := res.Schema["region_id"]
	assert.Equal(t, schema.TypeString, regionId.Type)
	assert.True(t, regionId.Optional)
	assert.True(t, regionId.Computed)
	assert.True(t, regionId.ForceNew)

	gatewayIp := res.Schema["gateway_ip"]
	assert.Equal(t, schema.TypeString, gatewayIp.Type)
	assert.True(t, gatewayIp.Optional)
	assert.True(t, gatewayIp.Computed)

	gatewayId := res.Schema["gateway_id"]
	assert.Equal(t, schema.TypeString, gatewayId.Type)
	assert.True(t, gatewayId.Computed)
	assert.False(t, gatewayId.Required)
	assert.False(t, gatewayId.Optional)

	status := res.Schema["status"]
	assert.Equal(t, schema.TypeString, status.Type)
	assert.True(t, status.Computed)

	needConfirm := res.Schema["need_confirm"]
	assert.Equal(t, schema.TypeString, needConfirm.Type)
	assert.True(t, needConfirm.Computed)
}
