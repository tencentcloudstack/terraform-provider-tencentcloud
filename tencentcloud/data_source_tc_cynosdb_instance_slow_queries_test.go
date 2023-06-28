package tencentcloud

import (
	"fmt"
	"testing"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
)

func TestAccTencentCloudCynosdbInstanceSlowQueriesDataSource_basic(t *testing.T) {
	t.Parallel()
	loc, _ := time.LoadLocation("Asia/Chongqing")
	startTime := time.Now().AddDate(0, 0, -7).In(loc).Format("2006-01-02 15:04:05")
	endTime := time.Now().In(loc).Format("2006-01-02 15:04:05")
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: fmt.Sprintf(testAccCynosdbInstanceSlowQueriesDataSource, startTime, endTime),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckTencentCloudDataSourceID("data.tencentcloud_cynosdb_instance_slow_queries.instance_slow_queries"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_instance_slow_queries.instance_slow_queries", "id"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_instance_slow_queries.instance_slow_queries", "instance_id"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_instance_slow_queries.instance_slow_queries", "start_time"),
					// resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_instance_slow_queries.instance_slow_queries", "end_time"),
					resource.TestCheckResourceAttr("data.tencentcloud_cynosdb_instance_slow_queries.instance_slow_queries", "username", "keep_dts"),
					resource.TestCheckResourceAttr("data.tencentcloud_cynosdb_instance_slow_queries.instance_slow_queries", "database", "tf_ci_test"),
					resource.TestCheckResourceAttr("data.tencentcloud_cynosdb_instance_slow_queries.instance_slow_queries", "order_by", "QueryTime"),
					resource.TestCheckResourceAttr("data.tencentcloud_cynosdb_instance_slow_queries.instance_slow_queries", "order_by_type", "desc"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_instance_slow_queries.instance_slow_queries", "slow_queries.#"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_instance_slow_queries.instance_slow_queries", "slow_queries.0.timestamp"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_instance_slow_queries.instance_slow_queries", "slow_queries.0.query_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_instance_slow_queries.instance_slow_queries", "slow_queries.0.sql_text"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_instance_slow_queries.instance_slow_queries", "slow_queries.0.user_host"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_instance_slow_queries.instance_slow_queries", "slow_queries.0.user_name"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_instance_slow_queries.instance_slow_queries", "slow_queries.0.database"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_instance_slow_queries.instance_slow_queries", "slow_queries.0.lock_time"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_instance_slow_queries.instance_slow_queries", "slow_queries.0.rows_examined"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_instance_slow_queries.instance_slow_queries", "slow_queries.0.rows_sent"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_instance_slow_queries.instance_slow_queries", "slow_queries.0.sql_template"),
					resource.TestCheckResourceAttrSet("data.tencentcloud_cynosdb_instance_slow_queries.instance_slow_queries", "slow_queries.0.sql_md5"),
				),
			},
		},
	})
}

const testAccCynosdbInstanceSlowQueriesDataSource = CommonCynosdb + `

data "tencentcloud_cynosdb_instance_slow_queries" "instance_slow_queries" {
  instance_id = var.cynosdb_cluster_instance_id
//   start_time = "%s"
//   end_time = "%s"
  username = "keep_dts"
  database = "tf_ci_test"
  order_by = "QueryTime"
  order_by_type = "desc"
}

`
