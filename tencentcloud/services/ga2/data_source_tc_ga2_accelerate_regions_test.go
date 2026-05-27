package ga2_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	ga2v20250115 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/ga2/v20250115"

	tccommon "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/common"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/ga2"
)

// go test ./tencentcloud/services/ga2/ -run "TestGa2AccelerateRegions" -v -count=1 -gcflags="all=-l"

// TestGa2AccelerateRegionsDataSource_Success tests Read calls API and sets accelerator_region_set
func TestGa2AccelerateRegionsDataSource_Success(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseGa2V20250115Client", ga2Client)

	patches.ApplyMethodFunc(ga2Client, "DescribeAccelerateRegionsWithContext", func(ctx context.Context, request *ga2v20250115.DescribeAccelerateRegionsRequest) (*ga2v20250115.DescribeAccelerateRegionsResponse, error) {
		resp := ga2v20250115.NewDescribeAccelerateRegionsResponse()
		resp.Response = &ga2v20250115.DescribeAccelerateRegionsResponseParams{
			AcceleratorRegionSet: []*ga2v20250115.AcceleratorRegionSet{
				{
					Name:            ptrGa2String("广州"),
					IsAvailable:     ptrGa2Int64(1),
					Region:          ptrGa2String("ap-guangzhou"),
					AreaName:        ptrGa2String("华南地区"),
					IsChinaMainland: ptrGa2Uint64(1),
					SupportIspType:  []*string{ptrGa2String("BGP"), ptrGa2String("ChinaTelecom")},
					IsTencentRegion: ptrGa2Uint64(1),
				},
				{
					Name:            ptrGa2String("硅谷"),
					IsAvailable:     ptrGa2Int64(1),
					Region:          ptrGa2String("na-siliconvalley"),
					AreaName:        ptrGa2String("北美地区"),
					IsChinaMainland: ptrGa2Uint64(0),
					SupportIspType:  []*string{ptrGa2String("BGP")},
					IsTencentRegion: ptrGa2Uint64(1),
				},
			},
		}
		return resp, nil
	})

	meta := &ga2MockMeta{client: mockClient}
	res := ga2.DataSourceTencentCloudGa2AccelerateRegions()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	regionSet := d.Get("accelerator_region_set").([]interface{})
	assert.Len(t, regionSet, 2)

	region0 := regionSet[0].(map[string]interface{})
	assert.Equal(t, "广州", region0["name"])
	assert.Equal(t, 1, region0["is_available"])
	assert.Equal(t, "ap-guangzhou", region0["region"])
	assert.Equal(t, "华南地区", region0["area_name"])
	assert.Equal(t, 1, region0["is_china_mainland"])
	assert.Equal(t, 1, region0["is_tencent_region"])
	ispTypes0 := region0["support_isp_type"].([]interface{})
	assert.Equal(t, "BGP", ispTypes0[0].(string))
	assert.Equal(t, "ChinaTelecom", ispTypes0[1].(string))

	region1 := regionSet[1].(map[string]interface{})
	assert.Equal(t, "硅谷", region1["name"])
	assert.Equal(t, "na-siliconvalley", region1["region"])
	assert.Equal(t, "北美地区", region1["area_name"])
	assert.Equal(t, 0, region1["is_china_mainland"])
}

// TestGa2AccelerateRegionsDataSource_APIError tests Read handles API error
func TestGa2AccelerateRegionsDataSource_APIError(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseGa2V20250115Client", ga2Client)

	patches.ApplyMethodFunc(ga2Client, "DescribeAccelerateRegionsWithContext", func(ctx context.Context, request *ga2v20250115.DescribeAccelerateRegionsRequest) (*ga2v20250115.DescribeAccelerateRegionsResponse, error) {
		return nil, fmt.Errorf("[TencentCloudSDKError] Code=InternalError, Message=Internal error")
	})

	meta := &ga2MockMeta{client: mockClient}
	res := ga2.DataSourceTencentCloudGa2AccelerateRegions()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	err := res.Read(d, meta)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "InternalError")
}

// TestGa2AccelerateRegionsDataSource_EmptyResult tests Read handles empty result
func TestGa2AccelerateRegionsDataSource_EmptyResult(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	ga2Client := &ga2v20250115.Client{}
	mockClient := &connectivity.TencentCloudClient{}
	patches.ApplyMethodReturn(mockClient, "UseGa2V20250115Client", ga2Client)

	patches.ApplyMethodFunc(ga2Client, "DescribeAccelerateRegionsWithContext", func(ctx context.Context, request *ga2v20250115.DescribeAccelerateRegionsRequest) (*ga2v20250115.DescribeAccelerateRegionsResponse, error) {
		resp := ga2v20250115.NewDescribeAccelerateRegionsResponse()
		resp.Response = &ga2v20250115.DescribeAccelerateRegionsResponseParams{
			AcceleratorRegionSet: []*ga2v20250115.AcceleratorRegionSet{},
		}
		return resp, nil
	})

	meta := &ga2MockMeta{client: mockClient}
	res := ga2.DataSourceTencentCloudGa2AccelerateRegions()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	regionSet := d.Get("accelerator_region_set").([]interface{})
	assert.Len(t, regionSet, 0)
}

// TestGa2AccelerateRegionsDataSource_Schema validates schema definition
func TestGa2AccelerateRegionsDataSource_Schema(t *testing.T) {
	res := ga2.DataSourceTencentCloudGa2AccelerateRegions()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Read)

	assert.Contains(t, res.Schema, "accelerator_region_set")
	assert.Contains(t, res.Schema, "result_output_file")

	regionSet := res.Schema["accelerator_region_set"]
	assert.Equal(t, schema.TypeList, regionSet.Type)
	assert.True(t, regionSet.Computed)

	resultOutputFile := res.Schema["result_output_file"]
	assert.Equal(t, schema.TypeString, resultOutputFile.Type)
	assert.True(t, resultOutputFile.Optional)
}

// ga2MockMeta implements tccommon.ProviderMeta
type ga2MockMeta struct {
	client *connectivity.TencentCloudClient
}

func (m *ga2MockMeta) GetAPIV3Conn() *connectivity.TencentCloudClient {
	return m.client
}

var _ tccommon.ProviderMeta = &ga2MockMeta{}

func ptrGa2String(s string) *string {
	return &s
}

func ptrGa2Int64(i int64) *int64 {
	return &i
}

func ptrGa2Uint64(i uint64) *uint64 {
	return &i
}
