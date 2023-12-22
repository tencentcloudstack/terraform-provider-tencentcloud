package dcdb_test

import (
	"fmt"
	"testing"
	"time"

	tcacctest "github.com/tencentcloudstack/terraform-provider-tencentcloud/tencentcloud/acctest"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDcdbSlowLogsDataSource_basic(t *testing.T) {
	t.Parallel()
	loc, _ := time.LoadLocation("Asia/Chongqing")
	startTime := time.Now().AddDate(0, -3, 0).In(loc).Format("2006-01-02 15:04:05")
	endTime := time.Now().In(loc).Format("2006-01-02 15:04:05")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			tcacctest.AccPreCheck(t)
		},
		Providers: tcacctest.AccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDcdbSlowLogsDataSource, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					tcacctest.AccCheckTencentCloudDataSourceID("data.tencentcloud_dcdb_slow_logs.slow_logs"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_slow_logs.slow_logs", "instance_id"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_slow_logs.slow_logs", "start_time", startTime),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_slow_logs.slow_logs", "end_time", endTime),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_slow_logs.slow_logs", "shard_id", "shard-1b5r04az"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_slow_logs.slow_logs", "db", "tf_test_db"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_slow_logs.slow_logs", "order_by", "query_time_sum"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_slow_logs.slow_logs", "order_by_type", "desc"),
					resource.TestCheckResourceAttr("data.tencentcloud_dcdb_slow_logs.slow_logs", "slave", "0"),

					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_slow_logs.slow_logs", "data.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_slow_logs.slow_logs", "data.0.check_sum"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_slow_logs.slow_logs", "data.0.db"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_slow_logs.slow_logs", "data.0.finger_print"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_slow_logs.slow_logs", "data.0.lock_time_avg"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_slow_logs.slow_logs", "data.0.lock_time_sum"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_slow_logs.slow_logs", "data.0.query_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_slow_logs.slow_logs", "data.0.query_time_avg"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_slow_logs.slow_logs", "data.0.query_time_sum"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_slow_logs.slow_logs", "data.0.rows_sent_sum"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_slow_logs.slow_logs", "data.0.ts_max"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_slow_logs.slow_logs", "data.0.ts_min"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_slow_logs.slow_logs", "data.0.user"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_slow_logs.slow_logs", "data.0.example_sql"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_slow_logs.slow_logs", "data.0.host"),

					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_slow_logs.slow_logs", "lock_time_sum"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_slow_logs.slow_logs", "query_count"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_dcdb_slow_logs.slow_logs", "query_time_sum"),
				),
			},
		},
	})
}

const testAccDcdbSlowLogsDataSource = tcacctest.CommonPresetDcdb + `

data "tencentcloud_dcdb_slow_logs" "slow_logs" {
	instance_id   = local.dcdb_id
	start_time    = "%s"
	end_time      = "%s"
	shard_id      = "shard-1b5r04az"
	db            = "tf_test_db"
	order_by      = "query_time_sum"
	order_by_type = "desc"
	slave         = 0
}

`
