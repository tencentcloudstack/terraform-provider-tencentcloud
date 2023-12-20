package dbbrain_test

import (
	"fmt"
	"testing"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDbbrainSlowLogsDataSource_basic(t *testing.T) {
	t.Parallel()
	startTime := "2023-05-15 16:32:00"
	endTime := "2023-05-15 17:03:00"
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDbbrainSlowLogsDataSource, tcacctest.DefaultDbBrainInstanceId, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dbbrain_slow_logs.slow_logs"),
					resource.TestCheckResourceAttr("data.tencentcloud_dbbrain_slow_logs.slow_logs", "instance_id", tcacctest.DefaultDbBrainInstanceId),
					resource.TestCheckResourceAttr("data.tencentcloud_dbbrain_slow_logs.slow_logs", "md5", "4961208426639258265"),
					resource.TestCheckResourceAttr("data.tencentcloud_dbbrain_slow_logs.slow_logs", "product", "mysql"),
					resource.TestCheckResourceAttr("data.tencentcloud_dbbrain_slow_logs.slow_logs", "start_time", startTime),
					resource.TestCheckResourceAttr("data.tencentcloud_dbbrain_slow_logs.slow_logs", "end_time", endTime),
					// resource.TestCheckResourceAttr("data.tencentcloud_dbbrain_slow_logs.slow_logs", "db", "mysql"),
					// resource.TestCheckResourceAttr("data.tencentcloud_dbbrain_slow_logs.slow_logs", "key", "mysql"),
					// resource.TestCheckResourceAttr("data.tencentcloud_dbbrain_slow_logs.slow_logs", "user", "mysql"),
					// resource.TestCheckResourceAttr("data.tencentcloud_dbbrain_slow_logs.slow_logs", "ip", "mysql"),
					// resource.TestCheckResourceAttr("data.tencentcloud_dbbrain_slow_logs.slow_logs", "time", "mysql"),

					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_slow_logs.slow_logs", "rows.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_slow_logs.slow_logs", "rows.0.timestamp"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_slow_logs.slow_logs", "rows.0.sql_text"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_slow_logs.slow_logs", "rows.0.database"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_slow_logs.slow_logs", "rows.0.user_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_slow_logs.slow_logs", "rows.0.user_host"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_slow_logs.slow_logs", "rows.0.query_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_slow_logs.slow_logs", "rows.0.lock_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_slow_logs.slow_logs", "rows.0.rows_examined"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_slow_logs.slow_logs", "rows.0.rows_sent"),
				),
			},
		},
	})
}

const testAccDbbrainSlowLogsDataSource = `

data "tencentcloud_dbbrain_slow_logs" "slow_logs" {
  product = "mysql"
  instance_id = "%s"
  md5 = "4961208426639258265"
  start_time = "%s"
  end_time = "%s"
}

`
