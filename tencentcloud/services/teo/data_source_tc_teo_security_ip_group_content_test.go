package teo_test

import (
	"testing"

	"github.com/agiledragon/gomonkey/v2"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/stretchr/testify/assert"
	teov20220901 "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/teo/v20220901"

	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/connectivity"
	"github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/services/teo"
)

// go test ./tencentcloud/services/teo/ -run "TestTeoSecurityIPGroupContentDataSource" -v -count=1 -gcflags="all=-l"

// TestTeoSecurityIPGroupContentDataSource_ReadSuccess tests successful read with IP group content data
func TestTeoSecurityIPGroupContentDataSource_ReadSuccess(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(&connectivity.TencentCloudClient{}, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeSecurityIPGroupContent", func(request *teov20220901.DescribeSecurityIPGroupContentRequest) (*teov20220901.DescribeSecurityIPGroupContentResponse, error) {
		resp := teov20220901.NewDescribeSecurityIPGroupContentResponse()
		resp.Response = &teov20220901.DescribeSecurityIPGroupContentResponseParams{
			IPTotalCount: ptrInt64(3),
			IPList: []*string{
				ptrString("1.2.3.4"),
				ptrString("5.6.7.0/24"),
				ptrString("10.0.0.1"),
			},
			RequestId: ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.DataSourceTencentCloudTeoSecurityIPGroupContent()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":  "zone-2qtuhspy7cr6",
		"group_id": 123,
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())

	ipTotalCount := d.Get("ip_total_count").(int)
	assert.Equal(t, 3, ipTotalCount)

	ipList := d.Get("ip_list").([]interface{})
	assert.Len(t, ipList, 3)
	assert.Equal(t, "1.2.3.4", ipList[0].(string))
	assert.Equal(t, "5.6.7.0/24", ipList[1].(string))
	assert.Equal(t, "10.0.0.1", ipList[2].(string))
}

// TestTeoSecurityIPGroupContentDataSource_ReadEmpty tests read with no IP data
func TestTeoSecurityIPGroupContentDataSource_ReadEmpty(t *testing.T) {
	patches := gomonkey.NewPatches()
	defer patches.Reset()

	teoClient := &teov20220901.Client{}
	patches.ApplyMethodReturn(&connectivity.TencentCloudClient{}, "UseTeoClient", teoClient)

	patches.ApplyMethodFunc(teoClient, "DescribeSecurityIPGroupContent", func(request *teov20220901.DescribeSecurityIPGroupContentRequest) (*teov20220901.DescribeSecurityIPGroupContentResponse, error) {
		resp := teov20220901.NewDescribeSecurityIPGroupContentResponse()
		resp.Response = &teov20220901.DescribeSecurityIPGroupContentResponseParams{
			IPTotalCount: ptrInt64(0),
			IPList:       []*string{},
			RequestId:    ptrString("fake-request-id"),
		}
		return resp, nil
	})

	meta := newMockMeta()
	res := teo.DataSourceTencentCloudTeoSecurityIPGroupContent()
	d := schema.TestResourceDataRaw(t, res.Schema, map[string]interface{}{
		"zone_id":  "zone-empty",
		"group_id": 456,
	})

	err := res.Read(d, meta)
	assert.NoError(t, err)
	assert.NotEmpty(t, d.Id())
}

// TestTeoSecurityIPGroupContentDataSource_Schema validates schema definition
func TestTeoSecurityIPGroupContentDataSource_Schema(t *testing.T) {
	res := teo.DataSourceTencentCloudTeoSecurityIPGroupContent()

	assert.NotNil(t, res)
	assert.NotNil(t, res.Read)

	assert.Contains(t, res.Schema, "zone_id")
	assert.Contains(t, res.Schema, "group_id")
	assert.Contains(t, res.Schema, "ip_total_count")
	assert.Contains(t, res.Schema, "ip_list")
	assert.Contains(t, res.Schema, "result_output_file")

	// zone_id is Required TypeString
	zoneId := res.Schema["zone_id"]
	assert.Equal(t, schema.TypeString, zoneId.Type)
	assert.True(t, zoneId.Required)

	// group_id is Required TypeInt
	groupId := res.Schema["group_id"]
	assert.Equal(t, schema.TypeInt, groupId.Type)
	assert.True(t, groupId.Required)

	// ip_total_count is Computed TypeInt
	ipTotalCount := res.Schema["ip_total_count"]
	assert.Equal(t, schema.TypeInt, ipTotalCount.Type)
	assert.True(t, ipTotalCount.Computed)

	// ip_list is Computed TypeList of TypeString
	ipList := res.Schema["ip_list"]
	assert.Equal(t, schema.TypeList, ipList.Type)
	assert.True(t, ipList.Computed)

	// result_output_file is Optional TypeString
	outputFile := res.Schema["result_output_file"]
	assert.Equal(t, schema.TypeString, outputFile.Type)
	assert.True(t, outputFile.Optional)
}
