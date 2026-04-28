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

// TestMultiPathGateway_Read_Success tests Read populates all fields
func TestMultiPathGateway_Read_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForMultiPathGateway().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeMultiPathGateway", func(request *teov20220901.DescribeMultiPathGatewayRequest) (*teov20220901.DescribeMultiPathGatewayResponse, error) {
		resp := teov20220901.NewDescribeMultiPathGatewayResponse()
		resp.Response = &teov20220901.DescribeMultiPathGatewayResponseParams{
			GatewayDetail: &teov20220901.MultiPathGateway{
				GatewayId:   ptrStrMPG("gw-abc123"),
				GatewayType: ptrStrMPG("cloud"),
				GatewayName: ptrStrMPG("test-gateway"),
				GatewayPort: ptrInt64MPG(8080),
				RegionId:    ptrStrMPG("ap-guangzhou"),
				Status:      ptrStrMPG("online"),
				NeedConfirm: ptrStrMPG("false"),
				Lines: []*teov20220901.MultiPathGatewayLine{
					{
						LineId:      ptrStrMPG("line-0"),
						LineType:    ptrStrMPG("direct"),
						LineAddress: ptrStrMPG("1.2.3.4:8080"),
					},
					{
						LineId:      ptrStrMPG("line-1"),
						LineType:    ptrStrMPG("proxy"),
						LineAddress: ptrStrMPG("5.6.7.8:8080"),
						ProxyId:     ptrStrMPG("proxy-001"),
						RuleId:      ptrStrMPG("rule-001"),
					},
				},
			},
			RequestId: ptrStrMPG("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForMultiPathGateway()
	res := teo.ResourceTencentCloudTeoMultiPathGateway()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-test123",
		"gateway_type": "cloud",
		"gateway_name": "test-gateway",
		"gateway_port": 8080,
	})
	d.SetId("zone-test123#gw-abc123")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "test-gateway", d.Get("gateway_name"))
	assert.Equal(t, "cloud", d.Get("gateway_type"))
	assert.Equal(t, 8080, d.Get("gateway_port"))
	assert.Equal(t, "ap-guangzhou", d.Get("region_id"))
	assert.Equal(t, "gw-abc123", d.Get("gateway_id"))
	assert.Equal(t, "online", d.Get("status"))
	assert.Equal(t, "false", d.Get("need_confirm"))

	lines := d.Get("lines").([]interface{})
	assert.Equal(t, 2, len(lines))
}

// TestMultiPathGateway_Read_NotFound tests Read handles resource not found
func TestMultiPathGateway_Read_NotFound(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForMultiPathGateway().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeMultiPathGateway", func(request *teov20220901.DescribeMultiPathGatewayRequest) (*teov20220901.DescribeMultiPathGatewayResponse, error) {
		resp := teov20220901.NewDescribeMultiPathGatewayResponse()
		resp.Response = &teov20220901.DescribeMultiPathGatewayResponseParams{
			GatewayDetail: nil,
			RequestId:     ptrStrMPG("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForMultiPathGateway()
	res := teo.ResourceTencentCloudTeoMultiPathGateway()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-test123",
		"gateway_type": "cloud",
		"gateway_name": "test-gateway",
		"gateway_port": 8080,
	})
	d.SetId("zone-test123#gw-abc123")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "", d.Id())
}

