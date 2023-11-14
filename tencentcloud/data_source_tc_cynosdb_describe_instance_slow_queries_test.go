package tencentcloud

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"testing"
)

func TestAccTencentCloudCynosdbDescribeInstanceSlowQueriesDataSource_basic(t *testing.T) {
	t.Parallel()
	resource.Test(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
		},
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: testAccCynosdbDescribeInstanceSlowQueriesDataSource,
				Check:  resource.ComposeTestCheckFunc(testAccCheckTencentCloudDataSourceID("data.tencentcloud_cynosdb_describe_instance_slow_queries.describe_instance_slow_queries")),
			},
		},
	})
}

const testAccCynosdbDescribeInstanceSlowQueriesDataSource = `

data "tencentcloud_cynosdb_describe_instance_slow_queries" "describe_instance_slow_queries" {
  instance_id = "cynosdbmysql-ins-123"
  start_time = "2022-01-01 12:00:00"
  end_time = "2022-01-01 14:00:00"
  username = "root"
  host = "10.10.10.10"
  database = "db1"
  order_by = "QueryTime"
  order_by_type = "desc"
  }

`
