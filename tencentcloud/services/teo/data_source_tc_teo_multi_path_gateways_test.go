package teo_test

import (
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

// go test ./tencentcloud/services/teo/ -run "TestTeoMultiPathGateways" -v -count=1 -gcflags="all=-l"

// TestTeoMultiPathGatewaysDataSource_Success tests Read calls API and sets gateways
func TestTeoMultiPathGatewaysDataSource_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseTeoClient", teoClient)

	// Patch DescribeMultiPathGateways to return two gateways
	patches.ApplyMethodFunc(teoClient, "DescribeMultiPathGateways", func(request *teov20220901.DescribeMultiPathGatewaysRequest) (*teov20220901.DescribeMultiPathGatewaysResponse, error) {
		assert.Equal(t, "zone-12345678", *request.ZoneId)
		resp := teov20220901.NewDescribeMultiPathGatewaysResponse()
		resp.Response = &teov20220901.DescribeMultiPathGatewaysResponseParams{
			Gateways: []*teov20220901.MultiPathGateway{
				{
					GatewayId:   ptrMPGString("gateway-001"),
					GatewayName: ptrMPGString("test-gateway-1"),
					GatewayType: ptrMPGString("cloud"),
					GatewayPort: ptrMPGInt64(8080),
					Status:      ptrMPGString("online"),
					GatewayIP:   ptrMPGString("1.2.3.4"),
					RegionId:    ptrMPGString("ap-guangzhou"),
					NeedConfirm: ptrMPGString("false"),
					Lines: []*teov20220901.MultiPathGatewayLine{
						{
							LineId:      ptrMPGString("line-0"),
							LineType:    ptrMPGString("direct"),
							LineAddress: ptrMPGString("1.2.3.4:8080"),
						},
					},
				},
				{
					GatewayId:   ptrMPGString("gateway-002"),
					GatewayName: ptrMPGString("test-gateway-2"),
					GatewayType: ptrMPGString("private"),
					GatewayPort: ptrMPGInt64(9090),
					Status:      ptrMPGString("offline"),
					GatewayIP:   ptrMPGString("5.6.7.8"),
					RegionId:    ptrMPGString("ap-shanghai"),
					NeedConfirm: ptrMPGString("true"),
				},
			},
			TotalCount: ptrMPGInt64(2),
			RequestId:  ptrMPGString("fake-request-id"),
		}
		return resp, nil
	})

	meta := &mpgMockMeta{client: mockClient}
	res := teo.DataSourceTencentCloudTeoMultiPathGateways()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-12345678",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.Equal(t, "zone-12345678", d.Id())

	gateways := d.Get("gateways").([]interface{})
	assert.Len(t, gateways, 2)

	gw0 := gateways[0].(map[string]interface{})
	assert.Equal(t, "gateway-001", gw0["gateway_id"])
	assert.Equal(t, "test-gateway-1", gw0["gateway_name"])
	assert.Equal(t, "cloud", gw0["gateway_type"])
	assert.Equal(t, int(8080), gw0["gateway_port"])
	assert.Equal(t, "online", gw0["status"])
	assert.Equal(t, "1.2.3.4", gw0["gateway_ip"])
	assert.Equal(t, "ap-guangzhou", gw0["region_id"])
	assert.Equal(t, "false", gw0["need_confirm"])

	lines := gw0["lines"].([]interface{})
	assert.Len(t, lines, 1)
	line0 := lines[0].(map[string]interface{})
	assert.Equal(t, "line-0", line0["line_id"])
	assert.Equal(t, "direct", line0["line_type"])
	assert.Equal(t, "1.2.3.4:8080", line0["line_address"])

	gw1 := gateways[1].(map[string]interface{})
	assert.Equal(t, "gateway-002", gw1["gateway_id"])
	assert.Equal(t, "private", gw1["gateway_type"])
	assert.Equal(t, "offline", gw1["status"])
	assert.Equal(t, "true", gw1["need_confirm"])
}

// TestTeoMultiPathGatewaysDataSource_APIError tests Read handles API error
func TestTeoMultiPathGatewaysDataSource_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeMultiPathGateways", func(request *teov20220901.DescribeMultiPathGatewaysRequest) (*teov20220901.DescribeMultiPathGatewaysResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid ZoneId")
	})

	meta := &mpgMockMeta{client: mockClient}
	res := teo.DataSourceTencentCloudTeoMultiPathGateways()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "invalid-zone-id",
	})

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestTeoMultiPathGatewaysDataSource_Schema validates schema definition
func TestTeoMultiPathGatewaysDataSource_Schema(t *testing.T) {
	res := teo.DataSourceTencentCloudTeoMultiPathGateways()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Read)

	assert.Contains(t, res.Schema, "zone_id")
	assert.Contains(t, res.Schema, "filters")
	assert.Contains(t, res.Schema, "gateways")
	assert.Contains(t, res.Schema, "result_output_file")

	zoneId := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneId.Type)
	assert.True(t, zoneId.Required)

	filters := res.Schema["filters"]
	assert.Equal(t, schema.TypeList, filters.Type)
	assert.True(t, filters.Optional)

	gateways := res.Schema["gateways"]
	assert.Equal(t, schema.TypeList, gateways.Type)
	assert.True(t, gateways.Computed)

	resultOutputFile := res.Schema["result_output_file"]
	assert.Equal(t, schema.TypeString, resultOutputFile.Type)
	assert.True(t, resultOutputFile.Optional)

	// Validate gateways elem schema
	gatewayElem := gateways.Elem.(*schema.Resource)
	assert.Contains(t, gatewayElem.Schema, "gateway_id")
	assert.Contains(t, gatewayElem.Schema, "gateway_name")
	assert.Contains(t, gatewayElem.Schema, "gateway_type")
	assert.Contains(t, gatewayElem.Schema, "gateway_port")
	assert.Contains(t, gatewayElem.Schema, "status")
	assert.Contains(t, gatewayElem.Schema, "gateway_ip")
	assert.Contains(t, gatewayElem.Schema, "region_id")
	assert.Contains(t, gatewayElem.Schema, "lines")
	assert.Contains(t, gatewayElem.Schema, "need_confirm")
}

// mpgMockMeta implements tccommon.ProviderMeta
type mpgMockMeta struct {
	client *connectivity.TencentCloudClient
}

func (m *mpgMockMeta) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mpgMockMeta{}

func ptrMPGString(s string) *string {
	return &s
}

func ptrMPGInt64(i int64) *int64 {
	return &i
}
