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

// go test ./tencentcloud/services/teo/ -run "TestTeoIPRegion" -v -count=1 -gcflags="all=-l"

// TestTeoIPRegionDataSource_Success tests Read calls API and sets ip_region_info
func TestTeoIPRegionDataSource_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	// Patch UseTeoClient to return a non-nil client
	teoClient := &teov20220901.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseTeoClient", teoClient)

	// Patch DescribeIPRegion to return success
	patches.ApplyMethodFunc(teoClient, "DescribeIPRegion", func(request *teov20220901.DescribeIPRegionRequest) (*teov20220901.DescribeIPRegionResponse, error) {
		resp := teov20220901.NewDescribeIPRegionResponse()
		resp.Response = &teov20220901.DescribeIPRegionResponseParams{
			IPRegionInfo: []*teov20220901.IPRegionInfo{
				{
					IP:          ptrIPRegionString("1.1.1.1"),
					IsEdgeOneIP: ptrIPRegionString("yes"),
				},
				{
					IP:          ptrIPRegionString("2.2.2.2"),
					IsEdgeOneIP: ptrIPRegionString("no"),
				},
			},
			RequestId: ptrIPRegionString("fake-request-id"),
		}
		return resp, nil
	})

	meta := &ipRegionMockMeta{client: mockClient}
	res := teo.DataSourceTencentCloudTeoIPRegion()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"ips": []interface{}{"1.1.1.1", "2.2.2.2"},
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	ipRegionInfo := d.Get("ip_region_info").([]interface{})
	assert.Len(t, ipRegionInfo, 2)

	info0 := ipRegionInfo[0].(map[string]interface{})
	assert.Equal(t, "1.1.1.1", info0["ip"])
	assert.Equal(t, "yes", info0["is_edge_one_ip"])

	info1 := ipRegionInfo[1].(map[string]interface{})
	assert.Equal(t, "2.2.2.2", info1["ip"])
	assert.Equal(t, "no", info1["is_edge_one_ip"])
}

// TestTeoIPRegionDataSource_APIError tests Read handles API error
func TestTeoIPRegionDataSource_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeIPRegion", func(request *teov20220901.DescribeIPRegionRequest) (*teov20220901.DescribeIPRegionResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InvalidParameter, Message=Invalid IP format")
	})

	meta := &ipRegionMockMeta{client: mockClient}
	res := teo.DataSourceTencentCloudTeoIPRegion()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"ips": []interface{}{"invalid-ip"},
	})

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InvalidParameter")
}

// TestTeoIPRegionDataSource_Schema validates schema definition
func TestTeoIPRegionDataSource_Schema(t *testing.T) {
	res := teo.DataSourceTencentCloudTeoIPRegion()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Read)

	assert.Contains(t, res.Schema, "ips")
	assert.Contains(t, res.Schema, "ip_region_info")
	assert.Contains(t, res.Schema, "result_output_file")

	iPs := res.Schema["ips"]
	assert.Equal(t, schema.TypeList, iPs.Type)
	assert.True(t, iPs.Required)
	assert.Equal(t, 100, iPs.MaxItems)

	ipRegionInfo := res.Schema["ip_region_info"]
	assert.Equal(t, schema.TypeList, ipRegionInfo.Type)
	assert.True(t, ipRegionInfo.Computed)

	resultOutputFile := res.Schema["result_output_file"]
	assert.Equal(t, schema.TypeString, resultOutputFile.Type)
	assert.True(t, resultOutputFile.Optional)
}

// ipRegionMockMeta implements tccommon.ProviderMeta
type ipRegionMockMeta struct {
	client *connectivity.TencentCloudClient
}

func (m *ipRegionMockMeta) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &ipRegionMockMeta{}

func ptrIPRegionString(s string) *string {
	return &s
}
