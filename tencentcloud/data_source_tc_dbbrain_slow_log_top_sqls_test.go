package tencentcloud

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudDbbrainSlowLogTopSqlsDataSource_basic(t *testing.T) {
	t.Parallel()
	loc, _ := time.LoadLocation("Asia/Chongqing")
	startTime := time.Now().AddDate(0, 0, -1).In(loc).Format("2006-01-02 15:04:05")
	endTime := time.Now().In(loc).Format("2006-01-02 15:04:05")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccDbbrainSlowLogTopSqlsDataSource, defaultDbBrainInstanceId, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_dbbrain_slow_log_top_sqls.test"),
					resource.TestCheckResourceAttr("data.tencentcloud_dbbrain_slow_log_top_sqls.test", "start_time", startTime),
					resource.TestCheckResourceAttr("data.tencentcloud_dbbrain_slow_log_top_sqls.test", "end_time", endTime),
					resource.TestCheckResourceAttr("data.tencentcloud_dbbrain_slow_log_top_sqls.test", "order_by", "ASC"),
					resource.TestCheckResourceAttr("data.tencentcloud_dbbrain_slow_log_top_sqls.test", "sort_by", "QueryTimeMax"),
					resource.TestCheckResourceAttr("data.tencentcloud_dbbrain_slow_log_top_sqls.test", "product", "mysql"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_slow_log_top_sqls.test", "rows.#"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_slow_log_top_sqls.test", "rows.0.lock_time"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_slow_log_top_sqls.test", "rows.0.lock_time_max"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_slow_log_top_sqls.test", "rows.0.lock_time_min"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_slow_log_top_sqls.test", "rows.0.rows_examined"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_slow_log_top_sqls.test", "rows.0.rows_examined_max"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_slow_log_top_sqls.test", "rows.0.rows_examined_min"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_slow_log_top_sqls.test", "rows.0.query_time"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_slow_log_top_sqls.test", "rows.0.query_time_max"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_slow_log_top_sqls.test", "rows.0.query_time_min"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_slow_log_top_sqls.test", "rows.0.rows_sent"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_slow_log_top_sqls.test", "rows.0.rows_sent_max"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_slow_log_top_sqls.test", "rows.0.rows_sent_min"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_slow_log_top_sqls.test", "rows.0.exec_times"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_slow_log_top_sqls.test", "rows.0.sql_template"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_slow_log_top_sqls.test", "rows.0.sql_text"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_dbbrain_slow_log_top_sqls.test", "rows.0.schema"),
				),
			},
		},
	})
}

const testAccDbbrainSlowLogTopSqlsDataSource = `

data "tencentcloud_dbbrain_slow_log_top_sqls" "test" {
  instance_id = "%s"
  start_time = "%s"
  end_time = "%s"
  sort_by = "QueryTimeMax"
  order_by = "ASC"
  product = "mysql"
}

`
