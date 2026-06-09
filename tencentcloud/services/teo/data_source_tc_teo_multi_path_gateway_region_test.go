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

// go test ./tencentcloud/services/teo/ -run "TestTeoMultiPathGatewayRegion" -v -count=1 -gcflags="all=-l"

// TestTeoMultiPathGatewayRegionDataSource_Success tests Read calls API and sets gateway_regions
func TestTeoMultiPathGatewayRegionDataSource_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	// Patch UseTeoClient to return a non-nil client
	teoClient := &teov20220901.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseTeoClient", teoClient)

	// Patch DescribeMultiPathGatewayRegions to return success
	patches.ApplyMethodFunc(teoClient, "DescribeMultiPathGatewayRegions", func(request *teov20220901.DescribeMultiPathGatewayRegionsRequest) (*teov20220901.DescribeMultiPathGatewayRegionsResponse, error) {
		resp := teov20220901.NewDescribeMultiPathGatewayRegionsResponse()
		resp.Response = &teov20220901.DescribeMultiPathGatewayRegionsResponseParams{
			GatewayRegions: []*teov20220901.GatewayRegion{
				{
					RegionId: ptrMPGRString("ap-guangzhou"),
					CNName:   ptrMPGRString("广州"),
					ENName:   ptrMPGRString("Guangzhou"),
				},
				{
					RegionId: ptrMPGRString("ap-shanghai"),
					CNName:   ptrMPGRString("上海"),
					ENName:   ptrMPGRString("Shanghai"),
				},
			},
			RequestId: ptrMPGRString("fake-request-id"),
		}
		return resp, nil
	})

	meta := &mpgrMockMeta{client: mockClient}
	res := teo.DataSourceTencentCloudTeoMultiPathGatewayRegion()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "zone-2noqxz9b6klw",
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	gatewayRegions := d.Get("gateway_regions").([]interface{})
	assert.Len(t, gatewayRegions, 2)

	region0 := gatewayRegions[0].(map[string]interface{})
	assert.Equal(t, "ap-guangzhou", region0["region_id"])
	assert.Equal(t, "广州", region0["cn_name"])
	assert.Equal(t, "Guangzhou", region0["en_name"])

	region1 := gatewayRegions[1].(map[string]interface{})
	assert.Equal(t, "ap-shanghai", region1["region_id"])
	assert.Equal(t, "上海", region1["cn_name"])
	assert.Equal(t, "Shanghai", region1["en_name"])
}

// TestTeoMultiPathGatewayRegionDataSource_APIError tests Read handles API error
func TestTeoMultiPathGatewayRegionDataSource_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeMultiPathGatewayRegions", func(request *teov20220901.DescribeMultiPathGatewayRegionsRequest) (*teov20220901.DescribeMultiPathGatewayRegionsResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid ZoneId")
	})

	meta := &mpgrMockMeta{client: mockClient}
	res := teo.DataSourceTencentCloudTeoMultiPathGatewayRegion()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id": "invalid-zone-id",
	})

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestTeoMultiPathGatewayRegionDataSource_Schema validates schema definition
func TestTeoMultiPathGatewayRegionDataSource_Schema(t *testing.T) {
	res := teo.DataSourceTencentCloudTeoMultiPathGatewayRegion()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Read)

	assert.Contains(t, res.Schema, "zone_id")
	assert.Contains(t, res.Schema, "gateway_regions")
	assert.Contains(t, res.Schema, "result_output_file")

	zoneId := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneId.Type)
	assert.True(t, zoneId.Required)

	gatewayRegions := res.Schema["gateway_regions"]
	assert.Equal(t, schema.TypeList, gatewayRegions.Type)
	assert.True(t, gatewayRegions.Computed)

	resultOutputFile := res.Schema["result_output_file"]
	assert.Equal(t, schema.TypeString, resultOutputFile.Type)
	assert.True(t, resultOutputFile.Optional)
}

// mpgrMockMeta implements tccommon.ProviderMeta
type mpgrMockMeta struct {
	client *connectivity.TencentCloudClient
}

func (m *mpgrMockMeta) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &mpgrMockMeta{}

func ptrMPGRString(s string) *string {
	return &s
}
