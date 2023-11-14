package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudDcdbSlowLogsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccDcdbSlowLogsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_dcdb_slow_logs.slow_logs")),
			},
		},
	})
}

const testAccDcdbSlowLogsDataSource = `

data "tencentcloud_dcdb_slow_logs" "slow_logs" {
  instance_id = ""
  start_time = ""
  shard_id = ""
  end_time = ""
  db = ""
  order_by = ""
  order_by_type = ""
  slave = 
        }

`