// TestMultiPathGateway_Create_Success tests Create with cloud type gateway
func TestMultiPathGateway_Create_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForMultiPathGateway().client, "UseTeoV20220901Client", teoClient)

	// Mock CreateMultiPathGatewayWithContext
	patches.ApplyMethodFunc(teoClient, "CreateMultiPathGatewayWithContext", func(ctx context.Context, request *teov20220901.CreateMultiPathGatewayRequest) (*teov20220901.CreateMultiPathGatewayResponse, error) {
		resp := teov20220901.NewCreateMultiPathGatewayResponse()
		resp.Response = &teov20220901.CreateMultiPathGatewayResponseParams{
			GatewayId: ptrStrMPG("gw-new123"),
			RequestId: ptrStrMPG("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeMultiPathGateway for Read after Create
	patches.ApplyMethodFunc(teoClient, "DescribeMultiPathGateway", func(request *teov20220901.DescribeMultiPathGatewayRequest) (*teov20220901.DescribeMultiPathGatewayResponse, error) {
		resp := teov20220901.NewDescribeMultiPathGatewayResponse()
		resp.Response = &teov20220901.DescribeMultiPathGatewayResponseParams{
			GatewayDetail: &teov20220901.MultiPathGateway{
				GatewayId:   ptrStrMPG("gw-new123"),
				GatewayType: ptrStrMPG("cloud"),
				GatewayName: ptrStrMPG("test-gateway"),
				GatewayPort: ptrInt64MPG(8080),
				RegionId:    ptrStrMPG("ap-guangzhou"),
				Status:      ptrStrMPG("creating"),
				NeedConfirm: ptrStrMPG("false"),
			},
			RequestId: ptrStrMPG("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForMultiPathGateway()
	res := teo.ResourceTencentCloudTeoMultiPathGateway()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-test123",
		"gateway_type": "cloud",
		"gateway_name": "test-gateway",
		"gateway_port": 8080,
		"region_id":    "ap-guangzhou",
	})

	err := res.Create(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-test123#gw-new123", d.Id())
}

// TestMultiPathGateway_Create_EmptyGatewayId tests Create returns error when GatewayId is empty
func TestMultiPathGateway_Create_EmptyGatewayId(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForMultiPathGateway().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "CreateMultiPathGatewayWithContext", func(ctx context.Context, request *teov20220901.CreateMultiPathGatewayRequest) (*teov20220901.CreateMultiPathGatewayResponse, error) {
		resp := teov20220901.NewCreateMultiPathGatewayResponse()
		resp.Response = &teov20220901.CreateMultiPathGatewayResponseParams{
			GatewayId: ptrStrMPG(""),
			RequestId: ptrStrMPG("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForMultiPathGateway()
	res := teo.ResourceTencentCloudTeoMultiPathGateway()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-test123",
		"gateway_type": "cloud",
		"gateway_name": "test-gateway",
		"gateway_port": 8080,
		"region_id":    "ap-guangzhou",
	})

	err := res.Create(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "GatewayId is empty")
}

// TestMultiPathGateway_Update_Success tests Update with gateway_name change
func TestMultiPathGateway_Update_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForMultiPathGateway().client, "UseTeoV20220901Client", teoClient)

	// Mock ModifyMultiPathGatewayWithContext
	patches.ApplyMethodFunc(teoClient, "ModifyMultiPathGatewayWithContext", func(ctx context.Context, request *teov20220901.ModifyMultiPathGatewayRequest) (*teov20220901.ModifyMultiPathGatewayResponse, error) {
		resp := teov20220901.NewModifyMultiPathGatewayResponse()
		resp.Response = &teov20220901.ModifyMultiPathGatewayResponseParams{
			RequestId: ptrStrMPG("fake-request-id"),
		}
		return resp, nil
	})

	// Mock DescribeMultiPathGateway for Read after Update
	patches.ApplyMethodFunc(teoClient, "DescribeMultiPathGateway", func(request *teov20220901.DescribeMultiPathGatewayRequest) (*teov20220901.DescribeMultiPathGatewayResponse, error) {
		resp := teov20220901.NewDescribeMultiPathGatewayResponse()
		resp.Response = &teov20220901.DescribeMultiPathGatewayResponseParams{
			GatewayDetail: &teov20220901.MultiPathGateway{
				GatewayId:   ptrStrMPG("gw-abc123"),
				GatewayType: ptrStrMPG("cloud"),
				GatewayName: ptrStrMPG("test-gateway-updated"),
				GatewayPort: ptrInt64MPG(8080),
				RegionId:    ptrStrMPG("ap-guangzhou"),
				Status:      ptrStrMPG("online"),
				NeedConfirm: ptrStrMPG("false"),
			},
			RequestId: ptrStrMPG("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForMultiPathGateway()
	res := teo.ResourceTencentCloudTeoMultiPathGateway()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-test123",
		"gateway_type": "cloud",
		"gateway_name": "test-gateway-updated",
		"gateway_port": 8080,
		"region_id":    "ap-guangzhou",
	})
	d.SetId("zone-test123#gw-abc123")

	err := res.Update(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "test-gateway-updated", d.Get("gateway_name"))
}

// TestMultiPathGateway_Delete_Success tests Delete operation
func TestMultiPathGateway_Delete_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForMultiPathGateway().client, "UseTeoV20220901Client", teoClient)

	// Mock DeleteMultiPathGateway
	patches.ApplyMethodFunc(teoClient, "DeleteMultiPathGateway", func(request *teov20220901.DeleteMultiPathGatewayRequest) (*teov20220901.DeleteMultiPathGatewayResponse, error) {
		resp := teov20220901.NewDeleteMultiPathGatewayResponse()
		resp.Response = &teov20220901.DeleteMultiPathGatewayResponseParams{
			RequestId: ptrStrMPG("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForMultiPathGateway()
	res := teo.ResourceTencentCloudTeoMultiPathGateway()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-test123",
		"gateway_type": "cloud",
		"gateway_name": "test-gateway",
		"gateway_port": 8080,
	})
	d.SetId("zone-test123#gw-abc123")

	err := res.Delete(d, meta)
	assert.NoError(t, err)
}

// TestMultiPathGateway_Read_APIError tests Read handles API error
func TestMultiPathGateway_Read_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForMultiPathGateway().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeMultiPathGateway", func(request *teov20220901.DescribeMultiPathGatewayRequest) (*teov20220901.DescribeMultiPathGatewayResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=ResourceNotFound, Message=Gateway not found")
	})

	meta := newMockMetaForMultiPathGateway()
	res := teo.ResourceTencentCloudTeoMultiPathGateway()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-test123",
		"gateway_type": "cloud",
		"gateway_name": "test-gateway",
		"gateway_port": 8080,
	})
	d.SetId("zone-test123#gw-abc123")

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "ResourceNotFound")
}

