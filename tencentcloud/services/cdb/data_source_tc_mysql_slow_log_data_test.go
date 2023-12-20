package cdb_test

import (
	"fmt"
	"testing"
	"time"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

// go test -i; go test -test.run TestAccTencentCloudMysqlSlowLogDataDataSource_basic -v
func TestAccTencentCloudNeedFixMysqlSlowLogDataDataSource_basic(t *testing.T) {
	t.Parallel()

	startTime := time.Now().AddDate(0, 0, -29).Unix()
	endTime := time.Now().Unix()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccMysqlSlowLogDataDataSource, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_mysql_slow_log_data.slow_log_data"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_slow_log_data.slow_log_data", "id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_slow_log_data.slow_log_data", "items.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_slow_log_data.slow_log_data", "items.0.database"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_slow_log_data.slow_log_data", "items.0.md5"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_slow_log_data.slow_log_data", "items.0.query_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_slow_log_data.slow_log_data", "items.0.rows_sent"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_slow_log_data.slow_log_data", "items.0.sql_template"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_slow_log_data.slow_log_data", "items.0.sql_text"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_slow_log_data.slow_log_data", "items.0.timestamp"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_slow_log_data.slow_log_data", "items.0.user_host"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_mysql_slow_log_data.slow_log_data", "items.0.user_name"),
				),
			},
		},
	})
}

const testAccMysqlSlowLogDataDataSourceVar = `
variable "instance_id" {
	default = "` + tcacctest.DefaultDbBrainInstanceId + `"
}
`

const testAccMysqlSlowLogDataDataSource = testAccMysqlSlowLogDataDataSourceVar + `

data "tencentcloud_mysql_slow_log_data" "slow_log_data" {
	instance_id = var.instance_id
	start_time = %v
	end_time = %v
	sort_by = "Timestamp"
	order_by = "ASC"
}

`
