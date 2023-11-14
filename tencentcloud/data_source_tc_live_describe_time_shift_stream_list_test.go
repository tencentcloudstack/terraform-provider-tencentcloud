package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudLiveDescribeTimeShiftStreamListDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccLiveDescribeTimeShiftStreamListDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_live_describe_time_shift_stream_list.describe_time_shift_stream_list")),
			},
		},
	})
}

const testAccLiveDescribeTimeShiftStreamListDataSource = `

data "tencentcloud_live_describe_time_shift_stream_list" "describe_time_shift_stream_list" {
  start_time = 
  end_time = 
  stream_name = ""
  domain = ""
  domain_group = ""
    }

`