// TestMultiPathGateway_Read_PrivateType tests Read with private type gateway
func TestMultiPathGateway_Read_PrivateType(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(newMockMetaForMultiPathGateway().client, "UseTeoV20220901Client", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeMultiPathGateway", func(request *teov20220901.DescribeMultiPathGatewayRequest) (*teov20220901.DescribeMultiPathGatewayResponse, error) {
		resp := teov20220901.NewDescribeMultiPathGatewayResponse()
		resp.Response = &teov20220901.DescribeMultiPathGatewayResponseParams{
			GatewayDetail: &teov20220901.MultiPathGateway{
				GatewayId:   ptrStrMPG("gw-xyz789"),
				GatewayType: ptrStrMPG("private"),
				GatewayName: ptrStrMPG("test-private-gw"),
				GatewayPort: ptrInt64MPG(9090),
				GatewayIP:   ptrStrMPG("10.0.0.1"),
				Status:      ptrStrMPG("online"),
				NeedConfirm: ptrStrMPG("true"),
			},
			RequestId: ptrStrMPG("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMetaForMultiPathGateway()
	res := teo.ResourceTencentCloudTeoMultiPathGateway()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":      "zone-test123",
		"gateway_type": "private",
		"gateway_name": "test-private-gw",
		"gateway_port": 9090,
	})
	d.SetId("zone-test123#gw-xyz789")

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "private", d.Get("gateway_type"))
	assert.Equal(t, "10.0.0.1", d.Get("gateway_ip"))
	assert.Equal(t, 9090, d.Get("gateway_port"))
	assert.Equal(t, "true", d.Get("need_confirm"))
}
