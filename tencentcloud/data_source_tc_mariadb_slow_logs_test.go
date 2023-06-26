package tencentcloud

import (
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMariadbSlowLogsDataSource_basic -v
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
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_mariadb_slow_logs.slow_logs"),
				),
			},
		},
	})
}

const testAccMariadbSlowLogsDataSource = `
data "tencentcloud_mariadb_slow_logs" "slow_logs" {
  instance_id   = "tdsql-9vqvls95"
  start_time    = "2023-06-01 14:55:20"
  order_by      = "query_time_sum"
  order_by_type = "desc"
  slave         = 0
}
`
