package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDbbrainDescribeSlowLogsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDbbrainDescribeSlowLogsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_dbbrain_describe_slow_logs.describe_slow_logs")),
			},
		},
	})
}

const testAccDbbrainDescribeSlowLogsDataSource = `

data "tencentcloud_dbbrain_describe_slow_logs" "describe_slow_logs" {
  product = ""
  instance_id = ""
  md5 = ""
  start_time = ""
  end_time = ""
  d_b = 
  key = 
  user = 
  ip = 
  time = 
  }

`
