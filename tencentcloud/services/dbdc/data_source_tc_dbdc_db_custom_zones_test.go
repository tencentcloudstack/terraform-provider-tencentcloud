package dbdc_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	dbdcv20201029 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/dbdc/v20201029"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/dbdc"
)

// go test ./tencentcloud/services/dbdc/ -run "TestDbdcDbCustomZonesDS" -v -count=1 -gcflags="all=-l"

func TestDbdcDbCustomZonesDS_ReadBasic(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dbdcClient := &dbdcv20201029.Client{}
	patches.ApplyMethodReturn(newMockMetaDbdcDS().client, "UseDbdcV20201029Client", dbdcClient)

	patches.ApplyMethodFunc(dbdcClient, "DescribeDBCustomZones", func(request *dbdcv20201029.DescribeDBCustomZonesRequest) (*dbdcv20201029.DescribeDBCustomZonesResponse, error) {
		resp := dbdcv20201029.NewDescribeDBCustomZonesResponse()
		resp.Response = &dbdcv20201029.DescribeDBCustomZonesResponseParams{
			ZoneSet: []*dbdcv20201029.ZoneInfo{
				{
					Zone:      ptrStr("ap-guangzhou-3"),
					ZoneState: ptrStr("SELL"),
				},
				{
					Zone:      ptrStr("ap-guangzhou-4"),
					ZoneState: ptrStr("SOLD_OUT"),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaDbdcDS()
	res := dbdc.DataSourceTencentCloudDbdcDbCustomZones()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	zoneSet := d.Get("zone_set").([]interface{})
	assert.Len(t, zoneSet, 2)

	zone0 := zoneSet[0].(map[string]interface{})
	assert.Equal(t, "ap-guangzhou-3", zone0["zone"].(string))
	assert.Equal(t, "SELL", zone0["zone_state"].(string))

	zone1 := zoneSet[1].(map[string]interface{})
	assert.Equal(t, "ap-guangzhou-4", zone1["zone"].(string))
	assert.Equal(t, "SOLD_OUT", zone1["zone_state"].(string))
}

func TestDbdcDbCustomZonesDS_ReadWithNilFields(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dbdcClient := &dbdcv20201029.Client{}
	patches.ApplyMethodReturn(newMockMetaDbdcDS().client, "UseDbdcV20201029Client", dbdcClient)

	patches.ApplyMethodFunc(dbdcClient, "DescribeDBCustomZones", func(request *dbdcv20201029.DescribeDBCustomZonesRequest) (*dbdcv20201029.DescribeDBCustomZonesResponse, error) {
		resp := dbdcv20201029.NewDescribeDBCustomZonesResponse()
		resp.Response = &dbdcv20201029.DescribeDBCustomZonesResponseParams{
			ZoneSet: []*dbdcv20201029.ZoneInfo{
				{
					Zone: ptrStr("ap-shanghai-2"),
				},
			},
		}
		return resp, nil
	})

	meta := newMockMetaDbdcDS()
	res := dbdc.DataSourceTencentCloudDbdcDbCustomZones()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	err := res.Read(d, meta)
	assert.NoError(t, err)

	zoneSet := d.Get("zone_set").([]interface{})
	assert.Len(t, zoneSet, 1)

	zone0 := zoneSet[0].(map[string]interface{})
	assert.Equal(t, "ap-shanghai-2", zone0["zone"].(string))
	// ZoneState is nil in the API response, Terraform SDK defaults it to empty string
	assert.Equal(t, "", zone0["zone_state"].(string))
}

func TestDbdcDbCustomZonesDS_ReadWithEmptyResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dbdcClient := &dbdcv20201029.Client{}
	patches.ApplyMethodReturn(newMockMetaDbdcDS().client, "UseDbdcV20201029Client", dbdcClient)

	patches.ApplyMethodFunc(dbdcClient, "DescribeDBCustomZones", func(request *dbdcv20201029.DescribeDBCustomZonesRequest) (*dbdcv20201029.DescribeDBCustomZonesResponse, error) {
		resp := dbdcv20201029.NewDescribeDBCustomZonesResponse()
		resp.Response = &dbdcv20201029.DescribeDBCustomZonesResponseParams{
			ZoneSet: []*dbdcv20201029.ZoneInfo{},
		}
		return resp, nil
	})

	meta := newMockMetaDbdcDS()
	res := dbdc.DataSourceTencentCloudDbdcDbCustomZones()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	err := res.Read(d, meta)
	// When response is empty, the datasource should return an error (NonRetryableError)
	assert.Error(t, err)
}

func TestDbdcDbCustomZonesDS_ReadWithNilResponse(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	dbdcClient := &dbdcv20201029.Client{}
	patches.ApplyMethodReturn(newMockMetaDbdcDS().client, "UseDbdcV20201029Client", dbdcClient)

	patches.ApplyMethodFunc(dbdcClient, "DescribeDBCustomZones", func(request *dbdcv20201029.DescribeDBCustomZonesRequest) (*dbdcv20201029.DescribeDBCustomZonesResponse, error) {
		resp := dbdcv20201029.NewDescribeDBCustomZonesResponse()
		resp.Response = nil
		return resp, nil
	})

	meta := newMockMetaDbdcDS()
	res := dbdc.DataSourceTencentCloudDbdcDbCustomZones()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{})

	err := res.Read(d, meta)
	// When Response is nil, the datasource should return an error (NonRetryableError)
	assert.Error(t, err)
}

func TestDbdcDbCustomZonesDS_Schema(t *testing.T) {
	res := dbdc.DataSourceTencentCloudDbdcDbCustomZones()

	assert.NotNil(t, res)
	assert.Contains(t, res.Schema, "zone_set")
	assert.Contains(t, res.Schema, "result_output_file")

	zoneSetSchema := res.Schema["zone_set"]
	assert.Equal(t, schema.TypeList, zoneSetSchema.Type)
	assert.True(t, zoneSetSchema.Computed)

	elemRes := zoneSetSchema.Elem.(*schema.Resource)
	assert.Contains(t, elemRes.Schema, "zone")
	assert.Contains(t, elemRes.Schema, "zone_state")

	zoneSchema := elemRes.Schema["zone"]
	assert.Equal(t, schema.TypeString, zoneSchema.Type)
	assert.True(t, zoneSchema.Computed)

	zoneStateSchema := elemRes.Schema["zone_state"]
	assert.Equal(t, schema.TypeString, zoneStateSchema.Type)
	assert.True(t, zoneStateSchema.Computed)

	resultOutputFileSchema := res.Schema["result_output_file"]
	assert.Equal(t, schema.TypeString, resultOutputFileSchema.Type)
	assert.True(t, resultOutputFileSchema.Optional)
}
