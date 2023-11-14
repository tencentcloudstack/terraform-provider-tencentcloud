package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudMariadbSlowLogsDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccMariadbSlowLogsDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_mariadb_slow_logs.slow_logs")),
			},
		},
	})
}

const testAccMariadbSlowLogsDataSource = `

data "tencentcloud_mariadb_slow_logs" "slow_logs" {
  instance_id = ""
  start_time = ""
  end_time = ""
  db = ""
  order_by = ""
  order_by_type = ""
  slave = 
        }

`
