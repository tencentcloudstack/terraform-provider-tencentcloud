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

// go test ./tencentcloud/services/teo/ -run "TestTeoMultiPathGateway" -v -count=1 -gcflags="all=-l"

// TestTeoMultiPathGatewayDataSource_Success tests Read calls API and sets gateways
func TestTeoMultiPathGatewayDataSource_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	// Patch UseTeoClient to return a non-nil client
	teoClient := &teov20220901.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseTeoClient", teoClient)

	// Patch DescribeMultiPathGateways to return success
	patches.ApplyMethodFunc(teoClient, "DescribeMultiPathGateways", func(request *teov20220901.DescribeMultiPathGatewaysRequest) (*teov20220901.DescribeMultiPathGatewaysResponse, error) {
		resp := teov20220901.NewDescribeMultiPathGatewaysResponse()
		resp.Response = &teov20220901.DescribeMultiPathGatewaysResponseParams{
			Gateways: []*teov20220901.MultiPathGateway{
				{
					GatewayId:   ptrMPGString("gw-abc123"),
					GatewayName: ptrMPGString("test-gateway"),
					GatewayType: ptrMPGString("cloud"),
					GatewayPort: ptrMPGInt64(8080),
					Status:      ptrMPGString("online"),
					GatewayIP:   ptrMPGString("1.2.3.4"),
					RegionId:    ptrMPGString("ap-guangzhou"),
					NeedConfirm: ptrMPGString("false"),
				},
			},
			TotalCount: ptrMPGInt64(1),
			RequestId:  ptrMPGString("fake-request-id"),
		}
		return resp, nil
	})

	meta := &mpgMockMeta{client: mockClient}
	res := teo.DataSourceTencentCloudTeoMultiPathGateway()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-2noq7st5t3t6",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	gateways := d.Get("gateways").([]interface{})
	assert.Len(t, gateways, 1)

	gw0 := gateways[0].(map[string]interface{})
	assert.Equal(t, "gw-abc123", gw0["gateway_id"])
	assert.Equal(t, "test-gateway", gw0["gateway_name"])
	assert.Equal(t, "cloud", gw0["gateway_type"])
	assert.Equal(t, 8080, gw0["gateway_port"])
	assert.Equal(t, "online", gw0["status"])
	assert.Equal(t, "1.2.3.4", gw0["gateway_ip"])
	assert.Equal(t, "ap-guangzhou", gw0["region_id"])
	assert.Equal(t, "false", gw0["need_confirm"])
}

// TestTeoMultiPathGatewayDataSource_APIError tests Read handles API error
func TestTeoMultiPathGatewayDataSource_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeMultiPathGateways", func(request *teov20220901.DescribeMultiPathGatewaysRequest) (*teov20220901.DescribeMultiPathGatewaysResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid ZoneId")
	})

	meta := &mpgMockMeta{client: mockClient}
	res := teo.DataSourceTencentCloudTeoMultiPathGateway()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "invalid-zone-id",
	})

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestTeoMultiPathGatewayDataSource_Schema validates schema definition
func TestTeoMultiPathGatewayDataSource_Schema(t *testing.T) {
	res := teo.DataSourceTencentCloudTeoMultiPathGateway()

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

func ptrMPGInt64(n int64) *int64 {
	return &n
}
