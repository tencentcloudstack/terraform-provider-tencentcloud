package dbdc_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	dbdcv20201029 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbdc/v20201029"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/dbdc"
)

// go test ./tencentcloud/services/dbdc/ -run "TestDbdcDbCustomRegionsDS" -v -count=1 -gcflags="all=-l"

func TestDbdcDbCustomRegionsDS_ReadBasic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dbdcClient := &dbdcv20201029.Client{}
	patches.ApplyMethodReturn(newMockMetaDbdcDS().client, "UseDbdcV20201029Client", dbdcClient)

	patches.ApplyMethodFunc(dbdcClient, "DescribeDBCustomRegions", func(request *dbdcv20201029.DescribeDBCustomRegionsRequest) (*dbdcv20201029.DescribeDBCustomRegionsResponse, error) {
		resp := dbdcv20201029.NewDescribeDBCustomRegionsResponse()
		resp.Response = &dbdcv20201029.DescribeDBCustomRegionsResponseParams{
			RegionSet: []*dbdcv20201029.RegionInfo{
				{
					Region:      ptrStr("ap-guangzhou"),
					RegionState: ptrStr("SELL"),
				},
				{
					Region:      ptrStr("ap-shanghai"),
					RegionState: ptrStr("SOLD_OUT"),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaDbdcDS()
	res := dbdc.DataSourceTencentCloudDbdcDbCustomRegions()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	regionSet := d.Get("region_set").([]interface{})
	assert.Len(t, regionSet, 2)

	region0 := regionSet[0].(map[string]interface{})
	assert.Equal(t, "ap-guangzhou", region0["region"].(string))
	assert.Equal(t, "SELL", region0["region_state"].(string))

	region1 := regionSet[1].(map[string]interface{})
	assert.Equal(t, "ap-shanghai", region1["region"].(string))
	assert.Equal(t, "SOLD_OUT", region1["region_state"].(string))
}

func TestDbdcDbCustomRegionsDS_ReadWithNilFields(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dbdcClient := &dbdcv20201029.Client{}
	patches.ApplyMethodReturn(newMockMetaDbdcDS().client, "UseDbdcV20201029Client", dbdcClient)

	patches.ApplyMethodFunc(dbdcClient, "DescribeDBCustomRegions", func(request *dbdcv20201029.DescribeDBCustomRegionsRequest) (*dbdcv20201029.DescribeDBCustomRegionsResponse, error) {
		resp := dbdcv20201029.NewDescribeDBCustomRegionsResponse()
		resp.Response = &dbdcv20201029.DescribeDBCustomRegionsResponseParams{
			RegionSet: []*dbdcv20201029.RegionInfo{
				{
					Region: ptrStr("ap-beijing"),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaDbdcDS()
	res := dbdc.DataSourceTencentCloudDbdcDbCustomRegions()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	err := res.Read(d, meta)
	assert.NoError(t, err)

	regionSet := d.Get("region_set").([]interface{})
	assert.Len(t, regionSet, 1)

	region0 := regionSet[0].(map[string]interface{})
	assert.Equal(t, "ap-beijing", region0["region"].(string))
	// RegionState is nil in the API response, Terraform SDK defaults it to empty string
	assert.Equal(t, "", region0["region_state"].(string))
}

func TestDbdcDbCustomRegionsDS_ReadWithEmptyResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dbdcClient := &dbdcv20201029.Client{}
	patches.ApplyMethodReturn(newMockMetaDbdcDS().client, "UseDbdcV20201029Client", dbdcClient)

	patches.ApplyMethodFunc(dbdcClient, "DescribeDBCustomRegions", func(request *dbdcv20201029.DescribeDBCustomRegionsRequest) (*dbdcv20201029.DescribeDBCustomRegionsResponse, error) {
		resp := dbdcv20201029.NewDescribeDBCustomRegionsResponse()
		resp.Response = &dbdcv20201029.DescribeDBCustomRegionsResponseParams{
			RegionSet: nil,
		}
		return resp, nil
	})

	meta := newMockMetaDbdcDS()
	res := dbdc.DataSourceTencentCloudDbdcDbCustomRegions()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	err := res.Read(d, meta)
	// When response RegionSet is nil, the datasource should return an error (NonRetryableError)
	assert.Error(t, err)
}

func TestDbdcDbCustomRegionsDS_Schema(t *testing.T) {
	res := dbdc.DataSourceTencentCloudDbdcDbCustomRegions()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "region_set")
	assert.Contains(t, res.Schema, "result_output_file")

	regionSetSchema := res.Schema["region_set"]
	assert.Equal(t, schema.TypeList, regionSetSchema.Type)
	assert.True(t, regionSetSchema.Computed)

	elemRes := regionSetSchema.Elem.(*schema.Resource)
	assert.Contains(t, elemRes.Schema, "region")
	assert.Contains(t, elemRes.Schema, "region_state")

	regionSchema := elemRes.Schema["region"]
	assert.Equal(t, schema.TypeString, regionSchema.Type)
	assert.True(t, regionSchema.Computed)

	regionStateSchema := elemRes.Schema["region_state"]
	assert.Equal(t, schema.TypeString, regionStateSchema.Type)
	assert.True(t, regionStateSchema.Computed)
}
